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
	var log string = fmt.Sprintf("New event %d at %s.\n"+
				                 "Host: %s\nUser: %s\nPwd: %s\n"+
								 "Cmd: %s\nPid: %s\nNotes: %s",
								 in.EventType, in.Date, in.HostId, in.User,
								 in.Pwd, in.Cmd, in.Pid, in.Notes)
	logger.LogInfo(log, tag)
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
