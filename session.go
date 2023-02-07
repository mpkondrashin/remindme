package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

const SessionCookieName = "Session"
const SessionMaxDiration = 15 * time.Minute

var counterLock = &sync.Mutex{}

const SessionsCleanupPeriod uint32 = 10

var SessionsCleanupPeriodCounter uint32 = 0

type SessionData struct {
	Data          map[string]interface{}
	LastOperation time.Time
}

type Sessions struct {
	sessions map[string]*SessionData
}

var sessions *Sessions
var lock = &sync.Mutex{}

func Session() *Sessions {
	if sessions == nil {
		lock.Lock()
		defer lock.Unlock()
		if sessions == nil {
			sessions = &Sessions{
				sessions: make(map[string]*SessionData),
			}
		}
	}
	return sessions
}

func (s *Sessions) Get(r *http.Request) (string, *SessionData) {
	counterLock.Lock()
	SessionsCleanupPeriodCounter++
	if SessionsCleanupPeriodCounter > SessionsCleanupPeriod {
		s.Expire()
		SessionsCleanupPeriodCounter = 0
	}
	counterLock.Unlock()
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == SessionCookieName {
			cookie = c
			break
		}
	}
	if cookie == nil {
		return "", nil
	}
	sessionID := cookie.Value
	session := s.sessions[sessionID]
	if session == nil {
		log.Printf("Session not found: %s", sessionID)
		return "", nil
	}
	threshold := time.Now().Add(-SessionMaxDiration)
	if session.LastOperation.Before(threshold) {
		lock.Lock()
		delete(s.sessions, sessionID)
		lock.Unlock()
		return "", nil
	}
	session.LastOperation = time.Now()
	return sessionID, session
}

func (s *Sessions) Start(w http.ResponseWriter) {
	uuid := uuid.New().String()
	newSessionData := &SessionData{
		Data:          make(map[string]interface{}),
		LastOperation: time.Now(),
	}
	lock.Lock()
	s.sessions[uuid] = newSessionData
	lock.Unlock()
	cookie := http.Cookie{
		Name:     SessionCookieName,
		Value:    uuid,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24,
		//		Expires:  time.Now().Add(SessionMaxDiration), this limits session duration. Could help if each
		// https://www.sohamkamani.com/golang/session-cookie-authentication/ - recreating session IDs each request
	}
	http.SetCookie(w, &cookie)
}

func (s *Sessions) End(w http.ResponseWriter, r *http.Request) {
	id, session := s.Get(r)
	if session == nil {
		return
	}
	cookie := http.Cookie{
		Name:     SessionCookieName,
		Value:    id,
		Path:     "/webui/pages",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, &cookie)
	lock.Lock()
	delete(s.sessions, id)
	lock.Unlock()
}

func (s *Sessions) Expire() {
	threshold := time.Now().Add(-SessionMaxDiration)
	lock.Lock()
	defer lock.Unlock()
	for id, session := range s.sessions {
		if session.LastOperation.After(threshold) {
			continue
		}
		delete(s.sessions, id)
	}
}
