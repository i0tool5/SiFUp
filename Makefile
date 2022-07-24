create-certificates:
	openssl req -x509 -nodes -newkey rsa:4096 -keyout tlscert/key.pem -out tlscert/cert.pem -sha256 -days 365