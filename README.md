# Diamond Rush Original Runtime

这是一个使用 Go + Ebitengine、依据本机 Java ME 原版源码与资源实现的 Diamond Rush 原作数据运行时。当前可玩范围是 Angkor World 0 的前五关，关卡直接读取 `decoded/world0/stage00.json` 到 `stage04.json`。

## 运行

需要 Go 1.24 或更高版本。在仓库根目录执行：

```bash
go run ./cmd/originalrush
```

游戏逻辑画布固定为 `240x320`，主循环为 `20 TPS`；桌面窗口可以缩放，但不会改变游戏逻辑尺寸。

## 按键

| 原手机键 | macOS 键盘 |
| --- | --- |
| `2` / `4` / `6` / `8` | 方向键、主键盘数字键或数字小键盘 |
| `5` | `5`、数字小键盘 `5` 或 `Enter` |
| `*` | `Shift+8`、数字小键盘 `*`、`R` 或 `Backspace` |
| 退出 | `Esc` |

`5` 是检查点及锤子/钩索的上下文动作键；`*` 用于主动召回检查点。

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

## 权威参考

本机 Java 参考项目位于：

```text
/Users/wanna/mine/github/wangle201210/DiamondRushSource
```

实现判断以原始 JAR 字节码、`i.java`、FreeJ2ME 运行结果和解码资源为依据。详细流程见 `AGENTS.md`，源码位置见 `docs/diamond-rush-source-mapping.md`。

## 验证

```bash
go test ./...
go build -o /tmp/originalrush-smoke ./cmd/originalrush
```

常用资源工具命令记录在 `AGENTS.md` 和 `decoded/README.md`。
