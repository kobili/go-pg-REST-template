package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"server/db"
)

type ReducedUserDetail struct {
	UserId    primitive.ObjectID `json:"userId"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}

func ListUsersHandler(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userEntities, err := db.GetUsers(client, req.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var users []ReducedUserDetail
		for _, userDetail := range userEntities {
			user := ReducedUserDetail{
				UserId:    userDetail.UserId,
				FirstName: userDetail.FirstName,
				LastName:  userDetail.LastName,
			}
			users = append(users, user)
		}

		resBody, err := json.Marshal(users)
		if err != nil {
			http.Error(w, fmt.Sprintf("ListUsersHandler - Error marshalling to json: %v", err), 500)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

type UserDetail struct {
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int32  `json:"age"`
}

func RetrieveUserHandler(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		user, err := db.GetUserById(client, req.Context(), userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("UserDetailhandler - Failed to retrieve user info: %v", err), 500)
			return
		}
		if user == nil {
			http.Error(w, fmt.Sprintf("User %s not found.", userId), 404)
		}

		resBody, err := json.Marshal(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("UserDetailHandler - Error marshalling to json: %v", err), 500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func CreateUserHandler(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var reqBody db.UpdateUserPayload
		err := json.NewDecoder(req.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("CreateUserHandler - Failed to decode request body: %v", err), 500)
			return
		}

		newUser, err := db.CreateUser(client, req.Context(), reqBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("CreateUserHandler - DB error: %v", err), 500)
			return
		}

		resBody, err := json.Marshal(newUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("CreateUserHandler - Error marshalling to json: %v", err), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func UpdateUserHandler(client *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		var reqBody db.UpdateUserPayload
		err := json.NewDecoder(req.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("UpdateUserHandler - Failed to decode request body: %v", err), 500)
			return
		}

		userDetail, err := db.UpdateUser(client, req.Context(), userId, reqBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("UpdateUserHandler - DB error: %v", err), 500)
			return
		}

		resBody, err := json.Marshal(userDetail)
		if err != nil {
			http.Error(w, fmt.Sprintf("UpdateUserHandler - Error marshalling to json: %v", err), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func DeleteUserHandler(sqlDB *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		db.DeleteUser(sqlDB, req.Context(), userId)

		w.WriteHeader(http.StatusNoContent)
	}

	return http.HandlerFunc(fn)
}
