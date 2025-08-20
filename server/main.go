package main

import (
	"context"
	"log"
	"net"
	"sort"
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

// --- Median calculation ---
func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values) // ensure sorted
	n := len(values)
	if n%2 == 1 {
		return values[n/2]
	}
	return (values[n/2-1] + values[n/2]) / 2
}

// --- gRPC Implementation ---
func (s *server) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	// CPU usage samples (takes 1s interval)
	cpuPercent, _ := cpu.Percent(time.Second, false)
	cpuMedian := median(cpuPercent)

	// RAM usage
	vm, _ := mem.VirtualMemory()
	ramUsage := vm.UsedPercent

	// Timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Network throughput
	netIO, _ := gopsnet.IOCounters(false)
	var netIn, netOut float64
	if len(netIO) > 0 {
		netIn = float64(netIO[0].BytesRecv) / 1024 // KB
		netOut = float64(netIO[0].BytesSent) / 1024
	}

	// Disk I/O
	diskIO, _ := disk.IOCounters()
	var readKB, writeKB float64
	for _, io := range diskIO {
		readKB = float64(io.ReadBytes) / 1024
		writeKB = float64(io.WriteBytes) / 1024
		break // take first disk only
	}

	return &pb.MetricsResponse{
		CpuUsage:    cpuMedian,
		RamUsage:    ramUsage,
		Timestamp:   timestamp,
		NetInKb:     netIn,
		NetOutKb:    netOut,
		DiskReadKb:  readKB,
		DiskWriteKb: writeKB,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSystemMonitorServer(s, &server{})

	log.Println("🚀 Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
