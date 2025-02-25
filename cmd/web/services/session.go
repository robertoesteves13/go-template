package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
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
	User         User
	CreatedAt   time.Time
	TimeToLive time.Duration
}

type SessionManager[User any] struct {
	servers []string
}

func NewSessionManager[User any](servers ...string) (*SessionManager[User], error) {
	mc := memcache.New(servers...)
	err := mc.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping memcache: [%v]", err)
	}

	return &SessionManager[User]{servers: servers}, nil
}

func (sm *SessionManager[User]) createSession(u User) (string, error) {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		panic("unreachable error on session.go: " + err.Error())
	}
	id := base64.URLEncoding.EncodeToString(b)

	info := SessionInfo[User]{
		User:         u,
		CreatedAt:   time.Now(),
		TimeToLive: 24 * time.Hour,
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(info)
	if err != nil {
		return "", fmt.Errorf("failed to encode session info: [%v]", err)
	}

	mc := memcache.New(sm.servers...)
	err = mc.Add(&memcache.Item{
		Key:        id,
		Value:      buf.Bytes(),
		Expiration: int32(info.TimeToLive.Seconds()),
	})

	if err != nil {
		return "", fmt.Errorf("failed to save to memcache: [%v]", err)
	}

	return id, nil
}

func (sm *SessionManager[User]) destroySession(id string) error {
	mc := memcache.New(sm.servers...)
	err := mc.Delete(id)

	if err != nil {
		return fmt.Errorf("failed to destroy session: [%v]", err)
	}

	return nil
}

func (sm *SessionManager[User]) getSessionInfo(id string) (SessionInfo[User], error) {
	var info SessionInfo[User]
	mc := memcache.New(sm.servers...)
	item, err := mc.Get(id)
	if err != nil {
		return info, fmt.Errorf("failed to get session: [%v]", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(item.Value))
	err = dec.Decode(&info)
	if err != nil {
		return info, fmt.Errorf("failed to decode session: [%v]", err)
	}

	return info, nil
}

type validateFunc[User any] func(*http.Request) (*User, error)

func (sm *SessionManager[User]) LoginRoute(router chi.Router, vf validateFunc[User]) {
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		if user, err := vf(r); user != nil {
			id, err := sm.createSession(*user)
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

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
			info, err := sm.getSessionInfo(id)
			if err == nil {
				if info.CreatedAt.Add(info.TimeToLive).Before(time.Now()) {
					sm.destroySession(id)
					http.SetCookie(w, &http.Cookie{
						Name:   "id",
						Value:  "",
						MaxAge: -1,
					})
				} else {
					ctx = context.WithValue(ctx, UserSession, info)
				}
			}
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
