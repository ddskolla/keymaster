package main

import (
	acl "github.com/hectane/go-acl"
	"os"
	"os/exec"
	"syscall"
)

func FixWindowsPerms(name string) error {
	// return winacl.Chmod(name, perms)
	var mode os.FileMode = 0700
	return acl.Apply(
		name,
		true,
		false,
		acl.GrantName((uint32(mode)&0700)<<23, "CREATOR OWNER"),
	)
}

func PlatformExec(cmd string, args []string, envv []string) error {
	// Does not actually "exec" on Windows, just hides the CMD window
	c := exec.Command(cmd, args[1:]...)
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c.Run()
}
