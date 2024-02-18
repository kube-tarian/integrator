package storeapps

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AppStoreAppConfigPath string `envconfig:"APP_STORE_APP_CONFIG_PATH" default:"/data/store-apps/conf"`
	AppStoreAppIconsPath  string `envconfig:"APP_STORE_APP_ICONS_PATH" default:"/data/store-apps/icons"`
	SyncAppStore          bool   `envconfig:"SYNC_APP_STORE" default:"true"`
	AppStoreConfigFile    string `envconfig:"APP_STORE_CONFIG_FILE" default:"/data/store-apps/app_list.yaml"`
}

type AppStoreConfig struct {
	EnabledApps []string `yaml:"enabledApps"`
}

type AppConfig struct {
	Name                string                 `yaml:"Name"`
	ChartName           string                 `yaml:"ChartName"`
	Category            string                 `yaml:"Category"`
	RepoName            string                 `yaml:"RepoName"`
	RepoURL             string                 `yaml:"RepoURL"`
	Namespace           string                 `yaml:"Namespace"`
	ReleaseName         string                 `yaml:"ReleaseName"`
	Version             string                 `yaml:"Version"`
	Description         string                 `yaml:"Description"`
	LaunchURL           string                 `yaml:"LaunchURL"`
	LaunchUIDescription string                 `yaml:"LaunchUIDescription"`
	LaunchUIIcon        string                 `yaml:"LaunchUIIcon"`
	LaunchUIValues      map[string]interface{} `yaml:"LaunchUIValues"`
	OverrideValues      map[string]interface{} `yaml:"OverrideValues"`
	CreateNamespace     bool                   `yaml:"CreateNamespace"`
	PrivilegedNamespace bool                   `yaml:"PrivilegedNamespace"`
	TemplateValues      map[string]interface{} `yaml:"TemplateValues"`
	Icon                string                 `yaml:"Icon"`
	PluginName          string                 `yaml:"PluginName"`
	PluginDescription   string                 `yaml:"PluginDescription"`
	APIEndpoint         string                 `yaml:"APIEndpoint"`
}

func SyncStoreApps(log logging.Logger, appStore store.ServerStore) error {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return err
	}

	if !cfg.SyncAppStore {
		log.Info("app store config synch disabled")
		return nil
	}

	configData, err := os.ReadFile(cfg.AppStoreConfigFile)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var config AppStoreConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	enabledApps := []AppConfig{}
	for _, appName := range config.EnabledApps {
		appData, err := os.ReadFile(cfg.AppStoreAppConfigPath + "/" + appName + ".yaml")
		if err != nil {
			return errors.WithMessagef(err, "failed to read store app config for %s", appName)
		}

		var appConfig AppConfig
		if err := yaml.Unmarshal(appData, &appConfig); err != nil {
			return errors.WithMessagef(err, "failed to unmarshall store app config for %s", appName)
		}

		if appConfig.Name == "" || appConfig.Version == "" {
			return fmt.Errorf("app name/version is missing for %s", appName)
		}

		storeAppConfig := &types.StoreAppConfig{
			AppName:             appConfig.Name,
			Version:             appConfig.Version,
			Category:            appConfig.Category,
			Description:         appConfig.Description,
			ChartName:           appConfig.ChartName,
			RepoName:            appConfig.RepoName,
			ReleaseName:         appConfig.ReleaseName,
			RepoURL:             appConfig.RepoURL,
			Namespace:           appConfig.Namespace,
			CreateNamespace:     appConfig.CreateNamespace,
			PrivilegedNamespace: appConfig.PrivilegedNamespace,
			LaunchURL:           appConfig.LaunchURL,
			LaunchUIDescription: appConfig.LaunchUIDescription,
			PluginName:          appConfig.PluginName,
			PluginDescription:   appConfig.PluginDescription,
			APIEndpoint:         appConfig.APIEndpoint,
		}

		if len(appConfig.LaunchUIIcon) != 0 {
			iconBytes, err := os.ReadFile(cfg.AppStoreAppIconsPath + "/" + appConfig.Icon)
			if err != nil {
				return fmt.Errorf("failed loading icon for app '%s', %v", appConfig.ReleaseName, err)
			}
			storeAppConfig.Icon = hex.EncodeToString(iconBytes)
		}

		if len(appConfig.OverrideValues) > 0 {
			marshaledOverride, err := yaml.Marshal(appConfig.OverrideValues)
			if err != nil {
				return errors.WithMessage(err, "override values marshal error")
			}
			storeAppConfig.OverrideValues = marshaledOverride
		}
		if len(appConfig.LaunchUIValues) > 0 {
			marshaledOLaunchUI, err := yaml.Marshal(appConfig.LaunchUIValues)
			if err != nil {
				return errors.WithMessage(err, "launchui values marshal error")
			}
			storeAppConfig.LaunchUIValues = marshaledOLaunchUI
		}

		storeAppConfig.TemplateValues = getAppTemplateValues(log, cfg, appName)

		if err := appStore.AddOrUpdateStoreApp(storeAppConfig); err != nil {
			return errors.WithMessagef(err, "failed to store app config for %s", appName)
		}
		enabledApps = append(enabledApps, appConfig)
	}

	deleteDisabledApps(log, appStore, enabledApps)
	return nil
}

func deleteDisabledApps(log logging.Logger, appStore store.ServerStore, enabledApps []AppConfig) {
	storedApps, err := appStore.GetAppsFromStore()
	if err != nil {
		log.Errorf("failed to get stored apps, %w", err)
		return
	}

	for _, app := range storedApps {
		if isAppEnabled(enabledApps, app.Name, app.Version) {
			continue
		}

		err := appStore.DeleteAppInStore(app.Name, app.Version)
		if err != nil {
			log.Errorf("failed to delete app %s:%s, %w", app.Name, app.Version, err)
		}
		log.Infof("Deleted store app %, version %s", app.Name, app.Version)
	}
}

func getAppTemplateValues(log logging.Logger, cfg *Config, appName string) []byte {
	templateValues, err := os.ReadFile(cfg.AppStoreAppConfigPath + "/values/" + appName + "_template.yaml")
	if err != nil {
		log.Infof("No template file for app %s", appName)
	}
	return templateValues
}

func isAppEnabled(apps []AppConfig, releaseName, version string) bool {
	for _, app := range apps {
		if app.ReleaseName == releaseName && app.Version == version {
			return true
		}
	}
	return false
}
