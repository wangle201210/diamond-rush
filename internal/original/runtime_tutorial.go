package original

const tutorialStageIndex = 13

const (
	TutorialTextBubble = iota + 1
	TutorialTextBottom
)

type TutorialPrompt struct {
	TextIndex int
	Placement int
	X         int
	Y         int
	Side      int
}

type tutorialCommandKind int

const (
	tutorialCommandCamera tutorialCommandKind = iota
	tutorialCommandMove
	tutorialCommandPrompt
	tutorialCommandWait
	tutorialCommandForeground
	tutorialCommandMovePrompt
	tutorialCommandCameraPrompt
	tutorialCommandPortraitFace
	tutorialCommandPortraitPosition
	tutorialCommandPortraitMark
	tutorialCommandPortraitHide
	tutorialCommandFlash
)

type tutorialCommand struct {
	kind         tutorialCommandKind
	x            int
	y            int
	duration     int
	direction    int
	textIndex    int
	placement    int
	promptY      int
	promptSide   int
	foregroundID RawID
	backgroundID RawID
	keep         bool
}

func tutorialCamera(x, y, duration int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandCamera, x: x, y: y, duration: duration}
}

func tutorialMove(direction int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandMove, direction: direction}
}

func tutorialPrompt(textIndex, placement, y, side int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandPrompt, textIndex: textIndex, placement: placement, promptY: y, promptSide: side}
}

func tutorialWait(duration int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandWait, duration: duration}
}

func tutorialForeground(x, y int, foregroundID, backgroundID RawID) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandForeground, x: x, y: y, foregroundID: foregroundID, backgroundID: backgroundID}
}

func tutorialMovePrompt(direction, textIndex int) tutorialCommand {
	return tutorialCommand{
		kind:       tutorialCommandMovePrompt,
		direction:  direction,
		textIndex:  textIndex,
		placement:  TutorialTextBottom,
		promptY:    197,
		promptSide: 0,
	}
}

func tutorialCameraPrompt(x, y, duration, textIndex int) tutorialCommand {
	return tutorialCommand{
		kind:       tutorialCommandCameraPrompt,
		x:          x,
		y:          y,
		duration:   duration,
		textIndex:  textIndex,
		placement:  TutorialTextBottom,
		promptY:    197,
		promptSide: 0,
	}
}

func tutorialPortraitFace(frame int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandPortraitFace, x: frame}
}

func tutorialPortraitPosition(x, y int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandPortraitPosition, x: x, y: y, duration: 5}
}

func tutorialPortraitMark(frame, cycles int, keep bool) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandPortraitMark, x: frame, duration: cycles * 4, keep: keep}
}

func tutorialPortraitHide() tutorialCommand {
	return tutorialCommand{kind: tutorialCommandPortraitHide}
}

func tutorialFlash(cycles int) tutorialCommand {
	return tutorialCommand{kind: tutorialCommandFlash, duration: cycles * 4}
}

var angkorTutorialScripts = map[int][]tutorialCommand{
	29: {
		tutorialCamera(6, 4, 5),
		tutorialPortraitFace(2),
		tutorialPortraitPosition(17, 50),
		tutorialPrompt(12, TutorialTextBubble, 90, 2),
		tutorialPortraitFace(1),
		tutorialPrompt(13, TutorialTextBubble, 90, 1),
		tutorialPortraitFace(3),
		tutorialPortraitMark(3, 3, true),
		tutorialPrompt(14, TutorialTextBubble, 90, 1),
	},
	10: {
		tutorialCamera(31, 6, 5),
		tutorialPortraitFace(7),
		tutorialPortraitPosition(57, 50),
		tutorialPrompt(0, TutorialTextBubble, 90, 2),
		tutorialMove(4),
		tutorialForeground(31, 7, 0, 10),
	},
	11: {
		tutorialPrompt(19, TutorialTextBottom, 197, 0),
		tutorialForeground(31, 7, EmptyRawID, EmptyRawID),
	},
	13: {
		tutorialCamera(38, 7, 20),
		tutorialPrompt(1, TutorialTextBottom, 197, 0),
		tutorialMove(1),
		tutorialMove(1),
		tutorialMove(2),
		tutorialMove(2),
		tutorialMove(2),
		tutorialCamera(40, 8, 20),
		tutorialForeground(42, 8, 31, 2),
		tutorialPrompt(2, TutorialTextBottom, 197, 0),
		tutorialPrompt(3, TutorialTextBottom, 197, 0),
		tutorialForeground(42, 8, EmptyRawID, EmptyRawID),
		tutorialForeground(36, 6, 31, 1),
		tutorialCamera(37, 7, 20),
		tutorialPrompt(4, TutorialTextBottom, 197, 0),
		tutorialPortraitFace(1),
		tutorialForeground(40, 7, EmptyRawID, EmptyRawID),
	},
	15: {
		tutorialForeground(37, 7, EmptyRawID, EmptyRawID),
	},
	16: {
		tutorialMovePrompt(4, 5),
		tutorialMove(2),
		tutorialMove(1),
		tutorialMove(1),
		tutorialMove(2),
		tutorialMove(2),
		tutorialMove(2),
		tutorialCamera(50, 8, 10),
		tutorialPrompt(6, TutorialTextBottom, 197, 0),
		tutorialCameraPrompt(36, 7, 50, 7),
	},
	17: {
		tutorialPrompt(8, TutorialTextBottom, 197, 0),
	},
	22: {
		tutorialPrompt(20, TutorialTextBottom, 197, 0),
		tutorialPrompt(21, TutorialTextBottom, 197, 0),
		tutorialCamera(26, 17, 10),
		tutorialPortraitFace(2),
		tutorialPortraitPosition(17, 50),
		tutorialPortraitMark(1, 2, true),
		tutorialPrompt(22, TutorialTextBubble, 90, 2),
		tutorialPortraitMark(1, 1, false),
		tutorialCamera(22, 4, 40),
		tutorialWait(5),
		tutorialPortraitFace(1),
		tutorialWait(20),
	},
	28: {
		tutorialMove(2),
		tutorialMove(2),
		tutorialMove(2),
		tutorialMove(2),
		tutorialMove(2),
		tutorialCamera(62, 3, 55),
		tutorialPortraitFace(2),
		tutorialPortraitPosition(17, 50),
		tutorialPortraitMark(2, 3, true),
		tutorialPrompt(9, TutorialTextBubble, 90, 1),
		tutorialPortraitHide(),
		tutorialMove(1),
		tutorialMove(4),
		tutorialMove(1),
		tutorialMove(1),
		tutorialMove(1),
		tutorialMove(4),
		tutorialMove(1),
		tutorialCamera(60, 1, 55),
		tutorialPortraitFace(0),
		tutorialPortraitPosition(17, 50),
		tutorialFlash(3),
		tutorialPortraitMark(1, 3, true),
		tutorialPrompt(10, TutorialTextBubble, 230, 2),
		tutorialPortraitMark(1, 1, false),
		tutorialPortraitFace(3),
		tutorialPrompt(11, TutorialTextBubble, 230, 2),
		tutorialPortraitHide(),
		tutorialWait(30),
		tutorialMove(2),
		tutorialWait(70),
	},
	30: {
		tutorialCamera(6, 17, 5),
		tutorialPortraitFace(2),
		tutorialPortraitPosition(17, 50),
		tutorialPortraitMark(1, 2, true),
		tutorialPrompt(17, TutorialTextBubble, 90, 1),
		tutorialPortraitFace(3),
		tutorialPrompt(18, TutorialTextBubble, 90, 2),
	},
	33: {
		tutorialCamera(6, 4, 5),
		tutorialPortraitFace(1),
		tutorialPortraitPosition(17, 50),
		tutorialPrompt(32, TutorialTextBubble, 90, 2),
		tutorialPortraitFace(3),
		tutorialPrompt(33, TutorialTextBubble, 90, 2),
	},
}

// Bavaria's foreground raw-0 events reference the same demo.f command
// stream as the Angkor tutorial. These four scripts are the only demo IDs
// authored in w1.bin.
var bavariaDemoScripts = map[int][]tutorialCommand{
	4: {
		tutorialCamera(13, 16, 30),
		tutorialWait(20),
		tutorialCamera(19, 16, 30),
	},
	6: {
		tutorialCamera(28, 18, 40),
		tutorialWait(20),
		tutorialCamera(26, 11, 40),
		tutorialWait(20),
	},
	19: {
		tutorialCamera(7, 42, 20),
		tutorialWait(20),
		tutorialCamera(13, 56, 45),
		tutorialWait(20),
	},
	34: {
		tutorialPortraitFace(1),
		tutorialPortraitPosition(17, 50),
		tutorialPrompt(34, TutorialTextBubble, 90, 2),
		tutorialPortraitFace(3),
		tutorialFlash(1),
		tutorialPortraitMark(0, 3, true),
		tutorialPrompt(35, TutorialTextBubble, 90, 2),
		tutorialPortraitFace(0),
		tutorialPortraitMark(4, 3, true),
		tutorialPrompt(36, TutorialTextBubble, 90, 2),
	},
}

func (rt *Runtime) initTutorial() {
	rt.TutorialScriptID = -1
	rt.TutorialTextIndex = -1
	rt.TutorialPortraitMark = -1
	rt.tutorialQueuedScript = -1
}

func (rt *Runtime) IsTutorialStage() bool {
	return rt != nil && rt.Stage != nil && rt.Stage.World == WorldAngkor && rt.Stage.Index == tutorialStageIndex
}

func (rt *Runtime) TutorialPrompt() (TutorialPrompt, bool) {
	if rt == nil || !rt.TutorialScriptActive || rt.TutorialTextIndex < 0 {
		return TutorialPrompt{}, false
	}
	return TutorialPrompt{
		TextIndex: rt.TutorialTextIndex,
		Placement: rt.TutorialTextPlacement,
		X:         rt.TutorialPromptX,
		Y:         rt.TutorialTextY,
		Side:      rt.TutorialTextSide,
	}, true
}

func (rt *Runtime) AdvanceTutorialPrompt() bool {
	if _, ok := rt.TutorialPrompt(); !ok {
		return false
	}
	if rt.tutorialPromptAcknowledged {
		return false
	}
	rt.tutorialPromptAcknowledged = true
	return true
}

func (rt *Runtime) SkipTutorialScript() bool {
	if rt == nil || !rt.TutorialScriptActive {
		return false
	}
	rt.TutorialTextIndex = -1
	rt.tutorialPromptAcknowledged = true
	rt.tutorialSkipping = true
	return true
}

func (rt *Runtime) TutorialCamera() (x, y, phase, elapsed, duration int, ok bool) {
	if rt == nil || !rt.TutorialCameraActive {
		return 0, 0, 0, 0, 0, false
	}
	x = clampRuntime(rt.TutorialCameraTarget.X*TileSize-108, 0, max(0, rt.Width()*TileSize-ScreenWidth))
	y = clampRuntime(rt.TutorialCameraTarget.Y*TileSize-108, 0, max(0, rt.Height()*TileSize-(ScreenHeight-80)))
	duration = max(1, rt.TutorialCameraDuration)
	elapsed = clampRuntime(rt.TutorialCameraTicks, 0, duration)
	return x, y, rt.TutorialCameraPhase, elapsed, duration, true
}

func (rt *Runtime) startTutorialForegroundEvent(eventID int) {
	if !rt.IsTutorialStage() {
		if rt.demoScriptAllowed(eventID) {
			rt.queueTutorialScript(eventID)
		}
		return
	}
	switch eventID {
	case 13:
		rt.tutorialResetFirst = true
	case 16:
		rt.tutorialResetSecond = true
	}
	switch eventID {
	case 10, 13, 16, 28, 29:
		rt.queueTutorialScript(eventID)
	}
}

func (rt *Runtime) queueTutorialScript(scriptID int) {
	if !rt.demoScriptAllowed(scriptID) {
		return
	}
	if rt.TutorialScriptActive {
		rt.tutorialQueuedScript = scriptID
		return
	}
	rt.startTutorialScript(scriptID)
}

func (rt *Runtime) startTutorialScript(scriptID int) bool {
	commands, ok := rt.demoScriptCommands(scriptID)
	if !ok || len(commands) == 0 || !rt.demoScriptAllowed(scriptID) || rt.IsTutorialStage() && rt.TutorialComplete {
		return false
	}
	rt.TutorialScriptActive = true
	rt.TutorialScriptID = scriptID
	rt.TutorialTextIndex = -1
	rt.tutorialCommandIndex = 0
	rt.tutorialCommandTicks = 0
	rt.tutorialCommandStarted = false
	rt.tutorialCommandMoveDone = false
	rt.tutorialMoveStarted = false
	rt.tutorialMoveAttempts = 0
	rt.tutorialPromptAcknowledged = false
	rt.tutorialSkipping = false
	rt.TutorialPortraitVisible = false
	rt.TutorialPortraitX = 0
	rt.TutorialPortraitY = 0
	rt.TutorialPortraitFace = 0
	rt.TutorialPortraitMark = -1
	rt.TutorialPortraitRevealTicks = 0
	rt.TutorialFlashVisible = false
	return true
}

func (rt *Runtime) tickTutorial() {
	if rt == nil || rt.Stage == nil {
		return
	}
	if rt.IsTutorialStage() {
		if !rt.tutorialRecallHintTriggered && rt.Player == (Point{X: 46, Y: 7}) && rt.playerSourceOffset() <= 0 {
			rt.TutorialRecallHintVisible = true
			rt.tutorialRecallHintTriggered = true
		}
		if rt.Player == (Point{X: 61, Y: 3}) && rt.playerSourceOffset() == 6 {
			rt.TutorialSealActivated = true
		}
	}
	rt.tickTutorialCamera()
	if !rt.TutorialScriptActive {
		if rt.tutorialQueuedScript >= 0 {
			scriptID := rt.tutorialQueuedScript
			rt.tutorialQueuedScript = -1
			rt.startTutorialScript(scriptID)
		}
		return
	}

	commands, ok := rt.demoScriptCommands(rt.TutorialScriptID)
	if !ok {
		rt.finishTutorialScript()
		return
	}
	for immediate := 0; immediate < 16 && rt.TutorialScriptActive; immediate++ {
		if rt.tutorialCommandIndex >= len(commands) {
			rt.finishTutorialScript()
			return
		}
		command := commands[rt.tutorialCommandIndex]
		if rt.tutorialSkipping {
			switch command.kind {
			case tutorialCommandCamera:
				rt.TutorialCameraActive = false
				rt.advanceTutorialCommand()
				continue
			case tutorialCommandPrompt, tutorialCommandWait:
				rt.advanceTutorialCommand()
				continue
			case tutorialCommandCameraPrompt:
				rt.TutorialCameraActive = false
				rt.advanceTutorialCommand()
				continue
			case tutorialCommandPortraitFace, tutorialCommandPortraitPosition, tutorialCommandPortraitMark, tutorialCommandPortraitHide, tutorialCommandFlash:
				rt.TutorialFlashVisible = false
				rt.advanceTutorialCommand()
				continue
			}
		}
		switch command.kind {
		case tutorialCommandCamera:
			if !rt.tutorialCommandStarted {
				rt.startTutorialCamera(command)
				rt.tutorialCommandStarted = true
				return
			}
			if rt.TutorialCameraActive {
				return
			}
			rt.advanceTutorialCommand()
		case tutorialCommandMove:
			if !rt.tutorialCommandStarted {
				rt.startTutorialMoveCommand()
			}
			if !rt.tickTutorialMoveCommand(command.direction) {
				return
			}
			rt.advanceTutorialCommand()
		case tutorialCommandPrompt:
			if !rt.tutorialCommandStarted {
				rt.startTutorialPrompt(command)
				return
			}
			if command.placement == TutorialTextBubble {
				if !rt.tutorialPromptAcknowledged {
					if rt.TutorialPromptX < 7 {
						rt.TutorialPromptX = min(7, rt.TutorialPromptX+30)
					}
					return
				}
				rt.TutorialPromptX += 30
				if rt.TutorialPromptX <= 240 {
					return
				}
			}
			if !rt.tutorialPromptAcknowledged {
				return
			}
			rt.advanceTutorialCommand()
		case tutorialCommandWait:
			rt.tutorialCommandStarted = true
			rt.tutorialCommandTicks++
			if rt.tutorialCommandTicks <= command.duration {
				return
			}
			rt.advanceTutorialCommand()
		case tutorialCommandForeground:
			rt.setTutorialForeground(command)
			rt.advanceTutorialCommand()
		case tutorialCommandMovePrompt:
			if !rt.tutorialCommandStarted {
				rt.startTutorialMoveCommand()
				if rt.tutorialSkipping {
					rt.tutorialPromptAcknowledged = true
				} else {
					rt.startTutorialPrompt(command)
				}
			}
			if !rt.tutorialCommandMoveDone {
				rt.tutorialCommandMoveDone = rt.tickTutorialMoveCommand(command.direction)
			}
			if !rt.tutorialCommandMoveDone || !rt.tutorialPromptAcknowledged {
				return
			}
			rt.advanceTutorialCommand()
		case tutorialCommandCameraPrompt:
			if !rt.tutorialCommandStarted {
				rt.startTutorialCamera(command)
				rt.startTutorialPrompt(command)
			}
			if rt.TutorialCameraActive || !rt.tutorialPromptAcknowledged {
				return
			}
			rt.advanceTutorialCommand()
		case tutorialCommandPortraitFace:
			rt.TutorialPortraitFace = command.x
			rt.advanceTutorialCommand()
		case tutorialCommandPortraitPosition:
			if !rt.tutorialCommandStarted {
				rt.TutorialPortraitX = command.x
				rt.TutorialPortraitY = command.y
				rt.TutorialPortraitRevealTicks = 0
				rt.tutorialCommandStarted = true
			}
			rt.tutorialCommandTicks++
			if rt.tutorialCommandTicks <= command.duration {
				rt.TutorialPortraitRevealTicks = min(5, rt.tutorialCommandTicks)
				return
			}
			rt.TutorialPortraitRevealTicks = 0
			rt.TutorialPortraitVisible = true
			rt.advanceTutorialCommand()
		case tutorialCommandPortraitMark:
			rt.tutorialCommandStarted = true
			rt.tutorialCommandTicks++
			if rt.tutorialCommandTicks <= command.duration {
				if (rt.tutorialCommandTicks/2)&1 == 0 {
					rt.TutorialPortraitMark = command.x
				} else {
					rt.TutorialPortraitMark = -1
				}
				return
			}
			if command.keep {
				rt.TutorialPortraitMark = command.x
			} else {
				rt.TutorialPortraitMark = -1
			}
			rt.advanceTutorialCommand()
		case tutorialCommandPortraitHide:
			rt.TutorialPortraitVisible = false
			rt.TutorialPortraitMark = -1
			rt.advanceTutorialCommand()
		case tutorialCommandFlash:
			rt.tutorialCommandStarted = true
			rt.tutorialCommandTicks++
			if rt.tutorialCommandTicks <= command.duration {
				rt.TutorialFlashVisible = (rt.tutorialCommandTicks/2)&1 == 0
				return
			}
			rt.TutorialFlashVisible = false
			rt.advanceTutorialCommand()
		}
	}
}

func (rt *Runtime) demoScriptAllowed(scriptID int) bool {
	if rt == nil || rt.Stage == nil {
		return false
	}
	switch rt.Stage.World {
	case WorldAngkor:
		switch scriptID {
		case 22:
			return rt.Stage.Index == 3
		case 30:
			return rt.Stage.Index == 2
		case 33:
			return rt.Stage.Index == 8
		default:
			return rt.IsTutorialStage()
		}
	case WorldBavaria:
		return scriptID == 4 && rt.Stage.Index == 3 ||
			scriptID == 6 && rt.Stage.Index == 8 ||
			scriptID == 34 && rt.Stage.Index == 9 ||
			scriptID == 19 && rt.Stage.Index == 12
	default:
		return false
	}
}

func (rt *Runtime) demoScriptCommands(scriptID int) ([]tutorialCommand, bool) {
	if rt != nil && rt.Stage != nil && rt.Stage.World == WorldBavaria {
		commands, ok := bavariaDemoScripts[scriptID]
		return commands, ok
	}
	commands, ok := angkorTutorialScripts[scriptID]
	return commands, ok
}

func (rt *Runtime) tickTutorialCamera() {
	if !rt.TutorialCameraActive {
		return
	}
	rt.TutorialCameraTicks++
	if rt.TutorialCameraTicks > rt.TutorialCameraDuration {
		rt.TutorialCameraActive = false
	}
}

func (rt *Runtime) startTutorialCamera(command tutorialCommand) {
	rt.TutorialCameraActive = true
	rt.TutorialCameraTarget = Point{X: command.x, Y: command.y}
	rt.TutorialCameraTicks = 0
	rt.TutorialCameraDuration = command.duration
	rt.TutorialCameraPhase++
}

func (rt *Runtime) startTutorialPrompt(command tutorialCommand) {
	rt.TutorialTextIndex = command.textIndex
	rt.TutorialTextPlacement = command.placement
	rt.TutorialTextY = command.promptY
	rt.TutorialTextSide = command.promptSide
	rt.TutorialPromptX = 6
	if command.placement == TutorialTextBubble {
		rt.TutorialPromptX = -240
	}
	rt.tutorialPromptAcknowledged = false
	rt.tutorialCommandStarted = true
}

func (rt *Runtime) startTutorialMoveCommand() {
	rt.tutorialCommandStarted = true
	rt.tutorialCommandMoveDone = false
	rt.tutorialMoveStarted = false
	rt.tutorialMoveAttempts = 0
}

func (rt *Runtime) tickTutorialMoveCommand(direction int) bool {
	if rt.tutorialCommandMoveDone {
		return true
	}
	if rt.tutorialMoveStarted {
		if rt.PlayerMotion.Remaining > 0 {
			return false
		}
		rt.tutorialCommandMoveDone = true
		return true
	}
	if rt.PlayerMotion.Remaining > 0 || !rt.canStartPlayerMove() {
		return false
	}
	dx, dy := tutorialDirection(direction)
	targetID, _ := rt.At(PlayerLayer, rt.Player.X+dx, rt.Player.Y+dy)
	rt.tutorialMoveAttempts++
	if rt.tryMove(dx, dy, true) {
		rt.tutorialMoveStarted = true
		return false
	}
	if (targetID == 0 || targetID == 9) && rt.tutorialMoveAttempts <= boulderPushAttempts+2 {
		return false
	}
	rt.tutorialCommandMoveDone = true
	return true
}

func tutorialDirection(direction int) (dx, dy int) {
	switch direction {
	case 1:
		return 0, -1
	case 2:
		return 1, 0
	case 3:
		return 0, 1
	case 4:
		return -1, 0
	default:
		return 0, 0
	}
}

func (rt *Runtime) setTutorialForeground(command tutorialCommand) {
	if command.x < 0 || command.y < 0 || command.x >= rt.Width() || command.y >= rt.Height() {
		return
	}
	idx := rt.index(command.x, command.y)
	rt.Foreground[idx] = command.foregroundID
	rt.Background[idx] = command.backgroundID
}

func (rt *Runtime) advanceTutorialCommand() {
	rt.tutorialCommandIndex++
	rt.tutorialCommandTicks = 0
	rt.tutorialCommandStarted = false
	rt.tutorialCommandMoveDone = false
	rt.tutorialMoveStarted = false
	rt.tutorialMoveAttempts = 0
	rt.tutorialPromptAcknowledged = false
	rt.TutorialTextIndex = -1
}

func (rt *Runtime) finishTutorialScript() {
	finishedScript := rt.TutorialScriptID
	rt.TutorialScriptActive = false
	rt.TutorialScriptID = -1
	rt.TutorialTextIndex = -1
	rt.TutorialFlashVisible = false
	rt.tutorialCommandIndex = 0
	rt.tutorialCommandStarted = false
	rt.tutorialPromptAcknowledged = false
	rt.tutorialSkipping = false
	if finishedScript == 28 && (rt.Player == (Point{X: 60, Y: 3}) || rt.Player == (Point{X: 61, Y: 3})) {
		rt.TutorialComplete = true
		rt.tutorialQueuedScript = -1
		return
	}
	if rt.tutorialQueuedScript >= 0 {
		scriptID := rt.tutorialQueuedScript
		rt.tutorialQueuedScript = -1
		rt.startTutorialScript(scriptID)
	}
}

func (rt *Runtime) restoreTutorialCheckpoint(first, second bool) {
	if !rt.IsTutorialStage() {
		return
	}
	rt.resetTutorialScriptState()
	rt.TutorialRecallHintVisible = false
	switch {
	case first:
		rt.tutorialResetFirst = false
		rt.tutorialResetSecond = second
		rt.clearTutorialCheckpointForeground(37, 7)
		rt.clearTutorialCheckpointForeground(39, 5)
		rt.startTutorialScript(15)
	case second:
		rt.tutorialResetSecond = false
		rt.clearTutorialCheckpointForeground(46, 7)
		rt.clearTutorialCheckpointForeground(50, 7)
		rt.startTutorialScript(17)
	}
}

func (rt *Runtime) resetTutorialScriptState() {
	rt.TutorialScriptActive = false
	rt.TutorialScriptID = -1
	rt.TutorialTextIndex = -1
	rt.TutorialCameraActive = false
	rt.TutorialCameraTicks = 0
	rt.TutorialCameraDuration = 0
	rt.TutorialPortraitVisible = false
	rt.TutorialPortraitMark = -1
	rt.TutorialPortraitRevealTicks = 0
	rt.TutorialFlashVisible = false
	rt.tutorialCommandIndex = 0
	rt.tutorialCommandTicks = 0
	rt.tutorialCommandStarted = false
	rt.tutorialCommandMoveDone = false
	rt.tutorialMoveStarted = false
	rt.tutorialMoveAttempts = 0
	rt.tutorialPromptAcknowledged = false
	rt.tutorialSkipping = false
	rt.tutorialQueuedScript = -1
}

func (rt *Runtime) clearTutorialCheckpointForeground(x, y int) {
	if x < 0 || y < 0 || x >= rt.Width() || y >= rt.Height() {
		return
	}
	idx := rt.index(x, y)
	rt.Foreground[idx] = EmptyRawID
	rt.Background[idx] = EmptyRawID
	if idx < len(rt.checkpoint.Foreground) {
		rt.checkpoint.Foreground[idx] = EmptyRawID
		rt.checkpoint.Background[idx] = EmptyRawID
	}
}
