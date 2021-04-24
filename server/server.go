package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	fdir string
	wfo  = "Waiting for connections on"
)

func errorWrapCli(err error) {
	if err != nil {
		log.Printf("[!] Error occured: %s\n", err)
	}
}

func errorWrapHttp(w http.ResponseWriter) {
	errTempl, err := ioutil.ReadFile("templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", errTempl)
}

func curDateTime() string {
	t := time.Now()
	h, m, s := t.Clock()
	y, M, d := t.Date()

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", y, M, d, h, m, s)
}

func createUploadDir(dirname string) {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		fmt.Println("[!] Creating uploads directory")
		err := os.Mkdir(dirname, 0744)
		errorWrapCli(err)
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[*] %s|%s|[%s]\n", r.RemoteAddr, r.UserAgent(), curDateTime())
	data, err := ioutil.ReadFile("templates/formtemplate.html")
	errorWrapCli(err)
	template := fmt.Sprintf("%s", data)
	fmt.Fprint(w, template)
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[*] File|%s|%s|[%s]\n", r.RemoteAddr, r.UserAgent(), curDateTime())
	err := r.ParseMultipartForm(1_000) //1KB for now
	if err != nil {
		errorWrapCli(err)
		errorWrapHttp(w)
		return
	}

	files := r.MultipartForm.File
	if len(files["upfile"]) == 0 {
		log.Println("[!] No files do save.")
		errorWrapHttp(w)
		return
	}
	for _, file := range files["upfile"] {
		cf, err := file.Open()
		if err != nil {
			errorWrapCli(err)
			errorWrapHttp(w)
			return
		}
		defer cf.Close()
		fileContent, _ := ioutil.ReadAll(cf)
		_ = ioutil.WriteFile(fdir+file.Filename, fileContent, 0644)
	}
	template, err := ioutil.ReadFile("templates/success.html")
	errorWrapCli(err)
	fmt.Fprintf(w, "%s", template)
}

// UploadServer serves
func UploadServer(addr string, saveDir string, tls bool) {
	fdir = saveDir

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/upload", handleFiles)
	fmt.Printf("[!] Starting server %s\n", curDateTime())
	createUploadDir(saveDir)
	if tls {
		fmt.Printf("[*] %s https://%s\n", wfo, addr)
		log.Fatal(http.ListenAndServeTLS(addr, "tlscert/cert.pem", "tlscert/key.pem", nil))
	} else {
		fmt.Printf("[*] %s http://%s\n", wfo, addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}
