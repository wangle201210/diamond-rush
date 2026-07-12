package original

func (rt *Runtime) windCellAllowed(x, y int) bool {
	if !rt.inBounds(x, y) {
		return false
	}
	idx := rt.index(x, y)
	playerID := rt.PlayerLayer[idx]
	if playerID != EmptyRawID && playerID >= 80 && playerID < 128 {
		return false
	}
	switch playerID {
	case 30, 10, 37, 34, 35:
		return false
	}
	switch rt.Foreground[idx] {
	case 14, 33, 15, 4, 16:
		return false
	}
	return true
}

// tickWindPodAt mirrors ajVoid(), which runs immediately after raw 47's
// gravity update and intentionally continues to use the object's scan cell.
func (rt *Runtime) tickWindPodAt(x, y int) {
	aboveY := y - 1
	if !rt.inBounds(x, aboveY) || !rt.inBounds(x, y+1) {
		return
	}
	waterSurface := rt.waterCellAt(x, y) != 0 && rt.waterCellAt(x, aboveY) == 0
	playerBeside := rt.isPlayerAt(x-1, y) || rt.isPlayerAt(x+1, y)
	belowID := rt.PlayerLayer[rt.index(x, y+1)]
	belowOccupied := belowID != EmptyRawID && belowID < 128
	aboveIdx := rt.index(x, aboveY)
	if rt.Foreground[aboveIdx] == 35 || !rt.windCellAllowed(x, aboveY) || playerBeside && rt.Pushing || !belowOccupied && !waterSurface {
		return
	}
	rt.Foreground[aboveIdx] = 35
	rt.ForegroundState[aboveIdx] = 18
}

func (rt *Runtime) tickWindForegroundAt(x, y int) {
	if y <= 0 || y+1 >= rt.Height() || rt.Hooking && rt.HookTarget == (Point{X: x, Y: y}) {
		return
	}
	idx := rt.index(x, y)
	if rt.Foreground[idx] != 35 && rt.Foreground[idx] != 37 {
		return
	}
	if rt.ForegroundState[idx] > 0 {
		rt.ForegroundState[idx] = max(0, rt.ForegroundState[idx]-6)
		return
	}

	aboveY := y - 1
	belowY := y + 1
	aboveIdx := rt.index(x, aboveY)
	oldAbove := rt.Foreground[aboveIdx]
	switch oldAbove {
	case 34, 37:
		rt.Foreground[aboveIdx] = 37
	default:
		if oldAbove != 35 && rt.windCellAllowed(x, aboveY) {
			rt.Foreground[aboveIdx] = 35
			rt.ForegroundState[aboveIdx] = 18
		}
	}

	objectID := rt.PlayerLayer[idx]
	if rt.PlayerLayer[aboveIdx] == EmptyRawID && !rt.isPlayerAt(x, aboveY) && oldAbove == 35 && objectID != 32 && objectID != 21 && objectID != EmptyRawID {
		rt.PlayerLayer[aboveIdx] = objectID
		rt.ObjectState[aboveIdx] = rt.ObjectState[idx]&^objectDirectionMask | 1
		rt.ObjectMotion[aboveIdx] = ObjectMotion{DY: -1, Remaining: 18}
		rt.PlayerLayer[idx] = EmptyRawID
		rt.ObjectState[idx] = 0
		rt.ObjectMotion[idx] = ObjectMotion{}
		rt.transferEnemyGateGroup(idx, aboveIdx)
	}

	belowIdx := rt.index(x, belowY)
	if rt.Foreground[belowIdx] != 35 && rt.PlayerLayer[belowIdx] != 47 {
		if rt.Foreground[idx] == 37 {
			rt.Foreground[idx] = 34
		} else {
			rt.Foreground[idx] = EmptyRawID
		}
	}
	if rt.PlayerLayer[idx] == EmptyRawID {
		rt.ForegroundState[idx] = 18
	}
}
