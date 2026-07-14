---
name: diamond-rush-asset-fidelity
description: Diagnose, implement, or review Diamond Rush assets and rendering against the local Java ME source and original .f resources. Use for sprite extraction, resource-slot or chunk mapping, frame/module composition, animation selection and timing, anchors, flips, alpha, atlas padding, multi-tile art, foreground ordering, camera/HUD placement, hidden chest payloads, missing art, or clipped/misaligned visuals in tools/drsprite and internal/originalgame. Do not infer gameplay rules from screenshots; use diamond-rush-gameplay-fidelity for state-machine changes.
---

# Diamond Rush Asset Fidelity

Trace the Java resource loader and draw call before choosing an image or frame. Treat screenshots as final visual confirmation, not as the source of mapping or timing rules.

## Start With Render State

1. Read the root `AGENTS.md`.
2. Identify world, stage, coordinates, layer, raw ID, Runtime state bits, animation phase, and expected draw order.
3. Open [references/reference-map.md](references/reference-map.md) and read only the relevant reference sections.
4. If the discrepancy comes from incorrect Runtime state, use `$diamond-rush-gameplay-fidelity` first.

Do not patch a wrong state with a special-case sprite override.

## Trace The Original Resource

1. Find the Java draw branch for the relevant raw ID/state.
2. Follow `aClassfArr[]`, `textures[]`, resource masks, and world initialization back to the exact `.f` file and chunk.
3. Inspect `decoded/sprites/manifest.json`, the chunk's animations JSON, modules metadata, and modules image.
4. Resolve:
   - animation index and source-tick frame selection
   - module composition and offsets
   - horizontal/vertical flip flags
   - frame dimensions and multi-cell span
   - whether the caller invokes `applyFrameOffset()`
   - clipping and render-pass order
5. Never assume raw ID equals frame ID or that `frames.png` is a fixed `24x24` atlas.

Use the original `.f` metadata for game sprites. Direct PNG loading is acceptable only for true whole-image chunks without sprite metadata, such as extracted splash images.

## Implement Without Runtime Side Effects

- Cache decoded subimages and solid-color images during loading. Do not allocate images or repeatedly create `SubImage` values in Draw.
- Preserve PNG alpha and transparent atlas padding.
- Apply animation-frame `x/y` offsets only where the Java caller explicitly applies them.
- Preserve Java clipping, extra scan margins, and post-foreground passes.
- Keep closed chest payloads hidden until the reward frame.
- Keep interpolation temporary; render code cannot write coordinates, state bits, timers, or camera values back to Runtime.
- Fix mappings at the resource/world/state level where possible. Avoid stage-coordinate exceptions unless Java does the same.

## Verify

1. Inspect the decoded frame/module composition before running the game.
2. Exercise every affected animation phase, direction, world variant, and clipping boundary.
3. Verify render order while objects cross cell boundaries.
4. Check that the fix does not introduce black edges, opaque padding, diagnostic blocks, payload spoilers, or world-to-world sprite leakage.
5. Use screenshots only after source mapping and frame selection are established.

For broad renderer changes, run:

```bash
go test ./internal/original ./internal/originalgame ./tools/drsprite
go build -o /tmp/originalrush-smoke ./cmd/originalrush
rm -f /tmp/originalrush-smoke
```

Run `go run ./tools/drsprite ...` only when extraction output itself changed. Do not churn decoded output for a runtime-only mapping fix.

## Completion Gate

Report the exact `.f` file, chunk, animation/frame rule, anchor/flip behavior, draw pass, and source location that justify the fix. Distinguish source-audited mappings from visuals checked only in the Go runtime.
