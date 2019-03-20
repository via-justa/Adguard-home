# Adguard Home with DNS over TLS support

[AdGuard home](https://www.github.com/AdguardTeam/AdGuardHome) is a free and open source, powerful network-wide Ads & trackers blocking DNS server. 

AdGuard Home supports [DNS over TLS (DoT)](https://en.wikipedia.org/wiki/DNS_over_TLS). This project aim to enable the use of AdGuard Home with DoT with [Let's Encrypt](https://letsencrypt.org/) certificates.

The project is leveraging [Caddy](https://caddyserver.com/) built-in Let's Encrypt client with [DigitalOcean as dynamic DNS provider](https://www.digitalocean.com/community/projects/digital-ocean-dynamic-dns)

To use different provider check the [Change provider](#Change-provider) section

## Prerequisites

* Access to a domain (or sub domain) hosted on DigitalOcen
* Docker and optionally docker-compose installed:
  ```
  $ apt-get update && apt-get upgrade
  $ curl -sSL https://get.docker.com | sh
  $ apt-get install docker-compose
  $ usermod -aG docker pi
  ```

## Usage

Clone the repo in a folder and CD to the repository
```
$ git clone git@github.com:via-justa/Adguard-Home.git && CD Adguard-Home
```
### Prepering the configuration
Create an API token as explained [here](https://www.digitalocean.com/docs/api/create-personal-access-token/).

Edit `./config/AdGuardHome.yaml` and set the following settings values:

field | type | default value
-----|------|------
letsencrypt.enabled | bool (true) | Use Caddy to generate the certificate. If `false` see [Using external certificate](#Using-external-certificate) |
letsencrypt.production| bool (true) | If `true` Let's Encrypt production server is used. If `false` Let's Encrypt Staging server is used |
letsencrypt.timeout | int (30) | Time to wait for certificate to generate on first use |
letsencrypt.email | string (none) | Email to register the certificate with |
letsencrypt.provider | string (digitalocean) | Provider to use with caddy, if changing see [Change provider](#Change-provider) section |
letsencrypt.provider_settings.DO_AUTH_TOKEN | string (none) | DigitalOcean API token with write permission to DNS zone |
auth_name | string (admin) | User for AgGuard-Home frontend |
auth_pass | string (ChangeM3) | Password for AgGuard-Home frontend |
tls.server_name | string (none) | FQDN of the server |

The rest of the settings can stay as they are. A full list of the settings can be found [here](https://github.com/AdguardTeam/AdGuardHome/wiki/Configuration#configuration-file)

### Start docker deployment

```
$ docker build . --tag adguardhome
$ docker run -d --rm --name adguardhome adguardhome
```

Or using docker-compose
```
$ cd adGuardHome
$ docker-compose up -d --build
```

:warning: Please note: Any configuration changes done via the Adguard-Home frontend will not be persistant.

## Change provider

Make sure you provider is [supported](https://caddyserver.com/docs/automatic-https#dns-challenge) and what are the required environment variables. 
Edit `dockerfile` and replace `github.com/caddyserver/dnsproviders/digitalocean` with your provider url like `github.com/caddyserver/dnsproviders/azure`
```
...
RUN go get \
  github.com/caddyserver/dnsproviders/digitalocean \
...
RUN sed -i '/\/\/ This is where other plugins get plugged in (imported)/a \        _ "github.com/caddyserver/dnsproviders/digitalocean" \n ...
...
```
Edit `./config/AdGuardHome.yaml`:
- Replace the value of `letsencrypt.provider` with the name of your provider.
- Remove `letsencrypt.provider_settings.DO_AUTH_TOKEN` and add the list of environment varables required by your provider,

Azure provider example:
```
letsencrypt:
  ...
  provider_settings: 
    AZURE_CLIENT_ID: xxxxx-xxxxx-xxxxx-xxxxx
    AZURE_CLIENT_SECRET: xxxxxxxxxxxxxxxxx
    AZURE_SUBSCRIPTION_ID: xxxxx-xxxxx-xxxxx-xxxxx
    AZURE_TENANT_ID: xxxxx-xxxxx-xxxxx-xxxxx

```
When the configuration is phrased each key value pair will be converted to an environment variable. 

## Using external certificate

If you like to manage the certificates externally, set `letsencrypt.enabled` to `false` and pass environment variables `CERT_FILE` with the path to a PEM crt file with a full certificate chain and `KEY_FILE` with the path to the corresponding key.
You'll need to mount the certificates when starting the container using `-v` (or `volumes:` if using docker-compose)