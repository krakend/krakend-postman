package postman

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/devopsfaith/krakend/config"
)

func ExampleParse() {
	c := Parse(config.ServiceConfig{
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
	// output:
	// sample
	// https://schema.getpostman.com/json/collection/v2.0.0/collection.json
	// 2
	// {Name:/foo Request:{URL:/foo Method:GET Header:[] Body:{Mode: Raw:} Description:}}
	// {Name:/bar Request:{URL:/bar Method:POST Header:[] Body:{Mode: Raw:} Description:}}
}

func ExampleHandleCollection() {
	c := Parse(config.ServiceConfig{
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

	ts := httptest.NewServer(http.HandlerFunc(HandleCollection(c)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	matched, err := regexp.Match(pattern, body)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !matched {
		fmt.Println(string(body))
	}
	fmt.Println("ok")

	// output:
	// ok
}

var pattern = `{"info":{"name":"sample","_postman_id":"","description":"collection parsed at (.*)","schema":"https:\/\/schema\.getpostman\.com\/json\/collection\/v2\.0\.0\/collection\.json"},"item":\[{"name":"\/foo","request":{"url":"\/foo","method":"GET","header":null,"body":{"mode":"","raw":""}}},{"name":"\/bar","request":{"url":"\/bar","method":"POST","header":null,"body":{"mode":"","raw":""}}}\]}`
