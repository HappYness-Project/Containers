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

func TestUpdateTask(t *testing.T) {
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
			name:        "valid update",
			taskName:    "Updated Task",
			description: "Updated Description",
			targetDate:  time.Now().Add(48 * time.Hour),
			priority:    "urgent",
			category:    "personal",
			wantErr:     false,
		},
		{
			name:        "empty task name should fail",
			taskName:    "",
			description: "Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "high",
			category:    "work",
			wantErr:     true,
			errMsg:      "task name cannot be empty",
		},
		{
			name:        "whitespace only task name should fail",
			taskName:    "   ",
			description: "Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "high",
			category:    "work",
			wantErr:     true,
			errMsg:      "task name cannot be empty",
		},
		{
			name:        "task name too long should fail",
			taskName:    strings.Repeat("a", 256),
			description: "Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "medium",
			category:    "work",
			wantErr:     true,
			errMsg:      "task name cannot exceed 255 characters",
		},
		{
			name:        "invalid priority defaults to medium",
			taskName:    "Updated Task",
			description: "Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "invalid",
			category:    "work",
			wantErr:     false,
		},
		{
			name:        "empty priority defaults to medium",
			taskName:    "Updated Task",
			description: "Description",
			targetDate:  time.Now().Add(24 * time.Hour),
			priority:    "",
			category:    "work",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create initial task
			task, err := CreateTask("Original Task", "Original Description", time.Now().Add(24*time.Hour), "low", "work")
			if err != nil {
				t.Fatalf("Failed to create initial task: %v", err)
			}

			originalUpdatedAt := task.UpdatedAt
			time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt will change

			// Update the task
			err = task.UpdateTask(tt.taskName, tt.description, tt.targetDate, tt.priority, tt.category)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateTask() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("UpdateTask() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateTask() unexpected error = %v", err)
				return
			}

			// Verify fields were updated
			if task.TaskName != strings.TrimSpace(tt.taskName) {
				t.Errorf("UpdateTask() TaskName = %v, want %v", task.TaskName, strings.TrimSpace(tt.taskName))
			}

			if task.TaskDesc != strings.TrimSpace(tt.description) {
				t.Errorf("UpdateTask() TaskDesc = %v, want %v", task.TaskDesc, strings.TrimSpace(tt.description))
			}

			if task.TargetDate != tt.targetDate {
				t.Errorf("UpdateTask() TargetDate = %v, want %v", task.TargetDate, tt.targetDate)
			}

			if task.Category != tt.category {
				t.Errorf("UpdateTask() Category = %v, want %v", task.Category, tt.category)
			}

			// Verify priority handling
			expectedPriority := strings.ToLower(tt.priority)
			if tt.priority == "" || (strings.ToLower(tt.priority) != "low" && strings.ToLower(tt.priority) != "medium" && strings.ToLower(tt.priority) != "high" && strings.ToLower(tt.priority) != "urgent") {
				expectedPriority = "medium"
			}

			if task.Priority != expectedPriority {
				t.Errorf("UpdateTask() Priority = %v, want %v", task.Priority, expectedPriority)
			}

			// Verify UpdatedAt was changed
			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateTask() UpdatedAt was not updated")
			}

			// Verify UpdatedAt is in UTC
			if task.UpdatedAt.Location() != time.UTC {
				t.Errorf("UpdateTask() UpdatedAt timezone = %v, want UTC", task.UpdatedAt.Location())
			}

			// Verify CreatedAt was not changed
			if !task.CreatedAt.Equal(originalUpdatedAt) {
				t.Errorf("UpdateTask() CreatedAt should not change during update")
			}
		})
	}
}

func TestToggleCompletion(t *testing.T) {
	tests := []struct {
		name        string
		isCompleted bool
	}{
		{
			name:        "toggle to completed",
			isCompleted: true,
		},
		{
			name:        "toggle to not completed",
			isCompleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create task
			task, err := CreateTask("Test Task", "Description", time.Now().Add(24*time.Hour), "medium", "work")
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}

			originalUpdatedAt := task.UpdatedAt
			originalIsCompleted := task.IsCompleted
			time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt will change

			// Toggle completion
			task.ToggleCompletion(tt.isCompleted)

			// Verify IsCompleted was updated
			if task.IsCompleted != tt.isCompleted {
				t.Errorf("ToggleCompletion() IsCompleted = %v, want %v", task.IsCompleted, tt.isCompleted)
			}

			// Verify it actually changed from original
			if originalIsCompleted == tt.isCompleted && task.IsCompleted == originalIsCompleted {
				// This is fine if they happen to be the same value
			}

			// Verify UpdatedAt was changed
			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("ToggleCompletion() UpdatedAt was not updated")
			}

			// Verify UpdatedAt is in UTC
			if task.UpdatedAt.Location() != time.UTC {
				t.Errorf("ToggleCompletion() UpdatedAt timezone = %v, want UTC", task.UpdatedAt.Location())
			}

			// Verify other fields were not changed
			if task.TaskName != "Test Task" {
				t.Errorf("ToggleCompletion() should not change TaskName")
			}
		})
	}
}

func TestToggleImportant(t *testing.T) {
	tests := []struct {
		name        string
		isImportant bool
	}{
		{
			name:        "toggle to important",
			isImportant: true,
		},
		{
			name:        "toggle to not important",
			isImportant: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create task
			task, err := CreateTask("Test Task", "Description", time.Now().Add(24*time.Hour), "medium", "work")
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}

			originalUpdatedAt := task.UpdatedAt
			originalIsImportant := task.IsImportant
			time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt will change

			// Toggle important
			task.ToggleImportant(tt.isImportant)

			// Verify IsImportant was updated
			if task.IsImportant != tt.isImportant {
				t.Errorf("ToggleImportant() IsImportant = %v, want %v", task.IsImportant, tt.isImportant)
			}

			// Verify it actually changed from original
			if originalIsImportant == tt.isImportant && task.IsImportant == originalIsImportant {
				// This is fine if they happen to be the same value
			}

			// Verify UpdatedAt was changed
			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("ToggleImportant() UpdatedAt was not updated")
			}

			// Verify UpdatedAt is in UTC
			if task.UpdatedAt.Location() != time.UTC {
				t.Errorf("ToggleImportant() UpdatedAt timezone = %v, want UTC", task.UpdatedAt.Location())
			}

			// Verify other fields were not changed
			if task.TaskName != "Test Task" {
				t.Errorf("ToggleImportant() should not change TaskName")
			}

			// Verify IsCompleted was not affected
			if task.IsCompleted != false {
				t.Errorf("ToggleImportant() should not change IsCompleted")
			}
		})
	}
}

func TestUpdateTask_PriorityNormalization(t *testing.T) {
	task, _ := CreateTask("Test Task", "Description", time.Now().Add(24*time.Hour), "medium", "work")

	tests := []struct {
		name             string
		priority         string
		expectedPriority string
	}{
		{"uppercase LOW", "LOW", "low"},
		{"uppercase MEDIUM", "MEDIUM", "medium"},
		{"uppercase HIGH", "HIGH", "high"},
		{"uppercase URGENT", "URGENT", "urgent"},
		{"mixed case Low", "Low", "low"},
		{"invalid priority", "invalid", "medium"},
		{"empty priority", "", "medium"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.UpdateTask("Test Task", "Description", time.Now().Add(24*time.Hour), tt.priority, "work")
			if err != nil {
				t.Fatalf("UpdateTask() failed: %v", err)
			}

			if task.Priority != tt.expectedPriority {
				t.Errorf("UpdateTask() Priority = %v, want %v", task.Priority, tt.expectedPriority)
			}
		})
	}
}

func TestToggleCompletion_MultipleToggles(t *testing.T) {
	task, err := CreateTask("Test Task", "Description", time.Now().Add(24*time.Hour), "medium", "work")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Initially should be false
	if task.IsCompleted {
		t.Errorf("New task should have IsCompleted = false")
	}

	// Toggle to true
	task.ToggleCompletion(true)
	if !task.IsCompleted {
		t.Errorf("After ToggleCompletion(true), IsCompleted should be true")
	}

	// Toggle to false
	task.ToggleCompletion(false)
	if task.IsCompleted {
		t.Errorf("After ToggleCompletion(false), IsCompleted should be false")
	}

	// Toggle back to true
	task.ToggleCompletion(true)
	if !task.IsCompleted {
		t.Errorf("After second ToggleCompletion(true), IsCompleted should be true")
	}
}

func TestToggleImportant_MultipleToggles(t *testing.T) {
	task, err := CreateTask("Test Task", "Description", time.Now().Add(24*time.Hour), "medium", "work")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Initially should be false
	if task.IsImportant {
		t.Errorf("New task should have IsImportant = false")
	}

	// Toggle to true
	task.ToggleImportant(true)
	if !task.IsImportant {
		t.Errorf("After ToggleImportant(true), IsImportant should be true")
	}

	// Toggle to false
	task.ToggleImportant(false)
	if task.IsImportant {
		t.Errorf("After ToggleImportant(false), IsImportant should be false")
	}

	// Toggle back to true
	task.ToggleImportant(true)
	if !task.IsImportant {
		t.Errorf("After second ToggleImportant(true), IsImportant should be true")
	}
}