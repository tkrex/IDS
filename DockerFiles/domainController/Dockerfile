# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
ARG domain=default

# Copy the local package files to the container's workspace.
#ADD . /go/src/github.com/tkrex/IDS

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get -u github.com/tkrex/IDS/domainController/...
RUN go install github.com/tkrex/IDS/domainController/...

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/IDSDomainController domain
EXPOSES 8080
