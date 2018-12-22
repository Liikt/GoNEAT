package framework

import (
	"github.com/liikt/GoNEAT/cppn"
	"github.com/liikt/GoNEAT/game"
)

var id = 1

type GlobalCTX struct {
	pool *cppn.Pool
}

type LocalCTX struct {
	game    *game.Game
	player1 *game.Player
	player2 *game.Player
}

func NewGlobalCtx(size int) *GlobalCTX {
	return &GlobalCTX{pool: cppn.InitPool(size)}
}

func NewLocalCTX() *LocalCTX {
	p1 := game.NewPlayer(id)
	id++
	p2 := game.NewPlayer(id)
	id++
	return &LocalCTX{game: game.NewAIGame(p1, p2), player1: p1, player2: p2}
}
