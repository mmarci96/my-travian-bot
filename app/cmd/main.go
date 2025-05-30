package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmarci96/my-travian-bot/app/auth"
	"github.com/mmarci96/my-travian-bot/app/internal/config"
)

func ExtractNoScriptContent(html string) (string, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	noscriptText := ""
	doc.Find("noscript").Each(func(i int, s *goquery.Selection) {
		noscriptText += strings.TrimSpace(s.Text()) + "\n"
	})

	return strings.TrimSpace(noscriptText), nil
}
func main() {
	fmt.Println("Starting up!")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	respBody, err := auth.CreateLoginRequest(cfg.Auth.Email, cfg.Auth.Password, cfg.Auth.ServerUrl)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	fmt.Println("Logged in, response body:", respBody)

	html, err := auth.FollowRedirect(respBody, cfg.Auth.ServerUrl)
	if err != nil {
		log.Fatalf("Redirect failed: %v", err)
	}
	res, err := ExtractNoScriptContent(html)

	if err != nil {
		log.Fatalf("Extract failed: %v", err)
	}
	println("Result: ", res)

	// resourcegoquery:= collector.ParseResourceHTML(html)
	// if err != nil {
	// 	log.Fatalf("Failed to parse resources: %v", err)
	// }
	//
	// fmt.Printf("Resources: %+v\n", resources)
}
