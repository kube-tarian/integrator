package defaultplugindeployer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/gitclient"
	pluginconfigstore "github.com/kube-tarian/kad/capten/common-pkg/pluginconfig-store"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	gitProjectEntityName = "git-project"
	pluginId             = "kad-default-apps"
)

type PluginStore struct {
	log     logging.Logger
	cfg     *Config
	dbStore *captenstore.Store
	pas     *pluginconfigstore.Store
	tc      *temporalclient.Client
}

func NewPluginStore(
	log logging.Logger,
	dbStore *captenstore.Store,
	pas *pluginconfigstore.Store,
	tc *temporalclient.Client,
) (*PluginStore, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return &PluginStore{
		log:     log,
		cfg:     cfg,
		dbStore: dbStore,
		pas:     pas,
		tc:      tc,
	}, nil
}

func (p *PluginStore) SyncPlugins() error {
	pluginStoreDir, err := p.clonePluginStoreProject()
	if err != nil {
		return err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginListFilePath := p.getPluginListFilePath(pluginStoreDir)
	p.log.Infof("Loading plugin data from %s", pluginListFilePath)
	pluginListData, err := os.ReadFile(pluginListFilePath)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var plugins PluginListData
	if err := yaml.Unmarshal(pluginListData, &plugins); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	addedPlugins := map[string]bool{}
	for _, pluginName := range plugins.Plugins {
		err := p.addPluginApp(p.cfg.PluginStoreProjectID, pluginStoreDir, pluginName)
		if err != nil {
			p.log.Errorf("%v", err)
			continue
		}
		addedPlugins[pluginName] = true
		p.log.Infof("stored plugin data for plugin %s for contorl-plane cluster", pluginName)
	}

	dbPlugins, err := p.dbStore.ReadPlugins(p.cfg.PluginStoreProjectID)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	for _, dbPlugin := range dbPlugins {
		if _, ok := addedPlugins[dbPlugin.PluginName]; !ok {
			if err = p.dbStore.DeletePlugin(p.cfg.PluginStoreProjectID, dbPlugin.PluginName); err != nil {
				p.log.Infof("failed to deleted plugin data for plugin %s for control-plane cluster", dbPlugin.PluginName)
			}
			p.log.Infof("deleted plugin data for plugin %s for control-plane cluster", dbPlugin.PluginName)
		}
	}

	return nil
}

func (p *PluginStore) clonePluginStoreProject() (pluginStoreDir string, err error) {
	pluginStoreDir, err = os.MkdirTemp(p.cfg.PluginsStoreProjectMount, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create plugin store tmp dir, err: %v", err)
		return
	}

	p.log.Infof("cloning plugin store project %s to %s", p.cfg.PluginStoreProjectURL, pluginStoreDir)
	gitClient := gitclient.NewGitClient()
	if err = gitClient.Clone(pluginStoreDir, p.cfg.PluginStoreProjectURL, p.cfg.DefaultPluginsGitAccessToken); err != nil {
		os.RemoveAll(pluginStoreDir)
		err = fmt.Errorf("failed to Clone plugin store project, err: %v", err)
		return
	}
	return
}

func (p *PluginStore) addPluginApp(gitProjectId, pluginStoreDir, pluginName string) error {
	appData, err := os.ReadFile(p.getPluginFilePath(pluginStoreDir, pluginName))
	if err != nil {
		return errors.WithMessagef(err, "failed to read store plugin %s", pluginName)
	}

	var pluginData Plugin
	if err := yaml.Unmarshal(appData, &pluginData); err != nil {
		return errors.WithMessagef(err, "failed to unmarshall store plugin %s", pluginName)
	}

	var iconData []byte
	if len(pluginData.Icon) != 0 {
		iconData, err = os.ReadFile(p.getPluginIconFilePath(pluginStoreDir, pluginName, pluginData.Icon))
		if err != nil {
			return errors.WithMessagef(err, "failed to read icon %s for plugin %s", pluginData.Icon, pluginName)
		}
	}

	if pluginData.PluginName == "" || len(pluginData.Versions) == 0 {
		return fmt.Errorf("app name/version is missing for %s", pluginName)
	}

	plugin := &captenstore.PluginData{
		PluginName:  pluginData.PluginName,
		Description: pluginData.Description,
		Category:    pluginData.Category,
		Versions:    pluginData.Versions,
		Icon:        iconData,
	}

	if err := p.dbStore.UpsertPluginData(gitProjectId, plugin); err != nil {
		return errors.WithMessagef(err, "failed to store plugin %s", pluginName)
	}
	return nil
}

func (p *PluginStore) GetPlugins() ([]*captenstore.PluginData, error) {
	plugins, err := p.dbStore.ReadPlugins(p.cfg.PluginStoreProjectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return []*captenstore.PluginData{}, nil
		}
		return nil, err
	}
	return plugins, nil
}

func (p *PluginStore) GetPluginData(pluginName string) (*captenstore.PluginData, error) {
	return p.dbStore.ReadPluginData(p.cfg.PluginStoreProjectID, pluginName)
}

func (p *PluginStore) GetPluginValues(pluginName, version string) ([]byte, error) {

	pluginStoreDir, err := p.clonePluginStoreProject()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginConfig, err := p.getPluginConfig(pluginStoreDir, pluginName, version)
	if err != nil {
		return nil, err
	}

	return p.getPluginValues(pluginConfig, pluginStoreDir, pluginName, version)
}

func (p *PluginStore) getPluginValues(pluginConfig *PluginConfig, pluginStoreDir, pluginName, version string) ([]byte, error) {
	pluginValuesPath := p.getPluginDeployValuesFilePath(
		pluginStoreDir, pluginName, version,
		pluginConfig.Deployment.ControlplaneCluster.ValuesFile)
	p.log.Infof("Loading %s plugin values from %s", pluginName, pluginValuesPath)
	pluginValuesData, err := os.ReadFile(pluginValuesPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read plugins values file")
	}

	return pluginValuesData, nil
}

func (p *PluginStore) getPluginConfig(pluginStoreDir, pluginName, version string) (*PluginConfig, error) {
	pluginConfigPath := p.getPluginConfigFilePath(pluginStoreDir, pluginName, version)
	pluginConfigData, err := os.ReadFile(pluginConfigPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read store config file")
	}

	pluginConfig := &PluginConfig{}
	if err := yaml.Unmarshal(pluginConfigData, pluginConfig); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshall store config file")
	}
	if pluginConfig.Deployment.ControlplaneCluster == nil {
		return nil, errors.WithMessage(err, "no deployment found")
	}
	return pluginConfig, nil
}

func (p *PluginStore) DeployPlugins() error {
	plugins, err := p.dbStore.ReadPlugins(p.cfg.PluginStoreProjectID)
	if err != nil {
		p.log.Errorf("failed to fetch plugins, %v", err)
		return err
	}

	for _, pluginData := range plugins {
		// Check the plugin installation status
		pluginDataConfig, err := p.pas.GetPluginConfig(pluginData.PluginName)
		if err != nil {
			p.log.Infof("plugin %s not installed in control-plane cluster, will be installed", pluginData.PluginName)

			err = p.deployPlugin(pluginData)
			if err != nil {
				p.log.Errorf("failed to deploy plugin %s, %v", pluginData.PluginName, err)
			}
		}

		// Install status is failure report the error and continue
		if pluginDataConfig.InstallStatus != string(model.AppUnInstalledStatus) {
			p.log.Errorf("plugin %s deployment failed, status: %s. Troubleshoot manually", pluginData.PluginName, pluginDataConfig.InstallStatus)
			p.log.Errorf("continuing...")
			continue
		}
	}

	return nil
}

func (p *PluginStore) deployPlugin(pluginData *captenstore.PluginData) error {
	// TODO: currently taking first version in the list for deployment
	// Pick the version as per configuration in plugin metadata
	version := pluginData.Versions[0]

	pluginStoreDir, err := p.clonePluginStoreProject()
	if err != nil {
		return err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginConfig, err := p.getPluginConfig(pluginStoreDir, pluginData.PluginName, version)
	if err != nil {
		return err
	}

	validCapabilities, invalidCapabilities := filterSupporttedCapabilties(pluginConfig.Capabilities)
	if len(invalidCapabilities) > 0 {
		p.log.Infof("skipped plugin %s invalid capabilities %v", pluginData.PluginName, invalidCapabilities)
	}

	values, err := p.getPluginValues(pluginConfig, pluginStoreDir, pluginData.PluginName, version)
	if err != nil {
		p.log.Infof("no values defined for plugin %s", pluginData.PluginName)
	}

	overrideValuesMapping, overrideValuesTemplateMapping, err := p.getOverrideTemplateValues()
	if err != nil {
		return err
	}

	apiEndpoint, uiEndpoint, err := p.getPluginDataAPIValues(pluginConfig, overrideValuesTemplateMapping)
	if err != nil {
		return err
	}

	overrideValues, err := yaml.Marshal(overrideValuesMapping)
	if err != nil {
		return err
	}

	values, err = replaceTemplateValuesInByteData(values, overrideValuesMapping)
	if err != nil {
		p.log.Errorf("failed to derive template values for plugin %s, %v", pluginData.PluginName, err)
		return nil
	}

	plugin := clusterpluginspb.Plugin{

		PluginName:          pluginData.PluginName,
		Description:         pluginData.Description,
		Category:            pluginData.Category,
		Icon:                pluginData.Icon,
		Version:             version,
		ChartName:           pluginConfig.Deployment.ControlplaneCluster.ChartName,
		ChartRepo:           pluginConfig.Deployment.ControlplaneCluster.ChartRepo,
		DefaultNamespace:    pluginConfig.Deployment.ControlplaneCluster.DefaultNamespace,
		PrivilegedNamespace: pluginConfig.Deployment.ControlplaneCluster.PrivilegedNamespace,
		ApiEndpoint:         apiEndpoint,
		UiEndpoint:          uiEndpoint,
		Capabilities:        validCapabilities,
		Values:              values,
		OverrideValues:      overrideValues,
		InstallStatus:       string(model.AppIntallingStatus),
	}

	if err := p.pas.UpsertPluginConfig(&pluginconfigstore.PluginConfig{
		Plugin: plugin,
	}); err != nil {
		p.log.Errorf("failed to update plugin config data for plugin %s, %v", pluginData.PluginName, err)
		return err
	}

	wd := workers.NewDeployment(p.tc, p.log)
	_, err = wd.SendEventV2(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppInstallAction), &plugin)
	if err != nil {
		// pluginConfig.InstallStatus = string(model.AppIntallFailedStatus)
		// if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
		// 	a.log.Errorf("failed to update plugin config for plugin %s, %v", pluginConfig.PluginName, err)
		// 	return
		// }
		p.log.Errorf("sendEventV2 failed, plugin: %s, reason: %v", pluginData.PluginName, err)
		return err
	}
	// TODO: workflow will update the final status
	// Write a periodic scheduler which will go through all apps not in installed status and check the status till either success or failed.
	// Make SendEventV2 asynchrounous so that periodic scheduler will take care of monitoring.

	return nil
}

func filterSupporttedCapabilties(pluginCapabilties []string) (validCapabilties, invalidCapabilities []string) {
	validCapabilties = []string{}
	invalidCapabilities = []string{}
	for _, pluginCapability := range pluginCapabilties {
		_, ok := supporttedCapabilities[pluginCapability]
		if ok {
			validCapabilties = append(validCapabilties, pluginCapability)
		} else {
			invalidCapabilities = append(invalidCapabilities, pluginCapability)
		}
	}
	return
}

func (p *PluginStore) getPluginDataAPIValues(pluginConfig *PluginConfig, overrideValues map[string]string) (string, string, error) {
	apiEndpoint, err := replaceTemplateValuesInString(pluginConfig.ApiEndpoint, overrideValues)
	if err != nil {
		return "", "", fmt.Errorf("failed to update template values in plguin data, %v", err)
	}

	uiEndpoint, err := replaceTemplateValuesInString(pluginConfig.UIEndpoint, overrideValues)
	if err != nil {
		return "", "", fmt.Errorf("failed to update template values in plguin data, %v", err)
	}
	return apiEndpoint, uiEndpoint, nil
}

func (p *PluginStore) getOverrideTemplateValues() (map[string]any, map[string]string, error) {
	clusterGlobalValues, err := p.getClusterGlobalValues()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cluster global values, %v", err)
	}

	overrideValues := map[string]string{}
	for key, val := range clusterGlobalValues {
		overrideValues[key] = fmt.Sprintf("%v", val)

	}

	return clusterGlobalValues, overrideValues, nil
}

func (p *PluginStore) getClusterGlobalValues() (map[string]interface{}, error) {
	var globalValues map[string]interface{}
	globalValuesYaml, err := credential.GetClusterGlobalValues(context.TODO())
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read cluster global values")
	}

	err = yaml.Unmarshal([]byte(globalValuesYaml), &globalValues)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal cluster values")
	}
	p.log.Debugf("globalValues: %+v", globalValues)
	return globalValues, nil
}

func replaceTemplateValuesInByteData(data []byte,
	values map[string]interface{}) (transformedData []byte, err error) {
	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}

func replaceTemplateValuesInString(data string, values map[string]string) (transformedData string, err error) {
	if len(data) == 0 {
		return
	}

	tmpl, err := template.New("templateVal").Parse(data)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.String()

	return
}

func prepareFilePath(parts ...string) string {
	return filepath.Join(parts...)
}

func (p *PluginStore) getPluginListFilePath(parentFolder string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, p.cfg.PluginsFileName)
}

func (p *PluginStore) getPluginFilePath(parentFolder, pluginName string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, p.cfg.PluginFileName)
}

func (p *PluginStore) getPluginIconFilePath(parentFolder, pluginName, iconFileName string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, iconFileName)
}

func (p *PluginStore) getPluginConfigFilePath(parentFolder, pluginName, version string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, version, p.cfg.PluginConfigFileName)
}

func (p *PluginStore) getPluginDeployValuesFilePath(parentFolder, pluginName, version, valuesFile string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, version, valuesFile)
}
