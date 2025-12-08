package repository

const (
	sqlGetAllContainers       = `SELECT id,name,description,is_active,usergroup_id FROM container.taskcontainer`
	sqlGetById                = `SELECT id,name,description,is_active,usergroup_id FROM container.taskcontainer WHERE id = $1`
	sqlGetContainersByGroupId = `SELECT id,name,description,is_active,usergroup_id FROM container.taskcontainer WHERE usergroup_id = $1`
	sqlCreateContainer        = `INSERT INTO container.taskcontainer(id, name, description, is_active, activity_level, type, usergroup_id)
								VALUES ($1,$2,$3,$4,$5,$6,$7);`
	sqlDeleteContainer              = `DELETE FROM container.taskcontainer WHERE id = $1;`
	sqlDeleteContainerByUsergroupId = `DELETE FROM container.taskcontainer WHERE usergroup_id = $1;`
)
