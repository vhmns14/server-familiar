package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"server-familiar/internal/model"
)

// CreateFamiliar adopts a familiar and opens its first life, in one
// transaction: a soul without a body is never a valid state.
func (db *DB) CreateFamiliar(ctx context.Context, nf NewFamiliar) (Familiar, error) {
	if nf.Name == "" {
		return Familiar{}, errors.New("store: familiar needs a name")
	}
	if nf.Species == "" {
		return Familiar{}, errors.New("store: familiar needs a species")
	}
	if nf.Token == "" {
		return Familiar{}, errors.New("store: familiar needs a token")
	}
	born := nf.BornAt
	if born.IsZero() {
		born = time.Now()
	}

	tx, err := db.sql.BeginTx(ctx, nil)
	if err != nil {
		return Familiar{}, fmt.Errorf("store: begin adoption: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
INSERT INTO familiar (name, species, stage, palette_seed, status, stitches, token_hash, born_at)
VALUES (?, ?, ?, ?, ?, 0, ?, ?)`,
		nf.Name, string(nf.Species), int(model.StageI), nf.PaletteSeed,
		StatusAlive, HashToken(nf.Token), toMillis(born),
	)
	if err != nil {
		return Familiar{}, fmt.Errorf("store: insert familiar %q: %w", nf.Name, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Familiar{}, fmt.Errorf("store: familiar id: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO life (
	familiar_id, species, started_at, host_hostname, host_os, host_arch,
	host_kernel, host_virt, host_cpu_cores, host_ram_bytes, host_disk_bytes,
	host_services
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, string(nf.Species), toMillis(born),
		nf.Host.Hostname, nf.Host.OS, nf.Host.Arch, nf.Host.Kernel,
		nf.Host.Virtualisation, nf.Host.CPUCores, nf.Host.RAMBytes,
		nf.Host.DiskTotalBytes, nf.Host.ServiceCount,
	); err != nil {
		return Familiar{}, fmt.Errorf("store: open first life for %q: %w", nf.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return Familiar{}, fmt.Errorf("store: commit adoption: %w", err)
	}

	return Familiar{
		ID:          id,
		Name:        nf.Name,
		Species:     nf.Species,
		Stage:       model.StageI,
		PaletteSeed: nf.PaletteSeed,
		Status:      StatusAlive,
		BornAt:      born.UTC(),
	}, nil
}

const familiarColumns = `
	f.id, f.name, f.species, f.stage, f.palette_seed, f.status, f.stitches,
	f.born_at, f.last_seen_at`

func scanFamiliar(row interface{ Scan(...any) error }) (Familiar, error) {
	var (
		f        Familiar
		species  string
		stage    int
		lastSeen sql.NullInt64
		bornAt   int64
	)
	err := row.Scan(&f.ID, &f.Name, &species, &stage, &f.PaletteSeed,
		&f.Status, &f.Stitches, &bornAt, &lastSeen)
	if err != nil {
		return Familiar{}, err
	}
	f.Species = model.Species(species)
	f.Stage = model.Stage(stage)
	f.BornAt = fromMillis(bornAt)
	f.LastSeenAt = timeFromNull(lastSeen)
	return f, nil
}

// FamiliarByID looks a familiar up by its id.
func (db *DB) FamiliarByID(ctx context.Context, id int64) (Familiar, error) {
	row := db.sql.QueryRowContext(ctx,
		`SELECT`+familiarColumns+` FROM familiar f WHERE f.id = ?`, id)
	f, err := scanFamiliar(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Familiar{}, ErrNotFound
	}
	if err != nil {
		return Familiar{}, fmt.Errorf("store: familiar %d: %w", id, err)
	}
	return f, nil
}

// FamiliarByToken looks a familiar up by its plaintext token. The token is
// hashed here; the caller never has to know how it is stored.
func (db *DB) FamiliarByToken(ctx context.Context, token string) (Familiar, error) {
	if token == "" {
		return Familiar{}, ErrNotFound
	}
	row := db.sql.QueryRowContext(ctx,
		`SELECT`+familiarColumns+` FROM familiar f WHERE f.token_hash = ?`,
		HashToken(token))
	f, err := scanFamiliar(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Familiar{}, ErrNotFound
	}
	if err != nil {
		return Familiar{}, fmt.Errorf("store: familiar by token: %w", err)
	}
	return f, nil
}

// CurrentLife returns the familiar's present body.
func (db *DB) CurrentLife(ctx context.Context, familiarID int64) (Life, error) {
	row := db.sql.QueryRowContext(ctx, `
SELECT id, familiar_id, species, started_at, ended_at, cause,
       host_hostname, host_os, host_arch, host_kernel, host_virt,
       host_cpu_cores, host_ram_bytes, host_disk_bytes, host_services
FROM life WHERE familiar_id = ? AND ended_at IS NULL`, familiarID)

	l, err := scanLife(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Life{}, ErrNotFound
	}
	if err != nil {
		return Life{}, fmt.Errorf("store: current life of familiar %d: %w", familiarID, err)
	}
	return l, nil
}

func scanLife(row interface{ Scan(...any) error }) (Life, error) {
	var (
		l       Life
		species string
		started int64
		ended   sql.NullInt64
	)
	err := row.Scan(&l.ID, &l.FamiliarID, &species, &started, &ended, &l.Cause,
		&l.Host.Hostname, &l.Host.OS, &l.Host.Arch, &l.Host.Kernel,
		&l.Host.Virtualisation, &l.Host.CPUCores, &l.Host.RAMBytes,
		&l.Host.DiskTotalBytes, &l.Host.ServiceCount)
	if err != nil {
		return Life{}, err
	}
	l.Species = model.Species(species)
	l.StartedAt = fromMillis(started)
	l.EndedAt = timeFromNull(ended)
	return l, nil
}

// InsertSample stores one reading against the familiar's current life and
// moves its last-seen timestamp forward.
func (db *DB) InsertSample(ctx context.Context, familiarID int64, s model.Sample) error {
	taken := s.TakenAt
	if taken.IsZero() {
		taken = time.Now()
	}

	tx, err := db.sql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: begin sample write: %w", err)
	}
	defer tx.Rollback()

	var lifeID int64
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM life WHERE familiar_id = ? AND ended_at IS NULL`,
		familiarID).Scan(&lifeID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("store: familiar %d has no current life: %w", familiarID, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("store: current life of familiar %d: %w", familiarID, err)
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO metric_sample (
	familiar_id, life_id, taken_at, cpu_percent, load_avg_1,
	mem_used_bytes, mem_total_bytes, disk_used_bytes, disk_total_bytes,
	uptime_seconds, net_rx_bytes, net_tx_bytes
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		familiarID, lifeID, toMillis(taken), s.CPUPercent, s.LoadAvg1,
		s.MemUsedBytes, s.MemTotalBytes, s.DiskUsedBytes, s.DiskTotalBytes,
		s.UptimeSeconds, s.NetRxBytes, s.NetTxBytes,
	); err != nil {
		return fmt.Errorf("store: insert sample for familiar %d: %w", familiarID, err)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE familiar SET last_seen_at = ? WHERE id = ?`,
		toMillis(taken), familiarID); err != nil {
		return fmt.Errorf("store: touch last-seen of familiar %d: %w", familiarID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit sample write: %w", err)
	}
	return nil
}

// TouchLastSeen records that the lab heard from a familiar, without storing a
// reading. Used when an agent checks in for work but has nothing new to say.
func (db *DB) TouchLastSeen(ctx context.Context, familiarID int64, at time.Time) error {
	if at.IsZero() {
		at = time.Now()
	}
	res, err := db.sql.ExecContext(ctx,
		`UPDATE familiar SET last_seen_at = ? WHERE id = ?`,
		toMillis(at), familiarID)
	if err != nil {
		return fmt.Errorf("store: touch last-seen of familiar %d: %w", familiarID, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// FamiliarWithVitals is a familiar, its current body, and its most recent
// reading — everything a creature card needs in one row.
type FamiliarWithVitals struct {
	Familiar Familiar
	Life     Life
	Latest   *model.Sample // nil when the familiar has never reported
}

// ListFamiliarsWithLatestSample returns every familiar with its current body
// and latest reading, oldest adoption first.
func (db *DB) ListFamiliarsWithLatestSample(ctx context.Context) ([]FamiliarWithVitals, error) {
	rows, err := db.sql.QueryContext(ctx, `
SELECT
	f.id, f.name, f.species, f.stage, f.palette_seed, f.status, f.stitches,
	f.born_at, f.last_seen_at,
	l.id, l.familiar_id, l.species, l.started_at, l.ended_at, l.cause,
	l.host_hostname, l.host_os, l.host_arch, l.host_kernel, l.host_virt,
	l.host_cpu_cores, l.host_ram_bytes, l.host_disk_bytes, l.host_services,
	s.taken_at, s.cpu_percent, s.load_avg_1, s.mem_used_bytes,
	s.mem_total_bytes, s.disk_used_bytes, s.disk_total_bytes,
	s.uptime_seconds, s.net_rx_bytes, s.net_tx_bytes
FROM familiar f
JOIN life l ON l.familiar_id = f.id AND l.ended_at IS NULL
LEFT JOIN metric_sample s ON s.id = (
	SELECT id FROM metric_sample
	WHERE familiar_id = f.id
	ORDER BY taken_at DESC, id DESC
	LIMIT 1
)
ORDER BY f.born_at, f.id`)
	if err != nil {
		return nil, fmt.Errorf("store: list familiars: %w", err)
	}
	defer rows.Close()

	var out []FamiliarWithVitals
	for rows.Next() {
		var (
			fv       FamiliarWithVitals
			fSpecies string
			fStage   int
			fBorn    int64
			fSeen    sql.NullInt64

			lSpecies string
			lStarted int64
			lEnded   sql.NullInt64

			sTaken     sql.NullInt64
			sCPU       sql.NullFloat64
			sLoad      sql.NullFloat64
			sMemUsed   sql.NullInt64
			sMemTotal  sql.NullInt64
			sDiskUsed  sql.NullInt64
			sDiskTotal sql.NullInt64
			sUptime    sql.NullInt64
			sRx        sql.NullInt64
			sTx        sql.NullInt64
		)
		err := rows.Scan(
			&fv.Familiar.ID, &fv.Familiar.Name, &fSpecies, &fStage,
			&fv.Familiar.PaletteSeed, &fv.Familiar.Status, &fv.Familiar.Stitches,
			&fBorn, &fSeen,
			&fv.Life.ID, &fv.Life.FamiliarID, &lSpecies, &lStarted, &lEnded,
			&fv.Life.Cause,
			&fv.Life.Host.Hostname, &fv.Life.Host.OS, &fv.Life.Host.Arch,
			&fv.Life.Host.Kernel, &fv.Life.Host.Virtualisation,
			&fv.Life.Host.CPUCores, &fv.Life.Host.RAMBytes,
			&fv.Life.Host.DiskTotalBytes, &fv.Life.Host.ServiceCount,
			&sTaken, &sCPU, &sLoad, &sMemUsed, &sMemTotal, &sDiskUsed,
			&sDiskTotal, &sUptime, &sRx, &sTx,
		)
		if err != nil {
			return nil, fmt.Errorf("store: scan familiar row: %w", err)
		}

		fv.Familiar.Species = model.Species(fSpecies)
		fv.Familiar.Stage = model.Stage(fStage)
		fv.Familiar.BornAt = fromMillis(fBorn)
		fv.Familiar.LastSeenAt = timeFromNull(fSeen)

		fv.Life.Species = model.Species(lSpecies)
		fv.Life.StartedAt = fromMillis(lStarted)
		fv.Life.EndedAt = timeFromNull(lEnded)

		if sTaken.Valid {
			fv.Latest = &model.Sample{
				TakenAt:        fromMillis(sTaken.Int64),
				CPUPercent:     sCPU.Float64,
				LoadAvg1:       sLoad.Float64,
				MemUsedBytes:   sMemUsed.Int64,
				MemTotalBytes:  sMemTotal.Int64,
				DiskUsedBytes:  sDiskUsed.Int64,
				DiskTotalBytes: sDiskTotal.Int64,
				UptimeSeconds:  sUptime.Int64,
				NetRxBytes:     sRx.Int64,
				NetTxBytes:     sTx.Int64,
			}
		}
		out = append(out, fv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list familiars: %w", err)
	}
	return out, nil
}

// RecentSamples returns up to limit of a familiar's readings, newest first.
func (db *DB) RecentSamples(ctx context.Context, familiarID int64, limit int) ([]model.Sample, error) {
	if limit <= 0 {
		limit = 60
	}
	rows, err := db.sql.QueryContext(ctx, `
SELECT taken_at, cpu_percent, load_avg_1, mem_used_bytes, mem_total_bytes,
       disk_used_bytes, disk_total_bytes, uptime_seconds, net_rx_bytes, net_tx_bytes
FROM metric_sample
WHERE familiar_id = ?
ORDER BY taken_at DESC, id DESC
LIMIT ?`, familiarID, limit)
	if err != nil {
		return nil, fmt.Errorf("store: recent samples of familiar %d: %w", familiarID, err)
	}
	defer rows.Close()

	var out []model.Sample
	for rows.Next() {
		var (
			s     model.Sample
			taken int64
		)
		if err := rows.Scan(&taken, &s.CPUPercent, &s.LoadAvg1,
			&s.MemUsedBytes, &s.MemTotalBytes, &s.DiskUsedBytes,
			&s.DiskTotalBytes, &s.UptimeSeconds, &s.NetRxBytes,
			&s.NetTxBytes); err != nil {
			return nil, fmt.Errorf("store: scan sample: %w", err)
		}
		s.TakenAt = fromMillis(taken)
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: recent samples of familiar %d: %w", familiarID, err)
	}
	return out, nil
}

// SetStatus moves a familiar between alive, sleeping, unreachable, and dead.
func (db *DB) SetStatus(ctx context.Context, familiarID int64, status string) error {
	switch status {
	case StatusAlive, StatusSleeping, StatusUnreachable, StatusDead:
	default:
		return fmt.Errorf("store: unknown status %q", status)
	}
	res, err := db.sql.ExecContext(ctx,
		`UPDATE familiar SET status = ? WHERE id = ?`, status, familiarID)
	if err != nil {
		return fmt.Errorf("store: set status of familiar %d: %w", familiarID, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
