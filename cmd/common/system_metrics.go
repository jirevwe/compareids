package common

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMetrics holds CPU and RAM utilization metrics
type SystemMetrics struct {
	CPUUsagePercent float64
	RAMUsageMB      float64
	RAMUsagePercent float64
	TotalRAMMB      float64
}

// AsMap formats the system metrics as a map of strings
func (s *SystemMetrics) AsMap() map[string]string {
	return map[string]string{
		"cpu_usage_percent": fmt.Sprintf("%.2f", s.CPUUsagePercent),
		"ram_usage_mb":      fmt.Sprintf("%.2f", s.RAMUsageMB),
		"ram_usage_percent": fmt.Sprintf("%.2f", s.RAMUsagePercent),
		"total_ram_mb":      fmt.Sprintf("%.2f", s.TotalRAMMB),
	}
}

// MeasureSystemResources measures CPU and RAM utilization during a function execution
// It returns the average CPU and RAM utilization during the execution
func MeasureSystemResources(fn func() error) (*SystemMetrics, error) {
	// Create channels to collect metrics
	cpuMetrics := make(chan float64, 100)
	ramMetrics := make(chan *mem.VirtualMemoryStat, 100)

	// Start a goroutine to collect metrics every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				// Measure CPU usage
				cpuPercent, err := cpu.Percent(0, false) // false = all cores combined
				if err == nil && len(cpuPercent) > 0 {
					cpuMetrics <- cpuPercent[0]
				}

				// Measure RAM usage
				memInfo, err := mem.VirtualMemory()
				if err == nil {
					ramMetrics <- memInfo
				}
			}
		}
	}(ticker)

	// Execute the function
	err := fn()

	// Stop the metrics collection
	ticker.Stop()
	close(cpuMetrics)
	close(ramMetrics)

	// Calculate average CPU usage
	var totalCPU float64
	var cpuCount int
	for cpuMetric := range cpuMetrics {
		totalCPU += cpuMetric
		cpuCount++
		if cpuCount >= cap(cpuMetrics) {
			break
		}
	}

	// Calculate average RAM usage
	var totalRAM float64
	var totalRAMPercent float64
	var totalRAMMB float64
	var ramCount int
	for ram := range ramMetrics {
		totalRAM += float64(ram.Used)
		totalRAMPercent += ram.UsedPercent
		totalRAMMB = float64(ram.Total) / 1024 / 1024 // Convert to MB
		ramCount++
		if ramCount >= cap(ramMetrics) {
			break
		}
	}

	// Calculate averages
	var avgCPU, avgRAM, avgRAMPercent float64
	if cpuCount > 0 {
		avgCPU = totalCPU / float64(cpuCount)
	}
	if ramCount > 0 {
		avgRAM = totalRAM / float64(ramCount) / 1024 / 1024 // Convert to MB
		avgRAMPercent = totalRAMPercent / float64(ramCount)
	}

	return &SystemMetrics{
		CPUUsagePercent: avgCPU,
		RAMUsageMB:      avgRAM,
		RAMUsagePercent: avgRAMPercent,
		TotalRAMMB:      totalRAMMB,
	}, err
}
