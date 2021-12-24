package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"log"
)

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
