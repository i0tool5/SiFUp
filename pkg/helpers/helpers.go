package helpers

import (
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/i0tool5/simpleuploader/pkg/templates"
)

func WrapCli(err error) {
	if err != nil {
		log.Printf("[!] Error occured: %s\n", err)
	}
}

func WrapHttp(rw http.ResponseWriter) {
	errTempl, err := templates.Temlpates.ReadFile("html/error.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(rw, string(errTempl))
}

// WrapBoth to Http response and stdout
func WrapBoth(rw http.ResponseWriter, err error) {
	WrapCli(err)
	WrapHttp(rw)
}

func CreateUploadsDir(dirname string, logger *slog.Logger) (err error) {
	if _, err = os.Stat(dirname); os.IsNotExist(err) {
		logger.Warn(
			"creating uploads directory",
			slog.Attr{
				Key:   "dirname",
				Value: slog.AnyValue(dirname),
			},
		)
		err = os.Mkdir(dirname, 0744)
		return err
	}
	return
}

func CreateFile(dir, fileName string, data []byte, logger *slog.Logger) (err error) {
	if err = CreateUploadsDir(dir, logger); err != nil {
		return
	}
	fn := dir + fileName // fn stands for full name
	err = os.WriteFile(fn, data, 0644)
	return
}

func PrepareTLSCert(cert, key []byte) (tls.Certificate, error) {
	return tls.X509KeyPair(cert, key)
}
