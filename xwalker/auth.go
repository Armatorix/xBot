package xwalker

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rotisserie/eris"
)

func loginFromCookiesFile(username string) (*XWalker, error) {
	f, err := os.ReadFile(username + "_cookies.txt")
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
		Username:   username,
		Playwright: pw,
		Page:       page,
	}, nil
}

func LoadOrLoginX(email, pass, user string) (*XWalker, error) {
	xd, err := loginFromCookiesFile(user)
	if err != nil {
		fmt.Println("No cookies file found, logging in with credentials", err)
	}
	if xd == nil {
		xd, err = loginX(email, pass, user)
		if err != nil {
			return nil, eris.Wrap(err, "failed to login with credentials")
		}

		err = xd.StoreCookiesToFile()
		if err != nil {
			return nil, eris.Wrap(err, "failed to store cookies to file")
		}
	}
	return xd, nil

}

func loginX(email, pass, user string) (*XWalker, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
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
			return nil, eris.Wrap(err, "failed to fill username field after unusual login activity")
		}
		// click button with text "Next"
		if err := page.Click("button:has-text('Next')"); err != nil {
			return nil, eris.Wrap(err, "failed to click 'Next' button after unusual login activity")
		}
		time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond)
	}
	// Wait for the password field to appear
	if _, err = page.WaitForSelector("input[name='password']"); err != nil {
		return nil, eris.Wrap(err, "failed to find password input field")
	}

	time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond) // Wait for the next page to load
	if err := page.Fill("input[name='password']", pass); err != nil {
		return nil, eris.Wrap(err, "failed to fill password field")
	}

	// button with text "Log in"
	if err := page.Click("button:has-text('Log in')"); err != nil {
		return nil, eris.Wrap(err, "failed to click 'Log in' button")
	}
	time.Sleep(2*time.Second + time.Duration(rand.Intn(10))*time.Millisecond) // Wait for the login to complete
	return &XWalker{
		Playwright: pw,
		Page:       page,
		Username:   user,
	}, nil
}

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

	f, err := os.Create(x.Username + "_cookies.txt")
	if err != nil {
		return fmt.Errorf("failed to create cookies file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(cookieData)
	if err != nil {
		return fmt.Errorf("failed to write cookies to file: %w", err)
	}
	return nil
}
