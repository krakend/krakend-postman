package postman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-contrib/uuid"
	"github.com/luraproject/lura/v2/config"
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
	schema := "http"
	if cfg.TLS != nil {
		schema = "https"
	}
	c := Collection{
		Info: Info{
			Name:        cfg.Name,
			PostmanID:   uuid.NewV4().String(),
			Description: fmt.Sprintf("collection parsed at %s", time.Now().String()),
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: []Item{},
		Variables: []Variable{
			{
				ID:    uuid.NewV4().String(),
				Key:   "HOST",
				Type:  "string",
				Value: fmt.Sprintf("localhost:%d", cfg.Port),
			},
			{
				ID:    uuid.NewV4().String(),
				Key:   "SCHEMA",
				Type:  "string",
				Value: schema,
			},
		},
	}
	for _, e := range cfg.Endpoints {
		item := Item{
			Name: e.Endpoint,
			Request: Request{
				URL: URL{
					Raw:      "{{SCHEMA}}://{{HOST}}" + e.Endpoint,
					Protocol: "{{SCHEMA}}",
					Host:     []string{"{{HOST}}"},
					Path:     []string{e.Endpoint[1:]},
				},
				Method: e.Method,
			},
		}
		c.Item = append(c.Item, item)
	}
	return c
}

type Collection struct {
	Variables []Variable `json:"variables"`
	Info      Info       `json:"info"`
	Item      []Item     `json:"item"`
}

type Variable struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Info struct {
	Name        string `json:"name"`
	PostmanID   string `json:"_postman_id"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

type Item struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
	// Response []interface{} `json:"response"`
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
