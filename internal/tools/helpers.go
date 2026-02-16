package tools

import (
	"encoding/json"
	"fmt"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// ptr returns a pointer to the given value
func ptr[T any](v T) *T {
	return &v
}

// errorResult creates an error CallToolResult
func errorResult(err error) *mcp.CallToolResult {
	isErr := true
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Error: %s", err.Error()),
			},
		},
		IsError: &isErr,
	}
}

// jsonResult marshals data to JSON and returns as text content
func jsonResult(data any) *mcp.CallToolResult {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errorResult(err)
	}
	return textResult(string(jsonBytes))
}

// textResult returns text as a CallToolResult
func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

// getStringArg extracts a string argument from the args map
func getStringArg(args map[string]any, key string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getIntArg extracts an integer argument (handles float64 from JSON)
func getIntArg(args map[string]any, key string) (int, error) {
	v, ok := args[key]
	if !ok {
		return 0, fmt.Errorf("missing required argument: %s", key)
	}
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	case json.Number:
		i, err := n.Int64()
		return int(i), err
	default:
		return 0, fmt.Errorf("argument %s must be a number, got %T", key, v)
	}
}

// getStringSliceArg extracts a string slice argument
func getStringSliceArg(args map[string]any, key string) []string {
	v, ok := args[key]
	if !ok {
		return nil
	}
	if arr, ok := v.([]interface{}); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// getAnyArg extracts an argument as-is
func getAnyArg(args map[string]any, key string) any {
	return args[key]
}
