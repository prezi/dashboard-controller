//http://www.golang-book.com/13/

package main

import ("net/http" ; "io" ; "log"; "encoding/json"; "io/ioutil")

type Message struct { // this will be the json: { "ID": "1", "URL": "http://google.com"}
    ID string
    URL string
}

func hello(res http.ResponseWriter, req *http.Request) {
    res.Header().Set(
        "Content-Type", 
        "text/html",
    )
    io.WriteString(
        res, 
        `<doctype html>
<html>
    <head>
        <title>Hello World</title>
    </head>
    <body>
        Hello World!
    </body>
</html>`,
    )
}

func write(rw http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))
    var t Message
    err = json.Unmarshal(body, &t)
    if err != nil {
        panic(err)
    }
    log.Println(t.URL)
    // decoder := json.NewDecoder(req.Body)
    // var t Message 
    // err := decoder.Decode(&t)
    // if err != nil {
    //     panic(err)
    // }
    // log.Println(t.Test)
}

func main() {
    http.HandleFunc("/", write)
    http.HandleFunc("/hello", hello)
    http.ListenAndServe(":4005", nil)
}