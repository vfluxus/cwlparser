# CWL PARSER

## Table of contents:
- [CWL PARSER](#cwl-parser)
  - [Table of contents:](#table-of-contents)
  - [1. DESCRIPTION](#1-description)
  - [2. CWL EXAMPLES](#2-cwl-examples)
    - [2.1.Workflow](#21workflow)
    - [2.2. CommandLineTool](#22-commandlinetool)
  - [3. CWL struct:](#3-cwl-struct)
    - [3.1. CWL CommandLineTool:](#31-cwl-commandlinetool)
    - [3.2. Workflow CWL:](#32-workflow-cwl)
  - [4. Convert to DAG:](#4-convert-to-dag)
    - [4.1. Workflow DAG Struct:](#41-workflow-dag-struct)
    - [4.2. How to convert:](#42-how-to-convert)
  - [5. Convert to run:](#5-convert-to-run)
    - [5.1. Run struct:](#51-run-struct)
    - [5.2. Convert Workflow DAG to Run:](#52-convert-workflow-dag-to-run)
      - [5.2.1. Convert Workflow DAG to Run:](#521-convert-workflow-dag-to-run)
      - [5.2.2. Convert step to task:](#522-convert-step-to-task)
  - [6. Not support these things yet:](#6-not-support-these-things-yet)
  - [7. Soon support:](#7-soon-support)

___

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

## 3. CWL struct:

### 3.1. CWL CommandLineTool:

- [CommandLineTool](commandlinetool/commandlinetool.go)
- [Base Command](commandlinetool/basecommand.go)
- [Arguments](commandlinetool/arguments.go)
- [Inputs](commandlinetool/inputs.go)
- [Output](commandlinetool/outputs.go)

### 3.2. Workflow CWL:

- [Workflow cwl](workflowcwl/workflowcwl.go)
- [Steps](workflowcwl/steps.go)
- [Inputs](workflowcwl/inputs.go)
- [Outputs](workflowcwl/outputs.go)
- [HTTP Form](workflowcwl/httpForm.go)

## 4. Convert to DAG:

### 4.1. Workflow DAG Struct:

- [workflow dag](workflowdag/workflowdag.go)
- [steps](workflowdag/step.go)

### 4.2. How to convert:
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
    -** #TODO: Add step inputs to arguments**

## 5. Convert to run:

### 5.1. Run struct:

- [Run](workflowrun/workflowrun.go)

### 5.2. Convert Workflow DAG to Run:

#### 5.2.1. Convert Workflow DAG to Run:
- ID generate:
  - WorkflowID
  - UserID
  - Retry time
  
- Easy converted:
  - Workflow ID
  - User ID
  - Status: default 0

#### 5.2.2. Convert step to task:
- ID generate:
  - RunID
  - Sequence
  - Task name

- ID special:
  - Start node: RunID - **bigbang**
  - End node: RunID - **ragnarok**
  
- Easy converted:
  - TaskID
  - StepID
  - ChildrenTaskID: init by StepID, replaced by TaskID
  - ParentTaskID: init by StepID, replaced by TaskID
  - Command: Join baseCommand
  - DockerImages:

- Convert Argument to param with regex struct
  - Add prefix
  - Add from for param
  - Add secondary files, value, prefix

## 6. Not support these things yet:
- Scatter
- $(runtime.something)
- add inputs value in step output

## 7. Soon support:
- Inputs index
