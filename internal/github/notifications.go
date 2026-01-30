// Package github provides a client for interacting with the GitHub Notifications API.
// It wraps the go-github library and converts responses to application-specific models.
package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/BenbenIO/gh-quoi/internal/model"
	"github.com/google/go-github/v55/github"
)

// ListOptions specifies filters for retrieving GitHub notifications.
type ListOptions struct {
	All           bool       // Include read notifications in the response
	Participating bool       // Only show notifications user is participating in
	Since         *time.Time // Only show notifications updated after this time
}

// DaysAgo returns a time.Time representing n days in the past from now.
// Useful for constructing the Since parameter in ListOptions.
func DaysAgo(n int) *time.Time {
	t := time.Now().AddDate(0, 0, -n)

	return &t
}

// ListNotifications retrieves all GitHub notifications matching the specified options.
// It automatically handles pagination to fetch all available notifications.
// Returns the notifications as application-specific model objects.
func (c *Client) ListNotifications(
	ctx context.Context,
	opts ListOptions,
) ([]model.Notification, error) {
	apiOpts := &github.NotificationListOptions{
		All:           opts.All,
		Participating: opts.Participating,
	}
	if opts.Since != nil {
		apiOpts.Since = *opts.Since
	}

	// Fetch notifications from GitHub API
	allNotifications := []*github.Notification{}

	for {
		notifications, resp, err := c.gh.Activity.ListNotifications(ctx, apiOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list notifications: %w", err)
		}

		allNotifications = append(allNotifications, notifications...)

		if resp.NextPage == 0 {
			break
		}

		apiOpts.Page = resp.NextPage
	}

	// Convert to internal model
	var results []model.Notification
	for _, n := range allNotifications {
		results = append(results, convertNotification(n))
	}

	return results, nil
}

// convertNotification normalizes go-github's Notification struct.
func convertNotification(n *github.Notification) model.Notification {
	repo := ""
	if n.Repository != nil && n.Repository.FullName != nil {
		repo = *n.Repository.FullName
	}

	subject := ""
	subjectType := ""
	apiURL := ""
	htmlURL := ""

	if n.Subject != nil {
		if n.Subject.Title != nil {
			subject = *n.Subject.Title
		}

		if n.Subject.Type != nil {
			subjectType = *n.Subject.Type
		}

		if n.Subject.URL != nil {
			apiURL = *n.Subject.URL
			htmlURL = apiToHTMLURL(apiURL)
		}
	}

	updated := ""

	if n.UpdatedAt != nil {
		t := n.UpdatedAt.Time
		updated = t.Format("2006-01-02 15h") // golang time format template
	}

	id := ""
	if n.ID != nil {
		id = *n.ID
	}

	reason := ""
	if n.Reason != nil {
		reason = *n.Reason
	}

	unread := false
	if n.Unread != nil {
		unread = *n.Unread
	}

	return model.Notification{
		ID:         id,
		Repository: repo,
		Title:      subject,
		Reason:     reason,
		Type:       subjectType,
		HTMLURL:    htmlURL,
		UpdatedAt:  updated,
		Unread:     unread,
	}
}

// apiToHTMLURL converts API URLs to browser URLs.
//
//	https://api.github.com/repos/owner/repo/issues/123
//
//	https://github.com/owner/repo/issues/123
func apiToHTMLURL(apiURL string) string {
	s := strings.Replace(apiURL, "https://api.", "https://", 1)
	s = strings.Replace(s, "/repos", "", 1)
	s = strings.Replace(s, "/pulls/", "/pull/", 1)

	return s
}
