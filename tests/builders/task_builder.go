package builders

import (
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/task/domain"
)

// TaskBuilder provides a fluent interface for building test tasks
type TaskBuilder struct {
	task *domain.Task
	err  error
}

// NewTaskBuilder creates a new task builder with sensible defaults
func NewTaskBuilder() *TaskBuilder {
	task, err := domain.CreateTask(
		"Test Task",
		"Test Description",
		time.Now().Add(24*time.Hour), // Due tomorrow
		"medium",
		"work",
	)

	return &TaskBuilder{
		task: task,
		err:  err,
	}
}

// WithName sets the task name
func (b *TaskBuilder) WithName(name string) *TaskBuilder {
	if b.task != nil {
		b.task.TaskName = name
	}
	return b
}

// WithDescription sets the task description
func (b *TaskBuilder) WithDescription(desc string) *TaskBuilder {
	if b.task != nil {
		b.task.TaskDesc = desc
	}
	return b
}

// WithPriority sets the task priority
func (b *TaskBuilder) WithPriority(priority string) *TaskBuilder {
	if b.task != nil {
		b.task.Priority = priority
	}
	return b
}

// WithCategory sets the task category
func (b *TaskBuilder) WithCategory(category string) *TaskBuilder {
	if b.task != nil {
		b.task.Category = category
	}
	return b
}

// WithTargetDate sets the target date
func (b *TaskBuilder) WithTargetDate(targetDate time.Time) *TaskBuilder {
	if b.task != nil {
		b.task.TargetDate = targetDate
	}
	return b
}

// WithTaskId sets a specific task ID (useful for testing)
func (b *TaskBuilder) WithTaskId(id string) *TaskBuilder {
	if b.task != nil {
		b.task.TaskId = id
	}
	return b
}

// Completed marks the task as completed
func (b *TaskBuilder) Completed() *TaskBuilder {
	if b.task != nil {
		b.task.ToggleCompletion(true)
	}
	return b
}

// Important marks the task as important
func (b *TaskBuilder) Important() *TaskBuilder {
	if b.task != nil {
		b.task.ToggleImportant(true)
	}
	return b
}

// LowPriority sets priority to low
func (b *TaskBuilder) LowPriority() *TaskBuilder {
	return b.WithPriority("low")
}

// HighPriority sets priority to high
func (b *TaskBuilder) HighPriority() *TaskBuilder {
	return b.WithPriority("high")
}

// UrgentPriority sets priority to urgent
func (b *TaskBuilder) UrgentPriority() *TaskBuilder {
	return b.WithPriority("urgent")
}

// DueTomorrow sets target date to tomorrow
func (b *TaskBuilder) DueTomorrow() *TaskBuilder {
	return b.WithTargetDate(time.Now().Add(24 * time.Hour))
}

// DueToday sets target date to today
func (b *TaskBuilder) DueToday() *TaskBuilder {
	return b.WithTargetDate(time.Now())
}

// Overdue sets target date to yesterday
func (b *TaskBuilder) Overdue() *TaskBuilder {
	return b.WithTargetDate(time.Now().Add(-24 * time.Hour))
}

// Build returns the built task or an error
func (b *TaskBuilder) Build() (*domain.Task, error) {
	return b.task, b.err
}

// MustBuild returns the built task or panics (useful in test setup)
func (b *TaskBuilder) MustBuild() *domain.Task {
	if b.err != nil {
		panic(b.err)
	}
	return b.task
}
