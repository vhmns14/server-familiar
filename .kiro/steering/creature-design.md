# Creature Design — Server Familiar

*Steering file: always loaded. These rules govern every creature
render, in the lab UI, the terrarium, certificates, and the Life Book.*

## The one rule

**The body IS the hardware.** A familiar is not a generic blob wearing
accessories. Its chassis, limbs, face, and proportions are built from
the class of machine it lives on: drive bays become teeth, status
LEDs become eyes, a cooling fan becomes a mouth, network cables become
tails and tentacles, antennas become ears. Read the machine first,
then draw the monster that machine would be.

## Species and hardware mapping

| Species | Hardware | Silhouette cues |
| --- | --- | --- |
| **Rackling** | Rack server | Wide flat 2U chassis body, LED row eyes, vent grille mouth |
| **Tower Brute** | Tower PC | Tall chunky tower body, one big glowing chest fan, fin on the head |
| **Board Sprite** | SBC / tiny VPS | Small green PCB body, circuit traces, one oversized eye, stick limbs, big antenna |
| **Cabinet Colony** | Shared rack / LXC | Very wide cabinet body tiled with small lit "neighbour" windows, short stubby arms |
| **Screenrunner** | Worker VM | Monitor for a head (the screen is the face), running legs, backpack of tangled cables |
| **Patchwork** | Multi-service stitched host | Mismatched panels from other species, seams everywhere, mixed eye types, cable veins |
| **Blade Imp** | Blade server | Long thin blade body, a single row of many small LED eyes |
| **Vaultback** | NAS / storage | Humped back and belly stacked with drive bays |
| **Routling** | Router / switch | Low wide body, a crown of antennas, port lights for a mouth |
| **Crunchbeast** | GPU rig | Open frame body covered in fans, heat shimmer / steam wisps |
| **Oldframe** | Mainframe | Massive slab body, tape-reel eyes, slow and monumental proportions |
| **Lapdog** | Laptop | Four-legged body, screen lid as head/ears, loyal dog posture |
| **Cloudpuff** | Cloud instance | No chassis at all — a soft hovering cloud body with LED eyes, faint outline only |
| **Minibox** | Mini-PC / NUC | Small neat cube body, single LED eye band, tidy and compact |

## Shared design language

- **LED eyes** on every species: colour encodes state (see mood rules).
- **Fan mouths** where the hardware has fans; fan blades rotate on
  hover in the web UI (CSS animation on the SVG group is fine).
- **Cable tails / tentacles** on anything networked; a Patchwork's
  cables glow like veins.
- **Antenna ears** on anything remote or wireless.
- **Frankenstein stitches** drawn as thick dashed seams across panel
  joins — and across resurrection points (see below).
- **Neck bolts** on every creature with a head. Non-negotiable; it is
  the Frankenstein signature.
- **Patchwork panels:** a panel in a second palette colour is allowed
  on any species after a migration (a part replaced in a past life),
  and is the default construction for the Patchwork species.
- Keep silhouettes readable at card size (~120 px). Detail serves the
  silhouette, never the reverse.

## Deterministic rendering

A creature render is a pure function of:

```text
render(species, stage, palette_seed, scars, stitches, mood, size_from_ram)
```

- **Species** picks the body plan. **Stage** (I / II / III) scales
  proportions and adds elder/juvenile features (smaller + bigger eyes
  when young; antenna grey, extra vents, slower blink when elder).
- **Palette seed** (stored on the familiar at adoption) picks its LED
  palette and panel tint deterministically. Same seed, same colours,
  forever — across pages, exports, and restarts.
- **Body size scales with RAM**: 512 MB is visibly tiny next to 16 GB.
  Use a log-ish scale so the range stays drawable.
- **Scars** are a stored list (kind + body location + when earned).
  Each scar kind has one canonical mark (crack, scorch, dent, patch
  plate, scratch cluster). **Once earned, a scar renders permanently**
  for the rest of that life — never fade, never remove scars in code.
- **Stitches** accumulate: stitch count is stored per familiar; each
  resurrection adds one seam at a deterministic location. Hermuse's
  record is 23 — the renderer must stay tidy at high counts (cycle
  seam placement along panel joins).
- **Mood changes expression and LED colour only** — eye shape
  (half-lidded when sleepy, wide when startled), mouth/fan speed, LED
  hue (calm green, fever amber/red, asleep dim violet, hungry pulsing).
  Mood never changes species, size, scars, or palette.

## Showcase consistency

The six demo familiars (see `product.md`) must render recognisably as
described there: Anls with his chest gauge pinned at **1 Mbps**; Alpi
with a tiny snow cap on the antenna; No. 1021 tiled with lit neighbour
windows; Hermuse half-CRT-headed with one LED eye and one X-stitched
eye and a router chest whose port lights beat like a heart. These
specifics come from per-familiar feature flags in the render inputs,
not from special-cased drawing code per name — add flags, not `if
name == "Hermuse"` branches.
