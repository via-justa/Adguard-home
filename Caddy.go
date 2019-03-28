package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func updateCaddyFile(errExit chan error, config *ConfigFile) {
	caddyFilePath := "/caddy/caddyfile"

	for key, val := range config.Letsencrypt.ProviderSettings {
		os.Setenv(key, val)
	}

	f, err := ioutil.ReadFile(caddyFilePath)
	checkERR(err)

	caddyfile := string(f)

	caddyfile = strings.Replace(caddyfile, "FQDN", config.TLS.ServerName, -1)
	caddyfile = strings.Replace(caddyfile, "EMAIL", config.Letsencrypt.Email, -1)
	caddyfile = strings.Replace(caddyfile, "PROVIDER", config.Letsencrypt.Provider, -1)

	b := []byte(caddyfile)
	ioutil.WriteFile(caddyFilePath, b, 644)
	log.Println("Caddyfile updated")

	go startCaddy(caddyFilePath, config, errExit)
}

func startCaddy(caddyFilePath string, config *ConfigFile, errExit chan error) {
	caddyExecutable := "/usr/local/bin/caddy"
	var ca string

	if config.Letsencrypt.Production == true {
		ca = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else {
		ca = "https://acme-v02.api.letsencrypt.org/directory"
	}

	//Start Caddy
	log.Println("Starting Caddy")
	cmd := exec.Command(caddyExecutable,
		"--conf", caddyFilePath,
		"--log", "stdout",
		"-http-port", "8081",
		"-https-port", "4434",
		"-agree",
		"-ca", ca)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	checkERR(err)

	errExit <- cmd.Wait()
}
