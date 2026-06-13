package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/user/taskapp/internal/models"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getStringParam(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func getIntParam(r *http.Request, key string, defaultVal int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultVal
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}
	return intVal
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

func validateCreateTask(req *models.CreateTaskRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}

	if string(req.Status) != "" {
		if req.Status != models.StatusPending && req.Status != models.StatusInProgress && req.Status != models.StatusCompleted {
			return fmt.Errorf("invalid status")
		}
	} else {
		req.Status = models.StatusPending
	}

	if string(req.Priority) != "" {
		if req.Priority != models.PriorityLow && req.Priority != models.PriorityMedium && req.Priority != models.PriorityHigh {
			return fmt.Errorf("invalid priority")
		}
	} else {
		req.Priority = models.PriorityMedium
	}

	return nil
}

func validateUpdateTask(req *models.UpdateTaskRequest) error {
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}

	if req.Status != nil {
		if *req.Status != models.StatusPending && *req.Status != models.StatusInProgress && *req.Status != models.StatusCompleted {
			return fmt.Errorf("invalid status")
		}
	}

	if req.Priority != nil {
		if *req.Priority != models.PriorityLow && *req.Priority != models.PriorityMedium && *req.Priority != models.PriorityHigh {
			return fmt.Errorf("invalid priority")
		}
	}

	return nil
}
