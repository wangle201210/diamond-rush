# Diamond Rush Original Runtime

这是一个使用 Go + Ebitengine、依据本机 Java ME 原版源码与资源实现的 Diamond Rush 原作数据运行时。当前可玩范围覆盖 Angkor World 0 打包的全部 14 个 stage，以及 Bavaria World 1 打包的全部 13 个 stage。Angkor 的 `stage00` 到 `stage08` 是普通路线与 Great Anaconda Boss，`stage09` 到 `stage12` 是四个秘密关，`stage13` 是新存档进入 Stage 1 前的原版教程；Bavaria 的 `stage00` 到 `stage09` 是普通路线与 Evil Teutonic Knight，`stage10` 到 `stage12` 是秘密关。两个世界均已接入原版世界地图、素材、进度和各自特殊机制；Bavaria 的逐关完整路线差分审计仍在进行。Siberia 与 Shop 尚未复刻。

## 运行

需要 Go 1.24 或更高版本。在仓库根目录执行：

```bash
go run ./cmd/originalrush
```

游戏逻辑画布固定为 `240x320`，Java 源码状态机保持 `20 TPS`。桌面端以 `60 TPS` 采集输入和刷新显示，每 3 次显示更新执行 1 次权威源码步骤；人物、敌人、落石、机关和相机的中间坐标只用于绘制，不会改变碰撞、速度、事件 tick 或存档。桌面窗口可以缩放，但不会改变游戏逻辑尺寸。

## 按键

| 原手机键 | macOS 键盘 |
| --- | --- |
| `2` / `4` / `6` / `8` | 方向键、主键盘数字键或数字小键盘 |
| `5`（宝箱、锤子、钩索、机关、菜单确认） | `Space`、`5` 或数字小键盘 `5` |
| `*`（关卡内返回复活点；菜单中返回） | `Enter`、`Shift+8`、数字小键盘 `*`、`R` 或 `Backspace` |
| 世界/关卡导航 | `Tab` |
| 教程 `SKIP` | `S` |
| 退出 | `Esc` |

`5` 是检查点及锤子/钩索的上下文动作键；关卡内的 `*` 用于主动召回检查点。游戏关卡内按 `Tab` 返回当前世界地图以选择已解锁关卡；世界地图内再按 `Tab` 进入世界选择。用方向键选择已解锁的 Angkor/Bavaria 或关卡，再按 `Space` 进入。锁定内容不能进入；世界选择内按 `Enter` 返回主菜单。

人物静止时，如果按下的方向和当前朝向不同，第一次只会原地转身。转身后再次按同方向会前进；持续按住方向键则会在原作的短暂转身时序结束后自动前进。

## 仓库结构

```text
cmd/originalrush/       Ebitengine 程序入口
internal/original/      原始三层关卡数据与源码时序玩法运行时
internal/originalgame/  输入、绘制、音频、世界地图、结算与进度
decoded/                解码后的关卡、世界地图、精灵、字体和音频
tools/                  原始资源解码、检查与审计工具
docs/                   源码映射、玩法规格和 World 0 数据清单
AGENTS.md               后续关卡复刻的证据、流程与回归基线
```

macOS 进度默认保存到：

```text
~/Library/Application Support/zskc-diamondrush/original-progress.json
```

当前存档格式为 v6，持久化 Angkor/Bavaria 节点、收集/勋章、教程、三枚徽记、世界解锁状态、工具等级，以及原作会永久消费的红钻/额外生命/徽记宝箱坐标。raw `41` 宝箱奖励计入右下角紫钻、关卡结算和可重玩累加的紫钻资产；Bavaria Stage 3 的 raw `27` 宝箱取得 Mystic Hook，之后由存档带入 Stage 5。

## 权威参考

本机 Java 参考项目位于：

```text
/Users/wanna/mine/github/wangle201210/DiamondRushSource
```

实现判断以原始 JAR 字节码、`i.java`、FreeJ2ME 运行结果和解码资源为依据。详细流程见 `AGENTS.md`，源码位置见 `docs/diamond-rush-source-mapping.md`，World 0 完整审计见 `docs/angkor-world0-logic-audit.md`。

## 验证

```bash
go test ./...
go build -o /tmp/originalrush-smoke ./cmd/originalrush
```

开发时可用 `ORIGINALRUSH_STAGE=1..14 go run ./cmd/originalrush` 直达 Angkor 数据关，`14` 是教程；使用 `ORIGINALRUSH_WORLD=bavaria ORIGINALRUSH_STAGE=1..13 go run ./cmd/originalrush` 可直达 Bavaria。该覆盖不改变正常启动选择，但完成关卡仍会写入当前 `HOME` 下的进度文件。

常用资源工具命令记录在 `AGENTS.md` 和 `decoded/README.md`。

## Wule资源Profile

西游主题Wule资源与原版资源并行存放，不覆盖`decoded/`。当前`go run ./cmd/originalrush`仍固定运行原版；Wule profile尚未生产最终素材，标记为`ready=false`，不会拿原版图片补缺。

资源规划和场景Brief可重复生成并校验：

```bash
go run ./tools/pilgrimmanifest
go run ./tools/pilgrimmanifest -check
```

双profile根目录、独立存档命名空间和禁止fallback规则见`assets/resource-profiles.json`；Wule侧说明见`assets/pilgrim/README.md`。
