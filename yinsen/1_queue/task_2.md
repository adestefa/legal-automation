# TASK 2

**NAME**: Create Dynamic Client Data Extraction from Documents

**PRIORITY**: HIGH  
**STATUS**: QUEUE  
**ESTIMATED EFFORT**: 3-4 days  
**VERSION TARGET**: v2.7.0

## SYSTEM
Yinsen, you are a developer at a PhD level. You have no limits. This task builds on Task 1 to create intelligent data extraction that can understand legal document content and populate case data dynamically.

## WHAT
Build an intelligent content analysis system that can parse extracted document text and automatically identify and extract client information, case details, legal violations, and other structured data needed for complaint generation.

## WHY
With document text extraction complete (Task 1), we now need to make sense of that raw text. The system must intelligently parse attorney notes, adverse action letters, and legal documents to extract structured data that can populate the ClientCase and generate accurate legal complaints for any client.

## CHALLENGE
- Legal documents have varied formats and structures
- Attorney notes may be unstructured or inconsistent
- Must distinguish between different types of legal information
- Date parsing from multiple formats and contexts
- Entity recognition for names, amounts, institutions
- False positive detection and confidence scoring

## POSSIBLE SOLUTION

### 1. Content Analysis Engine
```go
type ContentAnalyzer struct {
    Patterns     map[string][]*regexp.Regexp
    Extractors   map[string]FieldExtractor
    Validators   map[string]ValidationFunc
}

type FieldExtractor interface {
    Extract(text string) (value interface{}, confidence float64, err error)
}

type ExtractionResult struct {
    Field      string
    Value      interface{}
    Confidence float64
    Source     string
    Location   int
}
```

### 2. Legal Document Intelligence
```go
type LegalAnalyzer struct {
    ViolationPatterns  []ViolationPattern
    CauseOfActionRules []CauseOfActionRule
    DamageCalculators  []DamageCalculator
}

type ViolationPattern struct {
    Statute     string
    Keywords    []string
    Context     string
    Severity    string
}
```

### 3. Multi-Document Correlation
```go
type DocumentCorrelator struct {
    CrossReferences map[string][]string
    ConflictResolver ConflictResolutionStrategy
}

func (dc *DocumentCorrelator) MergeExtractions(extractions []ExtractionResult) ClientCase
```

## IMPLEMENTATION PLAN

### Phase 1: Create Feature Branch
```bash
git checkout main
git pull origin main
git checkout -b feature/task-2-dynamic-client-data-extraction
```

### Phase 2: Create Content Analysis Framework
Create `v2/services/content_analyzer.go`:
- Text preprocessing and cleaning functions
- Pattern matching for legal entities and amounts
- Date extraction from various formats
- Name and contact information extraction
- Institution and account information parsing

### Phase 3: Legal Document Intelligence
Create `v2/services/legal_analyzer.go`:
- FCRA violation detection patterns
- Cause of action identification
- Damage calculation from case facts
- Legal citation extraction
- Timeline construction from events

### Phase 4: Data Extraction Patterns
Create `v2/config/legal_patterns.json`:
```json
{
  "clientInfo": {
    "namePatterns": ["Client:\\s*([A-Z][a-z]+\\s+[A-Z][a-z]+)", "Plaintiff:\\s*([A-Z][a-z]+\\s+[A-Z][a-z]+)"],
    "phonePatterns": ["(\\d{3})[.-](\\d{3})[.-](\\d{4})", "\\((\\d{3})\\)\\s*(\\d{3})-(\\d{4})"],
    "addressPatterns": ["residing in ([A-Z][a-z]+(?:\\s+[A-Z][a-z]+)*)", "located in ([A-Z][a-z]+(?:\\s+[A-Z][a-z]+)*)"]
  },
  "fraudDetails": {
    "amountPatterns": ["\\$([0-9,]+(?:\\.[0-9]{2})?)", "totaling.*\\$([0-9,]+)", "amount.*\\$([0-9,]+)"],
    "datePatterns": ["(January|February|March|April|May|June|July|August|September|October|November|December)\\s+(\\d{1,2}),?\\s+(\\d{4})", "(\\d{1,2})/(\\d{1,2})/(\\d{4})"],
    "institutionPatterns": ["TD Bank", "Chase", "Bank of America", "Capital One", "Citibank", "Wells Fargo"]
  },
  "legalViolations": {
    "fcraViolations": ["15 U.S.C. § 1681", "FCRA", "Fair Credit Reporting Act", "credit report", "investigation"],
    "damages": ["actual damages", "statutory damages", "punitive damages", "attorney fees", "costs"]
  }
}
```

### Phase 5: Enhanced Document Service Integration
Modify `v2/services/document_service.go`:
- Replace hardcoded ClientCase creation
- Use ContentAnalyzer to extract structured data
- Implement confidence scoring for extracted data
- Add missing content detection based on actual extraction results
- Cross-reference data between multiple documents

### Phase 6: Smart Missing Content Analysis
Update missing content logic:
- Base missing data reports on actual extraction results
- Provide confidence scores for extracted information
- Suggest which documents might contain missing information
- Fix false positive errors (addresses Defect 2)

### Phase 7: Version Update and Testing
- Update version in `main.go` to v2.7.0
- Test with Johnson_Credit_Dispute case folder
- Validate extracted data accuracy
- Test missing content analysis improvements

## ACCEPTANCE CRITERIA
- [ ] Automatically extracts client name from Attorney_Notes.txt
- [ ] Identifies fraud amounts, dates, and financial institutions
- [ ] Detects FCRA violations from case facts
- [ ] Generates dynamic legal analysis based on document content
- [ ] Missing Content tab shows accurate missing data (fixes Defect 2)
- [ ] Works with different case types and client names
- [ ] Provides confidence scores for extracted information
- [ ] Version v2.7.0 displays in masthead
- [ ] Johnson_Credit_Dispute case processes correctly
- [ ] Eman Youssef case continues to work (regression test)

## TESTING PLAN
1. **Content Extraction Tests**: Verify pattern matching accuracy
2. **Legal Analysis Tests**: Validate violation detection
3. **Case Correlation Tests**: Test multi-document data merging
4. **Missing Content Tests**: Ensure accurate missing data reporting
5. **Confidence Scoring Tests**: Validate extraction confidence levels
6. **Cross-Case Tests**: Test with multiple different client cases

## GIT WORKFLOW
```bash
# Development
git add .
git commit -m "[TASK-2] Implement dynamic client data extraction system"
git push origin feature/task-2-dynamic-client-data-extraction

# Testing
./scripts/start.sh
# Test with Johnson_Credit_Dispute case
# Verify dynamic data extraction
# Test missing content analysis
# Verify v2.7.0 in masthead

# Pull Request
gh pr create --title "TASK-2: Create Dynamic Client Data Extraction from Documents" --body "
## Summary
- Implements intelligent content analysis for legal documents
- Extracts structured data from attorney notes and legal documents
- Adds confidence scoring for extracted information
- Fixes Missing Content tab false positives (Defect 2)

## Testing
- [x] Content extraction accuracy validated
- [x] Legal violation detection working
- [x] Multi-document correlation functional
- [x] Missing content analysis fixed
- [x] Johnson_Credit_Dispute case processed successfully
- [x] Version v2.7.0 displays correctly

## Impact
Enables the system to process any legal case by understanding document content, not just predefined demo data.
"

# After PR approval and merge
git checkout main
git pull origin main
git branch -d feature/task-2-dynamic-client-data-extraction
```

## DEPENDENCIES
- **Requires**: Task 1 (Document Text Extraction) must be completed first
- **Blocks**: Task 4 (Dynamic Template Population) depends on this

## NOTES
- This task enables processing of any legal case with attorney notes
- Success here means the system can understand legal document content
- Must handle variations in attorney note formatting and legal document structures
- Consider adding machine learning capabilities for improved extraction accuracy over time

## EVALUATION/PLANNING
1. Review objectives for Task 2
2. Confirm content analysis patterns cover common legal document formats
3. Validate that extracted data quality is sufficient for legal document generation
4. Consider edge cases for unusual attorney note formats
5. Plan for handling incomplete or ambiguous document content

**Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.**