package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rebay1982/redcaster/internal/config"
	"github.com/rebay1982/redcaster/internal/data"
	"github.com/rebay1982/redcaster/internal/game"
	"github.com/rebay1982/redcaster/internal/input"
	"github.com/rebay1982/redcaster/internal/render"
	"github.com/rebay1982/redcaster/internal/texture"

	rp "github.com/rebay1982/redpix"
)

func main() {
	appConfig := config.GetAppConfiguration()

	loader := data.NewDataLoader()
	levelData, err := loader.LoadLevelData(appConfig.DataFile)
	if err != nil {
		fmt.Printf("Failed to load data file %s. Aborting.\n", appConfig.DataFile)
		fmt.Printf("Caused by %v.\n", err)
		os.Exit(1)
	}

	inputHandler := input.NewInputHandler()
	game := game.NewGame(levelData, inputHandler)

	renderConfiguration := appConfig.RenderConfig

	textureManager := texture.NewTextureManager(renderConfiguration, levelData)
	renderer := render.NewRenderer(renderConfiguration, &game, &textureManager, levelData)

	winConfig := rp.WindowConfig{
		Title:     appConfig.WindowTitle,
		Width:     renderConfiguration.GetFbWidth(),
		Height:    renderConfiguration.GetFbHeight(),
		Resizable: true,
		VSync:     false,
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
