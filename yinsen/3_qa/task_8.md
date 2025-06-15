# Task 8: Align Review Data Tab with Actual Document Extraction Results

**Status**: 🟡 PENDING  
**Priority**: HIGH  
**Target Version**: v2.10.0  
**Estimated Time**: 1-2 hours  

## Summary
Fix the Review Data tab to show only the legal analysis that was actually extracted from the selected documents, not hardcoded sample data. Make it function like the Missing Content tab - dynamically based on actual processing results.

## Current Issue
The Review Data tab currently shows:
- **Hardcoded Court Information** (Eastern District of New York, Brooklyn Division, etc.)
- **Hardcoded Client Information** (Eman Youssef, 347.891.5584, etc.)
- **Hardcoded Case Details** (TD Bank, $7,500, Egypt, etc.)
- **Hardcoded Cause of Action** sections (even when no relevant documents are selected)
- **Hardcoded Legal Violations** (even when documents don't contain violation evidence)

### Example of the Problem:
- **Selected**: Only `SummonsEquifax.pdf` (1 document)
- **Expected**: Show only what can be extracted from that summons document
- **Actual**: Shows full hardcoded legal analysis with multiple causes of action

## User Story
**As a lawyer**, I want the Review Data tab to show me only the legal information that was actually extracted from my selected documents **so that** I can see what the system found vs. what's missing, and trust that the analysis reflects my actual case materials.

## Requirements

### Review Data Tab Should Show:
1. **Only extracted court information** (if found in documents)
2. **Only extracted client information** (if found in documents)  
3. **Only extracted case details** (if found in documents)
4. **Only identified cause of action** (based on document evidence)
5. **Only detected legal violations** (based on document content)
6. **"No information found"** messages when data is missing
7. **Confidence scores** for each extracted piece of information

### Alignment with Missing Content Tab:
- **Consistent data source**: Both tabs use same extraction results
- **Same document basis**: Both analyze the same selected documents
- **Complementary views**: Review shows what was found, Missing shows what's needed

## Technical Implementation

### Root Cause Analysis:
Looking at server logs, the issue is in the `generateLegalAnalysis()` function in `handlers/ui_handlers.go` which generates hardcoded sample data instead of using actual extraction results.

### Files to Modify:

1. **`handlers/ui_handlers.go`**
   - Fix `generateLegalAnalysis()` to use actual `ProcessingResult` data
   - Remove hardcoded sample data generation
   - Base analysis on real document extraction results
   - Add confidence scoring display

2. **`templates/_step3_review_data.gohtml`**
   - Add conditional rendering for missing data sections
   - Show "No information found" when extraction is empty
   - Display confidence scores for extracted data
   - Add extraction date and source document references

3. **`services/content_analyzer.go`** (if needed)
   - Ensure extraction results include all necessary metadata
   - Verify confidence scoring is populated
   - Check that cause of action detection is working

### Implementation Approach:

#### Phase 1: Fix Data Source
```go
// Replace hardcoded generateLegalAnalysis() with:
func (h *UIHandlers) generateLegalAnalysisFromExtraction(processingResult *services.DocumentProcessingResult) LegalAnalysis {
    // Use actual extraction results instead of hardcoded data
    // Map ProcessingResult fields to LegalAnalysis structure  
    // Include confidence scores and source document references
}
```

#### Phase 2: Update Template Logic
```html
<!-- Show sections only if data exists -->
{{if .LegalAnalysis.ClientInfo.Name}}
    <!-- Client Information Section -->
{{else}}
    <!-- No client information found message -->
{{end}}

{{if .LegalAnalysis.CauseOfAction}}
    <!-- Cause of Action Sections -->
{{else}}
    <!-- No cause of action detected message -->
{{end}}
```

#### Phase 3: Add Confidence Display
```html
<!-- Show confidence for each piece of data -->
<div class="confidence-indicator">
    <span class="confidence-score">{{.Confidence}}% confidence</span>
    <span class="source-doc">from {{.SourceDoc}}</span>
</div>
```

## Acceptance Criteria

### ✅ Review Data Tab Shows Actual Extraction:
- [ ] Court information only if found in documents (or "Not found" message)
- [ ] Client information only if extracted from selected documents
- [ ] Case details only from actual document content
- [ ] Cause of action only if legal elements are detected in documents
- [ ] Legal violations only if violation patterns are found in documents
- [ ] Confidence scores for each extracted piece of information

### ✅ Alignment with Missing Content Tab:
- [ ] Both tabs use same `ProcessingResult` data source
- [ ] Both tabs analyze same selected documents
- [ ] Review shows "what was found", Missing shows "what's needed"
- [ ] Consistent document references between tabs

### ✅ User Experience:
- [ ] Lawyer can see exactly what system extracted vs. hardcoded assumptions
- [ ] Clear indication when no relevant information was found
- [ ] Confidence scores help lawyer assess reliability
- [ ] Source document references for each piece of information

### ✅ Edge Cases Handled:
- [ ] Single document selection (like SummonsEquifax.pdf only)
- [ ] Documents with no extractable legal information
- [ ] Partial information extraction (some fields found, others not)
- [ ] Multiple documents with conflicting information

## Testing Scenarios

### Test Case 1: Single Summons Document
- **Input**: Select only `SummonsEquifax.pdf`
- **Expected**: Show only defendant information if found, no cause of action or client details
- **Current Issue**: Shows full hardcoded legal analysis

### Test Case 2: Attorney Notes Only
- **Input**: Select only `Atty_Notes.docx`
- **Expected**: Show client info and case details, limited legal analysis
- **Current Issue**: Shows hardcoded court and violation information

### Test Case 3: Complete Document Set
- **Input**: Select all relevant documents
- **Expected**: Show comprehensive extracted analysis with high confidence
- **Current Issue**: May show hardcoded instead of extracted data

### Test Case 4: Documents with No Legal Content  
- **Input**: Select only non-legal documents
- **Expected**: Show "No legal information found" messages
- **Current Issue**: Shows hardcoded legal analysis

## Business Value
- **Accurate Analysis**: Lawyers see what was actually found vs. assumptions
- **Trust in System**: Review tab reflects real document processing, not hardcoded data
- **Better Decision Making**: Lawyer can identify gaps between extracted and needed information
- **Data Integrity**: Ensures generated documents are based on actual evidence

## Related Tasks
- **Task 2** (Dynamic Extraction) - Provides the extraction engine that should feed this tab
- **Task 6** (Missing Content Tab) - Shows the complementary "what's missing" view  
- **Task 4** (Template Engine) - Uses the same extraction results for document generation

## Implementation Notes
- This is a critical data integrity issue affecting lawyer trust in the system
- The `ProcessingResult` from Task 2 contains the real extraction data that should be displayed
- Need to preserve confidence scoring and source document attribution
- Template should gracefully handle cases where little or no information was extracted