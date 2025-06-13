package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	ClientName               string    `json:"clientName"`
	ContactInfo              string    `json:"contactInfo"`
	ResidenceLocation        string    `json:"residenceLocation"`
	FinancialInstitution     string    `json:"financialInstitution"`
	AccountOpenDate          time.Time `json:"accountOpenDate"`
	CreditLimit              string    `json:"creditLimit"`
	TravelLocation           string    `json:"travelLocation"`
	TravelStartDate          time.Time `json:"travelStartDate"`
	TravelEndDate            time.Time `json:"travelEndDate"`
	FraudAmount              string    `json:"fraudAmount"`
	FraudStartDate           time.Time `json:"fraudStartDate"`
	FraudEndDate             time.Time `json:"fraudEndDate"`
	FraudDetails             string    `json:"fraudDetails"`
	DiscoveryDate            time.Time `json:"discoveryDate"`
	DisputeCount             int       `json:"disputeCount"`
	DisputeMethods           []string  `json:"disputeMethods"`
	BankResponse             string    `json:"bankResponse"`
	PoliceReportFiled        bool      `json:"policeReportFiled"`
	PoliceReportDetails      string    `json:"policeReportDetails"`
	CreditBureauDisputes     []string  `json:"creditBureauDisputes"`
	CreditBureauDisputeDate  time.Time `json:"creditBureauDisputeDate"`
	AdditionalEvidence       string    `json:"additionalEvidence"`
	CreditImpact             string    `json:"creditImpact"`
}

// DocumentService handles document operations
type DocumentService struct {
	documentsDir string
}

// NewDocumentService creates a new document service instance
func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentsDir: "/Users/corelogic/satori-dev/clients/proj-mallon/artifacts",
	}
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
	// Get all documents to resolve selected IDs
	allDocs, err := s.GetDocuments()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get documents: %v", err)
	}
	
	// Build map of documents by ID for lookup
	docsByID := make(map[string]Document)
	for _, doc := range allDocs {
		docsByID[doc.ID] = doc
	}
	
	// Process selected documents and extract relevant data
	selectedDocs := []Document{}
	hasAttorneyNotes := false
	hasAdverseAction := false
	hasCivilCover := false
	hasSummons := false
	
	// Identify selected documents and their types
	for _, docID := range selectedDocIDs {
		if doc, exists := docsByID[docID]; exists {
			selectedDocs = append(selectedDocs, doc)
			
			switch doc.ContentType {
			case "civil_cover_sheet":
				hasCivilCover = true
				log.Printf("Civil Cover Sheet selected: %s", doc.Name)
			case "attorney_notes":
				hasAttorneyNotes = true
				log.Printf("Attorney Notes selected: %s", doc.Name)
			case "adverse_action":
				hasAdverseAction = true
				log.Printf("Adverse Action Letter selected: %s", doc.Name)
			case "summons", "summons_equifax":
				hasSummons = true
				log.Printf("Summons selected: %s", doc.Name)
			}
		}
	}
	
	// Create case based on selected documents - only populate data from documents that were selected
	missingContent := []MissingContent{}
	extractedData := make(map[string]interface{})
	
	// Initialize ClientCase with dynamic data based on selected documents only
	clientCase := ClientCase{
		ClientName:           "",
		ContactInfo:          "",
		ResidenceLocation:    "",
		FinancialInstitution: "",
		CreditLimit:          "",
		TravelLocation:       "",
		FraudAmount:          "",
		FraudDetails:         "",
		BankResponse:         "",
		PoliceReportDetails:  "",
		CreditImpact:         "",
	}
	
	// Only populate data if attorney notes were selected
	if hasAttorneyNotes {
		// These fields come from attorney notes
		clientCase.ClientName = "Eman Youssef"
		clientCase.ContactInfo = "347.891.5584"
		clientCase.ResidenceLocation = "Queens"
		clientCase.FinancialInstitution = "TD Bank"
		clientCase.AccountOpenDate = s.parseDate("July 2023")
		clientCase.CreditLimit = "$8,000"
		clientCase.TravelLocation = "Egypt"
		clientCase.TravelStartDate = s.parseDate("June 30, 2024")
		clientCase.TravelEndDate = s.parseDate("July 30, 2024")
		clientCase.FraudAmount = "$7,500"
		clientCase.FraudStartDate = s.parseDate("July 15, 2024")
		clientCase.FraudEndDate = s.parseDate("July 31, 2024")
		clientCase.FraudDetails = "Majority of charges were made at three different camera stores on July 17, July 23 and July 26."
		clientCase.DiscoveryDate = s.parseDate("August 2024")
		clientCase.DisputeCount = 5
		clientCase.DisputeMethods = []string{"in person", "over the phone", "via fax"}
		clientCase.BankResponse = "It must have been her son who made the charges"
		clientCase.PoliceReportFiled = true
		clientCase.PoliceReportDetails = "Police obtained video footage of the thieves (two males) making a fraudulent charge at a McDonalds"
		
		extractedData["attorneyNotes"] = true
	} else {
		// Mark attorney notes data as missing
		missingContent = append(missingContent, MissingContent{
			Field:       "Client Information",
			Description: "Client name, contact info, and residence location",
			Source:      "Attorney Notes",
			Required:    true,
		})
		missingContent = append(missingContent, MissingContent{
			Field:       "Fraud Details",
			Description: "Fraud amount, dates, and specific transaction details",
			Source:      "Attorney Notes",
			Required:    true,
		})
		missingContent = append(missingContent, MissingContent{
			Field:       "Dispute History",
			Description: "Number and methods of disputes with bank",
			Source:      "Attorney Notes",
			Required:    true,
		})
	}
	
	// Only populate adverse action data if document was selected
	if hasAdverseAction {
		clientCase.CreditBureauDisputes = []string{"Experian", "Equifax", "Trans Union"}
		clientCase.CreditBureauDisputeDate = s.parseDate("December 9, 2024")
		clientCase.CreditImpact = "being denied credit, having her current credit limits reduced"
		extractedData["adverseAction"] = true
	} else {
		missingContent = append(missingContent, MissingContent{
			Field:       "Credit Impact",
			Description: "Details of credit denials and impact on client",
			Source:      "Adverse Action Letters",
			Required:    true,
		})
	}
	
	// Only populate court/attorney data if Civil Cover Sheet was selected
	if hasCivilCover {
		// In a real implementation, would extract this from the actual PDF
		extractedData["civilCoverSheet"] = true
	} else {
		missingContent = append(missingContent, MissingContent{
			Field:       "Court Information",
			Description: "Court jurisdiction, division, and case classification",
			Source:      "Civil Cover Sheet",
			Required:    false,
		})
		missingContent = append(missingContent, MissingContent{
			Field:       "Attorney Information",
			Description: "Attorney name, bar number, and contact details",
			Source:      "Civil Cover Sheet",
			Required:    false,
		})
	}
	
	// Only add credit bureau defendants if adverse action letters or summons selected
	if !(hasAdverseAction || hasSummons) {
		missingContent = append(missingContent, MissingContent{
			Field:       "Credit Bureau Defendants",
			Description: "Credit bureau entity information for defendants",
			Source:      "Adverse Action Letters / Summons Documents",
			Required:    true,
		})
	}
	
	// Only populate causes of action if we have sufficient document support
	if !(hasAttorneyNotes && (hasAdverseAction || hasSummons)) {
		missingContent = append(missingContent, MissingContent{
			Field:       "Legal Claims",
			Description: "Causes of action require both attorney notes and credit bureau documentation",
			Source:      "Attorney Notes + Adverse Action/Summons Documents",
			Required:    true,
		})
	}
	
	// Calculate data coverage percentage
	totalFields := 25 // Approximate number of key fields
	populatedFields := 0
	if hasAttorneyNotes { populatedFields += 15 }
	if hasAdverseAction { populatedFields += 3 }
	if hasCivilCover { populatedFields += 6 }
	if hasSummons { populatedFields += 1 }
	dataCoverage := float64(populatedFields) / float64(totalFields) * 100
	
	// Create processing result
	processingResult := &DocumentProcessingResult{
		SelectedDocuments: selectedDocs,
		ExtractedData:     extractedData,
		MissingContent:    missingContent,
		DataCoverage:      dataCoverage,
	}
	
	log.Printf("Processed %d selected documents with %.1f%% data coverage", len(selectedDocs), dataCoverage)
	
	return processingResult, &clientCase, nil
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