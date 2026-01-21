package config

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
)

type OauthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

func (o *OauthConfig) ToOauth() oauth2.Config {
	return oauth2.Config{
		ClientID:     o.ClientId,
		ClientSecret: o.ClientSecret,
		RedirectURL:  o.RedirectURL,
		Scopes:       o.Scopes,
		Endpoint:     o.Endpoint,
	}
}

func NewHernyaOauthConfig() (*OauthConfig, error) {
	clienId := os.Getenv("CLIENT_ID")
	if clienId == "" {
		return nil, errors.New("Error to check client id from .env file")
	}

	clienSecret := os.Getenv("CLIENT_SECRET")
	if clienSecret == "" {
		return nil, errors.New("Error to check client secret from .env file")
	}

	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		return nil, errors.New("Error to check redirect url from .env file")
	}

	return &OauthConfig{
		ClientId:     clienId,
		ClientSecret: clienSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"birthday",
			"email",
			"info",
			"avatar",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.yandex.ru/authorize",
			TokenURL: "https://oauth.yandex.ru/token",
		},
	}, nil
}
