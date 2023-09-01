package agent

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"reflect"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/credential"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (a *Agent) ConfigureAppSSO(
	ctx context.Context, req *agentpb.ConfigureAppSSORequest) (*agentpb.ConfigureAppSSOResponse, error) {
	if req.ReleaseName == "" {
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "release name empty",
		}, nil
	}
	a.log.Infof("Received request for ConfigureAppSSO, app %s", req.ReleaseName)

	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to GetAppConfig for release_name: %s err: %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err fetching appConfig").Error(),
		}, nil
	}

	if err := credential.StoreAppOauthCredential(ctx, req.ReleaseName, req.ClientId, req.ClientSecret); err != nil {
		a.log.Errorf("failed to store credential for ClientId: %s, %v", req.ClientId, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err saving SSO credentials in vault").Error(),
		}, nil
	}

	ssoOverwriteMapping := map[string]any{
		"ClientId":     req.ClientId,
		"ClientSecret": req.ClientSecret,
		"OAuthBaseURL": req.OAuthBaseURL,
	}

	templateValuesMapping, err := deriveTemplateValuesMapping(appConfig.Values.OverrideValues, appConfig.Values.TemplateValues)
	if err != nil {
		a.log.Errorf("failed to derivee template values, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to dervice template values").Error(),
		}, nil
	}

	launchUiMapping := map[string]any{}
	if err := yaml.Unmarshal(appConfig.Values.LaunchUIValues, &launchUiMapping); err != nil {
		a.log.Errorf("failed to Unmarshal LaunchUIValues: %s err: %v", string(appConfig.Values.LaunchUIValues), err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err Unmarshalling launchiUiValues").Error(),
		}, nil
	}

	launchUiMapping, err = replaceTemplateValues(launchUiMapping, ssoOverwriteMapping)
	if err != nil {
		a.log.Errorf("failed to replaceTemplateValues, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err replacing launchUiMapping").Error(),
		}, nil
	}

	finalOverrideValuesMapping := mergeRecursive(convertKey(templateValuesMapping), convertKey(launchUiMapping))
	marshaledOverrideValues, err := yaml.Marshal(finalOverrideValuesMapping)
	if err != nil {
		a.log.Errorf("failed to Marshal, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err marshalling overrideValues").Error(),
		}, nil
	}

	newAppConfig := *appConfig
	newAppConfig.Values.OverrideValues = marshaledOverrideValues
	newAppConfig.Config.InstallStatus = "Updating"

	if err := a.as.UpsertAppConfig(&newAppConfig); err != nil {
		a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err upserting new appConfig").Error(),
		}, nil
	}

	wd := workers.NewDeployment(a.tc, a.log)
	run, err := wd.SendEvent(context.TODO(), "update", installRequestFromSyncApp(&newAppConfig))
	if err != nil {
		newAppConfig.Config.InstallStatus = "Update Failed"
		if err := a.as.UpsertAppConfig(&newAppConfig); err != nil {
			a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
			return &agentpb.ConfigureAppSSOResponse{
				Status:        agentpb.StatusCode_INTERNRAL_ERROR,
				StatusMessage: errors.WithMessage(err, "err upserting new appConfig").Error(),
			}, nil
		}

		a.log.Errorf("failed to SendEvent, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err sending deployment event").Error(),
		}, nil
	}

	newAppConfig.Config.InstallStatus = "Updated"
	if err := a.as.UpsertAppConfig(&newAppConfig); err != nil {
		a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err upserting new appConfig").Error(),
		}, nil
	}

	return &agentpb.ConfigureAppSSOResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: fmt.Sprintf("app deployment scheduled, runId: %v", run.GetRunID()),
	}, nil

}

func replaceTemplateValues(templateData, values map[string]any) (transformedData map[string]any, err error) {
	yamlData, err := yaml.Marshal(templateData)
	if err != nil {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(yamlData))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = map[string]any{}
	err = yaml.Unmarshal(buf.Bytes(), &transformedData)
	if err != nil {
		return
	}
	return
}

// merge map[any]T and map[any]T where T => map[any]T | any
func mergeRecursive(original, override map[any]any) map[any]any {
	if override == nil {
		return original
	}
	if original == nil {
		original = map[any]any{}
	}
	for k, v := range override {
		// case 1: value not found in original
		if _, found := original[k]; !found {
			original[k] = v
			continue
		}

		// case 2: both are not maps
		if reflect.TypeOf(original[k]).Kind() != reflect.Map &&
			reflect.TypeOf(v).Kind() != reflect.Map {
			original[k] = v
			continue
		}

		// case 3: both are maps and v is not nil
		if reflect.TypeOf(v) != nil {
			original[k] = mergeRecursive(
				original[k].(map[any]any),
				v.(map[any]any),
			)
		}

	}
	return original
}

func convertKey(m map[string]any) map[any]any {
	ret := map[any]any{}
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func installRequestFromSyncApp(data *agentpb.SyncAppData) *model.ApplicationInstallRequest {
	values := make([]byte, len(data.Values.OverrideValues))
	copy(values, data.Values.OverrideValues)
	return &model.ApplicationInstallRequest{
		PluginName:     "helm",
		RepoName:       data.Config.RepoName,
		RepoURL:        data.Config.RepoURL,
		ChartName:      data.Config.ChartName,
		Namespace:      data.Config.Namespace,
		ReleaseName:    data.Config.ReleaseName,
		Version:        data.Config.Version,
		ClusterName:    "capten",
		OverrideValues: string(values),
		Timeout:        10,
	}
}
