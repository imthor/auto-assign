package selector

import (
	"testing"
)

func TestRoundRobin(t *testing.T) {
	users := []string{"alice", "bob", "charlie"}
	rr := &RoundRobin{}

	tests := []struct {
		name      string
		lastIndex int
		want      int
	}{
		{
			name:      "first selection",
			lastIndex: -1,
			want:      0,
		},
		{
			name:      "next in sequence",
			lastIndex: 0,
			want:      1,
		},
		{
			name:      "wrap around",
			lastIndex: 2,
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rr.SelectNext(users, tt.lastIndex, nil)
			if err != nil {
				t.Errorf("RoundRobin.SelectNext() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("RoundRobin.SelectNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandom(t *testing.T) {
	users := []string{"alice", "bob", "charlie"}
	r := &Random{}

	// Test multiple selections to ensure we get different indices
	seen := make(map[int]bool)
	for i := 0; i < 100; i++ {
		got, err := r.SelectNext(users, -1, nil)
		if err != nil {
			t.Errorf("Random.SelectNext() error = %v", err)
			return
		}
		if got < 0 || got >= len(users) {
			t.Errorf("Random.SelectNext() = %v, want index in range [0, %d]", got, len(users)-1)
			return
		}
		seen[got] = true
	}

	// Check if we got a good distribution of indices
	if len(seen) < 2 {
		t.Errorf("Random.SelectNext() did not produce enough variety in selections")
	}
}

func TestLeastAssigned(t *testing.T) {
	users := []string{"alice", "bob", "charlie"}
	la := &LeastAssigned{}

	tests := []struct {
		name   string
		counts map[string]int
		want   int
	}{
		{
			name:   "all equal",
			counts: map[string]int{"alice": 0, "bob": 0, "charlie": 0},
			want:   0,
		},
		{
			name:   "one least assigned",
			counts: map[string]int{"alice": 2, "bob": 1, "charlie": 2},
			want:   1,
		},
		{
			name:   "multiple least assigned",
			counts: map[string]int{"alice": 1, "bob": 1, "charlie": 2},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := la.SelectNext(users, -1, tt.counts)
			if err != nil {
				t.Errorf("LeastAssigned.SelectNext() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("LeastAssigned.SelectNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectorEdgeCases(t *testing.T) {
	selectors := []struct {
		name     string
		selector Selector
	}{
		{"RoundRobin", &RoundRobin{}},
		{"Random", &Random{}},
		{"LeastAssigned", &LeastAssigned{}},
	}

	for _, s := range selectors {
		t.Run(s.name, func(t *testing.T) {
			// Test empty users list
			_, err := s.selector.SelectNext([]string{}, -1, nil)
			if err == nil {
				t.Errorf("%s.SelectNext() with empty users list should return error", s.name)
			}

			// Test nil counts map
			users := []string{"alice", "bob", "charlie"}
			_, err = s.selector.SelectNext(users, -1, nil)
			if err != nil {
				t.Errorf("%s.SelectNext() with nil counts should not return error", s.name)
			}
		})
	}
}
