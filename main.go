package main

import (
	"github.com/vfluxus/cwlparser/workflowcwl"
	"github.com/vfluxus/cwlparser/workflowdag"
)

func ParseCWL(folder string, cwlfile string) (wfCWL *workflowcwl.WorkflowCWL, err error) {
	wfCWL = new(workflowcwl.WorkflowCWL)

	if err := wfCWL.Unmarshal(folder, cwlfile); err != nil {
		return nil, err
	}

	return wfCWL, nil
}

func CreateWorkflowDAG(wfCWL *workflowcwl.WorkflowCWL) (wfDAG *workflowdag.WorkflowDAG, err error) {
	wfDAG, err = workflowdag.ConvertFromCWL(wfCWL)

	if err != nil {
		return nil, err
	}

	return wfDAG, nil
}

func main() {

}
