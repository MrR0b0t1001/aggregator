package commands

import (
	"context"
	"fmt"
	"html"
	"log"
	"time"

	rss "github.com/MrR0b0t1001/aggregator/internal/RSS"
	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"
	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"
	"github.com/google/uuid"
)

const feedURL = "https://www.wagslane.dev/index.xml"

func HandlerAgg(s *cnfg.State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", html.UnescapeString(feed.Channel.Title))
	fmt.Printf("%v\n", feed.Channel.Link)
	fmt.Printf("%v\n", html.UnescapeString(feed.Channel.Description))

	for _, item := range feed.Channel.Item {
		fmt.Printf("%v\n", html.UnescapeString(item.Title))
		fmt.Printf("%v\n", item.Link)
		fmt.Printf("%v\n", html.UnescapeString(item.Description))
	}

	return nil
}

func HandlerAddFeed(s *cnfg.State, cmd Command) error {
	user, err := s.DB.GetUser(context.Background(), s.CurrState.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	arg := dbpk.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("Feed not created: %w", err)
	}

	log.Println("Feed created successfully")

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
