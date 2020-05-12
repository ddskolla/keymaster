package creds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type STSIssuer struct {
	STS stsiface.STSAPI
	RoleArn string
}

func NewSTSIssuer(STS stsiface.STSAPI, roleArn string) *STSIssuer {
	var issuer STSIssuer
	issuer.STS = STS
	issuer.RoleArn = roleArn
	return &issuer
}

func (i *STSIssuer) IssueFor(u *api.AuthInfo) ([]api.Cred, error) {
	var assumeRoleOutput *sts.AssumeRoleOutput
	var err error

	roleSessionName :=  u.Username + "-" + strconv.Itoa(int(time.Now().UnixNano() % 1e6))
	assumeRoleInput := sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(int64(u.ValidFor)),
		RoleArn:         &i.RoleArn,
		RoleSessionName: &roleSessionName,
	}
	assumeRoleOutput, err = i.STS.AssumeRole(&assumeRoleInput)
	if err != nil {
		return nil, errors.Wrapf(err, "error assuming role '%s'", i.RoleArn)
	}

	profileName := u.Environment + "-" + u.Role
	sExpiry := (*assumeRoleOutput.Credentials.Expiration).Unix()
	return []api.Cred{
		{
			Name:  profileName,
			Type:  "iam",
			Expiry: sExpiry,
			Value: &api.IAMCred{
				ProfileName:     profileName,
				RoleArn:         i.RoleArn,
				RoleSessionName: roleSessionName,
				AccessKeyId:     *assumeRoleOutput.Credentials.AccessKeyId,
				SecretAccessKey: *assumeRoleOutput.Credentials.SecretAccessKey,
				SessionToken:    *assumeRoleOutput.Credentials.SessionToken,
			},
		},
	}, nil
}
