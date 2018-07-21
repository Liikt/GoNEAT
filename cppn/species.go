package cppn

import (
	"fmt"
	"math"
	"math/rand"
)

type species struct {
	genomes        []*genome
	representative *genome
	topFitness     float64
	averageFitness float64
	staleness      int
}

func initSpecies(rep *genome, genomes []*genome, topFit, averageFit float64, staleness int) *species {
	sp := new(species)
	sp.staleness = staleness
	sp.representative = rep
	sp.genomes = genomes

	if rep != nil {
		sp.topFitness = rep.fitness
		sp.averageFitness = rep.fitness
	} else {
		sp.topFitness = topFit
		sp.averageFitness = averageFit
	}

	return sp
}

func (s *species) includes(g *genome) bool {
	return s.representative.shareFunction(g)
}

func (s *species) calcAverageFitness() float64 {
	ret := 0.0

	for _, g := range s.genomes {
		ret += g.fitness
	}

	return ret / float64(len(s.genomes))
}

func (s *species) calcTopFitness() float64 {
	ret := -1000000.0

	for _, g := range s.genomes {
		if g.fitness > ret {
			ret = g.fitness
		}
	}

	return ret
}

func (s *species) breed() *genome {
	child := new(genome)

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
	child.fitness = 0.0

	return child
}

func (s *species) cullSpecies(cutToOne bool) {
	sortGenomes(s)
	for _, x := range s.genomes {
		fmt.Printf("%v ", x.fitness)
	}
	fmt.Print("\n")

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
