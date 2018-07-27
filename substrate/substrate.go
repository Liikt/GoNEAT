package substrate

import (
	"fmt"
	"strconv"

	"github.com/liikt/GoNEAT/cppn"
)

type Substrate struct {
	neurons    [][]*neuron
	links      []*link
	width      int
	height     int
	inpNeurons map[string]chan float64
	outNeurons map[string]chan float64
}

func (s *Substrate) indexToVal(x, y int) (float64, float64) {
	stepX := 2. / float64(s.width)
	stepY := 2. / float64(s.height)

	return -1. + (stepX * float64(x)), 1. - (stepY * float64(y))
}

func BuildSubstrate(w, h int, inpNeurons, outNeurons [][]int) *Substrate {
	ret := &Substrate{width: w, height: h}

	for x := 0; x < ret.height; x++ {
		for y := 0; y < ret.width; y++ {
			ret.neurons[x][y] = &neuron{id: genID()}
		}
	}

	for i, l := range inpNeurons {
		if l[0] >= 0 && l[0] < ret.width && l[1] >= 0 && l[1] < ret.height {
			inpChan := make(chan float64)
			ret.neurons[l[0]][l[1]].id = "IN"
			ret.neurons[l[0]][l[1]].addIncLink(inpChan)
			ret.inpNeurons["IN"+strconv.Itoa(i)] = inpChan
		}
	}

	for i, l := range outNeurons {
		if l[0] >= 0 && l[0] < ret.width && l[1] >= 0 && l[1] < ret.height && ret.neurons[l[0]][l[1]].id != "IN" {
			outChan := make(chan float64)
			ret.neurons[l[0]][l[1]].id = "OUT"
			ret.neurons[l[0]][l[1]].addOutLink(outChan)
			ret.outNeurons["OUT"+strconv.Itoa(i)] = outChan
		}
	}

	return ret
}

func (s *Substrate) Populate(g *cppn.Genome) {
	isReachable := make(map[string][]string)
	for x1 := 0; x1 < s.height; x1++ {
		for y1 := 0; y1 < s.height; y1++ {
			if s.neurons[x1][y1].id == "OUT" {
				continue
			}

			x1Val, y1Val := s.indexToVal(x1, y1)
			n1ID := s.neurons[x1][y1].id

			for x2 := 0; x2 < s.height; x2++ {
				for y2 := 0; y2 < s.width; y2++ {
					if x1 == x2 && y1 == y2 || s.neurons[x2][y2].id == "IN" {
						continue
					}

					x2Val, y2Val := s.indexToVal(x2, y2)
					w := g.GetWeight(x1Val, y1Val, x2Val, y2Val)
					n2ID := s.neurons[x2][y2].id

					if w > -8. && w < 8. && !strIn(isReachable[n1ID], n2ID) {
						for _, n := range isReachable[n1ID] {
							if !strIn(isReachable[n2ID], n) {
								isReachable[n2ID] = append(isReachable[n2ID], n)
							}
						}
						isReachable[n2ID] = append(isReachable[n2ID], n1ID)

						srcChan := make(chan float64)
						dstChan := make(chan float64)
						newLink := &link{weight: w, src: srcChan, dst: dstChan}
						s.neurons[x1][y1].addOutLink(srcChan)
						s.neurons[x2][y2].addIncLink(dstChan)
						s.links = append(s.links, newLink)
					}
				}
			}
		}
	}
}

func (s *Substrate) Run(input []float64) []float64 {
	if len(input) != len(s.inpNeurons) {
		fmt.Println("Wrong input array size", len(input), len(s.inpNeurons))
		return []float64{}
	}

	for _, l := range s.links {
		go l.run()
	}

	for _, l := range s.neurons {
		for _, n := range l {
			go n.run()
		}
	}

	for i, inp := range input {
		s.inpNeurons["IN"+strconv.Itoa(i)] <- inp
	}

	ret := make([]float64, len(s.outNeurons))
	for i := 0; i < len(s.outNeurons); i++ {
		ret[i] = <-s.outNeurons["OUT"+strconv.Itoa(i)]
	}

	return ret
}
