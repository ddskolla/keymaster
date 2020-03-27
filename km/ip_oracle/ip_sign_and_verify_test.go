package ip_oracle

import (
	"testing"
	"time"
	"flag"

	"github.com/stretchr/testify/assert"
	jwt "github.com/dgrijalva/jwt-go"
)

var KmsTestKeyId = flag.String("KmsTestKeyId", "720ed5bb-33b2-4493-a4d6-68751f08d60a", "The ID of the test KMS key")

func TestSuccessfulCreateIPJWT(t *testing.T) {
	sm := NewSigningMethodKMS(*KmsTestKeyId)
	_, err := MakeIPJWT("192.168.5.6", sm)
	assert.NoError(t, err)
}

func TestUnsuccessfulCreateIPJWTWithInvalidKMSKey(t *testing.T) {
	sm := NewSigningMethodKMS("1234-abcd-5678-efgh")
	_, err := MakeIPJWT("1.1.1.1", sm)
	assert.Error(t, err, "NotFoundException: Invalid keyId 1234-abcd-5678-efgh")
}

func TestSuccessfulIPJWTVerifyViaKMS(t *testing.T) {
	t.Parallel()

	sm := NewSigningMethodKMS(*KmsTestKeyId)
	signedString, _ := MakeIPJWT("192.168.5.6", sm)
	sourceIp, err := VerifyIPJWT(signedString, sm)

	assert.NoError(t, err, "Verify IP via KMS should be successful!")
	assert.Equal(t, sourceIp, "192.168.5.6", "Verified IP should be same as the client source IP")
}

func TestUnsuccessfulIPJWTVerifyWithInvalidKMSKey(t *testing.T) {
	t.Parallel()

	sm := NewSigningMethodKMS(*KmsTestKeyId)
	signedString, _ := MakeIPJWT("192.168.5.6", sm)

	sm = NewSigningMethodKMS("1234-abcd-5678-efgh")
	_, err := VerifyIPJWT(signedString, sm)
	assert.Error(t, err, "NotFoundException: Invalid keyId 1234-abcd-5678-efgh")
}

func TestUnsuccessfulIPJWTVerifyAfterExpiry(t *testing.T) {
	t.Parallel()

	sm := NewSigningMethodKMS(*KmsTestKeyId)
	issuedAt := time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC)
	expiresAt := issuedAt.Add(time.Minute * 10)

	claims := IPOracleClaims{
		"1.2.3.4",
		jwt.StandardClaims{
		  IssuedAt: issuedAt.UnixNano() / int64(time.Second),
		  ExpiresAt: expiresAt.UnixNano() / int64(time.Second),
		},
	}
	token := jwt.NewWithClaims(sm, claims)
	signedString, _ := token.SignedString(sm)

	_, err := VerifyIPJWT(signedString, sm)
	assert.Contains(t, err.Error(), "token is expired by")
}