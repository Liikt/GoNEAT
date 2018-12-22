package cppn

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Pool struct {
	Population     int
	maxNodes       int
	Species        []*species
	Generation     int
	maxFitness     float64
	averageFitness float64
	topGenome      *Genome
}

func InitPool(popSize int) *Pool {
	pool := &Pool{
		Population:     popSize,
		maxNodes:       100000,
		Species:        make([]*species, 0),
		Generation:     1,
		maxFitness:     0.0,
		averageFitness: 0.0,
		topGenome:      nil,
	}

	genomeTemplate := initGenome(true)
	for x := 0; x < popSize; x++ {
		g := genomeTemplate.copy()
		g.mutate()
		pool.AddToSpecies(g)
	}

	return pool
}

func (p *Pool) AddToSpecies(g *Genome) {
	found := false
	for s := 0; s < len(p.Species); s++ {
		if p.Species[s].includes(g) {
			p.Species[s].genomes = append(p.Species[s].genomes, g)
			found = true
			break
		}
	}
	if !found {
		sp := initSpecies(g, []*Genome{g}, g.Fitness, g.Fitness, 0)
		p.Species = append(p.Species, sp)
	}
}

func (p *Pool) calcAdjustedFitness() {
	allGenomes := make([]*Genome, 0)
	for _, s := range p.Species {
		allGenomes = append(allGenomes, s.genomes...)
	}
	for _, s := range p.Species {
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

func (p *Pool) RemoveStaleSpecies() {
	tmp := make([]*species, 0)
	for _, s := range p.Species {
		if s.survives(p.averageFitness) {
			tmp = append(tmp, s)
		}
	}
	p.Species = tmp
}

func (p *Pool) CullSpecies(cutToOne bool) {
	for _, species := range p.Species {
		species.cullSpecies(cutToOne)
	}
}

func (p *Pool) GetGenomes() []*Genome {
	ret := make([]*Genome, 0)

	for _, s := range p.Species {
		for _, g := range s.genomes {
			if g != nil {
				ret = append(ret, g)
			}
		}
	}

	return ret
}

func (p *Pool) RankGlobally() {
	genomes := p.GetGenomes()
	sort.Slice(genomes, func(i, j int) bool {
		return genomes[i].Fitness < genomes[j].Fitness
	})
	for i, g := range genomes {
		g.globalRank = i
	}
}

func (p *Pool) CalcTotalAvgFitness() float64 {
	total := 0.0
	for _, s := range p.Species {
		total += s.AverageFitness
	}
	return total
}

func (p *Pool) RemoveWeakSpecies() {
	survived := make([]*species, 0)
	sum := p.CalcTotalAvgFitness()

	for _, s := range p.Species {
		if breed := math.Floor(s.AverageFitness / sum * float64(p.Population)); breed >= 1 {
			survived = append(survived, s)
		}
	}
	p.Species = survived
}

func (p *Pool) NewGeneration() {
	p.CullSpecies(false)
	p.RankGlobally()
	p.RemoveStaleSpecies()
	p.RankGlobally()
	for _, s := range p.Species {
		s.CalcAverageFitness()
	}
	p.RemoveWeakSpecies()
	sum := p.CalcTotalAvgFitness()
	children := make([]*Genome, 0)
	for _, s := range p.Species {
		breed := math.Floor(s.AverageFitness/sum*float64(p.Population)) - 1
		for x := 0; x < int(breed); x++ {
			children = append(children, s.Breed())
		}
	}
	p.CullSpecies(true)
	for x := len(p.Species) + len(children); x < p.Population; x++ {
		s := p.Species[rand.Intn(len(p.Species))]
		children = append(children, s.Breed())
	}
	for _, c := range children {
		p.AddToSpecies(c)
	}
	fmt.Println("Done with Generation", p.Generation)
	p.Generation++
}
