package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func welcomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/welcome")

	dat, err := ioutil.ReadFile("./public/welcome.html")
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not read welcome page: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(dat))
}
