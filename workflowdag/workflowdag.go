package workflowdag

import (
	"errors"

	"github.com/vfluxus/cwlparser/workflowcwl"
)

type WorkflowDAG struct {
	Name  string  `json:"name"`
	Steps []*Step `json:"steps"`
}

func ConvertFromCWL(wfCWL *workflowcwl.WorkflowCWL) (wfDAG *WorkflowDAG, err error) {
	wfDAG = &WorkflowDAG{
		Name: wfCWL.ID,
	}

	var (
		nameIDMap = make(map[string]int) // use for add id to parent and children name
	)
	id := 0
	for stepCWLIndex := range wfCWL.Steps {
		newStepDAG, err := ConvertStepCWLtoStepDAG(wfCWL, wfCWL.Steps[stepCWLIndex], id)
		if err != nil {
			return nil, err
		}

		if _, ok := nameIDMap[newStepDAG.WorkflowName]; ok {
			return nil, errors.New("Duplicate step workflowname")
		}
		nameIDMap[newStepDAG.WorkflowName] = id

		id++
		wfDAG.Steps = append(wfDAG.Steps, newStepDAG)
	}

	// add id to parent and children name
	for stepDAGIndex := range wfDAG.Steps {
		for parentIndex := range wfDAG.Steps[stepDAGIndex].ParentName {
			if parentID, ok := nameIDMap[wfDAG.Steps[stepDAGIndex].ParentName[parentIndex]]; ok {
				wfDAG.Steps[stepDAGIndex].ParentID = append(wfDAG.Steps[stepDAGIndex].ParentID, parentID)
				continue
			}
			return nil, errors.New("Parent not found")
		}

		for childrenIndex := range wfDAG.Steps[stepDAGIndex].ChildrenName {
			if childrenID, ok := nameIDMap[wfDAG.Steps[stepDAGIndex].ChildrenName[childrenIndex]]; ok {
				wfDAG.Steps[stepDAGIndex].ChildrenID = append(wfDAG.Steps[stepDAGIndex].ChildrenID, childrenID)
				continue
			}
			return nil, errors.New("Children not found")
		}
	}

	return wfDAG, nil
}
