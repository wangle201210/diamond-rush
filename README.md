# Diamond Rush Go Prototype

Go + Ebitengine + Tiled + self-contained tile/world logic prototype for a Diamond Rush style cave puzzle game.

## Source-Fidelity Rebuild Notes

The current Go code is a prototype and may be discarded or rewritten to match the original game more closely. The local Java reference project is:

```text
/Users/wanna/mine/github/wangle201210/DiamondRushSource
```

Read these documents before implementing new gameplay:

- `docs/diamond-rush-source-mapping.md` maps the Java decompilation, resources, input, stage loader, save data, and sound IDs to the Go remake.
- `docs/diamond-rush-remake-spec.md` defines the source-led rebuild plan, acceptance criteria, and tooling status.
- `docs/diamond-rush-original-gameplay.md` keeps the broader gameplay research notes.

Decoded original data and inspection outputs:

- `decoded/world0`, `decoded/world1`, `decoded/world2`: decoded `w0.bin`, `w1.bin`, and `w2.bin` stage JSON.
- `decoded/preview/index.html`: raw-ID stage inspector. Its PNGs are diagnostic color maps, not original art.
- `decoded/sprites/index.html`: extracted original `.f` sprite/frame sheets.
- `decoded/fonts`: deterministic FreeJ2ME 10px/12px font atlases and metrics used by the source-style overlays.
- `decoded/audio`: all 21 original MIDI payloads extracted from `snd.f` plus their manifest.
- `docs/angkor-world0-inventory.md`: implementation inventory for the original Angkor world data.

Useful tooling:

```bash
go run ./tools/drdecode -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/w0.bin -out decoded/world0
go run ./tools/drinspect -in decoded -out decoded/preview
go run ./tools/drsprite -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources -out decoded/sprites
go run ./tools/drsound -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/snd.f -out decoded/audio
go run ./tools/drworldaudit -in decoded/world0 -out docs/angkor-world0-inventory.md
javac tools/drfont/ExportFont.java && java -cp tools/drfont ExportFont decoded/fonts
```

## Run

```bash
go run ./cmd/diamondrush
```

Original-data remake debug entry:

```bash
go run ./cmd/originalrush
```

This entry currently opens only Angkor Stage 1 from the decoded original data. It uses the original `240x320` composition (`40px` top HUD, `240px` scrolling playfield, `40px` bottom HUD), `24px` tiles, extracted source sprites/modules, FreeJ2ME font metrics, original MIDI bank, and the source `20 TPS` movement/object cadence. Arrow keys or `2/4/6/8` move, `5`/Enter performs the context action, and `R`, Backspace, number-row `Shift+8` (`*`), or numpad `*` recalls the active checkpoint. Reaching the Stage 1 exit runs the source Loading/result sequence and holds its completed page instead of entering unfinished Stage 2/world-map flow.

The Stage 1 slice includes the source raw `79` automatic entrance and temporary closing door, checkpoints and room restore, vegetation timing, boulder/gem gravity and rolling, green snakes, horizontal fire, chest reward sequence, death/recall, source sound events, and persistent first-stage award history.

`cmd/originalrush` controls:

- Move: arrow keys or number row/numpad `2/4/6/8`
- Action/checkpoint reset: number row/numpad `5` or Enter
- Checkpoint recall: `R`, Backspace, number-row `Shift+8` (`*`), or numpad `*`
- Quit: Esc

Legacy `cmd/diamondrush` controls:

- Move: number row or numpad 2/4/6/8, arrow keys, or the on-screen D-pad
- Action: number row or numpad 5, Shift, or the on-screen center 5; checkpoint tiles reset the room, Mystic Hammer breaks cracked walls/stuns enemies, Mystic Hook pulls objects when available
- Checkpoint recall: Backspace or numpad `*`, costs one life
- Touch/mouse: on-screen phone-style D-pad in the lower-right HUD
- Restart: number row or numpad 9
- Pause: Esc
- Toggle audio: number row or numpad 0
- Title menu: 2/8 chooses Continue, New Game, Level Map, Options, Help, About, or Exit; 5/Enter selects
- Level select: 4/6 or left/right chooses unlocked levels, 5 or Enter starts, 1 buys HP, 3 buys armor, 7 buys lives, Esc returns
- Paused: 5/Enter/Esc resumes, 9 retries, Backspace returns to level select

## Current Mechanics

- Tiled `.tmx` level loading from `assets/levels/`
- Five polished levels with next/previous level controls
- Non-gameplay temple title screen with phone-style Continue/New Game/Options/Help/About/Exit menu and Angkor-style world-map level select flow
- Persistent unlocked-level progress and best-step records
- Persistent best-score records shown in level select and stage-clear results
- Red-diamond level gate: the final stage requires at least one saved red diamond or the discovered secret route
- Ancient-seal meta goal: the bundled pack opens the seal after the final stage is cleared with three total saved red diamonds
- Level titles, themes, and hints read from Tiled properties
- Per-level par steps and three-star ratings
- Score system with diamond, key, door, enemy, and par-step bonus points
- Treasure chests with Tiled-authored red-diamond, purple-diamond, key, potion, extra-life, tool, and score rewards
- Health, lives, checkpoint recall with a life cost, and checkpoint room reset with Action/5
- Compass HUD pointing toward the next inactive checkpoint
- Persistent Compass, Mystic Hammer, and Mystic Hook pickups placed in Tiled levels
- Secret exit tile with persistent per-level discovery marks
- Secret-exit route targeting for the bundled Angkor map, including a bypass around the final red-diamond seal
- Ancient-seal status shown on the title, world map, and final clear overlay
- Hidden wall tiles that reveal secret passages when entered
- Persistent no-damage, no-checkpoint-recall, and no-restart completion marks
- Persistent all-purple/all-red collection completion marks
- Boss-like guardian encounter that gates the final exit
- Purple diamond bank and a level-select shop for max-HP, armor, and lives upgrades
- Per-level gravity and enemy timing tuned from Tiled properties
- Per-enemy initial patrol directions tuned from Tiled object properties
- Stage-clear and try-again result overlays with gem, secret-route, all-collection, and clean-clear marks
- Pause overlay with resume, retry, and level-select actions
- Generated event sound effects
- Generated looping background music with mute toggle
- Nokia-style on-screen keypad controls plus desktop keyboard controls
- Phone-style scrolling viewport that follows the player
- Runtime sprite animation frames for player, enemies, and diamonds
- GPT-generated pixel art menu background
- Tile-grid player movement
- Interpolated player movement animation
- Dirt digging
- Diamond collection
- Required diamond count opening the exit
- Red diamonds collected from treasure chests, saved per stage, and used for late-stage entry gates
- Chests can contain red diamonds, purple diamonds, keys, potions, extra lives, tools, or score
- Checkpoint tiles that save a return point, support recall, and restore the room state captured at activation
- Health damage before life loss, with heavy crushes consuming remaining health
- Health potions restore HP and are placed in later stages
- Purple diamonds earned at stage clear can buy max-HP, armor, and starting-lives upgrades
- Horizontal boulder pushing
- Mystic Hook pickup and clear-line pulling for rocks, diamonds, and keys up to four tiles away
- Mystic Hammer pickup and strikes for cracked walls and adjacent enemies
- Boulder and diamond gravity with simple rolling
- Falling objects can crush the player and enemies
- Silver keys/doors and gold keys/doors
- Static spike traps, timed spike traps, and fire traps
- Switches that open bridge/gate tiles
- Cracked walls broken by falling rocks or diamonds
- Teleporters
- Lava hazards that burn players and falling objects
- Secret exits for hidden-route completion tracking
- Clean-run tracking for no-damage, no-checkpoint-recall, and no-restart clears
- All-gem tracking for purple and red diamond completion
- Horizontal patrol, vertical patrol, and chasing enemies
- Final guardian boss damaged by Mystic Hammer and falling rocks
- Restart and win/loss state
- GPT-generated pixel art source in `assets/sprites/generated_sources/`
- Runtime tileset in `assets/sprites/tileset.png`

## Level Curve

The bundled pack is intentionally limited to five stages and tuned like an early Diamond Rush world:

| Level | Title | Par | Focus |
| --- | --- | --- | --- |
| 1 | Angkor Gate | 85 | Collect diamonds, dig dirt, open the exit |
| 2 | Rolling Stones | 115 | Falling rocks and horizontal/vertical patrol enemies |
| 3 | Locked Shrine | 130 | Keys, locked doors, spikes, and tighter enemy corridors |
| 4 | Hidden Passage | 145 | Switches, bridge gates, cracked walls, teleporters, and lava |
| 5 | Temple Trial | 165 | Combined puzzle using all major mechanics plus a chasing enemy |

The test suite audits the pack's progression constraints: early levels do not introduce later mechanics, each level has exactly one exit, required diamonds do not exceed available diamonds, teleporters are paired, keyed levels include matching keys and doors, and the final level includes all enemy types.

## Verification

Use these commands for the current repeatable checks:

```bash
go test ./internal/world ./internal/level
go test -c ./internal/game -o /tmp/diamondrush-game.test && rm -f /tmp/diamondrush-game.test
go build -o /tmp/diamondrush-smoke ./cmd/diamondrush && rm -f /tmp/diamondrush-smoke
```

In macOS headless or no-primary-monitor contexts, `go test ./internal/game` can panic inside Ebitengine's GLFW monitor initialization before project tests run. Treat the commands above as the stable verification path unless the test runner has a real UI context.

## Project Layout

```text
cmd/diamondrush/      Ebitengine executable
internal/game/        Input, rendering, window, HUD, asset loading
internal/level/       Tiled TMX parsing
internal/world/       Pure tile/world gameplay logic
assets/levels/        Tiled maps and tileset metadata
assets/sprites/       Runtime sprites and generated source artwork
```

## Asset Notes

The current sprite source and menu background were generated with image generation. The sprite source was cropped into the 32px runtime sheet. The game loads `assets/sprites/tileset.png` and `assets/sprites/menu-background.png`; if the tileset is missing or invalid, it falls back to programmatic placeholder sprites.

Progress is saved as JSON under the OS user config directory:

```text
<config-dir>/zskc-diamondrush/progress.json
```

The source-data Stage 1 runtime keeps its isolated completion/award record at:

```text
<config-dir>/zskc-diamondrush/original-progress.json
```

Tileset GID mapping:

| GID | Tile |
| --- | --- |
| 0 | Empty |
| 1 | Wall |
| 2 | Dirt |
| 3 | Diamond |
| 4 | Rock |
| 5 | Closed exit |
| 6 | Key |
| 7 | Door |
| 8 | Open exit |
| 9 | Spike trap |
| 10 | Switch |
| 11 | Bridge/gate |
| 12 | Cracked wall |
| 13 | Teleporter |
| 14 | Lava |
| 15 | Treasure chest |
| 16 | Checkpoint |
| 17 | Health potion |
| 18 | Mystic Hammer pickup |
| 19 | Mystic Hook pickup |
| 20 | Compass pickup |
| 21 | Secret exit |
| 22 | Hidden wall / false wall |
| 23 | Timed spike trap |
| 24 | Fire trap |
| 25 | Gold key |
| 26 | Gold door |

## Remaining Fidelity Work

- Tune these five stages through repeated playtesting
- Tune enemy variants and continue playtesting exact original-game timings
