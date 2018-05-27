# Calls Protoc Library: AKA: Compiles Protobuf code.
build:
	# Updated to use go-micro plugin instead of grpc plugin.
	protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto 
	# Builds an image by the name consignment-service (Dot means that build process looks in current directory)
	#docker build -t dillonlpeterson/consignment .
	#docker push dillonlpeterson/consignment:latest 
	docker build -t consignment-service .
run: 
	docker run --net="host" \
		-p 50052 \
		-e MICRO_SERVER_ADDRESS=:50052 \
		-e MICRO_REGISTRY=mdns \
		consignment-service 



