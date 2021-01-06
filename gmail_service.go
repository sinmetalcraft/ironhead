package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/sinmetalcraft/ironhead/erroreportings"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
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
func (s *GmailService) GetMessage(ctx context.Context, userID string, startHistoryID uint64, labelID string) (*gmail.Message, error) {
	res, err := s.gus.History.List(userID).StartHistoryId(startHistoryID).LabelId(labelID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Watch.userID=%s : %w", userID, err)
	}

	if len(res.History) < 1 {
		return nil, newErrInvalidMessage("invalid gmail.history.list", map[string]interface{}{"history": res}, nil)
	}
	if len(res.History[0].Messages) < 1 {
		return nil, newErrInvalidMessage("invalid gmail.history.message", map[string]interface{}{"history.message": res.History[0]}, nil)
	}

	msg, err := s.gus.Messages.Get(userID, res.History[0].Messages[0].Id).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Message.Get.userID=%s : %w", userID, err)
	}
	return msg, nil
}

// GetMessageList is 指定した historyID以降のmessageを取得する
func (s *GmailService) GetMessageList(ctx context.Context, userID string, startHistoryID uint64, labelID string) ([]*gmail.Message, error) {
	res, err := s.gus.History.List(userID).StartHistoryId(startHistoryID).LabelId(labelID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Watch.userID=%s : %w", userID, err)
	}

	var result []*gmail.Message
	for _, h := range res.History {
		for _, m := range h.Messages {
			msg, err := s.gus.Messages.Get(userID, m.Id).Context(ctx).Do()
			var gAPI *googleapi.Error
			if errors.As(err, &gAPI) {
				if gAPI.Code == 404 {
					fmt.Printf("message %s is not found\n", m.Id)
					continue
					// 404の時はログ出力してスルー
				} else {
					return nil, fmt.Errorf("failed GmailUserService.Message.Get.userID=%s : %w", userID, err)
				}
			} else if err != nil {
				return nil, fmt.Errorf("failed GmailUserService.Message.Get.userID=%s : %w", userID, err)
			}
			result = append(result, msg)
		}
	}
	return result, nil
}

// GetMessageList is 指定した historyID以降のmessageを取得する
func (s *GmailService) GetMessageList2(ctx context.Context, userID string, labelID string) ([]*gmail.Message, error) {
	res, err := s.gus.Messages.List(userID).LabelIds(labelID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed GmailUserService.Messages.List.userID=%s : %w", userID, err)
	}

	var result []*gmail.Message
	for _, m := range res.Messages {
		msg, err := s.gus.Messages.Get(userID, m.Id).Context(ctx).Do()
		var gAPI *googleapi.Error
		if errors.As(err, &gAPI) {
			if gAPI.Code == 404 {
				fmt.Printf("message %s is not found\n", m.Id)
				continue
				// 404の時はログ出力してスルー
			} else {
				return nil, fmt.Errorf("failed GmailUserService.Message.Get.userID=%s : %w", userID, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed GmailUserService.Message.Get.userID=%s : %w", userID, err)
		}
		result = append(result, msg)
	}
	return result, nil
}

// GetErrorReportingInfo is startHistoryIDがErrorReportingのIDだと信じてGetする
func (s *GmailService) GetErrorReportingInfo(ctx context.Context, userID string, startHistoryID uint64, labelID string) (*erroreportings.Info, error) {
	msg, err := s.GetMessage(ctx, userID, startHistoryID, labelID)
	if err != nil {
		return nil, err
	}

	if len(msg.Payload.Parts) < 2 {
		return nil, newErrInvalidMessage("invalid error reporting mail format", map[string]interface{}{"payload": msg.Payload}, nil)
	}

	var plainText []byte
	var htmlText []byte
	if msg.Payload.Parts[0].MimeType == "text/plain" {
		plainText, err = base64.URLEncoding.DecodeString(msg.Payload.Parts[0].Body.Data)
		if err != nil {
			return nil, fmt.Errorf("invalid error reporting mail format.%s", msg.Payload.Parts[0].Body.Data)
		}
	}

	if msg.Payload.Parts[1].MimeType == "text/html" {
		htmlText, err = base64.URLEncoding.DecodeString(msg.Payload.Parts[1].Body.Data)
		if err != nil {
			return nil, fmt.Errorf("invalid error reporting mail format.%s", msg.Payload.Parts[1].Body.Data)
		}
	}

	return erroreportings.Parse(string(plainText), string(htmlText)), nil
}
