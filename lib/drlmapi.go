package lib

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"

	pb "github.com/brainupdaters/drlm-common/comms"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type DrlmapiConfig struct {
	Server   string
	Port     string
	Tls      bool
	Cert     string
	Key      string
	User     string
	Password string
}

func SetDrlmapiConfigDefaults() {
	viper.SetDefault("drlmapi.server", "godev")
	viper.SetDefault("drlmapi.port", 50051)
	viper.SetDefault("drlmapi.tls", false)
	viper.SetDefault("drlmapi.cert", "/tls/godev/godev.crt")
	viper.SetDefault("drlmapi.key", "/tls/godev/godev.key")
	viper.SetDefault("drlmapi.user", "drlmadmin")
	viper.SetDefault("drlmapi.password", "drlm3api")
}

type server struct {
}

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context, s *server) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")
		if clientLogin != Config.Drlmapi.User {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}
		if clientPassword != Config.Drlmapi.Password {
			return "", fmt.Errorf("bad password %s", clientPassword)
		}
		log.Printf("authenticated client: %s", clientLogin)
		return "42", nil
	}
	return "", fmt.Errorf("missing credentials")
}

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := info.Server.(*server)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}
	clientID, err := authenticateClient(ctx, s)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

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

func InitDrlmapi() {
	// create a listener on TCP port
	lis, err := net.Listen("tcp", Config.Drlmapi.Server+":"+Config.Drlmapi.Port)
	if err != nil {
		log.Fatal("failed to listen: " + err.Error())
	}

	// create a server instance
	s := server{}
	var creds credentials.TransportCredentials
	var opts []grpc.ServerOption

	if Config.Drlmapi.Tls {
		// Create the TLS credentials
		creds, err = credentials.NewServerTLSFromFile(Config.Drlmapi.Cert, Config.Drlmapi.Key)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		// Create an array of gRPC options with the credentials
		opts = []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor)}
	} else {
		opts = []grpc.ServerOption{grpc.UnaryInterceptor(unaryInterceptor)}
	}
	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// attach the DrlmApi service to the server
	pb.RegisterDrlmApiServer(grpcServer, &s)

	// start the server
	log.Info("starting HTTP/2 gRPC server on " + lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
