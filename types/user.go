package types

type User struct {
	Id         int64  `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Name       string `json:"name" db:"name"`
	Role       Role   `json:"role" db:"role"`
	RefreshEAT int64  `json:"refreshEAT"`
}
