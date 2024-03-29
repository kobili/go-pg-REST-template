package db

import (
	"context"
	"database/sql"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntity struct {
	UserId    string
	Email     string
	FirstName string
	LastName  string
	Age       int32
}

type UserDetail struct {
	UserId    primitive.ObjectID `json:"userId" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	FirstName string             `json:"firstName" bson:"first_name"`
	LastName  string             `json:"lastName" bson:"last_name"`
	Age       int32              `json:"age" bson:"age"`
}

func GetUsers(mongoClient *mongo.Client, ctx context.Context) ([]UserDetail, error) {
	coll := getUserCollection(*mongoClient)

	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	var users []UserDetail
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	return users, nil
}

func GetUserById(mongoClient *mongo.Client, ctx context.Context, userId string) (*UserDetail, error) {
	coll := getUserCollection(*mongoClient)

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userId %s. Expected a hex value: %w", userId, err)
	}

	var user UserDetail
	filter := bson.D{{Key: "_id", Value: objectId}}

	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user with id %s: %w", userId, err)
	}

	return &user, nil
}

type UpdateUserPayload struct {
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"first_name"`
	LastName  string `json:"lastName" bson:"last_name"`
	Age       int32  `json:"age" bson:"age"`
}

func (data UpdateUserPayload) toBson() bson.D {
	return bson.D{
		{Key: "email", Value: data.Email},
		{Key: "first_name", Value: data.FirstName},
		{Key: "last_name", Value: data.LastName},
		{Key: "age", Value: data.Age},
	}
}

func CreateUser(mongoClient *mongo.Client, ctx context.Context, data UpdateUserPayload) (*UserDetail, error) {
	coll := getUserCollection(*mongoClient)

	// TODO: Will need to prevent duplicate emails here
	insertResult, err := coll.InsertOne(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user into database: %w", err)
	}

	documentId, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("mongo InsertOne returned an invalid ObjectID")
	}

	return &UserDetail{
		UserId:    documentId,
		Email:     data.Email,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Age:       data.Age,
	}, nil
}

func UpdateUser(mongoClient *mongo.Client, ctx context.Context, userId string, data UpdateUserPayload) (*UserDetail, error) {
	coll := getUserCollection(*mongoClient)

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userId %s. Expected a hex value: %w", userId, err)
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err = coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data.toBson()}})
	if err != nil {
		return nil, fmt.Errorf("failed to update user %s: %w", userId, err)
	}

	return &UserDetail{
		UserId:    objectId,
		Email:     data.Email,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Age:       data.Age,
	}, nil
}

func DeleteUser(db *sql.DB, ctx context.Context, userId string) error {
	_, err := db.ExecContext(
		ctx,
		`DELETE FROM users WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return fmt.Errorf("DeleteUser - Could not delete user: %w", err)
	}

	return nil
}
