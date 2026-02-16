# tp-mcp-go

A Go-based MCP (Model Context Protocol) server for the Apptio Target Process API. This server replaces the TypeScript version and is built using the [Foxy Contexts](https://github.com/strowk/foxy-contexts) framework.

## Build Instructions

```bash
go build ./cmd/server/
```

This will generate a `server` executable in the current directory.

## Configuration

The server requires two environment variables:

- `TP_DOMAIN` - Your Target Process domain (e.g., `your-domain.tpondemand.com`)
- `TP_ACCESS_TOKEN` - Your Target Process API access token

You can set these in your shell environment or provide them when running the server.

## Usage with Claude Desktop / Cline / Goose

Add the following configuration to your MCP client settings:

```json
{
  "mcpServers": {
    "targetprocess": {
      "command": "/path/to/server",
      "env": {
        "TP_DOMAIN": "your-domain.tpondemand.com",
        "TP_ACCESS_TOKEN": "your-token"
      }
    }
  }
}
```

Replace `/path/to/server` with the absolute path to the built server executable.

## Usage with Claude Code (Docker)

```bash
claude mcp add targetprocess -- docker run -i --rm -e TP_DOMAIN=your-domain.tpondemand.com -e TP_ACCESS_TOKEN=your-token ghcr.io/kildall/tp-mcp-go:latest
```

Replace `your-domain.tpondemand.com` and `your-token` with your actual values.

## Available Tools

The server provides the following 10 tools for interacting with Target Process:

- **search** - Search entities with filters (status, assigned user, project, team, etc.) and pagination support
- **get_entity** - Retrieve a single entity by type and ID with optional field inclusion
- **create_entity** - Create a new entity with name, description, project, team, and custom fields
- **update_entity** - Update entity fields including name, description, status, and assignments
- **add_comment** - Add a private comment to an entity
- **list_comments** - List all comments on an entity
- **list_attachments** - List all attachments on an entity
- **download_attachment** - Download attachment content by ID

## MCP Resources

The server provides the following documentation resources:

- `docs://getting-started` - Introduction and quick start guide
- `docs://tool-reference` - Detailed reference for all tools
- `docs://examples` - Usage examples and common workflows
- `docs://query-guide` - Guide to filtering and querying entities
- `docs://authentication` - Authentication setup and troubleshooting

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

```
tp-mcp-go/
├── cmd/server/          # Main server entry point
├── internal/
│   ├── app/            # Application lifecycle and DI module
│   ├── client/         # Target Process API client
│   │   └── auth/       # Authentication handling
│   ├── config/         # Configuration management
│   ├── domain/         # Domain models and business logic
│   │   ├── entity/     # Entity types and models
│   │   ├── query/      # Query building and filtering
│   │   └── errors/     # Domain-specific errors
│   ├── tools/          # MCP tool implementations
│   └── testutil/       # Testing utilities and fixtures
├── go.mod              # Go module definition
└── README.md           # This file
```

## License

This project is part of the Apptio Target Process MCP integration.
