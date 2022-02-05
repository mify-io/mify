package api_path_to_api

import (
	"errors"
	"example.com/namespace/workspace1/go-services/internal/service2/generated/api"
	"example.com/namespace/workspace1/go-services/internal/service2/generated/core"
	"net/http"
)

type PathToApiApiService struct{}

// NewPathToApiApiService creates a default api service
func NewPathToApiApiService() openapi.PathToApiApiServicer {
	return &PathToApiApiService{}
}

// PathToApiGet - sample handler
func (s *PathToApiApiService) PathToApiGet(ctx *core.MifyRequestContext) (openapi.ServiceResponse, error) {
	// TODO: add handler logic

	//TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok
	//return openapi.Response(200, map[string]interface{}{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("PathToApiGet method not implemented")
}
