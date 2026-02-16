package docs

import "fmt"

type DocTopic struct {
	Key         string
	Title       string
	Description string
	Content     string
}

var topicRegistry = map[string]DocTopic{
	"overview": {
		Key:         "overview",
		Title:       "Overview",
		Description: "Getting started guide",
		Content:     GettingStarted,
	},
	"tools": {
		Key:         "tools",
		Title:       "Tool Reference",
		Description: "Complete tool reference",
		Content:     ToolReference,
	},
	"search": {
		Key:         "search",
		Title:       "Search Guide",
		Description: "How to use the search tool",
		Content:     searchContent,
	},
	"entities": {
		Key:         "entities",
		Title:       "Entity Management",
		Description: "Get, create, and update entities",
		Content:     entityContent,
	},
	"comments": {
		Key:         "comments",
		Title:       "Comments",
		Description: "Adding and listing comments",
		Content:     commentContent,
	},
	"attachments": {
		Key:         "attachments",
		Title:       "Attachments",
		Description: "Listing and downloading attachments",
		Content:     attachmentContent,
	},
	"inspect": {
		Key:         "inspect",
		Title:       "API Inspection",
		Description: "Inspecting entity types and metadata",
		Content:     inspectContent,
	},
	"authentication": {
		Key:         "authentication",
		Title:       "Authentication",
		Description: "Setting up authentication",
		Content:     Authentication,
	},
	"pagination": {
		Key:         "pagination",
		Title:       "Pagination",
		Description: "Cursor-based pagination",
		Content:     paginationContent,
	},
	"query-syntax": {
		Key:         "query-syntax",
		Title:       "Query Syntax",
		Description: "WHERE clause syntax guide",
		Content:     QueryGuide,
	},
	"examples": {
		Key:         "examples",
		Title:       "Examples",
		Description: "Usage examples",
		Content:     Examples,
	},
}

func GetTopic(key string) (*DocTopic, error) {
	topic, ok := topicRegistry[key]
	if !ok {
		return nil, fmt.Errorf("topic %q not found", key)
	}
	return &topic, nil
}

func ListTopics() []DocTopic {
	topics := make([]DocTopic, 0, len(topicRegistry))
	for _, t := range topicRegistry {
		topics = append(topics, t)
	}
	return topics
}
