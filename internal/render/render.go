// Package render draws familiars as procedural SVG, server-side.
//
// There are no image assets, sprite sheets, or icon fonts anywhere in this
// package. A render is a pure function of its Input: the same Input must
// produce byte-identical SVG forever, across pages, exports, and restarts.
package render

import (
	"time"

	"server-familiar/internal/model"
)

// Mood changes expression and LED colour only. It never changes species,
// size, scars, or palette.
type Mood string

const (
	MoodCalm     Mood = "calm"     // LEDs green, steady blink
	MoodFeverish Mood = "feverish" // LEDs amber to red, fan mouth fast
	MoodHungry   Mood = "hungry"   // LEDs pulsing
	MoodSleepy   Mood = "sleepy"   // half-lidded eyes
	MoodAsleep   Mood = "asleep"   // LEDs dim violet
	MoodStartled Mood = "startled" // eyes wide
)

// ScarKind is a canonical permanent mark. Once a scar is earned it renders for
// the rest of that life: scars are never faded or removed in code.
type ScarKind string

const (
	ScarCrack          ScarKind = "crack"
	ScarScorch         ScarKind = "scorch"
	ScarDent           ScarKind = "dent"
	ScarPatchPlate     ScarKind = "patch_plate"
	ScarScratchCluster ScarKind = "scratch_cluster"
)

// Scar is one permanent mark on a body.
type Scar struct {
	Kind     ScarKind
	Location string // body location token, e.g. chest, head, left_flank
	EarnedAt time.Time
}

// Trait is a per-familiar render flag, so showcase specifics stay data rather
// than name-matching branches in the drawing code: Anls's chest gauge pinned
// at 1 Mbps, Alpi's snow cap, Hermuse's X-stitched eye.
type Trait string

const (
	TraitPinnedGauge   Trait = "pinned_gauge"  // chest gauge stuck at its cap
	TraitSnowCap       Trait = "snow_cap"      // tiny snow cap on the antenna
	TraitNeighbours    Trait = "neighbours"    // lit co-tenant windows
	TraitStitchedEye   Trait = "stitched_eye"  // one eye X-stitched shut
	TraitHalfCRTHead   Trait = "half_crt_head" // half the head is a CRT
	TraitRouterChest   Trait = "router_chest"  // port lights beating like a heart
	TraitCableVeins    Trait = "cable_veins"   // cables glow like veins
	TraitCoolingFan    Trait = "cooling_fan"   // big glowing chest fan
	TraitSailFin       Trait = "sail_fin"      // sail / antenna fin on the head
	TraitCableBackpack Trait = "cable_backpack"
)

// Input is the complete render input tuple. Phase 1 fills in species, stage,
// palette seed, mood, and RAM; later phases supply scars, stitches, and
// traits without changing a single call site.
type Input struct {
	Species     model.Species
	Stage       model.Stage
	PaletteSeed int64
	Scars       []Scar
	Stitches    int
	Mood        Mood
	RAMBytes    int64 // body size scales with RAM, on a log-ish scale
	Traits      []Trait
	Size        int // viewport edge in px; 120 is card size
}

// TODO(phase-1, task 5.1): Creature(Input) string — v1 body plans for the six
// showcase species (Rackling, Tower Brute, Board Sprite, Cabinet Colony,
// Screenrunner, Patchwork), each with LED eyes and neck bolts.
// TODO(phase-1, task 5.2): size-from-RAM scaling and deterministic palettes.
// TODO(phase-1, task 5.3): golden tests for byte-identical output.
