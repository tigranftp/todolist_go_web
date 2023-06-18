package db

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"site/config"
	"site/types"
	"sync"
	"time"
)

// Store ...
type Store struct {
	mtx *sync.Mutex
	db  *sqlx.DB
	dsn string
}

// New ...
func New(config *config.Config) *Store {
	return &Store{
		mtx: &sync.Mutex{}, // Mutex for blocking operations like update/insert
		dsn: config.DSN,
	}
}

// Open ...
func (s *Store) Open() error {
	db, err := sqlx.Open("sqlite3", s.dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// Query ...
func (s *Store) Query(querySTR string, args ...any) (*sql.Rows, error) {
	return s.db.Query(querySTR, args...)
}

// Exec ...
func (s *Store) Exec(querySTR string, args ...any) (sql.Result, error) {
	return s.db.Exec(querySTR, args...)
}

// GetUser ...
func (s *Store) GetUser(username string, passwordHash string) (*types.User, error) {
	row := s.db.QueryRow(GetUserQuery, username, passwordHash)
	user := new(types.User)
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Role)
	if err == sql.ErrNoRows {
		return nil, errors.New("login or password is incorrect")
	}
	return user, err
}

// RegisterUser ...
func (s *Store) RegisterUser(username, passwordHash, name string, role types.Role) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, err := s.db.Exec(InsertUserQuery, name, username, passwordHash, role)
	if err != nil && err.Error() == "UNIQUE constraint failed: user.username" {
		return errors.New("user with this username already exists")
	}
	return err
}

// AddUserListItem ...
func (s *Store) AddUserListItem(userID int64, taskname, description string, done bool) (int64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	var doneInt int
	if done {
		doneInt = 1
	}
	resultOfSQL, err := s.db.Exec(InsertToDoListItemQuery, taskname, description, userID, doneInt)
	if err != nil {
		return 0, err
	}
	return resultOfSQL.LastInsertId()
}

// GetUserTodoList ...
func (s *Store) GetUserTodoList(userID int64) (types.ToDoList, error) {
	var res types.ToDoList
	rows, err := s.db.Query(GetUserToDoList, userID)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		newItem := types.ListItem{}
		if err := rows.Scan(&newItem.Id, &newItem.Taskname, &newItem.Description,
			&newItem.Done, &newItem.CreationDate); err != nil {
			return res, err
		}
		res = append(res, newItem)
	}
	err = rows.Err()
	return res, err
}

// UpdateTodoListItem ...
func (s *Store) UpdateTodoListItem(taskID, userID int64, taskname, description string, done bool) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	var doneInt int
	if done {
		doneInt = 1
	}
	_, err := s.db.Exec(UpdateStatusOfListItem, taskname, description, doneInt, taskID, userID)
	return err
}

// DeleteTaskByID ...
func (s *Store) DeleteTaskByID(taskID, userID int64) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, err := s.db.Exec(DeleteItemFromListQuery, taskID, userID)
	return err
}

// UpdateRefreshToken ...
func (s *Store) UpdateRefreshToken(userID int64, refreshToken string, refreshEAT int64) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, err := s.db.Exec(InsertRefreshTokenQuery, refreshToken, refreshEAT, userID)
	return err
}

// GetUserByRefresh ...
func (s *Store) GetUserByRefresh(refreshToken string) (*types.User, error) {
	row := s.db.QueryRow(GetUserByRefreshTokenQuery, refreshToken)
	user := new(types.User)
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Role, &user.RefreshEAT)
	if err == sql.ErrNoRows {
		return nil, errors.New("wrong refresh token")
	}
	if time.Now().Unix() > user.RefreshEAT {
		return nil, errors.New("refresh token expired")
	}
	return user, err
}

// Close ...
func (s *Store) Close() {
	s.db.Close()
}
