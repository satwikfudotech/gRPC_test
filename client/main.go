package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "grpc-system-monitor/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSystemMonitorClient(conn)

	for {
		resp, err := c.GetMetrics(context.Background(), &pb.MetricsRequest{})
		if err != nil {
			log.Fatalf("Error calling GetMetrics: %v", err)
		}

		log.Printf("📊 CPU: %.2f%% | RAM: %.2f%% | Time: %s | NetIn: %.2f KB | NetOut: %.2f KB | DiskRead: %.2f KB | DiskWrite: %.2f KB",
			resp.CpuUsage, resp.RamUsage, resp.Timestamp, resp.NetInKb, resp.NetOutKb, resp.DiskReadKb, resp.DiskWriteKb)

		time.Sleep(time.Duration(1000+rand.Intn(2000)) * time.Millisecond)
	}
}
