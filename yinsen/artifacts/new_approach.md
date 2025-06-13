# Core Idea of the New Approach:

Go Backend (main.go and modules): Handles all business logic, data fetching, state management (primarily via user sessions or request flow), and renders HTML views using Go's html/template package.

HTML Templates (e.g., templates/index.gohtml): Standard HTML files with Go template syntax to display dynamic data. These will contain HTMX attributes.

HTMX (in the HTML): Used to make requests to the Go backend when users interact with the UI. The backend responds with HTML fragments, which HTMX then swaps into the current page, avoiding full page reloads for many interactions.

driver.js: Becomes minimal or potentially unnecessary. HTMX handles the bulk of dynamic updates. Any remaining JavaScript would be for very specific, minor UI enhancements that don't involve core data flow.

1. main.go and Backend Modules (Conceptual Refactor)

Your main.go will set up routes. Handlers will now be responsible for rendering HTML.

# Suggested Modular Structure:

main.go: Router setup, global configurations, starts the server.

## handlers/: Contains Go files for different HTTP handlers.
- ui_handlers.go: Handles requests that render UI pages or HTMX fragments.
- api_handlers.go: If you still need JSON APIs for some specific purpose (though HTMX primarily uses HTML).
- auth_handlers.go: For login, logout.

### services/: Business logic.
- icloud_service.go: Logic for interacting with iCloud (fetching folders, documents - much of your existing logic can go here).
- document_service.go: Logic for processing documents, generating summaries, populating templates.

### models/: Struct definitions (ClientCase, Document, etc. - largely existing).

### templates/: Directory to store your *.gohtml template files.

### view/ (or render/): Helper functions for parsing and rendering Go HTML templates.


# **main.go (Simplified Router Setup Example):

```
package main

import (
	"log"
	"net/http"
	"yourapp/handlers" // Assuming your handlers are in this package
	"yourapp/auth"     // Assuming auth logic

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Serve static files (CSS, minimal JS if any)
	router.Static("/static", "./static") // Assuming a 'static' folder for CSS, etc.

	// Load HTML templates
	// You'd typically parse all templates at startup.
	// For simplicity, individual handlers might load their specific templates,
	// or you use a helper as suggested in `view/` module.
	// Example: router.LoadHTMLGlob("templates/**/*") if using Gin's built-in rendering.

	// Auth routes
	router.GET("/login", handlers.ShowLoginPage) // Renders login.gohtml
	router.POST("/login", handlers.HandleLoginPost)
	router.POST("/logout", handlers.HandleLogout)

	// Main application routes (protected)
	// Note: Gin's built-in template rendering expects templates in a specific layout
	// or you manage it fully yourself. For this example, assume handlers render directly.
	appGroup := router.Group("/")
	appGroup.Use(auth.Middleware()) // Your existing auth middleware
	{
		// Initial page load for the multi-step UI
		appGroup.GET("/", handlers.ShowDocumentAutomationPage) // Renders the main multi-step UI

		// HTMX endpoints for partial UI updates
		// These will return HTML fragments
		appGroup.GET("/ui/icloud-folders", handlers.GetICloudFoldersHTML)
		appGroup.GET("/ui/case-folders", handlers.GetCaseFoldersHTML)
		appGroup.GET("/ui/icloud-documents", handlers.GetICloudDocumentsHTML)
		appGroup.POST("/ui/select-documents", handlers.HandleSelectDocumentsAndShowTemplates) // Example
		appGroup.POST("/ui/select-template", handlers.HandleSelectTemplateAndShowReview)     // Example
		appGroup.POST("/ui/generate-document", handlers.HandleGenerateDocument)              // Example
		// ... other HTMX endpoints for each step/interaction
	}

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
```


# **handlers/ui_handlers.go (Conceptual Snippets):

```
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"yourapp/services" // Your service layer

	"github.com/gin-gonic/gin"
)

// Assume templates are parsed at startup or via a helper
var tmpl *template.Template

func init() {
	// Robust template parsing
	tmpl = template.Must(template.ParseGlob("templates/*.gohtml"))
    // For nested templates: template.Must(template.ParseGlob("templates/**/*.gohtml"))
}

// ShowDocumentAutomationPage renders the initial main page
func ShowDocumentAutomationPage(c *gin.Context) {
	// Initial data for the page (e.g., current step, user info)
	data := gin.H{
		"CurrentStep": 0,
		"Username":    c.GetString("username"), // From auth middleware
		// ... any other initial data
	}
	// c.HTML(http.StatusOK, "index.gohtml", data) // If using Gin's HTML rendering
    err := tmpl.ExecuteTemplate(c.Writer, "index.gohtml", data)
    if err != nil {
        log.Printf("Error executing template index.gohtml: %v", err)
        c.String(http.StatusInternalServerError, "Error rendering page")
    }
}

// GetICloudFoldersHTML handles HTMX request for iCloud parent folders
func GetICloudFoldersHTML(c *gin.Context) {
	// username := c.GetString("username") // from auth
	// For simplicity, skipping actual iCloud auth call details here
	// Assume icloudService.GetRootFolders(username) returns your ICloudDocument structs
	folders, err := services.GetRealICloudFolders("", "") // Pass appropriate params
	if err != nil {
		// Return an HTML fragment indicating error
		// c.HTML(http.StatusOK, "_error_fragment.gohtml", gin.H{"Error": "Could not load iCloud folders."})
        tmpl.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", gin.H{"Error": "Could not load iCloud folders."})
		return
	}

	data := gin.H{"Folders": folders}
	// c.HTML(http.StatusOK, "_icloud_folder_list.gohtml", data) // Renders only the list part
    tmpl.ExecuteTemplate(c.Writer, "_icloud_folder_list.gohtml", data)
}

// GetICloudDocumentsHTML handles HTMX request for documents in a folder
func GetICloudDocumentsHTML(c *gin.Context) {
    folderPath := c.Query("folder") // e.g., "/CASES/Yousef_Eman"
    if folderPath == "" {
        tmpl.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", gin.H{"Error": "Folder path is required."})
        return
    }

    // Fetch documents using your existing service logic
    // Note: getRealICloudDocuments in your main.go expects username, appPassword. Adapt as needed.
    docs, err := services.GetRealICloudDocuments("", "", folderPath) // Adapt with auth if necessary
    if err != nil {
        log.Printf("Error fetching documents for %s: %v", folderPath, err)
        tmpl.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", gin.H{"Error": "Could not load documents."})
        return
    }

    data := gin.H{"Documents": docs, "SelectedFolder": folderPath}
    // This template will be an HTML fragment containing the list of documents
    tmpl.ExecuteTemplate(c.Writer, "_document_list_fragment.gohtml", data)
}


// HandleSelectDocumentsAndShowTemplates processes selected documents and shows template selection UI
func HandleSelectDocumentsAndShowTemplates(c *gin.Context) {
    // 1. Get selected document IDs/paths from form submission (c.PostFormArray("selectedDocs"))
    // 2. Store them in session or pass them along.
    // 3. Fetch available templates (from services.GetTemplates())
    // 4. Prepare data for the template selection view/fragment.
    selectedDocs := c.PostFormArray("selectedDocs[]") // HTMX might send array like this
    log.Printf("Selected documents by user: %v", selectedDocs)
    // For now, just log and proceed to show templates
    // In a real app, you'd store these selectedDocs in the user's session or pass them to the next step

    templates, _ := services.GetTemplates() // Assuming this function exists

    data := gin.H{
        "Templates":    templates,
        "SelectedDocs": selectedDocs, // Pass along to next step if needed
        "CurrentStep":  2,
    }
    // Render a fragment that replaces the current step's content with template selection UI
    tmpl.ExecuteTemplate(c.Writer, "_template_selection_fragment.gohtml", data)
}

// ... other handlers for subsequent steps, each rendering an HTML fragment or full page ...

```



# **templates/index.gohtml (Replaces frontend/index.html)

```
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Legal Document Automation - Mallon Law (Go SSR + HTMX)</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    <link rel="stylesheet" href="/static/site.css">
</head>
<body class="bg-gray-50">
    <div class="container mx-auto px-4 py-8 max-w-6xl">
        <header class="mb-8">
            <h1 class="text-3xl font-bold text-gray-800">Legal Document Automation <span class="text-sm font-normal text-gray-500">v2.0 (Go+HTMX)</span></h1>
             </header>

        <div class="mb-8">
            </div>

        <div id="workflow-steps">
            {{if eq .CurrentStep 0}}
                {{template "_step0_case_setup.gohtml" .}}
            {{else if eq .CurrentStep 1}}
                {{template "_step1_document_selection.gohtml" .}}
            {{/* ... and so on for other steps ... */}}
            {{end}}
        </div>
    </div>

    </body>
</html>

{{define "_step0_case_setup.gohtml"}}
<div id="step-0-content" class="bg-white p-6 rounded-lg shadow-md mb-6">
    <h2 class="text-xl font-semibold mb-4">Step 0: Setup Case Folder</h2>
    <p class="text-gray-600 mb-6">Select your iCloud parent folder...</p>

    <div class="mb-6">
        <h3 class="text-lg font-medium mb-3">Select Parent Folder</h3>
        <button hx-get="/ui/icloud-folders"
                hx-target="#icloud-folder-list-container"
                hx-swap="innerHTML"
                class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
            Load iCloud Root Folders
        </button>
        <div id="icloud-folder-list-container" class="mt-4">
            </div>
    </div>

    <div id="case-folder-selection-container" class="mb-6">
        </div>
    <div class="mt-6 flex justify-end">
         <button hx-post="/ui/setup-case-complete" {{/* Or similar endpoint */}}
                hx-target="#workflow-steps"
                hx-swap="innerHTML"
                class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
            Continue to Document Selection
        </button>
    </div>
</div>
{{end}}


{{define "_icloud_folder_list.gohtml"}}
    {{if .Folders}}
        <div class="grid grid-cols-2 gap-4">
            {{range .Folders}}
                <div class="border rounded-lg p-4 cursor-pointer hover:bg-blue-50 transition"
                     hx-get="/ui/case-folders?parent={{.Path}}"  {{/* HTMX to load sub-folders */}}
                     hx-target="#case-folder-selection-container"
                     hx-swap="innerHTML"
                     title="Load subfolders for {{.Name}}">
                    <div class="font-medium text-sm text-black">{{.Name}}</div>
                    <div class="text-xs text-gray-500">Modified: {{.Modified.Format "Jan 2, 2006"}}</div>
                </div>
            {{end}}
        </div>
    {{else}}
        <p class="text-gray-500">No iCloud folders found or could not load.</p>
    {{end}}
{{end}}

{{define "_document_list_fragment.gohtml"}}
<h3 class="text-lg font-medium mb-3">Documents in {{.SelectedFolder}}</h3>
{{if .Documents}}
<form id="document-selection-form"
      hx-post="/ui/select-documents"
      hx-target="#workflow-steps" {{/* Target the main step area or a specific div for Step 2 */}}
      hx-swap="innerHTML">
    <ul class="divide-y divide-gray-200 border rounded-lg overflow-hidden">
        {{range .Documents}}
        <li class="p-3 hover:bg-gray-50 flex justify-between items-center">
            <div>
                <span class="text-black font-medium">{{.Name}}</span>
                <div class="text-xs text-gray-500">
                    Type: {{.Type}} | Size: {{.Size}} bytes | Modified: {{.Modified.Format "Jan 2, 2006"}}
                </div>
            </div>
            <input type="checkbox" name="selectedDocs[]" value="{{.Path}}" class="h-4 w-4 text-blue-600">
        </li>
        {{end}}
    </ul>
    <div class="mt-6 text-right">
        <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
            Continue with Selected Documents
        </button>
    </div>
</form>
{{else}}
<p class="text-gray-500 p-4">No documents found in this folder.</p>
{{end}}
{{end}}


{{define "_template_selection_fragment.gohtml"}}
<div id="step-2-content" class="bg-white p-6 rounded-lg shadow-md mb-6">
    <h2 class="text-xl font-semibold mb-4">Step 2: Select Complaint Template</h2>
    <p class="text-gray-600 mb-2">You selected {{len .SelectedDocs}} document(s).</p>
    <form hx-post="/ui/select-template" hx-target="#workflow-steps" hx-swap="innerHTML">
        {{range .SelectedDocs}}
        <input type="hidden" name="previouslySelectedDocs[]" value="{{.}}">
        {{end}}

        <div class="border rounded-lg overflow-hidden mt-4">
            <div class="bg-gray-50 p-3 border-b"><div class="text-sm font-medium">Available Templates</div></div>
            <ul class="divide-y divide-gray-200">
                {{range .Templates}}
                <li class="p-3 hover:bg-gray-50">
                    <label class="flex items-center cursor-pointer">
                        <input type="radio" name="selectedTemplateId" value="{{.ID}}" class="h-4 w-4 text-blue-600">
                        <div class="ml-3">
                            <span class="block font-medium text-black">{{.Name}}</span>
                            <span class="block text-sm text-gray-500">{{.Desc}}</span>
                        </div>
                    </label>
                </li>
                {{end}}
            </ul>
        </div>
        <div class="mt-6 flex justify-between">
            <button type="button"
                    hx-get="/" {{/* Or to the specific previous step URL that renders Step 1 */}}
                    hx-target="#workflow-steps"
                    class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50">
                Back
            </button>
            <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
                Process with Template
            </button>
        </div>
    </form>
</div>
{{end}}
```


# **templates/_error_fragment.gohtml (Example):

{{define "_error_fragment.gohtml"}}
<div class="p-4 bg-red-100 border border-red-300 text-red-700 rounded">
    <strong>Error:</strong> {{.Error}}
</div>
{{end}}


# **static/driver.js (Replaces frontend/driver.js)

This file becomes much simpler. You might only need it for:

Global HTMX event listeners if necessary (e.g., for loading indicators, error handling).
Very specific UI interactions that HTMX doesn't cover easily and aren't worth a server roundtrip.

```
// static/driver.js

document.body.addEventListener('htmx:beforeRequest', function(evt) {
    console.log('HTMX: Making a request to', evt.detail.pathInfo.requestPath);
    // You could show a global loading indicator here
});

document.body.addEventListener('htmx:afterSwap', function(evt) {
    console.log('HTMX: Content swapped for target', evt.detail.target.id);
    // You could hide a global loading indicator here
});

document.body.addEventListener('htmx:responseError', function(evt) {
    console.error('HTMX Error:', evt.detail.error);
    // evt.detail.xhr contains the XMLHttpRequest object
    // You could display a generic error message to the user
    const target = evt.detail.target || document.body;
    let errorMsg = "An error occurred. Please try again.";
    if (evt.detail.xhr && evt.detail.xhr.responseText) {
        // Attempt to parse error from response if backend sends structured errors
        try {
            const errData = JSON.parse(evt.detail.xhr.responseText);
            if (errData.error) errorMsg = errData.error;
        } catch (e) {
            // If not JSON, use a portion of the response text or default
            errorMsg = evt.detail.xhr.responseText.substring(0, 100) || errorMsg;
        }
    }
    // For simplicity, alerting. In a real app, you'd inject this into a dedicated error div.
    alert("HTMX Request Error: " + errorMsg);
});

// If you were to keep Alpine for very minor things, it would be initialized here
// document.addEventListener('alpine:init', () => {
//     Alpine.data('someSmallComponent', () => ({
//         isOpen: false,
//         toggle() { this.isOpen = !this.isOpen; }
//     }));
// });

console.log('Simple driver.js loaded. HTMX will handle most interactions.');
```


Workflow Example with Go SSR + HTMX (Loading Case Folders):

Initial Page Load: User navigates to /. Go handler ShowDocumentAutomationPage renders index.gohtml which includes _step0_case_setup.gohtml.
User Action: User clicks the "Load iCloud Root Folders" button in _step0_case_setup.gohtml.
HTMX Request: HTMX sees hx-get="/ui/icloud-folders" and makes a GET request to that URL.
Go Backend:
The GetICloudFoldersHTML handler is invoked.
It calls your services.GetRealICloudFolders() to get the folder data.
It prepares the data (gin.H{"Folders": folders}).
It renders the _icloud_folder_list.gohtml fragment with this data.
It sends the resulting HTML fragment (just the list of folders) back to the client.
HTMX Response Handling:
HTMX receives the HTML fragment.
It sees hx-target="#icloud-folder-list-container" and hx-swap="innerHTML".
It replaces the content of the div with ID icloud-folder-list-container with the received HTML fragment.
UI Update: The user now sees the list of iCloud root folders, without a full page reload.
This pattern would be repeated for subsequent interactions (selecting a parent folder loads case folders, selecting case folder loads documents, etc.).

Benefits of this Refactor:

Addresses the UI Defect: Rendering lists directly via Go templates and HTMX fragments bypasses complex client-side JavaScript reactivity for these core displays.
Robustness: More logic is in Go, which is strongly typed and often easier to test and maintain for backend-heavy teams. Less client-side JS means fewer browser compatibility or JS error issues breaking core functionality.
Simpler Debugging (for view logic): If a list isn't displaying correctly, you debug the Go handler and the Go template.
Clear Separation: Go handles business logic and view rendering; HTMX handles client-server communication for UI updates. driver.js becomes minimal.
This is a significant architectural shift but can lead to a more stable and maintainable application, especially if client-side JavaScript complexity is becoming a bottleneck. Remember to break down the refactor into manageable pieces, starting with one part of the workflow.