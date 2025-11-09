package monitoring

import (
	"os"
	"testing"
)

func TestNewPosthogMonitor(t *testing.T) {
	monitor := NewPosthogMonitor("hello", "not-an-endpoint")
	_, isType := any(monitor).(PosthogClient)
	if !isType {
		t.Error("Expecting PosthogMonitor to implement PosthogClient, but it does not")
	}
}

func TestPosthogMonitorGetClient(t *testing.T) {
	apiKey, okApi := os.LookupEnv("POSTHOG_API_KEY")
	endpoint, okEnd := os.LookupEnv("POSTHOG_ENDPOINT")
	if !okApi || !okEnd {
		t.Skip()
	} else {
		monitor := NewPosthogMonitor(apiKey, endpoint)
		_, err := monitor.GetClient()
		if err != nil {
			t.Errorf("Expecting no error, got %s", err.Error())
		}
	}
}

func TestPosthogSendEvent(t *testing.T) {
	apiKey, okApi := os.LookupEnv("POSTHOG_API_KEY")
	endpoint, okEnd := os.LookupEnv("POSTHOG_ENDPOINT")
	if !okApi || !okEnd {
		t.Skip()
	} else {
		monitor := NewPosthogMonitor(apiKey, endpoint)
		err := monitor.SendEvent("test", "testEvent", "testPosthogSendEvent", 0, false, "")
		if err != nil {
			t.Errorf("Expecting no error, got %s", err.Error())
		}
	}
}
