# simpleuploader
This small program is designed to make it easier to upload files to the system, for example during penetration testing. 

**Do not use standard files in the tlscert directory for security reasons.**

`Make` can be used to generate certificates, or you can use openssl (or other utility) to create them. E.x.:
```
openssl req -x509 -newkey rsa:4096 -keyout tlscers/key.pem -out tlscers/cert.pem -sha256 -days 365
```
