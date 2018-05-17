# Calls Protoc Library: AKA: Compiles Protobuf code.
build:
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/bitbucket.org/dillonlpeterson/shippy/consignment-service proto/consignment/consignment.proto 
	# Sets build architecture and builds to that
	GOOS=linux GOARCH=amd64 go build 
	# Builds an image by the name consignment-service (Dot means that build process looks in current directory)
	docker build -t consignment-service .
run: 
	docker run -p 50051:50051 consignment-service
	


