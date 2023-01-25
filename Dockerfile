FROM golang:1.19-alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" cmd/exporter/main.go

FROM scratch
COPY --from=builder /go/src/app/main /routeros-exporter
EXPOSE 9436
ENTRYPOINT ["/routeros-exporter"]
