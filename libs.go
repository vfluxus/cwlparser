package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJsonFormat(data interface{}) {
	printThis, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("Can not print in Json format, error", err)
	}
	fmt.Println(string(printThis))
}
