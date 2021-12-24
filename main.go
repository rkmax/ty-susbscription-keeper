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
	method            = flag.String("method", "subscriptions", "The APi Method to execute (default: playlists). playlists, subscriptions, rm-subscription, add-subscription")
	channelId         = flag.String("channel", "", "add-subscription | channelId to subscribe")
	subscriptionId    = flag.String("subscription", "", "rm-subscription | subscriptionId to be removed")
	subscriptionsFile = flag.String("subscription-file", "subscriptions.json", "save-subscription | Source or Destination file")
	launchWebServer   = flag.Bool("launch-web-server", true, "Launch web server for OAuth authentication flow")
)

const (
	maxResults             = int64(50)
	subscriptionListFormat = "%-43v\t%-24v\t%-20v\t%v\n"
)

func main() {
	flag.Parse()

	client := getClient(youtube.YoutubeScope, *launchWebServer)

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating Youtube client: %v", err)
	}

	switch *method {
	case "playlists":
		playList(service)
	case "subscriptions":
		listSubscriptionsCmd(service)
	case "add-subscription":
		addSubscriptionCmd(service)
	case "rm-subscription":
		deleteSubscriptionCmd(service)
	case "save-subscription":
		saveSubscriptionCmd(service)
	case "read-subscription":
		readSubscriptionCmd(service)
	}
}

func listSubscriptionsCmd(service *youtube.Service) {
	subscriptions, err := retrieveSubscriptionList(service)
	if err != nil {
		log.Fatalf("Unable to call service: %v", err)
	}
	for _, subscription := range subscriptions {
		fmt.Printf(subscriptionListFormat, subscription.Id, subscription.Snippet.ResourceId.ChannelId, subscription.Snippet.ResourceId.Kind, subscription.Snippet.Title)
	}
}

func addSubscriptionCmd(service *youtube.Service) {
	subscription, err := subscriptionAdd(service, *channelId)
	if err != nil {
		log.Fatalf("Unable to create subscription: %v", err)
	}

	fmt.Printf("Subscription created: %v", subscription.Id)
}

func deleteSubscriptionCmd(service *youtube.Service) {
	err := subscriptionDel(service, *subscriptionId)
	if err != nil {
		log.Fatalf("Unable to remove subscription: %v", err)
	}

	fmt.Printf("Subscription deleted")
}

func saveSubscriptionCmd(service *youtube.Service) {
	err := saveSubscriptionList(service, *subscriptionsFile)
	if err != nil {
		log.Fatalf("Unable to save subscriptions: %v", err)
	}

	fmt.Println("Subscription list save successfully")
}

func readSubscriptionCmd(service *youtube.Service) {
	subscriptions, err := readSubscriptionList(*subscriptionsFile)
	if err != nil {
		log.Fatalf("Unable to read subscription file %v", err)
	}

	for _, subscription := range subscriptions {
		fmt.Printf(subscriptionListFormat, subscription.Id, subscription.Snippet.ResourceId.ChannelId, subscription.Snippet.ResourceId.Kind, subscription.Snippet.Title)
	}
}
