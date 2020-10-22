package workflowrun

type Run struct {
	WorkflowID uint    `json:"workflow_id"`
	RunID      string  `json:"run_id"`
	User_ID    int     `json:"user_id"`
	Status     int     `json:"status"`
	Tasks      []*Task `json:"tasks"`
}

type Task struct {
	TaskID          string            `json:"task_id"`
	IsBoundary      bool              `json:"is_boundary"`
	StepID          string            `json:"step_id"`
	RunID           string            `json:"run_id"`
	UserID          int               `json:"user_id"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	ParentTasksID   []string          `json:"parent_tasks_id"`
	ChildrenTasksID []string          `json:"children_tasks_id"`
	OutputLocation  []string          `json:"output_location"`
	DockerImage     []string          `json:"docker_image"`
	QueueLevel      int               `json:"queue_level"`
	Status          int               `json:"status"`
}

type ParamWithRegex struct {
	From   []string `json:"from"`
	Prefix string   `json:"prefix"`
	Regex  []string `json:"regex"`
}
