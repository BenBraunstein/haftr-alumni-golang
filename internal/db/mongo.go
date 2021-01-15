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

func RetrieveUserByID(provideMongo *mongo.Database) RetrieveUserByIDFunc {
	return func(id string) (internal.User, error) {
		col := provideMongo.Collection(usersCollectionName)
		filter := bson.M{"id": id}

		var u internal.User
		if err := col.FindOne(context.Background(), filter).Decode(&u); err != nil {
			return internal.User{}, errors.Wrapf(err, "db - unable to find user with id=%v", id)
		}
		return u, nil
	}
}

func RetrieveAllAlumni(provideMongo *mongo.Database) RetrieveAllAlumniFunc {
	return func() ([]internal.Alumni, error) {
		col := provideMongo.Collection(alumnisCollectionName)
		filter := bson.M{}

		cur, err := col.Find(context.Background(), filter)
		if err != nil {
			return []internal.Alumni{}, errors.Wrap(err, "db - unable to find any alumnis")
		}
		var aa []internal.Alumni
		if err := cur.Decode(&aa); err != nil {
			return []internal.Alumni{}, errors.Wrap(err, "db - unable to decode alumnis response")
		}

		return aa, nil
	}
}
