package firebase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type AuthClient interface {
	GetAuthToken(ctx context.Context, r *http.Request) (*auth.Token, error)
}

type authClient struct {
	auth *auth.Client
}

func NewAuthClient(keyPath string) AuthClient {
	ctx := context.Background()
	var app *firebase.App
	var err error
	if keyPath != "" {
		opt := option.WithCredentialsFile(keyPath)
		app, err = firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatalln(errors.Wrap(err, fmt.Sprintf("failed initalize firebase. %+v", keyPath)))
		}
	} else {
		app, err = firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatalln(errors.Wrap(err, "failed initalize firebase"))
		}
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed initalize auth"))
	}

	return &authClient{auth: auth}
}

func (ac *authClient) GetAuthToken(ctx context.Context, r *http.Request) (*auth.Token, error) {
	idToken := ac.getIDToken(r)
	return ac.verifyIDToken(ctx, idToken)
}

func (ac *authClient) verifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return ac.auth.VerifyIDToken(ctx, idToken)
}

// GetIDToken retribe JWT from header
func (ac *authClient) getIDToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.Replace(authHeader, "Bearer ", "", 1)
}
