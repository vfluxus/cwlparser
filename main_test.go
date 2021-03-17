package cwlparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/vfluxus/cwlparser/commandlinetool"
	"github.com/vfluxus/cwlparser/libs"
	"github.com/vfluxus/cwlparser/workflowcwl"
	"github.com/vfluxus/cwlparser/workflowdag"
	"github.com/vfluxus/cwlparser/workflowrun"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"gopkg.in/yaml.v2"
)

const (
	argumentCwl  = "/home/tpp/go/src/github.com/vfluxus/cwlparser/test/arguments.cwl"
	applyBSQRCwl = "/home/tpp/go/src/github.com/vfluxus/cwlparser/test/ApplyBQSR.cwl"
)

func TestCmdLineTool(t *testing.T) {
	data, err := ioutil.ReadFile(applyBSQRCwl)
	if err != nil {
		t.Fatal(err)
	}

	var (
		cmdTool = new(commandlinetool.CommandLineTool)
	)
	if err := yaml.Unmarshal(data, cmdTool); err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(cmdTool)

	data2, err2 := ioutil.ReadFile(argumentCwl)
	if err2 != nil {
		t.Fatal(err2)
	}

	var (
		cmdTool2 = new(commandlinetool.CommandLineTool)
	)
	if err2 := yaml.Unmarshal(data2, cmdTool2); err2 != nil {
		t.Fatal(err2)
	}
	libs.PrintJsonFormat(cmdTool2)
}

func TestWorkflowCWL(t *testing.T) {
	data, err := ioutil.ReadFile("/home/tpp/go/src/github.com/vfluxus/demo-cwl/wes/wes.cwl")
	if err != nil {
		t.Fatal(err)
	}
	var (
		wfCwl = new(workflowcwl.WorkflowCWL)
	)
	if err := yaml.Unmarshal(data, wfCwl); err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(wfCwl)
}

func TestWorkflowCWLUnmarshal(t *testing.T) {
	newWorkflowCWL := new(workflowcwl.WorkflowCWL)
	if err := newWorkflowCWL.Unmarshal("/home/thanhpp/go/src/github.com/vfluxus/demo-cwl/CWL_scatter_example/", "pipeline_step4.cwl"); err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(newWorkflowCWL)
}

func TestConvertCWLToDAG(t *testing.T) {
	newWorkflowCWL := new(workflowcwl.WorkflowCWL)
	if err := newWorkflowCWL.Unmarshal("/home/thanhpp/go/src/github.com/vfluxus/demo-cwl/scatter/", "scatter-wf.cwl"); err != nil {
		t.Fatal(err)
	}
	newWorkflowDAG, err := workflowdag.ConvertFromCWL(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(newWorkflowDAG)
}

func TestMainCwlAndDag(t *testing.T) {
	var (
		folder  string = "/home/tpp/go/src/github.com/vfluxus/transformer/test/basic/"
		cwlfile string = "1st-workflow.cwl"
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(newWorkflowDAG)
}

func TestAddValueToStepInAndArg(t *testing.T) {
	var (
		folder    = "/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs/"
		cwlfile   = "wgs.cwl"
		inputPath = "wgs.yml"
		data      = make(map[string]interface{})
		err       error
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}

	inputFile, err := ioutil.ReadFile(folder + inputPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(inputFile, data); err != nil {
		t.Fatal(err)
	}

	if err := workflowdag.AddValueToStepInAndArg(data, newWorkflowDAG); err != nil {
		t.Fatal(err)
	}

	libs.PrintJsonFormat(newWorkflowDAG)
}

func TestAddOutputToInput(t *testing.T) {
	var (
		folder    = "/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs/"
		cwlfile   = "mash.cwl"
		inputPath = "mash.yml"
		data      = make(map[string]interface{})
		err       error
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}

	inputFile, err := ioutil.ReadFile(folder + inputPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(inputFile, data); err != nil {
		t.Fatal(err)
	}

	if err := workflowdag.AddValueToStepInAndArg(data, newWorkflowDAG); err != nil {
		t.Fatal(err)
	}

	if err := workflowdag.AddOutputToInput(newWorkflowDAG); err != nil {
		t.Fatal(err)
	}

	libs.PrintJsonFormat(newWorkflowDAG)
}

func TestConvertWorkflowDAGToRun(t *testing.T) {
	var (
		folder    = "/home/thanhpp/go/src/github.com/vfluxus/demo-cwl/scatter/"
		cwlfile   = "scatter-wf.cwl"
		inputPath = "input.yml"
		userID    = "thanhpp"
		retry     = 0
		data      = make(map[string]interface{})
		err       error
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}

	inputFile, err := ioutil.ReadFile(folder + inputPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(inputFile, data); err != nil {
		t.Fatal(err)
	}

	if err := workflowdag.AddValueToStepInAndArg(data, newWorkflowDAG); err != nil {
		t.Fatal(err)
	}

	if err := workflowdag.AddOutputToInput(newWorkflowDAG); err != nil {
		t.Fatal(err)
	}

	newRun, err := workflowrun.ConvertWorkflowDAGToRun(newWorkflowDAG, userID, retry)
	if err != nil {
		t.Fatal(err)
	}

	libs.PrintJsonFormat(newRun)
}

func TestCreateRunFromWorkflow(t *testing.T) {
	var (
		folder    = "/home/tpp/go/src/github.com/vfluxus/demo-cwl/wes/"
		cwlfile   = "wes.cwl"
		inputPath = "wes.yml"
		userID    = "thanhphanphu18"
		runID     = 1000
		data      = make(map[string]interface{})
		err       error
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}

	inputFile, err := ioutil.ReadFile(folder + inputPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(inputFile, data); err != nil {
		t.Fatal(err)
	}

	newRun, err := CreateRunFromWorkflow(newWorkflowDAG, data, userID, runID)
	if err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(newRun)
}

func TestCreateGraphViz(t *testing.T) {
	var (
		folder    = "/home/tpp/go/src/github.com/vfluxus/demo-cwl/wes/"
		cwlfile   = "wes.cwl"
		inputPath = "wes.yml"
		userID    = "0"
		retry     = 0
		data      = make(map[string]interface{})
		err       error
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}

	inputFile, err := ioutil.ReadFile(folder + inputPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(inputFile, data); err != nil {
		t.Fatal(err)
	}

	newRun, err := CreateRunFromWorkflow(newWorkflowDAG, data, userID, retry)
	if err != nil {
		t.Fatal(err)
	}

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			t.Fatal(err)
		}
		g.Close()
	}()

	var (
		nodeMap = make(map[string]*cgraph.Node)
	)
	for taskIndex := range newRun.Tasks {
		node, err := graph.CreateNode(newRun.Tasks[taskIndex].TaskID)
		if err != nil {
			t.Fatal(err)
		}
		nodeMap[newRun.Tasks[taskIndex].TaskID] = node
	}

	for taskIndex := range newRun.Tasks {
		for childIndex := range newRun.Tasks[taskIndex].ChildrenTasksID {
			_, err := graph.CreateEdge("", nodeMap[newRun.Tasks[taskIndex].TaskID], nodeMap[newRun.Tasks[taskIndex].ChildrenTasksID[childIndex]])
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
}

func TestCreateGraphvizDot(t *testing.T) {
	var (
		folder    = "/home/tpp/go/src/github.com/vfluxus/demo-cwl/wes/"
		cwlfile   = "wes.cwl"
		inputPath = "wes.yml"
		userID    = "thanhpp18@gmail.com"
		retry     = 0
		data      = make(map[string]interface{})
		err       error
	)

	newWorkflowCWL, err := ParseCWL(folder, cwlfile)
	if err != nil {
		t.Fatal(err)
	}
	// libs.PrintJsonFormat(newWorkflowCWL)

	newWorkflowDAG, err := CreateWorkflowDAG(newWorkflowCWL)
	if err != nil {
		t.Fatal(err)
	}

	inputFile, err := ioutil.ReadFile(folder + inputPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(inputFile, data); err != nil {
		t.Fatal(err)
	}

	newRun, err := CreateRunFromWorkflow(newWorkflowDAG, data, userID, retry)
	if err != nil {
		t.Fatal(err)
	}

	dot, err := CreateGraphvizDot(newRun)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dot)
}

func TestParseCWLInMem(t *testing.T) {
	var (
		jsonPath     string = "/home/thanhpp/go/src/github.com/vfluxus/demo-cwl/CWL_scatter_example/createjson.json"
		jsonWorkflow        = new(workflowcwl.HttpCWLForm)
	)

	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		t.Error(err)
		return
	}

	if err := json.Unmarshal(data, jsonWorkflow); err != nil {
		t.Error(err)
		return
	}

	wf, err := ParseCWLInMem(jsonWorkflow)
	if err != nil {
		t.Error(err)
		return
	}

	libs.PrintJsonFormat(wf)
}
