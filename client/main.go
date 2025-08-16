package main

import (
	"context"
	"log"
	"time"

	pb "grpc-system-monitor/proto"

	"google.golang.org/grpc"
)

func main() {
	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSystemMonitorClient(conn)

	for {
		resp, err := client.GetMetrics(context.Background(), &pb.MetricsRequest{})
		if err != nil {
			log.Fatalf("Error fetching metrics: %v", err)
		}

		log.Printf("📊 Metrics -> CPU: %.2f%% | RAM: %.2f%% | Time: %s | NetIn: %dB | NetOut: %dB | DiskRead: %dB | DiskWrite: %dB",
			resp.CpuUsage, resp.RamUsage, resp.Timestamp, resp.NetIn, resp.NetOut, resp.DiskRead, resp.DiskWrite)
		log.SetFlags(0) // removes default date/time prefix

		time.Sleep(5 * time.Second)
	}
}
