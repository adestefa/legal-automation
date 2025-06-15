package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ContentAnalyzer provides intelligent analysis of extracted document text
type ContentAnalyzer struct {
	Patterns              map[string]interface{}
	Extractors            map[string]FieldExtractor
	Validators            map[string]ValidationFunc
	DocumentClassifier    *DocumentClassifier
	AdverseActionParser   *AdverseActionParser
}

// FieldExtractor interface for extracting specific field types
type FieldExtractor interface {
	Extract(text string) (value interface{}, confidence float64, err error)
}

// ValidationFunc validates extracted data
type ValidationFunc func(value interface{}) bool

// ExtractionResult represents the result of extracting a specific field
type ExtractionResult struct {
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	Confidence float64     `json:"confidence"`
	Source     string      `json:"source"`
	Location   int         `json:"location"`
	Method     string      `json:"method"`
}

// LegalAnalysisResult contains comprehensive analysis of legal documents
type LegalAnalysisResult struct {
	ClientData              map[string]ExtractionResult `json:"clientData"`
	FraudDetails            map[string]ExtractionResult `json:"fraudDetails"`
	LegalViolations         []string                    `json:"legalViolations"`
	DocumentTypes           map[string]bool             `json:"documentTypes"`
	DocumentClassifications []DocumentClassification    `json:"documentClassifications"`
	AdverseActionLetters    []AdverseActionLetter       `json:"adverseActionLetters"`
	OverallConfidence       float64                     `json:"overallConfidence"`
	MissingFields           []string                    `json:"missingFields"`
	Suggestions             []string                    `json:"suggestions"`
}

// NewContentAnalyzer creates a new content analyzer with legal patterns
func NewContentAnalyzer() (*ContentAnalyzer, error) {
	analyzer := &ContentAnalyzer{
		Extractors: make(map[string]FieldExtractor),
		Validators: make(map[string]ValidationFunc),
	}
	
	// Load legal patterns
	err := analyzer.loadLegalPatterns()
	if err != nil {
		return nil, fmt.Errorf("failed to load legal patterns: %v", err)
	}
	
	// Initialize document classifier
	classifier, err := NewDocumentClassifier(analyzer)
	if err != nil {
		log.Printf("[CONTENT_ANALYZER] Warning: Could not initialize document classifier: %v", err)
		// Continue without classifier for backward compatibility
	} else {
		analyzer.DocumentClassifier = classifier
		log.Printf("[CONTENT_ANALYZER] Initialized with document classification engine")
	}

	// Initialize adverse action parser
	adverseActionParser, err := NewAdverseActionParser()
	if err != nil {
		log.Printf("[CONTENT_ANALYZER] Warning: Could not initialize adverse action parser: %v", err)
		// Continue without parser for backward compatibility
	} else {
		analyzer.AdverseActionParser = adverseActionParser
		log.Printf("[CONTENT_ANALYZER] Initialized with adverse action parsing engine")
	}
	
	// Initialize field extractors
	analyzer.initializeExtractors()
	
	// Initialize validators
	analyzer.initializeValidators()
	
	log.Printf("[CONTENT_ANALYZER] Initialized with legal document intelligence")
	return analyzer, nil
}

// loadLegalPatterns loads pattern definitions from JSON configuration
func (ca *ContentAnalyzer) loadLegalPatterns() error {
	patternsPath := "./config/legal_patterns_enhanced.json"
	
	// Fallback to original patterns if enhanced version doesn't exist
	if _, err := os.Stat(patternsPath); os.IsNotExist(err) {
		patternsPath = "./config/legal_patterns.json"
		log.Printf("[CONTENT_ANALYZER] Using original patterns file as fallback")
	}
	
	data, err := os.ReadFile(patternsPath)
	if err != nil {
		return fmt.Errorf("could not read patterns file: %v", err)
	}
	
	err = json.Unmarshal(data, &ca.Patterns)
	if err != nil {
		return fmt.Errorf("could not parse patterns JSON: %v", err)
	}
	
	log.Printf("[CONTENT_ANALYZER] Loaded legal patterns from %s", patternsPath)
	return nil
}

// initializeExtractors sets up field-specific extractors
func (ca *ContentAnalyzer) initializeExtractors() {
	ca.Extractors["clientName"] = &NameExtractor{patterns: ca.getPatterns("clientInfo", "namePatterns")}
	ca.Extractors["phoneNumber"] = &PhoneExtractor{patterns: ca.getPatterns("clientInfo", "phonePatterns")}
	ca.Extractors["fraudAmount"] = &AmountExtractor{patterns: ca.getPatterns("fraudDetails", "amountPatterns")}
	ca.Extractors["institution"] = &InstitutionExtractor{patterns: ca.getPatterns("fraudDetails", "institutionPatterns")}
	ca.Extractors["travelLocation"] = &TravelExtractor{patterns: ca.getPatterns("fraudDetails", "travelPatterns")}
}

// initializeValidators sets up field validators
func (ca *ContentAnalyzer) initializeValidators() {
	ca.Validators["clientName"] = func(value interface{}) bool {
		if str, ok := value.(string); ok {
			return len(strings.Fields(str)) >= 2 // At least first and last name
		}
		return false
	}
	
	ca.Validators["phoneNumber"] = func(value interface{}) bool {
		if str, ok := value.(string); ok {
			return regexp.MustCompile(`\d{10,}`).MatchString(strings.ReplaceAll(str, "[^0-9]", ""))
		}
		return false
	}
	
	ca.Validators["fraudAmount"] = func(value interface{}) bool {
		if str, ok := value.(string); ok {
			// Remove $ and commas, check if it's a valid number
			cleanAmount := strings.ReplaceAll(strings.ReplaceAll(str, "$", ""), ",", "")
			if amount, err := strconv.ParseFloat(cleanAmount, 64); err == nil {
				return amount > 0
			}
		}
		return false
	}
}

// getPatterns extracts pattern arrays from the configuration
func (ca *ContentAnalyzer) getPatterns(category, patternType string) []string {
	if categoryData, ok := ca.Patterns[category].(map[string]interface{}); ok {
		if patterns, ok := categoryData[patternType].([]interface{}); ok {
			result := make([]string, len(patterns))
			for i, pattern := range patterns {
				if str, ok := pattern.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}
	return []string{}
}

// AnalyzeLegalContent performs comprehensive analysis of legal document text
func (ca *ContentAnalyzer) AnalyzeLegalContent(text string, documentPath string) (*LegalAnalysisResult, error) {
	log.Printf("[CONTENT_ANALYZER] Analyzing document %s (%d chars)", documentPath, len(text))
	
	result := &LegalAnalysisResult{
		ClientData:              make(map[string]ExtractionResult),
		FraudDetails:            make(map[string]ExtractionResult),
		DocumentTypes:           make(map[string]bool),
		DocumentClassifications: []DocumentClassification{},
		AdverseActionLetters:    []AdverseActionLetter{},
	}
	
	// Classify document using new classification engine
	var primaryDocType DocumentType = DocumentTypeUnknown
	if ca.DocumentClassifier != nil {
		classification, err := ca.DocumentClassifier.ClassifyDocument(documentPath, text)
		if err == nil {
			result.DocumentClassifications = append(result.DocumentClassifications, *classification)
			primaryDocType = classification.PrimaryType
			log.Printf("[CONTENT_ANALYZER] Document classified as %s with %.1f%% confidence", 
				ca.DocumentClassifier.GetDocumentTypeName(classification.PrimaryType), 
				classification.Confidence*100)
		} else {
			log.Printf("[CONTENT_ANALYZER] Warning: Document classification failed: %v", err)
		}
	}

	// Process adverse action letters with specialized parser
	if primaryDocType == DocumentTypeAdverseActionLetter && ca.AdverseActionParser != nil {
		adverseActionLetter, err := ca.AdverseActionParser.ParseAdverseActionLetter(documentPath, text)
		if err == nil {
			result.AdverseActionLetters = append(result.AdverseActionLetters, *adverseActionLetter)
			
			// Merge adverse action violations into overall violations list
			for _, violation := range adverseActionLetter.ExtractedViolations {
				result.LegalViolations = append(result.LegalViolations, violation.ViolationType)
			}
			
			log.Printf("[CONTENT_ANALYZER] Adverse action letter parsed - %.1f%% confidence, %d violations",
				adverseActionLetter.ParsingConfidence*100, len(adverseActionLetter.ExtractedViolations))
		} else {
			log.Printf("[CONTENT_ANALYZER] Warning: Adverse action parsing failed: %v", err)
		}
	}
	
	// Extract client information with document type context
	ca.extractClientDataWithContext(text, result, documentPath)
	
	// Extract fraud details with document type context
	ca.extractFraudDetailsWithContext(text, result, documentPath)
	
	// Analyze legal violations with enhanced patterns
	ca.analyzeLegalViolationsEnhanced(text, result)
	
	// Determine document types (legacy support)
	ca.analyzeDocumentTypes(text, result)
	
	// Calculate overall confidence
	result.OverallConfidence = ca.calculateOverallConfidence(result)
	
	// Identify missing fields
	result.MissingFields = ca.identifyMissingFields(result)
	
	// Generate suggestions
	result.Suggestions = ca.generateSuggestions(result, documentPath)
	
	log.Printf("[CONTENT_ANALYZER] Analysis complete - %.1f%% confidence, %d missing fields", 
		result.OverallConfidence*100, len(result.MissingFields))
	
	return result, nil
}

// extractClientDataWithContext extracts client information with document type context
func (ca *ContentAnalyzer) extractClientDataWithContext(text string, result *LegalAnalysisResult, documentPath string) {
	// Determine document type from classification or path
	docType := ca.getDocumentTypeFromPath(documentPath)
	
	// Extract client name
	if extractor, ok := ca.Extractors["clientName"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.ClientData["clientName"] = ExtractionResult{
				Field:      "Client Name",
				Value:      value,
				Confidence: confidence,
				Source:     docType,
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract phone number
	if extractor, ok := ca.Extractors["phoneNumber"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.ClientData["phoneNumber"] = ExtractionResult{
				Field:      "Phone Number",
				Value:      value,
				Confidence: confidence,
				Source:     docType,
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract additional client data using enhanced patterns
	ca.extractWithPatternsEnhanced(text, "clientInfo", "addressPatterns", "clientAddress", result.ClientData, docType)
	ca.extractWithPatternsEnhanced(text, "clientInfo", "emailPatterns", "clientEmail", result.ClientData, docType)
}

// extractFraudDetailsWithContext extracts fraud information with document type context
func (ca *ContentAnalyzer) extractFraudDetailsWithContext(text string, result *LegalAnalysisResult, documentPath string) {
	docType := ca.getDocumentTypeFromPath(documentPath)
	
	// Extract fraud amount
	if extractor, ok := ca.Extractors["fraudAmount"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.FraudDetails["fraudAmount"] = ExtractionResult{
				Field:      "Fraud Amount",
				Value:      value,
				Confidence: confidence,
				Source:     docType,
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract financial institution
	if extractor, ok := ca.Extractors["institution"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.FraudDetails["institution"] = ExtractionResult{
				Field:      "Financial Institution",
				Value:      value,
				Confidence: confidence,
				Source:     docType,
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract additional fraud details using enhanced patterns
	ca.extractWithPatternsEnhanced(text, "fraudDetails", "datePatterns", "fraudDate", result.FraudDetails, docType)
	ca.extractWithPatternsEnhanced(text, "fraudDetails", "travelPatterns", "travelLocation", result.FraudDetails, docType)
}

// analyzeLegalViolationsEnhanced identifies violations using enhanced patterns
func (ca *ContentAnalyzer) analyzeLegalViolationsEnhanced(text string, result *LegalAnalysisResult) {
	textLower := strings.ToLower(text)
	violations := []string{}
	
	// Check for specific FCRA violations using enhanced patterns
	specificViolations := map[string][]string{
		"Negligent FCRA Noncompliance": ca.getEnhancedPatterns("legalViolations", "specificViolations", "negligentNoncompliance"),
		"Willful FCRA Noncompliance": ca.getEnhancedPatterns("legalViolations", "specificViolations", "willfulNoncompliance"),
		"Failure to Conduct Reasonable Investigation": ca.getEnhancedPatterns("legalViolations", "specificViolations", "failureToInvestigate"),
		"Inaccurate Credit Reporting": ca.getEnhancedPatterns("legalViolations", "specificViolations", "inaccurateReporting"),
		"Failure to Correct Disputed Information": ca.getEnhancedPatterns("legalViolations", "specificViolations", "failureToCorrect"),
		"Improper Disclosure of Credit Information": ca.getEnhancedPatterns("legalViolations", "specificViolations", "improperDisclosure"),
	}
	
	for violationType, patterns := range specificViolations {
		for _, pattern := range patterns {
			re, err := regexp.Compile("(?i)" + pattern)
			if err != nil {
				continue
			}
			if re.MatchString(textLower) {
				violations = append(violations, violationType)
				break // Only add each violation type once
			}
		}
	}
	
	// Check for general FCRA violations
	fcraViolations := ca.getPatterns("legalViolations", "fcraViolations")
	for _, violation := range fcraViolations {
		if strings.Contains(textLower, strings.ToLower(violation)) {
			violations = append(violations, violation)
		}
	}
	
	// Check for credit impact patterns
	creditImpacts := ca.getPatterns("legalViolations", "creditImpactPatterns")
	for _, impact := range creditImpacts {
		if strings.Contains(textLower, strings.ToLower(impact)) {
			violations = append(violations, "Credit Impact: " + impact)
		}
	}
	
	result.LegalViolations = ca.removeDuplicates(violations)
}

// Helper methods
func (ca *ContentAnalyzer) getDocumentTypeFromPath(documentPath string) string {
	filename := strings.ToLower(documentPath)
	if strings.Contains(filename, "attorney") || strings.Contains(filename, "atty") {
		return "Attorney Notes"
	}
	if strings.Contains(filename, "adverse") || strings.Contains(filename, "denial") {
		return "Adverse Action Letter"
	}
	if strings.Contains(filename, "summons") {
		return "Summons"
	}
	if strings.Contains(filename, "civil") {
		return "Civil Cover Sheet"
	}
	return "Unknown Document"
}

func (ca *ContentAnalyzer) getEnhancedPatterns(category, subcategory, patternType string) []string {
	if categoryData, ok := ca.Patterns[category].(map[string]interface{}); ok {
		if subData, ok := categoryData[subcategory].(map[string]interface{}); ok {
			if patterns, ok := subData[patternType].([]interface{}); ok {
				result := make([]string, len(patterns))
				for i, pattern := range patterns {
					if str, ok := pattern.(string); ok {
						result[i] = str
					}
				}
				return result
			}
		}
	}
	return []string{}
}

func (ca *ContentAnalyzer) extractWithPatternsEnhanced(text, category, patternType, field string, target map[string]ExtractionResult, docType string) {
	patterns := ca.getPatterns(category, patternType)
	
	for _, pattern := range patterns {
		re, err := regexp.Compile(`(?i)` + pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			value := strings.TrimSpace(matches[1])
			if value != "" {
				target[field] = ExtractionResult{
					Field:      field,
					Value:      value,
					Confidence: 0.8,
					Source:     docType,
					Method:     "Enhanced Pattern Matching",
					Location:   re.FindStringIndex(text)[0],
				}
				break
			}
		}
	}
}

func (ca *ContentAnalyzer) removeDuplicates(slice []string) []string {
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

// extractClientData extracts client-specific information
func (ca *ContentAnalyzer) extractClientData(text string, result *LegalAnalysisResult) {
	// Extract client name
	if extractor, ok := ca.Extractors["clientName"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.ClientData["clientName"] = ExtractionResult{
				Field:      "Client Name",
				Value:      value,
				Confidence: confidence,
				Source:     "Attorney Notes",
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract phone number
	if extractor, ok := ca.Extractors["phoneNumber"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.ClientData["phoneNumber"] = ExtractionResult{
				Field:      "Phone Number",
				Value:      value,
				Confidence: confidence,
				Source:     "Attorney Notes",
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract additional client details using direct pattern matching
	ca.extractWithPatterns(text, "clientInfo", "addressPatterns", "address", result.ClientData)
	ca.extractWithPatterns(text, "clientInfo", "emailPatterns", "email", result.ClientData)
}

// extractFraudDetails extracts fraud-related information
func (ca *ContentAnalyzer) extractFraudDetails(text string, result *LegalAnalysisResult) {
	// Extract fraud amount
	if extractor, ok := ca.Extractors["fraudAmount"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.FraudDetails["fraudAmount"] = ExtractionResult{
				Field:      "Fraud Amount",
				Value:      value,
				Confidence: confidence,
				Source:     "Attorney Notes",
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract financial institution
	if extractor, ok := ca.Extractors["institution"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.FraudDetails["institution"] = ExtractionResult{
				Field:      "Financial Institution",
				Value:      value,
				Confidence: confidence,
				Source:     "Attorney Notes",
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract travel information
	if extractor, ok := ca.Extractors["travelLocation"]; ok {
		if value, confidence, err := extractor.Extract(text); err == nil && value != nil {
			result.FraudDetails["travelLocation"] = ExtractionResult{
				Field:      "Travel Location",
				Value:      value,
				Confidence: confidence,
				Source:     "Attorney Notes",
				Method:     "Pattern Matching",
			}
		}
	}
	
	// Extract dates
	ca.extractWithPatterns(text, "fraudDetails", "datePatterns", "fraudDate", result.FraudDetails)
}

// extractWithPatterns extracts data using pattern lists
func (ca *ContentAnalyzer) extractWithPatterns(text, category, patternType, field string, target map[string]ExtractionResult) {
	patterns := ca.getPatterns(category, patternType)
	
	for _, pattern := range patterns {
		re, err := regexp.Compile(`(?i)` + pattern) // Case insensitive
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			value := strings.TrimSpace(matches[1])
			if value != "" {
				target[field] = ExtractionResult{
					Field:      field,
					Value:      value,
					Confidence: 0.8, // Default confidence for pattern matches
					Source:     "Document Text",
					Method:     "Regex Pattern",
					Location:   re.FindStringIndex(text)[0],
				}
				break // Use first match
			}
		}
	}
}

// analyzeLegalViolations identifies potential legal violations
func (ca *ContentAnalyzer) analyzeLegalViolations(text string, result *LegalAnalysisResult) {
	fcraViolations := ca.getPatterns("legalViolations", "fcraViolations")
	
	violations := []string{}
	textLower := strings.ToLower(text)
	
	for _, violation := range fcraViolations {
		if strings.Contains(textLower, strings.ToLower(violation)) {
			violations = append(violations, violation)
		}
	}
	
	result.LegalViolations = violations
}

// analyzeDocumentTypes determines what types of documents are present
func (ca *ContentAnalyzer) analyzeDocumentTypes(text string, result *LegalAnalysisResult) {
	if docTypes, ok := ca.Patterns["documentTypes"].(map[string]interface{}); ok {
		textLower := strings.ToLower(text)
		
		for docType, typeData := range docTypes {
			if indicators, ok := typeData.(map[string]interface{})["indicators"].([]interface{}); ok {
				for _, indicator := range indicators {
					if indicatorStr, ok := indicator.(string); ok {
						if strings.Contains(textLower, strings.ToLower(indicatorStr)) {
							result.DocumentTypes[docType] = true
							break
						}
					}
				}
			}
		}
	}
}

// calculateOverallConfidence calculates overall confidence in extracted data
func (ca *ContentAnalyzer) calculateOverallConfidence(result *LegalAnalysisResult) float64 {
	totalConfidence := 0.0
	totalFields := 0
	
	// Factor in client data confidence
	for _, extraction := range result.ClientData {
		totalConfidence += extraction.Confidence
		totalFields++
	}
	
	// Factor in fraud details confidence
	for _, extraction := range result.FraudDetails {
		totalConfidence += extraction.Confidence
		totalFields++
	}
	
	if totalFields == 0 {
		return 0.0
	}
	
	return (totalConfidence / float64(totalFields)) * 100
}

// identifyMissingFields identifies critical missing information
func (ca *ContentAnalyzer) identifyMissingFields(result *LegalAnalysisResult) []string {
	required := []string{"clientName", "fraudAmount", "institution"}
	missing := []string{}
	
	for _, field := range required {
		found := false
		
		// Check in client data
		if _, exists := result.ClientData[field]; exists {
			found = true
		}
		
		// Check in fraud details
		if _, exists := result.FraudDetails[field]; exists {
			found = true
		}
		
		if !found {
			missing = append(missing, field)
		}
	}
	
	return missing
}

// generateSuggestions creates actionable suggestions for improving data completeness
func (ca *ContentAnalyzer) generateSuggestions(result *LegalAnalysisResult, documentType string) []string {
	suggestions := []string{}
	
	if len(result.MissingFields) > 0 {
		suggestions = append(suggestions, "Consider adding Attorney Notes document for client information")
	}
	
	if !result.DocumentTypes["adverseAction"] {
		suggestions = append(suggestions, "Add Adverse Action letters to document credit impact")
	}
	
	if !result.DocumentTypes["civilCoverSheet"] {
		suggestions = append(suggestions, "Include Civil Cover Sheet for court information")
	}
	
	if result.OverallConfidence < 70 {
		suggestions = append(suggestions, "Review document quality - some information may be unclear")
	}
	
	return suggestions
}

// Field Extractor Implementations

// NameExtractor extracts client names from text
type NameExtractor struct {
	patterns []string
}

func (ne *NameExtractor) Extract(text string) (interface{}, float64, error) {
	for _, pattern := range ne.patterns {
		re, err := regexp.Compile(`(?i)` + pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			name := strings.TrimSpace(matches[1])
			if len(strings.Fields(name)) >= 2 { // At least first and last name
				return name, 0.9, nil
			}
		}
	}
	return nil, 0.0, fmt.Errorf("no name found")
}

// PhoneExtractor extracts phone numbers from text
type PhoneExtractor struct {
	patterns []string
}

func (pe *PhoneExtractor) Extract(text string) (interface{}, float64, error) {
	for _, pattern := range pe.patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			phone := strings.TrimSpace(matches[1])
			return phone, 0.8, nil
		}
	}
	return nil, 0.0, fmt.Errorf("no phone number found")
}

// AmountExtractor extracts monetary amounts from text
type AmountExtractor struct {
	patterns []string
}

func (ae *AmountExtractor) Extract(text string) (interface{}, float64, error) {
	for _, pattern := range ae.patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			amount := strings.TrimSpace(matches[1])
			// Validate it's a reasonable amount
			cleanAmount := strings.ReplaceAll(strings.ReplaceAll(amount, "$", ""), ",", "")
			if _, err := strconv.ParseFloat(cleanAmount, 64); err == nil {
				if !strings.HasPrefix(amount, "$") {
					amount = "$" + amount
				}
				return amount, 0.9, nil
			}
		}
	}
	return nil, 0.0, fmt.Errorf("no valid amount found")
}

// InstitutionExtractor extracts financial institution names
type InstitutionExtractor struct {
	patterns []string
}

func (ie *InstitutionExtractor) Extract(text string) (interface{}, float64, error) {
	for _, pattern := range ie.patterns {
		re, err := regexp.Compile(`(?i)` + pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			institution := strings.TrimSpace(matches[1])
			return institution, 0.8, nil
		}
	}
	return nil, 0.0, fmt.Errorf("no institution found")
}

// TravelExtractor extracts travel information from text
type TravelExtractor struct {
	patterns []string
}

func (te *TravelExtractor) Extract(text string) (interface{}, float64, error) {
	for _, pattern := range te.patterns {
		re, err := regexp.Compile(`(?i)` + pattern)
		if err != nil {
			continue
		}
		
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			travel := strings.TrimSpace(matches[1])
			return travel, 0.7, nil
		}
	}
	return nil, 0.0, fmt.Errorf("no travel information found")
}