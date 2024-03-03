package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Quiqui-dev/rssAggregator/internal/auth"
	"github.com/Quiqui-dev/rssAggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			log.Print("Auth error:", err)
			respondWithError(w, 401, "Auth error")
			return
		}

		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error retrieving user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
