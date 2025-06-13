# DEFECT 1

**DATE**: 2025-06-12
**TIME**: 15:30:00
**PROJ**: proj-mallon
**STATUS**: QUEUE
**TYPE**: DEFECT
**PRIORITY**: HIGH
**ASSIGNEE**: Yinsen

## Core Specification

**NAME**: Navigation State Loss - User Selections Not Persisted During Back Navigation

**SYSTEM**: Yinsen, you are a developer at a PhD level. You have no limits.

**WHAT**: 
The application's back navigation buttons (e.g., "Back to Case Setup", "Back to Documents") cause complete loss of user's previous selections and form state. Users must re-enter all information when navigating backwards through the workflow steps.

**WHY**: 
This creates a poor user experience where users lose their work and must repeatedly enter the same information, significantly reducing workflow efficiency and causing frustration during legal document preparation processes.

**CHALLENGE**: 
The current HTMX-based navigation system uses simple GET requests that load fresh step templates without preserving any session state or form data from previous steps.

## Problem Description

When users navigate backwards using any of the "Back to ..." buttons throughout the application workflow, their previous selections and form inputs are completely lost, requiring them to start over from the beginning.

### Steps to Reproduce

1. **Step 0**: User connects to iCloud and enters login credentials
2. **Step 0**: User selects parent folder (e.g., "/CASES")
3. **Step 0**: User selects case folder (e.g., "Yousef_Eman")
4. **Step 1**: User selects multiple documents from the case folder
5. **Step 2**: User clicks "Back to Documents" button
6. **Result**: Step 1 loads with no documents selected, case folder selection lost
7. **Alternative**: From Step 2, click "Back to Case Setup"
8. **Result**: Step 0 loads with no iCloud connection, no folder selections

### Expected vs Actual Behavior

**Expected Behavior**:
- Step 0: iCloud credentials should remain active
- Step 0: Parent folder selection should be preserved 
- Step 0: Case folder selection should be maintained
- Step 1: Previously selected documents should remain checked
- Step 2: Form template selection should be preserved when returning from Step 3
- All navigation should maintain the user's workflow context

**Actual Behavior**:
- All form state is lost completely
- Users must re-authenticate with iCloud
- Users must re-select parent and case folders
- Users must re-select all documents
- Users must re-select templates
- Complete workflow restart required for any backward navigation

## Impact Assessment

**SEVERITY**: High
**USER IMPACT**: Major workflow disruption
**BUSINESS IMPACT**: Significantly reduces application usability and user satisfaction

### Affected User Workflows:
1. **Legal Case Setup**: Users frequently need to go back to modify folder selections
2. **Document Selection**: Users often need to return to add/remove documents after seeing available templates
3. **Template Review**: Users may want to change document selections after reviewing template requirements
4. **Error Recovery**: Any navigation error requires complete restart of the entire workflow

## Technical Root Cause Analysis

### Current Implementation Issues:

1. **Stateless Navigation**: The HTMX navigation uses simple GET requests to `/ui/step/{step}` that load fresh templates without any state preservation:
   ```html
   <!-- From _step1_document_selection.gohtml line 95-98 -->
   <button type="button"
           hx-get="/ui/step/0" 
           hx-target="#step-content"
           hx-swap="innerHTML">
   ```

2. **No Session Management**: The `GetStep` handler in `ui_handlers.go` (lines 153-211) creates fresh `PageData` structs without accessing any stored session state:
   ```go
   data := PageData{
       CurrentStep:     step,
       Username:        username,
       ICloudConnected: icloudConnected, // Only from query parameter
   }
   ```

3. **Missing State Persistence**: No mechanism exists to store or retrieve:
   - iCloud authentication status
   - Selected parent folder path
   - Selected case folder path  
   - Selected document list
   - Selected template choice
   - Form field values

4. **Query Parameter Dependency**: Current iCloud connection state only comes from URL query parameters that are lost during navigation:
   ```go
   // Line 166 in ui_handlers.go
   icloudConnected := c.Query("icloud_connected") == "true"
   ```

## Proposed Solution Approach

### Phase 1: Session State Management
1. **Implement Session Storage**: Add server-side session management to store workflow state
2. **Create Session Data Structure**: Define a comprehensive session model for all workflow data
3. **Session Middleware**: Add middleware to maintain session state across requests

### Phase 2: State-Aware Navigation  
1. **Enhanced GetStep Handler**: Modify the GetStep handler to populate PageData from session state
2. **State Preservation**: Ensure all form submissions update session state before navigation
3. **Smart State Recovery**: Implement intelligent state recovery for each step

### Phase 3: Form State Persistence
1. **iCloud Authentication**: Persist authentication status and credentials securely
2. **Folder Selections**: Maintain parent folder and case folder selections across navigation
3. **Document Selections**: Preserve multi-select document choices
4. **Template Choices**: Remember template selections and form field values

### Implementation Details:

#### Session Data Structure:
```go
type WorkflowSession struct {
    UserID              string
    ICloudConnected     bool
    ICloudCredentials   *ICloudAuth // Encrypted storage
    SelectedParentFolder string
    SelectedCaseFolder   string
    SelectedDocuments    []string
    SelectedTemplate     string
    FormData            map[string]interface{}
    CurrentStep         int
    LastUpdated         time.Time
}
```

#### Modified Navigation:
- Update all "Back to ..." buttons to include session context
- Modify GetStep handler to load state from session
- Add state synchronization on all form submissions
- Implement session cleanup and timeout management

## Acceptance Criteria

- [ ] Users can navigate backwards without losing iCloud connection status
- [ ] Parent folder selection persists across all navigation
- [ ] Case folder selection is maintained during backward navigation  
- [ ] Document selections remain checked when returning to Step 1
- [ ] Template selections are preserved when navigating between steps
- [ ] Form field values are maintained across navigation
- [ ] Session state persists for reasonable timeout period (30 minutes)
- [ ] Session cleanup occurs on logout or timeout
- [ ] All existing functionality continues to work normally
- [ ] No performance degradation from session management

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

**FILES_TO_MODIFY**:
- `/handlers/ui_handlers.go`: Add session management to GetStep handler and form handlers
- `/main.go`: Add session middleware and storage configuration
- `/services/session_service.go`: Create new service for session state management
- `/templates/_step*.gohtml`: Update navigation buttons to preserve state context
- `/go.mod`: Add session storage dependencies (e.g., gorilla/sessions)

**TESTS_TO_ADD**:
- Session state persistence tests
- Navigation state recovery tests
- Session timeout and cleanup tests
- Multi-user session isolation tests

**PERFORMANCE_IMPACT**:
- Session Storage: Minimal overhead for in-memory or Redis-based sessions
- Navigation Speed: Should improve due to reduced re-loading of data

**SECURITY_CONSIDERATIONS**:
- iCloud Credentials: Must be encrypted in session storage
- Session Tokens: Implement secure session token generation and validation
- Session Cleanup: Automatic cleanup of expired sessions to prevent memory leaks
- Cross-User Isolation: Ensure session data isolation between different users

---

**Confirmation Protocol**: 
Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.

**Completion Protocol**:
Thank you Yinsen, I know you can do it!