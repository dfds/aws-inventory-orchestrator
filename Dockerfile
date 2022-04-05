FROM golang:1.18-alpine AS builder

RUN apk add --update --no-cache git

#ADD ./src /src
WORKDIR /app
COPY ./src/go.mod ./
COPY ./src/go.sum ./
RUN go mod download

COPY ./src/*.go ./
RUN go build -o /main
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO11MODULE=on go build -mod=mod -a -o /main .

FROM alpine:latest
COPY --from=builder /main /inventory-orchestrator
CMD ["./inventory-orchestrator"]