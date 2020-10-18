package main

import "github.com/vfluxus/cwlparser/workflowcwl"

func ParseCWL(folder string, cwlfile string) (wfCWL *workflowcwl.WorkflowCWL, err error) {
	wfCWL = new(workflowcwl.WorkflowCWL)

	if err := wfCWL.Unmarshal(folder, cwlfile); err != nil {
		return nil, err
	}

	return wfCWL, nil
}

func main() {

}
