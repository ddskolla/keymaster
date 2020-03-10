package kube

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"github.com/jonboulle/clockwork"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

const (
	RsaKeyBits     = 2048
	CertFileSuffix = ".cert"
	KeyFileSuffix  = ".key"
)

type Issuer struct {
	CAKeypair     *tls.Certificate
	CACert        *x509.Certificate
	CACertEncoded string
	Clock         clockwork.Clock
}

type UserKeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

type EncodedUserKeyPair struct {
	PublicKeyPEM  []byte
	PrivateKeyPEM []byte
}

func NewIssuer(certPem, keyPem []byte) (*Issuer, error) {
	issuer := Issuer{}
	caKeypair, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}
	caCert, err := x509.ParseCertificate(caKeypair.Certificate[0])
	if err != nil {
		return nil, err
	}
	issuer.CAKeypair = &caKeypair
	issuer.CACert = caCert
	issuer.Clock = clockwork.NewRealClock()
	return &issuer, nil
}

func RandomSerial() (*big.Int, error) {
	// CA serial values are supposed to be *guaranteed* unique. But 63 bits
	// of randomness should be good enough given reasonable birthday bounds.
	// There are not really any genuine operational consequences even if a
	// collision does occur.
	randomSerial, err := rand.Int(rand.Reader, big.NewInt((1<<63)-1))
	if err != nil {
		return nil, err
	}
	randomSerial.Add(randomSerial, big.NewInt(1)) // [0,maxUint63) -> (1,maxUint63]
	return randomSerial, nil
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyId generates SubjectKeyId used in Certificate
// Id is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
func GenerateSubjectKeyId(pub crypto.PublicKey) ([]byte, error) {
	var pubBytes []byte
	var err error
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubBytes, err = asn1.Marshal(rsaPublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only RSA public key is supported")
	}

	hash := sha1.Sum(pubBytes)
	return hash[:], nil
}

// Generate a signed certificate for the specified CN and OrganizationalUnits. These map to the
// username and roles/groups in kubernetes.
func (issuer *Issuer) GenerateUserKeyPair(cn string, orgs []string, validForSeconds int) (*UserKeyPair, error) {
	// Generate an RSA keypair for the user
	log.Printf("generating rsa keypair for: %s (%v)\n", cn, orgs)
	priv, err := rsa.GenerateKey(rand.Reader, RsaKeyBits)
	if err != nil {
		return nil, err
	}
	pub := &priv.PublicKey

	// The subjectKeyId extension is not really "critical" (we could set it to anything really),
	// but it's "nice".
	subjectKeyId, err := GenerateSubjectKeyId(pub)
	if err != nil {
		return nil, err
	}

	// User certificate details to sign
	now := issuer.Clock.Now()
	validFor := time.Duration(validForSeconds) * time.Second
	serial, err := RandomSerial()
	if err != nil {
		return nil, err
	}
	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: orgs,
		},
		NotBefore:    now.Add(time.Second * -60), // move nbf back 1min to avoid time skews
		NotAfter:     now.Add(validFor),
		SubjectKeyId: subjectKeyId,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Sign the certificate
	log.Println("signing user certificate")
	signedCert, err := x509.CreateCertificate(rand.Reader, cert, issuer.CACert, pub, issuer.CAKeypair.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &UserKeyPair{
		PrivateKey: x509.MarshalPKCS1PrivateKey(priv),
		PublicKey:  signedCert,
	}, nil
}

func (kp *UserKeyPair) Encode() *EncodedUserKeyPair {
	return &EncodedUserKeyPair{
		PublicKeyPEM:  pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: kp.PublicKey}),
		PrivateKeyPEM: pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: kp.PrivateKey}),
	}
}
