# TASK 34

**DATE**: 2025-06-12
**TIME**: 14:30:00
**PROJ**: proj-mallon
**STATUS**: QUEUE
**TYPE**: ENHANCEMENT
**PRIORITY**: MEDIUM
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Replace "Restart Setup" Button with Back Navigation System

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
Replace the "Restart Setup" button in Step 0 with a proper Back button that moves users to the previous step in the workflow. Implement Back buttons on all steps except the first (Step 0) and last (Step 5) pages to provide consistent navigation throughout the multi-step legal document generation process.

**WHY**: 
The current "Restart Setup" button on Step 0 forces users to completely restart the workflow, losing all progress and selections. This creates a poor user experience when users want to make minor corrections or adjustments to previous steps. A proper Back navigation system allows users to move backward through the workflow while preserving their progress, improving usability and reducing frustration during the legal document generation process.

**CHALLENGE**: 
- Step 0 (Case Setup) currently has a "Restart Setup" button that reloads the entire page
- Need to maintain state and context when navigating backward through steps
- Must ensure proper HTMX integration for seamless navigation
- Backend handlers need to support bidirectional navigation flow
- Template data must be preserved when moving between steps
- Need to handle edge cases where certain data may not be available when going back

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Remove the "Restart Setup" button from Step 0 (`_step0_case_setup.gohtml`)
2. Audit all step templates to ensure consistent Back button implementation:
   - Step 1: Already has "Back to Case Setup" - verify functionality
   - Step 2: Add "Back to Document Selection" button
   - Step 3: Add "Back to Template Selection" button  
   - Step 4: Add "Back to Review Data" button
   - Step 5: Keep existing "Back to Document" button (last step)
3. Update backend handlers in `ui_handlers.go` to support backward navigation
4. Ensure proper state preservation when navigating backwards
5. Test all navigation paths to verify data integrity

**EVALUATION/PLANNING**:
1. Yinsen shall review current navigation implementation across all step templates
2. Identify which templates are missing Back buttons and which handlers need updates
3. Document any state management requirements for backward navigation
4. Plan testing approach to verify data preservation during navigation

**ACCEPTANCE CRITERIA**:
- [ ] Step 0 no longer has "Restart Setup" button 
- [ ] Step 1 maintains existing "Back to Case Setup" functionality
- [ ] Step 2 has "Back to Document Selection" button that navigates to Step 1
- [ ] Step 3 has "Back to Template Selection" button that navigates to Step 2
- [ ] Step 4 has "Back to Review Data" button that navigates to Step 3
- [ ] Step 5 maintains existing "Back to Document" button
- [ ] All Back buttons use consistent HTMX patterns for navigation
- [ ] User selections and data are preserved when navigating backward
- [ ] Navigation flows work seamlessly without page reloads
- [ ] UI consistency maintained across all step templates

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
- /Users/corelogic/satori-dev/clients/proj-mallon/v2/templates/_step0_case_setup.gohtml: Remove "Restart Setup" button
- /Users/corelogic/satori-dev/clients/proj-mallon/v2/templates/_step2_template_selection.gohtml: Add Back button to Step 1
- /Users/corelogic/satori-dev/clients/proj-mallon/v2/templates/_step3_review_data.gohtml: Add Back button to Step 2
- /Users/corelogic/satori-dev/clients/proj-mallon/v2/templates/_step4_generate_document.gohtml: Add Back button to Step 3
- /Users/corelogic/satori-dev/clients/proj-mallon/v2/handlers/ui_handlers.go: Verify/update handlers support backward navigation

**TESTS_ADDED**:
- Manual testing: Navigation flow from Step 5 back to Step 0
- Manual testing: Data preservation when navigating backward
- Manual testing: HTMX functionality for all Back buttons

**PERFORMANCE_IMPACT**:
- Navigation: Improved user experience with faster step transitions
- State Management: Minimal impact, leveraging existing HTMX patterns

**SECURITY_CONSIDERATIONS**:
- Navigation: No additional security concerns, using existing HTMX patterns
- Data Integrity: Ensure backward navigation doesn't expose inappropriate data states

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!