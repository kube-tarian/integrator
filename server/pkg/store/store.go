package store

import (
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/kube-tarian/kad/server/pkg/types"

	"github.com/kube-tarian/kad/server/pkg/store/astra"
)

type ServerStore interface {
	CleanupDatabase() error
	InitializeDatabase() error
	GetClusterDetails(orgID, clusterID string) (*types.ClusterDetails, error)
	GetClusterForOrg(orgID string) (*types.ClusterDetails, error)
	GetClusters(orgID string) ([]types.ClusterDetails, error)
	AddCluster(orgID, clusterID, clusterName, endpoint string) error
	UpdateCluster(orgID, clusterID, clusterName, endpoint string) error
	DeleteCluster(orgID, clusterID string) error
	WritePluginStoreConfig(clusterId string, config *pluginstorepb.PluginStoreConfig) error
	ReadPluginStoreConfig(clusterId string, storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error)
	WritePluginData(gitProjectId string, pluginData *pluginstorepb.PluginData) error
	ReadPluginData(gitProjectId string, pluginName string) (*pluginstorepb.PluginData, error)
	ReadPlugins(gitProjectId string) ([]*pluginstorepb.Plugin, error)
	DeletePlugin(gitProjectId, pluginName string) error
}

func NewStore(log logging.Logger, db string) (ServerStore, error) {
	switch db {
	case "astra":
		return astra.NewStore(log)
	}
	return nil, fmt.Errorf("db: %s not found", db)
}
