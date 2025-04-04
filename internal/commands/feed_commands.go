package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	rss "github.com/MrR0b0t1001/aggregator/internal/RSS"
	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"
	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"
	"github.com/google/uuid"
)

const (
	feedURL = "https://www.wagslane.dev/index.xml"
)

func HandlerAgg(s *cnfg.State, cmd Command) error {
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func HandlerAddFeed(s *cnfg.State, cmd Command, user dbpk.User) error {
	feed, err := s.DB.CreateFeed(context.Background(), dbpk.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		log.Println("Hello ? ")
		return fmt.Errorf("Feed not created: %w", err)
	}

	log.Println("Feed created successfully")

	_, err = s.DB.CreateFeedFollow(context.Background(), dbpk.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		fmt.Println("Hello hello")
		return err
	}

	fmt.Printf("User: %v is now following %v", user.Name, feed.Name)
	fmt.Println(feed.Name, feed.Url)
	return nil
}

func HandlerFeeds(s *cnfg.State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%v\n%v\n", feed.Name, feed.Name_2)
	}

	return nil
}

func HandlerFollow(s *cnfg.State, cmd Command, user dbpk.User) error {
	feed, err := s.DB.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), dbpk.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%v\n%v\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func HandlerFollowing(s *cnfg.State, cmd Command, user dbpk.User) error {
	following, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range following {
		fmt.Printf("%v\n", follow.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *cnfg.State, cmd Command, user dbpk.User) error {
	feed, err := s.DB.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	if err := s.DB.UnfollowFeed(context.Background(), dbpk.UnfollowFeedParams{UserID: user.ID, FeedID: feed.ID}); err != nil {
		return err
	}

	log.Printf("Feed : %v Unfollowed\n", feed.Name)

	return nil
}

func scrapeFeeds(s *cnfg.State) {
	nFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Unable to get next feed: %v", err)
		return
	}

	feedItems, err := rss.FetchFeed(context.Background(), nFeed.Url)
	if err != nil {
		log.Printf("Error fetching feed with the specified url, %v", err)
		return
	}

	if err := s.DB.MarkFeedFetched(context.Background(), nFeed.ID); err != nil {
		log.Printf("Error during marking %v", err)
		return
	}

	for _, item := range feedItems.Channel.Item {
		fmt.Println(item.Title)
	}
}
