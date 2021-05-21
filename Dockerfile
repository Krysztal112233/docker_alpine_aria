FROM golang:1.16.4-alpine AS amanBuild
WORKDIR /go/src/app
COPY ./aman .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go build -v -ldflags "-s -w" ./...

FROM gruebel/upx:latest as upxAman
COPY --from=amanBuild /go/src/app/aman .
RUN upx -9 aman

FROM alpine
RUN apk update && apk add aria2
COPY --from=upxAman /aman /
ENTRYPOINT [ "./aman" ]   