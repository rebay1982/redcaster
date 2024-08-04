package main

import (
	rp "github.com/rebay1982/redpix"
	"time"
	//"fmt"
)

const (
	WINDOW_TITLE = "RedCaster"
	WINDOW_WIDTH = 640
	WINDOW_HEIGHT = 480
	FB_WIDTH = 640
	FB_HEIGHT = 480
)

type Game struct {
	playerX, payerY float64
	clear bool
}

type Renderer struct {
	game *Game
	frameBuffer []uint8
}

// TODO: Loadable from file.
var gameMap = [16][16]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 
}

func (g *Game) update() {
	time.Sleep(1 * time.Second)
	g.clear = !g.clear
}

func NewRenderer(game *Game) *Renderer{
	r := &Renderer{
		game: game,
		frameBuffer: make([]uint8, FB_WIDTH * FB_HEIGHT * 4, FB_WIDTH * FB_HEIGHT * 4),
	}

	return r
}

func (r Renderer) drawCeiling() {
	height := FB_HEIGHT >> 1
	for x := 0; x < FB_WIDTH; x++ {
		for y := FB_HEIGHT -1; y >= height; y-- {
			colorIndex := (x + y * FB_WIDTH) * 4
			r.frameBuffer[colorIndex + 2] = 0xFF	// Blue skies component
			r.frameBuffer[colorIndex + 3] = 0xFF  // Alpha
		}
	}
}

func (r Renderer) drawFloor() {
	height := FB_HEIGHT >> 1
	for x := 0; x < FB_WIDTH; x++ {
		for y := height; y >= 0; y-- {
			colorIndex := (x + y * FB_WIDTH) * 4
			r.frameBuffer[colorIndex + 1] = 0x7F	// Green grass component
			r.frameBuffer[colorIndex + 3] = 0xFF  // Alpha
		}
	}
}

func (r Renderer) drawVertical(x int) {
	height := FB_HEIGHT / 2
	startHeight := (FB_HEIGHT - height) >> 1

	for y := startHeight; y < (startHeight + height); y++ {
		colorIndex := (x + y * FB_WIDTH) * 4
		r.frameBuffer[colorIndex + 1] = 0xFF	// Green component
		r.frameBuffer[colorIndex + 2] = 0xFF	// Blue component
		r.frameBuffer[colorIndex + 3] = 0xFF  // Alpha
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
	game := Game{0, 0, false}
	renderer := NewRenderer(&game)

	config := rp.WindowConfig{
		Title: WINDOW_TITLE,
		Width: WINDOW_WIDTH,
		Height: WINDOW_HEIGHT,
		Resizable: true,
		VSync: true,
	}

	go func(){ for {game.update()}}() 

	rp.Init(config)
	rp.Run(nil, renderer.draw)
}
