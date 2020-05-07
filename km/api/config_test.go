package api

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLoadSampleConfigs(t *testing.T) {
	expected := Config{
		Name:    "fooproject_nonprod",
		Version: "1.0",
		Idp: []IdpConfig{
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
		Roles: []RoleConfig{
			{
				Name:            "cloudengineer",
				Credentials:     []string{"ssh-all", "kube", "aws-admin"},
				Workflow:        "cloudengineer",
				ValidForSeconds: 7200,
			},
			{
				Name:            "developer",
				Credentials:     []string{"ssh-jumpbox", "kube", "aws-ro"},
				Workflow:        "developer",
				ValidForSeconds: 7200,
			},
			{
				Name:            "deployment",
				Credentials:     []string{"kube", "aws-admin"},
				Workflow:        "deploy_with_approval",
				ValidForSeconds: 3600,
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
					CAKey: "s3://my-bucket/kubeca.key",
				},
			},
			{
				Name: "kube-admin",
				Type: "kubernetes",
				Config: &CredentialsConfigKube{
					CAKey: "s3://my-bucket/kubeca.key",
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
		Workflow: WorkflowConfig{
			BaseUrl: "https://workflow.int.btr.place/",
			Policies: []WorkflowPolicyConfig{
				{
					Name:                "deploy_with_identify",
					IdpName:             "nonprod",
					RequesterCanApprove: false,
					IdentifyRoles: map[string]int{
						"adfs_role_deployer": 1,
					},
				},
				{
					Name:                "deploy_with_approval",
					IdpName:             "nonprod",
					RequesterCanApprove: false,
					ApproverRoles: map[string]int{
						"adfs_role_approver": 1,
					},
				},
				{
					Name:                "deploy_with_identify_and_approval",
					IdpName:             "nonprod",
					RequesterCanApprove: false,
					IdentifyRoles: map[string]int{
						"adfs_role_deployer": 1,
					},
					ApproverRoles: map[string]int{
						"adfs_role_approver": 1,
					},
				},
				{
					Name:    "developer",
					IdpName: "nonprod",
					IdentifyRoles: map[string]int{
						"adfs_role_developer": 1,
					},
				},
				{
					Name:    "cloudengineer",
					IdpName: "nonprod",
					IdentifyRoles: map[string]int{
						"adfs_role_cloudengineer": 1,
					},
				},
			},
		},
		AccessControl: AccessControlConfig{
			IPOracle: IPOracleConfig{
				WhiteListCidrs: []string{"192.168.0.0/24", "172.16.0.0/12", "10.0.0.0/8"},
			},
		},
	}
	data, err := ioutil.ReadFile("./testdata/example_api_config.yaml")
	assert.NoError(t, err)
	var result Config
	err = yaml.Unmarshal([]byte(data), &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestConfig_FindCredentialByName(t *testing.T) {
	config := Config{
		Credentials: []CredentialsConfig{
			{
				Name: "ssh-all",
				Type: "ssh_ca",
				Config: &CredentialsConfigSSH{
					CAKey:      "s3://my-bucket/sshca.key",
					Principals: []string{"$idpuser", "core", "ec2-user"},
				},
			},
		},
	}
	credConfig := config.FindCredentialByName("ssh-all")
	assert.NotNil(t, credConfig)

	assert.Nil(t, config.FindCredentialByName("does-not-exist"))
}

func TestConfig_FindRoleByName(t *testing.T) {
	config := Config{
		Roles: []RoleConfig{
			{
				Name:            "developer",
				Credentials:     []string{"ssh-jumpbox", "kube", "aws-ro"},
				Workflow:        "developer",
				ValidForSeconds: 7200,
			},
		},
	}
	credConfig := config.FindRoleByName("developer")
	assert.NotNil(t, credConfig)

	assert.Nil(t, config.FindCredentialByName("does-not-exist"))
}

func TestIdpConfig_UnmarshalJSON(t *testing.T) {
	testCases := map[string]IdpConfig{
		"t1": {
			Type: "saml",
			Name: "my-idp",
			Config: &IdpConfigSaml{
				Issuer:      "foo",
				Audience:    "bar",
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
				CAKey: "my-ssh-ca-key",
			},
		},
		"kube1": {
			Name:   "kube-example",
			Type:   "kubernetes",
			Config: &CredentialsConfigKube{},
		},
		"iam_assumerole1": {
			Name:   "iam-assumerole-example",
			Type:   "iam_assume_role",
			Config: &CredentialsConfigIAMAssumeRole{},
		},
		"iam_user1": {
			Name:   "iam-user-example",
			Type:   "iam_user",
			Config: &CredentialsConfigIAMUser{},
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
