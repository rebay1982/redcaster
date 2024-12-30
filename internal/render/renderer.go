package render

import (
	"math"

	"github.com/rebay1982/redcaster/internal/config"
	"github.com/rebay1982/redcaster/internal/data"
)

type TextureManager interface {
	Reconfigure(config config.RenderConfiguration)
	GetTextureVertical(textureId int, renderHeight int, texColumnCoord float64) []uint8
	GetSkyTextureVertical(rAngle float64) []uint8
}

type GameManager interface {
	GetPlayerCoords() data.PlayerCoordData
	CheckWallCollision(x, y float64) (bool, int)
}

type Renderer struct {
	gameManager   GameManager
	config        config.RenderConfiguration
	frameBuffer   []uint8
	rAngleOffsets []float64
	ambientLight  float64
	// TODO: Create a rendering memory manager
	textureManager TextureManager
}

// NewRenderer The game is a pointer because we want updates (from game) to the player position to be accessible.
func NewRenderer(config config.RenderConfiguration, gMngr GameManager, tMngr TextureManager, levelData data.LevelData) *Renderer {
	r := &Renderer{
		gameManager:  gMngr,
		config:       config,
		frameBuffer:  make([]uint8, config.ComputeFrameBufferSize(), config.ComputeFrameBufferSize()),
		ambientLight: levelData.AmbientLight,
	}
	r.precomputeRayAngleOffsets()
	r.textureManager = tMngr

	return r
}

func (r *Renderer) ReconfigureRenderer(config config.RenderConfiguration) {
	r.config = config
	r.frameBuffer = make([]uint8, config.ComputeFrameBufferSize(), config.ComputeFrameBufferSize())
	r.precomputeRayAngleOffsets()

	r.textureManager.Reconfigure(config)
}

func (r *Renderer) precomputeRayAngleOffsets() {
	fov := r.config.GetFieldOfView()
	r.rAngleOffsets = make([]float64, r.config.GetFbWidth())

	fRad := (fov / 2) * math.Pi / 180

	oppositeRefLength := math.Tan(fRad)
	oppositeStep := oppositeRefLength / float64(r.config.GetFbWidth()>>1)

	for i := 0; i < r.config.GetFbWidth(); i++ {
		r.rAngleOffsets[i] = math.Atan(oppositeRefLength-float64(i)*oppositeStep) * 180 / math.Pi
	}
}

func (r Renderer) applyLightingEffects(colorComponent uint8) uint8 {
	return uint8(float64(colorComponent) * r.ambientLight)
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
	palyerCoords := r.gameManager.GetPlayerCoords()
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
			if collision, wall := r.gameManager.CheckWallCollision(x+rX, y-rY); collision {
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
			if collision, wall := r.gameManager.CheckWallCollision(x-rX-0.001, y+rY); collision {
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
			if collision, wall := r.gameManager.CheckWallCollision(x+rX, y-rY-0.001); collision {
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
			if collision, wall := r.gameManager.CheckWallCollision(x-rX, y+rY); collision {
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

func (r Renderer) computeWallRenderingDetails(x int) wallRenderingDetail {
	height := float64(r.config.GetFbHeight())

	rayAngle := r.computeRayAngle(x)
	playerCoords := r.gameManager.GetPlayerCoords()
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

	textureVertical := r.textureManager.GetTextureVertical(tId, h, tCoord)
	for y := renderHeightStart; y < renderHeightEnd; y++ {

		// Texture pixels need to be drawn from bottom up because of flipped OpenGL coordinate system.
		//	(0, 0) is bottom left in OpenGL vs being top left in more intuitive coordinate systems.
		fbPixIndex := (x + (r.config.GetFbHeight()-1-y)*r.config.GetFbWidth()) << 2
		textureIndex := y << 2

		// We devide by two if the orientation is a vertical wall.
		r.frameBuffer[fbPixIndex] = r.applyLightingEffects(textureVertical[textureIndex] >> o)
		r.frameBuffer[fbPixIndex+1] = r.applyLightingEffects(textureVertical[textureIndex+1] >> o)
		r.frameBuffer[fbPixIndex+2] = r.applyLightingEffects(textureVertical[textureIndex+2] >> o)
		r.frameBuffer[fbPixIndex+3] = textureVertical[textureIndex+3]
	}
}

func (r Renderer) drawCeiling(x int) {
	rAngle := r.computeRayAngle(x)
	skyVertTexture := r.textureManager.GetSkyTextureVertical(rAngle)
	halfHeight := r.config.GetFbHeight() >> 1
	for y := r.config.GetFbHeight() - 1; y >= halfHeight; y-- {
		skyTexIndex := ((r.config.GetFbHeight() - 1) - y) << 2
		fbIndex := (x + y*r.config.GetFbWidth()) << 2
		r.frameBuffer[fbIndex] = skyVertTexture[skyTexIndex]
		r.frameBuffer[fbIndex+1] = skyVertTexture[skyTexIndex+1]
		r.frameBuffer[fbIndex+2] = skyVertTexture[skyTexIndex+2]
		r.frameBuffer[fbIndex+3] = skyVertTexture[skyTexIndex+3]
	}
}

func (r Renderer) drawFloor() {
	height := r.config.GetFbHeight() >> 1
	for x := 0; x < r.config.GetFbWidth(); x++ {
		for y := height; y >= 0; y-- {
			colorIndex := (x + y*r.config.GetFbWidth()) * 4
			r.frameBuffer[colorIndex] = r.applyLightingEffects(0x33)
			r.frameBuffer[colorIndex+1] = r.applyLightingEffects(0x33)
			r.frameBuffer[colorIndex+2] = r.applyLightingEffects(0x33)
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
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
	r.drawFloor()

	// Draw walls
	for x := 0; x < r.config.GetFbWidth(); x++ {
		r.drawCeiling(x)
		r.drawVertical(x)
	}

	return r.frameBuffer
}
