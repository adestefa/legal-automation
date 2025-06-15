package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AdverseActionLetter represents a parsed adverse action letter
type AdverseActionLetter struct {
	DocumentPath        string                 `json:"documentPath"`
	LetterDate          time.Time              `json:"letterDate"`
	Creditor            CreditorInformation    `json:"creditor"`
	Consumer            ConsumerInformation    `json:"consumer"`
	ActionTaken         ActionDetails          `json:"actionTaken"`
	CreditBureau        CreditBureauInfo       `json:"creditBureau"`
	ConsumerRights      ConsumerRightsNotice   `json:"consumerRights"`
	ComplianceAnalysis  ComplianceAssessment   `json:"complianceAnalysis"`
	ExtractedViolations []SpecificViolation    `json:"extractedViolations"`
	RawContent          string                 `json:"rawContent"`
	ParsingConfidence   float64                `json:"parsingConfidence"`
}

type CreditorInformation struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	ContactPerson string `json:"contactPerson"`
	AccountNumber string `json:"accountNumber"`
	Department    string `json:"department"`
}

type ConsumerInformation struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	AccountNumber  string `json:"accountNumber"`
	ApplicationID  string `json:"applicationId"`
	SSN            string `json:"ssn"`
}

type ActionDetails struct {
	ActionType         string    `json:"actionType"`         // "denied", "reduced", "modified"
	ActionDate         time.Time `json:"actionDate"`
	SpecificAction     string    `json:"specificAction"`     // "credit application denied"
	ReasonCodes        []string  `json:"reasonCodes"`        // numeric reason codes
	ReasonDescriptions []string  `json:"reasonDescriptions"` // text descriptions
	CreditScore        int       `json:"creditScore"`
	ScoreRange         string    `json:"scoreRange"`
	RequestedAmount    string    `json:"requestedAmount"`
	ApprovedAmount     string    `json:"approvedAmount"`
}

type CreditBureauInfo struct {
	BureauName    string    `json:"bureauName"`
	BureauAddress string    `json:"bureauAddress"`
	BureauPhone   string    `json:"bureauPhone"`
	BureauWebsite string    `json:"bureauWebsite"`
	ReportDate    time.Time `json:"reportDate"`
	ReportNumber  string    `json:"reportNumber"`
	IsComplete    bool      `json:"isComplete"`
}

type ConsumerRightsNotice struct {
	FreeReportRightNoticed      bool   `json:"freeReportRightNoticed"`
	DisputeRightNoticed         bool   `json:"disputeRightNoticed"`
	ReinvestigationRightNoticed bool   `json:"reinvestigationRightNoticed"`
	AdditionalRightsNoticed     bool   `json:"additionalRightsNoticed"`
	RightsText                  string `json:"rightsText"`
	ComplianceScore             float64 `json:"complianceScore"`
}

type ComplianceAssessment struct {
	NoticeProvided              bool    `json:"noticeProvided"`
	ActionClearlyStated         bool    `json:"actionClearlyStated"`
	CreditBureauNamed           bool    `json:"creditBureauNamed"`
	CreditBureauAddressProvided bool    `json:"creditBureauAddressProvided"`
	CreditBureauPhoneProvided   bool    `json:"creditBureauPhoneProvided"`
	FreeReportRightDisclosed    bool    `json:"freeReportRightDisclosed"`
	DisputeRightDisclosed       bool    `json:"disputeRightDisclosed"`
	TimingCompliant             bool    `json:"timingCompliant"`
	OverallComplianceScore      float64 `json:"overallComplianceScore"`
	ComplianceIssues            []string `json:"complianceIssues"`
}

type SpecificViolation struct {
	ViolationType string  `json:"violationType"`
	Statute       string  `json:"statute"`
	Description   string  `json:"description"`
	Evidence      string  `json:"evidence"`
	Severity      string  `json:"severity"` // "critical", "significant", "minor"
	Confidence    float64 `json:"confidence"`
	Location      string  `json:"location"`
}

// AdverseActionParser handles parsing and analysis of adverse action letters
type AdverseActionParser struct {
	patterns map[string]interface{}
}

// NewAdverseActionParser creates a new adverse action parser
func NewAdverseActionParser() (*AdverseActionParser, error) {
	parser := &AdverseActionParser{
		patterns: make(map[string]interface{}),
	}

	if err := parser.loadPatterns(); err != nil {
		return nil, fmt.Errorf("failed to load adverse action patterns: %w", err)
	}

	log.Printf("[ADVERSE_ACTION_PARSER] Initialized with comprehensive extraction patterns")
	return parser, nil
}

// loadPatterns loads adverse action parsing patterns
func (aap *AdverseActionParser) loadPatterns() error {
	// Load from enhanced legal patterns
	data, err := os.ReadFile("./config/legal_patterns_enhanced.json")
	if err != nil {
		// Fallback to original patterns if enhanced version doesn't exist
		data, err = os.ReadFile("./config/legal_patterns.json")
		if err != nil {
			return fmt.Errorf("failed to read legal patterns: %w", err)
		}
	}

	if err := json.Unmarshal(data, &aap.patterns); err != nil {
		return fmt.Errorf("failed to parse patterns JSON: %w", err)
	}

	return nil
}

// ParseAdverseActionLetter performs comprehensive parsing of an adverse action letter
func (aap *AdverseActionParser) ParseAdverseActionLetter(documentPath, content string) (*AdverseActionLetter, error) {
	log.Printf("[ADVERSE_ACTION_PARSER] Parsing adverse action letter: %s (%d chars)", documentPath, len(content))

	letter := &AdverseActionLetter{
		DocumentPath:        documentPath,
		RawContent:         content,
		ExtractedViolations: []SpecificViolation{},
		LetterDate:         time.Now(), // Default, will be overridden if found
	}

	// Phase 1: Identify and validate as adverse action letter
	if !aap.isAdverseActionLetter(content) {
		return nil, fmt.Errorf("document does not appear to be an adverse action letter")
	}

	// Phase 2: Extract core components
	aap.extractLetterDate(content, letter)
	aap.extractCreditorInformation(content, letter)
	aap.extractConsumerInformation(content, letter)
	aap.extractActionDetails(content, letter)
	aap.extractCreditBureauInfo(content, letter)
	aap.extractConsumerRights(content, letter)

	// Phase 3: Perform compliance analysis
	aap.performComplianceAnalysis(letter)

	// Phase 4: Detect specific violations
	aap.detectSpecificViolations(letter)

	// Phase 5: Calculate overall parsing confidence
	letter.ParsingConfidence = aap.calculateParsingConfidence(letter)

	log.Printf("[ADVERSE_ACTION_PARSER] Parsing complete - %.1f%% confidence, %d violations detected",
		letter.ParsingConfidence*100, len(letter.ExtractedViolations))

	return letter, nil
}

// isAdverseActionLetter validates that the document is an adverse action letter
func (aap *AdverseActionParser) isAdverseActionLetter(content string) bool {
	contentLower := strings.ToLower(content)

	// Check for adverse action indicators
	adverseActionIndicators := []string{
		"adverse action",
		"notice of adverse action",
		"credit decision",
		"application.*denied",
		"application.*declined",
		"unable.*approve",
		"fair credit reporting act",
		"consumer reporting agency",
	}

	indicatorCount := 0
	for _, indicator := range adverseActionIndicators {
		re := regexp.MustCompile("(?i)" + indicator)
		if re.MatchString(contentLower) {
			indicatorCount++
		}
	}

	// Must have at least 2 indicators to be considered adverse action letter
	return indicatorCount >= 2
}

// extractLetterDate extracts the date of the letter
func (aap *AdverseActionParser) extractLetterDate(content string, letter *AdverseActionLetter) {
	datePatterns := []string{
		`(?i)(?:date|dated)[:.\s]*([A-Z][a-z]+ \d{1,2}, \d{4})`,
		`(\d{1,2}/\d{1,2}/\d{4})`,
		`(\d{4}-\d{1,2}-\d{1,2})`,
		`([A-Z][a-z]+ \d{1,2}, \d{4})`, // "January 15, 2024"
	}

	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			if parsedDate, err := aap.parseDate(matches[1]); err == nil {
				letter.LetterDate = parsedDate
				return
			}
		}
	}
}

// extractCreditorInformation extracts creditor details
func (aap *AdverseActionParser) extractCreditorInformation(content string, letter *AdverseActionLetter) {
	lines := strings.Split(content, "\n")

	// Extract creditor name (often in header or letterhead)
	letter.Creditor.Name = aap.extractCreditorName(content, lines)

	// Extract creditor address
	letter.Creditor.Address = aap.extractCreditorAddress(content)

	// Extract creditor phone
	letter.Creditor.Phone = aap.extractCreditorPhone(content)

	// Extract contact person/department
	letter.Creditor.ContactPerson = aap.extractContactPerson(content)
	letter.Creditor.Department = aap.extractDepartment(content)

	// Extract account number
	letter.Creditor.AccountNumber = aap.extractAccountNumber(content)
}

// extractCreditorName extracts the creditor's name
func (aap *AdverseActionParser) extractCreditorName(content string, lines []string) string {
	// Try common patterns
	namePatterns := []string{
		`(?i)from:\s*([A-Z][A-Za-z\s&.,]+(?:bank|credit|financial|card|capital|chase|wells|citi))`,
		`(?i)sincerely,\s*([A-Z][A-Za-z\s&.,]+)`,
		`(?i)((?:chase|wells fargo|bank of america|capital one|citibank|discover|american express|barclays)[A-Za-z\s&.,]*)`,
	}

	for _, pattern := range namePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 && len(strings.TrimSpace(matches[1])) > 2 {
			return strings.TrimSpace(matches[1])
		}
	}

	// Try first few lines for company name
	for i := 0; i < minInt(5, len(lines)); i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) > 5 && len(line) < 60 && aap.looksLikeCompanyName(line) {
			return line
		}
	}

	return ""
}

// extractCreditorAddress extracts the creditor's address
func (aap *AdverseActionParser) extractCreditorAddress(content string) string {
	addressPatterns := []string{
		`([0-9]+\s+[A-Za-z\s]+(?:Street|St|Avenue|Ave|Road|Rd|Boulevard|Blvd|Lane|Ln|Drive|Dr|Court|Ct)[A-Za-z\s,]*[A-Z]{2}\s+[0-9]{5}(?:-[0-9]{4})?)`,
		`(P\.?O\.?\s+Box\s+[0-9]+[A-Za-z\s,]*[A-Z]{2}\s+[0-9]{5}(?:-[0-9]{4})?)`,
		`([A-Za-z\s]+,\s*[A-Z]{2}\s+[0-9]{5}(?:-[0-9]{4})?)`,
	}

	for _, pattern := range addressPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractCreditorPhone extracts the creditor's phone number
func (aap *AdverseActionParser) extractCreditorPhone(content string) string {
	phonePatterns := []string{
		`(?i)(?:phone|tel|call|contact):\s*([0-9\-\(\)\s\.]{10,})`,
		`(?i)(?:customer service|questions).*?([0-9\-\(\)\s\.]{10,})`,
		`(\([0-9]{3}\)\s*[0-9]{3}[\-\.\s]*[0-9]{4})`,
		`([0-9]{3}[\-\.][0-9]{3}[\-\.][0-9]{4})`,
	}

	for _, pattern := range phonePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractContactPerson extracts contact person name
func (aap *AdverseActionParser) extractContactPerson(content string) string {
	contactPatterns := []string{
		`(?i)sincerely,\s*([A-Z][a-z]+ [A-Z][a-z]+)`,
		`(?i)contact:\s*([A-Z][a-z]+ [A-Z][a-z]+)`,
		`(?i)regards,\s*([A-Z][a-z]+ [A-Z][a-z]+)`,
	}

	for _, pattern := range contactPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractDepartment extracts department information
func (aap *AdverseActionParser) extractDepartment(content string) string {
	deptPatterns := []string{
		`(?i)(credit (?:card )?(?:services|department))`,
		`(?i)(underwriting (?:department|team))`,
		`(?i)(customer (?:service|care))`,
		`(?i)(risk (?:management|assessment))`,
	}

	for _, pattern := range deptPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractAccountNumber extracts account or application numbers
func (aap *AdverseActionParser) extractAccountNumber(content string) string {
	accountPatterns := []string{
		`(?i)(?:account|application|reference).*?(?:number|#|no)\.?\s*([A-Z0-9\-]{6,})`,
		`(?i)(?:account|app|ref)\.?\s*#?\s*([A-Z0-9\-]{6,})`,
	}

	for _, pattern := range accountPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractConsumerInformation extracts consumer details
func (aap *AdverseActionParser) extractConsumerInformation(content string, letter *AdverseActionLetter) {
	// Extract consumer name
	consumerPatterns := []string{
		`(?i)dear\s+([A-Z][a-z]+ [A-Z][a-z]+)`,
		`(?i)(?:to|for):\s*([A-Z][a-z]+ [A-Z][a-z]+)`,
		`(?i)applicant:\s*([A-Z][a-z]+ [A-Z][a-z]+)`,
	}

	for _, pattern := range consumerPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			letter.Consumer.Name = strings.TrimSpace(matches[1])
			break
		}
	}

	// Extract consumer address (if different patterns from creditor)
	consumerAddressPatterns := []string{
		`(?i)(?:dear|to)\s+[A-Z][a-z]+ [A-Z][a-z]+\s+([0-9]+\s+[A-Za-z\s]+[A-Z]{2}\s+[0-9]{5})`,
	}

	for _, pattern := range consumerAddressPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			letter.Consumer.Address = strings.TrimSpace(matches[1])
			break
		}
	}
}

// extractActionDetails extracts details about the adverse action taken
func (aap *AdverseActionParser) extractActionDetails(content string, letter *AdverseActionLetter) {
	contentLower := strings.ToLower(content)

	// Determine action type
	if strings.Contains(contentLower, "denied") || strings.Contains(contentLower, "decline") {
		letter.ActionTaken.ActionType = "denied"
	} else if strings.Contains(contentLower, "reduced") || strings.Contains(contentLower, "lowered") {
		letter.ActionTaken.ActionType = "reduced"
	} else if strings.Contains(contentLower, "modified") || strings.Contains(contentLower, "changed") {
		letter.ActionTaken.ActionType = "modified"
	}

	// Extract specific action description
	actionPatterns := []string{
		`(?i)(your (?:application|request)[^.]*(?:denied|declined|not approved)[^.]*)`,
		`(?i)(we (?:are unable|cannot)[^.]*(?:approve|accept)[^.]*)`,
		`(?i)(your (?:credit limit|line of credit)[^.]*(?:reduced|lowered|decreased)[^.]*)`,
	}

	for _, pattern := range actionPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			letter.ActionTaken.SpecificAction = strings.TrimSpace(matches[1])
			break
		}
	}

	// Extract reason codes and descriptions
	letter.ActionTaken.ReasonCodes = aap.extractReasonCodes(content)
	letter.ActionTaken.ReasonDescriptions = aap.extractReasonDescriptions(content)

	// Extract credit score if mentioned
	scorePattern := `(?i)(?:credit score|score).*?(\d{3})`
	re := regexp.MustCompile(scorePattern)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		if score, err := strconv.Atoi(matches[1]); err == nil {
			letter.ActionTaken.CreditScore = score
		}
	}

	// Extract requested/approved amounts
	amountPatterns := []string{
		`(?i)(?:requested|applied for).*?\$([0-9,]+)`,
		`(?i)(?:approved|granted).*?\$([0-9,]+)`,
	}

	for _, pattern := range amountPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			if strings.Contains(strings.ToLower(pattern), "request") {
				letter.ActionTaken.RequestedAmount = matches[1]
			} else {
				letter.ActionTaken.ApprovedAmount = matches[1]
			}
		}
	}
}

// extractReasonCodes extracts numeric reason codes
func (aap *AdverseActionParser) extractReasonCodes(content string) []string {
	var codes []string
	
	// Look for reason code patterns
	codePatterns := []string{
		`(?i)reason codes?\s*[:\-]?\s*([0-9,\s]+)`,
		`(?i)code\s*[:\-]?\s*([0-9]+)`,
	}

	for _, pattern := range codePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			// Split on commas and spaces to get individual codes
			codeStr := strings.ReplaceAll(matches[1], ",", " ")
			fields := strings.Fields(codeStr)
			for _, field := range fields {
				if len(field) >= 1 && len(field) <= 3 {
					codes = append(codes, field)
				}
			}
		}
	}

	return codes
}

// extractReasonDescriptions extracts reason descriptions
func (aap *AdverseActionParser) extractReasonDescriptions(content string) []string {
	var descriptions []string

	reasonPatterns := []string{
		`(?i)(?:reason|because|due to)[^.]*([^.]{20,100})`,
		`(?i)(insufficient (?:credit history|income)[^.]*)`,
		`(?i)(too many (?:inquiries|accounts)[^.]*)`,
		`(?i)(debt[- ]to[- ]income ratio[^.]*)`,
		`(?i)(credit score (?:too low|below)[^.]*)`,
	}

	for _, pattern := range reasonPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				desc := strings.TrimSpace(match[1])
				if len(desc) > 10 && len(desc) < 200 {
					descriptions = append(descriptions, desc)
				}
			}
		}
	}

	return aap.removeDuplicates(descriptions)
}

// extractCreditBureauInfo extracts credit bureau information
func (aap *AdverseActionParser) extractCreditBureauInfo(content string, letter *AdverseActionLetter) {
	contentLower := strings.ToLower(content)

	// Identify credit bureau
	if strings.Contains(contentLower, "equifax") {
		letter.CreditBureau.BureauName = "Equifax Information Services LLC"
		letter.CreditBureau.BureauPhone = "800-685-1111"
		letter.CreditBureau.BureauWebsite = "www.equifax.com"
	} else if strings.Contains(contentLower, "experian") {
		letter.CreditBureau.BureauName = "Experian Information Solutions Inc"
		letter.CreditBureau.BureauPhone = "888-397-3742"
		letter.CreditBureau.BureauWebsite = "www.experian.com"
	} else if strings.Contains(contentLower, "transunion") || strings.Contains(contentLower, "trans union") {
		letter.CreditBureau.BureauName = "TransUnion LLC"
		letter.CreditBureau.BureauPhone = "800-916-8800"
		letter.CreditBureau.BureauWebsite = "www.transunion.com"
	}

	// Extract bureau address if specifically mentioned
	addressPatterns := []string{
		`(?i)equifax.*?(P\.?O\.? Box [0-9]+[^,]*Atlanta[^,]*GA[^,]*[0-9]{5})`,
		`(?i)experian.*?(P\.?O\.? Box [0-9]+[^,]*Allen[^,]*TX[^,]*[0-9]{5})`,
		`(?i)trans\s?union.*?(P\.?O\.? Box [0-9]+[^,]*Chester[^,]*PA[^,]*[0-9]{5})`,
	}

	for _, pattern := range addressPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			letter.CreditBureau.BureauAddress = strings.TrimSpace(matches[1])
			break
		}
	}

	// Check if credit bureau information is complete
	letter.CreditBureau.IsComplete = letter.CreditBureau.BureauName != "" &&
		(letter.CreditBureau.BureauAddress != "" || letter.CreditBureau.BureauPhone != "")
}

// extractConsumerRights analyzes consumer rights disclosures
func (aap *AdverseActionParser) extractConsumerRights(content string, letter *AdverseActionLetter) {
	contentLower := strings.ToLower(content)

	// Check for free report right disclosure
	freeReportIndicators := []string{
		"free copy",
		"free credit report",
		"obtain.*free.*report",
		"no charge.*report",
	}

	for _, indicator := range freeReportIndicators {
		re := regexp.MustCompile("(?i)" + indicator)
		if re.MatchString(contentLower) {
			letter.ConsumerRights.FreeReportRightNoticed = true
			break
		}
	}

	// Check for dispute right disclosure
	disputeIndicators := []string{
		"dispute.*accuracy",
		"dispute.*information",
		"contest.*information",
		"challenge.*report",
		"right to dispute",
	}

	for _, indicator := range disputeIndicators {
		re := regexp.MustCompile("(?i)" + indicator)
		if re.MatchString(contentLower) {
			letter.ConsumerRights.DisputeRightNoticed = true
			break
		}
	}

	// Check for reinvestigation right
	reinvestigationIndicators := []string{
		"reinvestigation",
		"investigate.*dispute",
		"review.*dispute",
	}

	for _, indicator := range reinvestigationIndicators {
		re := regexp.MustCompile("(?i)" + indicator)
		if re.MatchString(contentLower) {
			letter.ConsumerRights.ReinvestigationRightNoticed = true
			break
		}
	}

	// Extract rights text
	rightsPatterns := []string{
		`(?i)(you have the right[^.]*\.(?:[^.]*\.){0,3})`,
		`(?i)(under the fair credit reporting act[^.]*\.(?:[^.]*\.){0,3})`,
		`(?i)(consumer rights[^.]*\.(?:[^.]*\.){0,3})`,
	}

	for _, pattern := range rightsPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			letter.ConsumerRights.RightsText = strings.TrimSpace(matches[1])
			break
		}
	}

	// Calculate consumer rights compliance score
	rightCount := 0
	if letter.ConsumerRights.FreeReportRightNoticed {
		rightCount++
	}
	if letter.ConsumerRights.DisputeRightNoticed {
		rightCount++
	}
	if letter.ConsumerRights.ReinvestigationRightNoticed {
		rightCount++
	}
	if letter.ConsumerRights.RightsText != "" {
		rightCount++
	}

	letter.ConsumerRights.ComplianceScore = float64(rightCount) / 4.0
}

// Helper functions

func (aap *AdverseActionParser) parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"January 2, 2006",
		"Jan 2, 2006",
		"1/2/2006",
		"01/02/2006",
		"2006-01-02",
		"2006-1-2",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func (aap *AdverseActionParser) looksLikeCompanyName(line string) bool {
	line = strings.TrimSpace(line)
	
	// Check for company indicators
	companyIndicators := []string{
		"bank", "financial", "credit", "card", "capital", "chase", "wells", "citi", "discover", "american express", "barclays",
		"corp", "inc", "llc", "ltd", "company", "services",
	}

	lineLower := strings.ToLower(line)
	for _, indicator := range companyIndicators {
		if strings.Contains(lineLower, indicator) {
			return true
		}
	}

	// Check if it looks like a proper company name (title case, reasonable length)
	words := strings.Fields(line)
	if len(words) >= 2 && len(words) <= 6 {
		properCaseCount := 0
		for _, word := range words {
			if len(word) > 0 && word[0] >= 'A' && word[0] <= 'Z' {
				properCaseCount++
			}
		}
		return properCaseCount >= len(words)/2
	}

	return false
}

func (aap *AdverseActionParser) removeDuplicates(slice []string) []string {
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

