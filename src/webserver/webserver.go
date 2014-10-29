package main

import (
	"fmt"
	"net/http"
	"log"
	"mime"
	"path/filepath"
	"html/template"
)

func statusCode(link string) (int) {
	response, err := http.Head(link)
	if (err!=nil) {
		return 0
	} else {
		return response.StatusCode
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if (r.URL.Path=="/"){
			http.Redirect(w,r,"/form.html",301)
		} else {
			mimetype:=mime.TypeByExtension(filepath.Ext(r.URL.Path[1:]))
			w.Header().Set("Content-type", mimetype)
			t, err := template.ParseFiles(r.URL.Path[1:]) 
			if (err!=nil){
				log.Fatal(err)
			} else {
				t.Execute(w, nil)
			}	
		}
	} else {
		URL:=r.FormValue("url")
		rb_ID:=r.FormValue("rb-id")
		fmt.Println(URL,rb_ID)
		fmt.Println(statusCode(URL))
		http.Redirect(w,r,"/form.html",301)
	}
	
}

func main() {

	http.HandleFunc("/",formHandler)

	http.ListenAndServe("localhost:4003", nil)
}
