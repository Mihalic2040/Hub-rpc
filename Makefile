server:
	cd ./examples/server && go run server.go
hand:
	cd ./examples/handlers && go run handler.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative src/proto/api/protocol.proto