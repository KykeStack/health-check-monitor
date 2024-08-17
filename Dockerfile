FROM golang:1.10-alpine

LABEL maintainer="Agustín Houlgrave <a.houlgrave@gmail.com>"

ARG APPLICATION_VERSION_ARG
ENV APPLICATION_VERSION $APPLICATION_VERSION_ARG

WORKDIR /go/src/health-check-monitor/

COPY . .

RUN apk --no-cache add git && \
	go get

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

ENTRYPOINT ["./main"]

EXPOSE 8001
