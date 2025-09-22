.PHONY: all proto client server

all: client server

proto: proto/echo.proto
	protoc --go_opt=module=github.com/ruokeqx/grpcdemo/grpcoveruds --go_out=. \
		--go-grpc_opt=module=github.com/ruokeqx/grpcdemo/grpcoveruds --go-grpc_out=. \
		proto/echo.proto

proto/echopb/echo.pb.go: proto/echo.proto
	protoc --go_opt=module=github.com/ruokeqx/grpcdemo/grpcoveruds --go_out=. \
		--go-grpc_opt=module=github.com/ruokeqx/grpcdemo/grpcoveruds --go-grpc_out=. \
		proto/echo.proto

client: client/client.go proto/echopb/echo.pb.go
	cd bin && go build ../client

server: server/server.go proto/echopb/echo.pb.go
	cd bin && go build ../server
