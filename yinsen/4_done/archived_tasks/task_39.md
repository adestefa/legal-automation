# TASK 39

**DATE**: 2025-06-08
**TIME**: 04:50:00
**PROJ**: Mallon Legal Assistant
**STATUS**: QUEUE
**TYPE**: TASK
**PRIORITY**: HIGH
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Remove 'Load iCloud Root Folder' Button and Replace with Back Button

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
On Step 0 (case folder setup), remove the large "Load iCloud Root Folder" button that takes up excessive room. Replace it with a smaller back button that returns to the previous step, matching the UI pattern used in other steps of the application. The changes you made are not showing the new back button and instead show we are not connected. PHOTO: /Users/corelogic/satori-dev/clients/proj-mallon/yinsen/artifacts/Step_0_setup_case_folder_bad.png Please fix this.

**WHY**: 
The current "Load iCloud Root Folder" button is too prominent and takes up unnecessary space in the interface. The functionality is not needed at this stage as users are already in the case folder selection process. A simple back button would provide better navigation and maintain consistency with other steps in the workflow.

**CHALLENGE**: 
- Ensure the back button maintains proper navigation flow
- Keep the UI balanced after removing the large button
- Maintain HTMX functionality for seamless navigation
- Ensure the change doesn't break existing folder selection functionality

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Locate and remove the "Load iCloud Root Folder" button HTML in _step0_case_setup.gohtml
2. Add a smaller back button that navigates to the appropriate previous step
3. Adjust layout and spacing to maintain visual balance
4. Test navigation flow to ensure smooth user experience

**EVALUATION/PLANNING**:
1. Yinsen shall review objectives for Task
2. Ask questions to clarify or provide options/feedback
3. Document any blockers and ways around them
4. Think like a hacker, be creative for optimal solutions

**ACCEPTANCE CRITERIA**:
- [ ] "Load iCloud Root Folder" button is removed from Step 0
- [ ] Back button is added that properly navigates to previous step
- [ ] UI remains balanced and professional looking
- [ ] Navigation flow works correctly with HTMX
- [ ] No regression in existing functionality

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

**FILES_MODIFIED**:
- dev/templates/_step0_case_setup.gohtml: Remove button and add back navigation

**TESTS_ADDED**:
- Manual testing of navigation flow

**PERFORMANCE_IMPACT**:
- None expected - UI change only

**SECURITY_CONSIDERATIONS**:
- None - UI change only

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!