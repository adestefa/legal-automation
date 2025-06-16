package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

// ConsistencyValidator validates consistency across multiple documents
type ConsistencyValidator struct {
	ValidationRules         []ConsistencyRule         `json:"validationRules"`
	ConflictDetector        ConflictDetector          `json:"conflictDetector"`
	ReliabilityScorer       DocumentReliabilityScorer `json:"reliabilityScorer"`
	ConflictResolver        ConflictResolver          `json:"conflictResolver"`
}

// ConsistencyRule defines rules for validating consistency
type ConsistencyRule struct {
	RuleID                  string                    `json:"ruleId"`
	RuleName                string                    `json:"ruleName"`
	RuleType                string                    `json:"ruleType"`        // "factual", "temporal", "legal"
	ValidationCriteria      []ValidationCriterion     `json:"validationCriteria"`
	ConsistencyThreshold    float64                   `json:"consistencyThreshold"`
	ConflictSeverity        string                    `json:"conflictSeverity"`
	AutoResolution          bool                      `json:"autoResolution"`
	ResolutionStrategy      ResolutionStrategy        `json:"resolutionStrategy"`
}

// ValidationCriterion defines specific validation criteria
type ValidationCriterion struct {
	CriterionID             string                    `json:"criterionId"`
	CriterionType           string                    `json:"criterionType"`
	FieldToValidate         string                    `json:"fieldToValidate"`
	ValidationMethod        string                    `json:"validationMethod"`
	ToleranceLevel          float64                   `json:"toleranceLevel"`
	RequiredDocumentTypes   []DocumentType            `json:"requiredDocumentTypes"`
	ValidationWeight        float64                   `json:"validationWeight"`
}

// ConflictDetector detects conflicts between documents
type ConflictDetector struct {
	DetectionRules          []ConflictDetectionRule   `json:"detectionRules"`
	SeverityClassifier      SeverityClassifier        `json:"severityClassifier"`
	PatternAnalyzer         ConflictPatternAnalyzer   `json:"patternAnalyzer"`
}

// ConflictDetectionRule defines rules for detecting conflicts
type ConflictDetectionRule struct {
	RuleID                  string                    `json:"ruleId"`
	ConflictType            string                    `json:"conflictType"`
	TriggerConditions       []TriggerCondition        `json:"triggerConditions"`
	ConflictWeight          float64                   `json:"conflictWeight"`
	ImpactAssessment        ImpactLevel               `json:"impactAssessment"`
}

// TriggerCondition defines conditions that trigger conflict detection
type TriggerCondition struct {
	ConditionType           string                    `json:"conditionType"`
	FieldName               string                    `json:"fieldName"`
	ComparisonOperator      string                    `json:"comparisonOperator"`
	ThresholdValue          interface{}               `json:"thresholdValue"`
	DocumentScope           []DocumentType            `json:"documentScope"`
}

// SeverityClassifier classifies the severity of conflicts
type SeverityClassifier struct {
	SeverityRules           []SeverityRule            `json:"severityRules"`
	ImpactWeights           map[string]float64        `json:"impactWeights"`
	CriticalThreshold       float64                   `json:"criticalThreshold"`
	SignificantThreshold    float64                   `json:"significantThreshold"`
}

// SeverityRule defines rules for classifying conflict severity
type SeverityRule struct {
	RuleID                  string                    `json:"ruleId"`
	ConflictTypes           []string                  `json:"conflictTypes"`
	SeverityLevel           ConflictSeverityLevel     `json:"severityLevel"`
	SeverityFactors         []SeverityFactor          `json:"severityFactors"`
	ImpactMultiplier        float64                   `json:"impactMultiplier"`
}

// SeverityFactor defines factors that affect conflict severity
type SeverityFactor struct {
	FactorType              string                    `json:"factorType"`
	FactorWeight            float64                   `json:"factorWeight"`
	FactorDescription       string                    `json:"factorDescription"`
}

// ConflictPatternAnalyzer analyzes patterns in conflicts
type ConflictPatternAnalyzer struct {
	PatternRules            []ConflictPatternRule     `json:"patternRules"`
	PatternDatabase         []ConflictPattern         `json:"patternDatabase"`
	LearningEngine          PatternLearningEngine     `json:"learningEngine"`
}

// ConflictPatternRule defines rules for detecting conflict patterns
type ConflictPatternRule struct {
	RuleID                  string                    `json:"ruleId"`
	PatternType             string                    `json:"patternType"`
	PatternDescription      string                    `json:"patternDescription"`
	DetectionCriteria       []PatternCriterion        `json:"detectionCriteria"`
	PatternSignificance     PatternSignificanceLevel  `json:"patternSignificance"`
}

// PatternCriterion defines criteria for pattern detection
type PatternCriterion struct {
	CriterionType           string                    `json:"criterionType"`
	MinOccurrences          int                       `json:"minOccurrences"`
	DocumentSpread          int                       `json:"documentSpread"`
	TemporalPattern         string                    `json:"temporalPattern"`
	SeverityPattern         string                    `json:"severityPattern"`
}

// ConflictPattern represents a detected conflict pattern
type ConflictPattern struct {
	PatternID               string                    `json:"patternId"`
	PatternType             string                    `json:"patternType"`
	PatternDescription      string                    `json:"patternDescription"`
	AffectedDocuments       []string                  `json:"affectedDocuments"`
	ConflictInstances       []string                  `json:"conflictInstances"`
	PatternStrength         float64                   `json:"patternStrength"`
	SystemicIndicator       bool                      `json:"systemicIndicator"`
	ResolutionComplexity    ComplexityLevel           `json:"resolutionComplexity"`
}

// PatternLearningEngine learns from conflict resolution patterns
type PatternLearningEngine struct {
	LearningDatabase        []ResolutionExample       `json:"learningDatabase"`
	SuccessPatterns         []SuccessPattern          `json:"successPatterns"`
	FailurePatterns         []FailurePattern          `json:"failurePatterns"`
	AdaptationRules         []AdaptationRule          `json:"adaptationRules"`
}

// DocumentReliabilityScorer scores document reliability
type DocumentReliabilityScorer struct {
	ScoringRules            []ReliabilityRule         `json:"scoringRules"`
	SourceCredibility       map[DocumentType]float64  `json:"sourceCredibility"`
	TemporalFactors         TemporalReliabilityFactors `json:"temporalFactors"`
	ConsistencyWeights      ConsistencyWeights        `json:"consistencyWeights"`
}

// ReliabilityRule defines rules for scoring document reliability
type ReliabilityRule struct {
	RuleID                  string                    `json:"ruleId"`
	DocumentType            DocumentType              `json:"documentType"`
	BaseReliability         float64                   `json:"baseReliability"`
	ReliabilityFactors      []ReliabilityFactor       `json:"reliabilityFactors"`
	AdjustmentRules         []ReliabilityAdjustment   `json:"adjustmentRules"`
}

// ReliabilityFactor defines factors affecting document reliability
type ReliabilityFactor struct {
	FactorType              string                    `json:"factorType"`
	FactorWeight            float64                   `json:"factorWeight"`
	FactorDescription       string                    `json:"factorDescription"`
	ImpactDirection         string                    `json:"impactDirection"`  // "positive", "negative"
}

// ReliabilityAdjustment defines adjustments to reliability scores
type ReliabilityAdjustment struct {
	AdjustmentType          string                    `json:"adjustmentType"`
	TriggerCondition        string                    `json:"triggerCondition"`
	AdjustmentValue         float64                   `json:"adjustmentValue"`
	AdjustmentReason        string                    `json:"adjustmentReason"`
}

// TemporalReliabilityFactors defines time-based reliability factors
type TemporalReliabilityFactors struct {
	DocumentAge             map[string]float64        `json:"documentAge"`
	RecencyBonus            float64                   `json:"recencyBonus"`
	StalenessePenalty       float64                   `json:"stalenessPenalty"`
	TemporalConsistency     float64                   `json:"temporalConsistency"`
}

// ConsistencyWeights defines weights for consistency calculations
type ConsistencyWeights struct {
	FactualConsistency      float64                   `json:"factualConsistency"`
	TemporalConsistency     float64                   `json:"temporalConsistency"`
	LegalConsistency        float64                   `json:"legalConsistency"`
	SourceConsistency       float64                   `json:"sourceConsistency"`
}

// ConflictResolver resolves conflicts between documents
type ConflictResolver struct {
	ResolutionStrategies    []ResolutionStrategy      `json:"resolutionStrategies"`
	PriorityMatrix          PriorityMatrix            `json:"priorityMatrix"`
	ResolutionHistory       []ResolutionRecord        `json:"resolutionHistory"`
	AutoResolutionEngine    AutoResolutionEngine      `json:"autoResolutionEngine"`
}

// ResolutionStrategy defines strategies for conflict resolution
type ResolutionStrategy struct {
	StrategyID              string                    `json:"strategyId"`
	StrategyName            string                    `json:"strategyName"`
	StrategyType            string                    `json:"strategyType"`
	ApplicableConflicts     []string                  `json:"applicableConflicts"`
	ResolutionSteps         []ResolutionStep          `json:"resolutionSteps"`
	SuccessRate             float64                   `json:"successRate"`
	AutomationLevel         AutomationLevel           `json:"automationLevel"`
}

// ResolutionStep defines individual steps in conflict resolution
type ResolutionStep struct {
	StepID                  string                    `json:"stepId"`
	StepType                string                    `json:"stepType"`
	StepDescription         string                    `json:"stepDescription"`
	RequiredInputs          []string                  `json:"requiredInputs"`
	ExpectedOutputs         []string                  `json:"expectedOutputs"`
	ValidationCriteria      []string                  `json:"validationCriteria"`
	FallbackAction          string                    `json:"fallbackAction"`
}

// PriorityMatrix defines priorities for conflict resolution
type PriorityMatrix struct {
	DocumentTypePriority    map[DocumentType]float64  `json:"documentTypePriority"`
	ConflictTypePriority    map[string]float64        `json:"conflictTypePriority"`
	SourceCredibilityMatrix map[string]float64        `json:"sourceCredibilityMatrix"`
	TemporalPriorityRules   []TemporalPriorityRule    `json:"temporalPriorityRules"`
}

// TemporalPriorityRule defines time-based priority rules
type TemporalPriorityRule struct {
	RuleID                  string                    `json:"ruleId"`
	TimeRange               string                    `json:"timeRange"`
	PriorityMultiplier      float64                   `json:"priorityMultiplier"`
	ApplicableDocumentTypes []DocumentType            `json:"applicableDocumentTypes"`
}

// ResolutionRecord records conflict resolution history
type ResolutionRecord struct {
	RecordID                string                    `json:"recordId"`
	ConflictID              string                    `json:"conflictId"`
	ResolutionStrategy      string                    `json:"resolutionStrategy"`
	ResolutionResult        ResolutionResult          `json:"resolutionResult"`
	ResolutionTimestamp     time.Time                 `json:"resolutionTimestamp"`
	ResolutionConfidence    float64                   `json:"resolutionConfidence"`
	PostResolutionValidation PostResolutionValidation `json:"postResolutionValidation"`
}

// AutoResolutionEngine handles automatic conflict resolution
type AutoResolutionEngine struct {
	AutoResolutionRules     []AutoResolutionRule      `json:"autoResolutionRules"`
	ConfidenceThresholds    ConfidenceThresholds      `json:"confidenceThresholds"`
	SafetyChecks            []SafetyCheck             `json:"safetyChecks"`
	EscalationCriteria      []EscalationCriterion     `json:"escalationCriteria"`
}

// AutoResolutionRule defines rules for automatic resolution
type AutoResolutionRule struct {
	RuleID                  string                    `json:"ruleId"`
	ConflictTypes           []string                  `json:"conflictTypes"`
	ResolutionLogic         ResolutionLogic           `json:"resolutionLogic"`
	ConfidenceRequirement   float64                   `json:"confidenceRequirement"`
	SafetyRequirements      []string                  `json:"safetyRequirements"`
}

// ResolutionLogic defines the logic for automatic resolution
type ResolutionLogic struct {
	LogicType               string                    `json:"logicType"`
	DecisionTree            []DecisionNode            `json:"decisionTree"`
	WeightingFactors        map[string]float64        `json:"weightingFactors"`
	FallbackStrategy        string                    `json:"fallbackStrategy"`
}

// DecisionNode represents a node in the decision tree
type DecisionNode struct {
	NodeID                  string                    `json:"nodeId"`
	Condition               string                    `json:"condition"`
	TrueAction              string                    `json:"trueAction"`
	FalseAction             string                    `json:"falseAction"`
	ConfidenceImpact        float64                   `json:"confidenceImpact"`
}

// ConfidenceThresholds defines thresholds for confidence-based decisions
type ConfidenceThresholds struct {
	AutoResolutionThreshold float64                   `json:"autoResolutionThreshold"`
	EscalationThreshold     float64                   `json:"escalationThreshold"`
	ManualReviewThreshold   float64                   `json:"manualReviewThreshold"`
	RejectionThreshold      float64                   `json:"rejectionThreshold"`
}

// SafetyCheck defines safety checks for automatic resolution
type SafetyCheck struct {
	CheckID                 string                    `json:"checkId"`
	CheckType               string                    `json:"checkType"`
	CheckDescription        string                    `json:"checkDescription"`
	CheckCriteria           []string                  `json:"checkCriteria"`
	FailureAction           string                    `json:"failureAction"`
}

// EscalationCriterion defines criteria for escalating conflicts
type EscalationCriterion struct {
	CriterionID             string                    `json:"criterionId"`
	TriggerCondition        string                    `json:"triggerCondition"`
	EscalationLevel         EscalationLevel           `json:"escalationLevel"`
	EscalationAction        string                    `json:"escalationAction"`
	NotificationRequired    bool                      `json:"notificationRequired"`
}

// Enums for type safety
type ConflictSeverityLevel string
type PatternSignificanceLevel string
type ComplexityLevel string
type AutomationLevel string
type ImpactLevel string
type EscalationLevel string

const (
	ConflictSeverityCritical     ConflictSeverityLevel = "critical"
	ConflictSeveritySignificant  ConflictSeverityLevel = "significant"
	ConflictSeverityMinor        ConflictSeverityLevel = "minor"

	PatternSignificanceHigh      PatternSignificanceLevel = "high"
	PatternSignificanceMedium    PatternSignificanceLevel = "medium"
	PatternSignificanceLow       PatternSignificanceLevel = "low"

	ComplexityHigh               ComplexityLevel = "high"
	ComplexityMedium             ComplexityLevel = "medium"
	ComplexityLow                ComplexityLevel = "low"

	AutomationFull               AutomationLevel = "full"
	AutomationPartial            AutomationLevel = "partial"
	AutomationManual             AutomationLevel = "manual"

	ImpactHigh                   ImpactLevel = "high"
	ImpactMedium                 ImpactLevel = "medium"
	ImpactLow                    ImpactLevel = "low"

	EscalationImmediate          EscalationLevel = "immediate"
	EscalationUrgent             EscalationLevel = "urgent"
	EscalationRoutine            EscalationLevel = "routine"
)

// Result types
type ResolutionResult struct {
	ResultType              string                    `json:"resultType"`
	ResolvedValue           string                    `json:"resolvedValue"`
	ResolutionMethod        string                    `json:"resolutionMethod"`
	ResolutionConfidence    float64                   `json:"resolutionConfidence"`
	ResolutionEvidence      []string                  `json:"resolutionEvidence"`
}

type PostResolutionValidation struct {
	ValidationPerformed     bool                      `json:"validationPerformed"`
	ValidationResults       []ValidationResult        `json:"validationResults"`
	ValidationConfidence    float64                   `json:"validationConfidence"`
	RecommendedActions      []string                  `json:"recommendedActions"`
}

type ValidationResult struct {
	ValidationType          string                    `json:"validationType"`
	ValidationPassed        bool                      `json:"validationPassed"`
	ValidationScore         float64                   `json:"validationScore"`
	ValidationNotes         string                    `json:"validationNotes"`
}

type ResolutionExample struct {
	ExampleID               string                    `json:"exampleId"`
	ConflictType            string                    `json:"conflictType"`
	ResolutionStrategy      string                    `json:"resolutionStrategy"`
	SuccessRating           float64                   `json:"successRating"`
	LessonsLearned          []string                  `json:"lessonsLearned"`
}

type SuccessPattern struct {
	PatternID               string                    `json:"patternId"`
	PatternDescription      string                    `json:"patternDescription"`
	SuccessFactors          []string                  `json:"successFactors"`
	ApplicationGuidance     []string                  `json:"applicationGuidance"`
}

type FailurePattern struct {
	PatternID               string                    `json:"patternId"`
	PatternDescription      string                    `json:"patternDescription"`
	FailureFactors          []string                  `json:"failureFactors"`
	AvoidanceGuidance       []string                  `json:"avoidanceGuidance"`
}

type AdaptationRule struct {
	RuleID                  string                    `json:"ruleId"`
	AdaptationTrigger       string                    `json:"adaptationTrigger"`
	AdaptationAction        string                    `json:"adaptationAction"`
	LearningWeight          float64                   `json:"learningWeight"`
}

// NewConsistencyValidator creates a new consistency validator
func NewConsistencyValidator() *ConsistencyValidator {
	validator := &ConsistencyValidator{
		ValidationRules:  []ConsistencyRule{},
		ConflictDetector: ConflictDetector{},
		ReliabilityScorer: DocumentReliabilityScorer{},
		ConflictResolver: ConflictResolver{},
	}
	
	// Load validation rules from configuration
	validator.loadValidationRules()
	
	// Initialize components
	validator.initializeConflictDetector()
	validator.initializeReliabilityScorer()
	validator.initializeConflictResolver()
	
	return validator
}

// loadValidationRules loads validation rules from configuration
func (cv *ConsistencyValidator) loadValidationRules() {
	configFile := "v2/config/consistency_validation_rules.json"
	
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Warning: Could not load validation rules from %s: %v", configFile, err)
		cv.createDefaultValidationRules()
		return
	}
	
	var config struct {
		ValidationRules []ConsistencyRule `json:"validationRules"`
	}
	
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("Error parsing validation rules: %v", err)
		cv.createDefaultValidationRules()
		return
	}
	
	cv.ValidationRules = config.ValidationRules
	log.Printf("Loaded %d validation rules", len(cv.ValidationRules))
}

// createDefaultValidationRules creates default validation rules
func (cv *ConsistencyValidator) createDefaultValidationRules() {
	cv.ValidationRules = []ConsistencyRule{
		{
			RuleID:               "client_name_consistency",
			RuleName:             "Client Name Consistency Validation",
			RuleType:             "factual",
			ConsistencyThreshold: 0.85,
			ConflictSeverity:     "significant",
			AutoResolution:       true,
			ValidationCriteria: []ValidationCriterion{
				{
					CriterionID:           "client_name_exact_match",
					CriterionType:         "exact_match",
					FieldToValidate:       "client_name",
					ValidationMethod:      "string_comparison",
					ToleranceLevel:        0.95,
					RequiredDocumentTypes: []DocumentType{DocTypeAttorneyNotes, DocTypeAdverseAction, DocTypeSummons},
					ValidationWeight:      1.0,
				},
			},
		},
		{
			RuleID:               "temporal_consistency",
			RuleName:             "Temporal Event Consistency",
			RuleType:             "temporal",
			ConsistencyThreshold: 0.8,
			ConflictSeverity:     "minor",
			AutoResolution:       false,
			ValidationCriteria: []ValidationCriterion{
				{
					CriterionID:           "event_date_proximity",
					CriterionType:         "temporal_proximity",
					FieldToValidate:       "event_dates",
					ValidationMethod:      "temporal_analysis",
					ToleranceLevel:        0.7,
					RequiredDocumentTypes: []DocumentType{DocTypeAdverseAction, DocTypeAttorneyNotes},
					ValidationWeight:      0.8,
				},
			},
		},
		{
			RuleID:               "legal_consistency",
			RuleName:             "Legal Analysis Consistency",
			RuleType:             "legal",
			ConsistencyThreshold: 0.75,
			ConflictSeverity:     "critical",
			AutoResolution:       false,
			ValidationCriteria: []ValidationCriterion{
				{
					CriterionID:           "violation_consistency",
					CriterionType:         "legal_consistency",
					FieldToValidate:       "identified_violations",
					ValidationMethod:      "legal_analysis",
					ToleranceLevel:        0.8,
					RequiredDocumentTypes: []DocumentType{DocTypeAdverseAction, DocTypeAttorneyNotes},
					ValidationWeight:      1.0,
				},
			},
		},
	}
	
	log.Printf("Created %d default validation rules", len(cv.ValidationRules))
}

// initializeConflictDetector initializes the conflict detector
func (cv *ConsistencyValidator) initializeConflictDetector() {
	cv.ConflictDetector = ConflictDetector{
		DetectionRules: []ConflictDetectionRule{
			{
				RuleID:       "factual_conflict",
				ConflictType: "factual_inconsistency",
				TriggerConditions: []TriggerCondition{
					{
						ConditionType:      "value_mismatch",
						FieldName:          "client_name",
						ComparisonOperator: "not_equal",
						ThresholdValue:     0.85,
						DocumentScope:      []DocumentType{DocTypeAttorneyNotes, DocTypeAdverseAction},
					},
				},
				ConflictWeight:   1.0,
				ImpactAssessment: ImpactHigh,
			},
		},
		SeverityClassifier: SeverityClassifier{
			CriticalThreshold:    0.9,
			SignificantThreshold: 0.6,
			ImpactWeights: map[string]float64{
				"factual_inconsistency": 0.8,
				"temporal_inconsistency": 0.6,
				"legal_inconsistency": 0.9,
			},
		},
	}
}

// initializeReliabilityScorer initializes the reliability scorer
func (cv *ConsistencyValidator) initializeReliabilityScorer() {
	cv.ReliabilityScorer = DocumentReliabilityScorer{
		SourceCredibility: map[DocumentType]float64{
			DocTypeAttorneyNotes: 0.95,
			DocTypeAdverseAction: 0.9,
			DocTypeSummons:      0.9,
			DocTypeCivilCover:   0.85,
			DocTypeCreditReport: 0.8,
			DocTypeCorrespondence: 0.7,
			DocTypeOther:        0.6,
		},
		TemporalFactors: TemporalReliabilityFactors{
			RecencyBonus:      0.1,
			StalenessePenalty: 0.05,
			TemporalConsistency: 0.15,
		},
		ConsistencyWeights: ConsistencyWeights{
			FactualConsistency:  0.3,
			TemporalConsistency: 0.2,
			LegalConsistency:   0.3,
			SourceConsistency:  0.2,
		},
	}
}

// initializeConflictResolver initializes the conflict resolver
func (cv *ConsistencyValidator) initializeConflictResolver() {
	cv.ConflictResolver = ConflictResolver{
		PriorityMatrix: PriorityMatrix{
			DocumentTypePriority: map[DocumentType]float64{
				DocTypeAttorneyNotes: 1.0,
				DocTypeAdverseAction: 0.9,
				DocTypeSummons:      0.85,
				DocTypeCivilCover:   0.8,
				DocTypeCreditReport: 0.75,
				DocTypeCorrespondence: 0.7,
				DocTypeOther:        0.5,
			},
			ConflictTypePriority: map[string]float64{
				"factual_inconsistency": 1.0,
				"legal_inconsistency":   0.95,
				"temporal_inconsistency": 0.8,
			},
		},
		AutoResolutionEngine: AutoResolutionEngine{
			ConfidenceThresholds: ConfidenceThresholds{
				AutoResolutionThreshold: 0.85,
				EscalationThreshold:     0.6,
				ManualReviewThreshold:   0.4,
				RejectionThreshold:      0.2,
			},
		},
	}
}

// ValidateDocumentConsistency validates consistency across multiple documents
func (cv *ConsistencyValidator) ValidateDocumentConsistency(
	documents []DocumentAnalysis,
) ConsistencyValidationResult {
	
	result := ConsistencyValidationResult{
		ValidationID:          fmt.Sprintf("validation_%d", time.Now().Unix()),
		DocumentCount:         len(documents),
		OverallConsistency:    0.0,
		FactConsistencies:     []FactConsistency{},
		DetectedConflicts:     []DocumentConflict{},
		ReliabilityScores:     make(map[string]float64),
		ResolutionSuggestions: []ConflictResolution{},
		ValidationSummary:     ValidationSummary{},
	}
	
	// Calculate reliability scores for each document
	for _, doc := range documents {
		reliabilityScore := cv.ReliabilityScorer.CalculateReliabilityScore(doc)
		result.ReliabilityScores[doc.DocumentPath] = reliabilityScore
	}
	
	// Validate key facts across documents
	keyFacts := cv.extractKeyFacts(documents)
	for _, fact := range keyFacts {
		consistency := cv.validateFactConsistency(fact, documents)
		result.FactConsistencies = append(result.FactConsistencies, consistency)
	}
	
	// Detect conflicts between documents
	conflicts := cv.ConflictDetector.DetectConflicts(documents)
	result.DetectedConflicts = conflicts
	
	// Generate resolution suggestions for conflicts
	for _, conflict := range conflicts {
		if conflict.ConflictSeverity != "minor" {
			resolution := cv.ConflictResolver.ResolveConflict(conflict, documents)
			result.ResolutionSuggestions = append(result.ResolutionSuggestions, resolution)
		}
	}
	
	// Calculate overall consistency score
	result.OverallConsistency = cv.calculateOverallConsistency(result)
	
	// Generate validation summary
	result.ValidationSummary = cv.generateValidationSummary(result)
	
	return result
}

// ConsistencyValidationResult represents the result of consistency validation
type ConsistencyValidationResult struct {
	ValidationID          string                    `json:"validationId"`
	DocumentCount         int                       `json:"documentCount"`
	OverallConsistency    float64                   `json:"overallConsistency"`
	FactConsistencies     []FactConsistency         `json:"factConsistencies"`
	DetectedConflicts     []DocumentConflict        `json:"detectedConflicts"`
	ReliabilityScores     map[string]float64        `json:"reliabilityScores"`
	ResolutionSuggestions []ConflictResolution      `json:"resolutionSuggestions"`
	ValidationSummary     ValidationSummary         `json:"validationSummary"`
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalFactsValidated   int                       `json:"totalFactsValidated"`
	ConsistentFacts       int                       `json:"consistentFacts"`
	ConflictingFacts      int                       `json:"conflictingFacts"`
	CriticalConflicts     int                       `json:"criticalConflicts"`
	SignificantConflicts  int                       `json:"significantConflicts"`
	MinorConflicts        int                       `json:"minorConflicts"`
	AutoResolvedConflicts int                       `json:"autoResolvedConflicts"`
	ManualReviewRequired  int                       `json:"manualReviewRequired"`
	AverageReliability    float64                   `json:"averageReliability"`
	RecommendedActions    []string                  `json:"recommendedActions"`
}

// extractKeyFacts extracts key facts that need validation across documents
func (cv *ConsistencyValidator) extractKeyFacts(documents []DocumentAnalysis) []string {
	keyFactTypes := []string{
		"client_name",
		"case_amount",
		"event_dates",
		"identified_violations",
		"defendant_names",
		"court_jurisdiction",
	}
	
	return keyFactTypes
}

// validateFactConsistency validates consistency of a specific fact across documents
func (cv *ConsistencyValidator) validateFactConsistency(factType string, documents []DocumentAnalysis) FactConsistency {
	consistency := FactConsistency{
		FactType:            factType,
		ConsistentValue:     "",
		SupportingDocuments: []string{},
		ConsistencyScore:    0.0,
		ConfidenceLevel:     0.0,
	}
	
	// Collect all values for this fact type across documents
	factValues := make(map[string][]string)
	
	for _, doc := range documents {
		for _, fact := range doc.ExtractedFacts {
			if fact.FactType == factType {
				if _, exists := factValues[fact.FactValue]; !exists {
					factValues[fact.FactValue] = []string{}
				}
				factValues[fact.FactValue] = append(factValues[fact.FactValue], doc.DocumentPath)
			}
		}
	}
	
	// Determine most consistent value
	if len(factValues) == 1 {
		// All documents agree
		for value, docs := range factValues {
			consistency.ConsistentValue = value
			consistency.SupportingDocuments = docs
			consistency.ConsistencyScore = 1.0
			consistency.ConfidenceLevel = 0.95
		}
	} else if len(factValues) > 1 {
		// Multiple values - find most common
		maxCount := 0
		mostCommonValue := ""
		var mostCommonDocs []string
		
		for value, docs := range factValues {
			if len(docs) > maxCount {
				maxCount = len(docs)
				mostCommonValue = value
				mostCommonDocs = docs
			}
		}
		
		consistency.ConsistentValue = mostCommonValue
		consistency.SupportingDocuments = mostCommonDocs
		consistency.ConsistencyScore = float64(maxCount) / float64(len(documents))
		consistency.ConfidenceLevel = consistency.ConsistencyScore * 0.8
	}
	
	return consistency
}

// DetectConflicts detects conflicts between documents
func (cd *ConflictDetector) DetectConflicts(documents []DocumentAnalysis) []DocumentConflict {
	conflicts := []DocumentConflict{}
	
	// Compare documents pairwise for conflicts
	for i := 0; i < len(documents); i++ {
		for j := i + 1; j < len(documents); j++ {
			pairConflicts := cd.detectPairwiseConflicts(documents[i], documents[j])
			conflicts = append(conflicts, pairConflicts...)
		}
	}
	
	// Classify conflict severity
	for i := range conflicts {
		conflicts[i].ConflictSeverity = cd.SeverityClassifier.ClassifyConflictSeverity(conflicts[i])
	}
	
	return conflicts
}

// detectPairwiseConflicts detects conflicts between two specific documents
func (cd *ConflictDetector) detectPairwiseConflicts(doc1, doc2 DocumentAnalysis) []DocumentConflict {
	conflicts := []DocumentConflict{}
	
	// Check for factual conflicts
	factualConflicts := cd.detectFactualConflicts(doc1, doc2)
	conflicts = append(conflicts, factualConflicts...)
	
	// Check for temporal conflicts
	temporalConflicts := cd.detectTemporalConflicts(doc1, doc2)
	conflicts = append(conflicts, temporalConflicts...)
	
	// Check for legal conflicts
	legalConflicts := cd.detectLegalConflicts(doc1, doc2)
	conflicts = append(conflicts, legalConflicts...)
	
	return conflicts
}

// detectFactualConflicts detects factual conflicts between documents
func (cd *ConflictDetector) detectFactualConflicts(doc1, doc2 DocumentAnalysis) []DocumentConflict {
	conflicts := []DocumentConflict{}
	
	// Compare facts of the same type
	for _, fact1 := range doc1.ExtractedFacts {
		for _, fact2 := range doc2.ExtractedFacts {
			if fact1.FactType == fact2.FactType && fact1.FactValue != fact2.FactValue {
				// Check if this represents a significant conflict
				similarity := cd.calculateSimilarity(fact1.FactValue, fact2.FactValue)
				if similarity < 0.7 { // Threshold for conflict detection
					conflict := DocumentConflict{
						ConflictID:           fmt.Sprintf("conflict_%s_%d", fact1.FactType, time.Now().Unix()),
						ConflictType:         "factual_inconsistency",
						ConflictingDocuments: []string{doc1.DocumentPath, doc2.DocumentPath},
						ConflictDescription:  fmt.Sprintf("Conflicting values for %s: '%s' vs '%s'", fact1.FactType, fact1.FactValue, fact2.FactValue),
						ConflictSeverity:     "pending_classification",
						ConflictingFacts: []ConflictingFact{
							{
								FactType: fact1.FactType,
								ConflictingValues: []ConflictingValue{
									{
										Value:           fact1.FactValue,
										SourceDocument:  doc1.DocumentPath,
										ConfidenceLevel: fact1.ConfidenceLevel,
									},
									{
										Value:           fact2.FactValue,
										SourceDocument:  doc2.DocumentPath,
										ConfidenceLevel: fact2.ConfidenceLevel,
									},
								},
								ConflictReason:   "Different values extracted for same fact type",
								ImpactAssessment: "Requires resolution for case consistency",
							},
						},
					}
					conflicts = append(conflicts, conflict)
				}
			}
		}
	}
	
	return conflicts
}

// detectTemporalConflicts detects temporal conflicts between documents
func (cd *ConflictDetector) detectTemporalConflicts(doc1, doc2 DocumentAnalysis) []DocumentConflict {
	conflicts := []DocumentConflict{}
	
	// Compare timeline events for conflicts
	for _, event1 := range doc1.Timeline {
		for _, event2 := range doc2.Timeline {
			if event1.EventType == event2.EventType {
				// Check for temporal inconsistency
				timeDiff := math.Abs(event1.EventDate.Sub(event2.EventDate).Hours())
				if timeDiff > 24*7 { // More than a week difference
					conflict := DocumentConflict{
						ConflictID:           fmt.Sprintf("temporal_conflict_%d", time.Now().Unix()),
						ConflictType:         "temporal_inconsistency",
						ConflictingDocuments: []string{doc1.DocumentPath, doc2.DocumentPath},
						ConflictDescription:  fmt.Sprintf("Temporal inconsistency for %s: %s vs %s", event1.EventType, event1.EventDate.Format("2006-01-02"), event2.EventDate.Format("2006-01-02")),
						ConflictSeverity:     "pending_classification",
					}
					conflicts = append(conflicts, conflict)
				}
			}
		}
	}
	
	return conflicts
}

// detectLegalConflicts detects legal conflicts between documents
func (cd *ConflictDetector) detectLegalConflicts(doc1, doc2 DocumentAnalysis) []DocumentConflict {
	conflicts := []DocumentConflict{}
	
	// Compare legal conclusions for conflicts
	for _, conclusion1 := range doc1.LegalConclusions {
		for _, conclusion2 := range doc2.LegalConclusions {
			if conclusion1.ConclusionType == conclusion2.ConclusionType {
				// Check for conflicting legal conclusions
				similarity := cd.calculateSimilarity(conclusion1.ConclusionText, conclusion2.ConclusionText)
				if similarity < 0.5 {
					conflict := DocumentConflict{
						ConflictID:           fmt.Sprintf("legal_conflict_%d", time.Now().Unix()),
						ConflictType:         "legal_inconsistency",
						ConflictingDocuments: []string{doc1.DocumentPath, doc2.DocumentPath},
						ConflictDescription:  fmt.Sprintf("Conflicting legal conclusions for %s", conclusion1.ConclusionType),
						ConflictSeverity:     "pending_classification",
					}
					conflicts = append(conflicts, conflict)
				}
			}
		}
	}
	
	return conflicts
}

// calculateSimilarity calculates similarity between two strings
func (cd *ConflictDetector) calculateSimilarity(str1, str2 string) float64 {
	// Simple similarity calculation - could be enhanced with more sophisticated algorithms
	str1Lower := strings.ToLower(str1)
	str2Lower := strings.ToLower(str2)
	
	if str1Lower == str2Lower {
		return 1.0
	}
	
	// Check for substring matches
	if strings.Contains(str1Lower, str2Lower) || strings.Contains(str2Lower, str1Lower) {
		shorter := len(str1Lower)
		if len(str2Lower) < shorter {
			shorter = len(str2Lower)
		}
		longer := len(str1Lower)
		if len(str2Lower) > longer {
			longer = len(str2Lower)
		}
		return float64(shorter) / float64(longer)
	}
	
	// Basic character overlap
	overlap := 0
	for _, char := range str1Lower {
		if strings.ContainsRune(str2Lower, char) {
			overlap++
		}
	}
	
	return float64(overlap) / float64(len(str1Lower))
}

// ClassifyConflictSeverity classifies the severity of a conflict
func (sc *SeverityClassifier) ClassifyConflictSeverity(conflict DocumentConflict) string {
	// Calculate severity score based on conflict type and impact
	severityScore := 0.0
	
	if weight, exists := sc.ImpactWeights[conflict.ConflictType]; exists {
		severityScore = weight
	} else {
		severityScore = 0.5 // Default weight
	}
	
	// Adjust based on confidence levels of conflicting facts
	if len(conflict.ConflictingFacts) > 0 {
		avgConfidence := 0.0
		for _, fact := range conflict.ConflictingFacts {
			for _, value := range fact.ConflictingValues {
				avgConfidence += value.ConfidenceLevel
			}
		}
		avgConfidence /= float64(len(conflict.ConflictingFacts) * len(conflict.ConflictingFacts[0].ConflictingValues))
		
		// Higher confidence in conflicting facts increases severity
		severityScore *= avgConfidence
	}
	
	// Classify based on thresholds
	if severityScore >= sc.CriticalThreshold {
		return "critical"
	} else if severityScore >= sc.SignificantThreshold {
		return "significant"
	} else {
		return "minor"
	}
}

// CalculateReliabilityScore calculates reliability score for a document
func (drs *DocumentReliabilityScorer) CalculateReliabilityScore(doc DocumentAnalysis) float64 {
	// Start with base reliability for document type
	baseReliability := 0.7 // Default
	if reliability, exists := drs.SourceCredibility[doc.DocumentType]; exists {
		baseReliability = reliability
	}
	
	// Apply consistency weights
	weights := drs.ConsistencyWeights
	consistencyScore := doc.ConfidenceScores.OverallConfidence
	
	// Calculate weighted reliability
	reliabilityScore := baseReliability * (
		weights.FactualConsistency*consistencyScore +
		weights.TemporalConsistency*consistencyScore +
		weights.LegalConsistency*consistencyScore +
		weights.SourceConsistency*1.0) // Source consistency assumed 1.0 for single document
	
	// Apply temporal factors (if document has timestamp information)
	reliabilityScore += drs.TemporalFactors.RecencyBonus * 0.1 // Simplified
	
	// Ensure score is within bounds
	if reliabilityScore > 1.0 {
		reliabilityScore = 1.0
	}
	if reliabilityScore < 0.0 {
		reliabilityScore = 0.0
	}
	
	return reliabilityScore
}

// ResolveConflict resolves a conflict between documents
func (cr *ConflictResolver) ResolveConflict(conflict DocumentConflict, documents []DocumentAnalysis) ConflictResolution {
	resolution := ConflictResolution{
		ResolutionID:         fmt.Sprintf("resolution_%d", time.Now().Unix()),
		ConflictID:           conflict.ConflictID,
		ResolutionMethod:     "priority_based",
		ResolvedValue:        "",
		ResolutionConfidence: 0.0,
		ResolutionReasoning:  "",
		ImpactOnCase:        "pending_assessment",
	}
	
	// Use priority matrix to resolve conflict
	if len(conflict.ConflictingFacts) > 0 {
		fact := conflict.ConflictingFacts[0]
		bestValue := ""
		highestPriority := 0.0
		
		for _, value := range fact.ConflictingValues {
			// Find document type for this value
			var docType DocumentType
			for _, doc := range documents {
				if doc.DocumentPath == value.SourceDocument {
					docType = doc.DocumentType
					break
				}
			}
			
			// Calculate priority score
			docPriority := cr.PriorityMatrix.DocumentTypePriority[docType]
			conflictPriority := cr.PriorityMatrix.ConflictTypePriority[conflict.ConflictType]
			confidenceFactor := value.ConfidenceLevel
			
			priorityScore := docPriority * conflictPriority * confidenceFactor
			
			if priorityScore > highestPriority {
				highestPriority = priorityScore
				bestValue = value.Value
			}
		}
		
		resolution.ResolvedValue = bestValue
		resolution.ResolutionConfidence = highestPriority
		resolution.ResolutionReasoning = fmt.Sprintf("Selected value based on document type priority and confidence level (score: %.2f)", highestPriority)
		
		// Assess impact on case
		if resolution.ResolutionConfidence > 0.8 {
			resolution.ImpactOnCase = "minimal_impact"
		} else if resolution.ResolutionConfidence > 0.6 {
			resolution.ImpactOnCase = "moderate_impact"
		} else {
			resolution.ImpactOnCase = "significant_impact_requires_review"
		}
	}
	
	return resolution
}

// calculateOverallConsistency calculates overall consistency score
func (cv *ConsistencyValidator) calculateOverallConsistency(result ConsistencyValidationResult) float64 {
	if len(result.FactConsistencies) == 0 {
		return 0.0
	}
	
	totalConsistency := 0.0
	for _, fact := range result.FactConsistencies {
		totalConsistency += fact.ConsistencyScore
	}
	
	averageConsistency := totalConsistency / float64(len(result.FactConsistencies))
	
	// Apply penalty for conflicts
	conflictPenalty := 0.0
	for _, conflict := range result.DetectedConflicts {
		switch conflict.ConflictSeverity {
		case "critical":
			conflictPenalty += 0.3
		case "significant":
			conflictPenalty += 0.2
		case "minor":
			conflictPenalty += 0.1
		}
	}
	
	// Normalize penalty
	if len(result.DetectedConflicts) > 0 {
		conflictPenalty /= float64(len(result.DetectedConflicts))
	}
	
	finalConsistency := averageConsistency - conflictPenalty
	if finalConsistency < 0.0 {
		finalConsistency = 0.0
	}
	
	return finalConsistency
}

// generateValidationSummary generates a summary of validation results
func (cv *ConsistencyValidator) generateValidationSummary(result ConsistencyValidationResult) ValidationSummary {
	summary := ValidationSummary{
		TotalFactsValidated:   len(result.FactConsistencies),
		ConsistentFacts:       0,
		ConflictingFacts:      0,
		CriticalConflicts:     0,
		SignificantConflicts:  0,
		MinorConflicts:        0,
		AutoResolvedConflicts: 0,
		ManualReviewRequired:  0,
		AverageReliability:    0.0,
		RecommendedActions:    []string{},
	}
	
	// Count consistent facts
	for _, fact := range result.FactConsistencies {
		if fact.ConsistencyScore > 0.8 {
			summary.ConsistentFacts++
		} else {
			summary.ConflictingFacts++
		}
	}
	
	// Count conflicts by severity
	for _, conflict := range result.DetectedConflicts {
		switch conflict.ConflictSeverity {
		case "critical":
			summary.CriticalConflicts++
		case "significant":
			summary.SignificantConflicts++
		case "minor":
			summary.MinorConflicts++
		}
	}
	
	// Count resolution status
	for _, resolution := range result.ResolutionSuggestions {
		if resolution.ResolutionConfidence > 0.8 {
			summary.AutoResolvedConflicts++
		} else {
			summary.ManualReviewRequired++
		}
	}
	
	// Calculate average reliability
	if len(result.ReliabilityScores) > 0 {
		totalReliability := 0.0
		for _, score := range result.ReliabilityScores {
			totalReliability += score
		}
		summary.AverageReliability = totalReliability / float64(len(result.ReliabilityScores))
	}
	
	// Generate recommended actions
	if summary.CriticalConflicts > 0 {
		summary.RecommendedActions = append(summary.RecommendedActions, "Immediate attention required for critical conflicts")
	}
	if summary.ManualReviewRequired > 0 {
		summary.RecommendedActions = append(summary.RecommendedActions, "Manual review required for unresolved conflicts")
	}
	if summary.AverageReliability < 0.7 {
		summary.RecommendedActions = append(summary.RecommendedActions, "Consider additional document verification")
	}
	if result.OverallConsistency > 0.8 {
		summary.RecommendedActions = append(summary.RecommendedActions, "High consistency - leverage for strong case narrative")
	}
	
	return summary
}

// GetValidationSummary returns a summary of the validator's capabilities
func (cv *ConsistencyValidator) GetValidationSummary() map[string]interface{} {
	summary := make(map[string]interface{})
	
	summary["validation_rules"] = len(cv.ValidationRules)
	summary["conflict_detection_rules"] = len(cv.ConflictDetector.DetectionRules)
	summary["document_types_supported"] = len(cv.ReliabilityScorer.SourceCredibility)
	summary["auto_resolution_enabled"] = cv.ConflictResolver.AutoResolutionEngine.ConfidenceThresholds.AutoResolutionThreshold > 0
	
	return summary
}