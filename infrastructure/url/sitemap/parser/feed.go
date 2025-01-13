package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
)

// RssFeed is responsible for parsing HTML content and extracting URLs.
type RssFeed struct{}

// RSS represents the structure to decode.
type RSS struct {
	Channel struct {
		Items []struct {
			Link string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

// NewRssFeed creates and returns a new RssFeed instance.
func NewRssFeed() *RssFeed { return &RssFeed{} }

// Parse decodes the RSS feed content from the provided body.
func (f *RssFeed) Parse(body io.Reader) ([]string, error) {
	var rss RSS
	if err := xml.NewDecoder(body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("parse RSS feed: %w", err)
	}

	if len(rss.Channel.Items) == 0 {
		return nil, fmt.Errorf("no items found in RSS feed")
	}

	urls := make([]string, len(rss.Channel.Items))
	for i, item := range rss.Channel.Items {
		parsedURL, err := url.Parse(item.Link)
		if err != nil {
			fmt.Printf("invalid URL in RSS feed: %s, error: %v", item.Link, err)
			continue
		}
		urls[i] = parsedURL.String()
	}

	return urls, nil
}
