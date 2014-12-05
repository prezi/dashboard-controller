package session

import (
	"github.com/gorilla/securecookie"
	"master/master"
	"master/master/website/hash"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var (
	FILE_PATH_TO_USER_AUTHENTICATION_DATA = master.GetRelativeFilePath("./user_authentication_data.txt")
	USER_AUTHENTICATION_MAP               = hash.InitializeUserAuthenticationMap(FILE_PATH_TO_USER_AUTHENTICATION_DATA)
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	password := r.FormValue("password")
	redirectTarget := "/"

	if name != "" && password != "" {
		if hash.IsHashMatchInUserAuthenticationMap(name, password, USER_AUTHENTICATION_MAP) {
			setSession(name, w)
			redirectTarget = "/internal"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func setSession(username string, response http.ResponseWriter) {
	value := map[string]string{
		"name": username,
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

func GetUsername(request *http.Request) (username string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
		}
	}
	return username
}
