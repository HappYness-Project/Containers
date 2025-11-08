package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskId      string    `json:"id"`
	TaskName    string    `json:"name"`
	TaskDesc    string    `json:"description"`
	TaskType    string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TargetDate  time.Time `json:"target_date"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	IsCompleted bool      `json:"is_completed"`
	IsImportant bool      `json:"is_important"`
}

func CreateTask(name, description string, targetDate time.Time, priority, category string) (*Task, error) {
	// Domain validation
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("task name cannot be empty")
	}

	if len(name) > 255 {
		return nil, errors.New("task name cannot exceed 255 characters")
	}

	// Validate priority
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
		"urgent": true,
	}

	if priority == "" || !validPriorities[strings.ToLower(priority)] {
		priority = "medium" // Default to medium if empty or invalid
	} else {
		priority = strings.ToLower(priority)
	}

	now := time.Now().UTC()

	task := &Task{
		TaskId:      uuid.New().String(),
		TaskName:    strings.TrimSpace(name),
		TaskDesc:    strings.TrimSpace(description),
		TaskType:    "",
		CreatedAt:   now,
		UpdatedAt:   now,
		TargetDate:  targetDate,
		Priority:    priority,
		Category:    category,
		IsCompleted: false,
		IsImportant: false,
	}

	return task, nil
}
