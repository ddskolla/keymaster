package main

import (
	"fmt"
	"github.com/bsycorp/keymaster/km/api"
	"github.com/bsycorp/keymaster/km/workflow"
	"log"
	"os"
	"runtime"
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
	log.Println(config)

	// Then create a workflow session
	// TODO: look this up from config
	workflow := workflow.Client("https://")

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
