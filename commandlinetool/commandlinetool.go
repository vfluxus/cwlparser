package commandlinetool

import (
	"errors"
)

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

type requirement struct {
	Class      string `yaml:"class"`
	DockerPull string `yaml:"dockerPull"`
	RamMin     int    `yaml:"ramMin"`
	CpuMin     int    `yaml:"cpuMin"`
}

// assert version v1.0
type version string

func (v *version) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		str version
	)
	// unmarshal
	if err := unmarshal(&str); err != nil {
		return err
	}
	// version assert
	if str != "v1.0" {
		return errors.New("Only support version v1.0")
	}

	*v = str

	return nil
}

// assert class CommandLineTool
type class string

func (c *class) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var (
		str class
	)
	// unmarshal
	if err := unmarshal(&str); err != nil {
		return err
	}
	// version assert
	if str != "CommandLineTool" {
		return errors.New("Not class CommandLineTool")
	}

	*c = str

	return nil
}
