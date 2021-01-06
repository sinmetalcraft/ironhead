# ironhead

Gmail Webhock

## Setup

[users.watch](https://developers.google.com/gmail/api/reference/rest/v1/users/watch) を設定する。
PubSubTopicには `gmail-api-push@system.gserviceaccount.com` を PubSub Publisherとして追加する。
GmailのLabelでフィルタリングするので、LabelごとにPubSubTopicを作った方がよい。

GmailのLabelIDを見るには [users.labels/list](https://developers.google.com/gmail/api/reference/rest/v1/users.labels/list) を実行する。