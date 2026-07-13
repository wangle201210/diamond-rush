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
- `docs/angkor-world0-logic-audit.md`：World 0 共用机制、14 关路线、存档经济与跨世界边界审计
- `decoded/world0/stageNN.json`：解出的三层关卡数据
- `decoded/preview/index.html`：原始 ID 检查器；其中彩色 PNG 只是诊断色块，不是游戏贴图
- `decoded/sprites/manifest.json`：`.f` 资源块解析结果
- `decoded/sprites/index.html`：原始精灵、帧和动画检查器
- `decoded/sprites/splash/`：从原版 `spl.f` 三个 PNG chunk 无损提取的标题背景、Logo 和版权图

## 证据优先级

1. 原始 `jars/diamondrush.jar` 字节码；当 `i.java` 出现不可能的循环、缺失 break 或其他反编译伪影时，用 `javap -p -c` 核对原 JAR，不能读取 `target/classes` 中重新编译的错误控制流。
2. `i.java` 中的加载、更新、碰撞、渲染和状态切换代码。
3. FreeJ2ME 中运行原版并记录的逐帧状态及操作结果。
4. `w0.bin`、`.f`、语言和音频资源解码结果。

`src/main/resources/w1.bin` 与实际运行的原 JAR 内 `w1.bin` 有 4 字节差异，全部位于 Bavaria `stage07` player layer。原 JAR 的 `(20,10)/(29,19)/(25,20)` 为空，`(7,17)` 为 raw `10`；`decoded/world1/stage07.json` 必须保留 JAR 版本，不能用源码目录资源重新生成后覆盖。

原 JAR 条目的 SHA-256：`w1.bin` 为 `951b998c82383c55144ed82c5c54a7dc70f638017929d46aa155e40b0a77674e`，`map_scotland.out` 为 `5c21ffc3ac32e6f571cba097eaf81f3a7044804d1e4ae19f2c381586eba543c0`。截至 2026-07-13，从这两个 JAR 条目重新解码所得的 13 个 stage JSON 和地图 JSON 与 `decoded/world1` 逐文件一致；manifest 只有记录的 source 路径不同。
5. 原作视频或截图，仅用于确认画面和用户可见结果。
6. 现有 Go 实现、旧原型和文字攻略只能作为线索，不能证明原作行为。

源码能回答的问题不要通过看图猜。截图无法证明碰撞条件、更新顺序、对象计时器、检查点恢复内容或状态位语义。

## 原始数据约定

- 逻辑画布为 `240x320`，上 HUD `40px`，游戏区 `240px`，下 HUD `40px`。
- 格子为 `24x24px`，主循环为 `20 TPS`。
- Ebitengine 外层以 `60 TPS` 采集输入和刷新显示，但完整源码步骤只能每 3 次 Update 执行一次；`g.tick`、`TickSourceFrame()`、菜单、地图、教学、结算、碰撞、伤害和存档均保持 `20 TPS`。`renderTick/renderPhase` 只允许用于显示。
- 60Hz 中间坐标必须由 `ObjectMotion`/相机源码状态临时计算，不能写回 Runtime。逻辑相机 `cameraX/cameraY` 仍传给 `SetViewport()`，平滑相机只由 Draw 使用；切关、复活、召回和脚本切镜不能继承旧插值状态。
- Draw 路径禁止创建 `ebiten.Image` 或重复构造 sprite `SubImage`；纯色矩形和精灵子图必须在加载期缓存。否则即使逻辑只有 20 TPS，也会因 GPU atlas 分配出现真实掉帧。
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

启动流程同样以源码为准：已有 JSON 进度时先显示原版标题菜单的 `Continue`/`New game`；新游戏必须经过 `No`/`Yes` 二次确认，只有 `Yes` 才能覆盖存档。当前存档不含关卡内即时快照，因此 Continue 回到最高已解锁地图节点；教程未完成则继续教程。

macOS 输入适配保持手机动作语义互斥：`Space`/`5`/数字小键盘 `5` 对应手机 `5`（交互），`Enter`/`R`/`Backspace`/数字小键盘 `*`/`Shift+8` 对应手机 `*`（回到复活点并失去一条命），教程跳过使用 `S`。不要让 Space 同时承担教程跳过，也不要让 Enter 同时承担交互和召回。桌面版所有面向玩家的按键提示统一显示 `SPACE`、`ENTER`、`S`，不要直接显示手机键位 `5` 或 `*`。

人物静止时按下不同于当前朝向的方向键，必须先原地转身：源码在 `i.java:1531-1532` 设置 `kInt|=0x1000`，再于 `10692-10699` 建立 `jInt=18`，后续按 `18→12→6→0` 消耗三个 20Hz 更新帧。点按只改变朝向；松开后再次按同方向才移动；持续按住则在转身完成后的下一帧移动。该门控只属于玩家输入层，raw `79` 自动入场、出口自动行走和 tutorial demo 脚本不能经过转身门控。复活点恢复后源码把朝向重置为右侧（`kInt=2`）。

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
6. raw `10` 在进入帧只标记激活，下一对象扫描才转换成 foreground raw `32`；8 帧消散按奇偶源码 tick 共持续到第 16 tick。Java 把该帧号存于 foreground 高位，与 `bIntArrArr` 对象状态分离；Go 必须使用独立 `ForegroundState`，否则落石进入 raw `32` 时会被画成草，并由草动画改写石头方向/旋转位。
7. Stage 1 宝箱在主角完全停稳后的对象扫描启动动画 `40`；箱盖状态为 `1@start -> 2@start+1 -> 3@start+4`，tick `37` 的动画序列 `12` 播放音效 `4`，tick `39` 的序列 `13` 发奖、播放 `cm.f` chunk `7` animation `0` 并显示头顶红钻，tick `67` 结束。关闭的 foreground raw `14/33` 必须遮蔽 player raw `2/4/5/6/7/24/26/27/41`，内容只能在奖励动画中出现。动画 `48` 只适用于 5 秒内连续的其他拾取，Stage 1 没有可先触发该条件的道具。
8. 离开检查点后按 `*` 会在动画启动帧播放音效 `2`，再播放 42 tick 的主角动画 `19`，结束后才扣命、增加 retry 并恢复；站在检查点按 `5` 或 `*` 立即重置、播放音效 `9` 且不扣命。
9. 紫钻拾取播放 `cm.f` chunk `7` animation `3`；紫钻、检查点和蛇均使用源码的直接动画序号节奏，而不是通用 frame-time 播放。

Stage 1 结算已按 Java `case 17` 复现 `bByte=35` 的 12 步 Loading、音效 `15`、`bByte=17` 分阶段动画和三种特效位移取帧。标题 y 为 `15/32`，四行 label/count 为 `75/91`、`131/147`、`187/203`、`243/259`，不得再使用会裁掉标题的旧坐标。文字使用 FreeJ2ME 导出的 `SansSerif Bold` 10/12px 字形与度量，通关、四类勋章位、生命、血量和工具等级持久化到独立 JSON。结算 Continue 会进入按 `map_angkor.out` 解出的世界地图，并按原节点关系解锁当前已落实的关卡。

## Stage 1 必须覆盖的回放

- raw `12` 在剩余 `1` 时仍绘制并阻挡，第 10 点配额清除后可通过；raw `5` 终点保持独立可触发
- 石头/紫钻垂直下落、滚落准备、落地、叠放和压到主角/蛇
- 蛇直行、遇阻转向、接触、像素重叠、受击方向和击退
- raw `23` 火焰完整动画周期及每个 reach 的伤害范围
- 宝箱关闭、开盖、主角动作、tick `39` 发奖、特效和动作结束
- 三个检查点的激活、Action/5 房间重置、`*` 主动召回、死亡复活
- raw `5` 触发、自动走出边界、祝贺文字、Loading 和完整结算阶段

## Angkor Stage 2-13 基线

Angkor 打包的全部 14 个 stage 必须使用真实关卡数据和 20 TPS 正式路线回归，不能用传送、停掉对象更新或直接改层数据证明“可通关”：

| 显示关卡 | 数据 | raw `12` 配额门 | raw `5` 终点 | 该关新增的核心机制 | 正式回归 |
| --- | --- | --- | --- | --- | --- |
| Stage 1 | `stage00.json`，26x21 | `(20,9)`，10 | `(22,9)` | 重力、滚落、蛇、火焰、宝箱、检查点、召回 | `TestRuntimeStage00CanBeCompletedAtSourceCadence` |
| Stage 2 | `stage01.json`，27x24 | `(22,3)`，15 | `(23,2)` | 更长的重力追逐与四个有序检查点 | `TestRuntimeStage01CanBeCompletedAtSourceCadence` |
| Stage 3 | `stage02.json`，27x26 | `(22,21)`，20 | `(23,21)` | 金银钥匙/锁、raw `30`、敌人分组门 | `TestRuntimeStage02CanBeCompletedAtSourceCadence` |
| Stage 4 | `stage03.json`，40x23 | `(36,5)`，25 | `(37,5)` | 压力机关、raw `24` 神秘锤、破墙和蛇眩晕 | `TestRuntimeStage03CanBeCompletedAtSourceCadence` |
| Stage 5 | `stage04.json`，51x24 | `(3,10)`，30 | `(2,10)` | 钩索、锤/钩组合谜题、压力门、分组蛇和爬虫 | `TestRuntimeStage04CanBeCompletedAtSourceCadence` |
| Stage 6 | `stage05.json`，30x75 | `(25,5)`，10 | `(26,5)` | 三段触发、demo 3 镜头脚本、坍塌火炬和追逐火海 | `TestRuntimeStage05CanBeCompletedAtSourceCadence` |
| Stage 7 | `stage06.json`，26x45 | `(17,4)`，35 | `(19,4)`；秘密出口 `(21,40)` raw `28` | 15/20 点宝箱、双向钩石压力门、深层弱墙与秘密路径 | `TestRuntimeStage06NormalExitCanBeCompletedAtSourceCadence`、`TestRuntimeStage06SecretExitCanBeCompletedAtSourceCadence` |
| Stage 8 | `stage07.json`，44x29 | `(6,19)`，40 | `(5,19)`；秘密出口 `(4,3)` raw `28` | 双银锁、落石宝箱、石头压力轴、Freeze Hammer 冻蛇隐藏门 | `TestRuntimeStage07NormalExitCanBeCompletedAtSourceCadence`、`TestRuntimeStage07SecretExitRequiresFreezeHammerAtSourceCadence` |
| Stage 9 | `stage08.json`，35x14 | 无 | 无；raw `53` 徽记宝箱 `(27,6)` | Great Anaconda 三血 Boss、落石命中窗口、尾击、组门和徽记结算 | `TestRuntimeStage08CanDefeatGreatAnacondaAndCollectSealAtSourceCadence` |
| Secret Stage 1 | `stage09.json`，50x30 | `(36,11)`，15 | `(42,11)` raw `28` | 五组压力石路线、金银钥匙、钩索阶梯和秘密关链路 | `TestRuntimeStage09SecretStageCanBeCompletedAtSourceCadence` |
| Secret Stage 2 | `stage10.json`，50x31 | `(42,27)`，40 | `(43,27)` raw `28` | 三把银钥匙、六石分流竖井、四锁联动和长距离推石 | `TestRuntimeStage10SecretStageCanBeCompletedAtSourceCadence` |
| Secret Stage 3 | `stage11.json`，46x31 | `(34,25)`，50 | `(35,25)` raw `28` | 四组敌人门、四把金钥匙、下层锤石迷宫 | `TestRuntimeStage11SecretStageCanBeCompletedAtSourceCadence` |
| Secret Stage 4 | `stage12.json`，46x31 | `(38,27)`，99 | `(39,27)` raw `28` | raw `30` 连锁破墙、双石钩推、单程右竖井和 99 点奖励室 | `TestRuntimeStage12SecretStageCanBeCompletedAtSourceCadence` |
| Tutorial | `stage13.json`，68x11 | 无 | 最终封印脚本，停在 `(61,3)` | `demo.f` 八段脚本、罗盘、检查点重置/召回教学和封印演出 | `TestRuntimeStage13TutorialCanBeCompletedAtSourceCadence` |

Stage 5 在原流程中是取得 Bavaria 神秘钩索后回访的关卡。解码数据明确把 raw `27` 放在 `world1/stage02`（Bavaria Stage 3）`(24,25)` 的 foreground raw `14` 宝箱中；`world0/stage04` 和 `world1/stage04` 均没有 raw `27`。当前实现会在 Bavaria Stage 3 开箱时写入工具等级 `2`，保存后由 Bavaria Stage 5 继承；Angkor Stage 5 的直接关卡入口仍补足该原作前置状态，但绝不能在任一 Stage 5 伪造钩索宝箱。

Stage 8 的秘密出口要求后期取得的 Freeze Hammer。当前 Angkor-only 切片同样缺少这段跨世界获取流程，因此加载 `stage07` 时注入工具等级 `8`；普通出口的正式路线测试仍只使用等级 `2`，证明普通通关不依赖冰冻能力。

## Stage 6 坍塌火炬源码规则

- `i.java` 将 Angkor `stage05` 定义为 `STAGE_ANGKOR_FALLING_TORCHES=5`，进入时设置特殊模式 `kByte=1`、火海高度 `alInt=816`、触发计数 `amInt=0`。
- 第一段触发要求石头完全停在 `(18,63)`；随后穿过 `(18..22,53)` 与 `(7..16,36)` 两组 foreground raw `1`。每组按原递归 `bVoid(...,1)` 整片清除，并把 `bjInt` 重置为 `120`、`amInt` 加一。
- `(7..16,37)` 的 foreground raw `0/background 3` 不是普通计数格。`demo.f` 脚本 `3` 依次执行：方向 `1` 自动上移一格、60 tick 把镜头移动到格 `(12,42)`、等待 20 tick。自动上移会穿过第二组 raw `1`，形成第三段触发；脚本期间禁止输入并暂停火海高度上升。
- 第三段触发启动 `mm0.f` chunk `1` animation `1` 的火炬坍塌，源码 animator 在末帧结束状态切到 animation `2`。火海顶部使用 chunk `0`，初始 animation `2` 完成后切到循环 animation `0`；火海内部平铺 animation `1`。
- `bjInt` 每源码帧递减；警告时镜头下移量为 `bjInt*aSInt%((bjInt>>1)+1)%12`，chunk `2` 绘制下落碎屑。火海启动后 `bjInt==10` 会回到 `60`，因此震动与碎屑循环出现。
- 火海世界 y 为 `stageHeight*24-alInt`。可上升时 `alInt` 每 tick 加一，并至少跟到当前视口底部，最大为 `1704`；当火线到达 `playerY*24+18` 且主角 `x<17` 时调用 `hurtHero(maxHealth,64,1)`，属于致命伤害。
- 检查点保存/恢复 `amInt/alInt` 等价状态；恢复会清零警告、重置火炬与火海启动动画，并按源码移除 `(18,63)` 的触发石头。对应时序、镜头、伤害、恢复与正式路线分别由 `runtime_stage05_test.go` 覆盖。

## Stage 7 普通/秘密出口源码规则

- `stage06` 同时包含 foreground raw `5` `(19,4)` 和 raw `28` `(21,40)`，两者 background 都是方向 `2`。初始化分别写入 Java `pByte/qByte`；进入 raw `5` 设置 `atBoolean=true`，进入 raw `28` 设置 `atBoolean=false`。
- 两个出口都先执行 `xBoolean` 自动行走，按 `jInt=18 -> 12 -> 6 -> 0` 一直走到 `x > width+5`。raw `5` 随后进入 `bByte=35` 的 Loading/成绩页；raw `28` 在普通关走 `bByte=20`，不会显示成绩页或把当前普通关标成已通关。
- `map_angkor.out` 的第一个 type `1` 节点是索引 `9`，因此 Java `dKInt=9`。`zVoid()` 从当前节点的连接里选择 type `1` 且 stage index 更大的节点；Stage 7 的 raw `28` 由索引 `6` 解锁并指向索引 `9`，而普通 raw `5` 解锁索引 `7`。
- `bByte=20` 显示语言文本 `61`：`Congratulations! You have unlocked a secret path!`。源码等待 `aSInt>30` 后重建世界地图并从当前节点移动到秘密节点。
- Go 进度 v6 使用 `StageUnlocked[14]` 逐节点位，不能用 `HighestUnlocked=9` 推导索引 `7/8` 也已解锁。v2 的顺序存档只在迁移时按旧语义展开；v4 引入 `TutorialComplete`，v5 增加 `RelicMask/WorldUnlocked`，v6 增加永久奖励宝箱坐标。秘密关索引 `9 -> 10 -> 11` 仍由地图 type `1` 连接逐个解锁。
- 地图方向输入在 `i.java:1582-1698` 检查节点运行时 bits `0..2` 是否仍为锁定值 `1`，不能检查 authored type bits `3..5`。`i.java:2052-2069` 会把所有已解锁节点（包括 type `1` 秘密节点）的运行时锁状态写成 `0`，所以秘密链解锁后必须允许 `11 -> 10 -> 9 -> 6` 双向返回；未解锁的普通或秘密节点都必须跳过。
- 世界地图标题属于地图上方 `y~=61` 的独立行。底部 `y=275` 状态条按源码使用 `ms.f` chunk `0` frame `12/11/10` 显示生命、全局紫钻和全局红钻；当前关红钻进度是节点附近的独立浮层。不要再把长秘密关标题、关卡计数和生命值塞进底部同一行。
- 普通路线必须取得 35 点：raw `1` 每颗减 1，`(17,19)` 与 `(21,11)` 的 raw `41` 宝箱分别减 20/15。钥匙房石头先从 `(17,8)` 向左钩落到 `(16,11)`，再从 `(18,11)` 向左钩到压力板 `(17,11)`。
- 秘密路线经 `(7,24)` 弱墙群和 `(3,32)` 弱墙进入底层。连续挖掉 `(4,42)/(5,42)` 支撑后，从右侧把右石逐格钩到 `(9,43)`，再把左石推到 `x=1` 保留回路；压力门 `(13,42)` 打开后才能到 raw `28`。两条路线均禁止测试传送和层数据改写。

## Stage 8 Freeze Hammer 与双路线规则

- Java `dInt(x,y)` 把 raw `1` 紫钻映射为冻结类型 `34`，raw `19/43` 蛇映射为类型 `37`；红蛇 raw `43` 另存状态位 `0x10000000`。命中后 player layer 统一变成 raw `9`，类型写在状态位 `22..27`。
- 工具等级 `>=8` 的锤击在命中帧检查中心及上下左右五格：紫钻只在正中心可冻结，蛇可在五格碰撞范围冻结。若中心已命中 raw `30` 或 raw `9`，源码设置动作处理标志并跳过这次五格扫描。
- raw `9` 走与石头/紫钻相同的 `aqVoid()` 重力、滚落、推动和钩索路径，可压机关、压伤主角并压死接触敌人。移动、检查点保存和恢复必须连同原始冻结类型一起转移；不能只替换贴图。
- 再次锤击 raw `9` 会解冻。类型 `34` 还原 raw `1`；类型 `37` 根据红蛇位还原 raw `19/43`，主角在其上方时方向为 `2`，否则为 `1`。
- 冻结紫钻使用 `decoded/sprites/gen0/chunk01-*`，冻结蛇使用 `decoded/sprites/gen1/chunk06-*`，均绘制 frame `0`。Angkor stage index `4/7` 的 raw `5/28` 出口使用目标图的特殊 frame `1`，其他关使用 frame `0`。
- 普通路线先取得两把银钥匙。第二钥匙室必须从 `(26,5)` 把 `(25,5)` 石头左推到石堆上，否则挖掉 `(26,6)` 后石头会滚到宝箱上方，在长开箱动画中反复砸死主角。
- 右区普通机关把 `(33,12)` 石头左推两次，使其落到 `(31,15)`，再从左向右推一次；石头沿 `(32,15..17)` 落到 pressure switch `1`，持续打开 `(34,16)` 的金钥匙区门。该路线只需工具等级 `2`。
- 隐藏路线打通 `(16,5)` 弱墙群和 `(16,2)` 单墙，在蛇回到 `(8,2)` 时冻结。冰块落到 `(8,3)` 后从右侧左推，最终落到 `(7,4)` 的 pressure switch `3`，打开 `(6,3)` 后方 raw `28`；地图从普通节点 `7` 跳到秘密节点 `12`，不顺带解锁普通节点 `8`。

## Stage 9 Great Anaconda 源码规则

- `stage08.json` 是 35x14 的 `STAGE_ANGKOR_GREAT_ANACONDA=8`。Java 初始化 `kByte=4`、`aoInt=0`、`aqInt=3`、`arInt=0`，落石竖井为 `PInt=12/QInt=15`、生成行 `OInt=2`、支撑行 `RInt=5`。
- Boss 三个身体站位的左列是 `10/13/16`。状态 `0 -> 6` 等主角进入 `x>=10`，延迟严格 `apInt>10`；状态 `1` 在 `>20` 后放置 `(bossX..bossX+1,8)` 的 raw `50`，在 `>40` 进入受击状态 `2`。
- 只有状态 `2` 扫描身体两列的 y `8..7`。raw `0` 会被移除并造成一次伤害；落得太早会在状态 `1` 被吞掉但不减血。无命中时 `>15` 切 animation `6`，`>30` 回收身体；命中后 animation `3`、音效 `10`，等待 `apInt>40`。
- 状态 `4` 的回收阈值是生命大于 1 时 10 tick、否则 5 tick。状态 `5` 调用 `acVoid()`：当 `(12,5)/(15,5)` 为空时在对应 y `2` 再生 raw `0`，设置 30 tick 震动，并在身体列 y `4` 放 raw `50`。
- 状态 `10` 在 `LInt=28` 切 animation `7`，`LInt>=50` 进入状态 `11`。尾击持续 12 tick，只检查主角 y `4`、x `bossX-3..bossX+4`，每轮只命中一次；随后状态 `9` 下沉，并按主角 x 选择下一列。
- raw `50` 可穿行，但主角移动偏移 `jInt<12` 且占同格时受伤。状态 `4/5/9/10` 还按身体像素上下边界处理接触伤害，不能把 Boss 简化成静态格子。
- foreground raw `26` `(9,5)` 激活 group `0`：先按来向关闭 `(8,5)`，镜头演出后关闭 `(18,5)`。Boss 第三次受击后进入状态 `7`，严格等待 `apInt>80` 才递减组计数并重新打开右门。
- 身体、平台、尾击和 Angkor 徽记分别使用 `b0.f` chunk `0`、`b0.f` chunk `1`、`gen1.f` chunk `0` animation `2`、`mmv.f` chunk `3`。生命条为屏幕 y `5` 的 `44x12` 黑底和每血 `12x8` 的 `#3bb78f` 段。
- raw `53` 映射徽记 bit `0`，不是按 `51/52/53` 顺序编号。奖励帧立即切人物 animation `47`（42 tick）并锁输入；源码全局计时 `aSInt>140` 后进入状态 `28`，执行 11 个 Loading 步骤，不走普通成绩页。
- 检查点 `(7,5)` 在战斗前。源码复活会把 Boss 重置到状态 `0`、列 `0`、3 血并隐藏尾巴，同时恢复关卡层与 group 计数。
- 正式路线依次把 `(12,5)` 左推到第一列、把 `(15,5)` 左推到第二列、再把再生的 `(15,5)` 右推到第三列；列选择时分别站在 x `13` 和 `16`。禁止测试直接改 `Health/Phase` 代替三次实际落石。

## 四个秘密关源码路线

- `map_angkor.out` 的秘密链是普通 Stage 7 节点 `6 -> 9 -> 10 -> 11`，另一路是普通 Stage 8 节点 `7 -> 12`。秘密关 raw `28` 仍进入普通成绩页并逐节点解锁；末端节点 `11/12` 返回自身，不能按数组下标顺序解锁普通关。
- `stage09` 依次用石头压住 `(9,7)`、`(12,17)`、`(24,23)` 和 `(26,15)` 等机关，取得银钥匙和金钥匙；15 点配额门与出口均在上层右侧。正式测试只需 Mystic Hook，游戏入口注入的 Freeze Hammer 也兼容该路线。
- `stage10` 的右侧钥匙竖井必须把连续六块石头分流到左右支撑，清空 x `43` 的 y `9..15` 后才能开第三把银钥匙宝箱；三把银钥匙依次打开 y `24..26` 的锁，联动打开 `(36,27)`。
- `stage11` 的四个 foreground raw `26` 竞技场分别控制敌人组门；四个 raw `14/33` 金钥匙箱在对应计数归零前可站入但不会开盖，也不能提前显示 payload。入口到下层的首次触发顺序是 group `1 -> 3`，下层迷宫后再走 group `0 -> 2`，因为 Java 只用当前 `cmInt` 递减计数。
- group `0` 把左右蛇冻结、推到落石下，再解冻碾压；group `1` 先把一条蛇冻入 `x=19,y=8`，把另一条固定在 `x=23,y=7`，跟随左落石到 `Remaining<=6` 后解冻，再对 `x=21`/右落石重复；group `2` 从中央分别钩入冰块和落石并用玩家阻止侧滚；group `3` 先在 `x=22,y=18` 放底冰，再把两条活蛇依次钩进同列后一次冻结，通过两次解冻连锁碾压，最后用 `(23,14)` 落石处理幸存冰蛇。出口四锁位于 x `37` 的 y `19..22`，右侧 raw `41` 宝箱补足 50 点配额。
- `stage12` 的关键次序是先锤掉 `(18..19,21)` raw `30`，再钩/推两块石头避免堵死竖井，锤掉 `(17..18,8)` 上层墙，经过检查点 `(19,4)/(25,3)` 收集安全配额，再一次性下右竖井进入 x `40..44` 奖励室。配额为 99，不能先下竖井后期待返回。
- 四关运行入口使用原作后期可达的前置状态：工具等级 `8`、最大生命 `8`。每关测试都保留对象扫描、重力、伤害、宝箱和检查点，不通过直接改层或传送完成。

## Tutorial Stage 13 源码规则

- `stage13` 是 World 0 包中的第 14 个数据关，但不属于 `map_angkor.out` 的 13 个地图节点。新存档先进入它；v1-v3 旧存档迁移时把教程标记为已完成，避免破坏已有进度。
- `demo.f` 由 `c.java` 的命令解释器驱动。教程严格执行脚本 `29 -> 10 -> 11 -> 13 -> 15 -> 16 -> 17 -> 28`；其文本索引序列分别为 `[12,13,14]`、`[0]`、`[19]`、`[1,2,3,4]`、空、`[5,6,7]`、`[8]`、`[9,10,11]`。同一解释器还运行 Stage 3 的脚本 `30`（文本 `17/18`）、Stage 4 神秘锤脚本 `22`（文本 `20/21/22`）和 Stage 9 Boss 入场脚本 `33`（文本 `32/33`）。解码结果落在 `decoded/demo-scripts.json`。
- 脚本命令包括镜头移动、人物自动移动、foreground 写回、等待、对白、头像表情/标点、白闪和最终封印路径。脚本活动时禁止普通输入；`SKIP` 只跳过可视等待和对白，仍执行会改变人物与关卡层的命令。
- raw `42` 是 Compass。`lVoid(42)` 设置 `aFInt=29/aGInt=0`；资源位 `0x80000000` 对应全局 gen chunk `31`，加载器把 `gen3.f` chunk `1` 放入 `textures[29]`。该 chunk 只有一个 `24x24` 模块，必须用于举起奖励；顶部方向箭头仍来自 `ui.f` chunk `2` 的 frame `3+direction`。
- 罗盘奖励触发脚本 `11` 和文本 `19`，同时启用导航。脚本 `13/15` 教 Action/5 在检查点重置；脚本 `16/17` 教 `*` 召回并实际扣一条命。
- 头像、标点和头像底图来自 `demoSpr.bin` 的 sprite `0/1/2`；最终 5x5 封印使用 `mmv.f` chunk `0` 的 frame `4..28`，激活覆盖使用 frame `1`。教程完成后持久化 `TutorialComplete` 并加载 `stage00`。
- 原版在 `(46,7)` 之后跟随人物的图标不是复活点，而是全局 chunk `30`（`gen3.f` chunk `0`）animation `0` 的蓝色 `*` 召回提示。桌面版保留 `(aSInt>>1)&3` 节奏、人物插值位置和 demo 脚本可见条件，但实际绘制 `ENTER` 键帽，不能继续显示手机 `*`；检查点本体仍固定在关卡 foreground 坐标。
- 最终 seal 在人物到 `(61,3)` 且 `jInt==6` 时触发 `cm.f` chunk `7` animation `5` 并隐藏人物；`mmv.f` chunk `0` frame `1` 此时仍必须继续覆盖绘制，不能因人物隐藏而一起消失。

## 锤子和钩索源码规则

- Action/电脑键 `5` 先处理当前检查点，其次处理工具；`*` 是主动召回，不是钩索键。
- foreground raw `2` 是可站入的当前格工具提示，不是碰撞墙。state `0` 在相邻 raw `30` 全部清除后递归消失；state `1` 只有玩家站在该格并按 Action/5、且工具等级 `>=2` 时清除连接块。提示图标按 state 选择锤/钩模块。
- 锤子要求工具等级 `>=1`，方向动画为上/右/下/左 `13/14/15/16`，第 3 tick 命中；向上动作共 11 tick，其余方向 12 tick。它使 raw `30` 进入破坏状态，使普通蛇以低 `0xf8` 位计时 `0x78` 眩晕，并直接移除符合源码条件的未标记 raw `43`。
- 工具等级 `>=8` 时同一命中帧启用 Freeze Hammer 五格扫描和 raw `9` 解冻；增强锤不是独立输入，也不能在 raw `30` 命中后继续冻结相邻对象。相邻移动蛇只有与中心命中格发生像素重叠才冻结；解冻蛇恢复方向后还要带源码 `0x78` 眩晕计时。
- 钩索要求工具等级 `>=2`，只沿水平方向搜索 2 到 3 格。JAR 候选集合为 raw `0/1/8/9/11/14/19/43/47/48`；raw `48` 状态位 `0x8` 时排除，关闭的 foreground raw `7` 以及中途非空对象会阻断。
- 钩索绳段是 player-layer raw `32`，相邻段从 timer `18` 开始，以 `18 -> 12 -> 6 -> 0` 伸展。源码按从左到右扫描，所以向右新建的第二段在同一扫描中立即减到 `12`，向左则保留 `18`。
- 命中后目标不是只移动一格：普通物理目标会逐格拉到主角相邻格；raw `1` 紫钻还会拉进主角所在格，整段动作结束后才收集。每次目标移动使用 18px 插值，残余绳段会在同一扫描中重新捕获目标并把该步 motion 归零。
- 原 JAR 只为 raw `0/8/9/47` 保存钩取前对象状态，并清除 `0x7000/0x200`；其他目标释放时写回 `-1`。因此活蛇被钩后不是恢复旧巡逻方向，而是按负 packed state 的目标 `(127,127)` 重新选择可行轴向。不要为所有候选统一保存/恢复原状态。
- 主角伸出动画为右 `20`、左 `22`，最终拉回动画为右 `21`、左 `23`。钩索期间禁止输入和受伤；命中/重新捕获发出音效 `12`，恢复检查点时必须清理 raw `32`、目标 motion 和钩索状态。

## 宝箱、紫钻与永久坐标

- raw `41` 不是独立分数。初始化时其 background 数值加入关卡紫钻总数 `aYInt`，领奖时同值加入 `aZInt`；Go 端必须同步增加 `VioletGems`，右下角紫钻 HUD、配额、结算分子和最终紫钻资产都读取这条计数链。
- Bavaria Stage 1 的 raw `41` 宝箱位于 `(25,14)`，background authored value 为 `10`。完整开箱后右下角紫钻必须从 `VioletGems` 增加 10，过关时同步进入 `BavariaStageVioletGems[0]` 和 `VioletGemBank`；`TestBavariaPurpleChestAndMysticHookPersistAcrossStages` 与 HUD 绘制回归覆盖这条链。
- player raw `6/7` 分别是额外生命和补血，`OVoid()` 使用 `textures[5]`；加载器把该槽绑定到 `cm.f` chunk `4` 的 module `0/1`（`1UP`/药水）。`gen0.f` chunk `8` 是 Bavaria 火焰陷阱，绝不能用作通用拾取物。raw `6` 增加 HUD 左上的生命数，raw `7` 补满 HP；它们不是紫钻或红钻进度。
- 满血打开 raw `7` 补血箱时，源码在动画开始前把奖励改成 raw `41`/10，因此头顶应显示紫钻和数字 10；生命达到 99 的 raw `6` 先降级成补血，再按是否满血决定是否转成 10 紫钻。
- 新档按每关 foreground raw `14/33` 建立坐标表。World 0 只有已领取红钻、真正增加生命的 raw `6`、以及徽记会移除坐标；钥匙、补血、raw `41`、工具和 Compass 可按源码重玩。普通已消费宝箱重进时保持打开且为空，Boss 已消费徽记箱重进时变成 10 紫钻奖励。
- 每次过关都把本局紫钻全部加入资产，允许重玩刷紫钻；红钻按本局新取得量累加，并依靠永久坐标防止同一红钻重复领取。检查点要快照当前红钻坐标列表；额外生命领奖还要同步修改检查点快照，避免复活后箱子看似关闭却没有内容。

## 门、爬虫与敌人组

- 初始化 raw `7` 门时，按同 group 的 foreground raw `6/8/9` 数量写入低 nibble。每个压力机关或钥匙锁只减一次，计数归零后才进入 `0x10 -> 0x20 -> 0x30` 开门动画；关门保留剩余计数。
- raw `11` 的 `amVoid()` 是贴墙巡逻，不是遇阻简单反向。必须保留前方/侧方角点探测、底到顶且左到右的同帧重复扫描，以及 `18 -> 13 -> 8 -> 3 -> -2` 位移计时。
- foreground raw `17` 是 group marker，raw `26` 才是玩家触发器。Java 只维护当前 `cmInt`；之后任何符合条件的敌人/容器移除都递减当前组，而不是按敌人出生位置保留所有权。Stage 11 的触发顺序因此属于谜题逻辑。

## Bavaria 完整性审计边界

- Bavaria 的可解性由原 JAR 的 authored 三层布局保证，不再维护逐关寻路回归。验收应直接核对 JAR `w1.bin` 解码字节、Java 初始化后的层/状态和实际出现 ID 的更新分支；测试只覆盖容易回归的独立状态机，不通过编写路线证明“没有死胡同”。
- raw `16` authored 底格在初始化时向上复制一格，逻辑、碰撞和销毁始终是两格对象；绘制只从底格执行一次。Java 该分支只把 animation frame 的 `x` 加到绘制坐标，不应用 frame `y`，不能把通用 animation offset 套到铁人上，否则整组会竖直错位。
- 原 JAR 字节码的 stage 初始化 case `12` 只保存 `abInt/acInt/aaInt`，不会清除 player raw `12`；可读 `i.java` 中的清除语句是反编译伪影。Bavaria Stage 4 的配额门必须保持阻挡，紫钻配额归零后才清除。
- Bavaria 专属伤害方向不能统一调用无方向 `Hurt`：raw `28` 尖刺使用 `hByteArr[kInt&7]`，raw `16` 长矛使用自身方向，raw `14` 移动机关使用运动方向。静止人物仍须保留 `kInt` 的最后朝向；恢复检查点时源码固定重置为向右的 `kInt=2`。
- raw `14` 的 cooldown `20` 是等待通路状态；下方和运动方向都堵塞时保持 `20`，通路出现后才从 `19` 倒数。反向机关贴住左侧对象时，原作偶数半帧还有 `(-1,+1)` 像素抖动。
- 爆炸 `bN()` 在水状态稳定 `xByte==3` 时会把邻近 raw `10` 置为清除状态；raw `30/37` 进入破坏状态，raw `8` 连锁爆炸，raw `16/19/43/49` 被移除。水体重排触发源码音效 `13`。
- fan pot 的 foreground `15/16` 与 player `34/35` 只在 `ceInt` 中间相位显示：`1..5` 或 `5..8`；稳定相位 `0/9` 不画罐体。相位 `5` 才交换 player/foreground 层。
- 水状态以 Java `eIntArrArr` 的三个 9-bit 子层保存：每层为 3-bit owner、4-bit shape、2-bit 横向偏移；`aLongArr/bLongArr` 对应 15 个 flow/source packed record。`WaterDepth` 现在只由 packed cells 同步生成，供碰撞和 HUD/渲染做占用查询，不能再作为权威状态写入。
- 水源按列优先扫描并逐个启动；下一水源必须在对象扫描前建立首个子层。phase `1..5`、basin fill、cleanup flow 和 fan reflow 均保存在检查点快照。fan pot 交换、raw `10` 变成 foreground `32`、raw `37` 完成破坏都会调用 `lVoid` 对应的重排入口并播放音效 `13`。
- 浮力直接使用 packed cell `cell != 0 && cell != 3`，不是按水深猜测。layer-0 shape `7/8` 才切游泳动画并屏蔽锤子/钩索；水中 gravity object 使用 `OVoid()` 的 8 帧上下浮动。raw `11` 爬虫入水进入 `0x100..0x400` 溺亡相位，追踪蛇选择垂直方向时避开水格，layer-0 有水会暂停顶石压伤计时。
- `ajVoid()` 只由 player raw `47` 在完成 `aqVoid()` 后调用，不是 raw `1`。raw `47` 会生成 foreground `35`；其 `18 -> 12 -> 6 -> 0` 链继续生成 `35/37/34`、把允许的对象向上搬运，并使用 `gen2.f` chunk `3` 与 `gen3.f` chunk `4` 原素材。当前 Bavaria authored 数据没有初始 raw `47`，但共享源码规则已实现。
- 13 个 stage 的 authored 低位 player/foreground ID 都必须有正式逻辑和绘制分支，不能落入诊断蓝块；地图分支、事件 `4/6/19/34`、Mystic Hook、水药、Evil Teutonic Knight 和封印进度链也必须接入。发布“一比一完成”结论前要完成 JAR 布局、Java 初始化状态、对象逐 tick 转换和资源映射审计，但不要求额外证明原布局是否有解。

## 全局封印世界选择

- Stage 9 徽记演出完成后的 11 步 Loading 进入全局封印选择，而不是直接回 Angkor 地图。界面使用 `mmv.f` 的封印底图/世界覆盖、世界地图图标、原 softkey 和源箭头动画。
- 四个位置为 Angkor、Bavaria、Siberia、Shop，移动图来自 `a_Config.sealMoveTargets`；Bavaria/Siberia 的红钻解锁价为 `10/25`。解锁闪烁、徽记飞入/白闪/特效和箭头 8 tick 插值均有 `worldselect_test.go` 覆盖。
- 进度 v6 持久化 `RelicMask`、`WorldUnlocked[3]`、Angkor/Bavaria 各关节点和永久奖励坐标。Angkor 进入 `map_angkor.out`，Bavaria 进入 `map_scotland.out` 并加载全部 13 个数据关；Siberia 与 Shop 仍是明确的不可用内容，不能描述成完整三世界游戏。Bavaria 必须完成全部实际 ID、共享状态机、Boss、地图和进度链的源码审计后才能宣称第二世界审计完成。

## Angkor 全部关卡素材映射

- raw `12` 配额门：`decoded/sprites/cm/chunk05-*`，两模块组合并叠加 HUD 数字。
- raw `11` 爬虫：`decoded/sprites/gen1/chunk04-*`，6 个模块；正常帧为 `(aSInt >> 1) % 3`，状态 phase 使用后续模块。
- foreground raw `6` 压力机关：`decoded/sprites/gen2/chunk09-*`，单个 `24x13` 模块，底部对齐并随压入量下移。
- `DVoid()` 先从滚动背景缓冲绘制地形，再以相对坐标 `-1..11` 扫描动态格；Go 端不能逐格交替画地面和物体，否则向右滚动的石头、向右喷射的火焰和跨格动画会被后一格地面覆盖。上下左右额外扫描范围也必须保留。
- 游戏内 sprite 必须通过 `*-animations.json` 与 `*-modules.png` 按 module/frame 元数据组合，不能把 `*-frames.png` 假定成固定 `24x24` atlas。Bavaria world chunk `2` 的 frame 最大高度为 `26`，diggable chunk `1` 的 cell 为 `40x41`（Angkor 为 `35x27`）；固定步长会让 `raw 124..129` 的 `2x3` 旗帜和破碎动画跨行错裁。只有没有 sprite 元数据的整屏启动图可直接加载普通 PNG。
- `drsprite` 必须把 `.f` 的 ARGB 调色板颜色预乘后写入 `image.RGBA`，并让 atlas padding 保持 alpha `0`。不能先铺不透明诊断底色再 `draw.Over`，否则透明或半透明边缘会被固化成深色黑边。运行时优先信任 PNG alpha；仅对旧的完全不透明诊断图兼容 `(20,22,28)` 色键。
- `.f` 的 animation-frame `x/y` 不是通用绘制锚点。Java `drawAnimationFrame(..., flags=0, offsetX=0, offsetY=0)` 会忽略它们；只有调用方先执行 `b_SpriteAnimator.applyFrameOffset()` 时才把它们加到 animator 坐标。当前显式使用该语义的是主角、火焰发射器、Stage 6 火炬/火海顶部和 Bavaria Boss；Anaconda、普通对象动画、火海内部平铺、奖励及 UI 动画都不得自动应用 `x/y`。
- foreground raw `20..23`、运行时 raw `32` 消散帧和 foreground `>=80` 属于源码后置前景扫描，应在主角/动态物体之后绘制；raw `255` 是空值，绝不能按 `>=80` 转成 world frame `175`。
- raw `22/23` 共用 `gen1.f` chunk `0` animation `0`。序列第 0 帧持续 20 tick 且为空，表示火焰完全收回；之后 reach 才按 frame index `1..10/11..20/21..` 扩展。raw `23` 水平翻转，不能把合法的收回阶段误判成缺素材。
- foreground raw `6` 的源码绘制会把 Graphics clip 限制在当前 `24x24` 格；下沉模块即使跨过格底也不能漏到下一行。
- raw `32` 钩索：颜色 `#d3d7e7` 的水平线加主角 sprite module `0/1`，不是独立位图。
- player raw `24/31/33/41/50/79` 在 Angkor 关卡中由前景容器、锁、Boss 或入口流程持有可见图形或本身不可见，不得回退绘制诊断蓝块；全部 14 个 stage 的 120 个关闭容器 payload 已审计，raw `2/4/5/6/7/24/26/27/40/41/42/51/52/53` 不得提前显示。
- Freeze Hammer raw `9` 根据保存的源类型选择 `gen0.f` chunk `1` 的冻结紫钻或 `gen1.f` chunk `6` 的冻结蛇，不得绘制成普通石头。
- Stage 6 火海、火炬与碎屑依次使用 `decoded/sprites/mm0/chunk00-*`、`chunk01-*`、`chunk02-*`，不能用通用火焰发射器素材代替。
- Compass 举起图标使用 `decoded/sprites/gen3/chunk01-modules.png` 的 module `0`，不能用 HUD 指针或文字代替。
- 教程人物资源通过 `drsprite -demo-sprites .../demoSpr.bin` 单独导出到 `decoded/sprites/tutorial/demoSpr`；不能把诊断色块当头像。

## 工具命令

```bash
go run ./tools/drdecode -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/w0.bin -out decoded/world0
go run ./tools/drinspect -in decoded -out decoded/preview
go run ./tools/drsprite -in /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources -out decoded/sprites
go run ./tools/drsprite -demo-sprites /Users/wanna/mine/github/wangle201210/DiamondRushSource/src/main/resources/demoSpr.bin -out decoded/sprites/tutorial
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
- 当前 Angkor 14 个数据关、Bavaria 13 个数据关及全局封印选择，与仍未实现的 Shop/Siberia 内容必须分开描述；Bavaria 的素材、地图、进度和机制实现也必须与尚未完成的逐关正式路线差分审计分开描述
- 文档中的“当前实现”与代码一致，不保留已被源码否定的规则
