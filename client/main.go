package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	pb "grpc-system-monitor/proto"

	"google.golang.org/grpc"
)

func saveToCSV(record []string) {
	fileName := "client_metrics.csv"

	// Open or create the CSV file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the record
	if err := writer.Write(record); err != nil {
		log.Fatalf("Failed to write record to CSV: %v", err)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSystemMonitorClient(conn)

	// Write CSV header only once
	saveToCSV([]string{"TIMESTAMP", "CPU", "RAM", "NET IN", "NET OUT", "DISK_READ", "DISK_WRITE"})

	for {
		resp, err := c.GetMetrics(context.Background(), &pb.MetricsRequest{})
		if err != nil {
			log.Fatalf("Error calling GetMetrics: %v", err)
		}

		log.Printf("📊 CPU: %.2f%% | RAM: %.2f%% | Time: %s | NetIn: %.2f KB | NetOut: %.2f KB | DiskRead: %.2f KB | DiskWrite: %.2f KB",
			resp.CpuUsage, resp.RamUsage, resp.Timestamp, resp.NetInKb, resp.NetOutKb, resp.DiskReadKb, resp.DiskWriteKb)

		// Save record in CSV (client side)
		saveToCSV([]string{
			resp.Timestamp,
			fmt.Sprintf("%.2f", resp.CpuUsage),
			fmt.Sprintf("%.2f", resp.RamUsage),
			fmt.Sprintf("%.2f", resp.NetInKb),
			fmt.Sprintf("%.2f", resp.NetOutKb),
			fmt.Sprintf("%.2f", resp.DiskReadKb),
			fmt.Sprintf("%.2f", resp.DiskWriteKb),
		})

		time.Sleep(time.Duration(1000+rand.Intn(2000)) * time.Millisecond)
	}
}
