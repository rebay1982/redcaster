package main

import (
	"math"
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

func Test_RendererCalculateRayAngle(t *testing.T) {
	testCases := []struct {
		name         string
		pAngle       float64
		fov          float64
		screenColumn int
		expected     float64
	}{
		{
			name:         "player_look_right_leftmost_column",
			pAngle:       0.0,
			fov:          64.0,
			screenColumn: 0,
			expected:     32.0,
		},
		{
			name:         "player_look_right_rightmost_column",
			pAngle:       0.0,
			fov:          64.0,
			screenColumn: FB_WIDTH - 1, // -1 because screen columns are 0 based.
			expected:     328.1,
		},
		{
			name:         "player_look_right_middle_column",
			pAngle:       0.0,
			fov:          64.0,
			screenColumn: FB_WIDTH>>1 - 1, // -1 because screen columns are 0 based.
			expected:     0.1,
		},
		{
			name:         "player_look_up_leftmost_column",
			pAngle:       90.0,
			fov:          64.0,
			screenColumn: 0, // -1 because screen columns are 0 based.
			expected:     122.0,
		},
		{
			name:         "player_look_up_rightmost_column",
			pAngle:       90.0,
			fov:          64.0,
			screenColumn: FB_WIDTH - 1, // -1 because screen columns are 0 based.
			expected:     58.1,
		},
		{
			name:         "player_look_up_middle_column",
			pAngle:       90.0,
			fov:          64.0,
			screenColumn: FB_WIDTH>>1 - 1, // -1 because screen columns are 0 based.
			expected:     90.1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := Game{
				playerAngle: tc.pAngle,
				fov:         tc.fov,
			}

			r := NewRenderer(&game)

			got := r.calculateRayAngle(tc.screenColumn)

			if !aproximately(tc.expected, got) {
				t.Errorf("Expected %f, got %f", tc.expected, got)
			}
		})
	}
}

func Test_RendererCalculateVerticalCollisionRayLength(t *testing.T) {
	game := Game{
		gameMap: [16][16]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
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
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}

	testCases := []struct {
		name     string
		pX, pY   float64
		rAngle   float64
		expected float64
	}{
		{
			name:     "1_1_pos_0_degrees",
			pX:       1.0,
			pY:       1.0,
			rAngle:   0.0,
			expected: 14.0,
		},
		{
			name:     "1_1_pos_180_degrees",
			pX:       1.0,
			pY:       1.0,
			rAngle:   180.0,
			expected: 0.0,
		},
		{
			name:     "1_1_pos_90_degrees_no_vert_hit",
			pX:       1.0,
			pY:       1.0,
			rAngle:   90.0,
			expected: 2048.0,
		},
		{
			name:     "1_1_pos_270_degrees_no_hit",
			pX:       1.0,
			pY:       1.0,
			rAngle:   270.0,
			expected: 2048.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer(&game)

			got := r.calculateVerticalCollisionRayLength(tc.pX, tc.pY, tc.rAngle)

			if !aproximately(tc.expected, got) {
				t.Errorf("Expected %f, got %f", tc.expected, got)
			}
		})
	}
}

func aproximately(x, y float64) bool {
	const tolerance = 0.000001
	epsilon := math.Nextafter(1.0, 2.0) - 1.0
	diff := math.Abs(x - y)

	return diff < math.Max(tolerance*math.Max(math.Abs(x), math.Abs(y)), epsilon*8)
}
