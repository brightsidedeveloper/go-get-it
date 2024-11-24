package main

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"
)

type Database struct {
	session *gocql.Session
}

// NewDatabase initializes and returns a new Database instance
func NewDatabase() (*Database, error) {
	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("PASSWORD must be set")
	}

	cluster := gocql.NewCluster(
		"node-0.aws-us-east-1.97e9c6d54b41a80c799d.clusters.scylla.cloud",
		"node-1.aws-us-east-1.97e9c6d54b41a80c799d.clusters.scylla.cloud",
		"node-2.aws-us-east-1.97e9c6d54b41a80c799d.clusters.scylla.cloud",
	)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "scylla",
		Password: password,
	}
	cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("AWS_US_EAST_1")

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	// Ensure the keyspace and table exist
	if err := initializeSchema(session); err != nil {
		return nil, err
	}

	return &Database{session: session}, nil
}

// Close closes the database session
func (db *Database) Close() {
	db.session.Close()
}

// initializeSchema ensures the required keyspace and table exist
func initializeSchema(session *gocql.Session) error {
	queries := []string{
		"CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 3}",
		"CREATE TABLE IF NOT EXISTS test.users (id UUID PRIMARY KEY, name TEXT, email TEXT)",
	}

	for _, query := range queries {
		if err := session.Query(query).Exec(); err != nil {
			return err
		}
	}
	return nil
}

// GetUsers retrieves all users from the database
func (db *Database) GetUsers() ([]map[string]interface{}, error) {
	query := db.session.Query("SELECT id, name, email FROM test.users")
	iter := query.Iter()
	return iter.SliceMap()
}
