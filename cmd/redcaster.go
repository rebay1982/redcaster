package main

import (
	rp "github.com/rebay1982/redpix"
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

func (g Game) update() {

}

func NewRenderer(game *Game) *Renderer{
	r := &Renderer{
		game: game,
		frameBuffer: make([]uint8, FB_WIDTH * FB_HEIGHT * 4, FB_WIDTH * FB_HEIGHT * 4),
	}

	return r
}

func (r Renderer) drawVertical(x int) {
	height := FB_HEIGHT / 6
	startHeight := (FB_HEIGHT - height) >> 1

	//fmt.Printf("h: %d, sh: %d", height, startHeight)

	for y := startHeight; y < (startHeight + height); y++ {
		colorIndex := (x + y * FB_WIDTH) * 4
		r.frameBuffer[colorIndex + 1] = 0xFF	// Green component
		r.frameBuffer[colorIndex + 3] = 0xFF  // Alpha
	}
}

func (r Renderer) draw() []uint8 {
	for x := 0; x < FB_WIDTH; x++ {
		r.drawVertical(x)
	}

	return r.frameBuffer
}

func main() {
	game := Game{0, 0}
	renderer := NewRenderer(&game)

	config := rp.WindowConfig{
		Title: WINDOW_TITLE,
		Width: WINDOW_WIDTH,
		Height: WINDOW_HEIGHT,
		Resizable: true,
		VSync: true,
	}

	rp.Init(config)
	rp.Run(game.update, renderer.draw)
}
