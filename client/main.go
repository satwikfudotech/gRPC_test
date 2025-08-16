package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-system-monitor/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPerformanceServiceClient(conn)

	for {
		res, err := client.GetVitals(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatalf("Error calling GetVitals: %v", err)
		}
		fmt.Printf("📊 CPU Usage: %.2f%% | RAM Usage: %.2f MB\n", res.CpuUsage, res.RamUsage)
		time.Sleep(2 * time.Second)
	}
}
