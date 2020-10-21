# CWL PARSER

- [CWL PARSER](#cwl-parser)
  - [1. DESCRIPTION](#1-description)
  - [2. CWL EXAMPLES](#2-cwl-examples)
    - [2.1.Workflow](#21workflow)
    - [2.2. CommandLineTool](#22-commandlinetool)
  - [3. Golang struct](#3-golang-struct)
    - [3.1. CommandLineTool](#31-commandlinetool)
    - [3.2. Workflow](#32-workflow)
  - [4. Convert to dag:](#4-convert-to-dag)
    - [4.1. Workflow DAG Struct](#41-workflow-dag-struct)
    - [4.2. How to convert](#42-how-to-convert)

## 1. DESCRIPTION
- Parse CWL to Golang struct
  - using gopkg.in/yaml.v2 with custom UnmarshalYAML
- Support
    - version: v1.0
    - class:
      - Workflow
      - CommandLineTool
  
## 2. CWL EXAMPLES

### 2.1.Workflow
```yaml
cwlVersion: v1.0
doc: string
class: Workflow
requirements:
    - class: string
    - class: string

inputs:
    input1:
        type: string
        secondaryFiles: string
    input2:
        type: 
            - type1
            - type2
        secondaryFiles: [string1, string2]

outputs:
    output1:
        type: string
        outputSource: string
    output2:
        type: [string]
        outputSource: [string]

steps:
    step1:
        run: string
        in:
            stepin1: string
            stepin2:
                source: string
                valueFrom: string
        out:
            [string]
```

### 2.2. CommandLineTool
```yaml
cwlVersion: v1.0
class: CommandLineTool

requirements:
    - class: DockerRequirement
        dockerPull: string
    - class: ResourceRequirement 
        ramMin: int
        cpuMin: int

inputs:
    input1:
        type: string
        SecondaryFiles: [string]
        inputBinding:
            position: int
    input2:
        type: 
            - type1
            - type2

outputs:
    output1:
        type: string
        outputBinding:
            glob: string
        secondaryFiles: [string]
    output2:
        type:
            - type1
            - type2
        outputBinding:
            glob: 
                - pattern1
                - pattern2

arguments:
    - position: int
    shellQuote: bool
    valueFrom: string

baseCommand: [string]
```

## 3. Golang struct

### 3.1. CommandLineTool
```go
type CommandLineTool struct {
	Version      string         `yaml:"cwlVersion"`
	Class        string         `yaml:"class"`
	ID           string         `yaml:"id"`
	Requirements []*requirement `yaml:"requirements"`
	Inputs       inputs         `yaml:"inputs"`
	BaseCommand  baseCommand    `yaml:"baseCommand"`
	Arguments    arguments      `yaml:"arguments"`
	Outputs      outputs        `yaml:"outputs"`
}

// requirement
type requirement struct {
	Class      string `yaml:"class"`
	DockerPull string `yaml:"dockerPull"`
	RamMin     int    `yaml:"ramMin"`
	CpuMin     int    `yaml:"cpuMin"`
}

// inputs
type inputs []*input

type input struct {
	Name           string `yaml:""`
	WorkflowName   string
	From           string
	Type           []string      `yaml:""`
	SecondaryFiles []string      `yaml:""`
	InputBinding   *inputBinding `yaml:"inputBinding"`
}

type inputBinding struct {
	Position int    `yaml:"position"`
	Prefix   string `yaml:"prefix"`
}

// baseCommand
type baseCommand []string

// arguments
type arguments []*argument

type argument struct {
	Position   int    `yaml:"position"`
	ShellQuote bool   `yaml:"shellQuote"`
	ValueFrom  string `yaml:"valueFrom"`
}

// outputs
type outputs []*output

type output struct {
	Name           string        `yaml:""`
	Type           []string      `yaml:"type"`
	OutputBinding  outputBinding `yaml:"outputBinding"`
	SecondaryFiles []string      `yaml:"secondaryFiles"`
}

type outputBinding struct {
	Glob []string `yaml:"glob"`
}
```

### 3.2. Workflow
```go
type WorkflowCWL struct {
	Version      string         `yaml:"cwlVersion"`
	Doc          string         `yaml:"doc"`
	ID           string         `yaml:"id"`
	Requirements []*requirement `yaml:"requirements"`
	Inputs       inputs         `yaml:"inputs"`
	Outputs      outputs        `yaml:"outputs"`
	Steps        steps          `yaml:"steps"`
}

// requirement
type requirement struct {
	Class string `yaml:"class"`
}

// inputs
type inputs []*input

type input struct {
	Name           string
	Type           []string
	SecondaryFiles []string
}

// outputs
type outputs []*output

type output struct {
	Name         string
	Type         []string
	OutputSource []string
}

// steps
type steps []*step

type Step struct {
	Name            string
	Run             string
	Scatter         string `yaml:"scatter"`
	Parents         []string
	Children        []string
	In              stepIns  `yaml:"in"`
	Out             stepOuts `yaml:"out"`
	CommandLineTool *commandlinetool.CommandLineTool
}
// step in
type stepIns []*stepIn

type stepIn struct {
	Name      string
	Source    string
	ValueFrom string
}

// step out
type stepOuts []*stepOut

type stepOut struct {
	Name string
}
```

## 4. Convert to dag:

### 4.1. Workflow DAG Struct

```go
// WorkflowDAG 
type WorkflowDAG struct {
	Name  string  `json:"name"`
	Steps []*Step `json:"steps"`
}

// Step
package workflowdag

type Step struct {
	ID           int
	Name         string
	WorkflowName string
	DockerImage  string
	Resource     struct {
		CPU int
		Ram int
	}
	BaseCommand  []string
	StepInput    []*stepInput
	StepOutput   []*stepOutput
	Arguments    []*Argument
	ParentName   []string
	ParentID     []int
	ChildrenName []string
	ChildrenID   []int
}

type stepInput struct {
	Name           string
	WorkflowName   string
	From           string
	Type           []string
	SecondaryFiles []string
	Value          []string
	NullFlag       bool
	Regex          bool
	Binding        *stepInputBinding
}

type stepInputBinding struct {
	Postition int
	Prefix    string
}

type stepOutput struct {
	Name           string
	WorkflowName   string
	Type           []string
	Patern         []string
	Regex          []string
	SecondaryFiles []string
}

type Argument struct {
    Postition int
    Index     int
	Input     *stepInput
	Prefix    string
}
```

### 4.2. How to convert
- Convert step CWL to step DAG
    -  Easy converted:
        - ID
        - Name: Run
        - WorkflowName: stepCWL.Name
        - Children
        - Parent
        - BaseCommand
    - Convert Requirements:
        - Docker images
        - CPU, RAM
    - Convert Step input:
        - Name, WorkflowName, From, Type, Secondary Files
        - Check input binding
    - Convert Step ouput:
        - Name, Type, SecondaryFiles, Regex
    - Add arguments:
        - Generate from arguments.ValueFrom
        - Link to stepInputs (if $(inputs.[stepInput]))
        - Each argument only have 1 step Inputs. 1 Argument.ValueFrom can be separated to multiple arguments
    - #TODO: Add step inputs to arguments
