FROM golang:1.25-alpine AS builder
WORKDIR /go/src/github.com/OwO-Network/DeepLX
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags="-s -w" -installsuffix cgo -o deeplx .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
RUN addgroup -S deeplx && adduser -h /app -G deeplx -SH deeplx
USER deeplx:deeplx
COPY --from=builder --chown=deeplx:deeplx /go/src/github.com/OwO-Network/DeepLX/deeplx /app/deeplx
EXPOSE 1188
ENTRYPOINT ["/app/deeplx"]
