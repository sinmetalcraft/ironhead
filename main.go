package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	var (
		cmd = flag.String("cmd", "default", "command")
	)
	flag.Parse()
	fmt.Printf("cmd=%s\n", *cmd)

	if *cmd == "save-token" {
		fmt.Println("save-token")
		saveToken()
		os.Exit(0)
	}

	watchGmail()

	const addr = ":8080"
	fmt.Printf("Start Listen %s", addr)

	http.HandleFunc("/notify/gmail", GmailNotifyPubSubHandler)
	http.HandleFunc("/", helloWorldHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Ironhead")
}

func saveToken() {
	ctx := context.Background()

	gcs, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed storage.NewClient err=%#v", err)
	}
	ts, err := NewTokenService(ctx, "sinmetal-ironhead-config", gcs)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ts.Close(); err != nil {
			panic(err)
		}
	}()

	if err := ts.SaveToken(ctx, gmail.GmailMetadataScope); err != nil {
		panic(err)
	}
}

func watchGmail() {
	ctx := context.Background()

	gcs, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed storage.NewClient err=%#v", err)
	}
	ts, err := NewTokenService(ctx, "sinmetal-ironhead-config", gcs)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ts.Close(); err != nil {
			panic(err)
		}
	}()

	client, err := ts.CreateHTTPClient(ctx, gmail.GmailMetadataScope)
	if err != nil {
		panic(err)
	}

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
