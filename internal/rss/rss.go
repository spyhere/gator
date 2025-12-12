package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	request.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var rssFeed RSSFeed
	if err := xml.Unmarshal(bts, &rssFeed); err != nil {
		return &RSSFeed{}, err
	}
	return &rssFeed, nil
}
