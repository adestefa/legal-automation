# Defect 1a: Session Infrastructure for Navigation State

## Overview
Add session management infrastructure to Go handlers to preserve user selections and form data across navigation events.

## Problem Statement
Currently, when users navigate between steps using the breadcrumb navigation, all form data and selections are lost. This sub-defect focuses on implementing the backend session infrastructure needed to store and retrieve user state.

## Scope
- Backend session state management
- Handler updates to save/retrieve session data
- Session storage implementation

## Technical Requirements
1. **Session Management**
   - Implement session middleware in Go
   - Create session storage mechanism (in-memory or file-based)
   - Add session initialization to main.go

2. **Handler Updates**
   - Update Step 1 handler to save:
     - Selected case folder
     - Case details entered
   - Update Step 2 handler to save:
     - Selected template
     - Any template-specific data
   - Update Step 3 handler to save:
     - Form field values
     - Document generation state

3. **Session Data Structure**
   ```go
   type SessionData struct {
       CaseFolder    string
       CaseDetails   map[string]string
       TemplateName  string
       TemplateData  map[string]interface{}
       FormData      map[string]string
       CurrentStep   int
   }
   ```

## Implementation Tasks
- [ ] Add session package/middleware
- [ ] Create session storage interface
- [ ] Update ui_handlers.go to save session data
- [ ] Add session retrieval logic to handlers
- [ ] Test session persistence across requests

## Success Criteria
- Session data persists across HTTP requests
- Handlers can save and retrieve session data
- Session cleanup/expiration is handled
- No memory leaks from session storage

## Dependencies
- None (foundational infrastructure)

## Testing Plan
1. Unit tests for session storage
2. Handler tests with session mocking
3. Integration test for session persistence

## Estimated Time
2-3 hours

## Priority
High - This is the foundation for fixing navigation state loss

## Related Files
- dev/main_v2.go
- dev/handlers/ui_handlers.go
- dev/services/session_service.go (to be created)