package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/bsycorp/keymaster/km/idp/saml"
	"github.com/bsycorp/keymaster/km/workflow"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

var roleFlag = flag.String("role", "", "target role")
var debugFlag = flag.Int("debug", 0, "enable debugging")
var debugLevel = 0

func main() {
	flag.Parse()

	// create km directory
	kmDirectory := fmt.Sprintf("%s/.km", UserHomeDir())
	if err := os.MkdirAll(kmDirectory, 0700); err != nil {
		log.Println("Failed to create ~/.km directory: ", err)
	}

	if *roleFlag == "" {
		log.Fatalln("Required argument role missing (need -role)")
	}
	debugLevel = *debugFlag
	// Draft workflow

	// First, get the config
	// target := "arn:aws:lambda:ap-southeast-2:062921715532:function:km2"
	//target := "arn:aws:apigateway:ap-southeast-2:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-southeast-2:218296299700:function:km2/invocations"
	target := "arn:aws:lambda:ap-southeast-2:218296299700:function:km-tools-bls-01"
	kmApi := api.NewClient(target)
	kmApi.Debug = debugLevel

	discoveryReq := new(api.DiscoveryRequest)
	_, err := kmApi.Discovery(discoveryReq)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error calling kmApi.Discovery"))
	}

	configReq := new(api.ConfigRequest)
	configResp, err := kmApi.GetConfig(configReq)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error calling kmApi.GetConfig"))
	}

	// Now start workflow to get nonce
	kmWorkflowStartResponse, err := kmApi.WorkflowStart(&api.WorkflowStartRequest{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "error calling kmApi.WorkflowStart"))
	}
	log.Println("Started workflow with km api")

	log.Println("Target role for authentication:", *roleFlag)
	targetRole := configResp.Config.FindRoleByName(*roleFlag)
	if targetRole == nil {
		log.Fatalf("Target role #{*roleFlag} not found in config")
	}
	workflowPolicyName := targetRole.Workflow
	configWorkflowPolicy := configResp.Config.Workflow.FindPolicyByName(workflowPolicyName)
	if configWorkflowPolicy == nil {
		log.Fatalf("workflow policy %s not found in config", workflowPolicyName)
	}
	workflowPolicy := workflow.Policy{
		Name:                configWorkflowPolicy.Name,
		IdpName:             configWorkflowPolicy.IdpName,
		RequesterCanApprove: configWorkflowPolicy.RequesterCanApprove,
		IdentifyRoles:       configWorkflowPolicy.IdentifyRoles,
		ApproverRoles:       configWorkflowPolicy.ApproverRoles,
	}

	workflowBaseUrl := configResp.Config.Workflow.BaseUrl
	log.Println("Using workflow engine:", workflowBaseUrl)
	workflowApi, err := workflow.NewClient(workflowBaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	workflowApi.Debug = debugLevel

	// And start a workflow session
	startResult, err := workflowApi.Create(context.Background(), &workflow.CreateRequest{
		IdpNonce: kmWorkflowStartResponse.IdpNonce,
		Requester: workflow.Requester{
			Name:     "Blair Strang",
			Username: "strangb",
			Email:    "blair.strang@auspost.com.au",
		},
		Source: workflow.Source{
			Description: "Deploy a new version 3.2 with amazing features",
			DetailsURI:  "https://gitlab.com/platform/keymaster",
		},
		Target: workflow.Target{
			EnvironmentName:         configResp.Config.Name,
			EnvironmentDiscoveryURI: "TBD",
		},
		Policy: workflowPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Now fix up the workflow URL
	fixedWorkflowUrl := "https://workflow.int.btr.place/workflow/" + startResult.WorkflowId
	log.Printf("------------------------------------------------------------------")
	log.Printf("******************************************************************")
	log.Printf("APPROVAL URL: %s", fixedWorkflowUrl)
	log.Printf("******************************************************************")
	log.Printf("------------------------------------------------------------------")

	// Poll for assertions
	var getAssertionsResult *workflow.GetAssertionsResponse
	for {
		getAssertionsResult, err = workflowApi.GetAssertions(context.Background(), &workflow.GetAssertionsRequest{
			WorkflowId:    startResult.WorkflowId,
			WorkflowNonce: startResult.WorkflowNonce,
		})
		if err != nil {
			log.Println(errors.Wrap(err, "error calling workflowApi.GetAssertions"))
		}
		log.Printf("workflow state: %s", getAssertionsResult.Status)
		if getAssertionsResult.Status == "CREATED" {
			time.Sleep(5 * time.Second)
		} else if getAssertionsResult.Status == "COMPLETED" {
			break
		} else if getAssertionsResult.Status == "REJECTED" {
			log.Fatal("Your change request was REJECTED by a workflow approver. Exiting.")
		} else {
			log.Fatal("unexpected assertions result status:", getAssertionsResult.Status)
		}
	}
	log.Printf("got: %d assertions from workflow", len(getAssertionsResult.Assertions))

	// Locally verify the assertions
	sp := &saml.AssertionProcessor{
		CAData: []byte(`
-----BEGIN CERTIFICATE-----
MIICnTCCAYUCBgFfA+Q72DANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdjbHVz
dGVyMB4XDTE3MTAxMDAxMjUxMFoXDTI3MTAxMDAxMjY1MFowEjEQMA4GA1UEAwwH
Y2x1c3RlcjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMdVYTd1h7fa
u6/uCgboFyFdoRSWFEHP0Iq9GUWA69g2x+QDqZikSv/JqPwtJBAm+dxdXfOd0RKT
4ypK09PUNy542kJ+Qwgzwif0ZIEKTYOVS8VvwzZv6BjzwDzSBS/LmdcK8WgRGwgh
62QgjIYQdGd+wrYN0tOQb6EzINWMs1bq9bFjeFegDG94p/MZ1YWRVXF6h/euq/ym
gJQc7yvUn5cy6l47tT1ARrCzpUF8Ss4eVhNlLDaz5WSzZ4P1Q+bPe4Iax//zMr/J
62aqmcf/YuVKIINLa5ML+QFW2B+mR0xky8jwWJiwU5gJzDzLoiNQZ3TJxcfvQaT1
PuC8ksM9bd0CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAvnrKy75SHGEAIPORf2QC
NxqWi6Qc/Pl1gHSGHd9nPcIn7u2dRmoq45XWAr55yVZqT/FWshOII504YuFJCQF5
fyOGKy00jVmaOEIPqyLRA0wf4AsZk607Y2CVZIl1JGwuYx5rHgZ2kf1M4Qxvnhl/
OUkMrW+VosBgIrqiKWd53Y5TnHaX/q+hYoa/GmRXq0JTJOX+5C11YX9G4rsI7o3c
MP19yto+e+d5myXu3POAvx4VG07LlWWk3cow2xuiw4zJbZVmK6KO2rMk66WJpfQu
EmyLmLPjKTmhoskvaHhvSoW6h06Uth3Lf6UHHsAkdzeU+mw0g2Zb2dPlDqz4IV4t
cg==
-----END CERTIFICATE-----`),
		Audience:     "keymaster-saml",
		UsernameAttr: "name",
		EmailAttr:    "name",
		GroupsAttr:   "groups",
		RedirectURI:  "https://workflow.int.btr.place/1/saml/approve",

	}
	err = sp.Init()
	if err != nil {
		log.Fatalln("saml init error:", err)
	}
	res, err := sp.Process(kmWorkflowStartResponse.IdpNonce, getAssertionsResult.Assertions)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)

	creds, err := kmApi.WorkflowAuth(&api.WorkflowAuthRequest{
		Username:     "gitlab", // TODO
		Role:         "deployment",
		IdpNonce:     kmWorkflowStartResponse.IdpNonce,
		IssuingNonce: kmWorkflowStartResponse.IssuingNonce,
		Assertions:   getAssertionsResult.Assertions,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "error calling kmApi.WorkflowAuth"))
	}
	log.Printf("GOT CREDENTIALS...")

	var iamCred *api.Cred
	for _, cred := range creds.Credentials {
		if cred.Type == "iam" {
			iamCred = &cred
			break
		}
	}
	if iamCred == nil {
		log.Fatal("Got creds but no IAM cred?")
	}
	iamCredValue, ok := iamCred.Value.(*api.IAMCred)
	if !ok {
		log.Fatal("oops IAM cred is wrong type?")
	}

	awsCredsFmt := `[%s]
aws_access_key_id = %s
aws_secret_access_key = %s
aws_session_token = %s
# Keymaster issued, expires: %s
`
	exp := time.Unix(iamCred.Expiry, 0)
	localAwsCreds := fmt.Sprintf(
		awsCredsFmt,
		iamCredValue.ProfileName,
		iamCredValue.AccessKeyId,
		iamCredValue.SecretAccessKey,
		iamCredValue.SessionToken,
		exp,
	)

	awsCredentialsPath := UserHomeDir() + "/.aws/credentials"
	existingCreds, err := ioutil.ReadFile(awsCredentialsPath)
	if err != nil {
		fmt.Printf("Failed to update local credentials: %v", err)
	} else {
		log.Printf("Found existing credentials file, appending..")
		awsCredentialsIni, err := ini.Load(existingCreds, []byte(localAwsCreds))
		if err != nil {
			fmt.Printf("Failed to read existing local credentials: %v", err)
		} else {
			err = awsCredentialsIni.SaveTo(awsCredentialsPath)
			if err != nil {
				fmt.Printf("Failed to update local credentials: %v", err)
			}
		}
	}

}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
