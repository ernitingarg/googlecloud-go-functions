package auth

import (
	"context"
	"net/http"
	"soteria-functions/client/firebase"
	"soteria-functions/env"

	"github.com/pkg/errors"
)

type AuthUsecase struct {
	authClient firebase.AuthClient
}

func NewAuthUsecase() *AuthUsecase {
	authClient := firebase.NewAuthClient(env.EnvVars.GCP.KeyPath)
	return &AuthUsecase{authClient: authClient}
}

func (a *AuthUsecase) GetUID(ctx context.Context, r *http.Request) (string, error) {
	token, err := a.authClient.GetAuthToken(ctx, r)
	if err != nil {
		return "", errors.Wrap(err, "failed GetUID")
	}
	return token.UID, nil
}
