package level

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type EnemySpawn struct {
	Point     Point
	Type      string
	Direction string
}

type BossSpawn struct {
	Point Point
	Type  string
	HP    int
}

type ChestReward struct {
	Point  Point
	Reward string
	Amount int
}

type Definition struct {
	Name             string
	Title            string
	Theme            string
	Hint             string
	Width            int
	Height           int
	TileWidth        int
	TileHeight       int
	Tiles            []int
	PlayerStart      Point
	EnemyStarts      []EnemySpawn
	BossStarts       []BossSpawn
	ChestRewards     []ChestReward
	RequiredDiamonds int
	ParSteps         int
	GravityTicks     int
	EnemyTicks       int
	HasHammer        bool
	HasHook          bool
}

type tmxMap struct {
	XMLName    xml.Name   `xml:"map"`
	Width      int        `xml:"width,attr"`
	Height     int        `xml:"height,attr"`
	TileWidth  int        `xml:"tilewidth,attr"`
	TileHeight int        `xml:"tileheight,attr"`
	Properties []tmxProp  `xml:"properties>property"`
	Layers     []tmxLayer `xml:"layer"`
	Objects    []tmxGroup `xml:"objectgroup"`
}

type tmxProp struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type tmxLayer struct {
	Name string  `xml:"name,attr"`
	Data tmxData `xml:"data"`
}

type tmxData struct {
	Encoding string `xml:"encoding,attr"`
	Value    string `xml:",chardata"`
}

type tmxGroup struct {
	Name    string      `xml:"name,attr"`
	Objects []tmxObject `xml:"object"`
}

type tmxObject struct {
	Name       string    `xml:"name,attr"`
	Type       string    `xml:"type,attr"`
	X          int       `xml:"x,attr"`
	Y          int       `xml:"y,attr"`
	Properties []tmxProp `xml:"properties>property"`
}

func LoadFile(name string) (*Definition, error) {
	path := "assets/levels/" + name
	raw, err := os.ReadFile(path)
	if err != nil {
		fallback := filepath.Join("..", "..", path)
		raw, err = os.ReadFile(fallback)
		if err != nil {
			return nil, fmt.Errorf("read level %q: %w", path, err)
		}
	}
	return Parse(name, raw)
}

func Parse(name string, raw []byte) (*Definition, error) {
	var parsed tmxMap
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("parse tmx: %w", err)
	}
	if parsed.Width <= 0 || parsed.Height <= 0 {
		return nil, fmt.Errorf("invalid map size %dx%d", parsed.Width, parsed.Height)
	}

	var terrain *tmxLayer
	for i := range parsed.Layers {
		if strings.EqualFold(parsed.Layers[i].Name, "terrain") {
			terrain = &parsed.Layers[i]
			break
		}
	}
	if terrain == nil {
		return nil, fmt.Errorf("level %q has no terrain layer", name)
	}

	tiles, err := parseCSVLayer(terrain.Data.Value)
	if err != nil {
		return nil, err
	}
	if len(tiles) != parsed.Width*parsed.Height {
		return nil, fmt.Errorf("terrain has %d tiles, expected %d", len(tiles), parsed.Width*parsed.Height)
	}

	def := &Definition{
		Name:             name,
		Title:            name,
		Width:            parsed.Width,
		Height:           parsed.Height,
		TileWidth:        parsed.TileWidth,
		TileHeight:       parsed.TileHeight,
		Tiles:            tiles,
		RequiredDiamonds: 8,
		ParSteps:         120,
		GravityTicks:     9,
		EnemyTicks:       18,
	}
	for _, prop := range parsed.Properties {
		switch prop.Name {
		case "required_diamonds":
			v, err := strconv.Atoi(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid required_diamonds %q: %w", prop.Value, err)
			}
			def.RequiredDiamonds = v
		case "par_steps":
			v, err := strconv.Atoi(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid par_steps %q: %w", prop.Value, err)
			}
			def.ParSteps = v
		case "gravity_ticks":
			v, err := strconv.Atoi(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid gravity_ticks %q: %w", prop.Value, err)
			}
			def.GravityTicks = v
		case "enemy_ticks":
			v, err := strconv.Atoi(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid enemy_ticks %q: %w", prop.Value, err)
			}
			def.EnemyTicks = v
		case "has_hook":
			v, err := strconv.ParseBool(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid has_hook %q: %w", prop.Value, err)
			}
			def.HasHook = v
		case "has_hammer":
			v, err := strconv.ParseBool(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid has_hammer %q: %w", prop.Value, err)
			}
			def.HasHammer = v
		case "title":
			def.Title = prop.Value
		case "theme":
			def.Theme = prop.Value
		case "hint":
			def.Hint = prop.Value
		}
	}

	spawnFound := false
	for _, group := range parsed.Objects {
		for _, obj := range group.Objects {
			pt := Point{X: obj.X / parsed.TileWidth, Y: obj.Y / parsed.TileHeight}
			switch strings.ToLower(strings.TrimSpace(obj.Type)) {
			case "player":
				if !strings.EqualFold(group.Name, "actors") {
					continue
				}
				def.PlayerStart = pt
				spawnFound = true
			case "enemy", "enemy_horizontal", "enemy_vertical", "enemy_chaser":
				if !strings.EqualFold(group.Name, "actors") {
					continue
				}
				def.EnemyStarts = append(def.EnemyStarts, EnemySpawn{
					Point:     pt,
					Type:      strings.ToLower(obj.Type),
					Direction: strings.ToLower(propertyValue(obj.Properties, "direction")),
				})
			case "boss", "boss_guardian":
				if !strings.EqualFold(group.Name, "actors") {
					continue
				}
				hp, err := intPropertyValue(obj.Properties, "hp", 3)
				if err != nil {
					return nil, fmt.Errorf("boss at %d,%d: %w", pt.X, pt.Y, err)
				}
				if hp <= 0 {
					return nil, fmt.Errorf("boss at %d,%d has invalid hp %d", pt.X, pt.Y, hp)
				}
				def.BossStarts = append(def.BossStarts, BossSpawn{
					Point: pt,
					Type:  strings.ToLower(obj.Type),
					HP:    hp,
				})
			case "chest":
				amount, err := intPropertyValue(obj.Properties, "amount", 1)
				if err != nil {
					return nil, fmt.Errorf("chest at %d,%d: %w", pt.X, pt.Y, err)
				}
				if amount <= 0 {
					return nil, fmt.Errorf("chest at %d,%d has invalid amount %d", pt.X, pt.Y, amount)
				}
				reward := strings.ToLower(strings.TrimSpace(propertyValue(obj.Properties, "reward")))
				if reward == "" {
					reward = "red_diamond"
				}
				def.ChestRewards = append(def.ChestRewards, ChestReward{
					Point:  pt,
					Reward: reward,
					Amount: amount,
				})
			}
		}
	}
	if !spawnFound {
		return nil, fmt.Errorf("level %q has no player object", name)
	}
	return def, nil
}

func parseCSVLayer(raw string) ([]int, error) {
	reader := csv.NewReader(bytes.NewBufferString(strings.TrimSpace(raw)))
	reader.FieldsPerRecord = -1
	var values []int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse csv layer: %w", err)
		}
		for _, field := range record {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}
			v, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("parse tile gid %q: %w", field, err)
			}
			values = append(values, v)
		}
	}
	return values, nil
}

func propertyValue(properties []tmxProp, name string) string {
	for _, prop := range properties {
		if strings.EqualFold(prop.Name, name) {
			return prop.Value
		}
	}
	return ""
}

func intPropertyValue(properties []tmxProp, name string, fallback int) (int, error) {
	raw := propertyValue(properties, name)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", name, raw, err)
	}
	return value, nil
}
