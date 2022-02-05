// THIS FILE IS AUTOGENERATED, DO NOT EDIT
// Generated by mify via OpenAPI Generator

package openapi_init

import (
	"example.com/namespace/workspace1/go-services/internal/service2/generated/api"
	"example.com/namespace/workspace1/go-services/internal/service2/generated/core"
	"github.com/go-chi/chi/v5"

	"example.com/namespace/workspace1/go-services/internal/service2/handlers/path/to/api"
)

func Routes(ctx *core.MifyServiceContext, routerConfig openapi.RouterConfig) chi.Router {

	PathToApiApiService := api_path_to_api.NewPathToApiApiService()
	PathToApiApiController := openapi.NewPathToApiApiController(ctx, PathToApiApiService)

	router := openapi.NewRouter(ctx, routerConfig, PathToApiApiController)
	return router
}
