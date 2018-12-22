package framework

import (
	"github.com/liikt/GoNEAT/cppn"
	"github.com/liikt/GoNEAT/game"
)

var id uint64 = 2

type GlobalCTX struct {
	pool     *cppn.Pool
	curChamp *game.Player
}

type LocalCTX struct {
	game   *game.Game
	player *game.Player
}

func NewGlobalCtx(size int) *GlobalCTX {
	return &GlobalCTX{pool: cppn.InitPool(size)}
}

func (gctx *GlobalCTX) NewLocalCTX() *LocalCTX {
	p := game.NewPlayer(id)
	id++
	return &LocalCTX{game: game.NewAIGame(gctx.curChamp, p), player: p}
}
