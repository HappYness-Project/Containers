package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/task/domain"
)

const dbTimeout = time.Second * 5

type TaskRepository interface {
	GetAllTasks() ([]domain.Task, error)
	GetAllTasksByGroupId(groupId int) ([]domain.Task, error)
	GetAllTasksByGroupIdOnlyImportant(groupId int) ([]domain.Task, error)
	GetTaskById(id string) (*domain.Task, error)
	GetTasksByContainerId(containerId string) ([]domain.Task, error)
	CreateTask(taskcontainerId string, task domain.Task) (domain.Task, error)
	UpdateTask(task domain.Task) error
	UpdateImportantTask(id string, isImportant bool) error
	DeleteTask(id string) error
	DoneTask(id string, isDone bool) error
}
type TaskRepo struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepo {
	return &TaskRepo{
		DB: db,
	}
}

func (m *TaskRepo) GetAllTasks() ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetAllTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		task, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (m *TaskRepo) GetTaskById(id string) (*domain.Task, error) {
	rows, err := m.DB.Query(sqlGetTaskById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var task *domain.Task
	for rows.Next() {
		task, err = scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	return task, nil
}

func (m *TaskRepo) GetTasksByContainerId(containerId string) ([]domain.Task, error) {
	rows, err := m.DB.Query(sqlGetAllTasksByContainerId, containerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		task, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (m *TaskRepo) CreateTask(containerId string, task domain.Task) (domain.Task, error) {
	_, err := m.DB.Exec(sqlCreateTask, task.TaskId, task.TaskName, task.TaskDesc, task.TaskType, task.CreatedAt, task.UpdatedAt, task.TargetDate, task.Priority, task.Category, task.IsCompleted, task.IsImportant)
	if err != nil {
		return task, fmt.Errorf("unable to insert into task table : %w", err)
	}
	_, err = m.DB.Exec(sqlCreateTaskForJoinTable,
		containerId, task.TaskId)

	if err != nil {
		return task, fmt.Errorf("unable to insert into taskcontainer_task table : %w", err)
	}

	return task, nil
}

func (m *TaskRepo) UpdateTask(task domain.Task) error {
	_, err := m.DB.Exec(sqlUpdateTask, task.TaskId, task.TaskName, task.TaskDesc, task.UpdatedAt, task.TargetDate, task.Priority, task.Category)
	if err != nil {
		return err
	}
	return nil
}

func (m *TaskRepo) DeleteTask(id string) error {
	_, err := m.DB.Exec(sqlDeleteTaskForJoinTable, id)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(sqlDeleteTask, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *TaskRepo) DoneTask(id string, isDone bool) error {
	_, err := m.DB.Exec(sqlUpdateTaskDoneField, isDone, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *TaskRepo) UpdateImportantTask(id string, isImportant bool) error {
	_, err := m.DB.Exec(sqlUpdateTaskImportantField, isImportant, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *TaskRepo) GetAllTasksByGroupId(groupId int) ([]domain.Task, error) {
	rows, err := m.DB.Query(sqlGetAllTasksByGroupId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		task, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (m *TaskRepo) GetAllTasksByGroupIdOnlyImportant(groupId int) ([]domain.Task, error) {
	rows, err := m.DB.Query(sqlGetAllTasksByGroupIdAndImportant, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		task, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func scanRowsIntoTask(rows *sql.Rows) (*domain.Task, error) {
	task := new(domain.Task)
	err := rows.Scan(
		&task.TaskId,
		&task.TaskName,
		&task.TaskDesc,
		&task.TaskType,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.TargetDate,
		&task.Priority,
		&task.Category,
		&task.IsCompleted,
		&task.IsImportant,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}
