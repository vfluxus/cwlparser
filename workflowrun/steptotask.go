package workflowrun

import (
	"github.com/vfluxus/cwlparser/workflowdag"
)

func convertFromStepDAGToTask(inputs map[string]interface{}, step *workflowdag.Step) (task *Task, err error) {
	task = &Task{
		StepID: step.ID,
	}

	return task, nil
}
