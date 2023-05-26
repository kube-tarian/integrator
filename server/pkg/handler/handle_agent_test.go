package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestAPIHandler_Close(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		customerId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.Close(tt.args.customerId)
		})
	}
}

func TestAPIHandler_CloseAll(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.CloseAll()
		})
	}
}

func TestAPIHandler_ConnectClient(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		customerId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			if err := a.ConnectClient(tt.args.customerId); (err != nil) != tt.wantErr {
				t.Errorf("ConnectClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAPIHandler_DeleteAgentClimondeploy(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.DeleteAgentClimondeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentCluster(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.DeleteAgentCluster(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentDeploy(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.DeleteAgentDeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentProject(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.DeleteAgentProject(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentRepository(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.DeleteAgentRepository(tt.args.c)
		})
	}
}

func TestAPIHandler_GetAgentEndpoint(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.GetAgentEndpoint(tt.args.c)
		})
	}
}

func TestAPIHandler_GetApiDocs(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.GetApiDocs(tt.args.c)
		})
	}
}

func TestAPIHandler_GetClient(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		customerId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *client.Agent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			if got := a.GetClient(tt.args.customerId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_GetStatus(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.GetStatus(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentApps(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}

	_ = log.New("debug")

	chartName := "argocd"
	name := "argocd"
	//override := ""
	releaseName := "test"
	repoName := "test"
	repoURL := "https://argocd.com"
	version := "v1.0.0"
	namespace := "capten"
	tools := []struct {
		ChartName   *string `json:"chartName,omitempty"`
		Name        *string `json:"name,omitempty"`
		Namespace   *string `json:"namespace,omitempty"`
		Override    *string `json:"override,omitempty"`
		ReleaseName *string `json:"releaseName,omitempty"`
		RepoName    *string `json:"repoName,omitempty"`
		RepoURL     *string `json:"repoURL,omitempty"`
		Version     *string `json:"version,omitempty"`
	}{
		{
			ChartName:   &chartName,
			Name:        &name,
			ReleaseName: &releaseName,
			RepoName:    &repoName,
			RepoURL:     &repoURL,
			Version:     &version,
			Namespace:   &namespace,
		},
	}

	apps := api.AgentAppsRequest{
		Apps: &tools,
	}

	jsonByte, err := json.Marshal(apps)
	require.NoError(t, err)
	fmt.Println(string(jsonByte))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("customer_id", "1")
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonByte))
	fmt.Println(c.Request.Body)
	agentConn, err := client.NewAgent(&types.AgentConfiguration{
		Address:    "127.0.0.1",
		Port:       9091,
		TlsEnabled: false,
	})

	require.NoError(t, err)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "post apps",
			fields: fields{agents: map[string]*client.Agent{
				"1": agentConn,
			}},
			args: args{c: c},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentApps(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentClimondeploy(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentClimondeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentCluster(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentCluster(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentDeploy(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentDeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentEndpoint(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentEndpoint(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentProject(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentProject(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentRepository(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentRepository(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentSecret(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PostAgentSecret(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentClimondeploy(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PutAgentClimondeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentDeploy(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PutAgentDeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentEndpoint(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PutAgentEndpoint(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentProject(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PutAgentProject(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentRepository(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.PutAgentRepository(tt.args.c)
		})
	}
}

func TestAPIHandler_getFileContent(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c        *gin.Context
		fileInfo map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			got, err := a.getFileContent(tt.args.c, tt.args.fileInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_sendResponse(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c   *gin.Context
		msg string
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.sendResponse(tt.args.c, tt.args.msg, tt.args.err)
		})
	}
}

func TestAPIHandler_setFailedResponse(t *testing.T) {
	type fields struct {
		agents map[string]*client.Agent
	}
	type args struct {
		c   *gin.Context
		msg string
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agents: tt.fields.agents,
			}
			a.setFailedResponse(tt.args.c, tt.args.msg, tt.args.err)
		})
	}
}

func TestNewAPIHandler(t *testing.T) {
	tests := []struct {
		name    string
		want    *APIHandler
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAPIHandler()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPIHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAgentConfig(t *testing.T) {
	type args struct {
		customerID string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.AgentConfiguration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAgentConfig(tt.args.customerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAgentConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAgentConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type args struct {
		resp *agentpb.JobResponse
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toString(tt.args.resp); got != tt.want {
				t.Errorf("toString() = %v, want %v", got, tt.want)
			}
		})
	}
}
