package workflowdag

import (
	"fmt"

	"github.com/vfluxus/cwlparser/libs"
)

// AddValueToStepInAndArg ...
// TODO: Index process ???
func AddValueToStepInAndArg(inputs map[string]interface{}, wfDAG *WorkflowDAG) (err error) {
	for stepIndex := range wfDAG.Steps {
		for stepInIndex := range wfDAG.Steps[stepIndex].StepInput {
			if wfDAG.Steps[stepIndex].StepInput[stepInIndex].From == "" && wfDAG.Steps[stepIndex].StepInput[stepInIndex].WorkflowName != "" {
				if value, ok := inputs[wfDAG.Steps[stepIndex].StepInput[stepInIndex].WorkflowName]; ok {
					val, err := checkTypeAndAddValue(value, wfDAG.Steps[stepIndex].StepInput[stepInIndex].Type)
					if err != nil {
						return err
					}
					if wfDAG.Steps[stepIndex].StepInput[stepInIndex].ValueFrom != nil {
						for valIndex := range val {
							val[valIndex] = wfDAG.Steps[stepIndex].StepInput[stepInIndex].ValueFrom.Prefix + val[valIndex] + wfDAG.Steps[stepIndex].StepInput[stepInIndex].ValueFrom.Postfix
						}
					}
					wfDAG.Steps[stepIndex].StepInput[stepInIndex].Value = val
					continue
				}
				return fmt.Errorf("Can not find %v in %v", wfDAG.Steps[stepIndex].StepInput[stepInIndex].WorkflowName, inputs)
			}
		}

		for argIndex := range wfDAG.Steps[stepIndex].Arguments {
			if wfDAG.Steps[stepIndex].Arguments[argIndex].Input == nil {
				continue
			}
			if wfDAG.Steps[stepIndex].Arguments[argIndex].Input.From == "" && wfDAG.Steps[stepIndex].Arguments[argIndex].Input.WorkflowName != "" {
				if value, ok := inputs[wfDAG.Steps[stepIndex].Arguments[argIndex].Input.WorkflowName]; ok {
					val, err := checkTypeAndAddValue(value, wfDAG.Steps[stepIndex].Arguments[argIndex].Input.Type)
					if err != nil {
						return err
					}
					if wfDAG.Steps[stepIndex].Arguments[argIndex].Input.ValueFrom != nil {
						for valIndex := range val {
							val[valIndex] = wfDAG.Steps[stepIndex].Arguments[argIndex].Input.ValueFrom.Prefix + val[valIndex] + wfDAG.Steps[stepIndex].Arguments[argIndex].Input.ValueFrom.Postfix
						}
					}
					wfDAG.Steps[stepIndex].Arguments[argIndex].Input.Value = val
					continue
				}
				return fmt.Errorf("Can not find %v in %v", wfDAG.Steps[stepIndex].Arguments[argIndex].Input.WorkflowName, inputs)
			}
		}
	}
	return nil
}

// checkTypeAndAddValue ... MISSING: CHECK TYPE
func checkTypeAndAddValue(val interface{}, valType []string) (valStr []string, err error) {
	switch valCast := val.(type) {
	case string:
		return append(valStr, valCast), nil

	case []string:
		return append(valStr, valCast...), nil

	case []interface{}:
		if err := libs.AppendStringSliceWithInterface(&valStr, valCast); err != nil {
			return nil, err
		}
		return valStr, nil

	case map[interface{}]interface{}:
		var (
			fileFlag bool = false
		)
		if class, ok := valCast["class"]; ok {
			switch classCast := class.(type) {
			case string:
				if classCast == "File" {
					fileFlag = true
				}

				if fileFlag {
					if path, ok := valCast["path"]; ok {
						if err := libs.AppendStringSliceWithInterface(&valStr, path); err != nil {
							return nil, err
						}
						return valStr, nil
					}
					return nil, fmt.Errorf("path not found in: %v", valCast)
				}

			default:
				return nil, fmt.Errorf("Not support class %v yet", class)
			}
		} else {
			return nil, fmt.Errorf("Not support input type: %v", valCast)
		}

	default:
		return nil, fmt.Errorf("Can not cast value %v - type: %T", val, val)
	}

	return valStr, nil
}
