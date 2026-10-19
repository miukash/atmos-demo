package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type StateID int

type CaptureAction int

const (
	CaptureNone CaptureAction = iota
	CaptureStart
	CaptureEnd
)

type Transition struct {
	From          StateID
	To            StateID
	Input         rune
	IsEpsilon     bool
	CaptureAction CaptureAction
}

type Fragment struct {
	Start  StateID
	Accept StateID
}

type NFA struct {
	NextState          StateID
	Transitions        []Transition
	StartState         StateID
	AcceptState        StateID
	TransitionsByState map[StateID][]Transition
}

type Thread struct {
	State        StateID
	CaptureStart int
	CaptureEnd   int
}

type CaptureResult struct {
	Start int
	End   int
}

func NewNFA() *NFA {
	return &NFA{
		TransitionsByState: make(map[StateID][]Transition),
	}
}

func (nfa *NFA) newState() StateID {
	id := nfa.NextState
	nfa.NextState++
	return id
}

func (nfa *NFA) addTransition(from, to StateID, input rune, isEpsilon bool) {
	nfa.addTransitionWithAction(from, to, input, isEpsilon, CaptureNone)
}

func (nfa *NFA) addTransitionWithAction(
	from, to StateID,
	input rune,
	isEpsilon bool,
	action CaptureAction,
) {
	if nfa.TransitionsByState == nil {
		nfa.TransitionsByState = make(map[StateID][]Transition)
	}

	t := Transition{
		From:          from,
		To:            to,
		Input:         input,
		IsEpsilon:     isEpsilon,
		CaptureAction: action,
	}

	nfa.Transitions = append(nfa.Transitions, t)
	nfa.TransitionsByState[from] = append(
		nfa.TransitionsByState[from], t,
	)
}

func (nfa *NFA) addEpsilon(from, to StateID) {
	nfa.addTransition(from, to, 0, true)
}

func (nfa *NFA) compileLiteral(ch rune) Fragment {
	start := nfa.newState()
	accept := nfa.newState()

	nfa.addTransition(start, accept, ch, false)

	return Fragment{Start: start, Accept: accept}
}

func (nfa *NFA) compileConcat(a, b Fragment) Fragment {
	nfa.addEpsilon(a.Accept, b.Start)

	return Fragment{
		Start:  a.Start,
		Accept: b.Accept,
	}
}

func (nfa *NFA) compileAlternate(a, b Fragment) Fragment {
	start := nfa.newState()
	accept := nfa.newState()

	nfa.addEpsilon(start, a.Start)
	nfa.addEpsilon(start, b.Start)
	nfa.addEpsilon(a.Accept, accept)
	nfa.addEpsilon(b.Accept, accept)

	return Fragment{Start: start, Accept: accept}
}

func (nfa *NFA) compileStar(a Fragment) Fragment {
	start := nfa.newState()
	accept := nfa.newState()

	nfa.addEpsilon(start, accept)
	nfa.addEpsilon(start, a.Start)
	nfa.addEpsilon(a.Accept, a.Start)
	nfa.addEpsilon(a.Accept, accept)

	return Fragment{Start: start, Accept: accept}
}

func (nfa *NFA) compileCapture(a Fragment) Fragment {
	start := nfa.newState()
	accept := nfa.newState()

	nfa.addTransitionWithAction(
		start, a.Start, 0, true, CaptureStart,
	)
	nfa.addTransitionWithAction(
		a.Accept, accept, 0, true, CaptureEnd,
	)

	return Fragment{Start: start, Accept: accept}
}

type threadKey struct {
	State        StateID
	CaptureStart int
	CaptureEnd   int
}

func (nfa *NFA) epsilonClosure(
	seeds []Thread,
	position int,
) []Thread {
	queue := append([]Thread(nil), seeds...)
	result := make([]Thread, 0, len(queue))
	visited := make(map[threadKey]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		key := threadKey{
			State:        current.State,
			CaptureStart: current.CaptureStart,
			CaptureEnd:   current.CaptureEnd,
		}

		if visited[key] {
			continue
		}
		visited[key] = true
		result = append(result, current)

		for _, t := range nfa.TransitionsByState[current.State] {
			if !t.IsEpsilon {
				continue
			}

			next := current
			next.State = t.To

			switch t.CaptureAction {
			case CaptureStart:
				// Keep the first start position across repetitions.
				if next.CaptureStart < 0 {
					next.CaptureStart = position
				}
				next.CaptureEnd = -1

			case CaptureEnd:
				next.CaptureEnd = position
			}

			queue = append(queue, next)
		}
	}

	return result
}

func (nfa *NFA) move(threads []Thread, ch rune) []Thread {
	next := make([]Thread, 0)

	for _, current := range threads {
		for _, t := range nfa.TransitionsByState[current.State] {
			if t.IsEpsilon || t.Input != ch {
				continue
			}

			thread := current
			thread.State = t.To
			next = append(next, thread)
		}
	}

	return next
}

// MatchCapture performs a full-string match and returns the capture range.
// Capture offsets are rune indices, with an exclusive End.
func (nfa *NFA) MatchCapture(input string) (CaptureResult, bool) {
	if !utf8.ValidString(input) {
		return CaptureResult{}, false
	}

	runes := []rune(input)
	threads := []Thread{
		{
			State:        nfa.StartState,
			CaptureStart: -1,
			CaptureEnd:   -1,
		},
	}

	for position := 0; ; position++ {
		threads = nfa.epsilonClosure(threads, position)

		if position == len(runes) {
			for _, thread := range threads {
				if thread.State == nfa.AcceptState {
					return CaptureResult{
						Start: thread.CaptureStart,
						End:   thread.CaptureEnd,
					}, true
				}
			}
			return CaptureResult{}, false
		}

		threads = nfa.move(threads, runes[position])
		if len(threads) == 0 {
			return CaptureResult{}, false
		}
	}
}

// Match preserves the original bool-returning API.
func (nfa *NFA) Match(input string) bool {
	_, ok := nfa.MatchCapture(input)
	return ok
}

// MaskCapture masks only the captured range.
// It returns the original input and false if there is no full match
// or no valid capture range.
func (nfa *NFA) MaskCapture(input string, mask rune) (string, bool) {
	result, ok := nfa.MatchCapture(input)
	if !ok || result.Start < 0 ||
		result.End < result.Start {
		return input, false
	}

	runes := []rune(input)
	if result.End > len(runes) {
		return input, false
	}

	for i := result.Start; i < result.End; i++ {
		runes[i] = mask
	}

	return string(runes), true
}

func (nfa *NFA) Mask(input string) string {
	masked, ok := nfa.MaskCapture(input, '*')
	if !ok {
		return input
	}
	return masked
}

type Parser struct {
	pattern    []rune
	pos        int
	nfa        *NFA
	groupCount int
}

func (p *Parser) peek() rune {
	if p.pos >= len(p.pattern) {
		return 0
	}
	return p.pattern[p.pos]
}

func (p *Parser) parseAtom() (Fragment, error) {
	ch := p.peek()

	if ch == 0 {
		return Fragment{}, fmt.Errorf("expected an expression")
	}

	switch ch {
	case '(':
		p.pos++
		p.groupCount++

		if p.groupCount > 1 {
			return Fragment{}, fmt.Errorf(
				"only one capture group is supported",
			)
		}

		inner, err := p.parseExpression()
		if err != nil {
			return Fragment{}, err
		}

		if p.peek() != ')' {
			return Fragment{}, fmt.Errorf(
				"missing closing parenthesis",
			)
		}
		p.pos++

		return p.nfa.compileCapture(inner), nil

	case ')', '|', '*':
		return Fragment{}, fmt.Errorf(
			"unexpected %q at position %d", ch, p.pos,
		)

	default:
		p.pos++
		return p.nfa.compileLiteral(ch), nil
	}
}

func (p *Parser) parseRepeat() (Fragment, error) {
	fragment, err := p.parseAtom()
	if err != nil {
		return Fragment{}, err
	}

	if p.peek() == '*' {
		p.pos++
		fragment = p.nfa.compileStar(fragment)

		if p.peek() == '*' {
			return Fragment{}, fmt.Errorf(
				"repeated '*' is not supported",
			)
		}
	}

	return fragment, nil
}

func (p *Parser) parseConcat() (Fragment, error) {
	if p.peek() == 0 || p.peek() == ')' || p.peek() == '|' {
		return Fragment{}, fmt.Errorf("expected an expression")
	}

	result, err := p.parseRepeat()
	if err != nil {
		return Fragment{}, err
	}

	for p.peek() != 0 && p.peek() != ')' && p.peek() != '|' {
		next, err := p.parseRepeat()
		if err != nil {
			return Fragment{}, err
		}
		result = p.nfa.compileConcat(result, next)
	}

	return result, nil
}

func (p *Parser) parseExpression() (Fragment, error) {
	left, err := p.parseConcat()
	if err != nil {
		return Fragment{}, err
	}

	for p.peek() == '|' {
		p.pos++

		right, err := p.parseConcat()
		if err != nil {
			return Fragment{}, err
		}

		left = p.nfa.compileAlternate(left, right)
	}

	return left, nil
}

func Compile(pattern string) (*NFA, error) {
	if pattern == "" {
		return nil, fmt.Errorf("empty pattern is not supported")
	}
	if !utf8.ValidString(pattern) {
		return nil, fmt.Errorf("pattern is not valid UTF-8")
	}

	nfa := NewNFA()
	parser := &Parser{
		pattern: []rune(pattern),
		nfa:     nfa,
	}

	fragment, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}

	if parser.pos != len(parser.pattern) {
		return nil, fmt.Errorf(
			"unexpected %q at position %d",
			parser.peek(), parser.pos,
		)
	}

	nfa.StartState = fragment.Start
	nfa.AcceptState = fragment.Accept

	return nfa, nil
}

// Optional convenience function for the example in the task.
func Mask(pattern, input string) (string, error) {
	nfa, err := Compile(pattern)
	if err != nil {
		return "", err
	}
	return nfa.Mask(input), nil
}

// Keep strings imported for callers who want a simple demo.
var _ = strings.Builder{}
