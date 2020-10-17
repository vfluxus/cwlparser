package workflowcwl

import (
	"fmt"

	"github.com/vfluxus/cwlparser/commandlinetool"

	"github.com/vfluxus/cwlparser/libs"
	"gopkg.in/yaml.v2"
)

type step struct {
	Name            string
	Run             string
	In              stepIns  `yaml:"in"`
	Out             stepOuts `yaml:"out"`
	CommandLineTool *commandlinetool.CommandLineTool
}

// ---------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- STEP IN ----------------------------------------------------------

type stepIn struct {
	Name      string
	Source    string
	ValueFrom string
}

type stepIns []*stepIn

func (si *stepIns) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		stepInMap = make(map[string]interface{})
		stepInSl  []*stepIn
	)
	if err := unmarshal(&stepInMap); err != nil {
		return err
	}

	for key := range stepInMap {
		newStepIn := new(stepIn)
		newStepIn.Name = key

		switch stepInCast := stepInMap[key].(type) {
		case map[interface{}]interface{}:
			if stepInSource, ok := stepInCast["source"]; ok {
				source, err := libs.CastToString(stepInSource)
				if err != nil {
					return err
				}
				newStepIn.Source = source
			}

			if stepInValueFrom, ok := stepInCast["valueFrom"]; ok {
				valueFrom, err := libs.CastToString(stepInValueFrom)
				if err != nil {
					return err
				}
				newStepIn.ValueFrom = valueFrom
			}

		case string:
			newStepIn.ValueFrom = stepInCast

		default:
			return fmt.Errorf("Can not cast stepIn. Data: %v. Type: %T", stepInCast, stepInCast)
		}

		stepInSl = append(stepInSl, newStepIn)
	}
	*si = stepInSl

	return nil
}

// ----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- STEP OUT ----------------------------------------------------------
type stepOut struct {
	Name string
}

type stepOuts []*stepOut

func (so *stepOuts) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		stepOutSl []*stepOut
		err1      error // error map[string]interface{}
		err2      error
	)
	// unmarshal to map[string]interface{}
	var (
		stepOutMap = make(map[string]interface{})
	)
	if err1 = unmarshal(&stepOutMap); err1 == nil {
		for key := range stepOutMap {
			newStepOut := new(stepOut)
			newStepOut.Name = key

			stepOutSl = append(stepOutSl, newStepOut)
		}

		*so = stepOutSl
		return nil
	}

	// unmarshal to []string
	var (
		strSl []string
	)
	if err2 = unmarshal(&strSl); err2 == nil {
		for strSlIndex := range strSl {
			newStepOut := new(stepOut)
			newStepOut.Name = strSl[strSlIndex]

			stepOutSl = append(stepOutSl, newStepOut)
		}

		*so = stepOutSl
		return nil
	}

	return fmt.Errorf("Can not unmarshal step out. map[string]interface{}: %v. []string: %v", err1, err2)
}

// -------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- STEPS ----------------------------------------------------------
type steps []*step

func (s *steps) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		stepMap = make(map[string]interface{})
		stepSl  []*step
	)
	if err := unmarshal(&stepMap); err != nil {
		return err
	}

	for key := range stepMap {
		newStep := new(step)
		newStep.Name = key

		stepByte, err := yaml.Marshal(stepMap[key])
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(stepByte, newStep); err != nil {
			return err
		}

		stepSl = append(stepSl, newStep)
	}

	*s = stepSl

	return nil
}
