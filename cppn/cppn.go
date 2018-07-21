package cppn

var (
	gamesToPlay        = 10
	deltaDisjoint      = 2.0
	weightScale        = 1.0
	deltaThreshold     = 0.7
	stalenessThreshold = 15
	functions          = []neuronFunction{
		sigmoig,
		sawtooth,
		sin,
		abs,
		square,
	}
)
