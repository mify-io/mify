FROM golang:1.16 AS build
WORKDIR /go/src
COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o service2 ./cmd/service2

FROM alpine:3.15 AS runtime
ENV SERVICE2_SERVER_ENDPOINT=:80
COPY --from=build /go/src/service2 ./
EXPOSE 80/tcp
ENTRYPOINT ["./service2"]
