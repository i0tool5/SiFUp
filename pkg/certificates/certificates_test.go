package certificates

import (
	"bytes"
	"testing"
)

func TestNewCertificationAuthority(t *testing.T) {
	ca, err := NewCertificationAuthority()
	if err != nil {
		t.Fatal(err)
	}

	if !ca.IsSigned() {
		t.Fatal("CA should be signed")
	}

	if !ca.IsCA() {
		t.Fatal("CA should be CA")
	}
}

func TestCertificateEncode(t *testing.T) {
	ca, err := NewCertificationAuthority()
	if err != nil {
		t.Fatal(err)
	}

	certPem, privPem := new(bytes.Buffer), new(bytes.Buffer)
	err = ca.Encode(certPem, privPem)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("cert: %s", certPem.String())
}
