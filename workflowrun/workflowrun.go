package workflowrun

type Run struct {
	RunID    int     `json:"run_id"`
	RunName  string  `json:"run_name"`
	UserName string  `json:"username"`
	Tasks    []*Task `json:"tasks"`
}

type Task struct {
	TaskID          string            `json:"task_id"`
	TaskName        string            `json:"task_name"`
	StepID          string            `json:"-"`
	UserName        string            `json:"username"`
	Command         string            `json:"command"`
	ScatterMethod   string            `json:"scatter_method"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	OutputRegex     []string          `json:"output_regex"`
	Output2ndFiles  []string          `json:"output_2nd_files"`
	ParentTasksID   []string          `json:"parent_tasks_id"`
	ChildrenTasksID []string          `json:"children_tasks_id"`
	OutputLocation  []string          `json:"output_location"`
	DockerImage     []string          `json:"docker_image"`
	IsBoundary      bool              `json:"is_boundary"`
	RunID           int               `json:"run_id"`
	QueueLevel      int               `json:"queue_level"`
	Status          int               `json:"status"`
}

type ParamWithRegex struct {
	Scatter        bool     `json:"scatter"`
	From           []string `json:"from"`
	SecondaryFiles []string `json:"secondary_files"`
	Regex          []string `json:"regex"`
	Prefix         string   `json:"prefix"`
}
