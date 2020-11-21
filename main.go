package cwlparser

import (
	"bytes"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"

	"github.com/vfluxus/cwlparser/workflowcwl"
	"github.com/vfluxus/cwlparser/workflowdag"
	"github.com/vfluxus/cwlparser/workflowrun"
)

func ParseCWL(folder string, cwlfile string) (wfCWL *workflowcwl.WorkflowCWL, err error) {
	wfCWL = new(workflowcwl.WorkflowCWL)

	if err := wfCWL.Unmarshal(folder, cwlfile); err != nil {
		return nil, err
	}

	return wfCWL, nil
}

func ParseCWLInMem(f *workflowcwl.HttpCWLForm) (wfCWL *workflowcwl.WorkflowCWL, err error) {
	wfCWL = new(workflowcwl.WorkflowCWL)

	if err := wfCWL.UnmarshalJson(f); err != nil {
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

func CreateRunFromWorkflow(wfDAG *workflowdag.WorkflowDAG, inputs map[string]interface{}, userID string, runID int) (run *workflowrun.Run, err error) {
	if err := workflowdag.AddValueToStepInAndArg(inputs, wfDAG); err != nil {
		return nil, err
	}

	if err := workflowdag.AddOutputToInput(wfDAG); err != nil {
		return nil, err
	}

	run, err = workflowrun.ConvertWorkflowDAGToRun(wfDAG, userID, runID)
	if err != nil {
		return nil, err
	}

	return run, nil
}

func CreateGraphvizDot(run *workflowrun.Run) (dot string, err error) {
	if run == nil {
		return "", nil
	}

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return "", err
	}
	defer func() {
		if err = graph.Close(); err != nil {
			return
		}
		g.Close()
	}()

	var (
		nodeMap = make(map[string]*cgraph.Node)
	)
	for taskIndex := range run.Tasks {
		node, err := graph.CreateNode(run.Tasks[taskIndex].TaskID)
		if err != nil {
			return "", err
		}
		nodeMap[run.Tasks[taskIndex].TaskID] = node
	}

	for taskIndex := range run.Tasks {
		for childIndex := range run.Tasks[taskIndex].ChildrenTasksID {
			_, err := graph.CreateEdge("", nodeMap[run.Tasks[taskIndex].TaskID], nodeMap[run.Tasks[taskIndex].ChildrenTasksID[childIndex]])
			if err != nil {
				return "", err
			}
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.XDOT, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
