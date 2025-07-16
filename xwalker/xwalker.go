package xwalker

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type XWalker struct {
	Playwright *playwright.Playwright
	Page       playwright.Page
}

func LoginFromCookiesFile() (*XWalker, error) {
	f, err := os.ReadFile("cookies.txt")
	if err != nil {
		return nil, err
	}

	cookiez := strings.Split(string(f), "; ")
	if len(cookiez) == 0 || (len(cookiez) == 1 && cookiez[0] == "") {
		return nil, fmt.Errorf("no cookies found in file")
	}

	var cookies []playwright.OptionalCookie
	for _, cookie := range cookiez {
		if cookie == "" {
			continue // Skip empty cookies
		}
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid cookie format: %s", cookie)
		}
		name := parts[0]
		value := parts[1]
		cookies = append(cookies, playwright.OptionalCookie{
			Name:  name,
			Value: value,
			URL:   playwright.String("https://x.com"), // Optional: specify the domain if needed
		})
	}
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Timeout:  playwright.Float(0), // Set a timeout for launching the browser
	})
	if err != nil {
		return nil, err
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		StorageState: &playwright.OptionalStorageState{
			Cookies: cookies,
		},
	})
	if err != nil {
		return nil, err
	}
	page, err := context.NewPage()
	if err != nil {
		return nil, err
	}
	page.SetDefaultTimeout(0) // Disable timeout for page operations
	// goto page and check if logged in
	return &XWalker{
		Playwright: pw,
		Page:       page,
	}, nil
}

func LoginX(email, pass, user string) (*XWalker, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	context, err := browser.NewContext()
	if err != nil {
		return nil, err
	}
	context.SetDefaultTimeout(0)

	page, err := context.NewPage()
	if err != nil {
		return nil, err
	}
	page.SetDefaultTimeout(0) // Disable timeout for page operations
	_, err = page.Goto("https://x.com/i/flow/login")
	if err != nil {
		return nil, err
	}
	// wait for the form to load
	if _, err = page.WaitForSelector("input[autocomplete='username']"); err != nil {
		return nil, err
	}
	// Fill in the username and password fields
	if err := page.Fill("input[autocomplete='username']", email); err != nil {
		return nil, err
	}

	// button with text "Next"
	if err := page.Click("button:has-text('Next')"); err != nil {
		return nil, err
	}

	time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond) // Wait for the next page to load
	// check page contains "There was unusual login activity on your account."
	if cont, _ := page.Content(); strings.Contains(cont, "There was unusual login activity on your account.") {
		// fill username again
		if err := page.Fill("input", user); err != nil {
			return nil, err
		}
		// click button with text "Next"
		if err := page.Click("button:has-text('Next')"); err != nil {
			return nil, err
		}
		time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond)
	}
	// Wait for the password field to appear
	if _, err = page.WaitForSelector("input[name='password']"); err != nil {
		return nil, err
	}

	time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond) // Wait for the next page to load
	if err := page.Fill("input[name='password']", pass); err != nil {
		return nil, err
	}

	// button with text "Log in"
	if err := page.Click("button:has-text('Log in')"); err != nil {
		return nil, err
	}
	time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond) // Wait for the login to complete
	return &XWalker{
		Playwright: pw,
		Page:       page,
	}, nil
}

func (x *XWalker) OpenFollowersPageAndUnsubN(n int) error {
	// Navigate to the followers page
	if _, err := x.Page.Goto("https://x.com/polski_wojt/following"); err != nil {
		return err
	}

	// Wait for the followers list to load
	if _, err := x.Page.WaitForSelector("button:has-text('Following')"); err != nil {
		return err
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(150))*time.Millisecond)

	// Unsubscribe from the first n followers
	for i := 0; i < n; i++ {
		fmt.Println("Unsubscribing from follower", i+1)
		// system press keyboard
		// find second button with thex Following text
		buttons, err := x.Page.QuerySelectorAll("button:has-text('Following')")
		if err != nil {
			return err
		}
		if len(buttons) < 2 {
			return fmt.Errorf("not enough buttons found")
		}
		// Click the second "Following" button
		if err := buttons[1].Click(); err != nil {
			return err
		}
		time.Sleep(time.Duration(rand.Intn(340)) * time.Millisecond) // Wait for the unfollow action to complete
		// Click the "Unfollow" button in the confirmation dialog
		// Check if has text "Unfollow"
		if unfollowButtons, err := x.Page.QuerySelectorAll("button:has-text('Unfollow')"); err != nil {
			return err
		} else if len(unfollowButtons) == 0 {
			return nil // No "Unfollow" button found, maybe already unfollowed
		}
		if err := x.Page.Click("button:has-text('Unfollow')"); err != nil {
			return err
		}
		time.Sleep(time.Second + time.Duration(rand.Intn(150))*time.Millisecond) // Wait
	}

	return nil
}

// TODO: handle if at some point the page changes, restart the process few times - then notify me

func (x *XWalker) StoreCookiesToFile() error {
	cookies, err := x.Page.Context().Cookies()
	if err != nil {
		return err
	}

	// Convert cookies to a string format or save them to a file
	cookieData := ""
	for _, cookie := range cookies {
		cookieData += fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value)
	}

	f, err := os.Create("cookies.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(cookieData)
	if err != nil {
		return err
	}
	return nil
}

func (x *XWalker) RefuseAllCookies() {
	x.Page.Goto("https://x.com")
	// Click the "Refuse" button for cookies if it exists
	ls, err := x.Page.Locator("button:has-text('Refuse non-essential cookies')").All()
	if err != nil {
		fmt.Println("Error finding cookie refusal button:", err)
		return
	}
	if len(ls) == 0 {
		fmt.Println("No cookie refusal button found, maybe already refused or not present")
		return
	}
	fmt.Println("Refusing non-essential cookies")
	if err := x.Page.Click("button:has-text('Refuse non-essential cookies')"); err != nil {
	}
	time.Sleep(time.Second + time.Duration(rand.Intn(150))*time.Millisecond) // Wait for the action to complete
}

func (x *XWalker) FollowUnfollowedFromHash(hash string, n int) error {
	q := url.QueryEscape(hash) // Ensure the hashtag is URL-encoded
	_, err := x.Page.Goto(fmt.Sprintf("https://x.com/search?q=%s&src=hashtag_click&f=user", q))
	if err != nil {
		return err
	}

	// Wait for the page to load and display the users
	if _, err := x.Page.WaitForSelector("button:has-text('Follow')"); err != nil {
		return err
	}

	totalFollowed := 0
	queryAttempts := 0
	// Find all "Follow" buttons
	// Follow the first n users
	for {
		time.Sleep(time.Second + time.Duration(rand.Intn(150))*time.Millisecond) // Wait for the follow action to complete
		buttons, err := x.Page.QuerySelectorAll("button:has-text('Follow')")
		if err != nil {
			return err
		}
		if len(buttons) == 0 {
			// scroll to the bottom of the page to load more users
			if _, err := x.Page.Evaluate("window.scrollTo(0, document.body.scrollHeight)"); err != nil {
				return err
			}
			queryAttempts++
			if queryAttempts > 5 {
				return nil
			}
			continue
		}

		// Click the first "Follow" button
		if err := buttons[0].Click(); err != nil {
			return err
		}

		totalFollowed++
		if totalFollowed >= n {
			break // Stop if we've followed enough users
		}
	}

	return nil
}
