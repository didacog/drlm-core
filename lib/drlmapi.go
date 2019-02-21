package lib

import (
	"bytes"
	"context"
	"fmt"
	"net"

	pb "github.com/brainupdaters/drlm-comm/drlmcomm"
	"github.com/olekukonko/tablewriter"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DrlmapiConfig struct {
	Port string
}

type server struct{}

func (s *server) LoginUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Info("Received login for user: " + in.User)

	return &pb.SessionReply{Message: "Hello " + in.User + " / nova sessio "}, nil
}

func (s *server) AddUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Info("Received add user - user: " + in.User + " / pass: " + in.Pass)
	u := User{User: in.User, Password: in.Pass}
	u.AddUser()
	return &pb.SessionReply{Message: "add user " + in.User + " to database"}, nil
}

func (s *server) DelUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Info("Received delete user - user: " + in.User)
	u := User{}
	u.LoadUser(in.User)
	u.Delete()
	return &pb.SessionReply{Message: "delete user " + in.User + " from database"}, nil
}

func (s *server) ListUser(ctx context.Context, in *pb.UserRequest) (*pb.SessionReply, error) {
	log.Info("Received list users")

	users := []User{}
	DBConn.Find(&users)

	buf := new(bytes.Buffer)

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Id", "Created At", "Updated At", "User", "Pass"})

	for i, _ := range users {
		table.Append([]string{fmt.Sprint(users[i].ID), users[i].CreatedAt.String(), users[i].UpdatedAt.String(), users[i].User, users[i].Password})
	}

	table.Render()

	return &pb.SessionReply{Message: "\n" + buf.String()}, nil
}

var s *grpc.Server

func InitDrlmapi(cfg DrlmapiConfig) {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal("failed to listen: " + err.Error())
	}
	s := grpc.NewServer()
	pb.RegisterDrlmApiServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: " + err.Error())
	}
}
