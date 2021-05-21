FROM golang:1.16.4-alpine AS amanBuild
WORKDIR /go/src/app
COPY ./aman .
RUN go env -w GO111MODULE=on
# RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go build -v -ldflags "-s -w" ./...

FROM gruebel/upx:latest as upxAman
COPY --from=amanBuild /go/src/app/aman .
RUN upx -9 aman

FROM alpine AS aria2cInstall 
RUN apk update && apk add aria2

FROM gruebel/upx:latest as upxAria2c
COPY --from=aria2cInstall /usr/bin/aria2c .
RUN upx -9 aria2c

FROM alpine
COPY --from=upxAman /aman .
COPY --from=upxAria2c /aria2c /usr/bin/aria2c
ENTRYPOINT [ "./aman" ]

EXPOSE 6800