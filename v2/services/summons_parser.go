package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type SummonsDocument struct {
	DocumentPath         string               `json:"documentPath"`
	CaseInformation      CaseDetails          `json:"caseInformation"`
	Plaintiff            PartyInformation     `json:"plaintiff"`
	Defendant            DefendantDetails     `json:"defendant"`
	ServiceDetails       ServiceInformation   `json:"serviceDetails"`
	CourtInformation     CourtDetails         `json:"courtInformation"`
	ResponseRequirements ResponseDetails      `json:"responseRequirements"`
	LegalAllegations     []Allegation         `json:"legalAllegations"`
	ReliefSought         []ReliefItem         `json:"reliefSought"`
	ComplianceIssues     []ComplianceIssue    `json:"complianceIssues"`
}

type CaseDetails struct {
	CaseNumber   string `json:"caseNumber"`
	CaseTitle    string `json:"caseTitle"`
	FilingDate   string `json:"filingDate"`
	JudgeName    string `json:"judgeName"`
	CaseType     string `json:"caseType"`
	CivilAction  bool   `json:"civilAction"`
}

type PartyInformation struct {
	Name         string    `json:"name"`
	Address      Address   `json:"address"`
	AttorneyName string    `json:"attorneyName"`
	AttorneyFirm string    `json:"attorneyFirm"`
	BarNumber    string    `json:"barNumber"`
}

type DefendantDetails struct {
	LegalName            string   `json:"legalName"`
	CorporateType        string   `json:"corporateType"`
	BusinessType         string   `json:"businessType"`
	RegisteredAgent      string   `json:"registeredAgent"`
	ServiceAddress       Address  `json:"serviceAddress"`
	BusinessAddress      Address  `json:"businessAddress"`
	StateOfIncorporation string   `json:"stateOfIncorporation"`
	FederalTaxID         string   `json:"federalTaxID"`
	Aliases              []string `json:"aliases"`
}

type ServiceInformation struct {
	ServiceMethod    string    `json:"serviceMethod"`
	ServiceAddress   Address   `json:"serviceAddress"`
	RegisteredAgent  string    `json:"registeredAgent"`
	ServiceDate      time.Time `json:"serviceDate"`
	ServiceCompleted bool      `json:"serviceCompleted"`
	ServiceNotes     string    `json:"serviceNotes"`
}

type ResponseDetails struct {
	ResponseDeadline    time.Time `json:"responseDeadline"`
	ResponseDays        int       `json:"responseDays"`
	DefaultWarning      bool      `json:"defaultWarning"`
	DefaultConsequences string    `json:"defaultConsequences"`
	AnswerRequirements  []string  `json:"answerRequirements"`
}

type CourtDetails struct {
	CourtName     string `json:"courtName"`
	District      string `json:"district"`
	Division      string `json:"division"`
	Address       Address `json:"address"`
	JudgeName     string `json:"judgeName"`
	CourtType     string `json:"courtType"`
	Jurisdiction  string `json:"jurisdiction"`
}

type Allegation struct {
	ClaimNumber    int      `json:"claimNumber"`
	AllegationType string   `json:"allegationType"`
	Statute        string   `json:"statute"`
	Description    string   `json:"description"`
	SpecificFacts  []string `json:"specificFacts"`
	LegalElements  []string `json:"legalElements"`
	DefendantRole  string   `json:"defendantRole"`
}

type ReliefItem struct {
	ReliefType         string  `json:"reliefType"`
	Description        string  `json:"description"`
	MonetaryAmount     float64 `json:"monetaryAmount"`
	Statute            string  `json:"relatedStatute"`
	DefendantLiability string  `json:"defendantLiability"`
}

type ComplianceIssue struct {
	IssueType   string `json:"issueType"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Remedy      string `json:"remedy"`
}

type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zipCode"`
	Country  string `json:"country"`
}

type SummonsParser struct {
	CreditBureauPatterns map[string]CreditBureauPattern
	CourtPatterns        CourtPatternSet
	ViolationPatterns    ViolationPatternSet
}

type CreditBureauPattern struct {
	NamePatterns              []string `json:"namePatterns"`
	AddressPatterns           []string `json:"addressPatterns"`
	RegisteredAgentPatterns   []string `json:"registeredAgentPatterns"`
	BusinessTypeIndicators    []string `json:"businessTypeIndicators"`
}

type CourtPatternSet struct {
	FederalCourtPatterns FederalCourtPattern `json:"federalCourtPatterns"`
	StateCourtPatterns   StateCourtPattern   `json:"stateCourtPatterns"`
}

type FederalCourtPattern struct {
	CourtIdentification []string `json:"courtIdentification"`
	DistrictPatterns    []string `json:"districtPatterns"`
	DivisionPatterns    []string `json:"divisionPatterns"`
	CaseNumberPatterns  []string `json:"caseNumberPatterns"`
}

type StateCourtPattern struct {
	StateCourtIdentification []string `json:"stateCourtIdentification"`
	CountyPatterns           []string `json:"countyPatterns"`
	CaseNumberPatterns       []string `json:"caseNumberPatterns"`
}

type ViolationPatternSet struct {
	FCRAViolations   map[string]FCRAViolationPattern `json:"fcraViolations"`
	StateViolations  map[string]StateViolationPattern `json:"stateViolations"`
}

type FCRAViolationPattern struct {
	Patterns []string `json:"patterns"`
	Elements []string `json:"elements"`
}

type StateViolationPattern struct {
	Patterns []string `json:"patterns"`
	Elements []string `json:"elements"`
	State    string   `json:"state"`
}

func NewSummonsParser() *SummonsParser {
	return &SummonsParser{
		CreditBureauPatterns: initCreditBureauPatterns(),
		CourtPatterns:        initCourtPatterns(),
		ViolationPatterns:    initViolationPatterns(),
	}
}

func (sp *SummonsParser) ParseSummons(extractedText string, documentPath string) (*SummonsDocument, error) {
	if extractedText == "" {
		return nil, fmt.Errorf("no text provided for summons parsing")
	}

	summons := &SummonsDocument{
		DocumentPath: documentPath,
	}

	text := strings.ToUpper(extractedText)

	sp.extractCaseInformation(text, summons)
	sp.extractPlaintiffInformation(text, summons)
	sp.extractDefendantInformation(text, summons)
	sp.extractCourtInformation(text, summons)
	sp.extractServiceInformation(text, summons)
	sp.extractResponseRequirements(text, summons)
	sp.extractLegalAllegations(text, summons)
	sp.extractReliefSought(text, summons)
	sp.validateComplianceIssues(summons)

	log.Printf("Parsed summons document: %s", documentPath)
	return summons, nil
}

func (sp *SummonsParser) extractCaseInformation(text string, summons *SummonsDocument) {
	caseNumberPatterns := []string{
		`CASE\s+NO\.?\s*:?\s*([0-9]{1,2}:[0-9]{4}-CV-[0-9]{5})`,
		`CIVIL\s+ACTION\s+NO\.?\s*:?\s*([0-9-]{8,})`,
		`DOCKET\s+NO\.?\s*:?\s*([0-9-]{8,})`,
		`INDEX\s+NO\.?\s*:?\s*([0-9/]{8,})`,
	}

	for _, pattern := range caseNumberPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			summons.CaseInformation.CaseNumber = matches[1]
			break
		}
	}

	caseTitlePattern := `([A-Z\s,]+)\s+V\.?\s+([A-Z\s,\.]+)`
	if re := regexp.MustCompile(caseTitlePattern); re.MatchString(text) {
		matches := re.FindStringSubmatch(text)
		if len(matches) > 2 {
			summons.CaseInformation.CaseTitle = strings.TrimSpace(matches[0])
		}
	}

	if strings.Contains(text, "CIVIL ACTION") || strings.Contains(text, "CIVIL LAWSUIT") {
		summons.CaseInformation.CivilAction = true
		summons.CaseInformation.CaseType = "Civil Rights"
	}
}

func (sp *SummonsParser) extractPlaintiffInformation(text string, summons *SummonsDocument) {
	plaintiffPatterns := []string{
		`PLAINTIFF[S]?\s*:?\s*([A-Z\s,\.]+)`,
		`([A-Z\s,\.]+),?\s+PLAINTIFF`,
	}

	for _, pattern := range plaintiffPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			summons.Plaintiff.Name = strings.TrimSpace(matches[1])
			break
		}
	}

	attorneyPattern := `ATTORNEY\s+FOR\s+PLAINTIFF[S]?\s*:?\s*([A-Z\s,\.]+)`
	if re := regexp.MustCompile(attorneyPattern); re.MatchString(text) {
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			summons.Plaintiff.AttorneyName = strings.TrimSpace(matches[1])
		}
	}
}

func (sp *SummonsParser) extractDefendantInformation(text string, summons *SummonsDocument) {
	for bureauName, patterns := range sp.CreditBureauPatterns {
		for _, namePattern := range patterns.NamePatterns {
			re := regexp.MustCompile(namePattern)
			if re.MatchString(text) {
				summons.Defendant.LegalName = sp.extractDefendantName(text, namePattern)
				summons.Defendant.BusinessType = "Credit Bureau"
				summons.Defendant.CorporateType = sp.determineCorporateType(summons.Defendant.LegalName)
				
				sp.extractRegisteredAgent(text, patterns, summons)
				sp.extractServiceAddress(text, patterns, summons)
				
				log.Printf("Identified %s defendant: %s", bureauName, summons.Defendant.LegalName)
				return
			}
		}
	}

	genericDefendantPatterns := []string{
		`DEFENDANT[S]?\s*:?\s*([A-Z\s,\.LLC]+)`,
		`([A-Z\s,\.LLC]+),?\s+DEFENDANT`,
	}

	for _, pattern := range genericDefendantPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			summons.Defendant.LegalName = strings.TrimSpace(matches[1])
			summons.Defendant.CorporateType = sp.determineCorporateType(summons.Defendant.LegalName)
			break
		}
	}
}

func (sp *SummonsParser) extractDefendantName(text, pattern string) string {
	re := regexp.MustCompile(pattern)
	if matches := re.FindStringSubmatch(text); len(matches) > 0 {
		return strings.TrimSpace(matches[0])
	}
	return ""
}

func (sp *SummonsParser) determineCorporateType(name string) string {
	name = strings.ToUpper(name)
	if strings.Contains(name, "LLC") {
		return "LLC"
	}
	if strings.Contains(name, "INC") || strings.Contains(name, "INCORPORATED") {
		return "Corporation"
	}
	if strings.Contains(name, "LP") || strings.Contains(name, "LIMITED PARTNERSHIP") {
		return "Limited Partnership"
	}
	return "Unknown"
}

func (sp *SummonsParser) extractRegisteredAgent(text string, patterns CreditBureauPattern, summons *SummonsDocument) {
	for _, agentPattern := range patterns.RegisteredAgentPatterns {
		re := regexp.MustCompile(agentPattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 0 {
			summons.Defendant.RegisteredAgent = strings.TrimSpace(matches[0])
			summons.ServiceDetails.RegisteredAgent = summons.Defendant.RegisteredAgent
			return
		}
	}
}

func (sp *SummonsParser) extractServiceAddress(text string, patterns CreditBureauPattern, summons *SummonsDocument) {
	for _, addressPattern := range patterns.AddressPatterns {
		re := regexp.MustCompile(addressPattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 0 {
			addressText := strings.TrimSpace(matches[0])
			summons.Defendant.ServiceAddress = sp.parseAddress(addressText)
			summons.ServiceDetails.ServiceAddress = summons.Defendant.ServiceAddress
			return
		}
	}
}

func (sp *SummonsParser) parseAddress(addressText string) Address {
	parts := strings.Split(addressText, ",")
	address := Address{}
	
	if len(parts) >= 1 {
		address.Street = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 {
		address.City = strings.TrimSpace(parts[1])
	}
	if len(parts) >= 3 {
		stateZip := strings.TrimSpace(parts[2])
		stateZipParts := strings.Fields(stateZip)
		if len(stateZipParts) >= 1 {
			address.State = stateZipParts[0]
		}
		if len(stateZipParts) >= 2 {
			address.ZipCode = stateZipParts[1]
		}
	}
	
	return address
}

func (sp *SummonsParser) extractCourtInformation(text string, summons *SummonsDocument) {
	for _, courtPattern := range sp.CourtPatterns.FederalCourtPatterns.CourtIdentification {
		re := regexp.MustCompile(courtPattern)
		if re.MatchString(text) {
			summons.CourtInformation.CourtType = "Federal"
			summons.CourtInformation.CourtName = "United States District Court"
			sp.extractDistrictInformation(text, summons)
			return
		}
	}

	for _, statePattern := range sp.CourtPatterns.StateCourtPatterns.StateCourtIdentification {
		re := regexp.MustCompile(statePattern)
		if re.MatchString(text) {
			summons.CourtInformation.CourtType = "State"
			summons.CourtInformation.CourtName = "State Court"
			return
		}
	}
}

func (sp *SummonsParser) extractDistrictInformation(text string, summons *SummonsDocument) {
	for _, districtPattern := range sp.CourtPatterns.FederalCourtPatterns.DistrictPatterns {
		re := regexp.MustCompile(districtPattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 0 {
			summons.CourtInformation.District = strings.TrimSpace(matches[0])
			break
		}
	}

	for _, divisionPattern := range sp.CourtPatterns.FederalCourtPatterns.DivisionPatterns {
		re := regexp.MustCompile(divisionPattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 0 {
			summons.CourtInformation.Division = strings.TrimSpace(matches[0])
			break
		}
	}
}

func (sp *SummonsParser) extractServiceInformation(text string, summons *SummonsDocument) {
	servicePatterns := []string{
		`SERVICE\s+OF\s+PROCESS`,
		`REGISTERED\s+AGENT`,
		`CERTIFIED\s+MAIL`,
	}

	for _, pattern := range servicePatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(text) {
			if strings.Contains(pattern, "REGISTERED") {
				summons.ServiceDetails.ServiceMethod = "Registered Agent"
			} else if strings.Contains(pattern, "CERTIFIED") {
				summons.ServiceDetails.ServiceMethod = "Certified Mail"
			} else {
				summons.ServiceDetails.ServiceMethod = "Personal Service"
			}
			break
		}
	}
}

func (sp *SummonsParser) extractResponseRequirements(text string, summons *SummonsDocument) {
	responsePatterns := []string{
		`(\d+)\s+DAYS?\s+TO\s+ANSWER`,
		`ANSWER\s+WITHIN\s+(\d+)\s+DAYS?`,
		`RESPOND\s+WITHIN\s+(\d+)\s+DAYS?`,
	}

	for _, pattern := range responsePatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			days := matches[1]
			if days == "20" {
				summons.ResponseRequirements.ResponseDays = 20
			} else if days == "30" {
				summons.ResponseRequirements.ResponseDays = 30
			}
			break
		}
	}

	if strings.Contains(text, "DEFAULT JUDGMENT") {
		summons.ResponseRequirements.DefaultWarning = true
		summons.ResponseRequirements.DefaultConsequences = "Default judgment may be entered against defendant"
	}
}

func (sp *SummonsParser) extractLegalAllegations(text string, summons *SummonsDocument) {
	allegations := []Allegation{}

	for statute, pattern := range sp.ViolationPatterns.FCRAViolations {
		for _, patternStr := range pattern.Patterns {
			re := regexp.MustCompile(patternStr)
			if re.MatchString(text) {
				allegation := Allegation{
					ClaimNumber:    len(allegations) + 1,
					AllegationType: "FCRA Violation",
					Statute:        statute,
					Description:    fmt.Sprintf("Violation of %s", statute),
					LegalElements:  pattern.Elements,
					DefendantRole:  sp.determineDefendantRole(summons.Defendant.BusinessType),
				}
				allegations = append(allegations, allegation)
				break
			}
		}
	}

	summons.LegalAllegations = allegations
}

func (sp *SummonsParser) determineDefendantRole(businessType string) string {
	switch businessType {
	case "Credit Bureau":
		return "Consumer Reporting Agency"
	case "Bank":
		return "Creditor/Furnisher"
	default:
		return "Defendant"
	}
}

func (sp *SummonsParser) extractReliefSought(text string, summons *SummonsDocument) {
	reliefItems := []ReliefItem{}

	if strings.Contains(text, "ACTUAL DAMAGES") || strings.Contains(text, "COMPENSATORY DAMAGES") {
		reliefItems = append(reliefItems, ReliefItem{
			ReliefType:  "Monetary",
			Description: "Actual and compensatory damages",
		})
	}

	if strings.Contains(text, "STATUTORY DAMAGES") {
		reliefItems = append(reliefItems, ReliefItem{
			ReliefType:     "Monetary",
			Description:    "Statutory damages",
			MonetaryAmount: 1000.0,
		})
	}

	if strings.Contains(text, "INJUNCTIVE RELIEF") {
		reliefItems = append(reliefItems, ReliefItem{
			ReliefType:  "Injunctive",
			Description: "Injunctive relief to prevent future violations",
		})
	}

	if strings.Contains(text, "ATTORNEY") && strings.Contains(text, "FEES") {
		reliefItems = append(reliefItems, ReliefItem{
			ReliefType:  "Monetary",
			Description: "Attorney fees and costs",
		})
	}

	summons.ReliefSought = reliefItems
}

func (sp *SummonsParser) validateComplianceIssues(summons *SummonsDocument) {
	issues := []ComplianceIssue{}

	if summons.Defendant.LegalName == "" {
		issues = append(issues, ComplianceIssue{
			IssueType:   "Defendant Identification",
			Description: "Unable to identify defendant name",
			Severity:    "High",
			Remedy:      "Manual review required",
		})
	}

	if summons.ServiceDetails.ServiceMethod == "" {
		issues = append(issues, ComplianceIssue{
			IssueType:   "Service Method",
			Description: "Service method not specified",
			Severity:    "Medium",
			Remedy:      "Verify proper service method",
		})
	}

	if summons.ResponseRequirements.ResponseDays == 0 {
		issues = append(issues, ComplianceIssue{
			IssueType:   "Response Deadline",
			Description: "Response deadline not found",
			Severity:    "Medium",
			Remedy:      "Check for answer requirements",
		})
	}

	summons.ComplianceIssues = issues
}

func initCreditBureauPatterns() map[string]CreditBureauPattern {
	return map[string]CreditBureauPattern{
		"equifax": {
			NamePatterns: []string{
				`EQUIFAX\s+INFORMATION\s+SERVICES,?\s+LLC`,
				`EQUIFAX\s+INFORMATION\s+SERVICES\s+LLC`,
				`EQUIFAX\s+INC\.?`,
				`EQUIFAX,?\s+LLC`,
			},
			AddressPatterns: []string{
				`1550\s+PEACHTREE.*ATLANTA.*GA`,
				`P\.?O\.?\s+BOX.*ATLANTA.*GEORGIA`,
				`EQUIFAX.*ATLANTA.*30309`,
			},
			RegisteredAgentPatterns: []string{
				`CORPORATION\s+SERVICE\s+COMPANY`,
				`CSC.*REGISTERED\s+AGENT`,
				`CT\s+CORPORATION.*ATLANTA`,
			},
		},
		"experian": {
			NamePatterns: []string{
				`EXPERIAN\s+INFORMATION\s+SOLUTIONS,?\s+INC\.?`,
				`EXPERIAN\s+INFORMATION\s+SOLUTIONS\s+INC\.?`,
				`EXPERIAN\s+INC\.?`,
				`EXPERIAN.*LLC`,
			},
			AddressPatterns: []string{
				`475\s+ANTON.*COSTA\s+MESA.*CA`,
				`P\.?O\.?\s+BOX.*ALLEN.*TX`,
				`EXPERIAN.*COSTA\s+MESA.*92626`,
			},
			RegisteredAgentPatterns: []string{
				`CORPORATION\s+SERVICE\s+COMPANY`,
				`CSC.*REGISTERED\s+AGENT`,
				`CT\s+CORPORATION.*CALIFORNIA`,
			},
		},
		"transunion": {
			NamePatterns: []string{
				`TRANS\s+UNION\s+LLC`,
				`TRANSUNION\s+LLC`,
				`TRANS\s+UNION\s+INFORMATION\s+SOLUTIONS`,
				`TRANSUNION.*LLC`,
			},
			AddressPatterns: []string{
				`555\s+WEST\s+ADAMS.*CHICAGO.*IL`,
				`P\.?O\.?\s+BOX.*CHESTER.*PA`,
				`TRANSUNION.*CHICAGO.*60661`,
			},
			RegisteredAgentPatterns: []string{
				`CORPORATION\s+SERVICE\s+COMPANY`,
				`ILLINOIS\s+CORPORATION\s+SERVICE`,
				`CSC.*REGISTERED\s+AGENT`,
			},
		},
	}
}

func initCourtPatterns() CourtPatternSet {
	return CourtPatternSet{
		FederalCourtPatterns: FederalCourtPattern{
			CourtIdentification: []string{
				`UNITED\s+STATES\s+DISTRICT\s+COURT`,
				`U\.S\.\s+DISTRICT\s+COURT`,
				`DISTRICT\s+COURT.*UNITED\s+STATES`,
			},
			DistrictPatterns: []string{
				`EASTERN\s+DISTRICT\s+OF\s+NEW\s+YORK`,
				`SOUTHERN\s+DISTRICT\s+OF\s+NEW\s+YORK`,
				`CENTRAL\s+DISTRICT\s+OF\s+CALIFORNIA`,
				`NORTHERN\s+DISTRICT\s+OF\s+ILLINOIS`,
			},
			DivisionPatterns: []string{
				`BROOKLYN\s+DIVISION`,
				`MANHATTAN\s+DIVISION`,
				`WHITE\s+PLAINS\s+DIVISION`,
			},
		},
		StateCourtPatterns: StateCourtPattern{
			StateCourtIdentification: []string{
				`SUPREME\s+COURT`,
				`SUPERIOR\s+COURT`,
				`CIRCUIT\s+COURT`,
			},
		},
	}
}

func initViolationPatterns() ViolationPatternSet {
	return ViolationPatternSet{
		FCRAViolations: map[string]FCRAViolationPattern{
			"15 U.S.C. § 1681e(b)": {
				Patterns: []string{
					`15\s+U\.S\.C\.?\s+§\s+1681E\(B\)`,
					`REASONABLE\s+PROCEDURES.*MAXIMUM\s+POSSIBLE\s+ACCURACY`,
					`FAILED\s+TO\s+FOLLOW\s+REASONABLE\s+PROCEDURES`,
				},
				Elements: []string{
					"duty to maintain reasonable procedures",
					"failure to assure maximum possible accuracy",
					"reporting of inaccurate information",
				},
			},
			"15 U.S.C. § 1681i": {
				Patterns: []string{
					`15\s+U\.S\.C\.?\s+§\s+1681I`,
					`REINVESTIGATION.*CONSUMER\s+DISPUTE`,
					`FAILED\s+TO.*REASONABLE\s+REINVESTIGATION`,
				},
				Elements: []string{
					"received consumer dispute",
					"failed to conduct reasonable reinvestigation",
					"continued reporting disputed information",
				},
			},
		},
	}
}