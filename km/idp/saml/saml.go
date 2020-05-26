package saml

import (
	"github.com/bsycorp/keymaster/km/idp/connector/saml"
	"github.com/dexidp/dex/connector"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AssertionProcessor struct {
	CAData   []byte
	Audience string
	// TODO: verify audience?
	// TODO: NameID stuff?
	// TODO: Encryption?
	UsernameAttr            string
	EmailAttr               string
	GroupsAttr              string
	RedirectURI             string
	DisableNameIDValidation bool

	conn     connector.Connector
	samlConn connector.SAMLConnector
}

type UserInfo struct {
	Username string
	Groups   []string
}

func (sp *AssertionProcessor) Init() error {
	c := saml.Config{
		EntityIssuer:            sp.Audience,
		CAData:                  []byte(sp.CAData),
		UsernameAttr:            sp.UsernameAttr,
		EmailAttr:               sp.EmailAttr,
		GroupsAttr:              sp.GroupsAttr,
		RedirectURI:             sp.RedirectURI,
		SSOURL:                  "UNUSED",
		DisableNameIDValidation: sp.DisableNameIDValidation,
	}
	conn, err := c.Open("saml", logrus.New())
	if err != nil {
		return errors.Wrap(err, "AssertionProcessor: open error")
	}
	var ok bool
	sp.samlConn, ok = conn.(connector.SAMLConnector)
	if !ok {
		return errors.New("AssertionProcessor: not a saml connector!")
	}
	return nil
}

func (sp *AssertionProcessor) Process(inResponseTo string, assertions []string) ([]UserInfo, error) {
	scopes := connector.Scopes{
		OfflineAccess: false,
		Groups:        true,
	}
	result := make([]UserInfo, 0, len(assertions))
	for _, samlResponse := range assertions {
		ident, err := sp.samlConn.HandlePOST(scopes, samlResponse, inResponseTo)
		if err != nil {
			return nil, errors.Wrap(err, "AssertionProcessor: invalid saml response")
		}
		userInfo := UserInfo{
			Username: ident.Username,
			Groups:   ident.Groups,
		}
		result = append(result, userInfo)
	}

	return result, nil
}
