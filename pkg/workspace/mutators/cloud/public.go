package cloud

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/pkg/cloudconfig"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace/mutators"
)

func UpdateCloudPublicity(mutContext *mutators.MutatorContext) error {
	frontends, err := mutContext.GetDescription().GetFrontendServices()
	if err != nil {
		return fmt.Errorf("error while listing frontend services: %w", err)
	}

	for _, frontend := range frontends {
		cfgPath := mutContext.GetDescription().GetMifySchemaAbsPath(frontend)
		frontConf, err := mifyconfig.ReadServiceCfg(cfgPath)
		if err != nil {
			return err
		}

		for clientTo := range frontConf.OpenAPI.Clients {
			cfgPath := mutContext.GetDescription().GetCloudSchemaAbsPath(clientTo)
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				continue
			}

			cloudConf, err := cloudconfig.ReadServiceCloudCfg(cfgPath)
			if err != nil {
				return err
			}

			if cloudConf.Publish {
				continue
			}

			cloudConf.Publish = true

			if err = cloudConf.WriteToFile(cfgPath); err != nil {
				return err
			}
		}
	}

	return nil
}
