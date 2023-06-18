package types

type Role int

const (
	AdminRole     Role = iota
	UserRole           = iota
	ModeratorRole      = iota
)
