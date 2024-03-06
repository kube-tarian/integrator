package util

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

type SyncPSQLSecretConfig struct {
	DBAdminCredIdentifier string `envconfig:"DB_ADMIN_CRED_IDENTIFIER" default:"psql-admin"`
	EntityName            string `envconfig:"DB_ENTITY_NAME" required:"true"`
	SecretName            string `envconfig:"PSQL_SECRET_NAME" required:"true" default:"postgresql"`
	Namespace             string `envconfig:"POD_NAMESPACE" required:"true"`
}

func SyncPSQLAdminSecret(log logging.Logger) error {
	conf := &SyncPSQLSecretConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("psql config read failed, %v", err)
	}

	k8sClient, err := k8s.NewK8SClient(log)
	if err != nil {
		return err
	}

	res, err := k8sClient.GetSecretData(conf.Namespace, conf.SecretName)
	if err != nil {
		return err
	}

	userName := res.Data["username"]
	password := res.Data["password"]

	err = credential.PutServiceUserCredential(context.Background(), conf.EntityName, conf.DBAdminCredIdentifier, userName, password)
	return err
}