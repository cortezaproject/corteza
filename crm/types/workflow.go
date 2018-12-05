package types

type (
	Workflow struct {
		ID      uint64          `json:"workflowID,string"`
		Name    string          `json:"name"`
		Tasks   WorkflowTaskSet `json:"tasks"`
		OnError WorkflowTaskSet `json:"onError"`
		Timeout int             `json:"timeout"`
	}

	// WorkflowTask denotes a step in the workflow
	//
	// When it comes to Body and Fallback, we may invoke anything here. If we'll be doing this
	// via os.Exec, then it makes sense to enforce some inputs and outputs, in the style of FAAS.
	// If we're integrating against airflow, we need to provide basic credentials that allow each
	// workflow step to proceed. The dependencies like a valid JWT sholud somehow be provisioned.
	WorkflowTask struct {
		Name     string `json:"name"`
		Body     string `json:"body"`
		Fallback string `json:"fallback"`
		Retries  []int  `json:"retries"`
		Timeout  int    `json:"timeout"`
	}

	WorkflowTaskSet []WorkflowTask
)
