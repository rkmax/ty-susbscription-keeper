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
	method         = flag.String("method", "playlists", "The APi Method to execute (default: playlists). playlists, subscriptions, rm-subscription, add-subscription")
	channelId      = flag.String("channel", "", "add-subscription | channelId to subscribe")
	subscriptionId = flag.String("subscription", "", "rm-subscription | subscriptionId to be removed")
)

const (
	maxResults = int64(50)
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

func subscriptionList(service *youtube.Service) {
	format := "%-43v\t%-24v\t%-20v\t%v\n"

	call := service.Subscriptions.List([]string{"snippet"})
	call.Mine(true)
	call.MaxResults(maxResults)

	fmt.Printf(format, "Id", "ResorceId", "Kind", "Title")
	err := call.Pages(context.Background(), func(response *youtube.SubscriptionListResponse) error {
		for _, subscription := range response.Items {
			fmt.Printf(format, subscription.Id, subscription.Snippet.ResourceId.ChannelId, subscription.Snippet.ResourceId.Kind, subscription.Snippet.Title)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Unable to call service: %v", err)
	}
}

func subscriptionAdd(service *youtube.Service) {
	if *channelId == "" {
		log.Fatalf("channelId is required")
	}

	subscription := youtube.Subscription{
		Snippet: &youtube.SubscriptionSnippet{
			ResourceId: &youtube.ResourceId{
				ChannelId: *channelId,
				Kind:      "youtube#channel",
			},
		},
	}

	call := service.Subscriptions.Insert([]string{"snippet"}, &subscription)
	response, err := call.Do()

	if err != nil {
		log.Fatalf("Unable to create subscription: %v", err)
	}

	fmt.Printf("Subscription created: %v", response.Id)
}

func subscriptionDel(service *youtube.Service) {
	if *subscriptionId == "" {
		log.Fatalf("resourceId is required")
	}

	call := service.Subscriptions.Delete(*subscriptionId)
	err := call.Do()

	if err != nil {
		log.Fatalf("Unable to remove subscription: %v", err)
	}

	fmt.Printf("Subscription deleted")
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
	case "add-subscription":
		subscriptionAdd(service)
	case "rm-subscription":
		subscriptionDel(service)
	}
}
