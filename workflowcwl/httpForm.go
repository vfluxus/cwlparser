package workflowcwl

type HttpCWLForm struct {
	RunID   int             `json:"run_id"`
	Name    string          `json:"name"`
	Content string          `json:"content"`
	Steps   []*HttpStepForm `json:"steps"`
}

type HttpStepForm struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}
