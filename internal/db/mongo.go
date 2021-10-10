package db

import (
	"context"
	"regexp"
	"strings"

	"github.com/BenBraunstein/haftr-alumni-golang/internal"
	"github.com/BenBraunstein/haftr-alumni-golang/pkg"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ReplaceUser(provideMongo *mongo.Database) ReplaceUserFunc {
	return func(u internal.User) error {
		col := provideMongo.Collection(usersCollectionName)
		filter := bson.M{"id": u.ID}

		_, err := col.ReplaceOne(context.Background(), filter, u)
		if err != nil {
			return errors.Wrapf(err, "db - unable to replace user with id=%v", u.ID)
		}
		return nil
	}
}

func InsertAlumni(provideMongo *mongo.Database) InsertAlumniFunc {
	return func(a internal.Alumni) error {
		col := provideMongo.Collection(alumnisCollectionName)
		_, err := col.InsertOne(context.Background(), a)
		return err
	}
}

func UpdateAlumni(provideMongo *mongo.Database) UpdateAlumniFunc {
	return func(id string, a internal.UpdateAlumniRequest) error {
		col := provideMongo.Collection(alumnisCollectionName)
		filter := bson.M{"id": id}

		update := bson.D{
			{"$set", a},
		}
		_, err := col.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return errors.Wrapf(err, "db - unable to update alumniId=%v", id)
		}

		return nil
	}
}

func RetrieveAlumniByID(provideMongo *mongo.Database) RetrieveAlumniByIDFunc {
	return func(id string) (internal.Alumni, error) {
		col := provideMongo.Collection(alumnisCollectionName)
		filter := bson.M{"id": id}

		var a internal.Alumni
		if err := col.FindOne(context.Background(), filter).Decode(&a); err != nil {
			return internal.Alumni{}, errors.Wrapf(err, "db - unable to find alumni with id=%v", id)
		}

		return a, nil
	}
}

func ChangeAlumniPrivacy(provideMongo *mongo.Database) ChangeAlumniPrivacyFunc {
	return func(id string, isPublic bool) error {
		col := provideMongo.Collection(alumnisCollectionName)
		filter := bson.M{"id": id}

		update := bson.D{
			{"$set", bson.D{{
				"isPublic", isPublic,
			}}},
		}

		_, err := col.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return errors.Wrapf(err, "db - unable to update privacy status for alumniId=%v", id)
		}

		return nil
	}
}

func RetrieveAllAlumni(provideMongo *mongo.Database) RetrieveAllAlumniFunc {
	return func(params pkg.QueryParams, isAdmin bool) ([]internal.Alumni, error) {
		col := provideMongo.Collection(alumnisCollectionName)
		filter := bson.M{
			"firstname": bson.M{"$regex": primitive.Regex{Pattern: regexp.QuoteMeta(params.Firstname), Options: "i"}},
			"lastname":  bson.M{"$regex": primitive.Regex{Pattern: regexp.QuoteMeta(params.Lastname), Options: "i"}},
		}

		if !isAdmin {
			filter["isPublic"] = true
		}

		var skip int64
		if params.Page > 0 {
			skip = (params.Page - 1) * params.Limit
		}

		opts := options.FindOptions{
			Limit: &params.Limit,
			Skip:  &skip,
		}

		if params.Limit == (-1) {
			opts.Limit = &zeroInt64
			opts.Skip = &zeroInt64
		}

		ctx := context.Background()
		cur, err := col.Find(ctx, filter, &opts)
		if err != nil {
			return []internal.Alumni{}, errors.Wrap(err, "db - unable to find any alumnis")
		}

		defer cur.Close(ctx)
		aa := []internal.Alumni{}
		for cur.Next(ctx) {
			var a internal.Alumni
			if err := cur.Decode(&a); err != nil {
				return []internal.Alumni{}, errors.Wrap(err, "db - error decoding alumni")
			}
			aa = append(aa, a)
		}

		return aa, cur.Err()
	}
}

func RetrieveEmailTemplateByName(provideMongo *mongo.Database) RetrieveEmailTemplateByNameFunc {
	return func(name string) (internal.EmailTemplate, error) {
		col := provideMongo.Collection(emailTemplatesCollectionName)
		filter := bson.M{"name": name}

		var et internal.EmailTemplate
		if err := col.FindOne(context.Background(), filter).Decode(&et); err != nil {
			return internal.EmailTemplate{}, errors.Wrapf(err, "db - unable to find email template with name=%v", name)
		}

		return et, nil
	}
}
