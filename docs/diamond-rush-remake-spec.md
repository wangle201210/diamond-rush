# Diamond Rush Go Remake Spec

Last updated: 2026-07-11

This is the working specification for the source-fidelity Go runtime built from the local Java reference. The target is behavior-first fidelity to the original Nokia Java ME Diamond Rush, backed by decoded original data and deterministic Go tests.

Reference project:

`/Users/wanna/mine/github/wangle201210/DiamondRushSource`

Source mapping:

`docs/diamond-rush-source-mapping.md`

Existing gameplay research:

`docs/diamond-rush-original-gameplay.md`

## Product Target

Build a Go + Ebitengine Diamond Rush remake that feels like the original Nokia Java ME game:

- 240x320 phone-screen composition.
- World map, title flow, shop, stage clear marks.
- Tile-based side-view adventure puzzle stages.
- Phone-key input model centered on `5`.
- Health, lives, checkpoint recall/reset, red/violet gems, chests, tools, hazards, enemies, bosses.
- Original-style authored stages, not a generic Boulder Dash board.

## Non-Goals

- Do not line-by-line port the Java decompilation.
- Do not replace decoded source data with guessed or self-authored stage behavior.
- Do not use a full-map desktop board as the default gameplay view.
- Do not make hook/hammer mouse-targeted.
- Do not use original copyrighted visual/audio assets directly in the final distributable without a separate asset decision.

## Fidelity Workflow

Every feature should go through this loop:

1. Observe the behavior in the Java reference runtime.
2. Locate or approximate the corresponding Java source region.
3. Write a short behavior note or update this spec.
4. Implement clean Go behavior.
5. Add a focused unit/integration test where the behavior is deterministic.
6. Compare manually against the Java runtime.

If observation and existing Go behavior conflict, trust the Java runtime/source.

## Required Tooling Before Rebuild

### `tools/drdecode`

Purpose: decode original stage packs and make data inspectable.

Input:

- `DiamondRushSource/src/main/resources/w0.bin`
- `DiamondRushSource/src/main/resources/w1.bin`
- `DiamondRushSource/src/main/resources/w2.bin`

Minimum output:

- `decoded/world0/stageXX.json`
- raw dimensions.
- raw player/background/foreground layer bytes.
- tile ID frequency tables.
- optional CSV preview layers.

Source basis:

- `i.java` lines 3407-3473 show the stage pack layout.

Acceptance:

- Decoding `w0.bin` produces every Angkor stage without manual sizes.
- Re-encoding or loading the decoded dimensions matches the Java loader dimensions.
- Unknown byte IDs are preserved, not guessed away.

### `tools/drinspect`

Purpose: inspect decoded stages and classify IDs.

Implemented features:

- Render raw byte layers with stable colors.
- Show tile coordinates and byte IDs.
- Export screenshots for documentation.
- Generate static per-stage HTML inspectors.
- Generate a raw ID frequency index.

Acceptance:

- We can open `decoded/preview/index.html`.
- We can open Angkor Stage 1 at `decoded/preview/world0/stage00.html`.
- We can hover cells to inspect raw `player`, `foreground`, and `background` byte IDs.
- The PNG output is understood as raw-ID diagnostics, not original-art rendering.

### `tools/drsprite`

Purpose: extract original `.f` sprite chunks into inspectable PNG and JSON assets.

Implemented features:

- Parse chunked `.f` resource files.
- Decode Gameloft sprite modules, frames, animations, palettes, and module pixel data.
- Export module sheets and frame sheets as PNG.
- Export frame/animation metadata as JSON.
- Generate `decoded/sprites/index.html` for browser inspection.

Current result:

- `94` sprite chunks exported from `121` chunks.
- Key outputs include:
  - `decoded/sprites/0/chunk02-frames.png`: Angkor world tile/frame sheet.
  - `decoded/sprites/o/chunk00-frames.png`: hero frame sheet.
  - `decoded/sprites/manifest.json`: parse status and metadata for all chunks.

Acceptance:

- Original hero and Angkor world sprite sheets are visible as real pixel art, not raw-ID color blocks.
- Non-sprite chunks are recorded as parse errors without breaking export.

### Full-Stage Visual Preview

Purpose: render decoded stages with original sprite art.

Current status:

- Not complete yet.
- `tools/drinspect` shows raw IDs as colors.
- `tools/drsprite` exports real sprites.
- The missing bridge is the Java render mapping from raw ID/runtime state to sprite file, chunk, frame, palette, offset, and animation state.

Acceptance:

- Angkor Stage 1 can be rendered as a static original-art composite close to the Java runtime's first visible screen.
- The renderer documents every ID mapping it uses and leaves unknown IDs visible, not silently blank.

### `internal/original` And `tools/drworldaudit`

Purpose: make decoded original data loadable by the Go codebase before mapping it into gameplay objects.

Implemented features:

- `internal/original` loads decoded stage JSON and world manifests.
- It preserves the original three layers:
  - player/object raw bytes.
  - background/state raw bytes.
  - foreground raw bytes.
- It exposes coordinate lookup, unique IDs, histograms, entrance markers, Java signed-byte conversion, and a small source-anchored ID role table.
- It has an initial runtime state that copies mutable layers, applies the Java-style entrance marker rule, indexes source-anchored checkpoints/goals/doors, and supports checkpoint layer snapshot/restore.
- `tools/drworldaudit` generates `docs/angkor-world0-inventory.md` from `decoded/world0`.

Acceptance:

- `go test ./internal/original` verifies the Angkor world has 14 stages, the expected dimensions, one entrance marker per stage, source-anchored IDs such as checkpoint foreground raw `4`, Java-style runtime entrance initialization, and checkpoint restoration of mutable layers.
- `docs/angkor-world0-inventory.md` is the current implementation checklist for first-world raw ID coverage.

### Asset Extraction Or Replacement

Initial path:

- Use generated placeholder art or existing temporary art while behavior is rebuilt.
- Keep object proportions and screen layout close to the Java runtime.

Later path:

- Use `tools/drsprite` PNG/JSON exports to guide proportions, silhouettes, animation categories, and source-fidelity debug rendering.
- Generate replacement sprite sheets only after the original roles and proportions are understood.

Acceptance:

- The gameplay logic must not depend on final art.
- Placeholder art must still preserve tile size, hitbox size, facing, and animation timing categories.

## Runtime Architecture

| Area | Package/path | Required behavior |
| --- | --- | --- |
| Executable | `cmd/originalrush` | Start the Ebitengine runtime. |
| Source runtime | `internal/original` | Preserve three original layers, object timers, source scan order, movement interpolation, tools, hazards and checkpoint state. |
| Game adapter | `internal/originalgame` | Phone-key mapping, 240x320 camera/HUD, source sprites, audio, stage flow, world map and JSON progress. |
| Decoded data | `decoded/` | Original world packs, map graph, sprite metadata/images, fonts and audio. |
| Extraction/audit | `tools/` | Reproducible decoding and source-data inspection. |

Keep gameplay facts in `internal/original`; Ebitengine and host-platform concerns belong in `internal/originalgame`.

## Original-Behavior Requirements

### Screen And Camera

Required:

- Internal logical screen is `240x320`.
- Stage tiles are 24 px, as implied by Java loader/rendering and camera math.
- Gameplay view scrolls around the player.
- HUD/status area must fit within the phone-screen composition.
- Desktop window may scale the 240x320 render, but gameplay logic and UI layout must not depend on desktop resolution.

Acceptance:

- The player can move across a stage larger than the viewport without seeing the whole map.
- Resizing the desktop window does not change tile physics or HUD layout.

### Input

Required:

- Direction keys update facing/movement.
- `5` is the single context/action key.
- `*` triggers checkpoint recall/death-style reset behavior.
- Soft keys are reserved for menu/back/skip behavior.

Original source anchors:

- `i.java` key constants lines 16-44.
- `keyPressed`/`keyReleased` lines 861-869.
- `getKeyFromKeyCode` lines 15811-15837.

Acceptance:

- A phone-key event stream can drive the whole game without mouse input.
- Desktop controls are configurable but map onto phone keys.

### Stage Data

Required:

- Decode original world packs before authoring final levels.
- Preserve three source layers:
  - player/object layer.
  - background/state layer.
  - foreground layer.
- Convert raw IDs only after they are observed/classified.

Acceptance:

- Angkor Stage 1 can be loaded from decoded source data into the Go runtime or into a debug viewer.
- Unknown IDs remain visible in tooling.

### Movement And Tile Physics

Required:

- Player movement is tile/grid based with original-feeling animation delay.
- Dirt/empty/wall/object passability must come from decoded object classes, not implementation guesses.
- Pushing/falling/rolling objects must be deterministic and match runtime observation.
- Falling objects can damage or kill the player according to original damage rules.

Acceptance:

- Angkor Stage 1 movement and first object interactions match the Java runtime at the level of tile outcomes.
- Falling object order is stable under tests.

### Health, Lives, Damage

Required:

- Player has health/energy, not only one-hit death.
- Damage may be partial.
- Hurt/temporary invulnerability state prevents repeated instant hits.
- Track hit count for stage-clear marks.
- Lives are separate from health.

Original source anchors:

- `hurtHero` starts at `i.java` line 761.
- Stage-clear mark checks include no-hit/no-restart-like criteria around lines 2183-2199.

Acceptance:

- A low-damage enemy reduces health without always ending the stage.
- Repeated contact does not drain all health in one frame unless original behavior does so.
- Stage results can identify no-damage clear.

### Checkpoints

Required:

- Checkpoint activation snapshots mutable room state.
- Stepping onto foreground raw `4` activates/saves the checkpoint automatically by its background/state order.
- `5` while standing on active checkpoint restores/resets checkpoint state.
- `*` while standing on active checkpoint restores/resets checkpoint state without life cost.
- `*` away from checkpoint triggers recall/death-style reset and consumes the appropriate resource.
- Snapshot includes player state, counters, and mutable layers, but not the extra-life count.

Original source anchors:

- Checkpoint `5` branch around `i.java` lines 1266-1271.
- `*` branch around lines 1500-1526.
- `XVoid()` starts around line 11018 and decrements `azInt` before restoring checkpoint state.
- `saveCheckpoint` starts line 15190.
- `restoreCheckpoint` starts line 15248.

Acceptance:

- Moving rocks/opening doors/collecting items after checkpoint can be reset according to original behavior.
- Extra lives collected after checkpoint persist across checkpoint restore because Java `restoreCheckpoint()` does not restore `azInt`.
- Re-entering the same or an earlier checkpoint does not overwrite the active snapshot.
- Recall away from checkpoint consumes an extra life and flags recall use; checkpoint reset does not.

### Action Key And Tools

Required:

- `5` resolves context before tools:
  1. checkpoint action if on checkpoint.
  2. hammer/local action if hammer is available.
  3. hook/line action if hook is available.
  4. generic confirm/interact in menus/cutscenes.
- Hammer and hook use facing/direction, not mouse targeting.
- Hook target rules must be decoded and observed before finalizing. Do not assume "only rocks."

Original source anchors:

- Main `KEY_OK` branch starts around `i.java` line 1266.
- Tool availability checks use `iByteArr[9]` in the action branch.
- Current Go runtime derives that tool level from collected source-special raw `24`/`27`/`26` effects; hammer/local action requires level `>=1`, hook requires level `>=2`.
- The horizontal hook scan at lines 1278-1344 searches 2-3 cells and recognizes raw `0/1/8/9/11/14/19/43/47/48`; closed doors, intervening objects, overlay blockers, and raw `48` state bit `0x8` constrain the cast.
- Original-JAR traces show raw `32` rope segments extending with `18 -> 12 -> 6 -> 0` timers. Physical targets are pulled all the way to the adjacent cell. Raw `1` is pulled one additional step into the hero cell and collected only after the hook finishes.
- On release, the JAR restores a saved state only for raw `0/8/9/47`, after clearing pending-turn and roll-preparation bits. Other hook candidates receive state `-1`; hooked snakes therefore resume through the source negative packed-target branch instead of recovering their old patrol state.
- The hammer/local neighborhood scan appears around lines 1345-1499. Hero animations `13/14/15/16` impact on tick `3` and last 11 ticks upward or 12 ticks in the other directions.
- Hook sound ID is `SOUND_SFX_HOOKING = 12`.

Acceptance:

- Hammer can break/stun only original-eligible targets.
- Hammer is unavailable before the raw `24`/equivalent tool level is acquired.
- Hook can pull verified physical candidates through every source step without moving the player; verified raw `1` collection occurs only after its final physical pull into the hero cell.
- Hook behavior matches at least three runtime-observed examples before being marked complete.

### Gems, Chests, Shop, And World Unlock

Required:

- Violet gems and red gems are separate currencies.
- Violet gems fund shop upgrades.
- Red gems unlock later worlds.
- Chests can carry important rewards and must not be mere score pickups.
- Stage result records all gem/secret/perfect marks.

Original source anchors:

- `a_Config.itemPrices = {150, 400, 1000, 3000}`.
- `a_Config.worldPrices = {0, 10, 25}`.
- Stage reward update around `i.java` lines 2100-2205.

Acceptance:

- Save data distinguishes total violet, total red, stage-collected violet, stage-collected red.
- The shop can spend violet gems without altering red gems.
- World unlock checks depend on red gem totals.

### World Map

Required:

- Three world nodes plus shop/seal layout:
  - Angkor
  - Bavaria/Scotland
  - Tibet/Siberia
  - Shop
- World movement follows the original seal graph.
- Stage nodes have completion marks.

Original source anchors:

- `a_Config.sealMoveDirections`.
- `map_angkor.out`, `map_scotland.out`, `map_tibet.out`.
- World map file selection around `i.java` lines 1940-1957 and 5139-5150.

Acceptance:

- Player can enter Angkor, see stage nodes, return to seal/shop, and preserve stage marks.

### Enemies, Hazards, Bosses

Required:

- Start with Angkor enemies/hazards only.
- Implement hazard damage through the same health/hurt pipeline.
- Boss stages are special cases; do not fake them as a normal enemy with more health.

Original source anchors:

- Special stage constants:
  - Angkor falling torches stage `5`.
  - Angkor Great Anaconda stage `8`.
  - Bavaria Evil Teutonic Knight stage `9`.
  - Siberia Yetti stage `10`.
- Special world/stage setup around `i.java` lines 3331-3379.

Acceptance:

- First rebuild milestone includes at least one normal hazard, one enemy, and the Angkor boss setup documented from runtime observation.

## Rebuild Phases

### Phase 0: Documentation And Decode

Deliverables:

- `docs/diamond-rush-source-mapping.md`.
- `docs/diamond-rush-remake-spec.md`.
- `tools/drdecode` can decode `w0.bin`.
- Debug output includes stage dimensions and raw layer IDs.

Exit criteria:

- No new gameplay feature work until Angkor Stage 1 raw data is inspectable.

### Phase 1: Minimal Original Stage Runtime

Deliverables:

- Load decoded Angkor stages through `internal/original`.
- Load decoded Angkor Stage 1 into Go.
- Render a 240x320 viewport.
- Move player on original grid data.
- Basic passability from classified source IDs.
- Save/progress skeleton.

Exit criteria:

- Go runtime can walk through the opening stage area using decoded data.

Current implementation:

- `cmd/originalrush` loads all 14 stages from the decoded Angkor pack. `stage00..08` are the normal route and boss, `stage09..12` are the four secret stages exposed through the decoded `map_angkor.out` graph, and the non-map `stage13` tutorial runs before Stage 1 for a fresh save.
- `internal/originalgame` renders the original `240x320` composition: a `40px` source HUD, a `240px` scrolling playfield with `24px` tiles, and a `40px` source HUD. It uses extracted Angkor floor/wall/boulder/vegetation art, source gem/checkpoint/goal/enemy/hazard art, and the original hero animation metadata.
- HUD frames, hero frames, snake frames, foreground effects, and flame frames are composed from their extracted source modules and JSON offsets/flip flags. This avoids the clipping present in convenience `frames.png` sheets for non-zero frame anchors.
- Rendering follows the Java pass order: first the scrolling tile background, then dynamic foreground/player objects over the relative `-1..11` source scan, then late foreground overlays. This keeps rolling rocks and horizontal flames from being overwritten by later floor cells, clips pressure switches to one tile, and preserves the authored foreground occlusion over the hero.
- Angkor Stage 1 runs at the Java loop's `20 TPS`. Its 60-tick title is an overlay rather than a simulation pause, and player-layer raw `79` automatically walks the hero four cells from `(0,17)` to the first checkpoint before accepting movement input.
- Stage initialization creates the source temporary foreground raw `7` entrance door at `(2,17)` with merged state `0x3f`. The fourth raw `79` auto-entry step runs `doorHeadClose`, changes it to blocking state `0x0f`, and emits source sound `14`.
- Player movement uses the source `18 -> 12 -> 6 -> 0` sub-tile offset. The camera follows the rendered position with the Java horizontal and vertical dead zones instead of hard-centering each tile.
- Player-layer raw `12` is a visible blocking quota gate. It remains in the runtime layer, renders `cm.f` chunk `5` with the remaining count, and clears only when raw `1`/`41` collection exhausts the quota. The readable decompile's initialization clear is contradicted by original-JAR runtime state and its own render branch. Raw `5` still does not inspect that quota, so an authored route that bypasses raw `12` may enter the exit independently.
- Entering Stage 1 raw `5` reads exit direction `2` from its background byte, shows `CONGRATULATIONS!`, and auto-walks right with `18 -> 12 -> 6 -> 0` interpolation. Only after the hero reaches `x = stageWidth + 6` does the game run the Java `bByte=35` 12-step Loading transition and enter the `bByte=17` result sequence. Continue then returns to the Angkor map and unlocks the next implemented node.
- The runtime walks through the opening corridor, commits ordered checkpoints on the source tick after movement settles, collects raw `1` violet gems through 24-pixel interpolation overlap, and plays `cm.f` chunk `7` animation `3`. Entered raw `10` vegetation remains for the movement frame, becomes foreground raw `32` on the next object scan, then advances its eight frames on source tick parity. Its foreground frame state is stored separately from gravity-object direction/rotation state, including in checkpoint snapshots.
- The `(19,2)` foreground raw `33` chest starts only on the settled-player object scan. Stage 1 uses hero animation `40` for 67 ticks, advances the lid through source states `1/2/3`, grants the red diamond on sequence index `13` at tick `39`, plays `cm.f` chunk `7` animation `0`, and displays the red diamond above the hero after the reward sequence.
- Raw `0` and raw `1` fall vertically on the first eligible object frame with source offsets `18 -> 12 -> 6 -> 0`. Rolling uses packed direction/rotation bits plus Java's `0x200` preparation state: offset `1` to `12`, 24 to 27 source frames of visible rocking, then a diagonal transfer with vertical offset `12 -> 6 -> 0`. Purple gems can be collected during interpolation; boulders preserve partial damage/crush behavior and remove a snake below before occupying its cell. Landing clears the packed low direction bits, preserves rotation/side markers, and emits source sound `14`.
- Green snakes preserve low/pending direction bits, `21 -> 18 -> ... -> 0` motion, same-pass rescans, source foreground passability, 24-pixel contact overlap, and directional knockback. Their extracted animation uses direct `(aSInt >> 1)` sequence selection. Horizontal fire reach remains tied to the extracted flame animation frame.
- `*` away from a checkpoint emits source sound `2`, then plays the complete 42-tick hero animation `19` before consuming a life and restoring; `5` or `*` while standing on a checkpoint resets immediately with sound `9` and no life cost. Lethal damage uses the 88-tick hurt/death transition and restores at full `4/4` health when an extra life remains.
- `TestRuntimeStage00CanBeCompletedAtSourceCadence` through `TestRuntimeStage13TutorialCanBeCompletedAtSourceCadence` replay every packed Angkor stage while advancing the unified 20 Hz object frame between player actions. They cover quotas up to 99, ordered checkpoints, gravity puzzles, locks, enemy gates, hammer/hook/freeze puzzles, the Stage 6 falling-torches chase, both exits in Stages 7 and 8, four secret stages, three actual falling-boulder hits against Great Anaconda, and the complete tutorial script chain.
- Stage 7 foreground raw `5` runs the normal result flow and unlocks normal node `7`. Foreground raw `28` keeps its separate `qByte/atBoolean=false` branch, skips results, displays the source 30-tick secret-path message, and unlocks map node `9` only.
- Stage 8's normal route uses the `(33,12)` boulder on pressure switch `(32,17)` to reach the gold key and normal exit. Its secret route requires Freeze Hammer raw `9`: freeze the snake at `(8,2)`, push the settled block onto `(7,4)`, and use raw `28` to unlock map node `12` only.
- Freeze Hammer stores source type separately while raw `9` participates in source gravity, rolling, pushing, hook, pressure, damage, enemy crush, movement interpolation, and checkpoint restore. Frozen violet/snake frames come from `gen0.f` chunk `1` and `gen1.f` chunk `6`.
- Stage 9 uses the source `kByte=4` Great Anaconda state machine: three body columns, strict `>10/>20/>40` phase thresholds, raw `50` contact zones, y `7/8` falling-boulder vulnerability, 50-tick charge, 12-tick tail strike, health-dependent retraction, boulder regeneration, and the `>80` death-door delay.
- The Great Anaconda body/platform/tail use extracted `b0.f` chunks `0/1` and `gen1.f` chunk `0` animation `2`. Its three-segment health bar, 30-tick regeneration shake, body contact bounds, gate-introduction camera, sounds, and checkpoint reset are source-mapped.
- Stage 9's raw `53` chest is hidden until opened and maps to Angkor seal bit `0`. The reward switches to hero animation `47`, locks input until source tick `>140`, and uses an 11-step seal loading transition instead of ordinary results or medal awards.
- Secret stages `09 -> 10 -> 11` follow the first raw-`28` map branch; `stage12` is the separate Stage 8 branch. Their route tests exercise pressure-rock chains, key/lock groups, enemy arenas, raw `30` wall clusters, the six-boulder key shaft, and terminal raw `28` exits without teleporting or mutating layer data. Stage 11 now solves all four enemy groups before opening their four source-locked key chests, including the three-snake ice stack and final boulder cascade.
- Tutorial scripts run in source order `29 -> 10 -> 11 -> 13 -> 15 -> 16 -> 17 -> 28`, including camera/movement/foreground commands, source text indices, portraits, checkpoint reset, life-cost recall, white flash, and the final seal walk. The same decoded interpreter runs Stage 3 script `30`, Stage 4 tool script `22`, and Stage 9 Boss script `33`. Input is locked while a script is active; skipping text does not skip state-changing commands.
- Tutorial portraits are extracted from `demoSpr.bin`; the final seal is composed from `mmv.f` chunk `0`. Compass raw `42` uses the original `gen3.f` chunk `1` 24x24 module for the overhead reward and enables the separate `ui.f` HUD direction frames.
- Progress save version 6 stores explicit stage-node unlock bits, `TutorialComplete`, `RelicMask`, per-world unlock bits, and source-equivalent consumed reward coordinates. This preserves the non-sequential Angkor graph and permanently opens only collected red-diamond, awarded-extra-life, and relic chests.
- The runtime collects and renders World 0's source-special pickups raw `24`/`26`/`27`, raw `42`, and raw `53`. Raw `53` includes its source seal bit, celebration, input lock, persistence, and transition; a later-world inventory/shop screen remains outside this Angkor slice.
- Raw `41` is a value-bearing violet-diamond reward: it increments the bottom-right violet HUD count, stage result count, quota, and saved bank. Full-health raw `7` and max-life/full-health raw `6` convert to raw `41` before their reward animation.
- Mystic Hook raw `27` is sourced from the foreground raw `14` chest at Bavaria Stage 3 `(24,25)`, not Angkor Stage 5. The Angkor-only Stage 5 entry models the post-Bavaria revisit prerequisite and must not invent a hook chest in World 0.
- The runtime treats player-layer raw `33` as a passable persistent marker with no fallback object sprite; foreground raw `33` owns the visible overlay.
- Foreground raw `7` uses merged high-nibble phases. Its low nibble is initialized to the number of same-group raw `6/8/9` activators still required; opening begins at `0x10` only when the count reaches zero, advances every third source tick to passable `0x20` and final `0x30`, and closing preserves the remaining count.
- Foreground raw `6` is a pressure switch that opens linked raw `7` doors while pressed or occupied and closes them when released. Its `gen2.f` chunk `9` module follows source depression interpolation.
- Foreground raw `0` is a source-anchored one-shot event. It triggers when movement interpolation reaches `jInt<=6`, records its decoded background/state and clears the cell; Stage 6 background `3` executes the decoded demo script rather than acting as a passive counter.
- Foreground raw `1` is a clearable cluster, recursively clearing connected raw `1` cells when entered.
- Foreground raw `2` is player-passable and exposes a current-cell tool prompt. State `0` blobs recursively clear once adjacent raw `30` walls are gone; state `1` blobs require the player to stand on the cell and press Action/5 with tool level `>=2`.
- Foreground raw `17` cells are group markers; raw `26` selects Java's single active `cmInt`. Enemy/container removal decrements that active group regardless of the object's original marker, then opens same-group raw `7` doors or unlocks raw `14/33` containers without exposing payloads early.
- The Angkor renderer uses extracted original art for every actually visible low raw ID across all 14 packed stages, including the `cm.f` chunk `5` quota marker, `gen1.f` chunk `4` crawler, `gen2.f` chunk `9` pressure switch, `gen0.f` chunk `1`/`gen1.f` chunk `6` frozen objects, code-drawn raw `32` hook rope, all three `mm0.f` falling-torches chunks, Stage 9 `b0.f`/`mmv.f` assets, the Compass module, and tutorial portrait/seal resources. All 120 closed-container payload cells are audited against early reveal.
- Stage title, checkpoint, congratulations, loading, and result text use deterministic atlases exported from FreeJ2ME's logical `SansSerif Bold` 10px/12px fonts. Source panel fill/border colors and the original y-offset behavior are preserved.
- `snd.f` is decoded into all 21 original standard-MIDI tracks. The Stage 1 runtime applies the JAR priority table and 50ms equal-priority guard, plays Angkor track `16`, and emits source IDs for door/boulder `14`, hurt `5`, checkpoint `9`, chest `3`/`4`, death `2`, and result `15`; macOS playback uses `AVMIDIPlayer`.
- Result rendering follows original-JAR bytecode coordinates and phase threshold behavior, including flat animation-sequence indexing for the three different award effect shifts. Stage 1 clear state and the four award bits are persisted so already-earned effects do not replay as newly earned.
- Stage 9's 11-step seal handoff enters the implemented four-position global seal selector. Angkor opens the decoded world map; Bavaria and Siberia use the source `10/25` red-diamond unlock thresholds and Shop participates in source navigation, while those three destinations explicitly remain outside the current Angkor content slice.

Remaining work outside the complete Angkor stage-data slice:

- Implement the shop, full RMS-equivalent economy/progression fields, Bavaria/Tibet maps, and the original cross-world trip that grants the hook before revisiting Stage 5.
- Wire the remaining sound IDs when their later-stage objects and global screens are implemented.

### Phase 2: Core Adventure Mechanics

Deliverables:

- Health/lives/hurt state.
- Violet/red gems and chests.
- Checkpoint save/restore and `*` recall.
- Falling objects and one normal enemy/hazard.

Exit criteria:

- All 14 packed Angkor stages have source-cadence route regressions with their actual decoded mechanics enabled, including both exits in Stages 7 and 8, all four secret stages, the Boss seal path, and the tutorial.

### Phase 3: Tools And Backtracking

Deliverables:

- Hammer behavior.
- Hook behavior from original-JAR per-tick runtime observation.
- Tool acquisition flags.
- Backtracking gates and secret/revisit rewards.

Exit criteria:

- At least three tool interactions match Java runtime observation.

### Phase 4: World Map, Shop, Completion

Deliverables:

- Angkor map nodes.
- Shop upgrades with original prices.
- Stage clear marks.
- Red-gem world unlock gate.

Exit criteria:

- Clear stages, spend violet gems, preserve red gems, and see map marks.

### Phase 5: Polished Five-Level Slice

Deliverables:

- Five tuned levels, ideally adapted from original Angkor flow rather than invented cave boards.
- Opening/menu/title flow.
- Source-mapped Great Anaconda boss, Angkor seal reward, and seal-specific transition.
- Replacement assets and audio cues.

Exit criteria:

- A player familiar with Diamond Rush recognizes the interaction model, screen composition, progression, and tool gating.

## Current Runtime Entry Points

Run the source-data implementation:

```bash
go run ./cmd/originalrush
```

Regenerate Angkor stage data when the decoder changes:

```bash
go run ./tools/drdecode \
  -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/w0.bin \
  -out decoded/world0
```

The generated files include:

```text
decoded/world0/manifest.json
decoded/world0/stage00.json
decoded/world0/stage01.json
...
```

`manifest.json` should include:

- source file path.
- world index.
- stage count.
- each stage width/height.
- per-layer byte histograms.

Each stage JSON should include:

- width.
- height.
- player layer raw bytes.
- background layer raw bytes.
- foreground layer raw bytes.
- optional notes section for manual annotations.
