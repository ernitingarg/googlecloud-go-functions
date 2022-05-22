include .env
export

run:
	go run cmd/main.go

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test ./...

.PHONY: set-project
set-project:
	gcloud config set project $(GCP_PROJECT)

.PHONY: deploy-fn
deploy-fn:
	gcloud functions deploy $(CF_NAME) \
	--env-vars-file .env.yaml \
	--runtime go111 \
	--trigger-http \
	--service-account=$(SERVICE_ACCOUNT) \
	$(CF_OPTIONS)

.PHONY: deploy-fn-pubsub
deploy-fn-pubsub:
	gcloud functions deploy $(CF_NAME) \
	--env-vars-file .env.yaml \
	--runtime go111 \
	--trigger-event google.pubsub.topic.publish \
	--service-account=$(SERVICE_ACCOUNT) \
	$(CF_OPTIONS)
