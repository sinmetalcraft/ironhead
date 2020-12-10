package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	watchGmail()

	const addr = ":8080"
	fmt.Printf("Start Listen %s", addr)

	http.HandleFunc("/", helloWorldHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Ironhead")
}

func watchGmail() {
	ctx := context.Background()

	clientSecret, err := getCredentialFile(ctx)
	if err != nil {
		log.Fatalf("failed getCredentialFile err=%#v", err)
	}

	config, err := google.ConfigFromJSON(clientSecret, gmail.GmailMetadataScope)
	if err != nil {
		log.Fatalf("failed google.ConfigFromJSON err=%#v", err)
	}

	client := getClient(config)
	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("failed gmail.NewService err=%+v", err)
	}

	s, err := NewGmailService(ctx, gmailService)
	if err != nil {
		log.Fatalf("failed NewGmailService err=%+v", err)
	}

	resp, err := s.Watch(ctx, "sinmetal@sinmetalcraft.jp", &gmail.WatchRequest{
		TopicName: "projects/sinmetal-ironhead/topics/gmail",
		LabelIds:  []string{"Label_6698128523804588152"},
	})
	if err != nil {
		log.Fatalf("failed Watch err=%+v", err)
	}
	fmt.Printf("%#v", resp)
}
