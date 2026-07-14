package original

const waterRecordCount = 15

// Diamond Rush stores three independent 8-pixel water layers in every map
// cell. Each layer uses three owner bits, four shape bits, and two motion bits.
// The two record arrays mirror the original JAR's a:[J and b:[J arrays.
type waterRuntimeState struct {
	Cells   []uint32
	Flows   [waterRecordCount]uint64
	Sources [waterRecordCount]uint64

	Phase       int
	StartNext   bool
	SourceIndex int
	SourceCount int

	DrainTick       int
	BasinTimer      int
	BasinRows       int
	RemainderLayer  int
	PivotX          int
	PivotY          int
	ReflowDirection int
	DrainFlowID     int
	DrainTargetY    int
}

func (state waterRuntimeState) clone() waterRuntimeState {
	state.Cells = append([]uint32(nil), state.Cells...)
	return state
}

type WaterRenderLayer struct {
	Visible bool
	Base    int
	Kind    int
	OffsetX int
	OffsetY int
	Flags   int
}

func packedGet(value uint64, start, width int) int {
	if width <= 0 {
		return 0
	}
	return int(value>>start) & ((1 << width) - 1)
}

func packedSet(value uint64, field, start, width int) uint64 {
	mask := uint64((1 << width) - 1)
	return value&^(mask<<start) | (uint64(field)&mask)<<start
}

func waterCellGet(cell uint32, layer, start, width int) int {
	return packedGet(uint64(cell), layer*9+start, width)
}

func waterCellSet(cell uint32, layer, field, start, width int) uint32 {
	return uint32(packedSet(uint64(cell), field, layer*9+start, width))
}

func (rt *Runtime) initWater() {
	if rt.Stage.World != WorldBavaria {
		return
	}
	rt.water = waterRuntimeState{
		Cells:      make([]uint32, rt.Width()*rt.Height()),
		Phase:      3,
		StartNext:  true,
		DrainTick:  -1,
		BasinTimer: -1,
	}
	rt.WaterSources = rt.WaterSources[:0]
	// afVoid() scans columns before rows; source order is observable because
	// only one reservoir starts flowing at a time.
	for x := 0; x < rt.Width(); x++ {
		for y := 0; y < rt.Height(); y++ {
			if rt.Stage.Player[rt.index(x, y)] != 38 {
				continue
			}
			rt.WaterSources = append(rt.WaterSources, Point{X: x, Y: y})
			amount := int(rt.Stage.Background[rt.index(x, y)])
			if amount == int(EmptyRawID) {
				amount = 0
			}
			if rt.allocWaterSource(amount, x, y, 0) >= 0 {
				rt.water.SourceCount++
			}
		}
	}
	rt.syncWaterState()
}

func (rt *Runtime) allocWaterSource(amount, x, y, direction int) int {
	for id := 0; id < waterRecordCount; id++ {
		if rt.waterSourceGet(id, 0, 4) != 0 {
			continue
		}
		rt.waterSourceSet(id, 1, 0, 4)
		rt.waterSourceSet(id, direction, 4, 2)
		rt.waterSourceSet(id, amount, 6, 7)
		rt.waterSourceSet(id, x, 13, 7)
		rt.waterSourceSet(id, y, 20, 7)
		rt.waterSourceSet(id, amount, 27, 7)
		return id
	}
	return -1
}

func (rt *Runtime) waterSourceGet(id, start, width int) int {
	if id < 0 || id >= waterRecordCount {
		return 0
	}
	return packedGet(rt.water.Sources[id], start, width)
}

func (rt *Runtime) waterSourceSet(id, field, start, width int) {
	if id < 0 || id >= waterRecordCount {
		return
	}
	rt.water.Sources[id] = packedSet(rt.water.Sources[id], field, start, width)
}

func (rt *Runtime) waterSourceAt(x, y int) int {
	for id := 0; id < waterRecordCount; id++ {
		if rt.waterSourceGet(id, 0, 4) != 0 && rt.waterSourceGet(id, 13, 7) == x && rt.waterSourceGet(id, 20, 7) == y {
			return id
		}
	}
	return -1
}

func (rt *Runtime) waterFlowGet(id, start, width int) int {
	if id < 1 || id > waterRecordCount {
		return 0
	}
	return packedGet(rt.water.Flows[id-1], start, width)
}

func (rt *Runtime) waterFlowSet(id, field, start, width int) {
	if id < 1 || id > waterRecordCount {
		return
	}
	rt.water.Flows[id-1] = packedSet(rt.water.Flows[id-1], field, start, width)
}

func (rt *Runtime) freeWaterFlow(id int) {
	if id >= 1 && id <= waterRecordCount {
		rt.water.Flows[id-1] = 0
	}
}

func (rt *Runtime) waterCellAt(x, y int) uint32 {
	if !rt.inBounds(x, y) || len(rt.water.Cells) != rt.Width()*rt.Height() {
		return 0
	}
	return rt.water.Cells[rt.index(x, y)]
}

func (rt *Runtime) setWaterCellField(x, y, layer, field, start, width int) {
	if !rt.inBounds(x, y) || layer < 0 || layer > 2 {
		return
	}
	idx := rt.index(x, y)
	rt.water.Cells[idx] = waterCellSet(rt.water.Cells[idx], layer, field, start, width)
}

func (rt *Runtime) clearWaterLayer(x, y, layer int) {
	if !rt.inBounds(x, y) || layer < 0 || layer > 2 {
		return
	}
	idx := rt.index(x, y)
	shift := layer * 9
	rt.water.Cells[idx] &^= uint32(0x1ff << shift)
}

func (rt *Runtime) allocWaterFlow(x, y, shape, direction, layer, sourceID int) int {
	id := 0
	for candidate := 1; candidate <= waterRecordCount; candidate++ {
		if rt.waterFlowGet(candidate, 28, 3) == 0 {
			id = candidate
			break
		}
	}
	if id == 0 {
		return -1
	}
	rt.water.Flows[id-1] = 0
	if shape == 7 {
		rt.clearWaterLayer(x, y, layer)
		rt.waterFlowSet(id, 7, 28, 3)
	} else {
		rt.waterFlowSet(id, 1, 28, 3)
		rt.waterFlowSet(id, sourceID, 54, 3)
		rt.setWaterCellField(x, y, layer, id, 0, 3)
		rt.setWaterCellField(x, y, layer, shape, 3, 4)
	}
	rt.waterFlowSet(id, x, 31, 6)
	rt.waterFlowSet(id, y, 37, 6)
	rt.waterFlowSet(id, layer, 43, 2)
	if direction < 0 {
		direction = 2
	}
	rt.waterFlowSet(id, direction, 45, 2)
	return id
}

func (rt *Runtime) startWaterSource(x, y int) {
	sourceID := rt.waterSourceAt(x, y)
	if sourceID < 0 {
		return
	}
	rt.waterSourceSet(sourceID, 2, 0, 4)
	if rt.inBounds(x, y) && rt.PlayerLayer[rt.index(x, y)] == 38 {
		rt.PlayerLayer[rt.index(x, y)] = EmptyRawID
	}
	rt.water.Phase = 1
	rt.allocWaterFlow(x, y+1, 0, 0, 0, sourceID)
}

// tickWaterSourceStart mirrors UVoid's alBoolean branch. The source starts
// before the object scan so gravity and wind observe the newly occupied
// sub-cell on the same source frame as the JAR.
func (rt *Runtime) tickWaterSourceStart() {
	if rt.Stage.World != WorldBavaria || len(rt.water.Cells) == 0 || !rt.water.StartNext {
		return
	}
	rt.water.StartNext = false
	if rt.water.SourceIndex < rt.water.SourceCount {
		sourceID := rt.water.SourceIndex
		rt.startWaterSource(rt.waterSourceGet(sourceID, 13, 7), rt.waterSourceGet(sourceID, 20, 7))
		rt.water.SourceIndex++
	} else {
		rt.water.SourceIndex = 0
		rt.water.SourceCount = 0
	}
	rt.syncWaterState()
}

func (rt *Runtime) tickWater() {
	if rt.Stage.World != WorldBavaria || len(rt.water.Cells) == 0 {
		return
	}
	rt.WaterTicks++

	if rt.water.DrainTick >= 0 && rt.water.BasinTimer >= 0 {
		rt.water.DrainTick++
		rt.tickWaterDrain()
	}
	switch rt.water.Phase {
	case 1:
		rt.water.Phase = 2
		rt.syncWaterState()
		return
	case 2:
		for id := 1; id <= waterRecordCount; id++ {
			switch rt.waterFlowGet(id, 28, 3) {
			case 3:
				rt.freeWaterFlow(id)
			case 2:
				if rt.water.BasinTimer >= 0 {
					rt.water.BasinTimer++
					rt.tickWaterBasin(id)
				}
			case 1, 6, 7:
				rt.tickWaterFlow(id)
			}
		}
	case 4:
		rt.cleanupWaterSources()
		rt.water.Phase = 2
	case 5:
		rt.startWaterReflow(rt.water.PivotX, rt.water.PivotY)
	}
	rt.syncWaterState()
}

func (rt *Runtime) tickWaterFlow(id int) {
	x := rt.waterFlowGet(id, 31, 6)
	y := rt.waterFlowGet(id, 37, 6)
	layer := rt.waterFlowGet(id, 43, 2)
	if !rt.inBounds(x, y) {
		rt.freeWaterFlow(id)
		return
	}
	offset := waterCellGet(rt.waterCellAt(x, y), layer, 7, 2)
	direction := rt.waterFlowGet(id, 45, 2)
	if direction > 1 {
		direction = -1
	}
	staticFlow := rt.waterFlowGet(id, 28, 3) == 7
	moveResult := 0

	if offset == 0 {
		nextLayer, nextY := layer+1, y
		if nextLayer > 2 {
			nextLayer = 0
			nextY++
		}
		if waterCellGet(rt.waterCellAt(x, nextY), nextLayer, 3, 4) == 7 {
			rt.water.BasinRows = 1
			if staticFlow {
				rt.freeWaterFlow(id)
				rt.water.Phase = 3
				rt.water.StartNext = true
				return
			}
			basinLayer := nextLayer - 1
			if basinLayer < 0 {
				basinLayer = 2
			}
			rt.water.BasinTimer = 0
			basinID := waterCellGet(rt.waterCellAt(x, nextY), nextLayer, 0, 3)
			rt.waterFlowSet(basinID, rt.waterFlowGet(id, 54, 3), 54, 3)
			rt.waterFlowSet(basinID, 2, 28, 3)
			rt.waterFlowSet(basinID, 0, 47, 2)
			rt.waterFlowSet(basinID, basinLayer, 26, 2)
			rt.waterFlowSet(basinID, x, 14, 6)
			rt.waterFlowSet(basinID, nextY, 20, 6)
			rt.waterFlowSet(basinID, 0, 57, 1)
			rt.freeWaterFlow(id)
			return
		}

		if layer == 2 && rt.waterFlowBlocked(x, y+1) {
			if direction == 0 {
				if staticFlow {
					rt.clearWaterLayer(x, y, layer)
				} else {
					rt.setWaterCellField(x, y, layer, 15, 3, 4)
				}
				moveResult = rt.advanceWaterFlow(id, x, y, direction, layer, 1, staticFlow)
				if moveResult < 0 {
					moveResult = rt.advanceWaterFlow(id, x, y, direction, layer, -1, staticFlow)
				} else if !rt.waterFlowBlocked(x-1, y) {
					shape := 5
					if staticFlow {
						shape = 7
					}
					rt.allocWaterFlow(x-1, y, shape, -1, 2, rt.waterFlowGet(id, 54, 3))
				}
			} else {
				moveResult = rt.advanceWaterFlow(id, x, y, direction, layer, direction, staticFlow)
			}
			if moveResult == -2 {
				below := rt.waterBackgroundBits(x, y+1)
				if below >= 0 && (below>>6)&1 == 1 && rt.waterCanFormBasin(x, y, direction) {
					rt.water.BasinRows = 0
					rt.water.BasinTimer = 0
					rt.waterFlowSet(id, 2, 28, 3)
					rt.waterFlowSet(id, below&0x3f, 0, 7)
					rt.waterFlowSet(id, 2, 26, 2)
					rt.waterFlowSet(id, x, 14, 6)
					rt.waterFlowSet(id, y+1, 20, 6)
				} else {
					rt.waterFlowSet(id, 3, 28, 3)
				}
				return
			}
			direction = moveResult
			x += direction
		} else if direction != 0 {
			direction = 0
			rt.waterFlowSet(id, 0, 45, 2)
		}
	}

	switch direction {
	case 0:
		if !staticFlow && waterCellGet(rt.waterCellAt(x, y), layer, 3, 4) == 0 {
			rt.setWaterCellField(x, y, layer, 3, 3, 4)
		}
		nextLayer := layer + 1
		if nextLayer > 2 {
			nextLayer = 0
			y++
			rt.waterFlowSet(id, y, 37, 6)
		}
		rt.waterFlowSet(id, nextLayer, 43, 2)
		if staticFlow {
			rt.setWaterCellField(x, y, nextLayer, 6, 3, 4)
		} else {
			rt.setWaterCellField(x, y, nextLayer, id, 0, 3)
			rt.setWaterCellField(x, y, nextLayer, 0, 3, 4)
		}
	case -1, 1:
		rt.setWaterCellField(x, y, layer, (offset+1)%3, 7, 2)
	}
}

func (rt *Runtime) advanceWaterFlow(id, x, y, requestedDirection, layer, move int, staticFlow bool) int {
	if rt.waterFlowBlocked(x+move, y) {
		return -2
	}
	if requestedDirection != move {
		direction := move
		if move < 0 {
			direction = 2
		}
		rt.waterFlowSet(id, direction, 45, 2)
	}
	x += move
	rt.waterFlowSet(id, x, 31, 6)
	rt.setWaterCellField(x, y, layer, id, 0, 3)
	if staticFlow {
		rt.clearWaterLayer(x, y, layer)
		return move
	}
	state := rt.waterFlowGet(id, 28, 3)
	shape, edge := 4, 9
	if move < 0 {
		shape, edge = 5, 12
	}
	if state == 6 && rt.water.PivotX == x && rt.water.PivotY == y {
		if move > 0 {
			shape = 11
		} else {
			shape = 14
		}
	}
	if !rt.waterHighBlockedAt(x, y+1) && !rt.waterSpecialBlockedAt(x+move, y) {
		shape = edge
	}
	rt.setWaterCellField(x, y, layer, shape, 3, 4)
	return move
}

func (rt *Runtime) tickWaterBasin(id int) {
	x := rt.waterFlowGet(id, 14, 6)
	y := rt.waterFlowGet(id, 20, 6)
	capacity := rt.waterFlowGet(id, 0, 7)
	filled := rt.waterFlowGet(id, 7, 7)
	phase := rt.waterFlowGet(id, 47, 2)
	sourceID := rt.waterFlowGet(id, 54, 3)
	remaining := rt.waterSourceGet(sourceID, 6, 7)

	if phase == 0 {
		rt.waterFlowSet(id, 1, 47, 2)
		y--
		if y < 0 {
			rt.waterFlowSet(id, 3, 28, 3)
			return
		}
		rt.waterFlowSet(id, y, 20, 6)
		if rt.waterFlowBlocked(x, y) {
			for x > 0 && rt.waterFlowBlocked(x, y) {
				x--
			}
		} else {
			for x+1 < rt.Width() && !rt.waterFlowBlocked(x+1, y) {
				x++
			}
		}
		rt.waterFlowSet(id, x, 14, 6)
		width := 1
		for x-width >= 0 && !rt.waterFlowBlocked(x-width, y) {
			width++
		}
		rt.waterFlowSet(id, width, 49, 5)

		available := remaining - width
		missing := 0
		rt.water.RemainderLayer = 0
		if available < 0 {
			missing = -available
			partialWidth := width + available
			extraLayer := 0
			if partialWidth*3-width != 0 && partialWidth*3 <= width*3/2 {
				extraLayer = 1
			}
			rt.water.RemainderLayer = missing*3/width + extraLayer
			if rt.water.RemainderLayer > 2 {
				rt.waterFlowSet(id, 1, 57, 1)
			}
			available = 0
		}
		newFilled := filled + width - missing
		if newFilled > capacity {
			delta := capacity - newFilled
			newFilled += delta
			available -= delta
		}
		rt.waterFlowSet(id, newFilled, 7, 7)
		rt.waterSourceSet(sourceID, max(0, available), 6, 7)
		return
	}

	width := rt.waterFlowGet(id, 49, 5)
	if width <= 0 || rt.water.BasinTimer == 0 || rt.water.BasinTimer%width != 0 {
		return
	}
	rt.water.BasinTimer = 0
	layer := rt.waterFlowGet(id, 26, 2)
	flag := rt.waterFlowGet(id, 57, 1)
	if flag != 1 {
		for column := 0; column < width; column++ {
			rt.setWaterCellField(x-column, y, layer, 7, 3, 4)
			rt.waterFlowSet(id, layer, 43, 2)
			rt.setWaterCellField(x-column, y, layer, id, 0, 3)
		}
		if rt.water.BasinRows > 0 {
			nextLayer, nextY := layer+1, y
			if nextLayer > 2 {
				nextLayer = 0
				nextY++
			}
			for column := 0; column < width; column++ {
				if !rt.waterFlowBlocked(x-column, nextY) {
					rt.setWaterCellField(x-column, nextY, nextLayer, 8, 3, 4)
				}
			}
		}
		rt.water.BasinRows++
	}

	filled = rt.waterFlowGet(id, 7, 7)
	remaining = rt.waterSourceGet(sourceID, 6, 7)
	if flag == 1 || ((filled >= capacity || remaining == 0) && layer == rt.water.RemainderLayer) {
		if flag == 1 {
			rt.waterFlowSet(id, 0, 57, 1)
		}
		rt.waterFlowSet(id, 3, 28, 3)
		if remaining == 0 {
			rt.waterSourceSet(sourceID, 3, 0, 4)
			rt.waterFlowSet(id, 5, 28, 3)
			if rt.water.DrainTick == -1 {
				rt.water.BasinTimer = -1
				rt.water.Phase = 4
			}
			return
		}
	}
	if layer == 0 {
		rt.waterFlowSet(id, 0, 47, 2)
		rt.waterFlowSet(id, 2, 26, 2)
	} else {
		rt.waterFlowSet(id, layer-1, 26, 2)
	}
}

func (rt *Runtime) cleanupWaterSources() {
	for sourceID := 0; sourceID < waterRecordCount; sourceID++ {
		if rt.waterSourceGet(sourceID, 0, 4) != 3 {
			continue
		}
		x := rt.waterSourceGet(sourceID, 13, 7)
		y := rt.waterSourceGet(sourceID, 20, 7)
		direction := rt.waterSourceGet(sourceID, 4, 2)
		rt.water.Sources[sourceID] = 0
		move, layer := 0, 0
		switch direction {
		case 0:
			y++
		case 1:
			x++
			move, layer = 1, 2
		case 2:
			x--
			move, layer = -1, 2
		}
		rt.allocWaterFlow(x, y, 7, move, layer, rt.waterSourceAt(x, y))
	}
}

func (rt *Runtime) triggerWaterReflow(x, y int) {
	if rt.Stage.World != WorldBavaria || len(rt.water.Cells) == 0 {
		return
	}
	direction := 0
	if rt.waterCellAt(x-1, y) != 0 {
		direction = -1
	} else if rt.waterCellAt(x+1, y) != 0 {
		direction = 1
	} else if rt.waterCellAt(x, y-1) != 0 {
		direction = 2
	}
	if direction == 0 {
		return
	}
	rt.water.Phase = 5
	rt.water.ReflowDirection = direction
	rt.water.PivotX = x
	rt.water.PivotY = y
	rt.emitSound(SoundWater)
}

func (rt *Runtime) startWaterReflow(pivotX, pivotY int) {
	direction := rt.water.ReflowDirection
	rt.water.DrainTick = 0
	cellX, cellY := pivotX+direction, pivotY
	if direction > 1 {
		cellX, cellY = pivotX, pivotY-1
	}
	flowID := waterCellGet(rt.waterCellAt(cellX, cellY), 2, 0, 3)
	if flowID <= 0 || flowID > waterRecordCount {
		rt.water.DrainTick = -1
		rt.water.BasinTimer = -1
		rt.water.Phase = 3
		rt.water.ReflowDirection = 0
		return
	}
	rt.water.DrainFlowID = flowID
	rt.waterFlowSet(flowID, 0, 47, 2)
	total := rt.waterFlowGet(flowID, 7, 7)
	sourceDirection, sourceX, targetY, remaining, sourceY := 0, 0, 0, 0, 0
	if direction == 2 {
		sourceX, targetY, remaining, sourceY = pivotX, pivotY, total, pivotY-1
		rt.water.PivotX, rt.water.PivotY = -1, -1
	} else {
		sourceX = pivotX + direction
		sourceDirection = 2
		if direction < 0 {
			sourceDirection = 1
		}
		targetY = pivotY + 1
		remaining = total - rt.moveWaterForReflow(sourceX, targetY, direction)
		remaining = max(0, remaining)
		sourceY = pivotY
	}
	rt.waterFlowSet(flowID, sourceDirection, 45, 2)
	rt.water.DrainTargetY = targetY
	sourceID := rt.allocWaterSource(remaining, sourceX, sourceY, sourceDirection)
	if sourceID < 0 {
		rt.water.Phase = 3
		rt.water.DrainTick = -1
		rt.water.BasinTimer = -1
		return
	}
	rt.waterSourceSet(sourceID, 2, 0, 4)
	rt.water.Phase = 1
	newFlow := -1
	switch direction {
	case 2:
		newFlow = rt.allocWaterFlow(sourceX, pivotY, 3, 0, 0, sourceID)
	case 1:
		newFlow = rt.allocWaterFlow(sourceX-1, pivotY, 14, -1, 2, sourceID)
	case -1:
		newFlow = rt.allocWaterFlow(sourceX+1, pivotY, 11, 1, 2, sourceID)
	}
	if newFlow > 0 {
		rt.waterFlowSet(newFlow, 6, 28, 3)
	}
	// Java leaves dzInt disabled until the new flow actually reaches a basin.
	rt.water.ReflowDirection = 0
}

func (rt *Runtime) moveWaterForReflow(x, y, direction int) int {
	if direction == 0 {
		return 0
	}
	for rt.inBounds(x, y) &&
		(rt.inBounds(x+direction, y-1) && rt.waterCellAt(x+direction, y-1) != 0 && rt.waterHighBlockedAt(x, y) ||
			rt.waterSpecialBlockedAt(x, y)) {
		x += direction
	}
	moved, rowWidth := 0, 0
	var flowID int
	for rt.inBounds(x, y) {
		first := true
		step := 0
		for rt.inBounds(x+step, y) && !rt.waterFlowBlocked(x+step, y) {
			if first {
				first = false
				flowID = rt.allocWaterFlow(x, y, 8, -2, 2, 0)
				if flowID < 0 {
					return moved
				}
				rt.waterFlowSet(flowID, 5, 28, 3)
				rt.waterFlowSet(flowID, x, 14, 6)
				rt.waterFlowSet(flowID, y, 20, 6)
			}
			rt.setWaterCellField(x+step, y, 2, flowID, 0, 3)
			step += direction
		}
		rowWidth = absInt(step)
		moved += rowWidth
		y++
		searched, exhausted := 0, false
		for rt.inBounds(x, y) && rt.waterCellAt(x, y) == 0 {
			x += direction
			if searched >= rowWidth || !rt.inBounds(x, y) {
				exhausted = true
				break
			}
			searched++
		}
		if exhausted || !rt.inBounds(x, y) {
			break
		}
	}
	if flowID > 0 {
		rt.waterFlowSet(flowID, moved, 7, 7)
	}
	return moved
}

func (rt *Runtime) tickWaterDrain() {
	id := rt.water.DrainFlowID
	if id <= 0 || id > waterRecordCount || rt.waterFlowGet(id, 28, 3) == 0 {
		return
	}
	x := rt.waterFlowGet(id, 14, 6)
	y := rt.waterFlowGet(id, 20, 6)
	phase := rt.waterFlowGet(id, 47, 2)
	if phase == 1 {
		if y == rt.water.DrainTargetY {
			rt.water.DrainTick = -1
			rt.water.BasinTimer = -1
			rt.freeWaterFlow(id)
			rt.water.Phase = 4
			return
		}
		width := rt.waterFlowGet(id, 49, 5)
		if width <= 0 || rt.water.DrainTick == 0 || rt.water.DrainTick%width != 0 {
			return
		}
		rt.water.DrainTick = 0
		layer := rt.waterFlowGet(id, 26, 2)
		for column := 0; column < width; column++ {
			rt.clearWaterLayer(x-column, y, layer)
			nextLayer, nextY := layer+1, y
			if nextLayer > 2 {
				nextLayer, nextY = 0, y+1
			}
			shape := waterCellGet(rt.waterCellAt(x-column, nextY), nextLayer, 3, 4)
			atDrainTarget := nextY == rt.water.DrainTargetY && layer == 2
			writeWake := shape != 12 && shape != 9
			if atDrainTarget {
				writeWake = shape != 0 && shape != 3
			}
			if writeWake {
				rt.setWaterCellField(x-column, nextY, nextLayer, 7, 3, 4)
			}
		}
		if layer == 2 {
			rt.waterFlowSet(id, 0, 47, 2)
			rt.waterFlowSet(id, 0, 26, 2)
			rt.waterFlowSet(id, y+1, 20, 6)
		} else {
			rt.waterFlowSet(id, layer+1, 26, 2)
		}
		return
	}

	rt.waterFlowSet(id, 1, 47, 2)
	width := rt.waterFlowGet(id, 49, 5)
	directionCode := rt.waterFlowGet(id, 45, 2)
	if directionCode != 0 {
		if directionCode == 2 {
			for x >= rt.water.PivotX && !rt.waterFlowBlocked(x+1, y) {
				x++
			}
		} else {
			for x <= rt.water.PivotX && !rt.waterFlowBlocked(x+1, y) {
				x++
			}
		}
	} else if !rt.waterFlowBlocked(x+1, y) {
		for !rt.waterFlowBlocked(x+1, y) {
			x++
		}
	} else {
		for steps := 0; steps < width && rt.waterFlowBlocked(x, y); steps++ {
			x--
		}
	}
	if rt.water.ReflowDirection != 2 && x <= rt.water.PivotX {
		for !rt.waterFlowBlocked(x+1, y) {
			x++
		}
	}
	rt.waterFlowSet(id, x, 14, 6)
	width = 1
	for x-width >= 0 && !rt.waterFlowBlocked(x-width, y) {
		width++
	}
	rt.waterFlowSet(id, width, 49, 5)
}

func (rt *Runtime) waterBackgroundBits(x, y int) int {
	if !rt.inBounds(x, y) {
		return -1
	}
	value := rt.Background[rt.index(x, y)]
	if value == EmptyRawID {
		return -1
	}
	return int(value)
}

func (rt *Runtime) waterCanFormBasin(x, y, direction int) bool {
	for {
		x -= direction
		if direction == 0 {
			if !rt.waterFlowBlocked(x-1, y) && !rt.waterFlowBlocked(x+1, y) {
				continue
			}
			return true
		}
		belowOpen := !rt.waterFlowBlocked(x, y+1)
		if rt.waterFlowBlocked(x, y) {
			return !belowOpen
		}
		if belowOpen {
			return false
		}
	}
}

func (rt *Runtime) waterHighBlockedAt(x, y int) bool {
	if !rt.inBounds(x, y) {
		return true
	}
	id := rt.PlayerLayer[rt.index(x, y)]
	return id >= 80 && id < 208
}

func (rt *Runtime) waterSpecialBlockedAt(x, y int) bool {
	if !rt.inBounds(x, y) {
		return true
	}
	id := rt.PlayerLayer[rt.index(x, y)]
	return id == 10 || id == 34 || id == 35 || id == 37
}

func (rt *Runtime) waterFlowBlocked(x, y int) bool {
	if !rt.inBounds(x, y) {
		return true
	}
	return rt.waterHighBlockedAt(x, y) || rt.waterSpecialBlockedAt(x, y)
}

func (rt *Runtime) waterStable() bool {
	return rt.Stage.World != WorldBavaria || rt.water.Phase == 3
}

func (rt *Runtime) syncWaterState() {
	// Shape 6 is the one-frame wake left by a cleanup flow. The JAR clears it
	// from the renderer once its owner record has advanced to another layer.
	// Doing that here keeps simulation deterministic when tests tick without a
	// Draw call, while preserving the same one-update lifetime.
	for idx, cell := range rt.water.Cells {
		x, y := idx%rt.Width(), idx/rt.Width()
		for layer := 0; layer < 3; layer++ {
			if waterCellGet(cell, layer, 3, 4) != 6 {
				continue
			}
			owner := waterCellGet(cell, layer, 0, 3)
			if owner == 0 || rt.waterFlowGet(owner, 28, 3) == 0 || rt.waterFlowGet(owner, 31, 6) != x || rt.waterFlowGet(owner, 37, 6) != y || rt.waterFlowGet(owner, 43, 2) != layer {
				rt.clearWaterLayer(x, y, layer)
				cell = rt.water.Cells[idx]
			}
		}
	}
	if len(rt.WaterDepth) != len(rt.water.Cells) {
		rt.WaterDepth = make([]uint8, len(rt.water.Cells))
	}
	for idx, cell := range rt.water.Cells {
		depth := uint8(0)
		for layer := 0; layer < 3; layer++ {
			if (cell>>(layer*9))&0x1ff != 0 {
				depth++
			}
		}
		rt.WaterDepth[idx] = depth
	}
	rt.WaterInitializing = rt.water.Phase != 3 || rt.water.StartNext || rt.water.SourceIndex < rt.water.SourceCount
}

func (rt *Runtime) WaterAt(x, y int) uint8 {
	if !rt.inBounds(x, y) || len(rt.WaterDepth) != rt.Width()*rt.Height() {
		return 0
	}
	return rt.WaterDepth[rt.index(x, y)]
}

func (rt *Runtime) UsesSwimmingAnimationAt(x, y int) bool {
	shape := waterCellGet(rt.waterCellAt(x, y), 0, 3, 4)
	return shape == 7 || shape == 8
}

func (rt *Runtime) waterBobOffsetAt(x, y, sourceTick int) int {
	if !rt.inBounds(x, y) {
		return 0
	}
	inWater := rt.waterCellAt(x, y) != 0
	if !inWater && y+1 < rt.Height() && rt.waterCellAt(x, y+1) != 0 {
		inWater = isRoundedGravitySupport(rt.PlayerLayer[rt.index(x, y+1)])
	}
	if !inWater {
		return 0
	}
	phase := (sourceTick >> 1) + x
	offset := phase % 4
	if phase/4&1 != 0 {
		offset = 4 - offset
	}
	return offset
}

func (rt *Runtime) WaterRenderLayerAt(x, y, layer int) WaterRenderLayer {
	if !rt.inBounds(x, y) || layer < 0 || layer > 2 {
		return WaterRenderLayer{}
	}
	cell := rt.waterCellAt(x, y)
	owner := waterCellGet(cell, layer, 0, 3)
	if owner == 0 {
		return WaterRenderLayer{}
	}
	shape := waterCellGet(cell, layer, 3, 4)
	baseShape := shape
	if rt.waterFlowGet(owner, 31, 6) == x && rt.waterFlowGet(owner, 37, 6) == y && rt.waterFlowGet(owner, 43, 2) == layer {
		switch shape {
		case 4:
			baseShape = 1
		case 5:
			baseShape = 2
		}
	}
	offset := waterCellGet(cell, layer, 7, 2) * 8
	if offset > 0 {
		if rt.waterFlowGet(owner, 45, 2) <= 1 {
			offset -= TileSize
		} else {
			offset = TileSize - offset
		}
	}
	return WaterRenderLayer{
		Visible: true,
		Base:    baseShape << 1,
		Kind:    baseShape,
		OffsetX: offset,
		OffsetY: layer * 8,
		Flags:   20,
	}
}

func (rt *Runtime) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < rt.Width() && y < rt.Height()
}
