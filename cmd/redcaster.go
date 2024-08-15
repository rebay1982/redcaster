package main

import (
	"fmt"

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

func (r Renderer) calculateHeight(x int) int {
	rayAngle := r.calculateRayAngle(x)

	// Vertical line collision check

	// While collision with wall (vertical)

	// Horizontal line collision check
	posX, posY := r.game.playerX, r.game.playerY
	if posY >= 0 && posX >= 0 && x >= 0 {
		return FB_HEIGHT >> 1
	}

	return FB_HEIGHT >> 1
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
	rayAngle := pAng + r.game.fov - xAngleRatio*float64(x)

	if rayAngle < 0 {
		rayAngle += 360
	}

	return rayAngle
}

func (r Renderer) checkWallCollision(x, y float64) bool {

	return false
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
	h := r.calculateHeight(x)
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
