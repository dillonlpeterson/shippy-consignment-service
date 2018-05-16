# Calls Protoc Library: AKA: Compiles Protobuf code.
build:
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/bitbucket.org/dillonlpeterson/shippy/consignment-service proto/consignment/consignment.proto 

