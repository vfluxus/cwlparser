package workflowcwl

import (
	"fmt"

	"github.com/vfluxus/cwlparser/libs"
)

type output struct {
	Name         string
	Type         []string
	OutputSource []string
}

type outputs []*output

func (o *outputs) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		outputMap = make(map[string]interface{})
		outputSl  []*output
	)

	if err := unmarshal(&outputMap); err != nil {
		return err
	}

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

			if outputSource, ok := outputCast["outputSource"]; ok {
				if err := libs.AppendStringSliceWithInterface(&newOutput.OutputSource, outputSource); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("Can not cast output. Data: %v, Type: %T", outputCast, outputCast)
		}

		outputSl = append(outputSl, newOutput)
	}
	*o = outputSl

	return nil
}
