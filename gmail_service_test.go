package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/k0kubun/pp"
	"google.golang.org/api/gmail/v1"
)

func TestPubSubPullGmailDebug(t *testing.T) {
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

func TestSetWatch(t *testing.T) {
	ctx := context.Background()

	gmailService := newGmailService(ctx)
	resp, err := gmailService.Watch(ctx, userID, &gmail.WatchRequest{
		TopicName: "projects/sinmetal-ironhead/topics/github",
		LabelIds:  []string{githubNotifyLabelID},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("HistoryID:%d,Expiration:%d\n", resp.HistoryID, resp.Expiration)
}

func TestPubSubPullGitHubDebug(t *testing.T) {
	ctx := context.Background()

	gmailService := newGmailService(ctx)

	pubsubClient, err := pubsub.NewClient(ctx, "sinmetal-ironhead")
	if err != nil {
		t.Fatal(err)
	}

	debugSub := pubsubClient.Subscription("github-debug")
	err = debugSub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		time.Sleep(3 * time.Second)
		fmt.Printf("%+v\n", message.Attributes)
		fmt.Printf("%s\n", string(message.Data))
		var d NotifyData
		if err := json.Unmarshal(message.Data, &d); err != nil {
			t.Fatal(err)
		}

		got, err := gmailService.GetMessageList(ctx, userID, d.HistoryID, githubNotifyLabelID)
		if err != nil {
			t.Fatal(err)
		}
		for _, msg := range got {
			pp.Printf("historyID:%s,msg.historyID:%s,Snippet:%s\n", d.HistoryID, msg.HistoryId, msg.Snippet)
		}

		message.Ack()
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGmailService_GetMessageList(t *testing.T) {
	ctx := context.Background()

	gmailService := newGmailService(ctx)

	got, err := gmailService.GetMessageList(ctx, userID, 12204190, githubNotifyLabelID)
	if err != nil {
		t.Fatal(err)
	}
	for _, msg := range got {
		pp.Printf("historyID:%s,msg.historyID:%s,Snippet:%s\n", 12204190, msg.HistoryId, msg.Snippet)
	}
}

func TestGmailService_GetMessageList2(t *testing.T) {
	ctx := context.Background()

	gmailService := newGmailService(ctx)

	got, err := gmailService.GetMessageList2(ctx, userID, githubNotifyLabelID)
	if err != nil {
		t.Fatal(err)
	}
	for _, msg := range got {
		pp.Printf("%s:%s\n", msg.Id, msg.Snippet)
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
