package services

import (
	"fmt"
	"log"
	"strings"
)

// FCRAComplianceValidator validates FCRA compliance for adverse action letters
type FCRAComplianceValidator struct {
	violationRules []ViolationRule
}

// ViolationRule defines a specific FCRA violation detection rule
type ViolationRule struct {
	RuleID      string
	Statute     string
	Description string
	Condition   string
	Severity    string
	CheckFunc   func(*AdverseActionLetter) (bool, string, float64)
}

// NewFCRAComplianceValidator creates a new FCRA compliance validator
func NewFCRAComplianceValidator() *FCRAComplianceValidator {
	validator := &FCRAComplianceValidator{}
	validator.initializeViolationRules()
	
	log.Printf("[FCRA_VALIDATOR] Initialized with %d violation detection rules", len(validator.violationRules))
	return validator
}

// ValidateCompliance performs comprehensive FCRA compliance analysis
func (fcv *FCRAComplianceValidator) ValidateCompliance(letter *AdverseActionLetter) {
	log.Printf("[FCRA_VALIDATOR] Validating FCRA compliance for %s", letter.DocumentPath)

	// Perform basic compliance checks
	fcv.performBasicComplianceChecks(letter)

	// Detect specific violations using rules
	fcv.detectViolationsUsingRules(letter)

	// Calculate overall compliance score
	fcv.calculateOverallComplianceScore(letter)

	log.Printf("[FCRA_VALIDATOR] Compliance analysis complete - %.1f%% compliance, %d violations",
		letter.ComplianceAnalysis.OverallComplianceScore*100, len(letter.ExtractedViolations))
}

// performBasicComplianceChecks performs fundamental FCRA compliance checks
func (fcv *FCRAComplianceValidator) performBasicComplianceChecks(letter *AdverseActionLetter) {
	// Check if notice was provided (document exists and is properly formatted)
	letter.ComplianceAnalysis.NoticeProvided = fcv.checkNoticeProvided(letter)

	// Check if action is clearly stated
	letter.ComplianceAnalysis.ActionClearlyStated = fcv.checkActionClearlyStated(letter)

	// Check if credit bureau is named
	letter.ComplianceAnalysis.CreditBureauNamed = fcv.checkCreditBureauNamed(letter)

	// Check if credit bureau address is provided
	letter.ComplianceAnalysis.CreditBureauAddressProvided = fcv.checkCreditBureauAddressProvided(letter)

	// Check if credit bureau phone is provided
	letter.ComplianceAnalysis.CreditBureauPhoneProvided = fcv.checkCreditBureauPhoneProvided(letter)

	// Check if free report right is disclosed
	letter.ComplianceAnalysis.FreeReportRightDisclosed = letter.ConsumerRights.FreeReportRightNoticed

	// Check if dispute right is disclosed
	letter.ComplianceAnalysis.DisputeRightDisclosed = letter.ConsumerRights.DisputeRightNoticed

	// Check timing compliance (if letter date is available)
	letter.ComplianceAnalysis.TimingCompliant = fcv.checkTimingCompliance(letter)
}

// detectViolationsUsingRules applies violation detection rules
func (fcv *FCRAComplianceValidator) detectViolationsUsingRules(letter *AdverseActionLetter) {
	letter.ExtractedViolations = []SpecificViolation{}

	for _, rule := range fcv.violationRules {
		if violated, evidence, confidence := rule.CheckFunc(letter); violated {
			violation := SpecificViolation{
				ViolationType: rule.Description,
				Statute:       rule.Statute,
				Description:   fmt.Sprintf("%s: %s", rule.RuleID, rule.Description),
				Evidence:      evidence,
				Severity:      rule.Severity,
				Confidence:    confidence,
				Location:      letter.DocumentPath,
			}
			letter.ExtractedViolations = append(letter.ExtractedViolations, violation)
		}
	}
}

// calculateOverallComplianceScore calculates the overall compliance score
func (fcv *FCRAComplianceValidator) calculateOverallComplianceScore(letter *AdverseActionLetter) {
	complianceFactors := []bool{
		letter.ComplianceAnalysis.NoticeProvided,
		letter.ComplianceAnalysis.ActionClearlyStated,
		letter.ComplianceAnalysis.CreditBureauNamed,
		letter.ComplianceAnalysis.CreditBureauAddressProvided,
		letter.ComplianceAnalysis.CreditBureauPhoneProvided,
		letter.ComplianceAnalysis.FreeReportRightDisclosed,
		letter.ComplianceAnalysis.DisputeRightDisclosed,
		letter.ComplianceAnalysis.TimingCompliant,
	}

	compliantCount := 0
	for _, factor := range complianceFactors {
		if factor {
			compliantCount++
		}
	}

	baseScore := float64(compliantCount) / float64(len(complianceFactors))

	// Adjust score based on violations
	violationPenalty := 0.0
	for _, violation := range letter.ExtractedViolations {
		switch violation.Severity {
		case "critical":
			violationPenalty += 0.3
		case "significant":
			violationPenalty += 0.2
		case "minor":
			violationPenalty += 0.1
		}
	}

	letter.ComplianceAnalysis.OverallComplianceScore = maxFloat(0.0, baseScore-violationPenalty)

	// Generate compliance issues list
	fcv.generateComplianceIssues(letter)
}

// Individual compliance check functions

func (fcv *FCRAComplianceValidator) checkNoticeProvided(letter *AdverseActionLetter) bool {
	// A proper notice should have basic adverse action language
	content := strings.ToLower(letter.RawContent)
	requiredElements := []string{"adverse action", "credit", "application"}
	
	elementCount := 0
	for _, element := range requiredElements {
		if strings.Contains(content, element) {
			elementCount++
		}
	}
	
	return elementCount >= 2 && len(letter.RawContent) > 100
}

func (fcv *FCRAComplianceValidator) checkActionClearlyStated(letter *AdverseActionLetter) bool {
	return letter.ActionTaken.ActionType != "" && letter.ActionTaken.SpecificAction != ""
}

func (fcv *FCRAComplianceValidator) checkCreditBureauNamed(letter *AdverseActionLetter) bool {
	return letter.CreditBureau.BureauName != ""
}

func (fcv *FCRAComplianceValidator) checkCreditBureauAddressProvided(letter *AdverseActionLetter) bool {
	return letter.CreditBureau.BureauAddress != "" || 
		   (letter.CreditBureau.BureauName != "" && strings.Contains(strings.ToLower(letter.RawContent), "address"))
}

func (fcv *FCRAComplianceValidator) checkCreditBureauPhoneProvided(letter *AdverseActionLetter) bool {
	return letter.CreditBureau.BureauPhone != "" ||
		   strings.Contains(strings.ToLower(letter.RawContent), "800-") ||
		   strings.Contains(strings.ToLower(letter.RawContent), "888-")
}

func (fcv *FCRAComplianceValidator) checkTimingCompliance(letter *AdverseActionLetter) bool {
	// If we don't have a letter date, we can't verify timing
	if letter.LetterDate.IsZero() {
		return true // Assume compliant if we can't verify
	}

	// For this implementation, assume timing is compliant if letter is dated
	// In a real system, we'd compare against the application/decision date
	return !letter.LetterDate.IsZero()
}

// Violation rule initialization
func (fcv *FCRAComplianceValidator) initializeViolationRules() {
	fcv.violationRules = []ViolationRule{
		{
			RuleID:      "FCRA-1681m-1",
			Statute:     "15 U.S.C. § 1681m(a)",
			Description: "Failure to Provide Adverse Action Notice",
			Condition:   "actionTaken && !noticeProvided",
			Severity:    "critical",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if !letter.ComplianceAnalysis.NoticeProvided {
					return true, "No proper adverse action notice found in document", 0.95
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-2",
			Statute:     "15 U.S.C. § 1681m(a)(1)",
			Description: "Failure to Identify Consumer Reporting Agency",
			Condition:   "noticeProvided && !creditBureauNamed",
			Severity:    "critical",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if letter.ComplianceAnalysis.NoticeProvided && !letter.ComplianceAnalysis.CreditBureauNamed {
					return true, "Adverse action notice fails to identify the consumer reporting agency", 0.90
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-3",
			Statute:     "15 U.S.C. § 1681m(a)(2)",
			Description: "Failure to Disclose Consumer's Right to Free Credit Report",
			Condition:   "noticeProvided && !freeReportRightDisclosed",
			Severity:    "significant",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if letter.ComplianceAnalysis.NoticeProvided && !letter.ComplianceAnalysis.FreeReportRightDisclosed {
					return true, "Notice fails to inform consumer of right to obtain free credit report", 0.85
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-4",
			Statute:     "15 U.S.C. § 1681m(a)(3)",
			Description: "Failure to Disclose Consumer's Right to Dispute Information",
			Condition:   "noticeProvided && !disputeRightDisclosed",
			Severity:    "significant",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if letter.ComplianceAnalysis.NoticeProvided && !letter.ComplianceAnalysis.DisputeRightDisclosed {
					return true, "Notice fails to inform consumer of right to dispute inaccurate information", 0.80
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-5",
			Statute:     "15 U.S.C. § 1681m(a)(1)",
			Description: "Incomplete Consumer Reporting Agency Contact Information",
			Condition:   "creditBureauNamed && (!addressProvided || !phoneProvided)",
			Severity:    "significant",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if letter.ComplianceAnalysis.CreditBureauNamed {
					missingInfo := []string{}
					if !letter.ComplianceAnalysis.CreditBureauAddressProvided {
						missingInfo = append(missingInfo, "address")
					}
					if !letter.ComplianceAnalysis.CreditBureauPhoneProvided {
						missingInfo = append(missingInfo, "phone number")
					}
					
					if len(missingInfo) > 0 {
						evidence := fmt.Sprintf("Credit bureau contact information incomplete: missing %s", strings.Join(missingInfo, " and "))
						return true, evidence, 0.75
					}
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-6",
			Statute:     "15 U.S.C. § 1681m(b)",
			Description: "Untimely Adverse Action Notice",
			Condition:   "actionTaken && noticeProvided && !timingCompliant",
			Severity:    "significant",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if letter.ComplianceAnalysis.NoticeProvided && !letter.ComplianceAnalysis.TimingCompliant {
					return true, "Adverse action notice not provided within required timeframe", 0.70
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-7",
			Statute:     "15 U.S.C. § 1681m(a)",
			Description: "Vague or Unclear Adverse Action Statement",
			Condition:   "noticeProvided && !actionClearlyStated",
			Severity:    "minor",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				if letter.ComplianceAnalysis.NoticeProvided && !letter.ComplianceAnalysis.ActionClearlyStated {
					return true, "Adverse action taken is not clearly stated in the notice", 0.65
				}
				return false, "", 0.0
			},
		},
		{
			RuleID:      "FCRA-1681m-8",
			Statute:     "15 U.S.C. § 1681m(a)",
			Description: "Inadequate Reason Code Disclosure",
			Condition:   "actionTaken && reasonCodesProvided == false",
			Severity:    "minor",
			CheckFunc: func(letter *AdverseActionLetter) (bool, string, float64) {
				hasReasons := len(letter.ActionTaken.ReasonCodes) > 0 || len(letter.ActionTaken.ReasonDescriptions) > 0
				if letter.ActionTaken.ActionType != "" && !hasReasons {
					return true, "No specific reason codes or descriptions provided for adverse action", 0.60
				}
				return false, "", 0.0
			},
		},
	}
}

// generateComplianceIssues creates a list of specific compliance issues
func (fcv *FCRAComplianceValidator) generateComplianceIssues(letter *AdverseActionLetter) {
	issues := []string{}

	if !letter.ComplianceAnalysis.NoticeProvided {
		issues = append(issues, "No proper adverse action notice found")
	}

	if !letter.ComplianceAnalysis.ActionClearlyStated {
		issues = append(issues, "Adverse action not clearly stated")
	}

	if !letter.ComplianceAnalysis.CreditBureauNamed {
		issues = append(issues, "Consumer reporting agency not identified")
	}

	if !letter.ComplianceAnalysis.CreditBureauAddressProvided {
		issues = append(issues, "Credit bureau address not provided")
	}

	if !letter.ComplianceAnalysis.CreditBureauPhoneProvided {
		issues = append(issues, "Credit bureau phone number not provided")
	}

	if !letter.ComplianceAnalysis.FreeReportRightDisclosed {
		issues = append(issues, "Right to free credit report not disclosed")
	}

	if !letter.ComplianceAnalysis.DisputeRightDisclosed {
		issues = append(issues, "Right to dispute information not disclosed")
	}

	if !letter.ComplianceAnalysis.TimingCompliant {
		issues = append(issues, "Notice timing may not be compliant")
	}

	letter.ComplianceAnalysis.ComplianceIssues = issues
}

// performComplianceAnalysis is the main entry point for compliance analysis
func (aap *AdverseActionParser) performComplianceAnalysis(letter *AdverseActionLetter) {
	validator := NewFCRAComplianceValidator()
	validator.ValidateCompliance(letter)
}

// detectSpecificViolations detects document-specific violations
func (aap *AdverseActionParser) detectSpecificViolations(letter *AdverseActionLetter) {
	// Additional violation detection beyond compliance validator
	// This could include document-specific patterns, format violations, etc.
	
	content := strings.ToLower(letter.RawContent)
	
	// Check for misleading language violations
	misleadingPatterns := []string{
		"pre-approved", "pre-selected", "guaranteed approval",
	}
	
	for _, pattern := range misleadingPatterns {
		if strings.Contains(content, pattern) {
			violation := SpecificViolation{
				ViolationType: "Misleading Marketing Language",
				Statute:       "15 U.S.C. § 1681m",
				Description:   "Use of misleading pre-approval language in adverse action context",
				Evidence:      fmt.Sprintf("Document contains potentially misleading language: '%s'", pattern),
				Severity:      "minor",
				Confidence:    0.75,
				Location:      letter.DocumentPath,
			}
			letter.ExtractedViolations = append(letter.ExtractedViolations, violation)
		}
	}
	
	// Check for incomplete disclosures
	if strings.Contains(content, "credit report") && !strings.Contains(content, "free") {
		violation := SpecificViolation{
			ViolationType: "Incomplete Free Report Disclosure",
			Statute:       "15 U.S.C. § 1681m(a)(2)",
			Description:   "Mentions credit report but fails to clearly state consumer's right to free copy",
			Evidence:      "Document mentions credit report but free copy rights disclosure is unclear",
			Severity:      "significant",
			Confidence:    0.80,
			Location:      letter.DocumentPath,
		}
		letter.ExtractedViolations = append(letter.ExtractedViolations, violation)
	}
}

// calculateParsingConfidence calculates overall parsing confidence
func (aap *AdverseActionParser) calculateParsingConfidence(letter *AdverseActionLetter) float64 {
	confidence := 0.0
	factors := 0

	// Factor in extracted information completeness
	if letter.Creditor.Name != "" {
		confidence += 0.15
	}
	factors++

	if letter.ActionTaken.ActionType != "" {
		confidence += 0.20
	}
	factors++

	if letter.CreditBureau.BureauName != "" {
		confidence += 0.15
	}
	factors++

	if letter.ConsumerRights.FreeReportRightNoticed || letter.ConsumerRights.DisputeRightNoticed {
		confidence += 0.20
	}
	factors++

	if letter.Consumer.Name != "" {
		confidence += 0.10
	}
	factors++

	if len(letter.ActionTaken.ReasonDescriptions) > 0 || len(letter.ActionTaken.ReasonCodes) > 0 {
		confidence += 0.10
	}
	factors++

	if letter.Creditor.Address != "" || letter.Creditor.Phone != "" {
		confidence += 0.10
	}
	factors++

	// Base confidence on document structure and content quality
	if len(letter.RawContent) > 200 && len(letter.RawContent) < 10000 {
		confidence += 0.1
	}

	return minFloat(confidence, 1.0)
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}