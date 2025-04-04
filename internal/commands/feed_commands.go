package commands

import (
	"context"
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	pq "github.com/lib/pq"

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

func HandlerBrowse(s *cnfg.State, cmd Command, user dbpk.User) error {
	postLimit := 2 // default limit

	if len(cmd.Args) > 1 {
		limit, err := strconv.Atoi(cmd.Args[1])
		if err != nil {
			return err
		}
		postLimit = limit
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), dbpk.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		log.Println(post.Title)
	}

	return nil
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

		parsedTime, err := parseTime(item.PubDate)
		if err != nil {
			log.Println(err)
			return
		}

		description := html.UnescapeString(item.Description)
		if len(description) > 255 {
			description = description[:255]
		}

		_, err = s.DB.CreatePost(context.Background(), dbpk.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       html.UnescapeString(item.Title),
			Url:         item.Link,
			Description: description,
			PublishedAt: parsedTime,
			FeedID:      nFeed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code != "23505" {
				log.Printf("Error occurred %v\n", err)
				return
			}
		}
	}

	log.Printf(
		"%v, Your posts have been saved.\n",
		strings.ToUpper(s.CurrState.CurrentUserName[:1])+s.CurrState.CurrentUserName[1:],
	)
}

func parseTime(publishedDate string) (time.Time, error) {
	dateLayouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		"2006-01-02 15:04",
		"02 Jan 2006",
	}

	for _, layout := range dateLayouts {
		parsedTime, err := time.Parse(layout, publishedDate)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("Could not parse date: %v", publishedDate)
}
