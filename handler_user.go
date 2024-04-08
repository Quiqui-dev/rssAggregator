package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Quiqui-dev/rssAggregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		DisplayName  string `json:"display_name"`
		Password     string `json:"password"`
		EmailAddress string `json:"email_address"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		CreateAt:     time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		DisplayName:  params.DisplayName,
		EmailAddress: params.EmailAddress,
		Password:     params.Password,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, 201, databseUserToUser(user))
}

func (apiConfig *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password     string `json:"password"`
		EmailAddress string `json:"email_address"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	user, err := apiConfig.DB.LogIn(r.Context(), database.LogInParams{
		EmailAddress: params.EmailAddress,
		Password:     params.Password,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error logging in user: %v", err))
		return
	}

	respondWithJSON(w, 200, databseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, 200, databseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching posts for user: %v", err))
		return
	}

	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
