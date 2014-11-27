package slave

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const testURL = "http://www.placekitten.com"
const testBrowser = false

// TODO: These tests can be greatly improved.
func TestBrowserHandler(t *testing.T) {
	OS := "Some Unknown OS"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		BrowserHandler(w, r, OS)
	}))

	client := &http.Client{}
	client.PostForm(testServer.URL, url.Values{"url": {testURL}})
}

func TestKillBrowserOS_X(t *testing.T) {
	if testBrowser {
		killBrowser("OS X")
	}
}

func TestKillBrowserLinux(t *testing.T) {
	if testBrowser {
		killBrowser("Linux")
	}
}

func TestOpenBrowserOS_X(t *testing.T) {
	if testBrowser {
		openBrowser("OS X", testURL)
	}
}

func TestOpenBrowserLinux(t *testing.T) {
	if testBrowser {
		openBrowser("Linux", testURL)
	}
}
