package originalgame

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wangle201210/zskc/internal/original"
)

const (
	bavariaExplosionFrames      = "decoded/sprites/gen0/chunk03-frames.png"
	bavariaExplosionModules     = "decoded/sprites/gen0/chunk03-modules.png"
	bavariaExplosionMetadata    = "decoded/sprites/gen0/chunk03-animations.json"
	bavariaExplosiveFrames      = "decoded/sprites/gen0/chunk05-frames.png"
	bavariaExplosiveModules     = "decoded/sprites/gen0/chunk05-modules.png"
	bavariaExplosiveMetadata    = "decoded/sprites/gen0/chunk05-animations.json"
	bavariaCrawlerTrapFrames    = "decoded/sprites/gen0/chunk08-frames.png"
	bavariaCrawlerTrapModules   = "decoded/sprites/gen0/chunk08-modules.png"
	bavariaCrawlerTrapMetadata  = "decoded/sprites/gen0/chunk08-animations.json"
	bavariaSpikeFrames          = "decoded/sprites/gen1/chunk01-frames.png"
	bavariaSpikeModules         = "decoded/sprites/gen1/chunk01-modules.png"
	bavariaSpikeMetadata        = "decoded/sprites/gen1/chunk01-animations.json"
	bavariaMovingHazardModules  = "decoded/sprites/gen1/chunk02-modules.png"
	bavariaMovingHazardMetadata = "decoded/sprites/gen1/chunk02-animations.json"
	bavariaSpearFrames          = "decoded/sprites/gen1/chunk03-frames.png"
	bavariaSpearModules         = "decoded/sprites/gen1/chunk03-modules.png"
	bavariaSpearMetadata        = "decoded/sprites/gen1/chunk03-animations.json"
	bavariaWindPodFrames        = "decoded/sprites/gen2/chunk03-frames.png"
	bavariaWindPodModules       = "decoded/sprites/gen2/chunk03-modules.png"
	bavariaWindPodMetadata      = "decoded/sprites/gen2/chunk03-animations.json"
	bavariaFanPotRedFrames      = "decoded/sprites/gen2/chunk04-frames.png"
	bavariaFanPotRedModules     = "decoded/sprites/gen2/chunk04-modules.png"
	bavariaFanPotBlueFrames     = "decoded/sprites/gen2/chunk04-palette01-frames.png"
	bavariaFanPotBlueModules    = "decoded/sprites/gen2/chunk04-palette01-modules.png"
	bavariaFanPotMetadata       = "decoded/sprites/gen2/chunk04-animations.json"
	bavariaBlastWallModules     = "decoded/sprites/gen2/chunk05-modules.png"
	bavariaBlastWallMetadata    = "decoded/sprites/gen2/chunk05-animations.json"
	bavariaWaterFrames          = "decoded/sprites/gen2/chunk06-frames.png"
	bavariaWaterModules         = "decoded/sprites/gen2/chunk06-modules.png"
	bavariaWaterMetadata        = "decoded/sprites/gen2/chunk06-animations.json"
	bavariaWaterPotionModules   = "decoded/sprites/gen2/chunk07-modules.png"
	bavariaWaterPotionMetadata  = "decoded/sprites/gen2/chunk07-animations.json"
	bavariaFanSwitchFrames      = "decoded/sprites/gen3/chunk09-frames.png"
	bavariaFanSwitchModules     = "decoded/sprites/gen3/chunk09-modules.png"
	bavariaFanSwitchMetadata    = "decoded/sprites/gen3/chunk09-animations.json"
	bavariaWindColumnFrames     = "decoded/sprites/gen3/chunk04-frames.png"
	bavariaWindColumnModules    = "decoded/sprites/gen3/chunk04-modules.png"
	bavariaWindColumnMetadata   = "decoded/sprites/gen3/chunk04-animations.json"
	bavariaKnightFrames         = "decoded/sprites/b1/chunk00-frames.png"
	bavariaKnightModules        = "decoded/sprites/b1/chunk00-modules.png"
	bavariaKnightMetadata       = "decoded/sprites/b1/chunk00-animations.json"
)

type bavariaSpriteSet struct {
	explosion    *spriteSheet
	explosive    *spriteSheet
	crawlerTrap  *spriteSheet
	spike        *spriteSheet
	movingHazard *spriteSheet
	spear        *spriteSheet
	windPod      *spriteSheet
	fanPotRed    *spriteSheet
	fanPotBlue   *spriteSheet
	blastWall    *spriteSheet
	water        *spriteSheet
	waterPotion  *spriteSheet
	fanSwitch    *spriteSheet
	windColumn   *spriteSheet
	knight       *spriteSheet
}

func loadBavariaSpriteSet() (bavariaSpriteSet, error) {
	var result bavariaSpriteSet
	var err error
	load := func(name string, target **spriteSheet, frames, modules, metadata string) error {
		*target, err = loadSpriteSheetWithModules(frames, modules, metadata)
		if err != nil {
			return fmt.Errorf("load Bavaria %s: %w", name, err)
		}
		return nil
	}
	for _, asset := range []struct {
		name                      string
		target                    **spriteSheet
		frames, modules, metadata string
	}{
		{"explosion", &result.explosion, bavariaExplosionFrames, bavariaExplosionModules, bavariaExplosionMetadata},
		{"explosive boulder", &result.explosive, bavariaExplosiveFrames, bavariaExplosiveModules, bavariaExplosiveMetadata},
		{"crawler trap", &result.crawlerTrap, bavariaCrawlerTrapFrames, bavariaCrawlerTrapModules, bavariaCrawlerTrapMetadata},
		{"spikes", &result.spike, bavariaSpikeFrames, bavariaSpikeModules, bavariaSpikeMetadata},
		{"moving hazard", &result.movingHazard, "", bavariaMovingHazardModules, bavariaMovingHazardMetadata},
		{"spear", &result.spear, bavariaSpearFrames, bavariaSpearModules, bavariaSpearMetadata},
		{"wind pod", &result.windPod, bavariaWindPodFrames, bavariaWindPodModules, bavariaWindPodMetadata},
		{"red fan pot", &result.fanPotRed, bavariaFanPotRedFrames, bavariaFanPotRedModules, bavariaFanPotMetadata},
		{"blue fan pot", &result.fanPotBlue, bavariaFanPotBlueFrames, bavariaFanPotBlueModules, bavariaFanPotMetadata},
		{"blast wall", &result.blastWall, "", bavariaBlastWallModules, bavariaBlastWallMetadata},
		{"water", &result.water, bavariaWaterFrames, bavariaWaterModules, bavariaWaterMetadata},
		{"water potion", &result.waterPotion, "", bavariaWaterPotionModules, bavariaWaterPotionMetadata},
		{"fan switch", &result.fanSwitch, bavariaFanSwitchFrames, bavariaFanSwitchModules, bavariaFanSwitchMetadata},
		{"wind column", &result.windColumn, bavariaWindColumnFrames, bavariaWindColumnModules, bavariaWindColumnMetadata},
		{"Evil Teutonic Knight", &result.knight, bavariaKnightFrames, bavariaKnightModules, bavariaKnightMetadata},
	} {
		if err := load(asset.name, asset.target, asset.frames, asset.modules, asset.metadata); err != nil {
			return bavariaSpriteSet{}, err
		}
	}
	return result, nil
}

func (g *Game) drawEvilTeutonicKnight(dst *ebiten.Image, camX, camY int) {
	if g.rt == nil || !g.rt.TeutonicKnight.Enabled || g.bavaria.knight == nil {
		return
	}
	boss := g.rt.TeutonicKnight
	if boss.State == original.TeutonicKnightStateComplete {
		return
	}
	g.bavaria.knight.drawAnimationWithFrameOffset(dst, boss.Animation, boss.AnimationTicks, boss.X-camX, boss.WorldY()-camY, 0)
	if boss.State == original.TeutonicKnightStateDefeated && g.bavaria.explosion != nil {
		g.bavaria.explosion.drawAnimationSequenceFrame(
			dst,
			0,
			g.tick,
			boss.DeathExplosionX(g.tick)-camX,
			boss.WorldY()+original.TileSize-camY,
			0,
		)
	}
}

func (g *Game) drawEvilTeutonicKnightHealth(screen *ebiten.Image) {
	if g.rt == nil || !g.rt.TeutonicKnight.Enabled || g.rt.TeutonicKnight.Health <= 0 {
		return
	}
	boss := g.rt.TeutonicKnight
	if boss.State == original.TeutonicKnightStateDormant || boss.State == original.TeutonicKnightStateDefeated || boss.State == original.TeutonicKnightStateComplete {
		return
	}
	const width = 4*14 + 2
	x := (original.ScreenWidth - width) / 2
	drawRect(screen, x, 5, width, 12, color.Black)
	for segment := 0; segment < boss.Health; segment++ {
		drawRect(screen, x+2+segment*14, 7, 12, 8, color.RGBA{0x3b, 0xb7, 0x8f, 0xff})
	}
}

func (g *Game) drawBavariaForeground(dst *ebiten.Image, id original.RawID, px, py int) bool {
	if g.worldIndex != original.WorldBavaria {
		return false
	}
	switch id {
	case 15:
		if g.rt.FanPhase > 0 && g.rt.FanPhase <= 5 {
			frame := clamp(g.rt.FanPhase*5/10, 0, 4)
			g.bavaria.fanPotBlue.drawFrame(dst, frame, px, py, 0)
		}
	case 16:
		if g.rt.FanPhase >= 5 && g.rt.FanPhase < 9 {
			frame := clamp(4-g.rt.FanPhase*5/10, 0, 4)
			g.bavaria.fanPotRed.drawFrame(dst, frame, px, py, 0)
		}
	case 34:
		g.bavaria.windColumn.drawAnimationSequenceFrame(dst, 2, 0, px, py, 0)
	case 35:
		g.bavaria.windPod.drawAnimation(dst, 1, g.tick, px, py, 0)
	case 37:
		g.bavaria.windColumn.drawAnimationSequenceFrame(dst, 2, 0, px, py, 0)
		g.bavaria.windPod.drawAnimation(dst, 1, g.tick, px, py, 0)
	default:
		return false
	}
	return true
}

func (g *Game) drawBavariaObject(dst *ebiten.Image, id original.RawID, x, y, px, py int) bool {
	if g.worldIndex != original.WorldBavaria {
		return false
	}
	state := g.objectStateAt(x, y)
	switch id {
	case 8:
		g.bavaria.explosive.drawFrame(dst, (g.tick>>1)&1, px, py, 0)
	case 14:
		g.drawBavariaMovingHazard(dst, x, y, state, px, py)
	case 16:
		below, _ := g.rt.At(original.PlayerLayer, x, y+1)
		if below == 16 {
			return true
		}
		animation := 0
		if state&7 == 4 {
			animation = 1
		}
		timer := (state >> 8) & 0xff
		elapsed := 0
		if timer > 0 {
			elapsed = 36 - timer
		}
		frame, ok := g.bavaria.spear.animationFrame(animation, elapsed)
		if ok {
			// Java raw 16 copies only animationFrames[x] into bNInt before
			// drawAnimationFrame; the animation-frame y offset is not applied.
			g.bavaria.spear.drawFrame(dst, frame.Frame, px+frame.X, py, frame.Flags)
		}
	case 18:
		frame := 0
		switch {
		case g.rt.FanPhase == 9:
			frame = 2
		case g.rt.FanDirection < 0:
			frame = 1
		case g.rt.FanDirection > 0:
			frame = 3
		}
		g.bavaria.fanSwitch.drawFrame(dst, frame, px, py, 0)
	case 28:
		g.drawBavariaSpike(dst, x, y, state, px, py)
	case 34:
		if g.rt.FanPhase >= 5 && g.rt.FanPhase < 9 {
			g.bavaria.fanPotBlue.drawFrame(dst, clamp(g.rt.FanPhase*5/10, 0, 4), px, py, 0)
		}
	case 35:
		if g.rt.FanPhase > 0 && g.rt.FanPhase <= 5 {
			g.bavaria.fanPotRed.drawFrame(dst, clamp(4-g.rt.FanPhase*5/10, 0, 4), px, py, 0)
		}
	case 36:
		animation := 0
		if state == 1 {
			animation = 1
		}
		g.bavaria.crawlerTrap.drawAnimationSequenceFrame(dst, animation, (g.tick >> 1), px, py, 0)
	case 37:
		g.bavaria.blastWall.drawModule(dst, clamp((state-1)*3/8, 0, 3), px, py)
	case 38, 40, 51, 52:
		// Water sources and chest payloads are represented by their foreground.
	case 47:
		g.bavaria.windPod.drawAnimation(dst, 0, g.tick, px, py, 0)
	case 54:
		g.bavaria.explosion.drawAnimation(dst, 0, state, px, py, 0)
	default:
		return false
	}
	return true
}

func (g *Game) drawBavariaMovingHazard(dst *ebiten.Image, x, y, state, px, py int) {
	reversed := state&8 != 0
	body := (g.tick >> 1) % 3
	bodyX, bodyY := px, py
	if reversed {
		body = 2 - body
		if state&7 != 3 && (g.tick>>1)&1 == 0 && x > 0 {
			if left, ok := g.rt.At(original.PlayerLayer, x-1, y); ok && left != original.EmptyRawID {
				bodyX--
				bodyY++
			}
		}
	}
	if state&7 != 3 {
		particle := (g.tick >> 1) % 5
		module := 3 + particle
		particleX := px - particle*4
		if reversed {
			module = 8 + particle
			particleX = px + 12 + particle*3
		}
		particleY := py + original.TileSize - g.bavaria.movingHazard.moduleHeight(module)
		g.bavaria.movingHazard.drawModule(dst, module, particleX, particleY)
	}
	g.bavaria.movingHazard.drawModule(dst, body, bodyX, bodyY)
}

func (g *Game) drawBavariaSpike(dst *ebiten.Image, x, y, state, px, py int) {
	extent := g.rt.SpikeExtentAt(x, y)
	if g.renderPhase > 0 {
		next := g.rt.SpikeExtentAtSourceTick(x, y, g.tick+1)
		extent += (next - extent) * g.renderPhase / renderStepsPerSource
	}
	segments := 1
	if extent > 0 {
		segments = (extent-1)/original.TileSize + 2
	}
	direction := -1
	frame := 3
	if state&7 == 3 {
		direction = 1
		frame = 0
	}
	for segment := 0; segment < segments; segment++ {
		// Source frames 0/3 are the purple crystal tip; 1/2 are the steel shaft.
		g.bavaria.spike.drawFrame(dst, frame+segment*direction, px+3, py+direction*(extent-segment*original.TileSize), 0)
	}
	coverOffset := sourceBavariaSpikeCoverOffset(direction)
	coverY := y + coverOffset
	if cover, ok := g.rt.At(original.PlayerLayer, x, coverY); ok && cover >= 80 {
		g.drawWorldFrame(dst, int(cover-80), px, py+coverOffset*original.TileSize)
	}
}

func sourceBavariaSpikeCoverOffset(direction int) int {
	return -direction
}

func (g *Game) drawBavariaWater(dst *ebiten.Image, camX, camY int) {
	if g.worldIndex != original.WorldBavaria || g.rt == nil || g.bavaria.water == nil {
		return
	}
	view := sourceStageCellViewport(camX, camY)
	for relY := view.firstRel; relY < view.lastRelY; relY++ {
		for relX := view.firstRel; relX < view.lastRelX; relX++ {
			x := view.firstX + relX
			y := view.firstY + relY
			px := view.offX + relX*original.TileSize
			py := view.offY + relY*original.TileSize
			for layerIndex := 0; layerIndex < 3; layerIndex++ {
				layer := g.rt.WaterRenderLayerAt(x, y, layerIndex)
				if !layer.Visible {
					continue
				}
				target := dst
				if layer.OffsetX != 0 {
					clip := image.Rect(px, py+layerIndex*8, px+original.TileSize, py+(layerIndex+1)*8).Intersect(dst.Bounds())
					if clip.Empty() {
						continue
					}
					target = dst.SubImage(clip).(*ebiten.Image)
				}
				if layer.Kind == 8 && layerIndex == 0 && g.rt.WaterAt(x, y-1) > 0 {
					g.bavaria.water.drawFrame(target, 33, px+layer.OffsetX, py, layer.Flags)
					break
				}
				if layer.Kind == 15 {
					g.bavaria.water.drawFrame(target, layer.Base+g.waterSpecialFrame, px+layer.OffsetX-8, py+layer.OffsetY+8, 36)
					g.waterSpecialFrame = (g.waterSpecialFrame + 1) % 3
					continue
				}
				animation := g.tick & 1
				if layer.Kind == 7 {
					animation = (g.tick >> 3) & 1
				}
				drawY := py + layer.OffsetY
				if layer.Kind == 14 || layer.Kind == 11 {
					drawY = py
				}
				g.bavaria.water.drawFrame(target, layer.Base+animation, px+layer.OffsetX, drawY, layer.Flags)
			}
		}
	}
}
