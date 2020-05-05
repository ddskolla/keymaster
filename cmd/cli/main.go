package main

import (
	"context"
	"fmt"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/bsycorp/keymaster/km/workflow"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {

	// create km directory
	kmDirectory := fmt.Sprintf("%s/.km", UserHomeDir())
	if err := os.MkdirAll(kmDirectory, 0700); err != nil {
		log.Println("Failed to create ~/.km directory: ", err)
	}

	// Draft workflow

	// First, get the config
	// target := "arn:aws:lambda:ap-southeast-2:062921715532:function:km2"
	//target := "arn:aws:apigateway:ap-southeast-2:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-southeast-2:218296299700:function:km2/invocations"
	target := "arn:aws:lambda:ap-southeast-2:218296299700:function:km-tools-bls-01"
	kmApi := api.NewClient(target)
	configReq := new(api.ConfigRequest)
	config, err := kmApi.GetConfig(configReq)
	if err != nil {
		log.Fatal(err)
	}

	// Now start workflow to get nonce
	kmWorkflowStartResponse, err := kmApi.WorkflowStart(&api.WorkflowStartRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Started workflow with km api")

	// Get the right policy
	workflowPolicyName := "deploy_with_approval"
	configWorkflowPolicy := config.Config.Workflow.FindPolicyByName(workflowPolicyName)
	if configWorkflowPolicy == nil {
		log.Fatalf("workflow policy %s not found in config", workflowPolicyName)
	}
	log.Println(configWorkflowPolicy)
	workflowPolicy := workflow.Policy{ // Blech
		Name:                configWorkflowPolicy.Name,
		IdpName:             configWorkflowPolicy.IdpName,
		RequesterCanApprove: configWorkflowPolicy.RequesterCanApprove,
		IdentifyRoles:       configWorkflowPolicy.IdentifyRoles,
		ApproverRoles:       configWorkflowPolicy.ApproverRoles,
	}

	// Then create a workflow client
	// TODO: the workflow URL should come from config
	workflowApi, err := workflow.NewClient("https://workflow.int.btr.place/1")
	if err != nil {
		log.Fatal(err)
	}

	// And start a workflow session
	startResult, err := workflowApi.Start(context.Background(), &workflow.StartRequest{
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
			EnvironmentName:         "apxyz-env-02",
			EnvironmentDiscoveryURI: "TBD",
		},
		Policy: workflowPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(startResult)

	// Now fix up the workflow URL
	fixedWorkflowUrl := "https://workflow.int.btr.place/workflow/" + startResult.WorkflowId
	log.Printf("------------------------------------------------------------------")
	log.Printf("******************************************************************")
	log.Printf("APPROVAL URL: %s", fixedWorkflowUrl)
	log.Printf("******************************************************************")
	log.Printf("------------------------------------------------------------------")

	// Poll for assertions
	for {
		getAssertionsResult, err := workflowApi.GetAssertions(context.Background(), &workflow.GetAssertionsRequest{
			WorkflowId: startResult.WorkflowId,
			Nonce:      startResult.Nonce,
		})
		if err != nil {
			log.Println(err)
		}
		if getAssertionsResult.Status == "CREATED" {
			log.Println("WATING FOR APPROVAL")
			time.Sleep(5 * time.Second)
		} else if getAssertionsResult.Status == "COMPLETED" {
			break
		} else {
			log.Fatal("unexpected assertions result status:", getAssertionsResult.Status)
		}
	}
	log.Println("GOT ASSERTIONS")

	creds, err := kmApi.WorkflowAuth(&api.WorkflowAuthRequest{
		Username: "gitlab",
		Role:     "deployment",
		Nonce:    kmWorkflowStartResponse.Nonce,
		// SAML assertions go here
	})
	if err != nil {
		log.Fatal(err)
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
