package main

import (
	"fmt"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

// GmailWatchHandler is gmail.Watch を実行する
// 1日に1回cronで実行する
func GmailWatchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := gmailService.Watch(ctx, userID, &gmail.WatchRequest{
		TopicName: "projects/sinmetal-ironhead/topics/gmail",
		LabelIds:  []string{tbfErrorReportingLabelID},
	})
	if err != nil {
		fmt.Printf("failed GmailService.Watch. err=%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("HistoryID:%d,Expiration:%d\n", resp.HistoryID, resp.Expiration)
}
