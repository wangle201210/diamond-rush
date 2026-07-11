package original

import (
	"slices"
	"testing"
)

func newStage13Route(t *testing.T) *stage07Route {
	t.Helper()
	stage := mustLoadOriginalStage(t, "stage13.json")
	rt, err := NewRuntime(stage)
	if err != nil {
		t.Fatal(err)
	}
	return &stage07Route{t: t, rt: rt}
}

func completeStage13TutorialScript(t *testing.T, route *stage07Route, scriptID int, wantPrompts []int) {
	t.Helper()
	rt := route.rt
	if !rt.TutorialScriptActive || rt.TutorialScriptID != scriptID {
		t.Fatalf("tutorial script active=%v id=%d, want active script %d", rt.TutorialScriptActive, rt.TutorialScriptID, scriptID)
	}
	prompts := make([]int, 0, len(wantPrompts))
	for tick := 0; tick < 1200 && rt.TutorialScriptActive; tick++ {
		if prompt, ok := rt.TutorialPrompt(); ok {
			if rt.AdvanceTutorialPrompt() {
				prompts = append(prompts, prompt.TextIndex)
			}
		}
		route.tick()
	}
	if rt.TutorialScriptActive {
		t.Fatalf("tutorial script %d did not complete at player=%+v command=%d", scriptID, rt.Player, rt.tutorialCommandIndex)
	}
	if !slices.Equal(prompts, wantPrompts) {
		t.Fatalf("tutorial script %d prompts=%v, want %v", scriptID, prompts, wantPrompts)
	}
}

func TestRuntimeStage13RecallHintAndSealUseSourceTriggerFrames(t *testing.T) {
	route := newStage13Route(t)
	rt := route.rt
	rt.Player = Point{X: 46, Y: 7}
	result := rt.TickSourceFrame(8, 1, 0)
	if !rt.TutorialRecallHintVisible || result.TutorialSealActivated {
		t.Fatalf("recall hint=%v seal event=%v, want visible/false at (46,7)", rt.TutorialRecallHintVisible, result.TutorialSealActivated)
	}

	rt.Player = Point{X: 61, Y: 3}
	rt.PlayerMotion = ObjectMotion{DX: 1, Remaining: 12}
	result = rt.TickSourceFrame(8, 2, 0)
	if rt.TutorialSealActivated || result.TutorialSealActivated {
		t.Fatal("tutorial seal activated before source movement offset 6")
	}
	rt.PlayerMotion.Remaining = 6
	result = rt.TickSourceFrame(8, 3, 0)
	if !rt.TutorialSealActivated || !result.TutorialSealActivated {
		t.Fatalf("seal active=%v event=%v, want one activation at movement offset 6", rt.TutorialSealActivated, result.TutorialSealActivated)
	}
	result = rt.TickSourceFrame(8, 4, 0)
	if result.TutorialSealActivated {
		t.Fatal("tutorial seal activation event repeated")
	}
}

func TestRuntimeStage13TutorialCanBeCompletedAtSourceCadence(t *testing.T) {
	route := newStage13Route(t)
	rt := route.rt
	if rt.CompassEnabled {
		t.Fatal("tutorial compass starts enabled")
	}
	route.tick()
	route.walkTo("automatic tutorial entrance before door", Point{X: 4, Y: 4})
	if !rt.CloseEntranceDoor() {
		t.Fatal("failed to close tutorial entrance door")
	}
	route.walkTo("tutorial entrance checkpoint", Point{X: 5, Y: 4})
	route.walkTo("opening tutorial event", Point{X: 6, Y: 4})
	route.waitUntil("opening tutorial portrait", 80, func() bool {
		prompt, ok := rt.TutorialPrompt()
		return ok && prompt.TextIndex == 12
	})
	if !rt.TutorialPortraitVisible || rt.TutorialPortraitFace != 2 || rt.TutorialPortraitX != 17 || rt.TutorialPortraitY != 50 {
		t.Fatalf("opening portrait visible=%v face=%d position=%d,%d", rt.TutorialPortraitVisible, rt.TutorialPortraitFace, rt.TutorialPortraitX, rt.TutorialPortraitY)
	}
	completeStage13TutorialScript(t, route, 29, []int{12, 13, 14})

	// The upper passage reaches event 10 without stepping on the compass chest.
	route.walkTo("left of compass-chest bypass", Point{X: 27, Y: 5})
	route.walkTo("right side of compass-chest bypass", Point{X: 30, Y: 5})
	route.walkTo("compass-chest prompt event", Point{X: 31, Y: 7})
	completeStage13TutorialScript(t, route, 10, []int{0})
	if rt.Player != (Point{X: 30, Y: 7}) {
		t.Fatalf("script 10 player=%+v, want automatic left move to 30,7", rt.Player)
	}

	route.walkTo("compass chest", Point{X: 28, Y: 6})
	route.waitUntil("compass chest reward", 180, func() bool {
		return rt.SpecialPickup42 && rt.TutorialScriptActive && rt.TutorialScriptID == 11
	})
	completeStage13TutorialScript(t, route, 11, []int{19})
	foreground, _ := rt.At(ForegroundLayer, 31, 7)
	if !rt.CompassEnabled || foreground != EmptyRawID {
		t.Fatalf("compass enabled=%v event31 foreground=%d", rt.CompassEnabled, foreground)
	}

	route.walkTo("tutorial reset checkpoint", Point{X: 36, Y: 7})
	route.walkTo("push-reset tutorial event", Point{X: 37, Y: 7})
	completeStage13TutorialScript(t, route, 13, []int{1, 2, 3, 4})
	route.walkTo("return to reset checkpoint", Point{X: 36, Y: 7})
	if !rt.ResetCheckpoint() {
		t.Fatal("failed to reset tutorial checkpoint with center action")
	}
	completeStage13TutorialScript(t, route, 15, nil)
	if foreground, _ := rt.At(ForegroundLayer, 37, 7); foreground != EmptyRawID {
		t.Fatalf("first reset event foreground=%d, want cleared", foreground)
	}

	route.walkTo("recall tutorial event", Point{X: 46, Y: 7})
	completeStage13TutorialScript(t, route, 16, []int{5, 6, 7})
	if !rt.RecallCheckpoint() {
		t.Fatal("failed to start tutorial star-key recall")
	}
	route.waitUntil("tutorial recall completes", 120, func() bool {
		return !rt.RecallPending && rt.TutorialScriptActive && rt.TutorialScriptID == 17
	})
	completeStage13TutorialScript(t, route, 17, []int{8})
	if rt.ExtraLives != 4 {
		t.Fatalf("tutorial recall lives=%d, want 4", rt.ExtraLives)
	}
	if foreground, _ := rt.At(ForegroundLayer, 46, 7); foreground != EmptyRawID {
		t.Fatalf("second reset event foreground=%d, want cleared", foreground)
	}

	route.walkTo("left of final seal tutorial event", Point{X: 56, Y: 8})
	if !rt.TryMove(1, 0) {
		t.Fatal("failed to enter final seal tutorial event")
	}
	route.waitUntil("final seal tutorial starts", 40, func() bool {
		return rt.TutorialScriptActive && rt.TutorialScriptID == 28
	})
	completeStage13TutorialScript(t, route, 28, []int{9, 10, 11})
	if !rt.TutorialComplete || !rt.TutorialSealActivated || rt.Player != (Point{X: 61, Y: 3}) {
		t.Fatalf("tutorial finish complete=%v seal=%v player=%+v", rt.TutorialComplete, rt.TutorialSealActivated, rt.Player)
	}
}
