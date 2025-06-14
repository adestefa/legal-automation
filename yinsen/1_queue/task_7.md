# Task 7: Clean Up Preview Document Tab for Professional Review

**Status**: 🟡 PENDING  
**Priority**: MEDIUM  
**Target Version**: v2.9.2  
**Estimated Time**: 30-45 minutes  

## Summary
Clean up the Preview Document tab in Step 3 to provide a professional, clean document preview focused purely on review. Remove technical metadata and file paths, but ADD BACK yellow highlighting on dynamically inserted text to help lawyers quickly identify what data was extracted and added by the system.

## Current Issue
The Preview Document tab currently shows:
- "Highlighted text extracted from source documents" section at the top
- File paths and document names in the preview header
- Technical metadata that distracts from document review
- Information that's already available in other tabs (Review Data, Missing Content)
- NO yellow highlighting on dynamically inserted text (this should be restored)

## User Story
**As a lawyer**, I want the Preview Document tab to show me a clean, professional preview of the generated complaint document with yellow highlighting on dynamically inserted text **so that** I can focus on reviewing the actual legal content and quickly identify which information was extracted and added by the system.

## Requirements

### Remove These Elements:
1. **"Highlighted text extracted from source documents"** section
2. **File paths** at the top of the document preview  
3. **Source document names** in the preview header
4. **Any technical metadata** that appears above the actual document content

### Keep These Elements:
1. **Actual complaint document content** (all legal sections)
2. **Professional document formatting** (headers, paragraphs, structure)
3. **Legal content** (jurisdiction, parties, causes of action, etc.)
4. **Document navigation** within the preview

### ADD These Elements:
1. **Yellow highlighting** on dynamically inserted text from document extraction
2. **Visual indication** of which data was added by the system vs. template text
3. **Reference highlighting** to help lawyers verify extracted information

## Technical Implementation

### Files to Modify:
1. **`handlers/ui_handlers.go`** - `PreviewDocument` function
   - Remove source document metadata from preview data
   - Clean up document header generation

2. **`templates/_document_preview.gohtml`** - Preview template
   - Remove highlighted text sections
   - Remove file path displays
   - Focus on clean document presentation

3. **`services/document_formatter.go`** - Preview formatting
   - Remove metadata annotations from preview output
   - ADD yellow highlighting for dynamically inserted text
   - Ensure clean, professional document formatting

### Implementation Approach:
1. **Identify metadata sources** - Find where file paths and source annotations are added
2. **Clean preview generation** - Remove technical elements from preview data structure
3. **Add highlighting logic** - Implement yellow highlighting for extracted/inserted text
4. **Update template** - Ensure preview template shows clean content with highlighting
5. **Test preview output** - Verify clean, professional appearance with useful highlighting

## Acceptance Criteria

### ✅ Preview Document Tab Shows:
- [ ] Clean legal document content only
- [ ] Professional court document formatting
- [ ] All legal sections (jurisdiction, parties, causes, etc.)
- [ ] No file paths or source document names
- [ ] No "highlighted text extracted from source documents" section
- [ ] No technical metadata
- [ ] **Yellow highlighting on dynamically inserted text** (NEW REQUIREMENT)

### ✅ Document Review Experience:
- [ ] Lawyer can focus on legal content
- [ ] Preview looks like a professional court filing
- [ ] Content is easy to read and review
- [ ] No technical distractions
- [ ] **Highlighted text helps lawyer quickly identify extracted data**
- [ ] **Clear visual distinction between template text and inserted information**

### ✅ Other Tabs Unaffected:
- [ ] Review Data tab still shows source information
- [ ] Missing Content tab still shows document analysis
- [ ] Selected documents info available in appropriate tabs

## Testing Plan

1. **Generate document preview** through complete workflow
2. **Verify Preview tab** shows clean document without metadata sections
3. **Check yellow highlighting** appears on dynamically inserted text
4. **Test highlighting accuracy** - ensure only extracted data is highlighted
5. **Check other tabs** still show source document information
6. **Test different case scenarios** to ensure consistent clean output with proper highlighting
7. **Verify professional appearance** suitable for lawyer review

## Business Value
- **Improved User Experience**: Lawyers get clean, focused document review
- **Professional Appearance**: Preview looks like actual court filing
- **Reduced Distraction**: No technical metadata cluttering the view
- **Better Workflow**: Clear separation between technical info and document review
- **Enhanced Verification**: Yellow highlighting helps lawyers quickly verify extracted data
- **Quality Assurance**: Visual indication of system-inserted vs. template content

## Related Tasks
- Task 6 (Missing Content Tab) - Completed ✅ 
- Task 4 (Dynamic Template Engine) - Provides the clean document structure
- Future Task 5 (iCloud Integration) - Will use this clean preview for saved documents

## Implementation Notes
- This is a UI/UX improvement task focused on user experience
- Should not affect underlying document generation logic
- Technical metadata should remain available in appropriate tabs
- Focus on making Preview tab lawyer-friendly for document review
- **Yellow highlighting should ONLY appear in Preview tab, not in final generated documents**
- Highlighting helps lawyers verify extracted data during review process
- Need to distinguish between template text and dynamically inserted content