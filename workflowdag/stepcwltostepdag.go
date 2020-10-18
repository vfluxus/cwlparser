package workflowdag

import (
	"errors"

	"github.com/vfluxus/cwlparser/workflowcwl"
)

func ConvertStepCWLtoStepDAG(stepCWL *workflowcwl.Step, id int) (stepDAG *Step, err error) {
	// check nil fields
	if stepCWL.CommandLineTool == nil {
		return nil, errors.New("CommandLineTool not exist")
	}

	if len(stepCWL.CommandLineTool.BaseCommand) == 0 {
		return nil, errors.New("BaseCommand not exist")
	}

	// easy convertable fields
	stepDAG = &Step{
		ID:           id,
		Name:         stepCWL.Name,
		ParentName:   stepCWL.Parents,
		ChildrenName: stepCWL.Children,
		BaseCommand:  stepCWL.CommandLineTool.BaseCommand,
	}

	// add requirements
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

	// step input
	for inputIndex := range stepCWL.CommandLineTool.Inputs {
		newStepInput := &stepInput{
			Name:           stepCWL.CommandLineTool.Inputs[inputIndex].Name,
			Type:           stepCWL.CommandLineTool.Inputs[inputIndex].Type,
			SecondaryFiles: stepCWL.CommandLineTool.Inputs[inputIndex].SecondaryFiles,
		}
		// add binding
		if stepCWL.CommandLineTool.Inputs[inputIndex].InputBinding != nil {
			newStepInput.Binding.Postition = stepCWL.CommandLineTool.Inputs[inputIndex].InputBinding.Position
			newStepInput.Binding.Prefix = stepCWL.CommandLineTool.Inputs[inputIndex].InputBinding.Prefix
		}

		stepDAG.StepInput = append(stepDAG.StepInput, newStepInput)
	}

	return stepDAG, nil
}
