package sendtablescs2

type fieldState struct {
	state []any
}

func newFieldState() *fieldState {
	return &fieldState{
		state: make([]any, 8),
	}
}

// maxFieldIndex bounds decoded field-path components and variable-collection sizes.
// Anything larger means a corrupt or desynced bitstream and would allocate GBs in one call.
const maxFieldIndex = 1 << 20

func validateFieldIndex(z int, fp *fieldPath) {
	if z < 0 || z > maxFieldIndex {
		_panicf("field path component %d out of range [0, %d] (path %v): corrupt or desynced demo bitstream", z, maxFieldIndex, fp.path[:fp.last+1])
	}
}

func (s *fieldState) get(fp *fieldPath) any {
	x := s
	z := 0

	for i := 0; i <= fp.last; i++ {
		z = fp.path[i]
		if len(x.state) < z+1 {
			return nil
		}
		if i == fp.last {
			return x.state[z]
		}
		if _, ok := x.state[z].(*fieldState); !ok {
			return nil
		}
		x = x.state[z].(*fieldState)
	}

	return nil
}

func (s *fieldState) set(fp *fieldPath, v any) {
	x := s
	z := 0

	for i := 0; i <= fp.last; i++ {
		z = fp.path[i]
		validateFieldIndex(z, fp)

		if y := len(x.state); y <= z {
			newCap := max(z+2, y*2)
			if z+2 > cap(x.state) {
				newSlice := make([]any, z+1, newCap)
				copy(newSlice, x.state)
				x.state = newSlice
			} else {
				// Re-slice to update the length without allocating new memory
				x.state = x.state[:z+1]
			}
		}

		if i == fp.last {
			if _, ok := x.state[z].(*fieldState); !ok {
				x.state[z] = v
			}
			return
		}

		if _, ok := x.state[z].(*fieldState); !ok {
			x.state[z] = newFieldState()
		}

		x = x.state[z].(*fieldState)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
