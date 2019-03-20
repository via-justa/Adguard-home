package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func watchFile(certificateFile string, privateKeyFile string, configFile string, changed chan time.Time, config *ConfigFile) {

	cert, err := os.Stat(certificateFile)
	checkERR(err)

	initState := cert.ModTime()

	for {
		cert, _ := os.Stat(certificateFile)
		newState := cert.ModTime()

		if newState != initState {
			updateConfig(certificateFile, privateKeyFile, configFile, config)
			initState = cert.ModTime()
			changed <- newState
		}

		time.Sleep(1 * time.Minute)
	}

}

func updateConfig(certificateFile string, privateKeyFile string, adGurardConfigFile string, config *ConfigFile) {
	var timer int
	timeout := config.Letsencrypt.Timeout
	for {
		if _, err := os.Stat(certificateFile); os.IsNotExist(err) {
			untilTimeout := timeout - timer
			if timer == 1 {
				log.Println("Looking for certificate file, it can take few seconds for the certificate to generate for te first time.")
			}
			log.Printf("Please wait... Time left until timeout: %v\n", untilTimeout)
			time.Sleep(5 * time.Second)
			timer = timer + 5
			if timer == timeout {
				log.Println("Certificate lookup timed out, pleasse make sure the FQDN selected is right and owned by you")
				os.Exit(1)
			}
		} else {
			log.Printf("Found new/updated certificate in path: %v\n", certificateFile)
			break
		}
	}

	newCertificate, err := ioutil.ReadFile(certificateFile)
	checkERR(err)
	cert := string(newCertificate)

	newKey, err := ioutil.ReadFile(privateKeyFile)
	checkERR(err)
	key := string(newKey)

	config.TLS.CertificateChain = cert
	config.TLS.PrivateKey = key

	toYaml, err := yaml.Marshal(&config)
	checkERR(err)

	err = ioutil.WriteFile(adGurardConfigFile, toYaml, 755)
	checkERR(err)

	log.Println("AdGuard config updated")

}

func startAdGuard(adguardconfigFile string, pid chan *os.Process, errExit chan error) {
	adguardExecutable := "/opt/adguardhome/AdGuardHome"

	log.Println("Starting AdGuard-Home...")
	cmd := exec.Command(adguardExecutable, "-h", "0.0.0.0", "-c", adguardconfigFile, "-w", "/opt/adguardhome/work")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	checkERR(err)

	go func() {
		pid <- cmd.Process
	}()

	errExit <- cmd.Wait()

}

func filePath(configurationFile string, config *ConfigFile) (string, string) {

	var certPath string
	var keyPath string

	if config.Letsencrypt.Enabled == false {
		if os.Getenv("CERT_FILE") == "" || os.Getenv("KEY_FILE") == "" {
			log.Printf("Both \"CERT_FILE\" and \"KEY_FILE\" environment variables need to be set when letsencrypt_enabled set to \"false\"")
		} else {
			certPath = os.Getenv("CERT_FILE")
			keyPath = os.Getenv("KEY_FILE")
		}
	} else {
		if config.Letsencrypt.Production == false {
			certPath = strings.Join([]string{"/root/.caddy/acme/acme-staging-v02.api.letsencrypt.org/sites/", config.TLS.ServerName, "/", config.TLS.ServerName, ".", "crt"}, "")
			keyPath = strings.Join([]string{"/root/.caddy/acme/acme-staging-v02.api.letsencrypt.org/sites/", config.TLS.ServerName, "/", config.TLS.ServerName, ".", "key"}, "")
		} else {
			certPath = strings.Join([]string{"/root/.caddy/acme/acme-v02.api.letsencrypt.org/sites/", config.TLS.ServerName, "/", config.TLS.ServerName, ".", "crt"}, "")
			keyPath = strings.Join([]string{"/root/.caddy/acme/acme-v02.api.letsencrypt.org/sites/", config.TLS.ServerName, "/", config.TLS.ServerName, ".", "key"}, "")
		}

	}

	return certPath, keyPath

}
