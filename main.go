package main

import (
	"flag"

	"github.com/i0tool5/simpleuploader/server"
)

const (
	defaultFileDir = "./uploaded/"
	defaultAddr    = "127.0.0.1:58000"
)

var (
	bindAddr string
	useTLS   bool
)

func init() {
	flag.StringVar(&bindAddr, "bind", defaultAddr, "Set host:port to listen on.")
	flag.BoolVar(&useTLS, "tls", false, "Should server use tls or not (default 'not')")
	flag.Parse()
}

func run() {
	listenAddr := defaultAddr

	if bindAddr != "" {
		listenAddr = bindAddr
	}

	server.UploadServer(listenAddr, defaultFileDir, useTLS)
}

func main() {
	run()
}
