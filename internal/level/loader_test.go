package level

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBundledLevels(t *testing.T) {
	wantTimings := map[string]struct {
		gravity int
		enemy   int
		chests  int
	}{
		"level01.tmx": {gravity: 10, enemy: 20},
		"level02.tmx": {gravity: 9, enemy: 18},
		"level03.tmx": {gravity: 9, enemy: 17, chests: 1},
		"level04.tmx": {gravity: 8, enemy: 16, chests: 1},
		"level05.tmx": {gravity: 8, enemy: 15, chests: 2},
	}
	for _, name := range []string{"level01.tmx", "level02.tmx", "level03.tmx", "level04.tmx", "level05.tmx"} {
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("..", "..", "assets", "levels", name))
			if err != nil {
				t.Fatal(err)
			}
			def, err := Parse(name, raw)
			if err != nil {
				t.Fatal(err)
			}
			if def.Width != 20 || def.Height != 15 {
				t.Fatalf("size = %dx%d, want 20x15", def.Width, def.Height)
			}
			if def.PlayerStart != (Point{X: 1, Y: 1}) {
				t.Fatalf("player start = %+v, want {1 1}", def.PlayerStart)
			}
			if len(def.Tiles) != def.Width*def.Height {
				t.Fatalf("tiles = %d, want %d", len(def.Tiles), def.Width*def.Height)
			}
			if def.Title == "" || def.Theme == "" || def.Hint == "" {
				t.Fatalf("metadata missing: title=%q theme=%q hint=%q", def.Title, def.Theme, def.Hint)
			}
			if def.ParSteps <= 0 {
				t.Fatalf("par steps = %d, want positive", def.ParSteps)
			}
			if def.GravityTicks != wantTimings[name].gravity || def.EnemyTicks != wantTimings[name].enemy {
				t.Fatalf("timing = gravity %d enemy %d, want gravity %d enemy %d", def.GravityTicks, def.EnemyTicks, wantTimings[name].gravity, wantTimings[name].enemy)
			}
			if len(def.ChestRewards) != wantTimings[name].chests {
				t.Fatalf("chest rewards = %d, want %d", len(def.ChestRewards), wantTimings[name].chests)
			}
			for _, enemy := range def.EnemyStarts {
				if enemy.Direction == "" {
					t.Fatalf("enemy at %+v has no configured direction", enemy.Point)
				}
			}
			for _, chest := range def.ChestRewards {
				if chest.Reward == "" || chest.Amount <= 0 {
					t.Fatalf("invalid chest reward: %+v", chest)
				}
			}
		})
	}
}

func TestParseEnemyDirectionProperty(t *testing.T) {
	raw := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<map orientation="orthogonal" width="3" height="3" tilewidth="32" tileheight="32">
 <layer name="terrain" width="3" height="3">
  <data encoding="csv">
1,1,1,
1,0,1,
1,1,1
  </data>
 </layer>
 <objectgroup name="actors">
  <object id="1" name="Player" type="player" x="32" y="32" width="32" height="32"/>
  <object id="2" name="Snake" type="enemy_horizontal" x="64" y="32" width="32" height="32">
   <properties>
    <property name="direction" value="left"/>
   </properties>
  </object>
 </objectgroup>
</map>`)
	def, err := Parse("direction-test.tmx", raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(def.EnemyStarts) != 1 {
		t.Fatalf("enemy count = %d, want 1", len(def.EnemyStarts))
	}
	if def.EnemyStarts[0].Direction != "left" {
		t.Fatalf("enemy direction = %q, want left", def.EnemyStarts[0].Direction)
	}
}

func TestParseChestRewardObject(t *testing.T) {
	raw := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<map orientation="orthogonal" width="3" height="3" tilewidth="32" tileheight="32">
 <layer name="terrain" width="3" height="3">
  <data encoding="csv">
1,1,1,
1,0,15,
1,1,1
  </data>
 </layer>
 <objectgroup name="actors">
  <object id="1" name="Player" type="player" x="32" y="32" width="32" height="32"/>
 </objectgroup>
 <objectgroup name="containers">
  <object id="2" name="Red Diamond Chest" type="chest" x="64" y="32" width="32" height="32">
   <properties>
    <property name="reward" value="red_diamond"/>
    <property name="amount" value="2"/>
   </properties>
  </object>
 </objectgroup>
</map>`)
	def, err := Parse("chest-test.tmx", raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(def.ChestRewards) != 1 {
		t.Fatalf("chest rewards = %d, want 1", len(def.ChestRewards))
	}
	chest := def.ChestRewards[0]
	if chest.Point != (Point{X: 2, Y: 1}) || chest.Reward != "red_diamond" || chest.Amount != 2 {
		t.Fatalf("chest = %+v", chest)
	}
}

func TestParseBossObject(t *testing.T) {
	raw := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<map orientation="orthogonal" width="3" height="3" tilewidth="32" tileheight="32">
 <layer name="terrain" width="3" height="3">
  <data encoding="csv">
1,1,1,
1,0,0,
1,1,1
  </data>
 </layer>
 <objectgroup name="actors">
  <object id="1" name="Player" type="player" x="32" y="32" width="32" height="32"/>
  <object id="2" name="Guardian" type="boss_guardian" x="64" y="32" width="32" height="32">
   <properties>
    <property name="hp" value="4"/>
   </properties>
  </object>
 </objectgroup>
</map>`)
	def, err := Parse("boss-test.tmx", raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(def.BossStarts) != 1 {
		t.Fatalf("boss count = %d, want 1", len(def.BossStarts))
	}
	if def.BossStarts[0].Point != (Point{X: 2, Y: 1}) || def.BossStarts[0].HP != 4 {
		t.Fatalf("boss = %+v, want point {2 1} hp 4", def.BossStarts[0])
	}
}
