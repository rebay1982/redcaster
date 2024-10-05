package main

import (
	"flag"
	"time"

	"github.com/rebay1982/redcaster/internal/game"
	"github.com/rebay1982/redcaster/internal/input"
	"github.com/rebay1982/redcaster/internal/render"

	rp "github.com/rebay1982/redpix"
)

const (
	WINDOW_TITLE  = "RedCaster"
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
	FOV           = 60.0
)

func initRendererConfiguration() render.RenderConfiguration {
	width := flag.Int("w", WINDOW_WIDTH, "Window width in pixels.")
	height := flag.Int("h", WINDOW_HEIGHT, "Window height in pixels.")
	fov := flag.Float64("f", FOV, "Field of view in degrees.")

	flag.Parse()
	return render.NewRenderConfiguration(*width, *height, *fov)
}

func main() {
	inputHandler := input.NewInputHandler()
	game := game.Game{
		PlayerX:     5.0,
		PlayerY:     5.0,
		PlayerAngle: 0.0,
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
	config := initRendererConfiguration()
	renderer := render.NewRenderer(config, &game)

	winConfig := rp.WindowConfig{
		Title:     WINDOW_TITLE,
		Width:     config.GetFbWidth(),
		Height:    config.GetFbHeight(),
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
