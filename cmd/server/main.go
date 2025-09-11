package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/KMankowski/CaveTrack-Caves/internal/pb"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("Program exiting with error", slog.String("error", err.Error()))
		os.Exit(1)
	} else {
		logger.Info("Program exiting successfully")
		os.Exit(0)
	}
}

func run(logger *slog.Logger) error {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	logger.Info("Listening...", slog.String("port", ":50051"))
	pb.RegisterCavesServiceServer(s, &server{})

	return s.Serve(l)
}
