---
name: diamond-rush-gameplay-fidelity
description: Reproduce, debug, or review Diamond Rush gameplay against the local Java ME original and original JAR. Use for stage layouts, raw IDs, object state machines, movement, gravity, enemies, hazards, tools, chests, scripts, bosses, checkpoints, input, save data, maps, progression, route solvability, source-tick timing, or fidelity claims in internal/original and internal/originalgame. Do not use as the primary workflow for sprite extraction or render-only discrepancies; use diamond-rush-asset-fidelity for those.
---

# Diamond Rush Gameplay Fidelity

Reconstruct behavior from the Java/JAR evidence chain, then verify the Go Runtime at the original 20 TPS cadence. Treat route tests as regression evidence, not as a substitute for source analysis.

## Start With Scope

1. Identify the world, data stage, displayed stage, coordinates, layer, raw IDs, and user-visible discrepancy.
2. Read the root `AGENTS.md`.
3. Open [references/reference-map.md](references/reference-map.md) and load only the task-relevant source/reference sections.
4. If the task is purely visual, switch to `$diamond-rush-asset-fidelity`. If behavior and visuals both change, finish the runtime state contract first.

Do not expand into another stage, world, Shop, or Siberia unless the requested behavior requires it.

## Establish Source Facts

1. Inspect the corresponding `decoded/worldN/stageNN.json`.
2. List all three layers' ID histograms, low-frequency IDs, and relevant coordinates. Keep authored data separate from initialization-time mutations.
3. Trace each relevant raw ID through `i.java`:
   - stage initialization
   - collision/passability
   - input or script trigger
   - object scan/update
   - damage and sound
   - checkpoint snapshot/restore
   - exit, map, and persistence transitions
4. Trace helper methods and packed state fields rather than inferring semantics from names.
5. When decompiled Java contains impossible control flow, inspect the original JAR with `javap -p -c`. Never use rebuilt `target/classes` as bytecode evidence.
6. Use FreeJ2ME only to resolve behavior that source/JAR inspection leaves ambiguous. Screenshots are not state-machine evidence.

Maintain a compact evidence table while working:

```text
layer/raw ID | coordinates | init | collision | update | timer/bits | result | Java/JAR anchor | trace
```

## Model The Runtime

Preserve Java's independent state:

- player, background, and foreground layers
- foreground high bits
- object state/direction bits
- motion timer and interpolation offset
- wake/dirty flags
- global source frame and special-stage state
- camera/script state
- checkpoint-restored data

Implement the smallest source-backed behavior needed for the task. Do not collapse these fields into a generic tile type, and do not add fallback mechanics merely to make a route pass.

Respect the original scan order and same-scan reprocessing behavior. Keep all gameplay, input gating, scripts, menus, results, damage, and persistence on source ticks. A 60 Hz display layer may interpolate snapshots but cannot mutate Runtime.

## Verify Behavior

For each relevant source tick, compare at least:

- hero tile, facing, action, animation frame, `jInt`, and pixel offset
- changed player/background/foreground cells
- relevant object state, timer, wake flag, direction, and interpolation
- enemy/hazard phase and collision result
- HP, lives, diamonds, keys, tools, hits, retries, checkpoint, mode, and progress
- animation, sound, Loading, result, and save transition start/end ticks

Use a real route at 20 TPS. Tests must not teleport, disable updates, edit layers, inject rewards, or directly modify boss health to claim completion.

Write only tests that protect confirmed source behavior, timing boundaries, restoration state, or a real regression. Prefer implementation plus one focused deterministic test over broad test scaffolding built from assumptions.

## Completion Gate

Before reporting completion:

1. Run focused package tests.
2. Run the real route test when stage logic changed.
3. Build `./cmd/originalrush` when runtime integration changed.
4. Update the relevant reference/spec when the confirmed contract changed.
5. State exactly which behaviors were source-audited, route-tested, visually checked, or still unverified.

Do not use “完全一致”, “一比一复原”, or “世界已完整复刻” without same-input Java/FreeJ2ME versus Go tick traces for all affected behavior.

## Commands

```bash
go test ./internal/original ./internal/originalgame
go build -o /tmp/originalrush-smoke ./cmd/originalrush
rm -f /tmp/originalrush-smoke
```

Use narrower `-run` filters while iterating. Run `go test ./...` only when the blast radius justifies it.
