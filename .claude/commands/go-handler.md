Description
Generate Go HTTP handlers optimized for Sensei server architecture
Parameters

endpoint: API endpoint path (e.g., "/api/documents")
method: HTTP method (GET, POST, PUT, DELETE)
function_name: Handler function name

Instructions
Create a Go HTTP handler following Yinsen's standards:
Code Requirements:

Use Go's built-in net/http package (no external dependencies)
Implement proper error handling with context
Include JSON response helpers
Add request validation
Follow Satori Tech naming conventions
Include performance considerations
Add appropriate logging

Structure:
gofunc {function_name}Handler(w http.ResponseWriter, r *http.Request) {
    // Request validation
    // Business logic
    // Error handling
    // JSON response
}
Standards:

Zero external dependencies unless absolutely necessary
Performance-first approach
Security by design
Clear error messages
Proper HTTP status codes

Provide complete, production-ready code that integrates with existing Sensei server patterns.
Context
Maintains consistency with bare-metal Go server architecture and Yinsen's dependency-minimal philosophy.