package postman

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luraproject/lura/v2/config"
)

const (
	namespace          = "documentation/postman"
	defaultDescription = "Collection parsed from KrakenD config"
	postmanJsonSchema  = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
)

// HandleCollection returns a simple http.HandleFunc exposing the POSTMAN collection description
func HandleCollection(c Collection) func(http.ResponseWriter, *http.Request) {
	b, err := json.Marshal(c)
	if err != nil {
		return func(rw http.ResponseWriter, _ *http.Request) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
	return func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(b)
	}
}

// Parse converts the received service config into a simple POSTMAN collection description
// @see https://schema.postman.com/collection/json/v2.1.0/draft-07/docs/index.html
func Parse(cfg config.ServiceConfig) Collection {
	serviceOpts, err := parseServiceOptions(&cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	c := Collection{
		Info: Info{
			Name:        serviceOpts.Name,
			PostmanID:   hash(serviceOpts.Name),
			Description: serviceOpts.Description,
			Schema:      postmanJsonSchema,
		},
		Item:      itemList{},
		Variables: parseVariables(&cfg),
	}
	if v, err := parseVersion(serviceOpts); err == nil {
		c.Info.Version = v
	}

	for _, e := range cfg.Endpoints {
		entry := newItem(e.Endpoint)
		entry.Request = &Request{
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
		opts, err := parseEndpointOptions(e)
		if err != nil {
			c.Item = append(c.Item, entry)
			continue
		}

		if opts.Name != "" {
			entry.Name = opts.Name
		}
		if opts.Description != "" {
			entry.Request.Description = opts.Description
		}

		var folder *Item
		if opts.Folder != "" && opts.Folder != separator {
			folder = createFolder(&c.Item, opts.Folder, findFolderOptions(serviceOpts, opts.Folder))
		}

		if folder != nil {
			folder.Item = append(folder.Item, entry)
		} else {
			c.Item = append(c.Item, entry)
		}
	}

	return c
}

func hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(input)))
}

type Collection struct {
	Variables []Variable `json:"variables"`
	Info      Info       `json:"info"`
	Item      itemList   `json:"item"`
}

type Info struct {
	Name        string   `json:"name"`
	PostmanID   string   `json:"_postman_id"`
	Description string   `json:"description,omitempty"`
	Schema      string   `json:"schema"`
	Version     *Version `json:"version,omitempty"`
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
