# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang:1.19-alpine3.17 as build

WORKDIR $GOPATH/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "linux" -o protocol-adapter -ldflags="-w" -a -v  ./cmd/protocol-adapter


FROM alpine:3.17
ENV GOPATH="/go/src"
WORKDIR /run

COPY --from=build $GOPATH/app/protocol-adapter .
COPY --from=build $GOPATH/app/cmd/protocol-adapter/conf/* conf/
EXPOSE 5141/tcp
EXPOSE 5141/udp
EXPOSE 5143/tcp
EXPOSE 4444/tcp

ENTRYPOINT ["/run/protocol-adapter", "--cf", "./conf"]
