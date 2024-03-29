package db

import (
	"context"
	"database/sql"
	"fmt"
)

type UserEntity struct {
	UserId    string
	Email     string
	FirstName string
	LastName  string
	Age       int32
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

func GetUserById(db *sql.DB, ctx context.Context, userId string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT * FROM users WHERE user_id = $1", userId).Scan(
		&user.UserId, &user.Email, &user.FirstName, &user.LastName, &user.Age,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserById(%s): %w", userId, err)
	}

	return &user, nil
}

type UpdateUserPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int32  `json:"age"`
}

func CreateUser(db *sql.DB, ctx context.Context, data UpdateUserPayload) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, first_name, last_name, age)
		VALUES ($1, $2, $3, $4) RETURNING *`,
		data.Email, data.FirstName, data.LastName, data.Age,
	).Scan(&user.UserId, &user.Email, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		return nil, fmt.Errorf("CreateUser - Could not create user: %w", err)
	}

	return &user, nil
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
