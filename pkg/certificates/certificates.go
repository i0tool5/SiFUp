package certificates

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"time"
)

type Certificator interface {
	PublicKey() crypto.PublicKey
	X509() *x509.Certificate
	SetBytes([]byte)
	Encode(io.ReadWriter, io.ReadWriter) error
	Sign(Certificator) error
	SetSigned()
	IsSigned() bool
	IsCA() bool
}

var _ Certificator = (*Certificate)(nil)

type Certificate struct {
	cert       *x509.Certificate
	privateKey *rsa.PrivateKey
	encoded    bool
	signed     bool
	bytes      []byte
}

// IsCA displays is this certificate is CA
func (c *Certificate) IsCA() bool {
	return c.cert.IsCA
}

// Encode certificate with PEM encoding.
// Puts encoded certificate data into certPemOut
// and encoded private key into privKeyPemOut
func (c *Certificate) Encode(
	certPemOut, privKeyPemOut io.ReadWriter) (err error) {

	if err = pem.Encode(certPemOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.bytes,
	}); err != nil {
		return
	}

	if err = pem.Encode(privKeyPemOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(c.privateKey),
	}); err != nil {
		return
	}

	c.encoded = true
	return
}

// PublicKey of the certificate
func (c *Certificate) PublicKey() crypto.PublicKey {
	return &c.privateKey.PublicKey
}

// X509 certificate instance
func (c *Certificate) X509() *x509.Certificate {
	return c.cert
}

func (c *Certificate) SetBytes(b []byte) {
	c.bytes = b
}

// Sign certificate. If cert is nil, then certificate is self signed
func (c *Certificate) Sign(cert Certificator) (
	err error) {

	var signable Certificator = c

	if cert != nil {
		signable = cert
	}

	var certBytes []byte
	certBytes, err = x509.CreateCertificate(
		rand.Reader,
		signable.X509(),
		c.X509(),
		signable.PublicKey(),
		c.privateKey,
	)
	if err != nil {
		return
	}

	signable.SetBytes(certBytes)
	signable.SetSigned()

	return
}

func (c *Certificate) SetSigned() {
	c.signed = true
}
func (c *Certificate) IsSigned() bool {
	return c.signed
}

// NewX509Certificate generates x509 cetificate template
func NewX509Certificate(
	sn *big.Int,
	subj pkix.Name,
	keyUsage x509.KeyUsage,
	extKeyUsage []x509.ExtKeyUsage) (cert *x509.Certificate) {

	now := time.Now()
	cert = &x509.Certificate{
		SerialNumber: sn,
		Subject:      subj,
		NotBefore:    now,
		// TODO: make configurable
		NotAfter:    now.AddDate(1, 0, 0),
		ExtKeyUsage: extKeyUsage,
		KeyUsage:    keyUsage,
	}

	return
}

// NewPrivateKey generates new rsa private key.
// If bitsCount is less than 4096, then it will be set to 4096
func NewPrivateKey(bitsCount int) (privateKey *rsa.PrivateKey, err error) {
	if bitsCount < 4096 {
		bitsCount = 4096
	}
	privateKey, err = rsa.GenerateKey(rand.Reader, bitsCount)
	return
}
