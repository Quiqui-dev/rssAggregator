package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Quiqui-dev/rssAggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	//godotenv.Load(".env")

	portStr := os.Getenv("PORT")

	if portStr == "" {
		log.Fatal("Port is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")

	log.Println("DB_URL:", dbUrl)
	if dbUrl == "" {
		log.Fatal("DB URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Could not open connection to database:", err)
	}

	dbConn := database.New(conn)
	apiCfg := apiConfig{
		DB: dbConn,
	}

	go startScraping(
		dbConn,
		10,
		time.Minute,
	)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/login", apiCfg.handleLoginUser)

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follow/{feed_follow_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	log.Printf("Server starting on port %v", portStr)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
