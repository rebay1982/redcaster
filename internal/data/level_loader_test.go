package data

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_DataLoader_DecodeLevelDataFile(t *testing.T) {
	testCases := []struct {
		name     string
		expected LevelData
		data     []byte
		err      bool
	}{
		{
			name: "basic_data_test",
			expected: LevelData{
				Name:   "test_data",
				Width:  3,
				Height: 3,
				Map: [][]int{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
				TextureMapping:   false,
				TextureFilenames: []string{},
				PlayerCoordData: PlayerCoordData{
					PlayerX:     1.0,
					PlayerY:     1.0,
					PlayerAngle: 45.0,
				},
			},
			data: []byte(`{
				"name": "test_data",
				"width": 3,
				"height": 3,
				"map": [
					[1, 1, 1],
					[1, 0, 1],
					[1, 1, 1]
				],
				"textureFilenames": [],
				"textures": [],
				"playerX": 1.0,
				"playerY": 1.0,
				"playerAngle": 45.0
			}`),
			err: false,
		},
		{
			name: "texture_data",
			expected: LevelData{
				Name:   "test_data",
				Width:  3,
				Height: 3,
				Map: [][]int{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
				TextureMapping: true,
				TextureFilenames: []string{
					"../../assets/test/test-black-pixel.png",
				},
				Textures: []TextureData{
					{
						Name:   "../../assets/test/test-black-pixel.png",
						Width:  1,
						Height: 1,
						Data:   []uint8{0x00, 0x00, 0x00, 0xFF},
					},
				},
				SkyTextureFilename: "../../assets/test/test-black-pixel.png",
				SkyTexture: TextureData{
					Name:   "../../assets/test/test-black-pixel.png",
					Width:  1,
					Height: 1,
					Data:   []uint8{0x00, 0x00, 0x00, 0xFF},
				},
				PlayerCoordData: PlayerCoordData{
					PlayerX:     1.0,
					PlayerY:     1.0,
					PlayerAngle: 45.0,
				},
			},
			data: []byte(`{
				"name": "test_data",
				"width": 3,
				"height": 3,
				"map": [
					[1, 1, 1],
					[1, 0, 1],
					[1, 1, 1]
				],
				"textures": [
					"../../assets/test/test-black-pixel.png"
				],
				"skyTexture": "../../assets/test/test-black-pixel.png",
				"playerX": 1.0,
				"playerY": 1.0,
				"playerAngle": 45.0
			}`),
			err: false,
		},
		{
			name:     "bad_data",
			expected: LevelData{},
			data: []byte(`{
				"data": "bad_data"
			}`),
			err: false,
		},
		{
			name:     "invalid_json",
			expected: LevelData{},
			data:     []byte(`This is not JSON`),
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dataLoader := DataLoader{}

			got, err := dataLoader.decodeLevelDataFile(tc.data)

			// Failed on error
			if tc.err && err == nil {
				t.Errorf("Expected err, got %v", err)
			}

			if !tc.err && err != nil {
				t.Errorf("Did not expect error, got %v", err)

			}

			if !cmp.Equal(tc.expected, got) {
				t.Errorf("Test failed\n%s\n", cmp.Diff(tc.expected, got))
			}
		})
	}
}
