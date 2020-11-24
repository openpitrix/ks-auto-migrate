FROM golang:1.13-alpine as builder

RUN apk add --no-cache git curl openssl

WORKDIR /go/src/github.com/openpitrix/ks-auto-migrate/
COPY . .

RUN mkdir -p /release_bin
RUN CGO_ENABLED=0 GOBIN=/release_bin go install -ldflags '-w -s' -tags netgo github.com/openpitrix/ks-auto-migrate/cmd/...


FROM alpine:3.7
RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /release_bin/* /usr/local/bin/
