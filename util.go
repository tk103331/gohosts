package main

import (
	"io/ioutil"
	"os"
)

func loadSystem() string {
	file := hostsFile()
	data, _ := ioutil.ReadFile(file)
	return string(data)
}

func saveSystem(content string) error {
	file := hostsFile()
	return ioutil.WriteFile(file, []byte(content), os.ModePerm)
}

func loadBackup() string {
	file := hostsFile() + ".bak"
	stat, err := os.Stat(file)
	if err == os.ErrNotExist {
		println("not found")
		return ""
	}
	if stat != nil && stat.IsDir() {
		println("is dir")
		return ""
	}
	data, _ := ioutil.ReadFile(file)
	return string(data)
}

func saveBackup(content string) error {
	file := hostsFile() + ".bak"
	return ioutil.WriteFile(file, []byte(content), os.ModePerm)
}
