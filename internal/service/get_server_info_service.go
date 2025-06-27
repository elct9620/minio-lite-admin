package service

import (
	"context"
	"fmt"

	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

type GetServerInfoService struct {
	minioClient *madmin.AdminClient
}

// ServerInfo represents basic MinIO server information
type ServerInfo struct {
	Mode         string `json:"mode"`
	Region       string `json:"region"`
	DeploymentID string `json:"deploymentId"`
}

// DiskUsage represents disk usage information extracted from ServerInfo
type DiskUsage struct {
	TotalCapacity     uint64     `json:"totalCapacity"`
	TotalUsedCapacity uint64     `json:"totalUsedCapacity"`
	TotalFreeCapacity uint64     `json:"totalFreeCapacity"`
	UsagePercentage   float64    `json:"usagePercentage"`
	OnlineDisks       int        `json:"onlineDisks"`
	OfflineDisks      int        `json:"offlineDisks"`
	HealingDisks      int        `json:"healingDisks"`
	PoolsCount        int        `json:"poolsCount"`
	ObjectsCount      uint64     `json:"objectsCount"`
	BucketsCount      uint64     `json:"bucketsCount"`
	DiskDetails       []DiskInfo `json:"diskDetails"`
}

// DiskInfo represents individual disk information
type DiskInfo struct {
	Endpoint       string  `json:"endpoint"`
	State          string  `json:"state"`
	TotalSpace     uint64  `json:"totalSpace"`
	UsedSpace      uint64  `json:"usedSpace"`
	AvailableSpace uint64  `json:"availableSpace"`
	Utilization    float64 `json:"utilization"`
	Healing        bool    `json:"healing"`
}

// CombinedServerInfo contains both server info and disk usage
type CombinedServerInfo struct {
	ServerInfo *ServerInfo `json:"serverInfo"`
	DiskUsage  *DiskUsage  `json:"diskUsage"`
}

func NewGetServerInfoService(minioClient *madmin.AdminClient) *GetServerInfoService {
	return &GetServerInfoService{
		minioClient: minioClient,
	}
}

func (s *GetServerInfoService) Execute(ctx context.Context) (*CombinedServerInfo, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Fetching MinIO server info")

	info, err := s.minioClient.ServerInfo(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch MinIO server info")
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	logger.Debug().Msg("Successfully fetched MinIO server info")

	// Extract basic server info
	serverInfo := &ServerInfo{
		Mode:         info.Mode,
		Region:       info.Region,
		DeploymentID: info.DeploymentID,
	}

	// Extract disk usage information
	diskUsage := s.extractDiskUsage(info)

	return &CombinedServerInfo{
		ServerInfo: serverInfo,
		DiskUsage:  diskUsage,
	}, nil
}

// extractDiskUsage extracts disk usage information from MinIO ServerInfo
func (s *GetServerInfoService) extractDiskUsage(info madmin.InfoMessage) *DiskUsage {
	var totalCapacity, totalUsed, totalFree uint64
	var onlineDisks, offlineDisks, healingDisks int
	var diskDetails []DiskInfo

	// Extract disk information from servers
	for _, server := range info.Servers {
		for _, disk := range server.Disks {
			diskDetails = append(diskDetails, DiskInfo{
				Endpoint:       disk.Endpoint,
				State:          disk.State,
				TotalSpace:     disk.TotalSpace,
				UsedSpace:      disk.UsedSpace,
				AvailableSpace: disk.AvailableSpace,
				Utilization:    disk.Utilization,
				Healing:        disk.Healing,
			})

			// Aggregate totals
			totalCapacity += disk.TotalSpace
			totalUsed += disk.UsedSpace
			totalFree += disk.AvailableSpace

			// Count disk states
			switch disk.State {
			case "ok":
				onlineDisks++
			case "offline":
				offlineDisks++
			default:
				// Consider other states as potentially problematic
				if disk.Healing {
					healingDisks++
				} else {
					offlineDisks++
				}
			}
		}
	}

	// Calculate usage percentage
	var usagePercentage float64
	if totalCapacity > 0 {
		usagePercentage = float64(totalUsed) / float64(totalCapacity) * 100
	}

	// Count pools - use TotalSets array length as each index represents a pool
	poolsCount := len(info.Backend.TotalSets)

	return &DiskUsage{
		TotalCapacity:     totalCapacity,
		TotalUsedCapacity: totalUsed,
		TotalFreeCapacity: totalFree,
		UsagePercentage:   usagePercentage,
		OnlineDisks:       onlineDisks,
		OfflineDisks:      offlineDisks,
		HealingDisks:      healingDisks,
		PoolsCount:        poolsCount,
		ObjectsCount:      info.Objects.Count,
		BucketsCount:      info.Buckets.Count,
		DiskDetails:       diskDetails,
	}
}
