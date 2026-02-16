package query

import (
	"fmt"
	"strings"
)

// EscapeValue doubles single quotes for safe use in WHERE clauses
func EscapeValue(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

// FormatStringCondition builds a WHERE condition for a string value
func FormatStringCondition(field, op, value string) string {
	return fmt.Sprintf("%s %s '%s'", field, op, EscapeValue(value))
}

// FormatNumberCondition builds a WHERE condition for a numeric value
func FormatNumberCondition(field, op string, value int) string {
	return fmt.Sprintf("%s %s %d", field, op, value)
}
