package sts

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/bsycorp/keymaster/km/creds/u"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Issuer struct {
	STS stsiface.STSAPI
	RoleArn string
	ProfileName string
	ValidForSeconds int
}

func NewIssuer(sess *session.Session, roleArn string, profileName string, validForSeconds int) *Issuer {
	var issuer Issuer
	issuer.STS = sts.New(sess)
	issuer.RoleArn = roleArn
	issuer.ProfileName = profileName
	issuer.ValidForSeconds = validForSeconds
	return &issuer
}

func (i *Issuer) IssueFor(u *u.UserInfo) (map[string]string, error) {
	var assumeRoleOutput *sts.AssumeRoleOutput
	var err error

	stsDuration := int64(i.ValidForSeconds)
	roleSessionName := u.Username + "-" + strconv.Itoa(int(time.Now().UnixNano() % 1e6))
	assumeRoleInput := sts.AssumeRoleInput{
		DurationSeconds: &stsDuration,
		RoleArn:         &i.RoleArn,
		RoleSessionName: &roleSessionName,
	}
	assumeRoleOutput, err = i.STS.AssumeRole(&assumeRoleInput)
	if err != nil {
		return nil, errors.Wrap(err, "Error in AWS assumerole")
	}

	awsCredsFormat :=
		`[%s]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token = %s
# Keymaster issued, expires: %s
`
	awsCredentials := fmt.Sprintf(
		awsCredsFormat,
		i.ProfileName,
		*assumeRoleOutput.Credentials.AccessKeyId,
		*assumeRoleOutput.Credentials.SecretAccessKey,
		*assumeRoleOutput.Credentials.SessionToken,
		*assumeRoleOutput.Credentials.Expiration,
	)
	// TODO: return typed credentials and let the client sort out
	return map[string]string{
		"~/.aws/credentials": awsCredentials,
	}, nil
}
