package originalgame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

const (
	startMenuContinue = iota
	startMenuNewGame
)

const (
	startMenuNo = iota
	startMenuYes
)

const startMenuRowHeight = 13

func (g *Game) enterStartMenu(hasProgress bool) {
	g.mode = gameModeStartMenu
	g.startMenuHasProgress = hasProgress
	g.startMenuSelection = startMenuNewGame
	if hasProgress {
		g.startMenuSelection = startMenuContinue
	}
	g.startMenuConfirmNew = false
	g.startMenuConfirmChoice = startMenuNo
	g.pendingMapTarget = -1
	g.message = "Main menu"
	if g.sounds != nil {
		g.sounds.Play(original.SoundTitleMusic)
	}
}

func (g *Game) updateStartMenu(action bool, dy int) {
	if g.startMenuConfirmNew {
		if dy != 0 {
			g.startMenuConfirmChoice = startMenuYes - g.startMenuConfirmChoice
		}
		if !action {
			return
		}
		if g.startMenuConfirmChoice == startMenuNo {
			g.startMenuConfirmNew = false
			return
		}
		g.startNewGame()
		return
	}

	if dy != 0 && g.startMenuHasProgress {
		g.startMenuSelection = startMenuNewGame - g.startMenuSelection
	}
	if !action {
		return
	}
	if g.startMenuSelection == startMenuContinue && g.startMenuHasProgress {
		g.continueGame()
		return
	}
	if g.startMenuHasProgress {
		g.startMenuConfirmNew = true
		g.startMenuConfirmChoice = startMenuNo
		return
	}
	g.startNewGame()
}

func (g *Game) continueGame() {
	if !g.progress.TutorialComplete {
		_ = g.switchWorld(original.WorldAngkor)
		g.loadStage(angkorTutorialStage)
		g.mode = gameModeStage
		return
	}
	world := g.progress.LastWorld
	if err := g.switchWorld(world); err != nil {
		g.message = err.Error()
		return
	}
	g.stageIndex = g.highestUnlockedMapStageForWorld(world)
	g.enterWorldMap()
	if g.sounds != nil {
		g.sounds.Play(worldMusic(world))
	}
}

func (g *Game) highestUnlockedMapStage() int {
	return g.highestUnlockedMapStageForWorld(original.WorldAngkor)
}

func (g *Game) highestUnlockedMapStageForWorld(world int) int {
	progress := g.progress.normalized()
	for stage := min(progress.highestUnlockedForWorld(world), worldStageCount(world)-1); stage >= 0; stage-- {
		if progress.stageUnlockedForWorld(world, stage) {
			if g.worldMap != nil {
				if _, ok := g.worldMap.nodeForStage(stage); !ok {
					continue
				}
			}
			return stage
		}
	}
	return 0
}

func (g *Game) startNewGame() {
	progress := newOriginalProgress()
	if g.progressPath != "" {
		if err := saveOriginalProgress(g.progressPath, progress); err != nil {
			g.message = err.Error()
			return
		}
	}
	g.progress = progress
	if err := g.switchWorld(original.WorldAngkor); err != nil {
		g.message = err.Error()
		return
	}
	g.startMenuHasProgress = true
	g.startMenuConfirmNew = false
	g.loadStage(angkorTutorialStage)
	g.mode = gameModeStage
	g.message = "New game"
}

func (g *Game) drawStartMenu(screen *ebiten.Image) {
	if g.startMenuConfirmNew {
		g.drawNewGameConfirmation(screen)
		return
	}
	screen.Fill(color.Black)
	if g.splashBackground != nil {
		screen.DrawImage(g.splashBackground, nil)
	}
	if g.splashLogo != nil {
		screen.DrawImage(g.splashLogo, nil)
	}

	items := []string{"New game"}
	selected := 0
	if g.startMenuHasProgress {
		items = []string{"Continue", "New game"}
		selected = g.startMenuSelection
	}
	g.drawStartMenuRows(screen, items, selected, original.ScreenHeight-startMenuRowHeight+2)
}

func (g *Game) drawNewGameConfirmation(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if g.fontSmall != nil {
		lines := []string{
			"Starting a new game",
			"deletes your progress.",
			"Are you sure?",
		}
		for index, line := range lines {
			g.fontSmall.drawText(screen, line, original.ScreenWidth/2, 123+index*15, true, color.White)
		}
	}
	g.drawStartMenuRows(screen, []string{"No", "Yes"}, g.startMenuConfirmChoice, 190)
}

func (g *Game) drawStartMenuRows(screen *ebiten.Image, items []string, selected, bottom int) {
	if len(items) == 0 || g.fontMedium == nil {
		return
	}
	top := bottom - len(items)*startMenuRowHeight - 1
	for y := top; y < bottom; y++ {
		if g.softkeys != nil {
			g.softkeys.drawFrame(screen, 4, y&1, y, 0)
		} else {
			drawRect(screen, 0, y, original.ScreenWidth, 1, color.RGBA{0x0c, 0x2f, 0x39, 0xff})
		}
	}
	drawRect(screen, 0, top-2, original.ScreenWidth, 1, color.Black)
	drawRect(screen, 0, top-1, original.ScreenWidth, 1, color.White)
	drawRect(screen, 0, bottom, original.ScreenWidth, 1, color.White)
	drawRect(screen, 0, bottom+1, original.ScreenWidth, 1, color.Black)

	selected = clamp(selected, 0, len(items)-1)
	rowTop := top + selected*startMenuRowHeight
	drawRect(screen, 0, rowTop, original.ScreenWidth, startMenuRowHeight, color.RGBA{0xce, 0x9b, 0x00, 0xff})
	for index, item := range items {
		itemTop := top + index*startMenuRowHeight
		g.fontMedium.drawText(screen, item, original.ScreenWidth/2, itemTop+12, true, color.White)
		if index != selected || g.softkeys == nil {
			continue
		}
		width := g.fontMedium.stringWidth(item)
		centerY := itemTop + startMenuRowHeight/2
		g.softkeys.drawFrame(screen, 2, original.ScreenWidth/2-width/2-8, centerY, 0)
		g.softkeys.drawFrame(screen, 2, original.ScreenWidth/2+width/2+8, centerY, 0)
	}
}
