package config

import (
	"flag"

	"github.com/rebay1982/redcaster/internal/render"
)

const (
	WINDOW_TITLE  = "RedCaster"
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
	FOV           = 60.0
	DATA_FILE     = "./assets/demo/demo.json"
)

type AppConfig struct {
	WindowTitle string
	RenderConfig render.RenderConfiguration
	DataFile string	
}

func GetAppConfiguration() AppConfig {

	width := flag.Int("w", WINDOW_WIDTH, "Window width in pixels.")
	height := flag.Int("h", WINDOW_HEIGHT, "Window height in pixels.")
	fov := flag.Float64("fov", FOV, "Field of view in degrees.")
	file := flag.String("f", DATA_FILE, "File containing game data.")

	flag.Parse()

	return AppConfig{
		WindowTitle: WINDOW_TITLE,
		RenderConfig: render.NewRenderConfiguration(*width, *height, *fov),
		DataFile: *file,
	}
}
