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

	loadedData.TextureMapping = len(loadedData.TextureFilenames) > 0
	if loadedData.TextureMapping {
		tl := NewTextureLoader()

		loadedData.Textures, err = tl.LoadTextureData(loadedData.TextureFilenames)
		if err != nil {
			return loadedData, err
		}
	}

	return loadedData, nil
}
