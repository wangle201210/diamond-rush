package originalgame

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path/filepath"

	"github.com/wangle201210/zskc/internal/original"
)

const (
	originalProgressVersion = 2
	angkorStageCount        = 14
	angkorReplicaStageCount = 5
)

type originalProgress struct {
	Version          int                    `json:"version"`
	HighestUnlocked  int                    `json:"highest_unlocked"`
	StageCleared     [angkorStageCount]bool `json:"stage_cleared"`
	StageAwards      [angkorStageCount]byte `json:"stage_awards"`
	StageVioletGems  [angkorStageCount]int  `json:"stage_violet_gems"`
	StageRedDiamonds [angkorStageCount]int  `json:"stage_red_diamonds"`
	VioletGemBank    int                    `json:"violet_gem_bank"`
	RedDiamondBank   int                    `json:"red_diamond_bank"`
	ExtraLives       int                    `json:"extra_lives"`
	MaxHealth        int                    `json:"max_health"`
	ToolLevel        int                    `json:"tool_level"`
}

// originalProgressDisk keeps the two Stage-1-only fields readable so saves
// created before the Angkor campaign flow was implemented migrate in place.
type originalProgressDisk struct {
	Version          int                    `json:"version"`
	HighestUnlocked  int                    `json:"highest_unlocked"`
	StageCleared     [angkorStageCount]bool `json:"stage_cleared"`
	StageAwards      [angkorStageCount]byte `json:"stage_awards"`
	StageVioletGems  [angkorStageCount]int  `json:"stage_violet_gems"`
	StageRedDiamonds [angkorStageCount]int  `json:"stage_red_diamonds"`
	VioletGemBank    int                    `json:"violet_gem_bank"`
	RedDiamondBank   int                    `json:"red_diamond_bank"`
	ExtraLives       int                    `json:"extra_lives"`
	MaxHealth        int                    `json:"max_health"`
	ToolLevel        int                    `json:"tool_level"`
	Stage1Cleared    bool                   `json:"stage1_cleared"`
	Stage1Awards     byte                   `json:"stage1_awards"`
}

func newOriginalProgress() originalProgress {
	return originalProgress{
		Version:         originalProgressVersion,
		HighestUnlocked: 0,
		ExtraLives:      5,
		MaxHealth:       4,
	}
}

func (p originalProgress) normalized() originalProgress {
	p.Version = originalProgressVersion
	p.HighestUnlocked = clamp(p.HighestUnlocked, 0, angkorReplicaStageCount-1)
	if p.MaxHealth < 4 {
		p.MaxHealth = 4
	}
	if p.ExtraLives > 99 {
		p.ExtraLives = 99
	}
	switch p.ToolLevel {
	case 0, 1, 2, 8:
	default:
		p.ToolLevel = 0
	}
	return p
}

func (p *originalProgress) recordStageResult(stageIndex int, rt *original.Runtime) byte {
	if p == nil || stageIndex < 0 || stageIndex >= angkorStageCount || rt == nil {
		return 0
	}
	awards := stageResultAwards(rt)
	newAwards := awards &^ p.StageAwards[stageIndex]
	p.StageAwards[stageIndex] |= awards
	p.StageCleared[stageIndex] = true
	if stageIndex+1 < angkorReplicaStageCount && p.HighestUnlocked < stageIndex+1 {
		p.HighestUnlocked = stageIndex + 1
	}

	if collected := rt.VioletGems; collected > p.StageVioletGems[stageIndex] {
		p.VioletGemBank += collected - p.StageVioletGems[stageIndex]
		p.StageVioletGems[stageIndex] = collected
	}
	if collected := rt.RedDiamonds; collected > p.StageRedDiamonds[stageIndex] {
		p.RedDiamondBank += collected - p.StageRedDiamonds[stageIndex]
		p.StageRedDiamonds[stageIndex] = collected
	}
	p.ExtraLives = min(99, rt.ExtraLives+bits.OnesCount8(uint8(newAwards)))
	p.MaxHealth = max(4, rt.MaxHealth)
	p.ToolLevel = maxToolLevel(p.ToolLevel, specialItemMaskToolLevel(rt.SpecialItemMask))
	*p = p.normalized()
	return newAwards
}

func specialItemMaskToolLevel(mask int) int {
	switch {
	case mask&8 != 0:
		return 8
	case mask&2 != 0:
		return 2
	case mask&1 != 0:
		return 1
	default:
		return 0
	}
}

func maxToolLevel(a, b int) int {
	order := func(level int) int {
		switch level {
		case 1:
			return 1
		case 2:
			return 2
		case 8:
			return 3
		default:
			return 0
		}
	}
	if order(b) > order(a) {
		return b
	}
	return a
}

func originalProgressPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		return filepath.Join(".", "diamondrush-original-progress.json")
	}
	return filepath.Join(configDir, "zskc-diamondrush", "original-progress.json")
}

func loadOriginalProgress(path string) (originalProgress, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if errors.Is(err, os.ErrNotExist) {
		return newOriginalProgress(), nil
	}
	if err != nil {
		return originalProgress{}, err
	}
	var disk originalProgressDisk
	if err := json.Unmarshal(data, &disk); err != nil {
		return originalProgress{}, fmt.Errorf("decode original progress: %w", err)
	}
	progress := originalProgress{
		Version:          disk.Version,
		HighestUnlocked:  disk.HighestUnlocked,
		StageCleared:     disk.StageCleared,
		StageAwards:      disk.StageAwards,
		StageVioletGems:  disk.StageVioletGems,
		StageRedDiamonds: disk.StageRedDiamonds,
		VioletGemBank:    disk.VioletGemBank,
		RedDiamondBank:   disk.RedDiamondBank,
		ExtraLives:       disk.ExtraLives,
		MaxHealth:        disk.MaxHealth,
		ToolLevel:        disk.ToolLevel,
	}
	if disk.Version < originalProgressVersion {
		progress = newOriginalProgress()
		progress.StageCleared[0] = disk.Stage1Cleared
		progress.StageAwards[0] = disk.Stage1Awards
		if disk.Stage1Cleared {
			progress.HighestUnlocked = 1
		}
	}
	return progress.normalized(), nil
}

func saveOriginalProgress(path string, progress originalProgress) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	progress = progress.normalized()
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(temporary, path); err != nil {
		_ = os.Remove(temporary)
		return err
	}
	return nil
}
