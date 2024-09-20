package entity

import (
	"context"
	"os"
)

type Config struct {
	Env       string
	ProjectID string

	// Spanner
	Instance string
	DBName   string
}

var config *Config

func InitConfig(ctx context.Context) {
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "local"
	}

	projectID, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		projectID = "test-project-id"
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbName = "test-db-name"
	}

	config = &Config{
		Env:       env,
		ProjectID: projectID,
		DBName:    dbName,
	}
}

func GetConfig() *Config {
	return config
}
