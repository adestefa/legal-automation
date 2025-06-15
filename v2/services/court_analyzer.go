package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type CourtAnalyzer struct {
	FederalCourts     map[string]FederalCourtInfo     `json:"federalCourts"`
	StateCourts       map[string]StateCourtInfo       `json:"stateCourts"`
	JurisdictionRules JurisdictionRuleSet             `json:"jurisdictionRules"`
	VenueRules        VenueRuleSet                    `json:"venueRules"`
}

type FederalCourtInfo struct {
	CourtName       string               `json:"courtName"`
	District        string               `json:"district"`
	Divisions       []CourtDivision      `json:"divisions"`
	Judges          []JudgeInfo          `json:"judges"`
	CaseNumberFormat string              `json:"caseNumberFormat"`
	Jurisdiction    JurisdictionDetails  `json:"jurisdiction"`
	Address         Address              `json:"address"`
}

type StateCourtInfo struct {
	CourtName        string              `json:"courtName"`
	County           string              `json:"county"`
	State            string              `json:"state"`
	Judges           []JudgeInfo         `json:"judges"`
	CaseNumberFormat string              `json:"caseNumberFormat"`
	Jurisdiction     JurisdictionDetails `json:"jurisdiction"`
	Address          Address             `json:"address"`
}

type CourtDivision struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Counties    []string `json:"counties"`
	Address     Address `json:"address"`
}

type JudgeInfo struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	Appointed string `json:"appointed"`
	Division  string `json:"division"`
}

type JurisdictionDetails struct {
	Type                  string   `json:"type"`
	GeographicScope       []string `json:"geographicScope"`
	SubjectMatterScope    []string `json:"subjectMatterScope"`
	PersonalJurisdiction  []string `json:"personalJurisdiction"`
	MinimumContacts       []string `json:"minimumContacts"`
}

type JurisdictionRuleSet struct {
	PersonalJurisdiction PersonalJurisdictionRules `json:"personalJurisdiction"`
	SubjectMatter        SubjectMatterRules        `json:"subjectMatter"`
	Venue                VenueAnalysisRules        `json:"venue"`
}

type PersonalJurisdictionRules struct {
	CorporateDefendants CorporateJurisdictionRules `json:"corporateDefendants"`
	MinimumContacts     MinimumContactsRules       `json:"minimumContacts"`
	LongArmStatutes     []LongArmStatute           `json:"longArmStatutes"`
}

type CorporateJurisdictionRules struct {
	IncorporationState    bool     `json:"incorporationState"`
	PrincipalPlaceOfBusiness bool  `json:"principalPlaceOfBusiness"`
	RegisteredToDoBusiness bool    `json:"registeredToDoBusiness"`
	SystematicContacts    []string `json:"systematicContacts"`
	ContinuousOperations  []string `json:"continuousOperations"`
}

type MinimumContactsRules struct {
	QuantityFactors   []string `json:"quantityFactors"`
	QualityFactors    []string `json:"qualityFactors"`
	RelatednessTest   []string `json:"relatednessTest"`
	FairnessFactors   []string `json:"fairnessFactors"`
}

type LongArmStatute struct {
	State       string   `json:"state"`
	Statute     string   `json:"statute"`
	Provisions  []string `json:"provisions"`
	Limitations []string `json:"limitations"`
}

type SubjectMatterRules struct {
	FederalQuestion   FederalQuestionRules   `json:"federalQuestion"`
	DiversityJurisdiction DiversityRules     `json:"diversityJurisdiction"`
	SupplementalJurisdiction SupplementalRules `json:"supplementalJurisdiction"`
}

type FederalQuestionRules struct {
	FederalStatutes    []string `json:"federalStatutes"`
	ConstitutionalClaims []string `json:"constitutionalClaims"`
	FederalRegulations []string `json:"federalRegulations"`
}

type DiversityRules struct {
	AmountInControversy float64  `json:"amountInControversy"`
	CompleteDiversity   bool     `json:"completeDiversity"`
	CitizenshipRules    []string `json:"citizenshipRules"`
}

type SupplementalRules struct {
	CommonNucleus     bool     `json:"commonNucleus"`
	SameTransaction   bool     `json:"sameTransaction"`
	RelatedClaims     []string `json:"relatedClaims"`
}

type VenueRuleSet struct {
	FederalVenue  FederalVenueRules  `json:"federalVenue"`
	StateVenue    StateVenueRules    `json:"stateVenue"`
	VenueFactors  VenueFactorRules   `json:"venueFactors"`
}

type FederalVenueRules struct {
	ResidenceRule      VenueRule `json:"residenceRule"`
	EventsRule         VenueRule `json:"eventsRule"`
	PropertyRule       VenueRule `json:"propertyRule"`
	FallbackRule       VenueRule `json:"fallbackRule"`
	TransferRules      []TransferRule `json:"transferRules"`
}

type StateVenueRules struct {
	ResidenceRule      VenueRule `json:"residenceRule"`
	BusinessRule       VenueRule `json:"businessRule"`
	CauseOfActionRule  VenueRule `json:"causeOfActionRule"`
	ConvenienceRule    VenueRule `json:"convenienceRule"`
}

type VenueRule struct {
	Description   string   `json:"description"`
	Statute       string   `json:"statute"`
	Requirements  []string `json:"requirements"`
	Limitations   []string `json:"limitations"`
}

type TransferRule struct {
	Type          string   `json:"type"`
	Statute       string   `json:"statute"`
	Requirements  []string `json:"requirements"`
	Factors       []string `json:"factors"`
}

type VenueFactorRules struct {
	ConvenienceFactors  []string `json:"convenienceFactors"`
	JusticeFactors      []string `json:"justiceFactors"`
	EconomicFactors     []string `json:"economicFactors"`
}

type VenueAnalysisRules struct {
	ProperVenue       []string `json:"properVenue"`
	ImproperVenue     []string `json:"improperVenue"`
	VenueWaiver       []string `json:"venueWaiver"`
	TransferCriteria  []string `json:"transferCriteria"`
}

type CourtAnalysisResult struct {
	CourtType            string                    `json:"courtType"`
	CourtName            string                    `json:"courtName"`
	District             string                    `json:"district"`
	Division             string                    `json:"division"`
	JurisdictionAnalysis JurisdictionAnalysisResult `json:"jurisdictionAnalysis"`
	VenueAnalysis        VenueAnalysisResult       `json:"venueAnalysis"`
	ComplianceIssues     []CourtComplianceIssue    `json:"complianceIssues"`
	Recommendations      []CourtRecommendation     `json:"recommendations"`
}

type JurisdictionAnalysisResult struct {
	PersonalJurisdiction  PersonalJurisdictionResult `json:"personalJurisdiction"`
	SubjectMatterProper   bool                       `json:"subjectMatterProper"`
	JurisdictionBasis     []string                   `json:"jurisdictionBasis"`
	JurisdictionIssues    []JurisdictionIssue        `json:"jurisdictionIssues"`
}

type PersonalJurisdictionResult struct {
	HasJurisdiction     bool     `json:"hasJurisdiction"`
	JurisdictionBasis   []string `json:"jurisdictionBasis"`
	MinimumContacts     bool     `json:"minimumContacts"`
	ConstitutionalTest  bool     `json:"constitutionalTest"`
	LongArmCoverage     bool     `json:"longArmCoverage"`
}

type VenueAnalysisResult struct {
	VenueProper         bool     `json:"venueProper"`
	VenueBasis          []string `json:"venueBasis"`
	AlternativeVenues   []string `json:"alternativeVenues"`
	VenueIssues         []VenueIssue `json:"venueIssues"`
}

type JurisdictionIssue struct {
	IssueType     string `json:"issueType"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
	Resolution    string `json:"resolution"`
}

type VenueIssue struct {
	IssueType     string `json:"issueType"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
	Resolution    string `json:"resolution"`
}

type CourtComplianceIssue struct {
	IssueType     string `json:"issueType"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
	Resolution    string `json:"resolution"`
	Rule          string `json:"rule"`
}

type CourtRecommendation struct {
	Type          string `json:"type"`
	Description   string `json:"description"`
	Priority      string `json:"priority"`
	ActionItems   []string `json:"actionItems"`
}

func NewCourtAnalyzer() (*CourtAnalyzer, error) {
	ca := &CourtAnalyzer{}
	
	if err := ca.loadCourtData(); err != nil {
		return nil, fmt.Errorf("failed to load court data: %v", err)
	}
	
	ca.initializeJurisdictionRules()
	ca.initializeVenueRules()
	
	return ca, nil
}

func (ca *CourtAnalyzer) AnalyzeCourt(extractedText string, defendant *DefendantDetails) (*CourtAnalysisResult, error) {
	if extractedText == "" {
		return nil, fmt.Errorf("no text provided for court analysis")
	}

	result := &CourtAnalysisResult{}
	text := strings.ToUpper(extractedText)

	ca.identifyCourtType(text, result)
	ca.analyzeJurisdiction(text, defendant, result)
	ca.analyzeVenue(text, defendant, result)
	ca.validateCompliance(result)
	ca.generateRecommendations(result)

	log.Printf("Court analysis completed for %s", result.CourtName)
	return result, nil
}

func (ca *CourtAnalyzer) identifyCourtType(text string, result *CourtAnalysisResult) {
	federalPatterns := []string{
		`UNITED\s+STATES\s+DISTRICT\s+COURT`,
		`U\.S\.\s+DISTRICT\s+COURT`,
		`DISTRICT\s+COURT.*UNITED\s+STATES`,
	}

	for _, pattern := range federalPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(text) {
			result.CourtType = "Federal"
			ca.identifyFederalCourt(text, result)
			return
		}
	}

	statePatterns := []string{
		`SUPREME\s+COURT`,
		`SUPERIOR\s+COURT`,
		`CIRCUIT\s+COURT`,
		`COUNTY\s+COURT`,
	}

	for _, pattern := range statePatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(text) {
			result.CourtType = "State"
			ca.identifyStateCourt(text, result)
			return
		}
	}

	result.CourtType = "Unknown"
	result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
		IssueType:   "Court Identification",
		Description: "Unable to identify court type",
		Severity:    "High",
		Resolution:  "Manual review required",
	})
}

func (ca *CourtAnalyzer) identifyFederalCourt(text string, result *CourtAnalysisResult) {
	districtPatterns := map[string]string{
		"EASTERN DISTRICT OF NEW YORK":     "E.D.N.Y.",
		"SOUTHERN DISTRICT OF NEW YORK":    "S.D.N.Y.",
		"NORTHERN DISTRICT OF NEW YORK":    "N.D.N.Y.",
		"WESTERN DISTRICT OF NEW YORK":     "W.D.N.Y.",
		"CENTRAL DISTRICT OF CALIFORNIA":   "C.D. Cal.",
		"NORTHERN DISTRICT OF CALIFORNIA":  "N.D. Cal.",
		"SOUTHERN DISTRICT OF CALIFORNIA":  "S.D. Cal.",
		"EASTERN DISTRICT OF CALIFORNIA":   "E.D. Cal.",
		"NORTHERN DISTRICT OF ILLINOIS":    "N.D. Ill.",
		"CENTRAL DISTRICT OF ILLINOIS":     "C.D. Ill.",
		"SOUTHERN DISTRICT OF ILLINOIS":    "S.D. Ill.",
	}

	for pattern, abbreviation := range districtPatterns {
		if strings.Contains(text, pattern) {
			result.CourtName = fmt.Sprintf("United States District Court, %s", abbreviation)
			result.District = pattern
			ca.identifyDivision(text, pattern, result)
			break
		}
	}

	if result.District == "" {
		result.CourtName = "United States District Court"
		result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
			IssueType:   "District Identification",
			Description: "Unable to identify specific federal district",
			Severity:    "Medium",
			Resolution:  "Review court heading for district information",
		})
	}
}

func (ca *CourtAnalyzer) identifyDivision(text, district string, result *CourtAnalysisResult) {
	divisionPatterns := map[string][]string{
		"EASTERN DISTRICT OF NEW YORK": {
			"BROOKLYN DIVISION",
			"CENTRAL ISLIP DIVISION",
		},
		"SOUTHERN DISTRICT OF NEW YORK": {
			"MANHATTAN DIVISION",
			"WHITE PLAINS DIVISION",
		},
	}

	if divisions, exists := divisionPatterns[district]; exists {
		for _, division := range divisions {
			if strings.Contains(text, division) {
				result.Division = division
				break
			}
		}
	}
}

func (ca *CourtAnalyzer) identifyStateCourt(text string, result *CourtAnalysisResult) {
	statePatterns := map[string]string{
		"NEW YORK SUPREME COURT":     "New York Supreme Court",
		"CALIFORNIA SUPERIOR COURT":  "California Superior Court",
		"ILLINOIS CIRCUIT COURT":     "Illinois Circuit Court",
		"FLORIDA CIRCUIT COURT":      "Florida Circuit Court",
	}

	for pattern, courtName := range statePatterns {
		if strings.Contains(text, pattern) {
			result.CourtName = courtName
			ca.identifyCounty(text, result)
			break
		}
	}

	if result.CourtName == "" {
		result.CourtName = "State Court"
		result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
			IssueType:   "State Court Identification",
			Description: "Unable to identify specific state court",
			Severity:    "Medium",
			Resolution:  "Review court heading for state court information",
		})
	}
}

func (ca *CourtAnalyzer) identifyCounty(text string, result *CourtAnalysisResult) {
	countyPatterns := []string{
		`(\w+)\s+COUNTY`,
		`COUNTY\s+OF\s+(\w+)`,
	}

	for _, pattern := range countyPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			result.District = fmt.Sprintf("%s County", matches[1])
			break
		}
	}
}

func (ca *CourtAnalyzer) analyzeJurisdiction(text string, defendant *DefendantDetails, result *CourtAnalysisResult) {
	result.JurisdictionAnalysis = JurisdictionAnalysisResult{}

	if result.CourtType == "Federal" {
		ca.analyzeFederalJurisdiction(text, defendant, result)
	} else if result.CourtType == "State" {
		ca.analyzeStateJurisdiction(text, defendant, result)
	}

	ca.analyzePersonalJurisdiction(defendant, result)
}

func (ca *CourtAnalyzer) analyzeFederalJurisdiction(text string, defendant *DefendantDetails, result *CourtAnalysisResult) {
	fcraPatterns := []string{
		`15\s+U\.S\.C\.?\s+§\s+1681`,
		`FAIR\s+CREDIT\s+REPORTING\s+ACT`,
		`FCRA`,
	}

	hasFederalQuestion := false
	for _, pattern := range fcraPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(text) {
			hasFederalQuestion = true
			result.JurisdictionAnalysis.JurisdictionBasis = append(
				result.JurisdictionAnalysis.JurisdictionBasis,
				"Federal Question Jurisdiction - FCRA Claims",
			)
			break
		}
	}

	if !hasFederalQuestion {
		ca.analyzeDiversityJurisdiction(defendant, result)
	}

	result.JurisdictionAnalysis.SubjectMatterProper = hasFederalQuestion || ca.hasDiversityJurisdiction(defendant)
}

func (ca *CourtAnalyzer) analyzeDiversityJurisdiction(defendant *DefendantDetails, result *CourtAnalysisResult) {
	if defendant.StateOfIncorporation != "" {
		result.JurisdictionAnalysis.JurisdictionBasis = append(
			result.JurisdictionAnalysis.JurisdictionBasis,
			fmt.Sprintf("Potential Diversity Jurisdiction - Defendant incorporated in %s", defendant.StateOfIncorporation),
		)
	}

	result.JurisdictionAnalysis.JurisdictionIssues = append(
		result.JurisdictionAnalysis.JurisdictionIssues,
		JurisdictionIssue{
			IssueType:   "Diversity Analysis",
			Description: "Requires verification of plaintiff citizenship and amount in controversy",
			Severity:    "Medium",
			Resolution:  "Confirm diversity requirements are met",
		},
	)
}

func (ca *CourtAnalyzer) hasDiversityJurisdiction(defendant *DefendantDetails) bool {
	return defendant.StateOfIncorporation != ""
}

func (ca *CourtAnalyzer) analyzeStateJurisdiction(text string, defendant *DefendantDetails, result *CourtAnalysisResult) {
	result.JurisdictionAnalysis.SubjectMatterProper = true
	result.JurisdictionAnalysis.JurisdictionBasis = append(
		result.JurisdictionAnalysis.JurisdictionBasis,
		"State Court General Jurisdiction",
	)

	if defendant.BusinessType == "Credit Bureau" {
		result.JurisdictionAnalysis.JurisdictionIssues = append(
			result.JurisdictionAnalysis.JurisdictionIssues,
			JurisdictionIssue{
				IssueType:   "Federal Preemption",
				Description: "FCRA claims may be better suited for federal court",
				Severity:    "Low",
				Resolution:  "Consider federal court filing advantages",
			},
		)
	}
}

func (ca *CourtAnalyzer) analyzePersonalJurisdiction(defendant *DefendantDetails, result *CourtAnalysisResult) {
	pjResult := PersonalJurisdictionResult{
		HasJurisdiction:   false,
		JurisdictionBasis: []string{},
		MinimumContacts:   false,
		ConstitutionalTest: false,
		LongArmCoverage:   false,
	}

	if defendant.BusinessType == "Credit Bureau" {
		pjResult.HasJurisdiction = true
		pjResult.MinimumContacts = true
		pjResult.ConstitutionalTest = true
		pjResult.LongArmCoverage = true
		pjResult.JurisdictionBasis = []string{
			"Systematic and continuous business contacts",
			"Credit reports furnished to forum state residents",
			"Registration or qualification to do business",
			"Substantial business operations nationwide",
		}
	}

	if defendant.StateOfIncorporation != "" {
		pjResult.JurisdictionBasis = append(pjResult.JurisdictionBasis,
			fmt.Sprintf("State of incorporation: %s", defendant.StateOfIncorporation))
	}

	result.JurisdictionAnalysis.PersonalJurisdiction = pjResult
}

func (ca *CourtAnalyzer) analyzeVenue(text string, defendant *DefendantDetails, result *CourtAnalysisResult) {
	venueResult := VenueAnalysisResult{
		VenueProper:       false,
		VenueBasis:        []string{},
		AlternativeVenues: []string{},
		VenueIssues:       []VenueIssue{},
	}

	if result.CourtType == "Federal" {
		ca.analyzeFederalVenue(defendant, &venueResult)
	} else if result.CourtType == "State" {
		ca.analyzeStateVenue(defendant, &venueResult)
	}

	result.VenueAnalysis = venueResult
}

func (ca *CourtAnalyzer) analyzeFederalVenue(defendant *DefendantDetails, venueResult *VenueAnalysisResult) {
	venueResult.VenueBasis = append(venueResult.VenueBasis,
		"28 U.S.C. § 1391(b) - General venue statute")

	if defendant.BusinessType == "Credit Bureau" {
		venueResult.VenueProper = true
		venueResult.VenueBasis = append(venueResult.VenueBasis,
			"Defendant conducts substantial business nationwide",
			"Consumer reporting activities in district",
		)

		venueResult.AlternativeVenues = []string{
			"District of defendant's incorporation",
			"District of defendant's principal place of business",
			"Any district where defendant conducts substantial business",
		}
	}

	if !venueResult.VenueProper {
		venueResult.VenueIssues = append(venueResult.VenueIssues, VenueIssue{
			IssueType:   "Venue Analysis",
			Description: "Venue determination requires further analysis",
			Severity:    "Medium",
			Resolution:  "Verify defendant's business activities in district",
		})
	}
}

func (ca *CourtAnalyzer) analyzeStateVenue(defendant *DefendantDetails, venueResult *VenueAnalysisResult) {
	venueResult.VenueProper = true
	venueResult.VenueBasis = append(venueResult.VenueBasis,
		"State venue statute",
		"Defendant conducts business in county",
	)

	if defendant.RegisteredAgent != "" {
		venueResult.VenueBasis = append(venueResult.VenueBasis,
			"Registered agent for service in state")
	}
}

func (ca *CourtAnalyzer) validateCompliance(result *CourtAnalysisResult) {
	if result.CourtName == "" || result.CourtName == "Unknown" {
		result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
			IssueType:   "Court Identification",
			Description: "Court not properly identified",
			Severity:    "High",
			Resolution:  "Review summons heading for court information",
			Rule:        "Fed. R. Civ. P. 3",
		})
	}

	if !result.JurisdictionAnalysis.SubjectMatterProper {
		result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
			IssueType:   "Subject Matter Jurisdiction",
			Description: "Subject matter jurisdiction not established",
			Severity:    "High",
			Resolution:  "Establish federal question or diversity jurisdiction",
			Rule:        "28 U.S.C. § 1331, § 1332",
		})
	}

	if !result.JurisdictionAnalysis.PersonalJurisdiction.HasJurisdiction {
		result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
			IssueType:   "Personal Jurisdiction",
			Description: "Personal jurisdiction over defendant questionable",
			Severity:    "High",
			Resolution:  "Establish minimum contacts and constitutional requirements",
			Rule:        "Due Process Clause",
		})
	}

	if !result.VenueAnalysis.VenueProper {
		result.ComplianceIssues = append(result.ComplianceIssues, CourtComplianceIssue{
			IssueType:   "Venue",
			Description: "Venue may not be proper",
			Severity:    "Medium",
			Resolution:  "Consider venue transfer or establish proper venue",
			Rule:        "28 U.S.C. § 1391",
		})
	}
}

func (ca *CourtAnalyzer) generateRecommendations(result *CourtAnalysisResult) {
	recommendations := []CourtRecommendation{}

	if result.CourtType == "Federal" && len(result.JurisdictionAnalysis.JurisdictionBasis) > 0 {
		recommendations = append(recommendations, CourtRecommendation{
			Type:        "Jurisdiction Strategy",
			Description: "Federal court filing recommended for FCRA claims",
			Priority:    "High",
			ActionItems: []string{
				"Emphasize federal question jurisdiction in complaint",
				"Include comprehensive FCRA violation allegations",
				"Consider supplemental state law claims",
			},
		})
	}

	if result.JurisdictionAnalysis.PersonalJurisdiction.HasJurisdiction {
		recommendations = append(recommendations, CourtRecommendation{
			Type:        "Service Strategy",
			Description: "Personal jurisdiction established - proceed with service",
			Priority:    "High",
			ActionItems: []string{
				"Serve registered agent in accordance with state law",
				"File proof of service promptly",
				"Monitor response deadline",
			},
		})
	}

	if len(result.ComplianceIssues) > 0 {
		recommendations = append(recommendations, CourtRecommendation{
			Type:        "Compliance Review",
			Description: "Address compliance issues before proceeding",
			Priority:    "High",
			ActionItems: []string{
				"Review and resolve all compliance issues",
				"Consider amended pleading if necessary",
				"Consult local rules and procedures",
			},
		})
	}

	result.Recommendations = recommendations
}

func (ca *CourtAnalyzer) loadCourtData() error {
	ca.FederalCourts = make(map[string]FederalCourtInfo)
	ca.StateCourts = make(map[string]StateCourtInfo)

	return nil
}

func (ca *CourtAnalyzer) initializeJurisdictionRules() {
	ca.JurisdictionRules = JurisdictionRuleSet{
		PersonalJurisdiction: PersonalJurisdictionRules{
			CorporateDefendants: CorporateJurisdictionRules{
				IncorporationState:       true,
				PrincipalPlaceOfBusiness: true,
				RegisteredToDoBusiness:   true,
				SystematicContacts: []string{
					"Regular business operations",
					"Continuous commercial activities",
					"Office or agent presence",
				},
			},
		},
		SubjectMatter: SubjectMatterRules{
			FederalQuestion: FederalQuestionRules{
				FederalStatutes: []string{
					"15 U.S.C. § 1681 (FCRA)",
					"42 U.S.C. § 1983 (Civil Rights)",
				},
			},
			DiversityJurisdiction: DiversityRules{
				AmountInControversy: 75000.00,
				CompleteDiversity:   true,
			},
		},
	}
}

func (ca *CourtAnalyzer) initializeVenueRules() {
	ca.VenueRules = VenueRuleSet{
		FederalVenue: FederalVenueRules{
			ResidenceRule: VenueRule{
				Description:  "District where defendant resides",
				Statute:      "28 U.S.C. § 1391(b)(1)",
				Requirements: []string{"Corporate residence in district"},
			},
			EventsRule: VenueRule{
				Description:  "District where substantial part of events occurred",
				Statute:      "28 U.S.C. § 1391(b)(2)",
				Requirements: []string{"Substantial business activity in district"},
			},
		},
	}
}