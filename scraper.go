package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Quiqui-dev/rssAggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrenty int,
	timeBetweenRequests time.Duration,
) {

	// start scraping
	log.Printf("Scraping on %v Go-routines every %s duration", concurrenty, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {

		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrenty),
		)

		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(wg, db, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)

	if err != nil {
		log.Println("Error fetching feed: err", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Printf("Could not parse date %v with error %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreateAt:    time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {

			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}

			log.Println("Failed to create post")
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
