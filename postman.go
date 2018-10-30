package postman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devopsfaith/krakend/config"
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
	c := Collection{
		Info: Info{
			Name:        cfg.Name,
			Description: fmt.Sprintf("collection parsed at %s", time.Now().String()),
			Schema:      "https://schema.getpostman.com/json/collection/v2.0.0/collection.json",
		},
		Item: []Item{},
	}
	for _, e := range cfg.Endpoints {
		item := Item{
			Name: e.Endpoint,
			Request: Request{
				URL:    e.Endpoint,
				Method: e.Method,
			},
		}
		c.Item = append(c.Item, item)
	}
	return c
}

type Collection struct {
	// Variables []interface{} `json:"variables"`
	Info Info   `json:"info"`
	Item []Item `json:"item"`
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
	URL         string   `json:"url"`
	Method      string   `json:"method"`
	Header      []Header `json:"header,omitempty"`
	Body        *Body    `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
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
