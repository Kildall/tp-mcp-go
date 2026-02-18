package docs

const GettingStarted = `# TP MCP Server - Getting Started

## Quick Start

1. Set environment variables:
   - TP_DOMAIN: Your Target Process domain (e.g., company.tpondemand.com)
   - TP_ACCESS_TOKEN: Your API access token

2. Run the server:
   TP_DOMAIN=company.tpondemand.com TP_ACCESS_TOKEN=your-token ./server

3. The server communicates via stdio using the MCP protocol.

## Available Tools

| Tool | Description |
|------|-------------|
| search | Search for entities across 17 entity types with filters |
| get_entity | Retrieve a single entity by type and ID |
| create_entity | Create a new entity |
| update_entity | Update an existing entity |
| add_comment | Add a comment to an entity |
| list_comments | List comments on an entity |
| list_attachments | List attachments on an entity |
| download_attachment | Download an attachment by ID |
| inspect_object | Inspect entity types and API metadata |
| get_documentation | Access this documentation |

## Available Documentation Topics

Use the get_documentation tool with one of these topics:
- overview: Getting started guide
- tools: Complete tool reference
- search: How to use the search tool
- entities: Get, create, and update entities
- comments: Adding and listing comments
- attachments: Listing and downloading attachments
- inspect: Inspecting entity types and metadata
- authentication: Setting up authentication
- pagination: Cursor-based pagination
- query-syntax: WHERE clause syntax guide
- examples: Usage examples

## Next Steps

- Review the authentication guide to set up your credentials
- Try the examples to see common use cases
- Explore the search tool to query your Target Process data
`

const ToolReference = `# Tool Reference

## search

Search for entities across 17 entity types with flexible filtering.

**Parameters:**
- entity_type (required): string - Type of entity to search (e.g., "UserStory", "Bug", "Task")
- where (optional): string - WHERE clause for filtering (e.g., "EntityState.Name eq 'Open'")
- take (optional): integer - Number of results to return (default: 25, max: 1000)
- skip (optional): integer - Number of results to skip (default: 0)
- orderByField (optional): string - Field name to sort by (e.g., "CreateDate", "Name", "Priority.Id"). Only single-field sorting is supported.
- orderByDirection (optional): enum - Sort direction: "asc" or "desc" (defaults to "asc")
- assignedUser (optional): string or number - Filter by assigned user; pass an email string (maps to AssignedUser.Email) or a numeric user ID (maps to AssignedUser.Id)
- project (optional): string or number - Filter by project name (string) or project ID (number)
- team (optional): string or number - Filter by team name (string) or team ID (number)
- feature (optional): string or number - Filter by feature name (string) or feature ID (number)
- include (optional): array - Related entities to include (e.g., ["AssignedUser", "EntityState"])

**Example:**
search(entity_type="UserStory", where="EntityState.Name eq 'Open'", take=10)

## get_entity

Retrieve a single entity by type and ID.

**Parameters:**
- entity_type (required): string - Type of entity (e.g., "UserStory", "Bug")
- id (required): integer - Entity ID
- include (optional): array - Related entities to include

**Example:**
get_entity(entity_type="UserStory", id=1234, include=["AssignedUser", "EntityState"])

## create_entity

Create a new entity in Target Process.

**Parameters:**
- entity_type (required): string - Type of entity to create
- data (required): object - Entity data (must include required fields)

**Example:**
create_entity(entity_type="Bug", data={"Name": "Login fails", "Project": {"Id": 100}})

## update_entity

Update an existing entity.

**Parameters:**
- entity_type (required): string - Type of entity to update
- id (required): integer - Entity ID
- data (required): object - Updated fields

**Example:**
update_entity(entity_type="UserStory", id=1234, data={"EntityState": {"Id": 5}})

## add_comment

Add a comment to an entity.

**Parameters:**
- entity_type (required): string - Type of entity
- entity_id (required): integer - Entity ID
- description (required): string - Comment text

**Example:**
add_comment(entity_type="Bug", entity_id=1234, description="Fixed in PR #456")

## list_comments

List all comments on an entity.

**Parameters:**
- entity_type (required): string - Type of entity
- entity_id (required): integer - Entity ID
- take (optional): integer - Number of comments to return (default: 25)
- skip (optional): integer - Number of comments to skip (default: 0)

**Example:**
list_comments(entity_type="UserStory", entity_id=1234, take=10)

## list_attachments

List all attachments on an entity.

**Parameters:**
- entity_type (required): string - Type of entity
- entity_id (required): integer - Entity ID
- take (optional): integer - Number of attachments to return (default: 25)
- skip (optional): integer - Number of attachments to skip (default: 0)

**Example:**
list_attachments(entity_type="Bug", entity_id=1234)

## download_attachment

Download an attachment by ID.

**Parameters:**
- attachment_id (required): integer - Attachment ID
- output_path (optional): string - Path to save file (default: current directory with original filename)

**Example:**
download_attachment(attachment_id=5678, output_path="./downloads/screenshot.png")

## inspect_object

Inspect entity types and API metadata.

**Parameters:**
- object_type (optional): string - Entity type to inspect (if omitted, lists all types)

**Example:**
inspect_object(object_type="UserStory")

## get_documentation

Access this documentation system.

**Parameters:**
- topic (optional): string - Documentation topic (if omitted, shows overview)

**Example:**
get_documentation(topic="search")
`

const Examples = `# Usage Examples

## Search for Open User Stories

Find all open user stories assigned to a specific user:

search(
  entity_type="UserStory",
  where="EntityState.Name eq 'Open' and AssignedUser.Id eq 123",
  take=50,
  include=["AssignedUser", "EntityState", "Project"]
)

## Create a Bug

Create a new bug report:

create_entity(
  entity_type="Bug",
  data={
    "Name": "Login button not working",
    "Description": "Users cannot log in when clicking the login button",
    "Project": {"Id": 100},
    "Priority": {"Id": 2},
    "Severity": {"Id": 3}
  }
)

## Update Story Status

Move a user story to "In Progress":

update_entity(
  entity_type="UserStory",
  id=1234,
  data={"EntityState": {"Id": 5}}
)

## Add a Comment

Add a status update comment:

add_comment(
  entity_type="Bug",
  entity_id=1234,
  description="Fixed in commit abc123. Deployed to staging for testing."
)

## Search by Date Range

Find bugs created in the last week:

search(
  entity_type="Bug",
  where="CreateDate gte '2026-02-08'",
  take=100
)

## Get Entity with Full Context

Retrieve a user story with all related data:

get_entity(
  entity_type="UserStory",
  id=1234,
  include=["AssignedUser", "EntityState", "Project", "Team", "Feature", "Comments"]
)

## List Recent Comments

Get the 10 most recent comments on a bug:

list_comments(
  entity_type="Bug",
  entity_id=5678,
  take=10
)

## Download Attachment

Download a screenshot attached to a bug:

download_attachment(
  attachment_id=9012,
  output_path="./bug-screenshots/screenshot.png"
)

## Inspect Entity Schema

Learn what fields are available on UserStory entities:

inspect_object(object_type="UserStory")

## Search with Sorting

Find the 25 most recently created tasks:

search(
  entity_type="Task",
  orderByField="CreateDate",
  orderByDirection="desc",
  take=25,
  include=["AssignedUser", "EntityState"]
)

## Complex Filter Query

Find high-priority bugs in a specific project that are not assigned:

search(
  entity_type="Bug",
  where="Project.Id eq 100 and Priority.Name eq 'High' and AssignedUser.Id is null",
  take=50
)

## Search with Multiple Status Values

Find stories in any active status:

search(
  entity_type="UserStory",
  where="EntityState.Name in ('Open','In Progress','Code Review') and Project.Id eq 100",
  take=50,
  include=["AssignedUser", "EntityState"]
)

## Update Multiple Fields

Update both status and assignee:

update_entity(
  entity_type="Task",
  id=7890,
  data={
    "EntityState": {"Id": 6},
    "AssignedUser": {"Id": 456}
  }
)
`

const QueryGuide = `# Query Syntax Guide

## WHERE Clause Basics

The search tool supports powerful filtering using WHERE clauses with OData syntax.

### Basic Operators

- eq: Equals
- ne: Not equals
- gt: Greater than
- ge: Greater than or equal
- lt: Less than
- le: Less than or equal
- in: Membership test (works with numbers and strings)
- contains: Substring match (case-insensitive)
- not contains: Excludes substring
- is null: Field is empty/unset
- is not null: Field has a value

### Examples

EntityState.Name eq 'Open'
Priority.Id eq 1
CreateDate gt '2026-01-01'
Effort ge 5

## Combining Filters

Use 'and' to combine multiple conditions:

EntityState.Name eq 'Open' and AssignedUser.Id eq 123
Priority.Name in ('High','Critical')
Project.Id eq 100 and EntityState.Name in ('Open','In Progress')

## Null Checks

Check for null or non-null values:

AssignedUser.Id is null
Description is not null

## Date Filtering

### Specific Date

CreateDate eq '2026-02-15'
ModifyDate gt '2026-02-01'

### Date Ranges

CreateDate gte '2026-02-01' and CreateDate lte '2026-02-15'
ModifyDate gt '2026-01-01'

### Common Date Fields

- CreateDate: When the entity was created
- ModifyDate: When the entity was last modified
- StartDate: Planned start date
- EndDate: Planned end date

## Navigating Relationships

Access related entity properties using dot notation:

AssignedUser.Email eq 'user@example.com'
Project.Name eq 'Mobile App'
EntityState.IsFinal eq 'false'
Team.Name eq 'Backend Team'

## Boolean Fields

Boolean values in the TP API must be quoted as strings:

EntityState.IsFinal eq 'false'
EntityState.IsFinal eq 'true'

IMPORTANT: Unquoted boolean values (e.g., IsFinal eq false) will cause a 400 error.

## String Operations

Contains (substring match):

Name contains 'login'
Description contains 'critical'

Case-insensitive comparison is automatic for string equality.

## Membership Test (in)

Match against multiple values:

Id in (155,156,157)
EntityState.Name in ('Open','In Progress','Planned')
Priority.Name in ('High','Urgent')

The 'in' operator is more reliable than chaining multiple 'or' conditions.

## Numeric Comparisons

Effort gt 8
Priority.Id eq 1
Id ge 1000

## Complex Examples

### High-Priority Unassigned Bugs

entity_type="Bug"
where="Priority.Name eq 'High' and AssignedUser.Id is null"

### Recently Modified Open Stories

entity_type="UserStory"
where="EntityState.Name eq 'Open' and ModifyDate gt '2026-02-01'"

### Tasks in Multiple Projects

entity_type="Task"
where="Project.Id in (100,101)"

### Stories with Effort Range

entity_type="UserStory"
where="Effort ge 3 and Effort le 8"

## Raw WHERE Clauses

For advanced scenarios, you can pass raw WHERE clause strings directly to the API. The server does not validate or transform these - they are passed as-is to Target Process.

This allows you to use any Target Process API filter syntax, including functions and operators not explicitly documented here.

## Common Gotchas

### Boolean values must be quoted
Boolean fields like IsFinal require single-quoted string values:
- CORRECT: EntityState.IsFinal eq 'false'
- WRONG: EntityState.IsFinal eq false (causes 400 error)

### Parentheses are limited
Parentheses around 'or' conditions cause 400 errors:
- FAILS: (EntityState.Name eq 'Open' or EntityState.Name eq 'Planned') and Project.Id eq 100
- WORKS: EntityState.Name in ('Open','Planned') and Project.Id eq 100

For 'and' conditions, use no-space parentheses syntax:
- WORKS: (EntityState.Name eq 'Open')and(Project.Id eq 100)

### Prefer 'in' over 'or' for multiple values
The 'in' operator is more reliable and concise:
- BETTER: Priority.Name in ('High','Urgent')
- FRAGILE: Priority.Name eq 'High' or Priority.Name eq 'Urgent'

### String values require single quotes
All string and date values must be wrapped in single quotes:
- CORRECT: EntityState.Name eq 'Open'
- WRONG: EntityState.Name eq Open (causes 400 error)

Numeric values must NOT be quoted:
- CORRECT: Id eq 123
- WRONG: Id eq '123'

## Tips

- Always quote string and boolean values with single quotes: Name eq 'value', IsFinal eq 'false'
- Do NOT quote numeric values: Id eq 123, Effort gt 5
- Use ISO date format in single quotes: CreateDate gt '2026-01-01'
- Navigate relationships with dots: AssignedUser.Email, EntityState.IsFinal
- Use 'in' for matching multiple values: Name in ('a','b','c')
- Combine conditions with 'and' (parentheses-free when possible)
- Use 'or' only for flat conditions without parentheses
- Use 'is null' / 'is not null' for empty field checks
`

const Authentication = `# Authentication

## Overview

The TP MCP Server uses access token authentication to connect to your Target Process instance.

## Requirements

You need two pieces of information:

1. **TP_DOMAIN**: Your Target Process domain
2. **TP_ACCESS_TOKEN**: Your API access token

## Getting Your Access Token

### Step 1: Log into Target Process

Navigate to your Target Process instance (e.g., https://company.tpondemand.com)

### Step 2: Access Settings

1. Click on your profile icon in the top-right corner
2. Select "Settings" or "Access Tokens"

### Step 3: Generate Access Token

1. Navigate to the "Access Tokens" section
2. Click "Create Access Token"
3. Give your token a descriptive name (e.g., "MCP Server")
4. Set expiration as needed (or "Never" for long-term use)
5. Click "Create"
6. Copy the token immediately - it will only be shown once

## Environment Variables

Set these environment variables before running the server:

### Linux/macOS

export TP_DOMAIN=company.tpondemand.com
export TP_ACCESS_TOKEN=your-access-token-here
./server

### Windows (PowerShell)

$env:TP_DOMAIN="company.tpondemand.com"
$env:TP_ACCESS_TOKEN="your-access-token-here"
.\server.exe

### Windows (Command Prompt)

set TP_DOMAIN=company.tpondemand.com
set TP_ACCESS_TOKEN=your-access-token-here
server.exe

## Domain Format

The domain should be your Target Process hostname without the protocol:

✓ Correct: company.tpondemand.com
✗ Incorrect: https://company.tpondemand.com
✗ Incorrect: https://company.tpondemand.com/

## Security Best Practices

1. **Never commit tokens**: Do not commit access tokens to version control
2. **Use .env files**: Store tokens in .env files that are gitignored
3. **Rotate regularly**: Rotate access tokens periodically
4. **Limit scope**: Use tokens with minimal required permissions
5. **Monitor usage**: Review access token usage in Target Process settings

## Troubleshooting

### "Authentication failed" error

- Verify TP_DOMAIN is correct (no https://, no trailing slash)
- Check that TP_ACCESS_TOKEN is valid and not expired
- Ensure the token has appropriate permissions

### "Connection refused" error

- Confirm your network can reach the Target Process instance
- Verify the domain name is spelled correctly
- Check for firewall or proxy issues

### Token expired

- Generate a new access token in Target Process settings
- Update the TP_ACCESS_TOKEN environment variable
- Restart the server
`

const searchContent = `# Search Guide

## Overview

The search tool is the most powerful tool in the TP MCP Server, allowing you to query any entity type with flexible filtering, sorting, and pagination.

## Supported Entity Types

The search tool supports 17 entity types:

- UserStory: User stories and requirements
- Bug: Bug reports
- Task: Tasks
- Feature: Features and epics
- Request: Customer requests
- Epic: Large initiatives
- TestCase: Test cases
- TestPlan: Test plans
- TestPlanRun: Test execution runs
- Release: Releases
- Iteration: Sprints/iterations
- TeamIteration: Team-specific iterations
- Build: Build records
- Project: Projects
- Team: Teams
- User: Users
- CustomActivity: Custom activities

## Basic Search

Search without filters returns the first 25 entities:

search(entity_type="UserStory")

## Filtering with WHERE

Use the where parameter to filter results:

search(
  entity_type="Bug",
  where="EntityState.Name eq 'Open'"
)

## Including Related Data

Use include to fetch related entities:

search(
  entity_type="UserStory",
  include=["AssignedUser", "EntityState", "Project"],
  take=10
)

## Pagination

Control pagination with take and skip:

search(entity_type="Task", take=50, skip=100)

## Sorting

Sort results with orderByField and orderByDirection:

search(
  entity_type="Bug",
  orderByField="CreateDate",
  orderByDirection="desc",
  take=25
)

Note: The TP API v1 supports single-field sorting only. Use orderByField for the field name and orderByDirection for asc/desc.

## Common Search Patterns

### Find my assigned work

search(
  entity_type="UserStory",
  where="AssignedUser.Email eq 'myemail@company.com' and EntityState.IsFinal eq 'false'"
)

### Find all non-final (active) items

search(
  entity_type="UserStory",
  where="EntityState.IsFinal eq 'false'",
  take=100
)

### Find items in specific statuses

search(
  entity_type="Bug",
  where="EntityState.Name in ('Open','In Progress','Code Review')",
  take=100,
  include=["AssignedUser", "EntityState"]
)

### Find work assigned to a user by ID (structured filter)

search(
  entity_type="UserStory",
  assignedUser=123
)

### Find items by project

search(
  entity_type="Bug",
  where="Project.Id eq 100",
  take=100
)

### Find recent items

search(
  entity_type="Task",
  where="CreateDate gt '2026-02-01'",
  orderByField="CreateDate",
  orderByDirection="desc"
)

### Complex filters

search(
  entity_type="UserStory",
  where="Project.Id eq 100 and EntityState.Name eq 'Open' and Priority.Name eq 'High'",
  include=["AssignedUser", "Feature"],
  orderByField="Priority.Id",
  orderByDirection="asc"
)

## Response Format

Search returns:

- items: Array of matching entities
- total: Total count of matches (if available)
- next_cursor: Cursor for pagination (if available)

## Performance Tips

1. Use specific WHERE clauses to reduce result set
2. Only include fields you need
3. Use take to limit results
4. Consider pagination for large datasets
5. Index-friendly filters: Id, CreateDate, EntityState
`

const entityContent = `# Entity Management

## Overview

The TP MCP Server provides three tools for entity management:
- get_entity: Retrieve a single entity
- create_entity: Create a new entity
- update_entity: Update an existing entity

## get_entity

Retrieve complete details for a single entity.

### Basic Usage

get_entity(entity_type="UserStory", id=1234)

### With Related Data

get_entity(
  entity_type="Bug",
  id=5678,
  include=["AssignedUser", "EntityState", "Project", "Comments"]
)

### Common Include Fields

- AssignedUser: Who the item is assigned to
- EntityState: Current state/status
- Project: Associated project
- Team: Associated team
- Feature: Parent feature (for stories)
- Priority: Priority level
- Comments: All comments
- Attachments: All attachments

## create_entity

Create new entities in Target Process.

### Required Fields

Each entity type has required fields. Common requirements:

**UserStory:**
- Name (string)
- Project (object with Id)

**Bug:**
- Name (string)
- Project (object with Id)

**Task:**
- Name (string)
- Project (object with Id)

### Examples

#### Create User Story

create_entity(
  entity_type="UserStory",
  data={
    "Name": "User can reset password",
    "Description": "As a user, I want to reset my password via email",
    "Project": {"Id": 100},
    "Priority": {"Id": 2},
    "Effort": 5
  }
)

#### Create Bug

create_entity(
  entity_type="Bug",
  data={
    "Name": "Login page shows 404",
    "Description": "Detailed reproduction steps...",
    "Project": {"Id": 100},
    "Priority": {"Id": 1},
    "Severity": {"Id": 2}
  }
)

#### Create Task

create_entity(
  entity_type="Task",
  data={
    "Name": "Write unit tests for login",
    "Project": {"Id": 100},
    "UserStory": {"Id": 1234},
    "Effort": 3
  }
)

## update_entity

Update existing entities.

### Partial Updates

You only need to provide fields you want to change:

update_entity(
  entity_type="UserStory",
  id=1234,
  data={"Name": "Updated story name"}
)

### Common Updates

#### Change Status

update_entity(
  entity_type="Bug",
  id=5678,
  data={"EntityState": {"Id": 5}}
)

#### Assign to User

update_entity(
  entity_type="Task",
  id=9012,
  data={"AssignedUser": {"Id": 123}}
)

#### Update Effort

update_entity(
  entity_type="UserStory",
  id=1234,
  data={"Effort": 8}
)

#### Multiple Fields

update_entity(
  entity_type="Bug",
  id=5678,
  data={
    "EntityState": {"Id": 6},
    "AssignedUser": {"Id": 456},
    "Priority": {"Id": 1}
  }
)

## Field Reference by Type

Entities reference related items by Id:

- Project: {"Id": 100}
- AssignedUser: {"Id": 123}
- EntityState: {"Id": 5}
- Priority: {"Id": 2}
- Team: {"Id": 10}

Use inspect_object to discover available fields for each entity type.

## Tips

1. Use get_entity to see current values before updating
2. Reference related entities by Id only
3. Use inspect_object to discover field names and types
4. Required fields vary by entity type - check validation errors
5. Some fields are read-only and cannot be updated
`

const commentContent = `# Comments

## Overview

Add and retrieve comments on any Target Process entity using two tools:
- add_comment: Add a new comment
- list_comments: Retrieve comments

## add_comment

Add a comment to any entity.

### Basic Usage

add_comment(
  entity_type="UserStory",
  entity_id=1234,
  description="This looks good to go!"
)

### Multi-line Comments

add_comment(
  entity_type="Bug",
  entity_id=5678,
  description="Reproduction steps:
1. Navigate to login page
2. Enter invalid credentials
3. Click submit
Result: Error message not displayed"
)

### Markdown Support

Comments support markdown formatting:

add_comment(
  entity_type="UserStory",
  entity_id=1234,
  description="## Testing Notes

- [x] Unit tests passing
- [x] Integration tests passing
- [ ] Performance tests pending

**Status:** Ready for review"
)

### Status Updates

add_comment(
  entity_type="Task",
  entity_id=9012,
  description="Completed implementation. Code review: https://github.com/company/repo/pull/456"
)

## list_comments

Retrieve comments from an entity.

### Basic Usage

list_comments(
  entity_type="UserStory",
  entity_id=1234
)

### Pagination

Get the 10 most recent comments:

list_comments(
  entity_type="Bug",
  entity_id=5678,
  take=10,
  skip=0
)

Get the next 10:

list_comments(
  entity_type="Bug",
  entity_id=5678,
  take=10,
  skip=10
)

## Comment Fields

Each comment includes:

- Id: Comment ID
- Description: Comment text
- CreateDate: When the comment was created
- Owner: User who created the comment
- ParentId: ID of the parent entity
- ParentType: Type of the parent entity

## Use Cases

### Code Review Feedback

add_comment(
  entity_type="UserStory",
  entity_id=1234,
  description="Code review feedback:
- Consider extracting validation logic
- Add error handling for edge cases
- Update unit tests for new scenario"
)

### QA Notes

add_comment(
  entity_type="Bug",
  entity_id=5678,
  description="Verified fix in staging environment. All test cases passing."
)

### Progress Updates

add_comment(
  entity_type="Task",
  entity_id=9012,
  description="50% complete. Database schema updated. Working on API endpoints next."
)

### Blocking Issues

add_comment(
  entity_type="UserStory",
  entity_id=1234,
  description="⚠️ BLOCKED: Waiting for API documentation from backend team."
)

## Tips

1. Use comments for status updates and collaboration
2. Markdown formatting is supported
3. Comments appear in chronological order
4. Use take/skip for pagination on entities with many comments
5. Comments are timestamped with CreateDate
`

const attachmentContent = `# Attachments

## Overview

Manage file attachments on Target Process entities:
- list_attachments: View attachments on an entity
- download_attachment: Download a specific attachment

## list_attachments

List all attachments on an entity.

### Basic Usage

list_attachments(
  entity_type="Bug",
  entity_id=5678
)

### With Pagination

list_attachments(
  entity_type="UserStory",
  entity_id=1234,
  take=10,
  skip=0
)

## Attachment Fields

Each attachment includes:

- Id: Attachment ID (use for download)
- Name: Original filename
- ContentType: MIME type (e.g., "image/png")
- Size: File size in bytes
- CreateDate: Upload timestamp
- Owner: User who uploaded the file

## download_attachment

Download an attachment by ID.

### Basic Usage

download_attachment(attachment_id=9012)

### Specify Output Path

download_attachment(
  attachment_id=9012,
  output_path="./downloads/screenshot.png"
)

### Download to Specific Directory

download_attachment(
  attachment_id=9012,
  output_path="./bug-screenshots/"
)

## Common Workflows

### List and Download All Attachments

1. List attachments:

list_attachments(entity_type="Bug", entity_id=5678)

2. Download each attachment:

download_attachment(attachment_id=9012, output_path="./attachment1.png")
download_attachment(attachment_id=9013, output_path="./attachment2.pdf")

### Download Screenshots from a Bug

1. List attachments:

list_attachments(entity_type="Bug", entity_id=5678)

2. Filter for images and download:

download_attachment(
  attachment_id=9012,
  output_path="./bug-5678-screenshot-1.png"
)

### Check Attachment Metadata

Use list_attachments to check:
- File size before downloading
- Content type to identify file format
- Upload date and owner

## Supported File Types

Target Process supports common file types:

- Images: PNG, JPEG, GIF, SVG
- Documents: PDF, DOC, DOCX, XLS, XLSX
- Archives: ZIP, RAR, TAR, GZ
- Code: TXT, JSON, XML, LOG
- And many more

## Tips

1. Check ContentType to verify file format
2. Check Size before downloading large files
3. Use descriptive output_path names
4. Organize downloads by entity type or ID
5. List attachments first to get IDs for download

## Limitations

The MCP server currently supports:
- Listing attachments
- Downloading attachments

Uploading attachments is not yet supported. Use the Target Process web UI to upload files.
`

const inspectContent = `# API Inspection

## Overview

The inspect_object tool helps you discover entity schemas, available fields, and API metadata.

## Basic Usage

### List All Entity Types

inspect_object()

Returns all available entity types with descriptions.

### Inspect Specific Entity

inspect_object(object_type="UserStory")

Returns detailed schema for UserStory including:
- Available fields
- Field types
- Required fields
- Relationships

## Use Cases

### Discover Available Fields

Before creating or updating entities, inspect to find required and optional fields:

inspect_object(object_type="Bug")

Returns fields like:
- Name (string, required)
- Description (string, optional)
- Project (reference, required)
- Priority (reference, optional)
- EntityState (reference, optional)

### Understand Relationships

Discover how entities relate to each other:

inspect_object(object_type="Task")

Shows relationships:
- UserStory: Parent user story
- AssignedUser: Assigned user
- Project: Associated project
- Team: Associated team

### Find Field Names for Queries

Discover exact field names for WHERE clauses:

inspect_object(object_type="UserStory")

Find filterable fields:
- EntityState.Name
- AssignedUser.Email
- Project.Id
- Priority.Name
- CreateDate

### Validate Entity Type

Check if an entity type exists before searching:

inspect_object(object_type="CustomEntity")

## Entity Categories

### Work Items
- UserStory
- Bug
- Task
- Feature
- Epic

### Testing
- TestCase
- TestPlan
- TestPlanRun

### Organization
- Project
- Team
- Release
- Iteration
- TeamIteration

### Other
- Request
- Build
- User
- CustomActivity

## Field Types

Common field types you'll encounter:

- String: Text values
- Integer: Numeric IDs and counts
- Date: Timestamps
- Boolean: True/false flags
- Reference: Links to other entities (use {"Id": value})

## Tips

1. Run without parameters first to see all entity types
2. Inspect before create to find required fields
3. Use to discover relationship field names
4. Check field types before setting values
5. Reference the output when building WHERE clauses

## Example Workflow

1. Discover entity type:

inspect_object()

2. Inspect specific type:

inspect_object(object_type="Bug")

3. Create entity with discovered fields:

create_entity(
  entity_type="Bug",
  data={
    "Name": "Login issue",
    "Project": {"Id": 100},
    "Priority": {"Id": 2}
  }
)
`

const paginationContent = `# Pagination

## Overview

The TP MCP Server supports cursor-based and offset-based pagination for search and list operations.

## Offset-Based Pagination

### search, list_comments, list_attachments

Use take and skip parameters:

- take: Number of items to return (default: 25, max: 1000)
- skip: Number of items to skip (default: 0)

### Examples

#### First Page

search(entity_type="UserStory", take=25, skip=0)

#### Second Page

search(entity_type="UserStory", take=25, skip=25)

#### Third Page

search(entity_type="UserStory", take=25, skip=50)

## Page Size Recommendations

- Small datasets: take=50
- Medium datasets: take=100
- Large datasets: take=500
- Maximum: take=1000

## Calculating Pages

Total items: 237
Page size: 25

- Page 1: skip=0, take=25
- Page 2: skip=25, take=25
- Page 3: skip=50, take=25
- ...
- Page 10: skip=225, take=25

## Response Metadata

Search responses include:

- items: Array of results
- total: Total count (if available)
- next_cursor: Cursor for next page (if available)

## Efficient Pagination

### Get Total Count First

search(entity_type="Bug", take=1)

Check total field to know how many pages.

### Iterate Through Pages

page_size = 100
page = 0

while True:
  result = search(
    entity_type="UserStory",
    take=page_size,
    skip=page * page_size
  )

  if not result.items:
    break

  # Process result.items
  page += 1

## Pagination with Filters

Pagination works with WHERE clauses:

search(
  entity_type="Bug",
  where="EntityState.Name eq 'Open'",
  take=50,
  skip=100
)

## Pagination with Sorting

Combine with orderByField/orderByDirection for consistent pagination:

search(
  entity_type="UserStory",
  orderByField="CreateDate",
  orderByDirection="desc",
  take=50,
  skip=0
)

## Performance Tips

1. Always use take to limit results
2. Use specific WHERE clauses to reduce total items
3. Sort by indexed fields (Id, CreateDate) for faster queries
4. Avoid very large skip values (>1000)
5. Consider filtering instead of paginating through everything

## Example: Export All Data

To export all entities of a type:

1. Start with skip=0
2. Request take=500 (larger batches)
3. Increment skip by 500
4. Continue until empty result

## Limits

- Maximum take: 1000
- Recommended skip: <1000 (performance degrades with large offsets)
- No hard limit on skip, but queries slow down
`
