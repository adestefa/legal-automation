# TASK 36

**DATE**: 2025-06-12
**TIME**: 15:45:00
**PROJ**: proj-mallon
**STATUS**: QUEUE
**TYPE**: ENHANCEMENT
**PRIORITY**: HIGH
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Dynamic Document Preview Based on Selected Documents

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
Modify the document preview functionality so that it only shows documents that were actually selected by the user in Step 1, rather than showing a static hardcoded list. Currently, no matter how many documents are selected, the document preview shows only "Attorney_Notes.txt, Adverse_Action_Letter_Cap_One.pdf, Civil_Cover_Sheet.txt, Complaint_Final.docx". The system should dynamically show only the selected documents with their extracted content, and display missing content information for documents not selected.

**WHY**: 
Users need accurate feedback about which documents are actually being used in their complaint generation. The current static preview is misleading and doesn't reflect their actual selections, making it impossible to verify that the right documents are being processed. This creates confusion and reduces trust in the system's accuracy.

**CHALLENGE**: 
- Document preview currently uses hardcoded document list
- Need to track user selections from Step 1 through to Step 3 preview
- Must distinguish between selected documents with content vs missing/unselected documents
- Preview tabs need to show document names with extracted content
- Missing content tab should list unselected documents and identify missing data

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Modify document selection tracking in Step 1 to persist selected document list
2. Update Step 3 preview logic to use actual selected documents instead of hardcoded list
3. Create document-specific content extraction display
4. Implement missing content analysis showing unselected documents
5. Update preview tabs to show selected document names with their content
6. Add missing content tab showing unselected documents and missing data

**EVALUATION/PLANNING**:
1. Yinsen shall review current document selection and preview implementation
2. Identify where document selections are stored and passed between steps
3. Determine how to modify preview generation to use dynamic document list
4. Plan missing content detection for unselected documents

**ACCEPTANCE CRITERIA**:
- [ ] Document preview shows only documents selected by user in Step 1
- [ ] Each selected document displays its name and extracted content in preview
- [ ] Missing content tab shows unselected documents from legal_artifacts folder
- [ ] Missing content analysis identifies what data is unavailable due to unselected documents
- [ ] No hardcoded document list in preview functionality
- [ ] Preview accurately reflects user's actual document selections

## Execution Tracking

**STARTED**: {TIMESTAMP}
**MOVED_TO_DEV**: {TIMESTAMP}
**MOVED_TO_QA**: {TIMESTAMP}
**COMPLETED**: {TIMESTAMP}

**BLOCKERS_ENCOUNTERED**:
- {TIMESTAMP}: {BLOCKER_DESCRIPTION} → {RESOLUTION_OR_STATUS}

**LESSONS_LEARNED**:
- {INSIGHT_OR_PATTERN_DISCOVERED}

**QA_FEEDBACK**:
- {TIMESTAMP}: {FEEDBACK_FROM_REVIEWER} → {ACTION_TAKEN}

## Technical Implementation

**FILES_MODIFIED**:
- main.go: Update document selection tracking and preview generation
- templates/_step3_review_data.gohtml: Modify preview tabs for dynamic content
- handlers/ui_handlers.go: Update preview logic to use selected documents

**TESTS_ADDED**:
- Manual testing: Select different document combinations and verify preview accuracy
- Validation: Ensure missing content properly identifies unselected documents

**PERFORMANCE_IMPACT**:
- Preview Generation: More accurate but similar performance
- Document Tracking: Minimal impact with proper state management

**SECURITY_CONSIDERATIONS**:
- Document Selection: Ensure only valid documents from legal_artifacts can be selected
- Content Display: Validate document content before preview display

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!