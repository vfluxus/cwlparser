package main

import (
	"io/ioutil"
	"testing"

	"github.com/vfluxus/cwlparser/workflowdag"

	"gopkg.in/yaml.v2"

	"github.com/vfluxus/cwlparser/commandlinetool"
	"github.com/vfluxus/cwlparser/libs"
	"github.com/vfluxus/cwlparser/workflowcwl"
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
	data, err := ioutil.ReadFile("/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs/wgs.cwl")
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

	data2, err2 := ioutil.ReadFile("/home/tpp/go/src/github.com/vfluxus/transformer/test/basic/1st-workflow.cwl")
	if err2 != nil {
		t.Fatal(err2)
	}
	var (
		wfCwl2 = new(workflowcwl.WorkflowCWL)
	)
	if err2 := yaml.Unmarshal(data2, wfCwl2); err != nil {
		t.Fatal(err2)
	}
	libs.PrintJsonFormat(wfCwl2)
}

func TestWorkflowCWLUnmarshal(t *testing.T) {
	newWorkflowCWL := new(workflowcwl.WorkflowCWL)
	if err := newWorkflowCWL.Unmarshal("/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs", "wgs.cwl"); err != nil {
		t.Fatal(err)
	}
	libs.PrintJsonFormat(newWorkflowCWL)
}

func TestConvertCWLToDAG(t *testing.T) {
	newWorkflowCWL := new(workflowcwl.WorkflowCWL)
	if err := newWorkflowCWL.Unmarshal("/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs", "wgs.cwl"); err != nil {
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
