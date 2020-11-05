package workflowcwl

type HttpCWLForm struct {
	Name    string          `json:"name"`
	Content []byte          `json:"string"`
	Steps   []*HttpStepForm `json:"steps"`
}

type HttpStepForm struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}
