package game

import (
	"testing"

	"github.com/rebay1982/redcaster/internal/data"
)

func Test_RendererCheckWallCollision(t *testing.T) {
	levelData := data.LevelData{
		Map: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
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
	g := NewGame(levelData, nil)

	testCases := []struct {
		name      string
		x, y      float64
		want      bool
		wantWType int
	}{
		{
			name:      "x_y_no_collision",
			x:         1.0,
			y:         1.0,
			want:      false,
			wantWType: 0,
		},
		{
			name:      "x_y_diff_collision",
			x:         11.0,
			y:         2.0,
			want:      true,
			wantWType: 1,
		},
		{
			name:      "x_y_diff_float_collision",
			x:         11.8,
			y:         2.9,
			want:      true,
			wantWType: 1,
		},
		{
			name:      "x_y_float_no_collision",
			x:         1.20,
			y:         1.45,
			want:      false,
			wantWType: 0,
		},
		{
			name:      "x_y_float_close_no_collision",
			x:         1.99,
			y:         1.99,
			want:      false,
			wantWType: 0,
		},
		{
			name:      "x_y_collision",
			x:         0.0,
			y:         0.0,
			want:      true,
			wantWType: 1,
		},
		{
			name:      "x_y_float_collision",
			x:         2.20,
			y:         2.45,
			want:      true,
			wantWType: 1,
		},
		{
			name:      "x_negative_out_of_bound",
			x:         -1.0,
			y:         0.0,
			want:      true,
			wantWType: 0,
		},
		{
			name:      "x_positive_out_of_bound",
			x:         16.0,
			y:         0.0,
			want:      true,
			wantWType: 0,
		},
		{
			name:      "y_negative_out_of_bound",
			x:         0.0,
			y:         -1.0,
			want:      true,
			wantWType: 0,
		},
		{
			name:      "y_positive_out_of_bound",
			x:         0.0,
			y:         16.0,
			want:      true,
			wantWType: 0,
		},
		{
			name:      "x_y_negative_out_of_bound",
			x:         -1.0,
			y:         -1.0,
			want:      true,
			wantWType: 0,
		},
		{
			name:      "x_y_positive_out_of_bound",
			x:         16.0,
			y:         16.0,
			want:      true,
			wantWType: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			got, gotWType := g.CheckWallCollision(tc.x, tc.y)

			if got != tc.want {
				t.Errorf("Expected %t, got %t", tc.want, got)
			}
			if gotWType != tc.wantWType {
				t.Errorf("Expected wall type %d, got %d", tc.wantWType, gotWType)
			}
		})
	}
}
