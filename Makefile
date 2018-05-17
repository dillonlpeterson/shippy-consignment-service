# Calls Protoc Library: AKA: Compiles Protobuf code.
build:
	# Updated to use go-micro plugin instead of grpc plugin.
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/bitbucket.org/dillonlpeterson/shippy/consignment-service proto/consignment/consignment.proto 
	# Sets build architecture and builds to that
	GOOS=linux GOARCH=amd64 go build 
	# Builds an image by the name consignment-service (Dot means that build process looks in current directory)
	docker build -t consignment-service .
run: 
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns consignment-service



