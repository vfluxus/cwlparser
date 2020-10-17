package main

import (
	"io/ioutil"

	"github.com/vfluxus/cwlparser/workflowcwl"

	"github.com/vfluxus/cwlparser/commandlinetool"
	"gopkg.in/yaml.v2"
)

const (
	argumentCwl  = "/home/tpp/go/src/github.com/vfluxus/cwlparser/test/arguments.cwl"
	applyBSQRCwl = "/home/tpp/go/src/github.com/vfluxus/cwlparser/test/ApplyBQSR.cwl"
)

func TestCmdLineTool() {
	data, err := ioutil.ReadFile(applyBSQRCwl)
	if err != nil {
		panic(err)
	}

	var (
		cmdTool = new(commandlinetool.CommandLineTool)
	)
	if err := yaml.Unmarshal(data, cmdTool); err != nil {
		panic(err)
	}
	PrintJsonFormat(cmdTool)

	data2, err2 := ioutil.ReadFile(argumentCwl)
	if err2 != nil {
		panic(err2)
	}

	var (
		cmdTool2 = new(commandlinetool.CommandLineTool)
	)
	if err2 := yaml.Unmarshal(data2, cmdTool2); err2 != nil {
		panic(err2)
	}
	PrintJsonFormat(cmdTool2)
}

func TestWorkflowCWL() {
	data, err := ioutil.ReadFile("/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs/wgs.cwl")
	if err != nil {
		panic(err)
	}
	var (
		wfCwl = new(workflowcwl.WorkflowCWL)
	)
	if err := yaml.Unmarshal(data, wfCwl); err != nil {
		panic(err)
	}
	PrintJsonFormat(wfCwl)

	data2, err2 := ioutil.ReadFile("/home/tpp/go/src/github.com/vfluxus/transformer/test/basic/1st-workflow.cwl")
	if err2 != nil {
		panic(err2)
	}
	var (
		wfCwl2 = new(workflowcwl.WorkflowCWL)
	)
	if err2 := yaml.Unmarshal(data2, wfCwl2); err != nil {
		panic(err2)
	}
	PrintJsonFormat(wfCwl2)
}

func TestWorkflowCWLUnmarshal() {
	newWorkflowCWL := new(workflowcwl.WorkflowCWL)
	if err := newWorkflowCWL.Unmarshal("/home/tpp/go/src/github.com/vfluxus/demo-cwl/wgs", "wgs.cwl"); err != nil {
		panic(err)
	}
	PrintJsonFormat(newWorkflowCWL)
}
func main() {
	TestWorkflowCWLUnmarshal()
}
