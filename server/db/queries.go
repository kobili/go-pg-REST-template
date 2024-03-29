package db

import (
	"context"
	"database/sql"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const MONGO_DB_NAME = "Go_REST"
const MONGO_COLLECTION_USER = "users"

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

func GetUsers(db *sql.DB, ctx context.Context) ([]UserEntity, error) {
	var users []UserEntity

	rows, err := db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserEntity
		if err := rows.Scan(&user.UserId, &user.Email, &user.FirstName, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("GetUsers: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUsers: %v", err)
	}
	return users, nil
}

func GetUserById(mongoClient *mongo.Client, ctx context.Context, userId string) (*UserDetail, error) {
	coll := mongoClient.Database(MONGO_DB_NAME).Collection(MONGO_COLLECTION_USER)

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

func CreateUser(mongoClient *mongo.Client, ctx context.Context, data UpdateUserPayload) (*UserDetail, error) {
	coll := mongoClient.Database(MONGO_DB_NAME).Collection(MONGO_COLLECTION_USER)

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

func UpdateUser(db *sql.DB, ctx context.Context, userId string, data UpdateUserPayload) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`UPDATE users
		SET email = $1,
			first_name = $2,
			last_name = $3,
			age = $4
		WHERE user_id = $5
		RETURNING *`,
		data.Email, data.FirstName, data.LastName, data.Age, userId,
	).Scan(&user.UserId, &user.Email, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser - Could not update user: %w", err)
	}

	return &user, nil
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
