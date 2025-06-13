# TASK 25

**DATE**: 2025-06-05
**TIME**: 13:20:00
**PROJ**: Mallon Legal Assistant
**STATUS**: DEV
**TYPE**: TASK
**PRIORITY**: HIGH
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Update Step 3 Preview Document Tab to Match Screenshot with Legal Document and Yellow Highlights

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
Update the Step 3 Preview Document tab to display a complete legal document with yellow highlighted text (not a summary) that matches the reference screenshot exactly. This must show the actual legal complaint document with yellow highlights indicating text extracted from source documents, providing attorneys with a full document preview before final generation.

**UI REFERENCE**:
Screenshot: /Users/corelogic/satori-dev/clients/proj-mallon/yinsen/artifacts/Step_3_preview_document_w_highlights_screen_ui.png

**WHY**: 
Attorneys need to see the complete legal document with highlighted source text before generating the final version. This preview allows them to verify document accuracy, formatting, and content extraction before proceeding to the editing phase (TASK:26). The current implementation shows document sections but needs to display the full legal document format with proper highlighting.

**CHALLENGE**: 
Need to implement a complete legal document template that displays the full FCRA complaint with proper legal formatting, yellow highlighting for extracted data, and source attribution. Must maintain the Go SSR + HTMX architecture while providing a document that matches legal standards and attorney expectations.

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Enhance document preview template to show complete legal document
2. Implement proper legal document formatting (Caption, Jurisdiction, Facts, etc.)
3. Add comprehensive yellow highlighting system for extracted data
4. Ensure document displays proper legal structure and language
5. Add source attribution tooltips for highlighted text
6. Match screenshot layout and styling exactly

**EVALUATION/PLANNING**:
1. Yinsen shall review objectives for Task
2. Ask questions to clarify or provide options/feedback  
3. Document any blockers and ways around them
4. Think like a hacker, be creative for optimal solutions

**ACCEPTANCE CRITERIA**:
- [ ] Step 3 Preview Document tab shows complete legal document (not summary)
- [ ] Document includes proper legal formatting and structure
- [ ] Yellow highlighting displays extracted text from source documents
- [ ] Document matches screenshot layout and styling exactly
- [ ] Source attribution available via tooltips or badges
- [ ] Legal document language and format suitable for attorney review
- [ ] HTMX integration maintains seamless tab switching
- [ ] Document ready for editing functionality (foundation for TASK:26)

## Execution Tracking

**STARTED**: 2025-06-05 13:20:00
**MOVED_TO_DEV**: 2025-06-05 13:20:00
**MOVED_TO_QA**: {TIMESTAMP}
**COMPLETED**: {TIMESTAMP}

**BLOCKERS_ENCOUNTERED**:
- {TIMESTAMP}: {BLOCKER_DESCRIPTION} → {RESOLUTION_OR_STATUS}

**LESSONS_LEARNED**:
- {INSIGHT_OR_PATTERN_DISCOVERED}

**QA_FEEDBACK**:
- {TIMESTAMP}: {FEEDBACK_FROM_REVIEWER} → {ACTION_TAKEN}

## Technical Implementation

**FILES_TO_MODIFY**:
- templates/_document_preview.gohtml: Complete rewrite for legal document format
- handlers/ui_handlers.go: Enhance preview document generation with full legal content
- CSS/Styling: Yellow highlighting system for extracted text

**TESTS_ADDED**:
- {TEST_FILE}: {TEST_DESCRIPTION}

**PERFORMANCE_IMPACT**:
- {METRIC}: {BEFORE} → {AFTER}

**SECURITY_CONSIDERATIONS**:
- {SECURITY_ASPECT}: {MITIGATION_OR_VALIDATION}

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!