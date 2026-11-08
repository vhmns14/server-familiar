// Command agent runs on each monitored server. It samples the host's vitals
// from /proc every ten seconds and posts them to the lab.
//
// It is a single static binary with no runtime dependencies, built with CGO
// off so it runs on musl hosts, and it only ever makes outbound connections.
//
// Configuration is environment variables and flags only:
//
//	AGENT_LAB_URL    the lab's base URL, e.g. https://lab.example.net
//	AGENT_TOKEN      this familiar's token (the fleet token before enrolment)
//	AGENT_PROC_ROOT  override /proc, for tests
//	AGENT_DISK_PATH  the filesystem to measure (default /)
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"server-familiar/internal/agent"
)

func main() {
	cfg := agent.ConfigFromEnv()
	flag.StringVar(&cfg.LabURL, "lab-url", cfg.LabURL, "the lab's base URL")
	flag.DurationVar(&cfg.Interval, "interval", cfg.Interval, "sampling interval")
	flag.StringVar(&cfg.DiskPath, "disk-path", cfg.DiskPath, "filesystem to measure")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("agent: ")

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// TODO(phase-1, tasks 3.1-3.2): sample /proc on cfg.Interval and post each
	// sample to cfg.LabURL, retrying with backoff over a bounded buffer.
	log.Printf("configured for %s, sampling every %s", cfg.LabURL, cfg.Interval)
	<-ctx.Done()
	log.Print("stopping")
}
