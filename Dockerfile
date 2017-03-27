FROM alpine

MAINTAINER Alexander Schittler <hello@damongant.de>

RUN apk --no-cache add go musl-dev ca-certificates

COPY . /go/src/github.com/Moonlington/FloSelfbot

RUN GOPATH=/go go build --ldflags '-extldflags "-static"' -o /usr/local/bin/FloSelfbot github.com/Moonlington/FloSelfbot
RUN rm -rf /go && apk del go musl-dev && mkdir /operator && chown -R operator:nobody /operator && chmod 700 /operator

USER operator
WORKDIR /operator
VOLUME /operator

ENTRYPOINT ["/usr/local/bin/FloSelfbot"]
