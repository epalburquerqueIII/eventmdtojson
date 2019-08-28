package main

import (
	"fmt"

	"./mdtojson"
)

func main() {

	json, err := mdtojson.ProcessRepo("http://localhost:1313/content/eventos/", "./dir")

	if json != "" {
		fmt.Printf(json)
	}
	if err != nil {
		panic(err.Error())

	}
}
