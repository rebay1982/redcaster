package main

import (
	"fmt"
	"math"

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
}

type Renderer struct {
	game        *Game
	frameBuffer []uint8
}

// TODO: Loadable from file.
var gameMap = [16][16]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func (g *Game) update() {

}

func (r Renderer) calculateHeight(x int) (int, bool) {
	rayAngle := r.calculateRayAngle(x)

	posX, posY := r.game.playerX, r.game.playerY
	vLength := r.calculateVerticalCollisionRayLength(posX, posY, rayAngle)
	hLength := r.calculateHorizontalCollisionRayLength(posX, posY, rayAngle)
	var rLength float64 = vLength

	if hLength < vLength {
		rLength = hLength
	}

	if rLength > 1 {
		return FB_HEIGHT >> 1, true
	}

	return FB_HEIGHT >> 1, true
}

/*
Reference for RayAngle:

					90
					|
	 180 ---+--- 0/360
	 				|
				 270
*/
func (r Renderer) calculateRayAngle(x int) float64 {
	pAng := r.game.playerAngle

	xAngleRatio := r.game.fov / float64(FB_WIDTH)
	rayAngle := pAng + (r.game.fov / 2) - xAngleRatio*float64(x)

	if rayAngle < 0 {
		rayAngle += 360
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

	if r.game.gameMap[ix][iy] > 0 {
		return true

	} else {
		return false
	}
}

func (r Renderer) calculateVerticalCollisionRayLength(x, y, rAngle float64) float64 {
	// Convert the angle (in degrees) to radians because that's what the math library expects.
	rRad := rAngle * math.Pi / 180.0
	rLength := 2048.0

	// Increment X.
	if rAngle < 90.0 || rAngle > 270 {
		for i := 1; i < 16; i++ {
			// Coordinates of ray FROM initial position x, y
			rX := float64(int(x)+i) - x
			rY := math.Tan(rRad) * rX
			//fmt.Printf("A i: %d, rX: %f, rY: %f\n", i, rX, rY)

			// Substract rY because 0 on the Y axis is at the top. When moving X to the right (inc), Y will decrement when the
			//   ray's angle is between 0 and 90.
			if r.checkWallCollision(x+rX, y-rY) {
				rLength = rX / math.Cos(rRad)
				break
			}
		}
	}

	// Decrement X.
	if rAngle > 90.0 && rAngle < 270.0 {
		for i := 0; i < 16; i++ {
			fmt.Printf("B i: %d\n", i)
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

func (r Renderer) calculateHorizontalCollisionRayLength(x, y, rAngle float64) float64 {

	return 0.0
}

func NewRenderer(game *Game) *Renderer {
	r := &Renderer{
		game:        game,
		frameBuffer: make([]uint8, FB_WIDTH*FB_HEIGHT*4, FB_WIDTH*FB_HEIGHT*4),
	}

	return r
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
	h, _ := r.calculateHeight(x)
	startHeight := (FB_HEIGHT - h) >> 1

	for y := startHeight; y < (startHeight + h); y++ {
		colorIndex := (x + y*FB_WIDTH) * 4
		r.frameBuffer[colorIndex+1] = 0xFF // Green component
		r.frameBuffer[colorIndex+2] = 0xFF // Blue component
		r.frameBuffer[colorIndex+3] = 0xFF // Alpha
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
	//fmt.Println("Rendering screen: ")
	for x := 0; x < FB_WIDTH; x++ {
		r.drawVertical(x)
	}

	return r.frameBuffer
}

func main() {
	game := Game{
		playerX:     5.0,
		playerY:     5.0,
		playerAngle: 0.0,
		fov:         64.0, // 64 because each pixel column (640) will be equal to 0.1 degrees.
		gameMap: [16][16]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	renderer := NewRenderer(&game)

	config := rp.WindowConfig{
		Title:     WINDOW_TITLE,
		Width:     WINDOW_WIDTH,
		Height:    WINDOW_HEIGHT,
		Resizable: true,
		VSync:     true,
	}

	go func() {
		for {
			game.update()
		}
	}()

	rp.Init(config)
	rp.Run(nil, renderer.draw)
}
