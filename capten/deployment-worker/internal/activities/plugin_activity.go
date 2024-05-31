package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	vaultcred "github.com/kube-tarian/kad/capten/common-pkg/vault-cred"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/captensdk"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/dbstorepreactions/postgresstore"
	"github.com/kube-tarian/kad/capten/model"
	v1 "k8s.io/api/core/v1"
)

const (
	postgresStoreInitializingStatus       = "postgres-" + "initializing"
	postgresStoreInitializedStatus        = "postgres-" + "initialized"
	postgresStoreInitializeFailedStatus   = "postgres-" + "initialize-faield"
	postgresStoreUninitializingStatus     = "postgres-" + "uninitializing"
	postgresStoreUninitializedStatus      = "postgres-" + "uninitialized"
	postgresStoreUninitializeFailedStatus = "postgres-" + "uninitialize-failed"
	vaultStoreInitializingStatus          = "vaultstore-" + "initializing"
	vaultStoreInitializedStatus           = "vaultstore-" + "initialized"
	vaultStoreInitializeFailedStatus      = "vaultstore-" + "initialize-failed"
	vaultStoreUnitializingStatus          = "vaultstore-" + "uninitializing"
	vaultStoreUninitializedStatus         = "vaultstore-" + "uninitialized"
	vaultStoreUninitializeFailedStatus    = "vaultstore-" + "uninitialize-failed"
	mtlsInitializingStatus                = "mtls-" + "initializing"
	mtlsInitializedStatus                 = "mtls-" + "initialized"
	mtlsInitializeFailedStatus            = "mtls-" + "initialize-failed"
	mtlsUnitializingStatus                = "mtls-" + "uninitializing"
	mtlsUnitializedStatus                 = "mtls-" + "uninitialized"
	mtlsUnitializeFailedStatus            = "mtls-" + "uninitialize-failed"
	deleteUnitiazingStatus                = "delete-" + "uninitializing"
	deleteSuccessStatus                   = "delete-" + "success"
	deleteFailedStatus                    = "delete-" + "failed"
	deployedStatus                        = "deployed"

	pluginConfigmapNameTemplate = "-init-config"
)

type Configuration struct {
	AgentAddress string `envconfig:"AGENT_ADDRESSES" required:"true"`
}

type PluginActivities struct {
	config    *Configuration
	as        *captenstore.Store
	k8sClient *k8s.K8SClient
}

func NewPluginActivities() (*PluginActivities, error) {
	conf := &Configuration{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("agent config read faile, %v", err)
	}

	as, err := captenstore.NewStore(logger)
	if err != nil {
		return nil, err
	}

	k8sclient, err := k8s.NewK8SClient(logger)
	if err != nil {
		return nil, err
	}

	return &PluginActivities{
		config:    conf,
		as:        as,
		k8sClient: k8sclient,
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionPostgresStoreActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	logger.Infof("Deploy postgres store started")
	err := p.updateStatus(req.PluginName, postgresStoreInitializingStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}

	pluginInitConfigmapName := req.PluginName + pluginConfigmapNameTemplate
	err = postgresstore.SetupPostgresDatabase(logger, req.PluginName, req.DefaultNamespace, pluginInitConfigmapName, p.k8sClient)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" %s\"}", err.Error())),
		}, err
	}

	err = p.updateStatus(req.PluginName, postgresStoreInitializedStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginUndeployPreActionPostgresStoreActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, postgresStoreUninitializingStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}

	err = p.as.DeleteClusterPluginConfig(req.PluginName)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}

	err = p.updateStatus(req.PluginName, postgresStoreUninitializedStatus)
	if err != nil {
		logger.Errorf("failed to update uninitialized status, %v", err)
	}

	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionVaultStoreActivity(
	ctx context.Context,
	req *clusterpluginspb.Plugin,
) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, vaultStoreInitializingStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}

	// Get vault token to access vault secret path
	vaultPaths := []string{"plugin/" + req.PluginName + "/*", "generic/*"}
	token, err := vaultcred.GetAppRoleToken(req.PluginName, vaultPaths)
	if err != nil {
		logger.Errorf("failed to get vault token for the path, %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"vault token status: %s\"}", err.Error())),
		}, err
	}

	// Create a secret with token data
	err = p.k8sClient.CreateOrUpdateSecret(ctx, req.DefaultNamespace, req.PluginName+"-vault-token", v1.SecretTypeOpaque, map[string][]byte{
		"token":       []byte(token),
		"secret-path": []byte(strings.Join(vaultPaths, ",")),
	}, nil)
	if err != nil {
		logger.Errorf("failed to create secret %s with vault token, %v", req.PluginName+"-vault-token", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"vault token secret status: %s\"}", err.Error())),
		}, err
	}

	pluginInitConfigmapName := req.PluginName + pluginConfigmapNameTemplate
	err = p.createUpdateConfigmap(ctx, req.DefaultNamespace, pluginInitConfigmapName, map[string]string{
		"vault-token-secret-name": req.PluginName + "-vault-token",
	})
	if err != nil {
		logger.Errorf("createupdate configmap failed: %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update configmap: %s\"}", pluginInitConfigmapName)),
		}, err
	}

	err = p.updateStatus(req.PluginName, vaultStoreInitializedStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginUndeployPreActionVaultStoreActivity(
	ctx context.Context,
	req *clusterpluginspb.Plugin,
) (*model.ResponsePayload, error) {
	// If any failure log error and should not return error
	err := p.updateStatus(req.PluginName, vaultStoreUnitializingStatus)
	if err != nil {
		logger.Errorf("failed to update undeploy status to vaultstore-uninitializing, %v", err)
	}

	// Delete App role
	err = vaultcred.DeleteAppRole(req.PluginName)
	if err != nil {
		logger.Errorf("failed to delete vault app role, %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	// Delete a secret with token data
	err = p.k8sClient.DeleteSecret(ctx, req.DefaultNamespace, req.PluginName+"-vault-token")
	if err != nil {
		logger.Errorf("failed to delete secret %s, %v", req.PluginName+"-vault-token", err)
	}

	err = p.updateStatus(req.PluginName, vaultStoreUninitializedStatus)
	if err != nil {
		logger.Errorf("failed to update undeploy status to vaultstore-uninitialized, %v", err)
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionMTLSActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, mtlsInitializingStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	// Write the mtls in the vault/conigmap
	captenSDKClient, err := captensdk.NewMTLSClient(logger)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	agentClusterCAIssuerName := os.Getenv("AGENT_CLUSTER_CA_ISSUER_NAME")
	if len(agentClusterCAIssuerName) == 0 {
		agentClusterCAIssuerName = "agent-ca-issuer"
	}

	pluginInitConfigmapName := req.PluginName + pluginConfigmapNameTemplate
	err = captenSDKClient.CreateCertificates(req.PluginName, req.DefaultNamespace, agentClusterCAIssuerName, pluginInitConfigmapName, p.k8sClient)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	err = p.updateStatus(req.PluginName, mtlsInitializedStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	logger.Infof("MTLS certificate creation finished")
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginUndeployPreActionMTLSActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, mtlsUnitializingStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	// Write the mtls in the vault/conigmap
	captenSDKClient, err := captensdk.NewMTLSClient(logger)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	err = captenSDKClient.DeleteCertificate(req.PluginName, req.DefaultNamespace)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	err = p.updateStatus(req.PluginName, mtlsUnitializedStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	logger.Infof("MTLS certificate deletion finished")
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

// PluginDeployPostActionActivity... Updates the plugin deployment as "installed"
func (p *PluginActivities) PluginDeployPostActionActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	pluginInitConfigmapName := req.PluginName + pluginConfigmapNameTemplate
	err := p.createUpdateConfigmap(ctx, req.DefaultNamespace, pluginInitConfigmapName, map[string]string{
		"capten-agent-address": p.config.AgentAddress,
	})
	if err != nil {
		logger.Errorf("update configmap failed to add agent address: %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update configmap failed, %s\"}", pluginInitConfigmapName)),
		}, err
	}

	err = p.updateStatus(req.PluginName, deployedStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

// PluginDeployPostActionActivity... Updates the plugin deployment as "installed"
func (p *PluginActivities) PluginUndeployPostActionActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	pluginInitConfigmapName := req.PluginName + pluginConfigmapNameTemplate
	err := p.k8sClient.DeleteConfigmap(ctx, req.DefaultNamespace, pluginInitConfigmapName)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"delete configmap %s faled\"}", pluginInitConfigmapName)),
		}, err
	}

	err = p.as.DeleteClusterPluginConfig(req.PluginName)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	// TODO: Is delete namespace to be invoked?

	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginUndeployActivity(ctx context.Context, req *model.DeployerDeleteRequest) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, deleteUnitiazingStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	resp, err := uninstallApplication(req)
	if err != nil {
		_ = p.updateStatus(req.ReleaseName, deleteFailedStatus)
		return &resp, err
	}

	err = p.updateStatus(req.ReleaseName, deleteSuccessStatus)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) updateStatus(releaseName, status string) error {
	plugin, err := p.as.GetClusterPluginConfig(releaseName)
	if err != nil {
		return fmt.Errorf("plugin application %s not found in database", releaseName)
	}
	plugin.InstallStatus = status
	p.as.UpsertClusterPluginConfig(plugin)
	return nil
}

func (p *PluginActivities) createUpdateConfigmap(ctx context.Context, namespace, cmName string, data map[string]string) error {
	err := p.k8sClient.CreateNamespace(ctx, namespace)
	if err != nil {
		logger.Errorf("Creation of namespace failed: %v", err)
		return fmt.Errorf("creation of namespace faield")
	}
	cm, err := p.k8sClient.GetConfigmap(ctx, namespace, cmName)
	if err != nil {
		logger.Infof("plugin configmap %s not found", cmName)
		err = p.k8sClient.CreateConfigmap(ctx, namespace, cmName, data, map[string]string{})
		if err != nil {
			return fmt.Errorf("failed to create configmap %v", cmName)
		}
	}
	// configmap found but data is empty/nil
	if cm == nil {
		cm = map[string]string{}
	}
	for k, v := range data {
		cm[k] = v
	}
	err = p.k8sClient.UpdateConfigmap(ctx, namespace, cmName, cm)
	if err != nil {
		return fmt.Errorf("plugin configmap %s not found", cmName)
	}
	return nil
}

func (p *PluginActivities) PluginDeployUpdateStatusActivity(ctx context.Context, pluginName, status string) (model.ResponsePayload, error) {
	err := p.updateStatus(pluginName, status)
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	return model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}
