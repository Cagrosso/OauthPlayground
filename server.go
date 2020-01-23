package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	port         string
	clientID     string
	clientSecret string
	httpClient   http.Client
)

type oAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	err := getConfiguration()
	if err != nil {
		panic(err)
	}

	httpClient = http.Client{}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/headers", headersHandler)
	http.HandleFunc("/oauth/redirect", oauthRedirectHandler)

	fmt.Printf("Listening on %v ...\n", port)
	http.ListenAndServe(port, nil)
}

func headersHandler(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func oauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirected...")

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	code := r.FormValue("code")

	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	req.Header.Set("accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	var t oAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
	w.WriteHeader(http.StatusFound)
}

func getConfiguration() error {
	var err error
	port, err = getEnv("PORT")
	if err != nil {
		return err
	}
	port = getPort(port)

	clientID, err = getEnv("CLIENT_ID")
	if err != nil {
		return err
	}

	clientSecret, err = getEnv("CLIENT_SECRET")
	if err != nil {
		return err
	}

	return nil
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
