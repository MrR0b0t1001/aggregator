package commands

import (
	"context"
	"fmt"
	"html"

	rss "github.com/MrR0b0t1001/aggregator/internal/RSS"
	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"
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
