# 西游主题高清素材逐项替换目录

## 1. 文档目标

本文是 [`journey-to-west-theme-replacement-spec.md`](./journey-to-west-theme-replacement-spec.md) 的制作级清单。总体规范回答“为什么替换、保持什么不变”，本文回答“每一个旧资源槽具体由什么新资源承接，或依据什么源码证据删除”。

整体世界观、悟乐（Wule）角色定义、三地因果和世界法则以总体规范第3.3节为准。任何单件素材brief若与该节冲突，必须先修正文档，不能在生成提示词中自行发明另一套设定。

第一阶段继续冻结现有关卡尺寸、三层数据、对象坐标、机关关系、Boss状态机、教程命令时序和世界地图连接。本文只决定视觉、音频、字体、Logo和可见文本的替代内容。

完成定义：

- `decoded/sprites/manifest.json` 中24个文件、121个块全部拥有唯一的 `Rxxx` 结论。
- 结论只能是“生成新资源”“重写为新数据”或“删除且禁止回退”，不能保留“以后再看”。
- 94个可解析精灵块按运行时语义替换；27个非精灵块按图片、脚本、文本或音频的真实格式处理。
- `.f` 内的module只是组成frame的切片，不视为独立美术概念。制作单位是具有独立语义的frame、animation或连续frame区间。
- 原清单共1432个module、1119个frame、188个animation。新资源不要求继续使用同样的module拆法，但必须覆盖本文列出的每个可达语义槽。
- `legacy_slot` 只用于迁移核对，不得成为发布资源名。
- 第2.6至2.8节必须恰好覆盖World 0的14关、World 1的13关和World 2的14关；每项资产生成时都要引用至少一个scene ID，不能只引用世界色板。

## 2. 全局制作合同

### 2.1 尺寸档位

| 档位 | 逻辑占用 | 高清母版和导出规则 | 默认锚点 |
| --- | --- | --- | --- |
| `T1` | 1格 `24x24` | 原生绘制为`96x96`；地形不留透明边 | 左上`(0,0)` |
| `T1O` | 1格对象 | 最小`96x96`，允许透明画布到`192x192`，最终按内容裁切 | 格中心`(48,48)`或脚底`(48,96)` |
| `T2V` | 纵向2格 | `96x192`，多格接缝位于`y=96` | 左上`(0,0)` |
| `T2H` | 横向2格 | `192x96`，多格接缝位于`x=96` | 左上`(0,0)` |
| `T2X2` | `2x2`格构件 | `192x192`，接缝位于`x/y=96`；先画整体再切四象限 | 左上`(0,0)` |
| `HERO` | 主角占1格、表现可越界 | 单帧允许包围盒`288x320`；透明裁切并记录pivot | 逻辑脚底`(12,24)`，高清`(48,96)` |
| `FX` | 视觉效果 | 单格效果建议`192x192`；横向火焰可到`384x192` | 逻辑发射点乘4 |
| `BOSS` | 多格对象 | 按源码占用格导出trimmed frame，最大边界写入manifest | 源码Boss原点乘4 |
| `HUD` | `240x40` | 每条HUD为`960x160`，图标另导出矢量或高清PNG | 屏幕坐标乘4 |
| `SCREEN` | `240x320` | `960x1280`全屏母版 | 左上`(0,0)` |
| `PORTRAIT` | 教程头像区 | `416x180`安全画布，人物不得被文本框裁切 | 左上`(0,0)` |

所有对象先在20 TPS逻辑坐标中更新，再在绘制层乘以4。美术不得改变命中、开门、发奖、伤害和状态切换tick。

### 2.2 文件合同

每个生成项至少输出：

```text
assets/pilgrim/<domain>/<asset>/atlas.png
assets/pilgrim/<domain>/<asset>/atlas.json
assets/pilgrim/<domain>/<asset>/source/<working-file>
```

`atlas.json` 必须记录：`asset_id`、`legacy_slot`、frame名、animation名、frame duration、pivot、visual bounds、逻辑事件tick和来源记录。静态全屏图可使用单张PNG；Logo同时保留可编辑矢量母版。

### 2.3 统一美术语言

- 风格：高清国风彩绘2D，清晰轮廓、受控纹理、统一左上光源，不使用低清像素图放大。
- 主角：悟乐（`Wule`，资源slug `wule`），一名原创年轻石猴行者；深棕短毛、青绿色短褂、朱红腰绳、乌金护腕，不采用现代影视或游戏的特定孙悟空造型。
- 花果山：青灰岩、翠藤、瀑水、暖金石灯和少量朱红封印。
- 火焰山：赤岩、黑铁、青铜、熔火、玉泉水和芭蕉风纹。
- 雪岭古道：蓝白冰岩、深灰古寺、朱红经幡、青玉寒光。
- 可交互对象必须比同格背景至少多一层轮廓或明度差；不能只靠红绿颜色区分。

#### 2.3.1 世界观制作摘要

- 悟乐（`Wule`）是从花果山废弃云路碑的香火中化生的年轻石猴行者，性格机敏、乐观，以观察和借力解决机关；所有角色资产统一使用`character_id=wule`。
- 云路连接三地洞天。`云路印`让驿站临时接续路线，`经卷灵印`长期稳定区域，两条进度互不替代，因此场景不能把Boss灵印画成下一世界的通行费用。
- 花果山是天然地脉洞府，火焰山是风、水、火协同的古代工程枢纽，雪岭是保存照妖镜和镇妖设施的古寺群；三地机关必须分别服从这三套设施来源。
- 土地神只能借固定神龛、头像和寻经针引导悟乐；救命毫毛负责召回，行囊坊用灵蕴珠强化护体毫光，不新增伙伴、法力值或第四种货币。
- 整体基调是轻快但有危险的解谜冒险。敌人可凶猛，受击和退场不使用写实血腥；悟乐也不能被塑造成傲慢、暴戾的战神。

完整设定及不可变叙事边界见总体规范第3.3节；本摘要只用于让美术、动画、音频和提示词在制作时快速对齐。

### 2.4 场景先于单件素材

本目录不是“旧PNG名称 -> 新PNG名称”的平面换表。每项素材生成前，必须先把它放回真实关卡数据和状态机中审查。场景语义的证据顺序为`stageNN.json`三层数据、Java初始化/更新/绘制分支、地图连接，最后才是运行画面检查。

每个关卡必须在`assets/pilgrim/manifests/stage-scenes.json`登记：

```text
scene_id, world, stage_index, map_node_type, width, height,
terrain_frames, foreground_pairs, dynamic_raw_ids, boss_mode,
entry_view, mechanic_views, reward_view, exit_view,
world_lore_role, local_conflict, infrastructure_system,
base_materials, accent_materials, light_source, ambient_fx,
asset_variants, removed_on_init_ids, forbidden_motifs, review_status
```

制作顺序固定如下：

1. 从该关三层数据生成纯几何占位图、raw ID坐标表和地形frame复用表，不读取旧截图作为新构图依据。
2. 对入口、首次出现的新机制、关键宝箱、检查点、出口或Boss区各制作一张`10x10`逻辑格场景板；场景板必须使用实际坐标遮罩，不能凭空重排房间。
3. 先在场景板中确定主材质、光源、危险色、可交互层级和多格构件，再生成单件母版。
4. 同一frame若出现在多个scene，纹样只能采用这些scene的共同语义。只属于单关的壁画、Boss标记或法宝符号，只能放在该关独占frame或纯视觉overlay中。
5. 多格旗幡、灯柱、神像、门、钉柱、持矛力士、封印和Boss必须先生成完整构图，再按源码原点切片。禁止逐格调用图像模型后拼接。
6. 关闭宝箱、未开启门、停用机关、危险预警和激活状态都要单独放回场景板检查，不能只验收一张静态美术图。
7. 最终用新atlas完整渲染该关，并在入口、所有检查点、奖励区、普通出口、秘密出口和Boss各阶段验收。单件PNG好看但全关不成立，仍视为失败。

世界或关卡变体只能改变绘制资源、灯光、粒子和纯视觉overlay。它们不得修改raw ID、碰撞、对象状态、事件tick或通关条件。渲染器按`scene_id`选择variant，缺少指定variant时构建失败，不能静默回退到另一世界或旧素材。

### 2.5 跨场景视觉语法

| 玩法职责 | 不可改变的视觉信息 | 场景化方式 |
| --- | --- | --- |
| 可滚落物 | 接近圆形的重心、8个旋转状态、与固定墙体不同的外轮廓 | 花果山镇山石、火焰山镇火石、雪岭玄冰石可换材质，重心和占格不变 |
| 可破坏墙 | 裂纹和受击阶段必须先于清除可读 | 三世界分别用岩、铜砖、冰寺砖；完整墙不能看起来像空洞或藤障 |
| 压阵盘 | 可被石头、冻结敌人等压住，下沉方向明确 | 盘面纹样按世界换材质，中心受力区保持相同轮廓 |
| 门和符锁 | 关闭时形成连续阻挡，开启方向和计数状态可读 | 石门、铜门、冰门共享结构语法，不共享同一材质贴图 |
| 经匣 | 关闭时完全遮住payload，开盖后才出现奖励光 | 外壳可随场景换木、铜、玉，不能在盖面画出内部奖励图案 |
| 方向危险物 | 发射口、刃尖、运动朝向必须由轮廓而非纯颜色表达 | 地火兽口、风火轮、长刺、寒光镜分别使用明确机械轴线 |
| 普通收集物 | 灵蕴珠是圆形紫青灵气，承担关内配额和商店货币 | 可换环境反光，不能改成钥匙、药品或路线信物轮廓 |
| 稀有收集物 | 云路印是朱砂方印，永久累计到`10/25`显现新云路 | 不画回复药效，不在拾取文案中使用“服用”或“消耗” |
| 主线信物 | 经卷灵印是Boss级大型印章，每世界一枚 | 花果山水石纹、火焰山芭蕉火纹、雪岭雪莲纹，尺寸和演出高于云路印 |
| 旧Crystal对白 | `Hidden/Fire Crystal`统一指向对应世界的经卷灵印 | 不建立独立经页、第四枚信物或额外计数；三枚灵印仍分别使用水石、芭蕉火和雪莲压印 |
| 工具进程 | 如意铁棒 -> 幌金绳加入工具能力 -> 同一铁棒玄冰淬炼 | 第三级保留原棒身识别点，奖励演出表现强化而非换一根武器 |
| 敌对猿类 | 必须一眼区别于悟乐 | 雪岭猿傀使用整脸木面具、符绳关节和白灰长毛，不做悟乐换色 |
| 水、火、冰 | 水可游、火有reach、冰冻物可滚落的功能要可预判 | 粒子可高清化，但不能用大面积发光遮住格边和对象状态 |

### 2.6 花果山场景矩阵

以下名称是新主题的场景brief，不是对原关卡名称的翻译。

| scene ID | 数据/地图身份 | 新场景名 | 代码与数据约束 | 生成和布景要求 |
| --- | --- | --- | --- | --- |
| `W0-TUTORIAL` | `stage13.json`，非地图教程 | 山门启行 | 寻经针、土地神龛、召回教学、脚本镜头和最终5x5封印 | 清晨山门到洞口逐步变亮；键帽和神龛是旅途设施，最终封印必须是场景中心而非贴在地上的普通花纹 |
| `W0-01` | `stage00.json`，普通Stage 1 | 水帘洞前庭 | 镇山石、灵蕴珠、蛇妖、左向地火口、经匣和三个神龛 | 瀑水冷光为主，地火口只出现在岩缝泄压点；入口、宝箱和出口保持清楚，藤障不能与蛇妖混色 |
| `W0-02` | `stage01.json`，普通Stage 2 | 藤根回廊 | 更长落石追逐、三条蛇和四个有序神龛 | 用连续树根和石梁强调纵深；落石路线保持高对比，装饰不得盖住支撑是否已挖空 |
| `W0-03` | `stage02.json`，普通Stage 3 | 日月符锁殿 | 金银符钥/锁、raw `30`裂墙、敌人组门和火口 | 日纹与月纹贯穿钥匙、锁和门组；竞技场门看起来属于同一镇妖系统，不能像随机木门 |
| `W0-04` | `stage03.json`，普通Stage 4 | 裂岩练功场 | 压阵盘、raw `24`如意铁棒、破墙和蛇眩晕 | 法宝经匣位于练功石台语境；裂墙有棒击痕，压阵盘嵌地，取得铁棒前后场景识别点保持一致 |
| `W0-05` | `stage04.json`，普通Stage 5 | 金绳回环洞 | 要求持久化工具等级`>=2`；标准来源是火焰山Stage 3幌金绳，提前进入雪岭触发的等级`2`兜底也可满足；含锤/绳组合、压力门、蛇和石蝎 | 洞内预先出现绳结石刻和远处不可达回环，解释回访价值；绝不能伪造幌金绳宝箱或把兜底表现成此关发奖 |
| `W0-06` | `stage05.json`，普通Stage 6 | 地火脉竖井 | 30x75纵向关、三段触发、脚本镜头、崩塌石灯和上升火潮 | 上层湿岩逐渐过渡到下层赤亮地火；完整石灯先设计再切片，火潮与普通火口不能共用同一视觉规模 |
| `W0-07` | `stage06.json`，普通Stage 7，含秘密出口 | 深根双岔洞 | 35点配额、钩石压力门、弱墙、普通/秘密双路线 | 根系在路线分叉处形成明确地标；秘密下层更幽暗但通道仍可读，普通界门和隐云界门轮廓不同 |
| `W0-08` | `stage07.json`，普通Stage 8，含冻结秘密出口 | 寒泉古阵 | 普通路线只需幌金绳；秘密路线需雪岭Stage 6后的玄冰铁棒冻结蛇并压阵 | 在花果山材质中加入局部寒泉和青玉旧阵，解释冰冻回访；不能把整关直接换成雪岭，也不能提前提示唯一解 |
| `W0-09` | `stage08.json`，Boss Stage 9 | 蟒王祭坛 | 碧鳞蟒王三列身体、落石命中窗、尾击、组门和经卷灵印 | 三个身体列对应祭坛三座破损石台；顶部石头表现为祭坛配重，只有易伤阶段出现明确腹甲开口，奖励台在战斗前不可抢眼 |
| `W0-S1` | `stage09.json`，秘密关链首关 | 五阵石窟 | 五组压力石、金银符钥、钩索阶梯和15点配额 | 五个阵区各有独立石刻编号但共用材质；钩索路线留出清楚水平视线，不加遮挡绳段的前景根系 |
| `W0-S2` | `stage10.json`，秘密关链第二关 | 符钥分流井 | 三把银钥匙、六石分流竖井、四锁联动和40点配额 | 竖井左右用日月刻线提示分流，不使用文字答案；六块石头必须始终与背景墙体分离 |
| `W0-S3` | `stage11.json`，秘密关链第三关 | 四门镇妖窟 | 四组敌人门、四把金钥匙、冻结/解冻与下层锤石迷宫 | 四个竞技场使用四种边框符号标识当前组，敌人仍共享同一物种美术；锁住的经匣不显示金符钥payload |
| `W0-S4` | `stage12.json`，独立秘密分支 | 九十九灵蕴宝库 | raw `30`连锁墙、双石钩推、单程右竖井和99点奖励室 | 上层是克制的石库，右侧奖励室才提升金色密度；单程竖井边缘明确但不画返回箭头，99点门保持功德封印语法 |

### 2.7 火焰山场景矩阵

| scene ID | 数据/地图身份 | 新场景名 | 代码与数据约束 | 生成和布景要求 |
| --- | --- | --- | --- | --- |
| `W1-01` | `stage00.json`，普通Stage 1 | 赤岩铜门 | raw `16`两格持矛铁人、raw `28`钉柱、裂墙、火蜥和敌人组门；承接已取得的如意铁棒 | 入口是山体外缘到铜门内廊；镇关力士、钉柱和门使用同一青铜铸造体系，力士必须用抬矛、警示底座和运动姿态表明敌意，不能像无害装饰佛像 |
| `W1-02` | `stage01.json`，普通Stage 2 | 黑铁机廊 | 石蝎、巡游风火轮、持矛力士、钉柱和大量裂墙 | 黑铁导轨贯穿移动物轨迹，安全地面保持暗哑；火花只出现在运动轴和撞击点，不满屏铺火 |
| `W1-03` | `stage02.json`，普通Stage 3 | 幌金绳藏库 | raw `27`位于`(24,25)`经匣，另有移动物、持矛力士、蛇妖和裂墙 | 藏库以金红绳结壁刻建立法宝预告；真正幌金绳只在开箱tick出现，盖上不得画绳结形payload剪影 |
| `W1-04` | `stage03.json`，普通Stage 4 | 风火巡廊 | 巡游风火轮、持矛力士、钉柱、敌人组门和配额 | 用导风槽把机关方向串起来；`2x2`灯龛保持完整，配额门嵌在主路而非像临时HUD贴片 |
| `W1-05` | `stage04.json`，普通Stage 5 | 阴阳风瓮殿 | 风脉阵枢、阴阳风瓮、地火虫穴、钉柱、锁和组门 | 日/月纹从阵枢延伸到两种风瓮，稳定相位和交换相位都要像同一装置；不得用纯红蓝方块表示状态 |
| `W1-06` | `stage05.json`，普通Stage 6 | 雷火丹炉窟 | 雷火炉石、爆裂符墙、风脉阵枢和阴阳风瓮 | 炉石是可滚圆物，爆裂墙是固定墙；白热裂纹强度分级，避免所有对象同时发同样橙光 |
| `W1-07` | `stage06.json`，普通Stage 7 | 芭蕉机关塔 | 36x43纵向复合关，集中使用炉石、风瓮、虫穴、钉柱、裂墙和阵枢 | 每层用青铜梁和芭蕉风纹建立垂直分区；机关密集但安全落脚点保持低纹理，不能让角色淹没在装饰中 |
| `W1-08` | `stage07.json`，普通Stage 8 | 玉泉暗渠 | packed water、避水珠、风阵、爆裂墙和巡游风火轮 | 场景由赤岩过渡到青绿暗渠；水源必须接到可解释的龙口或渠壁，避水珠与水面高光、灵蕴珠三者轮廓不同 |
| `W1-09` | `stage08.json`，普通Stage 9 | 风水合阵 | 水体、风瓮、阵枢、爆裂墙和敌人组门 | 用铜管、渠槽和风道把水重排与风机关放进同一工程系统；环境水汽不能遮挡风向和锁状态 |
| `W1-10` | `stage09.json`，Boss Stage 10 | 铁甲牛将关 | Evil Teutonic Knight源码状态机、两道组门和World 1经卷灵印 | 重做为铸坊守关台，铁甲牛将与黑铁/青铜环境同源但轮廓更亮；冲锋通道保持开阔，兵器尖端不改变命中范围 |
| `W1-S1` | `stage10.json`，从Stage 4分出的秘密关 | 余烬秘库 | 炉石、爆裂墙、风瓮、阵枢和符锁的短型组合 | 熄火后的灰黑库房只保留局部余烬；秘密感来自封条和低光，不靠把可走地面压成全黑 |
| `W1-S2` | `stage11.json`，从Stage 7分出的秘密关首关 | 百柱铸坊 | 60x20横向关、16组钉柱、71块裂墙及风水机关 | 强调长横向铸造线和重复机械节奏；每组危险柱保持完整接缝，避免旗幡、梁柱和钉柱互相截断 |
| `W1-S3` | `stage12.json`，秘密关链末关 | 地火竖炉 | 23x60纵向关、虫穴、风瓮、爆裂墙和多层组门 | 从上层冷却渠向下层炉心渐暖；竖向导航靠材质梯度和结构地标，不能用与玩法无关的巨大文字标牌 |

### 2.8 雪岭场景矩阵

`map_tibet.out`表明`stage00..10`是普通链，`stage11`从Stage 2分出，`stage12 -> stage13`从Stage 5分出。场景设计必须保留这两条秘密分支的视觉归属。

| scene ID | 数据/地图身份 | 新场景名 | 代码与数据约束 | 生成和布景要求 |
| --- | --- | --- | --- | --- |
| `W2-01` | `stage00.json`，普通Stage 1 | 冰钟乳长廊 | 34个raw `44`悬顶冰钟乳、4只raw `49`玄龟、裂墙和1处raw `38`冰泉源 | 顶部裂纹是主要预警，冰钟乳与普通冰墙使用不同边缘高光；冰泉从岩缝接入暗渠，玄龟石碑采用古寺文字纹样但不放可读现代字体 |
| `W2-02` | `stage01.json`，普通Stage 2 | 猿傀冰窟 | raw `45`雪岭猿傀、冰钟乳、风脉阵枢和阴阳风瓮 | 用断裂符绳、木面具架和青玉关节说明傀儡来源；悟乐与猿傀同屏时，体色、服装、脸部轮廓和姿势均不可混淆 |
| `W2-03` | `stage02.json`，普通Stage 3 | 寒针蜂巢 | 12只raw `46`寒针蜂、冰钟乳、猿傀、玄龟和多组门 | 蜂巢改为古寺檐下结霜妖巢，投射寒针与背景冰棱方向不同；飞行路径周围避免高密经幡遮挡 |
| `W2-04` | `stage03.json`，普通Stage 4 | 风雪悬桥 | raw `47`寒风葫芦首次成组出现，并混合冰钟乳、猿傀和寒针蜂 | 葫芦是圆形可动物，口朝上；生成的风柱贯穿上下格且能托举对象，悬桥背景不影响格边判断 |
| `W2-05` | `stage04.json`，普通Stage 5 | 日月寒风阵 | 阴阳风瓮、风脉阵枢、寒风葫芦、三类敌人和符锁 | 雪岭风瓮使用朱砂陶/青玉陶，不照搬火焰山铜瓮；日月纹保持跨世界规则一致，风雪粒子显示方向但不遮住机关相位 |
| `W2-06` | `stage05.json`，普通Stage 6 | 玄冰淬棒殿 | raw `26`位于`(32,22)`经匣；同时出现raw `48`寒光镜、葫芦、冰钟乳和玄龟 | 奖励台是古寺淬炼阵，表现已有如意铁棒被寒光赋能；经匣关闭时不可提前出现冰棒轮廓，镜柱与淬炼阵共享青玉光学材质 |
| `W2-07` | `stage06.json`，普通Stage 7 | 照妖镜冰廊 | 7组raw `48`两格镜柱、寒风葫芦、冰钟乳和多种敌人 | 镜头、底座和横向寒光必须是一套完整装置；光束端点清楚，中段透明度受控，不能把安全通道染成不可读的纯蓝块 |
| `W2-08` | `stage07.json`，普通Stage 8 | 经幡蜂道 | 敌人组门密集，含寒针蜂、镜柱、葫芦、冰钟乳和玄龟 | 经幡只放在不遮挡蜂针和门状态的后景；每个竞技场入口用相同镇妖符门语法，避免误判为场景边框 |
| `W2-09` | `stage08.json`，普通Stage 9 | 寒光迷宫 | 63块裂墙、6组镜柱、8个寒风葫芦和多种敌人；player raw `74`位于`(6,17)`，在JAR初始化默认分支中立即写回`-1` | 裂墙与实墙明度差固定，镜光用于建立路线层级而非铺满画面；长迷宫区使用重复冰玉刻度作为位置地标；raw `74`不生成贴图、不占位，也不赋予新机关语义 |
| `W2-10` | `stage09.json`，普通Stage 10 | 百丈冰桥 | 104x17超长横向关，镜柱、葫芦、冰钟乳及三类敌人连续组合 | 远景雪峰随区段变化提供进度感，前景桥面保持平静；所有悬顶危险仍从顶边接出，不能像漂浮奖励 |
| `W2-11` | `stage10.json`，Boss Stage 11 | 寒魄狮王殿 | Yeti源码Boss状态机、组门和World 2经卷灵印 | 古寺主殿形成开阔战斗轴，寒魄狮王维持寒阵；扑跃、震地和抛射各有独立轮廓，奖励雪莲灵印在胜利前保持封闭 |
| `W2-S1` | `stage11.json`，从Stage 2分出的秘密关 | 冻寺密库 | 大量裂墙、冰钟乳、猿傀、风瓮和符锁 | 以封存经架、冰封木梁和低亮金属构成密库；可破墙不能与书架装饰混淆，钥匙区使用日/月符号定位 |
| `W2-S2` | `stage12.json`，从Stage 5分出的秘密关首关 | 冰泉镇妖廊 | 冰钟乳、猿傀、寒针蜂、寒风葫芦、raw `38`冰泉源，并出现跨世界蛇类和风阵 | 解释为古寺收押各地妖物的镇妖廊，蛇类沿用花果山物种识别但加环境积霜；冰泉连接石槽，不把蛇误画成新的雪岭敌人 |
| `W2-S3` | `stage13.json`，秘密关链末关 | 雪印终室 | 冰钟乳、寒风葫芦、raw `38`冰泉源、风阵和终端奖励路径 | 场景从灰蓝冻寺过渡到银白印室，冰泉作为贯穿前两关的冷光引导；终点强调封印完成而非额外Boss，高亮只集中在奖励区 |

### 2.9 资产与场景绑定

`stage-scenes.json`之外还要自动生成`assets/pilgrim/manifests/asset-scene-bindings.json`。每个`Rxxx/Dxxx`条目至少包含：

```text
asset_id
direct_occurrences[]: scene_id, layer, raw_id, coordinates
runtime_occurrences[]: scene_id, parent_raw_id, state_or_event
progress_occurrences[]: save_field, threshold_or_fallback, reachable_scenes
ui_occurrences[]: screen, state
approved_variants[]
variant_selector: world / scene_id / ui_state
runtime_disposition: rendered / generated / removed_on_init / unreachable / unreachable_compat
```

绑定规则：

- 地形、基础石头和可挖层通过world manifest自动绑定该世界所有实际引用它们的scene，不把“属于某世界”等同于“每关都出现”。
- 通用门、经匣、钥匙、压阵盘、收集物和敌人通过layer/raw ID及源码绘制分支生成坐标；同一数字在不同layer含义不同，禁止只按raw数值合并。
- 爆炸、冰冻、风柱、开箱光、Boss攻击等运行时生成资产绑定到父对象和状态机事件。例如`R050`风柱继承raw `47`所在的World 2 scenes，`R061`的共用断冰环同时记录raw `44`和raw `47`事件来源。
- 工具图标和动作按进度可达性绑定。花果山Stage 5的幌金绳来自跨世界持久化状态，不得因该关没有raw `27`坐标就漏掉素材；进入火焰山/雪岭时源码还会分别把缺失工具补到等级`1/2`，manifest需记录`world_entry_fallback`，但不能为兜底状态伪造经匣坐标。
- `R013/R014`只绑定`W0-09`，`R070..R072`只绑定`W0-06`，`R015`只绑定`W1-10`，`R073`只绑定`W2-11`，`R074`只绑定`W0-TUTORIAL`。
- `R031/R048`的occurrence以完整`2x2`组合记录，不以四个孤立坐标生成四份无关brief。
- 关卡数据中被JAR初始化明确删除的值记录为`runtime_disposition: removed_on_init`。例如`W2-09`的player raw `74`必须留在数据审计中，但不得进入素材生成队列或被渲染成“缺失对象”。
- 品牌、HUD、菜单、地图、Tips和字体使用`ui_occurrences`，不能伪造关卡scene ID；所有其他资产若三种occurrence都为空则阻塞发布。

### 2.10 资源闭环状态和单一事实源

本文中的`Rxxx/Dxxx`表是制作输入，不直接充当运行时资源索引。进入资产生产前必须建立以下单一事实源；除`provenance.csv`外均位于`assets/pilgrim/manifests/`：

| 文件 | 职责 | 生成/维护方式 |
| --- | --- | --- |
| `assets.json` | 登记`R001..R121`和`D001..D015`的处置、输出、依赖、状态和哈希 | 人工批准语义，构建工具补文件哈希和状态 |
| `stage-scenes.json` | 41个真实关卡scene及材质、光源、机制、视口和环境效果 | 从三世界stage数据和地图连接生成，再人工审美审核 |
| `asset-scene-bindings.json` | 资产到直接坐标、运行时事件、进度状态和UI画面的反向绑定 | 从关卡数据、源码语义表和`assets.json`生成 |
| `runtime-resource-map.json` | 运行时语义键到atlas、音频、文本、字体和world/scene variant的唯一解析表 | 由已验证的`assets.json`生成 |
| `animation-events.json` | 动画名、20Hz总时长、逻辑事件tick和可插值区间 | 从源码事件审计和新atlas合同生成 |
| `audio-events.json` | 逻辑音频事件/兼容ID `0..20`到`R087..R107`、优先级、循环及duck规则 | 从第3.7节生成并由音频集成测试验证 |
| `text-key-map.json` | 旧文本索引、脚本索引和Go可见硬编码到新语义key的迁移关系 | 文本扫描与脚本解析生成，人工审校语义 |
| `font-glyphs.txt` | 两个首发locale、动态数字、标点、键帽和专名所需的去重字形集 | 每次文本构建重新生成 |
| `terrain-frame-map.csv` | 三世界151个地形frame的坐标、邻接、接缝和新frame映射 | 从关卡三层数据生成 |
| `procedural-fx.json` | 风带、水汽、风雪、余烬、瀑雾和碎屑等纯视觉环境效果及scene绑定 | `D013`，人工定参数并做遮挡验收 |
| `deletion-evidence.json` | 不生成资源的槽位、源码/JAR可达性证据和防回退测试 | 人工登记，构建测试验证 |
| `retained-content-exceptions.json` | 第一阶段有意保留的关卡层、地图连接和demo非可见命令的路径、哈希、用途及风险归属 | 人工批准范围，构建工具验证哈希和内容类型 |
| `original-source-lock.json` | 原版精灵、音频、字体和三世界数据逐文件SHA-256及整树哈希 | `tools/pilgrimmanifest`在生成前后核对；非显式刷新时变化即失败 |
| `closure-report.json` | 每类数量、缺失输出、未完成variant、文本缺键、旧资源引用和发布包扫描结果 | 只由构建工具生成，禁止手改 |
| `schemas/*.schema.json` | 上述JSON的字段、枚举和唯一性约束 | 与资源加载器同版本维护 |

`assets.json`每项至少包含：

```text
schema_version, asset_id, class, disposition, semantic_keys[], legacy_sources[],
character_id, functional_roles[], output_mode,
outputs[]: path, owner_asset_id, variant, locale, sha256,
locales[], scene_ids[], variants[], runtime_keys[], prompt_keys[], distribution_targets[],
primary_batch, dependencies[], provenance_refs[], status, blockers[]
```

- `class`只能是`R`或`D`；`disposition`只能是`replace`、`rewrite`、`delete`、`procedural`或`system`。
- 一个图集可以拥有多个`semantic_keys/runtime_keys`，但每个具体key只能由一个资产ID拥有；`legacy_sources`只参与迁移覆盖，不允许成为发布路径或资源查找键。
- 每个最终文件只能有一个`owner_asset_id`；其他条目通过`dependencies`引用，禁止两个ID共同写同一JSON、atlas或音频文件。目录和共享母版不视为可发布文件，仍需逐文件登记owner。
- `distribution_targets`只能从`runtime`、`source_archive`、`release_metadata`、`system`和`none_deleted`中选择。`packaged`表示已进入声明的目标包，不要求制作母版混入玩家运行包，也不能把运行时atlas只放进内部归档。
- 状态按`specified -> unblocked -> materialized -> integrated -> verified -> packaged`推进。`materialized`对图片/音频表示文件已经产生，对文本/配置表示数据已经写入，对删除项表示`deletion-evidence.json`已经具备证据。
- 依赖未满足时状态为`blocked`并列出`blockers`。只有所有条目达到`packaged`且`closure-report.json`无错误时，才可称资源生产闭环；文档中已有替换说明只代表`specified`。
- 每个输出进入`materialized`后记录SHA-256；同一ID的每个variant和locale分别记录哈希，不能用目录存在代替文件完整性。
- `provenance.csv`记录人工、生成和第三方内容来源；它自身以及构建生成的manifest、schema和`closure-report.json`不建立递归来源行，而由报告记录生成器版本、Git提交、输入哈希和构建时间。报告自身的发行哈希由外层包清单或签名记录。

运行时只通过`runtime-resource-map.json`解析语义资源。玩法层可以继续保留raw ID和兼容音频ID，但绘制、音频和文本加载器不得直接拼接`decoded/...`路径。variant选择器必须穷举所有可达`world/scene_id/ui_state`；缺键、未知variant、文件缺失或哈希不符均立即失败，不得回退到其他世界、默认旧图或原英文文本。

删除类至少登记三类证据：`R118`为源码/JAR不可达；`R080` frame `15`为不可达兼容索引并固定输出透明占位；`W2-09` player raw `74`为初始化删除的数据值。每项都要有静态扫描和运行时防回退测试，不能仅凭“当前没看到”判定删除。

第一阶段保留内容不能混入`assets.json`冒充新资源，必须规范化迁移到`content/pilgrim/`并逐文件进入`retained-content-exceptions.json`。每项至少记录`path`、`sha256`、`content_type`、`retained_reason`、`allowed_runtime_usage`、`forbidden_visible_payload`和`risk_owner`；`content_type`只允许`level_json`、`map_graph`和`demo_commands`。原`w*.bin`/JAR容器、PNG、音频、字体、Logo、头像、可见对白和旧资源metadata一律不能使用该例外。未列入、仍从`decoded/`加载或哈希变化的内容使发布构建失败。

仓库根下的`assets/resource-profiles.json`定义并行profile。`original`保持默认、只读现有`decoded/`并继续由`cmd/originalrush`启动；`wule`只读取`assets/pilgrim/`和`content/pilgrim/`，使用独立存档命名空间。两者均为`missing_policy=error`、`fallback_profile=null`，禁止任一方向的静默回退。Wule最终atlas与独立入口尚未接入前必须保持`ready=false`，不能用原版资源填充缺口后伪装成可运行。

当前已知的生产阻塞项如下；它们不影响制作规范自洽，但在解决前不能进入`packaged`：

| 资产/范围 | 阻塞条件 |
| --- | --- |
| `R108/D015`及品牌面 | 确定最终中英文发行名、bundle标识和存档命名空间，并完成商标与相似性审查 |
| `R110` | 确定权利主体、年份和两种locale的法律文本模板 |
| `D014`及悟乐相关角色面 | 批准唯一`character_id=wule`角色母版、Boss轮廓和三世界风格圣经后，才可生产`R085/R109/D009`及所有头像/插图 |
| `D004/D005` | 固定字体准确版本、文件哈希、许可证全文和游戏嵌入方式 |
| `R024/R116/D006/D007` | `zh-CN`与`en-US`键集合、占位符集合及容器宽度全部通过验证 |
| `R082..R084` | 完成无文字标题装饰框并由运行时绘制本地化世界名 |
| `D008/D011` | 只能在最终atlas、Logo和文本锁定后重新渲染，不得引用中间素材 |
| `D012` | 每个最终输出和第三方工具依赖都有来源、许可、编辑与审核记录 |
| `D013` | 每个程序化效果都有scene绑定、强度上限、前后景层级和禁用条件 |

`assets.json`至少声明以下依赖边；箭头右侧是前置项：

| 依赖项 | 必须依赖 | 原因 |
| --- | --- | --- |
| 除`D014`母版本身外，所有基于生成式流程制作的角色、敌人、Boss、场景和品牌视觉 | `D014`对应母版/风格条目 | 防止角色外形、材质和功能轮廓跨批次漂移 |
| `R108/R109/D006/D009/D011` | `D015` | Logo、标题、文本、应用图标和宣传面共享唯一产品身份 |
| `R110` | `D006`及已批准权利主体 | 法律事实与本地化模板分离，避免双写locale文件 |
| `D004/D005` | `D006/D007/R024/R116`的两种locale最终文本 | 字体atlas必须从最终字形并集生成 |
| `R082..R084` | `D004/D005/D006` | 装饰框先保留安全区，世界名由最终字体和文本绘制 |
| `D008` | 对应scene的全部`integrated`运行时atlas、字体、文本和`D013` | 预览必须由最终运行画面重渲染 |
| `D011` | `R108/R109/D006/D008/D015` | README、商店文案和截图不能引用中间品牌或旧画面 |

普通资产依赖图必须是DAG，出现循环即构建失败。`D012`、`deletion-evidence.json`、`retained-content-exceptions.json`和`closure-report.json`是末端聚合/审计门禁，不作为所有资产的反向前置依赖，以免形成自引用。

## 3. 121个原资源块逐项结论

表中 `M/F/A` 为原块module/frame/animation数量，仅用于覆盖核对。输出的新atlas可以重新拆帧。

### 3.1 三个世界的基础地形

| ID | legacy slot | 原用途和结构 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R001` | `0.f#0` `M8/F8/A0` | World 0普通石头8个旋转frame | 八面带浅金镇字符的“花果山镇山石”，8个旋转状态 | `worlds/huaguoshan/objects/mountain_stone/` | `T1O` / `P-W0-OBJECT` |
| `R002` | `0.f#1` `M13/F8/A1` | World 0可挖植被及消散 | 翠藤、蕨叶和苔藓组成的“水帘藤障”，含完整到消散8阶段 | `worlds/huaguoshan/terrain/diggable_vines/` | `T1` / `P-W0-TERRAIN` |
| `R003` | `0.f#2` `M30/F40/A0` | World 0 `raw 80..119`的40个地形frame | 40个同拓扑的青灰洞壁、石梁、根系、神龛和前景边角，逐frame映射到`frame_000..039` | `worlds/huaguoshan/terrain/tiles/` | `T1` / `P-W0-TERRAIN` |
| `R004` | `0.f#3` `M1/F0/A0` | World 0默认地面填充module | 无缝青灰洞底石板，带稀疏苔点，不含可交互符号 | `worlds/huaguoshan/terrain/base_floor.png` | `T1` / `P-W0-TERRAIN` |
| `R005` | `1.f#0` `M5/F8/A0` | World 1普通石头8个旋转frame | 八面赤铜包边的“镇火石”，与雷火炉石保持明显区别 | `worlds/flame_mountain/objects/fire_stone/` | `T1O` / `P-W1-OBJECT` |
| `R006` | `1.f#1` `M20/F6/A1` | World 1可挖障碍 | 焦藤、灰烬和脆裂盐壳组成的“焦藤障”，6阶段消散 | `worlds/flame_mountain/terrain/diggable_ash_vines/` | `T1` / `P-W1-TERRAIN` |
| `R007` | `1.f#2` `M39/F66/A0` | World 1 `raw 80..145`的66个地形frame | 66个同拓扑的赤岩、青铜梁、黑铁墙、玉泉渠、布幔和前景构件，逐frame映射到`frame_000..065` | `worlds/flame_mountain/terrain/tiles/` | `T1` / `P-W1-TERRAIN` |
| `R008` | `1.f#3` `M1/F0/A0` | World 1默认地面module | 无缝暗赤岩地板，嵌细青铜压条 | `worlds/flame_mountain/terrain/base_floor.png` | `T1` / `P-W1-TERRAIN` |
| `R009` | `2.f#0` `M8/F8/A0` | World 2普通石头8个旋转frame | 带青玉冰脉的“玄冰镇山石”，8个旋转状态 | `worlds/snow_road/objects/ice_stone/` | `T1O` / `P-W2-OBJECT` |
| `R010` | `2.f#1` `M9/F7/A1` | World 2可挖障碍 | 松雪、冰藤和经幡碎屑组成的“积雪障”，7阶段崩散 | `worlds/snow_road/terrain/diggable_snow/` | `T1` / `P-W2-TERRAIN` |
| `R011` | `2.f#2` `M33/F45/A0` | World 2 `raw 80..124`的45个地形frame | 45个同拓扑的冰岩、古寺砖、冻木梁、经幡和前景冰棱，逐frame映射到`frame_000..044` | `worlds/snow_road/terrain/tiles/` | `T1` / `P-W2-TERRAIN` |
| `R012` | `2.f#3` `M1/F0/A0` | World 2默认地面module | 无缝深蓝冻土地板，带少量薄霜，不画危险冰面符号 | `worlds/snow_road/terrain/base_floor.png` | `T1` / `P-W2-TERRAIN` |

`R003/R007/R011` 使用确定性槽位映射：旧frame `n` 只对应新 `frame_%03d`，不得合并、漏帧或更换raw ID。新图依据关卡层、邻接签名和四边接缝合同重新设计，不复制旧frame的像素级alpha、裂口轮廓或装饰构图。碰撞仍来自Runtime，PNG透明区不参与判定。

### 3.2 Boss、通用对象和脚本容器

| ID | legacy slot | 原用途和结构 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R013` | `b0.f#0` `M18/F23/A9` | Great Anaconda身体，9组状态动画 | 原创“碧鳞蟒王”：青碧鳞、金色腹甲、断角和藤纹，覆盖入场、探身、易伤、受击、回收、换列和败退 | `bosses/python_king/body/` | `BOSS` / `P-BOSS-PYTHON` |
| `R014` | `b0.f#1` `M2/F2/A0` | 蟒王战的两种平台frame | 蟒王祭坛的完整石台和破裂石台 | `bosses/python_king/platform/` | `T2H` / `P-BOSS-PYTHON` |
| `R015` | `b1.f#0` `M72/F78/A14` | Bavaria铁甲Boss，14组动画 | 原创“铁甲牛将”：黑铁层甲、青铜牛首面具、双手降魔槊；14个状态槽全部重做 | `bosses/iron_bull/` | `BOSS` / `P-BOSS-BULL` |
| `R016` | `cm.f#0` `M2/F0/A0` | 普通/特殊出口图块 | 普通出口为白金“祥云界门”，秘密出口为青玉“隐云界门”，轮廓和符纹均不同 | `common/goals/cloud_gates/` | `T1` / `P-COMMON` |
| `R017` | `cm.f#1` `M6/F6/A0` | 三世界通用门，原块3个palette | 两格连续的“镇妖符门”；花果山石门、火焰山铜门、雪岭冰门各一套材质，运行时按世界选atlas | `common/doors/seal_gate/huaguoshan/`；`common/doors/seal_gate/flame_mountain/`；`common/doors/seal_gate/snow_road/` | `T2V` / `P-COMMON`及对应世界地形brief |
| `R018` | `cm.f#2` `M4/F4/A0` | 紫钻和红钻两套palette | 紫钻替换为紫青圆形“灵蕴珠”；红钻替换为朱砂方形“云路印”，作为累计云路门槛而非药品；各自含静止高光帧 | `common/pickups/essence_and_route_seals/` | `T1O` / `P-COMMON` |
| `R019` | `cm.f#3` `M9/F4/A1` | 通用格子破坏/状态覆盖效果 | 金色符尘、灰色石屑和少量灵光组成的4帧“破障消散” | `common/fx/tile_dispel/` | `FX` / `P-FX` |
| `R020` | `cm.f#4` `M2/F0/A0` | module 0额外生命，module 1补血；结算奖励也复用 | module 0为一根带朱绳结的“救命毫毛”，module 1为适合高频投放的橙红“回元果”；获得完美奖励时仍显示救命毫毛 | `common/pickups/life_and_health/` | `T1O` / `P-COMMON` |
| `R021` | `cm.f#5` `M2/F2/A0` | 配额门主体和数字底座 | 圆形“功德封印”和悬挂数字木牌，剩余数量必须清晰 | `common/gates/merit_quota/` | `T1` / `P-COMMON` |
| `R022` | `cm.f#6` `M3/F8/A0` | 检查点未激活/激活动画 | 固定在关卡坐标的“土地神龛”，冷灰石龛到暖金香火共8帧 | `common/checkpoints/earth_shrine/` | `T1O` / `P-COMMON` |
| `R023` | `cm.f#7` `M17/F71/A6` | 6组拾取、冻结、奖励和封印效果 | 分拆为6组新特效：经匣开光、玄冰凝结、云路印朱金灵光、灵蕴珠收束、法宝星芒、五行封印激活 | `common/fx/pickup_effects/` | `FX` / `P-FX` |
| `R024` | `cr.f#0` 非精灵 | 原Credits字节流 | 重写为“制作团队、素材来源、字体与音频许可、开源组件、生成工具记录”五段滚动文本；中英文键结构一致 | `text/credits.zh-CN.json`；`text/credits.en-US.json` | 文本重写 / `P-TEXT` |
| `R025` | `demo.f#0` 非精灵 | 教程及关卡演出命令流 | 转写为`content/pilgrim/scripts/demo.json`；保留移动、等待、镜头和写层命令，全部对白ID改指向新西游文本 | `content/pilgrim/scripts/demo.json` | 数据重写，不生成图片 |
| `R026` | `demoui.f#0` `M10/F0/A0` | 对话框边、按键图，两套palette | 云纹卷轴对话框、土地神头像底座、`SPACE`/`ENTER`/`S`桌面键帽；中性青铜和Boss朱红两套主题 | `tutorial/dialog_ui/` | `HUD` / `P-TUTORIAL` |

### 3.3 `gen0.f` 到 `gen2.f` 通用机制

| ID | legacy slot | 原用途和结构 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R027` | `gen0.f#0` `M5/F1/A0` | Freeze Hammer条件加载的凝结/生成静态frame | 五瓣青白“玄冰凝结环”，用于冻结瞬间，不承担碰撞 | `common/fx/freeze_flash/` | `FX` / `P-FX-ICE` |
| `R028` | `gen0.f#1` `M7/F1/A0` | 冻结紫钻raw 9类型34 | 冰壳包裹的“冰封灵蕴珠”，仍保持圆形可滚轮廓 | `common/frozen/frozen_essence/` | `T1O` / `P-FX-ICE` |
| `R029` | `gen0.f#2` `M1/F0/A0` 两palette | 金钥匙/银钥匙拾取图 | 金色日纹“金符钥”和银色月纹“银符钥”，形状也必须不同 | `common/keys/talisman_keys/` | `T1O` / `P-COMMON` |
| `R030` | `gen0.f#3` `M15/F12/A2` | 爆炸爆发与消散 | “丹炉爆焰”：中心白热、朱红火环、黑灰石屑，两组动画分别为爆发和余烬 | `mechanisms/explosion/alchemy_blast/` | `FX` / `P-W1-MECH` |
| `R031` | `gen0.f#4` `M8/F8/A4` | 花果山foreground raw `20..23`；关卡数据主要按`[20,21;22,23]`组成`2x2`块，少量块按`[21,20;23,22]`水平镜像 | 先生成完整`192x192`“青藤镇洞灯龛”的暗/亮两状态，再切四象限；animation `0/1/2/3`依次绑定左上/右上/左下/右下，每象限2帧同步呼吸，不增加碰撞 | `worlds/huaguoshan/foreground/two_by_two_shrine/` | `T2X2` / `P-W0-TERRAIN` |
| `R032` | `gen0.f#5` `M2/F2/A0` | 可滚落、可爆的raw 8 | 带赤金裂纹和铜箍的“雷火炉石”，静止/蓄热2帧 | `worlds/flame_mountain/objects/thunderfire_stone/` | `T1O` / `P-W1-OBJECT` |
| `R033` | `gen0.f#6` `M1/F1/A0` | player raw 38初始化后的foreground raw 27水源；World 1与World 2均有authored实例 | 火焰山variant为青铜龙口“玉泉泉眼”，雪岭variant为石缝莲口“冰泉泉眼”；都明确朝向水体且不画奖励高光 | `common/water/spring_source/flame_mountain/`；`common/water/spring_source/snow_road/` | `T1` / `P-WATER` |
| `R034` | `gen0.f#7` `M10/F8/A1` | raw 30可破坏墙，原块2个palette但三世界均有authored实例 | “裂纹符墙”：完整、裂纹扩散、符纸脱落到清除8阶段；花果山岩墙、火焰山铜砖、雪岭冰寺砖各一套，裂纹拓扑和碰撞轮廓相同 | `common/walls/breakable_talisman_wall/huaguoshan/`；`common/walls/breakable_talisman_wall/flame_mountain/`；`common/walls/breakable_talisman_wall/snow_road/` | `T1` / `P-COMMON`及对应世界地形brief |
| `R035` | `gen0.f#8` `M8/F3/A2` | raw 36 crawler trap | “地火虫穴”：闭合石盖、张开火口、熄灭焦痕；两组动画覆盖休眠和触发 | `worlds/flame_mountain/hazards/fire_crawler_nest/` | `T1O` / `P-W1-MECH` |
| `R036` | `gen0.f#9` `M1/F2/A0` | raw 22/23火焰发射器左右frame；现有关卡数据仅World 0使用 | 左右朝向的“地火石兽口”，作为水帘洞地火脉泄压口；开口方向必须在轮廓上可辨 | `worlds/huaguoshan/hazards/earthfire_emitter/` | `T1` / `P-W0-COLLAPSE` |
| `R037` | `gen1.f#0` `M41/F27/A3` | 横向火焰、Boss通用抛射物、蟒王尾击 | animation 0“三昧火舌”，1“妖王法力冲击”，2“碧鳞蟒王尾击”；新包可拆成3个atlas | `common/hazards/flame_and_boss_fx/` | `FX/BOSS` / `P-FX` |
| `R038` | `gen1.f#1` `M2/F4/A0` | raw 28上下伸缩钉柱 | 连续的“镇妖钉柱”：柱身、紫色缺块全部替换成青铜伸缩套环，共4阶段 | `worlds/flame_mountain/hazards/spike_column/` | `T2V` / `P-W1-MECH` |
| `R039` | `gen1.f#2` `M13/F0/A0` | raw 14横向移动危险物和拖尾module | 带火纹方向叶片的“巡游风火轮”，含左右本体、火星和速度拖尾 | `worlds/flame_mountain/hazards/patrol_fire_wheel/` | `T1O/FX` / `P-W1-MECH` |
| `R040` | `gen1.f#3` `M20/F18/A2` | raw 16铁人/长矛危险物；authored底格初始化时复制上格，逻辑占纵向两格，左右方向各一组动画 | 至少`216x212`透明母版的“铜甲镇关力士”：青铜甲、黑铁底座和降魔长矛，左右各9阶段；脚底对齐底格，头和兵器进入上格，不能简化成独立伸缩刺 | `worlds/flame_mountain/hazards/bronze_spear_guard/` | `T2V`视觉越界 / `P-W1-MECH` |
| `R041` | `gen1.f#4` `M6/F0/A0` | raw 11贴墙crawler各阶段module | “石蝎精”：3个爬行姿态和3个石化碎裂阶段，四向由运行时翻转 | `common/enemies/stone_scorpion/` | `T1O` / `P-ENEMY` |
| `R042` | `gen1.f#5` `M47/F33/A5` 多palette | World 0 raw 19/43蛇，同时由World 2秘密`stage12`的raw 19复用 | 花果山variant为青藤蛇妖/赤练蛇妖两套独立花纹；雪岭variant保持同物种轮廓并增加薄霜、镇妖颈环，覆盖巡逻、转向、眩晕、受击和消散 | `common/enemies/vine_snakes/huaguoshan/`；`common/enemies/vine_snakes/snow_road/` | `T1O` / `P-W0-ENEMY`、`P-W2-ENEMY` |
| `R043` | `gen1.f#6` `M3/F1/A0` | 冻结raw 19/43蛇，raw 9类型37 | 半透明冰壳内的“冰封蛇俑”，红蛇身份用腹部纹样保留；世界光照只改外层反光，不改变冻前物种识别 | `common/enemies/frozen_snake/` | `T1O` / `P-FX-ICE` |
| `R044` | `gen1.f#7` `M13/F14/A3` 多palette | World 1复用raw 19/43的另一套敌人图 | 改为青铜背脊“火蜥妖”和银蓝“泉蜥妖”，覆盖爬行、转向和受击 | `worlds/flame_mountain/enemies/fire_lizards/` | `T1O` / `P-W1-ENEMY` |
| `R045` | `gen1.f#8` `M2/F1/A0` | 冻结World 1 raw 19/43 | “冰封火蜥俑”，保留日/月腹纹区分两种身份 | `worlds/flame_mountain/enemies/frozen_lizard/` | `T1O` / `P-FX-ICE` |
| `R046` | `gen1.f#9` `M3/F0/A0` | 三个工具等级图标；等级`1`来自raw `24`或火焰山入口兜底，等级`2`来自raw `27`或雪岭入口兜底，等级`8`只来自raw `26`淬炼 | module 0“如意铁棒”、1“幌金绳”、2“玄冰淬炼后的如意铁棒”；第三项必须保留同一棒身识别点并增加霜纹，表现强化而非另一件独立武器 | `common/tools/tool_icons/` | `T1O` / `P-COMMON` |
| `R047` | `gen2.f#0` `M2/F1/A0` | foreground raw 2锤/钩提示 | 地面“铁棒破障纹”和“金绳牵引纹”，不显示手机数字键 | `common/prompts/tool_floor_glyphs/` | `T1` / `P-COMMON` |
| `R048` | `gen2.f#1` `M8/F8/A4` | 火焰山foreground raw `20..23`；关卡数据主要按`[22,23;20,21]`组成`2x2`块，少量块水平镜像 | 先生成完整`192x192`“青铜风火灯龛”的暗/亮两状态，再按源码排布切四象限；animation `2/3`为左上/右上，`0/1`为左下/右下，每象限2帧同步摆动，禁止拆成四种无关装饰 | `worlds/flame_mountain/foreground/two_by_two_shrine/` | `T2X2` / `P-W1-TERRAIN` |
| `R049` | `gen2.f#2` `M4/F4/A1` | foreground raw 14普通/组锁宝箱 | 木铜结构“经匣”：闭合、松锁、开盖、空箱4状态；关闭时完全遮住payload | `common/chests/scripture_chest/` | `T1O` / `P-COMMON` |
| `R050` | `gen2.f#3` `M11/F8/A2` | raw 47风力pod；源码走重力、推动、钩取，并向上生成foreground `35/37/34`风柱；现有关卡仅World 2使用 | 圆腹青玉“寒风葫芦”：静止时葫芦口朝上，受支撑后喷出可托举物体的风雪柱；两组动画覆盖葫芦本体和持续上升气流 | `worlds/snow_road/mechanisms/frost_wind_gourd/` | `T1O` / `P-W2-MECH` |
| `R051` | `gen2.f#4` `M5/F5/A0` 两palette | raw 34/35红蓝fan pot；World 1与World 2均有authored实例 | “阴阳风瓮”保持日/月两种独立轮廓和5阶段开合；火焰山版为赤铜/青铜，雪岭版为朱砂陶/青玉陶，运行时按世界选atlas | `common/mechanisms/yin_yang_wind_pots/flame_mountain/`；`common/mechanisms/yin_yang_wind_pots/snow_road/` | `T1O` / `P-W1-MECH`、`P-W2-MECH` |
| `R052` | `gen2.f#5` `M4/F0/A0` | raw 37爆破墙4阶段module | 黑铁包边“爆裂符墙”，4个破损阶段 | `worlds/flame_mountain/walls/blast_seal_wall/` | `T1` / `P-W1-MECH` |
| `R053` | `gen2.f#6` `M17/F34/A0` | World 1/2共用的三层packed water全部形状 | 同一34-frame拓扑输出两套材质：火焰山“玉泉水”和雪岭“冰泉水”，覆盖水面、内部、左右流、源头、边角和特殊水格；运行时按世界选atlas | `common/water/tiles/flame_mountain/`；`common/water/tiles/snow_road/` | `T1` / `P-WATER` |
| `R054` | `gen2.f#7` `M1/F0/A0` | raw 40避水道具 | 青玉莲纹“避水珠”，外轮廓与灵蕴珠不同 | `common/pickups/water_breathing_orb/` | `T1O` / `P-COMMON` |
| `R055` | `gen2.f#8` `M8/F7/A0` 两palette | foreground raw 8/9金银锁 | 日纹金符锁、月纹银符锁，各含闭合、解锁和消散frame | `common/locks/talisman_locks/` | `T1` / `P-COMMON` |
| `R056` | `gen2.f#9` `M1/F0/A0` | foreground raw 6压力板；三个世界均有authored实例 | 黑白八卦“压阵盘”，花果山石环、火焰山青铜环、雪岭冰玉环各一套；中心受力区相同，静态module仍由运行时下压13px | `common/switches/bagua_pressure_plate/huaguoshan/`；`common/switches/bagua_pressure_plate/flame_mountain/`；`common/switches/bagua_pressure_plate/snow_road/` | `T1` / `P-COMMON`及对应世界地形brief |

### 3.4 `gen3.f`、`gen4.f`和第三世界特殊对象

| ID | legacy slot | 源码确认用途 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R057` | `gen3.f#0` `M7/F6/A3` | 教程跟随提示及foreground raw 31提示动画 | animation 0随人物的`ENTER`召回键帽，1固定`SPACE`交互键帽，2固定`S`跳过键帽；统一云纹圆牌 | `tutorial/key_hints/` | `T1O` / `P-TUTORIAL` |
| `R058` | `gen3.f#1` `M1/F0/A0` | raw 42 Compass举起图 | 24逻辑像素的铜制莲纹“寻经针” | `common/pickups/scripture_compass/` | `T1O` / `P-COMMON` |
| `R059` | `gen3.f#2` `M8/F3/A0` | raw 48两格机关的底部发射头，3种朝向frame | 左射、右射、双向三种“寒光照妖镜柱”发射头 | `worlds/snow_road/mechanisms/frost_mirror_emitter/` | `T1O` / `P-W2-MECH` |
| `R060` | `gen3.f#3` `M4/F5/A1` | foreground raw 33的特殊容器覆盖 | 玉石和朱漆结构“镇妖经匣”，5阶段封条解锁与开盖；关闭时遮住payload | `common/chests/sealed_relic_chest/` | `T1O` / `P-COMMON` |
| `R061` | `gen3.f#4` `M11/F7/A5` | raw 44悬顶物、坠落、撞碎，以及raw 47生成的foreground `34/37`共用覆盖 | “悬顶冰钟乳”覆盖静止、预警裂纹、下坠和撞碎；共用的animation `2`设计为断冰环，既能作为冰钟乳撞碎残留，也能接在寒风葫芦风柱顶部，不能做成只在一种上下文成立的完整物件 | `worlds/snow_road/hazards/icicle_and_wind_cap/` | `T1O/FX` / `P-W2-MECH` |
| `R062` | `gen3.f#5` `M49/F38/A11` | raw 45 Monkey全部状态 | 敌对“雪岭猿傀”：白灰长毛被符绳束缚、整脸木制镇妖面具、青玉关节；明确是古寺傀儡而非悟乐的同类换色，覆盖11组巡逻、攀跳、攻击、眩晕和消散 | `worlds/snow_road/enemies/snow_ape_puppet/` | `T1O` / `P-W2-ENEMY` |
| `R063` | `gen3.f#6` `M5/F1/A0` | raw 45冻结类型35 | “冰封猿傀”，整脸木面具、符绳关节及其与悟乐的差异仍可辨 | `worlds/snow_road/enemies/frozen_snow_ape_puppet/` | `T1O` / `P-FX-ICE` |
| `R064` | `gen3.f#7` `M35/F39/A18` | raw 46 Wasp及其raw 21投射物 | “寒针蜂妖”和青白“寒针”：0..11飞行/追击，12..14寒针破碎，15..17近身攻击 | `worlds/snow_road/enemies/frost_wasp/` | `T1O/FX` / `P-W2-ENEMY` |
| `R065` | `gen3.f#8` `M5/F1/A0` | raw 46冻结类型36 | 翅膀收拢在冰壳内的“冰封蜂妖” | `worlds/snow_road/enemies/frozen_frost_wasp/` | `T1O` / `P-FX-ICE` |
| `R066` | `gen3.f#9` `M9/F4/A1` | raw 18 fan switch；World 1与World 2均有authored实例 | 四阶段“风脉阵枢”：停、左风、右风、过载；火焰山版用芭蕉叶与青铜轮，雪岭版用经轮与冰玉风标，运行时按世界选atlas | `common/mechanisms/wind_switch/flame_mountain/`；`common/mechanisms/wind_switch/snow_road/` | `T1O` / `P-W1-MECH`、`P-W2-MECH` |
| `R067` | `gen4.f#0` `M12/F12/A2` | raw 48横向冻结线端点和凝结效果 | animation 0“寒光凝冰爆闪”，animation 1“照妖镜寒光束端点”；中间光束由运行时拉伸绘制 | `worlds/snow_road/mechanisms/frost_mirror_beam/` | `FX` / `P-FX-ICE` |
| `R068` | `gen4.f#1` `M5/F13/A4` | raw 49 Turtle巡逻和受击 | 背负小石碑的“玄龟妖”，4组缩壳、巡行、转向和受击动画 | `worlds/snow_road/enemies/stone_turtle/` | `T1O` / `P-W2-ENEMY` |
| `R069` | `gen4.f#2` `M1/F1/A0` | raw 49冻结类型39 | “冰封玄龟”，石碑轮廓和龟壳层次可辨 | `worlds/snow_road/enemies/frozen_stone_turtle/` | `T1O` / `P-FX-ICE` |

### 3.5 特殊关卡、世界灵印和地图

| ID | legacy slot | 原用途和结构 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R070` | `mm0.f#0` `M12/F21/A3` | Angkor坍塌关的火海顶部、内部和启动动画 | “水帘洞地火潮”：animation 0循环火线、1地火内部、2从裂隙喷发的启动过程 | `worlds/huaguoshan/collapse/rising_earthfire/` | `T1/FX` / `P-W0-COLLAPSE` |
| `R071` | `mm0.f#1` `M22/F24/A3` | 多格火炬坍塌主体 | “镇洞石灯崩塌”：animation 0稳定燃烧、1坍塌、2倒塌后余火；左右实例由运行时镜像 | `worlds/huaguoshan/collapse/stone_lantern/` | `BOSS` / `P-W0-COLLAPSE` |
| `R072` | `mm0.f#2` `M2/F0/A0` | 镜头震动期间的下落碎屑 | 两种“地火碎岩”，运行时随机平铺，不承担伤害判定 | `worlds/huaguoshan/collapse/debris/` | `FX` / `P-W0-COLLAPSE` |
| `R073` | `mm1.f#0` `M104/F44/A15` | Siberia Yeti Boss全部状态 | 原创“寒魄狮王”：雪白鬃毛、青玉护额、深灰甲片；15组入场、巡行、扑击、震地、抛射、受击、倒地和败退 | `bosses/frost_lion/` | `BOSS` / `P-BOSS-FROST-LION` |
| `R074` | `mmv.f#0` `M23/F29/A0` 两palette | 教程5x5 Seal及激活覆盖 | “五行封印阵”：frame 0完整未激活底阵、1激活金光、2白闪过渡、3完成印记、`4..28`为25格金木水火土拼图 | `tutorial/five_element_seal/` | `T1`拼成`480x480` / `P-TUTORIAL-SEAL` |
| `R075` | `mmv.f#1` `M1/F1/A1` | World 2灵印 | 六瓣雪莲形“雪岭经卷灵印”，青玉和银白配色 | `relics/snow_road_seal/` | `T1O` / `P-RELIC` |
| `R076` | `mmv.f#2` `M1/F1/A1` | World 1灵印 | 芭蕉叶包围火纹的“火焰山经卷灵印” | `relics/flame_mountain_seal/` | `T1O` / `P-RELIC` |
| `R077` | `mmv.f#3` `M1/F1/A1` | World 0灵印 | 水滴、山石和灵猴云纹组成的“花果山经卷灵印” | `relics/huaguoshan_seal/` | `T1O` / `P-RELIC` |
| `R078` | `mmv.f#4` `M3/F5/A1` | 灵印选择箭头和移动效果 | 五阶段“小筋斗云指针”，方向由运行时翻转，不使用原箭头轮廓 | `world_map/relic_selector/` | `T1O` / `P-MAP` |
| `R079` | `mmv.f#5` `M12/F4/A0` | 商店4级最大HP升级图标 | “护体毫光”四个等级：一至四层莲瓣光轮逐级完整，对应购买后最大HP `5..8` | `shop/protection_upgrades/` | `T1O` / `P-SHOP` |
| `R080` | `ms.f#0` `M20/F19/A0` | 世界地图节点、路线、悟乐标记、底栏计数及商店光标；逐frame合同见5.6 | 18个可达“云路地图图标”逐槽重做；frame `15`固定为`96x96`全透明兼容占位，不虚构Boss节点、世界入口或额外收集物 | `world_map/icons/` | `T1O` / `P-MAP` |
| `R081` | `ms.f#1` `M13/F1/A0` | 地图共用地面/路线底图 | 青绿山川、云海和金色云路组成的“通关文牒地图底板” | `world_map/cloud_route_base/` | `SCREEN` / `P-MAP` |
| `R082` | `ms.f#2` `M21/F1/A0` | World 0地图标题/组合frame | 无可读文字的瀑布山形、青灰节点底座和标题承托纹；运行时在安全区绘制`world.huaguoshan.name` | `world_map/headers/huaguoshan/` | `HUD` / `P-MAP` |
| `R083` | `ms.f#3` `M16/F1/A0` | World 1地图标题/组合frame | 无可读文字的赤峰、芭蕉风纹和标题承托纹；运行时绘制`world.flame_mountain.name` | `world_map/headers/flame_mountain/` | `HUD` / `P-MAP` |
| `R084` | `ms.f#4` `M27/F1/A0` | World 2地图标题/组合frame | 无可读文字的雪峰古寺、经幡和标题承托纹；运行时绘制`world.snow_road.name` | `world_map/headers/snow_road/` | `HUD` / `P-MAP` |

### 3.6 主角资源

| ID | legacy slot | 原用途和结构 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R085` | `o.f#0` `M222/F141/A49` 多palette | 主角49组动画 | 悟乐完整动作集，逐动画ID映射见第4节；全部frame使用同一角色母版，不保留原帽子、金发、蓝衣和锤子轮廓 | `hero/wule/` | `HERO` / `P-HERO` |
| `R086` | `o.f#1` `M19/F6/A1` | 主角着火覆盖动画 | “悟乐三昧火缠身”：6帧火焰沿脚到肩卷起，悟乐轮廓始终清楚 | `hero/wule/status_burning/` | `HERO/FX` / `P-HERO-FX` |

### 3.7 21个音频槽

音频全部从空白工程制作，禁止导入原MIDI。SFX输出`48 kHz/24-bit WAV`母版和运行时`OGG`；音乐输出无缝循环`OGG`和独立工程归档。

| ID | legacy slot | 源码事件 | 新音频的具体内容 | 建议长度 | 输出路径 |
| --- | --- | --- | --- | --- | --- |
| `R087` | `snd.f#0` | 开关/机关 | 木鱼轻击、铜簧弹回、短促符阵低频 | `0.25-0.45s` | `audio/sfx/switch.ogg` |
| `R088` | `snd.f#1` | 谜题或竞技场启动 | 小锣、三枚木鱼递进和法阵低鸣 | `0.8-1.4s` | `audio/sfx/riddle_start.ogg` |
| `R089` | `snd.f#2` | 死亡/主动召回开始 | 毫毛燃起、云烟后掠和一记低钟 | `0.9-1.5s` | `audio/sfx/recall_death.ogg` |
| `R090` | `snd.f#3` | 宝箱开锁 | 木匣摩擦、铜扣弹开 | `0.25-0.5s` | `audio/sfx/chest_open.ogg` |
| `R091` | `snd.f#4` | 宝箱发奖 | 三音原创灵光句，末音清晰落点 | `0.6-1.0s` | `audio/sfx/chest_reward.ogg` |
| `R092` | `snd.f#5` | 主角受伤 | 悟乐原创短呼、护体毫光破裂；建立独立声线，不模仿现代西游角色配音 | `0.2-0.45s` | `audio/sfx/hero_hurt.ogg` |
| `R093` | `snd.f#6` | 铁棒击中硬物 | 如意铁棒的乌金棒身撞岩、短金属回弹 | `0.15-0.35s` | `audio/sfx/staff_blocked.ogg` |
| `R094` | `snd.f#7` | 源码`SOUND_SFX_MINE`，同时用于爆炸、重型破坏和Boss败退阶段 | 中性低频爆裂、重击、碎石和灵力崩散；不能只带丹炉材质，否则花果山Boss死亡时失配 | `0.7-1.3s` | `audio/sfx/heavy_burst.ogg` |
| `R095` | `snd.f#8` | 门移动/机关工作 | 石门摩擦、符链滑动、木制机括 | `0.6-1.2s` | `audio/sfx/gate_working.ogg` |
| `R096` | `snd.f#9` | 检查点 | 土地神龛小铃、香火轻燃 | `0.6-1.0s` | `audio/sfx/checkpoint.ogg` |
| `R097` | `snd.f#10` | 敌人受击/Boss命中 | 妖气破裂、低沉短吼，不使用现成配音样本 | `0.25-0.55s` | `audio/sfx/enemy_hurt.ogg` |
| `R098` | `snd.f#11` | 可破坏墙 | 岩块崩裂、符纸撕开和细砂落地 | `0.35-0.7s` | `audio/sfx/wall_break.ogg` |
| `R099` | `snd.f#12` | 幌金绳命中/回收 | 编绳甩出、金环扣住、快速回收 | `0.25-0.5s` | `audio/sfx/golden_rope.ogg` |
| `R100` | `snd.f#13` | 水体 | 入水、玉泉涌动和短气泡 | `0.45-0.9s` | `audio/sfx/jade_water.ogg` |
| `R101` | `snd.f#14` | 石头落地/重门关闭 | 镇山石低沉撞地、少量碎石尾音 | `0.3-0.65s` | `audio/sfx/stone_land.ogg` |
| `R102` | `snd.f#15` | 关卡完成 | 竹笛上行四音、锣点和短弦收束 | `1.5-2.5s` | `audio/music/stage_clear.ogg` |
| `R103` | `snd.f#16` | World 0循环音乐 | 原创洞箫、阮、木鱼和瀑水质感的花果山主题 | `45-75s`无缝循环 | `audio/music/huaguoshan.ogg` |
| `R104` | `snd.f#17` | World 1循环音乐 | 原创堂鼓、低弦、埙和风道呼啸的火焰山主题 | `45-75s`无缝循环 | `audio/music/flame_mountain.ogg` |
| `R105` | `snd.f#18` | World 2循环音乐 | 原创钟磬、埙、低音弓弦和风雪氛围的雪岭主题 | `45-75s`无缝循环 | `audio/music/snow_road.ogg` |
| `R106` | `snd.f#19` | 标题音乐 | 原创五声音阶主题动机，以鼓点、竹笛和云海风声展开 | `50-90s`无缝循环 | `audio/music/title.ogg` |
| `R107` | `snd.f#20` | Game Over | 低钟、三音下行弦句和毫光熄灭声；不是商店音乐 | `2.0-4.0s` | `audio/music/game_over.ogg` |

### 3.8 标题、Tips文本和UI

| ID | legacy slot | 原用途和结构 | 新的具体替代物 | 输出路径 | 档位/提示词 |
| --- | --- | --- | --- | --- | --- |
| `R108` | `spl.f#0` PNG | 原Logo | 新发行名的原创中文主Logo和英文副标；祥云、云路碑纹和经卷结构，不使用原字形、宝石构图或现代西游作品的兵器标识 | `brand/logo/` | 矢量+高清PNG / `P-BRAND` |
| `R109` | `spl.f#1` PNG | 标题背景 | `960x1280`全屏“悟乐立于花果山云崖、远处水帘洞和云路”的原创彩绘场景 | `brand/title/title_background.png` | `SCREEN` / `P-BRAND-SCENE` |
| `R110` | `spl.f#2` PNG | 原版权行图片 | 删除图片；`brand/legal.json`只保存年份、权利主体和法律链接，并依赖`D006`提供`legal.copyright`本地化模板；主体未确定前不得发布 | `brand/legal.json` | 文本重写，不生成旧式版权图 |
| `R111` | `tips.f#0` `M16/F12/A4` | 方向移动Tips插图 | 方向键/WASD移动示意、四向石猴姿态和格子箭头；不画手机键盘 | `tutorial/tips/movement/` | `HUD` / `P-TUTORIAL` |
| `R112` | `tips.f#1` `M3/F2/A0` | 主动召回Tips | `ENTER`键帽、毫毛云烟和返回土地神龛的两阶段图 | `tutorial/tips/recall/` | `HUD` / `P-TUTORIAL` |
| `R113` | `tips.f#2` `M2/F2/A0` | 检查点Tips | 土地神龛未点亮/已点亮两图，并显示`SPACE`重置 | `tutorial/tips/checkpoint/` | `HUD` / `P-TUTORIAL` |
| `R114` | `tips.f#3` `M1/F1/A0` | 罗盘Tips | 顶部“寻经针”指向下一土地神龛的单图 | `tutorial/tips/compass/` | `HUD` / `P-TUTORIAL` |
| `R115` | `tips.f#4` `M8/F7/A0` | 托石过久会被压Tips | 石猴托起镇山石、手臂颤抖、石头下压到危险标记的7阶段图 | `tutorial/tips/crush_warning/` | `HUD` / `P-TUTORIAL` |
| `R116` | `tipst.f#0` 非精灵文本 | 8页Tips文字 | 重写为第7.2节列出的8页西游主题桌面按键文案；英文保持同页序和同占位符 | `text/tips.zh-CN.json`；`text/tips.en-US.json` | 文本重写 / `P-TEXT` |
| `R117` | `ui.f#0` `M21/F12/A2` | 启动Loading动画 | animation 0原创工作室“云印展开”，animation 1载入时“经卷翻页”；移除原厂牌表达 | `ui/boot/` | `SCREEN/FX` / `P-UI` |
| `R118` | `ui.f#1` `M82/F29/A0` | Java源码未发现直接加载的遗留UI图集 | 删除且不迁移；构建中禁止存在该旧块的兜底路径，reachability测试确认无调用 | 无输出 | 删除 |
| `R119` | `ui.f#2` `M38/F22/A0` | 顶/底HUD、数字、生命、钥匙和罗盘 | “通关文牒HUD”：石猴头像、救命毫毛、金银符钥、寻经针、护体毫光、灵蕴珠、云路印及0-9数字全套 | `ui/gameplay_hud/` | `HUD` / `P-UI` |
| `R120` | `ui.f#3` `M5/F5/A0` | 软键、菜单条和选择标记 | 桌面版`SPACE/ENTER/S`键帽、卷轴菜单条、左右莲瓣选择标记，不显示手机软键 | `ui/menus/controls_and_selection/` | `HUD` / `P-UI` |
| `R121` | `ui.f#4` `M1/F0/A0` | 结算勋章 | 金边青玉“功德莲印”，用于全灵蕴、全云路印、零受伤、零重试四个位置 | `ui/results/merit_medal/` | `T1O` / `P-UI` |

## 4. 主角49组动画逐槽映射

新动画允许增加视觉中间帧，但逻辑动作总时长和命中tick继续由Runtime控制。左右镜像可以共享动作设计，但必须输出独立pivot验证结果。

| animation | 原调用语义 | 悟乐动作 | 原序列帧数 | 制作要求 |
| ---: | --- | --- | ---: | --- |
| `0` | 向上静止 | 背向站立，铁棒斜收在身后 | 1 | 脚底不动 |
| `1` | 向右静止 | 右向警戒呼吸 | 8 | 循环，无水平漂移 |
| `2` | 向下静止 | 正面站立 | 1 | HUD可读轮廓 |
| `3` | 向左静止 | 左向警戒呼吸 | 8 | 与1同节奏但独立锚点 |
| `4` | 向上移动 | 背向轻步 | 8 | 配合`18 -> 12 -> 6 -> 0`插值 |
| `5` | 向右移动 | 右向轻步 | 4 | 脚步循环无滑步 |
| `6` | 向下移动 | 正面轻步 | 8 | 头饰不遮前方机关 |
| `7` | 向左移动 | 左向轻步 | 4 | 与5步幅一致 |
| `8` | 向右推物 | 右肩抵住镇山石 | 4 | 双手接触格边`x=24` |
| `9` | 向左推物 | 左肩抵住镇山石 | 4 | 双手接触格边`x=0` |
| `10` | 受伤 | 护体毫光碎裂、身体后仰 | 4 | 击退由Runtime完成 |
| `11` | 托住头顶重物 | 双臂撑起镇山石并颤抖 | 7 | 物体中心对齐头顶格 |
| `12` | 死亡 | 石猴化作救命毫毛与云烟 | 5 | 不使用眩晕星星 |
| `13` | 向上铁棒 | 背向举棒砸击上格 | 7 | 第3逻辑tick命中 |
| `14` | 向右铁棒 | 右向短促砸击 | 8 | 第3逻辑tick命中 |
| `15` | 向下铁棒 | 正面下砸 | 8 | 第3逻辑tick命中 |
| `16` | 向左铁棒 | 左向短促砸击 | 8 | 第3逻辑tick命中 |
| `17` | 向左解锁 | 面向左贴上符钥 | 2 | 消耗钥匙tick不变 |
| `18` | 向右解锁 | 面向右贴上符钥 | 2 | 消耗钥匙tick不变 |
| `19` | 主动召回 | 拔下毫毛、人物化云消散 | 16 | 总计42逻辑tick后复活 |
| `20` | 向右伸出钩索 | 右手甩出幌金绳 | 4 | 绳头从右手pivot发出 |
| `21` | 向右收回钩索 | 后坐发力向右回收 | 4 | 目标移动由Runtime完成 |
| `22` | 向左伸出钩索 | 左手甩出幌金绳 | 4 | 绳头从左手pivot发出 |
| `23` | 向左收回钩索 | 后坐发力向左回收 | 4 | 目标移动由Runtime完成 |
| `24` | 源码未发现直接setter的右向兼容槽A | 右向悬空挣扎第一循环 | 6 | 仍生成，禁止旧帧兜底 |
| `25` | 源码未发现直接setter的左向兼容槽A | 左向悬空挣扎第一循环 | 6 | 与24镜像但独立pivot |
| `26` | 源码未发现直接setter的右向兼容槽B | 右向下坠摆臂循环 | 6 | 仍生成，禁止旧帧兜底 |
| `27` | 源码未发现直接setter的左向兼容槽B | 左向下坠摆臂循环 | 6 | 与26镜像但独立pivot |
| `28` | 无支撑时向右铁棒 | 悬空右砸，腿部收起 | 7 | 不改变命中tick |
| `29` | 无支撑时向左铁棒 | 悬空左砸，腿部收起 | 7 | 不改变命中tick |
| `30` | 无支撑时向右伸绳 | 悬空右甩幌金绳 | 4 | 绳索pivot保持 |
| `31` | 无支撑时向左伸绳 | 悬空左甩幌金绳 | 4 | 绳索pivot保持 |
| `32` | 无支撑时向右收绳 | 悬空右收绳 | 4 | 不改变目标到达tick |
| `33` | 无支撑时向左收绳 | 悬空左收绳 | 4 | 不改变目标到达tick |
| `34` | 无支撑时向左静止 | 左向失足悬空单帧 | 1 | 与地面静止明显不同 |
| `35` | 无支撑时向右静止 | 右向失足悬空单帧 | 1 | 与34轮廓对应 |
| `36` | 水中向上 | 背向施展避水诀游动 | 4 | 身体中心不漂移 |
| `37` | 水中向右 | 右向划水 | 4 | 水波单独放FX层 |
| `38` | 水中向下 | 正面划水 | 4 | 不扩大碰撞 |
| `39` | 水中向左 | 左向划水 | 4 | 与37步频一致 |
| `40` | 长宝箱奖励 | 开经匣后双手举起法宝 | 16 | tick 37提示音，tick 39发奖 |
| `41` | 向上铁棒撞硬物 | 背向砸中硬物后反弹 | 8 | 保留硬物音效时点 |
| `42` | 向右铁棒撞硬物 | 右砸反弹、手臂震动 | 9 | 保留硬物音效时点 |
| `43` | 向下铁棒撞硬物 | 正面下砸反弹 | 8 | 保留硬物音效时点 |
| `44` | 向左铁棒撞硬物 | 左砸反弹、手臂震动 | 9 | 保留硬物音效时点 |
| `45` | 无支撑时向左撞硬物 | 悬空左砸反弹 | 8 | 不增加滞空时间 |
| `46` | 无支撑时向右撞硬物 | 悬空右砸反弹 | 8 | 不增加滞空时间 |
| `47` | 获得世界灵印 | 双手举起经卷灵印 | 1 | 42 tick动作和141 tick输入锁不变 |
| `48` | 短奖励 | 快速接住灵蕴珠、回元果或小道具 | 9 | 连续拾取窗口不变 |

动画`24..27`在当前Java主流程中没有找到直接设置点，因此本文不冒充其原语义已经确认。第一阶段仍给它们指定新画面，确保意外可达时不会回退原素材；后续若字节码可达性证明永久不可达，可以连同槽位一起删除。

## 5. 多语义资源的子槽合同

### 5.1 `R023` 拾取效果

| animation | 新效果 | 绑定事件 |
| ---: | --- | --- |
| `0` | 经匣开光，金色卷纹从箱口上升 | 宝箱tick 39发奖 |
| `1` | 玄冰凝结环向中心收束 | Freeze Hammer成功冻结 |
| `2` | 云路印朱金灵光 | 云路印/红色永久收集物 |
| `3` | 灵蕴珠紫青流光被吸入HUD | 紫色灵蕴拾取 |
| `4` | 法宝与功德莲印星芒 | 工具、生命或结算奖励 |
| `5` | 五行封印白金激活闪 | 教程最终seal触发 |

### 5.2 `R037` 火焰与Boss效果

| animation | 新效果 | 说明 |
| ---: | --- | --- |
| `0` | 三昧火舌 | 保留空帧及1、2、3格reach阶段，伤害仍由Runtime计算 |
| `1` | 妖王法力冲击 | 给特殊Boss抛射/冲击调用，使用中性赤金核心，世界专属色由粒子层叠加 |
| `2` | 碧鳞蟒王尾击 | 只服务蟒王尾击状态，视觉宽度不得改变命中区 |

### 5.3 Boss动画索引

Boss新动画必须保留原索引数量，不能只做一套待机和一套受击。

| 资源 | 索引合同 | 新动作分组 |
| --- | --- | --- |
| `R013` 碧鳞蟒王 | `0..8` | 潜伏待机、探身、易伤、受击、回收、换列、扑空、重伤和败退各占原状态槽 |
| `R015` 铁甲牛将 | `0/1`左右巡行，`2/3`左右长攻击，`4/5`左右静守，`6/7`左右受击，`8/9`左右冲锋，`10/11`左右战斗待机，`12`重伤，`13`败退 | 每个左右槽单独验证槊尖和脚底pivot |
| `R073` 寒魄狮王 | `0/1`左右待机，`2/3`左右巡行，`4/5`左右爪击，`6/7`左右扑跃，`8/9`左右受击，`10/11`左右恢复，`12`怒吼，`13/14`左右败退 | Boss原点、屏幕外入场和抛射点均写入manifest |

Boss动作名称是新主题的制作语义；实际切换仍按原状态机索引，不依据动画是否播放完自行改Boss阶段。

### 5.4 世界地形逐frame命名

三个世界的地形frame不在文档中重复151行相同规则，而使用以下无缺口公式：

```text
R003 old frame 00..39 -> huaguoshan/terrain/tiles/frame_000..039
R007 old frame 00..65 -> flame_mountain/terrain/tiles/frame_000..065
R011 old frame 00..44 -> snow_road/terrain/tiles/frame_000..044
```

每个frame必须在`terrain-frame-map.csv`中自动列出：`legacy_frame`、`scene_ids`、出现层、raw ID、坐标计数、邻接签名、`north/east/south/west`接缝合同、逻辑遮挡职责和`new_frame`。生成前只依据关卡数据计算直边、内角、外角与多格组合，不从旧PNG提取alpha模板。这样每个frame都有唯一替代文件，同时避免凭截图猜测它是墙角还是前景装饰。

### 5.5 水体和多格机关

- `R053`的34个水frame逐一输出`water/frame_000..033`，保留三段8px子层的拼接边。
- `R031`必须先按`[20,21;22,23]`合成花果山`2x2`灯龛，`R048`必须先按`[22,23;20,21]`合成火焰山`2x2`灯龛；验收还要覆盖数据中的水平镜像排列。
- `R017`门、`R038`钉柱必须先制作完整`96x192`母版；`R040`镇关力士按原最大约`54x53`逻辑visual bounds制作至少`216x212`透明母版，并以底格脚底为pivot。它们都不能逐格生成后硬拼。
- `R071`坍塌石灯、`R074`五行封印和三个Boss必须从完整构图切片，所有切片共用一份原点表。
- `R049/R060`两种经匣关闭frame的alpha必须覆盖同格奖励的完整visual bounds。

### 5.6 `R080` 云路地图19个frame

`ms.f#0`加载到Java `aClassfArr[17]`。以下映射来自世界地图构建、移动、底栏、灵印页和商店绘制调用，不按旧图外观猜测。新运行时可以拆成语义化图片，但兼容层必须保持这些调用结果：

| legacy frame | 源码职责 | 新图职责 | 生成/实现合同 |
| ---: | --- | --- | --- |
| `0` | 已解锁普通节点 | 白金莲瓣普通节点 | 节点中心清楚，不含完成标记 |
| `1` | 锁定普通节点及移动前的未解锁目标 | 灰青封符节点 | 轮廓与frame `0`相同，封条状态不能只靠降亮度 |
| `2` | 已解锁普通路线的重复线段 | 暖金云珠线段 | 允许沿任意角度重复铺设，不能带单向箭头 |
| `3` | 锁定普通路线的重复线段 | 灰封云珠线段 | 与frame `2`尺寸一致，保持地图连接可读 |
| `4` | 锁定秘密路线的重复线段 | 暗青封缄云珠 | 必须区别于普通锁定路线，但不能提前揭示终点奖励 |
| `5` | 秘密节点解锁动画完成前的目标占位 | 半隐青玉节点 | 只在解锁移动阶段出现，不作为常驻已解锁节点 |
| `6` | 地图行者向右/默认朝向 | 向右的小筋斗云悟乐标记 | 小尺寸仍能看出悟乐头部与云尾方向 |
| `7` | 地图行者向左 | 向左的小筋斗云悟乐标记 | 与frame `6`共享脚点和visual bounds |
| `8` | 已解锁秘密路线的重复线段 | 青玉隐云线段 | 与frame `4`同家族、与普通金色路线有形状差异 |
| `9` | 已解锁秘密节点 | 青玉隐云节点 | 不含完成标记，不能画成独立奖励物 |
| `10` | 关卡云路印进度与全局云路印累计 | 朱砂方形云路印 | 同一语义同时用于节点浮层和地图底栏 |
| `11` | 全局灵蕴珠银行；灵印页右下的行囊坊入口 | 紫青圆形灵蕴珠/行囊坊招牌 | 两处复用是有意的：入口用该货币图标表达“用灵蕴珠升级”，不解释为第四种收集物 |
| `12` | 救命次数 | 朱绳救命毫毛 | 与游戏HUD和结算中的生命图标使用同一母版 |
| `13` | 拥有三条以上连接的已解锁普通分岔节点 | 四向莲瓣云台 | 可连接性由轮廓表达，不能看成Boss专属节点 |
| `14` | 行囊坊四行升级项的当前选择光标 | 暖金莲瓣行选标 | 只服务商店列表，不加入关卡地图节点集合 |
| `15` | Java源码及原JAR可达调用点均未发现 | 无可视语义 | 固定输出`96x96`全透明兼容占位并标记`unreachable_compat`；以reachability测试锁定，禁止复用为新图标 |
| `16` | 世界选择/行囊坊背景上按24逻辑像素平铺的覆盖纹理 | 低对比淡墨云纹遮罩 | 必须可无缝平铺，不能遮住灵印、商店图标或文字 |
| `17` | 普通节点通关覆盖 | 暖金功德完成印 | 叠在frame `0/13`之上，不替换节点本体 |
| `18` | 秘密节点通关覆盖 | 青玉秘境完成印 | 叠在frame `9`之上，与frame `17`形状也有差异 |

Boss终点仍由`R075/R076/R077`在对应普通节点上叠加世界灵印。世界入口属于灵印页，不是`R080`的额外节点类型。

## 6. 121块之外的生产与发布资源

以下项目不在`decoded/sprites/manifest.json`的121块中，但同样属于必须纳入生产闭包的运行资源、制作归档或发行元数据。

| ID | 当前来源 | 新的具体替代物 | 输出路径 | 结论 |
| --- | --- | --- | --- | --- |
| `D001` | `demoSpr.bin sprite 0` `M22/F9/A1` | 原创土地神9种表情头像：平静、微笑、惊讶、严肃、担忧、指引、赞许、闭眼和告别 | `tutorial/portraits/earth_spirit/` | 生成 |
| `D002` | `demoSpr.bin sprite 1` `M15/F9/A0` | 9个原创对话标点：感叹、疑问、汗滴、灵光、方向、警告、确认、沉思和完成 | `tutorial/portraits/emotes/` | 生成 |
| `D003` | `demoSpr.bin sprite 2` `M3/F2/A1` | 土地神头像后的卷云底座，普通/高亮两状态 | `tutorial/portraits/backplate/` | 生成 |
| `D004` | `decoded/fonts/freej2me-small.*` | 固定版本`Source Han Sans SC` Regular/Bold的40px高清UI字形缓存；字形来自两个首发locale和动态文本并集 | `fonts/source-han-sans-sc/`、`fonts/ui-small.atlas`、`manifests/font-glyphs.txt` | 引入固定版本并归档许可证；重新栅格化 |
| `D005` | `decoded/fonts/freej2me-medium.*` | 同一固定字体版本的48px标题/对话字形缓存 | `fonts/ui-medium.atlas` | 重新栅格化，不复制旧度量图 |
| `D006` | 原127条本地化索引总迁移、关卡标题及Go可见硬编码 | 西游主题`zh-CN`与`en-US`基础UI文本包；`text-key-map.json`把每个旧索引唯一分派给本项或`R024/R116/D007`，不重复写Credits、Tips和demo文件 | `text/zh-CN.json`；`text/en-US.json`；`manifests/text-key-map.json` | 重写；缺键、重复owner或占位符不一致时构建失败 |
| `D007` | `decoded/demo-scripts.json`中的可见对白索引 | 第7.3节的新教程、法宝和Boss对白，提供两种locale；命令索引映射到语义key | `text/demo.zh-CN.json`；`text/demo.en-US.json` | 重写；命令时序保留 |
| `D008` | 原地图/关卡预览和截图 | 由新高清atlas重新渲染的关卡缩略图，不使用旧frame合成 | `ui/previews/` | 重新生成 |
| `D009` | 应用/窗口默认图标 | 原创“经卷中的金色猴毫”图标，16/32/64/128/256/512/1024px | `brand/app_icon/` | 生成 |
| `D010` | 鼠标指针或系统默认选择图 | 云路碑纹圆点指针和小筋斗云选择态；无自定义需求的平台继续用系统指针 | `ui/cursor/` | 生成或明确使用系统资源 |
| `D011` | README、Release、商店文案和截图 | 只使用新Logo、新角色、新世界观文案和新高清运行画面；中英文商店文案分别审校 | `marketing/` | 重新制作 |
| `D012` | 无统一来源记录和统一第三方声明 | 每项资产的作者、提示词、模型、编辑和许可记录；由记录生成第三方声明并归档许可证全文 | `assets/pilgrim/provenance.csv`；`licenses/`；`THIRD_PARTY_NOTICES.md` | 新建并作为发布门禁 |
| `D013` | 场景brief中的`ambient_fx`及“祥云风带”等非旧PNG效果 | 风带、水汽、风雪、余烬、瀑雾和碎屑的纯视觉粒子/程序图元参数；逐scene声明层级、强度、速度、遮挡上限和禁用状态 | `manifests/procedural-fx.json` | 新建；不得改变碰撞、对象速度或20Hz事件tick |
| `D014` | B00所需但尚无统一产物的角色/场景制作母版 | 悟乐四向与比例母版、土地神和三Boss轮廓、通用功能轮廓表、三世界材质/色板/光向圣经；只作为后续生成与人工修图依据 | `design/character-bible/`；`design/boss-bible/`；`design/world-style-bible/`；`design/functional-silhouettes.json` | 新建；`distribution_targets=source_archive` |
| `D015` | 分散在窗口标题、包配置、存档路径和发行脚本中的产品身份 | 最终中英文发行名、稳定产品slug、bundle/application ID、可执行文件名、存档命名空间、默认locale、支持/隐私链接的唯一配置 | `brand/product.json` | 新建；未完成商标和迁移审查前保持`blocked` |

`Source Han Sans SC`只是本文确定的字体家族；实际引入时必须固定准确版本、保存许可证全文和字体文件SHA-256，并由发布审查确认嵌入和分发方式。检查失败时不能退回FreeJ2ME字体图。

## 7. 文本替换基线

### 7.1 统一名词

| 旧显示概念 | 新显示文本 |
| --- | --- |
| Diamond Rush | 最终发行名未定；内部仅用`Project Pilgrim` |
| Explorer / Hero | 悟乐（`Wule`） |
| Angkor Wat | 花果山与水帘洞 |
| Bavaria | 火焰山与芭蕉洞 |
| Siberia | 雪岭雷音古道 |
| Diamonds | 灵蕴珠 |
| Red diamonds | 云路印 |
| Seal | 经卷灵印 |
| Hidden/Fire Crystal（旧对白） | 对应世界的经卷灵印；不是独立进度 |
| Magic circle | 土地神龛 |
| Magic Shop | 云路行囊坊 |
| Hammer | 如意铁棒 |
| Mystic Hook | 幌金绳 |
| Freeze Hammer | 玄冰铁棒 |
| Compass | 寻经针 |
| Lives | 救命毫毛 |
| Health | 护体毫光 |
| Health potion | 回元果 |
| Hits | 受伤次数 |
| Retries | 召回次数 |
| Secret Stage | 秘境 |
| Stage Complete | 关卡完成 |

### 7.2 八页Tips最终文案

`R116`按以下页顺序写入，键位通过绑定变量渲染，默认值为`WASD/方向键`、`SPACE`、`ENTER`和`S`：

1. `使用 {MOVE_KEYS} 移动。面向改变时，第一次按键只会转身；继续按住或再次按下才会前进。`
2. `被困时按 {RECALL_KEY} 召回到最近的土地神龛。离开神龛后召回会消耗一根救命毫毛。`
3. `踏入土地神龛即可记录当前进度。站在神龛上按 {ACTION_KEY}，可立即重置附近机关。`
4. `寻经针会指向下一座土地神龛。路线不明时，先观察顶部指针。`
5. `不要长时间托住镇山石。支撑耗尽后，落石会造成重伤。`
6. `在云路地图按 {MAP_BACK_KEY} 打开经卷灵印页。云路印累计达到10枚和25枚时，会分别显现新的云路。`
7. `金符钥与银符钥只能开启同纹符锁。隐云界门会在云路地图上显现秘境支路。`
8. `灵印页上的灵蕴珠入口通往云路行囊坊，可用灵蕴珠提升护体毫光。`

### 7.3 Demo脚本对白

| 脚本 | 新对白内容 |
| ---: | --- |
| `29` | `土地神：悟乐，花果山的云路被五行封印打乱了。` / `土地神：沿神龛前进，寻回散落的经卷灵印。` / `土地神：先活动筋骨，我会在前方等你。` |
| `10` | `土地神：使用 {MOVE_KEYS} 前进。改变方向时要先转身。` |
| `11` | `土地神：这是寻经针。顶部指针会指向下一座土地神龛。` |
| `13` | `土地神：踏入神龛，香火就会记住你此刻的旅程。` / `土地神：站在神龛上按 {ACTION_KEY}。` / `土地神：附近的石块和机关会恢复到记录时的状态。` / `土地神：这不会消耗救命毫毛。` |
| `15` | 无对白，只保留自动移动和镜头命令 |
| `16` | `土地神：离开神龛后，也能主动返回。` / `土地神：按 {RECALL_KEY} 施展毫毛替身。` / `土地神：召回会消耗一根救命毫毛。` |
| `17` | `土地神：现在按 {RECALL_KEY}，回到刚才的神龛。` |
| `28` | `土地神：前方是五行封印。` / `土地神：踏入阵心，让寻回的灵蕴回应云路。` / `土地神：花果山的试炼，从这里正式开始。` |
| `30` | `土地神：金符钥与银符钥只能开启对应的符锁。` / `土地神：先观察门上的日纹和月纹。` |
| `22` | `土地神：如意铁棒能击碎带裂纹的符墙。` / `土地神：面向障碍，按 {ACTION_KEY} 挥棒。` / `土地神：坚固墙面只会让铁棒反弹。` |
| `33` | `土地神：碧鳞蟒王正守着经卷灵印。` / `土地神：引落镇山石，在它探身露出弱点时出手。` |

其他菜单、结算、存档、世界解锁和错误文本统一从当前locale文本包读取，不允许在Go绘制代码中继续硬编码原显示文本。

### 7.4 本地化与字形闭包

- 首发locale固定为`zh-CN`和`en-US`。两者必须拥有完全相同的语义key集合；英文是按相同世界观独立审校的表达，不能保留原作英文句子作为缺省值。
- `text-key-map.json`必须覆盖原127条文本、Credits、8页Tips、全部demo对白索引、世界/关卡名、菜单、HUD、结算、商店、存档错误、窗口标题和Go代码中所有面向玩家的硬编码字符串。
- 同一key的占位符名称与数量必须一致，例如`{MOVE_KEYS}`、`{ACTION_KEY}`、`{RECALL_KEY}`和动态计数；未知占位符、漏占位符或字符串拼接出的半句均使构建失败。
- `R082..R084`只提供无文字世界标题装饰框，运行时分别绘制`world.huaguoshan.name`、`world.flame_mountain.name`和`world.snow_road.name`。除已批准的中英文Logo及固定桌面键帽外，不得把可本地化文字烘焙进PNG。
- `brand/legal.json`只存共享事实，`legal.copyright`在两个locale包中保存模板。Credits也必须提供`credits.zh-CN.json`和`credits.en-US.json`，不能由中文图片代替。
- `font-glyphs.txt`由两个locale的全部最终值、悟乐/Wule等专名、数字、标点、键帽和运行时格式字符生成。字体atlas缺任一字形、出现豆腐块或文本超出容器时构建/视觉验收失败。
- 已支持locale缺键时禁止逐键回退；构建阶段必须修齐。系统请求未支持locale时可以整体回退到产品默认`zh-CN`，但不得回退到原作文本。

## 8. GPT-image提示词库

所有提示词共用以下风格与独立创作限制：

```text
原创西游题材高清国风彩绘2D游戏资产，统一左上光源，清晰轮廓，
面向24x24逻辑格但以4倍分辨率制作，材质细节受控，
不得参考或模仿任何现代西游影视、动画、漫画、游戏，也不得参考Diamond Rush素材；
不要水印、伪3D截图或未经brief要求的文字与展示底座。
```

每次调用还必须附加以下场景上下文，缺少任一字段不得开始生成：

```text
asset_id: Rxxx or Dxxx
scene_id: 第2.6至2.8节中的场景；跨场景资产列出全部scene_id
character_id: 悟乐相关素材固定为wule；非角色素材写none
functional_role: 阻挡/可推/可钩/危险/奖励/纯装饰
logical_footprint: 占格、pivot、允许的visual overflow
authored_context: 来自stage JSON的层、关键坐标、相邻frame和复用次数
lore_context: 总体规范3.3中的区域职责、当前失序和设施系统
required_states: 静止、预警、激活、受击、打开、关闭等
scene_materials: 对应场景主材质与点缀材质
light_and_contrast: 光源、背景明度和目标可读性
output_mode: transparent_sprite / opaque_tile / opaque_scene_board / logo
forbidden_motifs: 该场景不能出现的世界材质、奖励暗示和现代作品特征
```

`transparent_sprite`使用透明背景；`opaque_tile`必须铺满需要实心的格区；`opaque_scene_board`和标题背景必须是不透明完整场景，不能套用透明精灵要求。Logo单独按品牌brief决定透明或单色版本。

场景板提示使用第2.6至2.8节的scene brief和关卡数据生成的几何遮罩，输出`960x960`概念板。它只用于检查材质、照明和视觉层级，最终格子位置仍由关卡JSON和正式atlas渲染，不直接把概念板当游戏背景。

| Prompt key | 追加brief |
| --- | --- |
| `P-HERO` | `原创主角悟乐（Wule），由花果山废弃云路碑香火化生的年轻石猴行者；深棕短毛、青绿色短褂、朱红腰绳、乌金护腕，短身敏捷比例，神情机敏乐观；按统一角色母版输出四向设定和动作关键姿势，不生成旧式探险帽，不作傲慢战神姿态，不参考孙悟空的现代改编造型。` |
| `P-HERO-FX` | `同一主角悟乐（Wule），严格保持脸型、毛色、短褂、腰绳和护腕一致；三昧火沿脚、腰、肩分层缠绕，悟乐轮廓始终清楚，火焰与角色分层。` |
| `P-W0-TERRAIN` | `花果山水帘洞地形组件，青灰岩、翠藤、苔藓、根系、瀑水和暖金石灯，正交格子视图，可无缝拼接。` |
| `P-W0-OBJECT` | `花果山可动物件，天然青灰石、浅金镇字符，圆形重心明确，不能与固定墙体混淆。` |
| `P-W0-ENEMY` | `花果山蛇妖，青藤或赤练纹样，夸张头部和可读朝向，危险但不血腥。` |
| `P-W0-COLLAPSE` | `水帘洞地火脉、崩塌石灯和碎岩，赤金火光与青灰洞壁形成冷暖对比，多格整体构图。` |
| `P-W1-TERRAIN` | `火焰山芭蕉洞地形组件，赤岩、黑铁、青铜、玉泉水渠和朱红布幔，正交格子视图。` |
| `P-W1-OBJECT` | `火焰山圆形可动物件，赤铜包边、炉火裂纹或镇火符，重心和旋转方向清楚。` |
| `P-W1-MECH` | `火焰山机关，青铜机括、黑铁导轨、芭蕉风纹和日月符号；包含整身持降魔长矛的铜甲镇关力士，结构连续、方向和敌意一眼可辨。` |
| `P-W1-ENEMY` | `原创火蜥妖与泉蜥妖，低伏爬行轮廓、青铜背脊、日月腹纹，不采用现代作品妖怪造型。` |
| `P-W2-TERRAIN` | `雪岭雷音古道地形组件，蓝白冰岩、深灰古寺砖、冻木梁、朱红经幡和薄霜，正交格子视图。` |
| `P-W2-OBJECT` | `雪岭可动物件，青玉冰脉和结霜石面，圆形重心明确，边缘与冰墙有亮度差。` |
| `P-W2-MECH` | `雪岭机关，寒光照妖镜、冰钟乳、青玉套环和古寺铜件，结构清晰，可按格拼接。` |
| `P-W2-ENEMY` | `原创雪岭猿傀、寒针蜂妖或负碑玄龟妖，统一木面具、青玉和朱绳识别体系；猿傀必须是整脸面具与符绳关节的敌对傀儡，不能像悟乐的换色版本。` |
| `P-COMMON` | `跨世界通用法器或机关，云纹、符纸、经卷、青玉和适量金属，交互状态可读，透明背景。` |
| `P-ENEMY` | `原创石蝎精，贴墙爬行，石质甲壳和少量青苔，四向轮廓清晰，死亡表现为石屑。` |
| `P-FX` | `国风法术粒子，符尘、灵光、石屑或火星，中心和发射方向清楚，透明背景，不遮挡整格。` |
| `P-FX-ICE` | `青白玄冰凝结、薄霜裂纹和少量青玉高光，冰内对象仍清楚可辨，透明背景。` |
| `P-WATER` | `正交水体瓦片，水面、内部、流向、边角和源头共享波纹尺度；按scene_id生成火焰山青绿玉泉/青铜渠壁或雪岭淡青冰泉/结霜石槽，两套拓扑相同但材质独立。` |
| `P-BOSS-PYTHON` | `原创碧鳞蟒王，青碧鳞、金腹甲、断角和藤纹，多格身体，威严而非写实恐怖。` |
| `P-BOSS-BULL` | `原创铁甲牛将，黑铁层甲、青铜牛首面具和双手降魔槊，不参考任何现代牛魔王形象。` |
| `P-BOSS-FROST-LION` | `原创寒魄狮王，雪白鬃毛、青玉护额、深灰甲片和冰霜爪痕，多格Boss比例。` |
| `P-RELIC` | `原创经卷灵印图标，传统印章与自然主题结合，单色轮廓也能区分世界，透明背景。` |
| `P-MAP` | `祥云航路地图组件，山川云海、通关文牒、莲瓣节点和金色路线，俯视平面构图。` |
| `P-SHOP` | `护体毫光四级升级图标，莲瓣光轮从一层到四层递增，青玉和暖金配色，透明背景。` |
| `P-TUTORIAL` | `云纹卷轴教程UI、土地神龛和桌面键帽，信息层级清楚，不出现手机数字键。` |
| `P-TUTORIAL-SEAL` | `完整5x5五行封印阵，金木水火土纹样形成一个整体，先生成完整构图再切25格。` |
| `P-UI` | `通关文牒式游戏HUD和菜单组件，青玉、朱红、乌金、白色文字，多色但克制，按钮状态清晰。` |
| `P-BRAND` | `原创西游格子解谜游戏字标，书法骨架、经卷和祥云结构，避免宝石字样和原游戏Logo构图。` |
| `P-BRAND-SCENE` | `主角悟乐（Wule）站在花果山云崖，外观严格遵循悟乐角色母版，远处清晰可见水帘洞和通向群山的云路；960x1280竖幅标题背景，主体可检查，不做模糊氛围图。` |
| `P-SCENE-BOARD` | `严格沿团队提供的10x10几何遮罩绘制不透明正交场景板，只决定材质、照明、危险色和层级，不增删墙体、通道、机关或对象。` |
| `P-TEXT` | 不调用图像模型；由编剧重写并走文本审校 |

GPT-image只负责概念母版或单个静态对象。动画关键帧、tile拼接、透明边清理、pivot和atlas打包必须人工完成。

## 9. 制作批次

1. `B00`设计冻结：完成`D014/D015`，确定发行名、悟乐（Wule）角色母版、三Boss设定、土地神头像、三世界色板、功能轮廓表和41个scene brief；此批只批准母版与产品身份配置，不导出最终atlas。
2. `B01`通用闭包：完成`R016`、`R018..R029`、`R037`、`R041`、`R043`、`R046..R047`、`R049`、`R054..R055`、`R057..R058`、`R085..R086`、`R108..R117`、`R119..R121`及`D001..D007`。与`B02`合并后必须让教程和花果山Stage 1完全脱离旧视听资源。
3. `B02`花果山：完成`R001..R004`、`R013..R014`、`R017`花果山variant、`R031`、`R034`花果山variant、`R036`、`R042`花果山variant、`R056`花果山variant、`R070..R072`、`R074`、`R077`。
4. `B03`火焰山：完成`R005..R008`、`R015`、`R017`火焰山variant、`R030`、`R032..R035`、`R038..R040`、`R044..R045`、`R048`、`R051..R053`、`R056`火焰山variant、`R066`火焰山variant、`R076`；其中`R033/R034/R051/R053/R056/R066`必须按火焰山scene生成，不得只换调色板名称。
5. `B04`雪岭：完成`R009..R012`、`R017`雪岭variant、`R033/R034/R042/R051/R053/R056`雪岭variant、`R050`、`R059..R069`、`R073`、`R075`；`R066`在该范围内完成雪岭variant。
6. `B05`地图与派生：完成`R078..R084`及`D008..D013`，重新渲染预览、应用图标、光标、宣传图、程序化环境效果和来源清单。`R080` frame `15`只做透明兼容占位，不进入图像生成队列。
7. `B06`音频：完成`R087..R107`共21槽原创音频，并统一响度、循环点、事件优先级和许可证记录。
8. `B07`集成清理：按审计结论删除不可达`R118`及所有旧资源回退路径，执行全关卡逐tick回归、全场景视觉验收和发布包扫描。

上述主批次覆盖`R001..R121`和`D001..D015`全部编号；重复出现的ID只表示同一资源的世界variant在后续批次补齐，不是复制一个通用PNG。构建清单为每个ID记录一个`primary_batch`，每个variant另记`completed_in_batch`，编号缺失或variant未完成时禁止进入`B07`。

每个批次必须同时更新`provenance.csv`和atlas manifest，不能先把图塞进代码再补来源。

## 10. 覆盖与验收

### 10.1 机械覆盖

- `R001..R121`必须连续、唯一且恰好121项。
- `D001..D015`全部有文件或明确的系统资源结论。
- `R085`的animation `0..48`恰好49项，无缺号。
- 三个地形映射分别恰好40、66、45个frame；水体恰好34个frame。
- 21个音频文件全部存在，ID `0..20`无缺号。
- `stage-scenes.json`恰好包含41个scene ID：World 0为14、World 1为13、World 2为14；每个stage索引只出现一次。
- 每个非品牌资产的manifest至少引用一个scene ID或明确标记为跨场景UI；引用不存在的scene ID时构建失败。
- 新资源manifest中不允许目标路径指向`decoded/sprites`、`decoded/audio`或`decoded/fonts`。
- `assets.json`恰好包含121个`R`项和15个`D`项，所有输出、variant、locale、依赖、分发目标和SHA-256可验证；发布时每项状态均为`packaged`。
- `runtime-resource-map.json`中的每个可达语义键唯一解析到新资源；`animation-events.json`和`audio-events.json`覆盖全部运行时事件且无悬空引用。
- `zh-CN`与`en-US`文本key及占位符集合完全一致，`font-glyphs.txt`覆盖所有最终文本；地图标题、Credits和法律文本不依赖烘焙图片。
- `procedural-fx.json`中的每个效果至少绑定一个真实scene或明确的跨场景UI状态，且不写入玩法状态。

### 10.2 视觉验收

- 在`960x1280`内部画布下检查所有关卡，不能出现旧PNG、诊断色块、纯色缺块或拼接断裂。
- 门、旗幡、钉柱、持矛力士、五行封印和Boss先检查完整构图，再检查切片。
- 宝箱关闭时，奖励在所有动画tick都不可见。
- 主角、敌人和可动物体在移动插值中pivot不跳动；视觉越界不能改变碰撞。
- 灵蕴珠、云路印、避水珠、金银符钥、救命毫毛和回元果在灰度下仍可区分。
- 每关至少验收入口、首次新机制、关键经匣、检查点和出口/Boss五类真实视口；不适用项在scene manifest中说明原因。
- 悟乐在游戏精灵、HUD头像、教程插图、地图标记和标题背景中必须匹配同一`character_id=wule`母版；脸型、毛色、短褂、腰绳和护腕不能随生成批次漂移。
- 每个关卡视口都能说明当地机关属于哪套设施系统；出现无法解释来源的孤立工业件、异世界材质或与本地冲突无关的装饰时，场景验收失败。
- 火焰山场景不得误用雪岭寒风葫芦，雪岭风瓮/阵枢不得直接复用火焰山铜材质；跨世界variant必须保持功能轮廓并更换场景材质。
- `R031/R048`按真实`2x2`坐标组合后必须形成完整灯龛，水平镜像时仍无断边、错位或不连续光源。
- 新角色和Logo通过现代西游作品相似性人工复核。

### 10.3 逻辑回归

- 同一输入脚本下，素材替换前后的Runtime状态逐20Hz tick一致。
- 60Hz显示只插值绘制坐标，不改变人物、蛇、落石、Boss或机关速度。
- 新动画增加帧时，锤击、开箱、发奖、冻结、伤害和通关事件tick保持不变。
- 缺失任何新资源时测试或构建立即失败，不能回退旧资源。
- 用干净环境关闭`decoded/sprites`、`decoded/audio`和`decoded/fonts`后，资源解析、构建、启动、教程及三个世界的代表关卡仍成功。

### 10.4 发布扫描

- 发布产物不含原PNG、MIDI、字体图、Logo、版权图、截图、精灵metadata或原对话。
- 搜索不到原品牌和旧世界名称的玩家可见文本。
- 当前决定保留的关卡数据单独登记为第一阶段例外，不能宣传为原创关卡或clean-room内容。
- 原JAR、Java反编译源码和研究资料只留在私有参考区。
- `closure-report.json`必须证明发行包内所有玩家可见及第三方资源都能反查到`assets.json`和`provenance.csv`，构建生成的manifest可反查生成器、提交和输入哈希，`deletion-evidence.json`中的删除项没有加载路径，保留关卡/地图/命令仅来自`retained-content-exceptions.json`。

当以上门禁全部通过时，可以说“121个原资源块及清单外发布素材均已有明确的新资源承接或审计删除结论”；这不等于保留的精确关卡设计已经没有法律风险。

## 11. 自洽审计结论

截至2026-07-13，本文档在**资源语义与闭包合同层面**已经形成可自洽的西游主题替换方案，但实际高清资源、双语文本、字体包、程序化效果、manifest和运行时variant尚未制作完成，不能把“文档闭环”表述为“素材已经完成”或“发布包已经闭环”。

本次本地机械审计快照：

| 审计项 | 目标 | 当前结果 | 判定 |
| --- | ---: | ---: | --- |
| 原资源主表 | 121个连续唯一`R`项 | 121/121 | 通过 |
| 清单外生产/发布项 | 15个连续唯一`D`项 | 15/15 | 通过 |
| 场景矩阵 | World 0/1/2为14/13/14 | 14/13/14，共41 | 通过 |
| 主角动画/音频槽 | 49/21 | 49/21 | 通过 |
| 提示词定义 | 所有被引用key均有定义 | 31个被引用key全部定义；`P-SCENE-BOARD`为独立场景板模板 | 通过 |
| 闭包manifest | 14个主文件及schemas | 已生成5个基础文件：资产、场景、绑定、地形、原资源锁 | 进行中 |
| 场景生产资料 | 41份Brief及真实关键视口mask | 41份Brief、423个三层`10x10`mask | 已生成，待人工审美审核 |
| Wule最终运行资源 | 所有可达美术、音频、字体和文本可独立加载 | `materialized_outputs`仍全部为空 | 未开始 |
| 正式运行时代码的旧视听路径 | 0 | `internal/originalgame`仍有195行`decoded/{sprites,audio,fonts}`引用 | 未闭合 |

`tools/`中另有3行`decoded`输出/说明路径，属于私有研究工具，可不进入发行包；若工具也随公开仓发布，则同样要隔离或改名，不能被正式运行时引用。

已闭合的主要关系：

- 经济只有三类真实职责：`灵蕴珠`同时服务关内配额、结算和商店扣费；`云路印`只永久累计并在`10/25`独立显现世界路线；三枚`经卷灵印`由各世界Boss流程另行持久化，不能被追加成世界解锁前置。旧Crystal对白并入经卷灵印，不再虚构“秘藏经页”。
- 标准工具进程为花果山Stage 4取得如意铁棒、火焰山Stage 3取得幌金绳、雪岭Stage 6把同一铁棒淬炼成玄冰铁棒；同时保留进入火焰山/雪岭时补足等级`1/2`的源码兜底。回访关卡只消费持久化能力，兜底也不伪造宝箱或玄冰升级演出。
- 高频补血物使用普通回元果，避免把大量投放包装成稀有仙果；经匣关闭态继续完全遮蔽奖励，叙事不会提前剧透payload。
- 花果山、火焰山和雪岭的水源、水体、蛇类、风瓮和阵枢按真实出现世界建立variant；功能轮廓一致，材质和场景归属不同。
- `2x2`灯龛、两格门、钉柱、镇关力士、崩塌石灯、旗幡、封印和Boss都先做完整母版再切片，生成模型不能把多格结构当成互不相关的小图。
- World 2 Stage 9的player raw `74`由JAR关卡初始化默认分支删除，不产生画面、碰撞或素材需求；它保留为数据审计项，不能被误判成漏图。
- `R080`的19个地图frame已按源码调用点逐槽定义，frame `15`固定为不可达透明兼容占位；灵蕴珠兼任行囊坊入口图标是货币语义复用，不是额外收集进度。
- 世界标题、Credits、Tips、demo、法律行和通用UI文字统一经过双语语义key及字体字形集，不再由中文烘焙图或Go硬编码承担。
- 祥云风带、水汽、风雪、余烬、瀑雾和碎屑由`D013`统一登记为纯视觉效果，已有scene归属和遮挡合同，不再游离于121块/D项之外。

开始批量生成前仍必须完成以下工程门禁：

1. `tools/pilgrimmanifest`已建立`assets.json`，并从41个stage JSON生成`stage-scenes.json`、423个关键`10x10`三层mask、`asset-scene-bindings.json`、`terrain-frame-map.csv`和388文件原资源锁；下一步先人工审核视口分类和场景美术约束，再冻结schemas。
2. 生成`runtime-resource-map.json`、`animation-events.json`和`audio-events.json`；为`R017/R033/R034/R042/R051/R053/R056/R066`实现穷举式world/scene variant选择，只生成多套图片而渲染器仍固定读一套会再次造成场景冲突。
3. 按第7.4节完成`zh-CN/en-US`、`text-key-map.json`和`font-glyphs.txt`；固定字体版本、权利主体与发行名。`R080` frame `15`按既定合同保留透明占位，后续frame不得索引偏移。
4. 按场景板先验收材质、光源、危险层级、多格拼接和`D013`环境效果遮挡，再制作动画、pivot和atlas；GPT-image输出不能直接作为最终精灵表。
5. 用全部真实关卡视口检查支撑关系、门状态、危险方向、经匣遮蔽和角色辨识度，并运行20Hz逐tick回归，证明高清资源与60Hz插值没有改变玩法。
6. 生成`deletion-evidence.json`和`closure-report.json`，在关闭旧资源目录的干净构建中完成来源、哈希、双语、variant、防回退和发行包扫描；全部条目达到`packaged`后才算生产闭环。

因此，当前文档足以指导“符合实际场景的素材重构”，不再是逐PNG换名清单；设计逻辑和第一批机器清单已经落盘，最终资源生产和Wule运行时尚未闭合。下一步应先审核41份场景Brief并完成`D014`悟乐/三Boss/三世界母版，再按scene分批生成素材，而不是无上下文批量调用图像模型。
