package cassandra

import (
	"os"
	"testing"

	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	"github.com/kube-tarian/kad/capten/common-pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	log := logging.NewLogger()
	setEnvConfig()

	err := cassandra.Create(log)
	assert.Nil(t, err)

	migrationClient, err := NewCassandraMigrate(log)
	assert.Nil(t, err)

	err = migrationClient.Run()
	assert.Nil(t, err)
}

func setEnvConfig() {
	os.Setenv("DB_ADDRESSES", "127.0.0.1:9042")
	os.Setenv("DB_ADMIN_USERNAME", "cassandra")
	os.Setenv("DB_SERVICE_USERNAME", "agent")
	os.Setenv("CASSANDRA_DB_NAME", "integrator")
	os.Setenv("DB_REPLICATION_FACTOR", `'datacenter1': 1`)
	os.Setenv("DB_ADMIN_PASSWD", "cassandra")
	os.Setenv("DB_SERVICE_PASSWD", "agent")
	os.Setenv("SOURCE_URI", "file://tests/migrations/cassandra")
}
