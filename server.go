package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cagrosso/OauthPlayground/helpers"
	"github.com/google/uuid"
)

var (
	serverConfig        helpers.ServerConfiguration
	httpClient          http.Client
	sessionTrackerCache map[string]sessionTracker
)

const (
	sessionTokenConst = "session_token"
)

type oAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	var err error
	serverConfig, err = helpers.GetConfiguration()
	if err != nil {
		panic(err)
	}

	httpClient = http.Client{}
	sessionTrackerCache = make(map[string]sessionTracker)

	http.HandleFunc("/oauth/redirect", oauthRedirectHandler)
	http.HandleFunc("/api/user", userInfoHandler)

	fmt.Printf("Listening on %v ...\n", serverConfig.Port)
	http.ListenAndServe(serverConfig.Port, nil)
}

func oauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/oauth/redirect")

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	code := r.FormValue("code")

	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		serverConfig.ClientID,
		serverConfig.ClientSecret,
		code,
	)
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

func userInfoHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/api/user")

	sessionToken, err := req.Cookie(sessionTokenConst)
	if err != nil {
		fmt.Fprintf(os.Stdout, "no cookie attached to request: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	authorizationToken := sessionTrackerCache[sessionToken.Value].AccessToken

	url, _ := url.Parse("https://api.github.com/user")
	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("token %s", authorizationToken))

	userRequest := http.Request{
		Method: "GET",
		Header: headers,
		URL:    url,
	}

	resp, err := httpClient.Do(&userRequest)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to make Github user api request: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not read body: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("Body: %+v\n", string(body))
	w.Write(body)
	w.WriteHeader(http.StatusOK)
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
		Path:    "/",
	})
}
