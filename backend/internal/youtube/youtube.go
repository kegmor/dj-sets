package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type YouTube struct {
	apiKey	string
    client  *http.Client
}

type youtubeAPIResponse struct {
    Items []youtubeItem `json:"items"`
}

type youtubeItem struct {
    Snippet YouTubeResponse `json:"snippet"`
}

type YouTubeResponse struct {
    Title       string `json:"title"`
    ChannelTitle string `json:"channelTitle"`
}

type Option func (*YouTube)

func NewYouTube(apiKey string, client *http.Client, options ...Option) *YouTube {
    yt := &YouTube{
        apiKey:     apiKey, 
        client:     client,            
    }

    for _, option := range options {
        option(yt)
    }

    return yt
}

func (yt *YouTube) GetVideoDetails(videoId string) (*YouTubeResponse, error) {
    url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet", videoId, yt.apiKey)
    req, err := yt.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("youtube api request failed %w", err)
    }
    if req.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("youtube api returned status %d", req.StatusCode)
    }
    defer req.Body.Close()

    var yti youtubeAPIResponse
    err = json.NewDecoder(req.Body).Decode(&yti)
    if err != nil {
        return nil, fmt.Errorf("failed to decode youtube api call %w", err)
    }
    if len(yti.Items) == 0 {
        return nil, fmt.Errorf("no video found for id %s", videoId)
    }
    
    yts := YouTubeResponse {
        Title: yti.Items[0].Snippet.Title,
        ChannelTitle: yti.Items[0].Snippet.ChannelTitle,
    }
    return &yts, nil
}
