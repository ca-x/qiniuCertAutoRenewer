package certkit

import (
	"testing"
)

func TestTLSEndpoint_GetCertificates(t *testing.T) {
	tlsEndpoint := NewTLSEndpoint("czyt.tech", "443")
	tls, err := tlsEndpoint.GetCertificates()
	if err != nil {
		t.Fatal(err)
	}
	for _, cert := range tls {
		if !cert.IsCA {
			t.Logf("certInfo:%+v,%+v\n", cert.DNSNames, cert.NotAfter)
		}
	}
}
