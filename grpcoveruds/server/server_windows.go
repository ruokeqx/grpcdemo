package main

import (
	"net"

	"github.com/Microsoft/go-winio"
)

const (
	socketPath = `\\.\pipe\grpc_namedpipe`
)

func NewIpcListener(path string) (net.Listener, error) {
	return winio.ListenPipe(path, nil)
}
