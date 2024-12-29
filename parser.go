package main

type state struct {
	start bool
	end bool
	transitions map[uint8][]*state
}

const epsilonChar uint8 = 0

func toNFA(ctx *parseContext) *state {
	startState, endState := tokenToNFA(&ctx.tokens[0])

	for i := 1;i< len(ctx.tokens);i++ {
		startNext, endNext := tokenToNFA(&ctx.tokens[i])
		endState.transitions[epsilonChar] = append(
			endState.transitions[epsilonChar],
			startNext,
		)

		endState = endNext
	}

	start := &state{
		transitions: map[uint8][]*state{
			epsilonChar: {startState},
		},
		start: true,
	}

	end := &state{
		transitions: map[uint8][]*state{},
		end: true,
	}

	endState.transitions[epsilonChar] = append(
		end.transitions[epsilonChar],
		end,
	)

	return start
}

func tokenToNFA(t *token) (*state, *state) {

	start := &state{
		transitions: map[uint8][]*state{},
	}

	end := &state{
		transitions: map[uint8][]*state{},
	}

	switch t.tokenType {
	case literal:
		ch := t.value.(uint8)
		start.transitions[ch] = []*state{end}
	case or:
		values := t.value.([]token)
		left := values[0]
		right := values[1]

		s1, e1 := tokenToNFA(&left)
		s2, e2 := tokenToNFA(&right)

		start.transitions[epsilonChar] = []*state{s1, s2}	
		e1.transitions[epsilonChar] = []*state{end}
		e2.transitions[epsilonChar] = []*state{end}
	case bracket:
		literals := t.value.(map[uint8]bool)
		for l := range literals {
			start.transitions[l] = []*state{end}
		}
	case group, groupUncaptured:
		tokens := t.value.([]token)
		start, end = tokenToNFA(&tokens[0])

		for i := 1;i<len(tokens);i++ {
			ts, te := tokenToNFA(&tokens[i])
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				ts,
			)
			end = te
		}
	case repeat:
		p := t.value.(repeatPayload)

		if p.min == 0 { 
			start.transitions[epsilonChar] = []*state{end}
		}

		var copyCount int 

		if p.max == repeatInfinity {
			if p.min == 0 {
				copyCount = 1
			} else {
				copyCount = p.min
			}
		} else {
			copyCount = p.max
		}

		from, to := tokenToNFA(&p.token) 
		start.transitions[epsilonChar] = append( 
			start.transitions[epsilonChar], 
			from,
		) 

		for i := 2; i <= copyCount; i++ { 
			s, e := tokenToNFA(&p.token)

			// connect the end of the previous one 
			// to the start of this one
			to.transitions[epsilonChar] = append(
				to.transitions[epsilonChar], 
				s,
			) 

			// keep track of the previous NFA's entry and exit states
			from = s 
			to = e   

			// after the minimum required amount of repetitions
			// the rest must be optional, thus we add an 
			// epsilon transition to the start of each NFA 
			// so that we can skip them if needed
			if i > p.min { 
				s.transitions[epsilonChar] = append(
					s.transitions[epsilonChar], 
					end,
				)
			}
		}

		to.transitions[epsilonChar] = append( 
			to.transitions[epsilonChar], 
			end,
		) 

		if p.max == repeatInfinity { 
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar], 
				from,
			)
		}
	default:
		panic("Unknown token type")	
	}

	return start, end
}