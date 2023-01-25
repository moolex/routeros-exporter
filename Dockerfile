FROM golang:1.19-alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o exporter cmd/exporter/main.go

FROM scratch
COPY --from=builder /go/src/app/exporter /exporter
EXPOSE 9436
ENTRYPOINT ["/exporter"]
