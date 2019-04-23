# Build the manager binary
FROM golang:1.11.0 as builder

# Copy in the go 
WORKDIR /go/src/github.com/ormazz/K8SpacketTrace
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o packet_trace_api github.com/ormazz/K8SpacketTrace/cmd/packetTrace-api

# Copy the agent into a thin image
FROM fedora:29
WORKDIR /root/
COPY --from=builder github.com/ormazz/K8SpacketTrace/cmd/packetTrace-api .
ENTRYPOINT ["./packetTrace-api"]