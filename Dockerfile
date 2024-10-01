FROM golang:1.23 as build

WORKDIR /go/src/github.com/skpr/zap
COPY . /go/src/github.com/skpr/zap

ENV CGO_ENABLED=0
RUN go build -o bin/zap-wrapper -ldflags='-extldflags "-static"' github.com/skpr/zap/cmd/zap-wrapper

FROM ghcr.io/zaproxy/zaproxy:stable
ENV ZAP_WRAPPER_DIRECTORY=/tmp
COPY --from=build /go/src/github.com/skpr/zap/bin/zap-wrapper /usr/local/bin/zap-wrapper
ENTRYPOINT ["/usr/local/bin/zap-wrapper"]