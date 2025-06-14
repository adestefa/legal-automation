package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// DocumentValidator validates generated legal documents for completeness and accuracy
type DocumentValidator struct {
	RequiredSections []string
	LegalPatterns    []ValidationPattern
}

// ValidationPattern represents a pattern to validate in legal documents
type ValidationPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Required    bool
	Description string
	Severity    string
}

// NewDocumentValidator creates a new document validator
func NewDocumentValidator() *DocumentValidator {
	validator := &DocumentValidator{
		RequiredSections: []string{
			"PARTIES",
			"FACTUAL ALLEGATIONS", 
			"CAUSES OF ACTION",
			"PRAYER FOR RELIEF",
		},
		LegalPatterns: []ValidationPattern{},
	}
	
	validator.loadValidationPatterns()
	
	log.Printf("[DOCUMENT_VALIDATOR] Initialized with %d required sections, %d validation patterns",
		len(validator.RequiredSections), len(validator.LegalPatterns))
	
	return validator
}

// ValidateDocument performs comprehensive validation of a generated legal document
func (dv *DocumentValidator) ValidateDocument(content string, clientCase *ClientCase) []ValidationIssue {
	var issues []ValidationIssue
	
	log.Printf("[DOCUMENT_VALIDATOR] Validating document for client: %s", clientCase.ClientName)
	
	// Check for required sections
	sectionIssues := dv.validateRequiredSections(content)
	issues = append(issues, sectionIssues...)
	
	// Check for required legal elements
	elementIssues := dv.validateLegalElements(content, clientCase)
	issues = append(issues, elementIssues...)
	
	// Check document formatting
	formatIssues := dv.validateFormatting(content)
	issues = append(issues, formatIssues...)
	
	// Check for placeholder text
	placeholderIssues := dv.validatePlaceholders(content)
	issues = append(issues, placeholderIssues...)
	
	// Check legal citations
	citationIssues := dv.validateCitations(content)
	issues = append(issues, citationIssues...)
	
	log.Printf("[DOCUMENT_VALIDATOR] Validation complete: %d issues found", len(issues))
	
	return issues
}

// validateRequiredSections checks that all required sections are present
func (dv *DocumentValidator) validateRequiredSections(content string) []ValidationIssue {
	var issues []ValidationIssue
	upperContent := strings.ToUpper(content)
	
	for _, section := range dv.RequiredSections {
		if !strings.Contains(upperContent, section) {
			issues = append(issues, ValidationIssue{
				Type:        "missing_section",
				Section:     section,
				Description: fmt.Sprintf("Required section '%s' is missing from the document", section),
				Severity:    "high",
				Suggestion:  fmt.Sprintf("Add the '%s' section to complete the document structure", section),
			})
		}
	}
	
	return issues
}

// validateLegalElements checks for required legal elements
func (dv *DocumentValidator) validateLegalElements(content string, clientCase *ClientCase) []ValidationIssue {
	var issues []ValidationIssue
	
	// Check for client name
	if clientCase.ClientName != "" && !strings.Contains(content, clientCase.ClientName) {
		issues = append(issues, ValidationIssue{
			Type:        "missing_client_info",
			Section:     "PARTIES",
			Description: "Client name is not mentioned in the document",
			Severity:    "high",
			Suggestion:  "Ensure the client name appears in the parties section",
		})
	}
	
	// Check for jurisdiction
	if clientCase.CourtJurisdiction == "" {
		issues = append(issues, ValidationIssue{
			Type:        "missing_jurisdiction",
			Section:     "HEADER",
			Description: "Court jurisdiction is not specified",
			Severity:    "medium",
			Suggestion:  "Specify the court jurisdiction in the document header",
		})
	}
	
	// Check for defendants
	if len(clientCase.Defendants) == 0 {
		issues = append(issues, ValidationIssue{
			Type:        "missing_defendants",
			Section:     "PARTIES",
			Description: "No defendants are specified in the case",
			Severity:    "high",
			Suggestion:  "Add defendant information to the parties section",
		})
	}
	
	// Check for fraud details
	if len(clientCase.FraudDetails) == 0 {
		issues = append(issues, ValidationIssue{
			Type:        "missing_fraud_details",
			Section:     "FACTUAL ALLEGATIONS",
			Description: "No fraud details are specified in the case",
			Severity:    "medium",
			Suggestion:  "Add specific fraud allegations to strengthen the case",
		})
	}
	
	// Check for FCRA citations
	if !dv.containsFCRACitations(content) {
		issues = append(issues, ValidationIssue{
			Type:        "missing_citations",
			Section:     "CAUSES OF ACTION",
			Description: "FCRA statutory citations are missing or incomplete",
			Severity:    "high",
			Suggestion:  "Include proper FCRA citations (15 U.S.C. § 1681 et seq.)",
		})
	}
	
	return issues
}

// validateFormatting checks document formatting requirements
func (dv *DocumentValidator) validateFormatting(content string) []ValidationIssue {
	var issues []ValidationIssue
	
	// Check for proper case number format
	caseNumberPattern := regexp.MustCompile(`Case No\. \d+`)
	if !caseNumberPattern.MatchString(content) {
		issues = append(issues, ValidationIssue{
			Type:        "formatting_issue",
			Section:     "HEADER",
			Description: "Case number format may be incorrect or missing",
			Severity:    "low",
			Suggestion:  "Ensure case number follows format 'Case No. [number]'",
		})
	}
	
	// Check for proper paragraph numbering
	paragraphPattern := regexp.MustCompile(`\d+\.\s+`)
	matches := paragraphPattern.FindAllString(content, -1)
	if len(matches) < 5 {
		issues = append(issues, ValidationIssue{
			Type:        "formatting_issue",
			Section:     "GENERAL",
			Description: "Document may lack proper paragraph numbering",
			Severity:    "low",
			Suggestion:  "Use numbered paragraphs for factual allegations",
		})
	}
	
	return issues
}

// validatePlaceholders checks for unreplaced placeholder text
func (dv *DocumentValidator) validatePlaceholders(content string) []ValidationIssue {
	var issues []ValidationIssue
	
	placeholders := []string{
		"[CLIENT NAME]",
		"[COURT]",
		"[CASE NUMBER]",
		"[DEFENDANT]",
		"[DATE]",
		"[AMOUNT]",
		"{{",
		"}}",
		"[TO BE DETERMINED]",
		"[FILL IN]",
	}
	
	for _, placeholder := range placeholders {
		if strings.Contains(content, placeholder) {
			issues = append(issues, ValidationIssue{
				Type:        "placeholder_text",
				Section:     "GENERAL",
				Description: fmt.Sprintf("Placeholder text '%s' found in document", placeholder),
				Severity:    "high",
				Suggestion:  "Replace placeholder text with actual case information",
			})
		}
	}
	
	return issues
}

// validateCitations checks for proper legal citations
func (dv *DocumentValidator) validateCitations(content string) []ValidationIssue {
	var issues []ValidationIssue
	
	// Check for FCRA citations
	fcraPattern := regexp.MustCompile(`15 U\.S\.C\. § 168[1-9]`)
	if !fcraPattern.MatchString(content) {
		issues = append(issues, ValidationIssue{
			Type:        "missing_citations",
			Section:     "CAUSES OF ACTION",
			Description: "Proper FCRA statutory citations are missing",
			Severity:    "medium",
			Suggestion:  "Include specific FCRA citations (e.g., 15 U.S.C. § 1681e, § 1681i)",
		})
	}
	
	// Check for section references
	sectionPattern := regexp.MustCompile(`§\s*\d+`)
	if !sectionPattern.MatchString(content) {
		issues = append(issues, ValidationIssue{
			Type:        "formatting_issue",
			Section:     "CITATIONS",
			Description: "Legal section symbols may be missing or improperly formatted",
			Severity:    "low",
			Suggestion:  "Use proper section symbols (§) in legal citations",
		})
	}
	
	return issues
}

// containsFCRACitations checks if the document contains proper FCRA citations
func (dv *DocumentValidator) containsFCRACitations(content string) bool {
	fcraPatterns := []string{
		"15 U.S.C.",
		"Fair Credit Reporting Act",
		"FCRA",
		"§ 1681",
	}
	
	upperContent := strings.ToUpper(content)
	foundCitations := 0
	
	for _, pattern := range fcraPatterns {
		if strings.Contains(upperContent, strings.ToUpper(pattern)) {
			foundCitations++
		}
	}
	
	// Require at least 2 different types of FCRA references
	return foundCitations >= 2
}

// loadValidationPatterns loads predefined validation patterns
func (dv *DocumentValidator) loadValidationPatterns() {
	patterns := []struct {
		name        string
		pattern     string
		required    bool
		description string
		severity    string
	}{
		{
			name:        "court_header",
			pattern:     `UNITED STATES DISTRICT COURT`,
			required:    true,
			description: "Federal court header",
			severity:    "high",
		},
		{
			name:        "plaintiff_designation",
			pattern:     `Plaintiff.*is.*consumer`,
			required:    true,
			description: "Plaintiff identified as consumer",
			severity:    "high",
		},
		{
			name:        "fcra_violation",
			pattern:     `violation.*Fair Credit Reporting Act`,
			required:    true,
			description: "FCRA violation allegation",
			severity:    "high",
		},
		{
			name:        "prayer_relief",
			pattern:     `WHEREFORE.*Plaintiff.*requests`,
			required:    true,
			description: "Prayer for relief section",
			severity:    "medium",
		},
		{
			name:        "attorney_signature",
			pattern:     `Respectfully submitted`,
			required:    true,
			description: "Attorney signature block",
			severity:    "medium",
		},
	}
	
	for _, p := range patterns {
		compiled, err := regexp.Compile(p.pattern)
		if err != nil {
			log.Printf("[DOCUMENT_VALIDATOR] Error compiling pattern '%s': %v", p.name, err)
			continue
		}
		
		dv.LegalPatterns = append(dv.LegalPatterns, ValidationPattern{
			Name:        p.name,
			Pattern:     compiled,
			Required:    p.required,
			Description: p.description,
			Severity:    p.severity,
		})
	}
}

// ValidateForCourtFiling validates document for court filing requirements
func (dv *DocumentValidator) ValidateForCourtFiling(content string) []ValidationIssue {
	var issues []ValidationIssue
	
	// Check for required court filing elements
	requiredElements := map[string]string{
		"UNITED STATES DISTRICT COURT": "Court identification header",
		"Case No.":                     "Case number",
		"COMPLAINT":                    "Document type identification",
		"Respectfully submitted":       "Attorney signature",
		"PRAYER FOR RELIEF":           "Prayer section",
	}
	
	upperContent := strings.ToUpper(content)
	
	for element, description := range requiredElements {
		if !strings.Contains(upperContent, element) {
			issues = append(issues, ValidationIssue{
				Type:        "court_filing_requirement",
				Section:     "GENERAL",
				Description: fmt.Sprintf("Missing required court filing element: %s", description),
				Severity:    "high",
				Suggestion:  fmt.Sprintf("Add '%s' to meet court filing requirements", element),
			})
		}
	}
	
	return issues
}

// GetValidationScore calculates an overall validation score (0-100)
func (dv *DocumentValidator) GetValidationScore(issues []ValidationIssue) float64 {
	if len(issues) == 0 {
		return 100.0
	}
	
	totalDeductions := 0.0
	
	for _, issue := range issues {
		switch issue.Severity {
		case "high":
			totalDeductions += 15.0
		case "medium":
			totalDeductions += 8.0
		case "low":
			totalDeductions += 3.0
		}
	}
	
	score := 100.0 - totalDeductions
	if score < 0 {
		score = 0
	}
	
	return score
}