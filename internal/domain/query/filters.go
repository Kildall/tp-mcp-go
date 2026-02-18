package query

import "tp-mcp-go/internal/domain/entity"

type SearchFilters struct {
	Status       string
	AssignedUser string
	Project      any    // string (name) or int (ID)
	Team         any    // string (name) or int (ID)
	Feature      any    // string (name) or int (ID)
	Priority     string
	DateFrom     string
	DateTo       string
	DateField    string // CreateDate, ModifyDate, StartDate, EndDate, PlannedStartDate, PlannedEndDate
}

type SearchRequest struct {
	EntityType entity.Type
	Filters    SearchFilters
	RawWhere   string
	Include    []string
	Take         int
	OrderByField string
	OrderByDesc  bool
	Cursor       string
}

type PaginatedResponse struct {
	Items      []map[string]any `json:"items"`
	Pagination PaginationMeta   `json:"pagination"`
}

type PaginationMeta struct {
	HasMore  bool   `json:"hasMore"`
	Cursor   string `json:"cursor,omitempty"`
	Returned int    `json:"returned"`
}
