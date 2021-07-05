FROM golang:1.16-alpine as builder

WORKDIR /build
ADD . /build/
RUN apk add gcc g++
RUN go build -ldflags "-w -s --extldflags '-static -fpic'" -tags "fts5 gin gorm" -o build/output/memo github.com/x2ox/memo/cmd/server

FROM alpine

RUN apk add tzdata
COPY --from=builder /build/build/output/memo /usr/bin/memo
WORKDIR /data/memo
VOLUME /data/memo
EXPOSE 8088
ENTRYPOINT ["memo"]
