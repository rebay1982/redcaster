package render

type RenderConfiguration struct {
	fbWidth     int
	fbHeight    int
	fieldOfView float64
}

func NewRenderConfiguration(width, height int, fov float64) RenderConfiguration {
	return RenderConfiguration{
		fbWidth:     width,
		fbHeight:    height,
		fieldOfView: fov,
	}
}

func (r RenderConfiguration) ComputeFrameBufferSize() int {
	return r.fbWidth * r.fbHeight * 4
}

func (r RenderConfiguration) GetFbWidth() int {
	return r.fbWidth
}

func (r RenderConfiguration) GetFbHeight() int {
	return r.fbHeight
}
