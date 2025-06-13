package services

import (
	"sync"
	"time"
)

// WorkflowState represents the user's progress through the legal document workflow
type WorkflowState struct {
	// Step 0: iCloud Setup
	ICloudConnected      bool     `json:"icloudConnected"`
	ICloudUsername       string   `json:"icloudUsername"`
	SelectedParentFolder string   `json:"selectedParentFolder"`
	SelectedCaseFolder   string   `json:"selectedCaseFolder"`
	
	// Step 1: Document Selection
	AvailableDocuments   []ICloudDocument `json:"availableDocuments"`
	SelectedDocuments    []string         `json:"selectedDocuments"`
	
	// Step 2: Template Selection  
	SelectedTemplate     string           `json:"selectedTemplate"`
	
	// Step 3: Review Data
	ProcessingResult     *DocumentProcessingResult `json:"processingResult,omitempty"`
	ClientCase           *ClientCase               `json:"clientCase,omitempty"`
	
	// Metadata
	CurrentStep          int               `json:"currentStep"`
	LastUpdated          time.Time         `json:"lastUpdated"`
	Username             string            `json:"username"`
}

// SessionService manages user session state
type SessionService struct {
	sessions map[string]*WorkflowState
	mutex    sync.RWMutex
	ttl      time.Duration
}

// NewSessionService creates a new session service
func NewSessionService(ttl time.Duration) *SessionService {
	service := &SessionService{
		sessions: make(map[string]*WorkflowState),
		mutex:    sync.RWMutex{},
		ttl:      ttl,
	}
	
	// Start cleanup goroutine
	go service.cleanupExpiredSessions()
	
	return service
}

// GetSession retrieves or creates a workflow state for the given session ID
func (s *SessionService) GetSession(sessionID string) *WorkflowState {
	s.mutex.RLock()
	session, exists := s.sessions[sessionID]
	s.mutex.RUnlock()
	
	if !exists {
		session = &WorkflowState{
			CurrentStep: 0,
			LastUpdated: time.Now(),
		}
		s.SetSession(sessionID, session)
	}
	
	return session
}

// SetSession saves the workflow state for the given session ID
func (s *SessionService) SetSession(sessionID string, state *WorkflowState) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	state.LastUpdated = time.Now()
	s.sessions[sessionID] = state
}

// UpdateSession partially updates session data
func (s *SessionService) UpdateSession(sessionID string, updateFunc func(*WorkflowState)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	session, exists := s.sessions[sessionID]
	if !exists {
		session = &WorkflowState{
			CurrentStep: 0,
			LastUpdated: time.Now(),
		}
		s.sessions[sessionID] = session
	}
	
	updateFunc(session)
	session.LastUpdated = time.Now()
}

// DeleteSession removes a session
func (s *SessionService) DeleteSession(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	delete(s.sessions, sessionID)
}

// cleanupExpiredSessions removes sessions older than TTL
func (s *SessionService) cleanupExpiredSessions() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for sessionID, state := range s.sessions {
			if now.Sub(state.LastUpdated) > s.ttl {
				delete(s.sessions, sessionID)
			}
		}
		s.mutex.Unlock()
	}
}

// GetSessionCount returns the number of active sessions
func (s *SessionService) GetSessionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.sessions)
}