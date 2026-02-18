package query

import "strings"

// BuildWhereClause constructs a WHERE clause from SearchFilters and optional raw WHERE clause
func BuildWhereClause(filters SearchFilters, rawWhere string) string {
	var conditions []string

	// Status
	if filters.Status != "" {
		conditions = append(conditions, FormatStringCondition("EntityState.Name", "eq", filters.Status))
	}

	// AssignedUser - string (email) or number (ID)
	if filters.AssignedUser != nil {
		switch v := filters.AssignedUser.(type) {
		case string:
			if v != "" {
				conditions = append(conditions, FormatStringCondition("AssignedUser.Email", "eq", v))
			}
		case int:
			conditions = append(conditions, FormatNumberCondition("AssignedUser.Id", "eq", v))
		case float64:
			conditions = append(conditions, FormatNumberCondition("AssignedUser.Id", "eq", int(v)))
		}
	}

	// Project - string (name) or number (ID)
	if filters.Project != nil {
		switch v := filters.Project.(type) {
		case string:
			conditions = append(conditions, FormatStringCondition("Project.Name", "eq", v))
		case int:
			conditions = append(conditions, FormatNumberCondition("Project.Id", "eq", v))
		case float64:
			// JSON unmarshaling produces float64 for numbers
			conditions = append(conditions, FormatNumberCondition("Project.Id", "eq", int(v)))
		}
	}

	// Team - string (name) or number (ID)
	if filters.Team != nil {
		switch v := filters.Team.(type) {
		case string:
			conditions = append(conditions, FormatStringCondition("Team.Name", "eq", v))
		case int:
			conditions = append(conditions, FormatNumberCondition("Team.Id", "eq", v))
		case float64:
			conditions = append(conditions, FormatNumberCondition("Team.Id", "eq", int(v)))
		}
	}

	// Feature - string (name) or number (ID)
	if filters.Feature != nil {
		switch v := filters.Feature.(type) {
		case string:
			conditions = append(conditions, FormatStringCondition("Feature.Name", "eq", v))
		case int:
			conditions = append(conditions, FormatNumberCondition("Feature.Id", "eq", v))
		case float64:
			conditions = append(conditions, FormatNumberCondition("Feature.Id", "eq", int(v)))
		}
	}

	// Priority
	if filters.Priority != "" {
		conditions = append(conditions, FormatStringCondition("Priority.Name", "eq", filters.Priority))
	}

	// Date range - use DateField (defaults to CreateDate)
	dateField := filters.DateField
	if dateField == "" {
		dateField = "CreateDate"
	}

	if filters.DateFrom != "" {
		conditions = append(conditions, FormatStringCondition(dateField, "gte", filters.DateFrom))
	}

	if filters.DateTo != "" {
		conditions = append(conditions, FormatStringCondition(dateField, "lte", filters.DateTo))
	}

	// Raw WHERE clause
	if rawWhere != "" {
		conditions = append(conditions, rawWhere)
	}

	if len(conditions) == 0 {
		return ""
	}

	return strings.Join(conditions, " and ")
}
