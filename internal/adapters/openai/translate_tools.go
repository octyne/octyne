package openai

import "github.com/octyne/octyne/internal/types"

func toTools(value *[]types.Tool) *[]Tool {
	if value == nil {
		return nil
	}
	tools := make([]Tool, len(*value))
	for i, tool := range *value {
		tools[i] = Tool{Type: tool.Type}
		if tool.Function != nil {
			tools[i].Function = &FunctionDefinition{Name: tool.Function.Name, Description: tool.Function.Description, Parameters: tool.Function.Parameters, Strict: tool.Function.Strict}
		}
		if tool.Custom != nil {
			tools[i].Custom = &CustomToolDefinition{Name: tool.Custom.Name, Description: tool.Custom.Description}
			if tool.Custom.Format != nil {
				tools[i].Custom.Format = &CustomToolFormat{Type: tool.Custom.Format.Type}
				if tool.Custom.Format.Grammar != nil {
					tools[i].Custom.Format.Grammar = &CustomToolGrammar{Definition: tool.Custom.Format.Grammar.Definition, Syntax: tool.Custom.Format.Grammar.Syntax}
				}
			}
		}
	}
	return &tools
}

func toToolChoice(value *types.ToolChoice) *ToolChoice {
	if value == nil {
		return nil
	}
	choice := &ToolChoice{Mode: value.Mode, Type: value.Type}
	if value.Function != nil {
		choice.Function = &NamedTool{Name: value.Function.Name}
	}
	if value.Custom != nil {
		choice.Custom = &NamedTool{Name: value.Custom.Name}
	}
	if value.AllowedTools != nil {
		choice.AllowedTools = &AllowedTools{Mode: value.AllowedTools.Mode, Tools: make([]ToolReference, len(value.AllowedTools.Tools))}
		for i, ref := range value.AllowedTools.Tools {
			choice.AllowedTools.Tools[i] = ToolReference{Type: ref.Type}
			if ref.Function != nil {
				choice.AllowedTools.Tools[i].Function = &NamedTool{Name: ref.Function.Name}
			}
			if ref.Custom != nil {
				choice.AllowedTools.Tools[i].Custom = &NamedTool{Name: ref.Custom.Name}
			}
		}
	}
	return choice
}

func toLegacyFunctions(value *[]types.LegacyFunctionDefinition) *[]LegacyFunctionDefinition {
	if value == nil {
		return nil
	}
	functions := make([]LegacyFunctionDefinition, len(*value))
	for i, function := range *value {
		functions[i] = LegacyFunctionDefinition{Name: function.Name, Description: function.Description, Parameters: function.Parameters}
	}
	return &functions
}

func toLegacyFunctionCall(value *types.LegacyFunctionCall) *LegacyFunctionCall {
	if value == nil {
		return nil
	}
	return &LegacyFunctionCall{Mode: value.Mode, Name: value.Name}
}
