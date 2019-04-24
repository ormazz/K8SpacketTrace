# Build the manager binary
FROM golang:1.11.0 as builder

# Copy in the go 
WORKDIR /go/src/github.com/ormazz/K8SpacketTrace
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o packetTraceAgent github.com/ormazz/K8SpacketTrace/cmd/packetTrace-agent

# Copy the agent into a thin image
FROM centos:latest
WORKDIR /root/
RUN yum install tcpdump iproute -y
#RUN dnf install tcpdump -y
COPY --from=builder /go/src/github.com/ormazz/K8SpacketTrace/packetTraceAgent .
ENTRYPOINT ["./packetTraceAgent"]