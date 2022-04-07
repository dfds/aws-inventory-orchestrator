FROM golang:1.18-alpine AS builder

RUN apk add --update --no-cache git

# Create required directories
RUN mkdir /src
RUN mkdir /src/orchestrator
RUN mkdir /src/runner
RUN mkdir /app
RUN mkdir ./monkeymagic

# Download Orchestrator modules
WORKDIR /src/orchestrator/
COPY ./orchestrator/go.mod ./
COPY ./orchestrator/go.sum ./
RUN go mod download

# Download Runner modules
WORKDIR /src/runner/
COPY ./runner/go.mod ./
COPY ./runner/go.sum ./
RUN go mod download

# Build Orchestrator
WORKDIR /src/orchestrator/
COPY ./orchestrator/ ./
RUN go build -o bin/orchestrator
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO11MODULE=on go build -mod=mod -a -o /main .

# Build Runner
WORKDIR /src/runner/
COPY ./runner/ ./
RUN go build -o bin/runner

# Copy binaries and run
FROM alpine:latest
COPY --from=builder /src/orchestrator/bin /app
COPY --from=builder /src/runner/bin /app
ENTRYPOINT ["/app/orchestrator"]
