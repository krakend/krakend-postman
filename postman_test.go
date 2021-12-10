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
	fmt.Printf("%+v\n", c.Item[0])
	fmt.Printf("%+v\n", c.Item[1])
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
	// {Name:/foo Request:{URL:{Raw:{{SCHEMA}}://{{HOST}}/foo Protocol:{{SCHEMA}} Host:[{{HOST}}] Path:[foo]} Method:GET Header:[] Body:<nil> Description:}}
	// {Name:/bar Request:{URL:{Raw:{{SCHEMA}}://{{HOST}}/bar Protocol:{{SCHEMA}} Host:[{{HOST}}] Path:[bar]} Method:POST Header:[] Body:<nil> Description:}}
	// 2
	// HOST
	// localhost:8080
	// string
	// SCHEMA
	// http
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
	fmt.Printf("%+v\n", c.Item[0])
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
	// {Name:/foo Request:{URL:{Raw:{{SCHEMA}}://{{HOST}}/foo Protocol:{{SCHEMA}} Host:[{{HOST}}] Path:[foo]} Method:GET Header:[] Body:<nil> Description:}}
	// 2
	// HOST
	// localhost:8080
	// string
	// SCHEMA
	// http
	// string
}
