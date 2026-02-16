package entity

import "testing"

func TestIsValidType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid UserStory", "UserStory", true},
		{"Valid Bug", "Bug", true},
		{"Valid Task", "Task", true},
		{"Valid Feature", "Feature", true},
		{"Valid Epic", "Epic", true},
		{"Valid PortfolioEpic", "PortfolioEpic", true},
		{"Valid Solution", "Solution", true},
		{"Valid Request", "Request", true},
		{"Valid Impediment", "Impediment", true},
		{"Valid TestCase", "TestCase", true},
		{"Valid TestPlan", "TestPlan", true},
		{"Valid Project", "Project", true},
		{"Valid Team", "Team", true},
		{"Valid Iteration", "Iteration", true},
		{"Valid TeamIteration", "TeamIteration", true},
		{"Valid Release", "Release", true},
		{"Valid Program", "Program", true},
		{"Case insensitive userstory", "userstory", true},
		{"Case insensitive USERSTORY", "USERSTORY", true},
		{"Case insensitive MixedCase", "UsErStOrY", true},
		{"Invalid type", "InvalidType", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidType(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidType(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Type
		shouldError bool
	}{
		{"Valid UserStory", "UserStory", TypeUserStory, false},
		{"Valid Bug", "Bug", TypeBug, false},
		{"Valid Task", "Task", TypeTask, false},
		{"Case insensitive userstory", "userstory", TypeUserStory, false},
		{"Case insensitive BUG", "BUG", TypeBug, false},
		{"Case insensitive tAsK", "tAsK", TypeTask, false},
		{"Case insensitive portfolioepic", "portfolioepic", TypePortfolioEpic, false},
		{"Case insensitive TEAMITERATION", "TEAMITERATION", TypeTeamIteration, false},
		{"Invalid type", "InvalidType", "", true},
		{"Empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseType(tt.input)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ParseType(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseType(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseType(%q) = %v; want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		name     string
		input    Type
		expected string
	}{
		{"UserStory plural", TypeUserStory, "UserStorys"},
		{"Bug plural", TypeBug, "Bugs"},
		{"Task plural", TypeTask, "Tasks"},
		{"Feature plural", TypeFeature, "Features"},
		{"Epic plural", TypeEpic, "Epics"},
		{"PortfolioEpic plural", TypePortfolioEpic, "PortfolioEpics"},
		{"TestCase plural", TypeTestCase, "TestCases"},
		{"TeamIteration plural", TypeTeamIteration, "TeamIterations"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pluralize(tt.input)
			if result != tt.expected {
				t.Errorf("Pluralize(%v) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
