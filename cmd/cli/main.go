package main

import (
	"fmt"
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

	// Invoke lambda
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
