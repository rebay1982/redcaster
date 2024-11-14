package data

import (
	"errors"
	"image"
	"image/png"
	"os"
)

type TextureLoader struct{}

func NewTextureLoader() TextureLoader {
	return TextureLoader{}
}

func (tl TextureLoader) LoadTextureData(filenames []string) ([]TextureData, error) {
	textureData := []TextureData{}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return textureData, err
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			return textureData, err
		}

		imgInfo := img.Bounds()
		rawTextureData, err := tl.getRawTextureData(img)
		if err != nil {
			return textureData, err
		}

		textureData = append(textureData, TextureData{
			Name:   filename,
			Width:  imgInfo.Max.X,
			Height: imgInfo.Max.Y,
			Data:   rawTextureData,
		})
	}

	return textureData, nil
}

func (tl TextureLoader) getRawTextureData(img image.Image) ([]byte, error) {
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		return nil, errors.New("Texture format is not RGBA")
	}

	return rgbaImg.Pix, nil
}
