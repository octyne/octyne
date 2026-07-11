package openai

import "encoding/json"

type FunctionDefinition struct {
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Parameters  *json.RawMessage `json:"parameters,omitempty"`
	Strict      *bool            `json:"strict,omitempty"`
}
type CustomToolGrammar struct {
	Definition string `json:"definition"`
	Syntax     string `json:"syntax"`
}
type CustomToolFormat struct {
	Type    string             `json:"type"`
	Grammar *CustomToolGrammar `json:"grammar,omitempty"`
}
type CustomToolDefinition struct {
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Format      *CustomToolFormat `json:"format,omitempty"`
}
type Tool struct {
	Type     string                `json:"type"`
	Function *FunctionDefinition   `json:"function,omitempty"`
	Custom   *CustomToolDefinition `json:"custom,omitempty"`
}
type NamedTool struct {
	Name string `json:"name"`
}
type ToolReference struct {
	Type     string     `json:"type"`
	Function *NamedTool `json:"function,omitempty"`
	Custom   *NamedTool `json:"custom,omitempty"`
}
type AllowedTools struct {
	Mode  string          `json:"mode"`
	Tools []ToolReference `json:"tools"`
}
type ToolChoice struct {
	Mode         *string       `json:"-"`
	Type         *string       `json:"type,omitempty"`
	Function     *NamedTool    `json:"function,omitempty"`
	Custom       *NamedTool    `json:"custom,omitempty"`
	AllowedTools *AllowedTools `json:"allowed_tools,omitempty"`
}

func (c ToolChoice) MarshalJSON() ([]byte, error) {
	if c.Mode != nil {
		return json.Marshal(*c.Mode)
	}
	type objectChoice ToolChoice
	return json.Marshal(objectChoice(c))
}

type LegacyFunctionDefinition struct {
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Parameters  *json.RawMessage `json:"parameters,omitempty"`
}
type LegacyFunctionCall struct {
	Mode *string `json:"-"`
	Name *string `json:"name,omitempty"`
}

func (c LegacyFunctionCall) MarshalJSON() ([]byte, error) {
	if c.Mode != nil {
		return json.Marshal(*c.Mode)
	}
	return json.Marshal(struct {
		Name *string `json:"name,omitempty"`
	}{Name: c.Name})
}
