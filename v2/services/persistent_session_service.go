package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionData represents persistent session data
type SessionData struct {
	SessionID       string                 `json:"sessionId"`
	WorkflowState   *WorkflowState        `json:"workflowState"`
	CreatedAt       time.Time             `json:"createdAt"`
	LastAccessed    time.Time             `json:"lastAccessed"`
	ExpiresAt       time.Time             `json:"expiresAt"`
}

// PersistentSessionService manages session state with file-based persistence
type PersistentSessionService struct {
	SessionDir    string
	BackupDir     string
	CleanupTicker *time.Ticker
	FileLock      sync.RWMutex
	ttl           time.Duration
}

// NewPersistentSessionService creates a new persistent session service
func NewPersistentSessionService(sessionDir string, ttl time.Duration) (*PersistentSessionService, error) {
	// Create session directories
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %v", err)
	}
	
	backupDir := filepath.Join(sessionDir, "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %v", err)
	}
	
	service := &PersistentSessionService{
		SessionDir: sessionDir,
		BackupDir:  backupDir,
		FileLock:   sync.RWMutex{},
		ttl:        ttl,
	}
	
	// Start cleanup goroutine
	service.CleanupTicker = time.NewTicker(time.Hour)
	go service.cleanupExpiredSessions()
	
	log.Printf("[INFO] Persistent session service initialized with directory: %s", sessionDir)
	return service, nil
}

// GetSession retrieves or creates a workflow state for the given session ID
func (s *PersistentSessionService) GetSession(sessionID string) *WorkflowState {
	sessionData := s.loadSessionData(sessionID)
	
	if sessionData == nil {
		// Create new session
		workflowState := &WorkflowState{
			CurrentStep: 0,
			LastUpdated: time.Now(),
		}
		s.saveSessionData(sessionID, workflowState)
		return workflowState
	}
	
	// Update last accessed time
	sessionData.LastAccessed = time.Now()
	s.saveSessionDataStruct(sessionData)
	
	return sessionData.WorkflowState
}

// SetSession saves the workflow state for the given session ID
func (s *PersistentSessionService) SetSession(sessionID string, state *WorkflowState) {
	state.LastUpdated = time.Now()
	s.saveSessionData(sessionID, state)
}

// UpdateSession partially updates session data with atomic file operations
func (s *PersistentSessionService) UpdateSession(sessionID string, updateFunc func(*WorkflowState)) {
	s.FileLock.Lock()
	defer s.FileLock.Unlock()
	
	sessionData := s.loadSessionDataUnsafe(sessionID)
	
	if sessionData == nil {
		// Create new session
		workflowState := &WorkflowState{
			CurrentStep: 0,
			LastUpdated: time.Now(),
		}
		updateFunc(workflowState)
		s.saveSessionDataUnsafe(sessionID, workflowState)
		return
	}
	
	// Update existing session
	updateFunc(sessionData.WorkflowState)
	sessionData.WorkflowState.LastUpdated = time.Now()
	sessionData.LastAccessed = time.Now()
	s.saveSessionDataStructUnsafe(sessionData)
}

// DeleteSession removes a session file
func (s *PersistentSessionService) DeleteSession(sessionID string) {
	s.FileLock.Lock()
	defer s.FileLock.Unlock()
	
	sessionFile := s.getSessionFilePath(sessionID)
	if err := os.Remove(sessionFile); err != nil && !os.IsNotExist(err) {
		log.Printf("[WARN] Failed to delete session file %s: %v", sessionFile, err)
	}
}

// RestoreSession loads session data from disk (used by middleware)
func (s *PersistentSessionService) RestoreSession(sessionID string) *SessionData {
	return s.loadSessionData(sessionID)
}

// SaveState persists state immediately (for critical state changes)
func (s *PersistentSessionService) SaveState(sessionID string, state interface{}) {
	if workflowState, ok := state.(*WorkflowState); ok {
		s.SetSession(sessionID, workflowState)
	}
}

// GetSessionCount returns the number of active sessions
func (s *PersistentSessionService) GetSessionCount() int {
	s.FileLock.RLock()
	defer s.FileLock.RUnlock()
	
	files, err := ioutil.ReadDir(s.SessionDir)
	if err != nil {
		log.Printf("[WARN] Failed to read session directory: %v", err)
		return 0
	}
	
	count := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			count++
		}
	}
	return count
}

// Private helper methods

func (s *PersistentSessionService) getSessionFilePath(sessionID string) string {
	return filepath.Join(s.SessionDir, fmt.Sprintf("%s.json", sessionID))
}

func (s *PersistentSessionService) getBackupFilePath(sessionID string) string {
	timestamp := time.Now().Format("20060102_150405")
	return filepath.Join(s.BackupDir, fmt.Sprintf("%s_%s.json", sessionID, timestamp))
}

func (s *PersistentSessionService) loadSessionData(sessionID string) *SessionData {
	s.FileLock.RLock()
	defer s.FileLock.RUnlock()
	return s.loadSessionDataUnsafe(sessionID)
}

func (s *PersistentSessionService) loadSessionDataUnsafe(sessionID string) *SessionData {
	sessionFile := s.getSessionFilePath(sessionID)
	
	data, err := ioutil.ReadFile(sessionFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[WARN] Failed to read session file %s: %v", sessionFile, err)
		}
		return nil
	}
	
	var sessionData SessionData
	if err := json.Unmarshal(data, &sessionData); err != nil {
		log.Printf("[ERROR] Failed to unmarshal session data for %s: %v", sessionID, err)
		// Create backup of corrupted file
		s.createCorruptedBackup(sessionID, data)
		return nil
	}
	
	// Check if session is expired
	if time.Now().After(sessionData.ExpiresAt) {
		log.Printf("[INFO] Session %s has expired, removing", sessionID)
		os.Remove(sessionFile)
		return nil
	}
	
	return &sessionData
}

func (s *PersistentSessionService) saveSessionData(sessionID string, state *WorkflowState) {
	s.FileLock.Lock()
	defer s.FileLock.Unlock()
	s.saveSessionDataUnsafe(sessionID, state)
}

func (s *PersistentSessionService) saveSessionDataUnsafe(sessionID string, state *WorkflowState) {
	now := time.Now()
	sessionData := &SessionData{
		SessionID:     sessionID,
		WorkflowState: state,
		CreatedAt:     now,
		LastAccessed:  now,
		ExpiresAt:     now.Add(s.ttl),
	}
	
	s.saveSessionDataStructUnsafe(sessionData)
}

func (s *PersistentSessionService) saveSessionDataStruct(sessionData *SessionData) {
	s.FileLock.Lock()
	defer s.FileLock.Unlock()
	s.saveSessionDataStructUnsafe(sessionData)
}

func (s *PersistentSessionService) saveSessionDataStructUnsafe(sessionData *SessionData) {
	sessionFile := s.getSessionFilePath(sessionData.SessionID)
	
	// Create backup of existing file before overwriting
	if _, err := os.Stat(sessionFile); err == nil {
		s.createBackupUnsafe(sessionData.SessionID)
	}
	
	// Marshal session data
	data, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		log.Printf("[ERROR] Failed to marshal session data for %s: %v", sessionData.SessionID, err)
		return
	}
	
	// Write to temporary file first (atomic operation)
	tempFile := sessionFile + ".tmp"
	if err := ioutil.WriteFile(tempFile, data, 0644); err != nil {
		log.Printf("[ERROR] Failed to write temporary session file %s: %v", tempFile, err)
		return
	}
	
	// Rename to final location (atomic operation on most filesystems)
	if err := os.Rename(tempFile, sessionFile); err != nil {
		log.Printf("[ERROR] Failed to rename session file %s: %v", sessionFile, err)
		os.Remove(tempFile) // Clean up temp file
		return
	}
}

func (s *PersistentSessionService) createBackupUnsafe(sessionID string) {
	sessionFile := s.getSessionFilePath(sessionID)
	backupFile := s.getBackupFilePath(sessionID)
	
	data, err := ioutil.ReadFile(sessionFile)
	if err != nil {
		return // Source file doesn't exist or can't be read
	}
	
	if err := ioutil.WriteFile(backupFile, data, 0644); err != nil {
		log.Printf("[WARN] Failed to create backup for session %s: %v", sessionID, err)
	}
}

func (s *PersistentSessionService) createCorruptedBackup(sessionID string, data []byte) {
	timestamp := time.Now().Format("20060102_150405")
	corruptedFile := filepath.Join(s.BackupDir, fmt.Sprintf("%s_corrupted_%s.json", sessionID, timestamp))
	
	if err := ioutil.WriteFile(corruptedFile, data, 0644); err != nil {
		log.Printf("[WARN] Failed to create corrupted file backup for session %s: %v", sessionID, err)
	} else {
		log.Printf("[INFO] Created backup of corrupted session file: %s", corruptedFile)
	}
}

func (s *PersistentSessionService) cleanupExpiredSessions() {
	defer s.CleanupTicker.Stop()
	
	for range s.CleanupTicker.C {
		s.FileLock.Lock()
		
		files, err := ioutil.ReadDir(s.SessionDir)
		if err != nil {
			log.Printf("[WARN] Failed to read session directory for cleanup: %v", err)
			s.FileLock.Unlock()
			continue
		}
		
		deletedCount := 0
		now := time.Now()
		
		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
				continue
			}
			
			sessionFile := filepath.Join(s.SessionDir, file.Name())
			data, err := ioutil.ReadFile(sessionFile)
			if err != nil {
				continue
			}
			
			var sessionData SessionData
			if err := json.Unmarshal(data, &sessionData); err != nil {
				// Remove corrupted files
				log.Printf("[WARN] Removing corrupted session file: %s", sessionFile)
				os.Remove(sessionFile)
				deletedCount++
				continue
			}
			
			if now.After(sessionData.ExpiresAt) {
				log.Printf("[INFO] Removing expired session: %s", sessionData.SessionID)
				os.Remove(sessionFile)
				deletedCount++
			}
		}
		
		activeCount := len(files) - deletedCount
		log.Printf("[INFO] Session cleanup: removed %d expired sessions, %d active sessions remaining", 
			deletedCount, activeCount)
		
		s.FileLock.Unlock()
	}
}

// Shutdown gracefully stops the cleanup ticker
func (s *PersistentSessionService) Shutdown() {
	if s.CleanupTicker != nil {
		s.CleanupTicker.Stop()
	}
	log.Printf("[INFO] Persistent session service shutdown complete")
}