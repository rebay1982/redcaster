package render

type coordinates struct {
	x float64
	y float64
}

type collisionDetail struct {
	rayStart  coordinates
	rayEnd    coordinates
	rayAngle  float64
	rayLength float64

	wallType        int
	wallOrientation int
}

type wallRenderingDetail struct {
	wallHeight      int
	wallTextureId   int
	wallOrientation int

	rayCollisionTextureCoordinate float64
}
