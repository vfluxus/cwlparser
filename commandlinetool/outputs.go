package commandlinetool

import (
	"fmt"

	"github.com/vfluxus/cwlparser/libs"
)

type outputs []*output

type output struct {
	Name           string        `yaml:""`
	Type           []string      `yaml:"type"`
	OutputBinding  outputBinding `yaml:"outputBinding"`
	SecondaryFiles []string      `yaml:"secondaryFiles"`
}

type outputBinding struct {
	Glob []string `yaml:"glob"`
}

func (o *outputs) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		outputMap = make(map[string]interface{})
	)
	if err := unmarshal(&outputMap); err != nil {
		return err
	}

	var (
		outputsSlice []*output
	)

	for key := range outputMap {
		newOutput := new(output)
		newOutput.Name = key

		switch outputCast := outputMap[key].(type) {
		case map[interface{}]interface{}:
			if outputType, ok := outputCast["type"]; ok {
				if err := libs.AppendStringSliceWithInterface(&newOutput.Type, outputType); err != nil {
					return err
				}
			}

			if outputSecondFile, ok := outputCast["secondaryFiles"]; ok {
				if err := libs.AppendStringSliceWithInterface(&newOutput.SecondaryFiles, outputSecondFile); err != nil {
					return err
				}
			}

			if outputBind, ok := outputCast["outputBinding"]; ok {
				switch outputBindCast := outputBind.(type) {
				case map[interface{}]interface{}:
					if glob, ok := outputBindCast["glob"]; ok {
						if err := libs.AppendStringSliceWithInterface(&newOutput.OutputBinding.Glob, glob); err != nil {
							return err
						}
					}
				}
			}
		default:
			return fmt.Errorf("Can not cast output. Data: %v, Type: %T", outputCast, outputCast)
		}

		outputsSlice = append(outputsSlice, newOutput)
	}
	*o = outputsSlice
	return nil
}
