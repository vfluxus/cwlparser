package main

import (
	"io/ioutil"

	"github.com/vfluxus/cwlparser/commandlinetool"
	"gopkg.in/yaml.v2"
)

const (
	argumentCwl  = "/home/tpp/go/src/github.com/vfluxus/cwlparser/test/arguments.cwl"
	applyBSQRCwl = "/home/tpp/go/src/github.com/vfluxus/cwlparser/test/ApplyBQSR.cwl"
)

func main() {
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
