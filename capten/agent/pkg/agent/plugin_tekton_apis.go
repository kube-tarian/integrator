package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	captenmodel "github.com/kube-tarian/kad/capten/model"
)

const (
	tektonConfigUseCase string = "tekton"
)

func (a *Agent) RegisterTektonProject(ctx context.Context, request *captenpluginspb.RegisterTektonProjectRequest) (
	*captenpluginspb.RegisterTektonProjectResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Register Tekton Git project %s request recieved", request.Id)

	tektonProject, err := a.as.GetTektonProjectForID(request.Id)
	if err != nil {
		a.log.Infof("failed to get git project %s, %v", request.Id, err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	if tektonProject.Status == string(model.TektonProjectConfigurationOngoing) ||
		tektonProject.Status == string(model.TektonProjectConfigured) {
		a.log.Infof("currently the tekton project configuration on-going/configured %s, %v", request.Id, tektonProject.Status)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_OK,
			StatusMessage: "tekton configuration on-going/configured",
		}, nil
	}

	tektonProject.Status = string(model.TektonProjectConfigurationOngoing)
	if err := a.as.UpsertTektonProject(tektonProject); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to tekton db got failed",
		}, err
	}

	// start the config-worker routine
	a.configureTektonGitRepo(tektonProject)

	a.log.Infof("Tekton Git project %s registration triggerred", request.Id)
	return &captenpluginspb.RegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully registered tekton",
	}, nil
}

func (a *Agent) UnRegisterTektonProject(ctx context.Context, request *captenpluginspb.UnRegisterTektonProjectRequest) (
	*captenpluginspb.UnRegisterTektonProjectResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("UnRegister Tekton Git project %s request recieved", request.Id)

	tektonProject, err := a.as.GetTektonProjectForID(request.Id)
	if err != nil {
		a.log.Infof("faile to get git project %s, %v", request.Id, err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	tektonProject.Status = string(model.TektonProjectAvailable)
	if err := a.as.UpsertTektonProject(tektonProject); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to tekton db got failed",
		}, err
	}

	a.log.Infof("UnRegister Tekton Git project %s request processed", request.Id)
	return &captenpluginspb.UnRegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully delete the tekton project",
	}, nil
}

func (a *Agent) GetTektonProjects(ctx context.Context, request *captenpluginspb.GetTektonProjectsRequest) (
	*captenpluginspb.GetTektonProjectsResponse, error) {
	a.log.Infof("Get Tekton Git projects request recieved")

	projects, err := a.as.GetTektonProjects()
	if err != nil {
		a.log.Errorf("failed to get tekton Project, %v", err)
		return &captenpluginspb.GetTektonProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get tekton Project",
		}, err
	}

	tekTonProjects := []*captenpluginspb.TektonProject{}
	for _, project := range projects {
		tekTonProject := &captenpluginspb.TektonProject{
			Id:            project.Id,
			GitProjectUrl: project.GitProjectUrl,
			Status:        project.Status,
		}
		tekTonProjects = append(tekTonProjects, tekTonProject)
	}

	a.log.Infof("Fetched %d Tekton Git projects", len(tekTonProjects))
	return &captenpluginspb.GetTektonProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the tekton projects",
		Projects:      tekTonProjects,
	}, nil
}

func (a *Agent) configureTektonGitRepo(req *model.TektonProject) {
	ci := captenmodel.UseCase{Type: tektonConfigUseCase, RepoURL: req.GitProjectUrl,
		VaultCredIdentifier: req.Id, PushToDefaultBranch: !a.createPr}
	wd := workers.NewConfig(a.tc, a.log)

	wkfId, err := wd.SendAsyncEvent(context.TODO(), &captenmodel.ConfigureParameters{Resource: tektonConfigUseCase}, ci)
	if err != nil {
		req.Status = string(model.TektonProjectConfigurationFailed)
		req.WorkflowId = "NA"
		if err := a.as.UpsertTektonProject(req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	a.log.Infof("Tekton Git project %s config workflow event %s created", wkfId)

	req.Status = string(model.TektonProjectConfigured)
	req.WorkflowId = wkfId
	req.WorkflowStatus = string(model.WorkFlowStatusStarted)
	if err := a.as.UpsertTektonProject(req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}

	go a.monitorWorkflow(req, wkfId)
	a.log.Infof("Tekton Git project %s registration completed", req.Id)
}

func (a *Agent) monitorWorkflow(req *model.TektonProject, wkfId string) {
	// during system reboot start monitoring, add it in map or somewhere.
	wd := workers.NewConfig(a.tc, a.log)
	wkfResp, err := wd.GetWorkflowInformation(context.TODO(), wkfId)
	if err != nil {
		req.Status = string(model.TektonProjectConfigurationFailed)
		req.WorkflowStatus = string(model.WorkFlowStatusFailed)
		if err := a.as.UpsertTektonProject(req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	a.log.Infof("Monitoring Tekton Git project %s config workflow event %s created", wkfId)

	req.Status = string(model.TektonProjectConfigured)
	req.WorkflowStatus = wkfResp.Status
	if err := a.as.UpsertTektonProject(req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}

	a.log.Infof("Tekton Git project %s monitoring completed", req.Id)
}
