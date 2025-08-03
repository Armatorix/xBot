package xwalker

import (
	"fmt"
	"strings"
)

// FromPostToRetweets constructs the URL for the retweets of the current post.
// It assumes the current page is a post and returns the URL for its retweets.
// The URL format is expected to be "https://x.com/username/status/post_id/retweets".
//
// Returns the retweets URL or an error if the current page URL is invalid.
func (x *XWalker) FromPostToRetweets() (string, error) {
	retweetsURL := strings.Split(x.Page.URL(), "/")
	if len(retweetsURL) < 6 {
		return "", fmt.Errorf("invalid URL format: %s", x.Page.URL())
	}

	retweetsURL = retweetsURL[:6]
	return strings.Join(retweetsURL, "/") + "/retweets", nil
}
