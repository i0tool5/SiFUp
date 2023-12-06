package main

import (
	"flag"
	"log/slog"
	"os"

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
	debug    bool
)

func init() {
	flag.StringVar(&bindAddr, "bind", defaultAddr, "Set host:port to listen on.")
	flag.StringVar(&saveDir, "save-to", defaultFileDir, "Set directory for uploaded files.")
	flag.BoolVar(&useTLS, "tls", false, "Should server use tls or not (default 'not')")
	flag.BoolVar(&debug, "debug", true, "Show debug level output. Might be verbose (default: true)")
	flag.Parse()
}

func run() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	err := helpers.CreateUploadsDir(saveDir, logger)
	if err != nil {
		logger.Error("can't create directory for files", err)
		os.Exit(1)
	}

	//
	// routing config
	//

	mplex := multiplexer.New()
	handlers := multiplexer.NewHandlers(logger, saveDir)
	middlewares := multiplexer.NewMiddleware(logger)

	if err := handlers.PrepareTemplates(); err != nil {
		logger.Error("can't prepare templates", err)
		os.Exit(1)
	}

	mplex.HandleFunc("/", middlewares.CliLogger(handlers.HandleMain))
	mplex.HandleFunc("/upload", middlewares.CliLogger(handlers.HandleFiles))
	mplex.HandleFunc("/fonts", middlewares.CliLogger(handlers.HandleFonts))

	//
	// server
	//

	srv := server.New(bindAddr, mplex)

	if useTLS {
		logger.Info(
			"starting TLS server",
			slog.Any("address", bindAddr),
		)
		srv.ListenAndServeTLS()
		logger.Info(
			"stoping TLS server",
			slog.Any("address", bindAddr),
		)
	} else {
		logger.Info(
			"starting server",
			slog.Any("address", bindAddr),
		)
		srv.ListenAndServe()
		logger.Info(
			"stoping server",
			slog.Any("address", bindAddr),
		)
	}
}

func main() {
	run()
}
