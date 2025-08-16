package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"

	pb "grpc-system-monitor/proto"
	"google.golang.org/grpc"
)

type performanceServer struct {
	pb.UnimplementedPerformanceServiceServer
}

func (s *performanceServer) GetVitals(ctx context.Context, req *pb.Empty) (*pb.VitalsResponse, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Dummy CPU usage (real CPU usage would need more tracking)
	cpuUsage := float32(25.0) // Hardcoded example
	ramUsage := float32(m.Alloc) / (1024 * 1024) // Convert to MB

	return &pb.VitalsResponse{
		CpuUsage: cpuUsage,
		RamUsage: ramUsage,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPerformanceServiceServer(grpcServer, &performanceServer{})

	fmt.Println("🚀 Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
