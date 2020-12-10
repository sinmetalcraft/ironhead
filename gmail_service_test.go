package main

import (
	"context"
	"testing"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func TestGmailService_Watch(t *testing.T) {
	ctx := context.Background()

	clientSecret, err := getCredentialFile(ctx)
	if err != nil {
		t.Fatal(err)
	}

	config, err := google.ConfigFromJSON(clientSecret, gmail.GmailMetadataScope)
	if err != nil {
		t.Fatal(err)
	}

	token := getTokenFromWeb(config)
	client := config.Client(ctx, token)
	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewGmailService(ctx, gmailService)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := s.Watch(ctx, "sinmetal@sinmetalcraft.jp", &gmail.WatchRequest{
		TopicName: "projects/sinmetal-ironhead/topics/gmail",
		LabelIds:  []string{"tbf-stackdriver-notifications"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", resp)
}
