package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// NotifyData is Gmail Notify で飛んでくる中身
type NotifyData struct {
	EmailAddress string `json:"emailAddress"`
	HistoryID    uint64 `json:"historyId"`
}

// GmailTBFErrorReportingNotifyPubSubHandler receives and processes a Pub/Sub push message.
func GmailTBFErrorReportingNotifyPubSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	msgs, err := gmailService.GetMessageList(ctx, userID, tbfErrorReportingLabelID)
	if err != nil {
		log.Printf("failed gmail.MessageList() %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, msg := range msgs {
		v, err := messageStore.Get(ctx, msg.Id)
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			// noop
		} else if err != nil {
			log.Printf("messageStore.Get: %v\n", err)
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
			return
		}
		if v.NotificationCompleteStatus == 1 {
			// すでに通知済み
			continue
		}

		// TODO Notify
		info, err := gmailService.GetErrorReportingInfo(ctx, userID, tbfErrorReportingLabelID)
		if errors.Is(err, ErrInvalidMessage) {
			log.Printf("invalid gmail.message: %s, %v\n", msg.Id, err)
			continue // Retryしても完了しないので、諦める
		} else if err != nil {
			log.Printf("gmailService.GetErrorReportingInfo: %v\n", err)
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
			return
		}
		fmt.Printf("%+v", info)

		err = messageStore.Put(ctx, &Message{
			ID:                         v.ID,
			NotificationCompleteStatus: 1,
		})
		if err != nil {
			log.Printf("messageStore.Put: %v\n", err)
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
			return
		}
	}
}
