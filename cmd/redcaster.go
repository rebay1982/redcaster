package main

import (
	"math"
	"time"

	rp "github.com/rebay1982/redpix"
)

const (
	WINDOW_TITLE  = "RedCaster"
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
	FB_WIDTH      = 640
	FB_HEIGHT     = 480
)

type Game struct {
	playerX, playerY float64
	playerAngle      float64
	fov              float64
	gameMap          [16][16]int
	inputHandler     *InputHandler
}

type Renderer struct {
	game          *Game
	frameBuffer   []uint8
	rAngleOffsets []float64
}

type InputHandler struct {
	input InputVector
}

type InputVector struct {
	PlayerForward  bool
	PlayerBackward bool
	PlayerLeft     bool
	PlayerRight    bool
}

func (i *InputHandler) handleInputEvent(e rp.InputEvent) {
	// Ignore repeated key presses.
	if e.Action == rp.IN_ACT_REPEATED {
		return
	}

	setInput := false
	// Released handling is implicit.
	if e.Action == rp.IN_ACT_PRESSED {
		setInput = true
	}

	switch e.Key {
	case rp.IN_PLAYER_FORWARD:
		i.input.PlayerForward = setInput
	case rp.IN_PLAYER_BACKWARD:
		i.input.PlayerBackward = setInput
	case rp.IN_PLAYER_LEFT:
		i.input.PlayerLeft = setInput
	case rp.IN_PLAYER_RIGHT:
		i.input.PlayerRight = setInput
	}
}

func (i InputHandler) GetInputVector() InputVector {
	return i.input
}

func NewInputHandler() *InputHandler {
	i := &InputHandler{
		input: InputVector{},
	}

	return i
}

func NewRenderer(game *Game) *Renderer {
	r := &Renderer{
		game:        game,
		frameBuffer: make([]uint8, FB_WIDTH*FB_HEIGHT*4, FB_WIDTH*FB_HEIGHT*4),
	}

	r.precomputeRayAngleOffsets()

	return r
}

func (r Renderer) computeWallHeight(x int) (int, bool) {
	height := float64(FB_HEIGHT)
	horizontalOrientation := false

	rayAngle := r.computeRayAngle(x)

	posX, posY := r.game.playerX, r.game.playerY
	vLength := r.computeVerticalCollisionRayLength(posX, posY, rayAngle)
	hLength := r.calculateHorizontalCollisionRayLength(posX, posY, rayAngle)
	var rLength float64 = vLength

	if hLength < vLength {
		horizontalOrientation = true
		rLength = hLength
	}

	rLength = r.fishEyeCompensation(r.game.playerAngle, rayAngle, rLength)
	if rLength >= 1 {
		height = height / rLength
	}

	return int(height), horizontalOrientation
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
	pAng := r.game.playerAngle

	//xAngleRatio := r.game.fov / float64(FB_WIDTH)
	//rayAngle := pAng + (r.game.fov / 2) - xAngleRatio*float64(x)
	rayAngle := pAng + r.rAngleOffsets[x]

	if rayAngle < 0.0 {
		rayAngle += 360.0
	}

	if rayAngle > 360.0 {
		rayAngle -= 360.0
	}
	return rayAngle
}

// TODO: Make bounds configurable (map/world size)
func (r Renderer) checkWallCollision(x, y float64) bool {
	// Ray is out of bounds, can happen when a cast ray is close to being parallel to vertical or horizontal when
	//   computing collisions with vertical or horizontal lines.
	if x < 0 || y < 0 {
		return true
	}

	if x > 15 || y > 15 {
		return true
	}

	ix := int(x)
	iy := int(y)

	if r.game.gameMap[iy][ix] > 0 {
		return true

	} else {
		return false
	}
}

func (r Renderer) computeVerticalCollisionRayLength(x, y, rAngle float64) float64 {
	// Convert the angle (in degrees) to radians because that's what the math library expects.
	rRad := rAngle * math.Pi / 180.0
	rLength := 2048.0

	if rAngle < 90.0 || rAngle > 270 {
		for i := 1; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rX := float64(int(x)+i) - x
			rY := math.Tan(rRad) * rX

			// Substract rY because 0 on the Y axis is at the top. When moving X to the right (inc), Y will decrement when the
			//   ray's angle is between 0 and 90.
			if r.checkWallCollision(x+rX, y-rY) {
				rLength = rX / math.Cos(rRad)
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
			if r.checkWallCollision(x-rX-0.001, y+rY) {

				rLength = rX / math.Cos(rRad)
				break
			}
		}
	}

	return math.Abs(rLength)
}

func (r Renderer) computeHorizontalCollisionRayLength(x, y, rAngle float64) float64 {
	// Convert the angle (in degrees) to radians because that's what the math library expects.
	rRad := rAngle * math.Pi / 180.0
	rLength := 2048.0

	if rAngle > 0.0 && rAngle < 180.0 {
		for i := 0; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rY := y - float64(int(y)-i)
			rX := (1 / math.Tan(rRad)) * rY

			// -0.001 hack on y-yR necessary because collision checking is done on integer values (ex: >= 1, < 2). Ray should
			//   be < 1 if player is standing right next to a wall in an adjacent square.
			if r.checkWallCollision(x+rX, y-rY-0.001) {
				rLength = rY / math.Sin(rRad)
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
			if r.checkWallCollision(x-rX, y+rY) {
				rLength = rY / math.Sin(rRad)
				break
			}
		}
	}

	return math.Abs(rLength)
}

// FishEyeCompensation compensates for the perspective effect of the ray's angle relative to the player's direction.
func (r Renderer) fishEyeCompensation(pAngle, rAngle, rLength float64) float64 {
	rAngleToPlayer := math.Abs(rAngle - pAngle)
	rRad := rAngleToPlayer * math.Pi / 180.0

	return math.Abs(rLength * math.Cos(rRad))
}

func (r *Renderer) precomputeRayAngleOffsets() {
	fov := r.game.fov
	r.rAngleOffsets = make([]float64, FB_WIDTH)

	fRad := (fov / 2) * math.Pi / 180

	oppositeRefLength := math.Tan(fRad)
	oppositeStep := oppositeRefLength / float64(FB_WIDTH>>1)

	for i := 0; i < FB_WIDTH; i++ {
		r.rAngleOffsets[i] = math.Atan(oppositeRefLength-float64(i)*oppositeStep) * 180 / math.Pi
	}
}

func (r Renderer) drawCeiling() {
	height := FB_HEIGHT >> 1
	for x := 0; x < FB_WIDTH; x++ {
		for y := FB_HEIGHT - 1; y >= height; y-- {
			colorIndex := (x + y*FB_WIDTH) * 4
			r.frameBuffer[colorIndex+2] = 0xFF // Blue skies component
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
		}
	}
}

func (r Renderer) drawFloor() {
	height := FB_HEIGHT >> 1
	for x := 0; x < FB_WIDTH; x++ {
		for y := height; y >= 0; y-- {
			colorIndex := (x + y*FB_WIDTH) * 4
			r.frameBuffer[colorIndex+1] = 0x7F // Green grass component
			r.frameBuffer[colorIndex+3] = 0xFF // Alpha
		}
	}
}

func (r Renderer) drawVertical(x int) {
	h, o := r.computeWallHeight(x)
	startHeight := (FB_HEIGHT - h) >> 1

	for y := startHeight; y < (startHeight + h); y++ {
		colorIndex := (x + y*FB_WIDTH) * 4
		colorIntensity := 0xFF
		if o {
			colorIntensity = 0xAA
		}
		r.frameBuffer[colorIndex] = uint8(colorIntensity) // Green component
		r.frameBuffer[colorIndex+1] = 0x00                // Green component
		r.frameBuffer[colorIndex+2] = 0x00                // Blue component
		r.frameBuffer[colorIndex+3] = 0xFF                // Alpha
	}
}

func (r *Renderer) clearFrameBuffer() {
	r.frameBuffer[0] = 0x00
	for i := 1; i < len(r.frameBuffer); i = i << 1 {
		copy(r.frameBuffer[i:], r.frameBuffer[:i])
	}
}

func (r Renderer) draw() []uint8 {
	r.clearFrameBuffer()

	r.drawCeiling()
	r.drawFloor()

	// Draw walls
	for x := 0; x < FB_WIDTH; x++ {
		r.drawVertical(x)
	}

	return r.frameBuffer
}

func main() {
	inputHandler := NewInputHandler()
	game := Game{
		playerX:     5.0,
		playerY:     5.0,
		playerAngle: 0.0,
		fov:         64.0, // 64 because each pixel column (640) will be equal to 0.1 degrees.
		gameMap: [16][16]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
			{1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		inputHandler: inputHandler,
	}
	renderer := NewRenderer(&game)

	config := rp.WindowConfig{
		Title:     WINDOW_TITLE,
		Width:     WINDOW_WIDTH,
		Height:    WINDOW_HEIGHT,
		Resizable: true,
		VSync:     true,
	}

	// Update goroutine.
	go func() {
		sinceLastCall := 0
		var start time.Time
		for {
			start = time.Now()
			game.update()
			sinceLastCall = int(time.Since(start).Nanoseconds())

			time.Sleep(time.Duration(1600000 - sinceLastCall))
		}
	}()

	rp.Init(config, renderer.draw, inputHandler.handleInputEvent)
	rp.Run()
}

// TODO: This is copied from the renderer to validate collisions.
func (g Game) checkWallCollision(x, y float64) bool {
	if x < 0 || y < 0 {
		return true
	}

	if x > 15 || y > 15 {
		return true
	}

	ix := int(x)
	iy := int(y)

	if g.gameMap[iy][ix] > 0 {
		return true

	} else {
		return false
	}
}

// Careful, we're updating game while the rendering loop is running. Might cause issues.
func (g *Game) update() {
	inputVector := g.inputHandler.input

	if inputVector.PlayerRight {
		g.playerAngle -= 0.3

		if g.playerAngle < 0.0 {
			g.playerAngle += 360.0
		}
	}

	if inputVector.PlayerLeft {
		g.playerAngle += 0.3

		if g.playerAngle > 360.0 {
			g.playerAngle -= 360.0
		}
	}

	pRad := g.playerAngle * math.Pi / 180.0
	deltaX := 0.01 * math.Cos(pRad)
	deltaY := 0.01 * math.Sin(pRad)
	if inputVector.PlayerForward {
		if !g.checkWallCollision(g.playerX+deltaX, g.playerY-deltaY) {
			g.playerX += deltaX
			g.playerY -= deltaY
		}
	}

	if inputVector.PlayerBackward {
		if !g.checkWallCollision(g.playerX-deltaX, g.playerY+deltaY) {
			g.playerX -= deltaX
			g.playerY += deltaY
		}
	}
	//fmt.Printf("angle: %f, x: %f, y: %f\n", g.playerAngle, g.playerX, g.playerY)
}
