package workflow

import (
	"encoding/base64"
	"log"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/mapping"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AddUser adds a new user
func AddUser(insertUser db.InsertUserFunc, retrieveUserByEmail db.RetrieveUserByEmailFunc, provideTime time.EpochProviderFunc, genUUID uuid.GenV4Func) AddUserFunc {
	return func(req pkg.UserRequest) (pkg.User, error) {
		log.Printf("Adding new user with email=%v", req.Email)

		u, err := retrieveUserByEmail(req.Email)
		if err == nil {
			return pkg.User{}, errors.Errorf("workflow - user already exists with email=%v", u.Email)
		}

		pw, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		if err != nil {
			return pkg.User{}, errors.Wrapf(err, "workflow - unable to hash password, email=%v", req.Email)
		}

		user := mapping.ToDbUser(req, pw, genUUID, provideTime)
		if err := insertUser(user); err != nil {
			return pkg.User{}, errors.Wrapf(err, "workflow - unable to insert user into db, email=%v", req.Email)
		}

		return mapping.ToDTOUser(user), nil
	}
}

// LoginUser logs in a new user
func LoginUser(retrieveUserByEmail db.RetrieveUserByEmailFunc) LoginUserFunc {
	return func(req pkg.UserRequest) (pkg.User, error) {
		log.Printf("Logging in user with email=%v", req.Email)

		user, err := retrieveUserByEmail(req.Email)
		if err != nil {
			return pkg.User{}, errors.Wrapf(err, "workflow - unable to find user with email=%v", req.Email)
		}

		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - password is invalid")
		}

		return mapping.ToDTOUser(user), nil
	}
}

// AutoLoginUser auto logs in a user
func AutoLoginUser(retrieveUserById db.RetrieveUserByIDFunc) AutoLoginUserFunc {
	return func(encodedId string) (pkg.User, error) {
		log.Printf("Attempting to auto login user with token=%v", encodedId)

		bb, err := base64.StdEncoding.DecodeString(encodedId)
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to decode token")
		}

		user, err := retrieveUserById(string(bb))
		if err != nil {
			return pkg.User{}, errors.Wrap(err, "workflow - unable to find user with given token")
		}

		return mapping.ToDTOUser(user), nil
	}
}
