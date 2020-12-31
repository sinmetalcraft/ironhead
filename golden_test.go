package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/storage"
)

const goldenBucket = "sinmetal-ironhead-golden"

func TestCreateGolden(t *testing.T) {
	ctx := context.Background()

	s := newGmailService(ctx)

	const userID = "sinmetal@sinmetalcraft.jp"
	const startHistoryID = 12158234
	const labelID = "Label_6698128523804588152"

	res, err := s.gus.History.List(userID).StartHistoryId(uint64(startHistoryID)).LabelId(labelID).Context(ctx).Do()
	if err != nil {
		t.Fatal(err)
	}
	//pp.Print(res)
	msg, err := s.gus.Messages.Get(userID, res.History[0].Messages[0].Id).Context(ctx).Do()
	if err != nil {
		t.Fatal(err)
	}
	writeGoldenFile(ctx, "error-reporting-mail-message.json", msg, t)
}

func writeGoldenFile(ctx context.Context, path string, body interface{}, t *testing.T) {
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := gcs.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	w := gcs.Bucket(goldenBucket).Object(path).NewWriter(ctx)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func readGoldenFile(ctx context.Context, path string, t *testing.T) []byte {
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := gcs.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	r, err := gcs.Bucket(goldenBucket).Object(path).NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
