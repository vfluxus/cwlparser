package workflowdag

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/vfluxus/cwlparser/commandlinetool"
	"github.com/vfluxus/cwlparser/workflowcwl"
)

func ConvertStepCWLtoStepDAG(wfCWL *workflowcwl.WorkflowCWL, stepCWL *workflowcwl.Step, id string) (stepDAG *Step, err error) {
	// check nil fields
	if stepCWL.CommandLineTool == nil {
		return nil, errors.New(stepCWL.Name + ": CommandLineTool not exist")
	}

	if len(stepCWL.CommandLineTool.BaseCommand) == 0 {
		return nil, errors.New(stepCWL.Name + ": BaseCommand not exist")
	}

	// easy convertable fields
	stepDAG = &Step{
		ID:            id,
		Name:          stepCWL.Run,
		WorkflowName:  stepCWL.Name,
		ScatterMethod: stepCWL.ScatterMethod,
		ParentName:    stepCWL.Parents,
		ChildrenName:  stepCWL.Children,
		BaseCommand:   stepCWL.CommandLineTool.BaseCommand,
	}

	// add requirements
	if err := convertRequirements(stepCWL, stepDAG); err != nil {
		return nil, err
	}

	// step input
	stepDAG.StepInput, err = convertStepInput(stepCWL)
	if err != nil {
		return nil, err
	}

	// step output
	stepDAG.StepOutput, err = convertStepOutput(stepCWL)
	if err != nil {
		return nil, err
	}

	// add arguments
	stepDAG.Arguments, err = addArugments(stepCWL, stepDAG)
	if err != nil {
		return nil, err
	}

	return stepDAG, nil
}

// convertRequirements ...
func convertRequirements(stepCWL *workflowcwl.Step, stepDAG *Step) (err error) {
	for requirementIndex := range stepCWL.CommandLineTool.Requirements {
		// docker pull
		if len(stepCWL.CommandLineTool.Requirements[requirementIndex].DockerPull) > 0 {
			stepDAG.DockerImage = stepCWL.CommandLineTool.Requirements[requirementIndex].DockerPull
			continue
		}
		// cpu and ram
		stepDAG.Resource.Ram = stepCWL.CommandLineTool.Requirements[requirementIndex].RamMin
		stepDAG.Resource.CPU = stepCWL.CommandLineTool.Requirements[requirementIndex].CpuMin
	}
	return nil
}

// convertStepInput ...
func convertStepInput(stepCWL *workflowcwl.Step) (newStepInputs []*stepInput, err error) {
	for inputIndex := range stepCWL.CommandLineTool.Inputs {
		newStepInput := &stepInput{
			Name:           stepCWL.CommandLineTool.Inputs[inputIndex].Name,
			WorkflowName:   stepCWL.CommandLineTool.Inputs[inputIndex].WorkflowName,
			From:           stepCWL.CommandLineTool.Inputs[inputIndex].From,
			Type:           stepCWL.CommandLineTool.Inputs[inputIndex].Type,
			SecondaryFiles: stepCWL.CommandLineTool.Inputs[inputIndex].SecondaryFiles,
		}
		// add binding
		if stepCWL.CommandLineTool.Inputs[inputIndex].InputBinding != nil {
			newBinding := &stepInputBinding{
				Postition: stepCWL.CommandLineTool.Inputs[inputIndex].InputBinding.Position,
				Prefix:    stepCWL.CommandLineTool.Inputs[inputIndex].InputBinding.Prefix,
			}
			newStepInput.Binding = newBinding
		}

		// add input value from
		for stepInIndex := range stepCWL.In {
			if stepCWL.In[stepInIndex].Name == newStepInput.Name {
				valFrom, err := separateSelfValueFrom(stepCWL.In[stepInIndex].ValueFrom)
				if err != nil {
					return nil, err
				}
				newStepInput.ValueFrom = valFrom
				break
			}
		}

		// scatter
		if stepCWL.Scatter == stepCWL.CommandLineTool.Inputs[inputIndex].WorkflowName {
			newStepInput.Scatter = true
		}

		// add value from
		newStepInputs = append(newStepInputs, newStepInput)
	}
	return newStepInputs, nil
}

// separateSelfValueFrom separate <prefix> + self + <postfix>
func separateSelfValueFrom(val string) (valFrom *valueFrom, err error) {
	if val == "" {
		return nil, nil
	}

	valFrom = &valueFrom{
		Raw: val,
	}

	var (
		switchFlag       bool // false prefix, true postfix
		valueFromBuilder = new(strings.Builder)
		prefixBuilder    = new(strings.Builder)
		postfixBuilder   = new(strings.Builder)
	)

	// remove all white spaces, (), $, \, ""
	for valIndex := range val {
		if val[valIndex] == ' ' || val[valIndex] == '(' || val[valIndex] == ')' || val[valIndex] == '$' || string(val[valIndex]) == "\"" {
			continue
		}
		valueFromBuilder.WriteByte(val[valIndex])
	}
	val = valueFromBuilder.String()

	// separete by +
	separatedByPlus := strings.Split(val, "+")

	for separateIndex := range separatedByPlus {
		if strings.Contains(separatedByPlus[separateIndex], "self") {
			switchFlag = true
			continue // skip self
		}
		if !switchFlag {
			prefixBuilder.WriteString(separatedByPlus[separateIndex])
			switchFlag = false
			continue
		}
		postfixBuilder.WriteString(separatedByPlus[separateIndex])
	}

	valFrom.Prefix = prefixBuilder.String()
	valFrom.Postfix = postfixBuilder.String()
	return valFrom, nil
}

// convertStepOutput ...
func convertStepOutput(stepCWL *workflowcwl.Step) (newStepOutputs []*stepOutput, err error) {
	for outputIndex := range stepCWL.CommandLineTool.Outputs {
		newStepOutput := &stepOutput{
			Name:           stepCWL.CommandLineTool.Outputs[outputIndex].Name,
			Type:           stepCWL.CommandLineTool.Outputs[outputIndex].Type,
			SecondaryFiles: stepCWL.CommandLineTool.Outputs[outputIndex].SecondaryFiles,
			Regex:          stepCWL.CommandLineTool.Outputs[outputIndex].OutputBinding.Glob,
		}

		newStepOutputs = append(newStepOutputs, newStepOutput)
	}

	return newStepOutputs, nil
}

// addArugments ...
func addArugments(stepCWL *workflowcwl.Step, stepDAG *Step) (args []*Argument, err error) {
	if stepCWL.CommandLineTool.Arguments == nil && stepCWL.CommandLineTool.Inputs == nil {
		return nil, errors.New("No arguments or inputs found")
	}
	// make map for sort -- // MAYBE: Move to parse cwl -> struct
	var (
		argsMap       = make(map[int][]*commandlinetool.Argument)
		argsPostition = make([]int, len(stepCWL.CommandLineTool.Arguments))
	)
	for argIndex := range stepCWL.CommandLineTool.Arguments {
		if argsMap[stepCWL.CommandLineTool.Arguments[argIndex].Position] == nil {
			argsPostition = append(argsPostition, stepCWL.CommandLineTool.Arguments[argIndex].Position)
		}
		argsMap[stepCWL.CommandLineTool.Arguments[argIndex].Position] = append(argsMap[stepCWL.CommandLineTool.Arguments[argIndex].Position], stepCWL.CommandLineTool.Arguments[argIndex])
	}
	sort.Ints(argsPostition)

	for position := range argsPostition {
		for postitionIndex := range argsMap[position] {
			var (
				trimedValueFrom = strings.TrimSpace(argsMap[position][postitionIndex].ValueFrom)
			)
			for {
				if len(trimedValueFrom) == 0 || trimedValueFrom == " " {
					break
				}
				arg := &Argument{
					Postition: argsMap[position][postitionIndex].Position,
				}
				trimedValueFrom, err = separateArgValueFrom(stepDAG, arg, trimedValueFrom)
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
			}
		}
	}
	return args, nil
}

// separateArgValueFrom
func separateArgValueFrom(stepDAG *Step, arg *Argument, valueFrom string) (valueFromLeft string, err error) {
	// assert empty valueFrom
	if valueFrom == "" {
		return "", nil
	}

	var (
		prefixBuilder  = new(strings.Builder)
		valueFromIndex int
	)
	for valueFromIndex = range valueFrom {
		// write whatever before $
		if valueFrom[valueFromIndex] != '$' {
			if err := prefixBuilder.WriteByte(valueFrom[valueFromIndex]); err != nil {
				return "", fmt.Errorf("Separating argument value from. Error: %v", err)
			}
			continue
		}
		break
	}
	arg.Prefix = prefixBuilder.String()

	if inputsIndex := strings.Index(valueFrom[valueFromIndex:], "$(inputs."); inputsIndex >= 0 {
		inputName := new(strings.Builder) // $(inputs.<inputName>)
		for i := inputsIndex + len("$(inputs."); string(valueFrom[valueFromIndex:][i]) != ")" && i < len(valueFrom[valueFromIndex:]); i++ {
			inputName.WriteByte(valueFrom[valueFromIndex:][i])
		}

		for stepInputIndex := range stepDAG.StepInput {
			if inputName.String() == stepDAG.StepInput[stepInputIndex].WorkflowName || inputName.String() == stepDAG.StepInput[stepInputIndex].Name {
				// create new to avoid duplicate when append regex
				newStepInput := *stepDAG.StepInput[stepInputIndex]
				arg.Input = &newStepInput
				break
			}
		}

		if arg.Input == nil {
			return "", fmt.Errorf("Can not find %s in step input", inputName)
		}

		// Check index
		if index := valueFromIndex + len("$(inputs.") + inputName.Len() + 1 + 1; index < len(valueFrom) && string(valueFrom[index-1]) == "[" {
			indexBuilder := new(strings.Builder)
			for i := index; i < len(valueFrom) && string(valueFrom[i]) != "]"; i++ {
				indexBuilder.WriteByte(valueFrom[i])
			}
			fmt.Println(indexBuilder.String())
			argIndex, err := strconv.Atoi(indexBuilder.String())
			if err == nil {
				arg.Index = argIndex
				if index+indexBuilder.Len()+1 < len(valueFrom) {
					valueFromLeft = valueFrom[index+indexBuilder.Len()+1:]
					return valueFromLeft, nil
				}
			}
		}
		// valueFromLeft after remove "<prefix> $(inputs.<inputName>)"
		valueFromLeft = valueFrom[valueFromIndex+len("$(inputs.")+inputName.Len()+1:] // remove "--prefix $(inputs.<inputName>)" from valueFrom

		return valueFromLeft, nil
	}

	return valueFromLeft, nil
}
