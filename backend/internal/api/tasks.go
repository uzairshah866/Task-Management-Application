package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/user/taskapp/internal/db"
	"github.com/user/taskapp/internal/models"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateTaskRequest
	lr := io.LimitReader(r.Body, 1<<20)
	if err := json.NewDecoder(lr).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := validateCreateTask(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task := &models.Task{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.CreateTask(task); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (s *Server) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	q := &models.ListTasksQuery{
		Status:    getStringParam(r, "status"),
		Priority:  getStringParam(r, "priority"),
		Search:    getStringParam(r, "search"),
		SortBy:    getStringParam(r, "sort_by"),
		SortOrder: getStringParam(r, "sort_order"),
		Page:      getIntParam(r, "page", 1),
		PageSize:  getIntParam(r, "page_size", 10),
	}

	tasks, total, err := db.ListTasks(userID, q)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	response := map[string]interface{}{
		"tasks":      tasks,
		"total":      total,
		"page":       q.Page,
		"page_size":  q.PageSize,
		"total_pages": (total + q.PageSize - 1) / q.PageSize,
	}

	respondJSON(w, http.StatusOK, response)
}

func (s *Server) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	taskID := mux.Vars(r)["id"]

	task, err := db.GetTask(taskID, userID)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get task")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (s *Server) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	taskID := mux.Vars(r)["id"]

	var req models.UpdateTaskRequest
	lr := io.LimitReader(r.Body, 1<<20)
	if err := json.NewDecoder(lr).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := validateUpdateTask(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := db.GetTask(taskID, userID)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get task")
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	task.UpdatedAt = time.Now()

	if err := db.UpdateTask(task); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (s *Server) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	taskID := mux.Vars(r)["id"]

	if err := db.DeleteTask(taskID, userID); err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "task deleted"})
}
