package tools

import (
	"context"
	"fmt"

	"tp-mcp-go/internal/client"
	"tp-mcp-go/internal/domain/entity"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// NewGetEntityTool creates a tool to get a single entity by ID
func NewGetEntityTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "get_entity",
			Description: ptr("Get a single Target Process entity by type and ID. " +
				"Returns the entity with all standard fields plus any additional fields requested via include."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"type": {
						"type":        "string",
						"description": "Entity type (e.g., UserStory, Bug, Task, Feature)",
						"enum":        entityTypeStrings(),
					},
					"id": {
						"type":        "integer",
						"description": "Entity ID",
					},
					"include": {
						"type":        "array",
						"description": "Additional fields to include in response (e.g., [Id,Name,Description,AssignedUser])",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
				Required: []string{"type", "id"},
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

			// Parse ID
			id, err := getIntArg(args, "id")
			if err != nil {
				return errorResult(err)
			}

			// Parse optional include
			include := getStringSliceArg(args, "include")

			// Call client
			result, err := c.GetEntity(context.Background(), entityType, id, include)
			if err != nil {
				return errorResult(err)
			}

			return jsonResult(result)
		},
	)
}

// NewCreateEntityTool creates a tool to create a new entity
func NewCreateEntityTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "create_entity",
			Description: ptr("Create a new Target Process entity. " +
				"Requires entity type and name. Optional fields include description, project, team, assignedUser, and customFields."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"type": {
						"type":        "string",
						"description": "Entity type to create (e.g., UserStory, Bug, Task, Feature)",
						"enum":        entityTypeStrings(),
					},
					"name": {
						"type":        "string",
						"description": "Entity name (required)",
					},
					"description": {
						"type":        "string",
						"description": "Entity description",
					},
					"project": {
						"type":        "object",
						"description": "Project reference with Id field (e.g., {\"Id\": 123})",
						"properties": map[string]interface{}{
							"Id": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"team": {
						"type":        "object",
						"description": "Team reference with Id field (e.g., {\"Id\": 456})",
						"properties": map[string]interface{}{
							"Id": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"assignedUser": {
						"type":        "object",
						"description": "Assigned user reference with Id field (e.g., {\"Id\": 789})",
						"properties": map[string]interface{}{
							"Id": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"customFields": {
						"type":        "object",
						"description": "Custom fields as key-value pairs to merge into the entity data",
					},
				},
				Required: []string{"type", "name"},
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

			// Build data map
			data := make(map[string]any)

			// Required: Name
			name := getStringArg(args, "name")
			if name == "" {
				return errorResult(fmt.Errorf("name parameter is required"))
			}
			data["Name"] = name

			// Optional: Description
			if description := getStringArg(args, "description"); description != "" {
				data["Description"] = description
			}

			// Optional: Project
			if project := getAnyArg(args, "project"); project != nil {
				data["Project"] = project
			}

			// Optional: Team
			if team := getAnyArg(args, "team"); team != nil {
				data["Team"] = team
			}

			// Optional: AssignedUser
			if assignedUser := getAnyArg(args, "assignedUser"); assignedUser != nil {
				data["AssignedUser"] = assignedUser
			}

			// Optional: CustomFields (merge into data map)
			if customFields := getAnyArg(args, "customFields"); customFields != nil {
				if cfMap, ok := customFields.(map[string]any); ok {
					for k, v := range cfMap {
						data[k] = v
					}
				}
			}

			// Call client
			result, err := c.CreateEntity(context.Background(), entityType, data)
			if err != nil {
				return errorResult(err)
			}

			return jsonResult(result)
		},
	)
}

// NewUpdateEntityTool creates a tool to update an existing entity
func NewUpdateEntityTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "update_entity",
			Description: ptr("Update an existing Target Process entity by type and ID. " +
				"Provide a fields object with key-value pairs to update."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"type": {
						"type":        "string",
						"description": "Entity type (e.g., UserStory, Bug, Task, Feature)",
						"enum":        entityTypeStrings(),
					},
					"id": {
						"type":        "integer",
						"description": "Entity ID to update",
					},
					"fields": {
						"type":        "object",
						"description": "Fields to update as key-value pairs (e.g., {\"Name\": \"New Name\", \"Description\": \"New Description\"})",
					},
				},
				Required: []string{"type", "id", "fields"},
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

			// Parse ID
			id, err := getIntArg(args, "id")
			if err != nil {
				return errorResult(err)
			}

			// Parse fields
			fields := getAnyArg(args, "fields")
			if fields == nil {
				return errorResult(fmt.Errorf("fields parameter is required"))
			}

			fieldsMap, ok := fields.(map[string]any)
			if !ok {
				return errorResult(fmt.Errorf("fields must be an object"))
			}

			// Call client
			result, err := c.UpdateEntity(context.Background(), entityType, id, fieldsMap)
			if err != nil {
				return errorResult(err)
			}

			return jsonResult(result)
		},
	)
}
