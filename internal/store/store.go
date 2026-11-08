// Package store is the lab's persistence layer: one SQLite file in WAL mode.
//
// Backing up the lab must stay "copy one file", so nothing here may reach for
// an external database, cache, or queue.
package store

import (
	"time"

	"server-familiar/internal/model"
)

// Familiar is the soul: one row that persists across bodies. Its current body
// is the open Life row that points at it.
type Familiar struct {
	ID          int64
	Name        string
	Species     model.Species
	Stage       model.Stage
	PaletteSeed int64
	Status      string // alive, sleeping, unreachable, dead
	Stitches    int
	BornAt      time.Time
	LastSeenAt  time.Time
	Host        model.HostProfile
}

// Life is one incarnation: a body that begins and ends. A migration ends one
// Life and opens the next against the same Familiar.
type Life struct {
	ID         int64
	FamiliarID int64
	Species    model.Species
	Host       model.HostProfile
	StartedAt  time.Time
	EndedAt    *time.Time // nil while this body is the current one
	Cause      string     // why this life ended: migration, trial expired, ...
}

// TODO(phase-1, task 2.1): migration runner and schema v1 (familiar, life,
// metric_sample) applied at lab startup, WAL mode on.
// TODO(phase-1, task 2.2): CreateFamiliarWithLife, InsertSample,
// ListFamiliarsWithLatestSample, TouchLastSeen.
//
// Blocked on one decision: CGO is off, so the standard library cannot talk to
// SQLite. The driver is the single third-party dependency this build expects
// (see .kiro/specs/phase-1-skeleton/design.md) and needs the keeper's approval
// before it is added.
