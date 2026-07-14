# Gameplay Reference Map

Use this index to avoid loading every historical fidelity note.

## Primary Evidence

- Java gameplay: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/i.java`
- Original bytecode: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/jars/diamondrush.jar`
- Config/economy: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/a_Config.java`
- Sound IDs: `/Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/j_SoundManager.java`
- Decoded stages: `decoded/world0/` and `decoded/world1/`

All repository-relative paths below are resolved from the repository root, not this file.

## Route By Task

Use `rg -n '^## ' docs/diamond-rush-fidelity-reference.md` to find section boundaries, then read only the needed range.

| Task | Read first |
| --- | --- |
| Data model, cadence, scan order, differential fields | `原始数据约定`, `每关复刻流程`, `差分验证字段` |
| Angkor Stage 1 | `Angkor Stage 1 基线`, `Stage 1 已落实的源码规则`, `Stage 1 必须覆盖的回放` |
| Angkor stage inventory and route tests | `Angkor Stage 2-13 基线`, `docs/angkor-world0-logic-audit.md` |
| Falling torches | `Stage 6 坍塌火炬源码规则` |
| Normal/secret exits and map nodes | `Stage 7 普通/秘密出口源码规则` |
| Freeze Hammer | `Stage 8 Freeze Hammer 与双路线规则` |
| Great Anaconda | `Stage 9 Great Anaconda 源码规则` |
| Angkor secret stages | `四个秘密关源码路线` |
| Tutorial/demo scripts | `Tutorial Stage 13 源码规则` |
| Hammer and hook | `锤子和钩索源码规则` |
| Chests, rewards, permanent pickups | `宝箱、紫钻与永久坐标` |
| Doors, pressure switches, enemy groups | `门、爬虫与敌人组` |
| Bavaria shared rules, water, potion, boss boundary | `Bavaria 完整性审计边界` |
| World switching and unlock progression | `全局封印世界选择` |
| Completion criteria | `单关完成定义` |

## Supporting Specifications

- `docs/diamond-rush-source-mapping.md`: Java methods, state fields, resource slots, and currently mapped behavior.
- `docs/diamond-rush-remake-spec.md`: product scope and acceptance contract.
- `docs/diamond-rush-original-gameplay.md`: broad research only; lower authority than source/JAR.
- `docs/angkor-world0-inventory.md`: generated World 0 ID inventory.
- `decoded/demo-scripts.json`: decoded demo command streams.

## Useful Inspection Commands

```bash
rg -n "case |raw|STAGE_|demo|checkpoint|save" /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/java/i.java
javap -classpath /Users/wanna/mine/github/wangle201210/DiamondRushSource/jars/diamondrush.jar -p -c i
go run ./tools/drworldaudit -in decoded/world0 -out docs/angkor-world0-inventory.md
```

Do not regenerate Bavaria `stage07` from `src/main/resources/w1.bin`; the checked-in decoded data intentionally follows the original JAR bytes.
