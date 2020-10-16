package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type CommandLineTool struct {
	Version      string         `yaml:"cwlVersion"`
	Class        string         `yaml:"class"`
	ID           string         `yaml:"id"`
	Requirements []*requirement `yaml:"requirements"`
	Inputs       inputs         `yaml:"inputs"`
}

type requirement struct {
	Class      string `yaml:"class"`
	DockerPull string `yaml:"dockerPull"`
	RamMin     string `yaml:"ramMin"`
	CpuMin     string `yaml:"cpuMin"`
}

type input struct {
	Name           string   `yaml:""`
	From           string   `yaml:""`
	Type           []string `yaml:""`
	SecondaryFiles []string `yaml:""`
}

type inputs []*input

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
				switch inputTypeCast := inputType.(type) {
				case string:
					newInput.Type = append(newInput.Type, inputTypeCast)

				case []string:
					newInput.Type = inputTypeCast

				default:
					return fmt.Errorf("Can not cast input type: %v - type %T", inputCast["type"], inputCast["type"])
				}
			}

			if inputSecondFile, ok := inputCast["secondaryFiles"]; ok {
				switch inputSecondFileCast := inputSecondFile.(type) {
				case []string:
					newInput.SecondaryFiles = inputSecondFileCast

				case string:
					newInput.SecondaryFiles = append(newInput.SecondaryFiles, inputSecondFileCast)

				case []interface{}:
					for inputSecondFileCastIndex := range inputSecondFileCast {
						switch elementSecondaryFileCast := inputSecondFileCast[inputSecondFileCastIndex].(type) {
						case string:
							newInput.SecondaryFiles = append(newInput.SecondaryFiles, elementSecondaryFileCast)
						}
					}

				default:
					return fmt.Errorf("Can not cast input secondary files: %v - type: %T", inputCast["secondaryFiles"], inputCast["secondaryFiles"])
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

func main() {
	data, err := ioutil.ReadFile("C:\\Go\\src\\github.com\\vfluxus\\cwlparser\\test\\ApplyBQSR.cwl")
	if err != nil {
		panic(err)
	}
	var (
		cmdTool = new(CommandLineTool)
	)
	if err := yaml.Unmarshal(data, cmdTool); err != nil {
		panic(err)
	}
	PrintJsonFormat(cmdTool)
}
