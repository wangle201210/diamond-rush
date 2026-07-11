package world

import (
	"fmt"
	"math"
	"strings"

	"github.com/wangle201210/zskc/internal/level"
)

type Tile uint8

const (
	Empty Tile = iota
	Wall
	Dirt
	Diamond
	Rock
	ExitClosed
	ExitOpen
	Key
	Door
	GoldKey
	GoldDoor
	Spike
	TimedSpike
	FireTrap
	Switch
	Bridge
	CrackedWall
	Teleporter
	Lava
	Chest
	Checkpoint
	Potion
	HammerPickup
	HookPickup
	CompassPickup
	SecretExit
	HiddenWall
)

type Direction struct {
	DX int
	DY int
}

var (
	Up    = Direction{DY: -1}
	Down  = Direction{DY: 1}
	Left  = Direction{DX: -1}
	Right = Direction{DX: 1}
)

type Enemy struct {
	X       int
	Y       int
	Dir     Direction
	Type    EnemyType
	Stunned int
}

type Boss struct {
	X     int
	Y     int
	HP    int
	MaxHP int
}

type EnemyType uint8

const (
	EnemyHorizontal EnemyType = iota
	EnemyVertical
	EnemyChaser
)

type Player struct {
	X int
	Y int
}

type Status int

const (
	Playing Status = iota
	Won
	Lost
)

type Event int

const (
	EventNone Event = iota
	EventStep
	EventDig
	EventDiamond
	EventKey
	EventDoor
	EventExitOpen
	EventDeath
	EventWin
	EventTrap
	EventSwitch
	EventBreak
	EventTeleport
	EventBurn
	EventHook
	EventChest
	EventRedDiamond
	EventCheckpoint
	EventRecall
	EventHammer
	EventDamage
	EventPotion
	EventReset
	EventToolHammer
	EventToolHook
	EventCompass
	EventSecretExit
	EventBossHit
	EventBossDefeat
	EventExtraLife
	EventReveal
)

type ChestContent struct {
	Reward string
	Amount int
}

type checkpointSnapshot struct {
	Valid         bool
	Tiles         []Tile
	Falling       []bool
	ChestContents map[int]ChestContent
	Player        Player
	Checkpoint    Player
	Enemies       []Enemy
	Bosses        []Boss
	Diamonds      int
	RedDiamonds   int
	Keys          int
	GoldKeys      int
	Score         int
	gravityPhase  int
	enemyPhase    int
	trapPhase     int
}

type World struct {
	Name             string
	Width            int
	Height           int
	Tiles            []Tile
	Falling          []bool
	ChestContents    map[int]ChestContent
	Player           Player
	Checkpoint       Player
	Enemies          []Enemy
	Bosses           []Boss
	RequiredDiamonds int
	HasHammer        bool
	HasHook          bool
	HasCompass       bool
	SecretExitFound  bool
	Damaged          bool
	RecallUsed       bool
	Diamonds         int
	RedDiamonds      int
	TotalDiamonds    int
	TotalRedDiamonds int
	Keys             int
	GoldKeys         int
	Lives            int
	MaxHealth        int
	Health           int
	MaxArmor         int
	Armor            int
	Score            int
	Steps            int
	Status           Status
	GravityTicks     int
	EnemyTicks       int
	lastEvents       []Event
	gravityPhase     int
	enemyPhase       int
	trapPhase        int
	checkpointState  checkpointSnapshot
}

const (
	defaultGravityTicks = 9
	defaultEnemyTicks   = 18
	defaultLives        = 3
	defaultHealth       = 3
	maxHookRange        = 4
	timedSpikeCycle     = 16
	timedSpikeActive    = 8
	fireTrapStart       = 8
	fireTrapEnd         = 12
	scoreDiamond        = 100
	scoreKey            = 250
	scoreDoor           = 50
	scoreGoldKey        = 400
	scoreGoldDoor       = 100
	scoreEnemy          = 500
	scoreChest          = 750
	scoreRedDiamond     = 1000
	scoreBoss           = 2000
)

func New(def *level.Definition) (*World, error) {
	if len(def.Tiles) != def.Width*def.Height {
		return nil, fmt.Errorf("level tile count mismatch")
	}
	w := &World{
		Name:             def.Name,
		Width:            def.Width,
		Height:           def.Height,
		Tiles:            make([]Tile, len(def.Tiles)),
		Falling:          make([]bool, len(def.Tiles)),
		ChestContents:    map[int]ChestContent{},
		Player:           Player{X: def.PlayerStart.X, Y: def.PlayerStart.Y},
		Checkpoint:       Player{X: def.PlayerStart.X, Y: def.PlayerStart.Y},
		RequiredDiamonds: def.RequiredDiamonds,
		HasHammer:        def.HasHammer,
		HasHook:          def.HasHook,
		Lives:            defaultLives,
		MaxHealth:        defaultHealth,
		Health:           defaultHealth,
		Status:           Playing,
		GravityTicks:     normalizeTicks(def.GravityTicks, defaultGravityTicks),
		EnemyTicks:       normalizeTicks(def.EnemyTicks, defaultEnemyTicks),
	}
	for i, gid := range def.Tiles {
		w.Tiles[i] = tileFromGID(gid)
		if w.Tiles[i] == Diamond {
			w.TotalDiamonds++
		}
	}
	for _, chest := range def.ChestRewards {
		if !w.inBounds(chest.Point.X, chest.Point.Y) {
			return nil, fmt.Errorf("chest at %d,%d out of bounds", chest.Point.X, chest.Point.Y)
		}
		if w.tileAt(chest.Point.X, chest.Point.Y) != Chest {
			return nil, fmt.Errorf("chest reward at %d,%d is not on a chest tile", chest.Point.X, chest.Point.Y)
		}
		if !validChestReward(chest.Reward) {
			return nil, fmt.Errorf("invalid chest reward %q at %d,%d", chest.Reward, chest.Point.X, chest.Point.Y)
		}
		w.ChestContents[w.idx(chest.Point.X, chest.Point.Y)] = ChestContent{
			Reward: chest.Reward,
			Amount: chest.Amount,
		}
		switch chest.Reward {
		case "red_diamond":
			w.TotalRedDiamonds += chest.Amount
		case "purple_diamond", "diamond":
			w.TotalDiamonds += chest.Amount
		}
	}
	for _, spawn := range def.EnemyStarts {
		enemyType := enemyTypeFromString(spawn.Type)
		dir, err := enemyDirection(enemyType, spawn.Direction)
		if err != nil {
			return nil, fmt.Errorf("enemy at %d,%d: %w", spawn.Point.X, spawn.Point.Y, err)
		}
		w.Enemies = append(w.Enemies, Enemy{
			X:    spawn.Point.X,
			Y:    spawn.Point.Y,
			Dir:  dir,
			Type: enemyType,
		})
	}
	for _, spawn := range def.BossStarts {
		if !w.inBounds(spawn.Point.X, spawn.Point.Y) {
			return nil, fmt.Errorf("boss at %d,%d out of bounds", spawn.Point.X, spawn.Point.Y)
		}
		if w.tileAt(spawn.Point.X, spawn.Point.Y) != Empty {
			return nil, fmt.Errorf("boss at %d,%d is not on an empty tile", spawn.Point.X, spawn.Point.Y)
		}
		w.Bosses = append(w.Bosses, Boss{
			X:     spawn.Point.X,
			Y:     spawn.Point.Y,
			HP:    spawn.HP,
			MaxHP: spawn.HP,
		})
	}
	w.saveCheckpointSnapshot()
	return w, nil
}

func validChestReward(reward string) bool {
	switch reward {
	case "", "score", "red_diamond", "purple_diamond", "diamond", "key", "gold_key", "potion", "health", "life", "extra_life", "hammer", "hook", "compass":
		return true
	default:
		return false
	}
}

func copyChestContents(contents map[int]ChestContent) map[int]ChestContent {
	copied := make(map[int]ChestContent, len(contents))
	for idx, content := range contents {
		copied[idx] = content
	}
	return copied
}

func enemyTypeFromString(value string) EnemyType {
	switch value {
	case "enemy_vertical":
		return EnemyVertical
	case "enemy_chaser":
		return EnemyChaser
	default:
		return EnemyHorizontal
	}
}

func defaultEnemyDirection(enemyType EnemyType) Direction {
	if enemyType == EnemyVertical {
		return Down
	}
	return Right
}

func enemyDirection(enemyType EnemyType, value string) (Direction, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return defaultEnemyDirection(enemyType), nil
	case "up":
		return Up, nil
	case "down":
		return Down, nil
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	default:
		return Direction{}, fmt.Errorf("invalid enemy direction %q", value)
	}
}

func tileFromGID(gid int) Tile {
	switch gid {
	case 1:
		return Wall
	case 2:
		return Dirt
	case 3:
		return Diamond
	case 4:
		return Rock
	case 5:
		return ExitClosed
	case 6:
		return Key
	case 7:
		return Door
	case 25:
		return GoldKey
	case 26:
		return GoldDoor
	case 8:
		return ExitOpen
	case 9:
		return Spike
	case 23:
		return TimedSpike
	case 24:
		return FireTrap
	case 10:
		return Switch
	case 11:
		return Bridge
	case 12:
		return CrackedWall
	case 13:
		return Teleporter
	case 14:
		return Lava
	case 15:
		return Chest
	case 16:
		return Checkpoint
	case 17:
		return Potion
	case 18:
		return HammerPickup
	case 19:
		return HookPickup
	case 20:
		return CompassPickup
	case 21:
		return SecretExit
	case 22:
		return HiddenWall
	default:
		return Empty
	}
}

func (w *World) Events() []Event {
	return w.lastEvents
}

func (w *World) ClearEvents() {
	w.lastEvents = w.lastEvents[:0]
}

func (w *World) Update(input Direction) {
	if w.Status != Playing {
		return
	}
	w.ClearEvents()
	if input != (Direction{}) {
		w.tryMovePlayer(input)
	}
	w.advance()
}

func (w *World) UpdateHook(dir Direction) bool {
	if w.Status != Playing {
		return false
	}
	w.ClearEvents()
	if !w.useHook(dir) {
		return false
	}
	w.advance()
	return true
}

func (w *World) advance() {
	w.updateTimedTraps()

	w.gravityPhase++
	if w.gravityPhase >= w.GravityTicks {
		w.gravityPhase = 0
		w.updateGravity()
	}

	w.enemyPhase++
	if w.enemyPhase >= w.EnemyTicks {
		w.enemyPhase = 0
		w.updateEnemies()
	}

	if w.tileAt(w.Player.X, w.Player.Y) == SecretExit {
		w.SecretExitFound = true
		w.Status = Won
		w.add(EventSecretExit)
		w.add(EventWin)
		return
	}
	if w.tileAt(w.Player.X, w.Player.Y) == ExitOpen {
		w.Status = Won
		w.add(EventWin)
	}
}

func (w *World) updateTimedTraps() {
	wasTimedSpikeActive := w.TimedSpikeActive()
	wasFireTrapActive := w.FireTrapActive()
	w.trapPhase = (w.trapPhase + 1) % timedSpikeCycle
	if !wasTimedSpikeActive && w.TimedSpikeActive() && w.tileAt(w.Player.X, w.Player.Y) == TimedSpike {
		w.hurtPlayer(1)
		w.add(EventTrap)
	}
	if !wasFireTrapActive && w.FireTrapActive() && w.tileAt(w.Player.X, w.Player.Y) == FireTrap {
		w.hurtPlayer(1)
		w.add(EventBurn)
	}
}

func (w *World) TimedSpikeActive() bool {
	return w.trapPhase < timedSpikeActive
}

func (w *World) FireTrapActive() bool {
	return w.trapPhase >= fireTrapStart && w.trapPhase < fireTrapEnd
}

func (w *World) tryMovePlayer(dir Direction) {
	nx, ny := w.Player.X+dir.DX, w.Player.Y+dir.DY
	if !w.inBounds(nx, ny) {
		return
	}

	tile := w.tileAt(nx, ny)
	if w.enemyAt(nx, ny) >= 0 || w.bossAt(nx, ny) >= 0 {
		w.hurtPlayer(1)
		return
	}

	switch tile {
	case Empty, ExitOpen, SecretExit:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventStep)
	case HiddenWall:
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventReveal)
	case Spike:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.hurtPlayer(1)
		w.add(EventTrap)
	case TimedSpike:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		if w.TimedSpikeActive() {
			w.hurtPlayer(1)
			w.add(EventTrap)
		} else {
			w.add(EventStep)
		}
	case FireTrap:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		if w.FireTrapActive() {
			w.hurtPlayer(1)
			w.add(EventBurn)
		} else {
			w.add(EventStep)
		}
	case Lava:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.hurtPlayer(1)
		w.add(EventBurn)
	case Teleporter:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.teleportPlayerFrom(nx, ny)
		w.add(EventTeleport)
	case Checkpoint:
		w.Player.X, w.Player.Y = nx, ny
		if w.Checkpoint != w.Player {
			w.Checkpoint = w.Player
			w.saveCheckpointSnapshot()
		}
		w.Steps++
		w.add(EventCheckpoint)
	case Potion:
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.healPlayer(1)
	case HammerPickup:
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.HasHammer = true
		w.Steps++
		w.add(EventToolHammer)
	case HookPickup:
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.HasHook = true
		w.Steps++
		w.add(EventToolHook)
	case CompassPickup:
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.HasCompass = true
		w.Steps++
		w.add(EventCompass)
	case Switch:
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.openBridges()
		w.add(EventSwitch)
	case Dirt:
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventDig)
	case Diamond:
		w.setTile(nx, ny, Empty)
		w.Falling[w.idx(nx, ny)] = false
		w.Diamonds++
		w.Score += scoreDiamond
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventDiamond)
		w.maybeOpenExit()
	case Key:
		w.setTile(nx, ny, Empty)
		w.Keys++
		w.Score += scoreKey
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventKey)
	case GoldKey:
		w.setTile(nx, ny, Empty)
		w.GoldKeys++
		w.Score += scoreGoldKey
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventKey)
	case Chest:
		w.setTile(nx, ny, Empty)
		w.add(EventChest)
		w.openChest(nx, ny)
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
	case Door:
		if w.Keys <= 0 {
			return
		}
		w.Keys--
		w.setTile(nx, ny, Empty)
		w.Score += scoreDoor
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventDoor)
	case GoldDoor:
		if w.GoldKeys <= 0 {
			return
		}
		w.GoldKeys--
		w.setTile(nx, ny, Empty)
		w.Score += scoreGoldDoor
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventDoor)
	case Rock:
		if dir.DY != 0 {
			return
		}
		bx, by := nx+dir.DX, ny
		if !w.inBounds(bx, by) || w.tileAt(bx, by) != Empty || w.enemyAt(bx, by) >= 0 || w.bossAt(bx, by) >= 0 {
			return
		}
		w.setTile(bx, by, Rock)
		w.Falling[w.idx(bx, by)] = false
		w.setTile(nx, ny, Empty)
		w.Player.X, w.Player.Y = nx, ny
		w.Steps++
		w.add(EventStep)
	case ExitClosed, Wall, Bridge, CrackedWall:
		return
	}
}

func (w *World) RecallCheckpoint() bool {
	if w.Status != Playing || w.Lives <= 1 {
		return false
	}
	if w.Player == w.Checkpoint {
		return false
	}
	w.ClearEvents()
	w.Lives--
	w.RecallUsed = true
	w.Player = w.Checkpoint
	w.Steps++
	w.add(EventRecall)
	w.advance()
	return true
}

func (w *World) CompassToCheckpoint() (Direction, int, bool) {
	bestDistance := w.Width + w.Height + 1
	var best Player
	found := false
	for i, tile := range w.Tiles {
		if tile != Checkpoint {
			continue
		}
		x, y := i%w.Width, i/w.Width
		if w.Checkpoint == (Player{X: x, Y: y}) {
			continue
		}
		distance := abs(w.Player.X-x) + abs(w.Player.Y-y)
		if distance < bestDistance {
			bestDistance = distance
			best = Player{X: x, Y: y}
			found = true
		}
	}
	if !found {
		return Direction{}, 0, false
	}
	return Direction{DX: sign(best.X - w.Player.X), DY: sign(best.Y - w.Player.Y)}, bestDistance, true
}

func (w *World) openChest(x, y int) {
	content := w.ChestContents[w.idx(x, y)]
	amount := content.Amount
	if amount <= 0 {
		amount = 1
	}
	switch content.Reward {
	case "red_diamond":
		w.RedDiamonds += amount
		w.Score += scoreRedDiamond * amount
		w.add(EventRedDiamond)
	case "purple_diamond", "diamond":
		w.Diamonds += amount
		w.Score += scoreDiamond * amount
		w.add(EventDiamond)
		w.maybeOpenExit()
	case "key":
		w.Keys += amount
		w.Score += scoreKey * amount
		w.add(EventKey)
	case "gold_key":
		w.GoldKeys += amount
		w.Score += scoreGoldKey * amount
		w.add(EventKey)
	case "potion", "health":
		w.healPlayer(amount)
	case "life", "extra_life":
		w.Lives += amount
		w.add(EventExtraLife)
	case "hammer":
		w.HasHammer = true
		w.add(EventToolHammer)
	case "hook":
		w.HasHook = true
		w.add(EventToolHook)
	case "compass":
		w.HasCompass = true
		w.add(EventCompass)
	default:
		w.Score += scoreChest * amount
	}
}

func (w *World) UseHook(dir Direction) bool {
	return w.useHook(dir)
}

func (w *World) UpdateAction(dir Direction) bool {
	if w.Status != Playing {
		return false
	}
	w.ClearEvents()
	if w.resetCheckpointRoom() {
		w.advance()
		return true
	}
	if !w.useHammer(dir) && !w.useHook(dir) {
		return false
	}
	w.advance()
	return true
}

func (w *World) saveCheckpointSnapshot() {
	w.checkpointState = checkpointSnapshot{
		Valid:         true,
		Tiles:         append([]Tile(nil), w.Tiles...),
		Falling:       append([]bool(nil), w.Falling...),
		ChestContents: copyChestContents(w.ChestContents),
		Player:        w.Checkpoint,
		Checkpoint:    w.Checkpoint,
		Enemies:       append([]Enemy(nil), w.Enemies...),
		Bosses:        append([]Boss(nil), w.Bosses...),
		Diamonds:      w.Diamonds,
		RedDiamonds:   w.RedDiamonds,
		Keys:          w.Keys,
		GoldKeys:      w.GoldKeys,
		Score:         w.Score,
		gravityPhase:  w.gravityPhase,
		enemyPhase:    w.enemyPhase,
		trapPhase:     w.trapPhase,
	}
}

func (w *World) resetCheckpointRoom() bool {
	if w.tileAt(w.Player.X, w.Player.Y) != Checkpoint || !w.checkpointState.Valid {
		return false
	}
	snapshot := w.checkpointState
	w.Tiles = append(w.Tiles[:0], snapshot.Tiles...)
	w.Falling = append(w.Falling[:0], snapshot.Falling...)
	w.ChestContents = copyChestContents(snapshot.ChestContents)
	w.Player = snapshot.Player
	w.Checkpoint = snapshot.Checkpoint
	w.Enemies = append(w.Enemies[:0], snapshot.Enemies...)
	w.Bosses = append(w.Bosses[:0], snapshot.Bosses...)
	w.Diamonds = snapshot.Diamonds
	w.RedDiamonds = snapshot.RedDiamonds
	w.Keys = snapshot.Keys
	w.GoldKeys = snapshot.GoldKeys
	w.Score = snapshot.Score
	w.gravityPhase = snapshot.gravityPhase
	w.enemyPhase = snapshot.enemyPhase
	w.trapPhase = snapshot.trapPhase
	w.Steps++
	w.add(EventReset)
	return true
}

func (w *World) UseHammer(dir Direction) bool {
	return w.useHammer(dir)
}

func (w *World) useHammer(dir Direction) bool {
	if w.Status != Playing || !w.HasHammer || dir == (Direction{}) {
		return false
	}
	tx, ty := w.Player.X+dir.DX, w.Player.Y+dir.DY
	if !w.inBounds(tx, ty) {
		return false
	}
	if w.tileAt(tx, ty) == CrackedWall {
		w.setTile(tx, ty, Empty)
		w.Falling[w.idx(tx, ty)] = false
		w.Steps++
		w.add(EventHammer)
		w.add(EventBreak)
		return true
	}
	if i := w.enemyAt(tx, ty); i >= 0 {
		w.Enemies[i].Stunned = 3
		w.Steps++
		w.add(EventHammer)
		return true
	}
	if i := w.bossAt(tx, ty); i >= 0 {
		w.damageBoss(i, 1)
		w.Steps++
		w.add(EventHammer)
		return true
	}
	return false
}

func (w *World) useHook(dir Direction) bool {
	if w.Status != Playing || !w.HasHook || dir == (Direction{}) {
		return false
	}
	targetX, targetY, ok := w.findHookTarget(dir)
	if !ok {
		return false
	}
	destX, destY := targetX-dir.DX, targetY-dir.DY
	tile := w.tileAt(targetX, targetY)
	w.setTile(destX, destY, tile)
	w.Falling[w.idx(destX, destY)] = false
	w.setTile(targetX, targetY, Empty)
	w.Falling[w.idx(targetX, targetY)] = false
	w.Steps++
	w.add(EventHook)
	return true
}

func (w *World) findHookTarget(dir Direction) (int, int, bool) {
	for distance := 1; distance <= maxHookRange; distance++ {
		x := w.Player.X + dir.DX*distance
		y := w.Player.Y + dir.DY*distance
		if !w.inBounds(x, y) {
			return 0, 0, false
		}
		if w.enemyAt(x, y) >= 0 || w.bossAt(x, y) >= 0 {
			return 0, 0, false
		}
		tile := w.tileAt(x, y)
		if tile == Empty {
			continue
		}
		if distance > 1 && hookable(tile) {
			return x, y, true
		}
		return 0, 0, false
	}
	return 0, 0, false
}

func hookable(tile Tile) bool {
	return tile == Rock || tile == Diamond || tile == Key || tile == GoldKey
}

func (w *World) teleportPlayerFrom(x, y int) {
	current := w.idx(x, y)
	for offset := 1; offset <= len(w.Tiles); offset++ {
		i := (current + offset) % len(w.Tiles)
		if w.Tiles[i] == Teleporter {
			if i != current {
				w.Player.X = i % w.Width
				w.Player.Y = i / w.Width
				return
			}
		}
	}
}

func (w *World) openBridges() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.tileAt(x, y) == Bridge {
				w.setTile(x, y, Empty)
			}
		}
	}
}

func (w *World) maybeOpenExit() {
	if w.Diamonds < w.RequiredDiamonds {
		return
	}
	if w.BossAlive() {
		return
	}
	opened := false
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.tileAt(x, y) == ExitClosed {
				w.setTile(x, y, ExitOpen)
				opened = true
			}
		}
	}
	if opened {
		w.add(EventExitOpen)
	}
}

func (w *World) updateGravity() {
	for y := w.Height - 2; y >= 0; y-- {
		for x := 0; x < w.Width; x++ {
			t := w.tileAt(x, y)
			if t != Rock && t != Diamond {
				continue
			}
			if w.tryFall(x, y) || w.tryRoll(x, y, -1) || w.tryRoll(x, y, 1) {
				continue
			}
			w.Falling[w.idx(x, y)] = false
		}
	}
}

func (w *World) tryFall(x, y int) bool {
	bx, by := x, y+1
	if !w.inBounds(bx, by) {
		return false
	}
	if w.Player.X == bx && w.Player.Y == by {
		if w.Falling[w.idx(x, y)] {
			w.killPlayer()
			return true
		}
		return false
	}
	if i := w.enemyAt(bx, by); i >= 0 && w.Falling[w.idx(x, y)] {
		w.removeEnemy(i)
		w.Score += scoreEnemy
	}
	if i := w.bossAt(bx, by); i >= 0 && w.Falling[w.idx(x, y)] {
		w.damageBoss(i, 2)
		w.setTile(x, y, Empty)
		w.Falling[w.idx(x, y)] = false
		return true
	}
	if w.tileAt(bx, by) == CrackedWall && w.Falling[w.idx(x, y)] {
		w.setTile(bx, by, Empty)
		w.setTile(x, y, Empty)
		w.Falling[w.idx(x, y)] = false
		w.add(EventBreak)
		return true
	}
	if w.tileAt(bx, by) == Lava {
		w.setTile(x, y, Empty)
		w.Falling[w.idx(x, y)] = false
		w.add(EventBurn)
		return true
	}
	if w.tileAt(bx, by) == FireTrap && (w.Falling[w.idx(x, y)] || w.FireTrapActive()) {
		w.setTile(x, y, Empty)
		w.Falling[w.idx(x, y)] = false
		w.add(EventBurn)
		return true
	}
	if w.tileAt(bx, by) != Empty || w.enemyAt(bx, by) >= 0 || w.bossAt(bx, by) >= 0 {
		return false
	}
	t := w.tileAt(x, y)
	w.setTile(bx, by, t)
	w.Falling[w.idx(bx, by)] = true
	w.setTile(x, y, Empty)
	w.Falling[w.idx(x, y)] = false
	return true
}

func (w *World) tryRoll(x, y int, dx int) bool {
	if !w.inBounds(x+dx, y) || !w.inBounds(x+dx, y+1) {
		return false
	}
	below := w.tileAt(x, y+1)
	if below != Rock && below != Diamond && below != Wall && below != Spike && below != TimedSpike && below != FireTrap && below != Bridge && below != CrackedWall && below != Lava {
		return false
	}
	if w.tileAt(x+dx, y) != Empty || w.tileAt(x+dx, y+1) != Empty {
		return false
	}
	if w.enemyAt(x+dx, y) >= 0 || (w.Player.X == x+dx && w.Player.Y == y) {
		return false
	}
	t := w.tileAt(x, y)
	w.setTile(x+dx, y, t)
	w.Falling[w.idx(x+dx, y)] = true
	w.setTile(x, y, Empty)
	w.Falling[w.idx(x, y)] = false
	return true
}

func (w *World) updateEnemies() {
	for _, b := range w.Bosses {
		if abs(b.X-w.Player.X)+abs(b.Y-w.Player.Y) == 1 {
			w.hurtPlayer(1)
			return
		}
	}
	for i := range w.Enemies {
		e := &w.Enemies[i]
		if e.Stunned > 0 {
			e.Stunned--
			continue
		}
		if math.Abs(float64(e.X-w.Player.X))+math.Abs(float64(e.Y-w.Player.Y)) == 1 {
			w.hurtPlayer(1)
			return
		}
		w.moveEnemy(e)
		if e.X == w.Player.X && e.Y == w.Player.Y {
			w.hurtPlayer(1)
			return
		}
	}
}

func (w *World) moveEnemy(e *Enemy) {
	if e.Type == EnemyChaser {
		w.moveChaser(e)
		return
	}
	nx, ny := e.X+e.Dir.DX, e.Y+e.Dir.DY
	if !w.canEnemyMove(nx, ny) {
		e.Dir = Direction{DX: -e.Dir.DX, DY: -e.Dir.DY}
		nx, ny = e.X+e.Dir.DX, e.Y+e.Dir.DY
	}
	if w.canEnemyMove(nx, ny) {
		e.X, e.Y = nx, ny
	}
}

func (w *World) moveChaser(e *Enemy) {
	for _, dir := range w.chaserDirections(e) {
		nx, ny := e.X+dir.DX, e.Y+dir.DY
		if w.canEnemyMove(nx, ny) || (w.Player.X == nx && w.Player.Y == ny) {
			e.Dir = dir
			e.X, e.Y = nx, ny
			return
		}
	}
}

func (w *World) chaserDirections(e *Enemy) []Direction {
	horizontal := Right
	if w.Player.X < e.X {
		horizontal = Left
	}
	vertical := Down
	if w.Player.Y < e.Y {
		vertical = Up
	}
	if abs(w.Player.X-e.X) >= abs(w.Player.Y-e.Y) {
		return []Direction{horizontal, vertical, Direction{DX: -horizontal.DX}, Direction{DY: -vertical.DY}}
	}
	return []Direction{vertical, horizontal, Direction{DY: -vertical.DY}, Direction{DX: -horizontal.DX}}
}

func (w *World) canEnemyMove(x, y int) bool {
	if !w.inBounds(x, y) || w.tileAt(x, y) != Empty || w.enemyAt(x, y) >= 0 || w.bossAt(x, y) >= 0 {
		return false
	}
	return !(w.Player.X == x && w.Player.Y == y)
}

func (w *World) killPlayer() {
	if w.Status != Playing {
		return
	}
	if w.MaxHealth <= 0 {
		w.MaxHealth = defaultHealth
	}
	if w.Health <= 0 {
		w.Health = w.MaxHealth
	}
	w.Damaged = true
	w.Health = 0
	w.add(EventDamage)
	w.consumeLife()
}

func (w *World) healPlayer(amount int) {
	if amount <= 0 {
		return
	}
	if w.MaxHealth <= 0 {
		w.MaxHealth = defaultHealth
	}
	w.Health = min(w.MaxHealth, w.Health+amount)
	w.add(EventPotion)
}

func (w *World) hurtPlayer(damage int) {
	if w.Status != Playing {
		return
	}
	if damage <= 0 {
		return
	}
	if w.MaxHealth <= 0 {
		w.MaxHealth = defaultHealth
	}
	if w.Health <= 0 {
		w.Health = w.MaxHealth
	}
	w.Damaged = true
	if w.Armor > 0 {
		absorbed := min(w.Armor, damage)
		w.Armor -= absorbed
		damage -= absorbed
		w.add(EventDamage)
		if damage <= 0 {
			return
		}
	}
	w.Health -= damage
	w.add(EventDamage)
	if w.Health > 0 {
		return
	}
	w.consumeLife()
}

func (w *World) consumeLife() {
	w.Lives--
	w.add(EventDeath)
	if w.Lives <= 0 {
		w.Health = 0
		w.Armor = 0
		w.Status = Lost
		return
	}
	w.Health = w.MaxHealth
	w.Armor = w.MaxArmor
	w.Player = w.Checkpoint
}

func (w *World) removeEnemy(i int) {
	w.Enemies = append(w.Enemies[:i], w.Enemies[i+1:]...)
}

func (w *World) damageBoss(i int, damage int) {
	if i < 0 || i >= len(w.Bosses) || damage <= 0 {
		return
	}
	w.Bosses[i].HP -= damage
	w.add(EventBossHit)
	if w.Bosses[i].HP > 0 {
		return
	}
	w.Bosses = append(w.Bosses[:i], w.Bosses[i+1:]...)
	w.Score += scoreBoss
	w.add(EventBossDefeat)
	w.maybeOpenExit()
}

func (w *World) BossAlive() bool {
	return len(w.Bosses) > 0
}

func (w *World) BossHealth() (int, int, bool) {
	if len(w.Bosses) == 0 {
		return 0, 0, false
	}
	return w.Bosses[0].HP, w.Bosses[0].MaxHP, true
}

func (w *World) TileAt(x, y int) Tile {
	return w.tileAt(x, y)
}

func (w *World) tileAt(x, y int) Tile {
	if !w.inBounds(x, y) {
		return Wall
	}
	return w.Tiles[w.idx(x, y)]
}

func (w *World) setTile(x, y int, tile Tile) {
	w.Tiles[w.idx(x, y)] = tile
}

func (w *World) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < w.Width && y < w.Height
}

func (w *World) idx(x, y int) int {
	return y*w.Width + x
}

func (w *World) enemyAt(x, y int) int {
	for i, e := range w.Enemies {
		if e.X == x && e.Y == y {
			return i
		}
	}
	return -1
}

func (w *World) bossAt(x, y int) int {
	for i, b := range w.Bosses {
		if b.X == x && b.Y == y {
			return i
		}
	}
	return -1
}

func (w *World) add(event Event) {
	w.lastEvents = append(w.lastEvents, event)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sign(v int) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	default:
		return 0
	}
}

func normalizeTicks(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}
