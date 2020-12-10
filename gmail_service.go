package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type GmailService struct {
	gus *gmail.UsersService
}

func NewGmailService(ctx context.Context, gmailService *gmail.Service) (*GmailService, error) {
	gmailUserService := gmail.NewUsersService(gmailService)
	return &GmailService{
		gus: gmailUserService,
	}, nil
}

type WatchResponse struct {
	HistoryID  int64
	Expiration int64
}

func (s *GmailService) Watch(ctx context.Context, userID string, req *gmail.WatchRequest) (*WatchResponse, error) {
	res, err := s.gus.Watch(userID, req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Watch.userID=%s,req=%#v : %w", userID, req, err)
	}
	return &WatchResponse{
		HistoryID:  int64(res.HistoryId),
		Expiration: res.Expiration,
	}, nil
}

func getCredentialFile(ctx context.Context) ([]byte, error) {
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := gcs.Close(); err != nil {
			fmt.Printf("failed gcs.Close() err=%v", err)
		}
	}()

	r, err := gcs.Bucket("sinmetal-ironhead-config").Object("client_secret.json").NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	unescapeAuthCode, err := url.QueryUnescape(authCode)
	if err != nil {
		log.Fatalf("url.QueryUnescape: %v", err)
	}

	tok, err := config.Exchange(context.Background(), unescapeAuthCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
