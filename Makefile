server:
	cd ./examples/server && go run server.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative src/proto/api/protocol.proto