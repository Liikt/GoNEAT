package cppn

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type genome struct {
	inputNodes    []chan float64
	outputNode    chan float64
	numNodes      int
	maxNeurons    int
	fitness       float64
	genes         map[string]*gene
	neurons       map[string]*neuron
	network       map[string][]string
	mutationRates map[string]float64
	staticRates   map[string]float64
}

func initGenome() *genome {
	g := new(genome)
	g.numNodes = 5
	g.fitness = 0.0
	g.genes = make(map[string]*gene)
	g.neurons = make(map[string]*neuron)
	g.network = make(map[string][]string)
	g.mutationRates = make(map[string]float64)
	g.staticRates = make(map[string]float64)
	g.mutationRates["MutateConnectionsChance"] = 0.65
	g.mutationRates["LinkMutationChance"] = 1.5
	g.mutationRates["NodeMutationChance"] = 0.5
	g.mutationRates["BiasMutationChance"] = 0.4
	g.mutationRates["DisableMutationChance"] = 0.4
	g.mutationRates["EnableMutationChance"] = 0.2
	g.mutationRates["WeightStepSize"] = 0.0
	g.mutationRates["FunctionChangeChance"] = 0.3
	g.staticRates["PerturbChance"] = 0.95
	g.staticRates["CrossoverChance"] = 0.75

	g.basicNetwork()

	return g
}

func (g *genome) basicNetwork() {
	outNode := initNeuron("OutN")
	g.neurons[outNode.id] = outNode
	g.network[outNode.id] = []string{}
	outChan := make(chan float64)
	outNode.addOutLink(outChan)
	g.outputNode = outChan
	for x := 0; x < 4; x++ {
		inNode := initNeuron("InN" + strconv.Itoa(x))
		g.neurons[inNode.id] = inNode
		g.network[inNode.id] = []string{}

		g.addLink(inNode, outNode, nil, "InN"+strconv.Itoa(x)+"Link")

		inChan := make(chan float64)
		inNode.addIncLink(inChan)
		g.inputNodes = append(g.inputNodes, inChan)
	}

	biasNode := initNeuron("Bias")
	g.neurons[biasNode.id] = biasNode
	g.network[biasNode.id] = []string{}

	g.addLink(biasNode, outNode, nil, "BiasLink")
}

func (g *genome) addLink(fromNode, intoNode *neuron, ge *gene, id string) int {
	if _, ok := g.network[fromNode.id]; !ok {
		fmt.Println("Just tried to add a link " + ge.id + " to a nonexisting outgoing node.")
		return -1
	} else if _, ok := g.network[intoNode.id]; !ok {
		fmt.Println("Just tried to add a link " + ge.id + " to a nonexisting ingoing node.")
		return -1
	}

	cin := make(chan float64)
	cout := make(chan float64)

	if ge == nil {
		ge = initGene(id, cout, cin, nil, nil, rand.Float64()-0.5, newInnovation())
	} else {
		ge.source = cout
		ge.destination = cin
	}

	if arr, ok := g.network[fromNode.id]; ok && containsString(arr, intoNode.id) {
		// fmt.Println("A link from", fromNode.id, "to", intoNode.id, "exists already.")
		return -1
	} else if strings.HasPrefix(intoNode.id, "In") {
		// fmt.Println("Can't create a link into an Input neuron")
		return -1
	} else if strings.HasPrefix(fromNode.id, "Out") {
		// fmt.Println("Can't create a link from an Output neuron")
		return -1
	} else if fromNode.id == intoNode.id {
		// fmt.Println("Can't create a recurrent link for", fromNode.id)
		return -1

	}

	g.network[ge.id] = []string{}
	g.genes[ge.id] = ge

	ge.sourceNeuron = fromNode
	ge.destinationNeuron = intoNode

	g.network[intoNode.id] = append(g.network[intoNode.id], ge.id)
	g.network[intoNode.id] = append(g.network[intoNode.id], fromNode.id)
	g.network[ge.id] = append(g.network[ge.id], intoNode.id)
	g.network[ge.id] = append(g.network[ge.id], fromNode.id)
	g.network[fromNode.id] = append(g.network[fromNode.id], ge.id)
	g.network[fromNode.id] = append(g.network[fromNode.id], intoNode.id)

	intoNode.addIncLink(cin)
	fromNode.addOutLink(cout)
	return 0
}

func (g *genome) containsLink(link *gene) bool {
	if _, ok := g.genes[link.id]; ok {
		return true
	}
	return false
}

func (g *genome) copy() *genome {
	newGenome := initGenome()
	newGenome.numNodes = g.numNodes
	newGenome.fitness = g.fitness
	newGenome.network = g.network
	newGenome.mutationRates = g.mutationRates
	newGenome.staticRates = g.staticRates

	for _, n := range g.neurons {
		newGenome.neurons[n.id] = initNeuron(n.id)
	}

	for _, ge := range g.genes {
		srcN := newGenome.neurons[ge.sourceNeuron.id]
		dstN := newGenome.neurons[ge.destinationNeuron.id]

		newGenome.addLink(srcN, dstN, ge.copy(), ge.id)
	}

	outChan := make(chan float64)
	newGenome.neurons["OutN"].outgoing = append(newGenome.neurons["OutN"].outgoing, outChan)
	newGenome.outputNode = outChan

	for x := 0; x < 4; x++ {
		inChan := make(chan float64)
		newGenome.neurons["InN"+strconv.Itoa(x)].incoming = append(newGenome.neurons["InN"+strconv.Itoa(x)].incoming, inChan)
		newGenome.inputNodes[x] = inChan
	}

	return newGenome
}

func (g *genome) GetWeight(x1, y1, x2, y2 float64) float64 {
	for _, gene := range g.genes {
		go gene.run()
	}

	for _, neuron := range g.neurons {
		go neuron.run()
	}

	g.inputNodes[0] <- x1
	g.inputNodes[1] <- y1
	g.inputNodes[2] <- x2
	g.inputNodes[3] <- y2

	return <-g.outputNode
}

func (g *genome) mutate() {
	for key, val := range g.mutationRates {
		i := rand.Intn(2)
		if i == 1 {
			g.mutationRates[key] = 0.95 * val
		} else {
			g.mutationRates[key] = val / 0.95
		}
	}

	funcArr := []mutateFunc{
		g.mutateLink,
		g.mutateEnableDisable,
		g.mutateWeights,
		g.mutateNode,
		g.mutateFunctions,
	}

	chanceArr := []float64{
		g.mutationRates["LinkMutationChance"],
		g.mutationRates["DisableMutationChance"],
		g.mutationRates["MutateConnectionsChance"],
		g.mutationRates["NodeMutationChance"],
		g.mutationRates["FunctionChangeChance"],
	}
	chanceNameArr := []string{
		"LinkMutationChance",
		"DisableMutationChance",
		"MutateConnectionsChance",
		"NodeMutationChance",
		"FunctionChangeChance",
	}

	secArr := []float64{
		g.mutationRates["BiasMutationChance"],
		g.mutationRates["EnableMutationChance"],
	}

	secNameArr := []string{
		"BiasMutationChance",
		"EnableMutationChance",
	}

	for x := 0; x < 5; x++ {
		if x < 2 {
			mutateWrapper(true, funcArr[x], secArr[x], secNameArr[x])
		}
		mutateWrapper(false, funcArr[x], chanceArr[x], chanceNameArr[x])
	}
}

func (g *genome) disjointGenes(other *genome) float64 {
	innovationSet1 := make([]int, len(g.genes))
	innovationSet2 := make([]int, len(other.genes))
	counter := 0

	for _, x := range g.genes {
		innovationSet1[counter] = x.innovation
		counter++
	}

	counter = 0
	for _, x := range other.genes {
		innovationSet2[counter] = x.innovation
		counter++
	}

	c := make(chan int)

	go findMissing(innovationSet1, innovationSet2, c)
	go findMissing(innovationSet2, innovationSet1, c)

	return float64(<-c+<-c) / math.Max(float64(len(g.genes)), float64(len(other.genes)))
}

func (g *genome) weightGap(other *genome) float64 {
	innoToObj1 := make(map[int]*gene)
	innoToObj2 := make(map[int]*gene)
	weightDiff := 0.0

	for _, x := range g.genes {
		innoToObj1[x.innovation] = x
	}

	for _, x := range other.genes {
		innoToObj2[x.innovation] = x
	}

	matching := matchingInts(indexGene(innoToObj1), indexGene(innoToObj2))

	for _, inno := range matching {
		weightDiff += math.Abs(innoToObj1[inno].weight - innoToObj2[inno].weight)
	}

	return weightDiff / float64(len(matching))
}

func (g *genome) shareFunction(other *genome) bool {
	geneticDifference := deltaDisjoint * g.disjointGenes(other)
	weightDifference := weightScale * g.weightGap(other)

	return (geneticDifference + weightDifference) < deltaThreshold
}

func (g *genome) crossover(other *genome) *genome {
	if other.fitness > g.fitness {
		g, other = other, g
	}

	child := initGenome()
	child.neurons["OutN"] = g.neurons["OutN"].copy()

	outChan := make(chan float64)
	child.neurons["OutN"].outgoing = append(child.neurons["OutN"].outgoing, outChan)
	child.outputNode = outChan

	for x := 0; x < 4; x++ {
		inChan := make(chan float64)
		child.neurons["InN"+strconv.Itoa(x)] = g.neurons["InN"+strconv.Itoa(x)].copy()
		child.neurons["InN"+strconv.Itoa(x)].incoming = append(child.neurons["InN"+strconv.Itoa(x)].incoming, inChan)
		child.inputNodes[x] = inChan
	}

	innovation2 := make(map[int]*gene)

	for _, gene := range other.genes {
		innovation2[gene.innovation] = gene
	}

	for _, gene := range g.genes {
		gene1 := gene
		if gene2, ok := innovation2[gene1.innovation]; ok && rand.Intn(2) == 1 && gene2.enabled {
			srcNeuron := gene2.sourceNeuron.copy()
			dstNeuron := gene2.destinationNeuron.copy()

			if _, ok := child.neurons[srcNeuron.id]; !ok {
				child.neurons[srcNeuron.id] = srcNeuron
				child.network[srcNeuron.id] = []string{}
			} else {
				srcNeuron = child.neurons[srcNeuron.id]
			}

			if _, ok := child.neurons[dstNeuron.id]; !ok {
				child.neurons[dstNeuron.id] = dstNeuron
				child.network[dstNeuron.id] = []string{}
			} else {
				dstNeuron = child.neurons[dstNeuron.id]
			}

			child.addLink(srcNeuron, dstNeuron, gene2, gene2.id)

		} else {
			srcNeuron := gene1.sourceNeuron.copy()
			dstNeuron := gene1.destinationNeuron.copy()

			if _, ok := child.neurons[srcNeuron.id]; !ok {
				child.neurons[srcNeuron.id] = srcNeuron
				child.network[srcNeuron.id] = []string{}
			} else {
				srcNeuron = child.neurons[srcNeuron.id]
			}

			if _, ok := child.neurons[dstNeuron.id]; !ok {
				child.neurons[dstNeuron.id] = dstNeuron
				child.network[dstNeuron.id] = []string{}
			} else {
				dstNeuron = child.neurons[dstNeuron.id]
			}

			child.addLink(srcNeuron, dstNeuron, gene1, gene1.id)
		}
	}

	child.numNodes = int(math.Max(float64(g.numNodes), float64(other.numNodes)))

	for mutation, val := range g.staticRates {
		child.staticRates[mutation] = val
	}

	for mutation, val := range g.mutationRates {
		child.mutationRates[mutation] = val
	}

	return child
}

func (g *genome) Debug() {
	for _, i := range g.neurons {
		fmt.Println(i.String())
	}

	fmt.Println()
	for _, i := range g.genes {
		fmt.Println(i.String())
	}
	return
}