package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

type MessageStore struct {
	ds *datastore.Client
}

func NewMessageStore(ctx context.Context, ds *datastore.Client) (*MessageStore, error) {
	return &MessageStore{
		ds,
	}, nil
}

type Message struct {
	// gmail.Message.ID
	ID string

	// 通知完了ステータス
	// 0:default
	// 1:done
	// 現状、通知完了した時にputするので、0の状態にDatastoreに保存されていることはない
	NotificationCompleteStatus int

	// 作成日時
	CreatedAt time.Time
}

func (s *MessageStore) Kind() string {
	return "Message"
}

func (s *MessageStore) Key(id string) *datastore.Key {
	return datastore.NameKey(s.Kind(), id, nil)
}

func (s *MessageStore) Put(ctx context.Context, model *Message) error {
	_, err := s.ds.Put(ctx, s.Key(model.ID), &model)
	if err != nil {
		return fmt.Errorf("failed Message.put() model is %+v : %w", model, err)
	}
	return nil
}

func (s *MessageStore) Get(ctx context.Context, id string) (*Message, error) {
	var model Message
	if err := s.ds.Get(ctx, s.Key(id), &model); err != nil {
		return nil, err
	}
	return &model, nil
}
