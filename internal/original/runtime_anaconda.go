package original

const (
	AnacondaPhaseDormant = iota
	AnacondaPhaseEmerge
	AnacondaPhaseVulnerable
	AnacondaPhaseHurt
	AnacondaPhaseRetract
	AnacondaPhaseRespawn
	AnacondaPhaseDelay
	AnacondaPhaseDefeated
	AnacondaPhaseComplete
	AnacondaPhaseDescend
	AnacondaPhaseTailCharge
	AnacondaPhaseTailStrike
)

const (
	anacondaLeftShaftX  = 12
	anacondaRightShaftX = 15
	anacondaSpawnY      = 2
	anacondaShaftY      = 5
	anacondaStrikeTopY  = 7
	anacondaStrikeY     = 8
	anacondaTailY       = 4
)

// GreatAnaconda mirrors the source fields used by Angkor stage index 8:
// aoInt, apInt, aqInt, arInt, LInt, MInt/NInt and both sprite animators.
type GreatAnaconda struct {
	Enabled            bool
	Phase              int
	PhaseTicks         int
	Health             int
	Column             int
	PreviousPhase      int
	CycleTicks         int
	Animation          int
	AnimationTicks     int
	TailAnimation      int
	TailAnimationTicks int
	TailVisible        bool
	TailHit            bool
	RumbleTicks        int
	BodyY              int
	Blocker            Point
	BlockerSet         bool
	Defeated           bool
	SealCollected      bool
	SealTicks          int
	StageComplete      bool
}

func newGreatAnaconda() GreatAnaconda {
	return GreatAnaconda{
		Enabled:       true,
		Health:        3,
		Animation:     0,
		TailAnimation: 0,
		BodyY:         1256,
	}
}

func (boss GreatAnaconda) X() int {
	return 10 + boss.Column*(2+boolInt(boss.Column > 0))
}

func (rt *Runtime) initGreatAnacondaStage() {
	if !rt.Anaconda.Enabled {
		return
	}
	// Doors with a raw-17 marker above or below are initialized in source
	// phase 3 (fully open). A marker above a door is consumed at load time.
	for y := 0; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 7 {
				continue
			}
			markerAbove := y > 0 && rt.Foreground[rt.index(x, y-1)] == 17
			markerBelow := y+1 < rt.Height() && rt.Foreground[rt.index(x, y+1)] == 17
			if !markerAbove && !markerBelow {
				continue
			}
			rt.Background[idx] = 0x30
			if markerAbove {
				rt.Foreground[rt.index(x, y-1)] = EmptyRawID
			}
		}
	}
}

func (rt *Runtime) tickGreatAnaconda(sourceTick int) (hits int, defeated bool) {
	boss := &rt.Anaconda
	if !boss.Enabled {
		return 0, false
	}
	if boss.RumbleTicks > 0 {
		boss.RumbleTicks--
	}
	// UVoid advances the shared gen1 animator once even while the tail is
	// hidden. The active boss branch advances it a second time below.
	boss.TailAnimationTicks++
	if boss.SealCollected && !boss.StageComplete {
		boss.SealTicks++
		if boss.SealTicks > angkorSealCompletionTicks {
			boss.StageComplete = true
		}
	}

	boss.PhaseTicks++
	rt.restoreAnacondaCeiling()
	nextAnimation := -1
	tailAnimationSet := false

	switch boss.Phase {
	case AnacondaPhaseDormant:
		if rt.Player.X >= 10 {
			boss.Phase = AnacondaPhaseDelay
			boss.PhaseTicks = 0
		}
	case AnacondaPhaseDelay:
		if boss.PhaseTicks > 10 {
			boss.Phase = AnacondaPhaseEmerge
			boss.PhaseTicks = 0
			nextAnimation = 2
		}
	case AnacondaPhaseEmerge:
		if boss.PhaseTicks > 40 {
			boss.Phase = AnacondaPhaseVulnerable
			boss.PhaseTicks = 0
			break
		}
		if boss.PhaseTicks > 20 {
			rt.clearAnacondaStrikeBoulders()
			rt.placeAnacondaBlocker(boss.X(), anacondaStrikeY)
		}
	case AnacondaPhaseVulnerable:
		if rt.clearAnacondaStrikeBoulders() {
			boss.Health--
			boss.PreviousPhase = AnacondaPhaseVulnerable
			rt.clearAnacondaBlocker()
			boss.Phase = AnacondaPhaseHurt
			nextAnimation = 3
			rt.emitSound(SoundEnemyHit)
			hits = 1
		}
		if boss.PhaseTicks > 15 && boss.Animation != 6 {
			nextAnimation = 6
		}
		if boss.PhaseTicks > 30 {
			boss.Phase = AnacondaPhaseRetract
			boss.PhaseTicks = 0
			nextAnimation = 0
			rt.clearAnacondaBlocker()
		}
	case AnacondaPhaseHurt:
		if boss.PhaseTicks > 40 {
			switch {
			case boss.Health <= 0:
				boss.Phase = AnacondaPhaseDefeated
				boss.PhaseTicks = 0
				boss.Defeated = true
				defeated = true
			case boss.PreviousPhase == AnacondaPhaseTailCharge:
				boss.Phase = AnacondaPhaseDescend
				boss.PhaseTicks = 0
			case boss.PreviousPhase == AnacondaPhaseVulnerable || boss.PreviousPhase == AnacondaPhaseEmerge:
				boss.Phase = AnacondaPhaseRetract
				boss.PhaseTicks = 0
				nextAnimation = 2
			}
		}
	case AnacondaPhaseRetract:
		threshold := 10
		if boss.Health <= 1 {
			threshold = 5
		}
		if boss.PhaseTicks >= threshold {
			boss.Phase = AnacondaPhaseRespawn
			boss.PhaseTicks = 0
			nextAnimation = 4
		} else if boss.PhaseTicks > threshold/2 && boss.Animation != 1 {
			nextAnimation = 1
		}
		rt.clearAnacondaBodyBoulders()
		rt.hurtFromAnacondaBody()
	case AnacondaPhaseRespawn:
		frameOffsetY := anacondaAnimationOffsetY(boss.Animation, boss.AnimationTicks)
		if boss.BodyY-40-frameOffsetY <= 112 {
			boss.CycleTicks = 0
			rt.respawnAnacondaBoulders()
			boss.Phase = AnacondaPhaseTailCharge
			rt.placeAnacondaBlocker(boss.X(), anacondaTailY)
		}
		rt.clearAnacondaBodyBoulders()
		rt.hurtFromAnacondaBody()
	case AnacondaPhaseTailCharge:
		boss.PhaseTicks--
		rt.hurtFromAnacondaBody()
		boss.CycleTicks++
		if boss.CycleTicks == 28 {
			nextAnimation = 7
		}
		if boss.CycleTicks >= 50 {
			boss.CycleTicks = 0
			boss.Phase = AnacondaPhaseTailStrike
			rt.clearAnacondaBlocker()
			nextAnimation = 8
			if boss.TailAnimation != 2 {
				boss.TailAnimation = 2
				boss.TailAnimationTicks = 0
			}
			tailAnimationSet = true
			boss.TailVisible = true
			boss.TailHit = false
		}
	case AnacondaPhaseTailStrike:
		boss.PhaseTicks--
		boss.CycleTicks++
		if boss.CycleTicks >= 12 {
			boss.CycleTicks = 0
			boss.Phase = AnacondaPhaseDescend
			nextAnimation = 4
			boss.TailHit = false
			boss.TailVisible = false
		} else if !boss.TailHit && rt.Player.Y == anacondaTailY && rt.Player.X >= boss.X()-3 && rt.Player.X <= boss.X()+4 {
			rt.HurtFromDirection(1, anacondaBodyKnockbackDirection(rt.Player.X, boss.X()))
			boss.TailHit = true
		}
	case AnacondaPhaseDescend:
		boss.PhaseTicks -= 2
		rt.hurtFromAnacondaBody()
		frameOffsetY := anacondaAnimationOffsetY(boss.Animation, boss.AnimationTicks)
		if boss.BodyY-40-frameOffsetY >= 280 {
			boss.Phase = AnacondaPhaseDelay
			boss.PhaseTicks = 0
			n := rt.Player.X - 10
			column := n / 3
			if n == column*3+2 {
				column += sourceTick % 50 / 25
			}
			boss.Column = clampRuntime(column, 0, 2)
		}
	case AnacondaPhaseDefeated:
		rt.clearAnacondaBlocker()
		if boss.PhaseTicks&0x6f == 1 {
			rt.emitSound(SoundBossDeath)
		}
		if boss.PhaseTicks > 80 {
			boss.Phase = AnacondaPhaseComplete
			rt.completeAnacondaGateGroup()
		}
	}

	if nextAnimation < 0 {
		boss.AnimationTicks++
	} else if nextAnimation != boss.Animation {
		boss.Animation = nextAnimation
		boss.AnimationTicks = 0
	}
	if boss.TailVisible && !tailAnimationSet {
		boss.TailAnimationTicks++
	}
	boss.updateBodyY()
	return hits, defeated
}

func (boss *GreatAnaconda) updateBodyY() {
	switch boss.Phase {
	case AnacondaPhaseEmerge:
		boss.BodyY = 256 - boss.PhaseTicks
	case AnacondaPhaseVulnerable, AnacondaPhaseHurt, AnacondaPhaseDefeated:
		boss.BodyY = 216
	case AnacondaPhaseRetract:
		boss.BodyY = 216 + boss.PhaseTicks*4
	case AnacondaPhaseRespawn:
		previous := boss.BodyY
		boss.BodyY = 241
		if (boss.BodyY-40-anacondaAnimationOffsetY(boss.Animation, boss.AnimationTicks))/TileSize == 3 {
			boss.BodyY = previous
		}
	case AnacondaPhaseDescend, AnacondaPhaseTailCharge, AnacondaPhaseTailStrike:
		boss.BodyY = 256 - (15 + boss.PhaseTicks*18)
	default:
		boss.BodyY = 1256
	}
}

func (rt *Runtime) restoreAnacondaCeiling() {
	for _, x := range [...]int{anacondaLeftShaftX, anacondaRightShaftX} {
		idx := rt.index(x, anacondaSpawnY)
		if rt.PlayerLayer[idx] == EmptyRawID {
			rt.PlayerLayer[idx] = 31
			rt.ObjectState[idx] = 0
			rt.ObjectMotion[idx] = ObjectMotion{}
		}
	}
}

func (rt *Runtime) clearAnacondaStrikeBoulders() bool {
	removed := false
	x := rt.Anaconda.X()
	for column := x; column <= x+1; column++ {
		for y := anacondaStrikeY; y >= anacondaStrikeTopY; y-- {
			if rt.PlayerLayer[rt.index(column, y)] != 0 {
				continue
			}
			rt.clearAnacondaObjectAt(column, y)
			removed = true
		}
	}
	return removed
}

func (rt *Runtime) clearAnacondaBodyBoulders() bool {
	top := (rt.Anaconda.BodyY - 40 - anacondaAnimationOffsetY(rt.Anaconda.Animation, rt.Anaconda.AnimationTicks)) / TileSize
	top = clampRuntime(top, 0, 10)
	removed := false
	for x := rt.Anaconda.X(); x <= rt.Anaconda.X()+1; x++ {
		for y := top; y <= 10; y++ {
			if rt.PlayerLayer[rt.index(x, y)] == 0 {
				rt.clearAnacondaObjectAt(x, y)
				removed = true
			}
		}
	}
	return removed
}

func (rt *Runtime) clearAnacondaObjectAt(x, y int) {
	idx := rt.index(x, y)
	if rt.PlayerLayer[idx] == 0 {
		rt.emitSound(SoundBreak)
	}
	rt.PlayerLayer[idx] = EmptyRawID
	rt.ObjectState[idx] = 0
	rt.ObjectMotion[idx] = ObjectMotion{}
	rt.FrozenOriginal[idx] = EmptyRawID
}

func (rt *Runtime) respawnAnacondaBoulders() {
	rt.Anaconda.RumbleTicks = 30
	for _, x := range [...]int{anacondaLeftShaftX, anacondaRightShaftX} {
		if rt.PlayerLayer[rt.index(x, anacondaShaftY)] != EmptyRawID {
			continue
		}
		idx := rt.index(x, anacondaSpawnY)
		rt.PlayerLayer[idx] = 0
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
	}
}

func (rt *Runtime) placeAnacondaBlocker(x, y int) {
	if rt.Anaconda.BlockerSet && rt.Anaconda.Blocker != (Point{X: x, Y: y}) {
		rt.clearAnacondaBlocker()
	}
	for column := x; column <= x+1; column++ {
		idx := rt.index(column, y)
		rt.PlayerLayer[idx] = 50
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
	}
	rt.Anaconda.Blocker = Point{X: x, Y: y}
	rt.Anaconda.BlockerSet = true
}

func (rt *Runtime) clearAnacondaBlocker() {
	if !rt.Anaconda.BlockerSet {
		return
	}
	point := rt.Anaconda.Blocker
	for x := point.X; x <= point.X+1; x++ {
		idx := rt.index(x, point.Y)
		if rt.PlayerLayer[idx] == 50 {
			rt.PlayerLayer[idx] = EmptyRawID
			rt.ObjectState[idx] = 0
			rt.ObjectMotion[idx] = ObjectMotion{}
		}
	}
	rt.Anaconda.Blocker = Point{}
	rt.Anaconda.BlockerSet = false
}

func (rt *Runtime) hurtFromAnacondaBody() bool {
	boss := &rt.Anaconda
	if rt.Player.X != boss.X() && rt.Player.X != boss.X()+1 {
		return false
	}
	frameOffsetY := anacondaAnimationOffsetY(boss.Animation, boss.AnimationTicks)
	top := boss.BodyY - 40 - frameOffsetY
	bottom := boss.BodyY + 256 - frameOffsetY
	playerY := rt.Player.Y*TileSize - rt.PlayerMotion.DY*rt.PlayerMotion.Remaining
	if playerY <= top || playerY >= bottom {
		return false
	}
	return rt.HurtFromDirection(1, anacondaBodyKnockbackDirection(rt.Player.X, boss.X()))
}

func (rt *Runtime) completeAnacondaGateGroup() {
	group := rt.ActiveEnemyGateGroup
	if group < 0 {
		return
	}
	if rt.EnemyGateCounters[group] > 0 {
		rt.EnemyGateCounters[group]--
	}
	if rt.EnemyGateCounters[group] == 0 {
		rt.openEnemyGateGroup(group)
	}
}

func (rt *Runtime) startEnemyGateDemo(triggerX, triggerY, group int) {
	behind := Point{X: triggerX - rt.PlayerMotion.DX, Y: triggerY}
	rt.closeDoorAt(behind.X, behind.Y)
	if !rt.enemyGateGroupValid(group) {
		return
	}
	target, ok := rt.enemyGateDemoTarget(group, behind)
	if !ok {
		return
	}
	rt.EnemyGateDemoActive = true
	rt.EnemyGateDemoPhase = 1
	rt.EnemyGateDemoTicks = 0
	rt.EnemyGateDemoOutboundTicks = rt.enemyGateDemoOutboundDuration(target)
	rt.EnemyGateDemoTarget = target
	rt.EnemyGateDemoTargetSet = true
}

func (rt *Runtime) enemyGateDemoOutboundDuration(target Point) int {
	targetX := clampRuntime(target.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	targetY := clampRuntime(target.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	startX := clampRuntime(rt.Player.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	startY := clampRuntime(rt.Player.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	if rt.viewportSet {
		startX = rt.viewportX
		startY = rt.viewportY
	}
	return max(1, max((absInt(targetX-startX)+7)/8, (absInt(targetY-startY)+7)/8))
}

func (rt *Runtime) enemyGateDemoTarget(group int, exclude Point) (Point, bool) {
	target := Point{}
	found := false
	for y := 1; y < rt.Height(); y++ {
		for x := 0; x < rt.Width(); x++ {
			idx := rt.index(x, y)
			if rt.Foreground[idx] != 17 || int(rt.Background[idx]) != group {
				continue
			}
			candidate := Point{X: x, Y: y}
			if rt.PlayerLayer[idx] != 18 {
				candidate.Y--
				switch rt.Foreground[rt.index(candidate.X, candidate.Y)] {
				case 7, 14, 33:
				default:
					continue
				}
			}
			if candidate == exclude {
				continue
			}
			target = candidate
			found = true
		}
	}
	return target, found
}

func (rt *Runtime) tickEnemyGateDemo() {
	if !rt.EnemyGateDemoActive {
		return
	}
	rt.EnemyGateDemoTicks++
	switch rt.EnemyGateDemoPhase {
	case 1:
		if rt.EnemyGateDemoTicks >= max(1, rt.EnemyGateDemoOutboundTicks) {
			rt.EnemyGateDemoPhase = 2
			rt.EnemyGateDemoTicks = 0
		}
	case 2:
		if rt.EnemyGateDemoTicks == 10 && rt.EnemyGateDemoTargetSet {
			rt.closeDoorAt(rt.EnemyGateDemoTarget.X, rt.EnemyGateDemoTarget.Y)
		}
		if rt.EnemyGateDemoTicks >= 40 {
			rt.EnemyGateDemoPhase = 3
			rt.EnemyGateDemoTicks = 0
			rt.EnemyGateMessageIndex = rt.EnemyGateMessages[rt.ActiveEnemyGateGroup]
			if rt.EnemyGateMessageIndex == 0 {
				rt.EnemyGateMessageIndex = 56
			}
			rt.EnemyGateMessageTicks = 80
		}
	case 3:
		if rt.EnemyGateDemoTicks >= rt.enemyGateDemoReturnTicks() {
			rt.EnemyGateDemoPhase = 4
			rt.EnemyGateDemoTicks = 0
		}
	case 4:
		if rt.EnemyGateDemoTicks >= 20 {
			rt.EnemyGateDemoActive = false
			rt.EnemyGateDemoPhase = 0
			rt.EnemyGateDemoTicks = 0
			rt.EnemyGateDemoOutboundTicks = 0
			rt.EnemyGateDemoTarget = Point{}
			rt.EnemyGateDemoTargetSet = false
		}
	}
}

func (rt *Runtime) EnemyGateDemoCamera() (x, y, phase, elapsed, duration int, ok bool) {
	if !rt.EnemyGateDemoActive || !rt.EnemyGateDemoTargetSet {
		return 0, 0, 0, 0, 0, false
	}
	gateX := clampRuntime(rt.EnemyGateDemoTarget.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	gateY := clampRuntime(rt.EnemyGateDemoTarget.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	playerX := clampRuntime(rt.Player.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	playerY := clampRuntime(rt.Player.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	phase = rt.EnemyGateDemoPhase
	switch phase {
	case 1:
		duration = max(1, rt.EnemyGateDemoOutboundTicks)
		return gateX, gateY, phase, clampRuntime(rt.EnemyGateDemoTicks, 0, duration), duration, true
	case 2:
		return gateX, gateY, phase, 1, 1, true
	case 3:
		duration = rt.enemyGateDemoReturnTicks()
		return playerX, playerY, phase, clampRuntime(rt.EnemyGateDemoTicks, 0, duration), duration, true
	case 4:
		return playerX, playerY, phase, 1, 1, true
	default:
		return 0, 0, 0, 0, 0, false
	}
}

func (rt *Runtime) EnemyGateMessage() (index, ticks int, ok bool) {
	if rt == nil || rt.EnemyGateMessageTicks <= 0 || rt.EnemyGateMessageIndex <= 0 {
		return 0, 0, false
	}
	return rt.EnemyGateMessageIndex, rt.EnemyGateMessageTicks, true
}

func (rt *Runtime) enemyGateDemoReturnTicks() int {
	if !rt.EnemyGateDemoTargetSet {
		return 1
	}
	gateX := clampRuntime(rt.EnemyGateDemoTarget.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	gateY := clampRuntime(rt.EnemyGateDemoTarget.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	playerX := clampRuntime(rt.Player.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	playerY := clampRuntime(rt.Player.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	return max(1, max((absInt(gateX-playerX)+4)/5, (absInt(gateY-playerY)+4)/5))
}

func (rt *Runtime) closeDoorAt(x, y int) bool {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return false
	}
	idx := rt.index(x, y)
	if rt.Foreground[idx] != 7 || int(rt.Background[idx])&0xf0 == 0 || rt.PlayerLayer[idx] == 32 {
		return false
	}
	rt.Background[idx] &= 0x0f
	rt.emitSound(SoundBoulder)
	rt.resolveClosedDoorOccupant(x, y)
	return true
}

func (rt *Runtime) playerCollisionDirection() int {
	switch {
	case rt.PlayerMotion.DX > 0:
		return 4
	case rt.PlayerMotion.DX < 0:
		return 2
	case rt.PlayerMotion.DY > 0:
		return 1
	case rt.PlayerMotion.DY < 0:
		return 3
	default:
		return 0
	}
}

func anacondaBodyKnockbackDirection(playerX, bossX int) int {
	if playerX == bossX {
		return 4
	}
	return 2
}

func anacondaAnimationOffsetY(animation, tick int) int {
	switch animation {
	case 4, 8:
		return 94
	case 7:
		sequence := animationSequenceForDurations(tick, [...]int{2, 2, 3, 15})
		if sequence == 3 {
			return 95
		}
		return 94
	default:
		return 0
	}
}

func animationSequenceForDurations(tick int, durations [4]int) int {
	total := 0
	for _, duration := range durations {
		total += duration
	}
	if total <= 0 {
		return 0
	}
	remaining := tick % total
	if remaining < 0 {
		remaining += total
	}
	for index, duration := range durations {
		if remaining < duration {
			return index
		}
		remaining -= duration
	}
	return len(durations) - 1
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
