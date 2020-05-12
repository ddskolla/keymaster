package main

import (
	"syscall"
)

func FixWindowsPerms(name string) error {
	return nil
}

func PlatformExec(cmd string, args []string, envv []string) error {
	return syscall.Exec(cmd, args, envv)
}
