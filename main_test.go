package main

import (
	"io"
	"net/http"
	"net/url"
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
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Get statistics about your (or your organization's) GitHub activity in the past year, share them with others and get inspired to build more!") {
		t.Error("Unexpected body in response")
	}
	req, err = http.NewRequest("POST", "/", nil)
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
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Don't have an account?") {
		t.Error("Unexpected body in response")
	}
	req, err = http.NewRequest("POST", "/signin", nil)
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
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Already have an account?") {
		t.Error("Unexpected body in response")
	}
	req, err = http.NewRequest("POST", "/signup", nil)
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
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Discover the treasures in your GitHub journey!") {
		t.Error("Unexpected body in response")
	}
	req, err = http.NewRequest("POST", "/search", nil)
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
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Not expecting error while reading response body, got %s", err.Error())
	}
	if !strings.Contains(string(data), "Ooops, this page does not exist!") {
		t.Error("Unexpected body in response")
	}
}

func TestPublishSocial(t *testing.T) {
	app := Setup()
	req, err := http.NewRequest("GET", "/urls/bsky", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Form = url.Values{}
	req.Form.Add("bskyInput", "input")

	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("Not expecting error while getting response, got %s", err.Error())
	}

	redirectURL := resp.Header.Get("HX-Redirect")
	expected := "https://bsky.app/intent/compose?text=input"
	if redirectURL != expected {
		t.Errorf("Expecting url '%s', got '%s'", expected, redirectURL)
	}
	req, err = http.NewRequest("GET", "/urls/x", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Form = url.Values{}
	req.Form.Add("xInput", "input")

	resp, err = app.Test(req, 5000)
	if err != nil {
		t.Fatalf("Not expecting error while getting response, got %s", err.Error())
	}

	redirectURL = resp.Header.Get("HX-Redirect")
	expected = "https://twitter.com/intent/tweet?text=input"
	if redirectURL != expected {
		t.Errorf("Expecting url '%s', got '%s'", expected, redirectURL)
	}
}
