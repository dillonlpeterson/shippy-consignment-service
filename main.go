package main

import (
	"log"

	"golang.org/x/net/context"

	vesselProto "bitbucket.org/dillonlpeterson/shippy/vessel-service/proto/vessel"

	pb "bitbucket.org/dillonlpeterson/shippy/consignment-service/proto/consignment"
	micro "github.com/micro/go-micro"
)

const (
	port = ":50051"
)

// Repository is interface for ConsignmentRepository
type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// ConsignmentRepository is a dummy Repo
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

// Create created a consignment object.
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// GetAll returns all existing Consignments
func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generate code itself for the exact method signatures etc. to
// give you a better idea
type service struct {
	repo         Repository
	VesselClient vesselProto.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Shoutout to our VesselService
	vesselResponse, err := s.VesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	// We set the VesselID as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	// Return the matching 'Response' message we created in our protobuf definition.
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	// Set-up our gRPC server.
	srv := micro.NewService(
		// Must match package name given in proto file!
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
