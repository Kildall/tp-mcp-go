package entity

import (
	"encoding/json"
	"testing"
)

func TestComment_JSONRoundTrip(t *testing.T) {
	owner := &User{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
	}
	general := &Ref{
		ID:           456,
		Name:         "US-789",
		ResourceType: "UserStory",
	}

	original := Comment{
		ID:          1,
		Description: "Test comment",
		CreateDate:  "2024-01-15T10:30:00",
		IsPrivate:   true,
		Owner:       owner,
		General:     general,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Comment: %v", err)
	}

	// Verify JSON tags produce correct field names
	jsonStr := string(jsonData)
	if !contains(jsonStr, `"Id":1`) {
		t.Error("Expected JSON to contain 'Id' field")
	}
	if !contains(jsonStr, `"Description":"Test comment"`) {
		t.Error("Expected JSON to contain 'Description' field")
	}
	if !contains(jsonStr, `"CreateDate":"2024-01-15T10:30:00"`) {
		t.Error("Expected JSON to contain 'CreateDate' field")
	}
	if !contains(jsonStr, `"IsPrivate":true`) {
		t.Error("Expected JSON to contain 'IsPrivate' field")
	}

	// Unmarshal back
	var unmarshaled Comment
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Comment: %v", err)
	}

	// Verify equality
	if unmarshaled.ID != original.ID {
		t.Errorf("ID mismatch: got %d, want %d", unmarshaled.ID, original.ID)
	}
	if unmarshaled.Description != original.Description {
		t.Errorf("Description mismatch: got %s, want %s", unmarshaled.Description, original.Description)
	}
	if unmarshaled.CreateDate != original.CreateDate {
		t.Errorf("CreateDate mismatch: got %s, want %s", unmarshaled.CreateDate, original.CreateDate)
	}
	if unmarshaled.IsPrivate != original.IsPrivate {
		t.Errorf("IsPrivate mismatch: got %v, want %v", unmarshaled.IsPrivate, original.IsPrivate)
	}
	if unmarshaled.Owner == nil || unmarshaled.Owner.ID != owner.ID {
		t.Errorf("Owner mismatch")
	}
	if unmarshaled.General == nil || unmarshaled.General.ID != general.ID {
		t.Errorf("General mismatch")
	}
}

func TestAttachment_OptionalFields(t *testing.T) {
	// Test with optional fields as nil
	attachment := Attachment{
		ID:             100,
		Name:           "document.pdf",
		Description:    nil,
		Date:           "2024-01-15T10:30:00",
		MimeType:       nil,
		Size:           1024,
		Uri:            "https://example.com/file.pdf",
		ThumbnailUri:   nil,
		UniqueFileName: "doc_12345.pdf",
		Owner:          nil,
		General:        nil,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(attachment)
	if err != nil {
		t.Fatalf("Failed to marshal Attachment: %v", err)
	}

	// Unmarshal back
	var unmarshaled Attachment
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Attachment: %v", err)
	}

	// Verify optional fields remain nil
	if unmarshaled.Description != nil {
		t.Error("Expected Description to be nil")
	}
	if unmarshaled.MimeType != nil {
		t.Error("Expected MimeType to be nil")
	}
	if unmarshaled.ThumbnailUri != nil {
		t.Error("Expected ThumbnailUri to be nil")
	}
	if unmarshaled.Owner != nil {
		t.Error("Expected Owner to be nil")
	}
	if unmarshaled.General != nil {
		t.Error("Expected General to be nil")
	}

	// Verify required fields
	if unmarshaled.ID != attachment.ID {
		t.Errorf("ID mismatch: got %d, want %d", unmarshaled.ID, attachment.ID)
	}
	if unmarshaled.Name != attachment.Name {
		t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, attachment.Name)
	}
	if unmarshaled.Size != attachment.Size {
		t.Errorf("Size mismatch: got %d, want %d", unmarshaled.Size, attachment.Size)
	}
}

func TestAPIResponse_TPFormat(t *testing.T) {
	// Sample TargetProcess API JSON response
	sampleJSON := `{
		"Items": [
			{
				"Id": 789,
				"Name": "User Story 1",
				"State": "Open"
			},
			{
				"Id": 790,
				"Name": "User Story 2",
				"State": "InProgress"
			}
		],
		"Next": "https://example.targetprocess.com/api/v1/UserStories?skip=100&take=100"
	}`

	var response APIResponse
	if err := json.Unmarshal([]byte(sampleJSON), &response); err != nil {
		t.Fatalf("Failed to unmarshal APIResponse: %v", err)
	}

	// Verify Items
	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}

	// Verify first item
	if response.Items[0]["Id"].(float64) != 789 {
		t.Errorf("Expected first item Id to be 789")
	}
	if response.Items[0]["Name"].(string) != "User Story 1" {
		t.Errorf("Expected first item Name to be 'User Story 1'")
	}

	// Verify Next URL
	expectedNext := "https://example.targetprocess.com/api/v1/UserStories?skip=100&take=100"
	if response.Next != expectedNext {
		t.Errorf("Next URL mismatch: got %s, want %s", response.Next, expectedNext)
	}

	// Test JSON field names (verify tags work correctly)
	marshaledJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal APIResponse: %v", err)
	}

	jsonStr := string(marshaledJSON)
	if !contains(jsonStr, `"Items"`) {
		t.Error("Expected JSON to contain 'Items' field")
	}
	if !contains(jsonStr, `"Next"`) {
		t.Error("Expected JSON to contain 'Next' field")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
