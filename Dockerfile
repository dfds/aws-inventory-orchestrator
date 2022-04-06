FROM golang:1.18-alpine AS builder

RUN apk add --update --no-cache git

# Download Orchestrator modules
WORKDIR /orchestrator/
COPY ./orchestrator/go.mod ./
COPY ./orchestrator/go.sum ./
RUN go mod download

# Download Runner modules
WORKDIR /runner/
COPY ./runner/go.mod ./
COPY ./runner/go.sum ./
RUN go mod download

# Build Orchestrator
WORKDIR /orchestrator/
COPY ./orchestrator/*.go ./
RUN go build -o /orchestrator
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO11MODULE=on go build -mod=mod -a -o /main .

# Build Runner
WORKDIR /runner/
COPY ./runner/*.go ./
RUN go build -o /runner

# Copy binaries and run
FROM alpine:latest
COPY --from=builder /orchestrator /orchestrator
COPY --from=builder /runner /runner
CMD ["./orchestrator/orchestrator"]