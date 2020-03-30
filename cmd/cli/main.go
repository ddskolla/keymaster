package main

import (
	"context"
	"fmt"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/bsycorp/keymaster/km/workflow"
	"github.com/davecgh/go-spew/spew"
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
	target := "arn:aws:lambda:ap-southeast-2:062921715532:function:km2"
	km := api.NewClient(target)
	configReq := new(api.ConfigRequest)
	config, err := km.GetConfig(configReq)
	if err != nil {
		log.Fatal(err)
	}

	// Now start workflow to get nonce
	// kmWorkflowStartResponse, err := km.WorkflowStart(&api.WorkflowStartRequest{})
	_, err = km.WorkflowStart(&api.WorkflowStartRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// Get the right policy
	workflowPolicyName := "deploy_with_approval"
	configWorkflowPolicy := config.Config.Workflow.FindPolicyByName(workflowPolicyName)
	if configWorkflowPolicy == nil {
		log.Fatalf("workflow policy %s not found in config", workflowPolicyName)
	}
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
			Username: "admstrangb",
			Email:    "blair.strang@auspost.com.au",
		},
		Source: workflow.Source{
			Description: "Deploy a new version 3.2 with amazing features",
			DetailsURI: "https://gitlab.com/platform/keymaster",
		},
		Target: workflow.Target{
			EnvironmentName: "apxyz-env-02",
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
	log.Printf("CLICK THIS LINK: %s", fixedWorkflowUrl)

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
		} else {
			spew.Dump(getAssertionsResult)
			break
		}
	}
	log.Println("GOT ASSERTIONS")

	// Now post these back to the km api
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
