package db

import (
	"context"
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(provideMongo *mongo.Database) InsertUserFunc {
	return func(u internal.User) error {
		col := provideMongo.Collection(usersCollectionName)
		_, err := col.InsertOne(context.Background(), u)
		return err
	}
}

func RetrieveUserByEmail(provideMongo *mongo.Database) RetrieveUserByEmailFunc {
	return func(email string) (internal.User, error) {
		col := provideMongo.Collection(usersCollectionName)
		filter := bson.M{"email": strings.ToLower(email)}

		var u internal.User
		if err := col.FindOne(context.Background(), filter).Decode(&u); err != nil {
			return internal.User{}, errors.Wrapf(err, "db - unable to find user with email=%v", email)
		}
		return u, nil
	}
}
