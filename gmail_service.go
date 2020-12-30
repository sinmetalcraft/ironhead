package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/sinmetalcraft/ironhead/erroreportings"
	"google.golang.org/api/gmail/v1"
)

var GmailServiceScope []string

func init() {
	GmailServiceScope = append(GmailServiceScope, gmail.GmailReadonlyScope)
}

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

// Watch is set watch request
// PubSub Topicに指定したLabelのメールが来たら、通知を行うWatchを設定する
// 1日1回呼ぶのがよいらしい
// https://developers.google.com/gmail/api/guides/push#watch_request
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

// GetMessage is historyIDから1件のメール・メッセージを取得する
func (s *GmailService) GetMessage(ctx context.Context, userID string, startHistoryID int64, labelID string) (*gmail.Message, error) {
	res, err := s.gus.History.List(userID).StartHistoryId(uint64(startHistoryID)).LabelId(labelID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Watch.userID=%s,req=%#v : %w", userID, err)
	}

	msg, err := s.gus.Messages.Get(userID, res.History[0].Messages[0].Id).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Message.Get.userID=%s,req=%#v : %w", userID, err)
	}
	return msg, nil
}

// GetErrorReportingInfo is startHistoryIDがErrorReportingのIDだと信じてGetする
func (s *GmailService) GetErrorReportingInfo(ctx context.Context, userID string, startHistoryID int64, labelID string) (*erroreportings.Info, error) {
	msg, err := s.GetMessage(ctx, userID, startHistoryID, labelID)
	if err != nil {
		return nil, err
	}

	if len(msg.Payload.Parts) < 2 {
		return nil, fmt.Errorf("invalid error reporting mail format")
	}

	var plainText []byte
	var htmlText []byte
	if msg.Payload.Parts[0].MimeType == "text/plain" {
		plainText, err = base64.URLEncoding.DecodeString(msg.Payload.Parts[0].Body.Data)
		if err != nil {
			return nil, fmt.Errorf("invalid error reporting mail format")
		}
	}

	if msg.Payload.Parts[1].MimeType == "text/html" {
		htmlText, err = base64.URLEncoding.DecodeString(msg.Payload.Parts[1].Body.Data)
		if err != nil {
			return nil, fmt.Errorf("invalid error reporting mail format")
		}
	}

	return erroreportings.Parse(string(plainText), string(htmlText)), nil
}
