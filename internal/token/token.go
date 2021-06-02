package token

import (
	"fmt"
	"os"
	gotime "time"

	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/common/uuid"
	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	userIdKey     = "user_id"
	adminKey      = "admin"
	expirationKey = "exp"
)

func getJwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func CreateUserToken(u internal.User, provideTime time.EpochProviderFunc) (string, error) {
	secret := getJwtSecret()

	exp, err := time.New(provideTime().ToISO8601().Val().Add(gotime.Hour * 36))
	if err != nil {
		return "", errors.Wrap(err, "token - unable to generate expiration")
	}

	m := jwt.MapClaims{}
	m["authorized"] = true
	m[userIdKey] = u.ID
	m[expirationKey] = exp.String()
	m[adminKey] = u.Admin

	return jwt.NewWithClaims(jwt.SigningMethodHS256, m).SignedString([]byte(secret))
}

func CheckUserToken(tokenString string, provideTime time.EpochProviderFunc) (uuid.V4, bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(getJwtSecret()), nil
	})
	if err != nil {
		return uuid.V4(""), false, err
	}

	m := token.Claims.(jwt.MapClaims)

	tTime, err := time.NewISO8601(m[expirationKey].(string))
	if err != nil {
		return uuid.V4(""), false, errors.Wrapf(err, "token - unable to retrieve expiration from JWT")
	}

	if !token.Valid || tTime.Val().Before(provideTime().ToISO8601().Val()) {
		return uuid.V4(""), false, fmt.Errorf("token - token is invalid or expired")
	}

	id := m[userIdKey].(string)
	admin := m[adminKey].(bool)

	return uuid.V4(id), admin, nil
}
