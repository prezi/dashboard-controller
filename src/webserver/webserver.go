package main

import (
	"fmt"
	"net/http"
	"html"
	"io/ioutil"
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

// func ()ServeHTTP(
// 	w http.ResponseWriter,
// 	r *http.Request){

// 	fmt.Fprint(w,"Hi!")
// }

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		 
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		html_file,err:=ioutil.ReadFile("form.html")
		if (err==nil){
			fmt.Fprint(w, String(html_file))
		}		
	})
		http.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	//http.Handle("/")
	//http.Handle("/",String("what"))
	//http.Handle("/string", String("I'm a frayed knot."))
	//http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})

	http.ListenAndServe("localhost:4000", nil)
}
