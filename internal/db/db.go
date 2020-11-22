package db

import "github.com/BenBraunstein/haftr-alumni-golang/internal"

const (
	usersCollectionName   = "users"
	alumnisCollectionName = "alumnis"
)

type InsertUserFunc func(u internal.User) error

type RetrieveUserByEmailFunc func(email string) (internal.User, error)

type RetrieveAllAlumniFunc func() ([]internal.Alumni, error)
