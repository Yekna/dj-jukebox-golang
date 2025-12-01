package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type YouTubeSearchResult struct {
	VideoID   string `json:"video_id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"youtube_url"`
}

func SearchYouTube(query, apiKey string, maxResults int) ([]YouTubeSearchResult, error) {
	if apiKey == "" {
		return []YouTubeSearchResult{}, nil
	}

	encodedQuery := url.QueryEscape(query)

	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&videoCategoryId=10&q=%s&maxResults=%d&key=%s",
		encodedQuery, maxResults, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var yt struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title      string `json:"title"`
				Thumbnails struct {
					High struct {
						URL string `json:"url"`
					} `json:"high"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&yt); err != nil {
		return nil, err
	}

	results := []YouTubeSearchResult{}

	for _, item := range yt.Items {
		log.Println(item)
		if item.ID.VideoID == "" {
			continue
		}
		results = append(results, YouTubeSearchResult{
			VideoID:   item.ID.VideoID,
			Title:     item.Snippet.Title,
			Thumbnail: item.Snippet.Thumbnails.High.URL,
			URL:       "https://www.youtube.com/watch?v=" + item.ID.VideoID,
		})
	}

	return results, nil
}
