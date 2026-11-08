package store

import (
	"fmt"
	"time"
)

// migration is one forward step in the schema. Migrations are append-only:
// once a version has shipped, its statements are never edited.
type migration struct {
	version int
	name    string
	stmt    string
}

// migrations are applied in order at lab startup.
var migrations = []migration{
	{
		version: 1,
		name:    "familiars, lives, and metric samples",
		stmt: `
CREATE TABLE familiar (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	name          TEXT    NOT NULL UNIQUE,
	species       TEXT    NOT NULL,
	stage         INTEGER NOT NULL DEFAULT 1,
	palette_seed  INTEGER NOT NULL,
	status        TEXT    NOT NULL DEFAULT 'alive',
	stitches      INTEGER NOT NULL DEFAULT 0,
	token_hash    TEXT    NOT NULL UNIQUE,
	born_at       INTEGER NOT NULL,
	last_seen_at  INTEGER
) STRICT;

-- One row per body. The open life (ended_at IS NULL) is the current one;
-- closed lives are the familiar's past lives.
CREATE TABLE life (
	id               INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	familiar_id      INTEGER NOT NULL REFERENCES familiar(id) ON DELETE CASCADE,
	species          TEXT    NOT NULL,
	started_at       INTEGER NOT NULL,
	ended_at         INTEGER,
	cause            TEXT    NOT NULL DEFAULT '',
	host_hostname    TEXT    NOT NULL DEFAULT '',
	host_os          TEXT    NOT NULL DEFAULT '',
	host_arch        TEXT    NOT NULL DEFAULT '',
	host_kernel      TEXT    NOT NULL DEFAULT '',
	host_virt        TEXT    NOT NULL DEFAULT '',
	host_cpu_cores   INTEGER NOT NULL DEFAULT 0,
	host_ram_bytes   INTEGER NOT NULL DEFAULT 0,
	host_disk_bytes  INTEGER NOT NULL DEFAULT 0,
	host_services    INTEGER NOT NULL DEFAULT 0
) STRICT;

-- A familiar has at most one open life at a time.
CREATE UNIQUE INDEX idx_life_one_current
	ON life(familiar_id) WHERE ended_at IS NULL;
CREATE INDEX idx_life_familiar ON life(familiar_id, started_at);

CREATE TABLE metric_sample (
	id                INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	familiar_id       INTEGER NOT NULL REFERENCES familiar(id) ON DELETE CASCADE,
	life_id           INTEGER NOT NULL REFERENCES life(id) ON DELETE CASCADE,
	taken_at          INTEGER NOT NULL,
	cpu_percent       REAL    NOT NULL DEFAULT 0,
	load_avg_1        REAL    NOT NULL DEFAULT 0,
	mem_used_bytes    INTEGER NOT NULL DEFAULT 0,
	mem_total_bytes   INTEGER NOT NULL DEFAULT 0,
	disk_used_bytes   INTEGER NOT NULL DEFAULT 0,
	disk_total_bytes  INTEGER NOT NULL DEFAULT 0,
	uptime_seconds    INTEGER NOT NULL DEFAULT 0,
	net_rx_bytes      INTEGER NOT NULL DEFAULT 0,
	net_tx_bytes      INTEGER NOT NULL DEFAULT 0
) STRICT;

-- "The latest sample for this familiar" is the lab's hottest query.
CREATE INDEX idx_sample_familiar_taken
	ON metric_sample(familiar_id, taken_at DESC);
`,
	},
}

// migrate applies every migration the database has not seen yet.
func (db *DB) migrate() error {
	_, err := db.sql.Exec(`
CREATE TABLE IF NOT EXISTS schema_migration (
	version     INTEGER NOT NULL PRIMARY KEY,
	name        TEXT    NOT NULL,
	applied_at  INTEGER NOT NULL
) STRICT;`)
	if err != nil {
		return fmt.Errorf("store: create migration table: %w", err)
	}

	var current int
	if err := db.sql.QueryRow(
		`SELECT COALESCE(MAX(version), 0) FROM schema_migration`,
	).Scan(&current); err != nil {
		return fmt.Errorf("store: read schema version: %w", err)
	}

	for _, m := range migrations {
		if m.version <= current {
			continue
		}
		tx, err := db.sql.Begin()
		if err != nil {
			return fmt.Errorf("store: begin migration %d: %w", m.version, err)
		}
		if _, err := tx.Exec(m.stmt); err != nil {
			tx.Rollback()
			return fmt.Errorf("store: apply migration %d (%s): %w", m.version, m.name, err)
		}
		if _, err := tx.Exec(
			`INSERT INTO schema_migration (version, name, applied_at) VALUES (?, ?, ?)`,
			m.version, m.name, toMillis(time.Now()),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("store: record migration %d: %w", m.version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("store: commit migration %d: %w", m.version, err)
		}
	}
	return nil
}

// SchemaVersion reports the highest migration applied.
func (db *DB) SchemaVersion() (int, error) {
	var v int
	err := db.sql.QueryRow(
		`SELECT COALESCE(MAX(version), 0) FROM schema_migration`,
	).Scan(&v)
	if err != nil {
		return 0, fmt.Errorf("store: read schema version: %w", err)
	}
	return v, nil
}
