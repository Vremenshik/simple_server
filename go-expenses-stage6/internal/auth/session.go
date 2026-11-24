package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]int64
}

func NewSessionManager() *SessionManager {
	return &SessionManager{sessions: make(map[string]int64)}
}

func (s *SessionManager) Create(userID int64) string {
	b := make([]byte, 16)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[sessionID] = userID
	s.mu.Unlock()
	return sessionID
}

func (s *SessionManager) Get(sessionID string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, ok := s.sessions[sessionID]
	return userID, ok
}

func (s *SessionManager) Delete(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

func (s *SessionManager) SetCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name: "session_id", Value: sessionID, Path: "/", HttpOnly: true,
		Expires: time.Now().Add(24 * time.Hour),
	})
}

func (s *SessionManager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", Path: "/", MaxAge: -1})
}
