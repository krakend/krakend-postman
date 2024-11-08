package postman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"

	"github.com/luraproject/lura/v2/config"
)

const (
	Namespace          = "documentation/postman"
	DefaultDescription = "Collection parsed at %s"
	PostmanJsonSchema  = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
)

// HandleCollection returns a simple http.HandleFunc exposing the POSTMAN collection description
func HandleCollection(c Collection) func(http.ResponseWriter, *http.Request) {
	b, err := json.Marshal(c)
	if err != nil {
		return func(rw http.ResponseWriter, r *http.Request) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(b)
	}
}

// Parse converts the received service config into a simple POSTMAN collection description
func Parse(cfg config.ServiceConfig) Collection {
	serviceOpts, err := ParseServiceOptions(&cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	c := Collection{
		Info: Info{
			Name:        serviceOpts.Name,
			PostmanID:   "fixed",
			Description: serviceOpts.Description,
			Schema:      PostmanJsonSchema,
		},
		Item:      Branch{},
		Variables: ParseVariables(&cfg),
	}
	if v, err := ParseVersion(serviceOpts); err == nil {
		c.Info.Version = *v
	}

	// Iterate the endpoint list and generate a map of paths
	var flattenPathsKeys []string
	flattenPaths := map[string][]string{}
	for _, e := range cfg.Endpoints {
		opts, err := ParseEndpointOptions(e)
		if err != nil {
			continue
		}

		flattenPaths[opts.Folder] = SlicePath(opts.Folder)
		flattenPathsKeys = append(flattenPathsKeys, opts.Folder)
	}

	sort.Strings(flattenPathsKeys)
	flattenPathsKeys = slices.Compact(flattenPathsKeys)

	// Now we build a first version of the tree based on the paths provided in the endpoint list
	// The folders are enriched with the provided service configuration
	for _, rawPath := range flattenPathsKeys {
		slicedPath := flattenPaths[rawPath]
		if len(slicedPath) == 0 {
			continue
		}

		branch := c.Item.FindItem(slicedPath[0])
		if branch == nil {
			branch = NewItem(slicedPath[0])
			folderOpts := FindFolderOptions(serviceOpts, rawPath)
			if folderOpts != nil {
				branch.Description = folderOpts.Description
			}
			c.Item = append(c.Item, branch)
		}

		for _, value := range slicedPath[1:] {
			child := branch.FindChild(value)
			if child == nil {
				child = NewItem(value)
				folderOpts := FindFolderOptions(serviceOpts, rawPath)
				if folderOpts != nil {
					child.Description = folderOpts.Description
				}
				branch.Item = append(branch.Item, child)
			}
			branch = child
		}
	}

	// The folder tree is now in place, let's traverse it to add the endpoints
	for _, e := range cfg.Endpoints {
		item := NewItem(e.Endpoint)
		item.Request = &Request{
			URL: URL{
				Raw:      "{{SCHEMA}}://{{HOST}}" + e.Endpoint,
				Protocol: "{{SCHEMA}}",
				Host:     []string{"{{HOST}}"},
				Path:     []string{e.Endpoint[1:]},
			},
			Method: e.Method,
		}

		// The endpoints that do not have options are added to the root of the collection
		// This simple check handles the backwards compatibility of the generator
		opts, err := ParseEndpointOptions(e)
		if err != nil {
			c.Item = append(c.Item, item)
			continue
		}

		if opts.Name != "" {
			item.Name = opts.Name
		}
		if opts.Description != "" {
			item.Request.Description = opts.Description
		}
		if opts.Folder != "" {
			node := c.Item.FindByPath(opts.Folder)
			node.Item = append(node.Item, item)
		} else {
			c.Item = append(c.Item, item)
		}
	}

	return c
}

type Collection struct {
	Variables []Variable `json:"variables"`
	Info      Info       `json:"info"`
	Item      Branch     `json:"item"`
}

type Info struct {
	Name        string  `json:"name"`
	PostmanID   string  `json:"_postman_id"`
	Description string  `json:"description,omitempty"`
	Schema      string  `json:"schema"`
	Version     Version `json:"version,omitempty"`
}

type Request struct {
	URL         URL      `json:"url"`
	Method      string   `json:"method"`
	Header      []Header `json:"header,omitempty"`
	Body        *Body    `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
}

type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
}

type Body struct {
	Mode string `json:"mode,omitempty"`
	Raw  string `json:"raw,omitempty"`
}

type Header []struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type Version struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Patch uint64 `json:"patch"`
}

type Variable struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
