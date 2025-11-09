package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHomeRoute(t *testing.T) {
	app := Setup()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	body := resp.Body
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Get statistics about your (or your organization's) GitHub activity in the past year, share them with others and get inspired to build more!") {
		t.Error("Unexpected body in response")
	}
	req, _ = http.NewRequest("POST", "/", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	if resp.StatusCode != 405 {
		t.Error("Expecting a POST request on a GET route to be blocked, it was not")
	}
}

func TestSigninRoute(t *testing.T) {
	app := Setup()
	req, err := http.NewRequest("GET", "/signin", nil)
	if err != nil {
		return
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	body := resp.Body
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Don't have an account?") {
		t.Error("Unexpected body in response")
	}
	req, _ = http.NewRequest("POST", "/signin", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	if resp.StatusCode != 405 {
		t.Error("Expecting a POST request on a GET route to be blocked, it was not")
	}
}

func TestSignUpRoute(t *testing.T) {
	app := Setup()
	req, err := http.NewRequest("GET", "/signup", nil)
	if err != nil {
		return
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	body := resp.Body
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Already have an account?") {
		t.Error("Unexpected body in response")
	}
	req, _ = http.NewRequest("POST", "/signup", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	if resp.StatusCode != 405 {
		t.Error("Expecting a POST request on a GET route to be blocked, it was not")
	}
}

func TestSearchoute(t *testing.T) {
	app := Setup()
	req, err := http.NewRequest("GET", "/search", nil)
	if err != nil {
		return
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	body := resp.Body
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Discover the treasures in your GitHub journey!") {
		t.Error("Unexpected body in response")
	}
	req, _ = http.NewRequest("POST", "/search", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	if resp.StatusCode != 405 {
		t.Error("Expecting a POST request on a GET route to be blocked, it was not")
	}
}

func Test404Route(t *testing.T) {
	app := Setup()
	req, err := http.NewRequest("GET", "/doesnotexist", nil)
	if err != nil {
		return
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Not expecting error while getting response, got %s", err.Error())
	}
	body := resp.Body
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Ooops, this page does not exist!") {
		t.Error("Unexpected body in response")
	}
}

func TestPublishSocialX(t *testing.T) {
	if _, ok := os.LookupEnv("POSTHOG_API_KEY"); !ok {
		t.Skip()
	} else {
		app := Setup()
		req := httptest.NewRequest("GET", "/urls/x?xInput=test+message", nil)

		resp, err := app.Test(req, 5000)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}

		if resp == nil {
			t.Fatal("Response is nil")
		}

		redirectURL := resp.Header.Get("HX-Redirect")
		// Check what you actually get
		t.Logf("Got redirect URL: %s", redirectURL)
	}
}

func TestPublishSocialBsky(t *testing.T) {
	if _, ok := os.LookupEnv("POSTHOG_API_KEY"); !ok {
		t.Skip()
	} else {
		app := Setup()
		req := httptest.NewRequest("GET", "/urls/bsky?bskyInput=test+message", nil)

		resp, err := app.Test(req, 5000)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}

		if resp == nil {
			t.Fatal("Response is nil")
		}

		redirectURL := resp.Header.Get("HX-Redirect")
		// Check what you actually get
		t.Logf("Got redirect URL: %s", redirectURL)
	}
}
