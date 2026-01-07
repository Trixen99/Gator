package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	RSSChannel Channel `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		var errored RSSFeed
		return &errored, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		var errored RSSFeed
		return &errored, err
	}

	defer resp.Body.Close()

	bytedata, err := io.ReadAll(resp.Body)
	if err != nil {
		var errored RSSFeed
		return &errored, err
	}

	var feed RSSFeed

	err = xml.Unmarshal(bytedata, &feed)
	if err != nil {
		return &feed, err
	}

	return &feed, nil

}

func (r RSSFeed) UnescapeStrings() {
	r.RSSChannel.Title = html.UnescapeString(r.RSSChannel.Title)
	r.RSSChannel.Description = html.UnescapeString(r.RSSChannel.Description)
	for i, item := range r.RSSChannel.Item {
		r.RSSChannel.Item[i].Title = html.UnescapeString(item.Title)
		r.RSSChannel.Item[i].Description = html.UnescapeString(item.Description)
	}
}
