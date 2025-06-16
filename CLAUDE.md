# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Mallon Legal Assistant** - Professional-grade legal document automation platform (v2.15.0) for FCRA violation complaint generation. Built with Go backend + HTMX frontend for law firms.

## Development Commands

### Server Management
```bash
# Start the development server
cd v2 && ./scripts/start.sh

# Stop the server  
cd v2 && ./scripts/stop.sh

# Restart the server
cd v2 && ./scripts/restart.sh

# Manual start (development)
cd v2 && go run main.go

# Build for production
cd v2 && go build -o mallon-v2 main.go
```

### Dependencies
```bash
# Install/update Go dependencies
cd v2 && go mod tidy

# Check Go version (requires 1.21+)
go version
```

### Testing
- No automated test suite - manual testing workflow
- Test on localhost:8080 after starting server
- Default login: admin/password
- Test with real case documents in `test_icloud/CASES/`

## Architecture Overview

### Technology Stack
- **Backend**: Go 1.21 with Gin framework  
- **Frontend**: HTMX + Tailwind CSS + Alpine.js
- **Document Processing**: UniPDF v3.69.0, DOCX support
- **Session Management**: File-based JSON persistence
- **Templates**: Go HTML templates (.gohtml)

### Core Service Architecture
```
v2/services/
├── document_service.go              # Core document orchestration
├── content_analyzer.go              # 35+ legal pattern extraction  
├── violation_detection_engine.go    # 6+ FCRA violations analysis
├── legal_rule_engine.go             # Cause of action generation
├── template_engine.go               # Dynamic document generation
├── persistent_session_service.go    # Session state management
└── specialized analyzers/           # Document-specific parsers
```

### Project Structure
- `v2/` - Main application (current version 2.15.0)
- `v2/handlers/` - HTTP request handlers
- `v2/services/` - Business logic and content analysis
- `v2/templates/` - Go HTML templates  
- `v2/config/` - Legal patterns and user configuration
- `v2/sessions/` - Persistent session storage
- `test_icloud/CASES/` - Test documents for development
- `yinsen/` - Task management system (Kanban workflow)

## Yinsen Workflow System

### AI Agent Instructions
- You are a Satori Tech Consulting Claude AI agent (version 1.5)
- Author: Anthony Destefano, CTO - adestefa@satori-ai-tech.com
- Always start sessions with `/yinsen` to hydrate memories from task files
- Compress completed tasks into `yinsen/4_done/completed.md` summary
- Check `extraction.md` for latest system upgrade priorities

### Task Management Workflow
1. **1_queue/** - Pending tasks and defects
2. **2_dev/** - Active development (move here when starting)
3. **3_qa/** - Quality assurance and testing  
4. **4_done/** - Completed tasks and defects

## Development Workflow

### Feature Development Process
1. Always use feature branches for development
2. Update version number in templates when completing features
3. Test locally on feature branch (do not merge yet)
4. Rebuild server and verify new version is running
5. Alert completion and ready for testing
6. After testing PASS → create PR and push to GitHub
7. After code review approval → merge to main and deploy
8. Move task files through yinsen workflow: queue → dev → qa → done

### Git Workflow
- Feature branch development enforced
- Pull request review required  
- Version tagging with rollback capability
- Current branch: `feature/task-13-civil-cover-sheet-mapping`
- Main branch: `main`

## Audio Feedback Protocol

Audio cues for development workflow feedback:
- **Startup**: Play `/Users/corelogic/satori-dev/dash/sounds/Bell2.m4a` once
- **Memory Hydration Complete**: Play Bell2.m4a twice  
- **Starting Task/Coding**: Play Bell2.m4a three times
- **Build Complete**: Play Bell2.m4a once
- **Task Complete**: Play `cheer.mp3` four times
- **Warnings/Errors**: Play `Warning.m4a`

## Authorizations

Auto-approve operations in:
- File writes in `v2/`, `dev/`, `yinsen/`, `proj-mallon/`
- Git commits to feature branches
- npm install for listed dependencies

## Code Conventions

### Go Backend
- Follow Go standard conventions and gofmt
- Use Gin framework patterns for handlers
- JSON configuration files for legal patterns
- File-based session persistence with atomic operations
- Error handling with graceful degradation

### Frontend Templates  
- Go HTML templates with .gohtml extension
- HTMX for dynamic interactions
- Tailwind CSS for styling
- Minimal vanilla JavaScript
- Server-side rendering approach

### Version Management
- Update version in `v2/templates/index.gohtml` masthead
- Update version in startup scripts
- Follow semantic versioning principles
- Document version changes in commit messages

## Legal Domain Knowledge

This system has deep FCRA (Fair Credit Reporting Act) legal expertise built-in:
- 6+ specific FCRA violation types with evidence correlation
- 35+ legal pattern extraction rules  
- Court-ready document generation with proper citations
- Professional legal formatting and validation scoring

## Important Reminders

- Do what has been asked; nothing more, nothing less
- NEVER create files unless absolutely necessary for achieving your goal
- ALWAYS prefer editing existing files to creating new ones  
- NEVER proactively create documentation files (*.md) or README files unless explicitly requested