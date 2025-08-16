package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "grpc-system-monitor/proto"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	gopsnet "github.com/shirou/gopsutil/v3/net"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSystemMonitorServer
}

func (s *server) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	// CPU
	cpuPercent, _ := cpu.Percent(0, false)

	// RAM
	vmStat, _ := mem.VirtualMemory()

	// Net I/O
	netIO, _ := gopsnet.IOCounters(false)

	// Disk I/O
	diskIO, _ := disk.IOCounters()

	// Get one disk (sda / C: for Windows)
	var readBytes, writeBytes uint64
	for _, io := range diskIO {
		readBytes = io.ReadBytes
		writeBytes = io.WriteBytes
		break
	}

	// Response
	return &pb.MetricsResponse{
		CpuUsage:  cpuPercent[0],
		RamUsage:  vmStat.UsedPercent,
		Timestamp: time.Now().Format(time.RFC3339),
		NetIn:     netIO[0].BytesRecv,
		NetOut:    netIO[0].BytesSent,
		DiskRead:  readBytes,
		DiskWrite: writeBytes,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051") // Laptop as server
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSystemMonitorServer(grpcServer, &server{})

	log.Println("Server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
