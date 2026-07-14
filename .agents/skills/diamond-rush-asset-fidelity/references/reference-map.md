# Asset And Rendering Reference Map

Use this index to load only the material needed for the current visual discrepancy.

## Primary Evidence

- Java sprite parser: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/f_Sprite.java`
- Java animator: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/b_SpriteAnimator.java`
- Java loader/draw code: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/i.java`
- Original resources: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources`
- Decoded manifest: `decoded/sprites/manifest.json`
- Visual inspector: `decoded/sprites/index.html`

All repository-relative paths below are resolved from the repository root, not this file.

## Read By Discrepancy

Use `rg -n '^## ' docs/diamond-rush-fidelity-reference.md` and read the relevant range.

| Discrepancy | Read first |
| --- | --- |
| Raw ID to sprite, draw pass, alpha, anchors, flips, flags, fire, tutorial art | `Angkor 全部关卡素材映射` |
| Stage-specific animation timing or special object art | The matching stage section in `docs/diamond-rush-fidelity-reference.md` |
| Java resource slots and draw helpers | `docs/diamond-rush-source-mapping.md` sprite/resource sections |
| Product layout, 240x320 viewport, HUD, interpolation | `docs/diamond-rush-remake-spec.md` |
| Individual chunks, frames, modules, animations | `decoded/sprites/manifest.json` and `decoded/sprites/index.html` |
| Splash/title image chunks | `decoded/sprites/splash/` |

## Known Non-Negotiable Rules

- Compose game frames from animation and module metadata; fixed `24x24` slicing is invalid.
- Animation-frame `x/y` is not a universal anchor. Apply it only when the Java caller uses `applyFrameOffset()`.
- Draw scrolling terrain first, then dynamic objects with the original extra scan margins, then post-foreground overlays.
- Raw `255` is empty and must never map to a world frame.
- World-specific foreground animations must switch resources when the world changes.
- Closed foreground containers conceal their player-layer payload.
- Preserve premultiplied source color handling, PNG alpha, and transparent atlas padding.
- Draw must not allocate `ebiten.Image` or mutate gameplay state.

## Extraction Commands

```bash
go run ./tools/drsprite -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources -out decoded/sprites
go run ./tools/drsprite -demo-sprites /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/demoSpr.bin -out decoded/sprites/tutorial
go test ./tools/drsprite
```

Do not regenerate all decoded assets just to inspect them. Use the checked-in manifest and inspector first.
