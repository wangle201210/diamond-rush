# Diamond Rush Original Runtime

这是一个使用 Go 和 Ebitengine、依据本机 Java ME 原版源码与资源实现的 Diamond Rush 运行时。当前覆盖 Angkor 和 Bavaria；Siberia 与 Shop 尚未实现。

## 运行与打包

开发环境需要 Go 1.24 或更高版本：

```bash
go run ./cmd/originalrush
```

为 Apple Silicon Mac 生成包含可执行文件和运行资源的发行包：

```bash
make package
```

产物位于 `dist/DiamondRush-macos-arm64.zip`。该构建不能在 Windows、Linux 或 Intel Mac 上直接运行。

## 按键

| 操作 | 桌面按键 |
| --- | --- |
| 移动 | 方向键、主键盘数字键或数字小键盘 |
| 交互、工具、确认 | `Space`、`5`、数字小键盘 `5` |
| 召回、返回 | `Enter`、`R`、`Backspace`、`Shift+8`、数字小键盘 `*` |
| 世界或关卡导航 | `Tab` |
| 跳过教程等待和对白 | `S` |
| 退出 | `Esc` |

人物静止时，首次按下与当前朝向不同的方向只会转身；再次按下或持续按住才会移动。

## 实现范围

- Angkor：14 个数据关，包括普通路线、Boss、秘密关和教程。
- Bavaria：13 个数据关，包括普通路线、Boss 和秘密关。
- 世界地图、进度、检查点、工具和两个世界的特殊机制已接入。
- Siberia 与 Shop 不在当前实现范围内。

详细的实现边界和验收条件见 [`docs/diamond-rush-remake-spec.md`](docs/diamond-rush-remake-spec.md)，源码与资源定位见 [`docs/diamond-rush-source-mapping.md`](docs/diamond-rush-source-mapping.md)。

## 仓库结构

```text
cmd/originalrush/       Ebitengine 程序入口
internal/original/      三层关卡数据与玩法运行时
internal/originalgame/  输入、绘制、音频、地图、结算与进度
decoded/                解码后的关卡、地图、精灵、字体和音频
tools/                  资源解码、检查与审计工具
docs/                   源码映射、规格和审计记录
```

macOS 存档默认位于：

```text
~/Library/Application Support/zskc-diamondrush/original-progress.json
```

## 开发验证

```bash
go test ./...
go build -o /tmp/originalrush-smoke ./cmd/originalrush
rm -f /tmp/originalrush-smoke
```

开发时可以用环境变量直达关卡：

```bash
ORIGINALRUSH_STAGE=1 go run ./cmd/originalrush
ORIGINALRUSH_WORLD=bavaria ORIGINALRUSH_STAGE=1 go run ./cmd/originalrush
```

完整的证据优先级、运行时约束和验证纪律见 [`AGENTS.md`](AGENTS.md)。
