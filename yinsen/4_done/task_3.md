# TASK 3

**NAME**: Fix Session Persistence on Refresh

**PRIORITY**: HIGH  
**STATUS**: QUEUE  
**ESTIMATED EFFORT**: 1-2 days  
**VERSION TARGET**: v2.8.0

## SYSTEM
Yinsen, you are a developer at a PhD level. You have no limits. This task is critical for user experience - lawyers cannot lose their work when refreshing the browser or navigating away from the application.

## WHAT
Replace the current in-memory session storage with persistent session management that survives browser refresh, server restarts, and navigation events. Ensure all workflow state (selected documents, extracted data, processing results) is preserved and restored reliably.

## WHY
Currently, any browser refresh or navigation away from the application loses all user progress. This is unacceptable for lawyers who may spend significant time selecting documents and reviewing extracted data. The system must be robust enough for production use where work is never lost.

## CHALLENGE
- Current sessions are stored in memory with TTL cleanup
- No persistence layer exists for workflow state
- Complex state objects (ProcessingResult, ClientCase) need serialization
- Must handle concurrent access and session conflicts
- Need to balance performance with persistence reliability
- Must maintain session security and prevent data leaks

## POSSIBLE SOLUTION

### 1. File-Based Session Storage
```go
type PersistentSessionService struct {
    SessionDir    string
    BackupDir     string
    CleanupTicker *time.Ticker
    FileLock      sync.RWMutex
}

type SessionData struct {
    SessionID       string                 `json:"sessionId"`
    WorkflowState   *WorkflowState        `json:"workflowState"`
    ProcessingResult *DocumentProcessingResult `json:"processingResult"`
    ClientCase      *ClientCase           `json:"clientCase"`
    CreatedAt       time.Time             `json:"createdAt"`
    LastAccessed    time.Time             `json:"lastAccessed"`
    ExpiresAt       time.Time             `json:"expiresAt"`
}
```

### 2. Session Restoration Middleware
```go
func SessionRestorationMiddleware(sessionService *PersistentSessionService) gin.HandlerFunc {
    return func(c *gin.Context) {
        sessionToken := getSessionToken(c)
        sessionData := sessionService.RestoreSession(sessionToken)
        
        if sessionData != nil {
            c.Set("restoredSession", true)
            c.Set("workflowState", sessionData.WorkflowState)
            c.Set("processingResult", sessionData.ProcessingResult)
            c.Set("clientCase", sessionData.ClientCase)
        }
        
        c.Next()
    }
}
```

### 3. Automatic State Persistence
```go
type StatePersister struct {
    SessionService *PersistentSessionService
    SaveInterval   time.Duration
}

func (sp *StatePersister) AutoSave(sessionID string, state interface{}) {
    go func() {
        sp.SessionService.SaveState(sessionID, state)
    }()
}
```

## IMPLEMENTATION PLAN

### Phase 1: Create Feature Branch
```bash
git checkout main
git pull origin main
git checkout -b feature/task-3-session-persistence
```

### Phase 2: Create Persistent Session Service
Create `v2/services/persistent_session_service.go`:
- File-based session storage implementation
- JSON serialization for complex state objects
- Atomic file operations for data safety
- Session cleanup and expiration handling
- Backup and recovery mechanisms

### Phase 3: Session Data Structures
Create `v2/models/session_models.go`:
- Enhanced SessionData structure
- JSON marshaling/unmarshaling for all state objects
- Version compatibility handling
- Data validation and corruption detection

### Phase 4: Update Main Application
Modify `v2/main.go`:
- Replace in-memory SessionService with PersistentSessionService
- Add session restoration middleware
- Configure session directory and cleanup intervals
- Add graceful shutdown to save pending sessions

### Phase 5: Update UI Handlers
Modify `v2/handlers/ui_handlers.go`:
- Add session restoration logic to GetStep and other handlers
- Implement automatic state saving after user actions
- Add UI indicators for session restoration
- Handle session corruption gracefully

### Phase 6: Session Directory Structure
Create session storage directory:
```
v2/sessions/
├── active/          # Current active sessions
├── backup/          # Backup copies for recovery
└── expired/         # Expired sessions (for cleanup)
```

### Phase 7: Enhanced Error Handling
- Session corruption detection and recovery
- Graceful fallback when session data is invalid
- User notifications for session restoration
- Logging for session persistence operations

### Phase 8: Version Update and Testing
- Update version in `main.go` to v2.8.0
- Test session persistence across browser refresh
- Verify state restoration maintains exact user progress
- Test session cleanup and expiration

## ACCEPTANCE CRITERIA
- [ ] Browser refresh preserves all workflow state
- [ ] Selected documents remain selected after refresh
- [ ] Extracted data persists through navigation
- [ ] Processing results are maintained across sessions
- [ ] Server restart preserves active sessions
- [ ] Session expiration works correctly (24 hour TTL)
- [ ] No data loss during normal operation
- [ ] Session corruption is handled gracefully
- [ ] Version v2.8.0 displays in masthead
- [ ] Performance impact is minimal (< 100ms overhead)

## TESTING PLAN
1. **Persistence Tests**: Verify data survives refresh/restart
2. **Corruption Tests**: Test with corrupted session files
3. **Concurrent Access Tests**: Multiple users with same session token
4. **Performance Tests**: Measure session save/load times
5. **Cleanup Tests**: Verify expired sessions are properly cleaned
6. **State Integrity Tests**: Ensure complex objects serialize correctly

## SESSION SECURITY CONSIDERATIONS
- Session tokens should be cryptographically secure
- Session data should not contain sensitive client information in plaintext
- File permissions should restrict access to session data
- Session cleanup should securely delete expired data

## GIT WORKFLOW
```bash
# Development
git add .
git commit -m "[TASK-3] Implement persistent session management"
git push origin feature/task-3-session-persistence

# Testing
./scripts/start.sh
# Test browser refresh during workflow
# Verify session restoration works
# Test with multiple concurrent sessions
# Verify v2.8.0 in masthead

# Pull Request
gh pr create --title "TASK-3: Fix Session Persistence on Refresh" --body "
## Summary
- Replaces in-memory sessions with file-based persistent storage
- Preserves all workflow state through browser refresh and server restart
- Adds session restoration with corruption handling
- Maintains performance with < 100ms overhead

## Testing
- [x] Browser refresh preserves complete workflow state
- [x] Server restart maintains active sessions
- [x] Session corruption handled gracefully
- [x] Performance impact minimal
- [x] Concurrent access works correctly
- [x] Version v2.8.0 displays correctly

## Impact
Eliminates data loss issues and provides production-ready session reliability for lawyer workflows.
"

# After PR approval and merge
git checkout main
git pull origin main
git branch -d feature/task-3-session-persistence
```

## DEPENDENCIES
- **Independent**: Can be developed parallel to Tasks 1-2
- **Benefits**: All subsequent tasks will benefit from reliable session management

## NOTES
- This task dramatically improves user experience reliability
- Success here enables confident production deployment
- Consider adding session analytics to monitor usage patterns
- File-based storage is simpler than database for this use case
- Could be upgraded to Redis/database later if needed

## EVALUATION/PLANNING
1. Review objectives for Task 3
2. Confirm file-based session storage is appropriate for current scale
3. Validate session security measures are sufficient
4. Consider backup and disaster recovery scenarios
5. Plan for monitoring session storage disk usage

**Stop. Confirm you understand. Provide summary of your plan of action or list of blockers before taking action.**