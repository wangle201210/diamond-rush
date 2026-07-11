package world

import (
	"testing"

	"github.com/wangle201210/zskc/internal/level"
)

func TestCollectDiamondsOpensExit(t *testing.T) {
	def := &level.Definition{
		Name:             "test",
		Width:            5,
		Height:           3,
		TileWidth:        32,
		TileHeight:       32,
		PlayerStart:      level.Point{X: 1, Y: 1},
		RequiredDiamonds: 1,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 3, 5, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if w.TotalDiamonds != 1 {
		t.Fatalf("total diamonds = %d, want 1", w.TotalDiamonds)
	}
	w.Update(Right)
	if w.Diamonds != 1 {
		t.Fatalf("diamonds = %d, want 1", w.Diamonds)
	}
	if w.Score != scoreDiamond {
		t.Fatalf("score = %d, want %d", w.Score, scoreDiamond)
	}
	if got := w.TileAt(3, 1); got != ExitOpen {
		t.Fatalf("exit tile = %v, want open", got)
	}
}

func TestSecretExitWinsAndMarksRoute(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 21, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Status != Won {
		t.Fatalf("status = %v, want won", w.Status)
	}
	if !w.SecretExitFound {
		t.Fatal("secret exit flag = false, want true")
	}
	wantEvents := []Event{EventStep, EventSecretExit, EventWin}
	if got := w.Events(); len(got) != len(wantEvents) {
		t.Fatalf("events = %v, want %v", got, wantEvents)
	} else {
		for i := range wantEvents {
			if got[i] != wantEvents[i] {
				t.Fatalf("events = %v, want %v", got, wantEvents)
			}
		}
	}
}

func TestHiddenWallRevealsAndLetsPlayerEnter(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 22, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if got := w.TileAt(2, 1); got != HiddenWall {
		t.Fatalf("tile = %v, want hidden wall", got)
	}
	w.Update(Right)
	if w.Player.X != 2 || w.Player.Y != 1 {
		t.Fatalf("player = (%d,%d), want (2,1)", w.Player.X, w.Player.Y)
	}
	if got := w.TileAt(2, 1); got != Empty {
		t.Fatalf("revealed tile = %v, want empty", got)
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventReveal {
		t.Fatalf("events = %v, want reveal", w.Events())
	}
}

func TestKeyAndDoorAddScore(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 6, 7, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	w.Update(Right)
	want := scoreKey + scoreDoor
	if w.Score != want {
		t.Fatalf("score = %d, want %d", w.Score, want)
	}
}

func TestGoldKeyAndDoorUseSeparateCounter(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       7,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1, 1, 1,
			1, 0, 6, 26, 25, 26, 1,
			1, 1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Keys != 1 || w.GoldKeys != 0 {
		t.Fatalf("keys = %d/%d, want silver=1 gold=0", w.Keys, w.GoldKeys)
	}
	w.Update(Right)
	if w.Player.X != 2 || w.TileAt(3, 1) != GoldDoor {
		t.Fatalf("silver key opened gold door: player=%+v tile=%v", w.Player, w.TileAt(3, 1))
	}
	w.setTile(3, 1, Empty)
	w.Update(Right)
	w.Update(Right)
	if w.GoldKeys != 1 {
		t.Fatalf("gold keys = %d, want 1", w.GoldKeys)
	}
	w.Update(Right)
	if w.Player.X != 5 || w.GoldKeys != 0 || w.TileAt(5, 1) != Empty {
		t.Fatalf("gold door did not open: player=%+v gold=%d tile=%v", w.Player, w.GoldKeys, w.TileAt(5, 1))
	}
	want := scoreKey + scoreGoldKey + scoreGoldDoor
	if w.Score != want {
		t.Fatalf("score = %d, want %d", w.Score, want)
	}
}

func TestGoldKeyDoesNotOpenSilverDoor(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 25, 7, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	w.Update(Right)
	if w.Player.X != 2 || w.TileAt(3, 1) != Door {
		t.Fatalf("gold key opened silver door: player=%+v tile=%v", w.Player, w.TileAt(3, 1))
	}
}

func TestChestAddsScoreWithoutOpeningExit(t *testing.T) {
	def := &level.Definition{
		Name:             "test",
		Width:            6,
		Height:           3,
		TileWidth:        32,
		TileHeight:       32,
		PlayerStart:      level.Point{X: 1, Y: 1},
		RequiredDiamonds: 1,
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 15, 5, 3, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Score != scoreChest {
		t.Fatalf("score = %d, want %d", w.Score, scoreChest)
	}
	if w.Diamonds != 0 {
		t.Fatalf("diamonds = %d, want 0", w.Diamonds)
	}
	if got := w.TileAt(3, 1); got != ExitClosed {
		t.Fatalf("exit tile = %v, want closed", got)
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventChest {
		t.Fatalf("events = %v, want chest event", w.Events())
	}
}

func TestChestCanContainRedDiamond(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		ChestRewards: []level.ChestReward{
			{Point: level.Point{X: 2, Y: 1}, Reward: "red_diamond", Amount: 1},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 15, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if w.TotalRedDiamonds != 1 {
		t.Fatalf("total red diamonds = %d, want 1", w.TotalRedDiamonds)
	}
	w.Update(Right)
	if w.RedDiamonds != 1 {
		t.Fatalf("red diamonds = %d, want 1", w.RedDiamonds)
	}
	if w.Score != scoreRedDiamond {
		t.Fatalf("score = %d, want %d", w.Score, scoreRedDiamond)
	}
	wantEvents := []Event{EventChest, EventRedDiamond}
	if got := w.Events(); len(got) != len(wantEvents) || got[0] != wantEvents[0] || got[1] != wantEvents[1] {
		t.Fatalf("events = %v, want %v", got, wantEvents)
	}
}

func TestChestCanContainExtraLife(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		ChestRewards: []level.ChestReward{
			{Point: level.Point{X: 2, Y: 1}, Reward: "extra_life", Amount: 1},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 15, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Lives != defaultLives+1 {
		t.Fatalf("lives = %d, want %d", w.Lives, defaultLives+1)
	}
	wantEvents := []Event{EventChest, EventExtraLife}
	if got := w.Events(); len(got) != len(wantEvents) || got[0] != wantEvents[0] || got[1] != wantEvents[1] {
		t.Fatalf("events = %v, want %v", got, wantEvents)
	}
}

func TestChestCanContainPotion(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		ChestRewards: []level.ChestReward{
			{Point: level.Point{X: 2, Y: 1}, Reward: "potion", Amount: 2},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 15, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Health = 1
	w.Update(Right)
	if w.Health != defaultHealth {
		t.Fatalf("health = %d, want healed to %d", w.Health, defaultHealth)
	}
	wantEvents := []Event{EventChest, EventPotion}
	if got := w.Events(); len(got) != len(wantEvents) || got[0] != wantEvents[0] || got[1] != wantEvents[1] {
		t.Fatalf("events = %v, want %v", got, wantEvents)
	}
}

func TestChestCanContainTool(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		ChestRewards: []level.ChestReward{
			{Point: level.Point{X: 2, Y: 1}, Reward: "hook", Amount: 1},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 15, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if !w.HasHook {
		t.Fatal("has hook = false, want chest reward to grant hook")
	}
	wantEvents := []Event{EventChest, EventToolHook}
	if got := w.Events(); len(got) != len(wantEvents) || got[0] != wantEvents[0] || got[1] != wantEvents[1] {
		t.Fatalf("events = %v, want %v", got, wantEvents)
	}
}

func TestChestRewardMustMatchChestTile(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		ChestRewards: []level.ChestReward{
			{Point: level.Point{X: 2, Y: 1}, Reward: "red_diamond", Amount: 1},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	if _, err := New(def); err == nil {
		t.Fatal("New() error = nil, want chest reward placement error")
	}
}

func TestPushRock(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Player.X != 2 || w.TileAt(3, 1) != Rock {
		t.Fatalf("push failed: player=(%d,%d), rock=%v", w.Player.X, w.Player.Y, w.TileAt(3, 1))
	}
}

func TestChestPurpleDiamondCountsTowardTotalDiamonds(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       6,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		ChestRewards: []level.ChestReward{
			{Point: level.Point{X: 2, Y: 1}, Reward: "purple_diamond", Amount: 2},
		},
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 15, 3, 0, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if w.TotalDiamonds != 3 {
		t.Fatalf("total diamonds = %d, want tile + chest total 3", w.TotalDiamonds)
	}
	w.Update(Right)
	if w.Diamonds != 2 {
		t.Fatalf("diamonds after chest = %d, want 2", w.Diamonds)
	}
}

func TestHookPullsRockTowardPlayer(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHook:     true,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 4, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.UseHook(Right) {
		t.Fatal("UseHook() = false, want true")
	}
	if got := w.TileAt(2, 1); got != Rock {
		t.Fatalf("adjacent tile = %v, want pulled rock", got)
	}
	if got := w.TileAt(3, 1); got != Empty {
		t.Fatalf("source tile = %v, want empty", got)
	}
	if w.Steps != 1 {
		t.Fatalf("steps = %d, want 1", w.Steps)
	}
}

func TestHookPullsDistantObjectOneTile(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       7,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHook:     true,
		Tiles: []int{
			1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 0, 4, 0, 1,
			1, 1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.UseHook(Right) {
		t.Fatal("UseHook() = false, want true")
	}
	if got := w.TileAt(3, 1); got != Rock {
		t.Fatalf("pulled tile = %v, want rock one step closer", got)
	}
	if got := w.TileAt(4, 1); got != Empty {
		t.Fatalf("old rock tile = %v, want empty", got)
	}
}

func TestHookCannotPullThroughBlockedLine(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       7,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHook:     true,
		Tiles: []int{
			1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 2, 4, 0, 1,
			1, 1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if w.UseHook(Right) {
		t.Fatal("UseHook() through dirt = true, want false")
	}
	if got := w.TileAt(4, 1); got != Rock {
		t.Fatalf("rock moved through blocker: got %v", got)
	}
}

func TestUpdateHookEmitsHookEvent(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHook:     true,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 4, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.UpdateHook(Right) {
		t.Fatal("UpdateHook() = false, want true")
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventHook {
		t.Fatalf("events = %v, want hook event", w.Events())
	}
}

func TestUpdateHookDoesNotAdvanceOnInvalidHook(t *testing.T) {
	def := &level.Definition{
		Name:         "test",
		Width:        5,
		Height:       4,
		TileWidth:    32,
		TileHeight:   32,
		PlayerStart:  level.Point{X: 1, Y: 1},
		HasHook:      true,
		GravityTicks: 1,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 0, 1,
			1, 0, 4, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if w.UpdateHook(Left) {
		t.Fatal("UpdateHook() = true, want invalid hook")
	}
	if got := w.TileAt(2, 2); got != Rock {
		t.Fatalf("rock moved after invalid hook: got %v", got)
	}
}

func TestHookPullsDiamondButDoesNotCollectIt(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHook:     true,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 3, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.UseHook(Right) {
		t.Fatal("UseHook() = false, want true")
	}
	if got := w.TileAt(2, 1); got != Diamond {
		t.Fatalf("adjacent tile = %v, want pulled diamond", got)
	}
	if w.Diamonds != 0 || w.Score != 0 {
		t.Fatalf("diamond was collected by hook: diamonds=%d score=%d", w.Diamonds, w.Score)
	}
}

func TestHookRequiresAbilityAndClearAdjacentTile(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 2, 4, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if w.UseHook(Right) {
		t.Fatal("UseHook() without hook ability = true, want false")
	}
	w.HasHook = true
	if w.UseHook(Right) {
		t.Fatal("UseHook() through blocked adjacent tile = true, want false")
	}
}

func TestHammerBreaksAdjacentCrackedWall(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHammer:   true,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 12, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.UseHammer(Right) {
		t.Fatal("UseHammer() = false, want true")
	}
	if got := w.TileAt(2, 1); got != Empty {
		t.Fatalf("cracked wall tile = %v, want empty", got)
	}
	if w.Steps != 1 {
		t.Fatalf("steps = %d, want 1", w.Steps)
	}
	wantEvents := []Event{EventHammer, EventBreak}
	if got := w.Events(); len(got) != len(wantEvents) || got[0] != wantEvents[0] || got[1] != wantEvents[1] {
		t.Fatalf("events = %v, want %v", got, wantEvents)
	}
}

func TestHammerStunsAdjacentEnemy(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		HasHammer:   true,
		EnemyTicks:  1,
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 2, Y: 1}, Type: "enemy_horizontal", Direction: "right"},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.UpdateAction(Right) {
		t.Fatal("UpdateAction() = false, want hammer hit")
	}
	if w.Enemies[0].Stunned != 2 {
		t.Fatalf("enemy stun = %d, want decremented to 2 after advance", w.Enemies[0].Stunned)
	}
	if w.Enemies[0].X != 2 || w.Enemies[0].Y != 1 {
		t.Fatalf("enemy moved while stunned: %+v", w.Enemies[0])
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventHammer {
		t.Fatalf("events = %v, want hammer event", w.Events())
	}
}

func TestBossBlocksExitUntilDefeated(t *testing.T) {
	def := &level.Definition{
		Name:             "test",
		Width:            6,
		Height:           3,
		TileWidth:        32,
		TileHeight:       32,
		PlayerStart:      level.Point{X: 1, Y: 1},
		RequiredDiamonds: 1,
		HasHammer:        true,
		BossStarts: []level.BossSpawn{
			{Point: level.Point{X: 3, Y: 1}, Type: "boss_guardian", HP: 2},
		},
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 3, 0, 5, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if got := w.TileAt(4, 1); got != ExitClosed {
		t.Fatalf("exit with boss alive = %v, want closed", got)
	}
	if !w.UpdateAction(Right) {
		t.Fatal("first hammer hit = false, want boss hit")
	}
	if hp, maxHP, ok := w.BossHealth(); !ok || hp != 1 || maxHP != 2 {
		t.Fatalf("boss hp = %d/%d ok=%v, want 1/2 true", hp, maxHP, ok)
	}
	if !w.UpdateAction(Right) {
		t.Fatal("second hammer hit = false, want boss defeat")
	}
	if w.BossAlive() {
		t.Fatal("boss alive after second hit, want defeated")
	}
	if got := w.TileAt(4, 1); got != ExitOpen {
		t.Fatalf("exit after boss defeat = %v, want open", got)
	}
}

func TestFallingRockDamagesBoss(t *testing.T) {
	def := &level.Definition{
		Name:         "test",
		Width:        5,
		Height:       5,
		TileWidth:    32,
		TileHeight:   32,
		PlayerStart:  level.Point{X: 1, Y: 1},
		GravityTicks: 1,
		BossStarts: []level.BossSpawn{
			{Point: level.Point{X: 2, Y: 3}, Type: "boss_guardian", HP: 2},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 0, 0, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Direction{})
	w.Update(Direction{})
	if w.BossAlive() {
		t.Fatal("boss alive after falling rock hit, want defeated")
	}
	if got := w.TileAt(2, 2); got != Empty {
		t.Fatalf("falling rock tile = %v, want removed after boss hit", got)
	}
}

func TestSpikeConsumesHealthBeforeLife(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 9, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Status != Playing {
		t.Fatalf("status = %v, want still playing", w.Status)
	}
	if w.Lives != defaultLives {
		t.Fatalf("lives = %d, want %d", w.Lives, defaultLives)
	}
	if w.Health != defaultHealth-1 {
		t.Fatalf("health = %d, want %d", w.Health, defaultHealth-1)
	}
	if !w.Damaged {
		t.Fatal("damaged = false after spike hit, want true")
	}
	if w.Player.X != 2 || w.Player.Y != 1 {
		t.Fatalf("player = (%d,%d), want spike tile", w.Player.X, w.Player.Y)
	}
	if len(w.Events()) != 2 || w.Events()[0] != EventDamage || w.Events()[1] != EventTrap {
		t.Fatalf("events = %v, want damage then trap", w.Events())
	}
}

func TestTimedSpikeDamagesWhenActive(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 23, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if !w.TimedSpikeActive() {
		t.Fatal("timed spike active = false, want initially active")
	}
	w.Update(Right)
	if w.Health != defaultHealth-1 {
		t.Fatalf("health = %d, want %d", w.Health, defaultHealth-1)
	}
	if len(w.Events()) < 2 || w.Events()[0] != EventDamage || w.Events()[1] != EventTrap {
		t.Fatalf("events = %v, want damage then trap", w.Events())
	}
}

func TestTimedSpikeDamagesWhenOpeningUnderPlayer(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 23, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.trapPhase = timedSpikeActive
	w.Update(Right)
	if w.Health != defaultHealth {
		t.Fatalf("health after inactive entry = %d, want %d", w.Health, defaultHealth)
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventStep {
		t.Fatalf("events after inactive entry = %v, want step", w.Events())
	}
	for w.trapPhase != timedSpikeCycle-1 {
		w.Update(Direction{})
	}
	w.Update(Direction{})
	if w.Health != defaultHealth-1 {
		t.Fatalf("health after spike opens = %d, want %d", w.Health, defaultHealth-1)
	}
	if len(w.Events()) < 2 || w.Events()[0] != EventDamage || w.Events()[1] != EventTrap {
		t.Fatalf("events after spike opens = %v, want damage then trap", w.Events())
	}
}

func TestFireTrapDamagesWhenActive(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 24, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.trapPhase = fireTrapStart
	if !w.FireTrapActive() {
		t.Fatal("fire trap active = false, want active")
	}
	w.Update(Right)
	if w.Health != defaultHealth-1 {
		t.Fatalf("health = %d, want %d", w.Health, defaultHealth-1)
	}
	if len(w.Events()) < 2 || w.Events()[0] != EventDamage || w.Events()[1] != EventBurn {
		t.Fatalf("events = %v, want damage then burn", w.Events())
	}
}

func TestFireTrapDamagesWhenIgnitingUnderPlayer(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 24, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Health != defaultHealth {
		t.Fatalf("health after inactive entry = %d, want %d", w.Health, defaultHealth)
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventStep {
		t.Fatalf("events after inactive entry = %v, want step", w.Events())
	}
	for w.trapPhase != fireTrapStart-1 {
		w.Update(Direction{})
	}
	w.Update(Direction{})
	if w.Health != defaultHealth-1 {
		t.Fatalf("health after fire ignites = %d, want %d", w.Health, defaultHealth-1)
	}
	if len(w.Events()) < 2 || w.Events()[0] != EventDamage || w.Events()[1] != EventBurn {
		t.Fatalf("events after fire ignites = %v, want damage then burn", w.Events())
	}
}

func TestFallingObjectBurnsInFireTrap(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      4,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 24, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Falling[w.idx(2, 1)] = true
	if !w.tryFall(2, 1) {
		t.Fatal("tryFall() = false, want falling rock to burn in fire trap")
	}
	if got := w.TileAt(2, 1); got != Empty {
		t.Fatalf("source tile = %v, want empty", got)
	}
	if got := w.TileAt(2, 2); got != FireTrap {
		t.Fatalf("fire trap tile = %v, want FireTrap", got)
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventBurn {
		t.Fatalf("events = %v, want burn", w.Events())
	}
}

func TestArmorAbsorbsDamageBeforeHealth(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 9, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.MaxArmor = 1
	w.Armor = 1
	w.Update(Right)
	if w.Armor != 0 {
		t.Fatalf("armor = %d, want 0", w.Armor)
	}
	if w.Health != defaultHealth {
		t.Fatalf("health = %d, want unchanged %d", w.Health, defaultHealth)
	}
	if w.Lives != defaultLives {
		t.Fatalf("lives = %d, want unchanged %d", w.Lives, defaultLives)
	}
}

func TestCrushingDamageBypassesArmor(t *testing.T) {
	def := &level.Definition{
		Name:         "test",
		Width:        5,
		Height:       5,
		TileWidth:    32,
		TileHeight:   32,
		PlayerStart:  level.Point{X: 2, Y: 3},
		GravityTicks: 1,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 0, 0, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.MaxArmor = 2
	w.Armor = 2
	w.Update(Direction{})
	w.Update(Direction{})
	if w.Lives != defaultLives-1 {
		t.Fatalf("lives = %d, want %d", w.Lives, defaultLives-1)
	}
	if w.Health != defaultHealth {
		t.Fatalf("health = %d, want reset to %d", w.Health, defaultHealth)
	}
	if w.Armor != 2 {
		t.Fatalf("armor = %d, want reset to 2", w.Armor)
	}
}

func TestHealthDepletedConsumesLifeAndRespawnsAtCheckpoint(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 9, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Health = 1
	w.Update(Right)
	if w.Status != Playing {
		t.Fatalf("status = %v, want still playing", w.Status)
	}
	if w.Lives != defaultLives-1 {
		t.Fatalf("lives = %d, want %d", w.Lives, defaultLives-1)
	}
	if w.Health != defaultHealth {
		t.Fatalf("health = %d, want reset to %d", w.Health, defaultHealth)
	}
	if w.Player.X != 1 || w.Player.Y != 1 {
		t.Fatalf("player = (%d,%d), want checkpoint", w.Player.X, w.Player.Y)
	}
}

func TestPotionRestoresHealthWithoutExceedingMax(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 17, 17, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Health = 2
	w.Update(Right)
	if w.Health != 3 {
		t.Fatalf("health = %d, want 3", w.Health)
	}
	if got := w.TileAt(2, 1); got != Empty {
		t.Fatalf("potion tile = %v, want empty after pickup", got)
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventPotion {
		t.Fatalf("events = %v, want potion event", w.Events())
	}
	w.Update(Right)
	if w.Health != 3 {
		t.Fatalf("health after max potion = %d, want capped at 3", w.Health)
	}
}

func TestToolPickupsGrantAbilities(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       6,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 20, 18, 19, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if !w.HasCompass || w.TileAt(2, 1) != Empty {
		t.Fatalf("compass pickup failed: has=%v tile=%v", w.HasCompass, w.TileAt(2, 1))
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventCompass {
		t.Fatalf("events = %v, want compass event", w.Events())
	}
	w.Update(Right)
	if !w.HasHammer || w.TileAt(3, 1) != Empty {
		t.Fatalf("hammer pickup failed: has=%v tile=%v", w.HasHammer, w.TileAt(3, 1))
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventToolHammer {
		t.Fatalf("events = %v, want hammer pickup event", w.Events())
	}
	w.Update(Right)
	if !w.HasHook || w.TileAt(4, 1) != Empty {
		t.Fatalf("hook pickup failed: has=%v tile=%v", w.HasHook, w.TileAt(4, 1))
	}
	if len(w.Events()) != 1 || w.Events()[0] != EventToolHook {
		t.Fatalf("events = %v, want hook pickup event", w.Events())
	}
}

func TestLastLifeLostEndsStage(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 9, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Lives = 1
	w.Health = 1
	w.Update(Right)
	if w.Status != Lost {
		t.Fatalf("status = %v, want lost", w.Status)
	}
}

func TestCheckpointActivationAndRecallCostsLife(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       7,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1, 1, 1,
			1, 0, 16, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Checkpoint != (Player{X: 2, Y: 1}) {
		t.Fatalf("checkpoint = %+v, want {2 1}", w.Checkpoint)
	}
	w.Update(Right)
	w.Update(Right)
	if !w.RecallCheckpoint() {
		t.Fatal("RecallCheckpoint() = false, want true")
	}
	if !w.RecallUsed {
		t.Fatal("recall used = false, want true")
	}
	if w.Player != (Player{X: 2, Y: 1}) {
		t.Fatalf("player = %+v, want checkpoint", w.Player)
	}
	if w.Lives != defaultLives-1 {
		t.Fatalf("lives = %d, want %d", w.Lives, defaultLives-1)
	}
	if len(w.Events()) == 0 || w.Events()[0] != EventRecall {
		t.Fatalf("events = %v, want recall event", w.Events())
	}
}

func TestCheckpointActionResetsStateSinceActivation(t *testing.T) {
	def := &level.Definition{
		Name:             "test",
		Width:            6,
		Height:           3,
		TileWidth:        32,
		TileHeight:       32,
		PlayerStart:      level.Point{X: 1, Y: 1},
		RequiredDiamonds: 9,
		GravityTicks:     99,
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 16, 3, 0, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	w.Update(Right)
	if w.Diamonds != 1 || w.TileAt(3, 1) != Empty {
		t.Fatalf("diamond not collected before reset: count=%d tile=%v", w.Diamonds, w.TileAt(3, 1))
	}
	w.Update(Left)
	if !w.UpdateAction(Right) {
		t.Fatal("UpdateAction() = false, want checkpoint reset")
	}
	if w.Player != (Player{X: 2, Y: 1}) {
		t.Fatalf("player = %+v, want checkpoint", w.Player)
	}
	if w.Diamonds != 0 || w.Score != 0 {
		t.Fatalf("progress = diamonds %d score %d, want reset", w.Diamonds, w.Score)
	}
	if got := w.TileAt(3, 1); got != Diamond {
		t.Fatalf("tile after reset = %v, want diamond restored", got)
	}
	if events := w.Events(); len(events) == 0 || events[0] != EventReset {
		t.Fatalf("events = %v, want reset event", events)
	}
}

func TestCompassPointsToNextCheckpoint(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       7,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1, 1, 1,
			1, 0, 16, 0, 0, 16, 1,
			1, 1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	dir, distance, ok := w.CompassToCheckpoint()
	if !ok || dir != Right || distance != 1 {
		t.Fatalf("compass = %+v %d %t, want right distance 1", dir, distance, ok)
	}
	w.Update(Right)
	dir, distance, ok = w.CompassToCheckpoint()
	if !ok || dir != Right || distance != 3 {
		t.Fatalf("compass after checkpoint = %+v %d %t, want right distance 3", dir, distance, ok)
	}
}

func TestSwitchOpensBridges(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       6,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 10, 11, 0, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if got := w.TileAt(3, 1); got != Empty {
		t.Fatalf("bridge tile = %v, want empty", got)
	}
}

func TestFallingRockBreaksCrackedWall(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      5,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 0, 0, 1,
			1, 0, 12, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 18; i++ {
		w.Update(Direction{})
	}
	if got := w.TileAt(2, 3); got != Empty {
		t.Fatalf("cracked wall tile = %v, want empty", got)
	}
	if got := w.TileAt(2, 2); got != Empty {
		t.Fatalf("rock above cracked wall = %v, want empty after impact", got)
	}
}

func TestTeleporterMovesPlayerToNextTeleporter(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       6,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1, 1,
			1, 0, 13, 13, 0, 1,
			1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Player.X != 3 || w.Player.Y != 1 {
		t.Fatalf("player = (%d,%d), want second teleporter", w.Player.X, w.Player.Y)
	}
}

func TestLavaConsumesHealth(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 14, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Right)
	if w.Status != Playing {
		t.Fatalf("status = %v, want still playing", w.Status)
	}
	if w.Lives != defaultLives {
		t.Fatalf("lives = %d, want %d", w.Lives, defaultLives)
	}
	if w.Health != defaultHealth-1 {
		t.Fatalf("health = %d, want %d", w.Health, defaultHealth-1)
	}
}

func TestFallingRockBurnsInLava(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      4,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 14, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 9; i++ {
		w.Update(Direction{})
	}
	if got := w.TileAt(2, 1); got != Empty {
		t.Fatalf("rock tile = %v, want empty after lava burn", got)
	}
	if got := w.TileAt(2, 2); got != Lava {
		t.Fatalf("lava tile = %v, want lava to remain", got)
	}
}

func TestCustomGravityTicksControlFallingSpeed(t *testing.T) {
	def := &level.Definition{
		Name:         "test",
		Width:        5,
		Height:       5,
		TileWidth:    32,
		TileHeight:   32,
		PlayerStart:  level.Point{X: 1, Y: 1},
		GravityTicks: 3,
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 0, 0, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		w.Update(Direction{})
	}
	if got := w.TileAt(2, 1); got != Rock {
		t.Fatalf("rock moved before custom gravity tick: got %v", got)
	}
	w.Update(Direction{})
	if got := w.TileAt(2, 2); got != Rock {
		t.Fatalf("rock tile after custom gravity tick = %v, want rock at y=2", got)
	}
}

func TestFallingRockCrushingEnemyAddsScore(t *testing.T) {
	def := &level.Definition{
		Name:         "test",
		Width:        5,
		Height:       5,
		TileWidth:    32,
		TileHeight:   32,
		PlayerStart:  level.Point{X: 1, Y: 1},
		GravityTicks: 1,
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 2, Y: 3}, Type: "enemy_vertical"},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 4, 0, 1,
			1, 0, 0, 0, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Direction{})
	w.Update(Direction{})
	if len(w.Enemies) != 0 {
		t.Fatalf("enemies = %d, want crushed enemy removed", len(w.Enemies))
	}
	if w.Score != scoreEnemy {
		t.Fatalf("score = %d, want %d", w.Score, scoreEnemy)
	}
}

func TestCustomEnemyTicksControlPatrolSpeed(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      5,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		EnemyTicks:  4,
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 3, Y: 1}, Type: "enemy_vertical"},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 1, 0, 1,
			1, 0, 1, 0, 1,
			1, 0, 1, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		w.Update(Direction{})
	}
	if w.Enemies[0].Y != 1 {
		t.Fatalf("enemy moved before custom enemy tick: y=%d", w.Enemies[0].Y)
	}
	w.Update(Direction{})
	if w.Enemies[0].X != 3 || w.Enemies[0].Y != 2 {
		t.Fatalf("enemy = (%d,%d), want custom-tick move to (3,2)", w.Enemies[0].X, w.Enemies[0].Y)
	}
}

func TestEnemySpawnDirectionControlsInitialPatrol(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		EnemyTicks:  1,
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 3, Y: 1}, Type: "enemy_horizontal", Direction: "left"},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 0, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	w.Update(Direction{})
	if w.Enemies[0].X != 2 || w.Enemies[0].Y != 1 {
		t.Fatalf("enemy = (%d,%d), want initial left move to (2,1)", w.Enemies[0].X, w.Enemies[0].Y)
	}
}

func TestInvalidEnemySpawnDirectionFails(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       3,
		Height:      3,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 1, Y: 1}, Type: "enemy_horizontal", Direction: "diagonal"},
		},
		Tiles: []int{
			1, 1, 1,
			1, 0, 1,
			1, 1, 1,
		},
	}
	if _, err := New(def); err == nil {
		t.Fatal("New() error = nil, want invalid enemy direction error")
	}
}

func TestVerticalEnemyPatrolsVertically(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       5,
		Height:      5,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 1},
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 3, Y: 1}, Type: "enemy_vertical"},
		},
		Tiles: []int{
			1, 1, 1, 1, 1,
			1, 0, 1, 0, 1,
			1, 0, 1, 0, 1,
			1, 0, 1, 0, 1,
			1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 18; i++ {
		w.Update(Direction{})
	}
	if w.Enemies[0].X != 3 || w.Enemies[0].Y != 2 {
		t.Fatalf("enemy = (%d,%d), want vertical move to (3,2)", w.Enemies[0].X, w.Enemies[0].Y)
	}
}

func TestChaserEnemyMovesTowardPlayer(t *testing.T) {
	def := &level.Definition{
		Name:        "test",
		Width:       7,
		Height:      5,
		TileWidth:   32,
		TileHeight:  32,
		PlayerStart: level.Point{X: 1, Y: 2},
		EnemyStarts: []level.EnemySpawn{
			{Point: level.Point{X: 5, Y: 2}, Type: "enemy_chaser"},
		},
		Tiles: []int{
			1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1,
		},
	}
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 18; i++ {
		w.Update(Direction{})
	}
	if w.Enemies[0].X != 4 || w.Enemies[0].Y != 2 {
		t.Fatalf("enemy = (%d,%d), want chaser move to (4,2)", w.Enemies[0].X, w.Enemies[0].Y)
	}
}
