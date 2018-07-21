package cppn

import (
	"math"

	"github.com/liikt/GoNEAT/Game"
)

type pool struct {
	poolSize       int
	numInputNodes  int
	numOutputNodes int
	maxNodes       int
	species        []*species
	generation     int
	maxFitness     float64
	averageFitness float64
	topGenome      *genome
	game           *Game.Game
}

func initPool(inNodes, outNodes, poolSize int) *pool {
	pool := &pool{
		poolSize:       poolSize,
		numInputNodes:  inNodes,
		numOutputNodes: outNodes,
		maxNodes:       100000,
		species:        make([]*species, 0),
		generation:     1,
		maxFitness:     0.0,
		averageFitness: 0.0,
		topGenome:      nil,
		game:           nil,
	}

	genomeTemplate := initGenome()
	for x := 0; x < poolSize; x++ {
		g := genomeTemplate.copy()
		g.mutate()
		pool.addToSpecies(g)
	}

	return pool
}

func (p *pool) addToSpecies(g *genome) {
	found := false
	for s := 0; s < p.poolSize; s++ {
		if p.species[s].includes(g) {
			p.species[s].genomes = append(p.species[s].genomes, g)
			found = true
			break
		}
	}
	if !found {
		sp := initSpecies(g, make([]*genome, 0), g.fitness, g.fitness, 0)
		p.species = append(p.species, sp)
	}
}

func (p *pool) calcAdjustedFitness() {
	allGenomes := make([]*genome, 0)
	for _, s := range p.species {
		allGenomes = append(allGenomes, s.genomes...)
	}
	for _, s := range p.species {
		for _, g := range s.genomes {
			counter := 0.0
			for _, ge := range allGenomes {
				if g.shareFunction(ge) {
					counter += 1.0
				}
			}
			g.fitness /= math.Pow(counter, 1.0/3.0)
		}
	}
}

func (p *pool) removeStaleSpecies() {
	tmp := make([]*species, 0)
	for _, s := range p.species {
		if s.survives(p.averageFitness) {
			tmp = append(tmp, s)
		}
	}
	p.species = tmp
}

func (p *pool) cullSpecies(cutToOne bool) {
	for _, species := range p.species {
		species.cullSpecies(cutToOne)
	}
}
