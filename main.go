package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Armatorix/xBot/x/xrand"
	"github.com/Armatorix/xBot/xwalker"
	"github.com/caarlos0/env/v11"
	"github.com/playwright-community/playwright-go"
)

// TODO:  handle "Nie możesz obecnie obserwować więcej osób."
// NOTE: stop sub immidietly

func main() {
	now := time.Now()
	// timezone warsaw
	timezone, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		timezone = time.UTC
	}
	now = now.In(timezone)
	fmt.Println("Current time in Warsaw:", now.Format("15:04:05"))
	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		fmt.Printf("Error parsing environment variables: %v\n", err)
		return
	}
	if !cfg.Localdev && (now.Hour() < 6 || now.Hour() > 23) {
		fmt.Println("xBot is not allowed to run at this hour. Please try again later.")
		return
	}
	err = playwright.Install()
	if err != nil {
		fmt.Printf("Error installing Playwright: %v\n ; continue", err)
	}

	if rand.Float64() < 0.02 {
		fmt.Println("Randomly skipping xBot execution")
		return
	}
	initSleep := time.Duration(rand.Intn(7)) * time.Minute
	fmt.Println("Starting xBot...\n Starting in", initSleep, "minutes")
	if cfg.Localdev {
		fmt.Println("Running in localdev mode")
	} else {
		time.Sleep(initSleep)
		fmt.Println("xBot started")
	}
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
		xd.Logout()
		return
	}
	fmt.Println("Followers:", followers)
	fmt.Println("Following:", following)

	{
		toUnsub := unfollowCount(now.Hour(), followers, following)
		fmt.Println("To unsubscribe:", toUnsub)

		if err = xd.OpenFollowersPageAndUnsubN(toUnsub); err != nil {
			fmt.Println("Error opening followers page and unsubscribing:", err)
		}
	}

	{
		// mass sub
		toFollow := followCount(now.Hour())

		err := xd.FollowRepostersFromTag(toFollow, xrand.SliceElement(cfg.Tags))
		if err != nil {
			fmt.Println("Error following from tag:", err)
		}
	}

	err = xd.StoreCookiesToFile()
	if err != nil {
		panic(err)
	}
}

func followCount(i int) int {
	if i < 6 || i > 23 {
		return 0
	}
	toFollow := (i/4 - 1)
	toFollow *= toFollow
	toFollow += rand.Intn(8)
	return max(toFollow, 0)
}

func unfollowCount(i int, followers, following int) int {
	if i < 6 || i > 23 {
		return 0
	}
	toUnsub := i/4 + rand.Intn(15)
	if followers < 200 {
		toUnsub = 0
	}
	if followers > int(float64(following)*1.1) {
		toUnsub = 0
	}
	return toUnsub
}
