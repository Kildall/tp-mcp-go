package query

import "testing"

func TestBuildWhereClause_EmptyFilters(t *testing.T) {
	filters := SearchFilters{}
	result := BuildWhereClause(filters, "")

	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestBuildWhereClause_SingleFilter(t *testing.T) {
	filters := SearchFilters{
		Status: "Open",
	}
	result := BuildWhereClause(filters, "")

	expected := "EntityState.Name eq 'Open'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_MultipleFilters(t *testing.T) {
	filters := SearchFilters{
		Status:       "Open",
		AssignedUser: "user@example.com",
		Priority:     "High",
	}
	result := BuildWhereClause(filters, "")

	expected := "EntityState.Name eq 'Open' and AssignedUser.Email eq 'user@example.com' and Priority.Name eq 'High'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_DateRangeWithCustomField(t *testing.T) {
	filters := SearchFilters{
		DateFrom:  "2024-01-01",
		DateTo:    "2024-12-31",
		DateField: "ModifyDate",
	}
	result := BuildWhereClause(filters, "")

	expected := "ModifyDate gte '2024-01-01' and ModifyDate lte '2024-12-31'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_DateRangeDefaultField(t *testing.T) {
	filters := SearchFilters{
		DateFrom: "2024-01-01",
		DateTo:   "2024-12-31",
	}
	result := BuildWhereClause(filters, "")

	expected := "CreateDate gte '2024-01-01' and CreateDate lte '2024-12-31'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_ProjectByName(t *testing.T) {
	filters := SearchFilters{
		Project: "MyProject",
	}
	result := BuildWhereClause(filters, "")

	expected := "Project.Name eq 'MyProject'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_ProjectByID_Int(t *testing.T) {
	filters := SearchFilters{
		Project: 123,
	}
	result := BuildWhereClause(filters, "")

	expected := "Project.Id eq 123"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_ProjectByID_Float64(t *testing.T) {
	filters := SearchFilters{
		Project: float64(456),
	}
	result := BuildWhereClause(filters, "")

	expected := "Project.Id eq 456"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_TeamByName(t *testing.T) {
	filters := SearchFilters{
		Team: "DevTeam",
	}
	result := BuildWhereClause(filters, "")

	expected := "Team.Name eq 'DevTeam'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_TeamByID(t *testing.T) {
	filters := SearchFilters{
		Team: 789,
	}
	result := BuildWhereClause(filters, "")

	expected := "Team.Id eq 789"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_TeamByID_Float64(t *testing.T) {
	filters := SearchFilters{
		Team: float64(42),
	}
	result := BuildWhereClause(filters, "")

	expected := "Team.Id eq 42"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_FeatureByName(t *testing.T) {
	filters := SearchFilters{
		Feature: "NewFeature",
	}
	result := BuildWhereClause(filters, "")

	expected := "Feature.Name eq 'NewFeature'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_FeatureByID(t *testing.T) {
	filters := SearchFilters{
		Feature: 999,
	}
	result := BuildWhereClause(filters, "")

	expected := "Feature.Id eq 999"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_FeatureByID_Float64(t *testing.T) {
	filters := SearchFilters{
		Feature: float64(99),
	}
	result := BuildWhereClause(filters, "")

	expected := "Feature.Id eq 99"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_RawWhereOnly(t *testing.T) {
	filters := SearchFilters{}
	result := BuildWhereClause(filters, "CustomField eq 'Value'")

	expected := "CustomField eq 'Value'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_FiltersAndRawWhere(t *testing.T) {
	filters := SearchFilters{
		Status: "Open",
	}
	result := BuildWhereClause(filters, "CustomField eq 'Value'")

	expected := "EntityState.Name eq 'Open' and CustomField eq 'Value'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_SingleQuoteEscaping(t *testing.T) {
	filters := SearchFilters{
		Status:   "O'Reilly",
		Priority: "Can't Wait",
	}
	result := BuildWhereClause(filters, "")

	expected := "EntityState.Name eq 'O''Reilly' and Priority.Name eq 'Can''t Wait'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_AssignedUserByEmail(t *testing.T) {
	filters := SearchFilters{
		AssignedUser: "user@example.com",
	}
	result := BuildWhereClause(filters, "")

	expected := "AssignedUser.Email eq 'user@example.com'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_AssignedUserByID_Int(t *testing.T) {
	filters := SearchFilters{
		AssignedUser: 789,
	}
	result := BuildWhereClause(filters, "")

	expected := "AssignedUser.Id eq 789"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_AssignedUserByID_Float64(t *testing.T) {
	filters := SearchFilters{
		AssignedUser: float64(42),
	}
	result := BuildWhereClause(filters, "")

	expected := "AssignedUser.Id eq 42"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBuildWhereClause_AssignedUserEmptyString(t *testing.T) {
	filters := SearchFilters{
		AssignedUser: "",
	}
	result := BuildWhereClause(filters, "")

	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestBuildWhereClause_ComplexCombination(t *testing.T) {
	filters := SearchFilters{
		Status:       "In Progress",
		AssignedUser: "dev@example.com",
		Project:      123,
		Team:         "Backend",
		Feature:      float64(456),
		Priority:     "Critical",
		DateFrom:     "2024-01-01",
		DateTo:       "2024-12-31",
		DateField:    "ModifyDate",
	}
	result := BuildWhereClause(filters, "CustomField eq 'test'")

	expected := "EntityState.Name eq 'In Progress' and AssignedUser.Email eq 'dev@example.com' and Project.Id eq 123 and Team.Name eq 'Backend' and Feature.Id eq 456 and Priority.Name eq 'Critical' and ModifyDate gte '2024-01-01' and ModifyDate lte '2024-12-31' and CustomField eq 'test'"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
