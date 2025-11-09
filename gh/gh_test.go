package gh

import (
	"os"
	"testing"
)

func TestNewGitYearClient(t *testing.T) {
	client := NewGitYearClient("thisIsNotaToken")
	_, isType := any(client).(GitHubClient)
	if !isType {
		t.Error("Expecting GitYearClient to implement GitHub client, but it does not")
	}
}

func TestGetUserStats(t *testing.T) {
	tok, ok := os.LookupEnv("GITHUB_AUTH_TOKEN")
	if !ok {
		t.Skip()
	} else {
		client := NewGitYearClient(tok)
		stats, err := client.GetUserStats("torvalds")
		if err != nil {
			t.Errorf("Not expecting any error, got %s", err.Error())
		}
		if stats.User != "torvalds" {
			t.Errorf("Expecting the user to be 'torvalds', got '%s'", stats.User)
		}
	}
}

func TestGetOrgStats(t *testing.T) {
	tok, ok := os.LookupEnv("GITHUB_AUTH_TOKEN")
	if !ok {
		t.Skip()
	} else {
		client := NewGitYearClient(tok)
		stats, err := client.GetOrgStats("indigo-notes")
		if err != nil {
			t.Errorf("Not expecting any error, got %s", err.Error())
		}
		if stats.Organization != "indigo-notes" {
			t.Errorf("Expecting the user to be 'indigo-notes', got '%s'", stats.Organization)
		}
	}
}
