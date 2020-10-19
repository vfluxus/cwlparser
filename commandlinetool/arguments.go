package commandlinetool

import (
	"fmt"
)

type arguments []*Argument

type Argument struct {
	Position   int    `yaml:"position"`
	ShellQuote bool   `yaml:"shellQuote"`
	ValueFrom  string `yaml:"valueFrom"`
}

func (args *arguments) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		err1  error // unmarshal to argument struct
		err2  error // unmarshal to []string
		argSl []*Argument
	)
	// unmarshal to argument struct
	if err1 = unmarshal(&argSl); err1 == nil {
		*args = argSl
		return nil
	}

	// unmarshall to []string
	var (
		strSl []string
	)
	if err2 = unmarshal(&strSl); err2 == nil {
		for strSlIndex := range strSl {
			newArg := &Argument{
				ValueFrom: strSl[strSlIndex],
			}
			argSl = append(argSl, newArg)
		}
		*args = argSl
		return nil
	}

	return fmt.Errorf("Arguments - Can not unmarshal to argument struct: %v. To []string: %v", err1, err2)
}
