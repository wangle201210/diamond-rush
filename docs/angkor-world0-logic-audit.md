# Angkor World 0 Logic Audit

Last audited: 2026-07-12

This document records the code-level completion audit for all 14 stages packed in `w0.bin`. The authority order is the original Java/JAR behavior, decoded stage/resource data, original runtime observation, and only then external guides or screenshots.

## Scope Result

- All low player and foreground IDs authored in World 0 have explicit runtime and render handling. No visible World 0 object should fall through to the diagnostic block renderer.
- Every normal stage, both normal/secret branches, all four secret stages, Great Anaconda, and the tutorial have source-cadence route tests.
- Shared mechanics are tested independently so a route cannot pass only because a test teleported, rewrote layers, or bypassed timing.
- The current boundary is cross-world content, not an unresolved World 0 object: Bavaria, Siberia, and Shop are not complete stages/screens in this repository.

## Shared Logic Matrix

| Area | Source behavior now enforced | Primary regression coverage |
| --- | --- | --- |
| Input | A changed direction turns first; movement starts on the next press/held tick. Space is phone `5`; Enter is checkpoint recall (`*`). Text input is disabled. | `game_test.go`, `tutorial_test.go` |
| Entrance/camera | raw `79` creates the open temporary entrance door; the final auto-step closes it. Camera follows source dead zones and scripted demos can lock it. | `runtime_test.go`, `game_test.go` |
| Checkpoints | Ordered automatic save after movement settles; Space/Enter on the checkpoint resets without a life; Enter away plays animation `19`, then consumes a life and restores. | `runtime_test.go`, `runtime_stage13_test.go` |
| Gravity | raw `0/1/9` use bottom-up, left-right scanning, immediate unsupported fall, `18 -> 12 -> 6 -> 0` interpolation, and separate slow roll preparation. | `runtime_test.go`, stage route tests |
| Boulder/gem roll | Rounded supports are raw `0/1/8/9`; roll transfer occurs at offset `12`; boulder landing emits sound `14`. | `runtime_test.go` |
| Digging/breakables | raw `10` becomes foreground raw `32` on the next object scan; raw `30` uses damage/chain state and eight-frame removal. | `runtime_test.go`, Stages 4/5/13 |
| Green/red snakes | Direction/pending direction, same-pass rescans, source contact overlap, red-snake three-hit state, chase target packing, and crush removal are retained. | `runtime_test.go`, Stages 5/8/12 |
| Crawler | raw `11` follows walls around corners with source scan order and `18 -> 13 -> 8 -> 3 -> -2` cadence; it does not merely reverse on collision. | `runtime_test.go` |
| Horizontal fire | raw `22/23` reach follows the extracted flame sequence, including the 20-tick empty/retracted frame. | `runtime_test.go`, renderer tests |
| Hammer | Direction animations `13..16`, impact tick `3`, raw `30` propagation, green-snake stun, red-snake damage, source sounds. | `runtime_test.go` |
| Hook | Horizontal 2-3-cell scan, blockers/candidate set, raw `32` rope timing, multi-cell pull, final raw `1` collection, and source target-state restoration. | `runtime_test.go`, Stages 5/7/8/10-13 |
| Freeze Hammer | Center/orthogonal scan, moving-snake pixel overlap, raw `9` gravity/push/hook/crush, source type transfer, thaw direction, and `0x78` post-thaw stun. | `runtime_freeze_test.go`, Stages 8/12 |
| Pressure/key doors | raw `7` low nibble starts as the count of same-group raw `6/8/9` activators. Each activation decrements once; zero starts `0x10 -> 0x20 -> 0x30`. | `runtime_test.go`, Stages 3/4/7/8/10/11 |
| Enemy arenas | raw `17` is a group marker and raw `26` selects the single active `cmInt`. Kills decrement that active group; zero opens doors/unlocks containers. | `runtime_enemy_gate_test.go`, Stage 12 |
| Locked chests | raw `14/33` hide payload art, remain inert while group-locked, then use the normal chest sequence after unlock. | `game_test.go`, `runtime_enemy_gate_test.go` |
| Chest animation | Start only after movement settles; source animations `40/48`, sounds `3/4`, reward ticks, overhead icon/effect, and final input release. | `runtime_test.go`, `game_test.go` |
| Violet values | raw `1` adds one. raw `41` adds its background value to quota, `VioletGems`, result count, HUD count, and saved bank. Full-health healing converts before animation to raw `41`/10. | `runtime_test.go`, `progress_test.go` |
| Permanent rewards | Save v6 records consumed red-diamond, awarded-extra-life, and relic chest coordinates only. Keys, healing, raw `41`, tools, and Compass remain replayable. | `runtime_test.go`, `progress_test.go` |
| Currency replay | Every clear adds that run's violet total, so violet can be farmed. Red adds only newly reachable pickups and accumulates partial runs through consumed coordinates. | `progress_test.go` |
| Results | 12-step Loading, sound `15`, phases/count-up/awards, Java coordinates (`15/32`, then `75/91`, `131/147`, `187/203`, `243/259`), and Continue flow. | `game_test.go` |
| Tutorial | Decoded scripts `29 -> 10 -> 11 -> 13 -> 15 -> 16 -> 17 -> 28`, fixed checkpoint art, desktop key prompts, Compass, recall, and final seal. | `runtime_stage13_test.go`, `tutorial_test.go` |
| Boss | Great Anaconda's three body columns, timing thresholds, tail strike, boulder regeneration, contact zones, three hits, gate delay, seal, and 11-step transition. | `runtime_stage08_test.go`, `game_test.go` |

## Stage Route Coverage

| Packed stage | Required covered route |
| --- | --- |
| `stage00` | Full Stage 1 route, quota, red chest, checkpoints, exit, and result. |
| `stage01` | Keys/locks, health conversion, falling gems/boulders, and exit. |
| `stage02` | Pressure doors, keys, Mystic Hammer route, and all quota mechanics. |
| `stage03` | Hammer pickup/script, red-snake combat, locks, and complete route. |
| `stage04` | Authored post-Bavaria hook revisit, hammer/hook puzzles, crawler, group snake, and exit. |
| `stage05` | Three collapse triggers, demo camera, falling torches, rising fire, checkpoint restore, and exit. |
| `stage06` | Normal exit and first secret exit as separate complete routes. |
| `stage07` | Hook-only normal route and Freeze Hammer secret route. |
| `stage08` | Full three-hit Great Anaconda fight and Angkor seal collection. |
| `stage09` | Secret Stage 1 pressure/key/hook route and raw `28` exit. |
| `stage10` | Secret Stage 2 six-boulder shaft, three silver locks, linked door, and exit. |
| `stage11` | All four active-group enemy arenas, four key chests/locks, quota, and exit. |
| `stage12` | raw `30` chains, hook/push route, one-way reward chamber, quota 99, and exit. |
| `stage13` | Complete original tutorial script sequence and transition to Stage 1. |

## Cross-World Prerequisites

- Mystic Hammer raw `24` is in Angkor Stage 4 at `(26,18)`.
- Mystic Hook raw `27` is not in Angkor Stage 5. It is in `world1/stage02` (Bavaria Stage 3) at `(24,25)` under foreground raw `14`. The Angkor-only Stage 5 entry supplies level `2` to model the original revisit state; do not add a fake hook chest to World 0.
- Freeze Hammer raw `26` is in `world2/stage05` at `(32,22)`. The Stage 8 secret-route and secret-stage entries similarly use a source-valid later-game prerequisite state.
- Reproducing those acquisition journeys requires implementing the corresponding Bavaria/Siberia stages. This is outside the completed World 0 stage-data slice and must not be described as implemented.

## Verification

Run these before declaring a World 0 change complete:

```bash
go test ./internal/original -count=1
go test ./internal/originalgame -count=1
go test ./... -count=1
go build -o originalrush ./cmd/originalrush
```

Any future discrepancy must be resolved from Java/JAR state and decoded data first. A passing route test is necessary but does not override contradictory source evidence.
