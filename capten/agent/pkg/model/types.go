package model

import (
	"encoding/json"
	"fmt"
)

type WorkFlowStatus string

const (
	WorkFlowStatusStarted    WorkFlowStatus = "started"
	WorkFlowStatusCompleted  WorkFlowStatus = "completed"
	WorkFlowStatusInProgress WorkFlowStatus = "in-progress"
	WorkFlowStatusFailed     WorkFlowStatus = "failed"
)

type ArgoCDProjectStatus string

const (
	ArgoCDProjectAvailable           ArgoCDProjectStatus = "available"
	ArgoCDProjectConfigured          ArgoCDProjectStatus = "configured"
	ArgoCDProjectConfigurationFailed ArgoCDProjectStatus = "configuration-failed"
)

type TektonProjectStatus string

const (
	TektonProjectAvailable            TektonProjectStatus = "available"
	TektonProjectConfigured           TektonProjectStatus = "configured"
	TektonProjectConfigurationOngoing TektonProjectStatus = "configuration-ongoing"
	TektonProjectConfigurationFailed  TektonProjectStatus = "configuration-failed"
)

type CrossplaneProjectStatus string

const (
	CrossplaneProjectAvailable            CrossplaneProjectStatus = "available"
	CrossplaneProjectConfigured           CrossplaneProjectStatus = "configured"
	CrossplaneProjectConfigurationOngoing CrossplaneProjectStatus = "configuration-ongoing"
	CrossplaneProjectConfigurationFailed  CrossplaneProjectStatus = "configuration-failed"
)

type AppConfig struct {
	AppName             string `json:"AppName,omitempty"`
	Version             string `json:"Version,omitempty"`
	Category            string `json:"Category,omitempty"`
	Description         string `json:"Description,omitempty"`
	ChartName           string `json:"ChartName,omitempty"`
	RepoName            string `json:"RepoName,omitempty"`
	ReleaseName         string `json:"ReleaseName,omitempty"`
	RepoURL             string `json:"RepoURL,omitempty"`
	Namespace           string `json:"Namespace,omitempty"`
	CreateNamespace     bool   `json:"CreateNamespace"`
	PrivilegedNamespace bool   `json:"PrivilegedNamespace"`
	Icon                string `json:"Icon,omitempty"`
	LaunchURL           string `json:"LaunchURL,omitempty"`
	LaunchUIDescription string `json:"LaunchUIDescription,omitempty"`
}

type ApplicationInstallRequest struct {
	PluginName     string `json:"PluginName,omitempty"`
	RepoName       string `json:"RepoName,omitempty"`
	RepoURL        string `json:"RepoURL,omitempty"`
	ChartName      string `json:"ChartName,omitempty"`
	Namespace      string `json:"Namespace,omitempty"`
	ReleaseName    string `json:"ReleaseName,omitempty"`
	Timeout        uint32 `json:"Timeout,omitempty"`
	Version        string `json:"Version,omitempty"`
	ClusterName    string `json:"ClusterName,omitempty"`
	OverrideValues string `json:"OverrideValues,omitempty"`
}

type ApplicationDeleteRequest struct {
	PluginName  string `json:"plugin_name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	ReleaseName string `json:"release_name,omitempty"`
	Timeout     uint32 `json:"timeout,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
}

type ResponsePayload struct {
	Status  string          `json:"status"`
	Message json.RawMessage `json:"message,omitempty"` // TODO: This will be enhanced along with plugin implementation
}

func (rsp *ResponsePayload) ToString() string {
	return fmt.Sprintf("Status: %s, Message: %s", rsp.Status, string(rsp.Message))
}

type ClusterGitoptsConfig struct {
	Usecase    string `json:"usecase,omitempty"`
	ProjectUrl string `json:"project_url,omitempty"`
	Status     string `json:"status,omitempty"`
}

type TektonProject struct {
	Id             string `json:"id,omitempty"`
	GitProjectId   string `json:"git_project_id,omitempty"`
	GitProjectUrl  string `json:"git_project_url,omitempty"`
	Status         string `json:"status,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
	WorkflowId     string `json:"workflow_id,omitempty"`
	WorkflowStatus string `json:"workflow_status,omitempty"`
}

type CrossplaneProject struct {
	Id             string `json:"id,omitempty"`
	GitProjectId   string `json:"git_project_id,omitempty"`
	GitProjectUrl  string `json:"git_project_url,omitempty"`
	Status         string `json:"status,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
	WorkflowId     string `json:"workflow_id,omitempty"`
	WorkflowStatus string `json:"workflow_status,omitempty"`
}

type ArgoCDProject struct {
	Id             string `json:"id,omitempty"`
	GitProjectId   string `json:"git_project_id,omitempty"`
	GitProjectUrl  string `json:"git_project_url,omitempty"`
	Status         string `json:"status,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
}

type CrossplaneProvider struct {
	Id              string `json:"id,omitempty"`
	CloudType       string `json:"cloud_type,omitempty"`
	ProviderName    string `json:"provider_name,omitempty"`
	CloudProviderId string `json:"cloud_provider_id,omitempty"`
	Status          string `json:"status,omitempty"`
}
