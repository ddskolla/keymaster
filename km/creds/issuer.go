package creds

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/bsycorp/keymaster/km/creds/iam/sts"
	"github.com/bsycorp/keymaster/km/creds/u"
	"github.com/pkg/errors"
	"log"
)

type issuer interface {
	IssueFor(u *u.UserInfo) (map[string]string, error)
}

type Issuer struct {
	issuers []issuer
}

func NewFromConfig(roleName string, config *api.Config) (*Issuer, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	roleConfig := config.FindRoleByName(roleName)
	if roleConfig == nil {
		return nil, errors.New("role not found: " + roleName)
	}
	var issuer Issuer
	for _, credName := range roleConfig.Credentials {
		credConfig := config.FindCredentialByName(credName)
		switch c := credConfig.Config.(type) {
		case *api.CredentialsConfigIAMAssumeRole:
			profileName := config.Name + "-" + roleName
			i := sts.NewIssuer(sess, c.TargetRole, profileName, roleConfig.ValidForSeconds)
			issuer.issuers = append(issuer.issuers, i)
		default:
			log.Printf("TODO: unimplemented cred config type for: %s", credName)
		}
	}
	return &issuer, nil
}

func (i *Issuer) IssueFor(u *u.UserInfo) (map[string]string, error) {
	allCreds := make(map[string]string)
	for _, iss := range i.issuers {
		creds, err := iss.IssueFor(u)
		if err != nil {
			errx := errors.Wrap(err, "error during credential issuance")
			log.Println(errx)
			return nil, errx
		}
		for k, v := range creds {
			// We're ignoring overwrites for now...
			// To be fixed in typed creds work
			allCreds[k] = v
		}
	}
	return allCreds, nil
}
