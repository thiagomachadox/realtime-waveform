package downloader

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/kkdai/youtube/v2"
)

var (
	youtubeRegex = regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.?be)/.+$`)
)

func DownloadAudio(url string) (io.ReadCloser, error) {
	if !youtubeRegex.MatchString(url) {
		return nil, fmt.Errorf("invalid YouTube URL")
	}

	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return nil, fmt.Errorf("error getting video info: %w", err)
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return nil, fmt.Errorf("no audio formats available")
	}

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return nil, fmt.Errorf("error getting audio stream: %w", err)
	}

	return stream, nil
}