package tools

import (
	"context"
	"fmt"

	"tp-mcp-go/internal/client"
	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// entityTypeStrings converts ValidTypes to a slice of interface{} for schema enum
func entityTypeStrings() []interface{} {
	types := make([]interface{}, len(entity.ValidTypes))
	for i, t := range entity.ValidTypes {
		types[i] = string(t)
	}
	return types
}

// NewSearchTool creates a search tool using Foxy Contexts DI
func NewSearchTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "search",
			Description: ptr("Search Target Process entities by type with optional filters. " +
				"ALWAYS prefer the structured filter parameters (status, assignedUser, project, team, feature, priority, " +
				"dateFrom, dateTo) over the raw 'where' parameter — they automatically build correct TP API syntax. " +
				"Only use 'where' for advanced queries not covered by filters. Returns paginated results with cursor."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"type": {
						"type":        "string",
						"description": "Entity type to search (e.g., UserStory, Bug, Task, Feature)",
						"enum":        entityTypeStrings(),
					},
					"status": {
						"type":        "string",
						"description": "Filter by entity state name (string, e.g., 'Open', 'In Progress', 'Done'). Maps to EntityState.Name.",
					},
					"assignedUser": {
						"description": "Filter by assigned user — pass a string for email (e.g., 'john@company.com') or a number for user ID (e.g., 789). " +
							"String maps to AssignedUser.Email, number maps to AssignedUser.Id. " +
							"To filter by login name, use the 'where' param with: AssignedUser.Login eq 'loginname'",
					},
					"project": {
						"description": "Filter by project — pass a string for project name (e.g., 'My Project') or a number for project ID (e.g., 123). " +
							"String maps to Project.Name, number maps to Project.Id.",
					},
					"team": {
						"description": "Filter by team — pass a string for team name (e.g., 'Backend Team') or a number for team ID (e.g., 456). " +
							"String maps to Team.Name, number maps to Team.Id.",
					},
					"feature": {
						"description": "Filter by feature — pass a string for feature name or a number for feature ID. " +
							"String maps to Feature.Name, number maps to Feature.Id.",
					},
					"priority": {
						"type":        "string",
						"description": "Filter by priority name (string, e.g., 'High', 'Medium', 'Low', 'Urgent'). Maps to Priority.Name.",
					},
					"dateFrom": {
						"type":        "string",
						"description": "Filter by date range start (ISO 8601 format: YYYY-MM-DD)",
					},
					"dateTo": {
						"type":        "string",
						"description": "Filter by date range end (ISO 8601 format: YYYY-MM-DD)",
					},
					"dateField": {
						"type":        "string",
						"description": "Date field to use for date filtering (default: CreateDate)",
						"enum":        []interface{}{"CreateDate", "ModifyDate", "StartDate", "EndDate", "PlannedStartDate", "PlannedEndDate"},
					},
					"where": {
						"type":        "string",
						"description": "Raw TP API WHERE clause for advanced filtering (combined with structured filters using 'and'). " +
							"Use TP API syntax: 'eq' for equals, 'ne' for not equals, 'gt'/'lt'/'gte'/'lte' for comparisons. " +
							"Example: \"EntityState.Name eq 'Open'\". Do NOT use '==' or '!=' — those are invalid.",
					},
					"include": {
						"type":        "array",
						"description": "Additional fields to include in response (e.g., [Id,Name,Description])",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"take": {
						"type":        "integer",
						"description": "Number of items to return (default: 100, min: 1, max: 1000)",
						"minimum":     1,
						"maximum":     1000,
					},
					"orderByField": {
						"type":        "string",
						"description": "Field to sort results by (e.g., 'CreateDate', 'Name', 'Priority.Id'). Only one sort field is supported per request.",
					},
					"orderByDirection": {
						"type":        "string",
						"description": "Sort direction: 'asc' for ascending (default), 'desc' for descending.",
						"enum":        []interface{}{"asc", "desc"},
					},
					"cursor": {
						"type":        "string",
						"description": "Pagination cursor from previous response. When provided, all other filter params are ignored.",
					},
				},
				Required: []string{"type"},
			},
		},
		func(args map[string]interface{}) *mcp.CallToolResult {
			// Parse and validate entity type
			typeStr := getStringArg(args, "type")
			if typeStr == "" {
				return errorResult(fmt.Errorf("type parameter is required"))
			}

			entityType, err := entity.ParseType(typeStr)
			if err != nil {
				return errorResult(err)
			}

			// Check if cursor is provided
			cursor := getStringArg(args, "cursor")

			// Build search request
			req := query.SearchRequest{
				EntityType: entityType,
				Cursor:     cursor,
			}

			// If cursor is provided, ignore all other filter params
			if cursor == "" {
				// Parse take with default and clamping
				take := 100 // default
				if t, err := getIntArg(args, "take"); err == nil {
					take = t
				}
				// Clamp to [1, 1000]
				if take < 1 {
					take = 1
				}
				if take > 1000 {
					take = 1000
				}
				req.Take = take

				// Parse filters
				req.Filters = query.SearchFilters{
					Status:       getStringArg(args, "status"),
					AssignedUser: getAnyArg(args, "assignedUser"),
					Project:      getAnyArg(args, "project"),
					Team:         getAnyArg(args, "team"),
					Feature:      getAnyArg(args, "feature"),
					Priority:     getStringArg(args, "priority"),
					DateFrom:     getStringArg(args, "dateFrom"),
					DateTo:       getStringArg(args, "dateTo"),
					DateField:    getStringArg(args, "dateField"),
				}

				// Parse optional params
				req.RawWhere = getStringArg(args, "where")
				req.Include = getStringSliceArg(args, "include")
				orderByField := getStringArg(args, "orderByField")
				orderByDirection := getStringArg(args, "orderByDirection")
				if orderByDirection != "" && orderByField == "" {
					return errorResult(fmt.Errorf("orderByDirection requires orderByField to be set"))
				}
				if orderByField != "" {
					req.OrderByField = orderByField
					req.OrderByDesc = orderByDirection == "desc"
				}
			}

			// Call client
			resp, err := c.SearchEntities(context.Background(), req)
			if err != nil {
				return errorResult(err)
			}

			// Return JSON result
			return jsonResult(resp)
		},
	)
}
