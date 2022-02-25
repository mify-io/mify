package client

import (
	"fmt"

	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/cloud"
)

func AddClient(mutContext *mutators.MutatorContext, fromService string, toService string) error {
	fmt.Printf("Adding client: %s to %s\n", fromService, toService)

	if err := addClientToMifyConfig(mutContext, fromService, toService); err != nil {
		return fmt.Errorf("error while updating mify config: %w", err)
	}

	if err := cloud.UpdateCloudPublicity(mutContext); err != nil {
		return fmt.Errorf("error while updating cloud config: %w", err)
	}

	return nil
}

func addClientToMifyConfig(mutContext *mutators.MutatorContext, fromService string, toService string) error {
	serviceConf, err := mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, fromService)
	if err != nil {
		return err
	}

	_, err = mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, toService)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		serviceConf.OpenAPI.Clients = map[string]mifyconfig.ServiceOpenAPIClientConfig{}
	}
	serviceConf.OpenAPI.Clients[toService] = mifyconfig.ServiceOpenAPIClientConfig{}
	err = mifyconfig.SaveServiceConfig(mutContext.GetDescription().BasePath, fromService, serviceConf)
	if err != nil {
		return err
	}

	return nil
}

func RemoveClient(mutContext *mutators.MutatorContext, name string, clientName string) error {
	fmt.Printf("Removing client: %s to %s\n", name, clientName)

	serviceConf, err := mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, name)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		return nil
	}

	delete(serviceConf.OpenAPI.Clients, clientName)
	err = mifyconfig.SaveServiceConfig(mutContext.GetDescription().BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}
