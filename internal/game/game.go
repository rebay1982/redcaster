package game

import (
	"math"
	"github.com/rebay1982/redcaster/internal/input"
)

type Game struct {
	PlayerX, PlayerY float64
	PlayerAngle      float64
	Fov              float64
	GameMap          [16][16]int
	InputHandler     *input.InputHandler
}

func (g *Game) Update() {
	inputVector := g.InputHandler.GetInputVector()

	if inputVector.PlayerRight {
		g.PlayerAngle -= 0.3

		if g.PlayerAngle < 0.0 {
			g.PlayerAngle += 360.0
		}
	}

	if inputVector.PlayerLeft {
		g.PlayerAngle += 0.3

		if g.PlayerAngle > 360.0 {
			g.PlayerAngle -= 360.0
		}
	}

	pRad := g.PlayerAngle * math.Pi / 180.0
	deltaX := 0.01 * math.Cos(pRad)
	deltaY := 0.01 * math.Sin(pRad)
	if inputVector.PlayerForward {
		if !g.CheckWallCollision(g.PlayerX+deltaX, g.PlayerY-deltaY) {
			g.PlayerX += deltaX
			g.PlayerY -= deltaY
		}
	}

	if inputVector.PlayerBackward {
		if !g.CheckWallCollision(g.PlayerX-deltaX, g.PlayerY+deltaY) {
			g.PlayerX -= deltaX
			g.PlayerY += deltaY
		}
	}
}

// CheckWallCollision returns true if there's a wall at the given coordinates.
func (g Game) CheckWallCollision(x, y float64) bool {
	if x < 0 || y < 0 {
		return true
	}

	if x > 15 || y > 15 {
		return true
	}

	ix := int(x)
	iy := int(y)

	if g.GameMap[iy][ix] > 0 {
		return true

	} else {
		return false
	}
}

