package main

import "net"

const (
	socketPath = "/tmp/grpc_uds.sock"
)

func NewIpcListener(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}
