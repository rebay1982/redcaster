package texture

import (
	"fmt"
	"unsafe"

	"github.com/rebay1982/redcaster/internal/config"
	"github.com/rebay1982/redcaster/internal/data"
	"github.com/rebay1982/redcaster/internal/utils"
)

type TextureManager struct {
	config                   config.RenderConfiguration
	textureData              []data.TextureData
	skyTextureData           []data.TextureData
	textureVerticalBuffer    []uint8
	skyTextureVerticalBuffer []uint8
}

func NewTextureManager(config config.RenderConfiguration, levelData data.LevelData) TextureManager {
	// Convert single sky texture to an array.
	skyTextures := []data.TextureData{}
	if levelData.SkyTextureFilename != "" {
		skyTextures = append(skyTextures, levelData.SkyTexture)
	}

	// Enable texture mapping by default.
	config.EnableTextureMapping()
	config.EnableSkyTextureMapping()

	manager := TextureManager{
		config:                   config,
		textureData:              levelData.Textures,
		skyTextureData:           skyTextures,
		textureVerticalBuffer:    make([]uint8, config.GetFbHeight()<<2), // *4 (4 bytes per pixel)
		skyTextureVerticalBuffer: make([]uint8, config.GetFbHeight()<<1), // /2 (half height) *4 (4 bytes per pixel)
	}

	// Only if we have texture data should we enable texture mapping, even if it was explicitly requested.
	if !(len(manager.textureData) > 0) {
		fmt.Println("WARN: Requested texture mapping, but no texture data found. Disabling texture mapping")
		manager.config.DisableTextureMapping()
	}

	if !(len(manager.skyTextureData) > 0) {
		fmt.Println("WARN: Requested sky texture mapping, but no texture data found. Disabling sky texture mapping")
		manager.config.DisableSkyTextureMapping()
	} else {
		manager.validateSkyTextureConfiguration()
	}

	return manager
}

func (tm *TextureManager) Reconfigure(config config.RenderConfiguration) {
	tm.config = config
}

func (tm TextureManager) validateSkyTextureConfiguration() {
	texData := tm.skyTextureData[0]

	utils.Assert(
		texData.Height >= (tm.config.GetFbHeight()>>1),
		fmt.Sprintf("Cannot use a sky texture whos height ([%d]) is less than half the configured frame buffer height ([%d]).\n",
			texData.Height,
			tm.config.GetFbHeight()>>1,
		),
	)

	totalVirtTexWidth := int(360 / tm.config.GetFieldOfView() * float64(tm.config.GetFbWidth()))
	utils.Assert(totalVirtTexWidth%texData.Width == 0,
		fmt.Sprintf("The sky texture width needs to be a multiple of [%d], which [%d] is not.\n",
			totalVirtTexWidth,
			texData.Width,
		),
	)
}

func (tm TextureManager) GetSkyTextureVertical(rAngle float64) []uint8 {
	skyVertBuffer := tm.skyTextureVerticalBuffer

	if tm.config.IsSkyTextureMappingEnabled() {
		skyTexData := tm.skyTextureData[0]
		skyTex := skyTexData.Data

		angle := 360 - rAngle // Flip the angle, the coordinate system is reverse to the FOV. Angles increment to the left,
		// while pixel positions decrement.

		// TODO: Compute this only once?
		virtTexWidth := 360 / tm.config.GetFieldOfView() * float64(tm.config.GetFbWidth())
		horizontalTexPosition := int(angle/360*virtTexWidth) % skyTexData.Width

		maxHeight := tm.config.GetFbHeight() >> 1
		for y := 0; y < maxHeight; y++ {
			vertBuffIndex := y << 2
			texIndex := (horizontalTexPosition + y*skyTexData.Width) << 2

			sTexSrc := (*uint32)(unsafe.Pointer(&skyTex[texIndex]))
			sTexVertBuffDst := (*uint32)(unsafe.Pointer(&skyVertBuffer[vertBuffIndex]))

			*sTexVertBuffDst = *sTexSrc
		}
	} else {
		for y := 0; y < len(skyVertBuffer); y += 4 {
			sTexVertBuffDst := (*uint32)(unsafe.Pointer(&skyVertBuffer[y]))
			*sTexVertBuffDst = 0xFFAA0000
		}
	}

	return skyVertBuffer
}

func (tm TextureManager) GetTextureVertical(textureId int, renderHeight int, texColumnCoord float64) []uint8 {
	texVertBuffer := tm.textureVerticalBuffer

	fullRH := renderHeight
	halfRH := fullRH >> 1
	fullTBH := len(texVertBuffer) >> 2
	halfTBH := fullTBH >> 1

	if tm.config.IsTextureMappingEnabled() {
		// Get texture data
		texture := tm.textureData[textureId-1]
		texHeight := texture.Height
		texWidth := texture.Width
		texColumn := int(float64(texWidth) * texColumnCoord)

		// Sampling ratio for the texture to texture vertical buffer
		texToTexVertBufferSampleRatio := float64(texHeight) / float64(renderHeight)

		// Samples the center of the vertical to outer edges. This way seems convoluted but actually simplifies the
		//	calculations quite a lot and always samples correctly whether the wall height to sample is smaller or larger than
		//  the frame buffer height (larger happens when the player is close up against a wall).
		for i := 0; i < halfRH && i < halfTBH; i++ {
			rhIndexNeg := (halfRH - i)
			rhIndexPos := (halfRH + i)
			tvbIndexNeg := (halfTBH - i)
			tvbIndexPos := (halfTBH + i)

			// Sample from texture
			textureRowNeg := int(float64(rhIndexNeg) * texToTexVertBufferSampleRatio)
			textureRowPos := int(float64(rhIndexPos) * texToTexVertBufferSampleRatio)

			// Sample from texture and write to texture vertical buffer.
			texPixIndex := (texColumn + (textureRowNeg * texWidth)) << 2
			tvbPixIndex := tvbIndexNeg << 2

			texSrc := (*uint32)(unsafe.Pointer(&texture.Data[texPixIndex]))
			texVertBuffDst := (*uint32)(unsafe.Pointer(&texVertBuffer[tvbPixIndex]))
			*texVertBuffDst = *texSrc

			texPixIndex = (texColumn + (textureRowPos * texWidth)) << 2
			tvbPixIndex = tvbIndexPos << 2

			texSrc = (*uint32)(unsafe.Pointer(&texture.Data[texPixIndex]))
			texVertBuffDst = (*uint32)(unsafe.Pointer(&texVertBuffer[tvbPixIndex]))
			*texVertBuffDst = *texSrc
		}
	} else {
		for i := 0; i < halfRH && i < halfTBH; i++ {
			tvbIndexNeg := (halfTBH - i)
			tvbIndexPos := (halfTBH + i)

			// Sample from texture and write to texture vertical buffer.
			tvbPixIndex := tvbIndexNeg << 2
			texVertBuffDst := (*uint32)(unsafe.Pointer(&texVertBuffer[tvbPixIndex]))
			*texVertBuffDst = 0xFFCCCCCC

			tvbPixIndex = tvbIndexPos << 2
			texVertBuffDst = (*uint32)(unsafe.Pointer(&texVertBuffer[tvbPixIndex]))
			*texVertBuffDst = 0xFFCCCCCC
		}
	}
	return texVertBuffer
}
