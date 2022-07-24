create-certificates:
	openssl req -x509 -nodes -newkey rsa:4096 \
		-keyout pkg/tlscert/files/key.pem \
		-out pkg/tlscert/files/cert.pem \
		-sha256 \
		-days 365

build:
	go build \
		-o cmd/simpleuploader/simpleuploader \
		cmd/simpleuploader/main.go
