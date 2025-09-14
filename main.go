package main

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

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

	xd, err := xwalker.LoadOrLoginX(ctx, cfg.Email, cfg.Password, cfg.User)
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

	if cfg.Cooldown || now.Day() > 26 {
		toUnsub := unfollowCount(now.Hour(), followers, following)
		fmt.Println("To unsubscribe:", toUnsub)

		if err = xd.OpenFollowersPageAndUnsubN(toUnsub / 2); err != nil {
			fmt.Println("Error opening followers page and unsubscribing:", err)
		}

		if err = xd.OpenFollowingPageAndUnsubN(toUnsub / 2); err != nil {
			fmt.Println("Error opening following page and unsubscribing:", err)
		}
	} else {
		fmt.Println("Skipping unsubscription, it's not after the 15th.")
	}

	{
		// mass sub
		toFollow := followCount(now.Hour(), followers)

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

func followCount(i int, followers int) int {
	if i < 7 {
		return 0 // No following in the first 7 hours
	}
	switch {
	case followers < 100:
		return (i / 10) + rand.Intn(i+1)
	case followers < 200:
		return (i / 6) + rand.Intn(i+1)
	case followers < 500:
		return (i / 4) + rand.Intn(i+1)
	default:
		return (i / 3) + rand.Intn(i+1)
	}
}

func unfollowCount(i int, followers, following int) int {
	if i < 7 {
		return 0 // No unfollowing in the first 7 hours
	}

	// do not unfollow if has 10% of following
	if followers > following*10 {
		return 0
	}

	switch {
	case followers > 2500:
		return 400 / (48 - 7*2)
	case followers < 100:
		return (i / 10) + rand.Intn(i+1)/2
	case followers < 200:
		return (i / 6) + rand.Intn(i+1)/2
	case followers < 500:
		return (i / 4) + rand.Intn(i+1)/2
	case following < 1000:
		return (i / 3) + rand.Intn(i+1)/2
	default:
		return (i / 2) + rand.Intn(i+1)/2

	}
}
