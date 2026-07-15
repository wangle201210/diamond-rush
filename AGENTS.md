# Diamond Rush Repository Instructions

## 目标与边界

本项目依据本机可运行的 Java ME 原版复刻 Diamond Rush 的可观察行为。Java 源码或原 JAR 能回答的问题，禁止依据旧 Go 原型、记忆、截图、攻略或 Boulder Dash 通用规则补全。

- 保留并继续维护 `originalrush` 运行链路；不要恢复或依赖已废弃的 `diamondrush/` 原型。
- 当前已纳入实现边界的是 Angkor 14 个数据关、Bavaria 13 个数据关和全局封印世界选择。
- Siberia 与 Shop 仍是未实现内容。除非用户明确要求，不要顺带推进世界 3、商店或无关重构。
- 功能实现优先，只为真实行为、回归风险和已确认边界编写必要测试；不要用大量测试脚手架代替功能。

## 指令路由

根文件只放每次任务都必须遵守的约束。按任务类型使用项目级 skill：

- 关卡布局、raw ID、碰撞、对象更新、输入、敌人、机关、Boss、教程、检查点、存档、地图、经济或逐 tick 对照：使用 `$diamond-rush-gameplay-fidelity`。
- 精灵提取、资源槽、chunk/frame/module、动画取帧、锚点、flip、透明度、绘制层级、HUD 或错位/缺图：使用 `$diamond-rush-asset-fidelity`。
- 同时改变逻辑和画面的任务依次使用两个 skill；先确认运行时状态，再处理显示。
- 详细历史结论在 `docs/diamond-rush-fidelity-reference.md`。只读取与当前任务相关的章节，不要把整份深度参考重新塞回根文件。

## 权威资料

- Java 参考项目：`/Users/wanna/mine/github/wangle201210/DiamondRushSource`
- 主游戏逻辑：`DiamondRushSource/src/main/java/i.java`
- 原始可运行字节码：`DiamondRushSource/jars/diamondrush.jar`
- 精灵解析：`DiamondRushSource/src/main/java/f_Sprite.java`
- 动画包装：`DiamondRushSource/src/main/java/b_SpriteAnimator.java`
- 音效 ID：`DiamondRushSource/src/main/java/j_SoundManager.java`
- 世界及经济常量：`DiamondRushSource/src/main/java/a_Config.java`
- 原始资源：`DiamondRushSource/src/main/resources`
- 原版启动脚本：`DiamondRushSource/run-diamond-rush.sh`
- 原版运行参数：`240x320`、缩放 `2`、Nokia 模式、`jars/diamondrush.jar`

仓库内主要参考：

- `docs/diamond-rush-source-mapping.md`：Java 类、状态机、资源和玩法源码定位
- `docs/diamond-rush-remake-spec.md`：实现范围、产品约束和验收条件
- `docs/diamond-rush-fidelity-reference.md`：逐关和素材的详细已验证规则
- `docs/angkor-world0-logic-audit.md`：Angkor 路线、共享机制和跨世界边界
- `decoded/world0/`、`decoded/world1/`：从原始关卡包解出的三层数据
- `decoded/sprites/manifest.json`、`decoded/sprites/index.html`：原始 sprite、帧和动画检查入口

运行原版：

```bash
cd /Users/wanna/mine/github/wangle201210/DiamondRushSource
./run-diamond-rush.sh
```

## 证据优先级

1. 原始 `diamondrush.jar` 字节码。遇到反编译出的不可能循环、缺失 `break` 或矛盾控制流时，用 `javap -p -c` 核对原 JAR，禁止使用重新编译的 `target/classes` 证明行为。
2. `i.java` 中的加载、更新、碰撞、绘制和状态切换代码。
3. FreeJ2ME 原版的同输入逐帧运行结果。
4. 原 JAR 内的 `w0.bin`、`w1.bin`、地图、`.f`、语言和音频资源。
5. 截图和视频只能确认用户可见结果，不能证明碰撞条件、更新顺序、计时器、检查点快照或状态位语义。
6. 现有 Go 代码、旧原型和攻略只能作为定位线索。

`src/main/resources` 下的 `w0.bin`、`w1.bin`、`w2.bin` 均与原 JAR 内同名条目存在差异，`decoded/worldN` 必须保留 JAR 版本，不能用源码目录资源重新生成后覆盖：

- `w1.bin` 有 4 字节差异，均在 Bavaria `stage07` player layer。原 JAR 的 `(20,10)/(29,19)/(25,20)` 为空，`(7,17)` 为 raw `10`。JAR 中 `w1.bin` SHA-256 为 `951b998c82383c55144ed82c5c54a7dc70f638017929d46aa155e40b0a77674e`，`map_scotland.out` 为 `5c21ffc3ac32e6f571cba097eaf81f3a7044804d1e4ae19f2c381586eba543c0`。
- `w0.bin` 有 16 字节差异，均在 Angkor `stage00`：JAR 版第二检查点在 `(19,7)`、出口在 `(21,9)`、`(19,9)` 为 foreground raw `30` 脚本触发格（脚本 `20`）。JAR 中 `w0.bin` SHA-256 为 `0b2eb7662959fd1e73b8ee435e96ca9cff87bd159dec845640954ed9658c0ffd`。详见 `docs/diamond-rush-fidelity-reference.md` 数据完整性一节。
- `w2.bin` 亦有差异，尚未逐字节审计；实现 Siberia 前必须先从 JAR 解码。

## 始终生效的运行时约束

- 逻辑画布为 `240x320`，上下 HUD 各 `40px`，格子为 `24x24px`。
- Java 源码主循环保持 `20 TPS`。Ebitengine 外层可用 `60 TPS` 采集输入和刷新显示，但完整源码步骤只能每 3 次 `Update` 执行一次。
- 菜单、地图、教学、结算、`g.tick`、`TickSourceFrame()`、碰撞、伤害和存档都属于 20 TPS 逻辑。显示插值不能写回 Runtime，也不能改变运动速度或触发时序。
- 切关、复活、召回和脚本切镜必须重置显示插值；逻辑相机仍由源码状态驱动。
- Draw 路径禁止创建 `ebiten.Image` 或重复构造 sprite `SubImage`；纯色图和精灵子图在加载期缓存。
- 保持 player、background、foreground 及其高位状态、对象 timer 和 wake/dirty 状态相互独立。禁止压成单一简化 tile type。
- Java `byte` 有符号；JSON raw `255` 对应 Java `-1`，通常为空对象。
- 动态对象活动区按源码从下到上、再从左到右扫描。一次扫描内对象是否再次被处理必须按源码或 JAR 验证。
- `ObjectMotion`、相机插值和 `renderPhase` 只服务显示；权威碰撞坐标、状态位和 timer 始终来自 20 TPS Runtime。

## 桌面输入与启动契约

- `Space`、手机 `5`、数字小键盘 `5`：交互/工具。
- `Enter`、`R`、`Backspace`、数字小键盘 `*`、`Shift+8`：手机 `*` 语义，主动召回到检查点并按规则扣命。
- `Tab`：桌面版已解锁世界/关卡导航；不能进入锁定节点。
- `S`：仅跳过教程可视等待和对白，不跳过会修改人物或关卡层的脚本命令。
- 上述动作必须互斥。面向玩家的提示统一显示 `SPACE`、`ENTER`、`TAB`、`S`，不显示手机键位。
- 静止时首次按下不同于当前朝向的方向键只转身；再次按同方向或持续按住后才移动。该门控不适用于入口、出口和 demo 脚本自动移动。
- 有存档时先显示 `Continue` / `New game`。新游戏必须经过 `No` / `Yes` 二次确认，只有 `Yes` 可以覆盖进度。Continue 回到最高已解锁地图节点；教程未完成则继续教程。

## 实现与验证纪律

1. 先从解码数据和 Java/JAR 确认事实，再改 Go；不要先从截图猜结论。
2. 先实现纯 Runtime 状态转换，再连接输入、相机、绘制和本地化。
3. 测试必须覆盖已确认的源行为和回归风险，不得通过传送、暂停对象更新、直接改层数据或直接改 Boss 血量来证明正式路线可通关。
4. 在 dirty worktree 中只修改任务相关文件；保留用户或其他任务已有改动。
5. 改变已记录的源行为、实现边界或验收方法时，同步更新相应文档或 skill reference。

Go 测试通过只能证明实现符合测试假设。没有 Java/FreeJ2ME 同输入逐源码 tick 对照时，不得宣称“完全一致”“一比一复原”或“该世界完整复刻”。

## 常用验证

针对原作运行时改动至少执行：

```bash
go test ./internal/original ./internal/originalgame
go build -o /tmp/originalrush-smoke ./cmd/originalrush
rm -f /tmp/originalrush-smoke
```

涉及解码器、资源工具或共享包时再执行：

```bash
go test ./internal/original ./internal/originalgame ./tools/drdecode ./tools/drinspect ./tools/drsprite ./tools/drsound ./tools/drworldaudit
go test ./...
```

macOS 无主显示器环境中，Ebitengine/GLFW 可能在项目测试前初始化失败；区分环境失败与玩法回归。仅修改文档、`AGENTS.md` 或 skill 时，不需要重启游戏或运行 Go 测试，但必须校验 skill 格式和仓库内链接。

## 完成定义

单关或单机制只有在以下条件满足后才可标记完成：

- 所有实际出现的三层 raw ID 均有源码/JAR 定位和明确语义。
- 碰撞、扫描顺序、timer、状态位、伤害、检查点恢复和退出路径已覆盖。
- 正常路线、死亡/召回路线、遗漏收集路线和关键边界操作可执行。
- 关键行为有 Java 与 Go 的同输入逐 tick 差分证据。
- 原作 frame、anchor、flip、层级、相机和 HUD 时序已核对。
- 文档中的“当前实现”与代码一致，并清楚区分已实现、已测试、已做源码审计和仍未实现。
