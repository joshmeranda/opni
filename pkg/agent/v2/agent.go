package v2

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/rancher/opni/pkg/ident/identserver"
	"github.com/rancher/opni/pkg/patch"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	controlv1 "github.com/rancher/opni/pkg/apis/control/v1"
	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	"github.com/rancher/opni/pkg/bootstrap"
	"github.com/rancher/opni/pkg/clients"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/health"
	"github.com/rancher/opni/pkg/health/annotations"
	"github.com/rancher/opni/pkg/ident"
	"github.com/rancher/opni/pkg/keyring"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/machinery"
	"github.com/rancher/opni/pkg/plugins"
	"github.com/rancher/opni/pkg/plugins/apis/apiextensions"
	"github.com/rancher/opni/pkg/plugins/hooks"
	"github.com/rancher/opni/pkg/plugins/meta"
	"github.com/rancher/opni/pkg/plugins/types"
	"github.com/rancher/opni/pkg/storage"
	"github.com/rancher/opni/pkg/trust"
	"github.com/rancher/opni/pkg/util"
	"github.com/rancher/opni/pkg/util/fwd"
)

var ErrRebootstrap = errors.New("re-bootstrap requested")

type Agent struct {
	AgentOptions

	config       v1beta1.AgentConfigSpec
	router       *gin.Engine
	Logger       *zap.SugaredLogger
	pluginLoader *plugins.PluginLoader

	tenantID         string
	identityProvider ident.Provider
	keyringStore     storage.KeyringStore
	gatewayClient    clients.GatewayClient
	trust            trust.Strategy

	loadedExistingKeyring bool
}

type AgentOptions struct {
	bootstrapper          bootstrap.Bootstrapper
	unmanagedPluginLoader *plugins.PluginLoader
	rebootstrap           bool
}

type AgentOption func(*AgentOptions)

func (o *AgentOptions) apply(opts ...AgentOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithBootstrapper(bootstrapper bootstrap.Bootstrapper) AgentOption {
	return func(o *AgentOptions) {
		o.bootstrapper = bootstrapper
	}
}

func WithUnmanagedPluginLoader(pluginLoader *plugins.PluginLoader) AgentOption {
	return func(o *AgentOptions) {
		o.unmanagedPluginLoader = pluginLoader
	}
}

func WithRebootstrap(rebootstrap bool) AgentOption {
	return func(o *AgentOptions) {
		o.rebootstrap = rebootstrap
	}
}

func New(ctx context.Context, conf *v1beta1.AgentConfig, opts ...AgentOption) (*Agent, error) {
	options := AgentOptions{}
	options.apply(opts...)
	level := logger.DefaultLogLevel.Level()
	if conf.Spec.LogLevel != "" {
		l, err := zap.ParseAtomicLevel(conf.Spec.LogLevel)
		if err != nil {
			return nil, fmt.Errorf("error parsing log level: %w", err)
		}
		level = l.Level()
	}
	lg := logger.New(logger.WithLogLevel(level)).Named("agent")
	lg.Debugf("using log level: %s", level)

	var pl *plugins.PluginLoader
	if options.unmanagedPluginLoader != nil {
		pl = options.unmanagedPluginLoader
	} else {
		pl = plugins.NewPluginLoader()
	}

	pl.Hook(hooks.OnLoadM(func(p types.CapabilityNodePlugin, m meta.PluginMeta) {
		lg.Infof("loaded capability node plugin %s", m.Module)
	}))

	router := gin.New()
	routerMutex := &sync.Mutex{}
	router.Use(logger.GinLogger(lg), gin.Recovery())
	pprof.Register(router)

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	pl.Hook(hooks.OnLoadM(func(p types.HTTPAPIExtensionPlugin, md meta.PluginMeta) {
		ctx, ca := context.WithTimeout(ctx, 10*time.Second)
		defer ca()
		cfg, err := p.Configure(ctx, apiextensions.NewInsecureCertConfig())
		if err != nil {
			lg.With(
				zap.String("plugin", md.Module),
				zap.Error(err),
			).Error("failed to configure routes")
			return
		}
		setupPluginRoutes(lg, routerMutex, router, cfg, md, []string{"/healthz", "/metrics"})
	}))

	initCtx, initCancel := context.WithTimeout(ctx, 10*time.Second)
	defer initCancel()

	ip, err := ident.GetProvider(conf.Spec.IdentityProvider)
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}
	id, err := ip.UniqueIdentifier(initCtx)
	if err != nil {
		return nil, fmt.Errorf("error getting unique identifier: %w", err)
	}

	broker, err := machinery.BuildKeyringStoreBroker(initCtx, conf.Spec.Storage)
	if err != nil {
		return nil, fmt.Errorf("error configuring keyring store broker: %w", err)
	}
	ks := broker.KeyringStore("agent", &corev1.Reference{
		Id: id,
	})

	var kr keyring.Keyring
	var loadedExistingKeyring bool
	var shouldBootstrap bool
	if existing, err := ks.Get(initCtx); err == nil {
		lg.Info("loaded existing keyring")
		kr = existing
		loadedExistingKeyring = true
	} else if errors.Is(err, storage.ErrNotFound) {
		lg.Info("no existing keyring found, starting bootstrap process")
		shouldBootstrap = true
	} else {
		return nil, fmt.Errorf("keyring store error: %w", err)
	}

	if options.rebootstrap {
		if conf.Spec.ContainsBootstrapCredentials() {
			lg.Info("attempting to re-bootstrap agent to generate new keyring")
			shouldBootstrap = true
		} else {
			lg.Warn("re-bootstrap requested, but no bootstrap credentials were provided in the config file")
		}
	}

	if shouldBootstrap {
		kr, err = options.bootstrapper.Bootstrap(initCtx, ip)
		if err != nil {
			return nil, fmt.Errorf("error during bootstrap: %w", err)
		}
		for {
			// Don't let this fail easily, otherwise we will lose the keyring forever.
			// Keep retrying until it succeeds.
			err = ks.Put(ctx, kr)
			if err != nil {
				lg.With(zap.Error(err)).Error("failed to persist keyring (retry in 1 second)")
				time.Sleep(1 * time.Second)
			} else {
				if options.rebootstrap {
					lg.Info("successfully replaced keyring")
				}
				break
			}
		}
		lg.Info("bootstrap completed successfully")
	}

	// Run post-bootstrap finalization. If this has previously succeeded,
	// it will likely do nothing.
	if err := options.bootstrapper.Finalize(initCtx); err != nil {
		lg.With(
			zap.Error(err),
		).Warn("error in post-bootstrap finalization")
	}
	trust, err := machinery.BuildTrustStrategy(conf.Spec.TrustStrategy, kr)
	if err != nil {
		return nil, fmt.Errorf("error building trust strategy: %w", err)
	}

	gatewayClient, err := clients.NewGatewayClient(ctx,
		conf.Spec.GatewayAddress, ip, kr, trust)
	if err != nil {
		return nil, fmt.Errorf("error configuring gateway client: %w", err)
	}
	controlv1.RegisterIdentityServer(gatewayClient, identserver.NewFromProvider(ip))

	hm := health.NewAggregator(health.WithStaticAnnotations(map[string]string{
		annotations.AgentVersion: annotations.Version2,
	}))
	controlv1.RegisterHealthServer(gatewayClient, hm)

	pl.Hook(hooks.OnLoadMC(func(hc controlv1.HealthClient, m meta.PluginMeta, cc *grpc.ClientConn) {
		client := controlv1.NewHealthClient(cc)
		hm.AddClient(path.Base(m.BinaryPath), client)
	}))

	pl.Hook(hooks.OnLoadMC(func(ext types.StreamAPIExtensionPlugin, md meta.PluginMeta, cc *grpc.ClientConn) {
		lg.With(
			zap.String("plugin", md.Module),
		).Debug("loaded stream api extension plugin")
		gatewayClient.RegisterSplicedStream(cc, md.ShortName())
	}))

	return &Agent{
		AgentOptions: options,
		config:       conf.Spec,
		router:       router,
		Logger:       lg,
		pluginLoader: pl,

		tenantID:         id,
		identityProvider: ip,
		keyringStore:     ks,
		trust:            trust,
		gatewayClient:    gatewayClient,

		loadedExistingKeyring: loadedExistingKeyring,
	}, nil
}

func (a *Agent) ListenAndServe(ctx context.Context) error {
	if a.unmanagedPluginLoader == nil {
		manifests, err := a.syncPlugins(ctx)
		if err != nil {
			if status.Code(err) == codes.Unauthenticated && a.loadedExistingKeyring {
				a.Logger.With(
					zap.Error(err),
				).Warn("The agent failed to authorize to the gateway using an existing keyring. " +
					"This could be due to a leftover keyring from a previous installation that was not deleted.",
				)
				if a.config.ContainsBootstrapCredentials() {
					a.Logger.Warn("Bootstrap credentials have been provided in the config file - " +
						"the agent will restart and attempt to re-bootstrap a new keyring using these credentials.")
					return ErrRebootstrap
				}
			}
			return fmt.Errorf("error syncing plugins: %w", err)
		}
		// eventually passed to runGatewayClient
		ctx = metadata.AppendToOutgoingContext(ctx,
			controlv1.ManifestDigestKey, manifests.Digest())

		done := make(chan struct{})
		a.pluginLoader.Hook(hooks.OnLoadingCompleted(func(numPlugins int) {
			a.Logger.Infof("loaded %d plugins", numPlugins)
			close(done)
		}))

		a.pluginLoader.LoadPlugins(ctx, a.config.Plugins, plugins.AgentScheme)

		select {
		case <-done:
		case <-ctx.Done():
			return ctx.Err()
		}
	} else {
		a.Logger.Info("using unmanaged plugin loader")
	}

	listener, err := net.Listen("tcp4", a.config.ListenAddress)
	if err != nil {
		return err
	}
	ctx, ca := context.WithCancel(ctx)

	e1 := lo.Async(func() error {
		return util.ServeHandler(ctx, a.router.Handler(), listener)
	})

	e2 := lo.Async(func() error {
		return a.runGatewayClient(ctx)
	})

	return util.WaitAll(ctx, ca, e1, e2)
}

func (a *Agent) ListenAddress() string {
	return a.config.ListenAddress
}

func setupPluginRoutes(
	lg *zap.SugaredLogger,
	mutex *sync.Mutex,
	router *gin.Engine,
	cfg *apiextensions.HTTPAPIExtensionConfig,
	pluginMeta meta.PluginMeta,
	reservedPrefixRoutes []string,
) {
	mutex.Lock()
	defer mutex.Unlock()

	sampledLogger := logger.New(
		logger.WithSampling(&zap.SamplingConfig{
			Initial:    1,
			Thereafter: 0,
		}),
	).Named("api")
	forwarder := fwd.To(cfg.HttpAddr,
		fwd.WithLogger(sampledLogger),
		fwd.WithDestHint(path.Base(pluginMeta.BinaryPath)),
	)
ROUTES:
	for _, route := range cfg.Routes {
		for _, reservedPrefix := range reservedPrefixRoutes {
			if strings.HasPrefix(route.Path, reservedPrefix) {
				lg.With(
					"route", route.Method+" "+route.Path,
					"plugin", pluginMeta.Module,
				).Warn("skipping route for plugin as it conflicts with a reserved prefix")
				continue ROUTES
			}
		}
		lg.With(
			"route", route.Method+" "+route.Path,
			"plugin", pluginMeta.Module,
		).Debug("configured route for plugin")
		router.Handle(route.Method, route.Path, forwarder)
	}
}

func (a *Agent) runGatewayClient(ctx context.Context) error {
	lg := a.Logger
	isRetry := false
	for ctx.Err() == nil {
		if isRetry {
			time.Sleep(1 * time.Second)
			lg.Info("attempting to reconnect...")
		} else {
			lg.Info("connecting to gateway...")
		}
		// this connects plugin extension servers to the agent's totem server.
		// clients on the other side of this stream will have access to gateway
		// services.
		_, errF := a.gatewayClient.Connect(ctx) // this unused cc can access all services
		if !errF.IsSet() {
			if isRetry {
				lg.Info("gateway reconnected")
			} else {
				lg.Info("gateway connected")
			}

			lg.With(
				zap.Error(errF.Get()), // this will block until an error is received
			).Warn("disconnected from gateway")
		} else {
			lg.With(
				zap.Error(errF.Get()),
			).Warn("error connecting to gateway")
		}

		switch util.StatusCode(errF.Get()) {
		case codes.FailedPrecondition:
			// this error will be returned if the agent needs to restart
			lg.Warn("encountered non-retriable error")
			return errF.Get()
		case codes.Unauthenticated:
			return errF.Get()
		}
		isRetry = true
	}
	lg.With(
		zap.Error(ctx.Err()),
	).Warn("shutting down gateway client")
	return ctx.Err()
}

func (a *Agent) syncPlugins(ctx context.Context) (_ *controlv1.ManifestMetadataList, retErr error) {
	a.Logger.Info("attempting to sync plugins with gateway")

	manifestClient := controlv1.NewPluginManifestClient(a.gatewayClient.ClientConn())
	// read local plugins on disk here
	localManifests, err := patch.GetFilesystemPlugins(a.config.Plugins, a.Logger)
	if err != nil {
		return nil, fmt.Errorf("hit an error getting filesystem plugins based on config %s", err)
	}

	patchList, err := manifestClient.SendManifestsOrKnownPatch(ctx, localManifests, grpc.UseCompressor("zstd"))
	if err != nil {
		return nil, fmt.Errorf("hit an error requesting plugins/patches from gateway %s", err)
	}
	a.Logger.Info("received patch manifest from gateway")

	done, err := patch.PatchWith(ctx, a.config.Plugins, patchList, a.Logger, manifestClient)
	if err != nil {
		return nil, fmt.Errorf("hit an error patching agent's plugins : %s", err)
	}

	go func() {
		<-done
		a.Logger.Debug("background patching tasks complete")
	}()

	// read local manifests again
	localManifests, err = patch.GetFilesystemPlugins(a.config.Plugins, a.Logger)
	if err != nil {
		return nil, fmt.Errorf("filesystem plugins sanity check failed after patching : %s", err)
	}

	return localManifests, nil
}
