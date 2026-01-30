package model

import (
	"testing"
)

func TestNotification_StructInitialization(t *testing.T) {
	want := Notification{
		ID:         "test-id",
		Repository: "owner/repo",
		Title:      "Test notification",
		Reason:     "mention",
		Type:       "Issue",
		HTMLURL:    "https://github.com/owner/repo/issues/1",
		UpdatedAt:  "2023-01-01 12h",
		Unread:     true,
	}

	if want.ID != "test-id" {
		t.Errorf("ID = %q, want %q", want.ID, "test-id")
	}

	if want.Repository != "owner/repo" {
		t.Errorf("Repository = %q, want %q", want.Repository, "owner/repo")
	}

	if want.Title != "Test notification" {
		t.Errorf("Title = %q, want %q", want.Title, "Test notification")
	}

	if want.Reason != "mention" {
		t.Errorf("Reason = %q, want %q", want.Reason, "mention")
	}

	if want.Type != "Issue" {
		t.Errorf("Type = %q, want %q", want.Type, "Issue")
	}

	if want.HTMLURL != "https://github.com/owner/repo/issues/1" {
		t.Errorf("HTMLURL = %q, want %q", want.HTMLURL, "https://github.com/owner/repo/issues/1")
	}

	if want.UpdatedAt != "2023-01-01 12h" {
		t.Errorf("UpdatedAt = %q, want %q", want.UpdatedAt, "2023-01-01 12h")
	}

	if !want.Unread {
		t.Errorf("Unread = %t, want %t", want.Unread, true)
	}
}

func TestNotification_ZeroValue(t *testing.T) {
	var got Notification

	if got.ID != "" {
		t.Errorf("ID = %q, want empty string", got.ID)
	}

	if got.Repository != "" {
		t.Errorf("Repository = %q, want empty string", got.Repository)
	}

	if got.Title != "" {
		t.Errorf("Title = %q, want empty string", got.Title)
	}

	if got.Reason != "" {
		t.Errorf("Reason = %q, want empty string", got.Reason)
	}

	if got.Type != "" {
		t.Errorf("Type = %q, want empty string", got.Type)
	}

	if got.HTMLURL != "" {
		t.Errorf("HTMLURL = %q, want empty string", got.HTMLURL)
	}

	if got.UpdatedAt != "" {
		t.Errorf("UpdatedAt = %q, want empty string", got.UpdatedAt)
	}

	if got.Unread {
		t.Errorf("Unread = %t, want false", got.Unread)
	}
}
