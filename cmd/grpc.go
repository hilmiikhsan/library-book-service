package cmd

import (
	"net"

	"github.com/hilmiikhsan/library-book-service/helpers"
	"google.golang.org/grpc"
)

func ServeGRPC() {
	server := grpc.NewServer()

	lis, err := net.Listen("tcp", ":"+helpers.GetEnv("GRPC_PORT", "6003"))
	if err != nil {
		helpers.Logger.Fatal("failed to listen grpc port: ", err)
	}

	helpers.Logger.Info("start listening grpc on port:" + helpers.GetEnv("GRPC_PORT", "6003"))
	if err := server.Serve(lis); err != nil {
		helpers.Logger.Fatal("failed to serve grpc port: ", err)
	}
}
