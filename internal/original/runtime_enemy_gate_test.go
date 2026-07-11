package original

import "testing"

func TestRuntimeAngkorEnemyArenaDoorsFollowRaw26AndRaw17Groups(t *testing.T) {
	tests := []struct {
		stage      string
		trigger    Point
		dx         int
		entrance   Point
		target     Point
		enemyCount int
	}{
		{stage: "stage02.json", trigger: Point{X: 14, Y: 21}, dx: 1, entrance: Point{X: 13, Y: 21}, target: Point{X: 21, Y: 21}, enemyCount: 2},
		{stage: "stage03.json", trigger: Point{X: 19, Y: 14}, dx: -1, entrance: Point{X: 20, Y: 14}, target: Point{X: 13, Y: 14}, enemyCount: 2},
		{stage: "stage04.json", trigger: Point{X: 28, Y: 7}, dx: -1, entrance: Point{X: 29, Y: 7}, target: Point{X: 22, Y: 7}, enemyCount: 1},
		{stage: "stage06.json", trigger: Point{X: 12, Y: 28}, dx: 1, entrance: Point{X: 11, Y: 28}, target: Point{X: 18, Y: 32}, enemyCount: 1},
		{stage: "stage07.json", trigger: Point{X: 32, Y: 9}, dx: 1, entrance: Point{X: 31, Y: 9}, target: Point{X: 38, Y: 9}, enemyCount: 2},
		{stage: "stage08.json", trigger: Point{X: 9, Y: 5}, dx: 1, entrance: Point{X: 8, Y: 5}, target: Point{X: 18, Y: 5}, enemyCount: 1},
	}
	for _, test := range tests {
		t.Run(test.stage, func(t *testing.T) {
			rt, err := NewRuntime(mustLoadOriginalStage(t, test.stage))
			if err != nil {
				t.Fatal(err)
			}
			openDoorForEnemyGateTest(rt, test.entrance)
			if !rt.foregroundDoorOpen(test.entrance.X, test.entrance.Y) || !rt.foregroundDoorOpen(test.target.X, test.target.Y) {
				t.Fatalf("arena doors do not start open: entrance=%#x target=%#x", doorStateForEnemyGateTest(rt, test.entrance), doorStateForEnemyGateTest(rt, test.target))
			}

			rt.Player = test.trigger
			rt.PlayerMotion = ObjectMotion{DX: test.dx, Remaining: 6}
			rt.activateEnemyGateTriggerAt(test.trigger.X, test.trigger.Y)
			if rt.ActiveEnemyGateGroup != 0 {
				t.Fatalf("active group=%d, want 0", rt.ActiveEnemyGateGroup)
			}
			if rt.foregroundDoorOpen(test.entrance.X, test.entrance.Y) {
				t.Fatalf("entrance door %+v remained open after raw26", test.entrance)
			}
			if !rt.foregroundDoorOpen(test.target.X, test.target.Y) {
				t.Fatalf("target door %+v closed before the source camera delay", test.target)
			}
			if !rt.EnemyGateDemoActive || !rt.EnemyGateDemoTargetSet || rt.EnemyGateDemoTarget != test.target {
				t.Fatalf("demo active=%v targetSet=%v target=%+v, want %+v", rt.EnemyGateDemoActive, rt.EnemyGateDemoTargetSet, rt.EnemyGateDemoTarget, test.target)
			}
			if rt.CanAcceptInput() {
				t.Fatal("input remained enabled during the enemy-arena camera demo")
			}

			for tick := 0; tick < rt.EnemyGateDemoOutboundTicks+9; tick++ {
				rt.tickEnemyGateDemo()
			}
			if !rt.foregroundDoorOpen(test.target.X, test.target.Y) {
				t.Fatalf("target door %+v closed before pan plus ten source ticks", test.target)
			}
			rt.tickEnemyGateDemo()
			if rt.foregroundDoorOpen(test.target.X, test.target.Y) {
				t.Fatalf("target door %+v remained open after pan plus ten source ticks", test.target)
			}

			enemies := enemyGatePointsForTest(rt, 0)
			if len(enemies) != test.enemyCount || rt.EnemyGateCounters[0] != test.enemyCount {
				t.Fatalf("group enemies=%v counter=%d, want %d", enemies, rt.EnemyGateCounters[0], test.enemyCount)
			}
			for i, enemy := range enemies {
				rt.decrementEnemyGateForObjectAt(enemy.X, enemy.Y)
				if i < len(enemies)-1 && rt.foregroundDoorOpen(test.target.X, test.target.Y) {
					t.Fatalf("target door opened after only %d/%d grouped enemies", i+1, len(enemies))
				}
			}
			if rt.EnemyGateCounters[0] != 0 {
				t.Fatalf("group counter=%d after all enemies, want 0", rt.EnemyGateCounters[0])
			}
			for _, door := range []Point{test.entrance, test.target} {
				if doorStateForEnemyGateTest(rt, door)&0xf0 == 0 {
					t.Errorf("group completion did not begin opening door %+v", door)
				}
			}
		})
	}
}

func TestRuntimeRaw26WithoutEnemyGroupStillClosesDoorBehindHero(t *testing.T) {
	tests := []struct {
		stage    string
		trigger  Point
		dx       int
		entrance Point
	}{
		{stage: "stage03.json", trigger: Point{X: 22, Y: 19}, dx: 1, entrance: Point{X: 21, Y: 19}},
		{stage: "stage13.json", trigger: Point{X: 54, Y: 8}, dx: 1, entrance: Point{X: 53, Y: 8}},
	}
	for _, test := range tests {
		t.Run(test.stage, func(t *testing.T) {
			rt, err := NewRuntime(mustLoadOriginalStage(t, test.stage))
			if err != nil {
				t.Fatal(err)
			}
			openDoorForEnemyGateTest(rt, test.entrance)
			rt.Player = test.trigger
			rt.PlayerMotion = ObjectMotion{DX: test.dx, Remaining: 6}
			rt.activateEnemyGateTriggerAt(test.trigger.X, test.trigger.Y)
			if rt.foregroundDoorOpen(test.entrance.X, test.entrance.Y) {
				t.Fatalf("one-way door %+v remained open", test.entrance)
			}
			if rt.ActiveEnemyGateGroup != -1 || rt.EnemyGateDemoActive {
				t.Fatalf("invalid group activated group=%d demo=%v", rt.ActiveEnemyGateGroup, rt.EnemyGateDemoActive)
			}
		})
	}
}

func TestRuntimeStage11EnemyArenaDemoTargetsLockedKeyChests(t *testing.T) {
	tests := []struct {
		group   int
		trigger Point
		chest   Point
	}{
		{group: 0, trigger: Point{X: 12, Y: 6}, chest: Point{X: 8, Y: 4}},
		{group: 1, trigger: Point{X: 17, Y: 6}, chest: Point{X: 20, Y: 2}},
		{group: 2, trigger: Point{X: 28, Y: 6}, chest: Point{X: 32, Y: 6}},
		{group: 3, trigger: Point{X: 24, Y: 18}, chest: Point{X: 18, Y: 13}},
	}
	for _, test := range tests {
		t.Run(string(rune('0'+test.group)), func(t *testing.T) {
			rt, err := NewRuntime(mustLoadOriginalStage(t, "stage11.json"))
			if err != nil {
				t.Fatal(err)
			}
			rt.Player = test.trigger
			rt.PlayerMotion = ObjectMotion{DX: 1, Remaining: 6}
			rt.activateEnemyGateTriggerAt(test.trigger.X, test.trigger.Y)
			if !rt.EnemyGateDemoActive || rt.EnemyGateDemoTarget != test.chest || rt.ActiveEnemyGateGroup != test.group {
				t.Fatalf("demo active=%v target=%+v group=%d, want chest %+v group %d", rt.EnemyGateDemoActive, rt.EnemyGateDemoTarget, rt.ActiveEnemyGateGroup, test.chest, test.group)
			}
			if !rt.ContainerLockedAt(test.chest.X, test.chest.Y) {
				t.Fatal("arena key chest did not start locked")
			}
			enemies := enemyGatePointsForTest(rt, test.group)
			for i, enemy := range enemies {
				rt.decrementEnemyGateForObjectAt(enemy.X, enemy.Y)
				if i < len(enemies)-1 && !rt.ContainerLockedAt(test.chest.X, test.chest.Y) {
					t.Fatalf("key chest unlocked after only %d/%d grouped enemies", i+1, len(enemies))
				}
			}
			if rt.ContainerLockedAt(test.chest.X, test.chest.Y) {
				t.Fatal("key chest remained locked after all grouped enemies")
			}
		})
	}
}

func TestRuntimeEnemyGateDemoSupportsSourceRaw18Target(t *testing.T) {
	rt, err := NewRuntime(mustLoadOriginalStage(t, "stage11.json"))
	if err != nil {
		t.Fatal(err)
	}
	marker := Point{X: 8, Y: 5}
	rt.Foreground[rt.index(marker.X, marker.Y-1)] = EmptyRawID
	rt.PlayerLayer[rt.index(marker.X, marker.Y)] = 18
	target, ok := rt.enemyGateDemoTarget(0, Point{X: -1, Y: -1})
	if !ok || target != marker {
		t.Fatalf("raw18 target=%+v ok=%v, want %+v", target, ok, marker)
	}
}

func TestRuntimeEnemyGateDemoUsesSourceCameraSpeedAndMessageLifetime(t *testing.T) {
	rt, err := NewRuntime(mustLoadOriginalStage(t, "stage02.json"))
	if err != nil {
		t.Fatal(err)
	}
	trigger := Point{X: 14, Y: 21}
	openDoorForEnemyGateTest(rt, Point{X: 13, Y: 21})
	rt.Player = trigger
	rt.PlayerMotion = ObjectMotion{DX: 1, Remaining: 6}
	rt.SetViewport(0, 0)
	rt.activateEnemyGateTriggerAt(trigger.X, trigger.Y)
	targetX := clampRuntime(21*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	targetY := clampRuntime(21*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	wantOutbound := max(1, max((targetX+7)/8, (targetY+7)/8))
	if rt.EnemyGateDemoOutboundTicks != wantOutbound {
		t.Fatalf("outbound ticks=%d, want %d at source 8 px/tick", rt.EnemyGateDemoOutboundTicks, wantOutbound)
	}
	_, _, phase, _, duration, ok := rt.EnemyGateDemoCamera()
	if !ok || phase != 1 || duration != wantOutbound {
		t.Fatalf("camera phase/duration/ok=%d/%d/%v, want 1/%d/true", phase, duration, ok, wantOutbound)
	}
	for rt.EnemyGateDemoPhase == 1 || rt.EnemyGateDemoPhase == 2 {
		rt.tickEnemyGateDemo()
	}
	if index, ticks, ok := rt.EnemyGateMessage(); !ok || index != 56 || ticks != 80 {
		t.Fatalf("normal gate message index/ticks/ok=%d/%d/%v, want 56/80/true", index, ticks, ok)
	}
	for tick := 0; tick < 79; tick++ {
		rt.TickStatus()
	}
	if _, ticks, ok := rt.EnemyGateMessage(); !ok || ticks != 1 {
		t.Fatalf("message after 79 ticks ticks/ok=%d/%v, want 1/true", ticks, ok)
	}
	rt.TickStatus()
	if _, _, ok := rt.EnemyGateMessage(); ok {
		t.Fatal("enemy-gate message remained after source 80 ticks")
	}
}

func TestRuntimeEnemyGateDoorCloseResolvesSourceDoorOccupants(t *testing.T) {
	t.Run("grouped snake", func(t *testing.T) {
		rt, err := NewRuntime(mustLoadOriginalStage(t, "stage00.json"))
		if err != nil {
			t.Fatal(err)
		}
		point := Point{X: 6, Y: 6}
		idx := rt.index(point.X, point.Y)
		rt.Foreground[idx] = 7
		rt.Background[idx] = 0x30
		rt.PlayerLayer[idx] = 19
		rt.EnemyGateGroup[idx] = 0
		rt.EnemyGateCounters[0] = 1
		rt.ActiveEnemyGateGroup = 0
		if !rt.closeDoorAt(point.X, point.Y) {
			t.Fatal("failed to close source enemy-gate door")
		}
		if rt.PlayerLayer[idx] != EmptyRawID || rt.EnemyGateCounters[0] != 0 || rt.Background[idx]&0xf0 != 0 {
			t.Fatalf("closed-door snake raw=%d counter=%d door=%#x, want empty/0/closed", rt.PlayerLayer[idx], rt.EnemyGateCounters[0], rt.Background[idx])
		}
	})

	t.Run("hero", func(t *testing.T) {
		rt, err := NewRuntime(mustLoadOriginalStage(t, "stage00.json"))
		if err != nil {
			t.Fatal(err)
		}
		point := Point{X: 6, Y: 6}
		idx := rt.index(point.X, point.Y)
		rt.Foreground[idx] = 7
		rt.Background[idx] = 0x30
		rt.PlayerLayer[idx] = EmptyRawID
		rt.Player = point
		if !rt.closeDoorAt(point.X, point.Y) || !rt.PlayerDead || rt.Background[idx]&0xf0 != 0 {
			t.Fatalf("hero door close closed=%v dead=%v state=%#x, want true/true/closed", rt.Background[idx]&0xf0 == 0, rt.PlayerDead, rt.Background[idx])
		}
		if events := rt.DrainSoundEvents(); len(events) != 3 || events[0] != SoundBoulder || events[1] != SoundHeroHurt || events[2] != SoundDeath {
			t.Fatalf("hero door-crush sounds=%v, want [%d %d %d]", events, SoundBoulder, SoundHeroHurt, SoundDeath)
		}
	})

	t.Run("ungrouped boulder decrements active group", func(t *testing.T) {
		rt, err := NewRuntime(mustLoadOriginalStage(t, "stage00.json"))
		if err != nil {
			t.Fatal(err)
		}
		point := Point{X: 6, Y: 6}
		idx := rt.index(point.X, point.Y)
		rt.Foreground[idx] = 7
		rt.Background[idx] = 0x30
		rt.PlayerLayer[idx] = 0
		rt.EnemyGateGroup[idx] = -1
		rt.EnemyGateCounters[3] = 1
		rt.ActiveEnemyGateGroup = 3
		if !rt.closeDoorAt(point.X, point.Y) {
			t.Fatal("failed to close door onto boulder")
		}
		if rt.PlayerLayer[idx] != EmptyRawID || rt.EnemyGateCounters[3] != 0 {
			t.Fatalf("door-crushed boulder raw=%d active counter=%d, want empty/0", rt.PlayerLayer[idx], rt.EnemyGateCounters[3])
		}
	})

	t.Run("hook rope", func(t *testing.T) {
		rt, err := NewRuntime(mustLoadOriginalStage(t, "stage00.json"))
		if err != nil {
			t.Fatal(err)
		}
		point := Point{X: 6, Y: 6}
		idx := rt.index(point.X, point.Y)
		rt.Foreground[idx] = 7
		rt.Background[idx] = 0x30
		rt.PlayerLayer[idx] = 32
		if rt.closeDoorAt(point.X, point.Y) || rt.Background[idx] != 0x30 {
			t.Fatalf("door closed through source hook rope: closed=%v state=%#x", rt.Background[idx]&0xf0 == 0, rt.Background[idx])
		}
	})
}

func enemyGatePointsForTest(rt *Runtime, group int) []Point {
	points := make([]Point, 0, rt.EnemyGateCounters[group])
	for idx, candidate := range rt.EnemyGateGroup {
		if candidate == group {
			points = append(points, Point{X: idx % rt.Width(), Y: idx / rt.Width()})
		}
	}
	return points
}

func openDoorForEnemyGateTest(rt *Runtime, point Point) {
	idx := rt.index(point.X, point.Y)
	rt.Background[idx] = RawID(0x30 | int(rt.Background[idx])&0x0f)
}

func doorStateForEnemyGateTest(rt *Runtime, point Point) RawID {
	state, _ := rt.At(BackgroundLayer, point.X, point.Y)
	return state
}
