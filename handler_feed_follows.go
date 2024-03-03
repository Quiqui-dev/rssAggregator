package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Quiqui-dev/rssAggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	feed, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreateAt:  time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feed))
}

func (apiConfig *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feed_follows, err := apiConfig.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving feed_follows for the given user: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feed_follows))
}

func (apiConfig *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feed_follow_id_str := chi.URLParam(r, "feed_follow_id")
	feed_follow_id, err := uuid.Parse(feed_follow_id_str)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing feed follow id: %v", err))
		return
	}

	err = apiConfig.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     feed_follow_id,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow id: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
