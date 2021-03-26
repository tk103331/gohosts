package main

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
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

func validateName(name string, group *HostsGroup) error {
	reg := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	match := reg.Match([]byte(name))
	if !match {
		return errors.New("Name can only contains [a-zA-Z0-9_]")
	}
	for _, item := range group.Items {
		if item.GetName() == name {
			return errors.New("Name already exist")
		}
	}
	return nil
}
