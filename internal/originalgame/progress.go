package originalgame

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"sort"

	"github.com/wangle201210/zskc/internal/original"
)

const (
	originalProgressVersion = 6
	angkorStageCount        = 14
	angkorReplicaStageCount = 14
	angkorFirstSecretStage  = 9
	angkorSealStage         = 8
)

type originalProgress struct {
	Version              int                                `json:"version"`
	HighestUnlocked      int                                `json:"highest_unlocked"`
	StageUnlocked        [angkorStageCount]bool             `json:"stage_unlocked"`
	StageCleared         [angkorStageCount]bool             `json:"stage_cleared"`
	StageAwards          [angkorStageCount]byte             `json:"stage_awards"`
	StageVioletGems      [angkorStageCount]int              `json:"stage_violet_gems"`
	StageRedDiamonds     [angkorStageCount]int              `json:"stage_red_diamonds"`
	StageConsumedRewards [angkorStageCount][]original.Point `json:"stage_consumed_rewards"`
	VioletGemBank        int                                `json:"violet_gem_bank"`
	RedDiamondBank       int                                `json:"red_diamond_bank"`
	ExtraLives           int                                `json:"extra_lives"`
	MaxHealth            int                                `json:"max_health"`
	ToolLevel            int                                `json:"tool_level"`
	TutorialComplete     bool                               `json:"tutorial_complete"`
	RelicMask            int                                `json:"relic_mask"`
	WorldUnlocked        [3]bool                            `json:"world_unlocked"`
}

// originalProgressDisk keeps the two Stage-1-only fields readable so saves
// created before the Angkor campaign flow was implemented migrate in place.
type originalProgressDisk struct {
	Version              int                                `json:"version"`
	HighestUnlocked      int                                `json:"highest_unlocked"`
	StageUnlocked        [angkorStageCount]bool             `json:"stage_unlocked"`
	StageCleared         [angkorStageCount]bool             `json:"stage_cleared"`
	StageAwards          [angkorStageCount]byte             `json:"stage_awards"`
	StageVioletGems      [angkorStageCount]int              `json:"stage_violet_gems"`
	StageRedDiamonds     [angkorStageCount]int              `json:"stage_red_diamonds"`
	StageConsumedRewards [angkorStageCount][]original.Point `json:"stage_consumed_rewards"`
	VioletGemBank        int                                `json:"violet_gem_bank"`
	RedDiamondBank       int                                `json:"red_diamond_bank"`
	ExtraLives           int                                `json:"extra_lives"`
	MaxHealth            int                                `json:"max_health"`
	ToolLevel            int                                `json:"tool_level"`
	TutorialComplete     bool                               `json:"tutorial_complete"`
	RelicMask            int                                `json:"relic_mask"`
	WorldUnlocked        [3]bool                            `json:"world_unlocked"`
	Stage1Cleared        bool                               `json:"stage1_cleared"`
	Stage1Awards         byte                               `json:"stage1_awards"`
}

func newOriginalProgress() originalProgress {
	progress := originalProgress{
		Version:         originalProgressVersion,
		HighestUnlocked: 0,
		ExtraLives:      5,
		MaxHealth:       4,
	}
	progress.StageUnlocked[0] = true
	progress.WorldUnlocked[0] = true
	return progress
}

func (p originalProgress) normalized() originalProgress {
	p.Version = originalProgressVersion
	p.HighestUnlocked = clamp(p.HighestUnlocked, 0, angkorStageCount-1)
	p.StageUnlocked[0] = true
	p.WorldUnlocked[0] = true
	p.RelicMask &= 0x07
	for stage, unlocked := range p.StageUnlocked {
		if unlocked && stage > p.HighestUnlocked {
			p.HighestUnlocked = stage
		}
	}
	for stage := range p.StageConsumedRewards {
		seen := make(map[original.Point]bool, len(p.StageConsumedRewards[stage]))
		points := p.StageConsumedRewards[stage][:0]
		for _, point := range p.StageConsumedRewards[stage] {
			if point.X < 0 || point.X > 255 || point.Y < 0 || point.Y > 255 || seen[point] {
				continue
			}
			seen[point] = true
			points = append(points, point)
		}
		sort.Slice(points, func(i, j int) bool {
			if points[i].Y != points[j].Y {
				return points[i].Y < points[j].Y
			}
			return points[i].X < points[j].X
		})
		p.StageConsumedRewards[stage] = points
	}
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

func (p originalProgress) stageUnlocked(stage int) bool {
	return stage >= 0 && stage < angkorStageCount && p.StageUnlocked[stage]
}

func (p *originalProgress) unlockStage(stage int) {
	if p == nil || stage < 0 || stage >= angkorStageCount {
		return
	}
	p.StageUnlocked[stage] = true
	p.HighestUnlocked = max(p.HighestUnlocked, stage)
}

func (p *originalProgress) recordStageResult(stageIndex int, rt *original.Runtime) byte {
	if p == nil || stageIndex < 0 || stageIndex >= angkorStageCount || rt == nil {
		return 0
	}
	awards := stageResultAwards(rt)
	newAwards := awards &^ p.StageAwards[stageIndex]
	p.StageAwards[stageIndex] |= awards
	p.StageCleared[stageIndex] = true
	// Angkor's normal route is stages 0..8. Secret stages are unlocked only
	// through foreground raw 28 and must not be inferred from a high index.
	if stageIndex < angkorFirstSecretStage-1 {
		p.unlockStage(stageIndex + 1)
	}
	p.recordStageCollections(stageIndex, rt, newAwards)
	*p = p.normalized()
	return newAwards
}

func (p *originalProgress) recordSecretExit(stageIndex, targetStage int, rt *original.Runtime) {
	if p == nil || stageIndex < 0 || stageIndex >= angkorStageCount || rt == nil {
		return
	}
	if stageIndex >= angkorFirstSecretStage {
		p.StageCleared[stageIndex] = true
	}
	p.unlockStage(stageIndex)
	p.unlockStage(targetStage)
	p.recordStageCollections(stageIndex, rt, 0)
	*p = p.normalized()
}

func (p *originalProgress) recordSecretStageResult(stageIndex, targetStage int, rt *original.Runtime) byte {
	newAwards := p.recordStageResult(stageIndex, rt)
	p.unlockStage(targetStage)
	*p = p.normalized()
	return newAwards
}

func (p *originalProgress) recordSealStageCompletion(stageIndex int, rt *original.Runtime) {
	if p == nil || stageIndex < 0 || stageIndex >= angkorStageCount || rt == nil {
		return
	}
	p.StageCleared[stageIndex] = true
	p.unlockStage(stageIndex)
	p.recordStageCollections(stageIndex, rt, 0)
	p.RelicMask |= rt.RelicMask & 0x07
	*p = p.normalized()
}

func (p *originalProgress) recordStageCollections(stageIndex int, rt *original.Runtime, newAwards byte) {
	p.VioletGemBank += max(0, rt.VioletGems)
	p.StageVioletGems[stageIndex] = max(p.StageVioletGems[stageIndex], rt.VioletGems)
	p.RedDiamondBank += max(0, rt.RedDiamonds)
	p.StageRedDiamonds[stageIndex] = min(rt.TotalRedDiamonds, p.StageRedDiamonds[stageIndex]+max(0, rt.RedDiamonds))
	p.ExtraLives = min(99, rt.ExtraLives+bits.OnesCount8(uint8(newAwards)))
	p.MaxHealth = max(4, rt.MaxHealth)
	p.ToolLevel = maxToolLevel(p.ToolLevel, specialItemMaskToolLevel(rt.SpecialItemMask))
	for _, point := range rt.PersistentRewardCoordinates() {
		if !containsPoint(p.StageConsumedRewards[stageIndex], point) {
			p.StageConsumedRewards[stageIndex] = append(p.StageConsumedRewards[stageIndex], point)
		}
	}
}

func containsPoint(points []original.Point, target original.Point) bool {
	for _, point := range points {
		if point == target {
			return true
		}
	}
	return false
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

func originalProgressExists(path string) (bool, error) {
	_, err := os.Stat(filepath.Clean(path))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
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
		Version:              disk.Version,
		HighestUnlocked:      disk.HighestUnlocked,
		StageUnlocked:        disk.StageUnlocked,
		StageCleared:         disk.StageCleared,
		StageAwards:          disk.StageAwards,
		StageVioletGems:      disk.StageVioletGems,
		StageRedDiamonds:     disk.StageRedDiamonds,
		StageConsumedRewards: disk.StageConsumedRewards,
		VioletGemBank:        disk.VioletGemBank,
		RedDiamondBank:       disk.RedDiamondBank,
		ExtraLives:           disk.ExtraLives,
		MaxHealth:            disk.MaxHealth,
		ToolLevel:            disk.ToolLevel,
		TutorialComplete:     disk.TutorialComplete,
		RelicMask:            disk.RelicMask,
		WorldUnlocked:        disk.WorldUnlocked,
	}
	if disk.Version < 2 {
		progress = newOriginalProgress()
		progress.StageCleared[0] = disk.Stage1Cleared
		progress.StageAwards[0] = disk.Stage1Awards
		if disk.Stage1Cleared {
			progress.HighestUnlocked = 1
			progress.StageUnlocked[1] = true
		}
	} else if disk.Version < 3 {
		// Version 2 used HighestUnlocked as a sequential range. Version 3 uses
		// explicit node bits so a raw-28 jump to stage 9 does not unlock 7 and 8.
		for stage := 0; stage <= clamp(disk.HighestUnlocked, 0, angkorStageCount-1); stage++ {
			progress.StageUnlocked[stage] = true
		}
	}
	if disk.Version < 4 {
		progress.TutorialComplete = true
	}
	if disk.Version < 5 && progress.StageCleared[8] {
		progress.RelicMask |= 1
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
