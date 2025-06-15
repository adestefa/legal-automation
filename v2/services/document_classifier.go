package services

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
)

type DocumentType int

const (
	DocumentTypeUnknown DocumentType = iota
	DocumentTypeAttorneyNotes
	DocumentTypeAdverseActionLetter
	DocumentTypeCivilCoverSheet
	DocumentTypeSummons
	DocumentTypeComplaint
	DocumentTypeDenialLetter
	DocumentTypeCorrespondence
	DocumentTypeCreditReport
	DocumentTypeFinancialStatement
	DocumentTypeDisputeLetter
)

var documentTypeNames = map[DocumentType]string{
	DocumentTypeUnknown:             "Unknown",
	DocumentTypeAttorneyNotes:       "Attorney Notes",
	DocumentTypeAdverseActionLetter: "Adverse Action Letter",
	DocumentTypeCivilCoverSheet:     "Civil Cover Sheet",
	DocumentTypeSummons:             "Summons",
	DocumentTypeComplaint:           "Complaint",
	DocumentTypeDenialLetter:        "Denial Letter",
	DocumentTypeCorrespondence:      "Correspondence",
	DocumentTypeCreditReport:        "Credit Report",
	DocumentTypeFinancialStatement:  "Financial Statement",
	DocumentTypeDisputeLetter:       "Dispute Letter",
}

type DocumentClassification struct {
	DocumentPath      string                 `json:"documentPath"`
	PrimaryType       DocumentType           `json:"primaryType"`
	SecondaryTypes    []DocumentType         `json:"secondaryTypes"`
	Confidence        float64                `json:"confidence"`
	ContentIndicators []ContentIndicator     `json:"contentIndicators"`
	ExtractedFields   map[string]interface{} `json:"extractedFields"`
	ValidationScore   float64                `json:"validationScore"`
}

type ContentIndicator struct {
	Pattern    string  `json:"pattern"`
	MatchType  string  `json:"matchType"`
	Confidence float64 `json:"confidence"`
	Location   string  `json:"location"`
}

type DocumentTypePattern struct {
	HeaderPatterns       []string          `json:"headerPatterns"`
	ContentPatterns      []string          `json:"contentPatterns"`
	StatutoryReferences  []string          `json:"statutoryReferences"`
	RequiredElements     []string          `json:"requiredElements"`
	DefendantPatterns    []string          `json:"defendantPatterns,omitempty"`
	StructureIndicators  map[string]string `json:"structureIndicators,omitempty"`
}

type DocumentTypePatterns struct {
	AdverseActionLetter DocumentTypePattern `json:"adverseActionLetter"`
	Summons            DocumentTypePattern `json:"summons"`
	AttorneyNotes      DocumentTypePattern `json:"attorneyNotes"`
	CivilCoverSheet    DocumentTypePattern `json:"civilCoverSheet"`
	Complaint          DocumentTypePattern `json:"complaint"`
	DenialLetter       DocumentTypePattern `json:"denialLetter"`
}

type DocumentClassifier struct {
	patterns       map[string]interface{}
	documentTypes  DocumentTypePatterns
	contentService *ContentAnalyzer
}

func NewDocumentClassifier(contentAnalyzer *ContentAnalyzer) (*DocumentClassifier, error) {
	classifier := &DocumentClassifier{
		contentService: contentAnalyzer,
		patterns:       make(map[string]interface{}),
	}

	if err := classifier.loadPatterns(); err != nil {
		return nil, fmt.Errorf("failed to load patterns: %w", err)
	}

	return classifier, nil
}

func (dc *DocumentClassifier) loadPatterns() error {
	data, err := os.ReadFile("./config/legal_patterns.json")
	if err != nil {
		return fmt.Errorf("failed to read legal patterns: %w", err)
	}

	if err := json.Unmarshal(data, &dc.patterns); err != nil {
		return fmt.Errorf("failed to parse legal patterns: %w", err)
	}

	dc.documentTypes = DocumentTypePatterns{
		AdverseActionLetter: DocumentTypePattern{
			HeaderPatterns: []string{
				"ADVERSE ACTION NOTICE",
				"NOTICE OF ADVERSE ACTION",
				"CREDIT DECISION NOTICE",
				"NOTICE OF ACTION TAKEN",
				"ADVERSE ACTION",
			},
			ContentPatterns: []string{
				"pursuant to.*Fair Credit Reporting Act",
				"credit report.*obtained from",
				"consumer reporting agency",
				"right to obtain.*copy.*credit report",
				"action.*taken.*credit",
				"declined.*credit",
				"denied.*application",
			},
			StatutoryReferences: []string{
				"15 U.S.C. § 1681m",
				"15 USC 1681m",
				"Fair Credit Reporting Act",
				"FCRA",
			},
			RequiredElements: []string{
				"creditor identification",
				"action taken",
				"credit bureau information",
			},
		},
		Summons: DocumentTypePattern{
			HeaderPatterns: []string{
				"SUMMONS",
				"CIVIL SUMMONS",
				"SUMMONS AND COMPLAINT",
				"UNITED STATES DISTRICT COURT",
			},
			ContentPatterns: []string{
				"YOU ARE HEREBY SUMMONED",
				"within.*days.*answer",
				"default judgment.*may be taken",
				"failure to respond",
				"appear and defend",
				"serve.*answer",
			},
			DefendantPatterns: []string{
				`v\.\s*(Equifax|Experian|Trans Union|TransUnion)`,
				`DEFENDANT[S]?:\s*(.*)`,
				`against.*defendant[s]?\s*(.*)`,
			},
			RequiredElements: []string{
				"court name",
				"case number",
				"defendant name",
				"response deadline",
			},
		},
		AttorneyNotes: DocumentTypePattern{
			HeaderPatterns: []string{
				"ATTORNEY NOTES",
				"CASE NOTES",
				"CLIENT MEETING",
				"LEGAL MEMORANDUM",
			},
			ContentPatterns: []string{
				"client.*meeting",
				"case.*summary",
				"legal.*strategy",
				"violation.*analysis",
				"damages.*calculation",
				"fraud.*amount",
			},
			RequiredElements: []string{
				"client information",
				"case details",
				"legal analysis",
			},
		},
		CivilCoverSheet: DocumentTypePattern{
			HeaderPatterns: []string{
				"CIVIL COVER SHEET",
				"JS 44",
				"CIVIL CASE COVER SHEET",
			},
			ContentPatterns: []string{
				"basis of jurisdiction",
				"nature of suit",
				"origin.*proceedings",
				"class action",
				"jury demand",
			},
			RequiredElements: []string{
				"plaintiff information",
				"defendant information",
				"case type",
				"jurisdiction basis",
			},
		},
	}

	return nil
}

func (dc *DocumentClassifier) ClassifyDocument(documentPath string, content string) (*DocumentClassification, error) {
	classification := &DocumentClassification{
		DocumentPath:      documentPath,
		PrimaryType:       DocumentTypeUnknown,
		SecondaryTypes:    []DocumentType{},
		Confidence:        0.0,
		ContentIndicators: []ContentIndicator{},
		ExtractedFields:   make(map[string]interface{}),
		ValidationScore:   0.0,
	}

	scores := make(map[DocumentType]float64)
	indicators := make(map[DocumentType][]ContentIndicator)

	filenameScore := dc.analyzeFilename(documentPath, scores, indicators)

	structureScore := dc.analyzeDocumentStructure(content, scores, indicators)

	contentScore := dc.analyzeContent(content, scores, indicators)

	dc.determinePrimaryType(scores, classification)

	dc.calculateConfidence(classification, filenameScore, structureScore, contentScore)

	for docType, typeIndicators := range indicators {
		if docType == classification.PrimaryType || dc.contains(classification.SecondaryTypes, docType) {
			classification.ContentIndicators = append(classification.ContentIndicators, typeIndicators...)
		}
	}

	return classification, nil
}

func (dc *DocumentClassifier) analyzeFilename(filepath string, scores map[DocumentType]float64, indicators map[DocumentType][]ContentIndicator) float64 {
	filename := strings.ToLower(filepath)
	baseScore := 0.0

	patterns := map[DocumentType][]string{
		DocumentTypeAttorneyNotes:       {"attorney", "atty", "notes", "legal_notes"},
		DocumentTypeAdverseActionLetter: {"adverse", "action", "denial", "declined"},
		DocumentTypeSummons:             {"summons", "civil_summons"},
		DocumentTypeCivilCoverSheet:     {"civil_cover", "js44", "cover_sheet"},
		DocumentTypeComplaint:           {"complaint", "legal_complaint"},
		DocumentTypeDenialLetter:        {"denial", "reject", "declined"},
		DocumentTypeCreditReport:        {"credit_report", "credit_bureau", "equifax", "experian", "transunion"},
		DocumentTypeDisputeLetter:       {"dispute", "contested", "challenge"},
	}

	for docType, typePatterns := range patterns {
		for _, pattern := range typePatterns {
			if strings.Contains(filename, pattern) {
				scores[docType] += 0.3
				indicators[docType] = append(indicators[docType], ContentIndicator{
					Pattern:    pattern,
					MatchType:  "filename",
					Confidence: 0.3,
					Location:   "filename",
				})
				baseScore = 0.3
			}
		}
	}

	return baseScore
}

func (dc *DocumentClassifier) analyzeDocumentStructure(content string, scores map[DocumentType]float64, indicators map[DocumentType][]ContentIndicator) float64 {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return 0.0
	}

	headerContent := strings.Join(lines[:min(10, len(lines))], " ")
	headerUpper := strings.ToUpper(headerContent)

	structureScore := 0.0

	checkHeaders := func(docType DocumentType, patterns []string) {
		for _, pattern := range patterns {
			if strings.Contains(headerUpper, pattern) {
				scores[docType] += 0.4
				indicators[docType] = append(indicators[docType], ContentIndicator{
					Pattern:    pattern,
					MatchType:  "header",
					Confidence: 0.4,
					Location:   "document header",
				})
				structureScore = max(structureScore, 0.4)
			}
		}
	}

	checkHeaders(DocumentTypeAdverseActionLetter, dc.documentTypes.AdverseActionLetter.HeaderPatterns)
	checkHeaders(DocumentTypeSummons, dc.documentTypes.Summons.HeaderPatterns)
	checkHeaders(DocumentTypeAttorneyNotes, dc.documentTypes.AttorneyNotes.HeaderPatterns)
	checkHeaders(DocumentTypeCivilCoverSheet, dc.documentTypes.CivilCoverSheet.HeaderPatterns)

	if strings.Contains(headerUpper, "UNITED STATES DISTRICT COURT") ||
		strings.Contains(headerUpper, "SUPERIOR COURT") ||
		strings.Contains(headerUpper, "CIRCUIT COURT") {
		scores[DocumentTypeSummons] += 0.2
		scores[DocumentTypeComplaint] += 0.2
		structureScore = max(structureScore, 0.2)
	}

	return structureScore
}

func (dc *DocumentClassifier) analyzeContent(content string, scores map[DocumentType]float64, indicators map[DocumentType][]ContentIndicator) float64 {
	contentLower := strings.ToLower(content)
	contentScore := 0.0

	checkContentPatterns := func(docType DocumentType, patterns []string, weight float64) {
		for _, pattern := range patterns {
			re, err := regexp.Compile("(?i)" + pattern)
			if err != nil {
				continue
			}
			
			matches := re.FindAllStringIndex(contentLower, -1)
			if len(matches) > 0 {
				matchWeight := weight * float64(min(len(matches), 3)) / 3.0
				scores[docType] += matchWeight
				
				indicators[docType] = append(indicators[docType], ContentIndicator{
					Pattern:    pattern,
					MatchType:  "content",
					Confidence: matchWeight,
					Location:   fmt.Sprintf("found %d times", len(matches)),
				})
				contentScore = max(contentScore, matchWeight)
			}
		}
	}

	checkContentPatterns(DocumentTypeAdverseActionLetter, dc.documentTypes.AdverseActionLetter.ContentPatterns, 0.3)
	checkContentPatterns(DocumentTypeSummons, dc.documentTypes.Summons.ContentPatterns, 0.3)
	checkContentPatterns(DocumentTypeAttorneyNotes, dc.documentTypes.AttorneyNotes.ContentPatterns, 0.3)
	checkContentPatterns(DocumentTypeCivilCoverSheet, dc.documentTypes.CivilCoverSheet.ContentPatterns, 0.3)

	checkStatutoryReferences := func(docType DocumentType, references []string) {
		for _, ref := range references {
			if strings.Contains(contentLower, strings.ToLower(ref)) {
				scores[docType] += 0.4
				indicators[docType] = append(indicators[docType], ContentIndicator{
					Pattern:    ref,
					MatchType:  "statutory",
					Confidence: 0.4,
					Location:   "statutory reference",
				})
				contentScore = max(contentScore, 0.4)
			}
		}
	}

	checkStatutoryReferences(DocumentTypeAdverseActionLetter, dc.documentTypes.AdverseActionLetter.StatutoryReferences)

	if strings.Contains(contentLower, "fcra") || strings.Contains(contentLower, "fair credit reporting act") {
		scores[DocumentTypeAdverseActionLetter] += 0.3
		scores[DocumentTypeComplaint] += 0.2
		scores[DocumentTypeDisputeLetter] += 0.2
	}

	return contentScore
}

func (dc *DocumentClassifier) determinePrimaryType(scores map[DocumentType]float64, classification *DocumentClassification) {
	var maxScore float64
	classification.PrimaryType = DocumentTypeUnknown

	for docType, score := range scores {
		if score > maxScore {
			maxScore = score
			classification.PrimaryType = docType
		}
	}

	threshold := maxScore * 0.7
	for docType, score := range scores {
		if docType != classification.PrimaryType && score >= threshold && score > 0.3 {
			classification.SecondaryTypes = append(classification.SecondaryTypes, docType)
		}
	}
}

func (dc *DocumentClassifier) calculateConfidence(classification *DocumentClassification, filenameScore, structureScore, contentScore float64) {
	baseConfidence := (filenameScore + structureScore + contentScore) / 3.0

	if classification.PrimaryType != DocumentTypeUnknown {
		typeScore := 0.0
		for _, indicator := range classification.ContentIndicators {
			typeScore = max(typeScore, indicator.Confidence)
		}
		baseConfidence = max(baseConfidence, typeScore)
	}

	multiplier := 1.0
	if len(classification.ContentIndicators) > 3 {
		multiplier = 1.2
	}
	if len(classification.ContentIndicators) > 5 {
		multiplier = 1.4
	}

	classification.Confidence = math.Min(baseConfidence*multiplier, 0.95)

	classification.ValidationScore = classification.Confidence
}

func (dc *DocumentClassifier) GetDocumentTypeName(docType DocumentType) string {
	if name, ok := documentTypeNames[docType]; ok {
		return name
	}
	return "Unknown"
}

func (dc *DocumentClassifier) contains(slice []DocumentType, item DocumentType) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}