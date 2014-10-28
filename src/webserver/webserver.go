package main

import (
	"fmt"
	"net/http"
	//"html"
	"io/ioutil"
	"log"
	"mime"
	"path/filepath"
)

type String string

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (t String) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, t)
}

func (t Struct) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, t.Greeting, t.Punct, t.Who)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path=="/"){
			http.Redirect(w,r,"/form.html",301)
		}

		mimetype:=mime.TypeByExtension(filepath.Ext(r.URL.Path[1:]))
		w.Header().Set("Content-type", mimetype)
		html_file,err:=ioutil.ReadFile(r.URL.Path[1:])

		if (err!=nil){
			log.Fatal(err)
		}else{
			fmt.Fprint(w, String(html_file))
		}		
	})
		http.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
		URL:=r.FormValue("url")
		rb_ID:=r.FormValue("rb-id")
		fmt.Println(URL,rb_ID)
		http.Redirect(w,r,"/form.html",301)
	})

	http.ListenAndServe("localhost:4000", nil)
}
