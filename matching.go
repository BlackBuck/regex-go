package main

const (
	startOfText uint8 = 1
	endOfText uint8 = 2
)

func getChar(input string, pos int) uint8 {
	if pos >= len(input) {
		return endOfText
	}

	if pos < 0 {
		return startOfText
	}

	return input[pos]
}

func (s *state) check(input string, pos int) bool {
	ch := getChar(input, pos)

	if ch == endOfText && s.end {
		return true
	}

	if states := s.transitions[ch]; len(states) > 0 {
		nextState := states[0]	
		if nextState.check(input, pos+1) {
			return true
		}
	}

	for _, state := range s.transitions[epsilonChar] {
		if state.check(input, pos) {
			return true
		}

		if ch == startOfText && state.check(input, pos+1) {
			return true
		}
	}

	return false;
}