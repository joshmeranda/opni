package gateway

import (
	"context"

	loggingv1beta1 "github.com/rancher/opni/apis/logging/v1beta1"
	"github.com/rancher/opni/pkg/resources"
	"github.com/rancher/opni/plugins/logging/pkg/apis/opensearch"
	"github.com/rancher/opni/plugins/logging/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (p *Plugin) GetDetails(ctx context.Context, cluster *opensearch.ClusterReference) (*opensearch.OpensearchDetails, error) {

	// Get the external URL
	var binding *loggingv1beta1.MulticlusterRoleBinding
	var opnimgmt *loggingv1beta1.OpniOpensearch

	opnimgmt = &loggingv1beta1.OpniOpensearch{}
	if err := p.k8sClient.Get(ctx, types.NamespacedName{
		Name:      p.opensearchCluster.Name,
		Namespace: p.storageNamespace,
	}, opnimgmt); err != nil {
		p.logger.Errorf("unable to fetch opniopensearch object: %v", err)
		return nil, err
	}

	labels := map[string]string{
		resources.OpniClusterID: cluster.AuthorizedClusterID,
	}
	secrets := &corev1.SecretList{}
	if err := p.k8sClient.List(ctx, secrets, client.InNamespace(p.storageNamespace), client.MatchingLabels(labels)); err != nil {
		p.logger.Errorf("unable to list secrets: %v", err)
		return nil, err
	}

	if len(secrets.Items) != 1 {
		p.logger.Error("no credential secrets found")
		return nil, errors.ErrGetDetailsInvalidList(cluster.AuthorizedClusterID)
	}

	return &opensearch.OpensearchDetails{
		Username: secrets.Items[0].Name,
		Password: string(secrets.Items[0].Data["password"]),
		ExternalURL: func() string {
			if binding != nil {
				return binding.Spec.OpensearchExternalURL
			}
			if opnimgmt != nil {
				return opnimgmt.Spec.ExternalURL
			}
			return ""
		}(),
		TracingEnabled: true,
	}, nil
}
