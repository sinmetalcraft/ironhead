package main

import (
	"context"
	"fmt"

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
