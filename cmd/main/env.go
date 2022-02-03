package main

import (
	"fmt"
	"os"
)

type envConfig struct {
	httpDevopsAddr      string
	httpSecretsAddr     string
	mongoConnString     string
	mongoDBName         string
	mongoCollectionName string
}

func loadEnvConfig() (*envConfig, error) {
	var err error

	cfg := &envConfig{}

	cfg.httpDevopsAddr, err = lookupEnv("HTTP_DEVOPS_ADDR")
	if err != nil {
		return nil, err
	}

	cfg.httpSecretsAddr, err = lookupEnv("HTTP_SECRETS_ADDR")
	if err != nil {
		return nil, err
	}

	cfg.mongoConnString, err = lookupEnv("MONGO_CONN_STRING")
	if err != nil {
		return nil, err
	}

	cfg.mongoDBName, err = lookupEnv("MONGO_DB_NAME")
	if err != nil {
		return nil, err
	}

	cfg.mongoCollectionName, err = lookupEnv("MONGO_COLLECTION_NAME")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func lookupEnv(name string) (string, error) {
	const provideEnvErrorMsg = `please provide "%s" environment variable`

	val, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf(provideEnvErrorMsg, name)
	}

	return val, nil
}
