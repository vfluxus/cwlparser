package libs

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func PrintJsonFormat(data interface{}) {
	printThis, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("Can not print in Json format, error", err)
	}
	fmt.Println(string(printThis))
}

func GetLowerLetters(raw string) (parsed string) {
	var (
		parsedBuilder = new(strings.Builder)
	)
	for i := range raw {
		if raw[i] >= 'a' && raw[i] <= 'z' {
			parsedBuilder.WriteByte(raw[i])
		}
		if raw[i] >= 'A' && raw[i] <= 'Z' {
			parsedBuilder.WriteByte(raw[i] + 32)
		}
	}
	return parsedBuilder.String()
}
