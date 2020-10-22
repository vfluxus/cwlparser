package workflowdag

type Step struct {
	ID           string
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
	ParentID     []string
	ChildrenName []string
	ChildrenID   []string
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
	ValueFrom      *valueFrom
	Binding        *stepInputBinding
}

type valueFrom struct {
	Raw     string
	Prefix  string
	Postfix string
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
