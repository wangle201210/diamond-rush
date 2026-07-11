# Diamond Rush Source-Fidelity Guide

本项目的目标是依据本机可运行的 Java ME 原版，复刻可观察到的 Diamond Rush 行为。不要依据旧版 Go 原型、记忆、截图或 Boulder Dash 的通用规则补全机制。只要 Java 源码中存在对应逻辑，就以源码和原版运行结果为准。

## 权威资料

- Java 参考项目：`/Users/wanna/mine/github/wangle201210/DiamondRushSource`
- 主游戏逻辑：`DiamondRushSource/src/main/java/i.java`
- 原始可运行字节码：`DiamondRushSource/jars/diamondrush.jar`
- 精灵解析：`DiamondRushSource/src/main/java/f_Sprite.java`
- 动画包装：`DiamondRushSource/src/main/java/b_SpriteAnimator.java`
- 音效 ID：`DiamondRushSource/src/main/java/j_SoundManager.java`
- 世界及经济常量：`DiamondRushSource/src/main/java/a_Config.java`
- 原始资源：`DiamondRushSource/src/main/resources`
- Angkor 关卡包：`DiamondRushSource/src/main/resources/w0.bin`
- Java 原版启动脚本：`DiamondRushSource/run-diamond-rush.sh`
- Java 原版运行参数：`240x320`、缩放 `2`、Nokia 模式、`jars/diamondrush.jar`

运行原版：

```bash
cd /Users/wanna/mine/github/wangle201210/DiamondRushSource
./run-diamond-rush.sh
```

仓库内资料：

- `docs/diamond-rush-source-mapping.md`：Java 类、状态机、资源和玩法源码定位
- `docs/diamond-rush-remake-spec.md`：复刻范围、当前实现和验收条件
- `docs/diamond-rush-original-gameplay.md`：较宽泛的玩法研究，可信度低于源码
- `docs/angkor-world0-inventory.md`：从 `w0.bin` 自动生成的 World 0 ID 清单
- `decoded/world0/stageNN.json`：解出的三层关卡数据
- `decoded/preview/index.html`：原始 ID 检查器；其中彩色 PNG 只是诊断色块，不是游戏贴图
- `decoded/sprites/manifest.json`：`.f` 资源块解析结果
- `decoded/sprites/index.html`：原始精灵、帧和动画检查器

## 证据优先级

1. 原始 `jars/diamondrush.jar` 字节码；当 `i.java` 出现不可能的循环、缺失 break 或其他反编译伪影时，用 `javap -p -c` 核对原 JAR，不能读取 `target/classes` 中重新编译的错误控制流。
2. `i.java` 中的加载、更新、碰撞、渲染和状态切换代码。
3. FreeJ2ME 中运行原版并记录的逐帧状态及操作结果。
4. `w0.bin`、`.f`、语言和音频资源解码结果。
5. 原作视频或截图，仅用于确认画面和用户可见结果。
6. 现有 Go 实现、旧原型和文字攻略只能作为线索，不能证明原作行为。

源码能回答的问题不要通过看图猜。截图无法证明碰撞条件、更新顺序、对象计时器、检查点恢复内容或状态位语义。

## 原始数据约定

- 逻辑画布为 `240x320`，上 HUD `40px`，游戏区 `240px`，下 HUD `40px`。
- 格子为 `24x24px`，主循环为 `20 TPS`。
- Java `byte` 有符号；解码 JSON 中 raw `255` 对应 Java `-1`，通常表示空对象。
- 关卡按 `player -> background -> foreground` 三层存储，每层均为 `width * height` 字节。
- `crtStagePlayerLayer[x][y]` 保存角色、可动物体、敌人和拾取物的 raw ID。
- `crtStageForegrondLayer[x][y]` 在初始化后常把低 8 位作为 foreground ID，高位作为运行状态；不要只保留 raw ID。
- `bIntArrArr[x][y]` 初始来自 background 层，运行时也承载对象状态和位标志。
- `bByteArrArr[x][y]` 是对象动作/插值计时状态，`cByteArrArr[x][y]` 用于唤醒、更新或脏状态。
- `aSInt` 是影响动画、火焰、滚落准备等逻辑的全局源码帧计数。
- 主角格间移动使用 `jInt=18`，通常每源码帧减 `6`，形成 `18 -> 12 -> 6 -> 0`。
- 动态对象活动区的源码扫描是从下到上、再从左到右。更新顺序和一次扫描内是否再次处理移动后的对象都必须按源码验证。

不要把背景、foreground 高位或对象计时器合并成一个简化的“tile type”。后续门、机关、宝箱、敌人和重力物体都依赖这些独立状态。

## 每关复刻流程

1. 先读取对应 `stageNN.json`，列出三层 ID 直方图、所有低频 ID 和关键坐标。
2. 在 `i.java` 中追踪每个实际出现的 raw ID：初始化分支、碰撞分支、对象更新、渲染和拾取/触发结果。
3. 记录完整状态转换，至少包括 ID、状态位、对象 timer、全局帧条件、移动插值、伤害及声音事件。
4. 从资源加载代码反查 `aClassfArr[]`/`textures[]` 的文件和 chunk，不能直接假定 raw ID 等于精灵帧 ID。
5. 使用 `decoded/sprites/*-animations.json` 和 modules 图组合原始帧；有非零锚点或 flip 时不要只裁 `frames.png`。
6. 先实现纯运行时逻辑和确定性测试，再接 Ebitengine 输入、相机和绘制。
7. 在 FreeJ2ME 和 Go 中执行相同输入脚本，逐源码 tick 对比状态，不只对比最终能否通关。
8. 只有该关实际出现的机制全部通过差分验证后，才开始下一关；不要提前实现未出现的对象来掩盖当前关差异。

建议每关维护一份 raw ID 到源码的表：

```text
layer/raw ID | coordinates | init branch | collision | update method | render branch | resource chunk | verified trace
```

## 差分验证字段

每个源码 tick 至少比较：

- 主角格坐标、面向、`jInt`/像素偏移、当前动作和动作帧
- player/foreground/background 的变更单元格
- `bIntArrArr`、`bByteArrArr`、`cByteArrArr` 中相关单元格
- 重力物体和敌人的方向、状态位、计时器及插值偏移
- 生命、HP、宝石、命中次数、重试次数、检查点和关卡模式
- 动画、声音、Loading、结算状态的开始 tick 和结束 tick

Go 单元测试通过只能说明 Go 实现符合测试中的假设，不能说明它符合原作。没有 Java/FreeJ2ME 的同输入逐帧对照，就不能使用“完全一致”或“一比一复原”的结论。

## Angkor Stage 1 基线

`decoded/world0/stage00.json` 的动态内容：

- player raw `0` 石头 13 个
- player raw `1` 紫钻 21 个
- player raw `2` 红钻 1 个，位于 `(19,2)` 的 foreground raw `33` 宝箱下
- player raw `10` 可穿行/消除植被 76 个
- player raw `12` 配额标记 1 个，位于 `(20,9)`，background 值为 `10`
- player raw `19` 绿蛇 4 个
- player raw `23` 向左火焰发射器 1 个
- player raw `79` 入口 1 个
- foreground raw `4` 检查点 3 个
- foreground raw `5` 终点 1 个，位于 `(22,9)`
- foreground raw `33` 宝箱 1 个
- foreground raw `20..23` 各 8 个动画前景格

爬虫、关卡数据内预置的机关门、锤子、钩索和特殊道具未出现在 Stage 1，不应成为 Stage 1 的完成阻塞项；唯一的门是 raw `79` 初始化时动态创建的入口临时门。其余机制在出现的后续关卡中再按同一流程审计。

## Stage 1 已落实的源码规则

截至 2026-07-11，第一关关卡内实际出现的核心规则已按源码/原 JAR 字节码落实，并有 Go tick 级回归测试：

1. raw `12` 是可见且阻挡的配额门，保留在 player layer，显示 `aaInt` 剩余数量；raw `1`/`41` 使数量递减，归零时才清除。可读反编译在 `i.java:3789` 显示一次错误的清除赋值，但原 JAR 运行状态和 `i.java:7271` 的绘制分支均证明该对象实际保留。它使用 `cm.f` chunk `5`。raw `5` 终点不检查配额，因此可以从不经过 `(20,9)` 的路线触发终点。
2. raw `79` 入口会在 `(入口.x-2, 入口.y)` 创建 foreground raw `7` 临时门，合并状态从可通行的 `0x3f` 开始；第四步自动入场前调用 `doorHeadClose` 将其变为阻挡的 `0x0f`，并播放音效 `14`。
3. raw `5` 读取 background 方向 `2`，主角以 `18 -> 12 -> 6 -> 0` 向右自动行走，直到 `x > stageWidth + 5` 才进入 `bByte=35` 等价 Loading。
4. 绿蛇使用低方向位、`0x7000` 待转向位和 `21 -> 18 -> ... -> 0` timer；移动采用源码扫描顺序和像素重叠，`hurtHero(1,48,direction)` 会按方向轮询空相邻格并产生 `jInt=18` 击退。
5. raw `0`/`1` 使用 packed direction/rotation/roll bits；垂直下落、24 到 27 tick 滚落准备、`12 -> 6 -> 0` 斜落、滚动紫钻像素收集、石头压伤及石头先压蛇后落下均有回归覆盖。`OVoid()` 的准备滚动位移是水平 offset、`offset²/24` 下沉和 `aSInt%3` 抖动；横向经过圆形支撑时也使用同一抛物位移。石头在 timer 归零的落地帧清方向、更新旋转并播放音效 `14`，普通支撑上的下一对象 tick 才清除左右标志。
6. raw `10` 在进入帧只标记激活，下一对象扫描才转换成 foreground raw `32`；8 帧消散按奇偶源码 tick 共持续到第 16 tick。
7. Stage 1 宝箱在主角完全停稳后的对象扫描启动动画 `40`；箱盖状态为 `1@start -> 2@start+1 -> 3@start+4`，tick `37` 的动画序列 `12` 播放音效 `4`，tick `39` 的序列 `13` 发奖、播放 `cm.f` chunk `7` animation `0` 并显示头顶红钻，tick `67` 结束。关闭的 foreground raw `14/33` 必须遮蔽 player raw `2/4/5/6/7/24/26/27/41`，内容只能在奖励动画中出现。动画 `48` 只适用于 5 秒内连续的其他拾取，Stage 1 没有可先触发该条件的道具。
8. 离开检查点后按 `*` 会在动画启动帧播放音效 `2`，再播放 42 tick 的主角动画 `19`，结束后才扣命、增加 retry 并恢复；站在检查点按 `5` 或 `*` 立即重置、播放音效 `9` 且不扣命。
9. 紫钻拾取播放 `cm.f` chunk `7` animation `3`；紫钻、检查点和蛇均使用源码的直接动画序号节奏，而不是通用 frame-time 播放。

Stage 1 结算已按原 JAR 字节码复现 `bByte=35` 的 12 步 Loading、音效 `15`、`bByte=17` 分阶段动画、精确行坐标和三种特效位移取帧。文字使用 FreeJ2ME 导出的 `SansSerif Bold` 10/12px 字形与度量，通关、四类勋章位、生命、血量和工具等级持久化到独立 JSON。结算 Continue 会进入按 `map_angkor.out` 解出的世界地图，并按原节点关系解锁前五关。

## Stage 1 必须覆盖的回放

- raw `12` 在剩余 `1` 时仍绘制并阻挡，第 10 点配额清除后可通过；raw `5` 终点保持独立可触发
- 石头/紫钻垂直下落、滚落准备、落地、叠放和压到主角/蛇
- 蛇直行、遇阻转向、接触、像素重叠、受击方向和击退
- raw `23` 火焰完整动画周期及每个 reach 的伤害范围
- 宝箱关闭、开盖、主角动作、tick `39` 发奖、特效和动作结束
- 三个检查点的激活、Action/5 房间重置、`*` 主动召回、死亡复活
- raw `5` 触发、自动走出边界、祝贺文字、Loading 和完整结算阶段

## Angkor Stage 2-5 基线

前五关必须使用以下真实关卡数据和 20 TPS 正式路线回归，不能用传送、停掉对象更新或直接改层数据证明“可通关”：

| 显示关卡 | 数据 | raw `12` 配额门 | raw `5` 终点 | 该关新增的核心机制 | 正式回归 |
| --- | --- | --- | --- | --- | --- |
| Stage 1 | `stage00.json`，26x21 | `(20,9)`，10 | `(22,9)` | 重力、滚落、蛇、火焰、宝箱、检查点、召回 | `TestRuntimeStage00CanBeCompletedAtSourceCadence` |
| Stage 2 | `stage01.json`，27x24 | `(22,3)`，15 | `(23,2)` | 更长的重力追逐与四个有序检查点 | `TestRuntimeStage01CanBeCompletedAtSourceCadence` |
| Stage 3 | `stage02.json`，27x26 | `(22,21)`，20 | `(23,21)` | 金银钥匙/锁、raw `30`、敌人分组门 | `TestRuntimeStage02CanBeCompletedAtSourceCadence` |
| Stage 4 | `stage03.json`，40x23 | `(36,5)`，25 | `(37,5)` | 压力机关、raw `24` 神秘锤、破墙和蛇眩晕 | `TestRuntimeStage03CanBeCompletedAtSourceCadence` |
| Stage 5 | `stage04.json`，51x24 | `(3,10)`，30 | `(2,10)` | 钩索、锤/钩组合谜题、压力门、分组蛇和爬虫 | `TestRuntimeStage04CanBeCompletedAtSourceCadence` |

Stage 5 在原流程中是取得 Bavaria 神秘钩索后回访的关卡。当前五关独立切片没有 Bavaria 节点，因此加载 `stage04` 时至少注入工具等级 `2`；这不是关卡内赠送钩索，而是还原原作进入该关完整收集路线时的跨世界前置状态。

## 锤子和钩索源码规则

- Action/电脑键 `5` 先处理当前检查点，其次处理工具；`*` 是主动召回，不是钩索键。
- 锤子要求工具等级 `>=1`，方向动画为上/右/下/左 `13/14/15/16`，第 3 tick 命中；向上动作共 11 tick，其余方向 12 tick。它使 raw `30` 进入破坏状态，使普通蛇以低 `0xf8` 位计时 `0x78` 眩晕，并直接移除符合源码条件的未标记 raw `43`。
- 钩索要求工具等级 `>=2`，只沿水平方向搜索 2 到 3 格。JAR 候选集合为 raw `0/1/8/9/11/14/19/43/47/48`；raw `48` 状态位 `0x8` 时排除，关闭的 foreground raw `7` 以及中途非空对象会阻断。
- 钩索绳段是 player-layer raw `32`，相邻段从 timer `18` 开始，以 `18 -> 12 -> 6 -> 0` 伸展。源码按从左到右扫描，所以向右新建的第二段在同一扫描中立即减到 `12`，向左则保留 `18`。
- 命中后目标不是只移动一格：普通物理目标会逐格拉到主角相邻格；raw `1` 紫钻还会拉进主角所在格，整段动作结束后才收集。每次目标移动使用 18px 插值，残余绳段会在同一扫描中重新捕获目标并把该步 motion 归零。
- 主角伸出动画为右 `20`、左 `22`，最终拉回动画为右 `21`、左 `23`。钩索期间禁止输入和受伤；命中/重新捕获发出音效 `12`，恢复检查点时必须清理 raw `32`、目标 motion 和钩索状态。

## 前五关素材映射

- raw `12` 配额门：`decoded/sprites/cm/chunk05-*`，两模块组合并叠加 HUD 数字。
- raw `11` 爬虫：`decoded/sprites/gen1/chunk04-*`，6 个模块；正常帧为 `(aSInt >> 1) % 3`，状态 phase 使用后续模块。
- foreground raw `6` 压力机关：`decoded/sprites/gen2/chunk09-*`，单个 `24x13` 模块，底部对齐并随压入量下移。
- raw `32` 钩索：颜色 `#d3d7e7` 的水平线加主角 sprite module `0/1`，不是独立位图。
- player raw `24/31/33/41/79` 在前五关由前景容器、锁、入口流程持有可见图形或本身不可见，不得回退绘制诊断蓝块。

## 工具命令

```bash
go run ./tools/drdecode -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/w0.bin -out decoded/world0
go run ./tools/drinspect -in decoded -out decoded/preview
go run ./tools/drsprite -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources -out decoded/sprites
go run ./tools/drsound -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/snd.f -out decoded/audio
go run ./tools/drworldaudit -in decoded/world0 -out docs/angkor-world0-inventory.md
javac tools/drfont/ExportFont.java && java -cp tools/drfont ExportFont decoded/fonts
go run ./cmd/originalrush
```

改动原作数据运行链路后执行：

```bash
go test ./internal/original ./internal/originalgame ./tools/drdecode ./tools/drinspect ./tools/drsprite ./tools/drsound ./tools/drworldaudit
go build -o /tmp/originalrush-smoke ./cmd/originalrush
go test ./...
rm -f /tmp/originalrush-smoke
```

macOS 无主显示器环境下，Ebitengine/GLFW 测试可能在项目测试执行前失败；不要把这类初始化失败误判为玩法回归。文档或清单单独修改时无需重启正在运行的游戏。

## 单关完成定义

一个关卡只有同时满足以下条件才算完成：

- 该关三层中所有实际出现的 raw ID 均有源码定位和明确语义
- 碰撞、更新顺序、计时器、状态位、伤害、重置和退出路径均有测试
- 所有关键回放通过 Java 与 Go 的逐 tick 差分
- 正常路线、死亡/召回路线、遗漏收集路线和边界操作均可完成
- 原作资源的帧、锚点、flip、层级、相机和 HUD 时序已核对
- 当前五关切片之外的教程 `stage13`、商店、Bavaria/Tibet 和 Angkor Stage 6-13 被明确列为范围外项目，而不是默认为已完成
- 文档中的“当前实现”与代码一致，不保留已被源码否定的规则
