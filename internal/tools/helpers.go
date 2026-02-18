package tools

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"

	"tp-mcp-go/internal/domain/errors"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// ptr returns a pointer to the given value
func ptr[T any](v T) *T {
	return &v
}

// errorResult creates an error CallToolResult with actionable messages for API errors
func errorResult(err error) *mcp.CallToolResult {
	var apiErr *errors.APIError
	if stderrors.As(err, &apiErr) {
		return apiErrorResult(apiErr)
	}

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

// apiErrorResult formats an APIError with actionable hints
func apiErrorResult(apiErr *errors.APIError) *mcp.CallToolResult {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TP API Error (HTTP %d): %s", apiErr.StatusCode, apiErr.Message))

	switch {
	case apiErr.StatusCode == 400:
		sb.WriteString("\n\nThis usually means the query syntax is invalid. Common issues:")
		sb.WriteString("\n- Boolean values must be single-quoted strings: EntityState.IsFinal eq 'true', EntityState.IsFinal eq 'false'")
		sb.WriteString("\n- Use 'in' for multiple values: EntityState.Name in ('Open','In Progress') or Id in (1,2,3)")
		sb.WriteString("\n- String values must be single-quoted: EntityState.Name eq 'Open'")
		sb.WriteString("\n- Verify field names are valid for this entity type (use inspect_object tool)")
		sb.WriteString("\n- Collection queries use .Any() syntax: Assignments.Any(GeneralUser.Id eq 123)")
		sb.WriteString("\n- Sorting: orderBy only accepts a field name (e.g., 'CreateDate'), not 'CreateDate desc' — use orderByField and orderByDirection parameters")
	case apiErr.StatusCode == 401:
		sb.WriteString("\n\nAuthentication failed. The access token may be invalid or expired.")
	case apiErr.StatusCode == 403:
		sb.WriteString("\n\nPermission denied. The current user may not have access to this resource.")
	case apiErr.StatusCode == 404:
		sb.WriteString("\n\nEntity not found. Verify the entity type and ID are correct.")
	case apiErr.StatusCode >= 500:
		sb.WriteString("\n\nServer error on the TP side. This is usually temporary — try the request again.")
	}

	if apiErr.Context != "" {
		sb.WriteString(fmt.Sprintf("\n\nRequest: %s", apiErr.Context))
	}

	isErr := true
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: sb.String(),
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
