package commandlinetool

import (
	"fmt"

	"github.com/vfluxus/cwlparser/libs"
)

type inputs []*input

type input struct {
	Name           string       `yaml:""`
	From           string       `yaml:""`
	Type           []string     `yaml:""`
	SecondaryFiles []string     `yaml:""`
	InputBinding   inputBinding `yaml:""`
}

type inputBinding struct {
	Position int
}

func (i *inputs) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	inputMap := make(map[string]interface{})
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

			if inputBind, ok := inputCast["position"]; ok {
				switch inputBindCast := inputBind.(type) {
				case map[interface{}]interface{}:
					if position, ok := inputBindCast["position"]; ok {
						switch positionCast := position.(type) {
						case int:
							newInput.InputBinding.Position = positionCast

						default:
							return fmt.Errorf("Can not cast input positon, data: %v, type: %T", position, position)
						}
					}

				default:
					return fmt.Errorf("Can not cast input bind, data: %v, type: %T", inputBind, inputBind)
				}
			}

		default:
			return fmt.Errorf("Can not cast input %v - type %T", inputMap[key], inputMap[key])
		}

		inputsSlice = append(inputsSlice, newInput)
	}

	*i = inputsSlice
	return nil
}
