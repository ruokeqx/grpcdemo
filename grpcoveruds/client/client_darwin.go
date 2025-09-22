package main

import (
	"context"
	"net"
)

const (
	socketPath = "/tmp/grpc_uds.sock"
	targetPath = "passthrough:///" + socketPath
)

// NewIpcConnection will connect to a Unix socket on the given endpoint.
func NewIpcConnection(ctx context.Context, endpoint string) (net.Conn, error) {
	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = context.WithTimeout(ctx, DefaultIpcDialTimeout)
	}

	var d net.Dialer
	return d.DialContext(ctx, "unix", endpoint)
}
