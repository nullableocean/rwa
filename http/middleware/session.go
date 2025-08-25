package middleware

import (
	"context"
	"net/http"
	"rwa/http/handlers"
	"rwa/internal/services"

	"github.com/gorilla/mux"
)

type SessionGuard struct {
	sesManager *services.SessionManager
}

func NewSessionGuard(sesManager *services.SessionManager) *SessionGuard {
	return &SessionGuard{
		sesManager: sesManager,
	}
}

func (sg *SessionGuard) GetAuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := handlers.GetTokenFromRequest(r)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("No Auth"))
				return
			}

			session, err := sg.sesManager.Get(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("No Auth - Bad Token"))
				return
			}

			newContext := context.WithValue(r.Context(), handlers.UserCtxKey, session.UserId)
			next.ServeHTTP(w, r.WithContext(newContext))
		})
	}
}
