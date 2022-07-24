package server

import (
	"crypto/tls"
	"embed"
	"log"
	"net"
	"net/http"

	"github.com/i0tool5/simpleuploader/pkg/helpers"
	"github.com/i0tool5/simpleuploader/pkg/tlscert"
)

type Server struct {
	certs    embed.FS
	listener *http.Server
}

func New(addr string, muxer *http.ServeMux) *Server {
	return &Server{
		certs: tlscert.CertFiles,
		listener: &http.Server{
			Addr:    addr,
			Handler: muxer,
		},
	}
}

func (s *Server) ListenAndServe() {
	if err := s.listener.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func (s *Server) ListenAndServeTLS() {
	cert, err := s.certs.ReadFile("files/cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	key, err := s.certs.ReadFile("files/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	certificate, err := helpers.PrepareTLSCert(cert, key)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", s.listener.Addr)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{
			certificate,
		},
	}

	tlsListener := tls.NewListener(ln, tlsCfg)

	if err := s.listener.Serve(tlsListener); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
