package original

const (
	TeutonicKnightStateWalkLeft = iota
	TeutonicKnightStateWalkRight
	TeutonicKnightStateStrikeLeft
	TeutonicKnightStateStrikeRight
	TeutonicKnightStateHurtLeft
	TeutonicKnightStateHurtRight
	TeutonicKnightStateChargeLeft
	TeutonicKnightStateChargeRight
	TeutonicKnightStateSlashLeft
	TeutonicKnightStateSlashRight
	TeutonicKnightStateRecoverLeft
	TeutonicKnightStateRecoverRight
	TeutonicKnightStateDefeated
	TeutonicKnightStateDormant
	TeutonicKnightStateUnused
	TeutonicKnightStateComplete
)

const (
	teutonicKnightInitialX      = 408
	teutonicKnightY             = 504
	teutonicKnightLeftBound     = 360
	teutonicKnightRightBound    = 504
	teutonicKnightBoulderTopY   = 21
	teutonicKnightBoulderBottom = 22
	teutonicKnightSpawnY        = 16
	teutonicKnightRestY         = 18
	teutonicKnightDeathTicks    = 100
)

// EvilTeutonicKnight mirrors the Bavaria stage-9 source fields aoInt,
// aqInt, atInt, apInt, aeBoolean, adBoolean and dhInt.
type EvilTeutonicKnight struct {
	Enabled        bool
	State          int
	Health         int
	X              int
	Animation      int
	AnimationTicks int
	IntroActivated bool
	VerticalAttack bool
	AttackOffset   int
	DeathTicks     int
	RumbleTicks    int
	Defeated       bool
	Complete       bool
}

var teutonicKnightAnimationFrameTimes = [...][]int{
	{3, 3, 3, 3, 3, 3, 3, 3},
	{2, 2, 2, 3, 2, 2, 2, 3},
	{1, 2, 6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 1, 1, 1, 1, 1},
	{1, 2, 6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 1, 1, 1, 1, 1},
	{3},
	{3},
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 10, 2, 2, 2, 1, 1},
	{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 10, 2, 2, 2, 1, 1},
	{4, 3, 3, 3, 4, 3, 3, 3},
	{4, 3, 3, 3, 4, 3, 3, 3},
	{4, 4, 3, 3, 2, 2, 1, 1, 90},
	{6, 1, 1, 1, 2, 2, 2, 2, 2, 1, 1, 1, 1, 2, 3, 2, 3, 3, 2, 2, 2, 2, 2, 3, 3, 4, 4, 5, 6, 3, 2, 2, 2, 2, 2, 2, 2, 2},
}

func newEvilTeutonicKnight() EvilTeutonicKnight {
	return EvilTeutonicKnight{
		Enabled:   true,
		State:     TeutonicKnightStateDormant,
		Health:    4,
		X:         teutonicKnightInitialX,
		Animation: TeutonicKnightStateRecoverLeft,
	}
}

func (boss EvilTeutonicKnight) AnimationFrame() (frame, frameTick int) {
	return teutonicKnightAnimationFrame(boss.Animation, boss.AnimationTicks)
}

func (boss EvilTeutonicKnight) WorldY() int {
	return teutonicKnightY
}

func (boss EvilTeutonicKnight) DeathExplosionX(sourceTick int) int {
	return boss.X + sourceTick*boss.DeathTicks%48
}

func (rt *Runtime) tickSealTransition() {
	if !rt.SealCollected || rt.SealStageComplete {
		return
	}
	rt.SealTicks++
	if rt.SealTicks > sealCompletionTicks {
		rt.SealStageComplete = true
	}
}

func (rt *Runtime) tickEvilTeutonicKnight(sourceTick int) (hits int, defeated bool) {
	boss := &rt.TeutonicKnight
	if !boss.Enabled {
		return 0, false
	}
	rt.tickSealTransition()
	if boss.RumbleTicks > 0 {
		boss.RumbleTicks--
	}
	if boss.State == TeutonicKnightStateComplete {
		return 0, false
	}
	if boss.State == TeutonicKnightStateDefeated {
		rt.tickEvilTeutonicKnightDeath(sourceTick)
		return 0, false
	}

	playerCenterX := rt.Player.X*TileSize + TileSize/2
	playerY := rt.Player.Y * TileSize
	bossCenterX := boss.X + TileSize
	nextAnimation := -1
	if boss.State == TeutonicKnightStateDormant {
		nextAnimation = TeutonicKnightStateDormant
		if playerCenterX > bossCenterX {
			boss.IntroActivated = true
		}
	}
	if !boss.VerticalAttack && (boss.State == TeutonicKnightStateWalkLeft || boss.State == TeutonicKnightStateWalkRight) {
		if boss.State == TeutonicKnightStateWalkLeft {
			boss.AttackOffset = -36
		} else {
			boss.AttackOffset = 36
		}
		if playerY < teutonicKnightY && (playerCenterX == bossCenterX+boss.AttackOffset || sourceTick%76 == 0) {
			boss.VerticalAttack = true
		}
	}

	repeats := 1
	if boss.Health > 0 && sourceTick%boss.Health == 0 {
		repeats = 2
	}
	if (boss.State == TeutonicKnightStateChargeLeft || boss.State == TeutonicKnightStateChargeRight) && sourceTick&0xb == 0 {
		repeats = 2
	}

	for step := 0; step < repeats; step++ {
		attackHits := false
		switch boss.State {
		case TeutonicKnightStateDormant:
			if boss.animationEnded() {
				boss.State = TeutonicKnightStateWalkLeft
				nextAnimation = TeutonicKnightStateWalkLeft
				boss.IntroActivated = false
			}
			if boss.IntroActivated {
				boss.tickAnimation()
			}
		case TeutonicKnightStateHurtLeft, TeutonicKnightStateHurtRight:
			if boss.animationEnded() {
				if boss.State == TeutonicKnightStateHurtLeft {
					boss.State = TeutonicKnightStateWalkLeft
				} else {
					boss.State = TeutonicKnightStateWalkRight
				}
				nextAnimation = boss.State
			}
			boss.VerticalAttack = false
			boss.AttackOffset = 0
		case TeutonicKnightStateRecoverLeft, TeutonicKnightStateRecoverRight:
			if playerCenterX > bossCenterX && bossCenterX < teutonicKnightRightBound {
				boss.State = TeutonicKnightStateWalkRight
				nextAnimation = TeutonicKnightStateWalkRight
			} else if playerCenterX < bossCenterX && bossCenterX > teutonicKnightLeftBound {
				boss.State = TeutonicKnightStateWalkLeft
				nextAnimation = TeutonicKnightStateWalkLeft
			}
			boss.VerticalAttack = false
			boss.AttackOffset = 0
		case TeutonicKnightStateChargeLeft:
			if playerY >= teutonicKnightY {
				if playerCenterX >= bossCenterX-2*TileSize {
					boss.State = TeutonicKnightStateSlashLeft
					nextAnimation = TeutonicKnightStateSlashLeft
				} else {
					boss.X -= 2
				}
			} else if bossCenterX >= teutonicKnightLeftBound {
				boss.State = TeutonicKnightStateWalkLeft
				nextAnimation = TeutonicKnightStateWalkLeft
			}
		case TeutonicKnightStateChargeRight:
			if playerY >= teutonicKnightY {
				if playerCenterX <= bossCenterX+2*TileSize {
					boss.State = TeutonicKnightStateSlashRight
					nextAnimation = TeutonicKnightStateSlashRight
				} else {
					boss.X += 2
				}
			} else if bossCenterX <= teutonicKnightRightBound {
				boss.State = TeutonicKnightStateWalkRight
				nextAnimation = TeutonicKnightStateWalkRight
			}
		case TeutonicKnightStateWalkLeft:
			if playerY >= teutonicKnightY && bossCenterX > teutonicKnightLeftBound {
				if playerCenterX < bossCenterX {
					boss.State = TeutonicKnightStateChargeLeft
					nextAnimation = TeutonicKnightStateChargeLeft
				} else {
					boss.X--
				}
			} else if boss.VerticalAttack {
				boss.State = TeutonicKnightStateStrikeLeft
				nextAnimation = TeutonicKnightStateStrikeLeft
			} else if bossCenterX <= teutonicKnightLeftBound {
				boss.State = TeutonicKnightStateWalkRight
				nextAnimation = TeutonicKnightStateWalkRight
			} else {
				boss.X--
			}
		case TeutonicKnightStateWalkRight:
			if playerY >= teutonicKnightY && bossCenterX < teutonicKnightRightBound {
				if playerCenterX < bossCenterX {
					boss.X++
				} else {
					boss.State = TeutonicKnightStateChargeRight
					nextAnimation = TeutonicKnightStateChargeRight
				}
			} else if boss.VerticalAttack {
				boss.State = TeutonicKnightStateStrikeRight
				nextAnimation = TeutonicKnightStateStrikeRight
			} else if bossCenterX >= teutonicKnightRightBound {
				boss.State = TeutonicKnightStateWalkLeft
				nextAnimation = TeutonicKnightStateWalkLeft
			} else {
				boss.X++
			}
		case TeutonicKnightStateSlashLeft:
			frame, _ := boss.AnimationFrame()
			if boss.animationEnded() {
				boss.State = TeutonicKnightStateRecoverLeft
				nextAnimation = TeutonicKnightStateRecoverLeft
			}
			attackHits = frame >= 4 && playerY >= teutonicKnightY && playerCenterX >= bossCenterX-2*TileSize && playerCenterX <= bossCenterX
		case TeutonicKnightStateSlashRight:
			frame, _ := boss.AnimationFrame()
			if boss.animationEnded() {
				boss.State = TeutonicKnightStateRecoverRight
				nextAnimation = TeutonicKnightStateRecoverRight
			}
			attackHits = frame >= 4 && playerY >= teutonicKnightY && playerCenterX >= bossCenterX && playerCenterX <= bossCenterX+2*TileSize
		case TeutonicKnightStateStrikeLeft, TeutonicKnightStateStrikeRight:
			frame, frameTick := boss.AnimationFrame()
			if frame == 5 && frameTick == 0 {
				boss.RumbleTicks = 30
			}
			if boss.animationEnded() {
				if boss.State == TeutonicKnightStateStrikeLeft {
					boss.State = TeutonicKnightStateRecoverLeft
				} else {
					boss.State = TeutonicKnightStateRecoverRight
				}
				nextAnimation = boss.State
				boss.VerticalAttack = false
				boss.AttackOffset = 0
			}
			attackHits = frame >= 7 && playerY < teutonicKnightY && playerCenterX == bossCenterX+boss.AttackOffset
		}
		if attackHits {
			rt.HurtFromDirection(1, 0)
		}
		if playerY >= teutonicKnightY && playerCenterX == bossCenterX-TileSize {
			rt.HurtFromDirection(1, rt.playerCollisionDirection())
		}
	}

	frame, _ := boss.AnimationFrame()
	if (boss.State == TeutonicKnightStateSlashLeft || boss.State == TeutonicKnightStateSlashRight) && frame == 5 {
		rt.spawnTeutonicKnightBoulders()
	}
	if boss.State != TeutonicKnightStateChargeLeft && boss.State != TeutonicKnightStateChargeRight {
		hits = rt.consumeTeutonicKnightBoulders(bossCenterX, &nextAnimation)
	}
	if boss.Health <= 0 {
		boss.State = TeutonicKnightStateDefeated
		boss.DeathTicks = 0
		boss.Defeated = true
		nextAnimation = TeutonicKnightStateDefeated
		defeated = true
	}
	if nextAnimation >= 0 {
		boss.setAnimation(nextAnimation)
	} else {
		boss.tickAnimation()
	}
	return hits, defeated
}

func (rt *Runtime) consumeTeutonicKnightBoulders(bossCenterX int, nextAnimation *int) int {
	boss := &rt.TeutonicKnight
	centerTile := bossCenterX / TileSize
	hits := 0
	for y := teutonicKnightBoulderTopY; y <= teutonicKnightBoulderBottom; y++ {
		for x := centerTile - 1; x <= centerTile+1; x++ {
			if x < 0 || x >= rt.Width() || y < 0 || y >= rt.Height() {
				continue
			}
			idx := rt.index(x, y)
			if rt.PlayerLayer[idx] != 0 {
				continue
			}
			falling := rt.ObjectState[idx]&objectDirectionMask == 3
			if falling && boss.State != TeutonicKnightStateDormant {
				boss.Health--
				hits++
				switch boss.State {
				case TeutonicKnightStateWalkLeft, TeutonicKnightStateStrikeLeft, TeutonicKnightStateHurtLeft,
					TeutonicKnightStateSlashLeft, TeutonicKnightStateRecoverLeft:
					boss.State = TeutonicKnightStateHurtLeft
					*nextAnimation = TeutonicKnightStateHurtLeft
				case TeutonicKnightStateWalkRight, TeutonicKnightStateStrikeRight, TeutonicKnightStateHurtRight,
					TeutonicKnightStateSlashRight, TeutonicKnightStateRecoverRight:
					boss.State = TeutonicKnightStateHurtRight
					*nextAnimation = TeutonicKnightStateHurtRight
				}
			}
			rt.PlayerLayer[idx] = 30
			rt.ObjectState[idx] = 4
			rt.ObjectMotion[idx] = ObjectMotion{}
			rt.FrozenOriginal[idx] = EmptyRawID
			rt.EnemyGateGroup[idx] = -1
			rt.emitSound(SoundBoulder)
		}
	}
	return hits
}

func (rt *Runtime) spawnTeutonicKnightBoulders() {
	for _, x := range [...]int{16, 19} {
		if x < 0 || x >= rt.Width() || teutonicKnightRestY >= rt.Height() || teutonicKnightSpawnY >= rt.Height() {
			continue
		}
		if rt.PlayerLayer[rt.index(x, teutonicKnightRestY)] != EmptyRawID {
			continue
		}
		idx := rt.index(x, teutonicKnightSpawnY)
		rt.PlayerLayer[idx] = 0
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
		rt.FrozenOriginal[idx] = EmptyRawID
		rt.EnemyGateGroup[idx] = -1
	}
}

func (rt *Runtime) tickEvilTeutonicKnightDeath(sourceTick int) {
	boss := &rt.TeutonicKnight
	if boss.DeathTicks > teutonicKnightDeathTicks {
		boss.DeathTicks++
		boss.State = TeutonicKnightStateComplete
		boss.Complete = true
		rt.completeActiveEnemyGateGroup()
	} else {
		boss.DeathTicks++
		rt.emitSound(SoundBossDeath)
		explosionX := boss.DeathExplosionX(sourceTick) / TileSize
		explosionY := (teutonicKnightY + TileSize) / TileSize
		if absInt(rt.Player.X-explosionX) <= 1 && absInt(rt.Player.Y-explosionY) <= 1 {
			rt.HurtFromDirection(1, 0)
		}
	}
	boss.tickAnimation()
}

func (boss *EvilTeutonicKnight) setAnimation(animation int) {
	if animation == boss.Animation {
		return
	}
	boss.Animation = animation
	boss.AnimationTicks = 0
}

func (boss *EvilTeutonicKnight) tickAnimation() {
	duration := teutonicKnightAnimationDuration(boss.Animation)
	if duration <= 0 {
		return
	}
	boss.AnimationTicks++
	if boss.AnimationTicks >= duration {
		boss.AnimationTicks = 0
	}
}

func (boss EvilTeutonicKnight) animationEnded() bool {
	duration := teutonicKnightAnimationDuration(boss.Animation)
	return duration > 0 && boss.AnimationTicks == duration-1
}

func teutonicKnightAnimationDuration(animation int) int {
	if animation < 0 || animation >= len(teutonicKnightAnimationFrameTimes) {
		return 0
	}
	duration := 0
	for _, frameDuration := range teutonicKnightAnimationFrameTimes[animation] {
		duration += max(1, frameDuration)
	}
	return duration
}

func teutonicKnightAnimationFrame(animation, tick int) (frame, frameTick int) {
	if animation < 0 || animation >= len(teutonicKnightAnimationFrameTimes) {
		return 0, 0
	}
	times := teutonicKnightAnimationFrameTimes[animation]
	duration := teutonicKnightAnimationDuration(animation)
	if len(times) == 0 || duration <= 0 {
		return 0, 0
	}
	tick %= duration
	if tick < 0 {
		tick += duration
	}
	for frame, frameDuration := range times {
		frameDuration = max(1, frameDuration)
		if tick < frameDuration {
			return frame, tick
		}
		tick -= frameDuration
	}
	return len(times) - 1, 0
}
