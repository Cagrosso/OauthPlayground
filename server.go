package main

import (
	"fmt"
	"net/http"
	"os"
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

func getEnv(name string) (string, error) {
	envValue := os.Getenv(name)
	if envValue == "" {
		return "", fmt.Errorf("Env '%v' not set", name)
	}

	return envValue, nil
}

func getPort(portNumber string) string {
	return fmt.Sprintf(":%v", portNumber)
}

func main() {
	port, err := getEnv("PORT")
	if err != nil {
		fmt.Printf("Failed to start, %+v\n", err)
		os.Exit(1)
	}
	port = getPort(port)

	clientID, err := getEnv("CLIENT_ID")
	if err != nil {
		fmt.Printf("Failed to start, %+v\n", err)
		os.Exit(1)
	}
	clientSecret, err := getEnv("CLIENT_SECRET")
	if err != nil {
		fmt.Printf("Failed to start, %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("clientID: %v\n", clientID)
	fmt.Printf("clientSecret: %v\n", clientSecret)

	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/headers", headers)

	fmt.Printf("Listening on %v ...\n", port)
	http.ListenAndServe(port, nil)
}
