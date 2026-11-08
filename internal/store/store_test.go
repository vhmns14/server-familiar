package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"server-familiar/internal/model"
)

// openTest opens a fresh database in a temporary directory.
func openTest(t *testing.T) *DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "familiar.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return db
}

// alpiHost is the smallest host in the showcase fleet: 0.5 core, 512 MB.
var alpiHost = model.HostProfile{
	Hostname:       "alpi",
	OS:             "alpine",
	Arch:           "amd64",
	Kernel:         "6.6.0",
	Virtualisation: "kvm",
	CPUCores:       1,
	RAMBytes:       512 << 20,
	DiskTotalBytes: 10 << 30,
	ServiceCount:   2,
}

func newAlpi(token string) NewFamiliar {
	return NewFamiliar{
		Name:        "Alpi",
		Species:     model.SpeciesBoardSprite,
		PaletteSeed: 20261105,
		Token:       token,
		Host:        alpiHost,
	}
}

func TestOpenAppliesMigrationsAndEnablesWAL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "familiar.db")

	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	version, err := db.SchemaVersion()
	if err != nil {
		t.Fatalf("SchemaVersion: %v", err)
	}
	if want := len(migrations); version != want {
		t.Errorf("schema version = %d, want %d", version, want)
	}

	var mode string
	if err := db.SQL().QueryRow(`PRAGMA journal_mode`).Scan(&mode); err != nil {
		t.Fatalf("read journal_mode: %v", err)
	}
	if mode != "wal" {
		t.Errorf("journal_mode = %q, want %q — backups and concurrent reads depend on it", mode, "wal")
	}

	var foreignKeys int
	if err := db.SQL().QueryRow(`PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil {
		t.Fatalf("read foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Errorf("foreign_keys = %d, want 1", foreignKeys)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Reopening an existing database must be a no-op, not a second apply.
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer reopened.Close()

	again, err := reopened.SchemaVersion()
	if err != nil {
		t.Fatalf("SchemaVersion after reopen: %v", err)
	}
	if again != version {
		t.Errorf("schema version after reopen = %d, want %d", again, version)
	}
}

func TestCreateFamiliarOpensFirstLife(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	born := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	nf := newAlpi("test-token-alpi")
	nf.BornAt = born

	created, err := db.CreateFamiliar(ctx, nf)
	if err != nil {
		t.Fatalf("CreateFamiliar: %v", err)
	}
	if created.ID == 0 {
		t.Error("created familiar has no id")
	}
	if created.Stage != model.StageI {
		t.Errorf("stage = %d, want %d (a new familiar is juvenile)", created.Stage, model.StageI)
	}
	if created.Status != StatusAlive {
		t.Errorf("status = %q, want %q", created.Status, StatusAlive)
	}

	loaded, err := db.FamiliarByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FamiliarByID: %v", err)
	}
	if loaded.Name != "Alpi" {
		t.Errorf("name = %q, want %q", loaded.Name, "Alpi")
	}
	if loaded.Species != model.SpeciesBoardSprite {
		t.Errorf("species = %q, want %q", loaded.Species, model.SpeciesBoardSprite)
	}
	if loaded.PaletteSeed != nf.PaletteSeed {
		t.Errorf("palette seed = %d, want %d — colours are decided once, at adoption",
			loaded.PaletteSeed, nf.PaletteSeed)
	}
	if !loaded.BornAt.Equal(born) {
		t.Errorf("born at = %s, want %s", loaded.BornAt, born)
	}
	if !loaded.LastSeenAt.IsZero() {
		t.Errorf("last seen = %s, want zero before the first sample", loaded.LastSeenAt)
	}
	if loaded.Stitches != 0 {
		t.Errorf("stitches = %d, want 0 — stitches are earned by dying", loaded.Stitches)
	}

	life, err := db.CurrentLife(ctx, created.ID)
	if err != nil {
		t.Fatalf("CurrentLife: %v", err)
	}
	if !life.IsCurrent() {
		t.Error("first life is not current")
	}
	if life.Host != alpiHost {
		t.Errorf("host profile = %+v, want %+v", life.Host, alpiHost)
	}
	if life.Species != model.SpeciesBoardSprite {
		t.Errorf("life species = %q, want %q", life.Species, model.SpeciesBoardSprite)
	}
}

func TestCreateFamiliarRejectsIncompleteInput(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	if _, err := db.CreateFamiliar(ctx, newAlpi("token-taken")); err != nil {
		t.Fatalf("seed Alpi: %v", err)
	}

	tests := []struct {
		name string
		in   NewFamiliar
	}{
		{
			name: "no name",
			in:   NewFamiliar{Species: model.SpeciesBoardSprite, Token: "t1"},
		},
		{
			name: "no species",
			in:   NewFamiliar{Name: "Nameless", Token: "t2"},
		},
		{
			name: "no token",
			in:   NewFamiliar{Name: "Tokenless", Species: model.SpeciesRackling},
		},
		{
			name: "duplicate name",
			in: NewFamiliar{
				Name: "Alpi", Species: model.SpeciesRackling, Token: "t3",
			},
		},
		{
			name: "duplicate token",
			in: NewFamiliar{
				Name: "Another", Species: model.SpeciesRackling, Token: "token-taken",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := db.CreateFamiliar(ctx, tc.in); err == nil {
				t.Fatal("CreateFamiliar succeeded, want an error")
			}
		})
	}
}

func TestTokenIsNeverStoredInTheClear(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	const token = "a-token-that-must-not-be-readable"
	created, err := db.CreateFamiliar(ctx, newAlpi(token))
	if err != nil {
		t.Fatalf("CreateFamiliar: %v", err)
	}

	var stored string
	if err := db.SQL().QueryRow(
		`SELECT token_hash FROM familiar WHERE id = ?`, created.ID,
	).Scan(&stored); err != nil {
		t.Fatalf("read token_hash: %v", err)
	}
	if stored == token {
		t.Fatal("token stored in the clear")
	}
	if stored != HashToken(token) {
		t.Errorf("token_hash = %q, want the SHA-256 of the token", stored)
	}

	found, err := db.FamiliarByToken(ctx, token)
	if err != nil {
		t.Fatalf("FamiliarByToken: %v", err)
	}
	if found.ID != created.ID {
		t.Errorf("found familiar %d, want %d", found.ID, created.ID)
	}
}

func TestLookupsMissNotFound(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	tests := []struct {
		name string
		call func() error
	}{
		{"familiar by id", func() error {
			_, err := db.FamiliarByID(ctx, 404)
			return err
		}},
		{"familiar by unknown token", func() error {
			_, err := db.FamiliarByToken(ctx, "nobody-has-this")
			return err
		}},
		{"familiar by empty token", func() error {
			_, err := db.FamiliarByToken(ctx, "")
			return err
		}},
		{"current life of nobody", func() error {
			_, err := db.CurrentLife(ctx, 404)
			return err
		}},
		{"touch nobody", func() error {
			return db.TouchLastSeen(ctx, 404, time.Now())
		}},
		{"set status of nobody", func() error {
			return db.SetStatus(ctx, 404, StatusSleeping)
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, ErrNotFound) {
				t.Errorf("error = %v, want ErrNotFound", err)
			}
		})
	}
}

func TestInsertSampleStoresVitalsAndTouchesLastSeen(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	created, err := db.CreateFamiliar(ctx, newAlpi("token-alpi"))
	if err != nil {
		t.Fatalf("CreateFamiliar: %v", err)
	}

	base := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	samples := []model.Sample{
		{TakenAt: base, CPUPercent: 11.5, MemUsedBytes: 200 << 20, MemTotalBytes: 512 << 20},
		{TakenAt: base.Add(10 * time.Second), CPUPercent: 42.25, MemUsedBytes: 260 << 20, MemTotalBytes: 512 << 20},
		{TakenAt: base.Add(20 * time.Second), CPUPercent: 93.75, MemUsedBytes: 300 << 20, MemTotalBytes: 512 << 20,
			DiskUsedBytes: 8 << 30, DiskTotalBytes: 10 << 30, UptimeSeconds: 86400,
			LoadAvg1: 1.5, NetRxBytes: 1 << 20, NetTxBytes: 2 << 20},
	}
	for i, s := range samples {
		if err := db.InsertSample(ctx, created.ID, s); err != nil {
			t.Fatalf("InsertSample %d: %v", i, err)
		}
	}

	reloaded, err := db.FamiliarByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FamiliarByID: %v", err)
	}
	wantSeen := samples[len(samples)-1].TakenAt
	if !reloaded.LastSeenAt.Equal(wantSeen) {
		t.Errorf("last seen = %s, want %s", reloaded.LastSeenAt, wantSeen)
	}

	recent, err := db.RecentSamples(ctx, created.ID, 10)
	if err != nil {
		t.Fatalf("RecentSamples: %v", err)
	}
	if len(recent) != len(samples) {
		t.Fatalf("got %d samples, want %d", len(recent), len(samples))
	}
	// Newest first.
	if !recent[0].TakenAt.Equal(wantSeen) {
		t.Errorf("newest sample taken at %s, want %s", recent[0].TakenAt, wantSeen)
	}
	newest := recent[0]
	if newest.CPUPercent != 93.75 {
		t.Errorf("cpu = %v, want 93.75", newest.CPUPercent)
	}
	if newest.UptimeSeconds != 86400 {
		t.Errorf("uptime = %d, want 86400", newest.UptimeSeconds)
	}
	if got := newest.DiskPercent(); got < 79 || got > 81 {
		t.Errorf("disk = %.1f%%, want about 80%%", got)
	}

	// A sample is always tied to the body that produced it.
	life, err := db.CurrentLife(ctx, created.ID)
	if err != nil {
		t.Fatalf("CurrentLife: %v", err)
	}
	var lifeID int64
	if err := db.SQL().QueryRow(
		`SELECT life_id FROM metric_sample ORDER BY taken_at DESC LIMIT 1`,
	).Scan(&lifeID); err != nil {
		t.Fatalf("read life_id: %v", err)
	}
	if lifeID != life.ID {
		t.Errorf("sample life_id = %d, want %d", lifeID, life.ID)
	}
}

func TestInsertSampleForUnknownFamiliar(t *testing.T) {
	db := openTest(t)
	err := db.InsertSample(context.Background(), 404, model.Sample{TakenAt: time.Now()})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestListFamiliarsWithLatestSample(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	base := time.Date(2026, 11, 5, 9, 0, 0, 0, time.UTC)

	anls, err := db.CreateFamiliar(ctx, NewFamiliar{
		Name:        "Anls",
		Species:     model.SpeciesRackling,
		PaletteSeed: 1,
		Token:       "token-anls",
		BornAt:      base,
		Host: model.HostProfile{
			Hostname: "anls", Virtualisation: "kvm", CPUCores: 2,
			RAMBytes: 2 << 30, DiskTotalBytes: 50 << 30,
		},
	})
	if err != nil {
		t.Fatalf("adopt Anls: %v", err)
	}

	// Adopted later, and never reports: its card must still appear.
	quiet, err := db.CreateFamiliar(ctx, NewFamiliar{
		Name:        "Vee",
		Species:     model.SpeciesScreenrunner,
		PaletteSeed: 2,
		Token:       "token-vee",
		BornAt:      base.Add(time.Hour),
		Host:        model.HostProfile{Hostname: "vee", RAMBytes: 4 << 30},
	})
	if err != nil {
		t.Fatalf("adopt Vee: %v", err)
	}

	for i, cpu := range []float64{10, 20, 77} {
		s := model.Sample{
			TakenAt:    base.Add(time.Duration(i) * 10 * time.Second),
			CPUPercent: cpu,
		}
		if err := db.InsertSample(ctx, anls.ID, s); err != nil {
			t.Fatalf("InsertSample %d: %v", i, err)
		}
	}

	list, err := db.ListFamiliarsWithLatestSample(ctx)
	if err != nil {
		t.Fatalf("ListFamiliarsWithLatestSample: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("got %d familiars, want 2", len(list))
	}

	// Ordered by adoption.
	if list[0].Familiar.Name != "Anls" || list[1].Familiar.Name != "Vee" {
		t.Errorf("order = %q, %q; want Anls, Vee",
			list[0].Familiar.Name, list[1].Familiar.Name)
	}

	withSample := list[0]
	if withSample.Latest == nil {
		t.Fatal("Anls has no latest sample")
	}
	if withSample.Latest.CPUPercent != 77 {
		t.Errorf("Anls cpu = %v, want 77 (the newest reading)", withSample.Latest.CPUPercent)
	}
	if withSample.Life.Host.RAMBytes != 2<<30 {
		t.Errorf("Anls RAM = %d, want %d — body size is drawn from it",
			withSample.Life.Host.RAMBytes, int64(2)<<30)
	}

	if list[1].Latest != nil {
		t.Errorf("Vee has a latest sample %+v, want nil — it never reported", list[1].Latest)
	}
	if list[1].Familiar.ID != quiet.ID {
		t.Errorf("second familiar id = %d, want %d", list[1].Familiar.ID, quiet.ID)
	}
}

func TestSetStatus(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	created, err := db.CreateFamiliar(ctx, newAlpi("token-alpi"))
	if err != nil {
		t.Fatalf("CreateFamiliar: %v", err)
	}

	tests := []struct {
		status  string
		wantErr bool
	}{
		{StatusAlive, false},
		{StatusSleeping, false},
		{StatusUnreachable, false},
		{StatusDead, false},
		{"feeling-fine", true},
		{"", true},
	}

	for _, tc := range tests {
		t.Run(tc.status, func(t *testing.T) {
			err := db.SetStatus(ctx, created.ID, tc.status)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("SetStatus(%q) succeeded, want an error", tc.status)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetStatus(%q): %v", tc.status, err)
			}
			reloaded, err := db.FamiliarByID(ctx, created.ID)
			if err != nil {
				t.Fatalf("FamiliarByID: %v", err)
			}
			if reloaded.Status != tc.status {
				t.Errorf("status = %q, want %q", reloaded.Status, tc.status)
			}
		})
	}
}

func TestAFamiliarCannotHaveTwoCurrentBodies(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	created, err := db.CreateFamiliar(ctx, newAlpi("token-alpi"))
	if err != nil {
		t.Fatalf("CreateFamiliar: %v", err)
	}

	// Reincarnation lands in Phase 4, but the schema must already refuse a
	// soul wearing two bodies at once.
	_, err = db.SQL().ExecContext(ctx, `
INSERT INTO life (familiar_id, species, started_at) VALUES (?, ?, ?)`,
		created.ID, string(model.SpeciesMinibox), toMillis(time.Now()))
	if err == nil {
		t.Fatal("opened a second current life, want a uniqueness error")
	}
}

func TestDeletingAFamiliarTakesItsLivesAndSamples(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	created, err := db.CreateFamiliar(ctx, newAlpi("token-alpi"))
	if err != nil {
		t.Fatalf("CreateFamiliar: %v", err)
	}
	if err := db.InsertSample(ctx, created.ID, model.Sample{TakenAt: time.Now()}); err != nil {
		t.Fatalf("InsertSample: %v", err)
	}

	if _, err := db.SQL().ExecContext(ctx, `DELETE FROM familiar WHERE id = ?`, created.ID); err != nil {
		t.Fatalf("delete familiar: %v", err)
	}

	for _, table := range []string{"life", "metric_sample"} {
		var n int
		if err := db.SQL().QueryRow(
			`SELECT COUNT(*) FROM `+table+` WHERE familiar_id = ?`, created.ID,
		).Scan(&n); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if n != 0 {
			t.Errorf("%s still holds %d rows for the deleted familiar", table, n)
		}
	}
}
