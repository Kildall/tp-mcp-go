package entity

// Comment represents a comment in TargetProcess
type Comment struct {
	ID          int     `json:"Id"`
	Description string  `json:"Description"`
	CreateDate  string  `json:"CreateDate"`
	IsPrivate   bool    `json:"IsPrivate"`
	Owner       *User   `json:"Owner,omitempty"`
	General     *Ref    `json:"General,omitempty"`
}

// Attachment represents a file attachment in TargetProcess
type Attachment struct {
	ID             int     `json:"Id"`
	Name           string  `json:"Name"`
	Description    *string `json:"Description"`
	Date           string  `json:"Date"`
	MimeType       *string `json:"MimeType"`
	Size           int     `json:"Size"`
	Uri            string  `json:"Uri"`
	ThumbnailUri   *string `json:"ThumbnailUri"`
	UniqueFileName string  `json:"UniqueFileName"`
	Owner          *User   `json:"Owner,omitempty"`
	General        *Ref    `json:"General,omitempty"`
}

// User represents a user in TargetProcess
type User struct {
	ID        int    `json:"Id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

// Ref represents a reference to another entity in TargetProcess
type Ref struct {
	ID           int    `json:"Id"`
	Name         string `json:"Name"`
	ResourceType string `json:"ResourceType"`
}

// APIResponse represents a generic TargetProcess API response
type APIResponse struct {
	Items []map[string]any `json:"Items,omitempty"`
	Next  string           `json:"Next,omitempty"`
}
