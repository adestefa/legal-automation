package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type DefendantAnalyzer struct {
	CreditBureauDB    map[string]interface{} `json:"creditBureauDB"`
	DefendantProfiles []DefendantProfile     `json:"defendantProfiles"`
	AnalysisRules     DefendantAnalysisRules `json:"analysisRules"`
}

type DefendantProfile struct {
	DefendantID       string                 `json:"defendantId"`
	LegalName         string                 `json:"legalName"`
	BusinessType      string                 `json:"businessType"`
	CorporateType     string                 `json:"corporateType"`
	Aliases           []string               `json:"aliases"`
	ContactInfo       DefendantContactInfo   `json:"contactInfo"`
	LegalStatus       DefendantLegalStatus   `json:"legalStatus"`
	BusinessOperations DefendantOperations   `json:"businessOperations"`
	ViolationHistory  []ViolationRecord      `json:"violationHistory"`
	ServiceHistory    []ServiceRecord        `json:"serviceHistory"`
	CaseInvolvement   []CaseInvolvement      `json:"caseInvolvement"`
}

type DefendantContactInfo struct {
	RegisteredAgent    RegisteredAgentInfo `json:"registeredAgent"`
	BusinessAddress    Address             `json:"businessAddress"`
	ServiceAddress     Address             `json:"serviceAddress"`
	MailingAddress     Address             `json:"mailingAddress"`
	ContactMethods     []ContactMethod     `json:"contactMethods"`
}

type RegisteredAgentInfo struct {
	AgentName       string    `json:"agentName"`
	AgentAddress    Address   `json:"agentAddress"`
	State           string    `json:"state"`
	RegistrationDate time.Time `json:"registrationDate"`
	Status          string    `json:"status"`
}

type DefendantLegalStatus struct {
	StateOfIncorporation string              `json:"stateOfIncorporation"`
	FederalTaxID         string              `json:"federalTaxId"`
	BusinessLicenses     []BusinessLicense   `json:"businessLicenses"`
	RegulatoryFilings    []RegulatoryFiling  `json:"regulatoryFilings"`
	CorporateStatus      string              `json:"corporateStatus"`
	LastUpdated          time.Time           `json:"lastUpdated"`
}

type DefendantOperations struct {
	BusinessScope        []string            `json:"businessScope"`
	GeographicPresence   []string            `json:"geographicPresence"`
	BusinessActivities   []BusinessActivity  `json:"businessActivities"`
	RegulatoryCompliance []ComplianceRecord  `json:"regulatoryCompliance"`
	IndustryClassification string            `json:"industryClassification"`
}

type BusinessLicense struct {
	LicenseType   string    `json:"licenseType"`
	IssuingState  string    `json:"issuingState"`
	LicenseNumber string    `json:"licenseNumber"`
	IssueDate     time.Time `json:"issueDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	Status        string    `json:"status"`
}

type RegulatoryFiling struct {
	FilingType    string    `json:"filingType"`
	Agency        string    `json:"agency"`
	FilingDate    time.Time `json:"filingDate"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
}

type BusinessActivity struct {
	ActivityType    string   `json:"activityType"`
	Description     string   `json:"description"`
	States          []string `json:"states"`
	Frequency       string   `json:"frequency"`
	Volume          string   `json:"volume"`
}

type ComplianceRecord struct {
	Regulation      string    `json:"regulation"`
	ComplianceDate  time.Time `json:"complianceDate"`
	Status          string    `json:"status"`
	Violations      []string  `json:"violations"`
	Remediation     string    `json:"remediation"`
}

type ViolationRecord struct {
	ViolationType     string    `json:"violationType"`
	Statute           string    `json:"statute"`
	ViolationDate     time.Time `json:"violationDate"`
	Description       string    `json:"description"`
	CaseNumber        string    `json:"caseNumber"`
	Resolution        string    `json:"resolution"`
	PenaltyAmount     float64   `json:"penaltyAmount"`
}

type ServiceRecord struct {
	ServiceDate     time.Time `json:"serviceDate"`
	ServiceMethod   string    `json:"serviceMethod"`
	ServiceAddress  Address   `json:"serviceAddress"`
	CaseNumber      string    `json:"caseNumber"`
	ServiceStatus   string    `json:"serviceStatus"`
	ProcessServer   string    `json:"processServer"`
	Notes           string    `json:"notes"`
}

type CaseInvolvement struct {
	CaseNumber      string    `json:"caseNumber"`
	CourtName       string    `json:"courtName"`
	CaseType        string    `json:"caseType"`
	FilingDate      time.Time `json:"filingDate"`
	DefendantRole   string    `json:"defendantRole"`
	CaseStatus      string    `json:"caseStatus"`
	Outcome         string    `json:"outcome"`
}

type ContactMethod struct {
	MethodType  string `json:"methodType"`
	Value       string `json:"value"`
	Purpose     string `json:"purpose"`
	Verified    bool   `json:"verified"`
}

type DefendantAnalysisRules struct {
	IdentificationRules  IdentificationRuleSet      `json:"identificationRules"`
	ComparisonRules      ComparisonRuleSet          `json:"comparisonRules"`
	ViolationRules       DefendantViolationRuleSet  `json:"violationRules"`
	ServiceRules         DefendantServiceRuleSet    `json:"serviceRules"`
}

type IdentificationRuleSet struct {
	NameMatching      NameMatchingRules      `json:"nameMatching"`
	AliasResolution   AliasResolutionRules   `json:"aliasResolution"`
	CorporateMatching CorporateMatchingRules `json:"corporateMatching"`
}

type NameMatchingRules struct {
	ExactMatch        bool     `json:"exactMatch"`
	FuzzyThreshold    float64  `json:"fuzzyThreshold"`
	IgnoreElements    []string `json:"ignoreElements"`
	CriticalElements  []string `json:"criticalElements"`
}

type AliasResolutionRules struct {
	CommonAliases     map[string][]string `json:"commonAliases"`
	CorporateVariants []string            `json:"corporateVariants"`
	BusinessNames     []string            `json:"businessNames"`
}

type CorporateMatchingRules struct {
	CorporateTypes    []string `json:"corporateTypes"`
	EntityIdentifiers []string `json:"entityIdentifiers"`
	RegistrationData  []string `json:"registrationData"`
}

type ComparisonRuleSet struct {
	MultiDefendant    MultiDefendantRules    `json:"multiDefendant"`
	BusinessRelations BusinessRelationRules  `json:"businessRelations"`
	ViolationPatterns ViolationPatternRules  `json:"violationPatterns"`
}

type MultiDefendantRules struct {
	GroupingCriteria  []string `json:"groupingCriteria"`
	SimilarityFactors []string `json:"similarityFactors"`
	RelationshipTypes []string `json:"relationshipTypes"`
}

type BusinessRelationRules struct {
	ParentSubsidiary  []string `json:"parentSubsidiary"`
	Affiliations      []string `json:"affiliations"`
	BusinessPartners  []string `json:"businessPartners"`
}

type ViolationPatternRules struct {
	CommonPatterns    []string `json:"commonPatterns"`
	ViolationClusters []string `json:"violationClusters"`
	TimelineAnalysis  []string `json:"timelineAnalysis"`
}

type DefendantViolationRuleSet struct {
	CommonPatterns    []string `json:"commonPatterns"`
	ViolationClusters []string `json:"violationClusters"`
	TimelineAnalysis  []string `json:"timelineAnalysis"`
}

type DefendantServiceRuleSet struct {
	ServiceMethods    DefendantServiceMethodRules    `json:"serviceMethods"`
	ServiceValidation DefendantServiceValidationRules `json:"serviceValidation"`
	ServiceCompliance DefendantServiceComplianceRules `json:"serviceCompliance"`
}

type DefendantServiceMethodRules struct {
	PreferredMethods []string `json:"preferredMethods"`
	AlternateMethods []string `json:"alternateMethods"`
	ProhibitedMethods []string `json:"prohibitedMethods"`
}

type DefendantServiceValidationRules struct {
	RequiredElements  []string `json:"requiredElements"`
	ValidationChecks  []string `json:"validationChecks"`
	ComplianceRules   []string `json:"complianceRules"`
}

type DefendantServiceComplianceRules struct {
	StatutoryRequirements []string `json:"statutoryRequirements"`
	TimeRequirements      []string `json:"timeRequirements"`
	DocumentationRules    []string `json:"documentationRules"`
}

type ServiceRequirement struct {
	RequirementType string   `json:"requirementType"`
	Description     string   `json:"description"`
	Requirements    []string `json:"requirements"`
	Deadline        string   `json:"deadline"`
}

type MultiDefendantAnalysis struct {
	TotalDefendants      int                        `json:"totalDefendants"`
	DefendantGroups      []DefendantGroup           `json:"defendantGroups"`
	CommonViolations     []ViolationType            `json:"commonViolations"`
	UniqueViolations     map[string][]ViolationType `json:"uniqueViolations"`
	ViolationPatterns    []ViolationPattern         `json:"violationPatterns"`
	ServiceAnalysis      ServiceAnalysisSummary     `json:"serviceAnalysis"`
	JurisdictionAnalysis JurisdictionAnalysisSummary `json:"jurisdictionAnalysis"`
	RecommendedStrategy  DefendantStrategy          `json:"recommendedStrategy"`
}

type DefendantGroup struct {
	GroupType         string             `json:"groupType"`
	GroupName         string             `json:"groupName"`
	Defendants        []DefendantSummary `json:"defendants"`
	CommonElements    []string           `json:"commonElements"`
	GroupViolations   []ViolationType    `json:"groupViolations"`
	ServiceRequirements []ServiceRequirement `json:"serviceRequirements"`
}

type DefendantSummary struct {
	DefendantID   string   `json:"defendantId"`
	LegalName     string   `json:"legalName"`
	BusinessType  string   `json:"businessType"`
	ServiceMethod string   `json:"serviceMethod"`
	ViolationCount int     `json:"violationCount"`
	PrimaryRole   string   `json:"primaryRole"`
}

type ViolationType struct {
	Statute         string   `json:"statute"`
	ViolationName   string   `json:"violationName"`
	Defendants      []string `json:"defendants"`
	CommonElements  []string `json:"commonElements"`
	DamagesClaimed  float64  `json:"damagesClaimed"`
	Frequency       int      `json:"frequency"`
}

type ViolationPattern struct {
	PatternType     string   `json:"patternType"`
	Description     string   `json:"description"`
	AffectedDefendants []string `json:"affectedDefendants"`
	ViolationTypes  []string `json:"violationTypes"`
	Timeline        []PatternEvent `json:"timeline"`
	Significance    string   `json:"significance"`
}

type PatternEvent struct {
	EventDate     time.Time `json:"eventDate"`
	EventType     string    `json:"eventType"`
	Description   string    `json:"description"`
	DefendantID   string    `json:"defendantId"`
}

type ServiceAnalysisSummary struct {
	ServiceableDefendants   int                    `json:"serviceableDefendants"`
	ServiceIssues          []ServiceIssue         `json:"serviceIssues"`
	ServiceRecommendations []ServiceRecommendation `json:"serviceRecommendations"`
	ServiceTimeline        []ServiceMilestone     `json:"serviceTimeline"`
}

type ServiceIssue struct {
	DefendantID   string `json:"defendantId"`
	IssueType     string `json:"issueType"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
	Resolution    string `json:"resolution"`
}

type ServiceRecommendation struct {
	DefendantID   string   `json:"defendantId"`
	Method        string   `json:"method"`
	Priority      string   `json:"priority"`
	Requirements  []string `json:"requirements"`
	Timeline      string   `json:"timeline"`
}

type ServiceMilestone struct {
	DefendantID   string    `json:"defendantId"`
	MilestoneType string    `json:"milestoneType"`
	DueDate       time.Time `json:"dueDate"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
}

type JurisdictionAnalysisSummary struct {
	JurisdictionProper      bool                      `json:"jurisdictionProper"`
	ProblematicDefendants   []JurisdictionIssue       `json:"problematicDefendants"`
	JurisdictionStrategies  []JurisdictionStrategy    `json:"jurisdictionStrategies"`
	VenueRecommendations    []VenueRecommendation     `json:"venueRecommendations"`
}

type JurisdictionStrategy struct {
	StrategyType    string   `json:"strategyType"`
	Description     string   `json:"description"`
	ApplicableDefendants []string `json:"applicableDefendants"`
	Requirements    []string `json:"requirements"`
	Benefits        []string `json:"benefits"`
}

type VenueRecommendation struct {
	VenueType       string   `json:"venueType"`
	Jurisdiction    string   `json:"jurisdiction"`
	ApplicableDefendants []string `json:"applicableDefendants"`
	Advantages      []string `json:"advantages"`
	Considerations  []string `json:"considerations"`
}

type DefendantStrategy struct {
	StrategyType        string              `json:"strategyType"`
	PrimaryApproach     string              `json:"primaryApproach"`
	DefendantPriorities []DefendantPriority `json:"defendantPriorities"`
	ServiceOrder        []string            `json:"serviceOrder"`
	LegalTheory         []LegalTheoryElement `json:"legalTheory"`
	ExpectedOutcomes    []OutcomeProjection `json:"expectedOutcomes"`
}

type DefendantPriority struct {
	DefendantID   string   `json:"defendantId"`
	Priority      string   `json:"priority"`
	Reasoning     []string `json:"reasoning"`
	ActionItems   []string `json:"actionItems"`
}

type LegalTheoryElement struct {
	TheoryType      string   `json:"theoryType"`
	Applicability   []string `json:"applicability"`
	LegalBasis      []string `json:"legalBasis"`
	ExpectedReliefs []string `json:"expectedReliefs"`
}

type OutcomeProjection struct {
	DefendantID     string  `json:"defendantId"`
	LikelihoodScore float64 `json:"likelihoodScore"`
	ProjectedOutcome string `json:"projectedOutcome"`
	FactorsBasis    []string `json:"factorsBasis"`
}

func NewDefendantAnalyzer() (*DefendantAnalyzer, error) {
	da := &DefendantAnalyzer{}
	
	if err := da.loadCreditBureauDatabase(); err != nil {
		return nil, fmt.Errorf("failed to load credit bureau database: %v", err)
	}
	
	da.initializeAnalysisRules()
	da.DefendantProfiles = []DefendantProfile{}
	
	return da, nil
}

func (da *DefendantAnalyzer) AnalyzeDefendants(summonsDocuments []*SummonsDocument) (*MultiDefendantAnalysis, error) {
	if len(summonsDocuments) == 0 {
		return nil, fmt.Errorf("no summons documents provided")
	}

	analysis := &MultiDefendantAnalysis{
		TotalDefendants: len(summonsDocuments),
		DefendantGroups: []DefendantGroup{},
	}

	da.extractDefendantProfiles(summonsDocuments)
	da.groupDefendants(analysis)
	da.analyzeViolationPatterns(summonsDocuments, analysis)
	da.analyzeServiceRequirements(analysis)
	da.analyzeJurisdictionStrategies(analysis)
	da.generateRecommendedStrategy(analysis)

	log.Printf("Multi-defendant analysis completed: %d defendants analyzed", analysis.TotalDefendants)
	return analysis, nil
}

func (da *DefendantAnalyzer) extractDefendantProfiles(summonsDocuments []*SummonsDocument) {
	profiles := []DefendantProfile{}

	for i, summons := range summonsDocuments {
		profile := DefendantProfile{
			DefendantID:  fmt.Sprintf("defendant_%d", i+1),
			LegalName:    summons.Defendant.LegalName,
			BusinessType: summons.Defendant.BusinessType,
			CorporateType: summons.Defendant.CorporateType,
			Aliases:      summons.Defendant.Aliases,
		}

		profile.ContactInfo = DefendantContactInfo{
			RegisteredAgent: RegisteredAgentInfo{
				AgentName:    summons.Defendant.RegisteredAgent,
				AgentAddress: summons.Defendant.ServiceAddress,
			},
			BusinessAddress: summons.Defendant.BusinessAddress,
			ServiceAddress:  summons.Defendant.ServiceAddress,
		}

		profile.LegalStatus = DefendantLegalStatus{
			StateOfIncorporation: summons.Defendant.StateOfIncorporation,
			FederalTaxID:        summons.Defendant.FederalTaxID,
			CorporateStatus:     "Active",
			LastUpdated:         time.Now(),
		}

		da.enrichDefendantProfile(&profile)
		profiles = append(profiles, profile)
	}

	da.DefendantProfiles = profiles
}

func (da *DefendantAnalyzer) enrichDefendantProfile(profile *DefendantProfile) {
	if profile.BusinessType == "Credit Bureau" {
		da.enrichCreditBureauProfile(profile)
	}

	profile.BusinessOperations = DefendantOperations{
		BusinessScope: []string{"National"},
		GeographicPresence: []string{"All 50 States"},
		BusinessActivities: []BusinessActivity{
			{
				ActivityType: "Consumer Reporting",
				Description:  "Credit report generation and distribution",
				States:       []string{"All States"},
				Frequency:    "Continuous",
				Volume:       "High",
			},
		},
		IndustryClassification: "Financial Services",
	}
}

func (da *DefendantAnalyzer) enrichCreditBureauProfile(profile *DefendantProfile) {
	creditBureauData := da.getCreditBureauData(profile.LegalName)
	if creditBureauData != nil {
		if data, ok := creditBureauData.(map[string]interface{}); ok {
			if corpStructure, exists := data["corporateStructure"]; exists {
				if corp, ok := corpStructure.(map[string]interface{}); ok {
					if state, exists := corp["stateOfIncorporation"]; exists {
						profile.LegalStatus.StateOfIncorporation = state.(string)
					}
					if entityType, exists := corp["entityType"]; exists {
						profile.CorporateType = entityType.(string)
					}
				}
			}

			if contactInfo, exists := data["contactInformation"]; exists {
				da.updateContactInformation(profile, contactInfo)
			}
		}
	}
}

func (da *DefendantAnalyzer) getCreditBureauData(defendantName string) interface{} {
	defendantName = strings.ToLower(defendantName)
	
	if strings.Contains(defendantName, "equifax") {
		if bureaus, ok := da.CreditBureauDB["creditBureaus"].(map[string]interface{}); ok {
			return bureaus["equifax"]
		}
	}
	
	if strings.Contains(defendantName, "experian") {
		if bureaus, ok := da.CreditBureauDB["creditBureaus"].(map[string]interface{}); ok {
			return bureaus["experian"]
		}
	}
	
	if strings.Contains(defendantName, "trans") || strings.Contains(defendantName, "union") {
		if bureaus, ok := da.CreditBureauDB["creditBureaus"].(map[string]interface{}); ok {
			return bureaus["transunion"]
		}
	}
	
	return nil
}

func (da *DefendantAnalyzer) updateContactInformation(profile *DefendantProfile, contactInfo interface{}) {
	if contact, ok := contactInfo.(map[string]interface{}); ok {
		if regAgent, exists := contact["registeredAgent"]; exists {
			if agent, ok := regAgent.(map[string]interface{}); ok {
				profile.ContactInfo.RegisteredAgent.AgentName = agent["name"].(string)
				if address, exists := agent["address"]; exists {
					profile.ContactInfo.RegisteredAgent.AgentAddress = da.parseAddressFromInterface(address)
				}
			}
		}
	}
}

func (da *DefendantAnalyzer) parseAddressFromInterface(addressData interface{}) Address {
	address := Address{}
	if addr, ok := addressData.(map[string]interface{}); ok {
		if street, exists := addr["street"]; exists {
			address.Street = street.(string)
		}
		if city, exists := addr["city"]; exists {
			address.City = city.(string)
		}
		if state, exists := addr["state"]; exists {
			address.State = state.(string)
		}
		if zipCode, exists := addr["zipCode"]; exists {
			address.ZipCode = zipCode.(string)
		}
	}
	return address
}

func (da *DefendantAnalyzer) groupDefendants(analysis *MultiDefendantAnalysis) {
	creditBureauGroup := DefendantGroup{
		GroupType: "Credit Bureaus",
		GroupName: "Consumer Reporting Agencies",
		Defendants: []DefendantSummary{},
		CommonElements: []string{
			"Consumer reporting activities",
			"FCRA compliance obligations",
			"National business operations",
		},
	}

	creditorGroup := DefendantGroup{
		GroupType: "Creditors",
		GroupName: "Financial Institutions",
		Defendants: []DefendantSummary{},
		CommonElements: []string{
			"Furnisher responsibilities",
			"Data accuracy obligations",
			"Consumer dispute handling",
		},
	}

	otherGroup := DefendantGroup{
		GroupType: "Other",
		GroupName: "Other Defendants",
		Defendants: []DefendantSummary{},
		CommonElements: []string{},
	}

	for _, profile := range da.DefendantProfiles {
		summary := DefendantSummary{
			DefendantID:   profile.DefendantID,
			LegalName:     profile.LegalName,
			BusinessType:  profile.BusinessType,
			ServiceMethod: da.determinePreferredServiceMethod(profile),
			ViolationCount: len(profile.ViolationHistory),
			PrimaryRole:   da.determinePrimaryRole(profile.BusinessType),
		}

		switch profile.BusinessType {
		case "Credit Bureau":
			creditBureauGroup.Defendants = append(creditBureauGroup.Defendants, summary)
		case "Bank", "Financial Institution":
			creditorGroup.Defendants = append(creditorGroup.Defendants, summary)
		default:
			otherGroup.Defendants = append(otherGroup.Defendants, summary)
		}
	}

	if len(creditBureauGroup.Defendants) > 0 {
		analysis.DefendantGroups = append(analysis.DefendantGroups, creditBureauGroup)
	}
	if len(creditorGroup.Defendants) > 0 {
		analysis.DefendantGroups = append(analysis.DefendantGroups, creditorGroup)
	}
	if len(otherGroup.Defendants) > 0 {
		analysis.DefendantGroups = append(analysis.DefendantGroups, otherGroup)
	}
}

func (da *DefendantAnalyzer) determinePreferredServiceMethod(profile DefendantProfile) string {
	if profile.ContactInfo.RegisteredAgent.AgentName != "" {
		return "Registered Agent Service"
	}
	return "Certified Mail"
}

func (da *DefendantAnalyzer) determinePrimaryRole(businessType string) string {
	switch businessType {
	case "Credit Bureau":
		return "Consumer Reporting Agency"
	case "Bank":
		return "Creditor/Furnisher"
	default:
		return "Defendant"
	}
}

func (da *DefendantAnalyzer) analyzeViolationPatterns(summonsDocuments []*SummonsDocument, analysis *MultiDefendantAnalysis) {
	violationMap := make(map[string]*ViolationType)
	violationsByDefendant := make(map[string][]string)

	for i, summons := range summonsDocuments {
		defendantID := fmt.Sprintf("defendant_%d", i+1)
		
		for _, allegation := range summons.LegalAllegations {
			statute := allegation.Statute
			
			if violation, exists := violationMap[statute]; exists {
				violation.Defendants = append(violation.Defendants, defendantID)
				violation.Frequency++
			} else {
				violationMap[statute] = &ViolationType{
					Statute:       statute,
					ViolationName: allegation.Description,
					Defendants:    []string{defendantID},
					CommonElements: allegation.LegalElements,
					DamagesClaimed: 1000.0,
					Frequency:     1,
				}
			}
			
			violationsByDefendant[defendantID] = append(violationsByDefendant[defendantID], statute)
		}
	}

	commonViolations := []ViolationType{}
	uniqueViolations := make(map[string][]ViolationType)

	for _, violation := range violationMap {
		if violation.Frequency > 1 {
			commonViolations = append(commonViolations, *violation)
		} else {
			defendantID := violation.Defendants[0]
			uniqueViolations[defendantID] = append(uniqueViolations[defendantID], *violation)
		}
	}

	analysis.CommonViolations = commonViolations
	analysis.UniqueViolations = uniqueViolations

	patterns := da.identifyViolationPatterns(violationsByDefendant)
	analysis.ViolationPatterns = patterns
}

func (da *DefendantAnalyzer) identifyViolationPatterns(violationsByDefendant map[string][]string) []ViolationPattern {
	patterns := []ViolationPattern{}

	fcraPattern := ViolationPattern{
		PatternType:    "FCRA Compliance Pattern",
		Description:    "Multiple defendants with FCRA violations",
		AffectedDefendants: []string{},
		ViolationTypes: []string{},
		Significance:   "High",
	}

	for defendantID, violations := range violationsByDefendant {
		for _, violation := range violations {
			if strings.Contains(violation, "1681") {
				fcraPattern.AffectedDefendants = append(fcraPattern.AffectedDefendants, defendantID)
				if !contains(fcraPattern.ViolationTypes, violation) {
					fcraPattern.ViolationTypes = append(fcraPattern.ViolationTypes, violation)
				}
			}
		}
	}

	if len(fcraPattern.AffectedDefendants) > 1 {
		patterns = append(patterns, fcraPattern)
	}

	return patterns
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (da *DefendantAnalyzer) analyzeServiceRequirements(analysis *MultiDefendantAnalysis) {
	serviceAnalysis := ServiceAnalysisSummary{
		ServiceableDefendants: 0,
		ServiceIssues:         []ServiceIssue{},
		ServiceRecommendations: []ServiceRecommendation{},
		ServiceTimeline:       []ServiceMilestone{},
	}

	for _, profile := range da.DefendantProfiles {
		if da.isServiceable(profile) {
			serviceAnalysis.ServiceableDefendants++
			
			recommendation := ServiceRecommendation{
				DefendantID: profile.DefendantID,
				Method:      da.determinePreferredServiceMethod(profile),
				Priority:    da.determineServicePriority(profile),
				Requirements: da.getServiceRequirements(profile),
				Timeline:    "Within 20 days of filing",
			}
			serviceAnalysis.ServiceRecommendations = append(serviceAnalysis.ServiceRecommendations, recommendation)
		} else {
			issue := ServiceIssue{
				DefendantID: profile.DefendantID,
				IssueType:   "Service Information Missing",
				Description: "Insufficient service information for defendant",
				Severity:    "High",
				Resolution:  "Research registered agent and service address",
			}
			serviceAnalysis.ServiceIssues = append(serviceAnalysis.ServiceIssues, issue)
		}

		milestone := ServiceMilestone{
			DefendantID:   profile.DefendantID,
			MilestoneType: "Service Deadline",
			DueDate:       time.Now().AddDate(0, 0, 20),
			Status:        "Pending",
			Notes:         "Standard service deadline",
		}
		serviceAnalysis.ServiceTimeline = append(serviceAnalysis.ServiceTimeline, milestone)
	}

	analysis.ServiceAnalysis = serviceAnalysis
}

func (da *DefendantAnalyzer) isServiceable(profile DefendantProfile) bool {
	return profile.ContactInfo.RegisteredAgent.AgentName != "" ||
		   profile.ContactInfo.ServiceAddress.Street != ""
}

func (da *DefendantAnalyzer) determineServicePriority(profile DefendantProfile) string {
	if profile.BusinessType == "Credit Bureau" {
		return "High"
	}
	return "Medium"
}

func (da *DefendantAnalyzer) getServiceRequirements(profile DefendantProfile) []string {
	requirements := []string{
		"Proper summons and complaint",
		"Proof of service form",
	}

	if profile.ContactInfo.RegisteredAgent.AgentName != "" {
		requirements = append(requirements, "Service on registered agent")
	} else {
		requirements = append(requirements, "Certified mail service")
	}

	return requirements
}

func (da *DefendantAnalyzer) analyzeJurisdictionStrategies(analysis *MultiDefendantAnalysis) {
	jurisdictionAnalysis := JurisdictionAnalysisSummary{
		JurisdictionProper:     true,
		ProblematicDefendants:  []JurisdictionIssue{},
		JurisdictionStrategies: []JurisdictionStrategy{},
		VenueRecommendations:   []VenueRecommendation{},
	}

	fcraStrategy := JurisdictionStrategy{
		StrategyType:        "Federal Question Jurisdiction",
		Description:         "FCRA claims provide federal question jurisdiction",
		ApplicableDefendants: []string{},
		Requirements:        []string{"FCRA violations alleged", "Federal court filing"},
		Benefits:           []string{"Consistent federal law application", "Experienced federal judges"},
	}

	for _, profile := range da.DefendantProfiles {
		if profile.BusinessType == "Credit Bureau" {
			fcraStrategy.ApplicableDefendants = append(fcraStrategy.ApplicableDefendants, profile.DefendantID)
		}
	}

	if len(fcraStrategy.ApplicableDefendants) > 0 {
		jurisdictionAnalysis.JurisdictionStrategies = append(jurisdictionAnalysis.JurisdictionStrategies, fcraStrategy)
	}

	venueRec := VenueRecommendation{
		VenueType:           "Federal District Court",
		Jurisdiction:        "Any district where defendants conduct business",
		ApplicableDefendants: fcraStrategy.ApplicableDefendants,
		Advantages:         []string{"Broad venue options", "Nationwide service availability"},
		Considerations:     []string{"Forum selection considerations", "Local rule variations"},
	}
	jurisdictionAnalysis.VenueRecommendations = append(jurisdictionAnalysis.VenueRecommendations, venueRec)

	analysis.JurisdictionAnalysis = jurisdictionAnalysis
}

func (da *DefendantAnalyzer) generateRecommendedStrategy(analysis *MultiDefendantAnalysis) {
	strategy := DefendantStrategy{
		StrategyType:    "Multi-Defendant FCRA Litigation",
		PrimaryApproach: "Coordinated federal court filing with phased service",
	}

	priorities := []DefendantPriority{}
	serviceOrder := []string{}

	for _, group := range analysis.DefendantGroups {
		if group.GroupType == "Credit Bureaus" {
			for _, defendant := range group.Defendants {
				priority := DefendantPriority{
					DefendantID: defendant.DefendantID,
					Priority:    "High",
					Reasoning:   []string{"Primary FCRA violator", "High damages potential"},
					ActionItems: []string{"Serve first", "Request expedited discovery"},
				}
				priorities = append(priorities, priority)
				serviceOrder = append(serviceOrder, defendant.DefendantID)
			}
		}
	}

	for _, group := range analysis.DefendantGroups {
		if group.GroupType != "Credit Bureaus" {
			for _, defendant := range group.Defendants {
				priority := DefendantPriority{
					DefendantID: defendant.DefendantID,
					Priority:    "Medium",
					Reasoning:   []string{"Secondary violator", "Furnisher liability"},
					ActionItems: []string{"Coordinate service timing", "Joint discovery"},
				}
				priorities = append(priorities, priority)
				serviceOrder = append(serviceOrder, defendant.DefendantID)
			}
		}
	}

	strategy.DefendantPriorities = priorities
	strategy.ServiceOrder = serviceOrder

	legalTheory := []LegalTheoryElement{
		{
			TheoryType:     "FCRA Violations",
			Applicability:  []string{"All credit bureau defendants"},
			LegalBasis:     []string{"15 U.S.C. § 1681e(b)", "15 U.S.C. § 1681i"},
			ExpectedReliefs: []string{"Statutory damages", "Actual damages", "Injunctive relief"},
		},
	}
	strategy.LegalTheory = legalTheory

	analysis.RecommendedStrategy = strategy
}

func (da *DefendantAnalyzer) loadCreditBureauDatabase() error {
	data, err := ioutil.ReadFile("config/credit_bureau_database.json")
	if err != nil {
		log.Printf("Could not load credit bureau database: %v", err)
		da.CreditBureauDB = make(map[string]interface{})
		return nil
	}

	err = json.Unmarshal(data, &da.CreditBureauDB)
	if err != nil {
		return fmt.Errorf("failed to parse credit bureau database: %v", err)
	}

	log.Printf("Credit bureau database loaded successfully")
	return nil
}

func (da *DefendantAnalyzer) initializeAnalysisRules() {
	da.AnalysisRules = DefendantAnalysisRules{
		IdentificationRules: IdentificationRuleSet{
			NameMatching: NameMatchingRules{
				ExactMatch:     false,
				FuzzyThreshold: 0.8,
				IgnoreElements: []string{"LLC", "Inc.", "Corporation"},
				CriticalElements: []string{"Equifax", "Experian", "TransUnion"},
			},
		},
		ComparisonRules: ComparisonRuleSet{
			MultiDefendant: MultiDefendantRules{
				GroupingCriteria: []string{"Business Type", "Industry", "Corporate Family"},
				SimilarityFactors: []string{"Legal Name", "Business Operations", "Regulatory Status"},
			},
		},
		ViolationRules: DefendantViolationRuleSet{
			CommonPatterns: []string{"FCRA Section 1681e", "FCRA Section 1681i"},
			ViolationClusters: []string{"Accuracy Violations", "Investigation Violations"},
		},
		ServiceRules: DefendantServiceRuleSet{
			ServiceMethods: DefendantServiceMethodRules{
				PreferredMethods: []string{"Registered Agent", "Certified Mail"},
				AlternateMethods: []string{"Personal Service", "Substituted Service"},
			},
		},
	}
}