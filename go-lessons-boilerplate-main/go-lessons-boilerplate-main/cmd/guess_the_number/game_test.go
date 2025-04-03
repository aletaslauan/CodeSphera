package main

import (
	"testing"
)

func TestGame_Guess(t *testing.T) {
	tests := []struct {
		name     string
		game     *Game
		guesses  []int
		expected []string
	}{
		{
			name:     "Win on first try",
			game:     NewGame(42, 3),
			guesses:  []int{42},
			expected: []string{"You have won!"},
		},
		{
			name:     "Win on last try",
			game:     NewGame(77, 3),
			guesses:  []int{50, 90, 77},
			expected: []string{"Too low!", "Too high!", "You have won!"},
		},
		{
			name:     "Lose all tries",
			game:     NewGame(10, 2),
			guesses:  []int{1, 2},
			expected: []string{"Too low!", "You have lost!"},
		},
		{
			name:     "Guess after losing",
			game:     NewGame(55, 2),
			guesses:  []int{60, 60, 55},
			expected: []string{"Too high!", "You have lost!", "You have lost!"},
		},
		{
			name:     "Guess after winning",
			game:     NewGame(33, 3),
			guesses:  []int{10, 33, 40},
			expected: []string{"Too low!", "You have won!", "You have won!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []string
			for _, guess := range tt.guesses {
				results = append(results, tt.game.Guess(guess))
			}
			for i := range results {
				if results[i] != tt.expected[i] {
					t.Errorf("[%s] Guess #%d = %q, want %q", tt.name, i+1, results[i], tt.expected[i])
				}
			}
		})
	}
}
