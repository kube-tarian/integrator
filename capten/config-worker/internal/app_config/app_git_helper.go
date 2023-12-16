package appconfig

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	"github.com/pkg/errors"

	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

const (
	tmpGitProjectCloneStr          = "clone*"
	gitProjectAccessTokenAttribute = "accessToken"
	gitUrlSuffix                   = ".git"
	kubeConfig                     = "kubeconfig"
	k8sEndpoint                    = "endpoint"
	k8sClusterCA                   = "clusterCA"
)

type Config struct {
	GitDefaultCommiterName   string `envconfig:"GIT_COMMIT_NAME" default:"capten-bot"`
	GitDefaultCommiterEmail  string `envconfig:"GIT_COMMIT_EMAIL" default:"capten-bot@intelops.dev"`
	GitVaultEntityName       string `envconfig:"GIT_VAULT_ENTITY_NAME" default:"git-project"`
	GitCloneDir              string `envconfig:"GIT_CLONE_DIR" default:"/gitCloneDir"`
	GitBranchName            string `envconfig:"GIT_BRANCH_NAME" default:"capten-template-bot"`
	ManagedClusterEntityName string `envconfig:"MANAGED_CLUSER_VAULT_ENTITY_NAME" default:"managedcluster"`
}

var logger = logging.NewLogger()

type AppGitConfigHelper struct {
	cfg       Config
	gitClient *git.GitClient
}

func NewAppGitConfigHelper() (*AppGitConfigHelper, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &AppGitConfigHelper{cfg: cfg, gitClient: git.NewClient()}, nil
}

func (ca *AppGitConfigHelper) GetAccessToken(ctx context.Context, projectId string) (string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		ca.cfg.GitVaultEntityName, projectId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			ca.cfg.GitVaultEntityName, projectId)
		return "", err
	}

	return cred[gitProjectAccessTokenAttribute], nil
}

func (ca *AppGitConfigHelper) CloneRepos(ctx context.Context, templateRepo, customerRepo, token string) (templateDir string,
	reqRepo string, err error) {
	templateDir, err = os.MkdirTemp(ca.cfg.GitCloneDir, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create template tmp dir, err: %v", err)
		return
	}

	if err = ca.gitClient.Clone(templateDir, templateRepo, token); err != nil {
		os.RemoveAll(templateDir)
		err = fmt.Errorf("failed to Clone template repo, err: %v", err)
		return
	}

	reqRepo, err = os.MkdirTemp(ca.cfg.GitCloneDir, tmpGitProjectCloneStr)
	if err != nil {
		os.RemoveAll(templateDir)
		err = fmt.Errorf("failed to create tmp dir for user repo, err: %v", err)
		return
	}

	if err = ca.gitClient.Clone(reqRepo, customerRepo, token); err != nil {
		os.RemoveAll(templateDir)
		os.RemoveAll(reqRepo)
		err = fmt.Errorf("failed to Clone user repo, err: %v", err)
		return
	}

	return
}

func (ca *AppGitConfigHelper) DeployMainApp(ctx context.Context, fileName string) (string, string, error) {
	k8sclient, err := k8s.NewK8SClient(logging.NewLogger())
	if err != nil {
		return "", "", fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	// For the testing change the reqrepo to template one
	ns, resName, err := k8sclient.DynamicClient.CreateResource(ctx, fileName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the k8s custom resource: %v", err)
	}

	return ns, resName, nil

}

func (ca *AppGitConfigHelper) SyncArgoCDApp(ctx context.Context, ns, resName string) error {
	client, err := argocd.NewClient(logger)
	if err != nil {
		return err
	}

	_, err = client.TriggerAppSync(ctx, ns, resName)
	if err != nil {
		return err
	}

	return nil
}

func (ca *AppGitConfigHelper) CreateCluster(ctx context.Context, id, clusterName string) (string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		ca.cfg.ManagedClusterEntityName, id)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			ca.cfg.GitVaultEntityName, id)
		return "", err
	}

	client, err := argocd.NewClient(logger)
	if err != nil {
		return "", err
	}

	err = client.CreateOrUpdateCluster(ctx, clusterName, cred[kubeConfig])
	if err != nil {
		return "", err
	}

	return cred[k8sEndpoint], nil
}

func (ca *AppGitConfigHelper) WaitForArgoCDToSync(ctx context.Context, ns, resName string) error {
	client, err := argocd.NewClient(logger)
	if err != nil {
		return err
	}

	synched := false
	for i := 0; i < 3; i++ {
		app, err := client.GetAppSyncStatus(ctx, ns, resName)
		if err != nil {
			return fmt.Errorf("app %s synch staus fetch failed", resName)
		}

		if app.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced {
			synched = true
			break
		}

		time.Sleep(30 * time.Second)
	}

	if !synched {
		return fmt.Errorf("app %s not synched", resName)
	}
	return nil
}

func (ca *AppGitConfigHelper) AddToGit(ctx context.Context, paramType, repoUrl, token string) error {
	if err := ca.gitClient.Commit(".", "configure requested app",
		ca.cfg.GitDefaultCommiterName, ca.cfg.GitDefaultCommiterEmail); err != nil {
		return fmt.Errorf("failed to commit the changes to user repo, err: %v", err)
	}

	defaultBranch, err := ca.gitClient.GetDefaultBranchName()
	if err != nil {
		return fmt.Errorf("failed to get default branch of user repo, err: %v", err)
	}

	if err := ca.gitClient.Push(defaultBranch, token); err != nil {
		return fmt.Errorf("failed to get push to default branch, err: %v", err)
	}
	return nil
}
