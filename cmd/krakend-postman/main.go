package main

import (
	"encoding/json"
	"flag"
	"fmt"

	postman "github.com/krakend/krakend-postman/v3"
	"github.com/luraproject/lura/v2/config"
)

func main() {
	conf := flag.String("c", "krakend.json", "the config file")
	flag.Parse()
	cfg, err := config.NewParser().Parse(*conf)
	if err != nil {
		fmt.Println("error parsing the config file:", err.Error())
		return
	}
	c, parseError := postman.Parse(cfg)
	if parseError != nil {
		fmt.Println("error parsing the collection: ", parseError.Error())
	}
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		fmt.Println("error marshaling the postma descriptor:", err.Error())
		return
	}
	fmt.Println(string(b))
}
