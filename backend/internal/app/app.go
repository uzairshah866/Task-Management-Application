package app

import (
	"log"

	"github.com/user/taskapp/internal/api"
	"github.com/user/taskapp/internal/auth"
	"github.com/user/taskapp/internal/db"
	"github.com/user/taskapp/internal/models"
)

func NewServer() *api.Server {
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return api.NewServer()
}

func RunMigrations() {
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			priority VARCHAR(50) NOT NULL DEFAULT 'medium',
			due_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at)`,
	}

	for i, migration := range migrations {
		if _, err := db.DB.Exec(migration); err != nil {
			log.Fatalf("Migration %d failed: %v", i, err)
		}
	}

	log.Println("Migrations completed successfully")
}

func SeedDatabase() {
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	user, err := auth.NewUser("test@example.com", "password123")
	if err != nil {
		log.Fatalf("Failed to create seed user: %v", err)
	}

	if err := db.CreateUser(user); err != nil {
		log.Printf("Seed user already exists, fetching: %v", err)
		existing, err := db.GetUserByEmail("test@example.com")
		if err != nil {
			log.Fatalf("Failed to get seed user: %v", err)
		}
		user = existing
	}

	sampleTasks := []*models.Task{
		{
			ID:          "task-1",
			UserID:      user.ID,
			Title:       "Setup project",
			Description: "Initialize project structure",
			Status:      models.StatusCompleted,
			Priority:    models.PriorityHigh,
		},
		{
			ID:          "task-2",
			UserID:      user.ID,
			Title:       "Implement API",
			Description: "Build REST API endpoints",
			Status:      models.StatusInProgress,
			Priority:    models.PriorityHigh,
		},
		{
			ID:          "task-3",
			UserID:      user.ID,
			Title:       "Create frontend",
			Description: "Build Next.js interface",
			Status:      models.StatusPending,
			Priority:    models.PriorityMedium,
		},
	}

	for _, task := range sampleTasks {
		if err := db.CreateTask(task); err != nil {
			log.Printf("Seed task %s already exists, skipping", task.ID)
		}
	}

	log.Println("Database seeded successfully")
}
