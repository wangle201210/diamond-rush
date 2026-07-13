package originalgame

import (
	"fmt"
	"path/filepath"

	"github.com/wangle201210/zskc/internal/original"
)

type worldVisualDefinition struct {
	frames            string
	frameModules      string
	frameMetadata     string
	boulder           string
	boulderModules    string
	boulderMetadata   string
	diggable          string
	diggableModules   string
	diggableMetadata  string
	floor             string
	floorMetadata     string
	mapHeaderSheet    string
	mapHeaderModules  string
	mapHeaderMetadata string
}

func worldVisualDefinitionFor(world int) worldVisualDefinition {
	if world == original.WorldBavaria {
		return worldVisualDefinition{
			frames:            bavariaWorldFrameSheet,
			frameModules:      bavariaWorldFrameModules,
			frameMetadata:     bavariaWorldFrameMetadata,
			boulder:           bavariaBoulderFrameSheet,
			boulderModules:    bavariaBoulderModules,
			boulderMetadata:   bavariaBoulderMetadata,
			diggable:          bavariaDiggableFrameSheet,
			diggableModules:   bavariaDiggableModules,
			diggableMetadata:  bavariaDiggableMetadata,
			floor:             bavariaFloorSheet,
			floorMetadata:     bavariaFloorMetadata,
			mapHeaderSheet:    bavariaMapHeaderSheet,
			mapHeaderModules:  bavariaMapHeaderModules,
			mapHeaderMetadata: bavariaMapHeaderMetadata,
		}
	}
	return worldVisualDefinition{
		frames:            angkorWorldFrameSheet,
		frameModules:      angkorWorldFrameModules,
		frameMetadata:     angkorWorldFrameMetadata,
		boulder:           angkorBoulderFrameSheet,
		boulderModules:    angkorBoulderModules,
		boulderMetadata:   angkorBoulderMetadata,
		diggable:          angkorDiggableFrameSheet,
		diggableModules:   angkorDiggableModules,
		diggableMetadata:  angkorDiggableMetadata,
		floor:             angkorFloorSheet,
		floorMetadata:     angkorFloorMetadata,
		mapHeaderSheet:    worldMapHeaderSheet,
		mapHeaderModules:  worldMapHeaderModules,
		mapHeaderMetadata: worldMapHeaderMetadata,
	}
}

func worldName(world int) string {
	switch world {
	case original.WorldBavaria:
		return "Bavaria"
	case original.WorldTibet:
		return "Siberia"
	default:
		return "Angkor Wat"
	}
}

func worldStageCount(world int) int {
	if world == original.WorldBavaria {
		return bavariaStageCount
	}
	return angkorStageCount
}

func worldFirstSecretStage(world int) int {
	if world == original.WorldBavaria {
		return bavariaFirstSecretStage
	}
	return angkorFirstSecretStage
}

func worldSealStage(world int) int {
	if world == original.WorldBavaria {
		return bavariaSealStage
	}
	return angkorSealStage
}

func worldMusic(world int) int {
	if world == original.WorldBavaria {
		return original.SoundBavariaMusic
	}
	return original.SoundAngkorMusic
}

func worldStageImplemented(world, stage int) bool {
	return stage >= 0 && stage < worldStageCount(world)
}

func (g *Game) switchWorld(world int) error {
	if g == nil {
		return fmt.Errorf("nil game")
	}
	if world != original.WorldAngkor && world != original.WorldBavaria {
		return fmt.Errorf("world %d is not replicated", world)
	}
	if g.pack != nil && g.pack.World == world {
		g.worldIndex = world
		return nil
	}
	dir := filepath.Join(g.worldRoot, fmt.Sprintf("world%d", world))
	pack, err := original.LoadWorldDir(dir)
	if err != nil {
		return fmt.Errorf("load %s stages: %w", worldName(world), err)
	}
	worldMap, err := loadWorldMap(filepath.Join(dir, "map.json"))
	if err != nil {
		return fmt.Errorf("load %s map: %w", worldName(world), err)
	}
	definition := worldVisualDefinitionFor(world)
	frames, err := loadSpriteSheetWithModules(definition.frames, definition.frameModules, definition.frameMetadata)
	if err != nil {
		return fmt.Errorf("load %s world frames: %w", worldName(world), err)
	}
	boulder, err := loadSpriteSheetWithModules(definition.boulder, definition.boulderModules, definition.boulderMetadata)
	if err != nil {
		return fmt.Errorf("load %s boulder: %w", worldName(world), err)
	}
	diggable, err := loadSpriteSheetWithModules(definition.diggable, definition.diggableModules, definition.diggableMetadata)
	if err != nil {
		return fmt.Errorf("load %s diggable tile: %w", worldName(world), err)
	}
	floor, err := loadModuleSpriteSheet(definition.floor, definition.floorMetadata)
	if err != nil {
		return fmt.Errorf("load %s floor: %w", worldName(world), err)
	}
	mapHeader, err := loadSpriteSheetWithModules(definition.mapHeaderSheet, definition.mapHeaderModules, definition.mapHeaderMetadata)
	if err != nil {
		return fmt.Errorf("load %s map header: %w", worldName(world), err)
	}

	g.pack = pack
	g.worldMap = worldMap
	g.worldFrames = frames
	g.boulder = boulder
	g.diggable = diggable
	g.floor = floor
	g.worldMapHeader = mapHeader
	g.worldDir = dir
	g.worldIndex = world
	return nil
}
