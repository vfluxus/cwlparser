package workflowdag

import (
	"encoding/json"
	"fmt"
	"strings"
)

type testArg struct {
	Prefix string
	Value  string
}

func ExtractValueFrom(valueFrom string) {
	// pre-process string
	valueFrom = strings.TrimSpace(valueFrom)
	valueFrom = strings.Join(strings.Fields(valueFrom), " ") // remove redundant spaces

	// separate
	sepValueFrom := strings.Split(valueFrom, " ")

	// var
	args := make([]testArg, 0, len(sepValueFrom))
	valueFlag := false

	for iSepValueFrom := range sepValueFrom {
		sepValueFrom[iSepValueFrom] += " "

		if strings.Contains(sepValueFrom[iSepValueFrom], "=") {
			valueFlag = false
			val := strings.Split(sepValueFrom[iSepValueFrom], "=")
			if len(val) != 2 {
				panic(fmt.Sprintf("Unacceptable value from %s", sepValueFrom[iSepValueFrom]))
			}
			args = append(args, testArg{
				Prefix: val[0] + "=",
				Value:  val[1],
			})
			continue
		}

		if strings.Contains(sepValueFrom[iSepValueFrom], "-") {
			valueFlag = true
			args = append(args, testArg{
				Prefix: sepValueFrom[iSepValueFrom],
			})
			continue
		}

		if valueFlag {
			valueFlag = false
			args[len(args)-1].Value = sepValueFrom[iSepValueFrom]
			continue
		}

		args = append(args, testArg{
			Value: sepValueFrom[iSepValueFrom],
		})
	}

	for i := range args {
		printArg, _ := json.MarshalIndent(args[i], "", "   ")
		fmt.Println(string(printArg))
	}
}
