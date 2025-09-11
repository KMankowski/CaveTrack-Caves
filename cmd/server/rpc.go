package main

import (
	"context"

	"github.com/KMankowski/CaveTrack-Caves/internal/pb"
)

type server struct {
	pb.UnimplementedCavesServiceServer
}

func (s *server) GetCaves(_ context.Context, r *pb.GetCavesRequest) (*pb.GetCavesResponse, error) {
	c := pb.Cave{Name: "hi"}
	caves := []*pb.Cave{&c}
	return &pb.GetCavesResponse{Caves: caves}, nil
}
