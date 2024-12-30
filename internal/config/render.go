package config

type RenderConfiguration struct {
	fbWidth           int
	fbHeight          int
	fieldOfView       float64
	textureMapping    bool
	skyTextureMapping bool
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

func (r RenderConfiguration) GetFieldOfView() float64 {
	return r.fieldOfView
}

func (r RenderConfiguration) IsTextureMappingEnabled() bool {
	return r.textureMapping
}

func (r RenderConfiguration) IsSkyTextureMappingEnabled() bool {
	return r.skyTextureMapping
}

func (r *RenderConfiguration) EnableTextureMapping() {
	r.textureMapping = true
}

func (r *RenderConfiguration) DisableTextureMapping() {
	r.textureMapping = false
}

func (r *RenderConfiguration) EnableSkyTextureMapping() {
	r.skyTextureMapping = true
}

func (r *RenderConfiguration) DisableSkyTextureMapping() {
	r.skyTextureMapping = false
}
