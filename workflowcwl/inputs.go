package workflowcwl

import (
	"fmt"

	"github.com/vfluxus/cwlparser/libs"
)

type input struct {
	Name           string
	Type           []string
	SecondaryFiles []string
}

type inputs []*input

func (i *inputs) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		inputMap = make(map[string]interface{})
	)
	if err := unmarshal(&inputMap); err != nil {
		return err
	}

	var (
		inputsSlice []*input
	)
	for key := range inputMap {
		newInput := new(input)
		newInput.Name = key

		switch inputCast := inputMap[key].(type) {
		case map[interface{}]interface{}:
			if inputType, ok := inputCast["type"]; ok {
				if err := libs.AppendStringSliceWithInterface(&newInput.Type, inputType); err != nil {
					return err
				}
			}

			if inputSecondFile, ok := inputCast["secondaryFiles"]; ok {
				if err := libs.AppendStringSliceWithInterface(&newInput.SecondaryFiles, inputSecondFile); err != nil {
					return err
				}
			}

		case string:
			newInput.Type = append(newInput.Type, inputCast)

		case []string:
			newInput.Type = append(newInput.Type, inputCast...)

		default:
			return fmt.Errorf("Can not cast input %v - type %T", inputMap[key], inputMap[key])
		}
		inputsSlice = append(inputsSlice, newInput)
	}
	*i = inputsSlice

	return nil
}
