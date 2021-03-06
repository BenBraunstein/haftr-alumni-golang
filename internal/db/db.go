package db

import (
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
)

const (
	usersCollectionName   = "users"
	alumnisCollectionName = "alumnis"
)

var (
	zeroInt64 = int64(0)
)

type InsertUserFunc func(u internal.User) error

type RetrieveUserByEmailFunc func(email string) (internal.User, error)

type RetrieveUserByIDFunc func(id string) (internal.User, error)

type ReplaceUserFunc func(u internal.User) error

type InsertAlumniFunc func(a internal.Alumni) error

type UpdateAlumniFunc func(id string, a internal.UpdateAlumniRequest) error

type RetrieveAlumniByIDFunc func(id string) (internal.Alumni, error)

type RetrieveAllAlumniFunc func(params pkg.QueryParams) ([]internal.Alumni, error)
