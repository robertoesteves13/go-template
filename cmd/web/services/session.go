package services

// TODO: Do these in the future
// - Alternative storage methods
// - Better redirection system
// - TOTP/email verification

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type SessionKey int

const (
	UserSession SessionKey = iota
)

func GetUserSession[User any](ctx context.Context) *SessionInfo[User] {
	if info, ok := ctx.Value(UserSession).(SessionInfo[User]); ok {
		return &info
	} else {
		return nil
	}
}

type SessionInfo[User any] struct {
	user         User
	created_at   time.Time
	time_to_live time.Duration
}

func (si *SessionInfo[User]) User() *User {
	return &si.user
}

type SessionManager[User any] struct {
	lock     sync.RWMutex
	memstore map[string]SessionInfo[User]
}

func NewSessionManager[User any]() *SessionManager[User] {
	return &SessionManager[User]{
		lock:     sync.RWMutex{},
		memstore: make(map[string]SessionInfo[User]),
	}
}

func (sm *SessionManager[User]) createSession(u User) string {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		panic("unreachable error on session.go: " + err.Error())
	}
	id := base64.URLEncoding.EncodeToString(b)

	sm.lock.Lock()
	sm.memstore[id] = SessionInfo[User]{
		user:       u,
		created_at: time.Now(),
		time_to_live: 24 * time.Hour,
	}
	sm.lock.Unlock()

	return id
}

func (sm *SessionManager[User]) destroySession(id string) {
	sm.lock.Lock()
	delete(sm.memstore, id)
	sm.lock.Unlock()
}

func (sm *SessionManager[User]) getSessionInfo(id string) SessionInfo[User] {
	sm.lock.RLock()
	info := sm.memstore[id]
	sm.lock.RUnlock()

	return info
}

type validateFunc[User any] func(*http.Request) (*User, error)

func (sm *SessionManager[User]) LoginRoute(router chi.Router, vf validateFunc[User]) {
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		if user, err := vf(r); user != nil {
			id := sm.createSession(*user)
			http.SetCookie(w, &http.Cookie{
				Name:     "id",
				Value:    id,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				// TODO: uncomment when http2 is supported
				// Secure: true,
				MaxAge: 3600 * 24,
			})

			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else if err != nil {
			log.Printf("%v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		} else {
			http.Error(w, "403 forbidden", http.StatusForbidden)
		}
	})
}

func (sm *SessionManager[User]) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session_id, err := r.Cookie("id")
		ctx := r.Context()
		if err == nil {
			id := session_id.Value
			info := sm.getSessionInfo(id)
			if info.created_at.Add(info.time_to_live).Before(time.Now()) {
				sm.destroySession(id)
				http.SetCookie(w, &http.Cookie{
					Name: "id",
					Value: "",
					MaxAge: -1,
				})
			} else {
				ctx = context.WithValue(ctx, UserSession, info)
			}
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
