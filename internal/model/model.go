// Package model holds the types shared between the agent and the lab.
//
// Nothing in this package may import another internal package: the agent
// binary depends on it, and the agent must stay small enough to run on a
// 512 MB Alpine host.
package model

import "time"

// Species is a creature body plan. Each species maps to a class of real
// hardware; see .kiro/steering/creature-design.md, where the mapping table
// is law.
type Species string

const (
	SpeciesRackling      Species = "rackling"       // rack server
	SpeciesTowerBrute    Species = "tower_brute"    // tower PC
	SpeciesBoardSprite   Species = "board_sprite"   // SBC / tiny VPS
	SpeciesCabinetColony Species = "cabinet_colony" // shared rack / LXC host
	SpeciesScreenrunner  Species = "screenrunner"   // worker VM
	SpeciesPatchwork     Species = "patchwork"      // multi-service stitched host
	SpeciesBladeImp      Species = "blade_imp"      // blade server
	SpeciesVaultback     Species = "vaultback"      // NAS / storage array
	SpeciesRoutling      Species = "routling"       // router / switch
	SpeciesCrunchbeast   Species = "crunchbeast"    // GPU rig
	SpeciesOldframe      Species = "oldframe"       // mainframe
	SpeciesLapdog        Species = "lapdog"         // laptop
	SpeciesCloudpuff     Species = "cloudpuff"      // cloud instance, no chassis
	SpeciesMinibox       Species = "minibox"        // mini-PC / NUC
)

// AllSpecies lists every species in the order of the design table.
func AllSpecies() []Species {
	return []Species{
		SpeciesRackling, SpeciesTowerBrute, SpeciesBoardSprite,
		SpeciesCabinetColony, SpeciesScreenrunner, SpeciesPatchwork,
		SpeciesBladeImp, SpeciesVaultback, SpeciesRoutling,
		SpeciesCrunchbeast, SpeciesOldframe, SpeciesLapdog,
		SpeciesCloudpuff, SpeciesMinibox,
	}
}

// Stage is a growth stage. Stage changes proportions and juvenile/elder
// features only; it never changes species, scars, or palette.
type Stage int

const (
	StageI   Stage = 1 // juvenile: smaller body, bigger eyes
	StageII  Stage = 2 // adult
	StageIII Stage = 3 // elder: grey antenna, extra vents, slower blink
)

// HostProfile describes the machine a familiar lives on. It is collected
// once at enrolment and refreshed when the agent notices a change, because a
// RAM or disk upgrade reshapes the body (Phase 4: evolution).
type HostProfile struct {
	Hostname       string `json:"hostname"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Kernel         string `json:"kernel"`
	Virtualisation string `json:"virtualisation"` // kvm, lxc, docker, none, ...
	CPUCores       int    `json:"cpu_cores"`
	RAMBytes       int64  `json:"ram_bytes"`
	DiskTotalBytes int64  `json:"disk_total_bytes"`
	ServiceCount   int    `json:"service_count"` // stitched services: feeds Patchwork
}

// Sample is one reading of a host's vitals, taken every 10 seconds.
type Sample struct {
	TakenAt        time.Time `json:"taken_at"`
	CPUPercent     float64   `json:"cpu_percent"`
	LoadAvg1       float64   `json:"load_avg_1"`
	MemUsedBytes   int64     `json:"mem_used_bytes"`
	MemTotalBytes  int64     `json:"mem_total_bytes"`
	DiskUsedBytes  int64     `json:"disk_used_bytes"`
	DiskTotalBytes int64     `json:"disk_total_bytes"`
	UptimeSeconds  int64     `json:"uptime_seconds"`
	NetRxBytes     int64     `json:"net_rx_bytes"` // cumulative since boot
	NetTxBytes     int64     `json:"net_tx_bytes"` // cumulative since boot
}

// MemPercent reports memory use as a percentage, or 0 when unknown.
func (s Sample) MemPercent() float64 {
	if s.MemTotalBytes <= 0 {
		return 0
	}
	return float64(s.MemUsedBytes) / float64(s.MemTotalBytes) * 100
}

// DiskPercent reports disk use as a percentage, or 0 when unknown.
func (s Sample) DiskPercent() float64 {
	if s.DiskTotalBytes <= 0 {
		return 0
	}
	return float64(s.DiskUsedBytes) / float64(s.DiskTotalBytes) * 100
}
