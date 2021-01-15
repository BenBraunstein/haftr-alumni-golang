package workflow

import "github.com/BenBraunstein/haftr-alumni-golang/pkg"

type AddUserFunc func(req pkg.UserRequest) (pkg.User, error)

type LoginUserFunc func(req pkg.UserRequest) (pkg.User, error)

type AutoLoginUserFunc func(encodedId string) (pkg.User, error)
