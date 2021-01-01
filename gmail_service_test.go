package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/k0kubun/pp"
)

func TestPubSubPull(t *testing.T) {
	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, "sinmetal-ironhead")
	if err != nil {
		t.Fatal(err)
	}

	debugSub := pubsubClient.Subscription("gmail-debug")
	err = debugSub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		fmt.Printf("%+v\n", message.Attributes)
		fmt.Printf("%s\n", string(message.Data))
		var d NotifyData
		if err := json.Unmarshal(message.Data, &d); err != nil {
			t.Fatal(err)
		}

	})
	if err != nil {
		t.Fatal(err)
	}
}

// 試しにErrorReportingのメッセージを取得してみる
func TestGmailService_GetMessage(t *testing.T) {
	ctx := context.Background()

	s := newGmailService(ctx)
	msg, err := s.GetMessage(ctx, "sinmetal@sinmetalcraft.jp", 12158234, tbfErrorReportingLabelID)
	if err != nil {
		t.Fatal(err)
	}
	pp.Print(msg)
}
