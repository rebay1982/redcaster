package render

import (
	"math"

	"testing"

	"github.com/rebay1982/redcaster/internal/config"
	"github.com/rebay1982/redcaster/internal/data"
	"github.com/rebay1982/redcaster/internal/game"
)

const (
	FB_WIDTH  = 640
	FB_HEIGHT = 480
)

func Test_RendererCalculateRayAngle(t *testing.T) {
	var tManager TextureManager = nil

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
			config := config.NewRenderConfiguration(FB_WIDTH, FB_HEIGHT, tc.fov)
			r := NewRenderer(config, &game, tManager, levelData)

			got := r.computeRayAngle(tc.screenColumn)

			if !approximately(tc.expected, got) {
				t.Errorf("Expected %f, got %f", tc.expected, got)
			}
		})
	}
}

func Test_RendererCalculateVerticalCollision(t *testing.T) {
	var tManager TextureManager = nil

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
		expected collisionDetail
	}{
		{
			name:   "1_1_pos_0_degrees",
			pX:     1.0,
			pY:     1.0,
			rAngle: 0.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 1.0,
				},
				rayAngle:  0.0,
				rayLength: 14.0,
				wallType:  1,
			},
		},
		{
			name:   "14_1_pos_0_degrees",
			pX:     14.9999999, // Completely against the right wall on row 1.
			pY:     1.0,
			rAngle: 0.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 14.9999999,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 1.0,
				},
				rayAngle:  0.0,
				rayLength: 0.0000001,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_15_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 15.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 6.464102,
				},
				rayAngle:  15.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_30_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 30.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 5.845299,
				},
				rayAngle:  30.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_45_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 45.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 5.0,
				},
				rayAngle:  45.0,
				rayLength: 2.828427,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_60_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 60.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 3.535898,
				},
				rayAngle:  60.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_75_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 75.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: -0.464102,
				},
				rayAngle:  75.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "1_1_pos_90_degrees_no_vert_hit",
			pX:     1.0,
			pY:     1.0,
			rAngle: 90.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 2048.0,
					y: 2048.0,
				},
				rayAngle:  90.0,
				rayLength: 2048,
				wallType:  0,
			},
		},
		{
			name:   "3_7_pos_105_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 105.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: -0.464102,
				},
				rayAngle:  105.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_120_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 120.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 3.535898,
				},
				rayAngle:  120.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_135_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 135.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 5.0,
				},
				rayAngle:  135.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_150_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 150.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 5.845299,
				},
				rayAngle:  150.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_165_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 165.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 6.464102,
				},
				rayAngle:  165.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "1_1_pos_180_degrees",
			pX:     1.0,
			pY:     1.0,
			rAngle: 180.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayAngle:  180.0,
				rayLength: 0.0,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_195_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 195.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 7.535898,
				},
				rayAngle:  195.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_210_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 210.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 8.154701,
				},
				rayAngle:  210.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_225_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 225.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 9.0,
				},
				rayAngle:  225.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_240_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 240.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 10.464102,
				},
				rayAngle:  240.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "3_7_pos_255_degrees",
			pX:     3.0,
			pY:     7.0,
			rAngle: 255.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 3.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 1.0,
					y: 14.464102,
				},
				rayAngle:  255.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "1_1_pos_270_degrees_no_hit",
			pX:     1.0,
			pY:     1.0,
			rAngle: 270.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 2048.0,
					y: 2048.0,
				},
				rayAngle:  270.0,
				rayLength: 2048.0,
				wallType:  0,
			},
		},
		{
			name:   "13_7_pos_285_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 285.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 14.464102,
				},
				rayAngle:  285.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_300_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 300.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 10.464102,
				},
				rayAngle:  300.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_315_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 315.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 9.0,
				},
				rayAngle:  315.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_330_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 330.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 8.154701,
				},
				rayAngle:  330.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_345_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 345.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 7.535898,
				},
				rayAngle:  345.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "13_7_pos_360_degrees",
			pX:     13.0,
			pY:     7.0,
			rAngle: 360.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 13.0,
					y: 7.0,
				},
				rayEnd: coordinates{
					x: 15.0,
					y: 7.0,
				},
				rayAngle:  360.0,
				rayLength: 2.0,
				wallType:  1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := config.NewRenderConfiguration(FB_WIDTH, FB_HEIGHT, 64.0)
			r := NewRenderer(config, &game, tManager, data.LevelData{})

			got := r.computeVerticalCollision(tc.pX, tc.pY, tc.rAngle)

			// Start
			if !approximately(tc.expected.rayStart.x, got.rayStart.x) ||
				!approximately(tc.expected.rayStart.y, got.rayStart.y) {
				t.Errorf("Expected Start (%f, %f), got (%f, %f)",
					tc.expected.rayStart.x,
					tc.expected.rayStart.y,
					got.rayStart.x,
					got.rayStart.y)
			}

			// End
			if !approximately(tc.expected.rayEnd.x, got.rayEnd.x) ||
				!approximately(tc.expected.rayEnd.y, got.rayEnd.y) {
				t.Errorf("Expected End (%f, %f), got (%f, %f)",
					tc.expected.rayEnd.x,
					tc.expected.rayEnd.y,
					got.rayEnd.x,
					got.rayEnd.y)
			}

			// Angle
			if !approximately(tc.expected.rayAngle, got.rayAngle) {
				t.Errorf("Expected Angle %f, got %f", tc.expected.rayAngle, got.rayAngle)
			}

			// Length
			if !approximately(tc.expected.rayLength, got.rayLength) {
				t.Errorf("Expected Length %f, got %f", tc.expected.rayLength, got.rayLength)
			}

			if tc.expected.wallType != got.wallType {
				t.Errorf("Expected Walltype of %d, got %d", tc.expected.wallType, got.wallType)
			}
		})
	}
}

func Test_RendererCalculateHorizontalCollision(t *testing.T) {
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
	var tManager TextureManager = nil

	testCases := []struct {
		name     string
		pX, pY   float64
		rAngle   float64
		expected collisionDetail
	}{
		{
			name:   "1_1_pos_0_degrees_no_horizontal_wall_hit",
			pX:     1.0,
			pY:     1.0,
			rAngle: 0.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 2048.0,
					y: 2048.0,
				},
				rayAngle:  0.0,
				rayLength: 2048.0,
				wallType:  0,
			},
		},
		{
			name:   "7_3_pos_15_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 15.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 14.464102,
					y: 1.0,
				},
				rayAngle:  15.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_30_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 30.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 10.464102,
					y: 1.0,
				},
				rayAngle:  30.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_45_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 45.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 9.0,
					y: 1.0,
				},
				rayAngle:  45.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_60_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 60.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 8.154701,
					y: 1.0,
				},
				rayAngle:  60.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_75_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 75.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 7.535898,
					y: 1.0,
				},
				rayAngle:  75.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_90_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 90.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 7.0,
					y: 1.0,
				},
				rayAngle:  90.0,
				rayLength: 2.0,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_105_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 105.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 6.464102,
					y: 1.0,
				},
				rayAngle:  105.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_120_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 120.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 5.845299,
					y: 1.0,
				},
				rayAngle:  120.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_135_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 135.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 5.0,
					y: 1.0,
				},
				rayAngle:  135.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_150_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 150.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: 3.535898,
					y: 1.0,
				},
				rayAngle:  150.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "7_3_pos_165_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 165.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: -0.464102,
					y: 1.0,
				},
				rayAngle:  165.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "1_1_pos_180_degrees_no_horizontal_wall_hit",
			pX:     1.0,
			pY:     1.0,
			rAngle: 180.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 1.0,
					y: 1.0,
				},
				rayEnd: coordinates{
					x: 2048.0,
					y: 2048.0,
				},
				rayAngle:  180.0,
				rayLength: 2048.0,
				wallType:  0,
			},
		},
		{
			name:   "7_3_pos_195_degrees",
			pX:     7.0,
			pY:     3.0,
			rAngle: 195.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 3.0,
				},
				rayEnd: coordinates{
					x: -0.464102,
					y: 5.0,
				},
				rayAngle:  195.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_210_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 210.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 3.535898,
					y: 15.0,
				},
				rayAngle:  210.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_225_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 225.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 5.0,
					y: 15.0,
				},
				rayAngle:  225.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_240_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 240.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 5.845299,
					y: 15.0,
				},
				rayAngle:  240.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_255_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 255.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 6.464102,
					y: 15.0,
				},
				rayAngle:  255.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_270_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 270.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 7.0,
					y: 15.0,
				},
				rayAngle:  270.0,
				rayLength: 2.0,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_285_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 285.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 7.535898,
					y: 15.0,
				},
				rayAngle:  285.0,
				rayLength: 2.070552361,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_300_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 300.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 8.154701,
					y: 15.0,
				},
				rayAngle:  300.0,
				rayLength: 2.309401077,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_315_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 315.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 9.0,
					y: 15.0,
				},
				rayAngle:  315.0,
				rayLength: 2.828427125,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_330_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 330.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 10.464102,
					y: 15.0,
				},
				rayAngle:  330.0,
				rayLength: 4.0,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_345_degrees",
			pX:     7.0,
			pY:     13.0,
			rAngle: 345.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 14.464102,
					y: 15.0,
				},
				rayAngle:  345.0,
				rayLength: 7.72740661,
				wallType:  1,
			},
		},
		{
			name:   "7_13_pos_360_degrees_no_horizontal_wall_hit",
			pX:     7.0,
			pY:     13.0,
			rAngle: 360.0,
			expected: collisionDetail{
				rayStart: coordinates{
					x: 7.0,
					y: 13.0,
				},
				rayEnd: coordinates{
					x: 2048.0,
					y: 2048.0,
				},
				rayAngle:  360.0,
				rayLength: 2048.0,
				wallType:  0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := config.NewRenderConfiguration(FB_WIDTH, FB_HEIGHT, 64.0)
			r := NewRenderer(config, &game, tManager, data.LevelData{})

			got := r.computeHorizontalCollision(tc.pX, tc.pY, tc.rAngle)

			// Start
			if !approximately(tc.expected.rayStart.x, got.rayStart.x) ||
				!approximately(tc.expected.rayStart.y, got.rayStart.y) {
				t.Errorf("Expected Start (%f, %f), got (%f, %f)",
					tc.expected.rayStart.x,
					tc.expected.rayStart.y,
					got.rayStart.x,
					got.rayStart.y)
			}

			// End
			if !approximately(tc.expected.rayEnd.x, got.rayEnd.x) ||
				!approximately(tc.expected.rayEnd.y, got.rayEnd.y) {
				t.Errorf("Expected End (%f, %f), got (%f, %f)",
					tc.expected.rayEnd.x,
					tc.expected.rayEnd.y,
					got.rayEnd.x,
					got.rayEnd.y)
			}

			// Angle
			if !approximately(tc.expected.rayAngle, got.rayAngle) {
				t.Errorf("Expected Angle %f, got %f", tc.expected.rayAngle, got.rayAngle)
			}

			// Length
			if !approximately(tc.expected.rayLength, got.rayLength) {
				t.Errorf("Expected Length %f, got %f", tc.expected.rayLength, got.rayLength)
			}

			if tc.expected.wallType != got.wallType {
				t.Errorf("Expected Walltype of %d, got %d", tc.expected.wallType, got.wallType)
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
