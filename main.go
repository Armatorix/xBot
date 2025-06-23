package main

import (
	"fmt"

	"github.com/Armatorix/xBot/xwalker"
	"github.com/caarlos0/env/v11"
	"github.com/playwright-community/playwright-go"
)

type Config struct {
	Email    string `env:"EMAIL"`
	Password string `env:"PASSWORD"`
	User     string `env:"USER"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("Error parsing environment variables: %v\n", err)
		return
	}
	playwright.Install()
	xd, err := xwalker.LoginFromCookiesFile()
	if err != nil {
		fmt.Println("No cookies file found, logging in with credentials", err)
	}
	if xd == nil {
		xd, err = xwalker.LoginX("dummy.alertz@gmail.com", "#oTZg6VweLY9psVzkXPy", "polski_wojt")
		if err != nil {
			panic(err)
		}
	}

	defer xd.Playwright.Stop()
	//scroll page to the bottom

	xd.RefuseAllCookies()

	xd.FollowUnfollowedFromHash("#pociągPrawych", 20)
	xd.OpenFollowersPageAndUnsubN(10)
	err = xd.StoreCookiesToFile()
	if err != nil {
		panic(err)
	}
}
