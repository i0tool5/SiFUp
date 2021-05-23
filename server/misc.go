package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func printReq(r *http.Request) {
	fmt.Printf("[*] %s|%s|[%s]\n", r.RemoteAddr, r.UserAgent(), curDateTime())
}

func MiddlewareCliLogger(f http.HandlerFunc) http.HandlerFunc {
	internal := func(rw http.ResponseWriter, r *http.Request) {
		printReq(r)
		f.ServeHTTP(rw, r)
	}
	return internal
}

type ErrorWrapper interface {
	WrapCli(error)
	WrapHttp(http.ResponseWriter)
	WrapBoth(http.ResponseWriter, error)
}

type WrapError struct{}

func (w WrapError) WrapCli(err error) {
	if err != nil {
		log.Printf("[!] Error occured: %s\n", err)
	}
}

func (w WrapError) WrapHttp(rw http.ResponseWriter) {
	errTempl, err := ioutil.ReadFile("templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(rw, string(errTempl))
}

func (w WrapError) WrapBoth(rw http.ResponseWriter, err error) {
	w.WrapCli(err)
	w.WrapHttp(rw)
}
