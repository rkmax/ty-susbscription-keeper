package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"log"
)

func playList(service *youtube.Service) {
	format := "%-34v\t%v\n"
	call := service.Playlists.List([]string{"snippet"})
	call.Mine(true)
	call.MaxResults(maxResults)

	fmt.Printf(format, "Id", "Title")
	err := call.Pages(context.Background(), func(response *youtube.PlaylistListResponse) error {
		for _, playlist := range response.Items {
			fmt.Printf(format, playlist.Id, playlist.Snippet.Title)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Unable to call service: %v", err)
	}
}
