package repository

const (
	sqlGetAllUsers     = `SELECT * FROM container.user`
	sqlGetUserByUserId = `SELECT id, user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id
						 FROM container.user
						 WHERE user_id = $1`
	sqlGetUserByEmail = `SELECT id, user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id
							FROM container.user
							WHERE email = $1`
	sqlGetUserByUsername = `SELECT id, user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id
							FROM container.user
							WHERE username = $1`
	sqlGetUsersByGroupId = `SELECT id, u.user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id from container.user u
							INNER JOIN container.usergroup_user ugu
							ON u.id = ugu.user_id
							WHERE ugu.usergroup_id = $1`
	sqlCreateUser = `INSERT INTO container.user(user_id, username, first_name, last_name, email, is_active, created_at, updated_at, default_group_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	sqlUpdateUser = `UPDATE container.user
							SET first_name=$2, last_name=$3, email=$4, default_group_id=$5, updated_at=$6
							WHERE id= $1`
	sqlGetUsersByGroupIdWithRoles = `SELECT u.id, u.user_id, u.username, u.first_name, u.last_name, u.email, u.is_active, u.created_at, u.updated_at, u.default_group_id, ugu.role, ugu.joined_at
									FROM container.user u
									INNER JOIN container.usergroup_user ugu ON u.id = ugu.user_id
									WHERE ugu.usergroup_id = $1`
	sqlGetUserRoleInGroup = `SELECT ugu.role FROM container.usergroup_user ugu
							INNER JOIN container.user u ON u.id = ugu.user_id
							WHERE u.user_id = $1 AND ugu.usergroup_id = $2`
)
