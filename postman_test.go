package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/luraproject/lura/v2/config"
)

func ExampleHandleCollection() {
	cfg := config.ServiceConfig{
		Port: 8080,
		Name: "sample",
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint: "/foo",
				Method:   "GET",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(HandleCollection(Parse(cfg))))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		fmt.Println(err.Error())
	}

	c := Collection{}
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		fmt.Println(err.Error())
	}
	res.Body.Close()

	fmt.Println(c.Info.Name)
	fmt.Println(c.Info.Schema)
	fmt.Println(len(c.Item))
	fmt.Printf("%+v\n", c.Item[0].Name)
	fmt.Printf("%+v\n", c.Item[0].Item)
	fmt.Printf("%+v\n", c.Item[0].Request.URL.Raw)
	fmt.Println(len(c.Variables))
	fmt.Printf("%+v\n", c.Variables[0].Key)
	fmt.Printf("%+v\n", c.Variables[0].Value)
	fmt.Printf("%+v\n", c.Variables[0].Type)
	fmt.Printf("%+v\n", c.Variables[1].Key)
	fmt.Printf("%+v\n", c.Variables[1].Value)
	fmt.Printf("%+v\n", c.Variables[1].Type)

	// output:
	// sample
	// https://schema.getpostman.com/json/collection/v2.1.0/collection.json
	// 1
	// /foo
	// []
	// {{SCHEMA}}://{{HOST}}/foo
	// 2
	// HOST
	// localhost:8080
	// string
	// SCHEMA
	// http
	// string
}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		in  string
		out string
	}{
		"Backwards compatibility": {
			in:  "./test/fixtures/compatibility.json",
			out: "./test/fixtures/compatibility.out.json",
		},
		"Basic collection info": {
			in:  "./test/fixtures/basic-collection-info.json",
			out: "./test/fixtures/basic-collection-info.out.json",
		},
		"Happy path for folder organization": {
			in:  "./test/fixtures/folder-happy-path.json",
			out: "./test/fixtures/folder-happy-path.out.json",
		},
		"Endpoint explicitly placed at root": {
			in:  "./test/fixtures/folder-endpoint-at-root.json",
			out: "./test/fixtures/folder-endpoint-at-root.out.json",
		},
		"Folders get created without service config": {
			in:  "./test/fixtures/folder-at-service-not-defined.json",
			out: "./test/fixtures/folder-at-service-not-defined.out.json",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := parseCfg(test.in)
			if err != nil {
				t.Error(err)
				return
			}
			c := Parse(cfg)

			b, _ := json.MarshalIndent(c, "", "\t")
			exp, _ := os.ReadFile(test.out)

			if !bytes.Equal(b, exp) {
				t.Errorf("unexpected output in %s:\n[GOT]\n%s\n\n[EXPECTED]\n%s", name, string(b), string(exp))
			}
		})
	}
}

func parseCfg(path string) (config.ServiceConfig, error) {
	cfg, _ := config.NewParser().Parse(path)
	return cfg, nil
}
