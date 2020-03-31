package creds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const validFor = 3600

type mockSTSClient struct {
	stsiface.STSAPI
	t *testing.T
}

func (m *mockSTSClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	assert.Equal(m.t, *input.DurationSeconds, int64(validFor))
	return &sts.AssumeRoleOutput{
		AssumedRoleUser: &sts.AssumedRoleUser{
			Arn:           aws.String("arn"),
			AssumedRoleId: aws.String("assumed-role-id"),
		},
		Credentials: &sts.Credentials{
			AccessKeyId:     aws.String("access-key-id"),
			Expiration:      aws.Time(time.Now().Add(5 * time.Minute)),
			SecretAccessKey: aws.String("secret-access-key"),
			SessionToken:    aws.String("session-token"),
		},
		PackedPolicySize: aws.Int64(32),
	}, nil
}

type mockSTSClientFail struct {
	stsiface.STSAPI
}

func (m *mockSTSClientFail) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return nil, errors.New("it didn't work")
}

func TestSTSIssuer(t *testing.T) {
	i := NewSTSIssuer(&mockSTSClient{t:t}, "my-super-role-arn")
	u := api.AuthInfo{
		Environment: "foo.io",
		Role:        "super-admin",
		Username:    "fred",
		ValidFor:    validFor,
	}
	result, err := i.IssueFor(&u)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "iam", result[0].Type)

	i.STS = &mockSTSClientFail{}
	result, err = i.IssueFor(&u)
	assert.Empty(t, result)
	assert.Error(t, err)
}
