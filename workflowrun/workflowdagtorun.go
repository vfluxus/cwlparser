package workflowrun

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vfluxus/cwlparser/libs"
	"github.com/vfluxus/cwlparser/workflowdag"
)

// generateRunID remove behind @ from userid, get letter & lowercases
func generateRunID(runid int, userid string) (runID string) {
	if id := strings.Index(userid, "@"); id > 0 {
		userid = userid[:id]
	}

	return strconv.Itoa(runid) + "-" + libs.GetLowerLetters(userid)
}

func generateTaskID(runID string, stepID string, stepWfName string) (taskID string) {
	return runID + "-" + stepID + "-" + libs.GetLowerLetters(stepWfName)
}

// ConvertWorkflowDAGToRun ...
func ConvertWorkflowDAGToRun(wfDAG *workflowdag.WorkflowDAG, userName string, runID int) (run *Run, err error) {
	run = &Run{
		RunID:    generateRunID(runID, userName),
		RunName:  wfDAG.Name,
		UserName: userName,
		Status:   0,
	}
	var (
		taskSl            = make([]*Task, len(wfDAG.Steps))
		stepNameTaskIDMap = make(map[string]string)
	)

	for stepIndex := range wfDAG.Steps {
		newTask, err := convertFromStepDAGToTask(wfDAG.Steps[stepIndex], generateTaskID(run.RunID, wfDAG.Steps[stepIndex].ID, wfDAG.Steps[stepIndex].WorkflowName), wfDAG.Steps[stepIndex].WorkflowName)
		if err != nil {
			return nil, err
		}
		newTask.RunID = run.RunID

		taskSl[stepIndex] = newTask
		stepNameTaskIDMap[wfDAG.Steps[stepIndex].WorkflowName] = newTask.TaskID
	}
	// replace step name in parent, child, param with regex with task id
	for taskIndex := range taskSl {
		for parentTaskIndex := range taskSl[taskIndex].ParentTasksID {
			if taskID, ok := stepNameTaskIDMap[taskSl[taskIndex].ParentTasksID[parentTaskIndex]]; ok {
				taskSl[taskIndex].ParentTasksID[parentTaskIndex] = taskID
				continue
			}
			return nil, fmt.Errorf("Can not find step name: %s. In map: %v", taskSl[taskIndex].ParentTasksID[parentTaskIndex], stepNameTaskIDMap)
		}

		for childrenTaskIndex := range taskSl[taskIndex].ChildrenTasksID {
			if taskID, ok := stepNameTaskIDMap[taskSl[taskIndex].ChildrenTasksID[childrenTaskIndex]]; ok {
				taskSl[taskIndex].ChildrenTasksID[childrenTaskIndex] = taskID
				continue
			}
			return nil, fmt.Errorf("Can not find step name: %s. In map: %v", taskSl[taskIndex].ChildrenTasksID[childrenTaskIndex], stepNameTaskIDMap)
		}

		for paramWithRegexIndex := range taskSl[taskIndex].ParamsWithRegex {
			if taskSl[taskIndex].ParamsWithRegex[paramWithRegexIndex].From != nil {
				for fromIndex := range taskSl[taskIndex].ParamsWithRegex[paramWithRegexIndex].From {
					if taskSl[taskIndex].ParamsWithRegex[paramWithRegexIndex].From[fromIndex] == "" {
						continue
					}
					if taskID, ok := stepNameTaskIDMap[taskSl[taskIndex].ParamsWithRegex[paramWithRegexIndex].From[fromIndex]]; ok {
						taskSl[taskIndex].ParamsWithRegex[paramWithRegexIndex].From[fromIndex] = taskID
						continue
					}
					return nil, fmt.Errorf("Can not find step name: %s. In map: %v", taskSl[taskIndex].ParamsWithRegex[paramWithRegexIndex].From[fromIndex], stepNameTaskIDMap)
				}
			}
		}
	}

	run.Tasks = taskSl

	addBoundary(run)

	return run, nil
}

// addBoundary add start node & end node
func addBoundary(run *Run) {
	var (
		start = &Task{
			TaskID:     run.RunID + "-" + "bigbang",
			IsBoundary: true,
		}
		end = &Task{
			TaskID:     run.RunID + "-" + "ragnarok",
			IsBoundary: true,
		}
	)

	for taskIndex := range run.Tasks {
		if len(run.Tasks[taskIndex].ParentTasksID) == 0 {
			run.Tasks[taskIndex].ParentTasksID = append(run.Tasks[taskIndex].ParentTasksID, start.TaskID)
			start.ChildrenTasksID = append(start.ChildrenTasksID, run.Tasks[taskIndex].TaskID)
		}

		if len(run.Tasks[taskIndex].ChildrenTasksID) == 0 {
			run.Tasks[taskIndex].ChildrenTasksID = append(run.Tasks[taskIndex].ChildrenTasksID, end.TaskID)
			end.ParentTasksID = append(end.ParentTasksID, run.Tasks[taskIndex].TaskID)
		}
	}

	run.Tasks = append(run.Tasks, []*Task{start, end}...)
}
