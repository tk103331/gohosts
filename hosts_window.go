// +build windows

package main

import "os"

func hostsFile() string {
	windir := "C:\\Windows"
	if dir, ok := os.LookupEnv("windir"); ok {
		windir = dir
	}
	return windir + "\\System32\\drivers\\etc\\hosts"
}
