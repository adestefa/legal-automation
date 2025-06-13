# Mallon Legal Document Automation - Prototype

This is a minimal working prototype demonstrating the workflow for automating legal complaint form generation based on source documents.

## Project Structure

```
dev/
├── backend/                # Go backend API
│   ├── main.go             # Main server code
│   └── go.mod              # Go module file
├── frontend/               # Frontend HTML/JS/CSS
│   └── index.html          # Single page application
└── templates/              # Template mapping files
    └── fcra-credit-card-fraud-mapping.json # Template schema
```

## Features Demonstrated

1. **Document Selection**: Browse and select source documents from a simulated iCloud connection
2. **Template Selection**: Choose a complaint form template to populate
3. **Information Extraction**: Extract key client information from source documents
4. **Preview Generation**: Generate a structured preview of extracted information
5. **Template Population**: Insert extracted data into the complaint form template
6. **Document Generation**: Create the final populated document

## Running the Prototype

### Prerequisites

- Go 1.21 or higher
- Basic knowledge of HTML, Go, and REST APIs

### Backend Setup

1. Navigate to the `backend` directory:
   ```
   cd dev/backend
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Start the server:
   ```
   go run main.go
   ```

The server will start on `http://localhost:8080`.

### Using the Application

1. Open a web browser and navigate to `http://localhost:8080`
2. Follow the step-by-step workflow:
   - Step 1: Select source documents
   - Step 2: Select a complaint form template
   - Step 3: Review extracted information
   - Step 4: Generate the populated document

## Implementation Notes

This prototype simulates several aspects of the full system:

- **Document Processing**: For the prototype, we use pre-extracted data instead of real-time processing
- **Template Population**: We demonstrate the UI flow but don't actually modify DOCX files
- **iCloud Integration**: File selection simulates the iCloud connection

In a production implementation, these would be replaced with:

- Real PDF/DOCX processing using UniDoc/UniPDF
- Actual template population with proper document manipulation
- True iCloud OAuth and API integration

## Next Steps

1. Implement actual document text extraction with UniDoc/UniPDF
2. Add pattern matching for information extraction from documents
3. Integrate with iCloud API for file access
4. Implement template population engine for DOCX files
5. Add Zilliz Cloud integration for semantic search
