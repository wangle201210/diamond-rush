package game

import (
	"testing"

	"github.com/wangle201210/zskc/internal/world"
)

func TestGeneratedAnimationFrames(t *testing.T) {
	assets := NewAssets(tileSize)
	if assets.MenuBackdrop == nil {
		t.Fatal("menu backdrop was not loaded")
	}
	if len(assets.PlayerFrames) != 4 {
		t.Fatalf("player frames = %d, want 4", len(assets.PlayerFrames))
	}
	if len(assets.DiamondFrames) != 3 {
		t.Fatalf("diamond frames = %d, want 3", len(assets.DiamondFrames))
	}
	for _, enemyType := range []world.EnemyType{world.EnemyHorizontal, world.EnemyVertical, world.EnemyChaser} {
		if len(assets.EnemyFrames[enemyType]) != 4 {
			t.Fatalf("enemy type %v frames = %d, want 4", enemyType, len(assets.EnemyFrames[enemyType]))
		}
	}
	if assets.PlayerFrame(0, false) == nil || assets.PlayerFrame(0, true) == nil {
		t.Fatal("player frame lookup returned nil")
	}
	if assets.DiamondFrame(0) == nil {
		t.Fatal("diamond frame lookup returned nil")
	}
	if assets.Tile(world.Chest) == nil {
		t.Fatal("chest tile lookup returned nil")
	}
	if assets.Tile(world.Checkpoint) == nil {
		t.Fatal("checkpoint tile lookup returned nil")
	}
	if assets.Tile(world.Potion) == nil {
		t.Fatal("potion tile lookup returned nil")
	}
	if assets.Tile(world.HammerPickup) == nil {
		t.Fatal("hammer pickup tile lookup returned nil")
	}
	if assets.Tile(world.HookPickup) == nil {
		t.Fatal("hook pickup tile lookup returned nil")
	}
	if assets.Tile(world.CompassPickup) == nil {
		t.Fatal("compass pickup tile lookup returned nil")
	}
	if assets.Tile(world.SecretExit) == nil {
		t.Fatal("secret exit tile lookup returned nil")
	}
	if assets.Tile(world.HiddenWall) == nil {
		t.Fatal("hidden wall tile lookup returned nil")
	}
	if assets.Tile(world.TimedSpike) == nil {
		t.Fatal("timed spike tile lookup returned nil")
	}
	if assets.Tile(world.FireTrap) == nil {
		t.Fatal("fire trap tile lookup returned nil")
	}
	if assets.Tile(world.GoldKey) == nil {
		t.Fatal("gold key tile lookup returned nil")
	}
	if assets.Tile(world.GoldDoor) == nil {
		t.Fatal("gold door tile lookup returned nil")
	}
	if assets.EnemyImage(world.EnemyChaser, 0) == nil {
		t.Fatal("enemy frame lookup returned nil")
	}
}
