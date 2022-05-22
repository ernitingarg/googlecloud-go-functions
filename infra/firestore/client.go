package firestore

import (
	"context"
	"log"
	"soteria-functions/env"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

func NewFirestoreClient(projectID env.ProjectID) *firestore.Client {
	ctx := context.Background()

	var firestoreClient *firestore.Client
	var err error

	// KEY_FILE_PATHが定義している理由はローカル起動時にfirestoreに接続するためのkeyfileが必要なため
	if keyPath := env.EnvVars.GCP.KeyPath; keyPath != "" {
		opt := option.WithCredentialsFile(keyPath)
		firestoreClient, err = firestore.NewClient(ctx, string(projectID), opt)
	} else {
		firestoreClient, err = firestore.NewClient(ctx, string(projectID))
	}
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed initalize firestore"))
	}
	return firestoreClient
}
