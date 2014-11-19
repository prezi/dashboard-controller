package slaveModule

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"net/url"
)

const testURL = "http://www.placekitten.com"

func TestBrowserHandler(t *testing.T) {
	OS := "Some Unknown OS"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		BrowserHandler(w, r, OS)
	}))

	client := &http.Client{}
	client.PostForm(testServer.URL, url.Values{"url": {testURL}})
}

func TestKillBrowserOS_X(t *testing.T) {
	killBrowser("OS X")
}

func TestKillBrowserLinux(t *testing.T) {
	killBrowser("Linux")
}

func TestOpenBrowserOS_X(t *testing.T) {
	openBrowser("OS X", testURL)
}

func TestOpenBrowserLinux(t *testing.T) {
	openBrowser("Linux", testURL)
}
