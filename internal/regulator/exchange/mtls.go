// Package exchange holds the secure data-exchange endpoints used
// between Joblantern and verified regulators. Mutual-TLS is optional
// but strongly recommended for any large bilateral feed.
package exchange

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
)

// Config bundles the trust material for an exchange listener.
type Config struct {
	ServerCert      tls.Certificate
	ClientCAPool    *x509.CertPool
	RequireClientCA bool
}

// TLSConfig assembles a *tls.Config that enforces the supplied
// client-CA policy. When RequireClientCA is true the listener will
// refuse any client whose certificate does not chain to the pool.
func (c Config) TLSConfig() (*tls.Config, error) {
	if c.ServerCert.Certificate == nil {
		return nil, errors.New("server certificate not loaded")
	}
	cfg := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{c.ServerCert},
	}
	if c.RequireClientCA {
		if c.ClientCAPool == nil {
			return nil, errors.New("client CA pool required")
		}
		cfg.ClientCAs = c.ClientCAPool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return cfg, nil
}
