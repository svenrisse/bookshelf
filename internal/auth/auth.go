package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	discordClientId := os.Getenv("DISCORD_CLIENT_ID")
	discordClientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	key := os.Getenv("SESSION_KEY")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(86400 * 30)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true

	gothic.Store = store

	goth.UseProviders(google.New(googleClientId, googleClientSecret, "http://localhost:4000/v1/auth/google/callback"))
	goth.UseProviders(discord.New(discordClientId, discordClientSecret, "http://localhost:4000/v1/auth/discord/callback"))
	goth.UseProviders(github.New(githubClientId, githubClientSecret, "http://localhost:4000/v1/auth/github/callback"))
}
