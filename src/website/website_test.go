package website

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"master/master"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"network"
)


type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

func TestIndexPageHandler(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		IndexPageHandler(w, request)
	}))
	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-type"))
}

func TestIndexPageHandlerWithWrongPath(t *testing.T) {
	VIEWS_PATH = "dummy"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		IndexPageHandler(w, request)
	}))
	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)
	assert.Equal(t, 404, resp.StatusCode)
	VIEWS_PATH = network.PROJECT_ROOT + "/src/website/views"
}

func sendGetToFormHandler(URL string) int {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		request.URL.Path = URL
		testSlaveNames := []string{"a", "b"}
		FormHandler(w, request, testSlaveNames)
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	return resp.StatusCode
}

func TestFormHandler(t *testing.T) {
	assert.Equal(t, 302, sendGetToFormHandler("/"))
}

func TestStatusMessageForAvailableServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	assert.Equal(t, true, master.IsURLValid(testServer.URL))
}

func TestStatusMessageForUnavailableServer(t *testing.T) {
	assert.Equal(t, false, master.IsURLValid(""))
}

func TestSubmitHandlerWithEmptyREsponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		SubmitHandler(w, request, map[string]master.Slave{})
	}))

	client := &http.Client{}
	resp, _ := client.Post(testServer.URL, "application/json", ioutil.NopCloser(strings.NewReader("")))
	byteContent, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(byteContent), "Failed"))
}

func TestDisplayFormPage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		displayFormPage(w, []string{"foo", "bar", "baz"}, "testUser")
	}))
	resp, _ := http.Get(testServer.URL)

	defer resp.Body.Close()
	byteContent, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(byteContent), "foo"))
	assert.True(t, strings.Contains(string(byteContent), "bar"))
	assert.True(t, strings.Contains(string(byteContent), "baz"))
	assert.True(t, strings.Contains(string(byteContent), "testUser"))
}

func TestSubmitHandlerWithWrongHttp(t *testing.T) {
	type FormData struct {
		URLToDisplay   string
		SlaveNames []string
	}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		SubmitHandler(w, request, map[string]master.Slave{})
	}))
	form := FormData{"dummy",[]string{"a", "b"}}
	b, err := json.Marshal(form)

	req, err := http.NewRequest("POST", testServer.URL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byteContent, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(byteContent), "Sadpanda"))
}

func TestSubmitHandlerWithNonExistentSlave(t *testing.T) {
	type FormData struct {
		URLToDisplay   string
		SlaveNames []string
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		SubmitHandler(w, request, map[string]master.Slave{})
	}))
	form := FormData{"http://www.google.com",[]string{"a", "b"}}
	b, err := json.Marshal(form)

	req, err := http.NewRequest("POST", testServer.URL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byteContent, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(byteContent), "refresh"))
}

func TestSubmitHandler(t *testing.T) {
	type FormData struct {
		URLToDisplay   string
		SlaveNames []string
	}

	var receivedUrl string

	testSlave := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receivedUrl = request.PostFormValue("url")
	}))

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		SubmitHandler(w, request, map[string]master.Slave{"a":master.Slave{
			URL : testSlave.URL,
			PreviouslyDisplayedURL: "",
			DisplayedURL : "",
		}})
	}))
	form := FormData{"http://www.google.com", []string{"a"}}
	b, _ := json.Marshal(form)

	req, _ := http.NewRequest("POST", testServer.URL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	byteContent, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(byteContent), "Slaves are updated"))
	assert.Equal(t, receivedUrl, "http://www.google.com")
}

func TestParseFromJSON(t *testing.T) {
	type FormData struct {
		URLToDisplay string
		SlaveNames   []string
	}
	testSlaveList := []string{"a", "b", "c"}
	testFormData := FormData{"testurl", testSlaveList}
	testJSON, err := json.Marshal(testFormData)

	returnedURL, returnedSlaveList, err := parseFromJSON(ioutil.NopCloser(bytes.NewReader(testJSON)))
	assert.Equal(t, "testurl", returnedURL)
	assert.Equal(t, []string{"a", "b", "c"}, returnedSlaveList)
	assert.Nil(t, err)
}

func TestSendConfirmationMessageToUser(t *testing.T) {
	testMessage := "testmessage"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sendConfirmationMessageToUser(w, testMessage)
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)
	respBodyContents, _ := ioutil.ReadAll(resp.Body)
	respBodyString := string(respBodyContents[:])
	assert.True(t, strings.Contains(respBodyString, testMessage))
	assert.Equal(t, "application/json", resp.Header.Get("Content-type"))
}

func TestCreateConfirmationMessage(t *testing.T) {
	msg := "testmessage"
	confirmationMessageJson, _ := createConfirmationMessage(msg)
	confirmationMessageJsonString := string(confirmationMessageJson[:])
	assert.True(t, strings.Contains(confirmationMessageJsonString, msg))
}

func TestCreateConfirmationMessageWithWrongPath(t *testing.T) {
	VIEWS_PATH = "dummy"
	msg := "testmessage"
	confirmationMessageJson, _ := createConfirmationMessage(msg)
	assert.Equal(t, len(confirmationMessageJson), 0)
	VIEWS_PATH = network.PROJECT_ROOT + "/src/website/views"
}
