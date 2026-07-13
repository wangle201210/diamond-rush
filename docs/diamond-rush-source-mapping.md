# Diamond Rush Source Mapping

Last updated: 2026-07-11

This document maps the local Java reference project to the Go remake work. The Java project is the current source-of-truth for behavior research:

`/Users/wanna/mine/github/wangle201210/DiamondRushSource`

Use it as a reference implementation and runtime oracle. Do not translate it line-by-line. The goal is to reproduce observable Diamond Rush behavior in Go with a clean domain model.

The readable `i.java` contains decompiler artifacts. When control flow is impossible, inspect the original `jars/diamondrush.jar` with `javap -p -c`; do not use `target/classes`, which was rebuilt from the flawed decompilation. The green-snake blocked-direction branch is a confirmed example: the readable source appears to loop, while the original JAR writes the reverse direction to `0x7000`, clears the low direction, and continues normally.

## Reference Project Status

The local Java project is a Nokia Java ME decompilation with a FreeJ2ME runtime wrapper.

Evidence:

- `README.md`: "Decompilation of Gamelofts' Diamond Rush (2006) video game for Nokia mobile devices."
- `src/main/resources/META-INF/MANIFEST.MF`:
  - `MIDlet-Name: Diamond Rush`
  - `MIDlet-Vendor: Nokia`
  - `MIDlet-Version: 1.1.7`
  - `MIDlet-1: Diamond Rush, /icon.png, GloftDIRU`
- Local run script: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/run-diamond-rush.sh`
- Working runtime settings:
  - Canvas: `240x320`
  - Scale: `2`
  - Phone mode: `Nokia`
  - JAR: `jars/diamondrush.jar`

Run command:

```bash
cd /Users/wanna/mine/github/wangle201210/DiamondRushSource
./run-diamond-rush.sh
```

## Source Layout

| Java path | Role | Go remake use |
| --- | --- | --- |
| `src/main/java/GloftDIRU.java` | MIDlet entry. Creates `i` and sets it as the current displayable. | Only useful for lifecycle understanding. Go should use an Ebitengine `Game` entry. |
| `src/main/java/i.java` | Main game class. Input, state machine, loading, world map, stage logic, rendering, persistence, collision, tools, hazards. | Primary behavior reference. Break into Go packages instead of copying the monolith. |
| `src/main/java/a_Config.java` | Small constants for shop prices, world unlock prices, seal map movement. | Use for economy/world-map constants. |
| `src/main/java/f_Sprite.java` | Gameloft sprite resource parser: modules, frames, animations, palettes. | Reference for extracting visual assets if needed; Go runtime should use generated PNG/JSON atlases. |
| `src/main/java/b_SpriteAnimator.java` | Animation wrapper over `f_Sprite`. | Reference for animation timing/state names after assets are decoded. |
| `src/main/java/j_SoundManager.java` | Sound ID table and MIDI playback priority. | Reference event IDs, priority groups, and extracted original MIDI playback. |
| `src/main/java/c.java` | Demo/cutscene interpreter-like class. | Reference opening/tutorial sequences if recreating presentation. |
| `src/main/java/org/recompile/freej2me/FreeJ2ME.java` | Desktop emulator wrapper. Locally patched for Mac input/repaint. | Runtime oracle only, not game logic. |
| `src/main/resources/w0.bin` | Angkor Wat stage pack. | Decode into the preserved three-layer Go/JSON format. |
| JAR `w1.bin` | Bavaria/Scotland stage pack. | Decoded into `decoded/world1`; all 13 source stages load through the Go runtime. The readable-source copy differs at four Stage 8 object cells, so the JAR entry is canonical. |
| `src/main/resources/w2.bin` | Tibet/Siberia stage pack. | Decode later after Angkor/Bavaria. |
| `src/main/resources/map_angkor.out` | Angkor world map metadata. | Reference stage-node layout and unlock flow. |
| `src/main/resources/map_scotland.out` | Bavaria/Scotland world map metadata. | Decoded into `decoded/world1/map.json` and used by the Bavaria world map. |
| `src/main/resources/map_tibet.out` | Tibet/Siberia world map metadata. | Later-world map reference. |
| `src/main/resources/*.f` | Gameloft packed sprite/text/sound resources. | Decode or replace with generated assets, but preserve visual roles and proportions. |
| `rms/Diamond Rush/*` | FreeJ2ME local RMS save files from testing. | Useful for observing persisted bytes; not a design source by itself. |

## Runtime And Input Anchors

`i.java` defines phone-key bitmasks at lines 16-44. Important constants:

| Original input | Java constant | Meaning for Go |
| --- | --- | --- |
| `2` / d-pad up | `KEY_UP` | Move or face up. |
| `8` / d-pad down | `KEY_DOWN` | Move or face down. |
| `4` / d-pad left | `KEY_LEFT` | Move or face left. |
| `6` / d-pad right | `KEY_RIGHT` | Move or face right. |
| `5` / center | `KEY_OK` | Context action: checkpoint reset, hammer, hook, interaction. |
| `*` | `KEY_RESPAWN` | Checkpoint recall / death-like reset path. |
| right soft key | `KEY_CANCEL` | Back/cancel/menu. |

`i.keyPressed` and `i.keyReleased` at lines 861-869 convert raw key codes via `getKeyFromKeyCode` and update a global `keysPressed` bitset. `getKeyFromKeyCode` at lines 15811-15837 maps negative Nokia d-pad codes, `42` for `*`, `35` for `#`, and ASCII numeric keys `48..57`.

Go implication:

- Keep a phone-first input abstraction instead of binding gameplay directly to desktop keys.
- Model the game in terms of `Up`, `Down`, `Left`, `Right`, `Action5`, `Star`, `Pound`, `LeftSoft`, `RightSoft`.
- Desktop mappings are adapter concerns only.
- The current macOS adapter maps `Space`/`5`/numpad `5` to `Action5`, and maps `Enter`/`R`/`Backspace`/numpad `*`/`Shift+8` to `Star`. Tutorial Skip uses `S`, so Space never has both Skip and interaction semantics.
- Player-facing desktop prompts must name the mapped keys as `SPACE`, `ENTER`, and `S`; source phone labels or sprites such as `Press 5`, the blue `5`, and the blue `*` must not be shown as keyboard instructions.
- A new direction does not move the hero immediately. `i.java:1531-1532` marks a stationary facing change with `kInt|=0x1000`; the player update at `10692-10699` turns in place with `jInt=18`, then counts `18 -> 12 -> 6 -> 0`. A released tap only turns, a second same-direction press moves, and a held direction moves on the update after the turn finishes. Apply this only to player input; scripted and entrance movement bypass it.

## Main State Machine

`i.run()` starts at line 890. The large `handleKeyPresses()` method starts at line 958 and drives most states.

Important identified states and regions:

- Stage gameplay input branch around lines 1189-1530.
- World-map branch around lines 1730-1760.
- World/stage-map resource loading around lines 1940-2006.
- Stage result and reward processing around lines 2100-2205.
- Stage loading branch around lines 3311-3484.
- Rendering entry `paint(Graphics)` starts at line 5377.

Go implication:

- Split the monolithic state machine into explicit modes:
  - `Boot`
  - `Title`
  - `WorldMap`
  - `Shop`
  - `Cutscene`
  - `StageLoading`
  - `StagePlaying`
  - `StageClear`
  - `GameOver`
- Keep mode transitions data-driven and testable.

## Stage Data Format

The stage loader is in `i.java` lines 3407-3473.

Observed format:

1. Open `worldFiles[crtWorldIndex]`.
2. Read one initial byte.
3. Loop over packs/groups.
4. Read `worldStageCount`.
5. For each stage:
   - Read 4-byte header.
   - Header bytes 0-1: little-endian width.
   - Header bytes 2-3: little-endian height.
   - Read player layer: `width * height` bytes.
   - Read background layer: `width * height` bytes.
   - Read foreground layer: `width * height` bytes.
6. If not the selected stage, skip `width * height * 3`.

Layer arrays:

| Java array | Meaning inferred from loader | Go target |
| --- | --- | --- |
| `crtStagePlayerLayer[x][y]` | Actors/player-layer objects: player, enemies, movable objects, many item/object IDs. | Dynamic object layer and object spawn layer. |
| `bIntArrArr[x][y]` / `bByteArrArr[x][y]` | Background tile/state plus timers/flags. Loaded from the second layer and then expanded. | Static tile map plus per-cell state. |
| `crtStageForegrondLayer[x][y]` | Foreground/overlay/interactable foreground. | Foreground tile map, doors, checkpoints, visual overlays. |
| `cByteArrArr[x][y]` | Dirty/redraw/timing helper initialized after load. | Usually not needed in Go, except as inspiration for dirty updates if optimizing. |

Go implication:

- Write a decoder for `w0.bin` first, not a hand-authored guess.
- Export decoded stages to an inspectable format before gameplay work:
  - `tools/drdecode` reads `w0.bin`.
  - Output JSON while preserving all three raw layers.
  - Preserve all raw byte IDs until they are named.
- Do not map raw IDs into a guessed tile taxonomy; preserve the source values and layer roles.

## Decoded Stage Packs

`tools/drdecode` has verified the Java loader format against all three local world packs.

Command:

```bash
go run ./tools/drdecode \
  -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/w0.bin \
  -out decoded/world0
```

Decoded outputs:

| World pack | Output directory | Stage count | Stage dimensions |
| --- | --- | ---: | --- |
| `w0.bin` | `decoded/world0` | 14 | `26x21`, `27x24`, `27x26`, `40x23`, `51x24`, `30x75`, `26x45`, `44x29`, `35x14`, `50x30`, `50x31`, `46x31`, `46x31`, `68x11` |
| `w1.bin` | `decoded/world1` | 13 | `45x24`, `40x24`, `33x30`, `43x30`, `45x33`, `43x33`, `36x43`, `53x27`, `34x32`, `35x25`, `37x24`, `60x20`, `23x60` |
| `w2.bin` | `decoded/world2` | 14 | `60x26`, `39x35`, `35x42`, `38x38`, `46x27`, `49x28`, `55x30`, `51x20`, `51x37`, `104x17`, `35x25`, `45x28`, `35x32`, `35x26` |

Each decoded stage file includes:

- `width`
- `height`
- `player`, `background`, and `foreground` raw layer arrays
- per-layer byte histograms

The raw byte IDs are intentionally not named yet. The next step is ID taxonomy: compare decoded layers against the Java runtime and classify object/tile meanings.

### Raw-ID Preview Output

`tools/drinspect` renders decoded stages as diagnostic color maps. These PNG files are not original art. Each color is a stable visualization of one raw byte ID, so the output is useful for checking coordinates, layer composition, and ID frequency before semantic naming.

Command:

```bash
go run ./tools/drinspect -in decoded -out decoded/preview
```

Important outputs:

- `decoded/preview/index.html`: static browser index for all 41 decoded stages.
- `decoded/preview/world0/stage00.html`: per-cell HTML inspector for Angkor Stage 1.
- `decoded/preview/world0/stage00-contactsheet.png`: background, foreground, player, and composite raw-ID panels.
- `decoded/preview/id-index.md`: all raw IDs by layer, total count, and stage appearances.
- `decoded/preview/palette.png`: raw byte `0..255` color legend.

The HTML cells show raw byte IDs. Hovering a cell shows `x/y` and the `player`, `foreground`, and `background` IDs for that coordinate. Raw byte `255` is Java signed byte `-1`, which usually means empty/no object after loading.

### Sprite Extraction Output

`tools/drsprite` extracts the original Gameloft `.f` sprite chunks into PNG sheets and JSON frame/animation metadata.

Command:

```bash
go run ./tools/drsprite \
  -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources \
  -out decoded/sprites
```

Current result:

- `94` sprite chunks exported from `121` `.f` chunks.
- `decoded/sprites/index.html`: browser index for real sprite sheets.
- `decoded/sprites/manifest.json`: chunk offsets, lengths, parse errors, and sprite metadata.
- `decoded/sprites/0/chunk02-frames.png`: Angkor world tile/frame sheet used by stage rendering.
- `decoded/sprites/o/chunk00-frames.png`: hero frame sheet.
- `decoded/sprites/tutorial/demoSpr/`: the three embedded tutorial sprite records exported from `demoSpr.bin` with the same module/frame metadata format.

Tutorial sprite extraction is separate because `demoSpr.bin` is not an `.f` container:

```bash
go run ./tools/drsprite \
  -demo-sprites /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/demoSpr.bin \
  -out decoded/sprites/tutorial
```

Some `.f` chunks are not sprite chunks, or are image/text/sound-like payloads; those remain listed in `manifest.json` with parse errors instead of stopping export.

Raw stage IDs still do not directly equal sprite frame IDs. The Java renderer first transforms stage bytes into runtime layer state, then chooses `textures[...]` or `aClassfArr[...]` by foreground/player/background state. Example anchors:

- `i.java` lines 4179-4528 load world, generated, and common sprite chunks into `aClassfArr` and `textures`.
- `i.java` lines 6435-6479 map several texture slots to `aClassfArr` indices.
- `i.java` lines 6886-7070 draw foreground objects.
- `i.java` lines 7072-7335 draw many player-layer objects.
- `i.java` lines 7739-7800 draw high foreground IDs and late overlays.

The broader-world fidelity step remains a render-mapping table: raw ID + runtime state -> sprite file/chunk/frame/palette. The actually visible low player IDs used by all 14 packed Angkor stages now have explicit handling in `internal/originalgame`. In particular, quota markers use `cm.f` chunk `5`, crawlers use `gen1.f` chunk `4`, pressure switches use `gen2.f` chunk `9`, frozen violet/snake objects use `gen0.f` chunk `1` and `gen1.f` chunk `6`, Stage 6 uses all three `mm0.f` chunks, Stage 9 uses `b0.f` plus `mmv.f` chunk `3`, and the tutorial uses `demoSpr.bin`, `gen3.f` chunk `1`, and `mmv.f` chunk `0`; diagnostic blue blocks must not leak through hidden payload IDs.

## World Map And Economy

`a_Config.java` contains key economy constants:

- Shop jacket/armor prices in purple gems: `{150, 400, 1000, 3000}`.
- World unlock red-gem prices: `{0, 10, 25}`.
- Seal positions and movement graph for Angkor, Bavaria, Siberia, Shop.

`i.java` lines 2100-2205 process stage-clear rewards and completion marks:

- Violet gem total is read from `iByteArr[4..5]` and incremented by stage-collected value.
- Red gem total is read from `iByteArr[6..7]` and incremented by stage-collected value.
- World unlock checks compare red gems against `a_Config.worldPrices`.
- Completion flags are written per world/stage with bit values.
- Extra achievement-style marks are awarded for all violet gems, all red gems, no hit/restart-like counters, and similar stage objectives.

Go implication:

- Persist violet and red gem totals separately.
- Stage result must track more than "cleared":
  - cleared
  - all violet gems
  - all red gems
  - no damage
  - no restart/checkpoint recall
  - secret/alternate route state
- Shop upgrades must use violet gems; world unlocks must use red gems.
- The implemented global seal selector follows the four source positions Angkor/Bavaria/Siberia/Shop, the directional graph from `a_Config`, and red-diamond thresholds `10/25`. It preserves the relic fly-in/flash/effect and world-unlock flash before accepting input. Angkor enters `map_angkor.out`, Bavaria enters `map_scotland.out`, and Siberia/Shop remain explicit unavailable-content states rather than fake screens.
- Save version 5 introduced `RelicMask` and `WorldUnlocked[3]`; current version 6 retains them and adds per-stage consumed reward coordinates. A pre-v5 save with Stage 9 already cleared migrates the Angkor relic bit.

## Secret Exit Raw 28 And Map Branches

`stage06` is the first Angkor stage containing both a normal foreground raw `5` exit and a foreground raw `28` secret exit. The two IDs share rendering and auto-walk cadence but not completion behavior:

- Stage initialization stores raw `5` background direction in `pByte` and raw `28` direction in `qByte` (`i.java` around lines 3645-3654).
- Entering raw `5` sets `xBoolean=true, atBoolean=true`; raw `28` sets `xBoolean=true, atBoolean=false` (around lines 12160-12176).
- The auto-exit branch selects `pByte` or `qByte`, creates `jInt=18`, and continues until the hero is more than five cells outside the stage (around lines 10833-10852 and 10893-10915).
- A normal exit, or an exit from an already-secret stage, enters result state `35`. A raw `28` exit from a normal stage instead enters state `20`, persists collected coordinates/currency, calls `zVoid()`, and skips the result screen.
- Map loading records the smallest type-`1` stage index as `dKInt`; Angkor's decoded map makes it `9`. `zVoid()` follows a linked type-`1` node with a larger stage index and writes unlock bit `0x40`. Thus stage index `6` branches to `9`, index `7` branches to `12`, and secret indices `9 -> 10 -> 11` form their own chain.
- State `20` draws localized text `61`, `Congratulations! You have unlocked a secret path!`, and waits until `aSInt > 30` before returning to the map.

Go persists explicit per-node unlocks in save version 6. `HighestUnlocked` remains compatibility/summary data only: a jump to index `9` must not unlock normal indices `7` and `8`. Version-2 saves migrate their old sequential range once when loaded; pre-v4 saves migrate with the tutorial already completed, and pre-v5 Stage 9 clears migrate the Angkor relic bit.

## Freeze Hammer And Player Raw 9

Stage 8's secret path is the first implemented route that requires tool level `8`. The source behavior is split across the hammer action, object update, rendering, and conditional resource loading:

- `dInt(x,y)` around `i.java:10953` maps violet raw `1` to generated sprite type `34`, monkey raw `45` to `35`, wasp raw `46` to `36`, snake raw `19/43` to `37`, and turtle raw `49` to `39`.
- `aBoolean(x,y)` replaces a supported target with player raw `9`, writes the generated type into state bits `22..27`, and preserves red-snake identity in bit `0x10000000`. The enhanced hammer scan at `i.java:11449-11624` checks the impact cell plus its four orthogonal neighbors; a violet gem qualifies only at the impact cell, while an adjacent moving snake must overlap the hammer impact cell in source pixels.
- A direct raw `30` hit or raw `9` thaw sets the action-handled flag and suppresses the five-cell scan for that swing. Hammering raw `9` restores the source object; generated type `37` chooses raw `19` or `43` from the red flag, initializes direction `2` when the hero stands above it and `1` otherwise, and applies the source `0x78` stun timer before it resumes movement.
- Object dispatch around `i.java:14114-14166` sends raw `9` to `aqVoid()`. Therefore it shares source gravity, rolling, push, hook, pressure-switch, landing, hero damage, and contact-enemy crush behavior with rounded gravity objects. The frozen source type must follow every cell transfer and checkpoint snapshot.
- `OVoid()` around `i.java:8884-8931` renders raw `9` from `aClassfArr[(state & 0xFC00000) >> 22]`, frame `0`. Conditional loading around `i.java:4157` provides generated sprite `34` from `gen0.f` chunk `1` and generated sprite `37` from `gen1.f` chunk `6` when tool level is at least `8`.

The Stage 8 normal route does not require freezing: move the source boulder from `(33,12)` to `(31,15)`, then push it right so it falls onto pressure switch `(32,17)`. The secret route freezes the cycling snake at `(8,2)`, lets raw `9` settle at `(8,3)`, and pushes it left onto switch `(7,4)` to open foreground raw `7` at `(6,3)`. Both paths are replayed at 20 TPS in `runtime_stage07_test.go`.

## Great Anaconda And Angkor Seal

Angkor `stage08` selects special mode `kByte=4` at `i.java:3345-3358`. Its dedicated update is `i.java:9976-10193`, rendering is `LVoid()`/`MVoid()` around `i.java:8670-8805`, and helpers are `abVoid()` through `acVoid()` around `i.java:11889-11975`.

- Source fields initialize as `aoInt=0`, `aqInt=3`, `arInt=0`, with boulder shafts x `12/15`, spawn y `2`, and support y `5`. `dInt()` maps columns `0/1/2` to body x `10/13/16`.
- State `1` creates passable-damaging raw `50` at body x, y `8` only after `apInt>20`; state `2` consumes raw `0` only in the two body columns and rows `8/7`. A hit decrements health once even if multiple boulders are found.
- Retraction uses 10 ticks above one health and 5 ticks at one health. `acVoid()` regenerates a missing shaft boulder at y `2` when its y `5` cell is empty and sets `bjInt=30` camera shake.
- `LInt=28/50` controls the charge animations; state `11` checks one tail hit for exactly 12 ticks on hero row `4`. State `9` chooses the next body column from `(playerX-10)/3`, preserving the source tie rule based on `aSInt%50/25`.
- The third hit enters state `7` only after the 40-tick hurt interval. State `7` waits for `apInt>80`, then calls `alVoid()` to decrement group `0` and open the exit door.
- `b0.f` chunks `0/1` are the body/platform sprites. Body animation durations are `1,3,1,5,12,4,8,22,2`; the secondary `gen1.f` chunk `0` animator supplies tail animation `2`. `mmv.f` chunk `3` is the Angkor seal.
- `lVoid(53)` maps raw `53` to seal index `DInt=0`, therefore saved mask bit `0`; raw `51/52` map to bits `1/2`. It switches the hero to one-frame animation `47` for 42 ticks and sets `hBoolean`, which clears input until `aSInt>140`.
- Seal completion enters state `28`, performs exactly 11 loading calls, then enters global state `27`; it does not use the ordinary 12-step stage-result screen or its four medal awards.

The Go implementation is isolated in `runtime_anaconda.go`. Unit tests lock every phase threshold, raw `50`, boulder regeneration, tail duration, death delay, checkpoint reset, seal mapping/timer, and source sprites. `TestRuntimeStage08CanDefeatGreatAnacondaAndCollectSealAtSourceCadence` performs all three hits through actual pushes before collecting `(27,6)`.

## Stage Result Animation

The stage-clear flow is split between update states and rendering code rather than being a static result page:

- Gameplay enters `bByte = 35` with `bsInt = 0` and `brInt = 12` around `i.java` lines 10893-10915.
- State `35` advances `bsInt` through 12 result/save steps at lines 2101-2244. `KVoid()` at lines 8559-8573 draws the black `Loading` screen and a bottom progress bar whose width is `(bsInt + 1) * 230 / 12`.
- After loading, Java enters `bByte = 17`, resets `aRInt = 0`, and plays result sound `15`.
- State `17` at lines 2247-2290 advances result phase `aRInt`. Phase `0` lasts 40 visible ticks, phase `1` lasts `max(40, collectedViolet * 2)`, phase `2` lasts 40 ticks, and phases `3` and `4` last 10 ticks each. Phase `5` remains until Continue is accepted. OK can skip each phase.
- The result renderer in Java `case 17` uses background color `0x261707`. During phase `0`, the stage title at y `15` slides in with offset `min(-100 + tick * 10, 0)` and `Complete!` at y `32` uses the same offset delayed by 240 pixels.
- Violet diamonds enter in phase `1`; their displayed count is `min(tick >> 1, collectedViolet)`. Red diamonds, hits, and retries then enter one row at a time with the same `-100 + tick * 10` slide.
- Java row label/count y coordinates are violet `75/91`, red `131/147`, hits `187/203`, and retries `243/259`. Hero icons use y `191` and `243`; the prompt is at the source bottom anchor.
- All-violet, all-red, zero-hit, and zero-retry conditions reveal the `ui.f` chunk `4` medal at y `86/142/198/254`. A newly earned mark also draws the `cm.f` chunk `4` spark at y `69/125/181/237` and the `cm.f` chunk `7` effect animation at that row.
- Award effects intentionally use three different flat-sequence rules from JAR bytecode: all-violet uses `tick << 1` with a half-count bound, all-red/no-hits use `tick >> 1` with a doubled bound, and no-retries uses `tick << 1` with a doubled bound that can cross animation boundaries.
- Hit and retry row icons are the first animation frames of hero animations `10` and `12`, respectively; they are not generic text glyphs.

Go implication/current Angkor implementation:

- `internal/originalgame` now preserves the 12-step Loading transition, result sound `15`, source phase durations and threshold-plus-one visible frame, slide equations, count-up cadence, JAR row coordinates, achievement reveal order, and extracted medal/spark/effect assets.
- Accepting Continue in phase `5` enters the decoded `map_angkor.out` world map and selects the next unlocked node. Stage 7 raw `28` travels from node `6` to node `9`; Stage 8 raw `28` travels from node `7` to node `12`. Both show the 30-tick secret-path message and skip the result page. Secret nodes `9 -> 10 -> 11` use ordinary results before following their raw-`28` link; terminal nodes `11/12` remain selected. Stage 9 instead follows the seal-specific 11-step transition above.
- Per-stage clear state and four result award bits persist in `<config>/zskc-diamondrush/original-progress.json`; an award effect is new only when its bit was absent before that stage clear. Lives, maximum health, and tool level carry between stages.
- Result text and stage overlays use exported FreeJ2ME `SansSerif Bold` 10px/12px atlases rather than the earlier handcrafted bitmap font.

## Persistence

`i.java` uses Java ME RMS:

- `DiamondRush` record store: game progress (`tVoid`, `uVoid`, lines 5028-5060).
- `Preferences` record store: options (`vVoid`, lines 5078-5087).
- New-game initialization starts around lines 5105-5135.
- A fresh save sets `iByteArr[3] = 5` (extra lives) and `iByteArr[8] = 4` (maximum health). `xVoid()` copies the saved life count into `azInt`, while stage initialization copies maximum health into the active `nByte` health value.
- Fresh-save construction records every foreground raw `14/33` coordinate in each stage. `aVoid(world,stage,x,y)` removes a coordinate when a permanent reward is collected; stage initialization uses `aBoolean(...)` to leave an already-consumed ordinary chest open and empty. In boss modes, a consumed relic chest becomes a replayable raw `41` reward worth 10.
- World 0 removes coordinates for collected red diamonds, successfully awarded extra lives, and relics. Keys, healing, value-bearing violet rewards, tools, and the tutorial Compass remain replayable under their separate progression rules.
- State `35` adds every run's `aZInt` to the saved violet bank. Violet diamonds are intentionally farmable on replay. Red diamonds add the current run's `bbInt`, while their removed chest coordinates prevent the same red diamond from being awarded twice and allow partial collections to accumulate over multiple runs.

Go implication:

- Save version 6 stores the source-equivalent consumed-reward coordinate list per stage in the structured JSON progress file.
- Preserve the original semantics, not the byte layout.
- Keep a debug importer later if RMS-byte comparison becomes useful.

## Checkpoints

Checkpoint behavior is backed by two different concepts:

1. `KEY_OK` on the active checkpoint tile restores/reset-checkpoints the room/state. The branch starts at lines 1266-1271.
2. `KEY_RESPAWN` (`*`) triggers checkpoint recall/death-style reset. The branch starts at lines 1500-1526.

Moving onto foreground raw `4` also activates a checkpoint. The movement helper reads the checkpoint foreground cell state from the high bits/background-derived state, compares it with the checkpoint progress counter, sets `shouldSaveCheckpoint`, advances the counter, and then the main update calls `saveCheckpoint()` once the hero is idle. This means activation is automatic on entry and revisiting the same or earlier checkpoint must not overwrite the saved room state.

`saveCheckpoint()` starts at line 15190 and snapshots:

- health/life-related bytes.
- player tile position.
- collected counters.
- multiple layer arrays.
- world-specific boss/hazard state when applicable.

`restoreCheckpoint()` starts at line 15248 and restores:

- camera around player.
- player position.
- health and counters.
- layer arrays and world-specific state.

`XVoid()` starts at line 11018. Death/recall decrements `azInt` before restoring the checkpoint, and `restoreCheckpoint()` does not restore `azInt`; extra lives are a player resource rather than checkpoint-local mutable room state.

Go implication:

- Checkpoints must snapshot both player state and mutable room tiles/objects.
- Do not include extra-life count in checkpoint snapshots; collected lives persist across room reset/recall.
- Checkpoint activation order comes from the raw `4` cell state/background value.
- `Action5` on a checkpoint is not identical to dying.
- `Action5` and `Star` while standing on a checkpoint reset the checkpoint state without consuming a life.
- `Star` away from a checkpoint first plays hero animation `19` for its full 42 source ticks. Only when that animation ends does `XVoid()` consume one extra life, increment retries, restore the checkpoint, and refill health.

## Action Button, Hammer, And Hook

The core `KEY_OK` branch starts at line 1266.

Source and original-JAR observations:

- If player stands on the saved checkpoint foreground tile, `KEY_OK` plays checkpoint sound and calls `restoreCheckpoint()`.
- Otherwise, the game checks tool/progression bytes (`iByteArr[9]`) and scans nearby cells.
- Tool level `>= 2` scans horizontally up to 3 cells. A target must be 2 or 3 cells away; an adjacent candidate cancels the cast. Closed foreground raw `7`, an intervening non-empty object, foreground raw `14`/`33`, or raw `48` with state bit `0x8` blocks the applicable path.
- The hook candidate switch includes player-layer raw `0`, `1`, `8`, `9`, `11`, `14`, `19`, `43`, `47`, and `48`.
- JAR traces confirm that raw `32` rope segments extend with timer `18 -> 12 -> 6 -> 0`; because source object scanning is left-to-right, a new rightward second segment is revisited in the same scan and starts visibly at `12`, while a leftward segment starts at `18`.
- On impact, a physical target is pulled cell by cell until adjacent to the hero, not merely one cell. A raw `1` violet gem receives one additional pull into the hero cell and is collected only after the hook action completes. Reacquisition by a remaining rope segment resets the first pull motion to `0` in the same source scan.
- `ZVoid()` restores a saved target state only for raw `0`, `8`, `9`, and `47`, with pending direction `0x7000` and roll preparation `0x200` cleared. Other candidates receive `-1`. For a hooked snake, `eVoid()` interprets that negative value as packed target `(127,127)` and chooses the available primary/fallback axis; restoring every target's old state is not source behavior.
- Hero animations are right cast/pull `20/21` and left cast/pull `22/23`. Hooking locks input and damage and emits sound `12` on impact/reacquisition.
- Tool level `>= 1` checks the local hammer neighborhood. Hammer animations are up/right/down/left `13/14/15/16`, impact occurs on tick `3`, and duration is 11 ticks upward or 12 ticks otherwise.
- The branch chooses/fixes facing direction and action animation/state through `aByte` and `kInt`.

Known action sound IDs:

- `SOUND_SFX_HAMMER_HIT_UNBREAKABLE = 6`
- `SOUND_SFX_CHECKPOINT = 9`
- `SOUND_SFX_ENEMY_HURT = 10`
- `SOUND_SFX_BREAK = 11`
- `SOUND_SFX_HOOKING = 12`

Go implication:

- Implement `Action5` as a mode/tool resolver:
  1. Checkpoint special action if standing on active checkpoint.
  2. Hammer-like adjacent/local action if hammer acquired.
  3. Hook-like horizontal ray/line action if hook acquired.
  4. Menu/dialog confirmation if not in gameplay.
- Do not implement hook as mouse-targeted.
- Do not assume hook can only grab rocks. The source branch scans several object IDs; target taxonomy must be decoded from stage/object IDs and runtime observation.
- Preserve raw `32` as real player-layer occupancy during the action, its directional state bits, per-segment motion timer, target interpolation, and cleanup during checkpoint restore.

## Angkor Diggable Vegetation

- Player-layer raw `10` is not a permanent wall in Angkor. Movement case `10` permits entry while `xByte == 3` and marks the cell state active.
- Object update case `10` then clears the player-layer cell and writes foreground raw `32`.
- Raw `32`'s removal frame is the high byte of `crtStageForegrondLayer`; it is not stored in `bIntArrArr`. A gravity object may enter the same cell while that foreground animation is active, so the remake keeps `ForegroundState` separate from `ObjectState` and snapshots both arrays.
- Foreground raw `32` advances animation `0` from `aClassfArr[16]` and clears itself after the final frame.
- World resource loading maps `aClassfArr[16]` to `0.f` chunk `1`, extracted as `decoded/sprites/0/chunk01-frames.png`.
- Falling-object helper `eBoolean` permits foreground raw `32`; snake helper `fBoolean` blocks it until the removal animation has cleared.
- Java calls `ahVoid()` before applying the hero movement branch in the same 20 Hz update. Entering raw `10` marks it active at the end of that frame; the next frame's object scan replaces it with foreground raw `32` before reducing the hero's movement offset.
- The Go runtime preserves that one-scan delay. Foreground raw `32` begins at sequence frame `0`, advances on even source ticks, and clears when the eighth frame is reached.

## Health, Damage, And Lives

`hurtHero` starts at line 761. The death/respawn path `XVoid()` starts at line 11019.

Observed behavior:

- A new game starts with five extra lives and four health segments.
- Damage is not always instant death.
- The function ignores hits during several animation/state cases.
- It decrements `nByte` by damage amount.
- It increments a hit counter (`bcInt`), relevant to perfect-clear marks.
- It sets knockback/hurt animation flags and plays `SOUND_SFX_HERO_HURT`.
- Some call sites pass damage values `1`, `2`, or `iByteArr[8]`, showing variable damage severity.
- `XVoid()` decrements `azInt`; if the result is still non-negative it calls `restoreCheckpoint()`, restores health to `iByteArr[8]`, and clears transient hero state. Otherwise it enters the game-over path.

Go implication:

- Model health/energy separately from lives.
- On lethal damage, consume one extra life and restore the checkpoint with full health when available.
- Track damage taken for stage result marks.
- Add invulnerability/hurt state, not just immediate tile reset.

## Falling Objects And Hazards

The source falling-object state machine is in `aqVoid()` and is shared by player-layer raw `0` boulders and raw `1` violet gems:

- `ahVoid()` starts at line 12803 and scans the active player neighborhood from bottom to top, then left to right.
- `eBoolean()` at line 14309 defines cells an object can enter; `iBoolean()` at line 15160 identifies rounded gravity supports raw `0`, `1`, `8`, and `9`; `kVoid()` at line 15174 wakes the changed cell and its eight neighbors.
- `aqVoid()` starts at line 15406 and contains collection, crush, vertical fall, roll preparation, interpolation, landing, and active-state updates.
- `OVoid()` at line 8884 renders the direction byte and object-state offsets.
- Sound `SOUND_SFX_BOULDER = 14`.
- A directly unsupported object moves on that object frame. `aqVoid()` stores it in the destination cell with an `18`-pixel reverse offset; subsequent frames reduce that offset by `6`, producing `18 -> 12 -> 6 -> 0` with no extra generic release timer.
- Side rolling is different. It is considered only while the object below satisfies `iBoolean()` (raw `0`, `1`, `8`, or `9`) and both cells on one side are empty according to `eBoolean()`.
- Roll preparation sets state bit `0x200` and starts its byte offset at `1`. While the offset is below `6`, it increments only when `(aSInt & 3) == 0`; from `6` onward it increments every source frame. At `12`, Java transfers the object into the diagonal cell with a vertical reverse offset of `12`.
- With the Stage 1 initial offset of `1`, the global-frame phase makes the first diagonal transfer occur after 24 to 27 source frames (1.20 to 1.35 seconds at 20 TPS). Raw `0` and raw `1` use the same delay.
- The renderer's `OVoid()` branch uses the preparation offset for horizontal rocking, `offset * offset / 24` vertical displacement, and `-1 + aSInt % 3` jitter before the transfer. The same arc is applied to horizontal movement over a stationary rounded support when the source-side neighbor condition passes.
- Bavaria water uses three packed 9-bit sub-layers per cell (owner `0..2`, shape `3..6`, offset `7..8`) plus 15 flow and 15 source records. Gravity reverses upward when the packed cell is nonzero except for literal cell value `3`; `OVoid()` adds the source eight-frame water bob to raw `0/1/8/9/47` rendering.
- On the timer-zero landing frame, Java updates the `0x38` rotation, clears the low direction, and emits sound `14` for a vertically landing boulder. The `0x400/0x800` side marker survives that frame but is cleared on the next stationary update over a non-rounded support.

Go implication:

- Keep vertical falling and roll preparation as separate states; do not apply one fixed delay to every unsupported object.
- Preserve bottom-to-top object scan order, the global-frame phase, and the source `12 -> 6 -> 0` post-roll vertical offset.
- Preserve partial damage and crush behavior.
- Decode object IDs before naming tiles permanently.

## Crawling Enemy

Source anchors:

- Stage init handles player-layer raw `11` by setting object state to `16` when the background/state byte is `1`, then clearing that background/state byte to empty.
- The object update switch around line 14009 dispatches raw `11` to `amVoid()`.
- `amVoid()` uses low state bits as movement direction and bit `0x10` to choose the alternate initial wall side. It probes the forward and side cells to follow a wall around convex and concave corners, rather than merely reversing on a blocked target, and clears the object after its high death/phase bits reach `4`.
- In Bavaria, any nonzero packed water cell starts the crawler's `0x100 -> 0x400` death sequence before another wall-follow step can start.
- When the crawler overlaps or reaches the player, it calls `hurtHero(1,64,0)`.
- The render branch around line 7428 draws raw `11` from `aClassfArr[6]`, loaded from `gen1.f` chunk `4`, with frame/offset choices derived from state and direction.

Go implication:

- Treat raw `11` as a blocking contact enemy that damages the player.
- Preserve a per-object state byte; low bits carry direction and bit `0x10` affects direction inference.
- Preserve the source wall-follow probes, corner turns, bottom-to-top/left-to-right same-scan movement, and `18 -> 13 -> 8 -> 3 -> -2` timer cadence.
- Render normal phase with module `(aSInt >> 1) % 3`; phases `1..3` use modules `3..5`, and phase `>=4` is hidden. Preserve source direction/motion offsets instead of drawing a rectangle.

## Snake Objects

Source anchors:

- Stage init identifies player-layer raw `19` as green snake and raw `43` as red snake.
- The object update switch around line 14065 dispatches raw `43` to `eVoid((byte)43)` and raw `19` to `eVoid((byte)19)`.
- `eVoid(byte)` starts at line 14632. It uses low state bits as direction, checks `fBoolean` for target-cell movement, reverses direction when blocked, and calls `hurtHero(1, 48, direction)` on player overlap/contact.
- `gByteArr` at line 17722 maps Java directions `1..4`; snake movement uses the negative of that direction vector.

Go implication:

- Preserve the low direction and the pending direction in state bits `0x7000`. A blocked snake clears its low direction and counts `21 -> 18 -> ... -> 0` before probing the reverse direction.
- Do not mark a moved snake as processed for the rest of the scan. A rightward or upward move reaches a cell that Java scans later in the same bottom-to-top, left-to-right pass, reducing its timer from `21` to `18` immediately.
- Use 24-pixel overlap with both object and hero interpolation offsets, and apply `hurtHero(1,48,direction)` directional knockback.
- Green-snake animation selects sequence `(aSInt >> 1) % frameCount` directly; it does not honor the `.f` animation-frame duration values.

## Horizontal Hazards

Source anchors:

- Stage init flags player-layer raw `22` and `23` for object updates.
- The render branch around line 7515 draws raw `22`/`23` with `aClassfArr[12]` and separate orientation offsets.
- The object update branch around line 14169 handles raw `23`, and line 14175 handles raw `22`.
- Raw `23` checks same-row player positions from `x` leftward over the current animation reach; raw `22` checks from `x` rightward. On hit, it calls `hurtHero(1,64,0)`.

Go implication:

- Track an animation/reach phase for raw `22`/`23`.
- Apply damage only over the current reach, not the full row.
- Render the emitter from the decoded hazard frame and the reach from `gen1.f` chunk `0` flame animation `0`, with raw `23` horizontally flipped; no diagnostic rectangle remains for these World 0 IDs.
- Animation sequence `0` begins with an empty frame lasting 20 source ticks. This is the authored fully retracted phase, after which sequence indices `1..10`, `11..20`, and `21+` produce reach `1`, `2`, and `3`.
- Java copies the tile map from its scrolling background buffer before scanning dynamic cells. The dynamic scan uses relative x/y `-1..11`; reproducing it as interleaved per-cell background/object drawing clips right-moving boulders and raw `22` flames under later floor tiles.
- Foreground raw `20..25`, runtime dig-removal raw `32`, and foreground world frames `>=80` are drawn in the later foreground scan after the hero. Pressure-switch raw `6` is explicitly clipped to its own `24x24` cell. Raw `255` remains empty and must be checked before the `>=80` frame mapping.

## Bonus Value Pickups

Source anchors:

- The readable decompile at `i.java:3789` appears to clear player-layer raw `12`, but this is a confirmed decompiler/control-flow artifact: original-JAR runtime inspection retains raw `12`, and the render branch at `i.java:7271` draws it while the quota remains. Java stores its coordinate in `abInt/acInt` and its background/state byte in `aaInt`.
- The object update switch around line 14109 dispatches raw `41` to `lVoid(41)`.
- `lVoid(41)` stores `bIntArrArr[x][y]` into `aHInt`, clears the player-layer object, and starts pickup animation `aFInt = 2`.
- Stage initialization adds each raw `41` background value to `aYInt`, the same total used by the violet-diamond result row.
- The pickup animation branch around line 11230 adds `aHInt` to `aZInt` and subtracts it from `aaInt`.
- When raw `7` health refill is collected at full health, Java changes the object to raw `41`, sets `bIntArrArr[x][y] = 10`, and calls `lVoid(41)`.
- Raw `1` violet gem collection in `aqVoid()` increments `aZInt` and decrements `aaInt`; when `aaInt` reaches zero, Java clears the stored `abInt/acInt` target.
- The raw `12` render branch composes `textures[18]` from `cm.f` chunk `5` and overlays `aaInt` using HUD digits. The marker occupies and blocks its player-layer cell until cleared.

Go implication:

- Decode raw `41` as a value-bearing violet-diamond pickup using the same cell's background/state value. It must increase the right-hand HUD violet count, the stage result numerator, and the saved violet bank.
- Bavaria Stage 1 supplies a concrete cross-world regression: the raw `41` chest at `(25,14)` has authored background value `10`; opening it must change the lower-right violet HUD to include those 10 gems and persist the same value in the Bavaria stage best and global bank.
- Decode raw `12` as a visible blocking quota gate. Keep it in the player layer until the remaining value reaches zero, then clear it.
- Decrement the quota for raw `1` and raw `41` collection.
- Keep foreground raw `5` exit logic independent: the exit does not query `aaInt`, although authored routes may place raw `12` in front of it.
- Preserve the full-health raw `7` conversion to a `10`-value bonus before the chest animation starts, so the overhead reward art is the violet diamond rather than the health item. At 99 lives, raw `6` first falls back to raw `7` and then follows the same full-health conversion.
- Render raw `41` as the violet-gem icon plus its authored numeric value and `cm.f` chunk `7` effect animation `3`; add every run's collected violet value to `VioletGemBank`.

## Red Diamonds Raw 2

Source anchors:

- Object update dispatches player-layer raw `2` to `lVoid(2)`.
- `lVoid(2)` starts the common pickup flow with `aFInt = 3`, the red-diamond pickup texture slot.
- The pickup completion branch for `aIInt == 2` increments Java `bbInt`.
- Java menu text names this resource "Red diamonds", and world unlock checks compare saved red diamonds against `a_Config.worldPrices`.
- Angkor Stage 1 places its only raw `2` at `(19,2)` under foreground raw `33`. `nVoid(33)` opens that foreground container only after hero movement reaches `jInt <= 0` and starts hero animation `40` (or the shorter `48` variant for closely spaced pickups).
- Extracted hero animation `40` lasts 67 source ticks. Its reward branch runs on sequence index `13`, reached at tick `39`; foreground raw `33` advances from closed frame `0` to final open frame `3` independently.
- `lVoid(2)` sets pickup texture slot `aFInt=3`. `NVoid()` draws that red-diamond frame one tile above the hero after animation `40` passes sequence index `13`, beginning at tick `41` for this animation.
- On sequence index `13`, `cVoid(playerX, playerY-1, 0)` also starts common effect animation `0` from `aClassfArr[9]`, loaded from `cm.f` chunk `7`; the extracted animation lasts 10 source ticks.

Go implication:

- Keep the chest closed while the hero is still interpolating into the cell.
- Do not draw the player-layer payload through a closed foreground raw `14`/`33` container. Raw `2/4/5/6/7/24/26/27/40/41/42/51/52/53` becomes visible only through the overhead reward/effect branch; all 120 such payload cells in `stage00` through `stage13` are covered by a container.
- Lock movement during animation `40`, award the red diamond at tick `39`, and restore input only after tick `67`.
- Render foreground raw `33` from its merged/runtime state, not from the separate background byte `255`; treating `255` as a frame index makes the closed chest disappear.
- Treat player-layer raw `2` as a red diamond collectible, separate from violet gems and bonus quota.
- During the 67-tick pickup action, render hero animation `40`, the 10-tick `cm.f` chunk `7` pickup effect from reward tick `39`, and the stationary red-diamond frame above the hero once the animation advances past sequence index `13`.
- Preserve it in checkpoint snapshots as a stage-level collection counter.
- Add each run's newly reachable red diamonds to `RedDiamondBank` and remove their source chest coordinates. Save v6 therefore supports collecting different subsets on separate runs without allowing the same red chest to be farmed; the bank drives the Bavaria/Siberia thresholds `10/25`.

## Special Source Pickups

Several rare player-layer pickups carry the cross-world tool progression used by authored revisit routes.

Source anchors:

- Stage init groups raw `24`, `26`, and `27` together, increments pickup/object counters, calls `iVoid(x,y)`, and object update dispatches them to `mVoid(24/26/27)`.
- `mVoid()` only collects these objects during the hero action/pickup animation, clears the player-layer cell, and sets `iByteArr[9]` to `1` for raw `24`, `8` for raw `26`, and `2` for raw `27`.
- The pickup animation branch maps raw `24` to `bmInt=22`, raw `27` to `bmInt=23`, and raw `26` to `bmInt=25`.
- Decoded source data places raw `24` in Angkor Stage 4 at `(26,18)`, raw `27` in Bavaria Stage 3 (`world1/stage02`) at `(24,25)`, and raw `26` in Siberia Stage 6. Each is under foreground raw `14`; Angkor Stage 5 contains no raw `27` chest.
- The `KEY_OK` action branch gates hammer/local tool checks behind `iByteArr[9] >= 1` and hook scanning behind `iByteArr[9] >= 2`.
- Stage init for raw `42` increments the same object/pickup counters, sets resource bits `0x80000000 | 0x40000000`, and calls `iVoid(x,y)`. Object update dispatches raw `42` to `lVoid(42)`, which sets `aFInt=29/aGInt=0`; the pickup animation sets `pBoolean=true` and `bmInt=11`.
- Object update dispatches raw `53` to `lVoid(53)`. The pickup animation stores `DInt=0`, sets bit `1 << DInt` in `iByteArr[2]`, calls `uVoid()`, and clears the object through the common pickup flow.

Go implication:

- Treat raw `24`, `26`, `27`, `42`, and `53` as passable source-special pickups. Raw `42` is the Compass and triggers tutorial script `11`; raw `53` preserves the special animation `47`, 141-tick input lock, seal bit `0`, and 11-step non-result transition.
- Preserve the Java state effects as runtime flags/masks and derive action tool level from raw `24` Mystic Mallet, raw `27` Mystic Hook, and raw `26` Freeze Hammer in source progression order `1/2/8`.
- Clear the source cell and include the flags in checkpoint snapshots.
- Draw all three tools from their source modules in overhead rewards. Bavaria Stage 3 now performs the real raw `27` acquisition and persists tool level `2`; Bavaria Stage 5 inherits it and contains no duplicate raw `27`. The Angkor Stage 5 direct-stage fallback supplies the same authored revisit prerequisite without inventing a false hook chest.

## Angkor Tutorial Stage 13 And Compass

`stage13` is the final entry in `w0.bin`, but `map_angkor.out` has only the 13 nodes for stages `0..12`. It is a pre-map demo/tutorial stage controlled by `demo.f` and `c.java` rather than an ordinary result route.

- Foreground event backgrounds dispatch scripts `29`, `10`, `13`, `16`, and `28`; checkpoint reset/recall dispatch scripts `15` and `17`; collecting raw `42` dispatches script `11`.
- The source order is `29 -> 10 -> 11 -> 13 -> 15 -> 16 -> 17 -> 28`. Prompt index sequences are `[12,13,14]`, `[0]`, `[19]`, `[1,2,3,4]`, none, `[5,6,7]`, `[8]`, and `[9,10,11]`.
- `c.java` commands move the camera and hero, write foreground cells, wait, show text, place/change/hide the portrait, animate punctuation, and flash white. Skipping suppresses presentational waits but must still execute movement and map writes.
- The portrait face, punctuation, and base records are sprite IDs `0`, `1`, and `2` in `demoSpr.bin`. They expose 9 face frames, 9 punctuation frames, and 2 base frames after extraction.
- Global generated chunk index `31` is `gen3.f` chunk `1`; the resource loader stores its one `24x24` module as `textures[29][0]`. `NVoid()` draws that module one tile above the hero for raw `42`. This is distinct from the HUD Compass direction, which uses `ui.f` chunk `2` frame `3+direction`.
- At tutorial event `(46,7)`, source flag `tBoolean` enables `aClassfArr[15]`, loaded from global generated chunk `30` (`gen3.f` chunk `0`). `LVoid()` draws animation `0`, sequence `(aSInt >> 1) & 3`, one tile above the interpolated hero position only while no demo script is active. It is the blue `*` recall-key cue, not a moving checkpoint. The desktop renderer preserves that timing, position, and visibility rule but substitutes an `ENTER` keycap so the displayed instruction matches the actual binding.
- The final 5x5 seal uses `mmv.f` chunk `0` frames `4..28`, with frame `1` as the activation overlay. At `(61,3)` and exactly `jInt == 6`, Java queues common effect animation `5`, sets `eBoolean=false` to hide the hero, but continues drawing overlay frame `1` while the hero remains on `(60,3)/(61,3)`. Script `28` then ends and transitions a fresh save into `stage00`.

Go keeps the script state in `internal/original/runtime_tutorial.go`, presentation in `internal/originalgame/tutorial.go`, and persists completion in progress version 6. `runtime_stage13_test.go` replays the complete tutorial route and asserts each prompt sequence, reset, life deduction, layer write, portrait start state, and final seal position. The same interpreter also runs Stage 3 script `30` (`17/18`), Stage 4 script `22` (`20/21/22`), and Stage 9 script `33` (`32/33`); `tools/drdemo` exports all decoded commands to `decoded/demo-scripts.json`.

## Passable Overlay Raw 33

Source anchors:

- Stage init preserves player-layer raw `33` unchanged, just like player-layer raw `31`.
- The movement/collision helper groups raw `33` with passable object IDs such as gems, keys, snakes, and special pickups.
- Unlike pickups, raw `33` has no object update dispatch and no pickup branch that clears the player-layer cell.
- The render switch at line 6905 groups **foreground** raw `14`/`33`; foreground raw `33` draws from `aClassfArr[22]` using its high state. It is not a render branch for player-layer raw `33`.
- World 0 raw `33` appears five times, often on foreground raw `7` door cells or near checkpoint/foreground triggers.

Go implication:

- Treat raw `33` as a passable overlay, not a collectible and not a hard block.
- Preserve it in the player layer when the hero moves onto its cell.
- Do not draw a fallback object sprite for player-layer raw `33`; render the owning foreground raw `33` overlay where present.

## One-Shot Foreground Event Raw 0

Source anchors:

- Stage init stores foreground raw `0` as `(background << 8) | 0` and counts it with other special foreground events.
- The foreground update branch for raw `0` checks whether the player is standing on the cell and the movement timer is near rest.
- When triggered, Java stores the merged foreground background/state in `bmInt` and clears the foreground cell.

Go implication:

- Treat foreground raw `0` as passable.
- On player entry, record the decoded background/state as the last foreground event and clear the foreground cell.
- Dispatch known World 0 values through the decoded demo interpreter: Stage 3 script `30`, Stage 6 camera event `3`, Stage 9 script `33`, and the tutorial event chain. Preserve `LastForegroundEvent` for deterministic verification.

## Clearable Foreground Cluster Raw 1

Source anchors:

- The movement helper checks foreground raw `1` before player-layer collision resolution.
- When the player enters a raw `1` cell, Java sets short-lived state counters and calls `bVoid(x, y, (byte) 1)`.
- `bVoid(x, y, (byte) 1)` recursively clears connected foreground raw `1` cells.

Go implication:

- Treat foreground raw `1` as passable.
- On player entry, recursively clear the connected raw `1` foreground blob.
- In Stage 6, reset the 120-tick falling-torch warning and advance the authored trigger count when each connected raw `1` region clears; the third trigger starts the decoded torch/fire animations.

## Special Foreground Barrier Raw 2

Source anchors:

- Stage init stores foreground raw `2` as `(background << 8) | 2` and marks the special-item foreground feature set.
- Original-JAR movement permits the hero to enter a raw `2` cell. The current-cell interaction branches on its background/state and the Java special-item byte `iByteArr[9]`; the readable decompile's blocking-looking control flow is not the observed result.
- State `0` checks for tool level `>= 1`; in World 0 this is paired with adjacent raw `30` breakable-wall gates.
- State `1` checks for tool level `>= 2`; the `KEY_OK` branch can then clear the connected raw `2` blob with `bVoid(x, y, (byte) 2)`.
- The foreground update branch for raw `2` with background/state `0` checks the four adjacent player-layer cells. If none is player-layer raw `30`, Java calls `bVoid(x, y, (byte) 2)`.
- `bVoid(x, y, (byte) 2)` recursively clears connected foreground raw `2` cells.

Go implication:

- Treat foreground raw `2` as player-passable and expose its hammer/hook prompt only while the hero occupies that cell.
- When a raw `2` cell with state `0` has no adjacent raw `30` breakable wall, recursively clear its connected raw `2` blob.
- Allow Action/5 to clear a state `1` raw `2` blob only after tool level `>= 2`.
- State `0` uses the hammer availability/module and state `1` uses the hook availability/module; render the current-cell prompt through the source tool asset rather than explanatory text.

## Foreground Doors Raw 7

Source anchors:

- Stage init records foreground raw `7` coordinates in `crtStageDoorXs/Ys` when the decoded background/state byte is not `-1`.
- Stage init stores the door state as `(background << 8) | 7` in the foreground layer.
- The movement helper checks foreground raw `7`; if `((state & 0xF0) >> 4) < 2`, it treats the door as blocking.
- Decoded background/state `255` maps to Java signed `-1`, so it behaves like an already-open/passable door state under that high-nibble check.
- Other routines, such as `doorHeadClose`, mutate the high state nibble to animate and close/open doors.
- Foreground raw `6` is a pressure-style door switch. Java reads its background/state as `doorI`; while the switch is pressed or occupied, it calls `hVoid(doorI)`, otherwise it calls `doorHeadClose(doorI)`.

Go implication:

- Keep foreground raw `7` as stateful, using the decoded background/state byte until the full merged foreground state is modeled.
- Block movement through raw `7` while the high-state nibble is less than `2`.
- Allow raw `7` cells with state `255` or high nibble `>=2`.
- During stage initialization, count all same-group foreground raw `6/8/9` activators and store that remaining count in each linked raw `7` door's low nibble. Each pressure/key activation decrements it; opening begins only when it reaches zero. Animation starts at phase `0x10`, advances every third source tick through passable `0x20` to `0x30`, and closing preserves the remaining activator count when the player is not in the door cell.
- Draw foreground raw `6` from `gen2.f` chunk `9` as a bottom-anchored `24x13` module, offset by the source player/object interpolation while it is depressed.

## Enemy Gate Triggers Raw 17

Source anchors:

- Stage init handles foreground raw `17` by reading its background/state byte as a group id.
- If the cell above raw `17` is an enemy (`19`, `36`, `43`, `45`, `46`, or `49`), Java increments a per-group counter, removes that raw `17` marker from the foreground layer, and records group metadata.
- Remaining raw `17` cells are kept as group markers. When a group counter reaches zero, Java scans raw `17` cells with that group id.
- Foreground raw `26` is the trigger switch for this mechanism. When the player stands on it, Java sets `cmInt` to the raw `26` background/state group, closes the door in front of the hero, validates the group index, plays the trigger sound, and clears the raw `26` cell.
- If a same-group raw `17` has foreground raw `7` immediately above it, Java starts that door opening by mutating the door state.
- If a same-group raw `17` has foreground raw `14` or `33` above it, that container is the authored group reward. It remains closed and its player-layer payload remains hidden/inert before the counter reaches zero; the completion scan writes the unlocked overlay state without replacing the payload background.

Go implication:

- Initialize per-group counters from raw `17` markers under enemies.
- Preserve the initialization metadata for diagnostics, but mirror Java `alVoid()`: an enemy/container removal decrements the single currently active `cmInt` group, regardless of where that enemy moved from.
- Treat raw `26` as a passable trigger switch that activates the current enemy-gate group and clears itself.
- Decrement only the currently active raw `26` trigger group. This makes Stage 11 trigger order part of the puzzle and intentionally does not retain per-enemy ownership after movement.
- When a group counter reaches zero, start same-group raw `7` doors at phase `0x10`; the shared door animation advances them to passable phase `0x20` on the next eligible third source tick.
- When a group counter reaches zero, preserve same-group foreground raw `14`/`33` overlays, clear their high state to `0`, and allow a hero already standing on the cell to start the normal chest sequence on the next object scan.
- Stage 11 exercises all four variants: groups `0/2` freeze and reposition snakes under hooked/pushed boulders; group `1` uses the two drop columns while temporarily freezing the other snake at a supported endpoint; group `3` stacks two hooked live snakes over a bottom ice block, freezes them together, then uses cascading thaws and the authored boulder. Its full route follows first-trigger order `1 -> 3 -> 0 -> 2` so Java's single active `cmInt` matches each kill.

## Sprite And Resource Handling

`f_Sprite.parseSprite` starts at line 60.

Important resource concepts:

- Modules: raw pixel blocks.
- Frames: module composition.
- Animations: ordered frame/time sequences.
- Palettes.
- Packed data formats include `I256`, `I16`, `I4`, `I2`, `I256RLE`, `I127RLE`.

Go implication:

- Do not load `.f` files directly at runtime unless needed.
- Preferred path:
  1. Write an extractor tool.
  2. Export PNG atlases and JSON animation metadata.
  3. Use those exported assets from Go.
- For generated GPT assets, keep proportions and animation categories aligned with the source.

## Sound Events

`j_SoundManager.java` defines 21 sound IDs at lines 15-35.

`tools/drsound` extracts all 21 standard-MIDI payloads from `snd.f` into `decoded/audio`. `internal/originalgame` loads the full bank, applies the original priority groups and 50ms equal-priority replacement guard, and uses `AVMIDIPlayer` on macOS. Stage 1 currently wires IDs `2`, `3`, `4`, `5`, `9`, `14`, `15`, and `16` at their source events; later IDs are retained for the stages/screens that use them.

| ID | Java name | Remake event |
| --- | --- | --- |
| 0 | `SOUND_SFX_SWITCH` | Switch / lever. |
| 1 | `SOUND_SFX_RIDDLE` | Puzzle/riddle event. |
| 2 | `SOUND_SFX_DEATH` | Death / recall. |
| 3 | `SOUND_SFX_CHEST_1` | Chest open variant 1. |
| 4 | `SOUND_SFX_CHEST_2` | Chest open variant 2. |
| 5 | `SOUND_SFX_HERO_HURT` | Player damaged. |
| 6 | `SOUND_SFX_HAMMER_HIT_UNBREAKABLE` | Hammer fails on target. |
| 7 | `SOUND_SFX_MINE` | Mine/explosion. |
| 8 | `SOUND_SFX_WORKING` | Tool/mechanism working. |
| 9 | `SOUND_SFX_CHECKPOINT` | Checkpoint. |
| 10 | `SOUND_SFX_ENEMY_HURT` | Enemy damaged. |
| 11 | `SOUND_SFX_BREAK` | Breakable block/object. |
| 12 | `SOUND_SFX_HOOKING` | Hook action. |
| 13 | `SOUND_SFX_WATER` | Water. |
| 14 | `SOUND_SFX_BOULDER` | Boulder/falling rock. |
| 15 | `SOUND_M_LEVEL_CLEAR` | Stage clear music. |
| 16 | `SOUND_M_ANGKOR_WAT` | Angkor music. |
| 17 | `SOUND_M_BAVARIA` | Bavaria music. |
| 18 | `SOUND_M_SIBERIA` | Siberia music. |
| 19 | `SOUND_M_TITLE` | Title music. |
| 20 | `SOUND_M_GAME_OVER` | Game over music. |

## Startup Menu And Save Selection

- `i.java:736-755` builds the source main-menu entries. With the disabled store link removed, the first two actions are text `16` (`Continue`) and text `40` (`New game`).
- `i.java:8131-8291` draws that menu over the three images in `spl.f`, uses `ui.f` chunk `3` frame `4` for the horizontal menu strip, frame `2` for both selection markers, and source color `0xCE9B00` for the selected row.
- `lang.xx` text `21` is `Starting a new game deletes your progress. Are you sure?`; source menu `7` contains text `24/25` (`No`/`Yes`). An existing save is not replaced until `Yes` is accepted.
- The Go runtime extracts the three lossless title images to `decoded/sprites/splash/`. Normal startup now enters this menu; `ORIGINALRUSH_STAGE` remains an explicit developer override that bypasses it.
- The JSON progress format has no in-stage snapshot. `Continue` therefore returns an unfinished save to the intro tutorial, or a completed-tutorial save to its highest explicitly unlocked Angkor map node.

## Go Display Cadence Boundary

- Java gameplay remains authoritative at `20 TPS`. `internal/originalgame.Game.Update` runs at `60 TPS`, captures desktop input, and calls the complete source update once every three adapter updates; it does not run only selected Runtime methods at the higher rate.
- One-shot Action, Recall, Skip and map-direction edges are consumed exactly once at the next source boundary. If the corresponding source state rejects input on that tick, the edge is discarded rather than replayed after an animation or script unlocks input. Held movement is sampled continuously and supplied to each source step.
- Runtime `ObjectMotion` remains packed in source values such as player/gravity `18 -> 12 -> 6 -> 0`, snake `21 -> 18 -> ... -> 0`, and crawler five-pixel decrements. `renderPhase=0..2` derives temporary intermediate offsets for drawing; those offsets never mutate layer, timer or collision state.
- `cameraX/cameraY` remain the source camera passed to `Runtime.SetViewport`. Normal following and tutorial/enemy-gate/foreground demo pans predict the next source camera only for 60 Hz display interpolation.
- Bavaria spike display extent similarly interpolates between the current and next source phase, while `SpikeSlowExtent/SpikeFastExtent` remain the only values used for collision.
- The renderer caches its one-pixel solid primitive and all sprite frame/module `SubImage` handles during loading. Normal cell and UI drawing therefore submits cached images instead of allocating Ebitengine atlas images per frame.
- `TestUpdateRunsOneSourceFramePerThreeRenderUpdates` checks the scheduler; `TestSixtyHzSchedulerPreservesTwentyHzRuntimeState` compares paced and direct runtime layers; `TestRenderInterpolationDoesNotMutateSourceMotion` verifies that intermediate motion is draw-only.

## Current Go Package Mapping

| Go package/path | Responsibilities | Java reference |
| --- | --- | --- |
| `cmd/originalrush` | Ebitengine startup. | `GloftDIRU`, `FreeJ2ME` only conceptually. |
| `internal/original` | Decoded stage loading, mutable layers, deterministic gameplay state, player/entities/tools/checkpoints. | Stage loader, gameplay branches and layer arrays in `i.java`. |
| `internal/originalgame` | Desktop phone-key adapter, rendering, camera/HUD, audio, world map, results and JSON progress. | `paint(Graphics)`, input handlers, `f_Sprite`, `b_SpriteAnimator`, `j_SoundManager`, RMS helpers. |
| `decoded/` | Generated world JSON, map graph, sprite metadata/images, fonts and audio. | `w*.bin`, `map_*.out`, `.f` resources and FreeJ2ME font output. |
| `tools/drdecode` | Decode world/stage files. | Stage loader lines 3407-3473. |
| `tools/drsprite` | Decode `.f` sprite resources. | `f_Sprite`. |
| `tools/drinspect`, `tools/drworldaudit` | Inspect raw IDs and generate source-data inventories. | Decoded stage layers. |

## Implementation Policy

Keep source-derived gameplay in `internal/original` and host/rendering concerns in `internal/originalgame`. Any mechanism that conflicts with the original JAR, Java source or a repeatable FreeJ2ME observation must be corrected and covered by a focused regression.
