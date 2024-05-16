package certkit

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
)

var (
	// Allow certificate that signed by unknown authority.
	// Controller only concerns expiration of certificate.
	defaultTLSConfig = tls.Config{InsecureSkipVerify: true}

	// DefaultPortNumber exposes default port number to testing
	DefaultPortNumber = 443
)

// TLSEndpoint expressses https endpoint that using TLS.
type TLSEndpoint struct {
	Hostname string
	Port     int
}

// NewTLSEndpoint creates new TLSEndpoint instance.
// If port number is empty, set DefaultPortNumber instead.
func NewTLSEndpoint(host string, port int) *TLSEndpoint {
	if port == 0 {
		port = DefaultPortNumber
	}

	return &TLSEndpoint{
		Hostname: host,
		Port:     port,
	}
}

// GetCertificates tries to get certificates from endpoint using tls.Dial
func (e *TLSEndpoint) GetCertificates() ([]*x509.Certificate, error) {

	// We cannot connect to Hostnames with wildcards, so replacing with cert-test.
	hostName := strings.Replace(e.Hostname, "*", "cert-test", -1)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", hostName, e.Port), &defaultTLSConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}
