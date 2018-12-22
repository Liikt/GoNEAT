package cppn

import (
	"math"
)

type Pool struct {
	poolSize       int
	maxNodes       int
	species        []*species
	generation     int
	maxFitness     float64
	averageFitness float64
	topGenome      *Genome
}

func InitPool(poolSize int) *Pool {
	pool := &Pool{
		poolSize:       poolSize,
		maxNodes:       100000,
		species:        make([]*species, 0),
		generation:     1,
		maxFitness:     0.0,
		averageFitness: 0.0,
		topGenome:      nil,
	}

	genomeTemplate := initGenome(true)
	for x := 0; x < poolSize; x++ {
		g := genomeTemplate.copy()
		g.mutate()
		pool.addToSpecies(g)
	}

	return pool
}

func (p *Pool) addToSpecies(g *Genome) {
	found := false
	for s := 0; s < len(p.species); s++ {
		if p.species[s].includes(g) {
			p.species[s].genomes = append(p.species[s].genomes, g)
			found = true
			break
		}
	}
	if !found {
		sp := initSpecies(g, []*Genome{g}, g.Fitness, g.Fitness, 0)
		p.species = append(p.species, sp)
	}
}

func (p *Pool) calcAdjustedFitness() {
	allGenomes := make([]*Genome, 0)
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
			g.Fitness /= math.Pow(counter, 1.0/3.0)
		}
	}
}

func (p *Pool) removeStaleSpecies() {
	tmp := make([]*species, 0)
	for _, s := range p.species {
		if s.survives(p.averageFitness) {
			tmp = append(tmp, s)
		}
	}
	p.species = tmp
}

func (p *Pool) cullSpecies(cutToOne bool) {
	for _, species := range p.species {
		species.cullSpecies(cutToOne)
	}
}

func (p *Pool) GetGenomes() []*Genome {
	ret := make([]*Genome, 0)

	for _, s := range p.species {
		for _, g := range s.genomes {
			if g != nil {
				ret = append(ret, g)
			}
		}
	}

	return ret
}
