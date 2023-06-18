package db

const (
	GetUserQuery = `SELECT id, username, name, role FROM user
						WHERE Username = (?) AND Password = (?)`
	// InsertUserQuery - for registration
	InsertUserQuery = `INSERT INTO user (Name, Username, Password, Role)
						VALUES ((?),(?),(?),(?))`
	InsertToDoListItemQuery = `INSERT INTO todolist_items (taskname, creation_date, description, user_id, done)
						VALUES ((?),datetime("now"),(?),(?), (?))
`
	GetUserToDoList = `
SELECT id, taskname, description, done, creation_date FROM todolist_items
						WHERE user_id = (?)
ORDER BY
    datetime(creation_date) ASC
    `

	UpdateStatusOfListItem = `
	UPDATE todolist_items
	SET taskname = (?),
		description = (?),
		done = (?)
	WHERE
		id = (?) and user_id = (?)
`
	DeleteItemFromListQuery = `
	DELETE FROM todolist_items
	WHERE id = (?) and user_id = (?)
`
	InsertRefreshTokenQuery = `
	UPDATE user
	SET refresh = (?),
		refreshEAT = (?)
	WHERE
		id = (?)
`
	GetUserByRefreshTokenQuery = `
	SELECT id, username, name, role, refreshEAT FROM user
		WHERE refresh = (?)
`
)
