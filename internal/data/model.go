package data

type LevelData struct {
	Name   string  `json:"name"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Map    [][]int `json:"map"`

	TextureMapping bool

	// Normal wall textures
	TextureFilenames []string `json:"textures"`
	Textures         []TextureData

	// Sky texture
	SkyTextureFilename string `json:"skyTexture"`
	SkyTexture         TextureData

	AmbientLight float64 `json:"ambientLight"`

	PlayerCoordData
}

type TextureData struct {
	Name   string
	Width  int
	Height int
	Data   []uint8
}

type PlayerCoordData struct {
	PlayerX     float64 `json:"playerX"`
	PlayerY     float64 `json:"playerY"`
	PlayerAngle float64 `json:"playerAngle"`
}

func (ld LevelData) GetPlayerCoordData() PlayerCoordData {
	return ld.PlayerCoordData
}

func (ld LevelData) GetMapData() [][]int {
	return ld.Map
}
