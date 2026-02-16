package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/testutil"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

func TestGetEntityReturnsResult(t *testing.T) {
	expectedResult := map[string]any{
		"Id":          123,
		"Name":        "Test Story",
		"Description": "Test description",
	}

	mock := &testutil.MockClient{
		GetEntityFn: func(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error) {
			if entityType != entity.TypeUserStory {
				t.Errorf("expected entityType UserStory, got %v", entityType)
			}
			if id != 123 {
				t.Errorf("expected id 123, got %d", id)
			}
			if len(include) != 2 || include[0] != "Id" || include[1] != "Name" {
				t.Errorf("expected include [Id, Name], got %v", include)
			}
			return expectedResult, nil
		},
	}

	tool := NewGetEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type":    "UserStory",
		"id":      float64(123),
		"include": []interface{}{"Id", "Name"},
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify JSON result
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var resultData map[string]any
	if err := json.Unmarshal([]byte(textContent.Text), &resultData); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if resultData["Id"] != float64(123) {
		t.Errorf("expected Id 123, got %v", resultData["Id"])
	}

	if resultData["Name"] != "Test Story" {
		t.Errorf("expected Name 'Test Story', got %v", resultData["Name"])
	}
}

func TestGetEntityInvalidType(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewGetEntityTool(mock)

	result := tool.Callback(map[string]interface{}{
		"type": "InvalidType",
		"id":   float64(123),
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for invalid entity type")
	}
}

func TestGetEntityMissingId(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewGetEntityTool(mock)

	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for missing id")
	}
}

func TestGetEntityHandlesError(t *testing.T) {
	mock := &testutil.MockClient{
		GetEntityFn: func(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	tool := NewGetEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
		"id":   float64(123),
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error")
	}
}

func TestCreateEntityBuildsDataMap(t *testing.T) {
	var capturedData map[string]any
	var capturedType entity.Type

	mock := &testutil.MockClient{
		CreateEntityFn: func(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error) {
			capturedType = entityType
			capturedData = data
			return map[string]any{"Id": 456}, nil
		},
	}

	tool := NewCreateEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type":        "UserStory",
		"name":        "New Story",
		"description": "Story description",
		"project":     map[string]interface{}{"Id": float64(10)},
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify entity type
	if capturedType != entity.TypeUserStory {
		t.Errorf("expected UserStory, got %v", capturedType)
	}

	// Verify data map
	if capturedData["Name"] != "New Story" {
		t.Errorf("expected Name 'New Story', got %v", capturedData["Name"])
	}

	if capturedData["Description"] != "Story description" {
		t.Errorf("expected Description 'Story description', got %v", capturedData["Description"])
	}

	project, ok := capturedData["Project"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Project to be map, got %T", capturedData["Project"])
	}

	if project["Id"] != float64(10) {
		t.Errorf("expected Project.Id 10, got %v", project["Id"])
	}
}

func TestCreateEntityWithCustomFields(t *testing.T) {
	var capturedData map[string]any

	mock := &testutil.MockClient{
		CreateEntityFn: func(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error) {
			capturedData = data
			return map[string]any{"Id": 789}, nil
		},
	}

	tool := NewCreateEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "Bug",
		"name": "Test Bug",
		"customFields": map[string]interface{}{
			"Severity": "High",
			"Priority": "P1",
		},
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify custom fields are merged
	if capturedData["Severity"] != "High" {
		t.Errorf("expected Severity 'High', got %v", capturedData["Severity"])
	}

	if capturedData["Priority"] != "P1" {
		t.Errorf("expected Priority 'P1', got %v", capturedData["Priority"])
	}

	// Name should still be present
	if capturedData["Name"] != "Test Bug" {
		t.Errorf("expected Name 'Test Bug', got %v", capturedData["Name"])
	}
}

func TestCreateEntityMissingName(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewCreateEntityTool(mock)

	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for missing name")
	}
}

func TestCreateEntityWithTeamAndAssignedUser(t *testing.T) {
	var capturedData map[string]any

	mock := &testutil.MockClient{
		CreateEntityFn: func(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error) {
			capturedData = data
			return map[string]any{"Id": 100}, nil
		},
	}

	tool := NewCreateEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type":         "Task",
		"name":         "Test Task",
		"team":         map[string]interface{}{"Id": float64(5)},
		"assignedUser": map[string]interface{}{"Id": float64(20)},
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify team
	team, ok := capturedData["Team"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Team to be map, got %T", capturedData["Team"])
	}
	if team["Id"] != float64(5) {
		t.Errorf("expected Team.Id 5, got %v", team["Id"])
	}

	// Verify assignedUser
	assignedUser, ok := capturedData["AssignedUser"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected AssignedUser to be map, got %T", capturedData["AssignedUser"])
	}
	if assignedUser["Id"] != float64(20) {
		t.Errorf("expected AssignedUser.Id 20, got %v", assignedUser["Id"])
	}
}

func TestUpdateEntityPassesFieldsCorrectly(t *testing.T) {
	var capturedType entity.Type
	var capturedId int
	var capturedData map[string]any

	mock := &testutil.MockClient{
		UpdateEntityFn: func(ctx context.Context, entityType entity.Type, id int, data map[string]any) (map[string]any, error) {
			capturedType = entityType
			capturedId = id
			capturedData = data
			return map[string]any{"Id": id, "Name": "Updated"}, nil
		},
	}

	tool := NewUpdateEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "Bug",
		"id":   float64(999),
		"fields": map[string]interface{}{
			"Name":        "Updated Bug",
			"Description": "Updated description",
			"EntityState": map[string]interface{}{"Id": float64(3)},
		},
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify entity type
	if capturedType != entity.TypeBug {
		t.Errorf("expected Bug, got %v", capturedType)
	}

	// Verify ID
	if capturedId != 999 {
		t.Errorf("expected id 999, got %d", capturedId)
	}

	// Verify fields
	if capturedData["Name"] != "Updated Bug" {
		t.Errorf("expected Name 'Updated Bug', got %v", capturedData["Name"])
	}

	if capturedData["Description"] != "Updated description" {
		t.Errorf("expected Description 'Updated description', got %v", capturedData["Description"])
	}

	entityState, ok := capturedData["EntityState"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected EntityState to be map, got %T", capturedData["EntityState"])
	}
	if entityState["Id"] != float64(3) {
		t.Errorf("expected EntityState.Id 3, got %v", entityState["Id"])
	}
}

func TestUpdateEntityMissingFields(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewUpdateEntityTool(mock)

	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
		"id":   float64(123),
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for missing fields")
	}
}

func TestUpdateEntityInvalidFieldsType(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewUpdateEntityTool(mock)

	result := tool.Callback(map[string]interface{}{
		"type":   "UserStory",
		"id":     float64(123),
		"fields": "not an object",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for invalid fields type")
	}
}

func TestUpdateEntityHandlesError(t *testing.T) {
	mock := &testutil.MockClient{
		UpdateEntityFn: func(ctx context.Context, entityType entity.Type, id int, data map[string]any) (map[string]any, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	tool := NewUpdateEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type":   "UserStory",
		"id":     float64(123),
		"fields": map[string]interface{}{"Name": "New Name"},
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error")
	}
}
