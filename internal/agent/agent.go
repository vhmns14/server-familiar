// Package agent samples host vitals from /proc and reports them to the lab.
//
// It reads /proc directly and never shells out: the agent has to run on hosts
// where no tool but the kernel can be assumed.
package agent

import (
	"errors"
	"os"
	"time"
)

// SampleInterval is how often the agent reads the host's vitals.
const SampleInterval = 10 * time.Second

// Config is the agent's complete configuration. It comes from environment
// variables and flags only; secrets are never read from a file in the repo.
type Config struct {
	LabURL   string        // AGENT_LAB_URL, e.g. https://lab.example.net
	Token    string        // AGENT_TOKEN: per-familiar token, or fleet token before enrolment
	Interval time.Duration // sampling interval; defaults to SampleInterval
	ProcRoot string        // AGENT_PROC_ROOT: overridable for tests; defaults to /proc
	DiskPath string        // AGENT_DISK_PATH: filesystem to measure; defaults to /
}

// ConfigFromEnv builds a Config from the environment, applying defaults.
func ConfigFromEnv() Config {
	cfg := Config{
		LabURL:   os.Getenv("AGENT_LAB_URL"),
		Token:    os.Getenv("AGENT_TOKEN"),
		Interval: SampleInterval,
		ProcRoot: os.Getenv("AGENT_PROC_ROOT"),
		DiskPath: os.Getenv("AGENT_DISK_PATH"),
	}
	if cfg.ProcRoot == "" {
		cfg.ProcRoot = "/proc"
	}
	if cfg.DiskPath == "" {
		cfg.DiskPath = "/"
	}
	return cfg
}

// Validate reports whether the agent has enough configuration to run.
func (c Config) Validate() error {
	if c.LabURL == "" {
		return errors.New("AGENT_LAB_URL is not set")
	}
	if c.Token == "" {
		return errors.New("AGENT_TOKEN is not set")
	}
	return nil
}

// TODO(phase-1, task 3.1): /proc parsers for CPU, memory, disk, uptime, net.
// TODO(phase-1, task 3.2): the sampling loop and the ingest client with
// exponential backoff (cap 5 min) over a bounded drop-oldest retry buffer.
