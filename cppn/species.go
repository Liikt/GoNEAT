package cppn

import (
	"math"
	"math/rand"
)

type species struct {
	genomes        []*Genome
	representative *Genome
	topFitness     float64
	AverageFitness float64
	staleness      int
}

func initSpecies(rep *Genome, genomes []*Genome, topFit, averageFit float64, staleness int) *species {
	sp := new(species)
	sp.staleness = staleness
	sp.representative = rep
	sp.genomes = genomes

	if rep != nil {
		sp.topFitness = rep.Fitness
		sp.AverageFitness = rep.Fitness
	} else {
		sp.topFitness = topFit
		sp.AverageFitness = averageFit
	}

	return sp
}

func (s *species) includes(g *Genome) bool {
	return s.representative.shareFunction(g)
}

func (s *species) calcTopFitness() float64 {
	ret := -1000000.0

	for _, g := range s.genomes {
		if g.Fitness > ret {
			ret = g.Fitness
		}
	}

	return ret
}

func (s *species) Breed() *Genome {
	child := new(Genome)

	if rand.Float64() < s.representative.staticRates["CrossoverChance"] {
		g1 := rand.Intn(len(s.genomes))
		g2 := rand.Intn(len(s.genomes))

		if g1 == g2 {
			child = s.genomes[g1].copy()
		} else {
			child = s.genomes[g1].crossover(s.genomes[g2])
		}
	} else {
		g := s.genomes[rand.Intn(len(s.genomes))]
		child = g.copy()
	}

	child.mutate()
	child.Fitness = 0.0

	return child
}

func (s *species) cullSpecies(cutToOne bool) {
	sortGenomes(s)

	cut := int(math.Ceil(float64(len(s.genomes)) / 2.0))

	if cutToOne {
		cut = 1
	}

	s.genomes = s.genomes[:cut]
}

func (s *species) survives(poolFitness float64) bool {
	newTopFitness := s.calcTopFitness()

	if newTopFitness > s.topFitness {
		s.topFitness = newTopFitness
		s.staleness = 0
	} else {
		s.staleness++
	}

	if s.staleness >= stalenessThreshold && s.topFitness < poolFitness {
		return false
	}
	return true
}

func (s *species) CalcAverageFitness() {
	total := 0.0
	for _, g := range s.genomes {
		total += g.Fitness
	}
	s.AverageFitness = total / float64(len(s.genomes))
}
