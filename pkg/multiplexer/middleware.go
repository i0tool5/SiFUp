package multiplexer

import (
	"log"
	"net/http"
)

func printReq(r *http.Request) {
	log.Printf("[*] %s | %s\n", r.RemoteAddr, r.UserAgent())
}

func MiddlewareCliLogger(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		printReq(r)
		f.ServeHTTP(rw, r)
	}
}
