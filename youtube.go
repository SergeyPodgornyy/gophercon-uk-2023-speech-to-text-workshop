package main

import (
	"context"
	"fmt"

	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

func YoutubeDL(ytlink string) error {
	client := youtube.Client{
		Debug: true,
	}
	ctx := context.Background()

	video, err := client.GetVideoContext(ctx, ytlink)
	if err != nil {
		return fmt.Errorf("getting video: %w", err)
	}
	d := downloader.Downloader{Client: client, OutputDir: "./"}

	return d.DownloadComposite(ctx, "", video, "hd1080", "mp4")
}
