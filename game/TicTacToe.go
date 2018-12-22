package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/liikt/GoNEAT/substrate"
)

type Player struct {
	ID          uint64
	AI          *substrate.Substrate
	GamesPlayed int
	GamesWon    int
}
type Game struct {
	Round   int
	Board   []uint64
	Convert map[uint64]string
	Player1 *Player
	Player2 *Player
}

type result struct {
	Continue bool
	Won      bool
	Round    int
}

func containsInt(arr []int, i int) bool {
	for _, x := range arr {
		if i == x {
			return true
		}
	}
	return false
}

func newGame() *Game {
	g := new(Game)
	g.Convert = map[uint64]string{0: " ", 1: "O", 2: "X"}
	g.Board = []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0}
	return g
}

func NewPlayer(id uint64) *Player {
	return &Player{ID: id}
}

func NewAIGame(p1, p2 *Player) *Game {
	if p2.ID < p1.ID {
		p1, p2 = p2, p1
	}
	g := new(Game)
	g.Round = 0
	g.Player1 = p1
	g.Player2 = p2
	g.Convert = map[uint64]string{0: " ", p1.ID: "O", p2.ID: "X"}
	g.Board = []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0}
	return g
}

func (g *Game) Reset() {
	g.Board = []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (g *Game) printField(player uint64) {
	fmt.Printf("%v|%v|%v\n", g.Convert[g.Board[0]], g.Convert[g.Board[1]], g.Convert[g.Board[2]])
	fmt.Printf("-+-+-\n")
	fmt.Printf("%v|%v|%v\n", g.Convert[g.Board[3]], g.Convert[g.Board[4]], g.Convert[g.Board[5]])
	fmt.Printf("-+-+-\n")
	fmt.Printf("%v|%v|%v\n", g.Convert[g.Board[6]], g.Convert[g.Board[7]], g.Convert[g.Board[8]])
	fmt.Printf("Player: %v\n", player+1)
}

func index(arr []int, el int) int {
	for x, v := range arr {
		if el == v {
			return x
		}
	}
	return -1
}

func (g *Game) checkWin(move int) bool {
	diagonal1 := []int{0, 4, 8}
	diagonal2 := []int{2, 4, 6}

	symbol := g.Board[move]

	if g.Board[(move+3)%9] == symbol && g.Board[(move+6)%9] == symbol {
		return true
	}

	if g.Board[(move+1)%3+(move/3)*3] == symbol && g.Board[(move+2)%3+(move/3)*3] == symbol {
		return true
	}

	if move%2 == 0 {
		if containsInt(diagonal1, move) &&
			g.Board[diagonal1[(index(diagonal1, move)+1)%3]] == symbol &&
			g.Board[diagonal1[(index(diagonal1, move)+2)%3]] == symbol {
			return true
		}
		if containsInt(diagonal2, move) &&
			g.Board[diagonal2[(index(diagonal2, move)+1)%3]] == symbol &&
			g.Board[diagonal2[(index(diagonal2, move)+2)%3]] == symbol {
			return true
		}
	}

	return false
}

func (g *Game) DoMove(p *Player) result {
	idx := 0
	// Starting AI which is random
	if p.ID == 1 {
		idx = rand.Intn(len(g.Board))
		for idx != 0 {
			idx = rand.Intn(len(g.Board))
		}
	} else {
		min := g.Board[0]
		for _, val := range g.Board {
			if min == 0 && val != 0 || val != 0 && val < min {
				min = val
			}
		}
		newBoard := make([]float64, len(g.Board))
		for i := range newBoard {
			if g.Board[i] == 0 {
				continue
			} else if g.Board[i] == min {
				newBoard[i] = 0.5
			} else {
				newBoard[i] = 1
			}
		}

		res := p.AI.Run(newBoard)
		max := res[0]

		for i, v := range res {
			if v > max {
				idx = i
				max = v
			}
		}
	}

	if g.Board[idx] != 0 {
		return result{Round: g.Round, Continue: false, Won: false}
	}

	g.Board[idx] = p.ID
	g.Round++
	if g.checkWin(idx) {
		return result{Round: g.Round, Continue: false, Won: true}
	}

	return result{Round: g.Round, Continue: true, Won: false}
}

func doMoveHuman() int {
	reader := bufio.NewReader(os.Stdin)
	sent := true
	num := -1

	for sent {
		fmt.Println("Give me a number between 0 and 8")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		n, ok := strconv.Atoi(text)

		if ok == nil && n <= 8 && n >= 0 {
			num = n
			sent = false
		} else if ok == nil {
			fmt.Println("The number was not between 0 and 8")
		} else {
			fmt.Println("You too stupid to give number?")
		}
	}
	return num
}

func (g *Game) runGame() (uint64, int) {
	turn := uint64(rand.Intn(2))

	for gameRound := 0; gameRound < 9; gameRound++ {
		g.printField(turn)
		move := doMoveHuman()
		fmt.Println("Player", turn+1, "chose:", move)
		if g.Board[move] != 0 {
			return turn - 2, gameRound + 1
		}
		g.Board[move] = turn + 1
		if g.checkWin(move) {
			return turn, gameRound + 1
		}
		turn = (turn + 1) % 2
	}
	return 2, 10
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := newGame()
	g.runGame()
}
