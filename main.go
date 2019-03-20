package main

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	configurationFile := "/opt/adguardhome/conf/AdGuardHome.yaml"

	pid := make(chan *os.Process)
	changed := make(chan time.Time)
	errExit := make(chan error)

	// collect configuration settings from file
	yamlFile, err := ioutil.ReadFile(configurationFile)
	checkERR(err)

	var config ConfigFile

	err = yaml.Unmarshal(yamlFile, &config)
	checkERR(err)

	certificateFile, privateKeyFile := filePath(configurationFile, &config)

	if config.Letsencrypt.Enabled == false {
		log.Printf("Certificates mangaed external to the progrem.\nSkipping certificate generation")
	} else {
		// Run Caddy to get certificates
		log.Println("Updating Caddy congiguration")
		updateCaddyFile(errExit, &config)
	}

	// Update the config at start
	log.Println("Updating AdGuard-Home configuration")
	updateConfig(certificateFile, privateKeyFile, configurationFile, &config)

	// Start ADGuars-Home
	go startAdGuard(configurationFile, pid, errExit)

	go watchFile(certificateFile, privateKeyFile, configurationFile, changed, &config)

	for {
		select {
		case <-errExit:
			os.Exit(19)
		case <-changed:
			p := <-pid
			p.Signal(syscall.SIGTERM)
			go startAdGuard(configurationFile, pid, errExit)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func checkERR(err error) {
	if err != nil {
		log.Panic(err)
	}
}
