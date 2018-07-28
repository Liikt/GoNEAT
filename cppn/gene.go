package cppn

import "fmt"

type gene struct {
	id                string
	source            chan float64
	sourceNeuron      *neuron
	destination       chan float64
	destinationNeuron *neuron
	weight            float64
	enabled           bool
	innovation        int
}

func initGene(id string, src, dest chan float64, srcNode, dstNode *neuron, weight float64, innov int) *gene {
	return &gene{
		id:                id,
		source:            src,
		sourceNeuron:      srcNode,
		destination:       dest,
		destinationNeuron: dstNode,
		weight:            weight,
		enabled:           true,
		innovation:        innov,
	}
}

func (g *gene) equal(o *gene) bool {
	if g.sourceNeuron.id == o.sourceNeuron.id && g.destinationNeuron.id == o.destinationNeuron.id {
		return true
	}
	return false
}

func (g *gene) copy() *gene {
	cpy := new(gene)
	cpy.id = g.id
	cpy.weight = g.weight
	cpy.enabled = g.enabled
	cpy.innovation = g.innovation
	return cpy
}

func (g *gene) run() {
	val := <-g.source
	if g.enabled {
		g.destination <- g.weight * val
	} else {
		g.destination <- 0.0 * val
	}
}

func (g *gene) String() string {
	ret := ""

	if g.sourceNeuron == nil {
		ret += "<nil> => "
	} else {
		ret += g.sourceNeuron.id + " => "
	}

	if g.destinationNeuron == nil {
		ret += "<nil>\n"
	} else {
		ret += g.destinationNeuron.id + "\n"
	}

	sid, did := "", ""
	if g.sourceNeuron != nil {
		sid = g.sourceNeuron.id
	}

	if g.destinationNeuron != nil {
		did = g.destinationNeuron.id
	}

	return fmt.Sprintln(ret, "ID:", g.id, "Weight:", g.weight, "Enabled", g.enabled, "Innovation:",
		g.innovation, "Src:", sid, "Dst:", did)
}
