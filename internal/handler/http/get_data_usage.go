package http

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

// DataUsageResponse represents the MinIO data usage response
type DataUsageResponse struct {
	TotalCapacity     uint64     `json:"totalCapacity"`
	TotalFreeCapacity uint64     `json:"totalFreeCapacity"`
	TotalUsedCapacity uint64     `json:"totalUsedCapacity"`
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

// GetDataUsageHandler handles MinIO data usage requests
func (s *Service) GetDataUsageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("Fetching MinIO data usage information")

	w.Header().Set("Content-Type", "application/json")

	// Use the same ServerInfo API call to get disk usage information
	combinedInfo, err := s.getServerInfoService.Execute(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get MinIO data usage info")
		http.Error(w, "Failed to get MinIO data usage info", http.StatusInternalServerError)
		return
	}

	diskUsage := combinedInfo.DiskUsage

	// Convert service DiskInfo to handler DiskInfo
	diskDetails := make([]DiskInfo, len(diskUsage.DiskDetails))
	for i, disk := range diskUsage.DiskDetails {
		diskDetails[i] = DiskInfo{
			Endpoint:       disk.Endpoint,
			State:          disk.State,
			TotalSpace:     disk.TotalSpace,
			UsedSpace:      disk.UsedSpace,
			AvailableSpace: disk.AvailableSpace,
			Utilization:    disk.Utilization,
			Healing:        disk.Healing,
		}
	}

	response := DataUsageResponse{
		TotalCapacity:     diskUsage.TotalCapacity,
		TotalFreeCapacity: diskUsage.TotalFreeCapacity,
		TotalUsedCapacity: diskUsage.TotalUsedCapacity,
		UsagePercentage:   diskUsage.UsagePercentage,
		OnlineDisks:       diskUsage.OnlineDisks,
		OfflineDisks:      diskUsage.OfflineDisks,
		HealingDisks:      diskUsage.HealingDisks,
		PoolsCount:        diskUsage.PoolsCount,
		ObjectsCount:      diskUsage.ObjectsCount,
		BucketsCount:      diskUsage.BucketsCount,
		DiskDetails:       diskDetails,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode MinIO data usage response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
