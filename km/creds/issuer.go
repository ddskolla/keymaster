package creds

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/pkg/errors"
	"log"
)

type issuer interface {
	IssueFor(u *api.AuthInfo) ([]api.Cred, error)
}

type Issuer struct {
	issuers []issuer
}

func NewFromConfig(role *api.RoleConfig, config *api.Config) (*Issuer, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	var issuer Issuer
	for _, credName := range role.Credentials {
		credConfig := config.FindCredentialByName(credName)
		switch c := credConfig.Config.(type) {
		case *api.CredentialsConfigIAMAssumeRole:
			i := NewSTSIssuer(sts.New(sess), c.TargetRole)
			issuer.issuers = append(issuer.issuers, i)
		default:
			log.Printf("TODO: unimplemented cred config type for: %s", credName)
		}
	}
	return &issuer, nil
}

func (i *Issuer) IssueFor(u *api.AuthInfo) ([]api.Cred, error) {
	allCreds := make([]api.Cred, 0)
	for _, iss := range i.issuers {
		creds, err := iss.IssueFor(u)
		if err != nil {
			errx := errors.Wrap(err, "error during credential issuance")
			log.Println(errx)
			return nil, errx
		}
		allCreds = append(allCreds, creds...)
	}
	return allCreds, nil
}
