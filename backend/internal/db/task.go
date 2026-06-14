package db

import (
	"fmt"
	"time"

	"github.com/user/taskapp/internal/log"
	"github.com/user/taskapp/internal/models"
)

var (
	validStatuses = map[string]bool{
		"pending": true, "in_progress": true, "completed": true,
	}
	validPriorities = map[string]bool{
		"low": true, "medium": true, "high": true,
	}
	validSortColumns = map[string]bool{
		"created_at": true, "updated_at": true, "due_date": true,
		"priority": true, "status": true, "title": true,
	}
)

func CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (id, user_id, title, description, status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := DB.Exec(
		query, task.ID, task.UserID, task.Title, task.Description,
		task.Status, task.Priority, task.DueDate, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func GetTask(id, userID string) (*models.Task, error) {
	task := &models.Task{}
	query := `
		SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks WHERE id = $1 AND user_id = $2
	`
	err := DB.QueryRow(query, id, userID).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description,
		&task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func ListTasks(userID string, q *models.ListTasksQuery) ([]*models.Task, int, error) {
	start := time.Now()
	defer func() {
		if time.Since(start) > 500*time.Millisecond {
			log.Warn("Slow query detected: ListTasks took %dms", time.Since(start).Milliseconds())
		}
	}()

	countQuery := `SELECT COUNT(*) FROM tasks WHERE user_id = $1`
	args := []interface{}{userID}
	argCount := 1

	listQuery := `
		SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks WHERE user_id = $1
	`

	if q.Status != nil && *q.Status != "" {
		if !isValidStatus(*q.Status) {
			return nil, 0, fmt.Errorf("invalid status value")
		}
		argCount++
		countQuery += fmt.Sprintf(" AND status = $%d", argCount)
		listQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *q.Status)
	}

	if q.Priority != nil && *q.Priority != "" {
		if !isValidPriority(*q.Priority) {
			return nil, 0, fmt.Errorf("invalid priority value")
		}
		argCount++
		countQuery += fmt.Sprintf(" AND priority = $%d", argCount)
		listQuery += fmt.Sprintf(" AND priority = $%d", argCount)
		args = append(args, *q.Priority)
	}

	if q.Search != nil && *q.Search != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		listQuery += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+*q.Search+"%")
	}

	sortBy := "created_at"
	if q.SortBy != nil && *q.SortBy != "" {
		if !isValidSortBy(*q.SortBy) {
			return nil, 0, fmt.Errorf("invalid sort column")
		}
		sortBy = *q.SortBy
	}
	sortOrder := "DESC"
	if q.SortOrder != nil && *q.SortOrder != "" {
		if *q.SortOrder != "ASC" && *q.SortOrder != "DESC" {
			return nil, 0, fmt.Errorf("invalid sort order")
		}
		sortOrder = *q.SortOrder
	}
	listQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 10
	}

	offset := (q.Page - 1) * q.PageSize
	listQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", q.PageSize, offset)

	var total int
	if err := DB.QueryRow(countQuery, args[:argCount]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := DB.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0, q.PageSize)
	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, rows.Err()
}

func UpdateTask(task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = $6
		WHERE id = $7 AND user_id = $8
	`
	result, err := DB.Exec(
		query, task.Title, task.Description, task.Status, task.Priority,
		task.DueDate, task.UpdatedAt, task.ID, task.UserID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTask(id, userID string) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func isValidStatus(s string) bool   { return validStatuses[s] }
func isValidPriority(p string) bool { return validPriorities[p] }
func isValidSortBy(col string) bool { return validSortColumns[col] }
