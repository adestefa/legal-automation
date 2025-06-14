# DEFECT 3

NAME: Step Icons Are No Longer Active
PRIORITY: HIGH
STATUS: QUEUE
GITHUB ISSUE: #8

## Issue Description
As users progress through the workflow, the step icons in the progress indicator are no longer updating to show the current active step. This regression occurred after the recent template fixes.

## Steps to Reproduce
1. Start at Step 0 (Case Setup)
2. Progress through Steps 1, 2, 3, etc.
3. Observe the step icons/indicators at the top of the page
4. Notice that step icons do not update to reflect current progress

## Expected Behavior
- Step icons should update as user progresses through workflow
- Current step should be highlighted/active
- Completed steps should show as completed
- Future steps should show as inactive/pending

## Actual Behavior
- Step icons remain static and do not update
- No visual indication of current workflow position
- Progress indicator appears broken

## Root Cause Analysis
This issue likely stems from the recent template fix where the out-of-band (OOB) progress steps update was removed from `_step_wrapper.gohtml` to resolve the 500 error. The OOB update was responsible for refreshing the progress indicator.

## Technical Details
- Related to removal of OOB swap in `_step_wrapper.gohtml`
- Progress steps update mechanism needs to be restored
- May need alternative approach to OOB if that was causing errors

## Impact
- Poor user experience - no sense of workflow progress
- Users cannot easily see where they are in the process
- Appears unprofessional and broken

## Environment
- Version: v2.5.39
- Affects all browsers
- Occurred after template wrapper simplification

## Acceptance Criteria
- [ ] Progress step icons update correctly as user advances
- [ ] Current step is visually highlighted
- [ ] Completed steps show completed state
- [ ] No 500 errors or template issues
- [ ] Works across all workflow steps