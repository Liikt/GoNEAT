package substrate

import (
	"math"
)

type neuron struct {
	id       string
	incoming []chan float64
	outgoing []chan float64
}

func (n *neuron) addIncLink(l chan float64) {
	n.incoming = append(n.incoming, l)
}

func (n *neuron) addOutLink(l chan float64) {
	n.outgoing = append(n.outgoing, l)
}

func sigmoig(inp float64) float64 {
	if inp <= -8.0 {
		return -1.0
	} else if inp >= 8.0 {
		return 1.0
	}
	return 2.0/(1.0+math.Exp(-4.9*inp)) - 1.0
}

func (n *neuron) run() {
	v := 0.0
	in := merge(n.incoming...)
	for i := range in {
		v += i
	}
	v = sigmoig(v)
	for _, c := range n.outgoing {
		c <- v
	}
}
