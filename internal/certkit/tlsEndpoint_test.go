package certkit

import (
	"reflect"
	"testing"
)

func TestTLSEndpoint_GetCertificates(t *testing.T) {
	tlsEndpoint := NewTLSEndpoint("fake.google.com", 443)
	endpointCerts, err := tlsEndpoint.GetCertificates()
	if err != nil {
		t.Fatal(err, reflect.TypeOf(err))
	}
	for _, cert := range endpointCerts {
		if !cert.IsCA {
			t.Logf("certInfo:%+v,%+v\n", cert.DNSNames, cert.NotAfter)
		}
	}
}
