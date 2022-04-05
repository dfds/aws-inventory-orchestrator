FROM golang:1.18-alpine AS builder

RUN apk add --update --no-cache git

ADD ./src /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO11MODULE=on go build -mod=mod -a -o /main .

FROM alpine:latest
COPY --from=builder /main /inventory-orchestrator
CMD ["./inventory-orchestrator"]