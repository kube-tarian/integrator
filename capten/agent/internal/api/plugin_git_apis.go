package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
)

const gitProjectEntityName = "git-project"

// add more labels here if needed
var whitelistedLabels = []string{"crossplane", "tekton"}

func (a *Agent) AddGitProject(ctx context.Context, request *captenpluginspb.AddGitProjectRequest) (
	*captenpluginspb.AddGitProjectResponse, error) {
	if err := validateArgs(request.ProjectUrl, request.AccessToken); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	if len(request.UserID) == 0 {
		request.UserID = "default"
	}

	a.log.Infof("Add Git project %s request received", request.ProjectUrl)

	if err := a.validateForWhitelistedLabels(ctx, request.Labels); err != nil {
		a.log.Infof("validateForWhitelistedLabels Git project err: %s", err.Error())
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "whilelisted label validation failed, label already present",
		}, nil
	}

	id := uuid.New()
	if err := a.storeGitProjectCredential(ctx, id.String(), request.UserID, request.AccessToken); err != nil {
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add gitProject credential in vault",
		}, nil
	}

	gitProject := captenpluginspb.GitProject{
		Id:         id.String(),
		ProjectUrl: request.ProjectUrl,
		Labels:     request.Labels,
	}
	if err := a.as.UpsertGitProject(&gitProject); err != nil {
		a.log.Errorf("failed to store git project to DB, %v", err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add gitProject in db",
		}, nil
	}

	a.log.Infof("Git project %s added with id %s", request.ProjectUrl, id.String())
	return &captenpluginspb.AddGitProjectResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateGitProject(ctx context.Context, request *captenpluginspb.UpdateGitProjectRequest) (
	*captenpluginspb.UpdateGitProjectResponse, error) {
	if err := validateArgs(request.ProjectUrl, request.UserID, request.AccessToken, request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Update Git project %s, %s request recieved", request.ProjectUrl, request.Id)

	if err := a.validateForWhitelistedLabels(ctx, request.Labels); err != nil {
		a.log.Infof("validateForWhitelistedLabels Git project err: %s", err.Error())
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "whilelisted label validation failed, label already present",
		}, nil
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, nil
	}

	if err := a.storeGitProjectCredential(ctx, request.Id, request.UserID, request.AccessToken); err != nil {
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add gitProject credential in vault",
		}, nil
	}

	gitProject := captenpluginspb.GitProject{
		Id:         id.String(),
		ProjectUrl: request.ProjectUrl,
		Labels:     request.Labels,
	}
	if err := a.as.UpsertGitProject(&gitProject); err != nil {
		a.log.Errorf("failed to update gitProject in db, %v", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update gitProject in db",
		}, nil
	}

	a.log.Infof("Git project %s, %s updated", request.ProjectUrl, request.Id)
	return &captenpluginspb.UpdateGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) DeleteGitProject(ctx context.Context, request *captenpluginspb.DeleteGitProjectRequest) (
	*captenpluginspb.DeleteGitProjectResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Delete Git project %s request recieved", request.Id)

	if err := a.deleteGitProjectCredential(ctx, request.Id); err != nil {
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete gitProject credential in vault",
		}, nil
	}

	if err := a.as.DeleteGitProjectById(request.Id); err != nil {
		a.log.Errorf("failed to delete gitProject from db, %v", err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete gitProject from db",
		}, nil
	}

	a.log.Infof("Git project %s deleted", request.Id)
	return &captenpluginspb.DeleteGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) GetGitProjects(ctx context.Context, request *captenpluginspb.GetGitProjectsRequest) (
	*captenpluginspb.GetGitProjectsResponse, error) {
	a.log.Infof("Get Git projects request recieved")
	res, err := a.as.GetGitProjects()
	if err != nil {
		a.log.Errorf("failed to get gitProjects from db, %v", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git projects",
		}, nil
	}

	for _, r := range res {
		accessToken, userID, err := a.getGitProjectCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetGitProjectsResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch git projects",
			}, nil
		}
		r.AccessToken = accessToken
		r.UserID = userID
	}

	a.log.Infof("Found %d git projects", len(res))
	return &captenpluginspb.GetGitProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Projects:      res,
	}, nil

}

func (a *Agent) GetGitProjectsForLabels(ctx context.Context, request *captenpluginspb.GetGitProjectsForLabelsRequest) (
	*captenpluginspb.GetGitProjectsForLabelsResponse, error) {
	if len(request.Labels) == 0 {
		a.log.Infof("request validation failed")
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Get Git projects with labels %v request recieved", request.Labels)

	res, err := a.as.GetGitProjectsByLabels(request.Labels)
	if err != nil {
		a.log.Errorf("failed to get gitProjects for labels from db, %v", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git projects",
		}, nil
	}

	for _, r := range res {
		accessToken, userID, err := a.getGitProjectCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetGitProjectsForLabelsResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch git projects",
			}, nil
		}
		r.AccessToken = accessToken
		r.UserID = userID
	}

	a.log.Infof("Found %d git projects for lables %v", len(res), request.Labels)
	return &captenpluginspb.GetGitProjectsForLabelsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Projects:      res,
	}, nil
}

func (a *Agent) getGitProjectCredential(ctx context.Context, id string) (string, string, error) {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, gitProjectEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to get crendential for %s, %v", credPath, err)
		return "", "", err
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, gitProjectEntityName, id)
	if err != nil {
		a.log.Errorf("failed to get credential for %s, %v", credPath, err)
		return "", "", err
	}
	return cred["accessToken"], cred["userID"], nil
}

func (a *Agent) storeGitProjectCredential(ctx context.Context, id string, userID string, accessToken string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, gitProjectEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}

	credentialMap := map[string]string{
		"accessToken": accessToken,
		"userID":      userID,
	}
	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, gitProjectEntityName,
		id, credentialMap)

	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("stored credential for entity %s", credPath)
	return nil
}

func (a *Agent) deleteGitProjectCredential(ctx context.Context, id string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, gitProjectEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to delete credential for %s, %v", credPath, err)
		return err
	}

	err = credAdmin.DeleteCredential(ctx, credentials.GenericCredentialType, gitProjectEntityName, id)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to delete credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("deleted credential for entity %s", credPath)
	return nil
}

func (a *Agent) validateForWhitelistedLabels(ctx context.Context, incomingLabels []string) error {
	var filtered []string
	for _, label := range incomingLabels {
		for _, whitelisted := range whitelistedLabels {
			if label == whitelisted {
				filtered = append(filtered, label)
			}
		}
	}
	if len(filtered) == 0 {
		// safe to add this project
		return nil
	}
	res, err := a.as.GetGitProjectsByLabels(filtered)
	if err != nil {
		a.log.Errorf("failed to get gitProjects for labels from db, %v", err)
		return err
	}

	if len(res) > 0 {
		return fmt.Errorf("project present with whitelisted labels: %v", filtered)
	}

	return nil
}
