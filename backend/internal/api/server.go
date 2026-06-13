package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/user/taskapp/internal/auth"
)

type contextKey struct{}

// userIDKey is the unexported context key for the authenticated user's ID.
var userIDKey = contextKey{}

// userIDFromContext retrieves the user ID set by AuthMiddleware.
// Returns ("", false) if the value is absent or wrong type.
func userIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok && id != ""
}

type Server struct {
	router      *mux.Router
	port        string
	corsOrigins string
}

func NewServer() *Server {
	router := mux.NewRouter()
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000"
	}

	s := &Server{
		router:      router,
		port:        port,
		corsOrigins: corsOrigins,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/api/health", s.HandleHealth).Methods("GET")

	s.router.HandleFunc("/api/auth/signup", s.HandleSignup).Methods("POST")
	s.router.HandleFunc("/api/auth/login", s.HandleLogin).Methods("POST")

	s.router.HandleFunc("/api/auth/me", s.AuthMiddleware(s.HandleGetMe)).Methods("GET")

	s.router.HandleFunc("/api/tasks", s.AuthMiddleware(s.HandleListTasks)).Methods("GET")
	s.router.HandleFunc("/api/tasks", s.AuthMiddleware(s.HandleCreateTask)).Methods("POST")
	s.router.HandleFunc("/api/tasks/{id}", s.AuthMiddleware(s.HandleGetTask)).Methods("GET")
	s.router.HandleFunc("/api/tasks/{id}", s.AuthMiddleware(s.HandleUpdateTask)).Methods("PATCH")
	s.router.HandleFunc("/api/tasks/{id}", s.AuthMiddleware(s.HandleDeleteTask)).Methods("DELETE")

}

func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.port)
	// Wrap at HTTP level so CORS middleware intercepts ALL requests,
	// including preflight OPTIONS that gorilla/mux would otherwise 405.
	return http.ListenAndServe(":"+s.port, s.corsMiddleware(s.router))
}

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := auth.VerifyToken(parts[1])
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			respondError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.corsOrigins)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}
