package cppn

import (
	"fmt"
	"math/rand"
)

type neuron struct {
	id       string
	incoming []chan float64
	outgoing []chan float64
	function neuronFunction
}

func initNeuron(id string) *neuron {
	return &neuron{id: id, function: functions[rand.Intn(len(functions))]}
}

func (n *neuron) copy() *neuron {
	newNeuron := new(neuron)
	newNeuron.id = n.id
	newNeuron.function = n.function

	return newNeuron
}

func (n *neuron) addIncLink(l chan float64) {
	n.incoming = append(n.incoming, l)
}

func (n *neuron) addOutLink(l chan float64) {
	n.outgoing = append(n.outgoing, l)
}

func (n *neuron) equal(o *neuron) bool {
	if len(n.incoming) != len(o.incoming) || len(n.outgoing) != len(o.outgoing) {
		return false
	}

	for x := range n.incoming {
		if n.incoming[x] != o.incoming[x] {
			return false
		}
	}

	for x := range n.outgoing {
		if n.outgoing[x] != o.outgoing[x] {
			return false
		}
	}

	return n.id == o.id
}

func (n *neuron) containsInpLink(link *gene) bool {
	for _, c := range n.incoming {
		if c == link.destination {
			return true
		}
	}
	return false
}

func (n *neuron) containsOutLink(link *gene) bool {
	for _, c := range n.outgoing {
		if c == link.destination {
			return true
		}
	}
	return false
}

func (n *neuron) run() {
	value := 1.0
	if n.id != "Bias" {
		value = 0.0
		in := merge(n.incoming...)
		for i := range in {
			value += i
		}
	}

	value = n.function(value)

	for _, out := range n.outgoing {
		out <- value
	}
}

func (n *neuron) String() string {
	return fmt.Sprintf("Neuron: %v", n.id)
}
