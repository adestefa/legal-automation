# TASK:37

**DATE**: 2025-06-06
**TIME**: 01:15:00
**PROJ**: Mallon Legal Assistant
**STATUS**: QUEUE
**TYPE**: ENHANCEMENT
**PRIORITY**: HIGH
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Match Professional Legal Document Format and Fix Edit Mode Formatting

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
1. Match the final professional legal document format to the reference screenshot in `/Users/corelogic/satori-dev/clients/proj-mallon/yinsen/artifacts/Screenshot_Document_Format_To_Match.png`
2. Fix the edit mode formatting issue where sections lose their line breaks and become run-on paragraphs when clicking the edit button
3. Implement consistent formatting between preview and edit modes

**WHY**: 
Attorneys expect legal documents to follow strict professional formatting standards. The current implementation needs refinement to match industry standards and provide a seamless experience between preview and editing. The current edit mode destroys formatting, making it difficult for attorneys to maintain the professional appearance of legal documents.

**CHALLENGE**: 
Maintaining consistent formatting between preview and edit modes while implementing the specific legal document formatting standards shown in the reference screenshot. Need to ensure proper preservation of paragraph breaks, section formatting, and specialized legal document elements across both modes.

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Analyze the reference screenshot to identify all required formatting elements
2. Update both document preview and editor templates to match the reference format
3. Implement consistent formatting handlers for both modes
4. Fix the edit mode issue where sections become run-on paragraphs
5. Add specialized legal document elements (lines, headers, etc.) according to the reference
6. Test across both preview and edit modes to ensure consistency

**EVALUATION/PLANNING**:
1. Compare current document formatting with reference screenshot
2. Identify specific formatting differences and required changes
3. Determine how to preserve line breaks and formatting in edit mode
4. Plan changes to CSS and JavaScript that will maintain consistency
5. Implement improvements to match professional legal document standards
6. Test thoroughly in both preview and edit modes

**ACCEPTANCE CRITERIA**:
- [ ] Document formatting precisely matches the reference screenshot
- [ ] Edit mode preserves all formatting, including line breaks and paragraph structure
- [ ] Consistent appearance between preview and edit modes
- [ ] Professional legal document elements (headers, lines, spacing) properly implemented
- [ ] All existing functionality (highlighting, saving, etc.) remains intact
- [ ] Version number incremented appropriately

## Execution Tracking

**STARTED**: 
**MOVED_TO_DEV**: 
**MOVED_TO_QA**: 
**COMPLETED**: 

**BLOCKERS_ENCOUNTERED**:
- 

**LESSONS_LEARNED**:
- 

**QA_FEEDBACK**:
- 

## Technical Implementation

**FILES_TO_MODIFY**:
- `/Users/corelogic/satori-dev/clients/proj-mallon/dev/templates/_document_preview.gohtml`
- `/Users/corelogic/satori-dev/clients/proj-mallon/dev/templates/_document_editor.gohtml`
- `/Users/corelogic/satori-dev/clients/proj-mallon/dev/templates/index.gohtml` (for version update)

**TESTS_NEEDED**:
- Preview mode formatting verification
- Edit mode formatting preservation
- Transition between modes without losing formatting
- Print and download format verification

**PERFORMANCE_IMPACT**:
- Minimal - primarily CSS and HTML changes with minor JavaScript enhancements

**SECURITY_CONSIDERATIONS**:
- Ensure no XSS vulnerabilities from content editable regions

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!