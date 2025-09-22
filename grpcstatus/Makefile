.PHONY: all proto client server

all: client server

proto: proto/tray/status.proto
	protoc --go_opt=module=github.com/ruokeqx/grpcdemo/grpcstatus --go_out=. \
		--go-grpc_opt=module=github.com/ruokeqx/grpcdemo/grpcstatus --go-grpc_out=. \
		proto/tray/status.proto

proto/tray/statuspb/status.pb.go: proto/tray/status.proto
	protoc --go_opt=module=github.com/ruokeqx/grpcdemo/grpcstatus --go_out=. \
		--go-grpc_opt=module=github.com/ruokeqx/grpcdemo/grpcstatus --go-grpc_out=. \
		proto/tray/status.proto

client: client/client.go proto/tray/statuspb/status.pb.go
	cd bin && go build ../client

server: server/server.go proto/tray/statuspb/status.pb.go
	cd bin && go build ../server
