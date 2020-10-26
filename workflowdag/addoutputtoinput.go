package workflowdag

import (
	"fmt"
)

// AddOutputToInput add output regex & secondary file to step input
func AddOutputToInput(wf *WorkflowDAG) (err error) {
	var (
		stepMap = make(map[string]*Step)
	)
	for stepIndex := range wf.Steps {
		if _, ok := stepMap[wf.Steps[stepIndex].WorkflowName]; ok {
			return fmt.Errorf("Duplicate step name: %s", wf.Steps[stepIndex].ID)
		}
		stepMap[wf.Steps[stepIndex].WorkflowName] = wf.Steps[stepIndex]
	}

	for stepIndex := range wf.Steps {
		for ArgumentIndex := range wf.Steps[stepIndex].Arguments {
			if wf.Steps[stepIndex].Arguments[ArgumentIndex].Input == nil {
				continue
			}
			if wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.From == "" {
				continue
			}
			if step, ok := stepMap[wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.From]; ok {
				for stepOutputIndex := range step.StepOutput {
					if step.StepOutput[stepOutputIndex].Name == wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.WorkflowName {
						// TODO: maybe type check, value from ??
						wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.Value = append(wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.Value, step.StepOutput[stepOutputIndex].Regex...)
						wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.SecondaryFiles = append(wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.SecondaryFiles, step.StepOutput[stepOutputIndex].SecondaryFiles...)
					}
				}
				continue
			}

			return fmt.Errorf("Can not find step: %s. In step map: %v", wf.Steps[stepIndex].Arguments[ArgumentIndex].Input.From, stepMap)
		}
	}
	return nil
}
