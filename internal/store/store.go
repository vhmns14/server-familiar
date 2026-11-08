// Package store is the lab's persistence layer: one SQLite file in WAL mode.
//
// Backing up the lab must stay "copy one file", so nothing here may reach for
// an external database, cache, or queue.
//
// Two conventions hold throughout:
//
// Timestamps are stored as INTEGER Unix milliseconds in UTC, never as text.
// Text timestamps sort wrongly the moment fractional seconds vary in width,
// and "give me the latest sample" has to be exact. Read one back in a shell
// with datetime(taken_at / 1000, 'unixepoch').
//
// A host profile belongs to a life, not to a familiar. A familiar is the
// soul; a body is a machine. When a server migrates, the familiar keeps its
// name, palette, scars, and stitches while a new life records the new host.
package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	// The lab's one third-party dependency: CGO is off, so the standard
	// library cannot speak SQLite. Only cmd/lab reaches this package, so the
	// agent binary stays dependency-free.
	_ "modernc.org/sqlite"

	"server-familiar/internal/model"
)

// ErrNotFound is returned when a lookup matches no row.
var ErrNotFound = errors.New("store: not found")

// Status values a familiar can hold. Phase 1 uses alive, sleeping, and
// unreachable; dead arrives with the graveyard.
const (
	StatusAlive       = "alive"
	StatusSleeping    = "sleeping"
	StatusUnreachable = "unreachable"
	StatusDead        = "dead"
)

// DB is an open lab database.
type DB struct {
	sql *sql.DB
}

// Open opens or creates the lab database at path and brings its schema up to
// date. The returned DB is safe for concurrent use.
func Open(path string) (*DB, error) {
	// WAL so a reader never blocks the writer; busy_timeout so a concurrent
	// writer waits instead of failing; foreign_keys so a deleted familiar
	// takes its lives and samples with it; synchronous NORMAL because losing
	// the last few seconds of telemetry in a crash costs nothing.
	dsn := fmt.Sprintf(
		"file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"+
			"&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)",
		path,
	)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", path, err)
	}
	// One writer keeps SQLite out of lock contention entirely. A fleet of a
	// few dozen agents sampling every ten seconds is nowhere near the limit
	// where this would matter.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(0)

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("store: ping %s: %w", path, err)
	}
	db := &DB{sql: sqlDB}
	if err := db.migrate(); err != nil {
		sqlDB.Close()
		return nil, err
	}
	return db, nil
}

// Close closes the database.
func (db *DB) Close() error { return db.sql.Close() }

// SQL exposes the underlying handle for packages that need it. Prefer the
// typed methods on DB; this exists so later phases can add queries without
// widening this package prematurely.
func (db *DB) SQL() *sql.DB { return db.sql }

// HashToken returns the hex SHA-256 of a token. Tokens are never stored in
// the clear: the lab keeps only the hash and compares against it.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func toMillis(t time.Time) int64 { return t.UTC().UnixMilli() }

func fromMillis(ms int64) time.Time { return time.UnixMilli(ms).UTC() }

func nullMillis(t time.Time) sql.NullInt64 {
	if t.IsZero() {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: toMillis(t), Valid: true}
}

func timeFromNull(n sql.NullInt64) time.Time {
	if !n.Valid {
		return time.Time{}
	}
	return fromMillis(n.Int64)
}

// Familiar is the soul: one row that persists across bodies.
type Familiar struct {
	ID          int64
	Name        string
	Species     model.Species
	Stage       model.Stage
	PaletteSeed int64
	Status      string
	Stitches    int
	BornAt      time.Time
	LastSeenAt  time.Time // zero until the first sample arrives
}

// Life is one incarnation: a body that begins and ends. A migration ends one
// life and opens the next against the same familiar.
type Life struct {
	ID         int64
	FamiliarID int64
	Species    model.Species
	Host       model.HostProfile
	StartedAt  time.Time
	EndedAt    time.Time // zero while this body is the current one
	Cause      string    // why this life ended: migration, trial expired, ...
}

// IsCurrent reports whether this life is the familiar's present body.
func (l Life) IsCurrent() bool { return l.EndedAt.IsZero() }

// NewFamiliar is everything needed to adopt a familiar and open its first
// life. The caller supplies the palette seed so colours are decided once, at
// adoption, and never drift afterwards.
type NewFamiliar struct {
	Name        string
	Species     model.Species
	PaletteSeed int64
	Token       string // stored only as a hash
	Host        model.HostProfile
	BornAt      time.Time // defaults to now
}
