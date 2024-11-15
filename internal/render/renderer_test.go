package render

import (
	"math"

	"testing"

	"github.com/rebay1982/redcaster/internal/data"
	"github.com/rebay1982/redcaster/internal/game"
)

const (
	FB_WIDTH  = 640
	FB_HEIGHT = 480
)

func Test_RendererCalculateRayAngle(t *testing.T) {
	testCases := []struct {
		name         string
		pAngle       float64
		fov          float64
		screenColumn int
		expected     float64
	}{
		{
			name:         "player_look_0_column_0",
			pAngle:       0.0,
			fov:          64.0,
			screenColumn: 0,
			expected:     32.0,
		},
		{
			name:         "player_look_0_column_max_width",
			pAngle:       0.0,
			fov:          64.0,
			screenColumn: FB_WIDTH - 1, // -1 because screen columns are 0 based.
			expected:     328.080535,   // Uses precomputed values for FOV of 64.
		},
		{
			name:         "player_look_0_column_half_width",
			pAngle:       0.0,
			fov:          64.0,
			screenColumn: FB_WIDTH >> 1,
			expected:     0.0,
		},
		{
			name:         "player_look_90_column_0",
			pAngle:       90.0,
			fov:          64.0,
			screenColumn: 0, // -1 because screen columns are 0 based.
			expected:     122.0,
		},
		{
			name:         "player_look_90_max_width",
			pAngle:       90.0,
			fov:          64.0,
			screenColumn: FB_WIDTH - 1, // -1 because screen columns are 0 based.
			expected:     58.080535,
		},
		{
			name:         "player_look_90_half_width",
			pAngle:       90.0,
			fov:          64.0,
			screenColumn: FB_WIDTH >> 1,
			expected:     90.0,
		},
		{
			name:         "player_look_350_column_0",
			pAngle:       350.0,
			fov:          64.0,
			screenColumn: 0,
			expected:     22.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			levelData := data.LevelData{
				PlayerCoordData: data.PlayerCoordData{
					PlayerAngle: tc.pAngle,
				},
			}
			game := game.NewGame(levelData, nil)

			config := NewRenderConfiguration(FB_WIDTH, FB_HEIGHT, tc.fov)
			r := NewRenderer(config, &game)

			got := r.computeRayAngle(tc.screenColumn)

			if !approximately(tc.expected, got) {
				t.Errorf("Expected %f, got %f", tc.expected, got)
			}
		})
	}
}

func Test_RendererCalculateVerticalCollisionRayLength(t *testing.T) {
	levelData := data.LevelData{
		Map: [][]int{
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
	game := game.NewGame(levelData, nil)

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
			name:     "14_1_pos_0_degrees",
			pX:       14.9999999, // Completely against the right wall on row 1.
			pY:       1.0,
			rAngle:   0.0,
			expected: 0.0000001,
		},
		{
			name:     "13_7_pos_15_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   15.0,
			expected: 2.070552361,
		},
		{
			name:     "13_7_pos_30_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   30.0,
			expected: 2.309401077,
		},
		{
			name:     "13_7_pos_45_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   45.0,
			expected: 2.828427125,
		},
		{
			name:     "13_7_pos_60_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   60.0,
			expected: 4.0,
		},
		{
			name:     "13_7_pos_75_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   75.0,
			expected: 7.72740661,
		},
		{
			name:     "1_1_pos_90_degrees_no_vert_hit",
			pX:       1.0,
			pY:       1.0,
			rAngle:   90.0,
			expected: 2048.0,
		},
		{
			name:     "3_7_pos_105_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   105.0,
			expected: 7.72740661,
		},
		{
			name:     "3_7_pos_120_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   120.0,
			expected: 4.0,
		},
		{
			name:     "3_7_pos_135_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   135.0,
			expected: 2.828427125,
		},
		{
			name:     "3_7_pos_150_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   150.0,
			expected: 2.309401077,
		},
		{
			name:     "3_7_pos_165_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   165.0,
			expected: 2.070552361,
		},
		{
			name:     "1_1_pos_180_degrees",
			pX:       1.0,
			pY:       1.0,
			rAngle:   180.0,
			expected: 0.0,
		},
		{
			name:     "3_7_pos_195_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   195.0,
			expected: 2.070552361,
		},
		{
			name:     "3_7_pos_210_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   210.0,
			expected: 2.309401077,
		},
		{
			name:     "3_7_pos_225_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   225.0,
			expected: 2.828427125,
		},
		{
			name:     "3_7_pos_240_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   240.0,
			expected: 4.0,
		},
		{
			name:     "3_7_pos_255_degrees",
			pX:       3.0,
			pY:       7.0,
			rAngle:   255.0,
			expected: 7.72740661,
		},
		{
			name:     "1_1_pos_270_degrees_no_hit",
			pX:       1.0,
			pY:       1.0,
			rAngle:   270.0,
			expected: 2048.0,
		},
		{
			name:     "13_7_pos_285_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   285.0,
			expected: 7.72740661,
		},
		{
			name:     "13_7_pos_300_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   300.0,
			expected: 4.0,
		},
		{
			name:     "13_7_pos_315_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   315.0,
			expected: 2.828427125,
		},
		{
			name:     "13_7_pos_330_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   330.0,
			expected: 2.309401077,
		},
		{
			name:     "13_7_pos_345_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   345.0,
			expected: 2.070552361,
		},
		{
			name:     "13_7_pos_360_degrees",
			pX:       13.0,
			pY:       7.0,
			rAngle:   360.0,
			expected: 2.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := NewRenderConfiguration(FB_WIDTH, FB_HEIGHT, 64.0)
			r := NewRenderer(config, &game)

			got := r.computeVerticalCollisionRayLength(tc.pX, tc.pY, tc.rAngle)

			if !approximately(tc.expected, got) {
				t.Errorf("Expected %f, got %f", tc.expected, got)
			}
		})
	}
}

func Test_RendererCalculateHorizontalCollisionRayLength(t *testing.T) {
	levelData := data.LevelData{
		Map: [][]int{
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
	game := game.NewGame(levelData, nil)

	testCases := []struct {
		name     string
		pX, pY   float64
		rAngle   float64
		expected float64
	}{
		{
			name:     "1_1_pos_0_degrees_no_horizontal_wall_hit",
			pX:       1.0,
			pY:       1.0,
			rAngle:   0.0,
			expected: 2048.0,
		},
		{
			name:     "7_3_pos_15_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   15.0,
			expected: 7.72740661,
		},
		{
			name:     "7_3_pos_30_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   30.0,
			expected: 4.0,
		},
		{
			name:     "7_3_pos_45_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   45.0,
			expected: 2.828427125,
		},
		{
			name:     "7_3_pos_60_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   60.0,
			expected: 2.309401077,
		},
		{
			name:     "7_3_pos_75_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   75.0,
			expected: 2.070552361,
		},
		{
			name:     "7_3_pos_90_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   90.0,
			expected: 2.0,
		},
		{
			name:     "7_3_pos_105_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   105.0,
			expected: 2.070552361,
		},
		{
			name:     "7_3_pos_120_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   120.0,
			expected: 2.309401077,
		},
		{
			name:     "7_3_pos_135_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   135.0,
			expected: 2.828427125,
		},
		{
			name:     "7_3_pos_150_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   150.0,
			expected: 4.0,
		},
		{
			name:     "7_3_pos_165_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   165.0,
			expected: 7.72740661,
		},
		{
			name:     "1_1_pos_180_degrees_no_horizontal_wall_hit",
			pX:       1.0,
			pY:       1.0,
			rAngle:   180.0,
			expected: 2048.0,
		},
		{
			name:     "7_3_pos_195_degrees",
			pX:       7.0,
			pY:       3.0,
			rAngle:   195.0,
			expected: 7.72740661,
		},
		{
			name:     "7_13_pos_210_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   210.0,
			expected: 4.0,
		},
		{
			name:     "7_13_pos_225_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   225.0,
			expected: 2.828427125,
		},
		{
			name:     "7_13_pos_240_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   240.0,
			expected: 2.309401077,
		},
		{
			name:     "7_13_pos_255_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   255.0,
			expected: 2.070552361,
		},
		{
			name:     "7_13_pos_270_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   270.0,
			expected: 2.0,
		},
		{
			name:     "7_13_pos_285_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   285.0,
			expected: 2.070552361,
		},
		{
			name:     "7_13_pos_300_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   300.0,
			expected: 2.309401077,
		},
		{
			name:     "7_13_pos_315_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   315.0,
			expected: 2.828427125,
		},
		{
			name:     "7_13_pos_330_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   330.0,
			expected: 4.0,
		},
		{
			name:     "7_13_pos_345_degrees",
			pX:       7.0,
			pY:       13.0,
			rAngle:   345.0,
			expected: 7.72740661,
		},
		{
			name:     "7_13_pos_360_degrees_no_horizontal_wall_hit",
			pX:       7.0,
			pY:       13.0,
			rAngle:   360.0,
			expected: 2048.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := NewRenderConfiguration(FB_WIDTH, FB_HEIGHT, 64.0)
			r := NewRenderer(config, &game)

			//fmt.Println(tc.name)
			got := r.computeHorizontalCollisionRayLength(tc.pX, tc.pY, tc.rAngle)

			if !approximately(tc.expected, got) {
				t.Errorf("Expected %f, got %f", tc.expected, got)
			}
		})
	}
}

func approximately(x, y float64) bool {
	const tolerance = 0.000001
	epsilon := math.Nextafter(1.0, 2.0) - 1.0
	diff := math.Abs(x - y)

	return diff < math.Max(tolerance*math.Max(math.Abs(x), math.Abs(y)), epsilon*8)
}
