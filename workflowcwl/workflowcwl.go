package workflowcwl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/vfluxus/cwlparser/commandlinetool"
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

// Unmarshal ...
func (wfCWL *WorkflowCWL) Unmarshal(folder string, file string) (err error) {
	if !strings.Contains(file, ".cwl") {
		return errors.New("Not cwl file")
	}
	// open workflow cwl file
	filePath := folder + "/" + file
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	// unmarshal
	if err := yaml.Unmarshal(fileData, wfCWL); err != nil {
		return err
	}
	// add children
	if err := wfCWL.addChildren(); err != nil {
		return err
	}
	// loop through step, read & unmarshal each file
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
		if err := addWorkflowNameAndFrom(wfCWL.Steps[stepIndex], newCmdLineTool); err != nil {
			return err
		}
		wfCWL.Steps[stepIndex].CommandLineTool = newCmdLineTool
	}

	return nil
}

// UnmarshalJson use in-mem data to unmarshal json
func (wfCWL *WorkflowCWL) UnmarshalJson(f *HttpCWLForm) (err error) {
	var (
		fileMap = make(map[string]*HttpStepForm)
	)

	// create map also, check for duplicate file name
	for stepIndex := range f.Steps {
		if _, ok := fileMap[f.Steps[stepIndex].Name]; ok {
			return fmt.Errorf("Duplicate file name %s", f.Steps[stepIndex].Name)
		}
		fileMap[f.Steps[stepIndex].Name] = f.Steps[stepIndex]
	}

	if err := yaml.Unmarshal([]byte(f.Content), wfCWL); err != nil {
		return err
	}

	// add children
	if err := wfCWL.addChildren(); err != nil {
		return err
	}

	// loop through step, read & unmarshal each step
	for stepIndex := range wfCWL.Steps {
		step, ok := fileMap[wfCWL.Steps[stepIndex].Run]
		if !ok {
			return fmt.Errorf("File not found %s", wfCWL.Steps[stepIndex].Run)
		}

		var (
			newCmdLineTool = new(commandlinetool.CommandLineTool)
		)
		if err := yaml.Unmarshal([]byte(step.Content), newCmdLineTool); err != nil {
			return err
		}
		if err := addWorkflowNameAndFrom(wfCWL.Steps[stepIndex], newCmdLineTool); err != nil {
			return err
		}
		wfCWL.Steps[stepIndex].CommandLineTool = newCmdLineTool
	}

	return nil
}

// addChildren create map[stepName]step and loop through all to add children
func (wfCWL *WorkflowCWL) addChildren() (err error) {
	if wfCWL == nil || len(wfCWL.Steps) < 1 {
		return errors.New("Empty workflow")
	}

	// create map and append to start node
	var (
		stepMap = make(map[string]*Step)
	)

	for stepIndex := range wfCWL.Steps {
		if _, ok := stepMap[wfCWL.Steps[stepIndex].Name]; ok {
			return fmt.Errorf("Duplicate step name: %v", wfCWL.Steps[stepIndex].Name)
		}
		stepMap[wfCWL.Steps[stepIndex].Name] = wfCWL.Steps[stepIndex]
	}

	for key := range stepMap {
		for parentIndex := range stepMap[key].Parents {
			stepMap[stepMap[key].Parents[parentIndex]].Children = append(stepMap[stepMap[key].Parents[parentIndex]].Children, stepMap[key].Name)
		}
	}

	return nil
}

func addWorkflowNameAndFrom(step *Step, cmdlinetool *commandlinetool.CommandLineTool) (err error) {
	var (
		stepInSourceMap = make(map[string]string)
	)
	for stepInIndex := range step.In {
		if _, ok := stepInSourceMap[step.In[stepInIndex].Name]; ok {
			return errors.New("Duplicate step input")
		}
		stepInSourceMap[step.In[stepInIndex].Name] = step.In[stepInIndex].Source
	}

	for cmdInputIndex := range cmdlinetool.Inputs {
		if source, ok := stepInSourceMap[cmdlinetool.Inputs[cmdInputIndex].Name]; ok {
			if sourceSl := strings.Split(source, "/"); len(sourceSl) == 2 {
				cmdlinetool.Inputs[cmdInputIndex].From = sourceSl[0]
				cmdlinetool.Inputs[cmdInputIndex].WorkflowName = sourceSl[1]
				continue
			}
			cmdlinetool.Inputs[cmdInputIndex].WorkflowName = source
		}
	}

	return nil
}
