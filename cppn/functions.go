package cppn

import "math"

type neuronFunction func(neuronValue float64) float64

func sigmoig(inp float64) float64 {
	if inp <= -8.0 {
		return -1.0
	} else if inp >= 8.0 {
		return 1.0
	}
	return 2.0/(1.0+math.Exp(-4.9*inp)) - 1.0
}

func sawtooth(inp float64) float64 {
	return float64(int(math.Floor(inp*100))%256) / 100
}

func sin(inp float64) float64 {
	return math.Sin(inp)
}

func abs(inp float64) float64 {
	return math.Abs(inp)
}

func square(inp float64) float64 {
	return math.Pow(inp, 2.0)
}
