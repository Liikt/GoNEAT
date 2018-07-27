package substrate

type link struct {
	src    chan float64
	dst    chan float64
	weight float64
}

func (l *link) set() {
	v := <-l.src
	l.dst <- v * l.weight
}
