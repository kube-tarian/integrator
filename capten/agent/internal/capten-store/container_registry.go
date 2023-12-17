package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/pkg/errors"
)

const (
	insertContainerRegistry        = "INSERT INTO %s.ContainerRegistry(id, registry_url, labels, last_update_time, registry_type) VALUES (?,?,?,?,?)"
	updateContainerRegistryById    = "UPDATE %s.ContainerRegistry SET %s WHERE id=?"
	deleteContainerRegistryById    = "DELETE FROM %s.ContainerRegistry WHERE id= ?"
	selectAllContainerRegistrys    = "SELECT id, registry_url, labels, last_update_time, registry_type FROM %s.ContainerRegistry"
	selectGetContainerRegistryById = "SELECT id, registry_url, labels, last_update_time, registry_type FROM %s.ContainerRegistry WHERE id=%s;"
)

func (a *Store) UpsertContainerRegistry(config *captenpluginspb.ContainerRegistry) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertContainerRegistry, a.keyspace), config.Id,
		config.RegistryUrl, config.Labels, config.LastUpdateTime, config.RegistryType)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		updatePlaceholders, values := formUpdateKvPairsForContainerRegistry(config)
		if updatePlaceholders == "" {
			return err
		}
		query := fmt.Sprintf(updateContainerRegistryById, a.keyspace, updatePlaceholders)
		args := append(values, config.Id)
		batch = a.client.Session().NewBatch(gocql.LoggedBatch)
		batch.Query(query, args...)
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) DeleteContainerRegistryById(id string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deleteContainerRegistryById,
		a.keyspace), id)
	err := deleteAction.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) GetContainerRegistryForID(id string) (*captenpluginspb.ContainerRegistry, error) {
	query := fmt.Sprintf(selectGetContainerRegistryById, a.keyspace, id)
	projects, err := a.executeContainerRegistrysSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf("project not found")
	}
	return projects[0], nil
}

func (a *Store) GetContainerRegistrys() ([]*captenpluginspb.ContainerRegistry, error) {
	query := fmt.Sprintf(selectAllContainerRegistrys, a.keyspace)
	return a.executeContainerRegistrysSelectQuery(query)
}

func (a *Store) executeContainerRegistrysSelectQuery(query string) ([]*captenpluginspb.ContainerRegistry, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	project := captenpluginspb.ContainerRegistry{}

	ret := make([]*captenpluginspb.ContainerRegistry, 0)
	for iter.Scan(
		&project.Id, &project.RegistryUrl,
		&project.Labels, &project.LastUpdateTime, &project.RegistryType,
	) {
		ContainerRegistry := &captenpluginspb.ContainerRegistry{
			Id:             project.Id,
			RegistryUrl:    project.RegistryUrl,
			Labels:         project.Labels,
			LastUpdateTime: project.LastUpdateTime,
			RegistryType:   project.RegistryType,
		}
		ret = append(ret, ContainerRegistry)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForContainerRegistry(config *captenpluginspb.ContainerRegistry) (updatePlaceholders string, values []interface{}) {
	params := []string{}

	if config.RegistryUrl != "" {
		params = append(params, "registry_url = ?")
		values = append(values, config.RegistryUrl)
	}

	if len(config.Labels) != 0 {
		params = append(params, "labels = ?")
		values = append(values, config.Labels)
	}

	if config.LastUpdateTime != "" {
		params = append(params, "last_update_time = ?")
		values = append(values, config.LastUpdateTime)
	}

	if config.RegistryType != "" {
		params = append(params, "registry_type = ?")
		values = append(values, config.RegistryType)
	}

	if len(params) == 0 {
		return "", nil
	}
	return strings.Join(params, ", "), values
}
