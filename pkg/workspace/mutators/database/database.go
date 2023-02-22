package database

import (
	"fmt"

	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace/mutators"
)

func AddPostgres(mutContext *mutators.MutatorContext, service string) error {
	fmt.Printf("Adding postgres to %s\n", service)

	serviceConf, err := mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, service)
	if err != nil {
		return err
	}

	serviceConf.Postgres.Enabled = true
	return mifyconfig.SaveServiceConfig(mutContext.GetDescription().BasePath, service, serviceConf)
}

func RemovePostgres(mutContext *mutators.MutatorContext, service string) error {
	fmt.Printf("Removing postgres to %s\n", service)

	serviceConf, err := mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, service)
	if err != nil {
		return err
	}

	serviceConf.Postgres = mifyconfig.PostgresConfig{}
	return mifyconfig.SaveServiceConfig(mutContext.GetDescription().BasePath, service, serviceConf)
}
