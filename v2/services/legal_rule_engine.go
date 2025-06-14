package services

import (
	"log"
	"strings"
)

// LegalRuleEngine handles legal rule application and cause of action determination
type LegalRuleEngine struct {
	FCRARules          []FCRARule
	CauseOfActionRules []CauseOfActionRule
	DamageRules        []DamageRule
}

// FCRARule represents a specific FCRA violation rule
type FCRARule struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Statute         string                    `json:"statute"`
	Elements        []string                  `json:"elements"`
	RequiredFacts   []string                  `json:"requiredFacts"`
	Penalties       []string                  `json:"penalties"`
	Condition       func(*ClientCase) bool    `json:"-"`
}

// CauseOfActionRule defines how to generate causes of action
type CauseOfActionRule struct {
	ID               string                    `json:"id"`
	Title            string                    `json:"title"`
	StatutoryBasis   string                    `json:"statutoryBasis"`
	Elements         []string                  `json:"elements"`
	FactRequirements []string                  `json:"factRequirements"`
	Applicability    func(*ClientCase) bool    `json:"-"`
}

// DamageRule defines damage calculation rules
type DamageRule struct {
	ID           string                    `json:"id"`
	Type         string                    `json:"type"`
	Description  string                    `json:"description"`
	Amount       string                    `json:"amount"`
	Condition    func(*ClientCase) bool    `json:"-"`
}

// CauseOfAction represents a generated cause of action
type CauseOfAction struct {
	Title           string   `json:"title"`
	StatutoryBasis  string   `json:"statutoryBasis"`
	Elements        []string `json:"elements"`
	FactBasis       []string `json:"factBasis"`
	Confidence      float64  `json:"confidence"`
	Damages         []string `json:"damages"`
}

// NewLegalRuleEngine creates a new legal rule engine
func NewLegalRuleEngine() *LegalRuleEngine {
	engine := &LegalRuleEngine{
		FCRARules:          []FCRARule{},
		CauseOfActionRules: []CauseOfActionRule{},
		DamageRules:        []DamageRule{},
	}
	
	engine.loadFCRARules()
	engine.loadCauseOfActionRules()
	engine.loadDamageRules()
	
	log.Printf("[LEGAL_RULE_ENGINE] Initialized with %d FCRA rules, %d cause of action rules, %d damage rules",
		len(engine.FCRARules), len(engine.CauseOfActionRules), len(engine.DamageRules))
	
	return engine
}

// DetermineCausesOfAction analyzes the case and determines applicable causes of action
func (lre *LegalRuleEngine) DetermineCausesOfAction(clientCase *ClientCase) []CauseOfAction {
	var causes []CauseOfAction
	
	log.Printf("[LEGAL_RULE_ENGINE] Analyzing case for client: %s", clientCase.ClientName)
	
	// Evaluate each cause of action rule
	for _, rule := range lre.CauseOfActionRules {
		if rule.Applicability(clientCase) {
			cause := lre.buildCauseOfAction(rule, clientCase)
			causes = append(causes, cause)
			log.Printf("[LEGAL_RULE_ENGINE] Added cause of action: %s (confidence: %.2f)", cause.Title, cause.Confidence)
		}
	}
	
	// If no specific causes found, add default FCRA violation
	if len(causes) == 0 {
		defaultCause := lre.createDefaultFCRACause(clientCase)
		causes = append(causes, defaultCause)
		log.Printf("[LEGAL_RULE_ENGINE] Added default FCRA cause of action")
	}
	
	return causes
}

// buildCauseOfAction constructs a cause of action from a rule and case facts
func (lre *LegalRuleEngine) buildCauseOfAction(rule CauseOfActionRule, clientCase *ClientCase) CauseOfAction {
	// Calculate confidence based on available facts
	confidence := lre.calculateConfidence(rule, clientCase)
	
	// Determine applicable damages
	damages := lre.determineApplicableDamages(clientCase)
	
	// Build fact basis from client case
	factBasis := lre.extractFactBasis(rule, clientCase)
	
	return CauseOfAction{
		Title:          rule.Title,
		StatutoryBasis: rule.StatutoryBasis,
		Elements:       rule.Elements,
		FactBasis:      factBasis,
		Confidence:     confidence,
		Damages:        damages,
	}
}

// calculateConfidence determines how well the case facts support the cause of action
func (lre *LegalRuleEngine) calculateConfidence(rule CauseOfActionRule, clientCase *ClientCase) float64 {
	requiredFactsFound := 0
	totalRequiredFacts := len(rule.FactRequirements)
	
	if totalRequiredFacts == 0 {
		return 0.8 // Default confidence if no specific requirements
	}
	
	for _, requirement := range rule.FactRequirements {
		if lre.hasRequiredFact(requirement, clientCase) {
			requiredFactsFound++
		}
	}
	
	baseConfidence := float64(requiredFactsFound) / float64(totalRequiredFacts)
	
	// Boost confidence if we have strong supporting evidence
	if len(clientCase.FraudDetails) > 0 {
		baseConfidence += 0.1
	}
	if len(clientCase.CreditBureauInteractions) > 0 {
		baseConfidence += 0.1
	}
	if clientCase.ClientName != "" && clientCase.ResidenceLocation != "" {
		baseConfidence += 0.05
	}
	
	// Cap at 0.95 to account for legal uncertainties
	if baseConfidence > 0.95 {
		baseConfidence = 0.95
	}
	
	return baseConfidence
}

// hasRequiredFact checks if the client case contains a required fact
func (lre *LegalRuleEngine) hasRequiredFact(requirement string, clientCase *ClientCase) bool {
	switch requirement {
	case "client_identity":
		return clientCase.ClientName != ""
	case "credit_report_use":
		return len(clientCase.CreditBureauInteractions) > 0
	case "inaccurate_information":
		return len(clientCase.FraudDetailsStructured) > 0 || clientCase.FraudDetails != ""
	case "dispute_filed":
		for _, interaction := range clientCase.CreditBureauInteractions {
			if strings.Contains(strings.ToLower(interaction.Type), "dispute") {
				return true
			}
		}
		return false
	case "inadequate_investigation":
		for _, interaction := range clientCase.CreditBureauInteractions {
			if strings.Contains(strings.ToLower(interaction.Response), "reinvestigation") ||
			   strings.Contains(strings.ToLower(interaction.Response), "investigation") {
				return true
			}
		}
		return false
	case "willful_violation":
		// Look for evidence of willful conduct
		return lre.hasEvidenceOfWillfulViolation(clientCase)
	case "damages_suffered":
		return len(clientCase.FraudDetailsStructured) > 0 || clientCase.FraudDetails != "" || clientCase.EstimatedDamages > 0
	default:
		return false
	}
}

// hasEvidenceOfWillfulViolation looks for evidence of willful FCRA violations
func (lre *LegalRuleEngine) hasEvidenceOfWillfulViolation(clientCase *ClientCase) bool {
	// Look for patterns that suggest willful violations
	for _, interaction := range clientCase.CreditBureauInteractions {
		response := strings.ToLower(interaction.Response)
		if strings.Contains(response, "refused") ||
		   strings.Contains(response, "ignored") ||
		   strings.Contains(response, "no response") {
			return true
		}
	}
	return false
}

// determineApplicableDamages determines what damages can be claimed
func (lre *LegalRuleEngine) determineApplicableDamages(clientCase *ClientCase) []string {
	var damages []string
	
	for _, rule := range lre.DamageRules {
		if rule.Condition(clientCase) {
			damages = append(damages, rule.Description)
		}
	}
	
	// Always include basic FCRA damages
	if len(damages) == 0 {
		damages = []string{
			"Statutory damages under 15 U.S.C. § 1681n",
			"Actual damages",
			"Attorney's fees and costs",
		}
	}
	
	return damages
}

// extractFactBasis extracts the factual basis for the cause of action
func (lre *LegalRuleEngine) extractFactBasis(rule CauseOfActionRule, clientCase *ClientCase) []string {
	var facts []string
	
	if clientCase.ClientName != "" {
		facts = append(facts, "Plaintiff is a consumer under the FCRA")
	}
	
	if len(clientCase.FraudDetailsStructured) > 0 || clientCase.FraudDetails != "" {
		facts = append(facts, "Inaccurate information appeared on Plaintiff's credit report")
	}
	
	for _, interaction := range clientCase.CreditBureauInteractions {
		if strings.Contains(strings.ToLower(interaction.Type), "dispute") {
			facts = append(facts, "Plaintiff disputed the inaccurate information")
		}
	}
	
	if len(clientCase.Defendants) > 0 {
		facts = append(facts, "Defendants failed to conduct reasonable investigation")
	}
	
	return facts
}

// createDefaultFCRACause creates a default FCRA cause of action when no specific rules apply
func (lre *LegalRuleEngine) createDefaultFCRACause(clientCase *ClientCase) CauseOfAction {
	return CauseOfAction{
		Title:          "Violation of the Fair Credit Reporting Act",
		StatutoryBasis: "15 U.S.C. § 1681e(b) and § 1681i",
		Elements: []string{
			"Plaintiff is a consumer within the meaning of the FCRA",
			"Defendants are consumer reporting agencies and/or furnishers of information",
			"Defendants failed to follow reasonable procedures to assure maximum possible accuracy",
			"Defendants failed to conduct reasonable reinvestigation of disputed information",
			"As a result, Plaintiff suffered damages",
		},
		FactBasis: []string{
			"Plaintiff discovered inaccurate information on credit report",
			"Plaintiff disputed the inaccurate information",
			"Defendants failed to properly investigate",
		},
		Confidence: 0.7,
		Damages: []string{
			"Statutory damages",
			"Actual damages",
			"Attorney's fees and costs",
		},
	}
}

// loadFCRARules loads the FCRA violation rules
func (lre *LegalRuleEngine) loadFCRARules() {
	lre.FCRARules = []FCRARule{
		{
			ID:      "fcra_1681e_accuracy",
			Name:    "Failure to Follow Reasonable Procedures",
			Statute: "15 U.S.C. § 1681e(b)",
			Elements: []string{
				"Defendant is a consumer reporting agency",
				"Defendant failed to follow reasonable procedures to assure maximum possible accuracy",
				"The procedures concerned information about a consumer",
				"Plaintiff suffered damages as a result",
			},
			RequiredFacts: []string{"credit_report_use", "inaccurate_information"},
			Penalties:     []string{"Actual damages", "Statutory damages up to $1,000", "Attorney's fees"},
			Condition: func(cc *ClientCase) bool {
				return len(cc.FraudDetailsStructured) > 0 || cc.FraudDetails != ""
			},
		},
		{
			ID:      "fcra_1681i_reinvestigation",
			Name:    "Failure to Reinvestigate",
			Statute: "15 U.S.C. § 1681i",
			Elements: []string{
				"Plaintiff disputed information with a consumer reporting agency",
				"The dispute was communicated in writing",
				"Defendant failed to conduct reasonable reinvestigation",
				"Defendant failed to record current status of disputed information",
			},
			RequiredFacts: []string{"dispute_filed", "inadequate_investigation"},
			Penalties:     []string{"Actual damages", "Statutory damages", "Attorney's fees"},
			Condition: func(cc *ClientCase) bool {
				for _, interaction := range cc.CreditBureauInteractions {
					if strings.Contains(strings.ToLower(interaction.Type), "dispute") {
						return true
					}
				}
				return false
			},
		},
	}
}

// loadCauseOfActionRules loads the cause of action generation rules
func (lre *LegalRuleEngine) loadCauseOfActionRules() {
	lre.CauseOfActionRules = []CauseOfActionRule{
		{
			ID:             "fcra_willful_violation",
			Title:          "Willful Violation of the Fair Credit Reporting Act",
			StatutoryBasis: "15 U.S.C. § 1681n",
			Elements: []string{
				"Defendants willfully failed to comply with the requirements of the FCRA",
				"Defendants' violations were knowing and intentional",
				"Plaintiff suffered damages as a direct result",
			},
			FactRequirements: []string{"willful_violation", "damages_suffered"},
			Applicability: func(cc *ClientCase) bool {
				// Apply if there's evidence of willful conduct
				return len(cc.CreditBureauInteractions) > 0 && (len(cc.FraudDetailsStructured) > 0 || cc.FraudDetails != "")
			},
		},
		{
			ID:             "fcra_negligent_violation",
			Title:          "Negligent Violation of the Fair Credit Reporting Act",
			StatutoryBasis: "15 U.S.C. § 1681o",
			Elements: []string{
				"Defendants negligently failed to comply with the requirements of the FCRA",
				"Plaintiff suffered actual damages as a result",
			},
			FactRequirements: []string{"credit_report_use", "inaccurate_information", "damages_suffered"},
			Applicability: func(cc *ClientCase) bool {
				// Apply for all cases with inaccurate information
				return len(cc.FraudDetailsStructured) > 0 || cc.FraudDetails != ""
			},
		},
		{
			ID:             "fcra_reinvestigation_failure",
			Title:          "Failure to Conduct Reasonable Reinvestigation",
			StatutoryBasis: "15 U.S.C. § 1681i",
			Elements: []string{
				"Plaintiff disputed inaccurate information in writing",
				"Defendants failed to conduct reasonable reinvestigation",
				"Defendants failed to delete or modify inaccurate information",
			},
			FactRequirements: []string{"dispute_filed", "inadequate_investigation"},
			Applicability: func(cc *ClientCase) bool {
				// Apply if disputes were filed
				for _, interaction := range cc.CreditBureauInteractions {
					if strings.Contains(strings.ToLower(interaction.Type), "dispute") {
						return true
					}
				}
				return false
			},
		},
	}
}

// loadDamageRules loads the damage calculation rules
func (lre *LegalRuleEngine) loadDamageRules() {
	lre.DamageRules = []DamageRule{
		{
			ID:          "statutory_willful",
			Type:        "statutory",
			Description: "Statutory damages for willful violations ($100-$1,000)",
			Amount:      "$100-$1,000",
			Condition: func(cc *ClientCase) bool {
				// Available for willful violations
				return len(cc.CreditBureauInteractions) > 0
			},
		},
		{
			ID:          "actual_damages",
			Type:        "actual",
			Description: "Actual damages including emotional distress and credit harm",
			Amount:      "To be determined at trial",
			Condition: func(cc *ClientCase) bool {
				// Always available
				return true
			},
		},
		{
			ID:          "punitive_damages",
			Type:        "punitive",
			Description: "Punitive damages for egregious conduct",
			Amount:      "To be determined at trial",
			Condition: func(cc *ClientCase) bool {
				// Available for willful violations with egregious conduct
				return len(cc.FraudDetailsStructured) > 1 // Multiple fraud instances
			},
		},
		{
			ID:          "attorney_fees",
			Type:        "fees",
			Description: "Reasonable attorney's fees and costs",
			Amount:      "Actual costs incurred",
			Condition: func(cc *ClientCase) bool {
				// Always available under FCRA
				return true
			},
		},
	}
}