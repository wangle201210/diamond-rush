package game

import (
	"os"
	"regexp"
	"testing"

	"github.com/wangle201210/zskc/internal/world"
)

func TestDirectionFromVirtualPad(t *testing.T) {
	width, height := viewportPixelWidth(), viewportPixelHeight()+hudHeight
	rect := virtualDPadRect(width, height)
	button := rect.Dx() / 3
	tests := []struct {
		name string
		x    int
		y    int
		want world.Direction
	}{
		{"up", rect.Min.X + button + button/2, rect.Min.Y + button/2, world.Up},
		{"down", rect.Min.X + button + button/2, rect.Min.Y + button*2 + button/2, world.Down},
		{"left", rect.Min.X + button/2, rect.Min.Y + button + button/2, world.Left},
		{"right", rect.Min.X + button*2 + button/2, rect.Min.Y + button + button/2, world.Right},
		{"center", rect.Min.X + button + button/2, rect.Min.Y + button + button/2, world.Direction{}},
		{"outside", rect.Min.X - 1, rect.Min.Y + button, world.Direction{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := directionFromVirtualPad(tt.x, tt.y, width, height); got != tt.want {
				t.Fatalf("directionFromVirtualPad() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestHookFromVirtualPad(t *testing.T) {
	width, height := viewportPixelWidth(), viewportPixelHeight()+hudHeight
	rect := virtualDPadRect(width, height)
	button := rect.Dx() / 3
	tests := []struct {
		name string
		x    int
		y    int
		want bool
	}{
		{"center", rect.Min.X + button + button/2, rect.Min.Y + button + button/2, true},
		{"up", rect.Min.X + button + button/2, rect.Min.Y + button/2, false},
		{"outside", rect.Min.X - 1, rect.Min.Y + button, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hookFromVirtualPad(tt.x, tt.y, width, height); got != tt.want {
				t.Fatalf("hookFromVirtualPad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCameraOriginFollowsPlayerAndClampsToMap(t *testing.T) {
	tests := []struct {
		name   string
		player world.Player
		wantX  int
		wantY  int
	}{
		{name: "top left", player: world.Player{X: 1, Y: 1}, wantX: 0, wantY: 0},
		{name: "center", player: world.Player{X: 10, Y: 7}, wantX: 3, wantY: 2},
		{name: "bottom right", player: world.Player{X: 18, Y: 13}, wantX: 5, wantY: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{world: &world.World{
				Width:  20,
				Height: 15,
				Player: tt.player,
			}}
			gotX, gotY := g.cameraOrigin()
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Fatalf("cameraOrigin() = (%d,%d), want (%d,%d)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestLayoutUsesPhoneStyleViewport(t *testing.T) {
	g := &Game{world: &world.World{Width: 20, Height: 15}}
	gotW, gotH := g.Layout(0, 0)
	wantW := viewportTilesX * tileSize
	wantH := viewportTilesY*tileSize + hudHeight
	if gotW != wantW || gotH != wantH {
		t.Fatalf("Layout() = (%d,%d), want (%d,%d)", gotW, gotH, wantW, wantH)
	}
}

func TestTitleMenuItemsMatchOriginalPhoneMenuShape(t *testing.T) {
	got := titleMenuItems()
	want := []string{"Continue", "New Game", "Level Map", "Options", "Help", "About", "Exit"}
	if len(got) != len(want) {
		t.Fatalf("titleMenuItems len = %d, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("titleMenuItems[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestActivateTitleHelpAndAboutUseMenuMessages(t *testing.T) {
	g := &Game{levels: []string{"one"}, progress: defaultProgress(1)}
	g.titleSelected = 4
	if err := g.activateTitleMenu(); err != nil {
		t.Fatal(err)
	}
	if g.message != "Help: 2/4/6/8 move, 5 action, * recall." {
		t.Fatalf("help message = %q", g.message)
	}
	g.titleSelected = 5
	if err := g.activateTitleMenu(); err != nil {
		t.Fatal(err)
	}
	if g.message != "About: Angkor five-stage Diamond Rush remake." {
		t.Fatalf("about message = %q", g.message)
	}
}

func TestNoTextLetterGameplayShortcuts(t *testing.T) {
	source, err := os.ReadFile("game.go")
	if err != nil {
		t.Fatal(err)
	}
	forbidden := regexp.MustCompile(`ebiten\.Key(?:W|A|S|D|U|I|O|M|P|R|L|N)\b`)
	if match := forbidden.Find(source); match != nil {
		t.Fatalf("text letter shortcut %q found; use phone-number keys or non-text keys to avoid IME interference", match)
	}
}
