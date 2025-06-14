package services

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// SectionType represents different types of document sections
type SectionType string

const (
	SectionTypeHeader         SectionType = "header"
	SectionTypeParties        SectionType = "parties"
	SectionTypeCausesOfAction SectionType = "causes_of_action"
	SectionTypeFacts          SectionType = "facts"
	SectionTypeDamages        SectionType = "damages"
	SectionTypePrayer         SectionType = "prayer"
)

// TemplateEngine manages dynamic legal document generation
type TemplateEngine struct {
	Templates       map[string]*DocumentTemplate
	RuleEngine      *LegalRuleEngine
	Formatter       *LegalDocumentFormatter
	Validator       *DocumentValidator
}

// DocumentTemplate represents a legal document template
type DocumentTemplate struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Sections        []TemplateSection `json:"sections"`
	RequiredFields  []string          `json:"requiredFields"`
	OptionalFields  []string          `json:"optionalFields"`
	LegalRules      []LegalRule       `json:"legalRules"`
}

// TemplateSection represents a section within a document template
type TemplateSection struct {
	Name            string      `json:"name"`
	Type            SectionType `json:"type"`
	ConditionalLogic string     `json:"conditionalLogic"`
	ContentTemplate string      `json:"contentTemplate"`
	Required        bool        `json:"required"`
	Order           int         `json:"order"`
}

// LegalRule defines legal requirements for document generation
type LegalRule struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Statute         string   `json:"statute"`
	Elements        []string `json:"elements"`
	RequiredFacts   []string `json:"requiredFacts"`
	Penalties       []string `json:"penalties"`
	ApplicableWhen  string   `json:"applicableWhen"`
}

// GeneratedDocument represents a dynamically generated legal document
type GeneratedDocument struct {
	Title           string                 `json:"title"`
	Content         string                 `json:"content"`
	Sections        []GeneratedSection     `json:"sections"`
	Metadata        DocumentMetadata       `json:"metadata"`
	ValidationIssues []ValidationIssue     `json:"validationIssues"`
}

// GeneratedSection represents a generated section of the document
type GeneratedSection struct {
	Name            string      `json:"name"`
	Type            SectionType `json:"type"`
	Content         string      `json:"content"`
	SourceFacts     []string    `json:"sourceFacts"`
	Confidence      float64     `json:"confidence"`
}

// DocumentMetadata contains metadata about the generated document
type DocumentMetadata struct {
	GeneratedAt     time.Time `json:"generatedAt"`
	TemplateID      string    `json:"templateId"`
	TemplateVersion string    `json:"templateVersion"`
	ClientCaseID    string    `json:"clientCaseId"`
	WordCount       int       `json:"wordCount"`
	Completeness    float64   `json:"completeness"`
}

// ValidationIssue represents a potential issue with the generated document
type ValidationIssue struct {
	Type        string `json:"type"`
	Section     string `json:"section"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Suggestion  string `json:"suggestion"`
}

// NewTemplateEngine creates a new template engine instance
func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		Templates:  make(map[string]*DocumentTemplate),
		RuleEngine: NewLegalRuleEngine(),
		Formatter:  NewLegalDocumentFormatter(),
		Validator:  NewDocumentValidator(),
	}
	
	// Load default templates
	engine.loadDefaultTemplates()
	
	log.Printf("[TEMPLATE_ENGINE] Initialized with %d templates", len(engine.Templates))
	return engine
}

// GenerateDocument creates a legal document from extracted case data
func (te *TemplateEngine) GenerateDocument(templateID string, clientCase *ClientCase) (*GeneratedDocument, error) {
	template, exists := te.Templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}
	
	log.Printf("[TEMPLATE_ENGINE] Generating document using template: %s for client: %s", templateID, clientCase.ClientName)
	
	// Apply legal rules to determine applicable causes of action
	applicableCauses := te.RuleEngine.DetermineCausesOfAction(clientCase)
	log.Printf("[TEMPLATE_ENGINE] Determined %d applicable causes of action", len(applicableCauses))
	
	// Generate document sections
	sections := make([]GeneratedSection, 0, len(template.Sections))
	var fullContent strings.Builder
	
	for _, sectionTemplate := range template.Sections {
		section, err := te.generateSection(sectionTemplate, clientCase, applicableCauses)
		if err != nil {
			log.Printf("[TEMPLATE_ENGINE] Error generating section %s: %v", sectionTemplate.Name, err)
			continue
		}
		
		if section != nil {
			sections = append(sections, *section)
			fullContent.WriteString(section.Content)
			fullContent.WriteString("\n\n")
		}
	}
	
	// Create document metadata
	metadata := DocumentMetadata{
		GeneratedAt:     time.Now(),
		TemplateID:      templateID,
		TemplateVersion: "1.0",
		ClientCaseID:    clientCase.ClientName, // Using client name as ID for now
		WordCount:       len(strings.Fields(fullContent.String())),
		Completeness:    te.calculateCompleteness(clientCase, template),
	}
	
	// Validate the generated document
	validationIssues := te.Validator.ValidateDocument(fullContent.String(), clientCase)
	
	document := &GeneratedDocument{
		Title:           fmt.Sprintf("FCRA Complaint - %s", clientCase.ClientName),
		Content:         fullContent.String(),
		Sections:        sections,
		Metadata:        metadata,
		ValidationIssues: validationIssues,
	}
	
	log.Printf("[TEMPLATE_ENGINE] Generated document: %d sections, %d words, %.1f%% complete", 
		len(sections), metadata.WordCount, metadata.Completeness*100)
	
	return document, nil
}

// generateSection creates content for a specific document section
func (te *TemplateEngine) generateSection(sectionTemplate TemplateSection, clientCase *ClientCase, causes []CauseOfAction) (*GeneratedSection, error) {
	// Check conditional logic
	if sectionTemplate.ConditionalLogic != "" {
		shouldInclude := te.evaluateConditionalLogic(sectionTemplate.ConditionalLogic, clientCase)
		if !shouldInclude {
			log.Printf("[TEMPLATE_ENGINE] Skipping section %s due to conditional logic", sectionTemplate.Name)
			return nil, nil
		}
	}
	
	var content string
	var sourceFacts []string
	confidence := 1.0
	
	switch sectionTemplate.Type {
	case SectionTypeHeader:
		content = te.generateHeaderSection(clientCase)
		sourceFacts = []string{"case_info"}
		
	case SectionTypeParties:
		content = te.generatePartiesSection(clientCase)
		sourceFacts = []string{"client_info", "defendant_info"}
		
	case SectionTypeCausesOfAction:
		content, sourceFacts = te.generateCausesOfActionSection(causes, clientCase)
		
	case SectionTypeFacts:
		content, sourceFacts = te.generateFactsSection(clientCase)
		
	case SectionTypeDamages:
		content = te.generateDamagesSection(clientCase, causes)
		sourceFacts = []string{"fraud_amounts", "damages"}
		
	case SectionTypePrayer:
		content = te.generatePrayerSection(causes)
		sourceFacts = []string{"legal_violations"}
		
	default:
		// Use template content with variable substitution
		content = te.substituteVariables(sectionTemplate.ContentTemplate, clientCase)
		sourceFacts = []string{"template"}
	}
	
	// Calculate confidence based on available data
	confidence = te.calculateSectionConfidence(sectionTemplate.Type, clientCase)
	
	return &GeneratedSection{
		Name:        sectionTemplate.Name,
		Type:        sectionTemplate.Type,
		Content:     content,
		SourceFacts: sourceFacts,
		Confidence:  confidence,
	}, nil
}

// generateHeaderSection creates the document header
func (te *TemplateEngine) generateHeaderSection(clientCase *ClientCase) string {
	header := fmt.Sprintf(`UNITED STATES DISTRICT COURT
%s

%s,
                                                    Plaintiff,
v.                                                 Case No. %s

%s,
                                                    Defendants.

COMPLAINT FOR VIOLATIONS OF THE FAIR CREDIT REPORTING ACT`,
		strings.ToUpper(clientCase.CourtJurisdiction),
		strings.ToUpper(clientCase.ClientName),
		clientCase.CaseNumber,
		te.formatDefendantsList(clientCase.Defendants))
	
	return header
}

// generatePartiesSection creates the parties section
func (te *TemplateEngine) generatePartiesSection(clientCase *ClientCase) string {
	parties := fmt.Sprintf(`PARTIES

1. Plaintiff %s is a natural person residing in %s.

2. Upon information and belief, Defendants are entities engaged in the business of consumer reporting and/or furnishing consumer information.`,
		clientCase.ClientName,
		clientCase.ResidenceLocation)
	
	// Add specific defendant information
	for i, defendant := range clientCase.Defendants {
		parties += fmt.Sprintf(`

%d. Upon information and belief, Defendant %s is a %s with its principal place of business at %s.`,
			i+3, defendant.Name, defendant.EntityType, defendant.Address)
	}
	
	return parties
}

// generateCausesOfActionSection creates the causes of action
func (te *TemplateEngine) generateCausesOfActionSection(causes []CauseOfAction, clientCase *ClientCase) (string, []string) {
	var content strings.Builder
	var sourceFacts []string
	
	content.WriteString("CAUSES OF ACTION\n\n")
	
	for i, cause := range causes {
		content.WriteString(fmt.Sprintf("COUNT %s\n", te.numberToRoman(i+1)))
		content.WriteString(fmt.Sprintf("%s\n\n", cause.Title))
		
		// Add elements of the cause of action
		for j, element := range cause.Elements {
			content.WriteString(fmt.Sprintf("%d. %s\n\n", j+1, element))
		}
		
		sourceFacts = append(sourceFacts, cause.StatutoryBasis)
	}
	
	return content.String(), sourceFacts
}

// generateFactsSection creates the factual allegations
func (te *TemplateEngine) generateFactsSection(clientCase *ClientCase) (string, []string) {
	var content strings.Builder
	var sourceFacts []string
	
	content.WriteString("FACTUAL ALLEGATIONS\n\n")
	
	paragraphNum := 1
	
	// Client background
	if clientCase.ClientName != "" {
		content.WriteString(fmt.Sprintf("%d. At all times relevant herein, Plaintiff %s was a consumer as defined by the Fair Credit Reporting Act, 15 U.S.C. § 1681 et seq.\n\n",
			paragraphNum, clientCase.ClientName))
		paragraphNum++
		sourceFacts = append(sourceFacts, "client_name")
	}
	
	// Fraud allegations
	if len(clientCase.FraudDetailsStructured) > 0 {
		content.WriteString(fmt.Sprintf("%d. Plaintiff became aware of fraudulent activity on their credit report involving unauthorized accounts and transactions.\n\n",
			paragraphNum))
		paragraphNum++
		
		for _, fraud := range clientCase.FraudDetailsStructured {
			content.WriteString(fmt.Sprintf("%d. Specifically, Plaintiff discovered fraudulent activity involving %s in the amount of approximately $%s.\n\n",
				paragraphNum, fraud.Institution, fraud.Amount))
			paragraphNum++
			sourceFacts = append(sourceFacts, fmt.Sprintf("fraud_%s", fraud.Institution))
		}
	}
	
	// Credit bureau interactions
	if len(clientCase.CreditBureauInteractions) > 0 {
		for _, interaction := range clientCase.CreditBureauInteractions {
			content.WriteString(fmt.Sprintf("%d. Plaintiff disputed the fraudulent information with %s on or about %s.\n\n",
				paragraphNum, interaction.Bureau, interaction.Date))
			paragraphNum++
			sourceFacts = append(sourceFacts, fmt.Sprintf("dispute_%s", interaction.Bureau))
		}
	}
	
	return content.String(), sourceFacts
}

// generateDamagesSection creates the damages section
func (te *TemplateEngine) generateDamagesSection(clientCase *ClientCase, causes []CauseOfAction) string {
	damages := `DAMAGES

As a direct and proximate result of Defendants' violations of the Fair Credit Reporting Act, Plaintiff has suffered and continues to suffer:

1. Actual damages including but not limited to:
   a. Time and effort spent disputing fraudulent information;
   b. Emotional distress and anxiety;
   c. Damage to credit reputation;
   d. Loss of credit opportunities.

2. Statutory damages as provided under 15 U.S.C. § 1681n and § 1681o.

3. Reasonable attorney's fees and costs as provided under 15 U.S.C. § 1681n(a)(3).`

	return damages
}

// generatePrayerSection creates the prayer for relief
func (te *TemplateEngine) generatePrayerSection(causes []CauseOfAction) string {
	prayer := `PRAYER FOR RELIEF

WHEREFORE, Plaintiff respectfully requests that this Court:

1. Enter judgment in favor of Plaintiff and against Defendants;

2. Award Plaintiff actual damages in an amount to be determined at trial;

3. Award Plaintiff statutory damages as provided under the Fair Credit Reporting Act;

4. Award Plaintiff reasonable attorney's fees and costs;

5. Grant such other relief as this Court deems just and proper.

Respectfully submitted,

_________________________
[Attorney Name]
[Attorney Title]
Attorney for Plaintiff
[Bar Number]
[Firm Name]
[Address]
[Phone]
[Email]`

	return prayer
}

// Helper methods

func (te *TemplateEngine) formatDefendantsList(defendants []Defendant) string {
	if len(defendants) == 0 {
		return "VARIOUS DEFENDANTS"
	}
	
	var names []string
	for _, defendant := range defendants {
		names = append(names, strings.ToUpper(defendant.Name))
	}
	
	if len(names) <= 3 {
		return strings.Join(names, ", ")
	}
	
	return strings.Join(names[:2], ", ") + ", et al."
}

func (te *TemplateEngine) evaluateConditionalLogic(logic string, clientCase *ClientCase) bool {
	// Simple conditional logic evaluation
	// In production, this would be a more sophisticated expression evaluator
	switch logic {
	case "has_fraud_details":
		return len(clientCase.FraudDetails) > 0
	case "has_credit_disputes":
		return len(clientCase.CreditBureauInteractions) > 0
	case "has_defendants":
		return len(clientCase.Defendants) > 0
	default:
		return true
	}
}

func (te *TemplateEngine) substituteVariables(template string, clientCase *ClientCase) string {
	content := template
	
	// Replace common variables
	content = strings.ReplaceAll(content, "{{.ClientName}}", clientCase.ClientName)
	content = strings.ReplaceAll(content, "{{.CourtJurisdiction}}", clientCase.CourtJurisdiction)
	content = strings.ReplaceAll(content, "{{.CaseNumber}}", clientCase.CaseNumber)
	content = strings.ReplaceAll(content, "{{.ResidenceLocation}}", clientCase.ResidenceLocation)
	
	return content
}

func (te *TemplateEngine) calculateSectionConfidence(sectionType SectionType, clientCase *ClientCase) float64 {
	switch sectionType {
	case SectionTypeHeader:
		if clientCase.ClientName != "" && clientCase.CourtJurisdiction != "" {
			return 0.95
		}
		return 0.7
	case SectionTypeParties:
		if clientCase.ClientName != "" && clientCase.ResidenceLocation != "" {
			return 0.9
		}
		return 0.6
	case SectionTypeFacts:
		if len(clientCase.FraudDetails) > 0 && len(clientCase.CreditBureauInteractions) > 0 {
			return 0.85
		}
		return 0.5
	default:
		return 0.8
	}
}

func (te *TemplateEngine) calculateCompleteness(clientCase *ClientCase, template *DocumentTemplate) float64 {
	completedFields := 0
	totalFields := len(template.RequiredFields)
	
	// Check required fields completion
	for _, field := range template.RequiredFields {
		switch field {
		case "client_name":
			if clientCase.ClientName != "" {
				completedFields++
			}
		case "court_jurisdiction":
			if clientCase.CourtJurisdiction != "" {
				completedFields++
			}
		case "defendants":
			if len(clientCase.Defendants) > 0 {
				completedFields++
			}
		case "fraud_details":
			if len(clientCase.FraudDetails) > 0 {
				completedFields++
			}
		default:
			completedFields++ // Assume completed if not specifically checked
		}
	}
	
	if totalFields == 0 {
		return 1.0
	}
	
	return float64(completedFields) / float64(totalFields)
}

func (te *TemplateEngine) numberToRoman(num int) string {
	romanNumerals := []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	if num > 0 && num <= len(romanNumerals) {
		return romanNumerals[num-1]
	}
	return fmt.Sprintf("%d", num)
}

// loadDefaultTemplates loads the default legal document templates
func (te *TemplateEngine) loadDefaultTemplates() {
	// FCRA Credit Card Fraud Complaint Template
	fcraTemplate := &DocumentTemplate{
		ID:          "fcra-credit-card-fraud",
		Name:        "FCRA Credit Card Fraud Complaint",
		Description: "Template for Fair Credit Reporting Act violations involving credit card fraud",
		Sections: []TemplateSection{
			{Name: "Header", Type: SectionTypeHeader, Required: true, Order: 1},
			{Name: "Parties", Type: SectionTypeParties, Required: true, Order: 2},
			{Name: "Facts", Type: SectionTypeFacts, Required: true, Order: 3, ConditionalLogic: "has_fraud_details"},
			{Name: "Causes of Action", Type: SectionTypeCausesOfAction, Required: true, Order: 4},
			{Name: "Damages", Type: SectionTypeDamages, Required: true, Order: 5},
			{Name: "Prayer", Type: SectionTypePrayer, Required: true, Order: 6},
		},
		RequiredFields: []string{"client_name", "court_jurisdiction", "defendants", "fraud_details"},
		OptionalFields: []string{"case_number", "attorney_info", "specific_damages"},
	}
	
	te.Templates[fcraTemplate.ID] = fcraTemplate
	log.Printf("[TEMPLATE_ENGINE] Loaded template: %s", fcraTemplate.Name)
}