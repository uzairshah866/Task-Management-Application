package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/user/taskapp/internal/auth"
	"github.com/user/taskapp/internal/db"
	"github.com/user/taskapp/internal/models"
)

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	lr := io.LimitReader(r.Body, 1<<20)
	if err := json.NewDecoder(lr).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := validateEmail(req.Email); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validatePassword(req.Password); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	existing, err := db.GetUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	if existing != nil {
		respondError(w, http.StatusConflict, "email already registered")
		return
	}

	user, err := auth.NewUser(req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	if err := db.CreateUser(user); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save user")
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	user.Password = ""
	respondJSON(w, http.StatusCreated, models.AuthResponse{Token: token, User: *user})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	lr := io.LimitReader(r.Body, 1<<20)
	if err := json.NewDecoder(lr).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusUnauthorized, "invalid credentials")
		} else {
			respondError(w, http.StatusInternalServerError, "database error")
		}
		return
	}

	if !auth.VerifyPassword(user.Password, req.Password) {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	user.Password = ""
	respondJSON(w, http.StatusOK, models.AuthResponse{Token: token, User: *user})
}

func (s *Server) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := db.GetUserByID(userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	user.Password = ""
	respondJSON(w, http.StatusOK, user)
}
