package data

import (
	"encoding/json"
	"os"
)

type DataLoader struct{}

func NewDataLoader() DataLoader {
	return DataLoader{}
}

func (dl DataLoader) LoadLevelData(filename string) (LevelData, error) {
	loadedData := LevelData{}

	dataFileContent, err := os.ReadFile(filename)
	if err != nil {
		return loadedData, err
	}

	return dl.decodeLevelDataFile(dataFileContent)
}

func (dl DataLoader) decodeLevelDataFile(content []byte) (LevelData, error) {
	loadedData := LevelData{}

	err := json.Unmarshal(content, &loadedData)
	if err != nil {
		return loadedData, err
	}

	if len(loadedData.TextureFilenames) > 0 {
		tl := NewTextureLoader()

		loadedData.Textures, err = tl.LoadTextureData(loadedData.TextureFilenames)
		if err != nil {
			return loadedData, err
		}

	}
	if loadedData.SkyTextureFilename != "" {
		tl := NewTextureLoader()

		skyTexture, err := tl.LoadTextureData([]string{loadedData.SkyTextureFilename})
		if err != nil {
			return loadedData, err
		}

		// Copy over the texture if we found one.
		if len(skyTexture) > 0 {
			loadedData.SkyTexture = skyTexture[0]
		}
	}

	return loadedData, nil
}
