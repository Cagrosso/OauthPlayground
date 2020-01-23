package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	port                string
	clientID            string
	clientSecret        string
	httpClient          http.Client
	sessionTrackerCache map[string]sessionTracker
)

const (
	sessionTokenConst = "session_token"
)

type oAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

type sessionTracker struct {
	AccessToken string
	TimeOut     time.Time
}

func main() {
	err := getConfiguration()
	if err != nil {
		panic(err)
	}

	httpClient = http.Client{}
	sessionTrackerCache = make(map[string]sessionTracker)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/welcome.html", welcomeHandler)
	http.HandleFunc("/headers", headersHandler)
	http.HandleFunc("/oauth/redirect", oauthRedirectHandler)

	fmt.Printf("Listening on %v ...\n", port)
	http.ListenAndServe(port, nil)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/")
	dat, err := ioutil.ReadFile("./public/index.html")
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not read index page: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(dat))
}

func welcomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/welcome")

	cookie, err := req.Cookie(sessionTokenConst)
	if err != nil {
		fmt.Fprintf(os.Stdout, "no cookie attached: %+v", err)
		//w.WriteHeader(http.StatusBadRequest)
	} else {
		sessionToken := cookie.Value
		fmt.Println("Cookie " + sessionToken)
	}

	dat, err := ioutil.ReadFile("./public/welcome.html")
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not read welcome page: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(dat))
}

func headersHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/headers")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func oauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/oauth/redirect")

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

	attachSessionCookieToResponseWriter(w, t.AccessToken)

	http.Redirect(w, r, "/welcome.html", http.StatusFound)
}

func attachSessionCookieToResponseWriter(w http.ResponseWriter, accessToken string) {
	newSession := sessionTracker{
		AccessToken: accessToken,
		TimeOut:     time.Now().Add(time.Minute * 15),
	}

	sessionToken := uuid.New().String()

	sessionTrackerCache[sessionToken] = newSession

	http.SetCookie(w, &http.Cookie{
		Name:    sessionTokenConst,
		Value:   sessionToken,
		Expires: newSession.TimeOut,
		Domain:  "localhost",
	})
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
