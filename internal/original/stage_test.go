package original

import (
	"path/filepath"
	"testing"
)

func TestLoadAngkorWorldPack(t *testing.T) {
	pack, err := LoadWorldDir(filepath.Join("..", "..", "decoded", "world0"))
	if err != nil {
		t.Fatal(err)
	}
	if pack.StageCount != 14 {
		t.Fatalf("stage count = %d, want 14", pack.StageCount)
	}
	if len(pack.Stages) != 14 {
		t.Fatalf("loaded stages = %d, want 14", len(pack.Stages))
	}
	wantSizes := []Point{
		{26, 21},
		{27, 24},
		{27, 26},
		{40, 23},
		{51, 24},
		{30, 75},
		{26, 45},
		{44, 29},
		{35, 14},
		{50, 30},
		{50, 31},
		{46, 31},
		{46, 31},
		{68, 11},
	}
	for i, want := range wantSizes {
		stage := pack.Stages[i]
		if stage.Width != want.X || stage.Height != want.Y {
			t.Fatalf("stage %02d size = %dx%d, want %dx%d", i, stage.Width, stage.Height, want.X, want.Y)
		}
		if err := stage.Validate(); err != nil {
			t.Fatalf("stage %02d invalid: %v", i, err)
		}
	}
}

func TestAngkorStagesHaveOneEntranceMarker(t *testing.T) {
	pack, err := LoadWorldDir(filepath.Join("..", "..", "decoded", "world0"))
	if err != nil {
		t.Fatal(err)
	}
	for _, stage := range pack.Stages {
		entrances := stage.EntranceMarkers()
		if len(entrances) != 1 {
			t.Fatalf("stage %02d entrance markers = %d, want 1: %v", stage.Index, len(entrances), entrances)
		}
		pt := entrances[0]
		if pt.X < 0 || pt.Y < 0 || pt.X >= stage.Width || pt.Y >= stage.Height {
			t.Fatalf("stage %02d entrance out of bounds: %+v", stage.Index, pt)
		}
	}
}

func TestStage00KnownRawIDs(t *testing.T) {
	stage, err := LoadStageFile(filepath.Join("..", "..", "decoded", "world0", "stage00.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got := stage.Histograms[ForegroundLayer][4]; got != 3 {
		t.Fatalf("stage00 foreground raw 4 count = %d, want 3 checkpoints", got)
	}
	if got := stage.Histograms[ForegroundLayer][5]; got != 1 {
		t.Fatalf("stage00 foreground raw 5 count = %d, want 1 goal marker", got)
	}
	if got := stage.Histograms[PlayerLayer][EntranceRawID]; got != 1 {
		t.Fatalf("stage00 player entrance count = %d, want 1", got)
	}
	if EmptyRawID.Signed() != -1 {
		t.Fatalf("raw 255 signed = %d, want -1", EmptyRawID.Signed())
	}
}

func TestKnownRolesFromJavaAnchors(t *testing.T) {
	tests := []struct {
		layer Layer
		id    RawID
		name  string
	}{
		{PlayerLayer, 0, "boulder"},
		{PlayerLayer, 1, "violet gem"},
		{PlayerLayer, 2, "red diamond"},
		{PlayerLayer, 4, "key for foreground raw 9"},
		{PlayerLayer, 5, "key for foreground raw 8"},
		{PlayerLayer, 6, "extra life"},
		{PlayerLayer, 7, "health refill"},
		{PlayerLayer, 10, "diggable tile"},
		{PlayerLayer, 11, "crawling enemy"},
		{PlayerLayer, 12, "bonus quota marker"},
		{PlayerLayer, 19, "green snake"},
		{PlayerLayer, 22, "right horizontal hazard"},
		{PlayerLayer, 23, "left horizontal hazard"},
		{PlayerLayer, 24, "special item pickup raw 24"},
		{PlayerLayer, 26, "special item pickup raw 26"},
		{PlayerLayer, 27, "special item pickup raw 27"},
		{PlayerLayer, 30, "breakable wall"},
		{PlayerLayer, 31, "keyed lock body"},
		{PlayerLayer, 33, "passable overlay raw 33"},
		{PlayerLayer, 41, "bonus value pickup"},
		{PlayerLayer, 42, "special pickup raw 42"},
		{PlayerLayer, 53, "artifact pickup raw 53"},
		{PlayerLayer, 79, "stage entrance marker"},
		{PlayerLayer, 80, "world tile/frame reference"},
		{ForegroundLayer, 0, "one-shot foreground event raw 0"},
		{ForegroundLayer, 1, "clearable foreground cluster raw 1"},
		{ForegroundLayer, 2, "special foreground barrier raw 2"},
		{ForegroundLayer, 4, "checkpoint"},
		{ForegroundLayer, 6, "pressure door switch"},
		{ForegroundLayer, 7, "stateful door"},
		{ForegroundLayer, 8, "lock consuming player raw 5 key"},
		{ForegroundLayer, 9, "lock consuming player raw 4 key"},
		{ForegroundLayer, 14, "foreground gate overlay raw 14"},
		{ForegroundLayer, 17, "enemy gate trigger"},
		{ForegroundLayer, 20, "animated foreground set"},
		{ForegroundLayer, 26, "enemy gate trigger switch"},
		{ForegroundLayer, 33, "foreground gate overlay raw 33"},
		{BackgroundLayer, 255, "empty"},
	}
	for _, tt := range tests {
		role, ok := KnownRole(tt.layer, tt.id)
		if !ok {
			t.Fatalf("KnownRole(%s,%d) not found", tt.layer, tt.id)
		}
		if role.Name != tt.name {
			t.Fatalf("KnownRole(%s,%d) = %q, want %q", tt.layer, tt.id, role.Name, tt.name)
		}
	}
}
