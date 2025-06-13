# TASK:30

**DATE**: 2025-06-05
**TIME**: 15:55:00
**PROJ**: Mallon Legal Assistant
**STATUS**: QUEUE
**TYPE**: MODIFICATION
**PRIORITY**: MEDIUM
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Hide 'Create Case Folder' Option for Demo

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
Temporarily hide the 'Create Case Folder' option in Step 0 for the upcoming demo. This should be done by commenting out the HTML feature in the template rather than removing it completely.

**WHY**: 
The 'Create Case Folder' functionality is not fully implemented for the demo tomorrow. To avoid confusion and present a clean interface to the client, this option should be temporarily hidden while preserving the code for future implementation.

**CHALLENGE**: 
Need to identify the correct location in the templates to comment out the feature without disrupting the layout or functionality of the rest of Step 0.

## Implementation Planning

**POSSIBLE SOLUTION**:
1. Locate the 'Create Case Folder' option in the Step 0 template file
2. Comment out the HTML section with appropriate markers for future restoration
3. Ensure the layout of Step 0 remains balanced and visually appealing without this option
4. Add a comment explaining that this is temporarily hidden for the demo

**EVALUATION/PLANNING**:
1. Review the `_step0_case_setup.gohtml` template to find the 'Create Case Folder' section
2. Determine the best way to comment out the section without affecting other elements
3. Test the UI to ensure it looks correct without the commented out section
4. Document the change in comments for future developers

**ACCEPTANCE CRITERIA**:
- [ ] 'Create Case Folder' option is not visible in Step 0
- [ ] The HTML code for the feature is preserved but commented out
- [ ] Step 0 layout remains visually balanced and professional
- [ ] Comment is added explaining why the feature is hidden
- [ ] No disruption to other Step 0 functionality

## Execution Tracking

**STARTED**: 2025-06-05 23:02:00
**MOVED_TO_DEV**: 2025-06-05 23:03:00
**MOVED_TO_QA**: 2025-06-05 23:04:00
**COMPLETED**: 2025-06-05 23:05:00

**BLOCKERS_ENCOUNTERED**:
- 

**LESSONS_LEARNED**:
- 

**QA_FEEDBACK**:
- 

## Technical Implementation

**FILES_MODIFIED**:
- `/Users/corelogic/satori-dev/clients/proj-mallon/dev/templates/_step0_case_setup.gohtml`
- `/Users/corelogic/satori-dev/clients/proj-mallon/yinsen/task_list.md`
- `/Users/corelogic/satori-dev/clients/proj-mallon/yinsen/history.md`

**TESTS_ADDED**:
- 

**PERFORMANCE_IMPACT**:
- 

**SECURITY_CONSIDERATIONS**:
- 

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!