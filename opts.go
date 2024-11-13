package postman

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/luraproject/lura/v2/config"
)

func parseServiceOptions(serviceConfig *config.ServiceConfig) (*serviceOptions, error) {
	opts := &serviceOptions{
		Name:        serviceConfig.Name,
		Description: defaultDescription,
	}

	raw, ok := serviceConfig.ExtraConfig[namespace].(map[string]interface{})
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

func findFolderOptions(serviceOpts *serviceOptions, name string) *folderOptions {
	for _, f := range serviceOpts.Folder {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

func parseEndpointOptions(endpointConfig *config.EndpointConfig) (*endpointOptions, error) {
	opts := &endpointOptions{}

	endpointCfg, ok := endpointConfig.ExtraConfig[namespace].(map[string]interface{})
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

func parseVersion(serviceOpts *serviceOptions) (*Version, error) {
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

func parseVariables(cfg *config.ServiceConfig) []Variable {
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
			ID:    hash("HOST"),
			Key:   "HOST",
			Type:  "string",
			Value: net.JoinHostPort(address, strconv.Itoa(cfg.Port)),
		},
		{
			ID:    hash("SCHEMA"),
			Key:   "SCHEMA",
			Type:  "string",
			Value: schema,
		},
	}
}

type serviceOptions struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Version     string          `json:"version,omitempty"`
	Folder      []folderOptions `json:"folder,omitempty"`
}

type folderOptions struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type endpointOptions struct {
	folderOptions
	Folder string `json:"folder,omitempty"`
}
