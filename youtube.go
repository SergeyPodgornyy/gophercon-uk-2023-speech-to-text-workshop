package main

import (
	"context"
	"fmt"

	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

func YoutubeDL(ytID string) (string, error) {
	client := youtube.Client{
		Debug: true,
	}
	ctx := context.Background()

	video, err := client.GetVideoContext(ctx, ytID)
	if err != nil {
		return "", fmt.Errorf("getting video: %w", err)
	}
	d := downloader.Downloader{Client: client, OutputDir: "/data"}
	outputfile := ytID + ".mp4"
	err = d.DownloadComposite(ctx, outputfile, video, "hd1080", "mp4")
	if err != nil {
		return "", fmt.Errorf("downloading video: %w", err)
	}

	return ytID, nil
}
