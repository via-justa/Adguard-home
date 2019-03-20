package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

func dynamicDNS(config *ConfigFile) {
	pat := config.DynamicDNS.APIToken

	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	ctx := context.TODO()

	// Get first 100 records
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	records, _, err := client.Domains.Records(ctx, extractDomain(config), opt)
	checkERR(err)

	for {
		// Check if record exist
		for i, record := range records {

			if record.Name == extractRecord(config) {

				if record.Data != getExternalIP() {

					editRequest := &godo.DomainRecordEditRequest{
						Type: record.Type,
						Name: record.Name,
						Data: getExternalIP(),
					}

					domainRecord, _, err := client.Domains.EditRecord(ctx, extractDomain(config), record.ID, editRequest)
					checkERR(err)

					fmt.Println(strings.Join([]string{"DNS record updated new IP is: ", domainRecord.Data}, ""))
				}

				// If record could not be found
			} else if i == len(records)-1 {

				createRequest := &godo.DomainRecordEditRequest{
					Type: "A",
					Name: extractRecord(config),
					Data: getExternalIP(),
				}

				domainRecord, _, err := client.Domains.CreateRecord(ctx, extractDomain(config), createRequest)
				checkERR(err)

				fmt.Println(strings.Join([]string{"DNS record created with IP: ", domainRecord.Data}, ""))
			}

		}

		// Run check every 1 hour
		time.Sleep(1 * time.Hour)

	}

}

func getExternalIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	checkERR(err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String()
}

func extractDomain(config *ConfigFile) string {

	return strings.SplitAfterN(config.TLS.ServerName, ".", 2)[1]
}

func extractRecord(config *ConfigFile) string {

	return strings.Split(config.TLS.ServerName, ".")[0]
}

// TokenSource struct to store API access token
type TokenSource struct {
	AccessToken string
}

// Token convert TokenSource to oauth2.Token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
