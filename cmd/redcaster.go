package main

import (
	"time"

	"github.com/rebay1982/redcaster/internal/input"
	"github.com/rebay1982/redcaster/internal/game"
	render "github.com/rebay1982/redcaster/internal/renderer"

	rp "github.com/rebay1982/redpix"
)

const (
	WINDOW_TITLE  = "RedCaster"
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
	FB_WIDTH      = 640
	FB_HEIGHT     = 480
)

func main() {
	inputHandler := input.NewInputHandler()
	game := game.Game{
		PlayerX:     5.0,
		PlayerY:     5.0,
		PlayerAngle: 0.0,
		Fov:         64.0, // 64 because each pixel column (640) will be equal to 0.1 degrees.
		GameMap: [16][16]int{
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
		InputHandler: inputHandler,
	}
	renderer := render.NewRenderer(&game)

	winConfig := rp.WindowConfig{
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

			// Careful, we're updating game while the rendering loop is running. Might cause issues.
			game.Update()
			sinceLastCall = int(time.Since(start).Nanoseconds())

			time.Sleep(time.Duration(1600000 - sinceLastCall))
		}
	}()

	rp.Init(winConfig, renderer.Draw, inputHandler.HandleInputEvent)
	rp.Run()
}
