// Package demo runs the lab against a simulated six-familiar fleet, so the
// product can be evaluated with no servers, no accounts, and no setup.
//
// Demo familiars are real rows written through the real store path, which is
// why every later phase works on them unchanged. A feature that cannot be
// shown in demo mode does not count as done.
package demo

// FleetSize is the number of showcase familiars demo mode seeds: Anls, Sail,
// Alpi, No. 1021, Vee, and Hermuse.
const FleetSize = 6

// TODO(phase-1, task 7.1): seed the six showcase familiars with the exact
// facts from .kiro/steering/product.md — Anls's egress pinned at ~1 Mbps,
// Alpi at 0.5 core / 512 MB, No. 1021's six co-tenants, Hermuse's 23
// stitches.
// TODO(phase-1, task 7.2): a telemetry simulator per familiar with plausible
// diurnal wobble, written through the real store functions.
// TODO(phase-1, task 7.3): LAB_DEMO=1 wiring and a visible DEMO MODE banner.
