package certificates

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
)

var _ Certificator = (*CertificationAuthority)(nil)

type CertificationAuthority struct {
	*Certificate
}

// NewCertificationAuthority
func NewCertificationAuthority() (ca *CertificationAuthority, err error) {
	var privKey *rsa.PrivateKey
	privKey, err = NewPrivateKey(4096)
	if err != nil {
		return
	}

	cert := &Certificate{
		cert:       newCertificationAuthority(),
		privateKey: privKey,
	}
	ca = &CertificationAuthority{
		Certificate: cert,
	}
	err = ca.Sign(nil) // self signed

	return
}

// newCertificationAuthority generates new certification authority
// certificate instance
func newCertificationAuthority() (ca *x509.Certificate) {
	ca = NewX509Certificate(
		big.NewInt(31337),
		pkix.Name{
			Country:       []string{""},
			Organization:  []string{""},
			Locality:      []string{""},
			Province:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		x509.KeyUsageCertSign|x509.KeyUsageDigitalSignature,
		[]x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
	)
	ca.IsCA = true

	return
}
