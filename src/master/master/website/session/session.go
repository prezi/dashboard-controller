// http://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout.html

// simple example for session management using secure cookies

// Serve two pages - an index page providing a login form and
// an internal page that is only accessible to authenticated users (= users that have used the login form).
// The internal page provides a possibility to log out.

// This has to be implemented using only the Golang standard packages and the Gorilla toolkit.

package session

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
)

const (
	USERNAME = "p"
	PASSWORD = "p"
)

// A secure cookie handler is initialized.
// The required parameters (hashKey and blockKey) are generated randomly.
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()

// func main() {
// 	router.HandleFunc("/", indexPageHandler)
// 	router.HandleFunc("/internal", internalPageHandler)

// 	router.HandleFunc("/login", loginHandler).Methods("POST")
// 	router.HandleFunc("/logout", logoutHandler).Methods("POST")

// 	http.Handle("/", router)
// 	http.ListenAndServe(":8080", nil)

// }

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexPage)
}

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	if userName != "" {
		fmt.Fprintf(w, internalPage, userName)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	redirectTarget := "/"

	if name != "" && password != "" {
		// check credentials
		if name == USERNAME && password == PASSWORD {
			setSession(name, w)
			redirectTarget = "/internal"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
