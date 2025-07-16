package main

import (
	"fmt"
	"math/rand"
	"time"

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
	now := time.Now()
	if now.Hour() < 6 || now.Hour() > 23 {
		fmt.Println("xBot is not allowed to run at this hour. Please try again later.")
		return
	}
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

	initSleep := time.Duration(rand.Intn(4)) * time.Minute
	fmt.Println("Starting xBot...\n Starting in", initSleep, "minutes")
	time.Sleep(initSleep)
	fmt.Println("xBot started")
	// Initialize xwalker with the provided configuration

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

	if now.Hour() == cfg.MassUnsubHour {
		toUnsub := max(0, following-followers+(rand.Intn(4)*followers/100))
		if followers < 100 {
			toUnsub = 0
		}
		fmt.Println("Followers:", followers, "Following:", following, "To unsubscribe:", toUnsub)

		if err = xd.OpenFollowersPageAndUnsubN(toUnsub); err != nil {
			fmt.Println("Error opening followers page and unsubscribing:", err)
		}
	}

	if now.Hour() >= cfg.SubFromHour && now.Hour() <= cfg.SubToHour {
		toFollow := 3 + rand.Intn(8)

		err := xd.FollowFromTag(toFollow, cfg.Tags[rand.Intn(len(cfg.Tags))])
		if err != nil {
			fmt.Println("Error following from tag:", err)
			return
		}
	}

	err = xd.StoreCookiesToFile()
	if err != nil {
		panic(err)
	}
}
