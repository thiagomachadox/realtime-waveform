package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/kkdai/youtube/v2"
)

var (
	youtubeRegex = regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.?be)/.+$`)
)

func DownloadAudio(url string) (io.ReadCloser, error) {
	log.Printf("Attempting to download audio from URL: %s", url)

	if !youtubeRegex.MatchString(url) {
		return nil, fmt.Errorf("invalid YouTube URL: %s", url)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	client := youtube.Client{
		HTTPClient: httpClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Getting video info...")
	video, err := client.GetVideoContext(ctx, url)
	if err != nil {
		log.Printf("Error getting video info: %v", err)
		return nil, fmt.Errorf("error getting video info: %w", err)
	}

	log.Println("Video info retrieved successfully")

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return nil, fmt.Errorf("no audio formats available")
	}

	log.Println("Getting audio stream...")
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		log.Printf("Error getting audio stream: %v", err)
		return nil, fmt.Errorf("error getting audio stream: %w", err)
	}

	log.Println("Audio stream obtained successfully")

	return stream, nil
}
