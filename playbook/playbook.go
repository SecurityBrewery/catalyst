package playbook

type Playbook struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Name        string     `json:"name"`
	Type        string     `json:"type,omitempty"`
	Description string     `json:"description"`
	Schema      JSONSchema `json:"schema,omitempty"`
}

type JSONSchema struct {
	Type       string                  `json:"type"`
	Properties map[string]JSONProperty `json:"properties"`
}

type JSONProperty struct {
	Title string `json:"title"`
	Type  string `json:"type"`
}
