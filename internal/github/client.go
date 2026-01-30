// Package github provides a client for interacting with the GitHub Notifications API.
// It wraps the go-github library and converts responses to application-specific models.
package github

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client and provides notification-specific operations.
type Client struct {
	gh *github.Client
}

// NewClient creates a new GitHub API client using the GITHUB_TOKEN environment variable.
// Returns an error if the token is not set or authentication fails.
func NewClient() (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("GITHUB_TOKEN environment variable not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &Client{gh: github.NewClient(tc)}, nil
}
