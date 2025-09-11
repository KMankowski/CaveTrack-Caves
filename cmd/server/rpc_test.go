package main

import (
	"context"
	"net"
	"testing"

	"github.com/KMankowski/CaveTrack-Caves/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestGetCaves(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterCavesServiceServer(s, &server{})
	go func() {
		s.Serve(listener)
	}()

	clientConn, err := grpc.DialContext(
		context.Background(),
		"bufconn",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	client := pb.NewCavesServiceClient(clientConn)
	response, err := client.GetCaves(context.Background(), &pb.GetCavesRequest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := "hi"
	got := response.Caves[0].Name
	if got != want {
		t.Errorf("want %s but got %s", want, got)
	}
	t.Logf("%s", got)
}
