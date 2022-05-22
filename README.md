# soteria-functions

## Usage

#### select target project for deployment

```sh
$ make set-project GCP_PROJECT=<PROJECT ID> # default is black-stream-292507

# dev
$ make set-project GCP_PROJECT=black-stream-292507
$ make set-project GCP_PROJECT=soteria-production
```

#### local deploy

実質UpdatePriceHistoryのみデプロイすれば問題ない

```sh
# 先にupdate-btc-price-history-topicのtopicをpubsubで作成する
# CF deploy
$ CF_NAME="UpdatePriceHistory" CF_OPTIONS="--trigger-resource update-btc-price-history-topic" make deploy-fn-pubsub
```

#### local run

- 事前準備
  - 環境変数`GCP_PROJECT`を定義（GCPのプロジェクトID）
  - localからfirestoreへ接続できるサービスアカウントの秘密鍵のファイルパス`KEY_FILE_PATH`
- cmd/main.goを事前に修正すれば、localから起動することが可能

```go
// main.go

// 新しいCFを追加時は↓のように追加すれば、ローカル起動が可能（http triggerのみ）
http.HandleFunc("/new-cloud-function", functions.NewCloudFunction)
```

```sh
# 事前準備
$ export KEY_FILE_PATH=./.key/service_account.json

# 実行
$ go run cmd/main.go
```
