package main

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
)

const (
	socketPath = `\\.\pipe\grpc_namedpipe`
	targetPath = "passthrough:///" + socketPath
)

// NewIpcConnection will connect to a named pipe with the given endpoint as name.
func NewIpcConnection(ctx context.Context, endpoint string) (net.Conn, error) {
	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = context.WithTimeout(ctx, DefaultIpcDialTimeout)
	}
	return winio.DialPipeContext(ctx, endpoint)
}
