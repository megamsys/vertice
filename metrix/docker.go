package metrix

import (
	"time"
)

type Stats struct {
	ContainerId     string
	Image           string
	MemoryUsage     uint64 //in bytes
	CPUUnitCost     string
	MemoryUnitCost  string
	AllocatedMemory int64
	AllocatedCpu    int64
	CPUStats        CPUStats //in percentage of total cpu used
	PreCPUStats     CPUStats
	NetworkIn       uint64
	NetworkOut      uint64
	AccountId       string
	AssemblyId      string
	QuotaId         string
	AssemblyName    string
	AssembliesId    string
	Status          string
	AuditPeriod     time.Time
}

type CPUStats struct {
	PercpuUsage       []uint64
	UsageInUsermode   uint64
	TotalUsage        uint64
	UsageInKernelmode uint64
	SystemCPUUsage    uint64
}
