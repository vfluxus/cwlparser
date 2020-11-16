package workflowrun

import (
	"strings"

	"github.com/vfluxus/cwlparser/workflowdag"
)

// convertFromStepDAGToTask 1 step -> 1 task, 1 arguments -> 1 param with regex
func convertFromStepDAGToTask(step *workflowdag.Step, taskID string, taskName string) (task *Task, err error) {
	task = &Task{
		TaskID:          taskID,
		TaskName:        taskName,
		IsBoundary:      false,
		StepID:          step.ID,
		ChildrenTasksID: step.ChildrenName, // for later replacement
		ParentTasksID:   step.ParentName,   // for later replacement
		Command:         strings.Join(step.BaseCommand, " ") + " ",
		DockerImage:     []string{step.DockerImage},
		OutputLocation:  nil,
		QueueLevel:      0,
		Status:          0,
	}
	var (
		taskParamWithRegex = make([]*ParamWithRegex, len(step.Arguments))
	)

	for argIndex := range step.Arguments {
		if step.Arguments[argIndex].Input == nil {
			newParamWithRegex := &ParamWithRegex{
				Prefix: step.Arguments[argIndex].Prefix,
			}
			taskParamWithRegex[argIndex] = newParamWithRegex
			continue
		}

		var (
			from []string
		)
		if step.Arguments[argIndex].Input.From != "" {
			from = append(from, step.Arguments[argIndex].Input.From)
		}
		newParamWithRegex := &ParamWithRegex{
			From:           from,
			SecondaryFiles: step.Arguments[argIndex].Input.SecondaryFiles,
			Regex:          step.Arguments[argIndex].Input.Value,
			Prefix:         step.Arguments[argIndex].Prefix,
		}
		taskParamWithRegex[argIndex] = newParamWithRegex
	}

	task.ParamsWithRegex = taskParamWithRegex

	return task, nil
}
