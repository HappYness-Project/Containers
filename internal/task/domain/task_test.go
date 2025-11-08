package domain

import (
	"strings"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name        string
		taskName    string
		description string
		targetDate  time.Time
		priority    string
		category    string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid task creation",
			taskName:    "Test Task",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "high",
			category:    "work",
			wantErr:     false,
		},
		{
			name:        "empty task name should fail",
			taskName:    "",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "medium",
			category:    "work",
			wantErr:     true,
			errMsg:      "task name cannot be empty",
		},
		{
			name:        "whitespace only task name should fail",
			taskName:    "   ",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "medium",
			category:    "work",
			wantErr:     true,
			errMsg:      "task name cannot be empty",
		},
		{
			name:        "task name too long should fail",
			taskName:    strings.Repeat("a", 256),
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "medium",
			category:    "work",
			wantErr:     true,
			errMsg:      "task name cannot exceed 255 characters",
		},
		{
			name:        "invalid priority defaults to medium",
			taskName:    "Test Task",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "invalid",
			category:    "work",
			wantErr:     false,
		},
		{
			name:        "empty priority defaults to medium",
			taskName:    "Test Task",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "",
			category:    "work",
			wantErr:     false,
		},
		{
			name:        "valid priority low",
			taskName:    "Test Task",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "LOW",
			category:    "work",
			wantErr:     false,
		},
		{
			name:        "valid priority urgent",
			taskName:    "Test Task",
			description: "Test Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "URGENT",
			category:    "work",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := CreateTask(tt.taskName, tt.description, tt.targetDate, tt.priority, tt.category)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateTask() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("CreateTask() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateTask() unexpected error = %v", err)
				return
			}

			// Verify task fields
			if task == nil {
				t.Error("CreateTask() returned nil task")
				return
			}

			if task.TaskId == "" {
				t.Error("CreateTask() TaskId is empty")
			}

			if task.TaskName != strings.TrimSpace(tt.taskName) {
				t.Errorf("CreateTask() TaskName = %v, want %v", task.TaskName, strings.TrimSpace(tt.taskName))
			}

			if task.TaskDesc != strings.TrimSpace(tt.description) {
				t.Errorf("CreateTask() TaskDesc = %v, want %v", task.TaskDesc, strings.TrimSpace(tt.description))
			}

			if task.TaskType != "" {
				t.Errorf("CreateTask() TaskType = %v, want empty string", task.TaskType)
			}

			if task.TargetDate != tt.targetDate {
				t.Errorf("CreateTask() TargetDate = %v, want %v", task.TargetDate, tt.targetDate)
			}

			if task.Category != tt.category {
				t.Errorf("CreateTask() Category = %v, want %v", task.Category, tt.category)
			}

			if task.IsCompleted != false {
				t.Errorf("CreateTask() IsCompleted = %v, want false", task.IsCompleted)
			}

			if task.IsImportant != false {
				t.Errorf("CreateTask() IsImportant = %v, want false", task.IsImportant)
			}

			// Verify priority handling
			expectedPriority := strings.ToLower(tt.priority)
			if tt.priority == "" || (tt.priority != "low" && tt.priority != "medium" && tt.priority != "high" && tt.priority != "urgent" && strings.ToLower(tt.priority) != "low" && strings.ToLower(tt.priority) != "medium" && strings.ToLower(tt.priority) != "high" && strings.ToLower(tt.priority) != "urgent") {
				expectedPriority = "medium"
			}

			if task.Priority != expectedPriority {
				t.Errorf("CreateTask() Priority = %v, want %v", task.Priority, expectedPriority)
			}

			// Verify timestamps are in UTC
			if task.CreatedAt.Location() != time.UTC {
				t.Errorf("CreateTask() CreatedAt timezone = %v, want UTC", task.CreatedAt.Location())
			}

			if task.UpdatedAt.Location() != time.UTC {
				t.Errorf("CreateTask() UpdatedAt timezone = %v, want UTC", task.UpdatedAt.Location())
			}

			// Verify CreatedAt and UpdatedAt are the same for new tasks
			if !task.CreatedAt.Equal(task.UpdatedAt) {
				t.Errorf("CreateTask() CreatedAt (%v) and UpdatedAt (%v) should be equal", task.CreatedAt, task.UpdatedAt)
			}

			// Verify timestamps are recent (within last second)
			now := time.Now().UTC()
			if now.Sub(task.CreatedAt) > time.Second {
				t.Errorf("CreateTask() CreatedAt is not recent: %v", task.CreatedAt)
			}
		})
	}
}

func TestCreateTask_PriorityValidation(t *testing.T) {
	validPriorities := []string{"low", "medium", "high", "urgent", "LOW", "MEDIUM", "HIGH", "URGENT"}
	
	for _, priority := range validPriorities {
		t.Run("priority_"+priority, func(t *testing.T) {
			task, err := CreateTask("Test Task", "Description", time.Now().Add(time.Hour), priority, "work")
			
			if err != nil {
				t.Errorf("CreateTask() with priority %s failed: %v", priority, err)
				return
			}
			
			expectedPriority := strings.ToLower(priority)
			if task.Priority != expectedPriority {
				t.Errorf("CreateTask() Priority = %v, want %v", task.Priority, expectedPriority)
			}
		})
	}
}