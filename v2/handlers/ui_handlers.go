package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"mallon-legal-v2/services"
	"github.com/gin-gonic/gin"
)

// UIHandlers contains the UI-related HTTP handlers
type UIHandlers struct {
	templates     *template.Template
	icloudService *services.ICloudService
	docService    *services.DocumentService
}

// PageData represents the data passed to templates
type PageData struct {
	CurrentStep          int
	Username             string
	ICloudConnected      bool
	SelectedParentFolder string
	SelectedCaseFolder   string
	CaseFolders          []services.ICloudDocument
	Documents            []services.ICloudDocument
	Templates            []services.Template
	Folders              []services.ICloudDocument
	ParentFolder         string
	Error                string
	RetryAction          string
	PreviewContent       PreviewDocument
	LegalAnalysis        LegalAnalysis
	DocumentHTML         template.HTML
	DocumentTitle        string
	DocumentFilename     string
	LastSaved            string
	ProcessingResult     *services.DocumentProcessingResult
	ClientCase           *services.ClientCase
	SelectedDocuments    []string
	
	// Session state for UI restoration
	SessionState         *services.WorkflowState
	IsReturningUser      bool
	SelectedTemplate     string
}

// LegalAnalysis represents the extracted legal information for Step 3
type LegalAnalysis struct {
	CauseOfAction  []CauseOfActionItem
	LegalViolations []LegalViolationItem
	SourceDocs     []string
	ExtractionDate string
}

// CauseOfActionItem represents a specific cause of action
type CauseOfActionItem struct {
	Title        string
	Description  string
	StatutoryBasis string
	SourceDoc    string
	Elements     []string
}

// LegalViolationItem represents a specific legal violation
type LegalViolationItem struct {
	Statute      string
	ViolationType string
	Description  string
	SourceDoc    string
	Penalties    string
}

// PreviewDocument represents a document with highlighted content
type PreviewDocument struct {
	Title         string
	Content       []PreviewSection
	SourceDocs    []string
	GeneratedDate string
}

// PreviewSection represents a section of the document with highlighted content
type PreviewSection struct {
	Title      string
	Content    string
	Highlights []HighlightedText
}

// HighlightedText represents text with source attribution
type HighlightedText struct {
	Text       string
	SourceDoc  string
	StartPos   int
	EndPos     int
}

// NewUIHandlers creates a new UI handlers instance
func NewUIHandlers() *UIHandlers {
	// Parse all templates
	tmpl := template.New("")
	
	// Add custom template functions
	tmpl.Funcs(template.FuncMap{
		"upper": strings.ToUpper,
		"formatSize": func(size int64) string {
			if size < 1024 {
				return fmt.Sprintf("%d B", size)
			} else if size < 1024*1024 {
				return fmt.Sprintf("%.1f KB", float64(size)/1024)
			} else {
				return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
			}
		},
		"ge": func(a, b int) bool { return a >= b },
		"eq": func(a, b int) bool { return a == b },
		"contains": func(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		"stringEq": func(a, b string) bool {
			return a == b
		},
		"add": func(a, b int) int {
			return a + b
		},
	})
	
	// Parse template files
	tmpl = template.Must(tmpl.ParseGlob("templates/*.gohtml"))
	
	return &UIHandlers{
		templates:     tmpl,
		icloudService: services.NewICloudService(),
		docService:    services.NewDocumentService(),
	}
}

// Session helper functions
func (h *UIHandlers) getSessionService(c *gin.Context) *services.SessionService {
	sessionService, exists := c.Get("sessionService")
	if !exists {
		log.Printf("Session service not found in context")
		return nil
	}
	return sessionService.(*services.SessionService)
}

func (h *UIHandlers) getSessionID(c *gin.Context) string {
	sessionID, exists := c.Get("sessionID")
	if !exists {
		log.Printf("Session ID not found in context")
		return "default_session"
	}
	return sessionID.(string)
}

func (h *UIHandlers) getWorkflowState(c *gin.Context) *services.WorkflowState {
	sessionService := h.getSessionService(c)
	if sessionService == nil {
		return &services.WorkflowState{CurrentStep: 0}
	}
	
	sessionID := h.getSessionID(c)
	return sessionService.GetSession(sessionID)
}

func (h *UIHandlers) updateWorkflowState(c *gin.Context, updateFunc func(*services.WorkflowState)) {
	sessionService := h.getSessionService(c)
	if sessionService == nil {
		return
	}
	
	sessionID := h.getSessionID(c)
	sessionService.UpdateSession(sessionID, updateFunc)
}

// ShowMainPage renders the main application page
func (h *UIHandlers) ShowMainPage(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Get current session state
	state := h.getWorkflowState(c)
	
	data := PageData{
		CurrentStep:          state.CurrentStep,
		Username:             username,
		ICloudConnected:      state.ICloudConnected,
		SelectedParentFolder: state.SelectedParentFolder,
		SelectedCaseFolder:   state.SelectedCaseFolder,
		SelectedDocuments:    state.SelectedDocuments,
		SelectedTemplate:     state.SelectedTemplate,
		SessionState:         state,
		IsReturningUser:      state.CurrentStep > 0,
	}
	
	err := h.templates.ExecuteTemplate(c.Writer, "index.gohtml", data)
	if err != nil {
		log.Printf("Error executing template index.gohtml: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering page")
	}
}

// GetStep renders a specific step
func (h *UIHandlers) GetStep(c *gin.Context) {
	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		step = 0
	}
	
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Get current session state
	state := h.getWorkflowState(c)
	
	// Validate session exists and handle missing sessions gracefully
	if state == nil {
		log.Printf("[WARNING] No session state found for GetStep request to step %d", step)
		// Create new state but indicate session was missing
		state = &services.WorkflowState{
			CurrentStep: 0,
			LastUpdated: time.Now(),
		}
		// Redirect to step 0 if trying to access later steps without session
		if step > 0 {
			log.Printf("[INFO] Redirecting to step 0 due to missing session state")
			c.Header("HX-Redirect", "/ui/step/0")
			c.String(http.StatusOK, "")
			return
		}
	}
	
	// Log session state for debugging
	log.Printf("[DEBUG] GetStep - Step: %d, Session CurrentStep: %d, HasCaseFolder: %v, HasDocuments: %d, HasTemplate: %v", 
		step, state.CurrentStep, state.SelectedCaseFolder != "", len(state.SelectedDocuments), state.SelectedTemplate != "")
	
	// Check if iCloud is connected (from query parameter or session)
	icloudConnected := c.Query("icloud_connected") == "true" || state.ICloudConnected
	
	// Determine if user is returning to a step they've been to before
	isReturningUser := state.CurrentStep > step
	
	// Update current step in session if moving forward
	if step > state.CurrentStep {
		h.updateWorkflowState(c, func(s *services.WorkflowState) {
			s.CurrentStep = step
		})
	}
	
	data := PageData{
		CurrentStep:          step,
		Username:             username,
		ICloudConnected:      icloudConnected,
		SelectedParentFolder: state.SelectedParentFolder,
		SelectedCaseFolder:   state.SelectedCaseFolder,
		SelectedDocuments:    state.SelectedDocuments,
		SelectedTemplate:     state.SelectedTemplate,
		SessionState:         state,
		IsReturningUser:      isReturningUser,
	}
	
	// Add step-specific data
	switch step {
	case 0:
		// If iCloud is connected, load root folders automatically
		if icloudConnected {
			folders, err := h.icloudService.GetRootFolders("", "")
			if err == nil {
				data.Folders = folders
				log.Printf("Loaded %d iCloud folders for connected user", len(folders))
			} else {
				log.Printf("[ERROR] Failed to load iCloud folders: %v", err)
				data.Error = "Could not load iCloud folders. Please try again."
			}
		}
		// If parent folder is selected, load case folders
		if state.SelectedParentFolder != "" {
			caseFolders, err := h.icloudService.GetSubfolders("", "", state.SelectedParentFolder)
			if err == nil {
				data.CaseFolders = caseFolders
				log.Printf("Loaded %d case folders from session state", len(caseFolders))
			} else {
				log.Printf("[ERROR] Failed to load case folders from %s: %v", state.SelectedParentFolder, err)
			}
		}
	case 1:
		// Validate prerequisites for step 1
		if state.SelectedCaseFolder == "" && !isReturningUser {
			log.Printf("[WARNING] Attempting to access step 1 without case folder selection")
			data.Error = "Please select a case folder first."
		}
		
		// Load documents for step 1 - prioritize from session state
		if state.SelectedCaseFolder != "" {
			documents, err := h.icloudService.GetDocuments("", "", state.SelectedCaseFolder)
			if err == nil {
				data.Documents = documents
				log.Printf("Loaded %d documents from session case folder: %s", len(documents), state.SelectedCaseFolder)
			} else {
				log.Printf("[WARNING] Failed to load documents from case folder %s: %v", state.SelectedCaseFolder, err)
				// Fallback to default loading
				documents, err := h.loadDocumentsForStep1(c)
				if err != nil {
					log.Printf("[ERROR] Fallback document loading failed: %v", err)
					data.Error = "Could not load documents. Please check your case folder selection."
				} else {
					data.Documents = documents
				}
			}
		} else {
			// Load documents using default method
			documents, err := h.loadDocumentsForStep1(c)
			if err != nil {
				log.Printf("[ERROR] Default document loading failed: %v", err)
				data.Error = "Could not load documents. Please select a case folder."
			} else {
				data.Documents = documents
			}
		}
	case 2:
		// Validate prerequisites for step 2
		if len(state.SelectedDocuments) == 0 && !isReturningUser {
			log.Printf("[WARNING] Attempting to access step 2 without document selection")
			data.Error = "Please select documents first."
		}
		
		// Load templates for step 2
		templates, err := h.docService.GetTemplates()
		if err != nil {
			log.Printf("[ERROR] Failed to load templates: %v", err)
			data.Error = "Could not load templates. Please try again."
		} else {
			data.Templates = templates
			log.Printf("Loaded %d templates for step 2", len(templates))
		}
	case 3:
		// Load available documents for Missing Content analysis early
		if state.SelectedCaseFolder != "" {
			documents, err := h.icloudService.GetDocuments("", "", state.SelectedCaseFolder)
			if err != nil {
				log.Printf("[ERROR] Failed to load documents for Missing Content analysis: %v", err)
			} else {
				data.Documents = documents
				log.Printf("[DEBUG] Loaded %d available documents for Missing Content analysis", len(documents))
			}
		}
		
		// Validate prerequisites for step 3
		if state.SelectedTemplate == "" && !isReturningUser {
			log.Printf("[WARNING] Attempting to access step 3 without template selection")
			data.Error = "Please select a template first."
			// Redirect back to step 2 with templates loaded
			templates, err := h.docService.GetTemplates()
			if err != nil {
				log.Printf("[ERROR] Failed to load templates: %v", err)
				templates = []services.Template{}
			}
			data.CurrentStep = 2
			data.Templates = templates
			break
		}
		
		// Check if we have processing results in session
		if state.ProcessingResult == nil || state.ClientCase == nil {
			log.Printf("[WARNING] Missing processing results in session for step 3")
			
			// If we have selected documents and template, try to reprocess
			if len(state.SelectedDocuments) > 0 && state.SelectedTemplate != "" {
				log.Printf("[INFO] Reprocessing documents for step 3")
				processingResult, clientCase, err := h.docService.ProcessSelectedDocuments(state.SelectedDocuments, state.SelectedTemplate)
				if err != nil {
					log.Printf("[ERROR] Failed to reprocess documents: %v", err)
					data.Error = "Failed to process documents. Please try again."
					// Redirect back to step 2
					templates, _ := h.docService.GetTemplates()
					data.CurrentStep = 2
					data.Templates = templates
					break
				}
				
				// Save the reprocessed results to session
				h.updateWorkflowState(c, func(state *services.WorkflowState) {
					state.ProcessingResult = processingResult
					state.ClientCase = clientCase
				})
				
				data.ProcessingResult = processingResult
				data.ClientCase = clientCase
			} else {
				// No way to recover, redirect to step 2
				log.Printf("[ERROR] Cannot proceed to step 3 without document processing")
				data.Error = "Session data missing. Please select your documents and template again."
				templates, _ := h.docService.GetTemplates()
				data.CurrentStep = 2
				data.Templates = templates
				break
			}
		} else {
			// Load from session
			data.ProcessingResult = state.ProcessingResult
			data.ClientCase = state.ClientCase
		}
		
		// Load legal analysis for step 3 with selected documents
		// Extract document names from paths for legal analysis
		selectedDocNames := make([]string, len(state.SelectedDocuments))
		for i, docPath := range state.SelectedDocuments {
			parts := strings.Split(docPath, "/")
			selectedDocNames[i] = parts[len(parts)-1]
		}
		legalAnalysis := h.generateLegalAnalysis(selectedDocNames)
		data.LegalAnalysis = legalAnalysis
		
		// Ensure we have selected documents list
		if len(state.SelectedDocuments) > 0 {
			data.SelectedDocuments = state.SelectedDocuments
		}
	}
	
	err = h.templates.ExecuteTemplate(c.Writer, "_step_wrapper.gohtml", data)
	if err != nil {
		log.Printf("Error executing template _step_wrapper.gohtml: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering step")
	}
}

// GetICloudFolders handles HTMX request for iCloud parent folders
func (h *UIHandlers) GetICloudFolders(c *gin.Context) {
	folders, err := h.icloudService.GetRootFolders("", "")
	if err != nil {
		// Return an HTML fragment indicating error
		data := PageData{
			Error:       "Could not load iCloud folders: " + err.Error(),
			RetryAction: "/ui/icloud-folders",
		}
		h.templates.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", data)
		return
	}

	data := PageData{
		Folders: folders,
	}
	h.templates.ExecuteTemplate(c.Writer, "_icloud_folder_list.gohtml", data)
}

// GetCaseFolders handles HTMX request for case subfolders
func (h *UIHandlers) GetCaseFolders(c *gin.Context) {
	parentFolder := c.Query("parent")
	if parentFolder == "" {
		data := PageData{
			Error: "Parent folder parameter required",
		}
		h.templates.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", data)
		return
	}
	
	caseFolders, err := h.icloudService.GetSubfolders("", "", parentFolder)
	if err != nil {
		log.Printf("Error accessing case folders: %v", err)
		data := PageData{
			Error:       "Could not load case folders: " + err.Error(),
			RetryAction: fmt.Sprintf("/ui/case-folders?parent=%s", parentFolder),
		}
		h.templates.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", data)
		return
	}
	
	data := PageData{
		CaseFolders:  caseFolders,
		ParentFolder: parentFolder,
	}
	h.templates.ExecuteTemplate(c.Writer, "_case_folder_list.gohtml", data)
}

// HandleICloudAuth handles iCloud authentication
func (h *UIHandlers) HandleICloudAuth(c *gin.Context) {
	username := c.PostForm("username")
	appPassword := c.PostForm("appPassword")
	
	if username == "" || appPassword == "" {
		// Return error modal
		data := PageData{
			Error: "Please enter both username and app password",
		}
		h.templates.ExecuteTemplate(c.Writer, "_icloud_auth_error.gohtml", data)
		return
	}
	
	// For prototype: simulate authentication
	// In production: implement actual iCloud API authentication
	log.Printf("iCloud auth attempt for user: %s", username)
	
	// Store iCloud connection state in session
	h.updateWorkflowState(c, func(state *services.WorkflowState) {
		state.ICloudConnected = true
		state.ICloudUsername = username
	})
	
	// Return success modal that will trigger a page refresh to Step 0 with connected state
	data := PageData{
		Username:        username,
		ICloudConnected: true,
	}
	
	// Return success modal that redirects to Step 0 with iCloud connected
	h.templates.ExecuteTemplate(c.Writer, "_icloud_auth_success.gohtml", data)
}

// ShowICloudSetup handles the iCloud setup modal
func (h *UIHandlers) ShowICloudSetup(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	data := PageData{
		Username:        username,
		ICloudConnected: false, // Would check actual session state
	}
	
	h.templates.ExecuteTemplate(c.Writer, "_icloud_setup_modal.gohtml", data)
}

// SelectParentFolder handles selecting an iCloud parent folder
func (h *UIHandlers) SelectParentFolder(c *gin.Context) {
	folderPath := c.PostForm("folderPath")
	if folderPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder path required"})
		return
	}
	
	// Store parent folder in session
	h.updateWorkflowState(c, func(state *services.WorkflowState) {
		state.SelectedParentFolder = folderPath
	})
	
	log.Printf("Selected parent folder: %s", folderPath)
	
	// Load case folders for the selected parent
	caseFolders, err := h.icloudService.GetSubfolders("", "", folderPath)
	if err != nil {
		log.Printf("Error loading case folders: %v", err)
		caseFolders = []services.ICloudDocument{} // Empty slice on error
	}
	
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Get current session state
	state := h.getWorkflowState(c)
	
	data := PageData{
		CurrentStep:          0,
		Username:             username,
		ICloudConnected:      true, // Assume connected if we got here
		SelectedParentFolder: folderPath,
		CaseFolders:          caseFolders,
		SessionState:         state,
		IsReturningUser:      state.CurrentStep > 0,
	}
	
	h.templates.ExecuteTemplate(c.Writer, "_step0_case_setup.gohtml", data)
}

// SelectCaseFolder handles selecting a specific case folder and navigates directly to Step 1
func (h *UIHandlers) SelectCaseFolder(c *gin.Context) {
	caseFolder := c.PostForm("caseFolder")
	if caseFolder == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Case folder required"})
		return
	}
	
	// Save to session state
	h.updateWorkflowState(c, func(state *services.WorkflowState) {
		state.SelectedCaseFolder = caseFolder
		state.CurrentStep = 1 // Move to document selection
	})
	
	log.Printf("Selected case folder: %s", caseFolder)
	
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Load documents from the selected case folder - only from iCloud
	documents, err := h.icloudService.GetDocuments("", "", caseFolder)
	if err != nil {
		log.Printf("Error loading documents from case folder %s: %v", caseFolder, err)
		// Return error without fallback - user must have valid iCloud access
		data := PageData{
			CurrentStep:        1,
			Username:           username,
			ICloudConnected:    true,
			SelectedCaseFolder: caseFolder,
			Error:             fmt.Sprintf("Could not load documents from case folder '%s'. Please check your iCloud Drive connection and ensure the folder exists.", caseFolder),
		}
		h.templates.ExecuteTemplate(c.Writer, "_step_wrapper.gohtml", data)
		return
	}
	
	// Log document loading results for debugging
	log.Printf("Successfully loaded %d documents from case folder %s", len(documents), caseFolder)
	if len(documents) == 0 {
		log.Printf("[INFO] Case folder %s is empty or contains no readable documents", caseFolder)
	}
	
	data := PageData{
		CurrentStep:        1,
		Username:           username,
		ICloudConnected:    true,
		SelectedCaseFolder: caseFolder,
		Documents:          documents, // This could be empty slice, which is fine
	}
	
	// Navigate directly to Step 1 with the loaded documents (or empty list)
	log.Printf("[DEBUG] Rendering step wrapper template for step %d with %d documents", data.CurrentStep, len(data.Documents))
	log.Printf("[DEBUG] Case folder: %s, ICloud connected: %t", data.SelectedCaseFolder, data.ICloudConnected)
	
	// Set content type to ensure proper HTML rendering
	c.Header("Content-Type", "text/html; charset=utf-8")
	
	// Capture the template output for debugging
	var templateBuffer strings.Builder
	err = h.templates.ExecuteTemplate(&templateBuffer, "_step_wrapper.gohtml", data)
	if err != nil {
		log.Printf("Error executing step wrapper template: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering page")
		return
	}
	
	output := templateBuffer.String()
	log.Printf("[DEBUG] ========== SERVER RESPONSE DEBUG ==========")
	log.Printf("[DEBUG] Template output length: %d characters", len(output))
	previewLen := 300
	if len(output) < previewLen {
		previewLen = len(output)
	}
	log.Printf("[DEBUG] Template output first 300 chars: %s", output[:previewLen])
	if len(output) > 300 {
		log.Printf("[DEBUG] Template output last 100 chars: %s", output[len(output)-100:])
	}
	log.Printf("[DEBUG] ========== END SERVER DEBUG ==========")
	
	// Write the captured output
	c.Writer.WriteString(output)
	
	log.Printf("[DEBUG] Successfully rendered step wrapper template")
}

// LoadDocuments handles loading documents from a folder
func (h *UIHandlers) LoadDocuments(c *gin.Context) {
	folder := c.Query("folder")
	if folder == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder parameter required"})
		return
	}
	
	documents, err := h.icloudService.GetDocuments("", "", folder)
	if err != nil {
		log.Printf("Error loading documents from %s: %v", folder, err)
		documents = []services.ICloudDocument{} // Empty slice on error
	}
	
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	data := PageData{
		CurrentStep:        1,
		Username:           username,
		ICloudConnected:    true,
		SelectedCaseFolder: folder,
		Documents:          documents,
	}
	
	h.templates.ExecuteTemplate(c.Writer, "_step1_document_selection.gohtml", data)
}

// SelectDocuments handles the document selection form submission
func (h *UIHandlers) SelectDocuments(c *gin.Context) {
	selectedDocs := c.PostFormArray("selectedDocs")
	caseFolder := c.PostForm("caseFolder")
	
	log.Printf("Selected documents: %v from folder: %s", selectedDocs, caseFolder)
	
	if len(selectedDocs) == 0 {
		// Return error or redirect back to step 1
		data := PageData{
			Error: "Please select at least one document",
		}
		h.templates.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", data)
		return
	}
	
	// Save selected documents to session state
	h.updateWorkflowState(c, func(state *services.WorkflowState) {
		state.SelectedDocuments = selectedDocs
		state.CurrentStep = 2 // Move to template selection
	})
	// For now, proceed to step 2 (template selection)
	
	templates, err := h.docService.GetTemplates()
	if err != nil {
		log.Printf("Error loading templates: %v", err)
		templates = []services.Template{}
	}
	
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	data := PageData{
		CurrentStep:        2,
		Username:           username,
		ICloudConnected:    true,
		SelectedCaseFolder: caseFolder,
		Templates:          templates,
		SelectedDocuments:  selectedDocs, // Store selected documents for next step
	}
	
	h.templates.ExecuteTemplate(c.Writer, "_step_wrapper.gohtml", data)
}

// SelectTemplate handles the template selection form submission
func (h *UIHandlers) SelectTemplate(c *gin.Context) {
	selectedTemplate := c.PostForm("selectedTemplate")
	selectedDocs := c.PostFormArray("selectedDocs") // Get selected documents from form
	
	log.Printf("Selected template: %s", selectedTemplate)
	log.Printf("Processing selected documents: %v", selectedDocs)
	
	// Save template selection to session state
	h.updateWorkflowState(c, func(state *services.WorkflowState) {
		state.SelectedTemplate = selectedTemplate
		state.SelectedDocuments = selectedDocs // Save selected documents
		state.CurrentStep = 3 // Move to review data
		// Clear any previous processing results when selecting new template
		state.ProcessingResult = nil
		state.ClientCase = nil
	})
	
	if selectedTemplate == "" {
		// Return error
		data := PageData{
			Error: "Please select a template",
		}
		h.templates.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", data)
		return
	}
	
	if len(selectedDocs) == 0 {
		// Return to Step 2 with error message and templates loaded
		templates, err := h.docService.GetTemplates()
		if err != nil {
			log.Printf("Error loading templates: %v", err)
			templates = []services.Template{}
		}
		
		username := c.GetString("username")
		if username == "" {
			username = "User"
		}
		
		data := PageData{
			CurrentStep:       2,
			Username:         username,
			ICloudConnected:  true,
			Templates:        templates,
			Error:           "No documents selected for processing. Please go back to Step 1 and select documents.",
		}
		h.templates.ExecuteTemplate(c.Writer, "_step_wrapper.gohtml", data)
		return
	}
	
	// Process selected documents using the document service (Task 8 implementation)
	processingResult, clientCase, err := h.docService.ProcessSelectedDocuments(selectedDocs, selectedTemplate)
	if err != nil {
		log.Printf("Error processing selected documents: %v", err)
		data := PageData{
			Error: "Error processing documents: " + err.Error(),
		}
		h.templates.ExecuteTemplate(c.Writer, "_error_fragment.gohtml", data)
		return
	}
	
	// Save processing results to session state
	h.updateWorkflowState(c, func(state *services.WorkflowState) {
		state.ProcessingResult = processingResult
		state.ClientCase = clientCase
	})
	
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Generate legal analysis for Step 3 with selected documents
	// Extract document names from selected document paths
	selectedDocNames := make([]string, len(selectedDocs))
	for i, docPath := range selectedDocs {
		parts := strings.Split(docPath, "/")
		selectedDocNames[i] = parts[len(parts)-1]
	}
	legalAnalysis := h.generateLegalAnalysis(selectedDocNames)
	
	data := PageData{
		CurrentStep:        3,
		Username:           username,
		ICloudConnected:    true,
		LegalAnalysis:      legalAnalysis,
		ProcessingResult:   processingResult,  // Add dynamic processing result
		ClientCase:         clientCase,        // Add dynamic client case data
		SelectedDocuments:  selectedDocs,      // Pass through selected documents
	}
	
	log.Printf("[DEBUG] SelectTemplate: Setting CurrentStep to %d with %d selected documents, %.1f%% coverage", 
		data.CurrentStep, len(selectedDocs), processingResult.DataCoverage)
	
	h.templates.ExecuteTemplate(c.Writer, "_step_wrapper.gohtml", data)
}

// PreviewDocument handles the document preview with highlighting
func (h *UIHandlers) PreviewDocument(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Get workflow state to access selected documents
	state := h.getWorkflowState(c)
	selectedDocs := state.SelectedDocuments
	if len(selectedDocs) == 0 {
		// Fallback to hardcoded for backwards compatibility during development
		selectedDocs = []string{"Attorney_Notes.txt", "Adverse_Action_Letter_Cap_One.pdf", "Civil_Cover_Sheet.txt", "Complaint_Final.docx"}
		log.Printf("[WARNING] No selected documents found in session, using fallback")
	}
	
	// Generate preview document with highlighted content using selected documents
	previewData := h.generatePreviewDocument(selectedDocs)
	log.Printf("[DEBUG] Generated preview document with %d sections for selected docs: %v", len(previewData.Content), selectedDocs)
	
	// Debug the content structure
	for i, section := range previewData.Content {
		log.Printf("[DEBUG] Section %d: Title='%s', Content length=%d", i, section.Title, len(section.Content))
	}
	
	data := PageData{
		Username:        username,
		ICloudConnected: true,
		PreviewContent:  previewData,
	}
	
	err := h.templates.ExecuteTemplate(c.Writer, "_document_preview.gohtml", data)
	if err != nil {
		log.Printf("[ERROR] Error executing template _document_preview.gohtml: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering document preview: "+err.Error())
		return
	}
	log.Printf("[SUCCESS] Document preview template executed successfully")
}

// generateLegalAnalysis creates legal analysis data from extracted information
func (h *UIHandlers) generateLegalAnalysis(selectedDocs []string) LegalAnalysis {
	// For prototype: simulate legal analysis based on FCRA credit card fraud case
	// In production: this would analyze actual extracted data from documents
	
	// Use actual selected documents, fallback to defaults if none provided
	sourceDocs := selectedDocs
	if len(sourceDocs) == 0 {
		sourceDocs = []string{
			"Attorney_Notes.txt",
			"Adverse_Action_Letter_Cap_One.pdf", 
			"Civil_Cover_Sheet.txt",
			"Complaint_Final.docx",
		}
	}
	
	return LegalAnalysis{
		ExtractionDate: "June 5, 2025",
		SourceDocs: sourceDocs,
		CauseOfAction: []CauseOfActionItem{
			{
				Title:          "Negligent Non-Compliance with FCRA",
				Description:    "Defendants negligently failed to follow reasonable procedures to assure maximum possible accuracy of consumer credit information",
				StatutoryBasis: "15 U.S.C. § 1681e(b)",
				SourceDoc:      "Complaint_Final.docx",
				Elements: []string{
					"Duty to maintain reasonable procedures",
					"Failure to assure maximum possible accuracy",
					"Reporting of inaccurate information",
					"Proximately caused damages to consumer",
				},
			},
			{
				Title:          "Willful Non-Compliance with FCRA",
				Description:    "Defendants willfully failed to conduct reasonable reinvestigation upon consumer dispute",
				StatutoryBasis: "15 U.S.C. § 1681i(a)",
				SourceDoc:      "Attorney_Notes.txt",
				Elements: []string{
					"Received consumer dispute",
					"Failed to conduct reasonable reinvestigation",
					"Willful or reckless disregard for consumer rights",
					"Continued reporting of disputed information",
				},
			},
			{
				Title:          "Failure to Provide Required Notices",
				Description:    "Defendants failed to provide adverse action notices as required by FCRA",
				StatutoryBasis: "15 U.S.C. § 1681m(a)",
				SourceDoc:      "Adverse_Action_Letter_Cap_One.pdf",
				Elements: []string{
					"Use of consumer report in adverse action",
					"Failure to provide timely notice",
					"Inadequate disclosure of consumer rights",
					"Damages from lack of notice",
				},
			},
		},
		LegalViolations: []LegalViolationItem{
			{
				Statute:       "15 U.S.C. § 1681e(b)",
				ViolationType: "Negligent Failure - Reasonable Procedures",
				Description:   "Consumer reporting agency failed to follow reasonable procedures to assure maximum possible accuracy of information concerning the consumer",
				SourceDoc:     "Complaint_Final.docx",
				Penalties:     "Actual damages, attorney fees, and costs",
			},
			{
				Statute:       "15 U.S.C. § 1681i(a)(1)(A)",
				ViolationType: "Willful Failure - Reinvestigation Duties",
				Description:   "Upon dispute, consumer reporting agency failed to conduct reasonable reinvestigation to determine whether the disputed information is inaccurate",
				SourceDoc:     "Attorney_Notes.txt",
				Penalties:     "Actual damages OR statutory damages $100-$1,000, plus attorney fees",
			},
			{
				Statute:       "15 U.S.C. § 1681i(a)(5)(A)",
				ViolationType: "Failure to Delete - Disputed Information",
				Description:   "Failed to promptly delete inaccurate or unverifiable information from consumer's file following dispute",
				SourceDoc:     "Attorney_Notes.txt",
				Penalties:     "Actual damages, attorney fees, and costs",
			},
			{
				Statute:       "15 U.S.C. § 1681c(a)(2)",
				ViolationType: "Reporting Prohibited Information",
				Description:   "Continued reporting of adverse account information beyond the permissible time periods",
				SourceDoc:     "Adverse_Action_Letter_Cap_One.pdf",
				Penalties:     "Actual damages, attorney fees, and costs",
			},
			{
				Statute:       "15 U.S.C. § 1681m(a)",
				ViolationType: "Adverse Action Notice Violations",
				Description:   "Failed to provide required adverse action notices with consumer reporting agency information",
				SourceDoc:     "Adverse_Action_Letter_Cap_One.pdf",
				Penalties:     "Actual damages, attorney fees, and costs",
			},
			{
				Statute:       "15 U.S.C. § 1681n",
				ViolationType: "Willful Non-Compliance - Civil Liability",
				Description:   "Pattern of willful non-compliance with FCRA requirements causing consumer harm",
				SourceDoc:     "Complaint_Final.docx",
				Penalties:     "Actual damages OR $100-$1,000 statutory damages, plus punitive damages and attorney fees",
			},
		},
	}
}

// generatePreviewDocument creates a complete legal document with highlighted content
func (h *UIHandlers) generatePreviewDocument(selectedDocs []string) PreviewDocument {
	// Generate complete FCRA legal complaint document using only selected documents
	// In production: this would process actual extracted data from documents
	
	return PreviewDocument{
		Title:         "COMPLAINT FOR VIOLATIONS OF THE FAIR CREDIT REPORTING ACT",
		GeneratedDate: time.Now().Format("January 2, 2006"),
		SourceDocs:    selectedDocs,
		Content: []PreviewSection{
			{
				Title:   "UNITED STATES DISTRICT COURT",
				Content: "EASTERN DISTRICT OF NEW YORK\n\nEMAN YOUSSEF,\n\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\tPlaintiff,\n\nv.\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\tCivil Action No. _______\n\nEQUIFAX INFORMATION SERVICES LLC,\nEXPERIAN INFORMATION SOLUTIONS INC.,\nTRANS UNION LLC,\n\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\tDefendants.\n\n____________________________________________________________________________",
			},
			{
				Title:   "COMPLAINT",
				Content: "Plaintiff Eman Youssef, by and through undersigned counsel, brings this action against Defendants for violations of the Fair Credit Reporting Act (\"FCRA\"), 15 U.S.C. § 1681 et seq., and alleges as follows:",
			},
			{
				Title:   "I. JURISDICTION AND VENUE",
				Content: "1. This Court has subject matter jurisdiction over this action pursuant to 15 U.S.C. § 1681p and 28 U.S.C. § 1331, as this action arises under federal law.\n\n2. Venue is proper in this District pursuant to 28 U.S.C. § 1391(b) because a substantial part of the events giving rise to the claims occurred in this judicial district, and Defendants conduct business in this District.\n\n3. This Court has personal jurisdiction over the Defendants because they conduct substantial business within this District and/or the acts giving rise to this lawsuit occurred within this District.",
			},
			{
				Title:   "II. PARTIES",
				Content: "4. Plaintiff Eman Youssef is an individual residing in Queens, New York. Plaintiff may be reached at 347.891.5584.\n\n5. Upon information and belief, Defendant EQUIFAX INFORMATION SERVICES LLC is a limited liability company organized and existing under the laws of Georgia, with its principal place of business located in Atlanta, Georgia. Equifax is a \"consumer reporting agency\" as that term is defined in 15 U.S.C. § 1681a(f).\n\n6. Upon information and belief, Defendant EXPERIAN INFORMATION SOLUTIONS INC. is a corporation organized and existing under the laws of Delaware, with its principal place of business located in Costa Mesa, California. Experian is a \"consumer reporting agency\" as that term is defined in 15 U.S.C. § 1681a(f).\n\n7. Upon information and belief, Defendant TRANS UNION LLC is a limited liability company organized and existing under the laws of Delaware, with its principal place of business located in Chicago, Illinois. Trans Union is a \"consumer reporting agency\" as that term is defined in 15 U.S.C. § 1681a(f).",
			},
			{
				Title:   "III. FACTUAL ALLEGATIONS",
				Content: "8. During the period from June 30, 2024 through July 30, 2024, Plaintiff was traveling in Egypt.\n\n9. While Plaintiff was traveling in Egypt, fraudulent charges totaling approximately $7,500 were made on Plaintiff's TD Bank credit card account.\n\n10. Upon discovering the fraudulent charges, Plaintiff immediately contacted TD Bank to report the unauthorized transactions and dispute the charges.\n\n11. Plaintiff filed a police report regarding the fraudulent transactions and provided all necessary documentation to TD Bank.\n\n12. Despite Plaintiff's timely notification and dispute of the fraudulent charges, the unauthorized accounts and/or adverse information related to these fraudulent transactions continue to appear on Plaintiff's consumer credit reports maintained by Defendants.\n\n13. The continued reporting of this fraudulent and inaccurate information has damaged Plaintiff's credit score and creditworthiness.\n\n14. As a result of Defendants' actions, Plaintiff has been denied credit and has suffered actual damages.",
			},
			{
				Title:   "IV. FIRST CAUSE OF ACTION",
				Content: "NEGLIGENT NON-COMPLIANCE WITH THE FCRA\n(15 U.S.C. § 1681e(b) and 15 U.S.C. § 1681o)\n\n15. Plaintiff incorporates by reference each and every allegation contained in the preceding paragraphs as if fully set forth herein.\n\n16. At all times relevant hereto, Defendants were \"consumer reporting agencies\" within the meaning of 15 U.S.C. § 1681a(f).\n\n17. Defendants owed a duty to Plaintiff to follow reasonable procedures to assure maximum possible accuracy of the information concerning Plaintiff in Plaintiff's consumer credit file.\n\n18. Defendants negligently violated this duty by failing to follow reasonable procedures to assure the maximum possible accuracy of the information in Plaintiff's credit file.\n\n19. As a direct and proximate result of Defendants' negligent violations of the FCRA, Plaintiff has suffered actual damages.",
			},
			{
				Title:   "V. SECOND CAUSE OF ACTION",
				Content: "WILLFUL NON-COMPLIANCE WITH THE FCRA\n(15 U.S.C. § 1681i(a) and 15 U.S.C. § 1681n)\n\n20. Plaintiff incorporates by reference each and every allegation contained in the preceding paragraphs as if fully set forth herein.\n\n21. Upon receiving notice of Plaintiff's dispute regarding the inaccurate information, Defendants were required to conduct a reasonable reinvestigation of the disputed information.\n\n22. Defendants willfully failed to conduct a reasonable reinvestigation as required by 15 U.S.C. § 1681i(a).\n\n23. Defendants' conduct was willful and in reckless disregard of Plaintiff's rights under the FCRA.\n\n24. As a direct and proximate result of Defendants' willful violations of the FCRA, Plaintiff has suffered actual damages and is entitled to statutory damages.",
			},
			{
				Title:   "VI. PRAYER FOR RELIEF",
				Content: "WHEREFORE, Plaintiff respectfully requests that this Court:\n\nA. Enter judgment in favor of Plaintiff and against Defendants;\n\nB. Award Plaintiff actual damages pursuant to 15 U.S.C. § 1681o and § 1681n;\n\nC. Award Plaintiff statutory damages in the amount of not less than $100 nor more than $1,000 for each willful violation pursuant to 15 U.S.C. § 1681n;\n\nD. Award Plaintiff punitive damages pursuant to 15 U.S.C. § 1681n;\n\nE. Award Plaintiff reasonable attorney's fees and costs pursuant to 15 U.S.C. § 1681o and § 1681n;\n\nF. Grant such other and further relief as this Court may deem just and proper.",
			},
			{
				Title:   "JURY DEMAND",
				Content: "Plaintiff hereby demands a trial by jury on all issues so triable.\n\n\nRespectfully submitted,\n\n_________________________\nKevin Mallon, Esq.\nAttorney for Plaintiff\nState Bar No. [Number]\n[Address]\n[Phone]\n[Email]",
			},
		},
	}
}

// ViewDocument handles the document viewing request
func (h *UIHandlers) ViewDocument(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Get client name from path param or default to Eman Youssef
	clientName := c.DefaultQuery("client", "Eman Youssef")
	clientNameLower := strings.ToLower(strings.Replace(clientName, " ", "_", -1))
	
	// Base document directory
	docDir := "/Users/corelogic/satori-dev/clients/proj-mallon/dev/saved_documents"
	
	// Ensure directory exists
	if _, err := os.Stat(docDir); os.IsNotExist(err) {
		log.Printf("[INFO] Creating saved_documents directory")
		if err := os.MkdirAll(docDir, 0755); err != nil {
			log.Printf("[ERROR] Failed to create saved_documents directory: %v", err)
			c.String(http.StatusInternalServerError, "Could not create document directory")
			return
		}
	}
	
	// Try to find the document in this priority order:
	// 1. Latest file (complaint_clientname_latest.html)
	// 2. Known timestamp file (from project history)
	// 3. Any file matching the pattern complaint_clientname_*.html
	// 4. Generate new document from preview
	
	// 1. Check for latest file
	latestPath := fmt.Sprintf("%s/complaint_%s_latest.html", docDir, clientNameLower)
	log.Printf("[INFO] Checking for latest document: %s", latestPath)
	
	var documentPath string
	var documentHTML []byte
	var err error
	
	if _, err := os.Stat(latestPath); err == nil {
		// Latest file exists
		documentPath = latestPath
		log.Printf("[INFO] Using latest document: %s", documentPath)
		
		// Read the document HTML
		documentHTML, err = os.ReadFile(documentPath)
		if err != nil {
			log.Printf("[ERROR] Error reading latest document file: %v", err)
			// Continue to next option
		}
	}
	
	if documentHTML == nil {
		// 2. Check for known timestamp file
		knownTimestampPath := fmt.Sprintf("%s/complaint_%s_20250605_010420.html", docDir, clientNameLower)
		if _, err := os.Stat(knownTimestampPath); err == nil {
			documentPath = knownTimestampPath
			log.Printf("[INFO] Using known timestamp document: %s", documentPath)
			
			// Read the document HTML
			documentHTML, err = os.ReadFile(documentPath)
			if err != nil {
				log.Printf("[ERROR] Error reading known timestamp document file: %v", err)
				// Continue to next option
			}
		}
	}
	
	if documentHTML == nil {
		// 3. Find any matching file
		pattern := fmt.Sprintf("%s/complaint_%s_*.html", docDir, clientNameLower)
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			// Use the most recent file (assuming timestamp naming)
			sort.Sort(sort.Reverse(sort.StringSlice(matches)))
			documentPath = matches[0]
			log.Printf("[INFO] Using most recent document: %s", documentPath)
			
			// Read the document HTML
			documentHTML, err = os.ReadFile(documentPath)
			if err != nil {
				log.Printf("[ERROR] Error reading matched document file: %v", err)
				// Continue to next option
			}
		}
	}
	
	if documentHTML == nil {
		// 4. Generate new document from preview content and save it
		log.Printf("[INFO] No existing document found, generating new document from preview content")
		
		// Generate preview content
		// Get selected documents from session state
		state := h.getWorkflowState(c)
		selectedDocs := state.SelectedDocuments
		if len(selectedDocs) == 0 {
			selectedDocs = []string{"Attorney_Notes.txt", "Adverse_Action_Letter_Cap_One.pdf", "Civil_Cover_Sheet.txt", "Complaint_Final.docx"}
		}
		previewContent := h.generatePreviewDocument(selectedDocs)
		
		// Create HTML document from preview content
		var legalDocHTML strings.Builder
		legalDocHTML.WriteString("<div class=\"legal-document\">")
		
		for _, section := range previewContent.Content {
			if section.Title != "" {
				legalDocHTML.WriteString(fmt.Sprintf("<div class=\"section-title\">%s</div>\n", section.Title))
			}
			legalDocHTML.WriteString(fmt.Sprintf("<div class=\"section-content\">%s</div>\n", section.Content))
		}
		
		legalDocHTML.WriteString("</div>")
		
		// Create timestamp for filename
		timestamp := time.Now().Format("20060102_150405")
		
		// Create paths for new document
		documentPath = fmt.Sprintf("%s/complaint_%s_%s.html", docDir, clientNameLower, timestamp)
		latestPath := fmt.Sprintf("%s/complaint_%s_latest.html", docDir, clientNameLower)
		
		// Create HTML document structure
		fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Legal Complaint - %s</title>
	<style>
		body { font-family: 'Times New Roman', serif; margin: 1in; line-height: 1.5; }
		.highlight { background-color: #fef08a; }
		.legal-document {
			font-family: Times New Roman, serif;
			line-height: 1.5;
			margin: 1in;
		}
		.header {
			text-align: center;
			margin-bottom: 24px;
		}
		.court-info {
			text-align: center;
			margin-bottom: 24px;
			text-transform: uppercase;
		}
		.case-info {
			text-align: center;
			margin-bottom: 24px;
		}
		.section-title {
			text-align: center;
			text-transform: uppercase;
			font-weight: bold;
			margin: 24px 0;
		}
		.paragraph {
			text-indent: 0.5in;
			margin-bottom: 12px;
		}
		.numbered-paragraph {
			margin-bottom: 12px;
		}
		.signature-block {
			margin-top: 48px;
		}
	</style>
</head>
<body>
	%s
</body>
</html>`, clientName, legalDocHTML.String())
		
		// Write the files
		err = os.WriteFile(documentPath, []byte(fullHTML), 0644)
		if err != nil {
			log.Printf("[ERROR] Error saving new document to %s: %v", documentPath, err)
			c.String(http.StatusInternalServerError, "Error creating new document: "+err.Error())
			return
		}
		
		err = os.WriteFile(latestPath, []byte(fullHTML), 0644)
		if err != nil {
			log.Printf("[WARNING] Error saving new document to latest path %s: %v", latestPath, err)
			// Don't fail the request if we can't save to the latest path
		}
		
		// Set document HTML and last saved time
		documentHTML = []byte(fullHTML)
		log.Printf("[SUCCESS] Generated and saved new document to %s", documentPath)
	}
	
	// Extract just the legal document div from the HTML if needed
	docContent := string(documentHTML)
	// Find the legal-document div
	legalDocStart := strings.Index(docContent, "<div class=\"legal-document\">") 
	legalDocEnd := strings.LastIndex(docContent, "</div>")
	
	// Extract just the legal document portion
	var legalDocHTML string
	if legalDocStart >= 0 && legalDocEnd > legalDocStart {
		legalDocHTML = docContent[legalDocStart:legalDocEnd+6] // +6 to include the closing div
		log.Printf("[INFO] Successfully extracted legal document div from HTML")
	} else {
		// If we can't extract just the legal doc div, use the whole document
		legalDocHTML = docContent
		log.Printf("[WARNING] Could not extract legal document div, using full document content")
	}
	
	// Create document filename for download
	timestamp := time.Now().Format("20060102_150405")
	documentFilename := fmt.Sprintf("complaint_%s_%s.html", clientNameLower, timestamp)
	
	data := PageData{
		Username:          username,
		ICloudConnected:   true,
		DocumentHTML:      template.HTML(legalDocHTML),
		DocumentTitle:     "Legal Complaint - " + clientName,
		DocumentFilename:  documentFilename,
	}
	
	log.Printf("[INFO] Rendering document viewer for %s", clientName)
	err = h.templates.ExecuteTemplate(c.Writer, "_document_viewer.gohtml", data)
	if err != nil {
		log.Printf("[ERROR] Error executing template _document_viewer.gohtml: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering document viewer: "+err.Error())
	}
}

// EditDocument handles the document editing request
func (h *UIHandlers) EditDocument(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		username = "User"
	}
	
	// Get client name from path param or default to Eman Youssef
	clientName := c.DefaultQuery("client", "Eman Youssef")
	clientNameLower := strings.ToLower(strings.Replace(clientName, " ", "_", -1))
	
	// Base document directory
	docDir := "/Users/corelogic/satori-dev/clients/proj-mallon/dev/saved_documents"
	
	// Ensure directory exists
	if _, err := os.Stat(docDir); os.IsNotExist(err) {
		log.Printf("[INFO] Creating saved_documents directory")
		if err := os.MkdirAll(docDir, 0755); err != nil {
			log.Printf("[ERROR] Failed to create saved_documents directory: %v", err)
			c.String(http.StatusInternalServerError, "Could not create document directory")
			return
		}
	}
	
	// Try to find the document in this priority order:
	// 1. Latest file (complaint_clientname_latest.html)
	// 2. Known timestamp file (from project history)
	// 3. Any file matching the pattern complaint_clientname_*.html
	// 4. Default fallback file
	// 5. Generate new document from preview
	
	// 1. Check for latest file
	latestPath := fmt.Sprintf("%s/complaint_%s_latest.html", docDir, clientNameLower)
	log.Printf("[INFO] Checking for latest document for editing: %s", latestPath)
	
	var documentPath string
	var lastSavedTime string
	var documentHTML []byte
	var err error
	
	if _, err := os.Stat(latestPath); err == nil {
		// Latest file exists
		documentPath = latestPath
		log.Printf("[INFO] Using latest document for editing: %s", documentPath)
		
		// Try to get file info for last modified time
		if fileInfo, err := os.Stat(documentPath); err == nil {
			lastSavedTime = fileInfo.ModTime().Format("3:04:05 PM")
		} else {
			lastSavedTime = time.Now().Format("3:04:05 PM")
		}
		
		// Read the document HTML
		documentHTML, err = os.ReadFile(documentPath)
		if err != nil {
			log.Printf("[ERROR] Error reading latest document file: %v", err)
			// Continue to next option
		}
	}
	
	if documentHTML == nil {
		// 2. Check for known timestamp file
		knownTimestampPath := fmt.Sprintf("%s/complaint_%s_20250605_010420.html", docDir, clientNameLower)
		if _, err := os.Stat(knownTimestampPath); err == nil {
			documentPath = knownTimestampPath
			log.Printf("[INFO] Using known timestamp document for editing: %s", documentPath)
			
			// Set last saved time from file timestamp
			lastSavedTime = "10:42 AM" // Hardcoded based on filename
			
			// Read the document HTML
			documentHTML, err = os.ReadFile(documentPath)
			if err != nil {
				log.Printf("[ERROR] Error reading known timestamp document file: %v", err)
				// Continue to next option
			}
		}
	}
	
	if documentHTML == nil {
		// 3. Find any matching file
		pattern := fmt.Sprintf("%s/complaint_%s_*.html", docDir, clientNameLower)
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			// Use the most recent file (assuming timestamp naming)
			sort.Sort(sort.Reverse(sort.StringSlice(matches)))
			documentPath = matches[0]
			log.Printf("[INFO] Using most recent document for editing: %s", documentPath)
			
			// Try to extract timestamp from filename
			baseFile := filepath.Base(documentPath)
			timestampParts := strings.Split(baseFile, "_")
			if len(timestampParts) > 2 {
				timestampStr := strings.TrimSuffix(timestampParts[len(timestampParts)-1], ".html")
				if t, err := time.Parse("20060102_150405", timestampStr); err == nil {
					lastSavedTime = t.Format("3:04:05 PM")
				} else {
					lastSavedTime = time.Now().Format("3:04:05 PM")
				}
			} else {
				lastSavedTime = time.Now().Format("3:04:05 PM")
			}
			
			// Read the document HTML
			documentHTML, err = os.ReadFile(documentPath)
			if err != nil {
				log.Printf("[ERROR] Error reading matched document file: %v", err)
				// Continue to next option
			}
		}
	}
	
	if documentHTML == nil {
		// 4. Generate new document from preview content and save it
		log.Printf("[INFO] No existing document found, generating new document from preview content")
		
		// Generate preview content
		// Get selected documents from session state
		state := h.getWorkflowState(c)
		selectedDocs := state.SelectedDocuments
		if len(selectedDocs) == 0 {
			selectedDocs = []string{"Attorney_Notes.txt", "Adverse_Action_Letter_Cap_One.pdf", "Civil_Cover_Sheet.txt", "Complaint_Final.docx"}
		}
		previewContent := h.generatePreviewDocument(selectedDocs)
		
		// Create HTML document from preview content
		var legalDocHTML strings.Builder
		legalDocHTML.WriteString("<div class=\"legal-document\">")
		
		for _, section := range previewContent.Content {
			if section.Title != "" {
				legalDocHTML.WriteString(fmt.Sprintf("<div class=\"section-title\">%s</div>\n", section.Title))
			}
			legalDocHTML.WriteString(fmt.Sprintf("<div class=\"section-content\">%s</div>\n", section.Content))
		}
		
		legalDocHTML.WriteString("</div>")
		
		// Create timestamp for filename
		timestamp := time.Now().Format("20060102_150405")
		
		// Create paths for new document
		documentPath = fmt.Sprintf("%s/complaint_%s_%s.html", docDir, clientNameLower, timestamp)
		latestPath := fmt.Sprintf("%s/complaint_%s_latest.html", docDir, clientNameLower)
		
		// Create HTML document structure
		fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Legal Complaint - %s</title>
	<style>
		body { font-family: 'Times New Roman', serif; margin: 1in; line-height: 1.5; }
		.highlight { background-color: #fef08a; }
		.legal-document {
			font-family: Times New Roman, serif;
			line-height: 1.5;
			margin: 1in;
		}
		.header {
			text-align: center;
			margin-bottom: 24px;
		}
		.court-info {
			text-align: center;
			margin-bottom: 24px;
			text-transform: uppercase;
		}
		.case-info {
			text-align: center;
			margin-bottom: 24px;
		}
		.section-title {
			text-align: center;
			text-transform: uppercase;
			font-weight: bold;
			margin: 24px 0;
		}
		.paragraph {
			text-indent: 0.5in;
			margin-bottom: 12px;
		}
		.numbered-paragraph {
			margin-bottom: 12px;
		}
		.signature-block {
			margin-top: 48px;
		}
	</style>
</head>
<body>
	%s
</body>
</html>`, clientName, legalDocHTML.String())
		
		// Write the files
		err = os.WriteFile(documentPath, []byte(fullHTML), 0644)
		if err != nil {
			log.Printf("[ERROR] Error saving new document to %s: %v", documentPath, err)
			c.String(http.StatusInternalServerError, "Error creating new document: "+err.Error())
			return
		}
		
		err = os.WriteFile(latestPath, []byte(fullHTML), 0644)
		if err != nil {
			log.Printf("[WARNING] Error saving new document to latest path %s: %v", latestPath, err)
			// Don't fail the request if we can't save to the latest path
		}
		
		// Set document HTML and last saved time
		documentHTML = []byte(legalDocHTML.String())
		lastSavedTime = "Just Now"
		log.Printf("[SUCCESS] Generated and saved new document to %s", documentPath)
	}
	
	// Extract just the legal document div from the HTML if needed
	docContent := string(documentHTML)
	// Find the legal-document div
	legalDocStart := strings.Index(docContent, "<div class=\"legal-document\">") 
	legalDocEnd := strings.LastIndex(docContent, "</div>")
	
	// Extract just the legal document portion
	var legalDocHTML string
	if legalDocStart >= 0 && legalDocEnd > legalDocStart {
		legalDocHTML = docContent[legalDocStart:legalDocEnd+6] // +6 to include the closing div
		log.Printf("[INFO] Successfully extracted legal document div from HTML for editing")
	} else {
		// If we can't extract just the legal doc div, use the whole document
		legalDocHTML = docContent
		log.Printf("[WARNING] Could not extract legal document div, using full document content for editing")
	}
	
	// Create document filename for download
	timestamp := time.Now().Format("20060102_150405")
	documentFilename := fmt.Sprintf("complaint_%s_%s.html", clientNameLower, timestamp)
	
	// Get selected documents from session state
	state := h.getWorkflowState(c)
	selectedDocs := state.SelectedDocuments
	if len(selectedDocs) == 0 {
		selectedDocs = []string{"Attorney_Notes.txt", "Adverse_Action_Letter_Cap_One.pdf", "Civil_Cover_Sheet.txt", "Complaint_Final.docx"}
	}
	
	// Generate preview content with source documents for the sidebar
	previewContent := h.generatePreviewDocument(selectedDocs)
	
	data := PageData{
		Username:          username,
		ICloudConnected:   true,
		DocumentHTML:      template.HTML(legalDocHTML),
		DocumentTitle:     "Legal Complaint - " + clientName,
		DocumentFilename:  documentFilename,
		PreviewContent:    previewContent,
		LastSaved:         lastSavedTime,
	}
	
	log.Printf("[INFO] Rendering document editor for %s with last saved time: %s", clientName, lastSavedTime)
	err = h.templates.ExecuteTemplate(c.Writer, "_document_editor.gohtml", data)
	if err != nil {
		log.Printf("[ERROR] Error executing template _document_editor.gohtml: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering document editor: "+err.Error())
	}
}

// SaveDocument handles saving the edited document
func (h *UIHandlers) SaveDocument(c *gin.Context) {
	// Define request struct
	type SaveRequest struct {
		Content     string `json:"content"`
		ClientName  string `json:"clientName"`
		DocumentType string `json:"documentType"`
	}
	
	// Parse request
	var req SaveRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("[ERROR] Error parsing save request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format: " + err.Error()})
		return
	}
	
	// Validate request
	if req.Content == "" {
		log.Printf("[ERROR] Empty document content in save request")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Empty document content"})
		return
	}
	
	if req.ClientName == "" {
		req.ClientName = "Eman Youssef" // Default client name
		log.Printf("[WARNING] Client name not provided, using default: %s", req.ClientName)
	}
	
	if req.DocumentType == "" {
		req.DocumentType = "complaint" // Default document type
		log.Printf("[WARNING] Document type not provided, using default: %s", req.DocumentType)
	}
	
	// Format client name for filename
	clientNameLower := strings.ToLower(strings.Replace(req.ClientName, " ", "_", -1))
	
	// Create timestamp for filename
	timestamp := time.Now().Format("20060102_150405")
	
	// Ensure saved_documents directory exists
	saveDir := "/Users/corelogic/satori-dev/clients/proj-mallon/dev/saved_documents"
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		log.Printf("[INFO] Creating saved_documents directory")
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			log.Printf("[ERROR] Failed to create saved_documents directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Could not create save directory"})
			return
		}
	}
	
	// Create path for new document
	documentPath := fmt.Sprintf("%s/%s_%s_%s.html", saveDir, req.DocumentType, clientNameLower, timestamp)
	log.Printf("[INFO] Saving document to: %s", documentPath)
	
	// Also save to a predictable path for easy lookup
	latestPath := fmt.Sprintf("%s/%s_%s_latest.html", saveDir, req.DocumentType, clientNameLower)
	
	// Create HTML document structure
	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Legal Complaint - %s</title>
	<style>
		body { font-family: 'Times New Roman', serif; margin: 1in; line-height: 1.5; }
		.highlight { background-color: #fef08a; }
		.legal-document {
			font-family: Times New Roman, serif;
			line-height: 1.5;
			margin: 1in;
		}
		.header {
			text-align: center;
			margin-bottom: 24px;
		}
		.court-info {
			text-align: center;
			margin-bottom: 24px;
			text-transform: uppercase;
		}
		.case-info {
			text-align: center;
			margin-bottom: 24px;
		}
		.section-title {
			text-align: center;
			text-transform: uppercase;
			font-weight: bold;
			margin: 24px 0;
		}
		.paragraph {
			text-indent: 0.5in;
			margin-bottom: 12px;
		}
		.numbered-paragraph {
			margin-bottom: 12px;
		}
		.signature-block {
			margin-top: 48px;
		}
	</style>
</head>
<body>
	<div class="legal-document">
		%s
	</div>
</body>
</html>`, req.ClientName, req.Content)
	
	// Write the file with timestamp
	err := os.WriteFile(documentPath, []byte(fullHTML), 0644)
	if err != nil {
		log.Printf("[ERROR] Error saving document to %s: %v", documentPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error saving document: " + err.Error()})
		return
	}
	
	// Also write to the latest path for easy access
	err = os.WriteFile(latestPath, []byte(fullHTML), 0644)
	if err != nil {
		log.Printf("[WARNING] Error saving document to latest path %s: %v", latestPath, err)
		// Don't fail the request if we can't save to the latest path
	}
	
	log.Printf("[SUCCESS] Document successfully saved to %s", documentPath)
	
	// Return success response with path
	c.JSON(http.StatusOK, gin.H{
		"success": true, 
		"path": documentPath,
		"latest_path": latestPath,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// Helper function to load documents for step 1
func (h *UIHandlers) loadDocumentsForStep1(c *gin.Context) ([]services.ICloudDocument, error) {
	// Get session state to check for selected case folder
	state := h.getWorkflowState(c)
	
	// Use selected case folder if available, otherwise require user to select one
	if state.SelectedCaseFolder == "" {
		return nil, fmt.Errorf("no case folder selected - please select a case folder first")
	}
	
	// Load documents only from iCloud - no test folder or backend fallback
	documents, err := h.icloudService.GetDocuments("", "", state.SelectedCaseFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to load documents from iCloud case folder %s: %v", state.SelectedCaseFolder, err)
	}
	
	return documents, nil
}