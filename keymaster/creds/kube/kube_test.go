package kube

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

const (
	OpenSSLCommand  = "/usr/bin/openssl"
	CaTestCertFile = "testdata/kube_ca.crt"
	CaTestCertKey = "testdata/kube_ca.key"
)


func MustLoadFile(s string) ([]byte) {
	res, err := ioutil.ReadFile(s)
	if err != nil {
		panic(err)
	}
	return res
}

// This isn't in stdlib coz it's "racy"; generally you want to create and open temp
// files in one step to avoid problems. But it's hard to avoid doing this kind of
// thing when you're shelling out to do work. You have to pass around names and not
// file handles.
//
// This is test code only.

func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)  // TODO: Not sure what else do if we can't read ...
	}
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func TestNewIssuer(t *testing.T) {
	// Good construction
	issuer, err := NewIssuer(MustLoadFile(CaTestCertFile), MustLoadFile(CaTestCertKey))
	assert.Nil(t, err)
	assert.NotNil(t, issuer)
}

func TestCreateUserKey(t *testing.T) {
	// Create an issuer
	issuer, err := NewIssuer(MustLoadFile(CaTestCertFile), MustLoadFile(CaTestCertKey))
	assert.Nil(t, err)
	assert.NotNil(t, issuer)

	// Use a mocked clock in the issuer to test NotBefore and NotAfter conditions
	mockedTime := time.Date(1984, time.April, 1, 16, 20, 1, 0, time.UTC)
	issuer.Clock = clockwork.NewFakeClockAt(mockedTime)

	// Sign a user key
	kp, err := issuer.GenerateUserKeyPair("admstrangb", []string{"super-admins", "something"}, 3600)
	assert.Nil(t, err)
	assert.NotNil(t, kp)

	userCertFile, userKeyFile, err := kp.Dump(TempFileName("ukp", ""))
	assert.Nil(t, err)

	defer os.Remove(userCertFile)
	defer os.Remove(userKeyFile)

	// Check the certificate with OpenSSL
	certDump, err := exec.Command(OpenSSLCommand, "x509", "-in", userCertFile, "-text").Output()
	assert.Nil(t, err)
	certDumpStr := string(certDump)

	// Check some key cert properties

	// We want to check for:
	//
	// Subject: O=super-admins, O=something, CN=admstrangb
	// Subject: O = super-admins + O = something, CN = admstrangb
	//
	// It turns out different versions of the openssl command-line tool output this line
	// in the cert differently. We at least check that the CN is what we expect and that at
	// a minimum the groups we asked for are there.
	// assert.Contains(t, certDumpStr, "Subject: O=super-admins, O=something, CN=admstrangb")
	//
	assert.Regexp(t, regexp.MustCompile("Subject:.*O *= *super-admins"), certDumpStr)
	assert.Regexp(t, regexp.MustCompile("Subject:.*O *= *something"), certDumpStr)
	assert.Regexp(t, regexp.MustCompile("Subject:.*CN *= *admstrangb"), certDumpStr)

	assert.Contains(t, certDumpStr, "Public Key Algorithm: rsaEncryption")
	assert.Contains(t, certDumpStr, "Signature Algorithm: sha256WithRSAEncryption")
	assert.Contains(t, certDumpStr, "Key: (2048 bit)")
	assert.Contains(t, certDumpStr, "Not Before: Apr  1 16:19:01 1984 GMT")
	assert.Contains(t, certDumpStr, "Not After : Apr  1 17:20:01 1984 GMT")

	// Since we mocked time way into the past, we expect openssl to say the cert has expired
	validateResult, _ := exec.Command(OpenSSLCommand, "verify", "-CAfile", CaTestCertFile, userCertFile).Output()
	// err may be "exit status 2" depending on openssl version
	// assert.Nil(t, err)

	// Another case where different versions of OpenSSL will output different things...
	assert.Regexp(t, regexp.MustCompile("(verification failed|certificate has expired)"), string(validateResult))
}

func TestValidUserKey(t *testing.T) {
	issuer, err := NewIssuer(MustLoadFile(CaTestCertFile), MustLoadFile(CaTestCertKey))
	assert.NotNil(t, issuer)
	assert.Nil(t, err)

	// Sign a user key
	kp, err := issuer.GenerateUserKeyPair("admstrangb", []string{"super-admins", "something"}, 3600)
	assert.Nil(t, err)
	assert.NotNil(t, kp)

	// Dump the keys
	userCertFile, userKeyFile, err := kp.Dump(TempFileName("ukp", ""))
	assert.Nil(t, err)

	defer os.Remove(userCertFile)
	defer os.Remove(userKeyFile)

	// Since we issued with current time, this should be:
	// a) Not expired
	// b) Valid with respect to the CA_TEST_CERT
	validateResult, err := exec.Command(OpenSSLCommand, "verify", "-CAfile", CaTestCertFile, userCertFile).Output()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s: OK\n", userCertFile), string(validateResult))
}

// Dump keypair to files, no cleanup is performed on error (aka partial key writes
// are not removed). Currently, this is only used for testing.
func (kp *UserKeyPair) Dump(path string) (certPath string, keyPath string, err error) {
	certPath = path + CertFileSuffix
	keyPath = path + KeyFileSuffix

	// Dump the cert
	certfile, err := os.Create(certPath)
	if err != nil {
		return
	}
	err = pem.Encode(certfile, &pem.Block{Type: "CERTIFICATE", Bytes: ukp.PublicKey})
	if err != nil {
		return
	}
	err = certfile.Close()
	if err != nil {
		return
	}

	// Dump the private key
	keyfile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	err = pem.Encode(keyfile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: ukp.PrivateKey})
	if err != nil {
		return
	}
	err = keyfile.Close()
	if err != nil {
		return
	}
	return
}
