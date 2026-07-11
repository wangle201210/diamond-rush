package original

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	EmptyRawID    RawID = 255
	EntranceRawID RawID = 79
	TileSize            = 24
	ScreenWidth         = 240
	ScreenHeight        = 320
)

type RawID uint8

type Layer string

const (
	PlayerLayer     Layer = "player"
	BackgroundLayer Layer = "background"
	ForegroundLayer Layer = "foreground"
)

type Point struct {
	X int
	Y int
}

type Stage struct {
	Index      int
	Group      int
	GroupIndex int
	Width      int
	Height     int
	Player     []RawID
	Background []RawID
	Foreground []RawID
	Histograms map[Layer]map[RawID]int
	Notes      []string
}

type WorldPack struct {
	Source      string
	InitialByte int
	StageCount  int
	Stages      []*Stage
}

type stageFile struct {
	Index      int                      `json:"index"`
	Group      int                      `json:"group"`
	GroupIndex int                      `json:"groupIndex"`
	Width      int                      `json:"width"`
	Height     int                      `json:"height"`
	Histograms map[string]map[RawID]int `json:"histograms"`
	Layers     map[string][]int         `json:"layers"`
	Notes      []string                 `json:"notes,omitempty"`
}

type manifestFile struct {
	Source      string `json:"source"`
	InitialByte int    `json:"initialByte"`
	StageCount  int    `json:"stageCount"`
	Stages      []struct {
		Index      int    `json:"index"`
		Group      int    `json:"group"`
		GroupIndex int    `json:"groupIndex"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		File       string `json:"file"`
	} `json:"stages"`
}

func LoadWorldDir(dir string) (*WorldPack, error) {
	raw, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var manifest manifestFile
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	pack := &WorldPack{
		Source:      manifest.Source,
		InitialByte: manifest.InitialByte,
		StageCount:  manifest.StageCount,
		Stages:      make([]*Stage, 0, manifest.StageCount),
	}
	for _, entry := range manifest.Stages {
		stage, err := LoadStageFile(filepath.Join(dir, entry.File))
		if err != nil {
			return nil, err
		}
		if stage.Index != entry.Index || stage.Width != entry.Width || stage.Height != entry.Height {
			return nil, fmt.Errorf("stage %s does not match manifest", entry.File)
		}
		pack.Stages = append(pack.Stages, stage)
	}
	if len(pack.Stages) != pack.StageCount {
		return nil, fmt.Errorf("manifest stage count %d, loaded %d", pack.StageCount, len(pack.Stages))
	}
	return pack, nil
}

func LoadStageFile(path string) (*Stage, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read stage %s: %w", path, err)
	}
	var file stageFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return nil, fmt.Errorf("parse stage %s: %w", path, err)
	}
	stage := &Stage{
		Index:      file.Index,
		Group:      file.Group,
		GroupIndex: file.GroupIndex,
		Width:      file.Width,
		Height:     file.Height,
		Player:     rawLayer(file.Layers[string(PlayerLayer)]),
		Background: rawLayer(file.Layers[string(BackgroundLayer)]),
		Foreground: rawLayer(file.Layers[string(ForegroundLayer)]),
		Histograms: map[Layer]map[RawID]int{},
		Notes:      append([]string(nil), file.Notes...),
	}
	for name, hist := range file.Histograms {
		stage.Histograms[Layer(name)] = hist
	}
	if err := stage.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return stage, nil
}

func (s *Stage) Validate() error {
	if s.Width <= 0 || s.Height <= 0 {
		return fmt.Errorf("invalid size %dx%d", s.Width, s.Height)
	}
	want := s.Width * s.Height
	for _, layer := range []Layer{PlayerLayer, BackgroundLayer, ForegroundLayer} {
		if got := len(s.Layer(layer)); got != want {
			return fmt.Errorf("%s layer has %d cells, want %d", layer, got, want)
		}
	}
	return nil
}

func (s *Stage) Layer(layer Layer) []RawID {
	switch layer {
	case PlayerLayer:
		return s.Player
	case BackgroundLayer:
		return s.Background
	case ForegroundLayer:
		return s.Foreground
	default:
		return nil
	}
}

func (s *Stage) At(layer Layer, x, y int) (RawID, bool) {
	if x < 0 || y < 0 || x >= s.Width || y >= s.Height {
		return 0, false
	}
	data := s.Layer(layer)
	if data == nil {
		return 0, false
	}
	return data[x+y*s.Width], true
}

func (s *Stage) Positions(layer Layer, id RawID) []Point {
	data := s.Layer(layer)
	points := make([]Point, 0)
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if data[x+y*s.Width] == id {
				points = append(points, Point{X: x, Y: y})
			}
		}
	}
	return points
}

func (s *Stage) EntranceMarkers() []Point {
	return s.Positions(PlayerLayer, EntranceRawID)
}

func (s *Stage) UniqueIDs(layer Layer) []RawID {
	hist := s.Histograms[layer]
	if hist == nil {
		hist = histogram(s.Layer(layer))
	}
	ids := make([]RawID, 0, len(hist))
	for id := range hist {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func (id RawID) Signed() int8 {
	return int8(id)
}

func (id RawID) IsEmpty() bool {
	return id == EmptyRawID
}

type IDRole struct {
	Layer      Layer
	ID         RawID
	Name       string
	Confidence string
	Source     string
}

func KnownRole(layer Layer, id RawID) (IDRole, bool) {
	role := IDRole{Layer: layer, ID: id, Confidence: "source-anchor"}
	if id == EmptyRawID {
		role.Name = "empty"
		role.Source = "Java signed byte -1"
		return role, true
	}
	switch layer {
	case PlayerLayer:
		switch {
		case id == 0:
			role.Name = "boulder"
			role.Source = "i.java doorHeadClose switch comment: case 0 // Boulder"
			return role, true
		case id == 10:
			role.Name = "diggable tile"
			role.Source = "i.java movement allows raw 10 while xByte == 3; object update clears it and writes foreground raw 32"
			return role, true
		case id == 11:
			role.Name = "crawling enemy"
			role.Source = "i.java stage init seeds raw 11 state; object update dispatches to amVoid() and contact calls hurtHero(1,64,0)"
			return role, true
		case id == 12:
			role.Name = "bonus quota marker"
			role.Source = "i.java stage init clears raw 12, stores x/y in abInt/acInt and background in aaInt"
			return role, true
		case id == 30:
			role.Name = "breakable wall"
			role.Source = "i.java dBoolean excludes raw 30; tool branch advances damage state and clears at >=16"
			return role, true
		case id == EntranceRawID:
			role.Name = "stage entrance marker"
			role.Source = "i.java case 79 during stage init"
			return role, true
		case id >= 80:
			role.Name = "world tile/frame reference"
			role.Source = "i.java render branch draws player IDs >=80 through world frames"
			role.Confidence = "render-anchor"
			return role, true
		case id == 1:
			role.Name = "violet gem"
			role.Source = "i.java dInt comment and doorHeadClose comment: case 1 // Violet gem"
			return role, true
		case id == 2:
			role.Name = "red diamond"
			role.Source = "i.java lVoid(2) uses pickup texture slot 3; pickup completion increments bbInt; menu text names this currency Red diamonds"
			return role, true
		case id == 4:
			role.Name = "key for foreground raw 9"
			role.Source = "i.java pickup branch increments aUInt; foreground raw 9 door branch consumes aUInt"
			return role, true
		case id == 5:
			role.Name = "key for foreground raw 8"
			role.Source = "i.java pickup branch increments aVInt; foreground raw 8 door branch consumes aVInt"
			return role, true
		case id == 6:
			role.Name = "extra life"
			role.Source = "i.java pickup branch increments azInt; death branch decrements azInt"
			return role, true
		case id == 7:
			role.Name = "health refill"
			role.Source = "i.java pickup branch calls aVoid((byte)127), clamped to max health"
			return role, true
		case id == 31:
			role.Name = "keyed lock body"
			role.Source = "i.java render branch draws player raw 31 with foreground animation state"
			role.Confidence = "render-anchor"
			return role, true
		case id == 33:
			role.Name = "passable overlay raw 33"
			role.Source = "i.java stage init preserves player raw 33, movement collision groups it with passable object IDs, and render draws it with aClassfArr[22]"
			return role, true
		case id == 48:
			role.Name = "moving/pullable object candidate"
			role.Source = "i.java hook scan and render branches check player ID 48"
			role.Confidence = "inferred"
			return role, true
		case id == 19:
			role.Name = "green snake"
			role.Source = "i.java doorHeadClose switch comment: case 19 // Green snake; object update calls eVoid((byte)19)"
			return role, true
		case id == 22:
			role.Name = "right horizontal hazard"
			role.Source = "i.java object update case 22 checks same-row player at x..x+animation reach and calls hurtHero(1,64,0)"
			return role, true
		case id == 23:
			role.Name = "left horizontal hazard"
			role.Source = "i.java object update case 23 checks same-row player at x..x-animation reach and calls hurtHero(1,64,0)"
			return role, true
		case id == 24:
			role.Name = "special item pickup raw 24"
			role.Source = "i.java stage init groups raw 24 with mVoid(24); pickup sets iByteArr[9]=1 and bmInt=22"
			return role, true
		case id == 26:
			role.Name = "special item pickup raw 26"
			role.Source = "i.java stage init groups raw 26 with mVoid(26); pickup sets iByteArr[9]=8 and bmInt=25"
			return role, true
		case id == 27:
			role.Name = "special item pickup raw 27"
			role.Source = "i.java stage init groups raw 27 with mVoid(27); pickup sets iByteArr[9]=2 and bmInt=23"
			return role, true
		case id == 41:
			role.Name = "bonus value pickup"
			role.Source = "i.java lVoid(41) stores bIntArrArr value in aHInt, then adds aHInt to aZInt"
			return role, true
		case id == 42:
			role.Name = "special pickup raw 42"
			role.Source = "i.java stage init calls iVoid for raw 42; object update calls lVoid(42), pickup sets pBoolean and bmInt=11"
			return role, true
		case id == 43:
			role.Name = "red snake"
			role.Source = "i.java doorHeadClose switch comment: case 43 // Red snake; object update calls eVoid((byte)43)"
			return role, true
		case id == 53:
			role.Name = "artifact pickup raw 53"
			role.Source = "i.java object update calls lVoid(53); pickup sets DInt=0 and stores bit 0 in iByteArr[2]"
			return role, true
		}
	case ForegroundLayer:
		switch id {
		case 0:
			role.Name = "one-shot foreground event raw 0"
			role.Source = "i.java foreground update case 0 records the merged background/state in bmInt when the player stands on it, then clears the foreground cell"
			return role, true
		case 4:
			role.Name = "checkpoint"
			role.Source = "i.java Action5 and respawn branches check foreground ID 4"
			return role, true
		case 5, 28:
			role.Name = "exit/goal marker candidate"
			role.Source = "i.java stage init stores foreground IDs 5 and 28 in goal bytes"
			role.Confidence = "inferred"
			return role, true
		case 1:
			role.Name = "clearable foreground cluster raw 1"
			role.Source = "i.java movement helper calls bVoid(x,y,(byte)1) when the player enters foreground raw 1, recursively clearing connected raw 1 cells"
			return role, true
		case 2:
			role.Name = "special foreground barrier raw 2"
			role.Source = "i.java foreground update case 2 blocks movement for special-item checks and recursively clears connected raw 2 cells when background 0 has no adjacent raw 30 breakable wall"
			return role, true
		case 6:
			role.Name = "pressure door switch"
			role.Source = "i.java foreground raw 6 opens hVoid(doorI) while pressed and calls doorHeadClose(doorI) when released"
			return role, true
		case 7:
			role.Name = "stateful door"
			role.Source = "i.java stores foreground ID 7 in stage door coordinates; movement blocks it while high state nibble is less than 2"
			return role, true
		case 8:
			role.Name = "lock consuming player raw 5 key"
			role.Source = "i.java foreground raw 8 branch consumes aVInt key counter"
			return role, true
		case 9:
			role.Name = "lock consuming player raw 4 key"
			role.Source = "i.java foreground raw 9 branch consumes aUInt key counter"
			return role, true
		case 14:
			role.Name = "foreground gate overlay raw 14"
			role.Source = "i.java render groups foreground raw 14/33; enemy-gate scan clears its high state when a matching raw17 group reaches zero"
			return role, true
		case 17:
			role.Name = "enemy gate trigger"
			role.Source = "i.java foreground raw 17 groups enemy counters by background state, then opens same-group raw 7 doors when the counter reaches zero"
			return role, true
		case 20, 21, 22, 23, 24, 25:
			role.Name = "animated foreground set"
			role.Source = "i.java render branch draws foreground IDs 20..25 as animations"
			role.Confidence = "render-anchor"
			return role, true
		case 26:
			role.Name = "enemy gate trigger switch"
			role.Source = "i.java foreground raw 26 activates cmInt from its background group, clears itself, and arms same-group raw17 enemy gates"
			return role, true
		case 33:
			role.Name = "foreground gate overlay raw 33"
			role.Source = "i.java render groups foreground raw 14/33; enemy-gate scan clears its high state when a matching raw17 group reaches zero"
			return role, true
		}
	}
	return IDRole{}, false
}

func rawLayer(values []int) []RawID {
	out := make([]RawID, len(values))
	for i, v := range values {
		out[i] = RawID(byte(v))
	}
	return out
}

func histogram(values []RawID) map[RawID]int {
	hist := map[RawID]int{}
	for _, v := range values {
		hist[v]++
	}
	return hist
}
