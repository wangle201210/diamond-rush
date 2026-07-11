package original

import (
	"fmt"
	"testing"
)

func TestRuntimeStage05TriggersFallingTorchesThroughSourceEvents(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage05.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	if !rt.IsFallingTorchStage() || rt.RisingFireHeight != fallingFireInitialHeight || rt.RisingFireAnimation != 2 {
		t.Fatalf("initial falling-torch state stage=%v height=%d animation=%d", rt.IsFallingTorchStage(), rt.RisingFireHeight, rt.RisingFireAnimation)
	}

	rt.SetForTest(PlayerLayer, 18, 63, 0)
	rt.ObjectMotion[rt.index(18, 63)] = ObjectMotion{}
	rt.TickSourceFrame(8, 1, 0)
	if rt.FallingTorchTriggers != 1 || rt.FallingTorchWarningTicks != fallingTorchWarningDuration {
		t.Fatalf("settled trigger rock count=%d warning=%d, want 1/%d", rt.FallingTorchTriggers, rt.FallingTorchWarningTicks, fallingTorchWarningDuration)
	}

	rt.Player = Point{X: 17, Y: 53}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetForTest(PlayerLayer, 17, 53, EmptyRawID)
	rt.SetForTest(PlayerLayer, 18, 53, EmptyRawID)
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter first raw1 trigger strip")
	}
	if rt.FallingTorchTriggers != 2 || rt.FallingTorchWarningTicks != fallingTorchWarningDuration {
		t.Fatalf("first strip count=%d warning=%d, want 2/%d", rt.FallingTorchTriggers, rt.FallingTorchWarningTicks, fallingTorchWarningDuration)
	}
	for x := 18; x <= 22; x++ {
		if id, _ := rt.At(ForegroundLayer, x, 53); id != EmptyRawID {
			t.Fatalf("first trigger strip x=%d remains raw%d", x, id)
		}
	}

	rt.Player = Point{X: 7, Y: 38}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetForTest(PlayerLayer, 7, 38, EmptyRawID)
	rt.SetForTest(PlayerLayer, 7, 37, EmptyRawID)
	rt.SetForTest(PlayerLayer, 7, 36, EmptyRawID)
	if !rt.TryMove(0, -1) {
		t.Fatal("failed to enter demo raw0 strip")
	}
	rt.PlayerMotion.Remaining = 6
	rt.TickSourceFrame(8, 2, 0)
	if !rt.ForegroundDemoActive || rt.ForegroundDemoID != 3 || rt.ForegroundEvents != 1 {
		t.Fatalf("demo active=%v id=%d events=%d, want true/3/1", rt.ForegroundDemoActive, rt.ForegroundDemoID, rt.ForegroundEvents)
	}
	rt.AdvancePlayerMotion()
	rt.TickSourceFrame(8, 3, 0)
	if rt.Player != (Point{X: 7, Y: 36}) || rt.FallingTorchTriggers != 3 {
		t.Fatalf("scripted step player=%+v triggers=%d, want (7,36)/3", rt.Player, rt.FallingTorchTriggers)
	}
	for x := 7; x <= 16; x++ {
		if id, _ := rt.At(ForegroundLayer, x, 36); id != EmptyRawID {
			t.Fatalf("second trigger strip x=%d remains raw%d", x, id)
		}
	}
	rt.AdvancePlayerMotion()
	rt.TickSourceFrame(8, 4, 0)
	if rt.FallingTorchAnimation != 1 {
		t.Fatalf("torch animation=%d, want collapse animation 1 after third trigger", rt.FallingTorchAnimation)
	}
}

func TestRuntimeStage05CollapseAndRisingFireCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage05.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.FallingTorchTriggers = 3
	rt.Player = Point{X: 18, Y: 41}
	rt.PlayerMotion = ObjectMotion{}
	rt.SetViewportY(900)

	for sourceTick := 1; sourceTick <= 37; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if rt.FallingTorchAnimation != 1 || rt.FallingTorchAnimationTicks != fallingTorchCollapseTicks {
		t.Fatalf("collapse after 37 ticks animation=%d ticks=%d, want 1/%d", rt.FallingTorchAnimation, rt.FallingTorchAnimationTicks, fallingTorchCollapseTicks)
	}
	rt.TickSourceFrame(8, 38, 0)
	if rt.FallingTorchAnimation != 2 || rt.RisingFireAnimation != 2 || rt.RisingFireAnimationTicks != 1 {
		t.Fatalf("post-collapse torch=%d fire=%d/%d, want 2/2/1", rt.FallingTorchAnimation, rt.RisingFireAnimation, rt.RisingFireAnimationTicks)
	}
	for sourceTick := 39; sourceTick <= 66; sourceTick++ {
		rt.TickSourceFrame(8, sourceTick, 0)
	}
	if rt.RisingFireAnimation != 0 || rt.RisingFireHeight != fallingFireInitialHeight {
		t.Fatalf("fire startup animation=%d height=%d, want 0/%d before first rise", rt.RisingFireAnimation, rt.RisingFireHeight, fallingFireInitialHeight)
	}
	rt.TickSourceFrame(8, 67, 0)
	if rt.RisingFireHeight != fallingFireInitialHeight+1 {
		t.Fatalf("first rising-fire height=%d, want %d", rt.RisingFireHeight, fallingFireInitialHeight+1)
	}

	rt.Player = Point{X: 16, Y: 41}
	result := rt.TickSourceFrame(8, 68, 0)
	if result.RisingFireHits != 1 || !rt.PlayerDead || rt.Health != 0 {
		t.Fatalf("rising fire hits=%d dead=%v health=%d, want 1/true/0", result.RisingFireHits, rt.PlayerDead, rt.Health)
	}
}

func TestRuntimeStage05DemoLocksCameraAndPausesFire(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage05.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.ForegroundDemoActive = true
	rt.ForegroundDemoID = 3
	rt.ForegroundDemoPhase = 1
	rt.ForegroundDemoTicks = 30
	rt.FallingTorchTriggers = 3
	rt.FallingTorchAnimation = 2
	rt.RisingFireAnimation = 0
	rt.RisingFireHeight = fallingFireInitialHeight
	rt.SetViewportY(600)

	x, y, elapsed, duration, ok := rt.ForegroundDemoCamera()
	if !ok || x != 180 || y != 900 || elapsed != 30 || duration != foregroundDemoPanTicks {
		t.Fatalf("demo camera=%d,%d elapsed=%d/%d ok=%v, want 180,900 30/%d true", x, y, elapsed, duration, ok, foregroundDemoPanTicks)
	}
	rt.TickSourceFrame(8, 1, 0)
	if rt.RisingFireHeight != fallingFireInitialHeight {
		t.Fatalf("fire rose during demo to %d, want paused at %d", rt.RisingFireHeight, fallingFireInitialHeight)
	}
}

func TestRuntimeStage05CheckpointRestoresSpecialState(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage05.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	rt.FallingTorchTriggers = 2
	rt.RisingFireHeight = 900
	rt.SetForTest(PlayerLayer, 18, 63, 0)
	rt.SaveSnapshot()
	rt.FallingTorchTriggers = 3
	rt.RisingFireHeight = 1200
	rt.FallingTorchAnimation = 2
	rt.ForegroundDemoActive = true
	if !rt.RestoreCheckpoint() {
		t.Fatal("failed to restore Stage 6 checkpoint")
	}
	rock, _ := rt.At(PlayerLayer, 18, 63)
	if rt.FallingTorchTriggers != 2 || rt.RisingFireHeight != 900 || rt.FallingTorchAnimation != 0 || rt.ForegroundDemoActive || rock != EmptyRawID {
		t.Fatalf("restored triggers=%d height=%d torch=%d demo=%v rock=%d, want 2/900/0/false/empty", rt.FallingTorchTriggers, rt.RisingFireHeight, rt.FallingTorchAnimation, rt.ForegroundDemoActive, rock)
	}
}

func TestRuntimeStage05CanBeCompletedAtSourceCadence(t *testing.T) {
	stage := mustLoadOriginalStage(t, "stage05.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	sourceTick := 0
	cameraY := clampRuntime(rt.Player.Y*TileSize-160, 0, rt.Height()*TileSize-(ScreenHeight-80))
	demoCameraLocked := false
	demoCameraStartY := 0
	updateCamera := func() {
		if _, targetY, elapsed, duration, ok := rt.ForegroundDemoCamera(); ok {
			if !demoCameraLocked {
				demoCameraLocked = true
				demoCameraStartY = cameraY
			}
			cameraY = (targetY*elapsed + demoCameraStartY*(duration-elapsed)) / duration
			return
		}
		demoCameraLocked = false
		playerY := rt.Player.Y*TileSize - rt.PlayerMotion.DY*rt.PlayerMotion.Remaining
		screenPlayerY := playerY + 40
		if screenPlayerY < cameraY+96 {
			cameraY = (cameraY - 96 + screenPlayerY) >> 1
		} else if screenPlayerY > cameraY+160 {
			cameraY = (cameraY - 160 + screenPlayerY) >> 1
		}
		cameraY = clampRuntime(cameraY, 0, rt.Height()*TileSize-(ScreenHeight-80))
	}
	tickUpdate := func() {
		sourceTick++
		rt.SetViewportY(cameraY)
		rt.TickSourceFrame(8, sourceTick, 0)
		rt.TickBreakables()
		rt.TickForegroundTriggers()
		if rt.PlayerMotion.Remaining > 0 {
			rt.AdvancePlayerMotion()
		}
		updateCamera()
	}
	busy := func() bool {
		return rt.PlayerMotion.Remaining > 0 || rt.HurtTicks > 0 || rt.ChestOpening || rt.LockOpening || rt.Hammering || rt.Hooking || rt.RecallPending || rt.ForegroundDemoActive
	}
	move := func(label string, dx, dy, count int) {
		t.Helper()
		for step := 1; step <= count; step++ {
			waits := 0
			for {
				if busy() {
					tickUpdate()
					continue
				}
				if rt.TryMove(dx, dy) {
					break
				}
				targetX, targetY := rt.Player.X+dx, rt.Player.Y+dy
				playerID, _ := rt.At(PlayerLayer, targetX, targetY)
				foregroundID, _ := rt.At(ForegroundLayer, targetX, targetY)
				if waits < 200 && ((playerID == 0 && dy == 0) || foregroundID == 7) {
					waits++
					tickUpdate()
					continue
				}
				nearby := make([]string, 0, 9)
				for ny := targetY - 1; ny <= targetY+1; ny++ {
					for nx := targetX - 1; nx <= targetX+1; nx++ {
						id, _ := rt.At(PlayerLayer, nx, ny)
						fg, _ := rt.At(ForegroundLayer, nx, ny)
						nearby = append(nearby, fmt.Sprintf("(%d,%d)=%d/%d", nx, ny, id, fg))
					}
				}
				t.Fatalf("%s step %d/%d failed tick=%d player=%+v target=(%d,%d) raw=%d foreground=%d health=%d triggers=%d fire=%d cameraY=%d remaining=%d nearby=%v", label, step, count, sourceTick, rt.Player, targetX, targetY, playerID, foregroundID, rt.Health, rt.FallingTorchTriggers, rt.RisingFireHeight, cameraY, rt.BonusRemaining, nearby)
			}
			for rt.PlayerMotion.Remaining > 0 {
				tickUpdate()
			}
			tickUpdate()
			if rt.PlayerDead {
				t.Fatalf("%s killed hero at tick=%d player=%+v fireY=%d cameraY=%d", label, sourceTick, rt.Player, rt.Height()*TileSize-rt.RisingFireHeight, cameraY)
			}
		}
	}
	waitUntil := func(label string, maxTicks int, condition func() bool) {
		t.Helper()
		for tick := 0; tick < maxTicks; tick++ {
			if condition() {
				return
			}
			tickUpdate()
		}
		t.Fatalf("%s not reached after %d ticks at tick=%d player=%+v health=%d triggers=%d fire=%d boulders=%v", label, maxTicks, sourceTick, rt.Player, rt.Health, rt.FallingTorchTriggers, rt.RisingFireHeight, runtimePointsWithRaw(rt, 0))
	}

	tickUpdate()
	move("automatic entrance first step", 1, 0, 1)
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close Stage 6 entrance door")
	}
	move("automatic entrance second step", 1, 0, 1)
	move("approach entrance boulder", 1, 0, 2)
	move("push entrance boulder right", 1, 0, 2)
	move("climb entrance dirt", 0, -1, 2)
	move("collect second entrance gem", 1, 0, 1)
	move("push ledge boulder over drop", 1, 0, 1)
	waitUntil("ledge boulder settles", 120, func() bool {
		id, _ := rt.At(PlayerLayer, 9, 69)
		return id == 0
	})
	move("enter lower shaft", 1, 0, 1)
	move("climb below falling pair", 0, -1, 1)
	move("stand below right falling boulder", 1, 0, 1)
	move("dig right falling-boulder support", 0, -1, 1)
	move("retreat from right falling boulder", 0, 1, 1)
	move("step out from under right falling boulder", -1, 0, 1)
	waitUntil("right falling boulder settles", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 10, 68)
		return id == 0
	})
	move("return to opened right shaft", 1, 0, 1)
	move("climb opened right shaft", 0, -1, 2)
	move("cross to gem shaft", -1, 0, 3)
	move("collect lower shaft gems", 0, -1, 3)
	move("approach pressure boulder", 1, 0, 6)
	move("push pressure boulder", 1, 0, 4)
	waitUntil("pressure boulder triggers warning", 160, func() bool {
		id, _ := rt.At(PlayerLayer, 18, 63)
		return id == 0 && rt.FallingTorchTriggers == 1
	})

	move("cross lower catwalk", 1, 0, 4)
	move("enter right ascent", 0, -1, 1)
	move("move to right ascent wall", 1, 0, 1)
	move("climb right ascent", 0, -1, 8)
	if rt.FallingTorchTriggers != 2 {
		t.Fatalf("first trigger strip count=%d, want 2", rt.FallingTorchTriggers)
	}
	move("cross middle ledge", -1, 0, 1)
	move("climb middle shaft", 0, -1, 5)
	move("cross upper connector left", -1, 0, 1)
	move("climb upper connector", 0, -1, 4)
	move("cross checkpoint connector right", 1, 0, 1)
	move("climb to checkpoint shelf", 0, -1, 1)
	move("enter checkpoint shaft", 1, 0, 1)
	move("climb checkpoint shaft", 0, -1, 2)
	move("activate checkpoint", -1, 0, 3)
	if rt.CheckpointProgress < 2 {
		t.Fatalf("checkpoint progress=%d, want upper Stage 6 checkpoint active", rt.CheckpointProgress)
	}
	move("cross to demo strip", -1, 0, 3)
	move("climb below demo strip", 0, -1, 3)
	move("trigger source demo 3", 0, -1, 1)
	waitUntil("falling-torches demo completes", 180, func() bool { return !rt.ForegroundDemoActive && rt.PlayerMotion.Remaining == 0 })
	if rt.Player != (Point{X: 16, Y: 36}) || rt.FallingTorchTriggers != 3 || rt.FallingTorchAnimation != 2 {
		t.Fatalf("post-demo player=%+v triggers=%d torch=%d, want (16,36)/3/2", rt.Player, rt.FallingTorchTriggers, rt.FallingTorchAnimation)
	}

	move("race left into upper maze", -1, 0, 1)
	move("race up first upper shaft", 0, -1, 16)
	move("race right connector", 1, 0, 1)
	move("race up second upper shaft", 0, -1, 7)
	move("race left connector", -1, 0, 1)
	move("race up third upper shaft", 0, -1, 3)
	move("race right to exit shaft", 1, 0, 1)
	move("race up exit shaft", 0, -1, 5)
	move("collect exit-row rewards", 1, 0, 8)
	move("pass quota and enter goal", 1, 0, 2)
	if !rt.ReachedGoal || !rt.BonusGateOpen || rt.BonusRemaining > 0 {
		t.Fatalf("Stage 6 goal reached=%v remaining=%d violet=%d player=%+v", rt.ReachedGoal, rt.BonusRemaining, rt.VioletGems, rt.Player)
	}
}
