package entity

import (
	"context"
	"log"
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
		projectID = "test-project"
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbName = "test-database"
	}

	instance, ok := os.LookupEnv("INSTANCE")
	if !ok {
		instance = "test-instance"
	}

	emulatorHost, ok := os.LookupEnv("SPANNER_EMULATOR_HOST")
	if !ok {
		emulatorHost = "localhost:9010"
	}
	os.Setenv("SPANNER_EMULATOR_HOST", emulatorHost)

	log.Print("env:", os.Getenv("SPANNER_EMULATOR_HOST"))
	config = &Config{
		Env:       env,
		ProjectID: projectID,
		DBName:    dbName,
		Instance:  instance,
	}
}

func GetConfig() *Config {
	return config
}
