package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/xerrors"
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

// GmailNotifyPubSubHandler receives and processes a Pub/Sub push message.
func GmailNotifyPubSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var m PubSubMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var d NotifyData
	if err := json.Unmarshal(m.Message.Data, &d); err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	info, err := gmailService.GetErrorReportingInfo(ctx, d.EmailAddress, d.HistoryID, tbfErrorReportingLabelID)
	if xerrors.Is(err, ErrInvalidMessage) {
		log.Printf("invalid gmail history: %+v, %v\n", d, err)
		w.WriteHeader(http.StatusOK) // Retryしても完了しないので、諦めて終わる
		return
	} else if err != nil {
		log.Printf("gmailService.GetErrorReportingInfo: %v\n", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	fmt.Printf("%+v", info)
}
