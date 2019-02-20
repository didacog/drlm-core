package lib

import (
	"context"
	"net"

	pb "github.com/brainupdaters/drlm-comm/drlmcomm"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DrlmapiConfig struct {
	Port string
}

type server struct{}

func (s *server) LoginUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Printf("Received login for user: %v", in.User)

	return &pb.SessionReply{Message: "Hello " + in.User + " / nova sessio "}, nil
}

func (s *server) AddUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Printf("Received add user - user: %v / pass: %v", in.User, in.Pass)
	u := User{User: in.User, Password: in.Pass}
	u.AddUser()
	return &pb.SessionReply{Message: "add user " + in.User + " to database"}, nil
}

func (s *server) DelUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Printf("Received delete user - user: %v", in.User)

	return &pb.SessionReply{Message: "delete user " + in.User + " from database"}, nil
}

func (s *server) ListUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Printf("Received list users")

	return &pb.SessionReply{Message: "lits users from database"}, nil
}

var s *grpc.Server

func InitDrlmapi(cfg DrlmapiConfig) {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDrlmApiServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
