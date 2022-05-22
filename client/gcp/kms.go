package gcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"soteria-functions/env"
	"soteria-functions/logger"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/pkg/errors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type KmsClient interface {
	Encryption(ringName string, keyName string, plainText string) (string, error)
}

type kmsClient struct {
	client     *kms.KeyManagementClient
	ctx        *context.Context
	projectID  string
	locationID string
}

func NewKmsClient(projectID env.ProjectID, locationID env.LocationID) KmsClient {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	return &kmsClient{client: client, projectID: projectID.String(), locationID: locationID.String(), ctx: &ctx}
}

func (kc *kmsClient) Encryption(ringName string, keyName string, plainText string) (string, error) {
	logger.Log.Info(fmt.Sprintf("ringName, keyName, plainText: %v %v %v", ringName, keyName, plainText))

	request := &kmspb.EncryptRequest{
		Name:      kc.cryptoKeyName(ringName, keyName),
		Plaintext: []byte(plainText),
	}
	response, err := kc.client.Encrypt(*kc.ctx, request)
	if err != nil {
		return "", errors.Wrap(err, "failed Encrypt")
	}
	logger.Log.Info(fmt.Sprintf("response: %v", response))
	encByte := response.GetCiphertext()
	return base64.StdEncoding.EncodeToString(encByte), nil
}

func (kc *kmsClient) keyRingName(ringName string) string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", kc.projectID, kc.locationID, ringName)
}

func (kc *kmsClient) cryptoKeyName(ringName string, keyName string) string {
	kringName := kc.keyRingName(ringName)
	return fmt.Sprintf("%s/cryptoKeys/%s", kringName, keyName)
}
