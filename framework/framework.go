package framework

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/liikt/GoNEAT/cppn"
	"github.com/liikt/GoNEAT/game"
	"github.com/liikt/GoNEAT/substrate"
)

var (
	id    uint64 = 2
	epoch        = 1

	inNodes = [][]int{
		{2, 2}, {2, 3}, {2, 4},
		{3, 2}, {3, 3}, {3, 4},
		{4, 2}, {4, 3}, {4, 4},
	}
	outNodes = [][]int{
		{2, 6}, {2, 7}, {2, 8},
		{3, 6}, {3, 7}, {3, 8},
		{4, 6}, {4, 7}, {4, 8},
	}
)

type GlobalCTX struct {
	pool     *cppn.Pool
	curChamp game.Player
	localctx []*LocalCTX
	numGames int
}

type LocalCTX struct {
	gctx        *GlobalCTX
	game        *game.Game
	player      *game.Player
	curChamp    *game.Player
	gamesPlayed int
	gamesWon    int
	gamesTied   int
}

func NewGlobalCtx(size int) *GlobalCTX {
	return &GlobalCTX{numGames: 10, pool: cppn.InitPool(size)}
}

func (gctx *GlobalCTX) NewLocalCTX() *LocalCTX {
	p := game.NewPlayer(id)
	id++
	var tmpChamp game.Player
	DeepCopy(&tmpChamp, gctx.curChamp)
	return &LocalCTX{gctx: gctx, game: game.NewAIGame(&tmpChamp, p), curChamp: &tmpChamp, player: p}
}

func (lctx *LocalCTX) runGame(wg *sync.WaitGroup) {
	var result *game.Result
	var curPlayer, waitingPlayer *game.Player
	numRoundsSurvived := 0
	for ; lctx.gamesPlayed < lctx.gctx.numGames; lctx.gamesPlayed++ {
		lctx.game.Reset()
		if coin := rand.Intn(2); coin == 1 {
			curPlayer = lctx.player
			waitingPlayer = lctx.curChamp
		} else {
			curPlayer = lctx.curChamp
			waitingPlayer = lctx.player
		}

		for result = lctx.game.DoMove(curPlayer); lctx.game.Round < 10 && result.Continue; result = lctx.game.DoMove(curPlayer) {
			curPlayer, waitingPlayer = waitingPlayer, curPlayer
		}

		if result.Won && curPlayer == lctx.player {
			fmt.Println("Player", lctx.player.ID, "won")
			lctx.gamesWon++
		} else if result.Continue && lctx.game.Round == 11 {
			fmt.Println("Player", lctx.player.ID, "tied")
			lctx.gamesTied++
		}
		numRoundsSurvived += lctx.game.Round - 1
	}
	lctx.player.Genome.Fitness = math.Sqrt(float64(lctx.gamesWon)/float64(lctx.gamesPlayed)) + math.Sqrt(float64(lctx.gamesTied)/float64(lctx.gamesPlayed))/4. + math.Sqrt(float64(numRoundsSurvived)/(float64(lctx.gamesPlayed)*4.))/10.
	if lctx.player.Genome.Fitness > 0.2 {
		fmt.Println("Player:", lctx.player.ID, "finished with", lctx.gamesWon, "wins and", lctx.gamesTied, "ties out of", lctx.gamesPlayed, "games.")
	}
	wg.Done()
}

func (gctx *GlobalCTX) startEpoch(champ game.Player) {
	gctx.localctx = make([]*LocalCTX, 0)
	genomes := gctx.pool.GetGenomes()
	fmt.Println("Started the building")
	ctxChan := make(chan *LocalCTX)
	for _, g := range genomes {
		go func(g *cppn.Genome) {
			lctx := gctx.NewLocalCTX()
			lctx.player.AI = *substrate.BuildSubstrate(10, 10, inNodes, outNodes, g)
			lctx.player.Genome = g
			lctx.game = game.NewAIGame(lctx.curChamp, lctx.player)
			ctxChan <- lctx
		}(g)
	}

	for idx := 0; idx < len(genomes); idx++ {
		if (idx+1)%10 == 0 {
			gctx.localctx = append(gctx.localctx, <-ctxChan)
			fmt.Println("Finished with genome number", idx+1)
		}
	}
}

func (gctx *GlobalCTX) Epoch() {
	for ; epoch <= 30; epoch++ {
		var champ game.Player
		if id == 2 {
			champ = *game.NewPlayer(1)
		} else if gctx.pool.ChangeChamp(0.75) {
			gctx.pool.RankGlobally()
			id++
			champ = *game.NewPlayer(id)
			champ.AI = *substrate.BuildSubstrate(10, 10, inNodes, outNodes, gctx.pool.GetGenomes()[0])
		} else {
			champ = gctx.curChamp
		}

		fmt.Println("============== Starting Epoch", epoch, "==============")

		if champ.Genome == nil {
			fmt.Println("First Round.")
		} else {
			fmt.Println("Champ Fitness:", champ.Genome.Fitness)
		}

		gctx.curChamp = champ
		gctx.startEpoch(champ)
		wg := new(sync.WaitGroup)
		wg.Add(len(gctx.localctx))

		for _, lctx := range gctx.localctx {
			go lctx.runGame(wg)
		}

		wg.Wait()
		fmt.Println("Old Stats")
		fmt.Println("Number of species:", len(gctx.pool.Species))
		fmt.Println("Number of genomes:", len(gctx.pool.GetGenomes()))
		fmt.Println("Innovations:", cppn.Innovation)
		gctx.pool.NewGeneration()
		fmt.Println("New Stats")
		fmt.Println("Number of species:", len(gctx.pool.Species))
		fmt.Println("Number of genomes:", len(gctx.pool.GetGenomes()))
		fmt.Println("Innovations:", cppn.Innovation)
		fmt.Println("============== Finished Epoch", epoch, "==============")
	}
}
