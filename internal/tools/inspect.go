package tools

import (
	"context"
	"fmt"

	"tp-mcp-go/internal/client"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// extractProperties extracts properties for a given entity type from metadata
func extractProperties(metadata any, entityType string) (any, error) {
	metaMap, ok := metadata.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected metadata format")
	}
	typeData, ok := metaMap[entityType]
	if !ok {
		return nil, fmt.Errorf("entity type %s not found in metadata", entityType)
	}
	return typeData, nil
}

// extractPropertyDetails extracts specific property details from metadata
func extractPropertyDetails(metadata any, entityType, property string) (any, error) {
	typeData, err := extractProperties(metadata, entityType)
	if err != nil {
		return nil, err
	}
	typeMap, ok := typeData.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected type data format")
	}
	propData, ok := typeMap[property]
	if !ok {
		return nil, fmt.Errorf("property %s not found for entity type %s", property, entityType)
	}
	return propData, nil
}

// NewInspectObjectTool creates an inspect_object tool using Foxy Contexts DI
func NewInspectObjectTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "inspect_object",
			Description: ptr("Inspect Target Process metadata, entity types, and properties. " +
				"Use this tool to discover available entity types, explore properties for a specific type, " +
				"get detailed information about a specific property, or examine the full API structure."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"action": {
						"type":        "string",
						"description": "Action to perform",
						"enum":        []interface{}{"list_types", "get_properties", "get_property_details", "discover_api_structure"},
					},
					"entityType": {
						"type":        "string",
						"description": "Entity type name (required for get_properties and get_property_details)",
					},
					"property": {
						"type":        "string",
						"description": "Property name (required for get_property_details)",
					},
				},
				Required: []string{"action"},
			},
		},
		func(args map[string]interface{}) *mcp.CallToolResult {
			action := getStringArg(args, "action")
			if action == "" {
				return errorResult(fmt.Errorf("action parameter is required"))
			}

			ctx := context.Background()

			switch action {
			case "list_types":
				types, err := c.GetValidEntityTypes(ctx)
				if err != nil {
					return errorResult(err)
				}
				return jsonResult(types)

			case "get_properties":
				entityType := getStringArg(args, "entityType")
				if entityType == "" {
					return errorResult(fmt.Errorf("entityType parameter is required for get_properties action"))
				}

				metadata, err := c.FetchMetadata(ctx)
				if err != nil {
					return errorResult(err)
				}

				properties, err := extractProperties(metadata, entityType)
				if err != nil {
					return errorResult(err)
				}

				return jsonResult(properties)

			case "get_property_details":
				entityType := getStringArg(args, "entityType")
				if entityType == "" {
					return errorResult(fmt.Errorf("entityType parameter is required for get_property_details action"))
				}

				property := getStringArg(args, "property")
				if property == "" {
					return errorResult(fmt.Errorf("property parameter is required for get_property_details action"))
				}

				metadata, err := c.FetchMetadata(ctx)
				if err != nil {
					return errorResult(err)
				}

				propDetails, err := extractPropertyDetails(metadata, entityType, property)
				if err != nil {
					return errorResult(err)
				}

				return jsonResult(propDetails)

			case "discover_api_structure":
				metadata, err := c.FetchMetadata(ctx)
				if err != nil {
					return errorResult(err)
				}
				return jsonResult(metadata)

			default:
				return errorResult(fmt.Errorf("unknown action: %s", action))
			}
		},
	)
}
