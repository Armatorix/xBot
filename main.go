package main

import (
	"fmt"
	"math/rand"

	"github.com/Armatorix/xBot/xwalker"
	"github.com/caarlos0/env/v11"
	"github.com/playwright-community/playwright-go"
)

type Config struct {
	Email         string   `env:"EMAIL"`
	Password      string   `env:"PASSWORD"`
	User          string   `env:"USERNAME"`
	Tags          []string `env:"TAGS"`
	SubFromHour   int      `env:"SUB_FROM_HOUR"`
	SubToHour     int      `env:"SUB_TO_HOUR"`
	MassUnsubHour int      `env:"MASS_UNSUB_HOUR"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("Error parsing environment variables: %v\n", err)
		return
	}
	err = playwright.Install()
	if err != nil {
		fmt.Printf("Error installing Playwright: %v\n ; continue", err)
	}

	xd, err := xwalker.LoadOrLoginX(cfg.Email, cfg.Password, cfg.User)
	if err != nil {
		fmt.Printf("Error loading or logging in to xwalker: %v\n", err)
		return
	}

	defer xd.Playwright.Stop()
	//scroll page to the bottom

	xd.RefuseAllCookies()

	followers, following, err := xd.FollowerAndFollowing()
	if err != nil {
		fmt.Println("Error getting followers and following:", err)
		return
	}

	toUnsub := max(0, following-followers+(rand.Intn(200)*followers/1000))
	if followers < 100 {
		toUnsub = 0
	}
	fmt.Println("Followers:", followers, "Following:", following, "To unsubscribe:", toUnsub)

	if err = xd.OpenFollowersPageAndUnsubN(toUnsub); err != nil {
		fmt.Println("Error opening followers page and unsubscribing:", err)
	}

	err = xd.StoreCookiesToFile()
	if err != nil {
		panic(err)
	}
}
