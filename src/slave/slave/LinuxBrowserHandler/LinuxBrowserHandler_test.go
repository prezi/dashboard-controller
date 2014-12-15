package LinuxBrowserHandler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"testing"
)

const testURL = "http://www.placekitten.com"

// set testBrowser to true for testing
const testBrowser = false

var browserProcess *exec.Cmd

// TODO: These tests can be greatly improved.
func TestBrowserHandler(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		BrowserHandler(w, r)

	}))


	client := &http.Client{}
	client.PostForm(testServer.URL, url.Values{"url": {testURL}})
}

func TestKillBrowserLinux(t *testing.T) {
	if testBrowser {
		killBrowser()
	}
}

func TestOpenBrowserLinux(t *testing.T) {
	if testBrowser {
		openBrowser(testURL)
	}
}
