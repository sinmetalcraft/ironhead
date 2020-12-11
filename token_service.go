package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
)

type TokenService struct {
	gcs    *storage.Client
	bucket string
}

func NewTokenService(ctx context.Context, bucket string, gcs *storage.Client) (*TokenService, error) {
	return &TokenService{
		gcs:    gcs,
		bucket: bucket,
	}, nil
}

func (s *TokenService) Close() error {
	return s.gcs.Close()
}

// SaveToken is RefreshToken を生成して Cloud Storage に保存する
func (s *TokenService) SaveToken(ctx context.Context, scope ...string) error {
	config, err := s.createConfig(ctx, scope...)
	if err != nil {
		return xerrors.Errorf("failed google.ConfigFromJSON err=%#v", err)
	}

	tok := s.getTokenFromWeb(config)
	if err := s.saveToken(ctx, tok); err != nil {
		return err
	}

	return nil
}

// CreateHTTPClient is Cloud Storage から Refresh Token 読み出して HTTP Client 作る
// scope は Refresh Token 作った時と同じのを指定しないといけない気がする
func (s *TokenService) CreateHTTPClient(ctx context.Context, scope ...string) (*http.Client, error) {
	config, err := s.createConfig(ctx, scope...)
	if err != nil {
		return nil, err
	}
	tok, err := s.getToken(ctx)
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, tok), nil
}

func (s *TokenService) createConfig(ctx context.Context, scope ...string) (*oauth2.Config, error) {
	clientSecret, err := s.getCredentialFile(ctx)
	if err != nil {
		log.Fatalf("failed getCredentialFile err=%#v", err)
	}

	config, err := google.ConfigFromJSON(clientSecret, scope...)
	if err != nil {
		log.Fatalf("failed google.ConfigFromJSON err=%#v", err)
	}
	return config, nil
}

func (s *TokenService) getCredentialFile(ctx context.Context) ([]byte, error) {
	r, err := s.gcs.Bucket(s.bucket).Object("client_secret.json").NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func (s *TokenService) getToken(ctx context.Context) (*oauth2.Token, error) {
	r, err := s.gcs.Bucket(s.bucket).Object("credential.json").NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		r.Close()
	}()
	tok := &oauth2.Token{}
	if err := json.NewDecoder(r).Decode(tok); err != nil {
		return nil, err
	}
	return tok, nil
}

// Saves a token to a file path.
func (s *TokenService) saveToken(ctx context.Context, token *oauth2.Token) (err error) {
	w := s.gcs.Bucket(s.bucket).Object("credential.json").NewWriter(ctx)
	defer func() {
		if errClose := w.Close(); errClose != nil {
			err = errClose // TODO うまいことWrapしたい気がする
		}
	}()
	if err := json.NewEncoder(w).Encode(token); err != nil {
		return err
	}

	return nil
}

func (s *TokenService) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
