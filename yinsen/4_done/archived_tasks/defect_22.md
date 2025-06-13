# DEFECT 22

**DATE**: 2025-06-04
**TIME**: 15:46:00
**PROJ**: Mallon Legal Assistant
**STATUS**: QUEUE
**TYPE**: DEFECT
**PRIORITY**: HIGH
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Restore Missing Preview Document Tab with Yellow Highlighting

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
The Preview Document tab that shows the generated document with yellow highlighted text from source documents was lost during the Go SSR + HTMX refactor. This tab was present in the original Alpine.js version and needs to be restored in the new architecture.

**WHY**: 
The preview functionality is critical for attorneys to review exactly what content was extracted from source documents and how it was integrated into the final complaint. The yellow highlighting shows the provenance of extracted text, which is essential for legal accuracy and compliance.

**CHALLENGE**: 
Need to recreate the preview document functionality in the Go SSR + HTMX architecture. Must implement text highlighting that shows which content came from which source documents, while maintaining the performance and simplicity of the server-side rendering approach.

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Add document preview endpoint to UI handlers
2. Implement text highlighting logic in Go backend
3. Create preview document template with highlighting CSS
4. Add preview tab to Step 3 or Step 4 interface
5. Ensure highlighting data is tracked during document generation

**EVALUATION/PLANNING**:
1. Yinsen shall review objectives for Task
2. Ask questions to clarify or provide options/feedback
3. Document any blockers and ways around them
4. Think like a hacker, be creative for optimal solutions

**ACCEPTANCE CRITERIA**:
- [ ] Preview Document tab is restored and accessible
- [ ] Generated document displays with proper formatting
- [ ] Source text is highlighted in yellow with attribution
- [ ] Preview updates when document content changes
- [ ] Highlighting shows which source document provided each piece of text
- [ ] Preview functionality works seamlessly in HTMX architecture

## Execution Tracking

**STARTED**: 2025-06-04 23:45:00
**MOVED_TO_DEV**: 2025-06-04 23:46:00
**MOVED_TO_QA**: 2025-06-05 00:02:00
**COMPLETED**: {TIMESTAMP}

**BLOCKERS_ENCOUNTERED**:
- 2025-06-04 23:50:00: Missing preview route registration → Fixed by adding /ui/preview-document endpoint
- 2025-06-04 23:55:00: Template highlighting logic needed → Implemented JavaScript-based highlighting with source attribution

**LESSONS_LEARNED**:
- Go SSR + HTMX architecture requires server-side data preparation for complex highlighting
- JavaScript highlighting provides better user experience than server-side HTML manipulation
- Template tab navigation works seamlessly with HTMX fragment loading

**QA_FEEDBACK**:
- {TIMESTAMP}: {FEEDBACK_FROM_REVIEWER} → {ACTION_TAKEN}

## Technical Implementation

**FILES_MODIFIED**:
- /dev/handlers/ui_handlers.go: Added PreviewDocument handler and preview data structures
- /dev/main_v2.go: Added /ui/preview-document route registration
- /dev/templates/_step3_review_data.gohtml: Added tab navigation and Preview Document tab
- /dev/templates/_document_preview.gohtml: Created new template with yellow highlighting functionality

**TESTS_ADDED**:
- Manual testing: Preview Document tab loads successfully
- Manual testing: Yellow highlighting displays with source attribution
- Manual testing: Tab switching works seamlessly

**PERFORMANCE_IMPACT**:
- Load time: Instant tab switching with HTMX
- Memory usage: Minimal impact with server-side rendering

**SECURITY_CONSIDERATIONS**:
- Input sanitization: Template data is server-controlled, no user input vulnerability
- Authentication: Preview functionality respects existing session management

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!