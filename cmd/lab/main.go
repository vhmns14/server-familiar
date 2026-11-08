// Command lab is the central Server Familiar laboratory: it stores the fleet's
// telemetry, derives each familiar's body and life from it, and serves the UI,
// the REST API, and the live vitals stream.
//
// Configuration is environment variables and flags only:
//
//	LAB_ADDR         listen address (default :8787)
//	LAB_DB           SQLite file path (default ./familiar.db)
//	LAB_DEMO         1 to run the simulated six-familiar fleet
//	LAB_FLEET_TOKEN  token an agent presents to enrol (required unless demo)
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"server-familiar/internal/store"
)

type config struct {
	addr string
	db   string
	demo bool
}

func loadConfig() config {
	cfg := config{
		addr: envOr("LAB_ADDR", ":8787"),
		db:   envOr("LAB_DB", "familiar.db"),
		demo: envBool("LAB_DEMO"),
	}
	flag.StringVar(&cfg.addr, "addr", cfg.addr, "listen address")
	flag.StringVar(&cfg.db, "db", cfg.db, "path to the SQLite database file")
	flag.BoolVar(&cfg.demo, "demo", cfg.demo, "run with the simulated six-familiar fleet")
	flag.Parse()
	return cfg
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string) bool {
	v, err := strconv.ParseBool(os.Getenv(key))
	return err == nil && v
}

func main() {
	cfg := loadConfig()
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("lab: ")

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config) error {
	// Opening the store applies any pending migrations, so the lab is always
	// started against an up-to-date schema.
	db, err := store.Open(cfg.db)
	if err != nil {
		return err
	}
	defer db.Close()

	version, err := db.SchemaVersion()
	if err != nil {
		return err
	}
	log.Printf("database %s ready at schema v%d", cfg.db, version)

	// TODO(phase-1): mount the ingest and enrolment handlers (task 4), the web
	// pages and SSE hub (task 6), and the demo simulator when cfg.demo is set
	// (task 7).
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errs := make(chan error, 1)
	go func() {
		log.Printf("listening on %s (db=%s demo=%t)", cfg.addr, cfg.db, cfg.demo)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		log.Print("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
