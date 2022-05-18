package cortex

import "github.com/rancher/opni/pkg/core"

func (p *Plugin) CanInstall() error {
	return nil
}

func (p *Plugin) Install(cluster *core.Reference) error {
	return nil
}

func (p *Plugin) Uninstall(clustre *core.Reference) error {
	return nil
}

func (p *Plugin) InstallerTemplate() string {
	return `helm install opni-monitoring-agent ` +
		`{{ arg "input" "Namespace" "+omitEmpty" "+default:opni-monitoring-agent" "+format:-n {{ value }}" }} ` +
		`oci://ghcr.io/kralicky/helm/opni-monitoring-agent --version=0.4.1 ` +
		`--set "token={{ .Token }},pin={{ .Pin }},address={{ .Address }}" ` +
		`{{ arg "toggle" "Install Prometheus Operator" "+omitEmpty" "+default:false" "+format:--set kube-prometheus-stack.enabled={{ value }}" }} ` +
		`--create-namespace`
}