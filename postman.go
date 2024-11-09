package postman

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luraproject/lura/v2/config"
)

const (
	Namespace          = "documentation/postman"
	DefaultDescription = "Collection parsed from KrakenD config"
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
// @see https://schema.postman.com/collection/json/v2.1.0/draft-07/docs/index.html
func Parse(cfg config.ServiceConfig) Collection {
	serviceOpts, err := ParseServiceOptions(&cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	c := Collection{
		Info: Info{
			Name:        serviceOpts.Name,
			PostmanID:   Hash(serviceOpts.Name),
			Description: serviceOpts.Description,
			Schema:      PostmanJsonSchema,
		},
		Item:      ItemList{},
		Variables: ParseVariables(&cfg),
	}
	if v, err := ParseVersion(serviceOpts); err == nil {
		c.Info.Version = v
	}

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

		var folder *Item
		if opts.Folder != "" && opts.Folder != "/" {
			folder = CreateFolder(&c.Item, opts.Folder, FindFolderOptions(serviceOpts, opts.Folder))
		}

		if folder != nil {
			folder.Item = append(folder.Item, item)
		} else {
			c.Item = append(c.Item, item)
		}
	}

	return c
}

func Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(input)))
}

type Collection struct {
	Variables []Variable `json:"variables"`
	Info      Info       `json:"info"`
	Item      ItemList   `json:"item"`
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
