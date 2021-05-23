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
	ew   WrapError
)

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
		ew.WrapCli(err)
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	template, err := ioutil.ReadFile("templates/formtemplate.html")
	ew.WrapCli(err)
	fmt.Fprint(w, string(template))
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10_000)
	if err != nil {
		ew.WrapBoth(w, err)
		return
	}

	files := r.MultipartForm.File
	if len(files["upfile"]) == 0 {
		log.Println("[!] No files to save.")
		ew.WrapHttp(w)
		return
	}
	for _, file := range files["upfile"] {
		cf, err := file.Open()
		if err != nil {
			ew.WrapBoth(w, err)
			return
		}
		defer cf.Close()
		fileContent, _ := ioutil.ReadAll(cf)
		_ = ioutil.WriteFile(fdir+file.Filename, fileContent, 0644)
	}
	template, err := ioutil.ReadFile("templates/success.html")
	ew.WrapCli(err)
	fmt.Fprint(w, string(template))
}

// UploadServer serves
func UploadServer(addr string, saveDir string, tls bool) {
	fdir = saveDir

	muxer := http.NewServeMux()
	handleMain := MiddlewareCliLogger(handleMain)
	handleFiles := MiddlewareCliLogger(handleFiles)
	muxer.HandleFunc("/", handleMain)
	muxer.HandleFunc("/upload", handleFiles)
	fmt.Printf("[!] Starting server %s\n", curDateTime())

	createUploadDir(saveDir)
	if tls {
		fmt.Printf("[*] %s https://%s\n", wfo, addr)
		log.Fatal(http.ListenAndServeTLS(addr, "tlscert/cert.pem", "tlscert/key.pem", muxer))
	} else {
		fmt.Printf("[*] %s http://%s\n", wfo, addr)
		log.Fatal(http.ListenAndServe(addr, muxer))
	}
}
