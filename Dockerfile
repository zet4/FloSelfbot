FROM alpine

MAINTAINER Alexander Schittler <hello@damongant.de>

COPY . /go/src/github.com/Moonlington/FloSelfbot

VOLUME /data
ENV GOPATH=/go

RUN apk --no-cache add go musl-dev ca-certificates git && go get -d github.com/Moonlington/FloSelfbot && go build --ldflags '-extldflags "-static"' -o /usr/local/bin/FloSelfbot github.com/Moonlington/FloSelfbot && rm -rf /go && apk del go musl-dev go && chown -R nobody:nobody /data && chmod 700 /data

USER nobody
WORKDIR /data

ENTRYPOINT ["/usr/local/bin/FloSelfbot"]
