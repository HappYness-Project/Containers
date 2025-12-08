package repository

const (
	sqlGetAllTasks              = `SELECT * FROM container.task`
	sqlGetTaskById              = `SELECT * FROM container.task WHERE id = $1`
	sqlGetAllTasksByContainerId = `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important
									FROM container.task t
									JOIN container.taskcontainer_task tct
									ON t.id = tct.task_id
									WHERE taskcontainer_id = $1`
	sqlGetAllTasksByGroupId = `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important from container.task t
										INNER JOIN container.taskcontainer_task tct
										ON t.id = tct.task_id
										WHERE tct.taskcontainer_id in (SELECT id FROM container.taskcontainer where usergroup_id = $1)`
	sqlGetAllTasksByGroupIdAndImportant = `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important from public.task t
											INNER JOIN container.taskcontainer_task tct
											ON t.id = tct.task_id
											WHERE tct.taskcontainer_id in (SELECT id FROM container.taskcontainer where usergroup_id = $1) AND t.is_important = true`

	sqlCreateTask = `INSERT INTO container.task(id, name, description,type, created_at, updated_at, target_date, priority, category, is_completed, is_important)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	sqlCreateTaskForJoinTable   = `INSERT INTO container.taskcontainer_task(taskcontainer_id, task_id) VALUES ($1, $2)`
	sqlDeleteTaskForJoinTable   = `DELETE FROM container.taskcontainer_task WHERE task_id=$1`
	sqlDeleteTask               = `DELETE FROM container.task WHERE id=$1`
	sqlUpdateTask               = `UPDATE container.task SET name=$2, description=$3, updated_at=$4, target_date=$5, priority=$6, category=$7 WHERE id=$1`
	sqlUpdateTaskDoneField      = `UPDATE container.task SET is_completed=$1 WHERE id = $2;`
	sqlUpdateTaskImportantField = `UPDATE container.task SET is_important=$1 WHERE id = $2;`
)
