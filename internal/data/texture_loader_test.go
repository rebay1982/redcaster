package data

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTextureLoader_GetTextureData(t *testing.T) {
	tl := NewTextureLoader()

	tests := []struct {
		name      string
		filenames []string
		want      []TextureData
		wantErr   bool
	}{
		{
			name:      "GetTextureData_single_texture",
			filenames: []string{"../../assets/test/test-vertical-small.png"},
			want: []TextureData{
				{
					Name:   "../../assets/test/test-vertical-small.png",
					Width:  2,
					Height: 2,
					Data:   []uint8{
						0x00, 0x00, 0x00, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF,
						0x00, 0x00, 0x00, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "GetTextureData_dual_texture",
			filenames: []string{
				"../../assets/test/test-vertical-small.png",
				"../../assets/test/test-horizontal-small.png",
			},
			want: []TextureData{
				{
					Name:   "../../assets/test/test-vertical-small.png",
					Width:  2,
					Height: 2,
					Data:   []uint8{
						0x00, 0x00, 0x00, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF,
						0x00, 0x00, 0x00, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF,
					},
				},
				{
					Name:   "../../assets/test/test-horizontal-small.png",
					Width:  2,
					Height: 2,
					Data:   []uint8{
						0x00, 0x00, 0x00, 0xFF,
						0x00, 0x00, 0x00, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "GetTextureData_bad_filename",
			filenames: []string{
				"../../assets/test/i-dont-exist.png",
			},
			want: []TextureData{},
			wantErr: true,
		},
		{
			name:      "GetTextureData_bad_format",
			filenames: []string{
				"../../assets/test/test-bad-format.png",
			},
			want: []TextureData{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tl.LoadTextureData(tc.filenames)

			if tc.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tc.wantErr && err != nil {
				t.Errorf("Not expecting error, got %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Failed to validate return value: -want +got:\n%s", diff)
			}
		})
	}
}
