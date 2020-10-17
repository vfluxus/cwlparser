package workflowcwl

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/vfluxus/cwlparser/commandlinetool"

	"gopkg.in/yaml.v2"
)

type WorkflowCWL struct {
	Version      string         `yaml:"cwlVersion"`
	Doc          string         `yaml:"doc"`
	ID           string         `yaml:"id"`
	Requirements []*requirement `yaml:"requirements"`
	Inputs       inputs         `yaml:"inputs"`
	Outputs      outputs        `yaml:"outputs"`
	Steps        steps          `yaml:"steps"`
}

type requirement struct {
	Class string `yaml:"class"`
}

func (wfCWL *WorkflowCWL) Unmarshal(folder string, file string) (err error) {
	if !strings.Contains(file, ".cwl") {
		return errors.New("Not cwl file")
	}

	filePath := folder + "/" + file
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(fileData, wfCWL); err != nil {
		return err
	}

	for stepIndex := range wfCWL.Steps {
		stepFilePath := folder + "/" + wfCWL.Steps[stepIndex].Run
		stepFileData, err := ioutil.ReadFile(stepFilePath)
		if err != nil {
			return err
		}

		var (
			newCmdLineTool = new(commandlinetool.CommandLineTool)
		)
		if err := yaml.Unmarshal(stepFileData, newCmdLineTool); err != nil {
			return err
		}
		wfCWL.Steps[stepIndex].CommandLineTool = newCmdLineTool
	}

	return nil
}
