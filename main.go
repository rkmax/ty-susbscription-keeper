package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	method = flag.String("method", "playlists", "The APi Method to execute (default: playlists)")
	mine   = flag.Bool("mine", false, "List playlist for authenticated user (default: false)")
)

const (
	maxResults = int64(50)
)

func playList(service *youtube.Service) {
	call := service.Playlists.List([]string{"snippet"})
	call.Mine(*mine)
	call.MaxResults(maxResults)

	err := call.Pages(context.Background(), func(response *youtube.PlaylistListResponse) error {
		for _, playlist := range response.Items {
			fmt.Println(playlist.Id, playlist.Snippet.Title)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Unable to call service: %v", err)
	}
}

func main() {
	flag.Parse()

	client := getClient(youtube.YoutubeScope)

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating Youtube client: %v", err)
	}

	switch *method {
	case "playlists":
		playList(service)
	case "subscriptions":
		subscriptionList(service)
	}
}

func subscriptionList(service *youtube.Service) {
	format := "%v\t%v\t%v\n"

	call := service.Subscriptions.List([]string{"snippet"})
	call.Mine(*mine)
	call.MaxResults(maxResults)

	fmt.Printf(format, "Id                                        ", "ResorceId               ", "Title")
	err := call.Pages(context.Background(), func(response *youtube.SubscriptionListResponse) error {
		for _, playlist := range response.Items {
			fmt.Printf(format, playlist.Id, playlist.Snippet.ResourceId.ChannelId, playlist.Snippet.Title)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Unable to call service: %v", err)
	}
}
