package xwalker

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rotisserie/eris"
)

func (x *XWalker) scrollDown() error {
	if _, err := x.Page.Evaluate("window.scrollTo(0, document.body.scrollHeight+" + strconv.Itoa(rand.Intn(400)) + ")"); err != nil {
		return eris.Wrap(err, "failed to scroll down the page")
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(100))*time.Millisecond)
	fmt.Println("Scrolled down successfully")
	return nil
}

func (x *XWalker) openProfilePage() error {
	if _, err := x.Page.Goto("https://x.com/" + x.Username); err != nil {
		return eris.Wrap(err, "failed to go to profile page")
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(150))*time.Millisecond)
	// Check if the page is loaded by looking for the profile header
	if _, err := x.Page.WaitForSelector("div:has-text('" + x.Username + "')"); err != nil {
		return fmt.Errorf("profile page did not load correctly: %w", err)
	}
	fmt.Println("Profile page opened successfully")
	return nil
}

func (x *XWalker) scrollDownX(v int) error {
	if _, err := x.Page.Evaluate("window.scrollTo(0, document.body.scrollHeight+" + strconv.Itoa(rand.Intn(100)+v) + ")"); err != nil {
		return eris.Wrap(err, "failed to scroll down the page")
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(350))*time.Millisecond)
	fmt.Println("Scrolled down successfully")
	return nil
}

func (x *XWalker) openFollowingPage() error {
	// Navigate to the followers page
	if _, err := x.Page.Goto("https://x.com/" + x.Username + "/following"); err != nil {
		return eris.Wrap(err, "failed to go to followers page")
	}

	// Wait for the followers list to load
	if _, err := x.Page.WaitForSelector("button:has-text('Obserwujesz')"); err != nil {
		return eris.Wrap(err, "failed to wait for followers list to load")
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(150))*time.Millisecond) // Wait for the page to load

	return nil
}

func (x *XWalker) refreshPage() error {
	if _, err := x.Page.Reload(); err != nil {
		return err
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(350))*time.Millisecond)
	// Wait for the page to reload
	if err := x.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	}); err != nil {
		return eris.Wrap(err, "failed to wait for page to reload after refreshing")
	}
	fmt.Println("Page reloaded successfully")
	return nil
}
