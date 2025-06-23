package main

import (
	"github.com/Armatorix/xBot/xwalker"
	"github.com/playwright-community/playwright-go"
)

func main() {
	playwright.Install()
	xd, err := xwalker.LoginX("", "", "")
	if err != nil {
		panic(err)
	}

	defer xd.Playwright.Stop()
	//scroll page to the bottom

	err = xd.OpenFollowersPageAndUnsubN(50)
	if err != nil {
		panic(err)
	}
}
