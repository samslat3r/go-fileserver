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


