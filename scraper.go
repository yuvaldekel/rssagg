package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yuvaldekel/rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Info: Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error: error fetching feeds: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error: Error marking feed as fetched: %v", err)
		return
	}

	RSSFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error: Error fetching feed: %v", err)

	}

	for _, item := range RSSFeed.Channel.Item {
		log.Printf("Info: Found post %v on feed %v", item.Title, feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(RSSFeed.Channel.Item))
}
