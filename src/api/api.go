package api

import (
	"fmt"
	"logger"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var grpcServer *grpc.Server

const tag string = "SERVER"

type Server struct {
	UnimplementedEventsServer
	events chan Event
}


func (s *Server) NewEvent(ctx context.Context, in *Event) (*Empty, error) {
	logger.LogInfo(fmt.Sprintf("[%s] New event %d", in.Date, in.EventType),
				   tag)
	s.events <-*in
	return &Empty{}, nil
}


func InitServer(port string, events chan Event) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to listen to tcp%s", port), tag)
		panic(fmt.Sprintf("Failed to listen to tcp %s", port))
	}

	grpcServer = grpc.NewServer()
	s := Server{events: events}
	logger.LogInfo("Registering server...", tag)
	RegisterEventsServer(grpcServer, &s)

	err = grpcServer.Serve(lis)
	if err != nil {
		logger.LogError("Failed to start grpc server", tag)
		panic("Failed to start grpc server")
	}

}
