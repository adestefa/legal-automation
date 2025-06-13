# TASK 19

**DATE**: 2025-06-04
**TIME**: 14:35:00
**PROJ**: Mallon Legal Assistant
**STATUS**: QUEUE
**TYPE**: FEATURE
**PRIORITY**: MEDIUM
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Add "Select All" Option to Document Selection Page

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
Add a "select all" option to the document selection page (Step 1) that allows users to quickly select all available documents in one click, with the ability to deselect all as well.

**WHY**: 
Users with multiple documents in a case folder would benefit from being able to select all documents at once rather than clicking each checkbox individually. This improves user experience and workflow efficiency.

**CHALLENGE**: 
- Must work with both iCloud documents and backend documents
- Should integrate seamlessly with Alpine.js reactive data binding
- Need to handle the toggle state properly (select all ↔ deselect all)
- Should maintain existing individual selection functionality
- Must update the selectedDocs array correctly for downstream processing

## Implementation Planning

**DESIGN CONSIDERATIONS**:
1. **UI Component**:
   - Add a "Select All" / "Deselect All" toggle button or checkbox
   - Position it prominently in the document selection area
   - Clear visual indication of current state (all selected vs. partial/none)

2. **Functionality**:
   - Select All: Check all document checkboxes and update selectedDocs array
   - Deselect All: Uncheck all checkboxes and clear selectedDocs array
   - Smart toggle: If all are selected, button shows "Deselect All"
   - If none or partial selection, button shows "Select All"

3. **Integration Points**:
   - Alpine.js selectedDocs array management
   - Both iCloud documents and backend documents support
   - Maintain compatibility with existing individual checkbox selection

**POSSIBLE SOLUTION**:
1. Add "Select All" functionality to Alpine.js component
2. Create selectAllDocuments() and deselectAllDocuments() methods
3. Add computed property to determine current selection state
4. Update UI to include select all toggle button
5. Ensure proper array manipulation for both document sources

**EVALUATION/PLANNING**:
1. Review current document selection UI structure
2. Implement Alpine.js methods for bulk selection
3. Add UI component with appropriate styling
4. Test with both iCloud and backend documents
5. Verify downstream workflow (Step 2 → Step 3) works correctly

**ACCEPTANCE CRITERIA**:
- [ ] "Select All" button/checkbox appears in Step 1 document selection
- [ ] Clicking "Select All" selects all visible documents
- [ ] Button text changes to "Deselect All" when all are selected
- [ ] Clicking "Deselect All" clears all selections
- [ ] Individual checkboxes remain functional alongside select all
- [ ] selectedDocs array is properly updated for downstream processing
- [ ] Works with both iCloud documents and backend documents
- [ ] Visual feedback clearly shows selection state

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
- `/dev/frontend/index.html`: Add "Select All" UI component to Step 1
- `/dev/frontend/driver.js`: Add selectAllDocuments() and deselectAllDocuments() functions
- `/dev/frontend/site.css`: Add styling for select all component (if needed)

**TESTS_ADDED**:
- Manual test: Select all documents in Yousef_Eman case folder
- Manual test: Deselect all after selecting individual documents
- Manual test: Verify select all works with backend documents (when no case folder selected)

**PERFORMANCE_IMPACT**:
- Minimal - simple array manipulation and UI updates

**SECURITY_CONSIDERATIONS**:
- None - client-side UI enhancement only

---

**Confirmation Protocol**: 
Stop. Confirm you understand. I will add "Select All" functionality to the document selection page that works with both iCloud and backend documents.

**Completion Protocol**:
Thank you Yinsen, this will greatly improve user workflow efficiency!
