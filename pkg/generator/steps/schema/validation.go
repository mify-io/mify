package schema

import (
	"fmt"
	"os"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/workspace"
)

func validateCtx(ctx *gencontext.GenContext) error {
	if ctx.GetServiceName() == workspace.DevRunnerName {
		return nil // Dev runner shouldn't have any scheme
	}
	if ctx.MustGetMifySchema().IsExternal {
		return nil // Dev runner shouldn't have any scheme
	}

	schemas := ctx.GetSchemaCtx().GetAllSchemas()
	serviceSchemas, ok := schemas[ctx.GetServiceName()]
	if !ok {
		schemasDirPath := ctx.GetWorkspace().GetSchemasAbsPath(ctx.GetServiceName())

		_, err := os.Stat(schemasDirPath)
		if os.IsNotExist(err) {
			return fmt.Errorf("schema directory for service '%s' wasn't found. "+
				"Did you created a valid service folder inside workspace schemas dir? "+
				"Expected service folder location: '%s'",
				ctx.GetServiceName(),
				schemasDirPath)
		} else if err != nil {
			return fmt.Errorf("error while reading service '%s' schema folder: %w",
				ctx.GetServiceName(),
				err)
		}
	}

	if serviceSchemas.GetMify() == nil {
		return fmt.Errorf("mify schema is missing for service '%s'. Expected location: %s",
			ctx.GetServiceName(),
			ctx.GetWorkspace().GetMifySchemaAbsPath(ctx.GetServiceName()))
	}

	return nil
}
