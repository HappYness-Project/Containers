package repository

const (
	sqlGetAllUsergroups      = `SELECT * FROM container.usergroup`
	sqlGetById               = `SELECT * FROM container.usergroup WHERE id = $1`
	sqlGetUserGroupsByUserId = `SELECT ug.id, ug.name, ug.description, ug.type, ug.thumbnailurl, ug.is_active
								FROM container.usergroup ug
								INNER JOIN container.usergroup_user ugu
								ON ug.id = ugu.usergroup_id
								WHERE ugu.user_id = $1`

	sqlCreateUserGroup = `INSERT INTO container.usergroup(name, description, type, thumbnailurl, is_active)
							VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	sqlAddUserToUserGroup      = `INSERT INTO container.usergroup_user(usergroup_id, user_id) VALUES ($1, $2)`
	sqlAddUserToUserGroupAdmin = `INSERT INTO container.usergroup_user(usergroup_id, user_id, role) VALUES ($1, $2, 'admin')`
	sqlRemoveUserFromUserGroup = `DELETE FROM container.usergroup_user WHERE usergroup_id = $1 AND user_id = $2`
	sqlUpdateUserRoleInGroup   = `UPDATE container.usergroup_user SET role = $3 WHERE usergroup_id = $1 AND user_id = $2`

	sqlDeleteUserGroup = `DELETE FROM container.usergroup WHERE id = $1`
)
