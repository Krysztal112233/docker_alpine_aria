FROM golang:1.16.4-alpine AS amanBuild
WORKDIR /go/src/app
COPY ./aman .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go build -v -ldflags "-s -w" ./...

FROM debian as upxAman
RUN apt update && apt install -y upx
COPY --from=amanBuild /go/src/app/aman .
RUN upx -9 aman

FROM alpine AS aria2cInstall 
RUN apk update && apk add aria2

FROM debian as upxAria2c
RUN apt update && apt install -y upx
COPY --from=aria2cInstall /usr/bin/aria2c .
RUN upx -9 aria2c

FROM alpine
COPY --from=upxAman /aman .
RUN apk add --update --no-cache aria2 && rm -rf /var/cache/apk/*
COPY --from=upxAria2c /aria2c /usr/bin/aria2c
ENTRYPOINT [ "./aman" ]

EXPOSE 8090 6800 6800/udp