package postman

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/luraproject/lura/v2/config"
)

func ParseServiceOptions(serviceConfig *config.ServiceConfig) (*ServiceOptions, error) {
	opts := &ServiceOptions{
		Name:        serviceConfig.Name,
		Description: DefaultDescription,
	}

	raw, ok := serviceConfig.ExtraConfig[Namespace].(map[string]interface{})
	if !ok {
		// Backwards compatibility: there's no specific service config, we should proceed without problems
		return opts, nil
	}

	tmp, err := json.Marshal(raw)
	if err != nil {
		return nil, errors.New("invalid service config")
	}

	if err := json.Unmarshal(tmp, opts); err != nil {
		return nil, errors.New("invalid service config")
	}

	return opts, nil
}

func FindFolderOptions(serviceOpts *ServiceOptions, name string) *FolderOptions {
	for _, f := range serviceOpts.Folder {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

func ParseEndpointOptions(endpointConfig *config.EndpointConfig) (*EndpointOptions, error) {
	opts := &EndpointOptions{}

	endpointCfg, ok := endpointConfig.ExtraConfig[Namespace].(map[string]interface{})
	if !ok {
		return nil, errors.New("ignored")
	}

	tmp, err := json.Marshal(endpointCfg)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint config: %s %s", endpointConfig.Method, endpointConfig.Endpoint)
	}

	if err := json.Unmarshal(tmp, opts); err != nil {
		return nil, fmt.Errorf("invalid endpoint config: %s %s", endpointConfig.Method, endpointConfig.Endpoint)
	}
	return opts, nil
}

func ParseVersion(serviceOpts *ServiceOptions) (*Version, error) {
	if serviceOpts.Version == "" {
		return nil, errors.New("ignored")
	}

	v, err := semver.NewVersion(serviceOpts.Version)
	if err != nil {
		return nil, errors.New("the provided version is not in semver format")
	}

	version := &Version{
		Major: v.Major(),
		Minor: v.Minor(),
		Patch: v.Patch(),
	}
	return version, nil
}

func ParseVariables(cfg *config.ServiceConfig) []Variable {
	address := "localhost"
	if cfg.Address != "" {
		address = cfg.Address
	}
	schema := "http"
	if cfg.TLS != nil && !cfg.TLS.IsDisabled {
		schema = "https"
	}
	return []Variable{
		{
			ID:    Hash("HOST"),
			Key:   "HOST",
			Type:  "string",
			Value: fmt.Sprintf("%s:%d", address, cfg.Port),
		},
		{
			ID:    Hash("SCHEMA"),
			Key:   "SCHEMA",
			Type:  "string",
			Value: schema,
		},
	}
}

type ServiceOptions struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Version     string          `json:"version,omitempty"`
	Folder      []FolderOptions `json:"folder,omitempty"`
}

type FolderOptions struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type EndpointOptions struct {
	FolderOptions
	Folder string `json:"folder,omitempty"`
}
