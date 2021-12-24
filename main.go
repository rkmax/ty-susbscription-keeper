package main

import (
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	method          = flag.String("method", "subscriptions", "The APi Method to execute (default: playlists). playlists, subscriptions, rm-subscription, add-subscription")
	channelId       = flag.String("channel", "", "add-subscription | channelId to subscribe")
	subscriptionId  = flag.String("subscription", "", "rm-subscription | subscriptionId to be removed")
	launchWebServer = flag.Bool("launch-web-server", true, "Launch web server for OAuth authentication flow")
)

const (
	maxResults = int64(50)
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
		subscriptionList(service)
	case "add-subscription":
		subscriptionAdd(service)
	case "rm-subscription":
		subscriptionDel(service)
	}
}
