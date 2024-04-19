grpc-server:
	go run server/main.go

grpc-client:
	go run client/main.go

all:
	protoc -I ./proto \
	--go_out ./proto --go_opt paths=source_relative \
	--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
	./proto/user/user.proto   

copyandpaste:
	export PATH=$PATH:$(go env GOPATH)/bin
