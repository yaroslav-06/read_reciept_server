package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    buf, err := ioutil.ReadFile("assets/pixel.png")
	if err != nil{
		log.Fatal(err)
	}
	log.Println(buf);
    handler := http.HandlerFunc(handleRequest)

    http.Handle("/image", handler)

    fmt.Println("Server started at port 2000")
    http.ListenAndServe(":2000", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

    buf, err := ioutil.ReadFile("assets/pixel.png")

    if err != nil {

        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "image/png")
    w.Write(buf)
}
