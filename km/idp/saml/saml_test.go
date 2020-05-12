package saml

import (
	"encoding/base64"
	"github.com/dexidp/dex/connector"
	"github.com/dexidp/dex/connector/saml"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const TestCA = `-----BEGIN CERTIFICATE-----
MIIDGTCCAgGgAwIBAgIJAKLbLcQajEf8MA0GCSqGSIb3DQEBCwUAMCMxDDAKBgNV
BAoMA0RFWDETMBEGA1UEAwwKY29yZW9zLmNvbTAeFw0xNzA0MDQwNzAwNTNaFw0z
NzAzMzAwNzAwNTNaMCMxDDAKBgNVBAoMA0RFWDETMBEGA1UEAwwKY29yZW9zLmNv
bTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKH3dKWbRqCIZD2m3aHI
4lfBT+u/4DECde74Ggq9WugdTucVQzDZUTaI7wzn17JM9hdPmXvaSRG9BaB1H3uO
ZCs/fmdhBERRhPvuEVfAZaFfQfR7vn7WvUzT7zwMLLB8+EHzL3fOSGM2QnCOMeUD
AB27Pb0fuBW43NXaTD9rwfFCHvo1UP+TBJIPnV65HMeMGIrtGLt7MZTPuPm3LnYA
faXLf2vWSzL5nAgnJvUgceZXmyuciBfXpt8c1jIsj4y3tBoRTRqaxuaW1Eo7WMKF
a7s6KvTBKErPKuzAoIcVB4ir6jm1ticAgB72SScKtPJJdEPemTXRNNzkiw7VbpY9
QacCAwEAAaNQME4wHQYDVR0OBBYEFNHyGYyY2+eZ1l7ZLPZsnc3GOtj/MB8GA1Ud
IwQYMBaAFNHyGYyY2+eZ1l7ZLPZsnc3GOtj/MAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADggEBAHVXB5QmZfki9QpKzoiBNfpQ/mo6XWhExLGBTJXEWJT3P7JP
oR4Z0+85bp0fUK338s+WjyqTn0U55Jtp0B65Qxy6ythkZat/6NPp/S7gto2De6pS
hSGygokQioVQnoYQeK0MXl2QbtrWwNiM4HC+9yohbUfjwv8yI7opwn/rjB6X/4De
oX2YzwTBJgoIXF7zMKYFF0DrKQjbTQr/a7kfNjq4930o7VhFph9Qpdv0EWM3svTd
esSffLKbWcabtyMtCr5QyEwZiozd567oWFWZYeHQyEtd+w6tAFmz9ZslipdQEa/j
1xUtrScuXt19sUfOgjUvA+VUNeMLDdpHUKHNW/Q=
-----END CERTIFICATE-----
`

func TestConnectorIntegration(t *testing.T) {
	// Just an integration test for the dex connector.
	// Not actually testing any of our own code here, if
	// this fails then the connector changed.
	c := saml.Config{
		CAData:       []byte(TestCA),
		UsernameAttr: "Name",
		EmailAttr:    "email",
		GroupsAttr:   "groups",
		RedirectURI:  "http://127.0.0.1:5556/dex/callback",
		SSOURL:       "http://foo.bar/",
		//InsecureSkipSignatureValidation: true,
	}
	conn, err := c.Open("saml", logrus.New())
	assert.NoError(t, err)
	assert.NotEmpty(t, conn)
	samlConn := conn.(connector.SAMLConnector)

	inResponseTo := "6zmm5mguyebwvajyf2sdwwcw6m"

	// You can't mock time from the outside like:
	// conn.now = func() time.Time { return now }

	// And monkey patching like so fails:
	// testNow := time.Date(2017, time.April, 04, 04, 35, 00, 4, time.UTC)
	// patch := monkey.Patch(time.Now, func() time.Time { return testNow })
	// defer patch.Unpatch()
	// With: "verify signature: response does not contain a valid signature element:
	//    Missing signature referencing the top-level element"

	// So good-resp.xml has been modified to not expire until 2037.

	resp, err := ioutil.ReadFile("testdata/good-resp.xml")
	if err != nil {
		t.Fatal(err)
	}
	samlResp := base64.StdEncoding.EncodeToString(resp)
	scopes := connector.Scopes{
		OfflineAccess: false,
		Groups:        true,
	}
	ident, err := samlConn.HandlePOST(scopes, samlResp, inResponseTo)
	assert.NoError(t, err)
	assert.NotEmpty(t, ident)
	assert.Equal(t, "eric.chiang+okta@coreos.com", ident.UserID)
	assert.Equal(t, "Eric", ident.Username)
	assert.Equal(t, "eric.chiang+okta@coreos.com", ident.Email)
	assert.Equal(t, true, ident.EmailVerified)
	assert.Equal(t, []string{"Everyone", "Admins"}, ident.Groups)
}

// TODO: ensure at least 1 failing test case here also
// TODO: ideally, verify all these cases: https://www.samltool.com/generic_sso_res.php
