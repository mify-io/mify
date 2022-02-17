package endpoints

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/mify-io/mify/internal/mify/util"
	"github.com/mify-io/mify/pkg/workspace"
)

type ServiceEndpoints struct {
	Api         string `yaml:"api_endpoint"`
	Maintenance string `yaml:"maintenance_endpoint"`
}

type EndpointsResolver struct {
	mu        sync.Mutex
	data      map[string]ServiceEndpoints
	workspace *workspace.Description
}

func NewEndpointsResolver(workspace *workspace.Description) *EndpointsResolver {
	return &EndpointsResolver{
		data:      make(map[string]ServiceEndpoints, 0),
		workspace: workspace,
	}
}

func (e *EndpointsResolver) ResolveEndpoints(serviceName string) (*ServiceEndpoints, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if data, ok := e.data[serviceName]; ok {
		return &data, nil
	}

	end, err := e.makeServiceEndpoints(serviceName)
	if err != nil {
		return nil, err
	}
	e.data[serviceName] = end

	return &end, nil
}

func (e *EndpointsResolver) makeServiceEndpoints(targetService string) (ServiceEndpoints, error) {
	tmpDir := e.workspace.GetServiceCacheDirectory(targetService)

	cacheFilePath := filepath.Join(tmpDir, ".service-endpoint.yaml")

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return ServiceEndpoints{}, fmt.Errorf("failed to create service cache directory: %w", err)
	}

	var cache ServiceEndpoints
	yd := util.NewYAMLData(cacheFilePath)
	err = yd.ReadFile(&cache)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ServiceEndpoints{}, fmt.Errorf("failed to read service gen cache: %w", err)
	}
	if err == nil && cache.Api != "" {
		return cache, nil
	}

	port, err := getFreePort()
	if err != nil {
		return ServiceEndpoints{}, fmt.Errorf("failed to get free port: %w", err)
	}
	cache.Api = fmt.Sprintf(":%d", port)

	port, err = getFreePort()
	if err != nil {
		return ServiceEndpoints{}, fmt.Errorf("failed to get free port: %w", err)
	}
	cache.Maintenance = fmt.Sprintf(":%d", port)

	err = yd.SaveFile(&cache)
	if err != nil {
		return ServiceEndpoints{}, fmt.Errorf("failed to save service gen cache: %w", err)
	}

	return cache, nil
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
