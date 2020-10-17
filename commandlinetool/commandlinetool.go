package commandlinetool

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
	RamMin     string `yaml:"ramMin"`
	CpuMin     string `yaml:"cpuMin"`
}
