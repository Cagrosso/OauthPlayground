package server

import (
	"fmt"
	"net/http"
)

func headersHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("/headers")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
