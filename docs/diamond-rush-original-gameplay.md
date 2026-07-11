# Diamond Rush Original Gameplay Design Notes

Last researched: 2026-07-11

This document is the gameplay baseline for the Go + Ebitengine remake. It captures the original Gameloft Java ME / BlackBerry Diamond Rush design as research notes, not as a final legal/content asset spec. Use it to decide what to implement, what to cut for the five-level pack, and what still needs frame-by-frame verification from captured gameplay.

## Source Confidence

- High confidence: facts repeated by multiple reference pages or visible in walkthrough material.
- Medium confidence: facts from one wiki or guide page, especially when consistent with known gameplay.
- Needs verification: timing, exact tile physics, exact map counts, enemy AI details, and secret-stage unlock edge cases.

Primary references:

- Wikipedia overview: https://en.wikipedia.org/wiki/Diamond_Rush
- Diamond Rush Wiki gameplay page: https://diamond-rush.fandom.com/wiki/Diamond_Rush
- Mobile Games Wiki page: https://mobilegames.fandom.com/wiki/Diamond_Rush
- Diamond Rush Guide index/review: https://dia-rush.blogspot.com/
- Angkor Wat Stage 5 guide: https://dia-rush.blogspot.com/2011/09/angkor-wat-stage-5.html
- Angkor Wat archive with Stage 1/3/4/5/7/8 notes: https://dia-rush.blogspot.com/2011_09_09_archive.html
- Angkor Wat Stage 8 guide: https://dia-rush.blogspot.com/2011/09/angkor-wat-stage-8_09.html
- Bavaria Stage 3 guide/comments: https://dia-rush.blogspot.com/2011/09/bavaria-stage-3.html
- Perfect walkthrough playlist: https://www.youtube.com/playlist?list=PL439E54997BE62088
- Full gameplay walkthrough: https://www.youtube.com/watch?v=MJ51oK75iWs

## Identity And Structure

- Genre: side-view action puzzle adventure, inspired by Boulder Dash but expanded with adventure progression, tools, health/lives, shops, secret exits, bosses, and world-specific hazards.
- Original platform target: Java ME feature phones, later BlackBerry.
- World structure: three themed regions:
  - Angkor Wat / Ancient Khmer temple jungle.
  - Bavaria dungeon / castle.
  - Tibet / Siberia icy caves, naming varies by version/source.
- Total scope: 40 stages across the worlds, including visible stages and secret stages.
- Each world culminates in a boss stage.
- Meta objective: recover major world diamonds/crystals and unlock the ancient seal.

## Core Loop

1. Enter a stage from the world map.
2. Explore a side-view tile maze with small-screen scrolling camera.
3. Collect required purple diamonds and/or required keys.
4. Avoid traps, falling objects, and enemies while solving tile puzzles.
5. Open or reach the stage exit.
6. Keep collected progress, red diamonds, tools, and unlocks.
7. Revisit earlier stages after acquiring tools to collect previously inaccessible rewards or secret exits.
8. Spend purple diamonds in shops for health/armor upgrades.

Important distinction from a pure Boulder Dash clone:

- The stage is not only "dig dirt and collect diamonds." Original Diamond Rush includes adventure gates, keys, chests, checkpoints, health, lives, tools, secret paths, shops, world map progression, and boss encounters.

## Controls

Original keypad/control behavior to preserve:

- Move: directional pad or numeric keypad `2/4/6/8`.
- Action: `5`.
  - Uses the currently available action tool in the facing direction.
  - Magic/Mystic Hammer attacks or breaks eligible blocks.
  - Mystic Hook / grappling hook acts in facing direction.
- Checkpoint recall: `*` costs one life and returns the player to the last checkpoint.
- Standing on a checkpoint and pressing `5` resets the room/checkpoint state, according to Diamond Rush Wiki.

Remake implication:

- Desktop keys can be modernized, but the primary touch/mobile layout should still read as a phone keypad with a meaningful center `5`.
- The action button must be directional/facing-based, not a mouse-like free target.

## HUD And Presentation

Original HUD elements reported by sources:

- Bottom status bar.
- Compass indicator.
- Lives remaining.
- Energy/health bar.
- Red and purple diamond counters.
- World/stage progress implied by map screen.

Title/menu options in source:

- New Game.
- Continue.
- Options.
- Help.
- About.
- Exit.

Visual rules for remake:

- Gameplay should be a small-screen scrolling viewport, not a full-map board.
- Menu/title should be a separate illustrated screen, not the gameplay map.
- The tone is temple/dungeon adventure with readable icons, not abstract cave-only boulders.

## Diamond Rush Identity Anchors

Use this as the first acceptance checklist before adding more levels or art. If these are weak, the remake will drift toward Boulder Dash even if the tile physics work.

| Anchor | Original design role | Five-level remake requirement |
| --- | --- | --- |
| Phone adventure framing | The game is built around Java ME phone screens, keypad actions, HUD bars, and stage/world menus. | Keep a constrained scrolling viewport, separate title/menu/map screens, and a visible phone-style control language. Do not show the whole map as the default play view. |
| Temple/dungeon expedition | The player is an explorer moving through authored ruins, not an abstract miner in random caves. | Angkor-style foreground art, carved walls, chests, doors, shrines, icons, and authored paths must dominate the first pack. Rocks support puzzles; they are not the visual identity. |
| Purple and red diamond split | Purple diamonds are common stage/currency resources; red diamonds are progression/completion rewards. | Track both separately, place red diamonds mostly in chests/secret branches, and use at least one red-diamond gate or seal. |
| Tool-gated backtracking | Compass, Mystic Hammer, Mystic Hook, and later freeze/ice tools create revisit value. | Show an unreachable reward before its tool, then let hammer/hook open optional collection or secret routing later. |
| Center `5` action | Tools use a facing-direction action, matching keypad play. | Keep Action/5 as a primary command. Hammer and hook must use facing direction rather than mouse targeting. |
| Checkpoint magic | Checkpoints are active objects with recall/reset behavior, not just respawn coordinates. | Support checkpoint activation, life-cost recall, and room reset from the checkpoint action. |
| Health/lives/shop | Original hazards are not just one-hit deaths; upgrades and lives are part of progression. | Use health, armor, lives, potions, and a small between-stage shop. Avoid making every trap instant death. |
| Chests and secrets | Chests contain real rewards and secrets alter routing/completion. | Decode container payloads from the original player/foreground layers and preserve the original reward timing. |
| Boss/world finale | Each world culminates in a boss stage. | The fifth level needs at least one guardian/boss puzzle gating the exit. |
| Authored stages | Levels are handcrafted with puzzle beats, tool reveals, and revisits. | Prefer five dense, tuned stages over many generic cave boards. |

## Mechanic Inventory Matrix

| System | Source confidence | Original behavior to preserve | Current five-level priority |
| --- | --- | --- | --- |
| Small scrolling viewport | High | Camera follows the player through a phone-sized playfield. | Required for all gameplay. |
| World/stage map | High | Stages are selected from themed map progression. | Required, even with only five nodes. |
| Purple diamonds | High | Common collectible and spendable/store resource. | Required. |
| Red diamonds | High | Scarce progression/completion collectible, often chest/secret related. | Required. |
| Chests | High | Containers for diamonds, keys, health, lives, and progression rewards. | Required. |
| Keys and locked doors | High | Gate local routes and exits. | Required by level 3. |
| Checkpoints | High | Activate checkpoint, recall with life cost, and reset room/checkpoint state. | Required by level 1 or 2. |
| Compass | Medium-high | Points toward magical checkpoint/progression marker. | Required as a HUD aid, exact art can be revised later. |
| Mystic Hammer | High | Breaks eligible blocks and stuns/attacks enemies via action button. | Required by level 3 or 4. |
| Mystic Hook / grapple | Medium-high | Pulls or grabs objects/items in facing direction; exact target list and range need footage verification. | Required by level 4. |
| Freeze / ice tool | Medium | Later-game tool for ice/secret routing; exact mechanics still uncertain. | Optional for five-level Angkor slice. |
| Health and armor | High | Damage drains energy/health; upgrades improve survivability. | Required. |
| Lives | High | Death/recall consumes lives; game-over behavior varies by version. | Required. |
| Shop/upgrades | Medium-high | Purple diamonds buy armor/health or similar survivability upgrades. | Required in compact form. |
| Boulder/rock physics | High | Falling/pushing/crushing objects create puzzle pressure. | Required, but visually secondary to adventure structure. |
| Fire/spikes/spears/traps | High | Hazards vary by world and often damage rather than only kill. | Required in at least static and timed forms. |
| Enemies | High | Snakes/spiders and world-specific enemies patrol/block routes and can be stunned/crushed. | Required, with at least patrol and chase variants. |
| Secret exits/stages | High | Hidden exits and paths create alternate progression and completion. | Required as a simplified secret route. |
| Bosses | High | Each world has a final boss encounter. | Required as fifth-level guardian. |

## Progression Currencies

### Purple Diamonds

- Common stage collectible.
- Used to satisfy per-stage exit requirements.
- Also used as store currency for armor/health upgrades.
- Some stages require revisiting with later tools to collect all purple diamonds.

### Red Diamonds

- Progression collectible, typically found in chests.
- Required to unlock later worlds/major doors.
- Sources mention a red-diamond requirement to unlock the next world; exact numbers should be verified per version.
- Red diamonds are central to 100% completion and secret-stage routing.

### Chests

- Chests may contain purple diamonds, red diamonds, keys, health potions, extra lives, or other useful items.
- Chests are not merely score objects; they are progression/reward containers.
- Secret chests can hold extra lives.

### Lives

- The game has a lives system.
- Losing all lives triggers game over.
- Some sources indicate game over removes purple diamonds while preserving other progress. Verify exact version behavior before implementation.
- Extra lives may be found in chests or awarded for strong stage clears.

### Health / Energy / Armor

- Player has an energy/health bar.
- Armor or health upgrades can be bought with purple diamonds.
- Damage can be partial, not always instant death.
- Some hazards/enemies kill quickly or instantly; exact rules need verification.

## Stage Completion And Rewards

Reported "perfect" style conditions:

- Collect all purple diamonds.
- Collect all red diamonds.
- Do not restart the stage.
- Avoid taking hits.

Expected remake requirements:

- Track per-stage completion flags separately:
  - clear stage.
  - all purple diamonds.
  - all red diamonds.
  - no damage.
  - no restart / no checkpoint recall.
  - secret exit found.
- Show completion marks on the stage/world map.

## Tools And Abilities

### Compass

- Obtained in Angkor Wat Stage 1 according to guide notes.
- Helps locate the next magical checkpoint.
- In the remake, this should be a HUD/navigation feature rather than a generic collectible.

### Mystic / Magic Hammer

- Obtained in Angkor Wat Stage 4 according to guide notes.
- Used to stun snakes/enemies and break eligible blocks.
- Earlier stages contain blocked rewards that cannot be fully collected until this tool is acquired.
- Action input is `5` while facing the target.

### Mystic Hook / Grappling Hook

- Available in Bavaria Stage 3 according to guide comments and Angkor Stage 5 notes.
- Used to grab/pull otherwise unreachable items or objects.
- Angkor Wat Stage 5 and secret stages require or recommend it for full collection.
- The hook is part of the reason the original has backtracking and is not just linear level solving.

### Freeze / Ice Hammer

- Required for Angkor Wat Stage 8's secret path according to guide notes.
- Guide comments imply it is obtained later than early Angkor; exact acquisition stage needs video verification.
- Likely interacts with water/ice/balls or freezing hazards, but exact rules need verification.

### Other Mentioned Utility

- A "grapple" is mentioned by the guide, likely the same as Mystic Hook.
- Shop armor/chain vest is mentioned in comments and wiki-style sources; exact item list needs capture verification.
- Magic potion / secret potion is mentioned in guide comments; exact function needs verification.

## Tile And Object Catalog

### Terrain

- Solid walls and temple/dungeon/ice-specific solid blocks.
- Diggable/clearable earth or soft blocks.
- Breakable blocks that require Mystic Hammer.
- Hidden/false walls and secret paths.
- Water appears in Bavaria-related guide comments and needs exact mechanics.
- Ice and falling ice appear in Tibet/Siberia material.

### Collectibles

- Purple diamonds.
- Red diamonds.
- Chest rewards.
- Silver keys.
- Gold keys.
- Health potions.
- Extra lives.
- Tool pickups.

### Gates

- Stage exits.
- World-lock doors that require red diamonds.
- Padlocks/exits opened by collecting required purple diamonds.
- Doors opened by silver/gold keys.
- Secret exits leading to hidden stages.

### Physics Objects

- Boulders/rocks.
- Round stone balls, including water-related ball puzzles mentioned in Bavaria comments.
- Falling ice in Tibet/Siberia.
- Objects can block, crush, trigger puzzles, or need tool interaction.

### Checkpoints

- Checkpoint circles/markers.
- Stepping on a checkpoint activates it.
- `*` returns to checkpoint at life cost.
- `5` on checkpoint resets the room, according to Diamond Rush Wiki.
- Compass points to the next ordered magical checkpoint, then the stage goal. The source's Angkor demo stage grants it immediately before Stage 1.

## Hazards And Traps

- Falling boulders/rocks crush the player.
- Fire traps.
- Giant spears.
- Spike-like traps.
- Falling ice in Tibet/Siberia.
- Water-related puzzle hazards or traversal restrictions in Bavaria.
- Standing under or near boulders may drain energy before death according to Mobile Games Wiki; exact rule needs video verification.

Implementation note:

- Do not model all hazards as instant death. Original uses health/energy, lives, and some partial-damage states.

## Enemy Catalog

Confirmed or repeatedly mentioned:

- Snakes.
- Spiders.

Also reported in Russian Wikipedia:

- Wolverine-like enemies.
- Turtles.
- Aggressive natives/tribal enemies.

Source confidence for enemy types varies; prioritize snakes/spiders first, then verify world-specific enemies from footage.

Expected enemy behavior categories:

- Patrol enemies.
- Area blockers.
- Enemies stunned or defeated by hammer.
- Enemies killed by falling rocks/objects.
- Bosses at end of each world.

## World And Level Design Patterns

### Angkor Wat

- Intro/tutorial world.
- Teaches compass/checkpoint, diamonds, keys, rocks, snakes, and hammer.
- The separate Angkor demo/tutorial stage (`stage13` in the packed source order) gives the Compass and presents tutorial dialog before Stage 1 (`stage00`).
- Stage 3 contains hammer-gated blocked rewards.
- Stage 4 gives Mystic Hammer.
- Stage 5 requires Mystic Hammer and Mystic Hook for full collection.
- Stage 7 introduces or exposes a secret path.
- Stage 8 has a Freeze Hammer requirement for a secret path.
- Secret stages exist and often expect later tools.

### Bavaria

- Dungeon/castle world.
- Contains Mystic Hook in Stage 3 according to guide comments.
- Includes secret stages.
- Guide comments mention water/ball/door puzzles; exact mechanics should be captured.

### Tibet / Siberia

- Icy cave world.
- Uses falling ice and snow/ice visual language.
- Likely introduces freeze/ice-specific mechanics and later-game hazards.
- Boss stage completes the final world and ancient seal objective.

## Secret Stages And Backtracking

- Secret exits and hidden paths are a major part of the original identity.
- Some stages should be cleared normally before taking secret paths, because secret exits can redirect map progression.
- Tools acquired later unlock full completion in earlier stages.
- A five-level remake should still include this pattern in miniature:
  - Level 1/2: visible unreachable reward.
  - Level 3: tool acquisition.
  - Level 4: return-style optional route or secret branch.
  - Level 5: combined test with secret reward.

## Bosses

- Each world ends with a boss stage.
- Current remake has a first boss-like guardian system in the five-level pack.
- For a five-level vertical slice, at least one final-stage boss or boss-like encounter is required to feel closer to Diamond Rush.
- Boss behavior should be puzzle/action hybrid, not pure HP combat.

## Store And Upgrade System

- Purple diamonds can buy armor/health upgrades.
- Store/upgrade UI is part of the original progression identity.
- Current remake only has score and best score, which is not equivalent.

Minimum remake implementation:

- Track purple diamond bank across stages.
- Add simple shop between levels.
- Offer at least one armor/health upgrade.
- Keep red diamonds separate from spendable purple diamonds.

## Five-Level Remake Target

The user's target is five polished levels, not all 40 original stages. The five-level pack should compress original progression rather than invent a Boulder Dash-only cave game.

Recommended five-level structure:

1. Angkor tutorial gate:
   - Compass/checkpoint, purple diamonds, exit requirement, one key door.
   - No complex enemies.
2. Rock and snake puzzle:
   - Boulder crush, hammer-gated optional reward visible but inaccessible.
   - First enemy.
3. Mystic Hammer acquisition:
   - Breakable blocks, stun enemy, chest with red diamond.
   - Shop becomes relevant afterward.
4. Mystic Hook / secret route:
   - Hook pickup or newly active hook.
   - Pull reward/object, hidden exit or secret chamber.
   - Backtracking-shaped optional completion.
5. Combined trial / boss-like finale:
   - Keys, red diamond chest, rocks, hammer, hook, checkpoint, traps, and a boss/puzzle guardian.
   - One secret reward requiring the earlier tool chain.

## Current Runtime Boundary

The maintained implementation is `cmd/originalrush` plus `internal/original` and `internal/originalgame`. It reads decoded original data and extracted source resources directly.

The current playable slice is Angkor `stage00` through `stage04`. `AGENTS.md` records the source-derived rules and named route regressions that are already implemented. This broader research document remains useful for feature discovery, but it is not evidence that an untested behavior matches the original.

## Verification Still Needed

Before claiming high fidelity:

- Capture or locate footage for each original world at native aspect ratios.
- Record exact HUD layout and counter names.
- Verify tool acquisition stages and exact item names across versions.
- Verify hook candidates and consequences again when later-world stages introduce IDs not exercised by the first five Angkor stages.
- Verify hammer targets again when later-world stages introduce new breakable classes or enemies.
- Verify health damage numbers by hazard/enemy.
- Verify checkpoint reset and life-cost semantics.
- Verify shop item names, costs, and upgrade effects.
- Verify boss mechanics for Angkor, Bavaria, and Tibet/Siberia.
