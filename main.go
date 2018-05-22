package main

import (
	"fmt"
	"log"
	"os"

	vesselProto "github.com/dillonlpeterson/shippy-vessel-service/proto/vessel"

	pb "github.com/dillonlpeterson/shippy-consignment-service/proto/consignment"
	micro "github.com/micro/go-micro"
)

const (
	port        = ":50051"
	defaultHost = "localhost:27017"
)

func main() {
	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}
	session, err := CreateSession(host)
	defer session.Close()

	if err != nil {
		// We're wrapping the error returned from our CreateSession
		// here to add some context to the error
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	// Set-up our gRPC server.
	srv := micro.NewService(
		// Must match package name given in proto file!
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &handler{session, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
