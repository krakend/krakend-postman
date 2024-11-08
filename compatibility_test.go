package postman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/luraproject/lura/v2/config"
)

func ExampleParse() {
	c := Parse(config.ServiceConfig{
		Port: 8080,
		Name: "sample",
		TLS:  &config.TLS{},
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint: "/foo",
				Method:   "GET",
			},
			{
				Endpoint: "/bar",
				Method:   "POST",
			},
		},
	})
	fmt.Println(c.Info.Name)
	fmt.Println(c.Info.Schema)
	fmt.Println(len(c.Item))
	fmt.Printf("%+v\n", c.Item[0].Name)
	fmt.Printf("%+v\n", c.Item[0].Item)
	fmt.Printf("%+v\n", c.Item[0].Request.URL.Raw)
	fmt.Printf("%+v\n", c.Item[1].Name)
	fmt.Printf("%+v\n", c.Item[1].Item)
	fmt.Printf("%+v\n", c.Item[1].Request.URL.Raw)
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
	// 2
	// /foo
	// []
	// {{SCHEMA}}://{{HOST}}/foo
	// /bar
	// []
	// {{SCHEMA}}://{{HOST}}/bar
	// 2
	// HOST
	// localhost:8080
	// string
	// SCHEMA
	// https
	// string
}

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
