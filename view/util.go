package view

import (
	"io/ioutil"
	"os"
)

func loadSystem() string {
	data, _ := ioutil.ReadFile("/etc/hosts")
	return string(data)
}

func loadBackup() string {
	file := "/etc/hosts.bak"
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
