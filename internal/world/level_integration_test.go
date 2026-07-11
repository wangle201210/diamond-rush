package world

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/wangle201210/zskc/internal/level"
)

type levelAudit struct {
	name         string
	wantRequired int
	wantDiamonds int
	forbidden    []Tile
	wantKeys     int
	wantDoors    int
	wantGoldKeys int
	wantGoldDoor int
	wantRocksMin int
	wantSpikes   int
	wantTimed    int
	wantFire     int
	wantSwitches int
	wantBridges  int
	wantCracked  int
	wantTeleport int
	wantLava     int
	wantChests   int
	wantCheckpts int
	wantPotions  int
	wantHammer   int
	wantHook     int
	wantCompass  int
	wantSecret   int
	wantHidden   int
	wantBosses   int
	wantEnemies  map[EnemyType]int
}

func TestBundledFiveLevelPackDesign(t *testing.T) {
	audits := []levelAudit{
		{name: "level01.tmx", wantRequired: 6, wantDiamonds: 8, wantCheckpts: 1, wantCompass: 1, forbidden: []Tile{Rock, Key, Door, GoldKey, GoldDoor, Spike, Switch, Bridge, CrackedWall, Teleporter, Lava, HammerPickup, HookPickup, SecretExit}},
		{name: "level02.tmx", wantRequired: 7, wantDiamonds: 8, wantCheckpts: 1, forbidden: []Tile{Key, Door, GoldKey, GoldDoor, Spike, Switch, Bridge, CrackedWall, Teleporter, Lava, SecretExit}, wantRocksMin: 5, wantEnemies: map[EnemyType]int{EnemyHorizontal: 1, EnemyVertical: 1}},
		{name: "level03.tmx", wantRequired: 7, wantDiamonds: 7, wantCheckpts: 1, wantPotions: 1, wantHammer: 1, wantChests: 1, forbidden: []Tile{GoldKey, GoldDoor, Switch, Bridge, Teleporter, Lava, HookPickup, SecretExit}, wantKeys: 1, wantDoors: 1, wantSpikes: 4, wantCracked: 1, wantRocksMin: 3, wantEnemies: map[EnemyType]int{EnemyHorizontal: 1, EnemyVertical: 1}},
		{name: "level04.tmx", wantRequired: 7, wantDiamonds: 7, wantKeys: 1, wantDoors: 1, wantSwitches: 1, wantBridges: 3, wantCracked: 1, wantTeleport: 2, wantLava: 3, wantChests: 1, wantCheckpts: 1, wantPotions: 1, wantHook: 1, wantSecret: 1, wantHidden: 1, wantRocksMin: 3, wantEnemies: map[EnemyType]int{EnemyHorizontal: 1, EnemyVertical: 1}},
		{name: "level05.tmx", wantRequired: 8, wantDiamonds: 8, wantGoldKeys: 1, wantGoldDoor: 1, wantSwitches: 1, wantBridges: 2, wantCracked: 1, wantTeleport: 2, wantLava: 4, wantChests: 2, wantCheckpts: 1, wantPotions: 1, wantSpikes: 1, wantTimed: 1, wantFire: 1, wantRocksMin: 4, wantBosses: 1, wantEnemies: map[EnemyType]int{EnemyHorizontal: 1, EnemyVertical: 1}},
	}
	for _, audit := range audits {
		t.Run(audit.name, func(t *testing.T) {
			def := mustLoadBundledWorldDef(t, audit.name)
			w, err := New(def)
			if err != nil {
				t.Fatal(err)
			}
			counts := countTiles(w)
			if def.RequiredDiamonds != audit.wantRequired {
				t.Fatalf("required diamonds = %d, want %d", def.RequiredDiamonds, audit.wantRequired)
			}
			assertTileCount(t, counts, Diamond, audit.wantDiamonds)
			assertTileCount(t, counts, Key, audit.wantKeys)
			assertTileCount(t, counts, Door, audit.wantDoors)
			assertTileCount(t, counts, GoldKey, audit.wantGoldKeys)
			assertTileCount(t, counts, GoldDoor, audit.wantGoldDoor)
			assertTileCount(t, counts, Spike, audit.wantSpikes)
			assertTileCount(t, counts, TimedSpike, audit.wantTimed)
			assertTileCount(t, counts, FireTrap, audit.wantFire)
			assertTileCount(t, counts, Switch, audit.wantSwitches)
			assertTileCount(t, counts, Bridge, audit.wantBridges)
			assertTileCount(t, counts, CrackedWall, audit.wantCracked)
			assertTileCount(t, counts, Teleporter, audit.wantTeleport)
			assertTileCount(t, counts, Lava, audit.wantLava)
			assertTileCount(t, counts, Chest, audit.wantChests)
			assertTileCount(t, counts, Checkpoint, audit.wantCheckpts)
			assertTileCount(t, counts, Potion, audit.wantPotions)
			assertTileCount(t, counts, HammerPickup, audit.wantHammer)
			assertTileCount(t, counts, HookPickup, audit.wantHook)
			assertTileCount(t, counts, CompassPickup, audit.wantCompass)
			assertTileCount(t, counts, SecretExit, audit.wantSecret)
			assertTileCount(t, counts, HiddenWall, audit.wantHidden)
			for _, tile := range audit.forbidden {
				if counts[tile] != 0 {
					t.Fatalf("%s should not appear in %s, found %d", tileName(tile), audit.name, counts[tile])
				}
			}
			if counts[Rock] < audit.wantRocksMin {
				t.Fatalf("rocks = %d, want at least %d", counts[Rock], audit.wantRocksMin)
			}
			if audit.wantTeleport > 0 && counts[Teleporter]%2 != 0 {
				t.Fatalf("teleporters = %d, want paired teleporters", counts[Teleporter])
			}
			if def.RequiredDiamonds > counts[Diamond] {
				t.Fatalf("required diamonds %d exceeds available %d", def.RequiredDiamonds, counts[Diamond])
			}
			if counts[ExitClosed] != 1 {
				t.Fatalf("closed exits = %d, want 1", counts[ExitClosed])
			}
			if len(w.Bosses) != audit.wantBosses {
				t.Fatalf("bosses = %d, want %d", len(w.Bosses), audit.wantBosses)
			}
			if def.ParSteps <= 0 {
				t.Fatalf("par steps = %d, want positive", def.ParSteps)
			}
			enemyCounts := countEnemies(w)
			for enemyType, want := range audit.wantEnemies {
				if enemyCounts[enemyType] != want {
					t.Fatalf("enemy %v count = %d, want %d; all=%v", enemyType, enemyCounts[enemyType], want, enemyCounts)
				}
			}
			if len(enemyCounts) != len(audit.wantEnemies) {
				t.Fatalf("unexpected enemy types: got=%v want=%v", enemyCounts, audit.wantEnemies)
			}
		})
	}
}

func TestBundledLevelSpikeTilesReachWorld(t *testing.T) {
	def := mustLoadBundledWorldDef(t, "level03.tmx")
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.TileAt(x, y) == Spike {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("level03 has no spike tiles in world")
	}
}

func TestBundledLevelSwitchTilesReachWorld(t *testing.T) {
	def := mustLoadBundledWorldDef(t, "level04.tmx")
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	var switches, bridges, cracked int
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			switch w.TileAt(x, y) {
			case Switch:
				switches++
			case Bridge:
				bridges++
			case CrackedWall:
				cracked++
			}
		}
	}
	if switches == 0 || bridges == 0 || cracked == 0 {
		t.Fatalf("level04 mechanic tiles: switches=%d bridges=%d cracked=%d", switches, bridges, cracked)
	}
}

func TestBundledLevelTeleporterAndLavaTilesReachWorld(t *testing.T) {
	def := mustLoadBundledWorldDef(t, "level04.tmx")
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	var teleporters, lava int
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			switch w.TileAt(x, y) {
			case Teleporter:
				teleporters++
			case Lava:
				lava++
			}
		}
	}
	if teleporters < 2 || lava == 0 {
		t.Fatalf("level04 mechanic tiles: teleporters=%d lava=%d", teleporters, lava)
	}
}

func TestBundledLevelFiveHasBossFinale(t *testing.T) {
	def := mustLoadBundledWorldDef(t, "level05.tmx")
	w, err := New(def)
	if err != nil {
		t.Fatal(err)
	}
	if len(w.Bosses) != 1 {
		t.Fatalf("level05 bosses = %d, want 1", len(w.Bosses))
	}
	found := map[EnemyType]bool{}
	for _, enemy := range w.Enemies {
		found[enemy.Type] = true
	}
	for _, enemyType := range []EnemyType{EnemyHorizontal, EnemyVertical} {
		if !found[enemyType] {
			t.Fatalf("level05 missing enemy type %v; found=%v", enemyType, found)
		}
	}
}

func mustLoadBundledWorldDef(t *testing.T, name string) *level.Definition {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "assets", "levels", name))
	if err != nil {
		t.Fatal(err)
	}
	def, err := level.Parse(name, raw)
	if err != nil {
		t.Fatal(err)
	}
	return def
}

func countTiles(w *World) map[Tile]int {
	counts := map[Tile]int{}
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			counts[w.TileAt(x, y)]++
		}
	}
	return counts
}

func countEnemies(w *World) map[EnemyType]int {
	counts := map[EnemyType]int{}
	for _, enemy := range w.Enemies {
		counts[enemy.Type]++
	}
	return counts
}

func assertTileCount(t *testing.T, counts map[Tile]int, tile Tile, want int) {
	t.Helper()
	if counts[tile] != want {
		t.Fatalf("%s count = %d, want %d", tileName(tile), counts[tile], want)
	}
}

func tileName(tile Tile) string {
	names := map[Tile]string{
		Empty:         "empty",
		Wall:          "wall",
		Dirt:          "dirt",
		Diamond:       "diamond",
		Rock:          "rock",
		ExitClosed:    "exit_closed",
		ExitOpen:      "exit_open",
		Key:           "key",
		Door:          "door",
		GoldKey:       "gold_key",
		GoldDoor:      "gold_door",
		Spike:         "spike",
		TimedSpike:    "timed_spike",
		FireTrap:      "fire_trap",
		Switch:        "switch",
		Bridge:        "bridge",
		CrackedWall:   "cracked_wall",
		Teleporter:    "teleporter",
		Lava:          "lava",
		Chest:         "chest",
		Checkpoint:    "checkpoint",
		Potion:        "potion",
		HammerPickup:  "hammer_pickup",
		HookPickup:    "hook_pickup",
		CompassPickup: "compass_pickup",
		SecretExit:    "secret_exit",
		HiddenWall:    "hidden_wall",
	}
	if name, ok := names[tile]; ok {
		return name
	}
	return fmt.Sprintf("tile_%d", tile)
}
