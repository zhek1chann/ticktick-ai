FROM golang:1.24 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.24 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /workdir
WORKDIR /workdir
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/app ./cmd

FROM ubuntu:noble
RUN apt-get update && apt-get install -y ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/app /workdir/app
WORKDIR /workdir

CMD ["./app"]
