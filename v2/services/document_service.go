package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Document represents a legal document in the system
type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // "pdf", "docx", etc.
	Path        string `json:"path"`
	ContentType string `json:"contentType"` // "attorney_notes", "adverse_action", etc.
	Size        int64  `json:"size"`
}

// Template represents a legal template
type Template struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Path string `json:"path"`
}

// MissingContent represents a missing required field
type MissingContent struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Required    bool   `json:"required"`
}

// DocumentProcessingResult contains the results of processing selected documents
type DocumentProcessingResult struct {
	SelectedDocuments []Document             `json:"selectedDocuments"`
	ExtractedData     map[string]interface{} `json:"extractedData"`
	MissingContent    []MissingContent       `json:"missingContent"`
	DataCoverage      float64                `json:"dataCoverage"` // Percentage of data extracted
}

// ClientCase represents the structured data extracted from documents
type ClientCase struct {
	// Basic client information
	ClientName               string    `json:"clientName"`
	ContactInfo              string    `json:"contactInfo"`
	ResidenceLocation        string    `json:"residenceLocation"`
	
	// Court and case information
	CourtJurisdiction        string    `json:"courtJurisdiction"`
	CaseNumber               string    `json:"caseNumber"`
	
	// Financial institution information
	FinancialInstitution     string    `json:"financialInstitution"`
	AccountOpenDate          time.Time `json:"accountOpenDate"`
	CreditLimit              string    `json:"creditLimit"`
	
	// Travel information
	TravelLocation           string    `json:"travelLocation"`
	TravelStartDate          time.Time `json:"travelStartDate"`
	TravelEndDate            time.Time `json:"travelEndDate"`
	
	// Fraud information (legacy fields for compatibility)
	FraudAmount              string    `json:"fraudAmount"`
	FraudStartDate           time.Time `json:"fraudStartDate"`
	FraudEndDate             time.Time `json:"fraudEndDate"`
	FraudDetails             string    `json:"fraudDetails"`
	
	// Enhanced fraud information
	FraudDetailsStructured   []FraudDetail `json:"fraudDetailsStructured"`
	
	// Discovery and dispute information
	DiscoveryDate            time.Time `json:"discoveryDate"`
	DisputeCount             int       `json:"disputeCount"`
	DisputeMethods           []string  `json:"disputeMethods"`
	BankResponse             string    `json:"bankResponse"`
	
	// Police report information
	PoliceReportFiled        bool      `json:"policeReportFiled"`
	PoliceReportDetails      string    `json:"policeReportDetails"`
	
	// Credit bureau information (legacy)
	CreditBureauDisputes     []string  `json:"creditBureauDisputes"`
	CreditBureauDisputeDate  time.Time `json:"creditBureauDisputeDate"`
	
	// Enhanced credit bureau information
	CreditBureauInteractions []CreditBureauInteraction `json:"creditBureauInteractions"`
	
	// Legal case information
	Defendants               []Defendant `json:"defendants"`
	EstimatedDamages         float64     `json:"estimatedDamages"`
	
	// Additional evidence and impact
	AdditionalEvidence       string    `json:"additionalEvidence"`
	CreditImpact             string    `json:"creditImpact"`
}

// DocumentService handles document operations
type DocumentService struct {
	documentsDir        string
	testCasesDir        string
	extractor           *DocumentExtractor
	contentAnalyzer     *ContentAnalyzer
	attorneyNotesAnalyzer *AttorneyNotesAnalyzer
	templateEngine      *TemplateEngine
	extractionPatterns  map[string]interface{}
}

// NewDocumentService creates a new document service instance
func NewDocumentService() *DocumentService {
	service := &DocumentService{
		documentsDir: "/Users/corelogic/satori-dev/clients/proj-mallon/artifacts",
		testCasesDir: "/Users/corelogic/satori-dev/clients/proj-mallon/test_icloud/CASES",
		extractor:    NewDocumentExtractor(),
	}
	
	// Initialize content analyzer
	analyzer, err := NewContentAnalyzer()
	if err != nil {
		log.Printf("[DOCUMENT_SERVICE] Warning: Could not initialize content analyzer: %v", err)
		// Continue without analyzer for now
	} else {
		service.contentAnalyzer = analyzer
	}
	
	// Initialize attorney notes analyzer
	attorneyAnalyzer, err := NewAttorneyNotesAnalyzer()
	if err != nil {
		log.Printf("[DOCUMENT_SERVICE] Warning: Could not initialize attorney notes analyzer: %v", err)
		// Continue without attorney analyzer for now
	} else {
		service.attorneyNotesAnalyzer = attorneyAnalyzer
		log.Printf("[DOCUMENT_SERVICE] Initialized with Attorney Notes Intelligence Engine")
	}
	
	// Initialize template engine
	service.templateEngine = NewTemplateEngine()
	log.Printf("[DOCUMENT_SERVICE] Initialized with dynamic template engine")
	
	// Load extraction patterns
	service.loadExtractionPatterns()
	
	return service
}

// loadExtractionPatterns loads the JSON patterns for data extraction
func (s *DocumentService) loadExtractionPatterns() {
	patternsPath := "/Users/corelogic/satori-dev/clients/proj-mallon/v2/config/extraction_patterns.json"
	
	data, err := os.ReadFile(patternsPath)
	if err != nil {
		log.Printf("[DOCUMENT_SERVICE] Warning: Could not load extraction patterns: %v", err)
		s.extractionPatterns = make(map[string]interface{})
		return
	}
	
	err = json.Unmarshal(data, &s.extractionPatterns)
	if err != nil {
		log.Printf("[DOCUMENT_SERVICE] Warning: Could not parse extraction patterns: %v", err)
		s.extractionPatterns = make(map[string]interface{})
		return
	}
	
	log.Printf("[DOCUMENT_SERVICE] Loaded extraction patterns successfully")
}

// GetDocuments returns all available documents
func (s *DocumentService) GetDocuments() ([]Document, error) {
	files, err := os.ReadDir(s.documentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read documents directory: %v", err)
	}
	
	var documents []Document
	for i, file := range files {
		if file.IsDir() {
			continue
		}
		
		info, err := file.Info()
		if err != nil {
			continue
		}
		
		ext := strings.ToLower(filepath.Ext(file.Name()))
		docType := "unknown"
		contentType := "unknown"
		
		// Determine document type based on extension
		switch ext {
		case ".pdf":
			docType = "pdf"
			if strings.Contains(file.Name(), "Adverse") {
				contentType = "adverse_action"
			} else if strings.Contains(file.Name(), "SummonsEquifax") {
				contentType = "summons_equifax"
			} else if strings.Contains(file.Name(), "Summons") {
				contentType = "summons"
			} else if strings.Contains(file.Name(), "Civil Cover") {
				contentType = "civil_cover_sheet"
			}
		case ".docx":
			docType = "docx"
			if strings.Contains(file.Name(), "Atty_Notes") {
				contentType = "attorney_notes"
			} else if strings.Contains(file.Name(), "Complaint") {
				contentType = "complaint_form"
			}
		}
		
		// Create document object
		doc := Document{
			ID:          fmt.Sprintf("doc_%d", i+1),
			Name:        file.Name(),
			Type:        docType,
			Path:        filepath.Join(s.documentsDir, file.Name()),
			ContentType: contentType,
			Size:        info.Size(),
		}
		
		documents = append(documents, doc)
	}
	
	return documents, nil
}

// GetTemplates returns all available templates
func (s *DocumentService) GetTemplates() ([]Template, error) {
	// For the prototype, we'll return a predefined list of templates
	templates := []Template{
		{
			ID:   "fcra-credit-card-fraud",
			Name: "FCRA Complaint - Credit Card Fraud",
			Desc: "For cases involving fraudulent credit card transactions",
			Path: filepath.Join(s.documentsDir, "Complaint_Final.docx"),
		},
		{
			ID:   "fcra-identity-theft",
			Name: "FCRA Complaint - Identity Theft",
			Desc: "For cases involving wider identity theft issues",
			Path: "../templates/fcra_identity_theft.docx",
		},
		{
			ID:   "fcra-inaccurate-reporting",
			Name: "FCRA Complaint - Inaccurate Reporting",
			Desc: "For cases involving credit report errors",
			Path: "../templates/fcra_inaccurate_reporting.docx",
		},
	}
	
	return templates, nil
}

// ProcessSelectedDocuments processes documents selected in Step 1 and generates dynamic case data
func (s *DocumentService) ProcessSelectedDocuments(selectedDocIDs []string, templateID string) (*DocumentProcessingResult, *ClientCase, error) {
	log.Printf("[DOCUMENT_SERVICE] Processing %d selected documents with DYNAMIC analysis engine", len(selectedDocIDs))
	
	if s.contentAnalyzer == nil {
		return nil, nil, fmt.Errorf("content analyzer not initialized")
	}
	
	// Initialize result structures
	selectedDocs := []Document{}
	extractedData := make(map[string]interface{})
	missingContent := []MissingContent{}
	allExtractedText := make(map[string]string)
	allAnalysisResults := make(map[string]*LegalAnalysisResult)
	
	// Initialize ClientCase with empty values
	clientCase := ClientCase{}
	
	// Track what types of documents we have
	documentTypes := make(map[string]bool)
	
	// Process each selected document with real text extraction and intelligent analysis
	for _, docPath := range selectedDocIDs {
		log.Printf("[DOCUMENT_SERVICE] Processing document: %s", docPath)
		
		// Extract text from the document
		content, err := s.extractor.ExtractText(docPath)
		if err != nil {
			log.Printf("[DOCUMENT_SERVICE] Error extracting text from %s: %v", docPath, err)
			continue
		}
		
		// Store extracted text for cross-reference
		fileName := filepath.Base(docPath)
		allExtractedText[fileName] = content.RawText
		
		// Create Document object
		fileInfo, _ := os.Stat(docPath)
		size := int64(0)
		if fileInfo != nil {
			size = fileInfo.Size()
		}
		
		doc := Document{
			ID:          fmt.Sprintf("doc_%d", len(selectedDocs)+1),
			Name:        fileName,
			Type:        strings.ToLower(filepath.Ext(fileName)),
			Path:        docPath,
			ContentType: s.determineContentType(fileName),
			Size:        size,
		}
		selectedDocs = append(selectedDocs, doc)
		
		// Track document types
		documentTypes[doc.ContentType] = true
		
		// Perform intelligent analysis of document content
		analysis, err := s.contentAnalyzer.AnalyzeLegalContent(content.RawText, doc.ContentType)
		if err != nil {
			log.Printf("[DOCUMENT_SERVICE] Error analyzing %s: %v", fileName, err)
			continue
		}
		
		// Enhanced attorney notes analysis
		if doc.ContentType == "attorney_notes" && s.attorneyNotesAnalyzer != nil {
			log.Printf("[DOCUMENT_SERVICE] Performing specialized attorney notes analysis for %s", fileName)
			attorneyAnalysis, err := s.attorneyNotesAnalyzer.AnalyzeAttorneyNotes(docPath, content.RawText)
			if err != nil {
				log.Printf("[DOCUMENT_SERVICE] Warning: Attorney notes analysis failed for %s: %v", fileName, err)
			} else {
				// Enhance the analysis with attorney intelligence
				analysis = s.enhanceAnalysisWithAttorneyIntelligence(analysis, attorneyAnalysis)
				log.Printf("[DOCUMENT_SERVICE] Enhanced %s with attorney intelligence - %.1f%% confidence", 
					fileName, attorneyAnalysis.ConfidenceScores.OverallConfidence)
			}
		}
		
		allAnalysisResults[fileName] = analysis
		log.Printf("[DOCUMENT_SERVICE] Analyzed %s - %.1f%% confidence, %d violations found", 
			fileName, analysis.OverallConfidence, len(analysis.LegalViolations))
	}
	
	// Correlate and merge analysis results into ClientCase
	s.correlateAnalysisResults(allAnalysisResults, &clientCase, extractedData)
	
	// Analyze missing content based on intelligent analysis
	missingContent = s.analyzeIntelligentMissingContent(&clientCase, documentTypes, allAnalysisResults)
	
	// Calculate data coverage based on populated fields
	dataCoverage := s.calculateIntelligentDataCoverage(&clientCase, allAnalysisResults)
	
	// Create processing result
	processingResult := &DocumentProcessingResult{
		SelectedDocuments: selectedDocs,
		ExtractedData:     extractedData,
		MissingContent:    missingContent,
		DataCoverage:      dataCoverage,
	}
	
	log.Printf("[DOCUMENT_SERVICE] DYNAMIC extraction complete - %d documents processed, %.1f%% data coverage", 
		len(selectedDocs), dataCoverage)
	
	return processingResult, &clientCase, nil
}

// determineContentType identifies the type of legal document based on filename
func (s *DocumentService) determineContentType(fileName string) string {
	fileName = strings.ToLower(fileName)
	
	if strings.Contains(fileName, "attorney") || strings.Contains(fileName, "atty") {
		return "attorney_notes"
	}
	if strings.Contains(fileName, "adverse") || strings.Contains(fileName, "denial") {
		return "adverse_action"
	}
	if strings.Contains(fileName, "civil") && strings.Contains(fileName, "cover") {
		return "civil_cover_sheet"
	}
	if strings.Contains(fileName, "summons") {
		if strings.Contains(fileName, "equifax") {
			return "summons_equifax"
		}
		return "summons"
	}
	if strings.Contains(fileName, "complaint") {
		return "complaint_form"
	}
	
	return "unknown"
}

// extractDataFromDocument extracts structured data from document text
func (s *DocumentService) extractDataFromDocument(text string, contentType string, clientCase *ClientCase, extractedData map[string]interface{}) {
	log.Printf("[DOCUMENT_SERVICE] Extracting data from %s document (%d chars)", contentType, len(text))
	
	switch contentType {
	case "attorney_notes":
		s.extractFromAttorneyNotes(text, clientCase, extractedData)
	case "adverse_action":
		s.extractFromAdverseAction(text, clientCase, extractedData)
	case "civil_cover_sheet":
		s.extractFromCivilCoverSheet(text, clientCase, extractedData)
	case "summons", "summons_equifax":
		s.extractFromSummons(text, clientCase, extractedData)
	default:
		log.Printf("[DOCUMENT_SERVICE] Unknown document type: %s", contentType)
	}
}

// extractFromAttorneyNotes extracts client and case information from attorney notes
func (s *DocumentService) extractFromAttorneyNotes(text string, clientCase *ClientCase, extractedData map[string]interface{}) {
	// Extract client name
	namePatterns := []string{
		`Client:\s*([A-Z][a-z]+\s+[A-Z][a-z]+)`,
		`Case for\s*([A-Z][a-z]+\s+[A-Z][a-z]+)`,
		`([A-Z][a-z]+\s+[A-Z][a-z]+)\s+Case`,
	}
	
	if name := s.extractFirstMatch(text, namePatterns); name != "" {
		clientCase.ClientName = name
		extractedData["clientName"] = name
		log.Printf("[DOCUMENT_SERVICE] Extracted client name: %s", name)
	}
	
	// Extract contact information
	phonePatterns := []string{
		`(\d{3}\.\d{3}\.\d{4})`,
		`(\d{3})-(\d{3})-(\d{4})`,
		`\((\d{3})\)\s*(\d{3})-(\d{4})`,
	}
	
	if phone := s.extractFirstMatch(text, phonePatterns); phone != "" {
		clientCase.ContactInfo = phone
		extractedData["contactInfo"] = phone
		log.Printf("[DOCUMENT_SERVICE] Extracted contact: %s", phone)
	}
	
	// Extract travel information
	travelPatterns := []string{
		`Travel Dates:\s*([A-Za-z0-9\s,-]+)`,
		`traveling.*in\s+([A-Z][a-z]+)`,
		`while.*was.*in\s+([A-Z][a-z]+)`,
	}
	
	if travel := s.extractFirstMatch(text, travelPatterns); travel != "" {
		clientCase.TravelLocation = travel
		extractedData["travelLocation"] = travel
		log.Printf("[DOCUMENT_SERVICE] Extracted travel location: %s", travel)
	}
	
	// Extract fraud amount
	amountPatterns := []string{
		`Fraud Amount:\s*\$([0-9,]+)`,
		`\$([0-9,]+)`,
		`([0-9,]+)\s*dollars?`,
	}
	
	if amount := s.extractFirstMatch(text, amountPatterns); amount != "" {
		clientCase.FraudAmount = "$" + amount
		extractedData["fraudAmount"] = "$" + amount
		log.Printf("[DOCUMENT_SERVICE] Extracted fraud amount: $%s", amount)
	}
	
	// Extract bank information
	bankPatterns := []string{
		`Bank:\s*([A-Z][a-z\s]+)`,
		`(TD Bank|Chase|Capital One|Citibank|Wells Fargo|Bank of America)`,
	}
	
	if bank := s.extractFirstMatch(text, bankPatterns); bank != "" {
		clientCase.FinancialInstitution = bank
		extractedData["financialInstitution"] = bank
		log.Printf("[DOCUMENT_SERVICE] Extracted bank: %s", bank)
	}
	
	// Mark attorney notes as processed
	extractedData["attorneyNotes"] = true
}

// extractFromAdverseAction extracts credit impact information from adverse action letters
func (s *DocumentService) extractFromAdverseAction(text string, clientCase *ClientCase, extractedData map[string]interface{}) {
	// Extract credit impact details
	impactPatterns := []string{
		`denied credit`,
		`credit.*reduced`,
		`credit.*limit.*reduced`,
		`application.*denied`,
	}
	
	impacts := []string{}
	for _, pattern := range impactPatterns {
		if matched, _ := regexp.MatchString(pattern, strings.ToLower(text)); matched {
			impacts = append(impacts, pattern)
		}
	}
	
	if len(impacts) > 0 {
		clientCase.CreditImpact = strings.Join(impacts, ", ")
		extractedData["creditImpact"] = clientCase.CreditImpact
		log.Printf("[DOCUMENT_SERVICE] Extracted credit impact: %s", clientCase.CreditImpact)
	}
	
	// Set default credit bureaus if adverse action exists
	clientCase.CreditBureauDisputes = []string{"Experian", "Equifax", "Trans Union"}
	extractedData["adverseAction"] = true
}

// extractFromCivilCoverSheet extracts court and legal information
func (s *DocumentService) extractFromCivilCoverSheet(text string, clientCase *ClientCase, extractedData map[string]interface{}) {
	// For now, just mark that we have civil cover sheet data
	// In a full implementation, would extract court jurisdiction, case classification, etc.
	extractedData["civilCoverSheet"] = true
	log.Printf("[DOCUMENT_SERVICE] Civil cover sheet processed")
}

// extractFromSummons extracts defendant information from summons documents
func (s *DocumentService) extractFromSummons(text string, clientCase *ClientCase, extractedData map[string]interface{}) {
	// Extract credit bureau names from summons
	bureauPatterns := []string{
		`Equifax`,
		`Experian`, 
		`Trans Union`,
		`TransUnion`,
	}
	
	bureaus := []string{}
	for _, pattern := range bureauPatterns {
		if matched, _ := regexp.MatchString(pattern, text); matched {
			bureaus = append(bureaus, pattern)
		}
	}
	
	if len(bureaus) > 0 {
		clientCase.CreditBureauDisputes = bureaus
		extractedData["creditBureaus"] = bureaus
		log.Printf("[DOCUMENT_SERVICE] Extracted credit bureaus: %v", bureaus)
	}
	
	extractedData["summons"] = true
}

// extractFirstMatch finds the first matching pattern in text
func (s *DocumentService) extractFirstMatch(text string, patterns []string) string {
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

// analyzeMissingContent determines what data is missing based on extraction results
func (s *DocumentService) analyzeMissingContent(clientCase *ClientCase, documentTypes map[string]bool, extractedData map[string]interface{}) []MissingContent {
	missing := []MissingContent{}
	
	// Check for missing client information
	if clientCase.ClientName == "" {
		missing = append(missing, MissingContent{
			Field:       "Client Name",
			Description: "Client name is required for legal documents",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	if clientCase.ContactInfo == "" {
		missing = append(missing, MissingContent{
			Field:       "Contact Information",
			Description: "Client contact information",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	if clientCase.FraudAmount == "" {
		missing = append(missing, MissingContent{
			Field:       "Fraud Amount",
			Description: "Amount of fraudulent charges",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	if clientCase.FinancialInstitution == "" {
		missing = append(missing, MissingContent{
			Field:       "Financial Institution",
			Description: "Bank or credit card company involved",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	// Check for missing document types
	if !documentTypes["attorney_notes"] {
		missing = append(missing, MissingContent{
			Field:       "Attorney Notes",
			Description: "Attorney notes containing case details",
			Source:      "Attorney Notes Document",
			Required:    true,
		})
	}
	
	if !documentTypes["adverse_action"] && clientCase.CreditImpact == "" {
		missing = append(missing, MissingContent{
			Field:       "Credit Impact Details",
			Description: "Documentation of credit denials or impacts",
			Source:      "Adverse Action Letters",
			Required:    true,
		})
	}
	
	return missing
}

// calculateDataCoverage calculates the percentage of required data that was extracted
func (s *DocumentService) calculateDataCoverage(clientCase *ClientCase, extractedData map[string]interface{}) float64 {
	totalFields := 10 // Key required fields
	populatedFields := 0
	
	if clientCase.ClientName != "" { populatedFields++ }
	if clientCase.ContactInfo != "" { populatedFields++ }
	if clientCase.FraudAmount != "" { populatedFields++ }
	if clientCase.FinancialInstitution != "" { populatedFields++ }
	if clientCase.TravelLocation != "" { populatedFields++ }
	if clientCase.CreditImpact != "" { populatedFields++ }
	if len(clientCase.CreditBureauDisputes) > 0 { populatedFields++ }
	if extractedData["attorneyNotes"] != nil { populatedFields++ }
	if extractedData["adverseAction"] != nil { populatedFields++ }
	if extractedData["civilCoverSheet"] != nil { populatedFields++ }
	
	return float64(populatedFields) / float64(totalFields) * 100
}

// parseDate parses date strings into time.Time
func (s *DocumentService) parseDate(dateStr string) time.Time {
	layouts := []string{
		"January 2, 2006",
		"January 2006",
	}
	
	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t
		}
	}
	
	// Return zero time if parsing fails
	return time.Time{}
}

// correlateAnalysisResults merges intelligent analysis results into ClientCase
func (s *DocumentService) correlateAnalysisResults(analysisResults map[string]*LegalAnalysisResult, clientCase *ClientCase, extractedData map[string]interface{}) {
	log.Printf("[DOCUMENT_SERVICE] Correlating analysis results from %d documents", len(analysisResults))
	
	// Merge client data from all documents (highest confidence wins)
	bestClientName := ""
	bestClientNameConfidence := 0.0
	bestPhone := ""
	bestPhoneConfidence := 0.0
	bestFraudAmount := ""
	bestFraudAmountConfidence := 0.0
	bestInstitution := ""
	bestInstitutionConfidence := 0.0
	bestTravelLocation := ""
	bestTravelConfidence := 0.0
	
	allViolations := []string{}
	creditBureaus := []string{}
	creditImpact := []string{}
	
	for fileName, analysis := range analysisResults {
		log.Printf("[DOCUMENT_SERVICE] Processing %s analysis results", fileName)
		
		// Extract client name (highest confidence)
		if clientName, exists := analysis.ClientData["clientName"]; exists {
			if clientName.Confidence > bestClientNameConfidence {
				bestClientName = clientName.Value.(string)
				bestClientNameConfidence = clientName.Confidence
				log.Printf("[DOCUMENT_SERVICE] Updated client name: %s (%.1f%% confidence)", bestClientName, clientName.Confidence)
			}
		}
		
		// Extract phone number (highest confidence)
		if phone, exists := analysis.ClientData["phoneNumber"]; exists {
			if phone.Confidence > bestPhoneConfidence {
				bestPhone = phone.Value.(string)
				bestPhoneConfidence = phone.Confidence
				log.Printf("[DOCUMENT_SERVICE] Updated phone: %s (%.1f%% confidence)", bestPhone, phone.Confidence)
			}
		}
		
		// Extract fraud amount (highest confidence)
		if fraudAmount, exists := analysis.FraudDetails["fraudAmount"]; exists {
			if fraudAmount.Confidence > bestFraudAmountConfidence {
				bestFraudAmount = fraudAmount.Value.(string)
				bestFraudAmountConfidence = fraudAmount.Confidence
				log.Printf("[DOCUMENT_SERVICE] Updated fraud amount: %s (%.1f%% confidence)", bestFraudAmount, fraudAmount.Confidence)
			}
		}
		
		// Extract institution (highest confidence)
		if institution, exists := analysis.FraudDetails["institution"]; exists {
			if institution.Confidence > bestInstitutionConfidence {
				bestInstitution = institution.Value.(string)
				bestInstitutionConfidence = institution.Confidence
				log.Printf("[DOCUMENT_SERVICE] Updated institution: %s (%.1f%% confidence)", bestInstitution, institution.Confidence)
			}
		}
		
		// Extract travel location (highest confidence)
		if travel, exists := analysis.FraudDetails["travelLocation"]; exists {
			if travel.Confidence > bestTravelConfidence {
				bestTravelLocation = travel.Value.(string)
				bestTravelConfidence = travel.Confidence
				log.Printf("[DOCUMENT_SERVICE] Updated travel location: %s (%.1f%% confidence)", bestTravelLocation, travel.Confidence)
			}
		}
		
		// Accumulate violations
		allViolations = append(allViolations, analysis.LegalViolations...)
		
		// Extract credit impact indicators
		if strings.Contains(strings.ToLower(fileName), "adverse") {
			creditImpact = append(creditImpact, "denied credit", "application denied")
		}
		
		// Extract credit bureaus from summons documents
		if strings.Contains(strings.ToLower(fileName), "summons") {
			if strings.Contains(strings.ToLower(fileName), "equifax") {
				creditBureaus = append(creditBureaus, "Equifax")
			}
			if strings.Contains(strings.ToLower(fileName), "experian") {
				creditBureaus = append(creditBureaus, "Experian")
			}
			if strings.Contains(strings.ToLower(fileName), "trans") || strings.Contains(strings.ToLower(fileName), "union") {
				creditBureaus = append(creditBureaus, "Trans Union")
			}
		}
	}
	
	// Populate ClientCase with best extracted data
	clientCase.ClientName = bestClientName
	clientCase.ContactInfo = bestPhone
	clientCase.FraudAmount = bestFraudAmount
	clientCase.FinancialInstitution = bestInstitution
	clientCase.TravelLocation = bestTravelLocation
	
	// Set credit impact and credit bureaus
	if len(creditImpact) > 0 {
		clientCase.CreditImpact = strings.Join(removeDuplicates(creditImpact), ", ")
	}
	if len(creditBureaus) > 0 {
		clientCase.CreditBureauDisputes = removeDuplicates(creditBureaus)
	} else {
		// Default credit bureaus if not specifically identified
		clientCase.CreditBureauDisputes = []string{"Experian", "Equifax", "Trans Union"}
	}
	
	// Set some reasonable defaults based on extracted data
	if clientCase.ClientName != "" {
		clientCase.ResidenceLocation = "United States" // Could be extracted from address patterns
	}
	if clientCase.FraudAmount != "" {
		clientCase.FraudDetails = fmt.Sprintf("Fraudulent charges totaling %s", clientCase.FraudAmount)
		// Set reasonable discovery date
		clientCase.DiscoveryDate = time.Now().AddDate(0, -3, 0) // 3 months ago
		clientCase.FraudStartDate = time.Now().AddDate(0, -6, 0) // 6 months ago  
		clientCase.FraudEndDate = time.Now().AddDate(0, -3, 0) // 3 months ago
	}
	
	// Store analysis data for UI
	extractedData["analysisResults"] = analysisResults
	extractedData["extractedViolations"] = removeDuplicates(allViolations)
	extractedData["dynamicExtraction"] = true
	extractedData["totalConfidence"] = (bestClientNameConfidence + bestFraudAmountConfidence + bestInstitutionConfidence) / 3
	
	log.Printf("[DOCUMENT_SERVICE] Correlation complete - Client: %s, Amount: %s, Institution: %s", 
		clientCase.ClientName, clientCase.FraudAmount, clientCase.FinancialInstitution)
}

// analyzeIntelligentMissingContent determines missing content based on intelligent analysis
func (s *DocumentService) analyzeIntelligentMissingContent(clientCase *ClientCase, documentTypes map[string]bool, analysisResults map[string]*LegalAnalysisResult) []MissingContent {
	missing := []MissingContent{}
	
	// Check for missing core client information
	if clientCase.ClientName == "" {
		missing = append(missing, MissingContent{
			Field:       "Client Name",
			Description: "Client name is required for legal documents",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	if clientCase.ContactInfo == "" {
		missing = append(missing, MissingContent{
			Field:       "Contact Information", 
			Description: "Client contact information",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	if clientCase.FraudAmount == "" {
		missing = append(missing, MissingContent{
			Field:       "Fraud Amount",
			Description: "Amount of fraudulent charges",
			Source:      "Attorney Notes or Adverse Action Letters",
			Required:    true,
		})
	}
	
	if clientCase.FinancialInstitution == "" {
		missing = append(missing, MissingContent{
			Field:       "Financial Institution",
			Description: "Bank or credit card company involved",
			Source:      "Attorney Notes or Documents",
			Required:    true,
		})
	}
	
	// Check for missing document types based on analysis
	foundAttorneyNotes := false
	foundAdverseAction := false
	for _, analysis := range analysisResults {
		if analysis.DocumentTypes["attorneyNotes"] {
			foundAttorneyNotes = true
		}
		if analysis.DocumentTypes["adverseAction"] {
			foundAdverseAction = true
		}
	}
	
	if !foundAttorneyNotes {
		missing = append(missing, MissingContent{
			Field:       "Attorney Notes",
			Description: "Attorney notes containing case details and client information",
			Source:      "Attorney Notes Document",
			Required:    true,
		})
	}
	
	if !foundAdverseAction && clientCase.CreditImpact == "" {
		missing = append(missing, MissingContent{
			Field:       "Credit Impact Details",
			Description: "Documentation of credit denials or impacts",
			Source:      "Adverse Action Letters",
			Required:    true,
		})
	}
	
	return missing
}

// calculateIntelligentDataCoverage calculates coverage based on intelligent analysis results
func (s *DocumentService) calculateIntelligentDataCoverage(clientCase *ClientCase, analysisResults map[string]*LegalAnalysisResult) float64 {
	totalFields := 8 // Core required fields for legal complaint
	populatedFields := 0
	
	if clientCase.ClientName != "" { populatedFields++ }
	if clientCase.ContactInfo != "" { populatedFields++ }
	if clientCase.FraudAmount != "" { populatedFields++ }
	if clientCase.FinancialInstitution != "" { populatedFields++ }
	if clientCase.TravelLocation != "" { populatedFields++ }
	if clientCase.CreditImpact != "" { populatedFields++ }
	if len(clientCase.CreditBureauDisputes) > 0 { populatedFields++ }
	if clientCase.FraudDetails != "" { populatedFields++ }
	
	coverage := float64(populatedFields) / float64(totalFields) * 100
	
	// Boost coverage based on analysis confidence
	totalConfidence := 0.0
	confidenceCount := 0
	for _, analysis := range analysisResults {
		if analysis.OverallConfidence > 0 {
			totalConfidence += analysis.OverallConfidence
			confidenceCount++
		}
	}
	
	if confidenceCount > 0 {
		avgConfidence := totalConfidence / float64(confidenceCount)
		// Adjust coverage based on confidence (higher confidence = better coverage)
		coverage = coverage * (avgConfidence / 100)
	}
	
	return coverage
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// GenerateComplaint generates a legal complaint using the dynamic template engine
func (s *DocumentService) GenerateComplaint(templateID string, clientCase *ClientCase) (*GeneratedDocument, error) {
	if s.templateEngine == nil {
		return nil, fmt.Errorf("template engine not initialized")
	}
	
	log.Printf("[DOCUMENT_SERVICE] Generating complaint using template: %s for client: %s", templateID, clientCase.ClientName)
	
	// Convert ClientCase to enhanced format for template engine
	enhancedClientCase := s.convertToEnhancedClientCase(clientCase)
	
	// Generate document using template engine
	document, err := s.templateEngine.GenerateDocument(templateID, enhancedClientCase)
	if err != nil {
		return nil, fmt.Errorf("failed to generate document: %v", err)
	}
	
	log.Printf("[DOCUMENT_SERVICE] Generated complaint: %d sections, %.1f%% complete", 
		len(document.Sections), document.Metadata.Completeness*100)
	
	return document, nil
}

// convertToEnhancedClientCase converts the basic ClientCase to the enhanced format
func (s *DocumentService) convertToEnhancedClientCase(basic *ClientCase) *ClientCase {
	// Create enhanced ClientCase with additional fields for template engine
	enhanced := &ClientCase{
		ClientName:        basic.ClientName,
		ContactInfo:       basic.ContactInfo,
		ResidenceLocation: s.determineResidenceLocation(basic),
		CourtJurisdiction: s.determineCourtJurisdiction(basic),
		CaseNumber:        s.generateCaseNumber(),
	}
	
	// Convert fraud details to structured format
	if basic.FraudAmount != "" || basic.FraudDetails != "" {
		enhanced.FraudDetailsStructured = []FraudDetail{
			{
				Institution: basic.FinancialInstitution,
				Amount:      basic.FraudAmount,
				Description: basic.FraudDetails,
				Date:        basic.FraudStartDate,
			},
		}
	}
	
	// Convert credit bureau interactions
	if len(basic.CreditBureauDisputes) > 0 {
		enhanced.CreditBureauInteractions = []CreditBureauInteraction{}
		for _, bureau := range basic.CreditBureauDisputes {
			enhanced.CreditBureauInteractions = append(enhanced.CreditBureauInteractions, CreditBureauInteraction{
				Bureau:   bureau,
				Type:     "dispute",
				Date:     basic.CreditBureauDisputeDate.Format("January 2, 2006"),
				Response: basic.BankResponse,
			})
		}
	}
	
	// Add defendants based on available information
	enhanced.Defendants = s.generateDefendants(basic)
	
	return enhanced
}

// determineResidenceLocation determines the client's residence location
func (s *DocumentService) determineResidenceLocation(clientCase *ClientCase) string {
	if clientCase.ResidenceLocation != "" {
		return clientCase.ResidenceLocation
	}
	
	// Default residence location if not specified
	return "New York, New York"
}

// determineCourtJurisdiction determines the appropriate court jurisdiction
func (s *DocumentService) determineCourtJurisdiction(clientCase *ClientCase) string {
	// For now, default to Southern District of New York
	// In production, this would be based on client location and case specifics
	return "SOUTHERN DISTRICT OF NEW YORK"
}

// generateCaseNumber generates a placeholder case number
func (s *DocumentService) generateCaseNumber() string {
	// In production, this would integrate with court systems
	return "[TO BE ASSIGNED]"
}

// generateDefendants creates defendant information based on the case
func (s *DocumentService) generateDefendants(clientCase *ClientCase) []Defendant {
	defendants := []Defendant{}
	
	// Add credit bureaus as defendants
	creditBureaus := []struct {
		name    string
		address string
	}{
		{"EXPERIAN INFORMATION SOLUTIONS, INC.", "475 Anton Blvd., Costa Mesa, CA 92626"},
		{"TRANS UNION LLC", "555 W. Adams Street, Chicago, IL 60661"},
		{"EQUIFAX INFORMATION SERVICES LLC", "1550 Peachtree Street, N.W., Atlanta, GA 30309"},
	}
	
	for _, bureau := range creditBureaus {
		defendants = append(defendants, Defendant{
			Name:       bureau.name,
			EntityType: "corporation",
			Address:    bureau.address,
		})
	}
	
	// Add financial institution if specified
	if clientCase.FinancialInstitution != "" {
		defendants = append(defendants, Defendant{
			Name:       strings.ToUpper(clientCase.FinancialInstitution),
			EntityType: "financial institution",
			Address:    "[ADDRESS TO BE DETERMINED]",
		})
	}
	
	return defendants
}

// Additional types needed for enhanced ClientCase
type FraudDetail struct {
	Institution string    `json:"institution"`
	Amount      string    `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type CreditBureauInteraction struct {
	Bureau   string `json:"bureau"`
	Type     string `json:"type"`
	Date     string `json:"date"`
	Response string `json:"response"`
}

type Defendant struct {
	Name       string `json:"name"`
	EntityType string `json:"entityType"`
	Address    string `json:"address"`
}

// enhanceAnalysisWithAttorneyIntelligence merges attorney notes analysis with standard content analysis
func (s *DocumentService) enhanceAnalysisWithAttorneyIntelligence(analysis *LegalAnalysisResult, attorneyAnalysis *AttorneyNotesAnalysis) *LegalAnalysisResult {
	if analysis == nil || attorneyAnalysis == nil {
		return analysis
	}
	
	log.Printf("[DOCUMENT_SERVICE] Enhancing analysis with attorney intelligence")
	
	// Enhance legal violations with attorney-documented violations
	if len(attorneyAnalysis.ViolationAssessment.IdentifiedViolations) > 0 {
		for _, violation := range attorneyAnalysis.ViolationAssessment.IdentifiedViolations {
			// Add attorney-documented violations to the analysis
			violationDesc := fmt.Sprintf("Attorney Analysis: %s (%s) - %s", 
				violation.ViolationType, violation.Statute, violation.ViolationDescription)
			analysis.LegalViolations = append(analysis.LegalViolations, violationDesc)
		}
		log.Printf("[DOCUMENT_SERVICE] Added %d attorney-documented violations", 
			len(attorneyAnalysis.ViolationAssessment.IdentifiedViolations))
	}
	
	// Enhance client data with attorney consultation information
	if attorneyAnalysis.ClientConsultation.ConsultationDate.Year() > 1 {
		// Add consultation information to client data
		if analysis.ClientData == nil {
			analysis.ClientData = make(map[string]ExtractionResult)
		}
		
		analysis.ClientData["ConsultationDate"] = ExtractionResult{
			Field:      "ConsultationDate",
			Value:      attorneyAnalysis.ClientConsultation.ConsultationDate.Format("January 2, 2006"),
			Confidence: attorneyAnalysis.ClientConsultation.Confidence,
			Source:     "Attorney Notes",
			Method:     "Attorney Intelligence Analysis",
		}
	}
	
	// Enhance with attorney's legal theory
	if attorneyAnalysis.LegalAnalysis.LegalTheory != "" {
		analysis.ClientData["LegalTheory"] = ExtractionResult{
			Field:      "LegalTheory",
			Value:      attorneyAnalysis.LegalAnalysis.LegalTheory,
			Confidence: attorneyAnalysis.LegalAnalysis.Confidence,
			Source:     "Attorney Legal Analysis",
			Method:     "Attorney Intelligence Analysis",
		}
	}
	
	// Enhance with case strategy
	if attorneyAnalysis.CaseStrategy.LegalStrategy != "" {
		analysis.ClientData["CaseStrategy"] = ExtractionResult{
			Field:      "CaseStrategy",
			Value:      attorneyAnalysis.CaseStrategy.LegalStrategy,
			Confidence: attorneyAnalysis.CaseStrategy.Confidence,
			Source:     "Attorney Strategic Analysis",
			Method:     "Attorney Intelligence Analysis",
		}
	}
	
	// Enhance with damage assessment
	if attorneyAnalysis.DamageAssessment.TotalDamageRange.EstimatedAmount > 0 {
		analysis.ClientData["EstimatedDamages"] = ExtractionResult{
			Field:      "EstimatedDamages",
			Value:      fmt.Sprintf("$%.2f", attorneyAnalysis.DamageAssessment.TotalDamageRange.EstimatedAmount),
			Confidence: attorneyAnalysis.DamageAssessment.Confidence,
			Source:     "Attorney Damage Analysis",
			Method:     "Attorney Intelligence Analysis",
		}
	}
	
	// Enhance with attorney's case strength assessment
	if attorneyAnalysis.LegalAnalysis.CaseStrength.OverallStrength != "" {
		analysis.ClientData["CaseStrength"] = ExtractionResult{
			Field:      "CaseStrength",
			Value:      attorneyAnalysis.LegalAnalysis.CaseStrength.OverallStrength,
			Confidence: attorneyAnalysis.LegalAnalysis.CaseStrength.LegalMerits,
			Source:     "Attorney Case Assessment",
			Method:     "Attorney Intelligence Analysis",
		}
	}
	
	// Enhance suggestions with attorney next steps
	if len(attorneyAnalysis.NextSteps) > 0 {
		for _, step := range attorneyAnalysis.NextSteps {
			analysis.Suggestions = append(analysis.Suggestions, 
				fmt.Sprintf("Attorney Action Item: %s (Priority: %s)", step.Task, step.Priority))
		}
	}
	
	// Update overall confidence with weighted average including attorney analysis
	attorneyWeight := 0.4  // Give attorney analysis significant weight
	standardWeight := 0.6
	
	weightedConfidence := (analysis.OverallConfidence * standardWeight) + 
						  (attorneyAnalysis.ConfidenceScores.OverallConfidence * attorneyWeight)
	analysis.OverallConfidence = weightedConfidence
	
	log.Printf("[DOCUMENT_SERVICE] Enhanced analysis - overall confidence: %.2f, violations: %d", 
		analysis.OverallConfidence, len(analysis.LegalViolations))
	
	return analysis
}