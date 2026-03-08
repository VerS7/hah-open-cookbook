package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	l "open-hah-cookbook/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			interceptWriter := &statusInterceptor{w, http.StatusOK}

			defer func() {
				s := interceptWriter.status

				msg := fmt.Sprintf("%s %s %d %s => %s",
					r.Method,
					r.RequestURI,
					s,
					r.RemoteAddr,
					time.Since(start),
				)

				switch {
				case s >= 400 && s < 500:
					l.Default.Error(msg)
					return
				case s >= 500:
					l.Default.Panic(msg)
					return
				default:
					l.Default.Info(msg)
					return
				}
			}()

			next.ServeHTTP(interceptWriter, r)
		},
	)
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					l.Default.Error(err)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}

func AllowCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Type, Content-Length")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Content-Disposition, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func (sr *UsersAPIServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token, err := getAuthTokenFromHeaders(r)
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, H{"error": "no authorization header provided"})
				return
			}

			if s := sr.HasSession(token); s != nil {
				next.ServeHTTP(w, r)
				return
			}

			WriteJSON(w, http.StatusUnauthorized, H{"error": "no authorization"})
		},
	)
}

func (sr *UsersAPIServer) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token, err := getAuthTokenFromHeaders(r)
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, H{"error": "no authorization header provided"})
				return
			}

			s := sr.HasSession(token)
			if s == nil {
				WriteJSON(w, http.StatusUnauthorized, H{"error": "no authorization"})
				return
			}

			if !s.IsAdmin {
				WriteJSON(w, http.StatusForbidden, H{"error": "no admin privileges"})
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

type statusInterceptor struct {
	http.ResponseWriter
	status int
}

func (s *statusInterceptor) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func getAuthTokenFromHeaders(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("'Authorization' header not provided")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
