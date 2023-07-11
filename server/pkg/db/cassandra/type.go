package cassandra

const (
	keyspace                        = "capten"
	createKeyspaceQuery             = "CREATE KEYSPACE IF NOT EXISTS capten WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};"
	createClusterEndpointTableQuery = "CREATE TABLE IF NOT EXISTS capten.cluster_endpoint (cluster_id uuid, org_id uuid, cluster_name text, endpoint text, PRIMARY KEY (cluster_id, org_id));"
	createOrgClusterTableQuery      = "CREATE TABLE IF NOT EXISTS capten.org_cluster (org_id uuid, cluster_ids set<uuid>, PRIMARY KEY (org_id));"
)
