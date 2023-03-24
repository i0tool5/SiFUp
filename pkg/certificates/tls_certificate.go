package certificates

var _ Certificator = (*TLSCertificate)(nil)

type TLSCertificate struct {
	*Certificate
}
