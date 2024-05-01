ARG  BUILDER_IMAGE=golang:1.22.2-alpine

FROM ${BUILDER_IMAGE} as builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /go/src/go-app/

# Copy the Go modules files
COPY go.mod go.sum /go/src/go-app/

# Download the Go modules
RUN go mod download

# Copy the rest of the project files
COPY . /go/src/go-app/

# Build
RUN go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/go-app ./cmd

# This results in a single layer image
FROM scratch

LABEL org.opencontainers.image.source=https://github.com/parinyapt/prinflix_backend

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/go-app /go/bin/go-app

# Use an unprivileged user.
USER appuser:appuser

# Run binary
ENTRYPOINT ["/go/bin/go-app"]
CMD ["-mode=production"]

# REF https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324