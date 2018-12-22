package framework

import (
	"github.com/liikt/GoNEAT/cppn"
	"github.com/liikt/GoNEAT/game"
	"github.com/liikt/GoNEAT/substrate"
)

var id uint64 = 2

type GlobalCTX struct {
	pool     *cppn.Pool
	curChamp *game.Player
	localctx []*LocalCTX
	numGames int
}

type LocalCTX struct {
	game   *game.Game
	player *game.Player
}

func NewGlobalCtx(size int) *GlobalCTX {
	return &GlobalCTX{numGames: 10, pool: cppn.InitPool(size)}
}

func (gctx *GlobalCTX) NewLocalCTX() *LocalCTX {
	p := game.NewPlayer(id)
	id++
	return &LocalCTX{game: game.NewAIGame(gctx.curChamp, p), player: p}
}

func (lctx *LocalCTX) runGame() {

}

func (gctx *GlobalCTX) runEpoch() {
	gctx.localctx = make([]*LocalCTX, 0)
	if id == 2 {
		gctx.curChamp = game.NewPlayer(1)
	}
	for _, g := range gctx.pool.GetGenomes() {
		lctx := gctx.NewLocalCTX()
		inNodes := [][]int{
			{2, 2}, {2, 3}, {2, 4},
			{3, 2}, {3, 3}, {3, 4},
			{4, 2}, {4, 3}, {4, 4},
		}
		outNodes := [][]int{
			{2, 6}, {2, 7}, {2, 8},
			{3, 6}, {3, 7}, {3, 8},
			{4, 6}, {4, 7}, {4, 8},
		}
		lctx.player.AI = substrate.BuildSubstrate(10, 10, inNodes, outNodes, g)
		gctx.localctx = append(gctx.localctx)
	}

}
