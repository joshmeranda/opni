package alerting

import (
	"context"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rancher/opni/pkg/alerting/shared"
	"github.com/rancher/opni/pkg/util/future"
	"github.com/rancher/opni/plugins/metrics/pkg/apis/cortexadmin"
	"github.com/rancher/opni/plugins/metrics/pkg/apis/cortexops"

	alertingv1 "github.com/rancher/opni/pkg/apis/alerting/v1"
	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	managementv1 "github.com/rancher/opni/pkg/apis/management/v1"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/machinery"
	"github.com/rancher/opni/pkg/plugins/apis/system"
	natsutil "github.com/rancher/opni/pkg/util/nats"
	"github.com/rancher/opni/plugins/alerting/pkg/alerting/alertstorage"
	"github.com/rancher/opni/plugins/alerting/pkg/alerting/drivers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const ApiExtensionBackoff = time.Second * 5

func (p *Plugin) UseManagementAPI(client managementv1.ManagementClient) {
	p.mgmtClient.Set(client)
	cfg, err := client.GetConfig(context.Background(),
		&emptypb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		p.Logger.With(
			"err", err,
		).Error("Failed to get mgmnt config")
		os.Exit(1)
	}
	objectList, err := machinery.LoadDocuments(cfg.Documents)
	if err != nil {
		p.Logger.With(
			"err", err,
		).Error("failed to load config")
		os.Exit(1)
	}
	objectList.Visit(func(config *v1beta1.GatewayConfig) {
		opt := shared.NewAlertingOptions{
			Namespace:             config.Spec.Alerting.Namespace,
			WorkerNodesService:    config.Spec.Alerting.WorkerNodeService,
			WorkerNodePort:        config.Spec.Alerting.WorkerPort,
			WorkerStatefulSet:     config.Spec.Alerting.WorkerStatefulSet,
			ControllerNodeService: config.Spec.Alerting.ControllerNodeService,
			ControllerNodePort:    config.Spec.Alerting.ControllerNodePort,
			ControllerClusterPort: config.Spec.Alerting.ControllerClusterPort,
			ConfigMap:             config.Spec.Alerting.ConfigMap,
			ManagementHookHandler: config.Spec.Alerting.ManagementHookHandler,
		}
		if os.Getenv(shared.LocalBackendEnvToggle) != "" {
			opt.WorkerNodesService = "http://localhost"
			opt.ControllerNodeService = "http://localhost"

		}

		p.configureAlertManagerConfiguration(
			p.Ctx,
			drivers.WithLogger(p.Logger.Named("alerting-manager")),
			drivers.WithAlertingOptions(&opt),
			drivers.WithManagementClient(client),
		)
	})
	go func() {
		p.watchGlobalCluster(client)
	}()

	go func() {
		p.watchGlobalClusterHealthStatus(client)
	}()
	<-p.Ctx.Done()
}

// UseKeyValueStore Alerting Condition & Alert Endpoints are stored in K,V stores
func (p *Plugin) UseKeyValueStore(client system.KeyValueStoreClient) {
	var (
		nc  *nats.Conn
		err error
	)
	p.storageNode = alertstorage.NewStorageNode(
		alertstorage.WithStorage(&alertstorage.StorageAPIs{
			Conditions:           system.NewKVStoreClient[*alertingv1.AlertCondition](client),
			Endpoints:            system.NewKVStoreClient[*alertingv1.AlertEndpoint](client),
			SystemTrackerStorage: future.New[nats.KeyValue](),
		}),
	)

	nc, err = natsutil.AcquireNATSConnection(
		p.Ctx,
		natsutil.WithLogger(p.Logger),
		natsutil.WithNatsOptions([]nats.Option{
			nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
				if s != nil {
					p.Logger.Error("nats : async error in %q/%q: %v", s.Subject, s.Queue, err)
				} else {
					p.Logger.Warn("nats : async error outside subscription")
				}
			}),
		}),
	)
	if err != nil {
		p.Logger.With("err", err).Error("fatal error connecting to NATs")
	}
	p.natsConn.Set(nc)
	mgr, err := p.natsConn.Get().JetStream()
	if err != nil {
		panic(err)
	}
	kv, err := mgr.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      shared.AgentDisconnectBucket,
		Description: "track system agent status updates",
		Storage:     nats.FileStorage,
	})
	if err != nil {
		panic(err)
	}
	p.storageNode.SetSystemTrackerStorage(kv)
	// spawn a reindexing task
	go func() {
		p.restartAgentDisconnectTrackers()
	}()
	<-p.Ctx.Done()
}

func (p *Plugin) restartAgentDisconnectTrackers() {
	lg := p.Logger.With("re-indexing", "agent-disconnect-trackers")
	ids, conds, err := p.storageNode.ListWithKeyConditionStorage(p.Ctx)
	if err != nil {
		lg.With("err", err).Error("failed to list alert conditions")
		return
	}
	for i, id := range ids {
		if s := conds[i].GetAlertType().GetSystem(); s != nil {
			// this checks that we won't crash when importing existing conditions from versions < 0.6
			if s.GetClusterId() != nil && s.GetTimeout() != nil {
				p.onSystemConditionCreate(id, s)
			} else {
				// delete invalid conditions that won't do anything
				_, err := p.DeleteAlertCondition(p.Ctx, &corev1.Reference{Id: id})
				if err != nil {
					lg.With("err", err).Error("failed to delete invalid condition")
				}
			}
		}
	}
	lg.Info("re-indexing complete")
}

func (p *Plugin) UseAPIExtensions(intf system.ExtensionClientInterface) {
	go func() {
		for {
			breakOut := false
			ccCortexAdmin, err := intf.GetClientConn(p.Ctx, "CortexAdmin")
			if err != nil {
				p.Logger.Errorf("alerting failed to get cortex admin client conn %s", err)
				p.adminClient = future.New[cortexadmin.CortexAdminClient]()
				UnregisterDatasource(shared.MonitoringDatasource)
			} else {
				adminClient := cortexadmin.NewCortexAdminClient(ccCortexAdmin)
				p.adminClient.Set(adminClient)
				RegisterDatasource(shared.MonitoringDatasource, NewAlertingMonitoringStore(p, p.Logger))
			}

			ccCortexOps, err := intf.GetClientConn(p.Ctx, "CortexOps")
			if err != nil {
				p.Logger.Errorf("alerting failed to get cortex ops client conn %s", err)
				p.cortexOpsClient = future.New[cortexops.CortexOpsClient]()
			} else {
				opsClient := cortexops.NewCortexOpsClient(ccCortexOps)
				p.cortexOpsClient.Set(opsClient)
			}

			select {
			case <-p.Ctx.Done():
				breakOut = true
			default:
				continue
			}
			if breakOut {
				break
			}
			time.Sleep(ApiExtensionBackoff)
		}
	}()
}
