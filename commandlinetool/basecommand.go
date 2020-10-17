package commandlinetool

import (
	"fmt"
)

type baseCommand []string

func (bc *baseCommand) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		strSlice []string
	)

	// unmarshal to string
	var (
		strBaseCommand string
		err1           error
	)
	if err1 = unmarshal(&strBaseCommand); err1 == nil {
		strSlice = append(strSlice, strBaseCommand)
		*bc = strSlice
		return nil
	}

	// unmarshal to []string
	var (
		strSlBaseCommand []string
		err2             error
	)
	if err2 = unmarshal(&strSlBaseCommand); err2 == nil {
		return nil
	}

	return fmt.Errorf("BaseCommand - Unmarshal to string error: %v. Unmarshal to []string error: %v", err1, err2)
}
