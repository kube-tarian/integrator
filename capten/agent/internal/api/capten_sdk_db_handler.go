package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/db"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
)

func (a *Agent) PostSetupdatabase(c *gin.Context) {
	a.log.Info("Creating new db for configuration")

	req := &api.SetupDatabaseRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	vaultPath, err := a.setupDatabase(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, api.SetupDatabaseResponse{
		Status:        api.OK,
		StatusMessage: "Database setup in postgres succesful",
		VaultPath:     vaultPath,
	})

}

func (a *Agent) setupDatabase(req *api.SetupDatabaseRequest) (vaultPath string, err error) {
	switch req.DbOemName {
	case db.POSTGRES.String():
		vaultPath, err = setupPostgresDatabase(a.log, req)
	default:
		err = fmt.Errorf("unsupported Database OEM %s", req.DbOemName)
	}
	return vaultPath, err
}

// Read the Postgres DB configuration
func readConfig() (*dbinit.Config, error) {
	var baseConfig dbinit.BaseConfig
	if err := envconfig.Process("", &baseConfig); err != nil {
		return nil, err
	}
	return &dbinit.Config{
		BaseConfig: baseConfig,
	}, nil
}

func setupPostgresDatabase(log logging.Logger, req *api.SetupDatabaseRequest) (vaultPath string, err error) {
	conf, err := readConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	conf.DBName = req.DbName
	conf.DBServiceUsername = req.ServiceUserName
	conf.Password = dbinit.GenerateRandomPassword(12)

	err = dbinit.CreatedDatabaseWithConfig(log, conf)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Insert into vault path plugin/<plugin-name>/<svc-entity> => plugin/test/postgres
	cred := credentials.PrepareServiceCredentialMap(credentials.ServiceCredential{
		UserName: conf.DBServiceUsername,
		Password: conf.Password,
		AdditionalData: map[string]string{
			"db-url":       conf.DBAddress,
			"db-name":      conf.DBName,
			"service-user": req.ServiceUserName,
		},
	})
	return fmt.Sprintf("%s/%s/%s", credentials.CertCredentialType, req.PluginName, conf.EntityName),
		credential.PutPluginCredential(context.TODO(), req.PluginName, conf.EntityName, cred)
}
