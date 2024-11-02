package game

import (
	"github.com/rebay1982/redcaster/internal/data"
	"github.com/rebay1982/redcaster/internal/input"
	"math"
)

type Game struct {
	playerCoords data.PlayerCoordData
	gameMap      [][]int
	inputHandler *input.InputHandler
}

func NewGame(levelData data.LevelData, inputHandler *input.InputHandler) Game {
	return Game{
		playerCoords: levelData.GetPlayerCoordData(),
		gameMap:      levelData.GetMapData(),
		inputHandler: inputHandler,
	}
}

func (g *Game) Update() {
	inputVector := g.inputHandler.GetInputVector()

	if inputVector.PlayerRight {
		g.playerCoords.PlayerAngle -= 0.3

		if g.playerCoords.PlayerAngle < 0.0 {
			g.playerCoords.PlayerAngle += 360.0
		}
	}

	if inputVector.PlayerLeft {
		g.playerCoords.PlayerAngle += 0.3

		if g.playerCoords.PlayerAngle > 360.0 {
			g.playerCoords.PlayerAngle -= 360.0
		}
	}

	pRad := g.playerCoords.PlayerAngle * math.Pi / 180.0
	deltaX := 0.01 * math.Cos(pRad)
	deltaY := 0.01 * math.Sin(pRad)
	if inputVector.PlayerForward {
		if !g.CheckWallCollision(g.playerCoords.PlayerX+deltaX, g.playerCoords.PlayerY-deltaY) {
			g.playerCoords.PlayerX += deltaX
			g.playerCoords.PlayerY -= deltaY
		}
	}

	if inputVector.PlayerBackward {
		if !g.CheckWallCollision(g.playerCoords.PlayerX-deltaX, g.playerCoords.PlayerY+deltaY) {
			g.playerCoords.PlayerX -= deltaX
			g.playerCoords.PlayerY += deltaY
		}
	}
}

// CheckWallCollision returns true if there's a wall at the given coordinates.
func (g Game) CheckWallCollision(x, y float64) bool {
	ix := int(x)
	iy := int(y)

	mapHeight := len(g.gameMap)
	mapWidth := len(g.gameMap[0]) // Assuming map is rectangular and all rows have the same length

	if ix < 0 || iy < 0 {
		return true
	}

	if ix >= mapWidth || iy >= mapHeight {
		return true
	}

	if g.gameMap[iy][ix] > 0 {
		return true

	} else {
		return false
	}
}

func (g Game) GetPlayerCoords() data.PlayerCoordData {
	return g.playerCoords
}
