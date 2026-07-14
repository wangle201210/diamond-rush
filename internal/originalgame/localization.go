package originalgame

import (
	"fmt"

	"github.com/wangle201210/zskc/internal/original"
)

type uiTextKey string

const (
	textWindowTitle          uiTextKey = "window.title"
	textMenuNewGame          uiTextKey = "menu.new_game"
	textMenuContinue         uiTextKey = "menu.continue"
	textMenuConfirmStart     uiTextKey = "menu.confirm_start"
	textMenuConfirmDelete    uiTextKey = "menu.confirm_delete"
	textMenuConfirmQuestion  uiTextKey = "menu.confirm_question"
	textMenuNo               uiTextKey = "menu.no"
	textMenuYes              uiTextKey = "menu.yes"
	textHint                 uiTextKey = "tutorial.hint"
	textPromptSkip           uiTextKey = "prompt.skip"
	textPromptContinue       uiTextKey = "prompt.continue"
	textPromptSelect         uiTextKey = "prompt.select"
	textPromptWorlds         uiTextKey = "prompt.worlds"
	textPromptMainMenu       uiTextKey = "prompt.main_menu"
	textLoading              uiTextKey = "common.loading"
	textCongratulations      uiTextKey = "common.congratulations"
	textCheckpoint           uiTextKey = "common.checkpoint"
	textObjectiveAnaconda    uiTextKey = "objective.anaconda"
	textObjectiveTorch       uiTextKey = "objective.torch"
	textObjectiveEnemies     uiTextKey = "objective.enemies"
	textSecretUnlockedFirst  uiTextKey = "secret.unlocked_first"
	textSecretUnlockedSecond uiTextKey = "secret.unlocked_second"
	textStage                uiTextKey = "stage.normal"
	textSecretStage          uiTextKey = "stage.secret"
	textComplete             uiTextKey = "result.complete"
	textDiamonds             uiTextKey = "result.diamonds"
	textRedDiamonds          uiTextKey = "result.red_diamonds"
	textHits                 uiTextKey = "result.hits"
	textRetries              uiTextKey = "result.retries"
	textWorldGoTo            uiTextKey = "world.go_to"
	textWorldNeedRedDiamonds uiTextKey = "world.need_red_diamonds"
	textWorldAngkor          uiTextKey = "world.angkor"
	textWorldBavaria         uiTextKey = "world.bavaria"
	textWorldSiberia         uiTextKey = "world.siberia"
	textWorldShop            uiTextKey = "world.shop"
)

var simplifiedChinese = map[uiTextKey]string{
	textWindowTitle:          "钻石狂潮原作运行版 - %s（世界%d）",
	textMenuNewGame:          "新游戏",
	textMenuContinue:         "继续游戏",
	textMenuConfirmStart:     "开始新游戏将删除",
	textMenuConfirmDelete:    "当前游戏进度。",
	textMenuConfirmQuestion:  "确定要开始吗？",
	textMenuNo:               "否",
	textMenuYes:              "是",
	textHint:                 "提示：",
	textPromptSkip:           "%s：跳过",
	textPromptContinue:       "%s：继续",
	textPromptSelect:         "%s：选择",
	textPromptWorlds:         "%s：封印/商店",
	textPromptMainMenu:       "%s：主菜单",
	textLoading:              "加载中",
	textCongratulations:      "恭喜！",
	textCheckpoint:           "复活点",
	textObjectiveAnaconda:    "击败巨蟒！",
	textObjectiveTorch:       "点燃火炬！",
	textObjectiveEnemies:     "击败所有敌人！",
	textSecretUnlockedFirst:  "恭喜！你已解锁",
	textSecretUnlockedSecond: "一条隐藏路线！",
	textStage:                "第%d关",
	textSecretStage:          "隐藏关卡%d",
	textComplete:             "完成！",
	textDiamonds:             "钻石",
	textRedDiamonds:          "红钻石",
	textHits:                 "受击次数",
	textRetries:              "重试次数",
	textWorldGoTo:            "按 %s 前往",
	textWorldNeedRedDiamonds: "需要 %d 颗红钻石",
	textWorldAngkor:          "吴哥窟",
	textWorldBavaria:         "巴伐利亚",
	textWorldSiberia:         "西伯利亚",
	textWorldShop:            "商店",
}

var allUITextKeys = []uiTextKey{
	textWindowTitle,
	textMenuNewGame,
	textMenuContinue,
	textMenuConfirmStart,
	textMenuConfirmDelete,
	textMenuConfirmQuestion,
	textMenuNo,
	textMenuYes,
	textHint,
	textPromptSkip,
	textPromptContinue,
	textPromptSelect,
	textPromptWorlds,
	textPromptMainMenu,
	textLoading,
	textCongratulations,
	textCheckpoint,
	textObjectiveAnaconda,
	textObjectiveTorch,
	textObjectiveEnemies,
	textSecretUnlockedFirst,
	textSecretUnlockedSecond,
	textStage,
	textSecretStage,
	textComplete,
	textDiamonds,
	textRedDiamonds,
	textHits,
	textRetries,
	textWorldGoTo,
	textWorldNeedRedDiamonds,
	textWorldAngkor,
	textWorldBavaria,
	textWorldSiberia,
	textWorldShop,
}

func tr(key uiTextKey, args ...any) string {
	text, ok := simplifiedChinese[key]
	if !ok {
		return string(key)
	}
	if len(args) == 0 {
		return text
	}
	return fmt.Sprintf(text, args...)
}

func worldDisplayName(world int) string {
	switch world {
	case original.WorldBavaria:
		return tr(textWorldBavaria)
	case original.WorldTibet:
		return tr(textWorldSiberia)
	default:
		return tr(textWorldAngkor)
	}
}

func windowTitleForWorld(world int) string {
	return tr(textWindowTitle, worldDisplayName(world), world+1)
}

var angkorTutorialTexts = [40]string{
	0:  "我得先看看那个宝箱。",
	1:  "推动石头时，",
	2:  "别堵住自己的路。",
	3:  "你可以让所有物体恢复原位，",
	4:  "只需回到最近的复活点并按 " + desktopActionKeyLabel + "。",
	5:  "如果你回不到复活点，",
	6:  "而道路又被堵住，",
	7:  "可随时按 " + desktopRecallKeyLabel + " 返回最近的复活点，",
	8:  "但会失去一条命。",
	9:  "这是一种封印吗？",
	10: "啊！封印有反应了！",
	11: "踩上去看看会发生什么……",
	12: "吴哥窟的宏伟神庙……",
	13: "我终于进来了！",
	14: "出发吧！",
	15: "看，前面有一把魔法锁！",
	16: "收集指定数量的钻石即可开锁。",
	17: "糟了，这扇门上锁了！",
	18: "钥匙一定就在附近……",
	19: "找到罗盘了！它会帮我找到出口。",
	20: "找到神秘锤了！",
	21: "按 " + desktopActionKeyLabel + " 使用。",
	22: "太好了！现在能砸碎那些脆弱的墙了。",
	23: "找到神秘钩索了！",
	24: "隔着一段距离按 " + desktopActionKeyLabel + "，可将物体拉向自己。",
	25: "有意思……也许可以用这个……",
	26: "找到神秘药水了！",
	27: "现在可以在水下呼吸了。",
	28: "找到冰冻锤了！",
	29: "按 " + desktopActionKeyLabel + " 冻结物体。",
	30: "把东西冻住？",
	31: "试试看吧！",
	32: "吴哥窟的最后一间密室！据说火焰水晶就藏在这里……",
	33: "但我有种不祥的预感……",
	34: "银色钻石就在这里……我敢肯定！",
	35: "嗯……附近有一股黑暗力量……",
	36: "大干一场吧！",
	37: "寒冰钻石一定就在西伯利亚的最后一间密室里！",
	38: "成败在此一举！",
	39: "完成本隐藏关需要神秘药水。请先在巴伐利亚第8关取得，再返回这里。",
}

func tutorialText(index int) string {
	if index < 0 || index >= len(angkorTutorialTexts) {
		return ""
	}
	return angkorTutorialTexts[index]
}
