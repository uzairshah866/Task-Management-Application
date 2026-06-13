package api

import (
	"testing"

	"github.com/user/taskapp/internal/models"
)

func TestValidateCreateTask(t *testing.T) {
	tests := []struct {
		name      string
		req       models.CreateTaskRequest
		wantErr   bool
		wantTitle string
		wantStatus models.Status
		wantPriority models.Priority
	}{
		{
			name:    "empty title returns error",
			req:     models.CreateTaskRequest{Title: ""},
			wantErr: true,
		},
		{
			name:    "whitespace-only title returns error",
			req:     models.CreateTaskRequest{Title: "   "},
			wantErr: true,
		},
		{
			name:    "invalid status returns error",
			req:     models.CreateTaskRequest{Title: "Test", Status: "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid priority returns error",
			req:     models.CreateTaskRequest{Title: "Test", Priority: "invalid"},
			wantErr: true,
		},
		{
			name:         "valid title sets default status and priority",
			req:          models.CreateTaskRequest{Title: "Test task"},
			wantErr:      false,
			wantStatus:   models.StatusPending,
			wantPriority: models.PriorityMedium,
		},
		{
			name:         "explicit status and priority preserved",
			req:          models.CreateTaskRequest{Title: "Test", Status: models.StatusInProgress, Priority: models.PriorityHigh},
			wantErr:      false,
			wantStatus:   models.StatusInProgress,
			wantPriority: models.PriorityHigh,
		},
		{
			name:         "all valid statuses accepted",
			req:          models.CreateTaskRequest{Title: "Test", Status: models.StatusCompleted},
			wantErr:      false,
			wantStatus:   models.StatusCompleted,
		},
		{
			name:         "all valid priorities accepted",
			req:          models.CreateTaskRequest{Title: "Test", Priority: models.PriorityLow},
			wantErr:      false,
			wantPriority: models.PriorityLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateTask(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.wantStatus != "" && tt.req.Status != tt.wantStatus {
					t.Errorf("status = %v, want %v", tt.req.Status, tt.wantStatus)
				}
				if tt.wantPriority != "" && tt.req.Priority != tt.wantPriority {
					t.Errorf("priority = %v, want %v", tt.req.Priority, tt.wantPriority)
				}
			}
		})
	}
}

func TestValidateUpdateTask(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	statusPtr := func(s models.Status) *models.Status { return &s }
	priorityPtr := func(p models.Priority) *models.Priority { return &p }

	tests := []struct {
		name    string
		req     models.UpdateTaskRequest
		wantErr bool
	}{
		{
			name:    "empty title pointer returns error",
			req:     models.UpdateTaskRequest{Title: strPtr("")},
			wantErr: true,
		},
		{
			name:    "whitespace-only title returns error",
			req:     models.UpdateTaskRequest{Title: strPtr("   ")},
			wantErr: true,
		},
		{
			name:    "invalid status returns error",
			req:     models.UpdateTaskRequest{Status: statusPtr("invalid")},
			wantErr: true,
		},
		{
			name:    "invalid priority returns error",
			req:     models.UpdateTaskRequest{Priority: priorityPtr("invalid")},
			wantErr: true,
		},
		{
			name:    "nil fields are valid (partial update)",
			req:     models.UpdateTaskRequest{},
			wantErr: false,
		},
		{
			name:    "valid title update",
			req:     models.UpdateTaskRequest{Title: strPtr("New title")},
			wantErr: false,
		},
		{
			name:    "valid status update",
			req:     models.UpdateTaskRequest{Status: statusPtr(models.StatusCompleted)},
			wantErr: false,
		},
		{
			name:    "valid priority update",
			req:     models.UpdateTaskRequest{Priority: priorityPtr(models.PriorityHigh)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateTask(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
