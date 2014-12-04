package session

import (
	// "fmt"
	"github.com/gorilla/securecookie"
	"net/http"
)

const (
	USERNAME = "Prezi"
	PASSWORD = "prezi"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

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

// func IndexPageHandler(w http.ResponseWriter, r *http.Request) {

// 	template, err := template.ParseFiles(path.Join(VIEWS_PATH, "login.html"))
// 	network.ErrorHandler(err, "Error encountered while parsing website template files: %v.")
// 	template.Execute(w, "Login Error Message Here.")
// 	// fmt.Fprintf(w, indexPage)
// }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	password := r.FormValue("password")
	redirectTarget := "/"

	if name != "" && password != "" {
		if name == USERNAME && password == PASSWORD {
			setSession(name, w)
			redirectTarget = "/internal"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

// func LoginHandler(responseWriter http.ResponseWriter, request *http.Request) {
// 	if request.Method == "GET" {
// 		template, err := template.ParseFiles(path.Join(VIEWS_PATH, "login.html"))
// 		network.ErrorHandler(err, "Error encountered while parsing website template files: %v.")
// 		template.Execute(responseWriter, "Login Error Message Here.")
// 	}
// 	if request.Method == "POST" {
// 		username := request.FormValue("username")
// 		password := request.FormValue("password")
// 		fmt.Println(username, password)
// 		http.Redirect(responseWriter, request, "/", 301)
// 	}
// }

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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
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

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}
