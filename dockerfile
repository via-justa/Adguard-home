## build caddy
FROM golang:1.12.1 as caddy-builder

ARG caddy_tag="v0.11.3"

RUN \
  mkdir -p /go/src/github.com/mholt \
  && cd /go/src/github.com/mholt \
  && git clone https://github.com/mholt/caddy.git \
  && cd caddy \
  && git checkout tags/${caddy_tag}

RUN go get \
  github.com/caddyserver/dnsproviders/digitalocean \
  github.com/caddyserver/builds

RUN \
  sed -i '/\/\/ This is where other plugins get plugged in (imported)/a \        _ "github.com/caddyserver/dnsproviders/digitalocean"' /go/src/github.com/mholt/caddy/caddy/caddymain/run.go \
  && sed -i"" 's/var EnableTelemetry = true/var EnableTelemetry = false/g' /go/src/github.com/mholt/caddy/caddy/caddymain/run.go

WORKDIR /go/src/github.com/mholt/caddy/caddy

RUN go run build.go


## Build Adguard
FROM golang:1.12.1 AS adguard-builder

ARG Adguard_tag="v0.93"

RUN \
  apt-get update \
  && apt install --upgrade curl software-properties-common -y \
  && curl -sL https://deb.nodesource.com/setup_11.x | bash - \
  && apt install --upgrade git make nodejs -y

RUN \
  git clone https://github.com/AdguardTeam/AdGuardHome.git \
  && cd AdGuardHome \
  && git checkout tags/${Adguard_tag} \
  && make


## build manager
FROM golang:1.12.1 as manager-builder

ADD ./ /go/src/github.com/via-justa/adguard-home/

WORKDIR /go/src/github.com/via-justa/adguard-home

RUN \
  CGO_ENABLED=0 go get \
  && CGO_ENABLED=0 go build


## Put it all together 
FROM alpine:latest

COPY --from=caddy-builder /go/src/github.com/mholt/caddy/caddy/caddy /usr/local/bin/caddy
COPY --from=adguard-builder /go/AdGuardHome/AdGuardHome /opt/adguardhome/AdGuardHome
COPY --from=manager-builder /go/src/github.com/via-justa/adguard-home/adguard-home /usr/local/bin/manager

RUN \
  apk --no-cache --update add ca-certificates \
  && rm -rf /var/cache/apk/* && mkdir -p /opt/adguardhome/conf

ADD config/caddyfile /caddy/caddyfile
ADD config/AdGuardHome.yaml /opt/adguardhome/conf/AdGuardHome.yaml

VOLUME ["/opt/adguardhome/work", "/root/.caddy"]

EXPOSE 53/tcp 53/udp 67/tcp 67/udp 68/tcp 68/udp 80/tcp 443/tcp 853/tcp 853/udp 3000/tcp

ENTRYPOINT ["/usr/local/bin/manager"]