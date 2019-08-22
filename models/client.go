package models

import (
    "context"
    "github.com/digitalocean/godo"
    "golang.org/x/oauth2"
    "os"
)

type TokenSource struct {
    AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
    token := &oauth2.Token{AccessToken: t.AccessToken}

    return token, nil
}

func MakeClient() *godo.Client {
    tokenSource := &TokenSource{AccessToken: os.Getenv("DO_API_KEY")}
    oauthClient := oauth2.NewClient(context.Background(), tokenSource)
    client := godo.NewClient(oauthClient)

    return client
}

