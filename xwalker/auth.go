package xwalker

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
	"github.com/rotisserie/eris"
)

var (
	headless = false // Set to false if you want to see the browser actions
	// 30 minutes as milliseconds
	timeout = float64(30 * 60 * 1000)
)

func loginFromCookiesFile(ctx context.Context, username string) (*XWalker, error) {
	f, err := os.ReadFile(username + "_cookies.txt")
	if err != nil {
		return nil, eris.Wrap(err, "failed to read cookies file")
	}

	cookiez := strings.Split(string(f), "; ")
	if len(cookiez) == 0 || (len(cookiez) == 1 && cookiez[0] == "") {
		return nil, fmt.Errorf("no valid cookies found in file %s", username+"_cookies.txt")
	}

	var cookies []playwright.OptionalCookie
	for _, cookie := range cookiez {
		if cookie == "" {
			continue // Skip empty cookies
		}
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) != 2 {
			return nil, eris.Errorf("invalid cookie format: %s", cookie)
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
		return nil, eris.Wrap(err, "failed to start Playwright")
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
		Timeout:  playwright.Float(timeout), // Set a timeout for launching the browser
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to launch browser")
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		StorageState: &playwright.OptionalStorageState{
			Cookies: cookies,
		},
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to create new browser context with cookies")
	}
	page, err := context.NewPage()
	if err != nil {
		return nil, eris.Wrap(err, "failed to create new page in browser context")
	}
	page.SetDefaultTimeout(timeout) // Disable timeout for page operations
	// goto page and check if logged in
	return &XWalker{
		Username:   username,
		Playwright: pw,
		Page:       page,
	}, nil
}

func LoadOrLoginX(ctx context.Context, email, pass, user string) (*XWalker, error) {
	xd, err := loginFromCookiesFile(ctx, user)
	if err != nil {
		fmt.Println("No cookies file found, logging in with credentials", err)
	}

	if xd != nil {
		if _, err := xd.Page.Goto("https://x.com"); err != nil {
			return nil, eris.Wrap(err, "failed to go to x.com")
		}
		// TODO: wait for load state and check if user logged with cookies

	} else {
		xd, err = loginX(ctx, email, pass, user)
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

func loginX(ctx context.Context, email, pass, user string) (*XWalker, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return nil, err
	}

	context, err := browser.NewContext()
	if err != nil {
		return nil, err
	}
	context.SetDefaultTimeout(timeout)

	page, err := context.NewPage()
	if err != nil {
		return nil, err
	}
	page.SetDefaultTimeout(timeout) // Disable timeout for page operations
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

	sleep2N(2) // Wait for the next page to load
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
		sleep2N(1)
	}
	// Wait for the password field to appear
	if _, err = page.WaitForSelector("input[name='password']"); err != nil {
		return nil, eris.Wrap(err, "failed to find password input field")
	}

	sleep2N(2) // Wait for the next page to load
	if err := page.Fill("input[name='password']", pass); err != nil {
		return nil, eris.Wrap(err, "failed to fill password field")
	}

	// button with text "Log in"
	if err := page.Click("button:has-text('Log in')"); err != nil {
		return nil, eris.Wrap(err, "failed to click 'Log in' button")
	}

	sleep2N(2) // Wait for the login to complete
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

func (x *XWalker) Logout() {
	x.Page.Goto("https://x.com/logout")
	sleep2N(4)
	x.Page.Click("button:has-text('Wyloguj się')") // Click the logout button

	sleep2N(4)
	os.Remove(x.Username + "_cookies.txt") // Remove the cookies file
	fmt.Println("Logged out and cookies file removed.")
	return
}
