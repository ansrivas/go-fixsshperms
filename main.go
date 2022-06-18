package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type StringSlice []string

const (
	PermSSHDir     = os.FileMode(0700)
	PermAuthKeys   = os.FileMode(0644)
	PermConfig     = os.FileMode(0644)
	PermKnownHosts = os.FileMode(0644)
	PermPubKeys    = os.FileMode(0644)
	PermPrivKeys   = os.FileMode(0600)
)

func updatePermissions(files []string, permission os.FileMode) {
	for _, fpath := range files {
		err := os.Chmod(fpath, permission)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Updated permissions for [%s] with [%v]", fpath, permission)
		time.Sleep(time.Millisecond * 500)
	}

}

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func getKeysPath(basePath string) ([]string, []string) {

	pubKeys := []string{}
	privKeys := []string{}
	fileList := []string{}
	err := filepath.Walk(basePath, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileList {
		if strings.Contains(file, "id_rsa") || strings.Contains(file, "id_ed25519") {

			if strings.HasSuffix(file, ".pub") {
				pubKeys = append(pubKeys, file)
			}
			if !strings.HasSuffix(file, ".pub") {
				privKeys = append(privKeys, file)
			}

		}
	}
	return pubKeys, privKeys
}

func main() {
	basePath := getUserHome()
	sshDir := filepath.Join(basePath, ".ssh")
	authKeys := filepath.Join(sshDir, "authorized_keys")
	config := filepath.Join(sshDir, "config")
	knownHosts := filepath.Join(sshDir, "known_hosts")

	log.Printf("Using %s as the ssh directory", sshDir)

	pubKeys, privKeys := getKeysPath(basePath)
	updatePermissions(pubKeys, PermPubKeys)
	updatePermissions(privKeys, PermPrivKeys)
	updatePermissions(StringSlice{sshDir}, PermSSHDir)
	updatePermissions(StringSlice{authKeys}, PermAuthKeys)
	updatePermissions(StringSlice{config}, PermConfig)
	updatePermissions(StringSlice{knownHosts}, PermKnownHosts)
}
