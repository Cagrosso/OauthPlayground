package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/")
	dat, err := ioutil.ReadFile("./public/index.html")
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not read index page: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(dat))
}
