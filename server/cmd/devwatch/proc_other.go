//go:build !linux

package main

import "syscall"

// childAttr puts the child in its own process group, so a shutdown reaches
// whatever the server itself started and the terminal's Ctrl-C arrives here
// rather than racing us to the child.
func childAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}
