package server

import (
	openaicompat "github.com/octyne/octyne/internal/compat/openai"
	"github.com/octyne/octyne/internal/types"
)

func toCanonicalTools(value *[]openaicompat.Tool) *[]types.Tool {
	if value == nil {
		return nil
	}
	tools := make([]types.Tool, len(*value))
	for i, tool := range *value {
		tools[i] = types.Tool{Type: tool.Type}
		if tool.Function != nil {
			tools[i].Function = &types.FunctionDefinition{Name: tool.Function.Name, Description: tool.Function.Description, Parameters: tool.Function.Parameters, Strict: tool.Function.Strict}
		}
		if tool.Custom != nil {
			tools[i].Custom = &types.CustomToolDefinition{Name: tool.Custom.Name, Description: tool.Custom.Description}
			if tool.Custom.Format != nil {
				tools[i].Custom.Format = &types.CustomToolFormat{Type: tool.Custom.Format.Type}
				if tool.Custom.Format.Grammar != nil {
					tools[i].Custom.Format.Grammar = &types.CustomToolGrammar{Definition: tool.Custom.Format.Grammar.Definition, Syntax: tool.Custom.Format.Grammar.Syntax}
				}
			}
		}
	}
	return &tools
}

func toCanonicalToolChoice(value *openaicompat.ToolChoice) *types.ToolChoice {
	if value == nil {
		return nil
	}
	choice := &types.ToolChoice{Mode: value.Mode, Type: value.Type}
	if value.Function != nil {
		choice.Function = &types.NamedTool{Name: value.Function.Name}
	}
	if value.Custom != nil {
		choice.Custom = &types.NamedTool{Name: value.Custom.Name}
	}
	if value.AllowedTools != nil {
		choice.AllowedTools = &types.AllowedTools{Mode: value.AllowedTools.Mode, Tools: make([]types.ToolReference, len(value.AllowedTools.Tools))}
		for i, ref := range value.AllowedTools.Tools {
			choice.AllowedTools.Tools[i] = types.ToolReference{Type: ref.Type}
			if ref.Function != nil {
				choice.AllowedTools.Tools[i].Function = &types.NamedTool{Name: ref.Function.Name}
			}
			if ref.Custom != nil {
				choice.AllowedTools.Tools[i].Custom = &types.NamedTool{Name: ref.Custom.Name}
			}
		}
	}
	return choice
}

func toCanonicalLegacyFunctions(value *[]openaicompat.LegacyFunctionDefinition) *[]types.LegacyFunctionDefinition {
	if value == nil {
		return nil
	}
	functions := make([]types.LegacyFunctionDefinition, len(*value))
	for i, function := range *value {
		functions[i] = types.LegacyFunctionDefinition{Name: function.Name, Description: function.Description, Parameters: function.Parameters}
	}
	return &functions
}

func toCanonicalLegacyFunctionCall(value *openaicompat.LegacyFunctionCall) *types.LegacyFunctionCall {
	if value == nil {
		return nil
	}
	return &types.LegacyFunctionCall{Mode: value.Mode, Name: value.Name}
}
