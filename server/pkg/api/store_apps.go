package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/types"
)

func (s *Server) AddStoreApp(ctx context.Context, request *serverpb.AddStoreAppRequest) (
	*serverpb.AddStoreAppResponse, error) {

	if request.AppConfig.AppName == "" || request.AppConfig.Version == "" {
		s.log.Errorf("failed to add app config to store, %v", "App name/version is missing")
		return &serverpb.AddStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed add app config to store, app name/version is missing",
		}, nil
	}

	config := &types.StoreAppConfig{
		ReleaseName:         request.AppConfig.ReleaseName,
		AppName:             request.AppConfig.AppName,
		Version:             request.AppConfig.Version,
		Category:            request.AppConfig.Category,
		Description:         request.AppConfig.Description,
		ChartName:           request.AppConfig.ChartName,
		RepoName:            request.AppConfig.RepoName,
		RepoURL:             request.AppConfig.RepoURL,
		Namespace:           request.AppConfig.Namespace,
		CreateNamespace:     request.AppConfig.CreateNamespace,
		PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
		Icon:                request.AppConfig.Icon,
		LaunchURL:           request.AppConfig.LaunchURL,
		LaunchUIDescription: request.AppConfig.LaunchUIDescription,
		OverrideValues:      request.AppValues.OverrideValues,
		LaunchUIValues:      request.AppValues.LaunchUIValues,
	}

	if err := s.serverStore.AddOrUpdateApp(config); err != nil {
		s.log.Errorf("failed to add app config to store, %v", err)
		return &serverpb.AddStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed add app config to store",
		}, nil
	}

	return &serverpb.AddStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly added to store",
	}, nil
}

func (s *Server) UpdateStoreApp(ctx context.Context, request *serverpb.UpdateStoreAppRequest) (
	*serverpb.UpdateStoreAppRsponse, error) {
	if request.AppConfig.AppName == "" || request.AppConfig.Version == "" {
		s.log.Errorf("failed to update app config in store, %v", "App name/version is missing")
		return &serverpb.UpdateStoreAppRsponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config in store, app name/version is missing",
		}, nil
	}

	config := &types.StoreAppConfig{
		ReleaseName:         request.AppConfig.ReleaseName,
		AppName:             request.AppConfig.AppName,
		Version:             request.AppConfig.Version,
		Category:            request.AppConfig.Category,
		Description:         request.AppConfig.Description,
		ChartName:           request.AppConfig.ChartName,
		RepoName:            request.AppConfig.RepoName,
		RepoURL:             request.AppConfig.RepoURL,
		Namespace:           request.AppConfig.Namespace,
		CreateNamespace:     request.AppConfig.CreateNamespace,
		PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
		Icon:                request.AppConfig.Icon,
		LaunchURL:           request.AppConfig.LaunchURL,
		LaunchUIDescription: request.AppConfig.LaunchUIDescription,
		OverrideValues:      request.AppValues.OverrideValues,
		LaunchUIValues:      request.AppValues.LaunchUIValues,
	}

	if err := s.serverStore.AddOrUpdateApp(config); err != nil {
		s.log.Errorf("failed to update app config in store, %v", err)
		return &serverpb.UpdateStoreAppRsponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config in store",
		}, nil
	}

	return &serverpb.UpdateStoreAppRsponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly updated",
	}, nil
}

func (s *Server) DeleteStoreApp(ctx context.Context, request *serverpb.DeleteStoreAppRequest) (
	*serverpb.DeleteStoreAppResponse, error) {
	if request.AppName == "" || request.Version == "" {
		s.log.Errorf("failed to delete app config from store, %v", "App name/version is missing")
		return &serverpb.DeleteStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete app config from store, app name/version is missing",
		}, nil
	}

	if err := s.serverStore.DeleteAppInStore(request.AppName, request.Version); err != nil {
		s.log.Errorf("failed to delete app config from store, %v", err)
		return &serverpb.DeleteStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete app config from store",
		}, nil
	}

	return &serverpb.DeleteStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly deleted",
	}, nil

}

func (s *Server) GetStoreApp(ctx context.Context, request *serverpb.GetStoreAppRequest) (
	*serverpb.GetStoreAppResponse, error) {
	if request.AppName == "" || request.Version == "" {
		s.log.Errorf("failed to get app config from store, %v", "App name/version is missing")
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config from store, app name/version is missing",
		}, nil
	}
	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get app config from store, %v", err)
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config from store",
		}, nil
	}

	appConfig := &serverpb.StoreAppConfig{
		AppName:             config.Name,
		Version:             config.Version,
		Category:            config.Category,
		Description:         config.Description,
		ChartName:           config.ChartName,
		RepoName:            config.RepoName,
		RepoURL:             config.RepoURL,
		Namespace:           config.Namespace,
		CreateNamespace:     config.CreateNamespace,
		PrivilegedNamespace: config.PrivilegedNamespace,
		Icon:                config.Icon,
		LaunchURL:           config.LaunchUIURL,
		LaunchUIDescription: config.LaunchUIDescription,
		ReleaseName:         config.ReleaseName,
	}

	return &serverpb.GetStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly fetched from store",
		AppConfig:     appConfig,
	}, nil

}

func (s *Server) GetStoreApps(ctx context.Context, request *serverpb.GetStoreAppsRequest) (
	*serverpb.GetStoreAppsResponse, error) {

	configs, err := s.serverStore.GetAppsFromStore()
	if err != nil {
		s.log.Errorf("failed to get app config's from store, %v", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config's from store",
		}, nil
	}

	appConfigs := []*serverpb.StoreAppConfig{}
	for _, config := range *configs {
		appConfigs = append(appConfigs, &serverpb.StoreAppConfig{
			AppName:             config.Name,
			Version:             config.Version,
			Category:            config.Category,
			Description:         config.Description,
			ChartName:           config.ChartName,
			RepoName:            config.RepoName,
			RepoURL:             config.RepoURL,
			Namespace:           config.Namespace,
			CreateNamespace:     config.CreateNamespace,
			PrivilegedNamespace: config.PrivilegedNamespace,
			Icon:                config.Icon,
			LaunchURL:           config.LaunchUIURL,
			LaunchUIDescription: config.LaunchUIDescription,
			ReleaseName:         config.ReleaseName,
		})
	}

	return &serverpb.GetStoreAppsResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config's are sucessfuly fetched from store",
		AppConfigs:    appConfigs,
	}, nil
}

func (s *Server) GetStoreAppValues(ctx context.Context, request *serverpb.GetStoreAppValuesRequest) (
	*serverpb.GetStoreAppValuesResponse, error) {
	if request.AppName == "" || request.Version == "" {
		s.log.Errorf("failed to get store app values, %v", "App name/version is missing")
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get store app values, app name/version is missing",
		}, nil
	}
	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get store app values, %v", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get store app values",
		}, nil
	}

	appConfig := &serverpb.StoreAppConfig{
		AppName:             config.Name,
		Version:             config.Version,
		Category:            config.Category,
		Description:         config.Description,
		ChartName:           config.ChartName,
		RepoName:            config.RepoName,
		RepoURL:             config.RepoURL,
		Namespace:           config.Namespace,
		CreateNamespace:     config.CreateNamespace,
		PrivilegedNamespace: config.PrivilegedNamespace,
		Icon:                config.Icon,
		LaunchURL:           config.LaunchUIURL,
		LaunchUIDescription: config.LaunchUIDescription,
		ReleaseName:         config.ReleaseName,
	}

	return &serverpb.GetStoreAppValuesResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "store app values sucessfuly fetched",
		AppConfig:     appConfig,
	}, nil

}

func (s *Server) DeployStoreApp(ctx context.Context, request *serverpb.DeployStoreAppRequest) (
	*serverpb.DeployStoreAppResponse, error) {
	if request.AppConfig.AppName == "" || request.AppConfig.Version == "" {
		s.log.Errorf("failed to get store app values, %v", "App name/version is missing")
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get store app values, app name/version is missing",
		}, nil
	}

	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}
	clusterId := metadataMap[clusterIDAttribute]
	if orgId == "" {
		s.log.Errorf("cluster Id is missing in the request")
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "cluster Id is missing",
		}, nil

	}

	agent, err := s.agentHandeler.GetAgent(clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy the app",
		}, nil
	}

	req := &agentpb.InstallAppRequest{
		AppConfig: &agentpb.AppConfig{
			AppName:             request.AppConfig.AppName,
			Version:             request.AppConfig.Version,
			ReleaseName:         request.AppConfig.ReleaseName,
			Category:            request.AppConfig.Category,
			Description:         request.AppConfig.Description,
			ChartName:           request.AppConfig.ChartName,
			RepoName:            request.AppConfig.RepoName,
			RepoURL:             request.AppConfig.RepoURL,
			Namespace:           request.AppConfig.Namespace,
			CreateNamespace:     request.AppConfig.CreateNamespace,
			PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
			Icon:                []byte(request.AppConfig.Icon),
			LaunchURL:           request.AppConfig.LaunchURL,
			LaunchUIDescription: request.AppConfig.LaunchUIDescription,
		},
	}

	_, err = agent.GetClient().InstallApp(ctx, req)
	if err != nil {
		s.log.Errorf("failed to deploy app, %v", err)
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy the app",
		}, nil
	}

	return &serverpb.DeployStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app is successfully deployed",
	}, nil

}
