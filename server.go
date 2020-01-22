package main

import (
	"fmt"
	"net/http"
)

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hit /")
	fmt.Fprintf(w, "Hello, World!\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hit /headers")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	port := ":8080"
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/headers", headers)

	fmt.Printf("Listening on %v...\n", port)
	http.ListenAndServe(port, nil)
}
