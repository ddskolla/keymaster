package api

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestLoadSampleConfigs(t *testing.T) {
	expected := ApiConfig{
		Name:        "fooproject_nonprod",
		Idp:         []IdpConfig{
			{
				Name: "nonprod",
				Type: "saml",
				Config: &IdpConfigSaml{
					Issuer:      "foo_saml_nonprod",
					Audience:    "fooproject_nonprod",
					Certificate: "pem_block",
				},
			},
		},
		Roles:       []RoleConfig{
			{
				Name: "cloudengineer",
				Credentials: []string{"ssh-all", "kube", "aws-admin"},
				Workflow: "cloudengineer",
				ValidForSeconds: 7200,
			},
			{
				Name: "developer",
				Credentials: []string{"ssh-jumpbox", "kube", "aws-ro"},
				Workflow: "developer",
				ValidForSeconds: 7200,
			},
			{
				Name: "deployment",
				Credentials: []string{"kube", "aws-admin"},
				Workflow: "deploy_with_identify_and_approval",
				CredentialDelivery: RoleCredentialDeliveryConfig{
					KmsWrapWith: "arn:aws:kms:ap-southeast-2:062921715532:key/95a6a059-8281-4280-8500-caf8cc217367",
				},
			},
		},
		Credentials: []CredentialsConfig{
			{
				Name: "ssh-jumpbox",
				Type: "ssh_ca",
				Config: &CredentialsConfigSSH{
					CAKey:      "s3://my-bucket/sshca.key",
					Principals: []string{"$idpuser"},
				},
			},
			{
				Name: "ssh-all",
				Type: "ssh_ca",
				Config: &CredentialsConfigSSH{
					CAKey:      "s3://my-bucket/sshca.key",
					Principals: []string{"$idpuser", "core", "ec2-user"},
				},
			},
			{
				Name: "kube-user",
				Type: "kubernetes",
				Config: &CredentialsConfigKube{
					CAKey:      "s3://my-bucket/kubeca.key",
				},
			},
			{
				Name: "kube-admin",
				Type: "kubernetes",
				Config: &CredentialsConfigKube{
					CAKey:      "s3://my-bucket/kubeca.key",
				},
			},
			{
				Name: "aws-ro",
				Type: "iam_assume_role",
				Config: &CredentialsConfigIAMAssumeRole{
					TargetRole: "arn:aws:iam::062921715666:role/ReadOnly",
				},
			},
			{
				Name: "aws-admin",
				Type: "iam_assume_role",
				Config: &CredentialsConfigIAMAssumeRole{
					TargetRole: "Administrator",
				},
			},
		},
		Workflow:    WorkflowConfig{
			BaseUrl: "https://workflow.bsy.place/1/",
			Policies: []WorkflowPolicyConfig{
				{
					Name: "deploy_with_identify",
					RequesterCanApprove: false,
					IdentifyRoles: map[string]int{
						"adfs_role_deployer": 1,
					},
				},
				{
					Name: "deploy_with_approval",
					RequesterCanApprove: false,
					ApproverRoles: map[string]int{
						"adfs_role_approver": 1,
					},
				},
				{
					Name: "deploy_with_identify_and_approval",
					RequesterCanApprove: false,
					IdentifyRoles: map[string]int{
						"adfs_role_deployer": 1,
					},
					ApproverRoles: map[string]int{
						"adfs_role_approver": 1,
					},
				},
				{
					Name: "developer",
					IdentifyRoles: map[string]int{
						"adfs_role_developer": 1,
					},
				},
				{
					Name: "cloudengineer",
					IdentifyRoles: map[string]int{
						"adfs_role_cloudengineer": 1,
					},
				},
			},
		},
		AccessControl: AccessControlConfig{ WhiteListCidrs: []string{
			"192.168.0.0/24",
			"172.16.0.0/12",
			"10.0.0.0/8",
		} },
	}
	data, err := ioutil.ReadFile("./testdata/example_api_config.yaml")
	assert.NoError(t, err)
	var result ApiConfig
	err = yaml.Unmarshal([]byte(data), &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestRequest_UnmarshalJSON(t *testing.T) {
	testCases := map[string]Request{
		"config": {
			Type: "config",
			Payload: &ConfigRequest{},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		log.Println(string(b))
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 Request
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}

func TestIdpConfig_UnmarshalJSON(t *testing.T) {
	testCases := map[string]IdpConfig{
		"t1": {
			Type: "saml",
			Name: "my-idp",
			Config: &IdpConfigSaml{
				Issuer: "foo",
				Audience: "bar",
				Certificate: "pem-goes-here",
			},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 IdpConfig
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}

func TestCredentialsConfig_UnmarshalJSON(t *testing.T) {
	testCases := map[string]CredentialsConfig{
		"ssh1": {
			Name: "ssh-example",
			Type: "ssh_ca",
			Config: &CredentialsConfigSSH{
				CAKey:"my-ssh-ca-key",
			},
		},
		"kube1": {
			Name: "kube-example",
			Type: "kubernetes",
			Config: &CredentialsConfigKube{
			},
		},
		"iam_assumerole1": {
			Name: "iam-assumerole-example",
			Type: "iam_assume_role",
			Config: &CredentialsConfigIAMAssumeRole{
			},
		},
		"iam_user1": {
			Name: "iam-user-example",
			Type: "iam_user",
			Config: &CredentialsConfigIAMUser{
			},
		},
	}

	// Unmarshal c -> c2, check c == c2
	for _, c := range testCases {
		b, err := json.Marshal(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, b)

		var c2 CredentialsConfig
		err = json.Unmarshal(b, &c2)
		assert.NoError(t, err)

		assert.Equal(t, c, c2)
	}
}