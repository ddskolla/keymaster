package ssh

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/jonboulle/clockwork"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
)

const (
	KeyBits            = 2048
	MaxValidForSeconds = 7 * 24 * 3600
)

type UserInfo struct {
	Identity        string
	Principals      []string
	ValidForSeconds int
}

type Credentials struct {
	Certificate []byte
	PrivateKey  []byte
}

type Issuer struct {
	Random io.Reader
	Clock  clockwork.Clock
}

func (issuer *Issuer) GenerateKeyPair(user *UserInfo) (ssh.PublicKey, *rsa.PrivateKey, error) {
	if user.ValidForSeconds < 0 || user.ValidForSeconds > MaxValidForSeconds {
		return nil, nil, errors.New("Invalid issuance period")
	}
	if issuer.Random == nil {
		return nil, nil, errors.New("No random source? what happened?")
	}

	// Generate user private key
	log.Println("Generating RSA private key")
	privateKey, err := rsa.GenerateKey(issuer.Random, KeyBits)
	if err != nil {
		return nil, nil, err
	}

	// Generate user public key
	log.Println("Getting public key")
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}

func (issuer *Issuer) CreateSignedCertificate(ca ssh.Signer, publicKey ssh.PublicKey, privateKey *rsa.PrivateKey, user *UserInfo, extensions map[string]string, options map[string]string) (*Credentials, error) {
	// Create a signed SSH certificate for the user
	// As per: https://www.ietf.org/mail-archive/web/secsh/current/msg00327.html
	now := uint64(issuer.Clock.Now().Unix())
	userCert := &ssh.Certificate{
		CertType:        ssh.UserCert,
		KeyId:           user.Identity,
		ValidPrincipals: user.Principals,
		ValidAfter:      now,
		ValidBefore:     now + uint64(user.ValidForSeconds),
		Key:             publicKey,
		Permissions: ssh.Permissions{
			CriticalOptions: options,
			Extensions: extensions,
		},
	}
	// Sign the user's key with the CA key
	log.Println("Signing SSH certificate")
	if err := userCert.SignCert(issuer.Random, ca); err != nil {
		return nil, err
	}

	// Marshal the user's certificate
	log.Println("Marshalling SSH certificate")
	userCertBytes := ssh.MarshalAuthorizedKey(userCert)

	// Marshal the user's private key
	log.Println("Marshalling SSH private key")
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return nil, err
	}

	sshCreds := Credentials{
		Certificate: userCertBytes,
		PrivateKey:  private.Bytes(),
	}

	log.Println("Successfully issued SSH credentials")
	return &sshCreds, nil
}
