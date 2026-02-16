package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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

func LoadEnv() {
	currentDir, _ := os.Getwd()
	fmt.Printf("📍 Текущая рабочая директория: %s", currentDir)

	// Проверяем, существует ли .env в этой директории
	envPath := filepath.Join(currentDir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Printf("❌ .env не найден по пути: %s", envPath)
	} else {
		fmt.Printf("✅ .env найден по пути: %s", envPath)
	}

	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("⚠️ Ошибка загрузки .env: %v", err)
	}
}

func NewHernyaOauthConfig() *OauthConfig {
	LoadEnv()
	clienId := os.Getenv("CLIENT_ID")
	if clienId == "" {
		fmt.Print("wwwww")
		return nil
	}

	clienSecret := os.Getenv("CLIENT_SECRET")
	if clienSecret == "" {
		fmt.Print("wwwwwsssss")
		return nil
	}

	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		fmt.Print("sssss")
		return nil
	}

	return &OauthConfig{
		ClientId:     clienId,
		ClientSecret: clienSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"login:birthday",
			"login:email",
			"login:info",
			"login:avatar",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.yandex.ru/authorize",
			TokenURL: "https://oauth.yandex.ru/token",
		},
	}
}
