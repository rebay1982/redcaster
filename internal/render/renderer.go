package render

import (
	"fmt"
	"github.com/rebay1982/redcaster/internal/data"
	"github.com/rebay1982/redcaster/internal/game"
	"math"
)

type Renderer struct {
	game          *game.Game
	config        RenderConfiguration
	frameBuffer   []uint8
	rAngleOffsets []float64
	// TODO: Create a rendering memory manager
	textureData           []data.TextureData
	textureVerticalBuffer []uint8 // Reusable texture vertical buffer to sample textures to.
}

func NewRenderer(config RenderConfiguration, game *game.Game, textureData []data.TextureData) *Renderer {
	// Enable texture mapping by default.
	config.textureMapping = true
	r := &Renderer{
		game:                  game,
		config:                config,
		frameBuffer:           make([]uint8, config.ComputeFrameBufferSize(), config.ComputeFrameBufferSize()),
		textureData:           textureData,
		textureVerticalBuffer: make([]byte, config.fbHeight<<2),
	}

	// Only if we have texture data should we enable texture mapping, even if it was explicitly requested.
	if !(len(textureData) > 0) {
		fmt.Println("WARN: Requested texture mapping, but no texture data found. Disabling texture mapping")
		r.config.textureMapping = false
	}

	r.precomputeRayAngleOffsets()

	return r
}

func (r Renderer) ReconfigureRenderer(config RenderConfiguration) {
	r.config = config
	r.frameBuffer = make([]uint8, config.ComputeFrameBufferSize(), config.ComputeFrameBufferSize())
	r.precomputeRayAngleOffsets()
}

func (r *Renderer) precomputeRayAngleOffsets() {
	fov := r.config.fieldOfView
	r.rAngleOffsets = make([]float64, r.config.GetFbWidth())

	fRad := (fov / 2) * math.Pi / 180

	oppositeRefLength := math.Tan(fRad)
	oppositeStep := oppositeRefLength / float64(r.config.GetFbWidth()>>1)

	for i := 0; i < r.config.GetFbWidth(); i++ {
		r.rAngleOffsets[i] = math.Atan(oppositeRefLength-float64(i)*oppositeStep) * 180 / math.Pi
	}
}

/*
Reference for RayAngle:

					90
					|
	 180 ---+--- 0/360
	 				|
				 270
*/
func (r Renderer) computeRayAngle(x int) float64 {
	palyerCoords := r.game.GetPlayerCoords()
	pAng := palyerCoords.PlayerAngle

	rayAngle := pAng + r.rAngleOffsets[x]

	if rayAngle < 0.0 {
		rayAngle += 360.0
	}

	if rayAngle > 360.0 {
		rayAngle -= 360.0
	}
	return rayAngle
}

func (r Renderer) computeVerticalCollision(x, y, rAngle float64) collisionDetail {
	// Note: We don't check for rAngle == 90 or 270. This is because hen checking vertical wall collisions, the ray will
	//       never intersect with one when it is projected at 90 or 270 degrees.

	// Convert the angle (in degrees) to radians because that's what the math library expects.
	rRad := rAngle * math.Pi / 180.0
	rLength := 2048.0
	wType := 0
	startCoords := coordinates{
		x: x,
		y: y,
	}
	endCoords := coordinates{
		x: 2048.0,
		y: 2048.0,
	}

	if rAngle < 90.0 || rAngle > 270 {
		for i := 1; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rX := float64(int(x)+i) - x
			rY := math.Tan(rRad) * rX

			// Substract rY because 0 on the Y axis is at the top. When moving X to the right (inc), Y will decrement when the
			//   ray's angle is between 0 and 90.
			if collision, wall := r.game.CheckWallCollision(x+rX, y-rY); collision {
				rLength = rX / math.Cos(rRad)
				wType = wall
				endCoords.x = x + rX
				endCoords.y = y - rY

				break
			}
		}
	}

	if rAngle > 90.0 && rAngle < 270.0 {
		for i := 0; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rX := x - float64(int(x)-i)
			rY := math.Tan(rRad) * rX

			// Substract rX because 0 on the X axis is at the far left. When the ray's angle is between 90 and 270, X is
			//   moving to the left thus decrementing X.

			// -0.001 hack on x-xR necessary because collision checking is done on integer values (ex: >= 1, < 2). Ray should
			//   be < 1 if player is standing right next to a wall in an adjacent square.
			if collision, wall := r.game.CheckWallCollision(x-rX-0.001, y+rY); collision {
				rLength = rX / math.Cos(rRad)
				wType = wall
				endCoords.x = x - rX
				endCoords.y = y + rY

				break
			}
		}
	}

	return collisionDetail{
		rayStart:        startCoords,
		rayEnd:          endCoords,
		rayAngle:        rAngle,
		rayLength:       math.Abs(rLength),
		wallType:        wType,
		wallOrientation: 0, // Vertical wall collision
	}
}

func (r Renderer) computeHorizontalCollision(x, y, rAngle float64) collisionDetail {
	// Note: We don't check for rAngle == 0 or 180. This is because hen checking horizontal wall collisions, the ray will
	//       never intersect with one when it is projected at 0 or 180 degrees.

	// Convert the angle (in degrees) to radians because that's what the math library expects.
	rRad := rAngle * math.Pi / 180.0
	rLength := 2048.0
	wType := 0
	startCoords := coordinates{
		x: x,
		y: y,
	}
	endCoords := coordinates{
		x: 2048.0,
		y: 2048.0,
	}

	if rAngle > 0.0 && rAngle < 180.0 {
		for i := 0; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rY := y - float64(int(y)-i)
			rX := (1 / math.Tan(rRad)) * rY

			// -0.001 hack on y-yR necessary because collision checking is done on integer values (ex: >= 1, < 2). Ray should
			//   be < 1 if player is standing right next to a wall in an adjacent square.
			if collision, wall := r.game.CheckWallCollision(x+rX, y-rY-0.001); collision {
				rLength = rY / math.Sin(rRad)
				wType = wall
				endCoords.x = x + rX
				endCoords.y = y - rY

				break
			}
		}
	}

	if rAngle > 180.0 && rAngle < 360.0 {
		for i := 1; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rY := float64(int(y)+i) - y
			rX := (1 / math.Tan(rRad)) * rY

			// Substract rX because the Tangent is negative from 270 to 360 and positive from 180 to 270, which is the
			//	 opposite of our reference coordinate system. (it is negative from 180 to 270 and positive from 270 to 360).
			if collision, wall := r.game.CheckWallCollision(x-rX, y+rY); collision {
				rLength = rY / math.Sin(rRad)
				wType = wall
				endCoords.x = x - rX
				endCoords.y = y + rY

				break
			}
		}
	}

	return collisionDetail{
		rayStart:        startCoords,
		rayEnd:          endCoords,
		rayAngle:        rAngle,
		rayLength:       math.Abs(rLength),
		wallType:        wType,
		wallOrientation: 1, // Horizontal wall collision
	}
}

// FishEyeCompensation compensates for the perspective effect of the ray's angle relative to the player's direction.
func (r Renderer) fishEyeCompensation(pAngle, rAngle, rLength float64) float64 {
	rAngleToPlayer := math.Abs(rAngle - pAngle)
	rRad := rAngleToPlayer * math.Pi / 180.0

	return math.Abs(rLength * math.Cos(rRad))
}

func (r Renderer) drawCeiling() {
	height := r.config.GetFbHeight() >> 1
	for x := 0; x < r.config.GetFbWidth(); x++ {
		for y := r.config.GetFbHeight() - 1; y >= height; y-- {
			colorIndex := (x + y*r.config.GetFbWidth()) * 4
			r.frameBuffer[colorIndex] = 0x00
			r.frameBuffer[colorIndex+1] = 0x00
			r.frameBuffer[colorIndex+2] = 0x33 // Blue skies component
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
		}
	}
}

func (r Renderer) drawFloor() {
	height := r.config.GetFbHeight() >> 1
	for x := 0; x < r.config.GetFbWidth(); x++ {
		for y := height; y >= 0; y-- {
			colorIndex := (x + y*r.config.GetFbWidth()) * 4
			r.frameBuffer[colorIndex] = 0x33
			r.frameBuffer[colorIndex+1] = 0x33
			r.frameBuffer[colorIndex+2] = 0x33
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
		}
	}
}

func (r Renderer) computeWallRenderingDetails(x int) wallRenderingDetail {
	height := float64(r.config.GetFbHeight())

	rayAngle := r.computeRayAngle(x)
	playerCoords := r.game.GetPlayerCoords()
	vCollision := r.computeVerticalCollision(playerCoords.PlayerX, playerCoords.PlayerY, rayAngle)
	hCollision := r.computeHorizontalCollision(playerCoords.PlayerX, playerCoords.PlayerY, rayAngle)

	wallType := vCollision.wallType
	wallOrientation := vCollision.wallOrientation
	collisionRayLength := vCollision.rayLength

	// We're only really interested in the factional part of collision coordinate because textures are mapped between
	//	0 and 1. It has no value to keep the absolute world value of the collision.
	var relCollisionTexCoord float64

	// Validate if were computing the texture collision coordinate for WEST vertical walls. If so, flip the texture
	//	coordinate so that the normal of the wall is facing towards the player and the texture renders in the correct
	//	orientation. Failing to do this results in mirrored texture on the vertical axis.
	if vCollision.rayEnd.x < vCollision.rayStart.x {
		frac := vCollision.rayEnd.y - float64(int(vCollision.rayEnd.y))
		relCollisionTexCoord = 0.999999 - frac
	} else {
		relCollisionTexCoord = vCollision.rayEnd.y - float64(int(vCollision.rayEnd.y))
	}

	if hCollision.rayLength < vCollision.rayLength {
		wallType = hCollision.wallType
		wallOrientation = hCollision.wallOrientation
		collisionRayLength = hCollision.rayLength

		// Validate if were computing the texture collision coordinate for SOUTH horitontal walls. If so, flip the texture
		//	coordinate so that the normal of the wall is facing towards the player and the texture renders in the correct
		//	orientation. Failing to do this results in mirrored texture on the vertical axis.
		if hCollision.rayEnd.y > hCollision.rayStart.y {
			frac := hCollision.rayEnd.x - float64(int(hCollision.rayEnd.x))
			relCollisionTexCoord = 0.999999 - frac
		} else {
			relCollisionTexCoord = hCollision.rayEnd.x - float64(int(hCollision.rayEnd.x))
		}
	}

	// Fix the projection
	rLength := r.fishEyeCompensation(playerCoords.PlayerAngle, rayAngle, collisionRayLength)

	// Height will exceed the frame buffer height if we're closer than a ray length of 1 from the wall. This can be locked
	//	down to FBHeight when texture mapping is diabled since we're applying solid colours.
	height = height / rLength

	return wallRenderingDetail{
		wallHeight:                    int(height),
		wallDistance:                  rLength,
		wallTextureId:                 wallType,
		wallOrientation:               wallOrientation,
		rayCollisionTextureCoordinate: relCollisionTexCoord,
	}
}

func (r Renderer) getTextureVerticalToRender(textureId int, renderHeight int, texColumnCoord float64) []uint8 {
	texVertBuffer := r.textureVerticalBuffer

	// Get texture data
	texture := r.textureData[textureId-1]
	texHeight := texture.Height
	texWidth := texture.Width
	texColumn := int(float64(texWidth) * texColumnCoord)

	// Sampling ratio for the texture to texture vertical buffer
	texToTexVertBufferSampleRatio := float64(texHeight) / float64(renderHeight)

	fullRH := renderHeight
	halfRH := fullRH >> 1
	fullTBH := len(texVertBuffer) >> 2
	halfTBH := fullTBH >> 1

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
		texVertBuffer[tvbPixIndex] = texture.Data[texPixIndex]
		texVertBuffer[tvbPixIndex+1] = texture.Data[texPixIndex+1]
		texVertBuffer[tvbPixIndex+2] = texture.Data[texPixIndex+2]
		texVertBuffer[tvbPixIndex+3] = texture.Data[texPixIndex+3]

		texPixIndex = (texColumn + (textureRowPos * texWidth)) << 2
		tvbPixIndex = tvbIndexPos << 2
		texVertBuffer[tvbPixIndex] = texture.Data[texPixIndex]
		texVertBuffer[tvbPixIndex+1] = texture.Data[texPixIndex+1]
		texVertBuffer[tvbPixIndex+2] = texture.Data[texPixIndex+2]
		texVertBuffer[tvbPixIndex+3] = texture.Data[texPixIndex+3]
	}

	return texVertBuffer
}

func (r Renderer) drawVertical(x int) {
	renderingDetails := r.computeWallRenderingDetails(x)
	h := renderingDetails.wallHeight
	o := renderingDetails.wallOrientation
	tId := renderingDetails.wallTextureId
	tCoord := renderingDetails.rayCollisionTextureCoordinate

	renderHeightStart := (r.config.GetFbHeight() - h) >> 1
	renderHeightEnd := (renderHeightStart + h)

	if renderHeightStart < 0 {
		renderHeightStart = 0
		renderHeightEnd = r.config.GetFbHeight()
	}

	// TODO: Make drawing the vertical independent of texture mapping being enabled or disabled.
	//			 ie, make the "getTextureVerticalToRender" handle this.
	if r.config.IsTextureMappingEnabled() {
		textureColumn := r.getTextureVerticalToRender(tId, h, tCoord)

		for y := renderHeightStart; y < renderHeightEnd; y++ {

			// Texture pixels need to be drawn from bottom up because of flipped OpenGL coordinate system.
			//	(0, 0) is bottom left in OpenGL vs being top left in more intuitive coordinate systems.
			pixIndex := (x + (r.config.GetFbHeight()-1-y)*r.config.GetFbWidth()) << 2
			textureIndex := y << 2

			// We devide by two if the orientation is a vertical wall.
			r.frameBuffer[pixIndex] = textureColumn[textureIndex] >> o
			r.frameBuffer[pixIndex+1] = textureColumn[textureIndex+1] >> o
			r.frameBuffer[pixIndex+2] = textureColumn[textureIndex+2] >> o
			r.frameBuffer[pixIndex+3] = textureColumn[textureIndex+3]
		}
	} else {
		for y := renderHeightStart; y < (renderHeightStart + h); y++ {
			// Special case when the height of the wall is > than the framebuffer.
			// Can and will happen when the player is close enough to a wall that the collision ray length is < 1.
			// Skip if outside top portion of framebuffer
			if y < 0 {
				continue
			}

			// Quit when done drawing whole frame buffer.
			if y >= r.config.GetFbHeight() {
				break
			}

			colorIndex := (x + y*r.config.GetFbWidth()) * 4
			colorIntensity := 0xCC
			if o == 1 {
				colorIntensity = 0x88
			}
			r.frameBuffer[colorIndex] = uint8(colorIntensity)
			r.frameBuffer[colorIndex+1] = uint8(colorIntensity)
			r.frameBuffer[colorIndex+2] = uint8(colorIntensity)
			r.frameBuffer[colorIndex+3] = 0xFF
		}
	}
}

func (r *Renderer) clearFrameBuffer() {
	r.frameBuffer[0] = 0x00
	for i := 1; i < len(r.frameBuffer); i = i << 1 {
		copy(r.frameBuffer[i:], r.frameBuffer[:i])
	}
}

// Draw draws the game to the frame buffer.
func (r Renderer) Draw() []uint8 {
	//r.clearFrameBuffer()

	r.drawCeiling()
	r.drawFloor()

	// Draw walls
	for x := 0; x < r.config.GetFbWidth(); x++ {
		r.drawVertical(x)
	}

	return r.frameBuffer
}
