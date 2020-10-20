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
	Input     *stepInput
	Prefix    string
}
