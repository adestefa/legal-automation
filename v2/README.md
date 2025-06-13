# Mallon Legal V2 - Clean Application

This is the clean V2 version of the Mallon Legal document generation application, extracted from the mixed v1/v2 development environment.

## Directory Structure

```
v2/
├── main.go                 # Main application entry point (from main_v2.go)
├── go.mod                 # Go module definition (from go_v2.mod)
├── go.sum                 # Go dependencies
├── handlers/              # HTTP handlers for UI and API
│   └── ui_handlers.go     # HTMX-based UI handlers
├── services/              # Business logic services
│   ├── document_service.go # Document processing service
│   └── icloud_service.go   # iCloud integration service
├── templates/             # HTML templates for SSR
├── frontend/              # Static web assets
├── scripts/               # Utility scripts
│   ├── start.sh          # Start server
│   ├── stop.sh           # Stop server
│   └── restart.sh        # Restart server
└── docs/                  # Documentation
    ├── DEPLOYMENT_READY.md
    └── REFACTOR_COMPLETE.md
```

## Features

- Dynamic document processing (Task 8)
- Document editing capabilities
- Go Server-Side Rendering (SSR) with HTMX
- iCloud integration for document management
- Legal document generation with template mapping

## Quick Start

1. Navigate to the v2 directory:
   ```bash
   cd /Users/corelogic/satori-dev/clients/proj-mallon/v2
   ```

2. Start the application:
   ```bash
   ./scripts/start.sh
   ```

3. Access the application at: http://localhost:8080

## Version

- Server Version: v2.5.28
- Go Version: 1.21
- Framework: Gin (Go web framework)
- Frontend: HTMX + Server-Side Rendering

## Key Dependencies

- github.com/gin-gonic/gin v1.10.0
- github.com/unidoc/unioffice v1.29.0
- golang.org/x/crypto v0.21.0

## Notes

- This is a clean extraction of the v2 codebase
- All v1 legacy code has been removed
- Hardcoded paths have been updated for the new directory structure
- The application maintains compatibility with existing artifacts and test data