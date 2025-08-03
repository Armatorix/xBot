package xwalker

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/Armatorix/xBot/x/xrand"
	"github.com/playwright-community/playwright-go"
	"github.com/rotisserie/eris"
)

type XWalker struct {
	Playwright *playwright.Playwright
	Page       playwright.Page
	Username   string
}

func (x *XWalker) OpenFollowersPageAndUnsubN(n int) error {
	// Navigate to the followers page
	if err := x.openFollowingPage(); err != nil {
		fmt.Println("Error opening followers page:", err)
		return eris.Wrap(err, "failed to open followers page")
	}

	for range rand.Intn(3) + 1 {
		if err := x.scrollDownX(600); err != nil {
			return eris.Wrap(err, "failed to scroll down on followers page")
		}
	}

	// Unsubscribe from the first n followers
	for i := 0; i < n; i++ {
		fmt.Println("Unsubscribing from follower", i+1)

		// system press keyboard
		// find second button with thex Obserwujesz text
		buttons, err := x.Page.QuerySelectorAll("button:has-text('Obserwujesz')")
		if err != nil {
			return eris.Wrap(err, "failed to query 'Obserwujesz' buttons")
		}
		if len(buttons) < 2 {
			fmt.Println("Not enough 'Obserwujesz' buttons found, maybe already unsubscribed or not present")
			if err := x.openFollowingPage(); err != nil {
				fmt.Println("Error reopening followers page:", err)
				return eris.Wrap(err, "failed to reopen followers page after not finding 'Obserwujesz' buttons")
			}
			i--
			continue
		}
		// Click the second "Obserwujesz" button
		if err := buttons[rand.Intn(len(buttons))].Click(); err != nil {
			return eris.Wrap(err, "failed to click 'Obserwujesz' button")
		}
		sleep2N(1) // Wait for the unfollow action to complete
		// Click the "Przestań obserwować" button in the confirmation dialog
		// Check if has text "Przestań obserwować"
		if unfollowButtons, err := x.Page.QuerySelectorAll("button:has-text('Przestań obserwować')"); err != nil {
			return eris.Wrap(err, "failed to query 'Przestań obserwować' buttons")
		} else if len(unfollowButtons) == 0 {
			fmt.Println("No 'Przestań obserwować' button found, maybe already unsubscribed or not present")
			if err := x.refreshPage(); err != nil {
				return eris.Wrap(err, "failed to refresh page after not finding 'Przestań obserwować' button")
			}
			i--
			continue
		}
		if err := x.Page.Click("button:has-text('Przestań obserwować')"); err != nil {
			return eris.Wrap(err, "failed to click 'Przestań obserwować' button")
		}

		sleep2N(1)
		if rand.Intn(20) < 2 { // 10% chance to scroll down
			if err := x.scrollDown(); err != nil {
				return eris.Wrap(err, "failed to scroll down after unsubscribing")
			}
		} else if rand.Intn(20) < 2 { // 10% chance to click on a random link
			links, err := x.Page.QuerySelectorAll("a")
			if err != nil {
				return eris.Wrap(err, "failed to query all links on the page")
			}
			if len(links) > 0 {
				randomIndex := rand.Intn(len(links))
				if err := links[randomIndex].Click(); err != nil {
					return eris.Wrap(err, "failed to click on a random link")
				}
				sleep2N(2) // Wait for the click to complete
				// Go back to the followers page
				if _, err := x.Page.GoBack(); err != nil {
					return eris.Wrap(err, "failed to go back after clicking a random link")
				}
				sleep2N(3) // Wait for the page to load
				// if page is not followers page, go to the followers page again
				if err := x.openFollowingPage(); err != nil {
					fmt.Println("Error reopening followers page:", err)
					return eris.Wrap(err, "failed to reopen followers page after clicking a random link")
				}
			}
		}

	}

	return nil
}

func (x *XWalker) RefuseAllCookies() {
	x.Page.Goto("https://x.com")
	sleep2N(4) // Wait for the page to load
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
	sleep2N(2) // Wait for the action to complete
}

func (x *XWalker) FollowerAndFollowing() (int, int, error) {
	// Navigate to the followers page
	if err := x.openProfilePage(); err != nil {
		return 0, 0, fmt.Errorf("error opening followers page: %w", err)
	}

	// find <a> to /{username}/following
	fmt.Println("Finding following and followers counts for user:", x.Username)
	followingLink, err := x.Page.QuerySelector("a[href='/" + x.Username + "/following']")
	if err != nil {
		return 0, 0, fmt.Errorf("error finding following link: %w", err)
	}

	// Get the text content of the following link
	fmt.Println("Getting following count for user:", x.Username)
	followingText, err := followingLink.TextContent()
	if err != nil {
		return 0, 0, fmt.Errorf("error getting following text: %w", err)
	}
	// remove uunicode characters
	followingText = strings.Map(func(r rune) rune {
		if r == '\u00A0' {
			return -1 // Remove these characters
		}
		return r
	}, followingText)
	// Extract the number of following from the text
	parts := strings.Split(followingText, " ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected following text format: %s", followingText)
	}

	following, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error converting following count to integer: %w", err)
	}

	// find <a> to /{username}/verified_followers
	fmt.Println("Finding followers count for user:", x.Username)
	followersLink, err := x.Page.QuerySelector("a[href='/" + x.Username + "/verified_followers']")
	if err != nil {
		return 0, 0, fmt.Errorf("error finding followers link: %w", err)
	}
	// Get the text content of the followers link
	followersText, err := followersLink.TextContent()
	if err != nil {
		return 0, 0, fmt.Errorf("error getting followers text: %w", err)
	}
	// Extract the number of followers from the text
	parts = strings.Split(followersText, " ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected followers text format: %s", followersText)
	}

	followers, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error converting followers count to integer: %w", err)
	}

	return followers, following, nil
}

// follow users from a tag,
// if there is no users to follow, open a random user from the tag, then open their followers page
// and then do the same thing
func (x *XWalker) FollowUsersFromTag(n int, tag string) error {
	totalFollowed := 0
	if n <= 0 {
		return fmt.Errorf("number of users to follow must be greater than 0")
	}

	fmt.Println("Following", n, "users from tag:", tag)
	if _, err := x.Page.Goto(fmt.Sprintf("https://x.com/search?q=%s&src=typed_query&f=user", url.QueryEscape("#"+tag))); err != nil {
		return eris.Wrap(err, "failed to go to tag page")
	}

	// Wait for the page to load and display the users
	if _, err := x.Page.WaitForSelector("span:has-text('Użytkownicy')"); err != nil {
		return eris.Wrap(err, "failed to wait for users to load on tag page")
	}

	for totalFollowed < n {
		sleep2N(3)

		buttons, err := x.Page.QuerySelectorAll("button div span span:text-is('Obserwuj')")
		if err != nil {
			return eris.Wrap(err, "failed to query 'Follow' buttons")
		}
		for range 5 {
			if len(buttons) != 0 {
				break
			}

			fmt.Println("No 'Follow' buttons found, trying to scroll down")
			if err := x.scrollDownX(1800); err != nil {
				return eris.Wrap(err, "failed to scroll down to find more 'Follow' buttons")
			}
			// Open a random user from the tag

			buttons, err = x.Page.QuerySelectorAll("button:text-is('Obserwuj')")
			if err != nil {
				return eris.Wrap(err, "failed to query 'Follow' buttons")
			}
		}

		if len(buttons) == 0 {
			fmt.Println("Still no 'Follow' buttons found, maybe already followed or not present")
			// open random user from the tag - find buttons with data-testid="UserCell"
			userCells, err := x.Page.QuerySelectorAll("button[data-testid='UserCell']")
			if err != nil {
				return eris.Wrap(err, "failed to query user cells")
			}
			// click random user cell
			if len(userCells) == 0 {
				return fmt.Errorf("no user cells found on the tag page")
			}
			if err := userCells[rand.Intn(len(userCells))].Click(); err != nil {
				return eris.Wrap(err, "failed to click on a random user cell")
			}
			sleep2N(3) // Wait for the page to load
			buttons, err = x.Page.QuerySelectorAll("*:text-is('Obserwujących')")
			if err != nil {
				return eris.Wrap(err, "failed to query 'Followers' buttons")
			}
			sleep2N(2) // Wait for the page to load
			if len(buttons) == 0 {
				return fmt.Errorf("no 'Followers' buttons found on the user page")
			}
			// Click on the first 'Followers' button
			if err := buttons[0].Click(); err != nil {
				return eris.Wrap(err, "failed to click on 'Followers' button")
			}
			sleep2N(1)
			continue
		}

		// Click the first "Follow" button
		if err := buttons[rand.Intn(len(buttons))].Click(); err != nil {
			return eris.Wrap(err, "failed to click 'Follow' button")
		}
		if x.hasFollowLimitsReachedNote() {
			return fmt.Errorf("follow limits reached")
		}

		totalFollowed++
		fmt.Println("Followed", totalFollowed, "users from tag:", tag)
	}

	return nil
}

// follow users from a tag,
// if there is no users to follow, open a random user from the tag, then open their followers page
// and then do the same thing
func (x *XWalker) FollowRepostersFromTag(n int, tag string) error {
	totalFollowed := 0
	if n <= 0 {
		return fmt.Errorf("number of users to follow must be greater than 0")
	}

	fmt.Println("Following", n, "users from tag:", tag)
	if _, err := x.Page.Goto(fmt.Sprintf("https://x.com/search?q=%s&src=typed_query", url.QueryEscape(tag))); err != nil {
		return eris.Wrap(err, "failed to go to tag page")
	}

	// Wait for the page to load and display the users
	if _, err := x.Page.WaitForSelector("article[data-testid='tweet']"); err != nil {
		return eris.Wrap(err, "failed to wait for users to load on tag page")
	}

	// query elements "article" with with test-id="tweet"
	articles, err := x.Page.QuerySelectorAll("article[data-testid='tweet']")
	if err != nil {
		return eris.Wrap(err, "failed to query 'Follow' buttons")
	}
	// open first article
	if len(articles) == 0 {
		return fmt.Errorf("no articles found on the tag page")
	}

	if err := articles[0].Click(); err != nil {
		return eris.Wrap(err, "failed to click on the first article")
	}
	sleep2N(3) // Wait for the page to load

	gotoUrl, err := x.FromPostToRetweets()
	if err != nil {
		return eris.Wrap(err, "failed to get retweets URL from post")
	}

	fmt.Println("Going to retweets page:", gotoUrl)
	if _, err := x.Page.Goto(gotoUrl); err != nil {
		return eris.Wrap(err, "failed to go to retweets page")
	}

	sleep2N(3)
	// Wait for the retweets page to load
	if _, err := x.Page.WaitForSelector("button:has-text('Obserwuj')"); err != nil {
		return eris.Wrap(err, "failed to wait for retweets page to load")
	}

	failedAttempts := 0
	for totalFollowed < n {
		sleep2N(1)

		// with test-id="*-follow"
		buttons, err := x.Page.QuerySelectorAll("button:text-is('Obserwuj')")
		if err != nil {
			return eris.Wrap(err, "failed to query 'Follow' buttons")
		}

		if len(buttons) != 0 {
			if err := xrand.SliceElement(buttons).Click(); err != nil {
				return eris.Wrap(err, "failed to click 'Follow' button")
			}
			failedAttempts = 0 // Reset failed attempts after a successful follow
			totalFollowed++
			sleep2N(1) // Wait for the follow action to complete
			fmt.Println("Followed", totalFollowed, "users from tag:", tag)
		} else if len(buttons) == 0 {
			failedAttempts++
			fmt.Println("Still no 'Follow' buttons found, maybe already followed or not present")
			if err := x.scrollDownX(300); err != nil {
				return eris.Wrap(err, "failed to scroll down after not finding 'Follow' buttons")
			}
		}

		if failedAttempts >= 5 {
			return fmt.Errorf("too many failed attempts to find 'Follow' buttons, stopping")
		}
		if x.hasFollowLimitsReachedNote() {
			return fmt.Errorf("follow limits reached")
		}
	}

	return nil
}
