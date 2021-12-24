package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"os"
)

func retrieveSubscriptionList(service *youtube.Service) ([]*youtube.Subscription, error) {
	call := service.Subscriptions.List([]string{"snippet"})
	call.Mine(true)
	call.MaxResults(maxResults)
	var subscriptions []*youtube.Subscription

	err := call.Pages(context.Background(), func(response *youtube.SubscriptionListResponse) error {
		subscriptions = append(subscriptions, response.Items[:]...)
		return nil
	})

	return subscriptions, err
}

func saveSubscriptionList(service *youtube.Service, filename string) error {
	subscriptions, err := retrieveSubscriptionList(service)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(subscriptions)
}

func readSubscriptionList(filename string) ([]*youtube.Subscription, error) {
	var subscriptions []*youtube.Subscription
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return subscriptions, err
	}

	err = json.Unmarshal(b, &subscriptions)

	return subscriptions, err
}

func subscriptionAdd(service *youtube.Service, channelId string) (*youtube.Subscription, error) {
	if channelId == "" {
		return nil, errors.New("channelId is required")
	}

	subscription := youtube.Subscription{
		Snippet: &youtube.SubscriptionSnippet{
			ResourceId: &youtube.ResourceId{
				ChannelId: channelId,
				Kind:      "youtube#channel",
			},
		},
	}

	call := service.Subscriptions.Insert([]string{"snippet"}, &subscription)
	return call.Do()
}

func subscriptionDel(service *youtube.Service, subscriptionId string) error {
	if subscriptionId == "" {
		return errors.New("subscriptionId is required")
	}

	call := service.Subscriptions.Delete(subscriptionId)
	return call.Do()
}
