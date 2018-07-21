package Game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	board []int
}

var convert = map[int]string{0: " ", 1: "O", 2: "X"}

func ContainsInt(arr []int, i int) bool {
	for _, x := range arr {
		if i == x {
			return true
		}
	}
	return false
}

func InitGame() *Game {
	g := new(Game)
	g.board = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	return g
}

func (g *Game) Reset() {
	g.board = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (g *Game) PrintField(player int) {
	fmt.Printf("%v|%v|%v\n", convert[g.board[0]], convert[g.board[1]], convert[g.board[2]])
	fmt.Printf("-+-+-\n")
	fmt.Printf("%v|%v|%v\n", convert[g.board[3]], convert[g.board[4]], convert[g.board[5]])
	fmt.Printf("-+-+-\n")
	fmt.Printf("%v|%v|%v\n", convert[g.board[6]], convert[g.board[7]], convert[g.board[8]])
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

	symbol := g.board[move]

	if g.board[(move+3)%9] == symbol && g.board[(move+6)%9] == symbol {
		return true
	}

	if g.board[(move+1)%3+(move/3)*3] == symbol && g.board[(move+2)%3+(move/3)*3] == symbol {
		return true
	}

	if move%2 == 0 {
		if ContainsInt(diagonal1, move) &&
			g.board[diagonal1[(index(diagonal1, move)+1)%3]] == symbol &&
			g.board[diagonal1[(index(diagonal1, move)+2)%3]] == symbol {
			return true
		}
		if ContainsInt(diagonal2, move) &&
			g.board[diagonal2[(index(diagonal2, move)+1)%3]] == symbol &&
			g.board[diagonal2[(index(diagonal2, move)+2)%3]] == symbol {
			return true
		}
	}

	return false
}

func Get() int {
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

func (g *Game) RunGame() (int, int) {
	turn := rand.Intn(2)

	for gameRound := 0; gameRound < 9; gameRound++ {
		g.PrintField(turn)
		move := Get()
		fmt.Println("Player", turn+1, "chose:", move)
		if g.board[move] != 0 {
			return turn - 2, gameRound + 1
		}
		g.board[move] = turn + 1
		if g.checkWin(move) {
			return turn, gameRound + 1
		}
		turn = (turn + 1) % 2
	}
	return 2, 10
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := InitGame()
	g.RunGame()
}
