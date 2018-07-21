package cppn

import "sort"

type genomeSorter struct {
	genes []*genome
	by    func(p1, p2 *genome) bool
}

type byFunc func(p1, p2 *genome) bool

func (s *genomeSorter) Len() int {
	return len(s.genes)
}

func (s *genomeSorter) Less(i, j int) bool {
	return s.by(s.genes[i], s.genes[j])
}

func (s *genomeSorter) Swap(i, j int) {
	s.genes[i], s.genes[j] = s.genes[j], s.genes[i]
}

func (by byFunc) Sort(genomes []*genome) {
	ps := &genomeSorter{
		genes: genomes,
		by:    by,
	}
	sort.Sort(ps)
}

func sortGenomes(g *species) {
	f := func(g1, g2 *genome) bool {
		return g1.fitness > g2.fitness
	}
	byFunc(f).Sort(g.genomes)
}
