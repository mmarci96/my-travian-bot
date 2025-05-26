package main

import (
	"fmt"
	"log"

	"github.com/mmarci96/my-travian-bot/app/auth"
	"github.com/mmarci96/my-travian-bot/app/internal/config"
	"github.com/mmarci96/my-travian-bot/app/internal/redirector"
)

func main() {
	fmt.Println("Starting up!")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	token, err := auth.CreateLoginRequest(cfg.Auth.Email, cfg.Auth.Password, cfg.Auth.ServerUrl)
	if err != nil {
		fmt.Println("Error logging in: ", err)
	}

	fmt.Println("Logged in, response body: ", token)
	err = redirector.FollowRedirect(token, cfg.Auth.ServerUrl)
	if err != nil {
		log.Fatalf("Redirect failed: %v", err)
	}
}
