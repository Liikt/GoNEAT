package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/liikt/GoNEAT/framework"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	ctx := framework.NewGlobalCtx(100)
	fmt.Println("Starting the Simulation")
	ctx.Epoch()
}
