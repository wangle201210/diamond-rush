package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Progress struct {
	UnlockedLevel     int    `json:"unlocked_level"`
	BestSteps         []int  `json:"best_steps"`
	BestScores        []int  `json:"best_scores"`
	RedDiamonds       []int  `json:"red_diamonds"`
	PurpleDiamonds    []int  `json:"purple_diamonds"`
	SecretExits       []bool `json:"secret_exits"`
	AllPurpleClears   []bool `json:"all_purple_clears"`
	AllRedClears      []bool `json:"all_red_clears"`
	NoDamageClears    []bool `json:"no_damage_clears"`
	NoRecallClears    []bool `json:"no_recall_clears"`
	NoRestartClears   []bool `json:"no_restart_clears"`
	PurpleBank        int    `json:"purple_bank"`
	MaxHealthUpgrades int    `json:"max_health_upgrades"`
	ArmorUpgrades     int    `json:"armor_upgrades"`
	LifeUpgrades      int    `json:"life_upgrades"`
	HasCompass        bool   `json:"has_compass"`
	HasHammer         bool   `json:"has_hammer"`
	HasHook           bool   `json:"has_hook"`
	AncientSealOpen   bool   `json:"ancient_seal_open"`
}

func defaultProgress(levelCount int) Progress {
	return Progress{
		UnlockedLevel:   1,
		BestSteps:       make([]int, levelCount),
		BestScores:      make([]int, levelCount),
		RedDiamonds:     make([]int, levelCount),
		PurpleDiamonds:  make([]int, levelCount),
		SecretExits:     make([]bool, levelCount),
		AllPurpleClears: make([]bool, levelCount),
		AllRedClears:    make([]bool, levelCount),
		NoDamageClears:  make([]bool, levelCount),
		NoRecallClears:  make([]bool, levelCount),
		NoRestartClears: make([]bool, levelCount),
	}
}

func loadProgress(path string, levelCount int) (Progress, error) {
	progress := defaultProgress(levelCount)
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return progress, nil
	}
	if err != nil {
		return progress, fmt.Errorf("read progress: %w", err)
	}
	if err := json.Unmarshal(raw, &progress); err != nil {
		return defaultProgress(levelCount), nil
	}
	normalizeProgress(&progress, levelCount)
	return progress, nil
}

func saveProgress(path string, progress Progress) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create progress dir: %w", err)
	}
	raw, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return fmt.Errorf("encode progress: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return fmt.Errorf("write progress: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit progress: %w", err)
	}
	return nil
}

func normalizeProgress(progress *Progress, levelCount int) {
	if progress.UnlockedLevel < 1 {
		progress.UnlockedLevel = 1
	}
	if progress.UnlockedLevel > levelCount {
		progress.UnlockedLevel = levelCount
	}
	for len(progress.BestSteps) < levelCount {
		progress.BestSteps = append(progress.BestSteps, 0)
	}
	if len(progress.BestSteps) > levelCount {
		progress.BestSteps = progress.BestSteps[:levelCount]
	}
	for len(progress.BestScores) < levelCount {
		progress.BestScores = append(progress.BestScores, 0)
	}
	if len(progress.BestScores) > levelCount {
		progress.BestScores = progress.BestScores[:levelCount]
	}
	for len(progress.RedDiamonds) < levelCount {
		progress.RedDiamonds = append(progress.RedDiamonds, 0)
	}
	if len(progress.RedDiamonds) > levelCount {
		progress.RedDiamonds = progress.RedDiamonds[:levelCount]
	}
	for len(progress.PurpleDiamonds) < levelCount {
		progress.PurpleDiamonds = append(progress.PurpleDiamonds, 0)
	}
	if len(progress.PurpleDiamonds) > levelCount {
		progress.PurpleDiamonds = progress.PurpleDiamonds[:levelCount]
	}
	for len(progress.SecretExits) < levelCount {
		progress.SecretExits = append(progress.SecretExits, false)
	}
	if len(progress.SecretExits) > levelCount {
		progress.SecretExits = progress.SecretExits[:levelCount]
	}
	for len(progress.AllPurpleClears) < levelCount {
		progress.AllPurpleClears = append(progress.AllPurpleClears, false)
	}
	if len(progress.AllPurpleClears) > levelCount {
		progress.AllPurpleClears = progress.AllPurpleClears[:levelCount]
	}
	for len(progress.AllRedClears) < levelCount {
		progress.AllRedClears = append(progress.AllRedClears, false)
	}
	if len(progress.AllRedClears) > levelCount {
		progress.AllRedClears = progress.AllRedClears[:levelCount]
	}
	for len(progress.NoDamageClears) < levelCount {
		progress.NoDamageClears = append(progress.NoDamageClears, false)
	}
	if len(progress.NoDamageClears) > levelCount {
		progress.NoDamageClears = progress.NoDamageClears[:levelCount]
	}
	for len(progress.NoRecallClears) < levelCount {
		progress.NoRecallClears = append(progress.NoRecallClears, false)
	}
	if len(progress.NoRecallClears) > levelCount {
		progress.NoRecallClears = progress.NoRecallClears[:levelCount]
	}
	for len(progress.NoRestartClears) < levelCount {
		progress.NoRestartClears = append(progress.NoRestartClears, false)
	}
	if len(progress.NoRestartClears) > levelCount {
		progress.NoRestartClears = progress.NoRestartClears[:levelCount]
	}
	for i, steps := range progress.BestSteps {
		if steps < 0 {
			progress.BestSteps[i] = 0
		}
	}
	for i, score := range progress.BestScores {
		if score < 0 {
			progress.BestScores[i] = 0
		}
	}
	for i, redDiamonds := range progress.RedDiamonds {
		if redDiamonds < 0 {
			progress.RedDiamonds[i] = 0
		}
	}
	for i, purpleDiamonds := range progress.PurpleDiamonds {
		if purpleDiamonds < 0 {
			progress.PurpleDiamonds[i] = 0
		}
	}
	if progress.PurpleBank < 0 {
		progress.PurpleBank = 0
	}
	if progress.MaxHealthUpgrades < 0 {
		progress.MaxHealthUpgrades = 0
	}
	if progress.ArmorUpgrades < 0 {
		progress.ArmorUpgrades = 0
	}
	if progress.LifeUpgrades < 0 {
		progress.LifeUpgrades = 0
	}
}

func totalRedDiamonds(progress Progress) int {
	total := 0
	for _, redDiamonds := range progress.RedDiamonds {
		if redDiamonds > 0 {
			total += redDiamonds
		}
	}
	return total
}

func redDiamondRequirement(levelIndex int) int {
	switch levelIndex {
	case 4:
		return 1
	default:
		return 0
	}
}

func ancientSealRequirement(levelCount int) int {
	if levelCount >= 5 {
		return 3
	}
	return 0
}

func sealStatusText(progress Progress, levelCount int) string {
	requirement := ancientSealRequirement(levelCount)
	if requirement <= 0 {
		return "Seal open"
	}
	have := totalRedDiamonds(progress)
	if have > requirement {
		have = requirement
	}
	if progress.AncientSealOpen {
		return fmt.Sprintf("Seal OPEN Red %d/%d", have, requirement)
	}
	return fmt.Sprintf("Seal sealed Red %d/%d", have, requirement)
}

func maxHealthForProgress(progress Progress) int {
	return 3 + progress.MaxHealthUpgrades
}

func maxHealthUpgradeCost(progress Progress) int {
	return 8 + progress.MaxHealthUpgrades*4
}

func maxArmorForProgress(progress Progress) int {
	return progress.ArmorUpgrades
}

func armorUpgradeCost(progress Progress) int {
	return 6 + progress.ArmorUpgrades*3
}

func maxLivesForProgress(progress Progress) int {
	return 3 + progress.LifeUpgrades
}

func lifeUpgradeCost(progress Progress) int {
	return 10 + progress.LifeUpgrades*5
}

func starRating(bestSteps int, parSteps int) int {
	if bestSteps <= 0 {
		return 0
	}
	if parSteps <= 0 || bestSteps <= parSteps {
		return 3
	}
	if bestSteps <= parSteps*13/10 {
		return 2
	}
	return 1
}

func ratingText(stars int) string {
	switch {
	case stars >= 3:
		return "***"
	case stars == 2:
		return "**-"
	case stars == 1:
		return "*--"
	default:
		return "---"
	}
}

func completionScore(baseScore int, steps int, parSteps int) int {
	if parSteps > steps {
		return baseScore + (parSteps-steps)*10
	}
	return baseScore
}

func progressPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		return filepath.Join(".", "diamondrush-progress.json")
	}
	return filepath.Join(configDir, "zskc-diamondrush", "progress.json")
}
