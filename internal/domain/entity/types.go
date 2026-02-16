package entity

import (
	"fmt"
	"strings"
)

// Type represents a Target Process entity type
type Type string

// Entity type constants
const (
	TypeUserStory      Type = "UserStory"
	TypeBug            Type = "Bug"
	TypeTask           Type = "Task"
	TypeFeature        Type = "Feature"
	TypeEpic           Type = "Epic"
	TypePortfolioEpic  Type = "PortfolioEpic"
	TypeSolution       Type = "Solution"
	TypeRequest        Type = "Request"
	TypeImpediment     Type = "Impediment"
	TypeTestCase       Type = "TestCase"
	TypeTestPlan       Type = "TestPlan"
	TypeProject        Type = "Project"
	TypeTeam           Type = "Team"
	TypeIteration      Type = "Iteration"
	TypeTeamIteration  Type = "TeamIteration"
	TypeRelease        Type = "Release"
	TypeProgram        Type = "Program"
)

// ValidTypes contains all valid entity types
var ValidTypes = []Type{
	TypeUserStory,
	TypeBug,
	TypeTask,
	TypeFeature,
	TypeEpic,
	TypePortfolioEpic,
	TypeSolution,
	TypeRequest,
	TypeImpediment,
	TypeTestCase,
	TypeTestPlan,
	TypeProject,
	TypeTeam,
	TypeIteration,
	TypeTeamIteration,
	TypeRelease,
	TypeProgram,
}

// IsValidType checks if the given string is a valid entity type (case-insensitive)
func IsValidType(t string) bool {
	lower := strings.ToLower(t)
	for _, vt := range ValidTypes {
		if strings.ToLower(string(vt)) == lower {
			return true
		}
	}
	return false
}

// ParseType converts a string to a Type with case-insensitive matching
func ParseType(t string) (Type, error) {
	lower := strings.ToLower(t)
	for _, vt := range ValidTypes {
		if strings.ToLower(string(vt)) == lower {
			return vt, nil
		}
	}
	return "", fmt.Errorf("invalid entity type: %s", t)
}

// Pluralize returns the plural form of the entity type
func Pluralize(t Type) string {
	return string(t) + "s"
}
