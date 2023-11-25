package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	appconfig "github.com/kube-tarian/kad/capten/config-worker/internal/app_config"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
)

const (
	configClusterEndpoint = "configClusterEndpoint"
)

type Config struct {
	PluginConfigFile        string `envconfig:"CROSSPLANE_PLUGIN_CONFIG_FILE" default:"/crossplane_plugin_config.json"`
	CloudProviderEntityName string `envconfig:"CLOUD_PROVIDER_ENTITY_NAME" default:"cloud-provider"`
}

type CrossPlaneApp struct {
	helper       *appconfig.AppGitConfigHelper
	pluginConfig *crossplanePluginConfig
	cfg          Config
}

func NewCrossPlaneApp() (*CrossPlaneApp, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	helper, err := appconfig.NewAppGitConfigHelper()
	if err != nil {
		return nil, err
	}

	pluginConfig, err := readCrossPlanePluginConfig(cfg.PluginConfigFile)
	if err != nil {
		return nil, err
	}
	return &CrossPlaneApp{pluginConfig: pluginConfig, helper: helper, cfg: cfg}, err
}

func readCrossPlanePluginConfig(pluginFile string) (*crossplanePluginConfig, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData crossplanePluginConfig
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &pluginData, nil
}

func getAppNameNamespace(ctx context.Context, fileName string) (string, string, error) {
	k8sclient, err := k8s.NewK8SClient(logging.NewLogger())
	if err != nil {
		return "", "", fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", "", err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return "", "", err
	}

	// For the testing change the reqrepo to template one
	ns, resName, err := k8sclient.DynamicClient.GetNameNamespace(jsonData)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the k8s custom resource: %v", err)
	}

	return ns, resName, nil

}

func (cp *CrossPlaneApp) Configure(ctx context.Context, reqType string, req json.RawMessage) (status string, err error) {
	switch reqType {
	case configClusterEndpoint:
		reqLocal := &model.CrossplaneClusterEndpoint{}
		if err = json.Unmarshal(req, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req, %v", err)
			err = fmt.Errorf("failed to unmarshall the crossplane req")
			return
		}
		status, err = cp.configureClusterEndpoint(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project, %v", err)
			err = fmt.Errorf("failed to configure crossplane project")
		}
	default:
		reqLocal := &model.CrossplaneUseCase{}
		if err = json.Unmarshal(req, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req, %v", err)
			err = fmt.Errorf("failed to unmarshall the crossplane req")
			return
		}
		status, err = cp.configureProjectAndApps(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project, %v", err)
			err = fmt.Errorf("failed to configure crossplane project")
		}
	}
	return
}

func (cp *CrossPlaneApp) configureClusterEndpoint(ctx context.Context, req *model.CrossplaneClusterEndpoint) (status string, err error) {
	logger.Infof("configuring the cluster endpoint for %s", req.RepoURL)
	err = cp.helper.CreateCluster(ctx, req.Namespace, req.Endpoint, req.Name)
	// if err != nil {
	// 	return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to CreateCluster in argocd app")
	// }

	logger.Infof("CreateCluster argocd err: ", err)
	accessToken, err := cp.helper.GetAccessToken(ctx, req.Id)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get token from vault")
	}

	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)
	templateRepo, customerRepo, err := cp.helper.CloneRepos(ctx, cp.pluginConfig.TemplateGitRepo, req.RepoURL, accessToken)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(templateRepo)
	defer os.RemoveAll(customerRepo)

	fileName := filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.File)
	// replace cluster endpoint
	err = updateClusterEndpointDetials(fileName, req)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace the file")
	}

	err = cp.helper.AddToGit(ctx, configClusterEndpoint, req.RepoURL, accessToken)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}

	logger.Infof("added cloned project %s changed to git", req.RepoURL)
	ns, resName, err := getAppNameNamespace(ctx, filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.MainAppGitPath))
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get name and namespace from")
	}

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

func (cp *CrossPlaneApp) configureProjectAndApps(ctx context.Context, req *model.CrossplaneUseCase) (status string, err error) {
	accessToken, err := cp.helper.GetAccessToken(ctx, req.VaultCredIdentifier)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get token from vault")
	}

	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)
	templateRepo, customerRepo, err := cp.helper.CloneRepos(ctx, cp.pluginConfig.TemplateGitRepo, req.RepoURL, accessToken)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(templateRepo)
	defer os.RemoveAll(customerRepo)

	err = cp.synchProviders(req, templateRepo, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to update configs to repo")
	}
	logger.Infof("added provider config resources to cloned project %s", req.RepoURL)

	// update git project url
	if err := replaceCaptenUrls(customerRepo, cp.pluginConfig.TemplateGitRepo, req.RepoURL); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace template url")
	}
	logger.Infof("updated resource configurations in cloned project %s", req.RepoURL)

	err = cp.helper.AddToGit(ctx, req.Type, req.RepoURL, accessToken)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}
	logger.Infof("added cloned project %s changed to git", req.RepoURL)

	err = cp.deployArgoCDApps(ctx, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to depoy argoCD apps")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func (cp *CrossPlaneApp) synchProviders(req *model.CrossplaneUseCase, templateDir, reqRepo string) error {
	err := cp.createProviderConfigs(filepath.Join(templateDir, cp.pluginConfig.ProviderConfigSyncPath), req)
	if err != nil {
		return fmt.Errorf("failed to create provider config, %v", err)
	}

	err = copy.Copy(filepath.Join(templateDir, cp.pluginConfig.CrossplaneConfigSyncPath),
		filepath.Join(reqRepo, cp.pluginConfig.CrossplaneConfigSyncPath),
		copy.Options{
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				return copy.Replace
			}})
	if err != nil {
		return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
	}
	return nil
}

func (cp *CrossPlaneApp) deployArgoCDApps(ctx context.Context, customerRepo string) (err error) {
	logger.Infof("%d main apps to deploy", len(cp.pluginConfig.ArgoCDApps))

	for _, argoApp := range cp.pluginConfig.ArgoCDApps {
		appPath := filepath.Join(customerRepo, argoApp.MainAppGitPath)
		err = cp.deployArgoCDApp(ctx, appPath, argoApp.ChildAppNames, argoApp.SynchApp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *CrossPlaneApp) deployArgoCDApp(ctx context.Context, appPath string, childApps []string, synchApp bool) (err error) {
	ns, resName, err := cp.helper.DeployMainApp(ctx, appPath)
	if err != nil {
		return errors.WithMessage(err, "failed to deploy main app")
	}
	logger.Infof("deployed provider config main-app %s", resName)

	if synchApp {
		err = cp.helper.SyncArgoCDApp(ctx, ns, resName)
		if err != nil {
			return errors.WithMessage(err, "failed to sync argocd app")
		}
		logger.Infof("synched provider config main-app %s", resName)

		err = cp.helper.WaitForArgoCDToSync(ctx, ns, resName)
		if err != nil {
			return errors.WithMessage(err, "failed to fetch argocd app")
		}

		err = cp.syncArgoCDChildApps(ctx, ns, childApps)
		if err != nil {
			return errors.WithMessage(err, "failed to synch argocd child app")
		}
		logger.Infof("synched provider config child apps")
	}
	return nil
}

func (cp *CrossPlaneApp) syncArgoCDChildApps(ctx context.Context, namespace string, apps []string) error {
	for _, appName := range apps {
		err := cp.helper.SyncArgoCDApp(ctx, namespace, appName)
		if err != nil {
			return fmt.Errorf("failed to sync app %s, %v", appName, err)
		}
		logger.Infof("synched provider config child-app %s", appName)

		err = cp.helper.WaitForArgoCDToSync(ctx, namespace, appName)
		if err != nil {
			return fmt.Errorf("failed to get sync status of app %s, %v", appName, err)
		}
	}
	return nil
}

func (cp *CrossPlaneApp) createProviderConfigs(dir string, req *model.CrossplaneUseCase) error {
	logger.Infof("processing %d crossplane providers to generate provider config", len(req.CrossplaneProviders))
	for _, provider := range req.CrossplaneProviders {
		providerName := strings.ToLower(provider.ProviderName)
		providerFile := filepath.Join(dir, fmt.Sprintf("%s-provider.yaml", providerName))
		dir := filepath.Dir(providerFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create dir %s, %v", dir, err)
		}

		file, err := os.Create(providerFile)
		if err != nil {
			return fmt.Errorf("failed to create file %s, %v", providerFile, err)
		}
		defer file.Close()

		providerConfigString, err := cp.createProviderConfigResource(provider, req)
		if err != nil {
			return fmt.Errorf("failed prepare provider %s config: %v", providerName, err)
		}

		if _, err := file.WriteString(providerConfigString); err != nil {
			return fmt.Errorf("failed to write provider %s config to %s, %v", providerName, providerFile, err)
		}
		logger.Infof("crossplane provider %s config written to %s", providerName, providerFile)
	}
	return nil
}

func (cp *CrossPlaneApp) createProviderConfigResource(provider agentmodel.CrossplaneProvider, req *model.CrossplaneUseCase) (string, error) {
	cloudType := strings.ToLower(provider.CloudType)
	pkg, found := cp.pluginConfig.ProviderPackages[cloudType]
	if !found {
		return "", fmt.Errorf("plugin package not found")
	}

	secretPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, cp.cfg.CloudProviderEntityName, provider.CloudProviderId)

	switch provider.CloudType {
	case "AWS":
		providerConfigString := fmt.Sprintf(
			crossplaneAWSProviderTemplate,
			cloudType, secretPath, secretPath,
			cloudType, pkg, cloudType,
		)
		return providerConfigString, nil
	case "GCP":
		providerConfigString := fmt.Sprintf(
			crossplaneGCPProviderTemplate,
			cloudType, secretPath, secretPath,
			cloudType, pkg, cloudType,
		)
		return providerConfigString, nil
	case "AZUR":
		providerConfigString := fmt.Sprintf(
			crossplaneAzureProviderTemplate,
			cloudType, secretPath, secretPath,
			cloudType, pkg, cloudType,
		)
		return providerConfigString, nil
	default:
		return "", fmt.Errorf("cloud type %s not supported", provider.CloudType)
	}

}

func replaceCaptenUrls(dir string, src, target string) error {
	if !strings.HasSuffix(src, ".git") {
		src += ".git"
	}

	if !strings.HasSuffix(target, ".git") {
		target += ".git"
	}

	fileList := []string{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			fileList = append(fileList, path)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, filePath := range fileList {
		err := replaceInFile(filePath, src, target)
		if err != nil {
			logger.Errorf("Error replacing in %s: %v\n", filePath, err)
		}
	}
	return nil
}

func replaceInFile(filePath, target, replacement string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	newData := strings.Replace(string(data), target, replacement, -1)
	err = os.WriteFile(filePath, []byte(newData), 0644)
	if err != nil {
		return err
	}
	return nil
}

func updateClusterEndpointDetials(filename string, req *model.CrossplaneClusterEndpoint) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return err
	}

	var argoCDAppValue model.ArgoCDAppValue

	err = json.Unmarshal(jsonData, &argoCDAppValue)
	if err != nil {
		return err
	}

	clusters := *argoCDAppValue.Clusters
	for index := range clusters {
		cluster := &clusters[index]
		if cluster.Name == req.Name {
			logger.Infof("udpated the req endpoint details to %s for name %s ", req.Endpoint, req.Name)
			cluster.Server = req.Endpoint

			break
		}
	}

	argoCDAppValue.Clusters = &clusters

	jsonBytes, err := json.Marshal(argoCDAppValue)
	if err != nil {
		return err
	}

	yamlBytes, err := k8s.ConvertJsonToYaml(jsonBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, yamlBytes, os.ModeAppend)

	return err
}
