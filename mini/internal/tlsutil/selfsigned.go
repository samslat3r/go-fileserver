package tlsutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"net"
	"os"
	"time"
)

//EnsureSelfSigned checks if the self-signed certificate and key exist at the specified paths.
func EnsureSelfSigned(certPath, keyPath, host string) error {
	if exists(certPath) && exists(keyPath) {
		return nil // If they are already there
	}
	return generate(certPath, keyPath, host)
}

func generate(certPath, keyPath, host string) error {
	priv, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	tmpl := &x509.Certificate{
		SerialNumber:		  bigInt(),
		Subject:			  pkix.Name{Organization: []string{"Mini Self-Signed"}},
		NotBefore:			  time.Now(),
		NotAfter:			  time.Now().Add(30 * 24 * time.Hour),
		KeyUsage:			  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:		  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range []string{host, "localhost"} {
		if ip := net.ParseIP(h); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}

	der, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)

	certOut, _ := os.Create(certPath)
	defer certOut.Close()
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der})

	keyOut, _ := os.Create(keyPath)
	defer keyOut.Close()
	b, _ := x509.MarshalPKCS8PrivateKey(priv) 
	_ = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: b}) 
	
	return nil
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }
func bigInt() *big.Int {
	n, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	return n
}
func Config() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP512, tls.CurveP384, tls.CurveP256},
		PreferCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384, tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_EECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
	}
}
