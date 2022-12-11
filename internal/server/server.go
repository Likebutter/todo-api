package server

import (
	"context"
	"log"

	api "github.com/Likebutter/todo-api/api/v1"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type grpcServer struct {
	api.UnimplementedActivitiesApiServer
	Activities *Activities
}

func (s *grpcServer) Get(ctx context.Context, req *api.GetActivityRequest) (*api.Activity, error) {
	resp, err := s.Activities.Get(int(req.Id))

	if err == ErrIDNotFound {
		return nil, status.Error(codes.NotFound, "id was not found")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (s *grpcServer) Insert(ctx context.Context, activity *api.Activity) (*api.Activity, error) {
	res, err := s.Activities.Insert(activity)

	if err != nil {
		log.Printf("Error:%s", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *grpcServer) List(ctx context.Context, req *api.ListActivitiesRequest) (*api.Activities, error) {
	activities, err := s.Activities.List(int(req.Offset))

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.Activities{Activities: activities}, nil
}

func NewGRPCServer() *grpc.Server {
	var acc *Activities
	var err error

	if acc, err = NewActivities(); err != nil {
		log.Fatal(err)
	}

	gsrv := grpc.NewServer()

	srv := grpcServer{Activities: acc}

	api.RegisterActivitiesApiServer(gsrv, &srv)

	return gsrv
}
