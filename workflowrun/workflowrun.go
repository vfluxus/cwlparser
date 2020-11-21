package workflowrun

type Run struct {
	WorkflowID string  `json:"workflow_id"`
	RunID      string  `json:"run_id"`
	RunName    string  `json:"run_name"`
	UserName   string  `json:"username"`
	Status     int     `json:"status"`
	Tasks      []*Task `json:"tasks"`
}

type Task struct {
	TaskID          string            `json:"task_id"`
	TaskName        string            `json:"task_name"`
	IsBoundary      bool              `json:"is_boundary"`
	StepID          string            `json:"step_id"`
	RunID           string            `json:"run_id"`
	UserName        string            `json:"username"`
	Command         string            `json:"command"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	ParentTasksID   []string          `json:"parent_tasks_id"`
	ChildrenTasksID []string          `json:"children_tasks_id"`
	OutputLocation  []string          `json:"output_location"`
	DockerImage     []string          `json:"docker_image"`
	QueueLevel      int               `json:"queue_level"`
	Status          int               `json:"status"`
}

type ParamWithRegex struct {
	From           []string `json:"from"`
	SecondaryFiles []string `json:"secondary_files"`
	Regex          []string `json:"regex"`
	Prefix         string   `json:"prefix"`
}
