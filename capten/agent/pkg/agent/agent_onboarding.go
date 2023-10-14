package agent

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	topmodel "github.com/kube-tarian/kad/capten/model"
)

const (
	CredEntityNameOnboarding string = "onboarding"
)

func (a *Agent) SetClusterGitoptsProject(ctx context.Context, request *agentpb.SetClusterGitoptsProjectRequest) (*agentpb.SetClusterGitoptsProjectResponse, error) {
	a.log.Infof("Set Cluster Gitopts Project request for usecase %s, project %s", request.Usecase, request.ProjectUrl)
	gitOptsConfig := &model.ClusterGitoptsConfig{
		Usecase:    request.Usecase,
		ProjectUrl: request.ProjectUrl,
		Status:     "started",
	}

	if err := a.as.AddOrUpdateOnboardingIntegration(gitOptsConfig); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &agentpb.SetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Cluster Gitopts Project Set failed",
		}, err
	}

	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, CredEntityNameOnboarding, request.Usecase)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentails client for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &agentpb.SetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, CredEntityNameOnboarding,
		request.Usecase, request.Credential)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store credentail for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &agentpb.SetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}
	a.log.Audit("security", "storecred", "success", "system", "credentail stored for %s", credPath)
	a.log.Infof("Set Cluster Gitopts Project stored for usecase %s, project %s", request.Usecase, request.ProjectUrl)

	// start the config-worker routine
	go a.configureGitRepo(*gitOptsConfig)

	return &agentpb.SetClusterGitoptsProjectResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Set Cluster Gitopts Project successful",
	}, nil
}

func (a *Agent) GetClusterGitoptsProject(ctx context.Context, request *agentpb.GetClusterGitoptsProjectRequest) (
	*agentpb.GetClusterGitoptsProjectResponse, error) {
	a.log.Infof("Get Cluster Gitopts Project request for usecase %s", request.Usecase)

	resp, err := a.as.GetOnboardingIntegration(request.Usecase)
	if err != nil {
		a.log.Errorf("failed to get the Cluster Gitopts Project, %v", err)
		return &agentpb.GetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the Cluster Gitopts Project",
		}, err
	}
	a.log.Infof("Successfully fetched the the Cluster Gitopts Project. Project Url - %s", request.Usecase)

	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, CredEntityNameOnboarding, request.Usecase)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Errorf("failed to get credentail for %s, %v", credPath, err)
		return &agentpb.GetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, CredEntityNameOnboarding,
		request.Usecase)
	if err != nil {
		a.log.Errorf("failed to get credentail for %s, %v", credPath, err)
		return &agentpb.GetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	a.log.Infof("Get Cluster Gitopts Project request for usecase %s successful, project %s", request.Usecase, resp.ProjectUrl)
	return &agentpb.GetClusterGitoptsProjectResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Successfully fetched the onboarding integration",
		ClusterGitoptsConfig: &agentpb.ClusterGitoptsConfig{
			Usecase:    resp.Usecase,
			ProjectUrl: resp.ProjectUrl,
			Status:     resp.Status,
			Credential: cred,
		},
	}, nil
}

func (a *Agent) DeleteClusterGitoptsProject(ctx context.Context, request *agentpb.DeleteClusterGitoptsProjectRequest) (
	*agentpb.DeleteClusterGitoptsProjectResponse, error) {
	a.log.Infof("Delete Cluster Gitopts Project request for usecase %s and project %s", request.Usecase, request.ProjectUrl)

	if err := a.as.DeleteOnboardingIntegration(request.Usecase, request.ProjectUrl); err != nil {
		a.log.Errorf("failed to delete onboarding integration, %v", err)
		return &agentpb.DeleteClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete the onboarding integration",
		}, err
	}

	a.log.Infof("Delete Cluster Gitopts Project request for usecase %s and project %s successful", request.Usecase, request.ProjectUrl)
	return &agentpb.DeleteClusterGitoptsProjectResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Successfully deleted the onboarding integration",
	}, nil
}

func (a *Agent) configureGitRepo(req model.ClusterGitoptsConfig) {
	ci := topmodel.UseCase{Type: req.Usecase, RepoURL: req.ProjectUrl, VaultCredIdentifier: req.Usecase}
	wd := workers.NewConfig(a.tc, a.log)
	_, err := wd.SendEvent(context.TODO(), &topmodel.ConfigureParameters{Resource: req.Usecase}, ci)
	if err != nil {
		req.Status = "failed"
		if err := a.as.AddOrUpdateOnboardingIntegration(&req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.ProjectUrl, err)
		return
	}

	req.Status = "completed"
	if err := a.as.AddOrUpdateOnboardingIntegration(&req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}
	a.log.Infof("Configure onboarding Gitopts Project completed for usecase %s and project %s", req.Usecase, req.ProjectUrl)
}
