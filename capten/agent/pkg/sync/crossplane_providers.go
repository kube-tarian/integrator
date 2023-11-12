package sync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	providerNamePrefix = "provider-"
)

type FetchCrossPlaneProviders struct {
	log    logging.Logger
	client *k8s.K8SClient
	db     *captenstore.Store
	creds  credentials.CredentialAdmin
}

func NewFetchCrossPlaneProviders() (*FetchCrossPlaneProviders, error) {
	log := logging.NewLogger()
	db, err := captenstore.NewStore(log)
	if err != nil {
		// ignoring store failure until DB user creation working
		// return nil, err
		log.Errorf("failed to initialize store, %v", err)
	}

	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	credAdmin, err := credentials.NewCredentialAdmin(context.TODO())
	if err != nil {
		log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client")
		return nil, err
	}

	return &FetchCrossPlaneProviders{log: log, client: k8sclient, db: db, creds: credAdmin}, nil
}

func (fetch *FetchCrossPlaneProviders) Run() {
	fetch.log.Info("started to sync CrossplaneProvider resources")

	objList, err := fetch.client.DynamicClient.ListAllNamespaceResource(context.TODO(), schema.GroupVersionResource{Group: "pkg.crossplane.io", Version: "v1", Resource: "providers"})
	if err != nil {
		fetch.log.Error("Failed to fetch all the resource, err:", err)

		return
	}

	providers, err := json.Marshal(objList)
	if err != nil {
		fetch.log.Error("Failed to marshall the data, err:", err)

		return
	}

	var providerObj model.ProviderList
	err = json.Unmarshal(providers, &providerObj)
	if err != nil {
		fetch.log.Error("Failed to un-marshall the data, err:", err)

		return
	}

	fetch.UpdateCrossplaneProvider(providerObj.Items)

	fetch.log.Info("succesfully sync-ed CrossplaneProvider resources")
}

func (fetch *FetchCrossPlaneProviders) UpdateCrossplaneProvider(clObj []model.Provider) {
	prvList, err := fetch.db.GetCrossplaneProviders()
	if err != nil {
		fetch.log.Error("Failed to GetCrossplaneProviders, err:", err)

		return
	}

	prvMap := make(map[string]*captenpluginspb.CrossplaneProvider)
	for _, prov := range prvList {
		prvMap[providerNamePrefix+prov.ProviderName] = prov
	}

	for _, obj := range clObj {
		for _, status := range obj.Status.Conditions {
			if status.Type != model.TypeHealthy {
				continue
			}

			prvObj, ok := prvMap[obj.Name]
			if !ok {
				fetch.log.Infof("Provider name %s is not found in the db, skipping the update", obj.Name)
				continue
			}

			provider := model.CrossplaneProvider{
				Id:              prvObj.Id,
				Status:          string(status.Status),
				CloudType:       prvObj.CloudType,
				CloudProviderId: prvObj.CloudProviderId,
				ProviderName:    prvObj.ProviderName,
			}

			if err := fetch.db.UpdateCrossplaneProvider(&provider); err != nil {
				fetch.log.Errorf("failed to update provider %s details in db, err: ", prvObj.ProviderName, err)
				continue
			}

			fetch.log.Infof("successfully updated the details for %s", prvObj.ProviderName)

		}
	}
}
