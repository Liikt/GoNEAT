package cppn

import (
	"math/rand"
	"strings"
)

type mutateFunc func(bool)

func mutateWrapper(flag bool, f mutateFunc, chance float64, funcName string) {
	for chance > 0 {
		if rand.Float64() < chance {
			f(flag)
		}
		chance--
	}
}

func (g *Genome) mutateEnableDisable(enable bool) {
	candidates := make([]string, 0)

	for _, gene := range g.genes {
		if gene.enabled == enable {
			candidates = append(candidates, gene.id)
		}
	}

	if len(candidates) == 0 {
		return
	}

	g.genes[candidates[rand.Intn(len(candidates))]].enabled = !g.genes[candidates[rand.Intn(len(candidates))]].enabled
}

func (g *Genome) mutateLink(forceBias bool) {
	ret := -1

	for ok := false; !ok; ok = ret == 0 {
		neuron1 := g.randomNeuron(false)
		neuron2 := g.randomNeuron(true)
		for neuron1 == neuron2 {
			neuron1 = g.randomNeuron(false)
			neuron2 = g.randomNeuron(true)
		}

		if strings.HasPrefix(neuron1, "In") && strings.HasPrefix(neuron2, "In") {
			return
		}

		if strings.HasPrefix(neuron2, "In") {
			neuron2, neuron1 = neuron1, neuron2
		}

		if forceBias && neuron2 != "Bias" {
			neuron1 = "Bias"
		} else if neuron2 == "Bias" {
			neuron1, neuron2 = neuron2, neuron1
		}

		ret = g.addLink(g.neurons[neuron1], g.neurons[neuron2], nil, genID())
	}
}

func (g *Genome) mutateWeights(_ bool) {
	c := make(chan bool)
	for _, ge := range g.genes {
		go func(ge *gene, c chan bool) {
			if rand.Float64() < g.staticRates["PerturbChance"] {
				ge.weight += rand.Float64()*2*g.mutationRates["WeightStepSize"] - g.mutationRates["WeightStepSize"]
			} else {
				ge.weight = rand.Float64()/2.0 - 0.25
			}
			c <- true
		}(ge, c)
	}

	for range g.genes {
		<-c
	}
}

func (g *Genome) mutateNode(_ bool) {
	geneIndex := rand.Intn(len(g.genes))
	nameArr := keyGene(g.genes)

	if !g.genes[nameArr[geneIndex]].enabled {
		return
	}

	newGenes := make([]*gene, 2)

	for x := 0; x < 2; x++ {
		newGenes[x] = g.genes[nameArr[geneIndex]].copy()
		newGenes[x].innovation = newInnovation()
		newGenes[x].id = genID()
	}

	newNeuron := initNeuron(genID())
	g.network[newNeuron.id] = []string{}
	g.neurons[newNeuron.id] = newNeuron

	g.addLink(g.genes[nameArr[geneIndex]].sourceNeuron, newNeuron, newGenes[0], newGenes[0].id)
	g.addLink(newNeuron, g.genes[nameArr[geneIndex]].destinationNeuron, newGenes[1], newGenes[1].id)

	g.genes[nameArr[geneIndex]].enabled = false
	g.numNodes++
}

func (g *Genome) mutateFunctions(_ bool) {
	for _, neuron := range g.neurons {
		if rand.Intn(3) == 0 {
			neuron.function = functions[rand.Intn(len(functions))]
		}
	}
}

func (g *Genome) randomNeuron(noIn bool) string {
	labelsAdList := keysNeuron(g.neurons)
	i := rand.Intn(len(g.neurons))

	for !strings.HasPrefix(labelsAdList[i], "In") && noIn {
		i = rand.Intn(len(g.neurons))
	}

	return labelsAdList[i]
}
