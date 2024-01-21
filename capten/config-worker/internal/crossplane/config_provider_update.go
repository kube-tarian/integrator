package crossplane

import (
	"context"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

const (
	CrossPlaneResource  = "crossplane"
	CrossplaneNamespace = "crossplane-system"
)

func (cp *CrossPlaneApp) configureConfigProviderUpdate(ctx context.Context, req *model.CrossplaneProviderUpdate) (status string, err error) {
	logger.Infof("configuring config provider %s update", req.CloudType)

	cloudType := strings.ToLower(req.CloudType)
	syncPath := fmt.Sprintf("/infra/crossplane/argocd-apps/templates/package-k8s/%s-package/%s-k8s-package.yaml", cloudType, cloudType)

	ns, resName, err := getAppNameNamespace(ctx, syncPath)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get name and namespace from")
	}

	fmt.Println("ns => " + ns)
	fmt.Println("resname => " + resName)

	err = cp.helper.SyncArgoCDApp(ctx, ns, resName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to sync argocd app")
	}
	logger.Infof("synched provider config main-app %s", resName)

	err = cp.helper.WaitForArgoCDToSync(ctx, ns, resName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to fetch argocd app")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}