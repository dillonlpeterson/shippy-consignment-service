# Calls Protoc Library: AKA: Compiles Protobuf code.
build:
	# Updated to use go-micro plugin instead of grpc plugin.
	protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto 
	# Builds an image by the name consignment-service (Dot means that build process looks in current directory)
	docker build -t shippy-consignment-service:latest .
	#docker push dillonlpeterson/shippy-consignment-service:latest 
run: 
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns shippy-consignment-service 



