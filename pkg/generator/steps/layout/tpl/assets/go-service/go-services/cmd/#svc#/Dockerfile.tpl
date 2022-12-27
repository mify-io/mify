FROM golang:1.19 AS build
WORKDIR /go/src
COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o {{.ServiceName}} ./cmd/{{.ServiceName}}

FROM alpine:3.15 AS runtime
ENV {{.GetApiEndpointEnvName}}=:80
COPY --from=build /go/src/{{.ServiceName}} ./
EXPOSE 80/tcp
ENTRYPOINT ["./{{.ServiceName}}"]
