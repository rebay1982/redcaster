package render

import (
	"testing"

	"github.com/rebay1982/redcaster/internal/data"
)

func Test_RendererValidateSkyTextureConfiguration(t *testing.T) {
	testCases := []struct {
		name           string
		config         RenderConfiguration
		skyTextureData []data.TextureData
		expectPanic    bool
	}{
		{
			name: "valid_texture",
			config: RenderConfiguration{
				fbWidth:     100,
				fbHeight:    50,
				fieldOfView: 90.0,
			},
			skyTextureData: []data.TextureData{
				{
					Width:  100,
					Height: 100,
					Data:   []uint8{},
				},
			},
			expectPanic: false,
		},
		{
			name: "invalid_texture_height",
			config: RenderConfiguration{
				fbWidth:     100,
				fbHeight:    200,
				fieldOfView: 90.0,
			},
			skyTextureData: []data.TextureData{
				{
					Width:  100,
					Height: 50,
					Data:   []uint8{},
				},
			},
			expectPanic: true,
		},
		{
			name: "invalid_texture_width",
			config: RenderConfiguration{
				fbWidth:     100,
				fbHeight:    100,
				fieldOfView: 90.0,
			},
			skyTextureData: []data.TextureData{
				{
					Width:  123,
					Height: 50,
					Data:   []uint8{},
				},
			},
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if !tc.expectPanic && r != nil {
					t.Errorf("Not expecting a panic, recovered %v", r)
				}

				if tc.expectPanic && r == nil {
					t.Errorf("Expected a panic, none recovered.")
				}
			}()

			texMngr := TextureManager{
				config:         tc.config,
				skyTextureData: tc.skyTextureData,
			}

			texMngr.validateSkyTextureConfiguration()
		})
	}
}
