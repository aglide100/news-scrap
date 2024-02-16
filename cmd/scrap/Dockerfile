FROM golang:1.20-alpine3.17 AS builder

RUN apk update && apk add upx && apk --no-cache add ca-certificates

WORKDIR "/go/src/github.com/aglide100/hate-ai"

COPY go.mod .
COPY go.sum .

RUN mkdir -p /opt/bin/scrap/

RUN go mod tidy \
    && go get -u -d -v ./...

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o /opt/bin/scrap/main cmd/scrap/main.go \    
    && upx --lzma /opt/bin/scrap/main

FROM debian:stretch-slim AS runtime
COPY --from=builder /opt/bin/scrap /opt/bin/scrap
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "/opt/bin/scrap/main" ]