Description
Create versioned project backup before beginning development work
Parameters

version: Current project version (will auto-increment patch if not specified)

Instructions
Create a safety backup following Satori Tech protocols:

Identify Project: Determine current project name and version from readme.md or main files
Generate Backup Name: Format as proj-{name}_YYYY-MM-DD_vX.X.X.zip
Create Archive: Zip the entire project directory
Store Location: Place in RELEASES/ folder or appropriate backup directory
Version Increment: Update version number in project files (patch increment)
Document: Log backup creation in history.md

Confirm backup creation with file size and location. This ensures recovery capability before any development changes.
Context
Critical safety protocol for preserving working state before modifications, especially important for production systems.