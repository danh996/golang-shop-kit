gen:
	rm -rf pb && protoc --proto_path=proto proto/*.proto  --go_out=:. --go-grpc_out=:.

client-run:
	go run ./client/main.go
server-run:
	go run ./server/main.go