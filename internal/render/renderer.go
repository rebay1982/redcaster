package render

import (
	"github.com/rebay1982/redcaster/internal/game"
	"math"
)

type Renderer struct {
	game          *game.Game
	config        RenderConfiguration
	frameBuffer   []uint8
	rAngleOffsets []float64
}

func NewRenderer(config RenderConfiguration, game *game.Game) *Renderer {
	r := &Renderer{
		game:        game,
		config:      config,
		frameBuffer: make([]uint8, config.ComputeFrameBufferSize(), config.ComputeFrameBufferSize()),
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

func (r Renderer) computeWallRenderingDetails(x int) wallRenderingDetail {
	// func (r Renderer) computeWallRenderingDetails(x int) (int, int, float64, bool) {
	height := float64(r.config.GetFbHeight())

	rayAngle := r.computeRayAngle(x)
	playerCoords := r.game.GetPlayerCoords()
	vCollision := r.computeVerticalCollision(playerCoords.PlayerX, playerCoords.PlayerY, rayAngle)
	hCollision := r.computeHorizontalCollision(playerCoords.PlayerX, playerCoords.PlayerY, rayAngle)

	wallType := vCollision.wallType
	wallOrientation := vCollision.wallOrientation
	collisionTextCoord := vCollision.rayEnd.y
	collisionRayLength := vCollision.rayLength

	if hCollision.rayLength < vCollision.rayLength {
		wallType = hCollision.wallType
		wallOrientation = hCollision.wallOrientation
		collisionTextCoord = hCollision.rayEnd.x
		collisionRayLength = hCollision.rayLength
	}

	// Fix the projection
	rLength := r.fishEyeCompensation(playerCoords.PlayerAngle, rayAngle, collisionRayLength)

	// If we're close to the wall, just display the complete height.
	if rLength >= 1 {
		height = height / rLength
	}

	return wallRenderingDetail{
		wallHeight:                    int(height),
		wallTextureId:                 wallType,
		wallOrientation:               wallOrientation,
		rayCollisionTextureCoordinate: collisionTextCoord,
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

func (r Renderer) computeTextureColumnCoordinate(textureId int, collisionY float64) int {

	return 0
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
			r.frameBuffer[colorIndex+2] = 0xCC // Blue skies component
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
		}
	}
}

func (r Renderer) drawFloor() {
	height := r.config.GetFbHeight() >> 1
	for x := 0; x < r.config.GetFbWidth(); x++ {
		for y := height; y >= 0; y-- {
			colorIndex := (x + y*r.config.GetFbWidth()) * 4
			r.frameBuffer[colorIndex] = 0x00
			r.frameBuffer[colorIndex+1] = 0x77 // Green grass component
			r.frameBuffer[colorIndex+2] = 0x00
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
		}
	}
}

func (r Renderer) drawVertical(x int) {
	renderingDetails := r.computeWallRenderingDetails(x)
	h := renderingDetails.wallHeight
	o := renderingDetails.wallOrientation

	startHeight := (r.config.GetFbHeight() - h) >> 1

	for y := startHeight; y < (startHeight + h); y++ {
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
