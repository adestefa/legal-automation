# Mallon Legal Assistant

**Version**: v2.5.41  
**Status**: Active Development  
**Last Updated**: June 13, 2025

## Overview

Mallon Legal Assistant is an advanced legal document automation system designed to streamline the process of generating legal complaints from client documents. The system provides a step-by-step workflow for document selection, data extraction, and automated legal document generation.

## Features

- **Document Selection**: Choose from various legal documents and client files
- **Data Extraction**: Automated extraction of client information from selected documents
- **Template Population**: Dynamic generation of legal complaints with client-specific data
- **Preview & Review**: Multi-tab interface for reviewing structured data and document preview
- **User Management**: JSON-based user administration system
- **Session Management**: Persistent user selections during navigation

## Technology Stack

- **Backend**: Go (Golang) with Gin framework
- **Frontend**: HTML5, Tailwind CSS, Alpine.js, HTMX
- **Templates**: Go template engine
- **Session Management**: Server-side session storage
- **Version Control**: Git with feature branch workflow

## Project Structure

```
proj-mallon/
├── dev/                    # Development implementation
│   ├── backend/           # Go server application
│   ├── frontend/          # HTML/CSS/JS interface
│   ├── templates/         # Go template files
│   ├── legal_artifacts/   # Sample legal documents
│   ├── start.sh          # Server start script
│   ├── stop.sh           # Server stop script
│   └── restart.sh        # Server restart script
├── yinsen/               # Project management & task tracking
│   ├── 1_queue/          # Pending tasks and defects
│   ├── 2_dev/            # Active development
│   ├── 3_qa/             # Quality assurance
│   └── 4_done/           # Completed tasks
└── README.md             # This file
```

## Quick Start

### Prerequisites
- Go 1.19 or higher
- Git

### Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd proj-mallon
   ```

2. **Start the application**:
   ```bash
   ./start.sh
   ```

3. **Access the application**:
   Open your browser and navigate to `http://localhost:8080`

4. **Stop the application**:
   ```bash
   ./stop.sh
   ```

### Development Workflow

1. **Create feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make changes and test locally**:
   ```bash
   ./restart.sh  # Restart server after changes
   ```

3. **Update version number** in templates before committing

4. **Create pull request** when ready for review

## Current Development Status

### Completed Features ✅
- Document selection workflow
- User session management
- Dynamic document preview
- User.json-based administration
- Step-by-step navigation with back buttons
- Direct case folder navigation

### Active Development 🔄
- **DEFECT 2**: Missing Content Tab error resolution
- **TASK 8**: Enhanced document generation engine
- **TASK 33**: iCloud document save functionality

### Version History
- **v2.5.41**: Current stable version with all recent enhancements
- **v2.5.27**: Direct case folder navigation enhancement
- **v2.5.26**: UI improvements and navigation consistency
- **v2.5.25**: Step icons activation fixes
- **v0.1.0**: Initial working prototype

## API Endpoints

- `GET /` - Main application interface
- `POST /api/generate-summary` - Generate document summary
- `POST /api/select-case-folder` - Case folder selection
- `GET /api/documents` - List available documents
- `POST /api/session/set` - Set session data
- `GET /api/session/get` - Retrieve session data

## Testing

The application includes comprehensive testing through:
- Local development server testing
- Feature branch validation
- Pull request review process
- Quality assurance workflow

## Contributing

1. Check for active defects before starting new tasks
2. Use feature branches for all development
3. Update version numbers for releases
4. Test locally before creating pull requests
5. Follow the Yinsen task management workflow

## Support

For issues and feature requests, please create GitHub issues or contact the development team.

## License

All rights reserved - Satori Tech Consulting 2025