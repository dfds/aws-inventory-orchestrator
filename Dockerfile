FROM golang:1.19-alpine AS builder

RUN apk add --update --no-cache git

# Create required directories
RUN mkdir /src
RUN mkdir /src/orchestrator
RUN mkdir /src/runner
RUN mkdir /app

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

# Define variables
ARG USER=appuser
ENV HOME /home/$USER

# Copy binaries into the image
COPY --from=builder /src/orchestrator/bin $HOME/app
COPY --from=builder /src/runner/bin $HOME/app

# install sudo as root
RUN apk add --update sudo

# add new user
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

# set user and working directory
USER $USER
WORKDIR $HOME

# By default execute the orchestrator
ENTRYPOINT ["./app/orchestrator"]
