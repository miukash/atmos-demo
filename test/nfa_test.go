package main

import "testing"

func TestCompileLiteral(t *testing.T) {
	nfa := &NFA{}

	fragment := nfa.compileLiteral('a')

	if fragment.Start != 0 {
		t.Errorf("Start = %d, want 0", fragment.Start)
	}

	if fragment.Accept != 1 {
		t.Errorf("Accept = %d, want 1", fragment.Accept)
	}

	if len(nfa.Transitions) != 1 {
		t.Fatalf(
			"Transitions = %d, want 1",
			len(nfa.Transitions),
		)
	}

	transition := nfa.Transitions[0]

	if transition.From != fragment.Start {
		t.Errorf("From = %d, want %d",
			transition.From, fragment.Start)
	}

	if transition.To != fragment.Accept {
		t.Errorf("To = %d, want %d",
			transition.To, fragment.Accept)
	}

	if transition.Input != 'a' {
		t.Errorf("Input = %q, want 'a'", transition.Input)
	}

	if transition.IsEpsilon {
		t.Error("IsEpsilon = true, want false")
	}

}

func TestCompileConcat(t *testing.T) {
	nfa := &NFA{}

	fragmentA := nfa.compileLiteral('a')
	fragmentB := nfa.compileLiteral('b')

	fragmentC := nfa.compileConcat(fragmentA, fragmentB)

	if fragmentC.Start != fragmentA.Start {
		t.Errorf("Start = %d, want %d",
			fragmentC.Start, fragmentA.Start)
	}

	if fragmentC.Accept != fragmentB.Accept {
		t.Errorf("Accept = %d, want %d",
			fragmentC.Accept, fragmentB.Accept)
	}

	if len(nfa.Transitions) != 3 {
		t.Fatalf(
			"Transitions = %d, want 3",
			len(nfa.Transitions),
		)
	}
}

func TestCompileAlternate(t *testing.T) {
	nfa := &NFA{}

	fragmentA := nfa.compileLiteral('a')
	fragmentB := nfa.compileLiteral('b')

	fragmentC := nfa.compileAlternate(fragmentA, fragmentB)

	if len(nfa.Transitions) != 6 {
		t.Fatalf(
			"Transitions = %d, want 6",
			len(nfa.Transitions),
		)
	}

	if fragmentC.Start == fragmentA.Start ||
		fragmentC.Start == fragmentB.Start {
		t.Errorf("Start = %d, want new state", fragmentC.Start)
	}

	if fragmentC.Accept == fragmentA.Accept ||
		fragmentC.Accept == fragmentB.Accept {
		t.Errorf("Accept = %d, want new state", fragmentC.Accept)
	}
}

func TestCompileStar(t *testing.T) {
	nfa := &NFA{}

	fragmentA := nfa.compileLiteral('a')

	fragmentB := nfa.compileStar(fragmentA)

	if len(nfa.Transitions) != 5 {
		t.Fatalf(
			"Transitions = %d, want 5",
			len(nfa.Transitions),
		)
	}

	if fragmentB.Start == fragmentA.Start {
		t.Errorf("Start = %d, want new state", fragmentB.Start)
	}

	if fragmentB.Accept == fragmentA.Accept {
		t.Errorf("Accept = %d, want new state", fragmentB.Accept)
	}
	if fragmentB.Start == fragmentB.Accept {
		t.Errorf("Start = %d, Accept = %d, want different states",
			fragmentB.Start, fragmentB.Accept)
	}
}

// func TestMove(t *testing.T) {
// 	nfa := &NFA{}

// 	fragmentA := nfa.compileLiteral('a')
// 	fragmentB := nfa.compileLiteral('b')

// 	nfa.compileConcat(fragmentA, fragmentB)

// 	states := []StateID{fragmentA.Start}

// 	nextStates := nfa.move(states, 'a')

// 	if len(nextStates) != 1 {
// 		t.Fatalf(
// 			"NextStates = %d, want 1",
// 			len(nextStates),
// 		)
// 	}

// 	if nextStates[0] != fragmentA.Accept {
// 		t.Errorf("NextState = %d, want %d",
// 			nextStates[0], fragmentA.Accept)
// 	}
// }

func TestMatch(t *testing.T) {
	tests := []struct {
		name  string
		build func(*NFA) Fragment
		input string
		want  bool
	}{
		{
			name: "literal match",
			build: func(nfa *NFA) Fragment {
				return nfa.compileLiteral('a')
			},
			input: "a",
			want:  true,
		},
		{
			name: "literal rejects extra character",
			build: func(nfa *NFA) Fragment {
				return nfa.compileLiteral('a')
			},
			input: "aa",
			want:  false,
		},
		{
			name: "concat match",
			build: func(nfa *NFA) Fragment {
				a := nfa.compileLiteral('a')
				b := nfa.compileLiteral('b')
				return nfa.compileConcat(a, b)
			},
			input: "ab",
			want:  true,
		},
		{
			name: "alternate match",
			build: func(nfa *NFA) Fragment {
				a := nfa.compileLiteral('a')
				b := nfa.compileLiteral('b')
				return nfa.compileAlternate(a, b)
			},
			input: "b",
			want:  true,
		},
		{
			name: "star matches empty",
			build: func(nfa *NFA) Fragment {
				a := nfa.compileLiteral('a')
				return nfa.compileStar(a)
			},
			input: "",
			want:  true,
		},
		{
			name: "star matches repeated characters",
			build: func(nfa *NFA) Fragment {
				a := nfa.compileLiteral('a')
				return nfa.compileStar(a)
			},
			input: "aaa",
			want:  true,
		},
		{
			name: "star rejects different character",
			build: func(nfa *NFA) Fragment {
				a := nfa.compileLiteral('a')
				return nfa.compileStar(a)
			},
			input: "b",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa := &NFA{}
			fragment := tt.build(nfa)

			nfa.StartState = fragment.Start
			nfa.AcceptState = fragment.Accept

			if got := nfa.Match(tt.input); got != tt.want {
				t.Fatalf(
					"Match(%q) = %v, want %v",
					tt.input,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCompileAndMatch(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"a", "a", true},
		{"a", "b", false},
		{"ab", "ab", true},
		{"ab", "abc", false},
		{"a|b", "a", true},
		{"a|b", "b", true},
		{"a|b", "c", false},
		{"a*", "", true},
		{"a*", "aaa", true},
		{"a*", "b", false},
		{"ab|c*", "ab", true},
		{"ab|c*", "ccc", true},
		{"ab|c*", "c", true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"/"+tt.input, func(t *testing.T) {
			nfa, err := Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Compile returned error: %v", err)
			}

			if got := nfa.Match(tt.input); got != tt.want {
				t.Errorf(
					"Match(%q) = %v, want %v",
					tt.input, got, tt.want,
				)
			}
		})
	}
}

func TestCompileInvalidPatterns(t *testing.T) {
	tests := []string{
		"a)",
		"a|",
		"*a",
		"a**",
		"(",
	}

	for _, pattern := range tests {
		t.Run(pattern, func(t *testing.T) {
			if _, err := Compile(pattern); err == nil {
				t.Errorf("Compile(%q) expected an error", pattern)
			}
		})
	}
}

func TestCompiler(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"a", "a", true},
		{"a", "b", false},
		{"ab", "ab", true},
		{"ab", "abc", false},
		{"a|b", "a", true},
		{"a|b", "b", true},
		{"a|b", "c", false},
		{"a*", "", true},
		{"a*", "aaa", true},
		{"a*", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"/"+tt.input, func(t *testing.T) {
			nfa, err := Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Compile returned error: %v", err)
			}

			if got := nfa.Match(tt.input); got != tt.want {
				t.Errorf(
					"Match(%q) = %v, want %v",
					tt.input, got, tt.want,
				)
			}
		})
	}
}

func TestCompilerGroups(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		// 括弧で囲まれた連結
		{
			name:    "group matches",
			pattern: "(ab)",
			input:   "ab",
			want:    true,
		},
		{
			name:    "group does not partially match",
			pattern: "(ab)",
			input:   "a",
			want:    false,
		},

		// 括弧の中で選択
		{
			name:    "group alternate left",
			pattern: "(ab|cd)",
			input:   "ab",
			want:    true,
		},
		{
			name:    "group alternate right",
			pattern: "(ab|cd)",
			input:   "cd",
			want:    true,
		},
		{
			name:    "group alternate no match",
			pattern: "(ab|cd)",
			input:   "ef",
			want:    false,
		},

		// 括弧と繰り返し
		{
			name:    "group repeated",
			pattern: "(ab)*",
			input:   "abab",
			want:    true,
		},
		{
			name:    "group repeated zero times",
			pattern: "(ab)*",
			input:   "",
			want:    true,
		},

		// 括弧の外との連結
		{
			name:    "group in concatenation",
			pattern: "x(ab)y",
			input:   "xaby",
			want:    true,
		},
		{
			name:    "group in concatenation no match",
			pattern: "x(ab)y",
			input:   "xabz",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Compile(%q) returned error: %v", tt.pattern, err)
			}

			if got := nfa.Match(tt.input); got != tt.want {
				t.Errorf(
					"Match(%q) = %v, want %v",
					tt.input, got, tt.want,
				)
			}
		})
	}
}

func TestCompilerInvalidGroups(t *testing.T) {
	patterns := []string{
		"(ab",  // 閉じ括弧がない
		"ab)",  // 対応する開き括弧がない
		"()",   // 空のグループは未対応
		"(a|)", // 選択の右側が空
	}

	for _, pattern := range patterns {
		t.Run(pattern, func(t *testing.T) {
			_, err := Compile(pattern)
			if err == nil {
				t.Errorf("Compile(%q) expected an error", pattern)
			}
		})
	}
}

func TestCaptureMask(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    string
	}{
		{"x(ab)y", "xaby", "x**y"},
		{"(abc)", "abc", "***"},
		{"x(ab)y", "zaby", "zaby"}, // Full matchではない
		{"x(ab)y", "xabyz", "xabyz"},
		{"x(a|b)y", "xay", "x*y"},
		{"x(a|b)y", "xby", "x*y"},
		{"x(ab)*y", "xababy", "x****y"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"/"+tt.input, func(t *testing.T) {
			nfa, err := Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Compile: %v", err)
			}

			got := nfa.Mask(tt.input)
			if got != tt.want {
				t.Errorf("Mask() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMatchPreservesBoolAPI(t *testing.T) {
	nfa, err := Compile("x(ab)y")
	if err != nil {
		t.Fatal(err)
	}

	if !nfa.Match("xaby") {
		t.Fatal("expected a match")
	}
	if nfa.Match("xabc") {
		t.Fatal("expected no full match")
	}
}

func TestCaptureRange(t *testing.T) {
	nfa, err := Compile("x(ab)y")
	if err != nil {
		t.Fatal(err)
	}

	got, ok := nfa.MatchCapture("xaby")
	if !ok {
		t.Fatal("expected a match")
	}

	if got.Start != 1 || got.End != 3 {
		t.Fatalf("got range [%d,%d), want [1,3)", got.Start, got.End)
	}
}

func TestUnicodeCaptureRange(t *testing.T) {
	nfa, err := Compile("x(猫)y")
	if err != nil {
		t.Fatal(err)
	}

	got, ok := nfa.MatchCapture("x猫y")
	if !ok {
		t.Fatal("expected a match")
	}

	if got.Start != 1 || got.End != 2 {
		t.Fatalf("got range [%d,%d), want [1,2)", got.Start, got.End)
	}

	if masked := nfa.Mask("x猫y"); masked != "x*y" {
		t.Fatalf("Mask() = %q, want %q", masked, "x*y")
	}
}

func TestInvalidPatterns(t *testing.T) {
	patterns := []string{
		"",
		"(ab",
		"ab)",
		"a|",
		"*a",
		"a**",
		"(a)(b)",
	}

	for _, pattern := range patterns {
		t.Run(pattern, func(t *testing.T) {
			if _, err := Compile(pattern); err == nil {
				t.Fatalf("Compile(%q) expected an error", pattern)
			}
		})
	}
}
