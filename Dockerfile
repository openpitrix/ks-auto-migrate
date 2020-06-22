FROM openpitrix/openpitrix-builder as builder

WORKDIR /go/src/openpitrix.io/Jobs/
COPY . .

RUN mkdir -p /release_bin
RUN GOPROXY=https://goproxy.io CGO_ENABLED=0 GOBIN=/release_bin go install -ldflags '-w -s' -tags netgo openpitrix.io/Jobs/cmd/...


FROM alpine:3.7
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /release_bin/* /usr/local/bin/