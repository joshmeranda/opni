package gateway

import (
	"github.com/rancher/opni/pkg/capabilities/wellknown"
	streamext "github.com/rancher/opni/pkg/plugins/apis/apiextensions/stream"
	"github.com/rancher/opni/plugins/metrics/pkg/apis/node"
	"github.com/rancher/opni/plugins/metrics/pkg/apis/remotewrite"
)

func (p *Plugin) StreamServers() []streamext.Server {
	return []streamext.Server{
		{
			Desc:              &remotewrite.RemoteWrite_ServiceDesc,
			Impl:              &p.cortexRemoteWrite,
			RequireCapability: wellknown.CapabilityMetrics,
		},
		{
			Desc:              &node.NodeMetricsCapability_ServiceDesc,
			Impl:              &p.metrics,
			RequireCapability: wellknown.CapabilityMetrics,
		},
	}
}
