package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	pb "github.com/dillonlpeterson/shippy-consignment-service/proto/consignment"
	userService "github.com/dillonlpeterson/shippy-user-service/proto/auth"
	vesselProto "github.com/dillonlpeterson/shippy-vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
)

const (
	port        = ":50051"
	defaultHost = "localhost:27017"
)

var (
	srv micro.Service
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
	srv := k8s.NewService(
		// Must match package name given in proto file!
		micro.Name("shippy.consignment"),
		micro.Version("latest"),
		// Our auth middleware
		micro.WrapHandler(AuthWrapper),
	)

	vesselClient := vesselProto.NewVesselServiceClient("shippy.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &handler{session, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

// AuthWrapper is a high-order function which takes a HandlerFunc
// and returns a function, which takes a context, request and response interface.
// The token is extracted from the context set in our consignment-cli, that
// token is then sent over to the user service to be validated.
// If valid, the call is passed along to the handler. If not,
// an error is returned.
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, res interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, res)
		}

		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]
		fmt.Println("Authenticating with token: ", token)

		// Auth here
		authClient := userService.NewAuthClient("shippy.auth", srv.Client())
		_, err := authClient.ValidateToken(context.Background(), &userService.Token{
			Token: token,
		})
		if err != nil {
			return err
		}

		return fn(ctx, req, res)
	}
}
