package main

import (
	"testing"
)

func Test_RendererCheckWallCollision(t *testing.T) {
	game := Game{
		gameMap: [16][16]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}

	testCases := []struct {
		name     string
		x, y     float64
		expected bool
	}{
		{
			name:     "x_y_no_collision",
			x:        1.0,
			y:        1.0,
			expected: false,
		},
		{
			name:     "x_y_float_no_collision",
			x:        1.20,
			y:        1.45,
			expected: false,
		},
		{
			name:     "x_y_float_close_no_collision",
			x:        1.99,
			y:        1.99,
			expected: false,
		},
		{
			name:     "x_y_collision",
			x:        0.0,
			y:        0.0,
			expected: true,
		},
		{
			name:     "x_y_float_collision",
			x:        2.20,
			y:        2.45,
			expected: true,
		},
		{
			name:     "x_negative_out_of_bound",
			x:        -1.0,
			y:        0.0,
			expected: true,
		},
		{
			name:     "x_positive_out_of_bound",
			x:        16.0,
			y:        0.0,
			expected: true,
		},
		{
			name:     "y_negative_out_of_bound",
			x:        0.0,
			y:        -1.0,
			expected: true,
		},
		{
			name:     "y_positive_out_of_bound",
			x:        0.0,
			y:        16.0,
			expected: true,
		},
		{
			name:     "x_y_negative_out_of_bound",
			x:        -1.0,
			y:        -1.0,
			expected: true,
		},
		{
			name:     "x_y_positive_out_of_bound",
			x:        16.0,
			y:        16.0,
			expected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			r := NewRenderer(&game)

			got := r.checkWallCollision(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, got)
			}
		})
	}
}
