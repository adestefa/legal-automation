# Task 7: Clean Up Preview Document Tab for Professional Review

**Status**: 🟡 PENDING  
**Priority**: MEDIUM  
**Target Version**: v2.9.2  
**Estimated Time**: 30-45 minutes  

## Summary
Clean up the Preview Document tab in Step 3 to provide a professional, clean document preview focused purely on review. Remove technical metadata, file paths, and highlighted text annotations that are not needed for document review.

## Current Issue
The Preview Document tab currently shows:
- "Highlighted text extracted from source documents" section at the top
- File paths and document names in the preview header
- Technical metadata that distracts from document review
- Information that's already available in other tabs (Review Data, Missing Content)

## User Story
**As a lawyer**, I want the Preview Document tab to show me a clean, professional preview of the generated complaint document **so that** I can focus on reviewing the actual legal content without distractions from technical metadata.

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
   - Ensure clean, professional document formatting

### Implementation Approach:
1. **Identify metadata sources** - Find where file paths and source annotations are added
2. **Clean preview generation** - Remove technical elements from preview data structure
3. **Update template** - Ensure preview template only shows legal document content
4. **Test preview output** - Verify clean, professional appearance

## Acceptance Criteria

### ✅ Preview Document Tab Shows:
- [ ] Clean legal document content only
- [ ] Professional court document formatting
- [ ] All legal sections (jurisdiction, parties, causes, etc.)
- [ ] No file paths or source document names
- [ ] No "highlighted text" annotations
- [ ] No technical metadata

### ✅ Document Review Experience:
- [ ] Lawyer can focus on legal content
- [ ] Preview looks like a professional court filing
- [ ] Content is easy to read and review
- [ ] No technical distractions

### ✅ Other Tabs Unaffected:
- [ ] Review Data tab still shows source information
- [ ] Missing Content tab still shows document analysis
- [ ] Selected documents info available in appropriate tabs

## Testing Plan

1. **Generate document preview** through complete workflow
2. **Verify Preview tab** shows clean document only
3. **Check other tabs** still show source document information
4. **Test different case scenarios** to ensure consistent clean output
5. **Verify professional appearance** suitable for lawyer review

## Business Value
- **Improved User Experience**: Lawyers get clean, focused document review
- **Professional Appearance**: Preview looks like actual court filing
- **Reduced Distraction**: No technical metadata cluttering the view
- **Better Workflow**: Clear separation between technical info and document review

## Related Tasks
- Task 6 (Missing Content Tab) - Completed ✅ 
- Task 4 (Dynamic Template Engine) - Provides the clean document structure
- Future Task 5 (iCloud Integration) - Will use this clean preview for saved documents

## Implementation Notes
- This is a UI/UX improvement task focused on user experience
- Should not affect underlying document generation logic
- Technical metadata should remain available in appropriate tabs
- Focus on making Preview tab lawyer-friendly for document review