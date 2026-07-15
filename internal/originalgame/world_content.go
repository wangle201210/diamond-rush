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
	foregroundEffects spriteVisualDefinition
	// Java loads gen chunk 15 (gen1.f #5) for raw 19/43 enemies except in
	// world 1, which loads chunk 17 (gen1.f #7) instead; the frozen variant
	// pairs chunk 16 with 15 and chunk 18 with 17 (i.java:3887-3902,
	// 4157-4166, JAR-verified).
	enemyGreen  spriteVisualDefinition
	enemyRed    spriteVisualDefinition
	enemyFrozen spriteVisualDefinition
}

type spriteVisualDefinition struct {
	sheet    string
	modules  string
	metadata string
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
			foregroundEffects: spriteVisualDefinition{
				sheet:    bavariaForegroundEffectSheet,
				modules:  bavariaForegroundEffectModules,
				metadata: bavariaForegroundEffectMetadata,
			},
			enemyGreen: spriteVisualDefinition{
				sheet:    "decoded/sprites/gen1/chunk07-frames.png",
				modules:  "decoded/sprites/gen1/chunk07-modules.png",
				metadata: "decoded/sprites/gen1/chunk07-animations.json",
			},
			enemyRed: spriteVisualDefinition{
				sheet:    "decoded/sprites/gen1/chunk07-palette01-frames.png",
				modules:  "decoded/sprites/gen1/chunk07-palette01-modules.png",
				metadata: "decoded/sprites/gen1/chunk07-animations.json",
			},
			enemyFrozen: spriteVisualDefinition{
				sheet:    "decoded/sprites/gen1/chunk08-frames.png",
				modules:  "decoded/sprites/gen1/chunk08-modules.png",
				metadata: "decoded/sprites/gen1/chunk08-animations.json",
			},
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
		foregroundEffects: spriteVisualDefinition{
			sheet:    foregroundEffectSheet,
			modules:  foregroundEffectModules,
			metadata: foregroundEffectMetadata,
		},
		enemyGreen: spriteVisualDefinition{
			sheet:    snakeSheet,
			modules:  snakeModuleSheet,
			metadata: snakeMetadata,
		},
		enemyRed: spriteVisualDefinition{
			sheet:    redSnakeSheet,
			modules:  redSnakeModuleSheet,
			metadata: snakeMetadata,
		},
		enemyFrozen: spriteVisualDefinition{
			sheet:    frozenSnakeSheet,
			modules:  frozenSnakeModules,
			metadata: frozenSnakeMetadata,
		},
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
		setWindowTitleForWorld(world)
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
	foregroundEffects, err := loadSpriteSheetWithModules(definition.foregroundEffects.sheet, definition.foregroundEffects.modules, definition.foregroundEffects.metadata)
	if err != nil {
		return fmt.Errorf("load %s foreground effects: %w", worldName(world), err)
	}
	mapHeader, err := loadSpriteSheetWithModules(definition.mapHeaderSheet, definition.mapHeaderModules, definition.mapHeaderMetadata)
	if err != nil {
		return fmt.Errorf("load %s map header: %w", worldName(world), err)
	}
	enemyGreen, err := loadSpriteSheetWithModules(definition.enemyGreen.sheet, definition.enemyGreen.modules, definition.enemyGreen.metadata)
	if err != nil {
		return fmt.Errorf("load %s enemy sprite: %w", worldName(world), err)
	}
	enemyRed, err := loadSpriteSheetWithModules(definition.enemyRed.sheet, definition.enemyRed.modules, definition.enemyRed.metadata)
	if err != nil {
		return fmt.Errorf("load %s red enemy sprite: %w", worldName(world), err)
	}
	enemyFrozen, err := loadSpriteSheetWithModules(definition.enemyFrozen.sheet, definition.enemyFrozen.modules, definition.enemyFrozen.metadata)
	if err != nil {
		return fmt.Errorf("load %s frozen enemy sprite: %w", worldName(world), err)
	}

	g.pack = pack
	g.worldMap = worldMap
	g.worldFrames = frames
	g.boulder = boulder
	g.diggable = diggable
	g.floor = floor
	g.foregroundEffects = foregroundEffects
	g.worldMapHeader = mapHeader
	g.snakes = enemyGreen
	g.redSnakes = enemyRed
	g.frozenSnake = enemyFrozen
	g.worldDir = dir
	g.worldIndex = world
	setWindowTitleForWorld(world)
	return nil
}
