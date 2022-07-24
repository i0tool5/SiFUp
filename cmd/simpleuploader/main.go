package main

import (
	"flag"
	"log"

	"github.com/i0tool5/simpleuploader/pkg/helpers"
	"github.com/i0tool5/simpleuploader/pkg/multiplexer"
	"github.com/i0tool5/simpleuploader/pkg/server"
)

const (
	defaultFileDir = "./uploaded/"
	defaultAddr    = "127.0.0.1:58000"
)

var (
	bindAddr string
	saveDir  string
	useTLS   bool
)

func init() {
	flag.StringVar(&bindAddr, "bind", defaultAddr, "Set host:port to listen on.")
	flag.StringVar(&saveDir, "save-to", defaultFileDir, "Set directory for uploaded files.")
	flag.BoolVar(&useTLS, "tls", false, "Should server use tls or not (default 'not')")
	flag.Parse()
}

func run() {
	helpers.CreateUploadsDir(saveDir)

	//
	// routing config
	//

	mplex := multiplexer.New()
	handlers := multiplexer.NewHandlers(saveDir)

	if err := handlers.PrepareTemplates(); err != nil {
		log.Fatal(err)
	}

	mplex.HandleFunc("/", multiplexer.MiddlewareCliLogger(handlers.HandleMain))
	mplex.HandleFunc("/upload", multiplexer.MiddlewareCliLogger(handlers.HandleFiles))

	//
	// server
	//

	srv := server.New(bindAddr, mplex)

	if useTLS {
		log.Printf("Starting TLS server on %s\n", bindAddr)
		srv.ListenAndServeTLS()
		log.Println("TLS server done")
	} else {
		log.Printf("Starting server on %s\n", bindAddr)
		srv.ListenAndServe()
		log.Println("Server done")
	}
}

func main() {
	run()
}
